// Package bangumi_platform 基于 Bangumi API 的 Platform 接口实现（Phase 2 换锚 Wave B）。
//
// 数据源：GET /v0/users/{username}/collections（收藏）、GET /v0/subjects/{id}（条目）、
// GET /v0/subjects/{id}/subjects（关系树）、POST/PATCH/DELETE 收藏写操作。
//
// 降级策略（M4 决策，Bangumi 数据面缺口）：
//   - GetAnimeAiringSchedule：无逐集播出时间戳 → 空表 + TODO；
//   - GetViewerStats：无服务端统计端点 → 由本地收藏集合合成最小统计 + TODO；
//   - GetStudioDetails：Bangumi 无制作公司 ID 体系 → infobox「制作」名称 +
//     伪 ID 注册表 + TODO；
//   - UpdateEntryRepeat：Bangumi 收藏无重看次数字段 → no-op；
//   - 条目级起止日期（startedAt/completedAt）：收藏 upsert 请求体无此字段 → 忽略。
package bangumi_platform

import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"seanime/internal/api/bangumi"
	"seanime/internal/customsource"
	"seanime/internal/database/db"
	"seanime/internal/extension"
	"seanime/internal/hook"
	"seanime/internal/media"
	"seanime/internal/platforms/platform"
	"seanime/internal/platforms/shared_platform"
	"seanime/internal/util"
	"sync"

	"github.com/rs/zerolog"
	"github.com/samber/mo"
)

type (
	BangumiPlatform struct {
		logger     *zerolog.Logger
		username   mo.Option[string]
		// selfUsername 经 GetMe 解析的本人 username 句柄（懒加载缓存）。
		// Bangumi 的 username 句柄 ≠ 昵称，GET /v0/users/{句柄}/collections 必须用句柄。
		selfUsername mo.Option[string]
		client     *bangumi.Client
		cacheLayer *shared_platform.CacheLayer

		animeCollection    mo.Option[*media.AnimeCollection]
		rawAnimeCollection mo.Option[*media.AnimeCollection]
		mangaCollection    mo.Option[*media.MangaCollection]
		rawMangaCollection mo.Option[*media.MangaCollection]

		helper *shared_platform.PlatformHelper
		db     *db.Database

		// mappingService 可选注入的 bangumi→anidb 名称映射服务（编排层装配）。
		// nil 时 GetAnimeByMalID 仅依赖 reverseIndex 中已确认的映射。
		mappingService media.MappingService
		// reverseIndex 反向 ID 索引（anidb/mal → bangumi），机会式填充，见 media.ReverseIDIndex。
		reverseIndex *media.ReverseIDIndex

		// studioNames 制作公司伪 ID → 名称注册表（GetAnimeDetails 时填充，
		// GetStudioDetails 时反查；进程内有效，重启后由详情页重建）。
		studioNames sync.Map // map[int]string

		collectionMu sync.Mutex // 保护 collection 选项的刷新
	}
)

// NewBangumiPlatform 创建 Bangumi 平台实现。
// client 由核心层构造（带 token 与 filecache）；cacheDir 用于 CacheLayer 的持久缓存桶。
func NewBangumiPlatform(client *bangumi.Client, cacheDir string, extensionBankRef *util.Ref[*extension.UnifiedBank], logger *zerolog.Logger, db *db.Database, logoutFunc ...func()) platform.Platform {
	bp := &BangumiPlatform{
		client:       client,
		cacheLayer:   shared_platform.NewCacheLayer(client, cacheDir, logoutFunc...),
		logger:       logger,
		username:     mo.None[string](),
		helper:       shared_platform.NewPlatformHelper(extensionBankRef, db, logger),
		db:           db,
		reverseIndex: media.NewReverseIDIndex(),
	}

	return bp
}

// SetMappingService 注入 bangumi→anidb 映射服务（编排层装配，additive）。
func (bp *BangumiPlatform) SetMappingService(svc media.MappingService) {
	bp.mappingService = svc
}

var _ platform.Platform = (*BangumiPlatform)(nil)

func (bp *BangumiPlatform) ClearCache() {
	bp.helper.ClearCache()
	// TODO: bangumi.Client 的 filecache 无对外清空接口，条目缓存由 24h TTL 自然过期
}

func (bp *BangumiPlatform) Close() {
	bp.helper.Close()
	// 停掉 client 的节流 ticker（Close 后 client 不可再用，调用时机为应用关闭）
	if bp.client != nil {
		bp.client.Close()
	}
}

func (bp *BangumiPlatform) GetBangumiClient() *bangumi.Client {
	return bp.client
}

// GetReverseIDIndex 暴露反向索引（编排层/清偿工具回填用，additive）。
func (bp *BangumiPlatform) GetReverseIDIndex() *media.ReverseIDIndex {
	return bp.reverseIndex
}

