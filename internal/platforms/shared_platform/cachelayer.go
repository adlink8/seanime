package shared_platform

import (
	"context"
	"errors"
	"fmt"
	"seanime/internal/api/bangumi"
	"seanime/internal/events"
	"seanime/internal/media"
	"seanime/internal/util"
	"seanime/internal/util/filecache"
	"seanime/internal/util/result"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"
	"github.com/samber/lo"
)

// devnote: I got lazy and used global variables

// 换锚说明（Phase 2 Wave B）：
//   - CacheLayer 由「AniList GraphQL 客户端包装」改为「*bangumi.Client 包装」；
//   - 职责不变：网络优先 + filecache 持久桶 + API 不可用时回退缓存 +
//     收藏/进度写操作离线队列；
//   - bangumi.Client 自身已有 GET 24h 缓存与节流，本层与其互补不冲突；
//   - 客户端依赖收敛为 BangumiAPI 最小接口（*bangumi.Client 天然满足，
//     测试可注入假实现）。

var ShouldCache = atomic.Bool{}
var IsWorking = atomic.Bool{}
var BangumiClient = atomic.Value{}

type failureRecord struct {
	timestamp time.Time
	err       error
}

var (
	failureTracking      = make([]failureRecord, 0)
	failureTrackingMutex sync.RWMutex
)

const (
	failureWindow     = 30 * time.Second // time window to consider failures
	failureThreshold  = 4                // number of failures needed to mark as down
	cleanupInterval   = 5 * time.Minute  // how often to clean up old failure records
	maxFailureRecords = 50               // maximum number of failure records to keep
)

const (
	AnimeCollectionBucket         = "anime-collection"
	MangaCollectionBucket         = "manga-collection"
	BaseAnimeBucket               = "base-anime"
	BaseMangaBucket               = "base-manga"
	CompleteAnimeBucket           = "complete-anime"
	AnimeDetailsBucket            = "anime-details"
	MangaDetailsBucket            = "manga-details"
	ViewerStatsBucket             = "viewer-stats"
	StudioDetailsBucket           = "studio-details"
	PendingMediaListUpdatesBucket = "pending-media-list-updates"

	maxNonCollectionCacheEntries      = 10
	maxNonCollectionMediaCacheEntries = 50
	// Collection update interval (refresh collection tracking every 30 minutes)
	collectionUpdateInterval = 30 * time.Minute

	// collectionPageLimit 收藏分页单页大小（Bangumi 上限 50）
	collectionPageLimit = 50
)

// BangumiAPI CacheLayer 依赖的 Bangumi 客户端最小接口。
// *bangumi.Client 天然满足；测试注入假实现。
type BangumiAPI interface {
	GetSubject(ctx context.Context, subjectID int) (*bangumi.Subject, error)
	GetRelatedSubjects(ctx context.Context, subjectID int) ([]bangumi.RelatedSubject, error)
	GetUserCollectionsByUser(ctx context.Context, username string, opts bangumi.UserCollectionsOpts) (*bangumi.UserCollectionsResult, error)
	UpsertCollection(ctx context.Context, subjectID int, body bangumi.CollectionUpsertBody) error
	PatchCollection(ctx context.Context, subjectID int, body bangumi.CollectionUpsertBody) error
	DeleteCollection(ctx context.Context, subjectID int) error
}

func init() {
	ShouldCache.Store(true)
	IsWorking.Store(true)

	go func() {
		// Every 10 seconds, check if the Bangumi client is working
		for {
			time.Sleep(time.Second * 10)
			if !ShouldCache.Load() {
				IsWorking.Store(true)
				continue
			}
			if IsWorking.Load() {
				continue
			}
			ref, ok := BangumiClient.Load().(*util.Ref[BangumiAPI])
			if !ok || ref == nil {
				IsWorking.Store(true)
				continue
			}
			client := ref.Get()
			// 探活：拉一个固定 subject。
			// 注意：bangumi.Client 自带 24h GET 缓存，探活可能命中缓存而「假阳」；
			// 但此时上层 GET 也能从缓存回包，功能等价可用，故可接受。
			_, err := client.GetSubject(context.Background(), 1)
			if err != nil {
				IsWorking.Store(false)
			} else {
				clearFailureTracking()
				events.GlobalWSEventManager.SendEvent(events.InfoToast, "The Bangumi API is back online")
				IsWorking.Store(true)
			}
		}
	}()

	// periodic cleanup of old failure records
	go func() {
		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()
		for range ticker.C {
			cleanupOldFailures()
		}
	}()
}

