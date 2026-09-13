package testmocks

import (
	"context"
	"fmt"
	"seanime/internal/api/bangumi"
	"seanime/internal/media"
)

type FakePlatformBuilder struct {
	platform *FakePlatform
}

type FakePlatform struct {
	animeByID                   map[int]*media.Anime
	mangaByID                   map[int]*media.Manga
	animeCollection             *media.AnimeCollection
	rawAnimeCollection          *media.AnimeCollection
	animeCollectionWithRel      *media.AnimeCollectionWithRelations
	mangaCollection             *media.MangaCollection
	rawMangaCollection          *media.MangaCollection
	animeAiringSchedule         *media.AnimeAiringSchedule
	viewerStats                 *media.ViewerStats
	animeCollectionErr          error
	rawAnimeCollectionErr       error
	animeCollectionWithRelErr   error
	mangaCollectionErr          error
	rawMangaCollectionErr       error
	animeAiringScheduleErr      error
	viewerStatsErr              error
	updateEntryProgressErr      error
	animeCalls                  map[int]int
	mangaCalls                  map[int]int
	animeCollectionCalls        int
	rawAnimeCollectionCalls     int
	animeCollectionWithRelCalls int
	mangaCollectionCalls        int
	rawMangaCollectionCalls     int
	animeAiringScheduleCalls    int
	viewerStatsCalls            int
	updateEntryCalls            []FakeUpdateEntryCall
	updateEntryProgressCalls    []FakeUpdateEntryProgressCall
	bangumiClient               *bangumi.Client
}

type FakeUpdateEntryCall struct {
	MediaID     int
	Status      *media.MediaListStatus
	ScoreRaw    *int
	Progress    *int
	StartedAt   *media.FuzzyDateInput
	CompletedAt *media.FuzzyDateInput
}

type FakeUpdateEntryProgressCall struct {
	MediaID       int
	Progress      int
	TotalEpisodes *int
}

func NewFakePlatformBuilder() *FakePlatformBuilder {
	return &FakePlatformBuilder{
		platform: &FakePlatform{
			animeByID:  make(map[int]*media.Anime),
			mangaByID:  make(map[int]*media.Manga),
			animeCalls: make(map[int]int),
			mangaCalls: make(map[int]int),
		},
	}
}

func (b *FakePlatformBuilder) WithAnime(anime *media.Anime) *FakePlatformBuilder {
	if anime != nil {
		b.platform.animeByID[anime.ID] = anime
	}
	return b
}

func (b *FakePlatformBuilder) WithManga(manga *media.Manga) *FakePlatformBuilder {
	if manga != nil {
		b.platform.mangaByID[manga.ID] = manga
	}
	return b
}

func (b *FakePlatformBuilder) WithAnimeCollection(collection *media.AnimeCollection) *FakePlatformBuilder {
	b.platform.animeCollection = collection
	return b
}

func (b *FakePlatformBuilder) WithAnimeCollectionError(err error) *FakePlatformBuilder {
	b.platform.animeCollectionErr = err
	return b
}

func (b *FakePlatformBuilder) WithRawAnimeCollection(collection *media.AnimeCollection) *FakePlatformBuilder {
	b.platform.rawAnimeCollection = collection
	return b
}

func (b *FakePlatformBuilder) WithAnimeCollectionWithRelations(collection *media.AnimeCollectionWithRelations) *FakePlatformBuilder {
	b.platform.animeCollectionWithRel = collection
	return b
}

func (b *FakePlatformBuilder) WithMangaCollection(collection *media.MangaCollection) *FakePlatformBuilder {
	b.platform.mangaCollection = collection
	return b
}

func (b *FakePlatformBuilder) WithAnimeAiringSchedule(schedule *media.AnimeAiringSchedule) *FakePlatformBuilder {
	b.platform.animeAiringSchedule = schedule
	return b
}

func (b *FakePlatformBuilder) WithViewerStats(stats *media.ViewerStats) *FakePlatformBuilder {
	b.platform.viewerStats = stats
	return b
}

func (b *FakePlatformBuilder) WithUpdateEntryProgressError(err error) *FakePlatformBuilder {
	b.platform.updateEntryProgressErr = err
	return b
}

// WithBangumiClient 注入 Bangumi client，使 handler 测试可经 httptest 观察到
// 真实的出站通路（legacy vs v0）。不注入时 GetBangumiClient() 返回 nil。
func (b *FakePlatformBuilder) WithBangumiClient(client *bangumi.Client) *FakePlatformBuilder {
	b.platform.bangumiClient = client
	return b
}

func (b *FakePlatformBuilder) Build() *FakePlatform {
	return b.platform
}

func (f *FakePlatform) AnimeCalls(mediaID int) int {
	return f.animeCalls[mediaID]
}

func (f *FakePlatform) MangaCalls(mediaID int) int {
	return f.mangaCalls[mediaID]
}

func (f *FakePlatform) AnimeCollectionCalls() int {
	return f.animeCollectionCalls
}

func (f *FakePlatform) UpdateEntryProgressCalls() []FakeUpdateEntryProgressCall {
	ret := make([]FakeUpdateEntryProgressCall, len(f.updateEntryProgressCalls))
	copy(ret, f.updateEntryProgressCalls)
	return ret
}

func (f *FakePlatform) SetUsername(string) {}

func (f *FakePlatform) UpdateEntryCalls() []FakeUpdateEntryCall {
	ret := make([]FakeUpdateEntryCall, len(f.updateEntryCalls))
	copy(ret, f.updateEntryCalls)
	return ret
}