func (bp *BangumiPlatform) SetUsername(username string) {
	// 存储的仅是展示用昵称（前端 username 语义），不作 API 路径参数；
	// API 路径参数由 getUsername(ctx) 经 GetMe 解析的真实句柄提供。
	bp.username = mo.Some(username)
}

// Username 返回本人 username 句柄（GetMe 懒解析 + 进程内缓存），解析失败返回空串。
// 供 handlers 等直接调 client 的方使用：Bangumi 的 username 句柄 ≠ 昵称，
// GET /v0/users/{句柄}/collections 必须用此句柄，不能用 App.GetUsername()（昵称，404）。
func (bp *BangumiPlatform) Username(ctx context.Context) string {
	if u, ok := bp.getUsername(ctx); ok {
		return u
	}
	return ""
}

func (bp *BangumiPlatform) getUsername(ctx context.Context) (string, bool) {
	if bp.selfUsername.IsPresent() {
		return bp.selfUsername.MustGet(), true
	}
	if bp.client == nil {
		return "", false
	}

	// 首次调用经 GetMe 解析本人句柄并缓存（进程内有效）
	me, err := bp.client.GetMe(ctx)
	if err != nil {
		bp.logger.Warn().Err(err).Msg("bangumi platform: GetMe 解析 username 失败")
		return "", false
	}
	bp.selfUsername = mo.Some(me.Username)
	return me.Username, true
}

func (bp *BangumiPlatform) UpdateEntry(ctx context.Context, mediaID int, status *media.MediaListStatus, scoreRaw *int, progress *int, startedAt *media.FuzzyDateInput, completedAt *media.FuzzyDateInput) error {
	bp.logger.Trace().Msg("bangumi platform: Updating entry")

	// Use shared hook handling
	return bp.helper.TriggerUpdateEntryHooks(ctx, mediaID, status, scoreRaw, progress, startedAt, completedAt, func(event *platform.PreUpdateEntryEvent) error {
		// Check if this is a custom source entry (after hooks have been triggered)
		if handled, err := bp.helper.HandleCustomSourceUpdateEntry(ctx, mediaID, event.Status, event.ScoreRaw, event.Progress, event.StartedAt, event.CompletedAt); handled {
			return err
		}

		// Bangumi 收藏请求体无条目级起止日期字段，日期参数被忽略
		return bp.cacheLayer.UpdateMediaListEntry(ctx, *event.MediaID, event.Status, event.ScoreRaw, event.Progress, event.StartedAt, event.CompletedAt)
	})
}

func (bp *BangumiPlatform) UpdateEntryProgress(ctx context.Context, mediaID int, progress int, totalCount *int) error {
	bp.logger.Trace().Msg("bangumi platform: Updating entry progress")

	// Use shared hook handling
	return bp.helper.TriggerUpdateEntryProgressHooks(ctx, mediaID, progress, totalCount, func(event *platform.PreUpdateEntryProgressEvent) error {
		// Check if this is a custom source entry (after hooks have been triggered)
		if handled, err := bp.helper.HandleCustomSourceUpdateEntryProgress(ctx, mediaID, *event.Progress, event.TotalCount); handled {
			return err
		}

		realTotalCount := 0
		if totalCount != nil && *totalCount > 0 {
			realTotalCount = *totalCount
		}

		// Check if the anime is in the repeating list
		// If it is, set the status to repeating
		// （Bangumi 无重看状态，REPEATING 写入时由映射层保守转为「看过」）
		if bp.rawAnimeCollection.IsPresent() {
			for _, list := range bp.rawAnimeCollection.MustGet().MediaListCollection.Lists {
				if list.Status != nil && *list.Status == media.MediaListStatusRepeating {
					if list.Entries != nil {
						for _, entry := range list.Entries {
							if entry.GetMedia().GetID() == mediaID {
								*event.Status = media.MediaListStatusRepeating
								break
							}
						}
					}
				}
			}
		}
		if realTotalCount > 0 && *event.Progress >= realTotalCount {
			*event.Status = media.MediaListStatusCompleted
		}

		if realTotalCount > 0 && *event.Progress > realTotalCount {
			*event.Progress = realTotalCount
		}

		return bp.cacheLayer.UpdateMediaListEntryProgress(
			ctx,
			*event.MediaID,
			*event.Progress,
			event.Status,
		)
	})
}

func (bp *BangumiPlatform) UpdateEntryRepeat(ctx context.Context, mediaID int, repeat int) error {
	bp.logger.Trace().Msg("bangumi platform: Updating entry repeat")

	// Use shared hook handling
	return bp.helper.TriggerUpdateEntryRepeatHooks(ctx, mediaID, repeat, func(event *platform.PreUpdateEntryRepeatEvent) error {
		// Check if this is a custom source entry (after hooks have been triggered)
		if handled, err := bp.helper.HandleCustomSourceUpdateEntryRepeat(ctx, mediaID, *event.Repeat); handled {
			return err
		}

		// Bangumi 收藏无重看次数字段 → no-op（保留 hook 语义）
		return bp.cacheLayer.UpdateMediaListEntryRepeat(ctx, *event.MediaID, *event.Repeat)
	})
}

