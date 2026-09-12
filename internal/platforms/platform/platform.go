package platform

import (
	"context"
	"seanime/internal/api/bangumi"
	"seanime/internal/media"
)

// Platform 平台层接口（Phase 2 换锚：数据锚点由 AniList 迁至 Bangumi）。
//
// 类型约定：所有领域类型均来自 internal/media 镜像层；
// mediaID 一律为 Bangumi subject ID。
type Platform interface {
	SetUsername(username string)
	// UpdateEntry updates the entry for the given media ID
	UpdateEntry(context context.Context, mediaID int, status *media.MediaListStatus, scoreRaw *int, progress *int, startedAt *media.FuzzyDateInput, completedAt *media.FuzzyDateInput) error
	// UpdateEntryProgress updates the entry progress for the given media ID
	UpdateEntryProgress(context context.Context, mediaID int, progress int, totalEpisodes *int) error
	// UpdateEntryRepeat updates the entry repeat number for the given media ID
	UpdateEntryRepeat(context context.Context, mediaID int, repeat int) error
	// DeleteEntry deletes the entry for the given media ID
	DeleteEntry(context context.Context, mediaID int, entryID int) error
	// GetAnime gets the anime for the given media ID
	GetAnime(context context.Context, mediaID int) (*media.Anime, error)
	// GetAnimeByMalID gets the anime by MAL ID
	GetAnimeByMalID(context context.Context, malID int) (*media.Anime, error)
	// GetAnimeWithRelations gets the anime with relations for the given media ID
	// This is used for scanning purposes in order to build the relation tree
	GetAnimeWithRelations(context context.Context, mediaID int) (*media.CompleteAnime, error)
	// GetAnimeDetails gets the anime details for the given media ID
	// These details are only fetched by the anime page
	GetAnimeDetails(context context.Context, mediaID int) (*media.AnimeDetails, error)
	// GetManga gets the manga for the given media ID
	GetManga(context context.Context, mediaID int) (*media.Manga, error)
	// GetAnimeCollection gets the anime collection without custom lists
	// This should not make any API calls and instead should be based on GetRawAnimeCollection
	GetAnimeCollection(context context.Context, bypassCache bool) (*media.AnimeCollection, error)
	// GetRawAnimeCollection gets the anime collection with custom lists
	GetRawAnimeCollection(context context.Context, bypassCache bool) (*media.AnimeCollection, error)
	// GetMangaDetails gets the manga details for the given media ID
	// These details are only fetched by the manga page
	GetMangaDetails(context context.Context, mediaID int) (*media.MangaDetails, error)
	// GetAnimeCollectionWithRelations gets the anime collection with relations
	// This is used for scanning purposes in order to build the relation tree
	GetAnimeCollectionWithRelations(context context.Context) (*media.AnimeCollectionWithRelations, error)
	// GetMangaCollection gets the manga collection without custom lists
	// This should not make any API calls and instead should be based on GetRawMangaCollection
	GetMangaCollection(context context.Context, bypassCache bool) (*media.MangaCollection, error)
	// GetRawMangaCollection gets the manga collection with custom lists
	GetRawMangaCollection(context context.Context, bypassCache bool) (*media.MangaCollection, error)
	// AddMediaToCollection adds the media to the collection
	AddMediaToCollection(context context.Context, mIds []int) error
	// GetStudioDetails gets the studio details for the given studio ID
	GetStudioDetails(context context.Context, studioID int) (*media.StudioDetails, error)
	// GetBangumiClient gets the Bangumi API client
	// （替代原 GetAnilistClient；Wave C 消费方按需改写）
	GetBangumiClient() *bangumi.Client
	// RefreshAnimeCollection refreshes the anime collection
	RefreshAnimeCollection(context context.Context) (*media.AnimeCollection, error)
	// RefreshMangaCollection refreshes the manga collection
	RefreshMangaCollection(context context.Context) (*media.MangaCollection, error)
	// GetViewerStats gets the viewer stats
	GetViewerStats(context context.Context) (*media.ViewerStats, error)
	// GetAnimeAiringSchedule gets the schedule for airing anime in the collection
	GetAnimeAiringSchedule(context context.Context) (*media.AnimeAiringSchedule, error)
	ClearCache()
	Close()
}