type (
	// CacheLayer is a "network-first" wrapper around a Bangumi client that caches fetched data in cache files.
	// It detects when the API client is not working and falls back to the cached data instead.
	// When the API client not working, it will still send the requests in the background and transition back to working state when the API client is working again.
	// Entry/progress updates are queued when the API client is not working; other mutations return an error.
	CacheLayer struct {
		bangumiClientRef       *util.Ref[BangumiAPI]
		fileCacher             *filecache.Cacher
		buckets                map[string]filecache.PermanentBucket
		logger                 *zerolog.Logger
		collectionMediaIDs     *result.Map[int, struct{}] // Track which media IDs are in collections
		lastCollectionUpdate   time.Time                  // When collections were last fetched
		logoutFunc             func()                     // called when an invalid token is detected
		pendingUpdateSyncMutex sync.Mutex
	}
)

// addFailureRecord adds a new failure record to the tracking
func addFailureRecord(err error) {
	failureTrackingMutex.Lock()
	defer failureTrackingMutex.Unlock()

	now := time.Now()
	failureTracking = append(failureTracking, failureRecord{
		timestamp: now,
		err:       err,
	})

	// keep only the most recent records
	if len(failureTracking) > maxFailureRecords {
		failureTracking = failureTracking[len(failureTracking)-maxFailureRecords:]
	}
}

// getRecentFailureCount returns the number of failures within the failure window
func getRecentFailureCount() int {
	failureTrackingMutex.RLock()
	defer failureTrackingMutex.RUnlock()

	now := time.Now()
	cutoff := now.Add(-failureWindow)
	count := 0

	for _, record := range failureTracking {
		if record.timestamp.After(cutoff) {
			count++
		}
	}

	return count
}

// cleanupOldFailures removes failure records older than the failure window
func cleanupOldFailures() {
	failureTrackingMutex.Lock()
	defer failureTrackingMutex.Unlock()

	now := time.Now()
	cutoff := now.Add(-failureWindow)
	validRecords := make([]failureRecord, 0, len(failureTracking))

	for _, record := range failureTracking {
		if record.timestamp.After(cutoff) {
			validRecords = append(validRecords, record)
		}
	}

	failureTracking = validRecords
}

// clearFailureTracking clears all failure records (called when API comes back online)
func clearFailureTracking() {
	failureTrackingMutex.Lock()
	defer failureTrackingMutex.Unlock()
	failureTracking = failureTracking[:0]
}

// NewCacheLayer returns a new instance of the global cache layer.
// An optional logoutFunc can be passed to perform server-side cleanup when an invalid token is detected.
func NewCacheLayer(client BangumiAPI, cacheDir string, logoutFunc ...func()) *CacheLayer {
	var clientRef *util.Ref[BangumiAPI]
	if client != nil {
		clientRef = util.NewRef[BangumiAPI](client)
		// 统一存储具体类型 *util.Ref[BangumiAPI]：atomic.Value 对后续 Store 的动态类型一致性有要求，
		// 直接存接口值会因不同实现类型（如测试注入的假客户端）触发 panic。
		BangumiClient.Store(clientRef)
	}

	fileCacher, err := filecache.NewCacher(cacheDir)
	if err != nil {
		// 缓存目录不可用时降级为直连模式（不落缓存）
		fileCacher = nil
	}

	buckets := make(map[string]filecache.PermanentBucket)
	buckets[AnimeCollectionBucket] = filecache.NewPermanentBucket(AnimeCollectionBucket)
	buckets[MangaCollectionBucket] = filecache.NewPermanentBucket(MangaCollectionBucket)
	buckets[BaseAnimeBucket] = filecache.NewPermanentBucket(BaseAnimeBucket)
	buckets[BaseMangaBucket] = filecache.NewPermanentBucket(BaseMangaBucket)
	buckets[CompleteAnimeBucket] = filecache.NewPermanentBucket(CompleteAnimeBucket)
	buckets[AnimeDetailsBucket] = filecache.NewPermanentBucket(AnimeDetailsBucket)
	buckets[MangaDetailsBucket] = filecache.NewPermanentBucket(MangaDetailsBucket)
	buckets[ViewerStatsBucket] = filecache.NewPermanentBucket(ViewerStatsBucket)
	buckets[StudioDetailsBucket] = filecache.NewPermanentBucket(StudioDetailsBucket)
	buckets[PendingMediaListUpdatesBucket] = filecache.NewPermanentBucket(PendingMediaListUpdatesBucket)

	logger := util.NewLogger()

	var logout func()
	if len(logoutFunc) > 0 {
		logout = logoutFunc[0]
	}

	cl := &CacheLayer{
		bangumiClientRef:   clientRef,
		fileCacher:         fileCacher,
		buckets:            buckets,
		logger:             logger,
		collectionMediaIDs: result.NewMap[int, struct{}](),
		logoutFunc:         logout,
	}

	cl.startQueuedUpdateSync()

	return cl
}