func (bp *BangumiPlatform) DeleteEntry(ctx context.Context, mediaID, entryId int) error {
	bp.logger.Trace().Msg("bangumi platform: Deleting entry")

	return bp.helper.TriggerDeleteEntryHooks(ctx, mediaID, entryId, func(event *platform.PreDeleteEntryEvent) error {
		if handled, err := bp.helper.HandleCustomSourceDeleteEntry(ctx, *event.MediaID, *event.EntryID); handled {
			return err
		}

		// Bangumi 收藏无条目 ID 概念，按 subject ID 删除（entryId 即 mediaID）
		return bp.cacheLayer.DeleteEntry(ctx, *event.MediaID)
	})
}

func (bp *BangumiPlatform) GetAnime(ctx context.Context, mediaID int) (*media.Anime, error) {
	bp.logger.Trace().Int("mediaId", mediaID).Msg("bangumi platform: Fetching anime")

	if cachedAnime, ok := bp.helper.GetCachedBaseAnime(mediaID); ok {
		bp.logger.Trace().Msg("bangumi platform: Returning anime from cache")
		return bp.helper.TriggerGetAnimeEvent(cachedAnime)
	}

	// Check if this is a custom source entry
	if m, isCustom, err := bp.helper.HandleCustomSourceAnime(ctx, mediaID); isCustom {
		if err != nil {
			return nil, err
		}

		triggeredMedia, err := bp.helper.TriggerGetAnimeEvent(m)
		if err != nil {
			return nil, err
		}

		bp.helper.SetCachedBaseAnime(mediaID, triggeredMedia)
		return triggeredMedia, nil
	}

	// Get from Bangumi
	ret, err := bp.cacheLayer.BaseAnimeByID(ctx, mediaID)
	if err != nil {
		return nil, err
	}

	bp.recordReverseMapping(ctx, mediaID)

	triggeredMedia, err := bp.helper.TriggerGetAnimeEvent(ret)
	if err != nil {
		return nil, err
	}

	bp.helper.SetCachedBaseAnime(mediaID, triggeredMedia)
	return triggeredMedia, nil
}

// recordReverseMapping 机会式填充反向索引：条目命中后尝试用映射服务确认
// bangumi → anidb 对应关系（GetAnimeByMalID 的查询基础）。
func (bp *BangumiPlatform) recordReverseMapping(ctx context.Context, mediaID int) {
	if bp.mappingService == nil || bp.client == nil {
		return
	}
	subject, err := bp.client.GetSubject(ctx, mediaID)
	if err != nil {
		return
	}
	anidbID, ok := bp.mappingService.ResolveBangumiToAniDB(subject.ID, media.NameSet{
		Native:  subject.Name,
		Chinese: subject.NameCN,
	})
	if ok {
		bp.reverseIndex.Put(subject.ID, anidbID, 0)
	}
}

func (bp *BangumiPlatform) GetAnimeByMalID(ctx context.Context, malID int) (*media.Anime, error) {
	bp.logger.Trace().Msg("bangumi platform: Fetching anime by MAL ID")

	// 反向查找链：mal → bangumi（已确认映射）→ mal → anidb（animap 直查）→ anidb → bangumi
	bangumiID, ok := bp.reverseIndex.BangumiByMal(malID)
	if !ok && bp.mappingService != nil {
		if anidbID, ok2 := bp.mappingService.ResolveMalToAniDB(malID); ok2 {
			bangumiID, ok = bp.reverseIndex.BangumiByAniDB(anidbID)
		}
	}
	if !ok || bangumiID <= 0 {
		return nil, fmt.Errorf("bangumi platform: no confirmed bangumi mapping for MAL ID %d", malID)
	}

	ret, err := bp.cacheLayer.BaseAnimeByID(ctx, bangumiID)
	if err != nil {
		return nil, err
	}

	return bp.helper.TriggerGetAnimeEvent(ret)
}

func (bp *BangumiPlatform) GetAnimeDetails(ctx context.Context, mediaID int) (*media.AnimeDetails, error) {
	bp.logger.Trace().Int("mediaId", mediaID).Msg("bangumi platform: Fetching anime details")

	// Check if this is a custom source entry
	if m, isCustom, err := bp.helper.HandleCustomSourceAnimeDetails(ctx, mediaID); isCustom {
		if err != nil {
			return nil, err
		}
		return bp.helper.TriggerGetAnimeDetailsEvent(m)
	}

	// Get from Bangumi
	ret, err := bp.cacheLayer.AnimeDetailsByID(ctx, mediaID)
	if err != nil {
		return nil, err
	}

	// 制作公司无官方 ID，分配稳定伪 ID（名称 FNV-1a 哈希）并登记，供 GetStudioDetails 反查
	if ret.Studios != nil {
		for _, node := range ret.Studios.Nodes {
			if node == nil || node.Name == "" {
				continue
			}
			node.ID = studioPseudoID(node.Name)
			bp.studioNames.Store(node.ID, node.Name)
		}
	}

	return bp.helper.TriggerGetAnimeDetailsEvent(ret)
}