func (f *FakePlatform) UpdateEntry(_ context.Context, mediaID int, status *media.MediaListStatus, scoreRaw *int, progress *int, startedAt *media.FuzzyDateInput, completedAt *media.FuzzyDateInput) error {
	call := FakeUpdateEntryCall{MediaID: mediaID}
	if status != nil {
		statusCopy := *status
		call.Status = &statusCopy
	}
	if scoreRaw != nil {
		scoreCopy := *scoreRaw
		call.ScoreRaw = &scoreCopy
	}
	if progress != nil {
		progressCopy := *progress
		call.Progress = &progressCopy
	}
	if startedAt != nil {
		startedAtCopy := *startedAt
		call.StartedAt = &startedAtCopy
	}
	if completedAt != nil {
		completedAtCopy := *completedAt
		call.CompletedAt = &completedAtCopy
	}
	f.updateEntryCalls = append(f.updateEntryCalls, call)
	return nil
}

func (f *FakePlatform) UpdateEntryProgress(_ context.Context, mediaID int, progress int, totalEpisodes *int) error {
	call := FakeUpdateEntryProgressCall{}
	call.MediaID = mediaID
	call.Progress = progress
	if totalEpisodes != nil {
		totalEpisodesCopy := *totalEpisodes
		call.TotalEpisodes = &totalEpisodesCopy
	}
	f.updateEntryProgressCalls = append(f.updateEntryProgressCalls, call)
	return f.updateEntryProgressErr
}

func (f *FakePlatform) UpdateEntryRepeat(context.Context, int, int) error {
	return nil
}

func (f *FakePlatform) DeleteEntry(context.Context, int, int) error {
	return nil
}

func (f *FakePlatform) GetAnime(_ context.Context, mediaID int) (*media.Anime, error) {
	f.animeCalls[mediaID]++
	anime, ok := f.animeByID[mediaID]
	if !ok {
		return nil, fmt.Errorf("anime %d not found", mediaID)
	}
	return anime, nil
}

func (f *FakePlatform) GetAnimeByMalID(context.Context, int) (*media.Anime, error) {
	return nil, nil
}

func (f *FakePlatform) GetAnimeWithRelations(context.Context, int) (*media.CompleteAnime, error) {
	return nil, nil
}

func (f *FakePlatform) GetAnimeDetails(context.Context, int) (*media.AnimeDetails, error) {
	return nil, nil
}

func (f *FakePlatform) GetManga(_ context.Context, mediaID int) (*media.Manga, error) {
	f.mangaCalls[mediaID]++
	manga, ok := f.mangaByID[mediaID]
	if !ok {
		return nil, fmt.Errorf("manga %d not found", mediaID)
	}
	return manga, nil
}

func (f *FakePlatform) GetAnimeCollection(context.Context, bool) (*media.AnimeCollection, error) {
	f.animeCollectionCalls++
	if f.animeCollectionErr != nil {
		return nil, f.animeCollectionErr
	}
	if f.animeCollection == nil {
		f.animeCollection = &media.AnimeCollection{}
	}
	return f.animeCollection, nil
}

func (f *FakePlatform) GetRawAnimeCollection(context.Context, bool) (*media.AnimeCollection, error) {
	f.rawAnimeCollectionCalls++
	if f.rawAnimeCollectionErr != nil {
		return nil, f.rawAnimeCollectionErr
	}
	return f.rawAnimeCollection, nil
}

func (f *FakePlatform) GetMangaDetails(context.Context, int) (*media.MangaDetails, error) {
	return nil, nil
}

func (f *FakePlatform) GetAnimeCollectionWithRelations(context.Context) (*media.AnimeCollectionWithRelations, error) {
	f.animeCollectionWithRelCalls++
	if f.animeCollectionWithRelErr != nil {
		return nil, f.animeCollectionWithRelErr
	}
	return f.animeCollectionWithRel, nil
}

func (f *FakePlatform) GetMangaCollection(context.Context, bool) (*media.MangaCollection, error) {
	f.mangaCollectionCalls++
	if f.mangaCollectionErr != nil {
		return nil, f.mangaCollectionErr
	}
	return f.mangaCollection, nil
}

func (f *FakePlatform) GetRawMangaCollection(context.Context, bool) (*media.MangaCollection, error) {
	f.rawMangaCollectionCalls++
	if f.rawMangaCollectionErr != nil {
		return nil, f.rawMangaCollectionErr
	}
	return f.rawMangaCollection, nil
}

func (f *FakePlatform) AddMediaToCollection(context.Context, []int) error {
	return nil
}

func (f *FakePlatform) GetStudioDetails(context.Context, int) (*media.StudioDetails, error) {
	return nil, nil
}

// GetBangumiClient 返回 Bangumi client（FakePlatform 默认无 client，测试按需注入）。
func (f *FakePlatform) GetBangumiClient() *bangumi.Client {
	return f.bangumiClient
}

func (f *FakePlatform) RefreshAnimeCollection(context.Context) (*media.AnimeCollection, error) {
	return nil, nil
}

func (f *FakePlatform) RefreshMangaCollection(context.Context) (*media.MangaCollection, error) {
	return nil, nil
}

func (f *FakePlatform) GetViewerStats(context.Context) (*media.ViewerStats, error) {
	f.viewerStatsCalls++
	if f.viewerStatsErr != nil {
		return nil, f.viewerStatsErr
	}
	return f.viewerStats, nil
}

func (f *FakePlatform) GetAnimeAiringSchedule(context.Context) (*media.AnimeAiringSchedule, error) {
	f.animeAiringScheduleCalls++
	if f.animeAiringScheduleErr != nil {
		return nil, f.animeAiringScheduleErr
	}
	return f.animeAiringSchedule, nil
}

func (f *FakePlatform) ClearCache() {}

func (f *FakePlatform) Close() {}