// client 返回底层 Bangumi 客户端（可能为 nil：测试直构 CacheLayer 时）。
func (c *CacheLayer) client() BangumiAPI {
	if c.bangumiClientRef == nil {
		return nil
	}
	return c.bangumiClientRef.Get()
}

// fetchAllUserCollections 全分页拉取指定用户的某类收藏
// （循环 offset 直到取满 total 或不足一页）。
func (c *CacheLayer) fetchAllUserCollections(ctx context.Context, username string, subjectType int) (*bangumi.UserCollectionsResult, error) {
	all := make([]bangumi.UserCollection, 0)
	offset := 0
	total := 0
	for {
		res, err := c.client().GetUserCollectionsByUser(ctx, username, bangumi.UserCollectionsOpts{
			SubjectType: &subjectType,
			Limit:       collectionPageLimit,
			Offset:      offset,
		})
		if err != nil {
			return nil, err
		}
		all = append(all, res.Data...)
		total = res.Total
		if len(res.Data) < collectionPageLimit || (total > 0 && len(all) >= total) {
			break
		}
		offset += len(res.Data)
	}
	return &bangumi.UserCollectionsResult{Data: all, Total: total}, nil
}

// checkAndUpdateWorkingState checks if the API client is working and updates the state
func (c *CacheLayer) checkAndUpdateWorkingState(err error) {
	if err != nil {
		// Skip context.Canceled errors, not indicative of API issues
		if errors.Is(err, context.Canceled) {
			return
		}

		// skip 404 errors（subject 不存在是正常业务态）
		if strings.Contains(err.Error(), "404") {
			return
		}
		// skip 429 errors（client 已做退避重试，重试耗尽后不视为宕机信号）
		if strings.Contains(err.Error(), "429") {
			return
		}

		// handle invalid token
		if isBangumiAuthError(err) {
			events.GlobalWSEventManager.SendEvent(events.ServerLoggedOutAnilist, "Your Bangumi session has expired. Please log in again.")
			if c.logoutFunc != nil {
				go c.logoutFunc()
			}
			return
		}

		// Add failure to tracking
		addFailureRecord(err)

		// Only mark as down if we have enough recent failures and are currently marked as working
		if IsWorking.Load() {
			recentFailures := getRecentFailureCount()
			if recentFailures >= failureThreshold {
				c.logger.Warn().
					Err(err).
					Int("recent_failures", recentFailures).
					Dur("within_window", failureWindow).
					Msg("bangumi cache: Multiple API failures detected, switching to cache-only mode.")
				events.GlobalWSEventManager.SendEvent(events.WarningToast,
					fmt.Sprintf("The Bangumi API is experiencing issues (%d failures in %v), switching to cache-only mode.",
						recentFailures, failureWindow))
				IsWorking.Store(false)
			} else {
				c.logger.Debug().
					Err(err).
					Int("recent_failures", recentFailures).
					Int("threshold", failureThreshold).
					Msg("bangumi cache: API failure recorded, monitoring for more failures")
			}
		}
	} else {
		// clear failure tracking and mark as working if not already
		if !IsWorking.Load() {
			c.logger.Info().Msg("bangumi cache: API client is working again, switching back to network-first mode.")
			events.GlobalWSEventManager.SendEvent(events.InfoToast, "The Bangumi API is back online")
			IsWorking.Store(true)
		}
		clearFailureTracking()
	}
}

func isBangumiAuthError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, bangumi.ErrUnauthorized) {
		return true
	}
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "invalid token") || strings.Contains(errStr, "unauthorized") || strings.Contains(errStr, "user not found")
}