// studioPseudoID 由公司名生成稳定非零伪 ID。
// Bangumi 无制作公司 ID 体系；哈希 ID 仅用于前端路由键，不落库。
func studioPseudoID(name string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(name))
	id := int(h.Sum32()) & 0x7fffffff // 保证非负
	if id == 0 {
		id = 1
	}
	return id
}

func (bp *BangumiPlatform) GetAnimeWithRelations(ctx context.Context, mediaID int) (*media.CompleteAnime, error) {
	bp.logger.Trace().Int("mediaId", mediaID).Msg("bangumi platform: Fetching anime with relations")

	if cachedAnime, ok := bp.helper.GetCachedCompleteAnime(mediaID); ok {
		bp.logger.Trace().Msg("bangumi platform: Cache HIT for anime with relations")
		return cachedAnime, nil
	}

	// Check if this is a custom source entry
	if m, isCustom, err := bp.helper.HandleCustomSourceAnimeWithRelations(ctx, mediaID); isCustom {
		if err != nil {
			return nil, err
		}
		bp.helper.SetCachedCompleteAnime(mediaID, m)
		return m, nil
	}

	// Get from Bangumi（subject + related subjects 一次组装）
	ret, err := bp.cacheLayer.CompleteAnimeByID(ctx, mediaID)
	if err != nil {
		return nil, err
	}

	bp.helper.SetCachedCompleteAnime(mediaID, ret)
	return ret, nil
}

func (bp *BangumiPlatform) GetManga(ctx context.Context, mediaID int) (*media.Manga, error) {
	bp.logger.Trace().Int("mediaId", mediaID).Msg("bangumi platform: Fetching manga")

	if cachedManga, ok := bp.helper.GetCachedBaseManga(mediaID); ok {
		bp.logger.Trace().Msg("bangumi platform: Returning manga from cache")
		return bp.helper.TriggerGetMangaEvent(cachedManga)
	}

	// Check if this is a custom source entry
	if m, isCustom, err := bp.helper.HandleCustomSourceManga(ctx, mediaID); isCustom {
		if err != nil {
			return nil, err
		}

		triggeredMedia, err := bp.helper.TriggerGetMangaEvent(m)
		if err != nil {
			return nil, err
		}

		bp.helper.SetCachedBaseManga(mediaID, triggeredMedia)
		return triggeredMedia, nil
	}

	// Get from Bangumi
	ret, err := bp.cacheLayer.BaseMangaByID(ctx, mediaID)
	if err != nil {
		return nil, err
	}

	triggeredMedia, err := bp.helper.TriggerGetMangaEvent(ret)
	if err != nil {
		return nil, err
	}

	bp.helper.SetCachedBaseManga(mediaID, triggeredMedia)
	return triggeredMedia, nil
}

func (bp *BangumiPlatform) GetMangaDetails(ctx context.Context, mediaID int) (*media.MangaDetails, error) {
	bp.logger.Trace().Msg("bangumi platform: Fetching manga details")

	// Check if this is a custom source entry
	if m, isCustom, err := bp.helper.HandleCustomSourceMangaDetails(ctx, mediaID); isCustom {
		return m, err
	}

	// Get from Bangumi
	return bp.cacheLayer.MangaDetailsByID(ctx, mediaID)
}

//----------------------------------------------------------------------------------------------------------------------
// 收藏

// GetAnimeCollection 获取动画收藏（不走 custom lists 过滤——Bangumi 无 custom lists）。
func (bp *BangumiPlatform) GetAnimeCollection(ctx context.Context, bypassCache bool) (*media.AnimeCollection, error) {
	if !bypassCache && bp.animeCollection.IsPresent() {
		event := new(platform.GetCachedAnimeCollectionEvent)
		event.AnimeCollection = bp.animeCollection.MustGet()
		err := hook.GlobalHookManager.OnGetCachedAnimeCollection().Trigger(event)
		if err != nil {
			return nil, err
		}
		return event.AnimeCollection, nil
	}

	if _, ok := bp.getUsername(ctx); !ok {
		return nil, nil
	}

	err := bp.refreshAnimeCollection(ctx)
	if err != nil {
		return nil, err
	}

	event := new(platform.GetAnimeCollectionEvent)
	event.AnimeCollection = bp.animeCollection.MustGet()

	err = hook.GlobalHookManager.OnGetAnimeCollection().Trigger(event)
	if err != nil {
		return nil, err
	}

	return event.AnimeCollection, nil
}

