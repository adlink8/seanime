package shared_platform

import (
	"context"
	"errors"
	"seanime/internal/api/bangumi"
	"seanime/internal/events"
	"seanime/internal/media"
	"seanime/internal/util"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// cacheLayerTestClient 满足 BangumiClient 接口的假客户端。
type cacheLayerTestClient struct {
	cacheDir           string
	animeCollectionRes *bangumi.UserCollectionsResult
	mangaCollectionRes *bangumi.UserCollectionsResult
	upsertCalls        []cacheLayerUpsertCall
	patchCalls         []cacheLayerPatchCall
	deleteCalls        []int
}

type cacheLayerUpsertCall struct {
	SubjectID int
	Body      bangumi.CollectionUpsertBody
}

type cacheLayerPatchCall struct {
	SubjectID int
	Body      bangumi.CollectionUpsertBody
}

func (c *cacheLayerTestClient) GetSubject(_ context.Context, subjectID int) (*bangumi.Subject, error) {
	return &bangumi.Subject{ID: subjectID}, nil
}

func (c *cacheLayerTestClient) GetRelatedSubjects(_ context.Context, _ int) ([]bangumi.RelatedSubject, error) {
	return nil, nil
}

func (c *cacheLayerTestClient) GetUserCollectionsByUser(_ context.Context, _ string, _ bangumi.UserCollectionsOpts) (*bangumi.UserCollectionsResult, error) {
	if c.animeCollectionRes != nil {
		return c.animeCollectionRes, nil
	}
	if c.mangaCollectionRes != nil {
		return c.mangaCollectionRes, nil
	}
	return &bangumi.UserCollectionsResult{}, nil
}

func (c *cacheLayerTestClient) UpsertCollection(_ context.Context, subjectID int, body bangumi.CollectionUpsertBody) error {
	c.upsertCalls = append(c.upsertCalls, cacheLayerUpsertCall{SubjectID: subjectID, Body: body})
	return nil
}

func (c *cacheLayerTestClient) PatchCollection(_ context.Context, subjectID int, body bangumi.CollectionUpsertBody) error {
	c.patchCalls = append(c.patchCalls, cacheLayerPatchCall{SubjectID: subjectID, Body: body})
	return nil
}

func (c *cacheLayerTestClient) DeleteCollection(_ context.Context, subjectID int) error {
	c.deleteCalls = append(c.deleteCalls, subjectID)
	return nil
}

func TestCacheLayerLogsOutOnInvalidToken(t *testing.T) {
	previousEventManager := events.GlobalWSEventManager
	events.GlobalWSEventManager = &events.GlobalWSEventManagerWrapper{}
	t.Cleanup(func() {
		events.GlobalWSEventManager = previousEventManager
		clearFailureTracking()
	})

	logoutCalled := make(chan struct{}, 1)
	cacheLayer := &CacheLayer{
		logoutFunc: func() {
			logoutCalled <- struct{}{}
		},
	}

	cacheLayer.checkAndUpdateWorkingState(errors.New("unauthorized: invalid token"))

	select {
	case <-logoutCalled:
	case <-time.After(time.Second):
		t.Fatal("expected invalid token error to trigger logout")
	}
	require.Zero(t, getRecentFailureCount())
}

func TestCacheLayerQueuesProgressUpdateAndPatchesAnimeCache(t *testing.T) {
	client := &cacheLayerTestClient{
		animeCollectionRes: newTestUserCollections(bangumi.SubjectAnime, 101, 321, bangumi.CollectionDoing, 2),
	}
	cacheLayer := newTestCacheLayer(t, client)

	// get the collection in cache
	_, err := cacheLayer.AnimeCollection(context.Background(), "-")
	require.NoError(t, err)

	IsWorking.Store(false)
	err = cacheLayer.UpdateMediaListEntryProgress(context.Background(), 101, 6, new(media.MediaListStatusCompleted))
	require.NoError(t, err)
	require.Empty(t, client.patchCalls)

	// the cached entry should move lists immediately so refetches see the local change.
	cached := getCachedAnimeCollection(t, cacheLayer)
	entry, found := cached.GetListEntryFromAnimeId(101)
	require.True(t, found)
	require.Equal(t, 6, *entry.GetProgress())
	require.Equal(t, media.MediaListStatusCompleted, *entry.GetStatus())
	require.True(t, animeListContains(cached, media.MediaListStatusCompleted, 101))
	require.False(t, animeListContains(cached, media.MediaListStatusCurrent, 101))

	queued := getQueuedUpdate(t, cacheLayer, 101)
	require.Equal(t, 101, queued.MediaID)
	require.Equal(t, 6, *queued.Progress)
	require.Equal(t, media.MediaListStatusCompleted, *queued.Status)
	require.False(t, queued.FullUpdate)
}

func TestCacheLayerQueuesEntryUpdateAndSyncsWhenOnline(t *testing.T) {
	client := &cacheLayerTestClient{
		mangaCollectionRes: newTestUserCollections(bangumi.SubjectBook, 202, 654, bangumi.CollectionDoing, 4),
	}
	cacheLayer := newTestCacheLayer(t, client)

	// seed manga cache, then queue an edit while the api is marked down
	_, err := cacheLayer.MangaCollection(context.Background(), "-")
	require.NoError(t, err)

	IsWorking.Store(false)
	startedAt := &media.FuzzyDateInput{Year: new(2025), Month: new(1), Day: new(2)}
	completedAt := &media.FuzzyDateInput{Year: new(2025), Month: new(2), Day: new(3)}
	err = cacheLayer.UpdateMediaListEntry(context.Background(), 202, new(media.MediaListStatusCompleted), new(85), new(12), startedAt, completedAt)
	require.NoError(t, err)
	require.Empty(t, client.upsertCalls)

	cached := getCachedMangaCollection(t, cacheLayer)
	entry, found := cached.GetListEntryFromMangaId(202)
	require.True(t, found)
	require.Equal(t, 12, *entry.GetProgress())
	require.Equal(t, float64(85), *entry.GetScore())
	require.Equal(t, media.MediaListStatusCompleted, *entry.GetStatus())
	require.True(t, mangaListContains(cached, media.MediaListStatusCompleted, 202))
	require.False(t, mangaListContains(cached, media.MediaListStatusCurrent, 202))

	// when the api is healthy again, the queued full edit is flushed once and removed
	IsWorking.Store(true)
	cacheLayer.syncQueuedUpdates(context.Background())
	require.Len(t, client.upsertCalls, 1)
	require.Equal(t, 202, client.upsertCalls[0].SubjectID)
	// 状态换算：COMPLETED → 看过(2)；评分换算：85 → 9（0-100 → 0-10）；进度 → ep_status
	require.Equal(t, bangumi.CollectionDone, *client.upsertCalls[0].Body.Type)
	require.Equal(t, 9, *client.upsertCalls[0].Body.Rate)
	require.Equal(t, 12, *client.upsertCalls[0].Body.EpStatus)
	requireNoQueuedUpdate(t, cacheLayer, 202)
}

func TestCacheLayerLiveProgressUpdateClearsQueuedUpdate(t *testing.T) {
	client := &cacheLayerTestClient{
		animeCollectionRes: newTestUserCollections(bangumi.SubjectAnime, 101, 321, bangumi.CollectionDoing, 2),
	}
	cacheLayer := newTestCacheLayer(t, client)

	_, err := cacheLayer.AnimeCollection(context.Background(), "-")
	require.NoError(t, err)

	// first update is queued while the api is marked down
	IsWorking.Store(false)
	err = cacheLayer.UpdateMediaListEntryProgress(context.Background(), 101, 6, new(media.MediaListStatusCompleted))
	require.NoError(t, err)
	queued := getQueuedUpdate(t, cacheLayer, 101)
	require.Equal(t, 6, *queued.Progress)

	// a later online update should win and remove the stale queued state
	IsWorking.Store(true)
	err = cacheLayer.UpdateMediaListEntryProgress(context.Background(), 101, 7, new(media.MediaListStatusCurrent))
	require.NoError(t, err)
	requireNoQueuedUpdate(t, cacheLayer, 101)

	cacheLayer.syncQueuedUpdates(context.Background())
	require.Len(t, client.patchCalls, 1)
	require.Equal(t, 7, *client.patchCalls[0].Body.EpStatus)
}

func TestCacheLayerLiveEntryUpdateClearsQueuedUpdate(t *testing.T) {
	client := &cacheLayerTestClient{
		mangaCollectionRes: newTestUserCollections(bangumi.SubjectBook, 202, 654, bangumi.CollectionDoing, 4),
	}
	cacheLayer := newTestCacheLayer(t, client)

	_, err := cacheLayer.MangaCollection(context.Background(), "-")
	require.NoError(t, err)

	// queue an older edit while Bangumi is unavailable
	IsWorking.Store(false)
	err = cacheLayer.UpdateMediaListEntry(context.Background(), 202, new(media.MediaListStatusCompleted), new(80), new(12), nil, nil)
	require.NoError(t, err)
	queued := getQueuedUpdate(t, cacheLayer, 202)
	require.Equal(t, 80, *queued.ScoreRaw)

	// the successful online edit replaces it and should prevent stale replay
	IsWorking.Store(true)
	err = cacheLayer.UpdateMediaListEntry(context.Background(), 202, new(media.MediaListStatusCurrent), new(90), new(13), nil, nil)
	require.NoError(t, err)
	requireNoQueuedUpdate(t, cacheLayer, 202)

	cacheLayer.syncQueuedUpdates(context.Background())
	require.Len(t, client.upsertCalls, 1)
	require.Equal(t, 90, *client.upsertCalls[0].Body.Rate*10)
	require.Equal(t, 13, *client.upsertCalls[0].Body.EpStatus)
}

func newTestCacheLayer(t *testing.T, client *cacheLayerTestClient) *CacheLayer {
	t.Helper()
	ShouldCache.Store(true)
	IsWorking.Store(true)
	clearFailureTracking()

	BangumiClient.Store(util.NewRef[BangumiAPI](client))
	cacheLayer := NewCacheLayer(client, t.TempDir())

	t.Cleanup(func() {
		ShouldCache.Store(true)
		IsWorking.Store(true)
		clearFailureTracking()
	})

	return cacheLayer
}

// newTestUserCollections 构造单条收藏条目的用户收藏响应。
func newTestUserCollections(subjectType int, mediaID int, subjectID int, collectionType int, epStatus int) *bangumi.UserCollectionsResult {
	return &bangumi.UserCollectionsResult{
		Total: 1,
		Data: []bangumi.UserCollection{
			{
				SubjectID: subjectID,
				Type:      collectionType,
				EpStatus:  epStatus,
				Subject: &bangumi.Subject{
					ID:   mediaID,
					Type: subjectType,
					Name: "Test Subject",
				},
			},
		},
	}
}

func getCachedAnimeCollection(t *testing.T, cacheLayer *CacheLayer) *media.AnimeCollection {
	t.Helper()
	var cached media.AnimeCollection
	found, err := cacheLayer.fileCacher.GetPerm(cacheLayer.buckets[AnimeCollectionBucket], cacheLayer.generateCacheKey("collection", "-"), &cached)
	require.NoError(t, err)
	require.True(t, found)
	return &cached
}

func getCachedMangaCollection(t *testing.T, cacheLayer *CacheLayer) *media.MangaCollection {
	t.Helper()
	var cached media.MangaCollection
	found, err := cacheLayer.fileCacher.GetPerm(cacheLayer.buckets[MangaCollectionBucket], cacheLayer.generateCacheKey("collection", "-"), &cached)
	require.NoError(t, err)
	require.True(t, found)
	return &cached
}

func getQueuedUpdate(t *testing.T, cacheLayer *CacheLayer, mediaID int) queuedMediaListUpdate {
	t.Helper()
	var queued queuedMediaListUpdate
	found, err := cacheLayer.fileCacher.GetPerm(cacheLayer.buckets[PendingMediaListUpdatesBucket], strconv.Itoa(mediaID), &queued)
	require.NoError(t, err)
	require.True(t, found)
	return queued
}

func requireNoQueuedUpdate(t *testing.T, cacheLayer *CacheLayer, mediaID int) {
	t.Helper()
	var queued queuedMediaListUpdate
	found, err := cacheLayer.fileCacher.GetPerm(cacheLayer.buckets[PendingMediaListUpdatesBucket], strconv.Itoa(mediaID), &queued)
	require.NoError(t, err)
	require.False(t, found)
}

func animeListContains(collection *media.AnimeCollection, status media.MediaListStatus, mediaID int) bool {
	for _, list := range collection.GetMediaListCollection().GetLists() {
		if list.GetStatus() == nil || *list.GetStatus() != status {
			continue
		}
		for _, entry := range list.GetEntries() {
			if entry.GetMedia().GetID() == mediaID {
				return true
			}
		}
	}
	return false
}

func mangaListContains(collection *media.MangaCollection, status media.MediaListStatus, mediaID int) bool {
	for _, list := range collection.GetMediaListCollection().GetLists() {
		if list.GetStatus() == nil || *list.GetStatus() != status {
			continue
		}
		for _, entry := range list.GetEntries() {
			if entry.GetMedia().GetID() == mediaID {
				return true
			}
		}
	}
	return false
}