// generateCacheKey generates a cache key from the given parameters
func (c *CacheLayer) generateCacheKey(params ...interface{}) string {
	var keyParts []string
	for _, param := range params {
		if param == nil {
			keyParts = append(keyParts, "nil")
			continue
		}
		switch v := param.(type) {
		case *int:
			if v != nil {
				keyParts = append(keyParts, strconv.Itoa(*v))
			} else {
				keyParts = append(keyParts, "nil")
			}
		case *string:
			if v != nil {
				keyParts = append(keyParts, *v)
			} else {
				keyParts = append(keyParts, "nil")
			}
		case string:
			keyParts = append(keyParts, v)
		case int:
			keyParts = append(keyParts, strconv.Itoa(v))
		default:
			keyParts = append(keyParts, fmt.Sprintf("%v", param))
		}
	}
	return lo.Reduce(keyParts, func(acc, item string, _ int) string {
		if acc == "" {
			return item
		}
		return acc + "-" + item
	}, "")
}

// isInCollection checks if a media ID is in the user's collection
func (c *CacheLayer) isInCollection(mediaID int) bool {
	// Update collection tracking if needed
	c.updateCollectionTracking()
	_, ok := c.collectionMediaIDs.Get(mediaID)
	return ok
}

// updateCollectionTracking updates the collection media IDs tracking
func (c *CacheLayer) updateCollectionTracking() {
	if time.Since(c.lastCollectionUpdate) < collectionUpdateInterval {
		return
	}

	go func() {
		defer func() {
			c.lastCollectionUpdate = time.Now()
		}()

		if c.client() == nil {
			return
		}

		// 动画（subject_type=2）与书籍（subject_type=1）两个分区都跟踪
		for _, subjectType := range []int{bangumi.SubjectAnime, bangumi.SubjectBook} {
			res, err := c.fetchAllUserCollections(context.Background(), "-", subjectType)
			if err != nil {
				continue
			}
			for _, uc := range res.Data {
				c.collectionMediaIDs.Set(uc.SubjectID, struct{}{})
			}
		}
	}()
}

// networkFirstGet performs a network-first get operation with caching
func networkFirstGet[T any](c *CacheLayer, bucketName string, cacheKey string, networkFn func() (*T, error)) (*T, error) {
	if !ShouldCache.Load() || c.fileCacher == nil {
		return networkFn()
	}

	bucket := c.buckets[bucketName]

	// Try network first if API is working
	if IsWorking.Load() {
		res, err := networkFn()
		c.checkAndUpdateWorkingState(err)

		if err == nil && res != nil {
			// Cache the successful result
			if err := c.fileCacher.SetPerm(bucket, cacheKey, res); err != nil {
				c.logger.Warn().Err(err).Msg("bangumi cache: Failed to cache result")
			}
			return res, nil
		}
	} else {
		// If API is not working, try it in the background to check if it's back
		go func() {
			res, err := networkFn()
			c.checkAndUpdateWorkingState(err)
			if err == nil && res != nil {
				// Cache the result for future use
				if err := c.fileCacher.SetPerm(bucket, cacheKey, res); err != nil {
					c.logger.Warn().Err(err).Msg("bangumi cache: Failed to cache background result")
				}
			}
		}()
	}

	// Fall back to cache
	var cached T
	found, err := c.fileCacher.GetPerm(bucket, cacheKey, &cached)
	if err != nil {
		return nil, fmt.Errorf("cache lookup failed: %w", err)
	}
	if !found {
		return nil, fmt.Errorf("no cached data available")
	}

	c.logger.Debug().Str("bucket", bucketName).Str("key", cacheKey).Msg("bangumi cache: Serving from cache")
	return &cached, nil
}

// boundedCacheSet caches data with a limit on non-collection entries
func (c *CacheLayer) boundedCacheSet(bucketName string, cacheKey string, data interface{}, mediaID int) error {
	if !ShouldCache.Load() || c.fileCacher == nil {
		return nil
	}

	bucket := c.buckets[bucketName]

	// Always cache collection media
	if c.isInCollection(mediaID) {
		return c.fileCacher.SetPerm(bucket, cacheKey, data)
	}

	// For non-collection media, enforce the limit
	allData, err := filecache.GetAll[interface{}](c.fileCacher, filecache.NewBucket(bucket.Name(), 0))
	if err != nil {
		return err
	}

	// If we're at the limit, remove the oldest entry (simple FIFO for now)
	if len(allData) >= maxNonCollectionMediaCacheEntries {
		// Remove the first key we find (this is a simple implementation)
		for key := range allData {
			if err := c.fileCacher.DeletePerm(bucket, key); err == nil {
				break
			}
		}
	}

	return c.fileCacher.SetPerm(bucket, cacheKey, data)
}