// GetRawAnimeCollection 获取带 custom lists 的动画收藏。
// Bangumi 无 custom lists 概念，本方法与 GetAnimeCollection 等价
// （保留独立入口以维持接口兼容；raw 副本保留原始 lists 顺序）。
func (bp *BangumiPlatform) GetRawAnimeCollection(ctx context.Context, bypassCache bool) (*media.AnimeCollection, error) {
	if !bypassCache && bp.rawAnimeCollection.IsPresent() {
		event := new(platform.GetCachedRawAnimeCollectionEvent)
		event.AnimeCollection = bp.rawAnimeCollection.MustGet()
		err := hook.GlobalHookManager.OnGetCachedRawAnimeCollection().Trigger(event)
		if err != nil {
			return nil, err
		}
		return event.AnimeCollection, nil
	}

	if _, ok := bp.getUsername(ctx); !ok {
		return nil, nil
	}

	err := bp.refreshAnimeCollection(ctx)
	if err != nil {
		return nil, err
	}

	event := new(platform.GetRawAnimeCollectionEvent)
	event.AnimeCollection = bp.rawAnimeCollection.MustGet()

	err = hook.GlobalHookManager.OnGetRawAnimeCollection().Trigger(event)
	if err != nil {
		return nil, err
	}

	return event.AnimeCollection, nil
}

func (bp *BangumiPlatform) RefreshAnimeCollection(ctx context.Context) (*media.AnimeCollection, error) {
	if _, ok := bp.getUsername(ctx); !ok {
		return nil, nil
	}

	err := bp.refreshAnimeCollection(ctx)
	if err != nil {
		return nil, err
	}

	event := new(platform.GetAnimeCollectionEvent)
	event.AnimeCollection = bp.animeCollection.MustGet()

	err = hook.GlobalHookManager.OnGetAnimeCollection().Trigger(event)
	if err != nil {
		return nil, err
	}

	event2 := new(platform.GetRawAnimeCollectionEvent)
	event2.AnimeCollection = bp.rawAnimeCollection.MustGet()

	err = hook.GlobalHookManager.OnGetRawAnimeCollection().Trigger(event2)
	if err != nil {
		return nil, err
	}

	return event.AnimeCollection, nil
}

func (bp *BangumiPlatform) refreshAnimeCollection(ctx context.Context) error {
	username, ok := bp.getUsername(ctx)
	if !ok {
		return errors.New("bangumi: Username is not set")
	}

	bp.collectionMu.Lock()
	defer bp.collectionMu.Unlock()

	// Else, get the collection from Bangumi（subject_type=2 全分页，CacheLayer 内完成）
	collection, err := bp.cacheLayer.AnimeCollection(ctx, username)
	if err != nil {
		return err
	}

	// Merge the custom entries into the collection
	bp.helper.MergeCustomSourceAnimeEntries(collection)

	// Save the raw collection to App（Bangumi 无 custom lists，raw 即全量副本）
	raw := &media.AnimeCollection{
		MediaListCollection: &media.AnimeCollection_MediaListCollection{},
	}
	listsCopy := make([]*media.AnimeCollection_MediaListCollection_Lists, len(collection.MediaListCollection.Lists))
	copy(listsCopy, collection.MediaListCollection.Lists)
	raw.MediaListCollection.Lists = listsCopy
	bp.rawAnimeCollection = mo.Some(raw)

	// 过滤无状态的 lists（对齐 AniList 行为；Bangumi 适配器产出的 lists 均有状态，此为防御性过滤）
	collection.MediaListCollection.Lists = bp.helper.FilterOutCustomAnimeLists(collection.MediaListCollection.Lists)

	// Save the collection to App
	bp.animeCollection = mo.Some(collection)

	return nil
}

// GetMangaCollection 获取漫画收藏（subject_type=1 全分页）。
// 漫画与轻小说同属 Bangumi 书籍分区，二者都进入 manga collection（轻小说域验收点）。
func (bp *BangumiPlatform) GetMangaCollection(ctx context.Context, bypassCache bool) (*media.MangaCollection, error) {
	if !bypassCache && bp.mangaCollection.IsPresent() {
		event := new(platform.GetCachedMangaCollectionEvent)
		event.MangaCollection = bp.mangaCollection.MustGet()
		err := hook.GlobalHookManager.OnGetCachedMangaCollection().Trigger(event)
		if err != nil {
			return nil, err
		}
		return event.MangaCollection, nil
	}

	if _, ok := bp.getUsername(ctx); !ok {
		return nil, nil
	}

	err := bp.refreshMangaCollection(ctx)
	if err != nil {
		return nil, err
	}

	event := new(platform.GetMangaCollectionEvent)
	event.MangaCollection = bp.mangaCollection.MustGet()

	err = hook.GlobalHookManager.OnGetMangaCollection().Trigger(event)
	if err != nil {
		return nil, err
	}

	return event.MangaCollection, nil
}

