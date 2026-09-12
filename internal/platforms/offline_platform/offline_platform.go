package offline_platform

import (
	"context"
	"errors"
	"seanime/internal/api/bangumi"
	"seanime/internal/local"
	"seanime/internal/media"
	"seanime/internal/platforms/platform"

	"github.com/rs/zerolog"
)

var (
	ErrNoLocalAnimeCollection   = errors.New("no local anime collection")
	ErrorNoLocalMangaCollection = errors.New("no local manga collection")
	// ErrMediaNotFound means the media wasn't found in the local collection
	ErrMediaNotFound = errors.New("media not found")
	// ErrActionNotSupported means the action isn't valid on the local platform
	ErrActionNotSupported = errors.New("action not supported")
)

// OfflinePlatform used when offline.
// It provides the same API as the bangumi_platform.BangumiPlatform but some methods are no-op.
type OfflinePlatform struct {
	logger       *zerolog.Logger
	localManager local.Manager
}

func NewOfflinePlatform(localManager local.Manager, logger *zerolog.Logger) (platform.Platform, error) {
	ap := &OfflinePlatform{
		logger:       logger,
		localManager: localManager,
	}

	return ap, nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (lp *OfflinePlatform) SetUsername(username string) {
	// no-op
}

func (lp *OfflinePlatform) Close() {
	// no-op
}

func (lp *OfflinePlatform) ClearCache() {
	// no-op
}

func rearrangeAnimeCollectionLists(animeCollection *media.AnimeCollection) {
	removedEntries := make([]*media.AnimeCollection_MediaListCollection_Lists_Entries, 0)
	for _, list := range animeCollection.MediaListCollection.Lists {
		if list.GetStatus() == nil || list.GetEntries() == nil {
			continue
		}
		var indicesToRemove []int
		for idx, entry := range list.GetEntries() {
			if entry.GetStatus() == nil {
				continue
			}
			// Mark for removal if status differs
			if *list.GetStatus() != *entry.GetStatus() {
				indicesToRemove = append(indicesToRemove, idx)
				removedEntries = append(removedEntries, entry)
			}
		}
		// Remove entries in reverse order to avoid re-slicing issues
		for i := len(indicesToRemove) - 1; i >= 0; i-- {
			idx := indicesToRemove[i]
			list.Entries = append(list.Entries[:idx], list.Entries[idx+1:]...)
		}
	}

	// Add removed entries to the correct list
	for _, entry := range removedEntries {
		for _, list := range animeCollection.MediaListCollection.Lists {
			if list.GetStatus() == nil {
				continue
			}
			if *list.GetStatus() == *entry.GetStatus() {
				list.Entries = append(list.Entries, entry)
			}
		}
	}
}

func rearrangeMangaCollectionLists(mangaCollection *media.MangaCollection) {
	removedEntries := make([]*media.MangaCollection_MediaListCollection_Lists_Entries, 0)
	for _, list := range mangaCollection.MediaListCollection.Lists {
		if list.GetStatus() == nil || list.GetEntries() == nil {
			continue
		}
		var indicesToRemove []int
		for idx, entry := range list.GetEntries() {
			if entry.GetStatus() == nil {
				continue
			}
			// Mark for removal if status differs
			if *list.GetStatus() != *entry.GetStatus() {
				indicesToRemove = append(indicesToRemove, idx)
				removedEntries = append(removedEntries, entry)
			}
		}
		// Remove entries in reverse order to avoid re-slicing issues
		for i := len(indicesToRemove) - 1; i >= 0; i-- {
			idx := indicesToRemove[i]
			list.Entries = append(list.Entries[:idx], list.Entries[idx+1:]...)
		}
	}

	// Add removed entries to the correct list
	for _, entry := range removedEntries {
		for _, list := range mangaCollection.MediaListCollection.Lists {
			if list.GetStatus() == nil {
				continue
			}
			if *list.GetStatus() == *entry.GetStatus() {
				list.Entries = append(list.Entries, entry)
			}
		}
	}
}

// UpdateEntry updates the entry for the given media ID.
// It doesn't add the entry if it doesn't exist.
func (lp *OfflinePlatform) UpdateEntry(ctx context.Context, mediaID int, status *media.MediaListStatus, scoreRaw *int, progress *int, startedAt *media.FuzzyDateInput, completedAt *media.FuzzyDateInput) error {
	if lp.localManager.GetLocalAnimeCollection().IsPresent() {
		animeCollection := lp.localManager.GetLocalAnimeCollection().MustGet()

		// Find the entry
		for _, list := range animeCollection.MediaListCollection.Lists {
			for _, entry := range list.GetEntries() {
				if entry.GetMedia().GetID() == mediaID {
					// Update the entry
					if status != nil {
						entry.Status = status
					}
					if scoreRaw != nil {
						entry.Score = new(float64(*scoreRaw))
					}
					if progress != nil {
						entry.Progress = progress
					}
					if startedAt != nil {
						entry.StartedAt = &media.AnimeCollection_MediaListCollection_Lists_Entries_StartedAt{
							Year:  startedAt.Year,
							Month: startedAt.Month,
							Day:   startedAt.Day,
						}
					}
					if completedAt != nil {
						entry.CompletedAt = &media.AnimeCollection_MediaListCollection_Lists_Entries_CompletedAt{
							Year:  completedAt.Year,
							Month: completedAt.Month,
							Day:   completedAt.Day,
						}
					}

					// Save the collection
					rearrangeAnimeCollectionLists(animeCollection)
					lp.localManager.UpdateLocalAnimeCollection(animeCollection)
					lp.localManager.SetHasLocalChanges(true)
					return nil
				}
			}
		}
	}

	if lp.localManager.GetLocalMangaCollection().IsPresent() {
		mangaCollection := lp.localManager.GetLocalMangaCollection().MustGet()

		// Find the entry
		for _, list := range mangaCollection.MediaListCollection.Lists {
			for _, entry := range list.Entries {
				if entry.GetMedia().GetID() == mediaID {
					// Update the entry
					if status != nil {
						entry.Status = status
					}
					if scoreRaw != nil {
						entry.Score = new(float64(*scoreRaw))
					}
					if progress != nil {
						entry.Progress = progress
					}
					if startedAt != nil {
						entry.StartedAt = &media.MangaCollection_MediaListCollection_Lists_Entries_StartedAt{
							Year:  startedAt.Year,
							Month: startedAt.Month,
							Day:   startedAt.Day,
						}
					}
					if completedAt != nil {
						entry.CompletedAt = &media.MangaCollection_MediaListCollection_Lists_Entries_CompletedAt{
							Year:  completedAt.Year,
							Month: completedAt.Month,
							Day:   completedAt.Day,
						}
					}

					// Save the collection
					rearrangeMangaCollectionLists(mangaCollection)
					lp.localManager.UpdateLocalMangaCollection(mangaCollection)
					lp.localManager.SetHasLocalChanges(true)
					return nil
				}
			}
		}
	}

	return ErrMediaNotFound
}

func (lp *OfflinePlatform) UpdateEntryProgress(ctx context.Context, mediaID int, progress int, totalEpisodes *int) error {
	if lp.localManager.GetLocalAnimeCollection().IsPresent() {
		animeCollection := lp.localManager.GetLocalAnimeCollection().MustGet()

		// Find the entry
		for _, list := range animeCollection.MediaListCollection.Lists {
			for _, entry := range list.GetEntries() {
				if entry.GetMedia().GetID() == mediaID {
					// Update the entry
					entry.Progress = &progress
					if totalEpisodes != nil {
						entry.Media.Episodes = totalEpisodes
					}

					// Save the collection
					rearrangeAnimeCollectionLists(animeCollection)
					lp.localManager.UpdateLocalAnimeCollection(animeCollection)
					lp.localManager.SetHasLocalChanges(true)
					return nil
				}
			}
		}
	}

	if lp.localManager.GetLocalMangaCollection().IsPresent() {
		mangaCollection := lp.localManager.GetLocalMangaCollection().MustGet()

		// Find the entry
		for _, list := range mangaCollection.MediaListCollection.Lists {
			for _, entry := range list.Entries {
				if entry.GetMedia().GetID() == mediaID {
					// Update the entry
					entry.Progress = &progress
					if totalEpisodes != nil {
						entry.Media.Chapters = totalEpisodes
					}

					// Save the collection
					rearrangeMangaCollectionLists(mangaCollection)
					lp.localManager.UpdateLocalMangaCollection(mangaCollection)
					lp.localManager.SetHasLocalChanges(true)
					return nil
				}
			}
		}
	}

	return ErrMediaNotFound
}

func (lp *OfflinePlatform) UpdateEntryRepeat(ctx context.Context, mediaID int, repeat int) error {
	if lp.localManager.GetLocalAnimeCollection().IsPresent() {
		animeCollection := lp.localManager.GetLocalAnimeCollection().MustGet()

		// Find the entry
		for _, list := range animeCollection.MediaListCollection.Lists {
			for _, entry := range list.GetEntries() {
				if entry.GetMedia().GetID() == mediaID {
					// Update the entry
					entry.Repeat = &repeat

					// Save the collection
					rearrangeAnimeCollectionLists(animeCollection)
					lp.localManager.UpdateLocalAnimeCollection(animeCollection)
					lp.localManager.SetHasLocalChanges(true)
					return nil
				}
			}
		}
	}

	if lp.localManager.GetLocalMangaCollection().IsPresent() {
		mangaCollection := lp.localManager.GetLocalMangaCollection().MustGet()

		// Find the entry
		for _, list := range mangaCollection.MediaListCollection.Lists {
			for _, entry := range list.Entries {
				if entry.GetMedia().GetID() == mediaID {
					// Update the entry
					entry.Repeat = &repeat

					// Save the collection
					rearrangeMangaCollectionLists(mangaCollection)
					lp.localManager.UpdateLocalMangaCollection(mangaCollection)
					lp.localManager.SetHasLocalChanges(true)
					return nil
				}
			}
		}
	}

	return ErrMediaNotFound
}

// DeleteEntry isn't supported for the local platform, always returns an error.
func (lp *OfflinePlatform) DeleteEntry(ctx context.Context, mediaID, entryId int) error {
	return ErrActionNotSupported
}

func (lp *OfflinePlatform) GetAnime(ctx context.Context, mediaID int) (*media.Anime, error) {
	if lp.localManager.GetLocalAnimeCollection().IsPresent() {
		animeCollection := lp.localManager.GetLocalAnimeCollection().MustGet()

		// Find the entry
		for _, list := range animeCollection.MediaListCollection.Lists {
			for _, entry := range list.Entries {
				if entry.GetMedia().GetID() == mediaID {
					return entry.Media, nil
				}
			}
		}
	}

	return nil, ErrMediaNotFound
}

func (lp *OfflinePlatform) GetAnimeByMalID(ctx context.Context, malID int) (*media.Anime, error) {
	if lp.localManager.GetLocalAnimeCollection().IsPresent() {
		animeCollection := lp.localManager.GetLocalAnimeCollection().MustGet()

		// Find the entry
		for _, list := range animeCollection.MediaListCollection.Lists {
			for _, entry := range list.Entries {
				if entry.GetMedia().GetIDMal() != nil && *entry.GetMedia().GetIDMal() == malID {
					return entry.Media, nil
				}
			}
		}
	}

	return nil, ErrMediaNotFound
}

// GetAnimeDetails isn't supported for the local platform, always returns an empty struct.
func (lp *OfflinePlatform) GetAnimeDetails(ctx context.Context, mediaID int) (*media.AnimeDetails, error) {
	return &media.AnimeDetails{}, nil
}

// GetAnimeWithRelations isn't supported for the local platform, always returns an error.
func (lp *OfflinePlatform) GetAnimeWithRelations(ctx context.Context, mediaID int) (*media.CompleteAnime, error) {
	return nil, ErrActionNotSupported
}

func (lp *OfflinePlatform) GetManga(ctx context.Context, mediaID int) (*media.Manga, error) {
	if lp.localManager.GetLocalMangaCollection().IsPresent() {
		mangaCollection := lp.localManager.GetLocalMangaCollection().MustGet()

		// Find the entry
		for _, list := range mangaCollection.MediaListCollection.Lists {
			for _, entry := range list.Entries {
				if entry.GetMedia().GetID() == mediaID {
					return entry.Media, nil
				}
			}
		}
	}

	return nil, ErrMediaNotFound
}

// GetMangaDetails isn't supported for the local platform, always returns an empty struct.
func (lp *OfflinePlatform) GetMangaDetails(ctx context.Context, mediaID int) (*media.MangaDetails, error) {
	return &media.MangaDetails{}, nil
}

func (lp *OfflinePlatform) GetAnimeCollection(ctx context.Context, bypassCache bool) (*media.AnimeCollection, error) {
	if lp.localManager.GetLocalAnimeCollection().IsPresent() {
		return lp.localManager.GetLocalAnimeCollection().MustGet(), nil
	} else {
		return nil, ErrNoLocalAnimeCollection
	}
}

func (lp *OfflinePlatform) GetRawAnimeCollection(ctx context.Context, bypassCache bool) (*media.AnimeCollection, error) {
	if lp.localManager.GetLocalAnimeCollection().IsPresent() {
		return lp.localManager.GetLocalAnimeCollection().MustGet(), nil
	} else {
		return nil, ErrNoLocalAnimeCollection
	}
}

// RefreshAnimeCollection is a no-op, always returns the local anime collection.
func (lp *OfflinePlatform) RefreshAnimeCollection(ctx context.Context) (*media.AnimeCollection, error) {
	animeCollection, ok := lp.localManager.GetLocalAnimeCollection().Get()
	if !ok {
		return nil, ErrNoLocalAnimeCollection
	}

	return animeCollection, nil
}

func (lp *OfflinePlatform) GetAnimeCollectionWithRelations(ctx context.Context) (*media.AnimeCollectionWithRelations, error) {
	return nil, ErrActionNotSupported
}

func (lp *OfflinePlatform) GetMangaCollection(ctx context.Context, bypassCache bool) (*media.MangaCollection, error) {
	if lp.localManager.GetLocalMangaCollection().IsPresent() {
		return lp.localManager.GetLocalMangaCollection().MustGet(), nil
	} else {
		return nil, ErrorNoLocalMangaCollection
	}
}

func (lp *OfflinePlatform) GetRawMangaCollection(ctx context.Context, bypassCache bool) (*media.MangaCollection, error) {
	if lp.localManager.GetLocalMangaCollection().IsPresent() {
		return lp.localManager.GetLocalMangaCollection().MustGet(), nil
	} else {
		return nil, ErrorNoLocalMangaCollection
	}
}

func (lp *OfflinePlatform) RefreshMangaCollection(ctx context.Context) (*media.MangaCollection, error) {
	mangaCollection, ok := lp.localManager.GetLocalMangaCollection().Get()
	if !ok {
		return nil, ErrorNoLocalMangaCollection
	}

	return mangaCollection, nil
}

// AddMediaToCollection isn't supported for the local platform, always returns an error.
func (lp *OfflinePlatform) AddMediaToCollection(ctx context.Context, mIds []int) error {
	return ErrActionNotSupported
}

// GetStudioDetails isn't supported for the local platform, always returns an empty struct
func (lp *OfflinePlatform) GetStudioDetails(ctx context.Context, studioID int) (*media.StudioDetails, error) {
	return &media.StudioDetails{}, nil
}

// GetBangumiClient 离线模式无 Bangumi 客户端，返回 nil。
func (lp *OfflinePlatform) GetBangumiClient() *bangumi.Client {
	return nil
}

func (lp *OfflinePlatform) GetViewerStats(ctx context.Context) (*media.ViewerStats, error) {
	return nil, ErrActionNotSupported
}

func (lp *OfflinePlatform) GetAnimeAiringSchedule(ctx context.Context) (*media.AnimeAiringSchedule, error) {
	return nil, ErrActionNotSupported
}