// updateCollectionTrackingFromAnimeCollection updates collection tracking from anime collection
func (c *CacheLayer) updateCollectionTrackingFromAnimeCollection(collection *media.AnimeCollection) {
	if !ShouldCache.Load() || collection == nil || collection.MediaListCollection == nil {
		return
	}

	for _, list := range collection.MediaListCollection.Lists {
		if list != nil {
			for _, entry := range list.Entries {
				if entry != nil && entry.Media != nil {
					c.collectionMediaIDs.Set(entry.Media.ID, struct{}{})
				}
			}
		}
	}
	c.lastCollectionUpdate = time.Now()
}

func (c *CacheLayer) updateCollectionTrackingFromAnimeCollectionWithRelations(collection *media.AnimeCollectionWithRelations) {
	if !ShouldCache.Load() || collection == nil || collection.MediaListCollection == nil {
		return
	}

	for _, list := range collection.MediaListCollection.Lists {
		if list != nil {
			for _, entry := range list.Entries {
				if entry != nil && entry.Media != nil {
					c.collectionMediaIDs.Set(entry.Media.ID, struct{}{})
				}
			}
		}
	}
	c.lastCollectionUpdate = time.Now()
}

func (c *CacheLayer) updateCollectionTrackingFromMangaCollection(collection *media.MangaCollection) {
	if !ShouldCache.Load() || collection == nil || collection.MediaListCollection == nil {
		return
	}

	for _, list := range collection.MediaListCollection.Lists {
		if list != nil {
			for _, entry := range list.Entries {
				if entry != nil && entry.Media != nil {
					c.collectionMediaIDs.Set(entry.Media.ID, struct{}{})
				}
			}
		}
	}
	c.lastCollectionUpdate = time.Now()
}

// invalidateMediaCaches invalidates caches for a specific media ID
func (c *CacheLayer) invalidateMediaCaches(mediaID int) {
	if !ShouldCache.Load() || c.fileCacher == nil {
		return
	}

	mediaIDStr := strconv.Itoa(mediaID)

	// Delete from all media-specific buckets
	buckets := []string{
		BaseAnimeBucket,
		CompleteAnimeBucket,
		AnimeDetailsBucket,
		BaseMangaBucket,
		MangaDetailsBucket,
	}

	for _, bucketName := range buckets {
		bucket := c.buckets[bucketName]
		if err := c.fileCacher.DeletePerm(bucket, mediaIDStr); err != nil {
			c.logger.Debug().Err(err).Str("bucket", bucketName).Int("mediaID", mediaID).Msg("bangumi cache: Failed to invalidate cache entry")
		}
	}
}

// invalidateCollectionCaches invalidates all collection caches
func (c *CacheLayer) invalidateCollectionCaches() {
	if !ShouldCache.Load() || c.fileCacher == nil {
		return
	}

	collectionBuckets := []string{
		AnimeCollectionBucket,
		MangaCollectionBucket,
	}

	for _, bucketName := range collectionBuckets {
		bucket := c.buckets[bucketName]
		if err := c.fileCacher.EmptyPerm(bucket); err != nil {
			c.logger.Warn().Err(err).Str("bucket", bucketName).Msg("bangumi cache: Failed to invalidate collection cache")
		}
	}

	// Reset collection tracking
	c.collectionMediaIDs.Clear()
	c.lastCollectionUpdate = time.Time{}
}

// extractBaseAnimeFromCollection attempts to extract Anime data from cached anime collection
func (c *CacheLayer) extractBaseAnimeFromCollection(mediaID int) *media.Anime {
	if c.fileCacher == nil {
		return nil
	}

	bucket := c.buckets[AnimeCollectionBucket]
	cacheKey := c.generateCacheKey("collection", nil)
	var animeCollection media.AnimeCollection
	found, err := c.fileCacher.GetPerm(bucket, cacheKey, &animeCollection)
	if err == nil && found && animeCollection.MediaListCollection != nil {
		for _, list := range animeCollection.MediaListCollection.Lists {
			if list != nil {
				for _, entry := range list.Entries {
					if entry != nil && entry.Media != nil && entry.Media.ID == mediaID {
						return entry.Media
					}
				}
			}
		}
	}

	return nil
}