// GetRawMangaCollection 获取带 custom lists 的漫画收藏。
// Bangumi 无 custom lists，与 GetMangaCollection 等价（接口兼容保留）。
func (bp *BangumiPlatform) GetRawMangaCollection(ctx context.Context, bypassCache bool) (*media.MangaCollection, error) {
	bp.logger.Trace().Msg("bangumi platform: Fetching raw manga collection")

	if !bypassCache && bp.rawMangaCollection.IsPresent() {
		bp.logger.Trace().Msg("bangumi platform: Returning raw manga collection from cache")
		event := new(platform.GetCachedRawMangaCollectionEvent)
		event.MangaCollection = bp.rawMangaCollection.MustGet()
		err := hook.GlobalHookManager.OnGetCachedRawMangaCollection().Trigger(event)
		if err != nil {
			return nil, err
		}
		return event.MangaCollection, nil
	}

	if _, ok := bp.getUsername(ctx); !ok {
		return nil, nil
	}

	err := bp.refreshMangaCollection(ctx)
	if err != nil {
		return nil, err
	}

	event := new(platform.GetRawMangaCollectionEvent)
	event.MangaCollection = bp.rawMangaCollection.MustGet()

	err = hook.GlobalHookManager.OnGetRawMangaCollection().Trigger(event)
	if err != nil {
		return nil, err
	}

	return event.MangaCollection, nil
}

func (bp *BangumiPlatform) RefreshMangaCollection(ctx context.Context) (*media.MangaCollection, error) {
	if _, ok := bp.getUsername(ctx); !ok {
		return nil, nil
	}

	err := bp.refreshMangaCollection(ctx)
	if err != nil {
		return nil, err
	}

	event := new(platform.GetMangaCollectionEvent)
	event.MangaCollection = bp.mangaCollection.MustGet()

	err = hook.GlobalHookManager.OnGetMangaCollection().Trigger(event)
	if err != nil {
		return nil, err
	}

	event2 := new(platform.GetRawMangaCollectionEvent)
	event2.MangaCollection = bp.rawMangaCollection.MustGet()

	err = hook.GlobalHookManager.OnGetRawMangaCollection().Trigger(event2)
	if err != nil {
		return nil, err
	}

	return event.MangaCollection, nil
}

func (bp *BangumiPlatform) refreshMangaCollection(ctx context.Context) error {
	username, ok := bp.getUsername(ctx)
	if !ok {
		return errors.New("bangumi: Username is not set")
	}

	bp.collectionMu.Lock()
	defer bp.collectionMu.Unlock()

	collection, err := bp.cacheLayer.MangaCollection(ctx, username)
	if err != nil {
		return err
	}

	// Merge the custom entries into the collection
	bp.helper.MergeCustomSourceMangaEntries(collection)

	// Save the raw collection to App
	raw := &media.MangaCollection{
		MediaListCollection: &media.MangaCollection_MediaListCollection{},
	}
	listsCopy := make([]*media.MangaCollection_MediaListCollection_Lists, len(collection.MediaListCollection.Lists))
	copy(listsCopy, collection.MediaListCollection.Lists)
	raw.MediaListCollection.Lists = listsCopy
	bp.rawMangaCollection = mo.Some(raw)

	// 注意：与 AniList 版不同，这里不过滤轻小说（NOVEL）——
	// Bangumi 书籍分区不区分漫画/轻小说，轻小说域要求保留在 manga collection 中。

	collection.MediaListCollection.Lists = bp.helper.FilterOutCustomMangaLists(collection.MediaListCollection.Lists)

	// Save the collection to App
	bp.mangaCollection = mo.Some(collection)

	return nil
}

// GetAnimeCollectionWithRelations 获取带关系树的动画收藏（scanner 建树用）。
//
// Bangumi 无「一次查询全量带关系」的聚合端点，只能逐条目
// subject + related subjects 组树。节流：请求由 client 全局限流器串行化
// （0.5s/req），大集合耗时与条目数成正比；逐条处理并响应 ctx 取消，
// 完成数每 50 条打一条日志。
func (bp *BangumiPlatform) GetAnimeCollectionWithRelations(ctx context.Context) (*media.AnimeCollectionWithRelations, error) {
	bp.logger.Trace().Msg("bangumi platform: Fetching anime collection with relations")

	raw, err := bp.GetRawAnimeCollection(ctx, false)
	if err != nil {
		return nil, err
	}
	if raw == nil || raw.MediaListCollection == nil {
		return nil, nil
	}

	out := &media.AnimeCollectionWithRelations{
		MediaListCollection: &media.AnimeCollectionWithRelations_MediaListCollection{},
	}

	processed := 0
	for _, list := range raw.MediaListCollection.Lists {
		if list == nil {
			continue
		}
		newList := &media.AnimeCollectionWithRelations_MediaListCollection_Lists{
			Status:       list.Status,
			Name:         list.Name,
			IsCustomList: list.IsCustomList,
		}
		for _, entry := range list.Entries {
			if entry == nil || entry.Media == nil {
				continue
			}
			if err := ctx.Err(); err != nil {
				return nil, err
			}

			complete, err := bp.cacheLayer.CompleteAnimeByID(ctx, entry.Media.ID)
			if err != nil {
				// 单条失败不阻塞整树构建（降级为无关系条目）
				bp.logger.Warn().Err(err).Int("mediaId", entry.Media.ID).Msg("bangumi platform: Failed to fetch relations for collection entry")
				complete = completeAnimeFromAnime(entry.Media)
			}

			newList.Entries = append(newList.Entries, &media.AnimeCollectionWithRelations_MediaListCollection_Lists_Entries{
				CompletedAt: entry.CompletedAt,
				ID:          entry.ID,
				Media:       complete,
				Notes:       entry.Notes,
				Private:     entry.Private,
				Progress:    entry.Progress,
				Repeat:      entry.Repeat,
				Score:       entry.Score,
				StartedAt:   entry.StartedAt,
				Status:      entry.Status,
			})

			processed++
			if processed%50 == 0 {
				bp.logger.Info().Int("processed", processed).Msg("bangumi platform: Building collection relation tree")
			}
		}
		out.MediaListCollection.Lists = append(out.MediaListCollection.Lists, newList)
	}

	return out, nil
}

