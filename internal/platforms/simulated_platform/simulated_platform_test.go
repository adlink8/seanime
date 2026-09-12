package simulated_platform

import (
	"context"
	"errors"
	"seanime/internal/api/bangumi"
	"seanime/internal/extension"
	"seanime/internal/local"
	"seanime/internal/media"
	"seanime/internal/platforms/shared_platform"
	"seanime/internal/util/limiter"
	"seanime/internal/testutil"
	"seanime/internal/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// refreshingFixtureClient 满足 shared_platform.BangumiAPI 的假客户端，
// 按 subject ID 返回预置条目，并记录刷新调用。
type refreshingFixtureClient struct {
	animeSubjects map[int]*bangumi.Subject
	bookSubjects  map[int]*bangumi.Subject
	animeCalls    []int
	mangaCalls    []int
}

func newRefreshingFixtureClient(animeByID map[int]*bangumi.Subject, mangaByID map[int]*bangumi.Subject) *refreshingFixtureClient {
	return &refreshingFixtureClient{
		animeSubjects: animeByID,
		bookSubjects:  mangaByID,
	}
}

func (c *refreshingFixtureClient) GetSubject(_ context.Context, subjectID int) (*bangumi.Subject, error) {
	if s, ok := c.animeSubjects[subjectID]; ok {
		c.animeCalls = append(c.animeCalls, subjectID)
		return s, nil
	}
	if s, ok := c.bookSubjects[subjectID]; ok {
		c.mangaCalls = append(c.mangaCalls, subjectID)
		return s, nil
	}

	return nil, errors.New("unexpected subject request")
}

func (c *refreshingFixtureClient) GetRelatedSubjects(_ context.Context, _ int) ([]bangumi.RelatedSubject, error) {
	return nil, nil
}

func (c *refreshingFixtureClient) GetUserCollectionsByUser(_ context.Context, _ string, _ bangumi.UserCollectionsOpts) (*bangumi.UserCollectionsResult, error) {
	return &bangumi.UserCollectionsResult{}, nil
}

func (c *refreshingFixtureClient) UpsertCollection(_ context.Context, _ int, _ bangumi.CollectionUpsertBody) error {
	return nil
}

func (c *refreshingFixtureClient) PatchCollection(_ context.Context, _ int, _ bangumi.CollectionUpsertBody) error {
	return nil
}

func (c *refreshingFixtureClient) DeleteCollection(_ context.Context, _ int) error {
	return nil
}

func TestRefreshAnimeCollectionRefreshesMutableEntries(t *testing.T) {
	sp, manager, client := newTestSimulatedPlatform(t, newRefreshingFixtureClient(
		map[int]*bangumi.Subject{
			// 放送日期在过去 → 推导为 RELEASING（可刷新）。
			// 注：Bangumi 纯 Subject 推导不出 FINISHED，故「已完结不刷新」的语义
			// 通过把已完结条目放进非 Current 列表来表达（见下方 collection）。
			101: newTestSubject(101, bangumi.SubjectAnime, "anime current fresh"),
			// 以下仅作哨兵：若实现误刷新这些条目，calls 断言会失败。
			102: newTestSubject(102, bangumi.SubjectAnime, "anime paused fresh"),
			103: newTestSubject(103, bangumi.SubjectAnime, "anime planning fresh"),
		},
		nil,
	))

	// keep a mix of mutable and settled entries so refresh only fetches the ones that can change
	manager.SaveSimulatedAnimeCollection(&media.AnimeCollection{
		MediaListCollection: &media.AnimeCollection_MediaListCollection{
			Lists: []*media.AnimeCollection_MediaListCollection_Lists{
				newAnimeCollectionList(media.MediaListStatusCurrent,
					newAnimeCollectionEntry(newTestSubject(101, bangumi.SubjectAnime, "anime current stale"), media.MediaListStatusCurrent),
				),
				newAnimeCollectionList(media.MediaListStatusPaused,
					newAnimeCollectionEntry(newTestSubject(102, bangumi.SubjectAnime, "anime paused stale"), media.MediaListStatusPaused),
				),
				newAnimeCollectionList(media.MediaListStatusPlanning,
					newAnimeCollectionEntry(newTestSubject(103, bangumi.SubjectAnime, "anime planning stale"), media.MediaListStatusPlanning),
				),
				newAnimeCollectionList(media.MediaListStatusCompleted,
					newAnimeCollectionEntry(newTestSubject(104, bangumi.SubjectAnime, "anime completed stale"), media.MediaListStatusCompleted),
					newAnimeCollectionEntry(newTestSubject(105, bangumi.SubjectAnime, "anime settled stale"), media.MediaListStatusCompleted),
				),
			},
		},
	})

	_, err := sp.RefreshAnimeCollection(context.Background())
	require.NoError(t, err)

	// 仅 Current + RELEASING 的 101 应被刷新（见 shouldRefreshSimulatedMedia）。
	require.ElementsMatch(t, []int{101}, client.animeCalls)

	collection := manager.GetSimulatedAnimeCollection().MustGet()
	currentEntry, found := collection.GetListEntryFromAnimeId(101)
	require.True(t, found)
	require.Equal(t, "anime current fresh", *currentEntry.GetMedia().GetTitle().GetEnglish())

	pausedEntry, found := collection.GetListEntryFromAnimeId(102)
	require.True(t, found)
	require.Equal(t, "anime paused stale", *pausedEntry.GetMedia().GetTitle().GetEnglish())

	planningEntry, found := collection.GetListEntryFromAnimeId(103)
	require.True(t, found)
	require.Equal(t, "anime planning stale", *planningEntry.GetMedia().GetTitle().GetEnglish())

	completedEntry, found := collection.GetListEntryFromAnimeId(104)
	require.True(t, found)
	require.Equal(t, "anime completed stale", *completedEntry.GetMedia().GetTitle().GetEnglish())

	settledEntry, found := collection.GetListEntryFromAnimeId(105)
	require.True(t, found)
	require.Equal(t, "anime settled stale", *settledEntry.GetMedia().GetTitle().GetEnglish())
}

func TestRefreshMangaCollectionRefreshesMutableEntries(t *testing.T) {
	sp, manager, client := newTestSimulatedPlatform(t, newRefreshingFixtureClient(
		nil,
		map[int]*bangumi.Subject{
			201: newTestSubject(201, bangumi.SubjectBook, "manga current fresh"),
			// 哨兵：不应被请求。
			202: newTestSubject(202, bangumi.SubjectBook, "manga paused fresh"),
			203: newTestSubject(203, bangumi.SubjectBook, "manga planning fresh"),
		},
	))

	// refresh should skip paused, planning, dropped or already settled manga entries.
	manager.SaveSimulatedMangaCollection(&media.MangaCollection{
		MediaListCollection: &media.MangaCollection_MediaListCollection{
			Lists: []*media.MangaCollection_MediaListCollection_Lists{
				newMangaCollectionList(media.MediaListStatusCurrent,
					newMangaCollectionEntry(newTestSubject(201, bangumi.SubjectBook, "manga current stale"), media.MediaListStatusCurrent),
				),
				newMangaCollectionList(media.MediaListStatusPaused,
					newMangaCollectionEntry(newTestSubject(202, bangumi.SubjectBook, "manga paused stale"), media.MediaListStatusPaused),
				),
				newMangaCollectionList(media.MediaListStatusPlanning,
					newMangaCollectionEntry(newTestSubject(203, bangumi.SubjectBook, "manga planning stale"), media.MediaListStatusPlanning),
				),
				newMangaCollectionList(media.MediaListStatusCompleted,
					newMangaCollectionEntry(newTestSubject(205, bangumi.SubjectBook, "manga settled stale"), media.MediaListStatusCompleted),
				),
				newMangaCollectionList(media.MediaListStatusDropped,
					newMangaCollectionEntry(newTestSubject(204, bangumi.SubjectBook, "manga dropped stale"), media.MediaListStatusDropped),
				),
			},
		},
	})

	_, err := sp.RefreshMangaCollection(context.Background())
	require.NoError(t, err)

	// 仅 Current + RELEASING 的 201 应被刷新（见 shouldRefreshSimulatedMedia）。
	require.ElementsMatch(t, []int{201}, client.mangaCalls)

	collection := manager.GetSimulatedMangaCollection().MustGet()
	currentEntry, found := collection.GetListEntryFromMangaId(201)
	require.True(t, found)
	require.Equal(t, "manga current fresh", *currentEntry.GetMedia().GetTitle().GetEnglish())

	pausedEntry, found := collection.GetListEntryFromMangaId(202)
	require.True(t, found)
	require.Equal(t, "manga paused stale", *pausedEntry.GetMedia().GetTitle().GetEnglish())

	planningEntry, found := collection.GetListEntryFromMangaId(203)
	require.True(t, found)
	require.Equal(t, "manga planning stale", *planningEntry.GetMedia().GetTitle().GetEnglish())

	droppedEntry, found := collection.GetListEntryFromMangaId(204)
	require.True(t, found)
	require.Equal(t, "manga dropped stale", *droppedEntry.GetMedia().GetTitle().GetEnglish())

	settledEntry, found := collection.GetListEntryFromMangaId(205)
	require.True(t, found)
	require.Equal(t, "manga settled stale", *settledEntry.GetMedia().GetTitle().GetEnglish())
}

func newTestSimulatedPlatform(t *testing.T, client *refreshingFixtureClient) (*SimulatedPlatform, local.Manager, *refreshingFixtureClient) {
	t.Helper()

	// Windows 下 manager 的 sqlite 文件连接未显式关闭会锁住 TempDir 导致清理失败；
	// local 包在 TEST_ENV=true 时使用内存数据库，不落盘。
	t.Setenv("TEST_ENV", "true")

	env := testutil.NewTestEnv(t)
	logger := env.Logger()
	database := env.MustNewDatabase(logger)
	// Windows 下 TempDir 清理会因 sqlite 连接未释放而失败，先显式关闭连接。
	t.Cleanup(func() {
		if sqlDB, err := database.Gorm().DB(); err == nil {
			_ = sqlDB.Close()
		}
	})
	manager := local.NewTestManager(t, database)

	shared_platform.ShouldCache.Store(true)
	shared_platform.IsWorking.Store(true)
	t.Cleanup(func() {
		shared_platform.ShouldCache.Store(true)
		shared_platform.IsWorking.Store(true)
	})

	sp := &SimulatedPlatform{
		logger:           logger,
		localManager:     manager,
		client:           client,
		cacheLayer:       shared_platform.NewCacheLayer(client, t.TempDir()),
		refreshRateLimit: limiter.NewLimiter(1*time.Millisecond, 1),
		helper:           shared_platform.NewPlatformHelper(util.NewRef(extension.NewUnifiedBank()), database, logger),
	}

	return sp, manager, client
}

// newTestSubject 构造测试用 Bangumi subject（放送日期取过去时间，状态推导为 RELEASING）。
func newTestSubject(id int, subjectType int, name string) *bangumi.Subject {
	return &bangumi.Subject{
		ID:   id,
		Type: subjectType,
		Name: name,
		Date: "2020-01-01",
	}
}

func newAnimeCollectionList(status media.MediaListStatus, entries ...*media.AnimeCollection_MediaListCollection_Lists_Entries) *media.AnimeCollection_MediaListCollection_Lists {
	return &media.AnimeCollection_MediaListCollection_Lists{
		Status:       &status,
		Name:         new(string(status)),
		IsCustomList: new(false),
		Entries:      entries,
	}
}

func newAnimeCollectionEntry(subject *bangumi.Subject, status media.MediaListStatus) *media.AnimeCollection_MediaListCollection_Lists_Entries {
	return &media.AnimeCollection_MediaListCollection_Lists_Entries{
		Media:    media.AnimeFromSubject(bangumi.SubjectToMedia(subject)),
		Progress: new(0),
		Score:    new(0.0),
		Repeat:   new(0),
		Status:   &status,
	}
}

func newMangaCollectionList(status media.MediaListStatus, entries ...*media.MangaCollection_MediaListCollection_Lists_Entries) *media.MangaCollection_MediaListCollection_Lists {
	return &media.MangaCollection_MediaListCollection_Lists{
		Status:       &status,
		Name:         new(string(status)),
		IsCustomList: new(false),
		Entries:      entries,
	}
}

func newMangaCollectionEntry(subject *bangumi.Subject, status media.MediaListStatus) *media.MangaCollection_MediaListCollection_Lists_Entries {
	return &media.MangaCollection_MediaListCollection_Lists_Entries{
		Media:    media.MangaFromSubject(bangumi.SubjectToMedia(subject)),
		Progress: new(0),
		Score:    new(0.0),
		Repeat:   new(0),
		Status:   &status,
	}
}