// extractBaseMangaFromCollection attempts to extract Manga data from cached manga collection
func (c *CacheLayer) extractBaseMangaFromCollection(mediaID int) *media.Manga {
	if c.fileCacher == nil {
		return nil
	}

	bucket := c.buckets[MangaCollectionBucket]
	cacheKey := c.generateCacheKey("collection", nil)
	var mangaCollection media.MangaCollection
	found, err := c.fileCacher.GetPerm(bucket, cacheKey, &mangaCollection)
	if err == nil && found && mangaCollection.MediaListCollection != nil {
		for _, list := range mangaCollection.MediaListCollection.Lists {
			if list != nil {
				for _, entry := range list.Entries {
					if entry != nil && entry.Media != nil && entry.Media.ID == mediaID {
						return entry.Media
					}
				}
			}
		}
	}

	return nil
}

//----------------------------------------------------------------------------------------------------------------------
// 查询方法

// AnimeCollection 获取动画收藏（subject_type=2 全分页 → media.AnimeCollection）。
func (c *CacheLayer) AnimeCollection(ctx context.Context, username string) (*media.AnimeCollection, error) {
	cacheKey := c.generateCacheKey("collection", username)
	res, err := networkFirstGet(c, AnimeCollectionBucket, cacheKey, func() (*media.AnimeCollection, error) {
		data, err := c.fetchAllUserCollections(ctx, username, bangumi.SubjectAnime)
		if err != nil {
			return nil, err
		}
		return bangumi.AnimeCollectionFromUserCollections(data), nil
	})

	if err == nil && res != nil && c.applyQueuedUpdatesToAnimeCollection(res) {
		if c.fileCacher != nil {
			if err := c.fileCacher.SetPerm(c.buckets[AnimeCollectionBucket], cacheKey, res); err != nil {
				c.logger.Warn().Err(err).Msg("bangumi cache: Failed to apply queued updates to anime collection cache")
			}
		}
	}

	// Update collection tracking with the fetched data
	if err == nil && res != nil {
		go c.updateCollectionTrackingFromAnimeCollection(res)
	}

	return res, err
}

// MangaCollection 获取漫画/书籍收藏（subject_type=1 全分页 → media.MangaCollection）。
func (c *CacheLayer) MangaCollection(ctx context.Context, username string) (*media.MangaCollection, error) {
	cacheKey := c.generateCacheKey("collection", username)
	res, err := networkFirstGet(c, MangaCollectionBucket, cacheKey, func() (*media.MangaCollection, error) {
		data, err := c.fetchAllUserCollections(ctx, username, bangumi.SubjectBook)
		if err != nil {
			return nil, err
		}
		return bangumi.MangaCollectionFromUserCollections(data), nil
	})

	if err == nil && res != nil && c.applyQueuedUpdatesToMangaCollection(res) {
		if c.fileCacher != nil {
			if err := c.fileCacher.SetPerm(c.buckets[MangaCollectionBucket], cacheKey, res); err != nil {
				c.logger.Warn().Err(err).Msg("bangumi cache: Failed to apply queued updates to manga collection cache")
			}
		}
	}

	// Update collection tracking with the fetched data
	if err == nil && res != nil {
		go c.updateCollectionTrackingFromMangaCollection(res)
	}

	return res, err
}

// BaseAnimeByID 获取单个条目（subject → media.Anime）。
func (c *CacheLayer) BaseAnimeByID(ctx context.Context, id int) (*media.Anime, error) {
	cacheKey := c.generateCacheKey(id)
	res, err := networkFirstGet(c, BaseAnimeBucket, cacheKey, func() (*media.Anime, error) {
		subject, err := c.client().GetSubject(ctx, id)
		if err != nil {
			return nil, err
		}
		return media.AnimeFromSubject(bangumi.SubjectToMedia(subject)), nil
	})

	// If network and direct cache failed, try to extract from collection cache
	if err != nil {
		if collectionResult := c.extractBaseAnimeFromCollection(id); collectionResult != nil {
			c.logger.Debug().Int("mediaID", id).Msg("bangumi cache: Extracted Anime from collection cache")
			return collectionResult, nil
		}
	}

	// If successful, update bounded cache for non-collection media
	if err == nil && res != nil {
		go func() {
			if err := c.boundedCacheSet(BaseAnimeBucket, cacheKey, res, id); err != nil {
				c.logger.Warn().Err(err).Msg("bangumi cache: Failed to update bounded cache")
			}
		}()
	}

	return res, err
}