// completeAnimeFromAnime 由 Anime 构造无关系的 CompleteAnime（降级兜底）。
func completeAnimeFromAnime(a *media.Anime) *media.CompleteAnime {
	if a == nil {
		return nil
	}
	return &media.CompleteAnime{
		ID:              a.ID,
		IDMal:           a.IDMal,
		SiteURL:         a.SiteURL,
		Status:          a.Status,
		Season:          a.Season,
		SeasonYear:      a.SeasonYear,
		Type:            a.Type,
		Format:          a.Format,
		BannerImage:     a.BannerImage,
		Episodes:        a.Episodes,
		Synonyms:        a.Synonyms,
		IsAdult:         a.IsAdult,
		CountryOfOrigin: a.CountryOfOrigin,
		MeanScore:       a.MeanScore,
		Description:     a.Description,
		Genres:          a.Genres,
		Duration:        a.Duration,
		Trailer:         a.Trailer,
		Title:           a.Title,
		CoverImage:      a.CoverImage,
		StartDate:       a.StartDate,
		EndDate:         a.EndDate,
		NameCN:          a.NameCN,
	}
}

//----------------------------------------------------------------------------------------------------------------------
// 其余接口方法

// AddMediaToCollection 批量将条目加入「想看」。
// Bangumi 为逐条 upsert，请求由 client 节流器串行化。
func (bp *BangumiPlatform) AddMediaToCollection(ctx context.Context, mIds []int) error {
	bp.logger.Trace().Msg("bangumi platform: Adding media to collection")
	if len(mIds) == 0 {
		bp.logger.Debug().Msg("bangumi: No media added to planning list")
		return nil
	}

	for _, id := range mIds {
		if customsource.IsExtensionId(id) {
			_, err := bp.helper.HandleCustomSourceUpdateEntry(ctx,
				id,
				new(media.MediaListStatusPlanning),
				new(0),
				new(0),
				nil,
				nil,
			)
			if err != nil {
				bp.logger.Error().Msg("bangumi: An error occurred while adding media to planning list: " + err.Error())
			}
			continue
		}

		err := bp.cacheLayer.UpdateMediaListEntry(
			ctx,
			id,
			new(media.MediaListStatusPlanning),
			new(0),
			new(0),
			nil,
			nil,
		)
		if err != nil {
			bp.logger.Error().Msg("bangumi: An error occurred while adding media to planning list: " + err.Error())
		}
	}

	bp.logger.Debug().Any("count", len(mIds)).Msg("bangumi: Media added to planning list")
	return nil
}

// GetStudioDetails 获取制作公司详情。
//
// 降级实现：Bangumi 无制作公司 ID 体系。公司名在 GetAnimeDetails 时以
// 伪 ID 登记；本方法反查注册表构造最小结构，Media 列表留空。
// TODO: 如需「公司作品列表」，可基于收藏集合按 infobox 制作公司过滤（M4）。
func (bp *BangumiPlatform) GetStudioDetails(ctx context.Context, studioID int) (*media.StudioDetails, error) {
	bp.logger.Trace().Msg("bangumi platform: Fetching studio details")

	name := ""
	if v, ok := bp.studioNames.Load(studioID); ok {
		if s, ok := v.(string); ok {
			name = s
		}
	}

	ret := &media.StudioDetails{
		Studio: &media.StudioDetails_Studio{
			ID:                studioID,
			IsAnimationStudio: true,
			Name:              name,
			Media:             &media.StudioDetails_Studio_Media{}, // TODO: 无公司→作品反查端点，留空
		},
	}

	return bp.helper.TriggerGetStudioDetailsEvent(ret)
}