// BaseMangaByID 获取单个条目（subject → media.Manga）。
func (c *CacheLayer) BaseMangaByID(ctx context.Context, id int) (*media.Manga, error) {
	cacheKey := c.generateCacheKey(id)
	res, err := networkFirstGet(c, BaseMangaBucket, cacheKey, func() (*media.Manga, error) {
		subject, err := c.client().GetSubject(ctx, id)
		if err != nil {
			return nil, err
		}
		return media.MangaFromSubject(bangumi.SubjectToMedia(subject)), nil
	})

	// If network and direct cache failed, try to extract from collection cache
	if err != nil {
		if collectionResult := c.extractBaseMangaFromCollection(id); collectionResult != nil {
			c.logger.Debug().Int("mediaID", id).Msg("bangumi cache: Extracted Manga from collection cache")
			return collectionResult, nil
		}
	}

	// If successful, update bounded cache for non-collection media
	if err == nil && res != nil {
		go func() {
			if err := c.boundedCacheSet(BaseMangaBucket, cacheKey, res, id); err != nil {
				c.logger.Warn().Err(err).Msg("bangumi cache: Failed to update bounded cache")
			}
		}()
	}

	return res, err
}

// CompleteAnimeByID 获取带关系树的完整条目（subject + related subjects）。
func (c *CacheLayer) CompleteAnimeByID(ctx context.Context, id int) (*media.CompleteAnime, error) {
	cacheKey := c.generateCacheKey(id)
	res, err := networkFirstGet(c, CompleteAnimeBucket, cacheKey, func() (*media.CompleteAnime, error) {
		subject, err := c.client().GetSubject(ctx, id)
		if err != nil {
			return nil, err
		}
		related, err := c.client().GetRelatedSubjects(ctx, id)
		if err != nil {
			// 关系列表失败不阻塞主条目返回（降级为无关系）
			related = nil
		}
		return bangumi.CompleteAnimeFromSubject(subject, related), nil
	})

	// If successful, update bounded cache for non-collection media
	if err == nil && res != nil {
		go func() {
			if err := c.boundedCacheSet(CompleteAnimeBucket, cacheKey, res, id); err != nil {
				c.logger.Warn().Err(err).Msg("bangumi cache: failed to update bounded cache")
			}
		}()
	}

	return res, err
}

// AnimeDetailsByID 获取 anime 详情（infobox 尽力映射）。
func (c *CacheLayer) AnimeDetailsByID(ctx context.Context, id int) (*media.AnimeDetails, error) {
	cacheKey := c.generateCacheKey(id)
	res, err := networkFirstGet(c, AnimeDetailsBucket, cacheKey, func() (*media.AnimeDetails, error) {
		subject, err := c.client().GetSubject(ctx, id)
		if err != nil {
			return nil, err
		}
		return bangumi.AnimeDetailsFromSubject(subject), nil
	})

	// If successful, update bounded cache for non-collection media
	if err == nil && res != nil {
		go func() {
			if err := c.boundedCacheSet(AnimeDetailsBucket, cacheKey, res, id); err != nil {
				c.logger.Warn().Err(err).Msg("bangumi cache: failed to update bounded cache")
			}
		}()
	}

	return res, err
}

// MangaDetailsByID 获取 manga 详情（infobox 尽力映射）。
func (c *CacheLayer) MangaDetailsByID(ctx context.Context, id int) (*media.MangaDetails, error) {
	cacheKey := c.generateCacheKey(id)
	res, err := networkFirstGet(c, MangaDetailsBucket, cacheKey, func() (*media.MangaDetails, error) {
		subject, err := c.client().GetSubject(ctx, id)
		if err != nil {
			return nil, err
		}
		return bangumi.MangaDetailsFromSubject(subject), nil
	})

	// If successful, update bounded cache for non-collection media
	if err == nil && res != nil {
		go func() {
			if err := c.boundedCacheSet(MangaDetailsBucket, cacheKey, res, id); err != nil {
				c.logger.Warn().Err(err).Msg("bangumi cache: failed to update bounded cache")
			}
		}()
	}

	return res, err
}