// GetViewerStats 获取观看统计。
//
// 降级实现：Bangumi 无服务端统计端点，由本地收藏集合合成最小统计。
// 结构齐全、数值保守（MinutesWatched 等无数据源的字段为 0）。
// TODO: 分钟数需单集时长数据，Bangumi 不提供，保持 0（M4 评估）。
func (bp *BangumiPlatform) GetViewerStats(ctx context.Context) (*media.ViewerStats, error) {
	if bp.username.IsAbsent() {
		return nil, errors.New("bangumi: Username is not set")
	}

	bp.logger.Trace().Msg("bangumi platform: Building viewer stats from local collection")

	stats := &media.ViewerStats{
		Viewer: &media.ViewerStats_Viewer{
			Statistics: &media.ViewerStats_Viewer_Statistics{},
		},
	}

	// 动画统计
	if animeCollection, err := bp.GetAnimeCollection(ctx, false); err == nil && animeCollection != nil {
		stats.Viewer.Statistics.Anime = buildAnimeStatsFromCollection(animeCollection)
	}

	// 书籍统计
	if mangaCollection, err := bp.GetMangaCollection(ctx, false); err == nil && mangaCollection != nil {
		stats.Viewer.Statistics.Manga = buildMangaStatsFromCollection(mangaCollection)
	}

	return stats, nil
}

// buildAnimeStatsFromCollection 由动画收藏合成最小统计。
func buildAnimeStatsFromCollection(collection *media.AnimeCollection) *media.ViewerStats_Viewer_Statistics_Anime {
	if collection == nil || collection.MediaListCollection == nil {
		return nil
	}

	anime := &media.ViewerStats_Viewer_Statistics_Anime{
		Statuses: make([]*media.UserStatusStats, 0),
		Scores:   make([]*media.UserScoreStats, 0),
	}

	statusCount := make(map[media.MediaListStatus]int)
	scoreSum, scoreCount := 0, 0

	for _, list := range collection.MediaListCollection.Lists {
		if list == nil {
			continue
		}
		for _, entry := range list.Entries {
			if entry == nil {
				continue
			}
			anime.Count++
			if entry.Progress != nil {
				anime.EpisodesWatched += *entry.Progress
			}
			if entry.Status != nil {
				statusCount[*entry.Status]++
			}
			if entry.Score != nil && *entry.Score > 0 {
				scoreSum += int(*entry.Score)
				scoreCount++
			}
		}
	}

	for _, status := range media.AllMediaListStatus {
		if c, ok := statusCount[status]; ok {
			anime.Statuses = append(anime.Statuses, &media.UserStatusStats{Status: new(status), Count: c})
		}
	}
	if scoreCount > 0 {
		anime.MeanScore = float64(scoreSum) / float64(scoreCount)
	}

	return anime
}

// buildMangaStatsFromCollection 由书籍收藏合成最小统计。
func buildMangaStatsFromCollection(collection *media.MangaCollection) *media.ViewerStats_Viewer_Statistics_Manga {
	if collection == nil || collection.MediaListCollection == nil {
		return nil
	}

	manga := &media.ViewerStats_Viewer_Statistics_Manga{
		Statuses: make([]*media.UserStatusStats, 0),
		Scores:   make([]*media.UserScoreStats, 0),
	}

	statusCount := make(map[media.MediaListStatus]int)
	scoreSum, scoreCount := 0, 0

	for _, list := range collection.MediaListCollection.Lists {
		if list == nil {
			continue
		}
		for _, entry := range list.Entries {
			if entry == nil {
				continue
			}
			manga.Count++
			if entry.Progress != nil {
				manga.ChaptersRead += *entry.Progress
			}
			if entry.Status != nil {
				statusCount[*entry.Status]++
			}
			if entry.Score != nil && *entry.Score > 0 {
				scoreSum += int(*entry.Score)
				scoreCount++
			}
		}
	}

	for _, status := range media.AllMediaListStatus {
		if c, ok := statusCount[status]; ok {
			manga.Statuses = append(manga.Statuses, &media.UserStatusStats{Status: new(status), Count: c})
		}
	}
	if scoreCount > 0 {
		manga.MeanScore = float64(scoreSum) / float64(scoreCount)
	}

	return manga
}

// GetAnimeAiringSchedule 获取放送时间表。
//
// 降级实现（M4 决策）：Bangumi 无逐集播出时间戳（仅日粒度 airdate），
// 无法填充 AniList 形状的逐集 schedule；返回空表，页面留待 M4 重构。
// TODO(M4): 可用 GET /v0/subjects/{id}/episodes 的 airdate 做日粒度近似。
func (bp *BangumiPlatform) GetAnimeAiringSchedule(ctx context.Context) (*media.AnimeAiringSchedule, error) {
	if bp.username.IsAbsent() {
		return nil, errors.New("bangumi: Username is not set")
	}

	bp.logger.Trace().Msg("bangumi platform: Airing schedule is not available on Bangumi, returning empty schedule")
	return &media.AnimeAiringSchedule{}, nil
}