//----------------------------------------------------------------------------------------------------------------------
// 变更方法

// UpdateMediaListEntry 更新收藏条目（状态/评分/进度）。
// Bangumi 无条目级起止日期字段，startedAt/completedAt 被忽略（TODO: 确认服务端后续支持）。
func (c *CacheLayer) UpdateMediaListEntry(ctx context.Context, mediaID int, status *media.MediaListStatus, scoreRaw *int, progress *int, startedAt *media.FuzzyDateInput, completedAt *media.FuzzyDateInput) error {
	// Mutations require the API to be working
	if !IsWorking.Load() {
		return c.queueMediaListEntryUpdate(mediaID, status, scoreRaw, progress, startedAt, completedAt)
	}

	err := c.sendMediaListEntryUpdate(ctx, mediaID, status, scoreRaw, progress, startedAt, completedAt)
	c.checkAndUpdateWorkingState(err)
	if err != nil && shouldQueueMediaListUpdate(err) {
		return c.queueMediaListEntryUpdate(mediaID, status, scoreRaw, progress, startedAt, completedAt)
	}

	// Invalidate relevant caches on successful mutation
	if err == nil {
		c.invalidateMediaCaches(mediaID)
		c.invalidateCollectionCaches()
	}

	return err
}

// UpdateMediaListEntryProgress 更新章节进度。
func (c *CacheLayer) UpdateMediaListEntryProgress(ctx context.Context, mediaID int, progress int, status *media.MediaListStatus) error {
	// Mutations require the API to be working
	if !IsWorking.Load() {
		return c.queueMediaListEntryProgressUpdate(mediaID, progress, status)
	}

	err := c.sendMediaListEntryProgressUpdate(ctx, mediaID, progress, status)
	c.checkAndUpdateWorkingState(err)
	if err != nil && shouldQueueMediaListUpdate(err) {
		return c.queueMediaListEntryProgressUpdate(mediaID, progress, status)
	}

	// Invalidate relevant caches on successful mutation
	if err == nil {
		c.invalidateMediaCaches(mediaID)
		c.invalidateCollectionCaches()
	}

	return err
}

// UpdateMediaListEntryRepeat 更新重看次数。
// Bangumi 收藏无重看次数字段，此操作为 no-op（保留接口兼容）。
// TODO: 如 Bangumi 后续支持重看标记，在此接入。
func (c *CacheLayer) UpdateMediaListEntryRepeat(ctx context.Context, mediaID int, repeat int) error {
	c.logger.Debug().Int("mediaID", mediaID).Int("repeat", repeat).Msg("bangumi cache: Repeat count is not supported by Bangumi, skipping")
	return nil
}

// DeleteEntry 删除收藏条目。
func (c *CacheLayer) DeleteEntry(ctx context.Context, mediaID int) error {
	// Mutations require the API to be working
	if !IsWorking.Load() {
		return fmt.Errorf("bangumi cache: API client is not working, mutation operations are not available")
	}

	err := c.client().DeleteCollection(ctx, mediaID)
	c.checkAndUpdateWorkingState(err)

	// Invalidate collection caches on successful deletion
	if err == nil {
		c.invalidateMediaCaches(mediaID)
		c.invalidateCollectionCaches()
	}

	return err
}

func (c *CacheLayer) sendMediaListEntryUpdate(ctx context.Context, mediaID int, status *media.MediaListStatus, scoreRaw *int, progress *int, startedAt *media.FuzzyDateInput, completedAt *media.FuzzyDateInput) error {
	err := c.client().UpsertCollection(ctx, mediaID, collectionUpsertBodyFromValues(status, scoreRaw, progress))
	if err == nil {
		c.deleteQueuedUpdate(mediaID)
	}
	return err
}

func (c *CacheLayer) sendMediaListEntryProgressUpdate(ctx context.Context, mediaID int, progress int, status *media.MediaListStatus) error {
	// 进度更新走 PATCH（局部更新，避免覆盖用户已有的评分/评论）
	body := bangumi.CollectionUpsertBody{}
	if status != nil {
		if t, ok := bangumi.MediaListStatusToBangumiType(*status); ok {
			body.Type = &t
		}
	}
	body.EpStatus = &progress
	err := c.client().PatchCollection(ctx, mediaID, body)
	if err == nil {
		c.deleteQueuedUpdate(mediaID)
	}
	return err
}
