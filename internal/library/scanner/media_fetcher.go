package scanner

import (
	"context"
	"errors"
	"seanime/internal/api/animeofflinedb"
	"seanime/internal/api/mal"
	"seanime/internal/api/metadata"
	"seanime/internal/api/metadata_provider"
	"seanime/internal/customsource"
	"seanime/internal/hook"
	"seanime/internal/library/anime"
	medialib "seanime/internal/media"
	"seanime/internal/platforms/platform"
	"seanime/internal/util"
	"seanime/internal/util/limiter"
	"seanime/internal/util/parallel"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	lop "github.com/samber/lo/parallel"
)

// MediaFetcher holds all medialib.Anime that will be used for the comparison process
type MediaFetcher struct {
	AllMedia                     []*anime.NormalizedMedia
	CollectionMediaIds           []int
	UnknownMediaIds              []int // Media IDs that are not in the user's collection
	AnimeCollectionWithRelations *medialib.AnimeCollectionWithRelations
	ScanLogger                   *ScanLogger
}

type MediaFetcherOptions struct {
	Enhanced                   bool
	EnhanceWithOfflineDatabase bool
	PlatformRef                *util.Ref[platform.Platform]
	MetadataProviderRef        *util.Ref[metadata_provider.Provider]
	LocalFiles                 []*anime.LocalFile
	CompleteAnimeCache         *medialib.CompleteAnimeCache
	Logger                     *zerolog.Logger
	AnilistRateLimiter         *limiter.Limiter
	DisableAnimeCollection     bool
	ScanLogger                 *ScanLogger
	// used for adding custom sources
	OptionalAnimeCollection *medialib.AnimeCollection
}

// NewMediaFetcher
// Calling this method will kickstart the fetch process
// When enhancing is false, MediaFetcher.AllMedia will be all medialib.Anime from the user's AniList collection.
// When enhancing is true, MediaFetcher.AllMedia will be medialib.Anime for each unique, parsed anime title and their relations.
func NewMediaFetcher(ctx context.Context, opts *MediaFetcherOptions) (ret *MediaFetcher, retErr error) {
	defer util.HandlePanicInModuleWithError("library/scanner/NewMediaFetcher", &retErr)

	if opts.PlatformRef.IsAbsent() ||
		opts.LocalFiles == nil ||
		opts.CompleteAnimeCache == nil ||
		opts.MetadataProviderRef.IsAbsent() ||
		opts.Logger == nil ||
		opts.AnilistRateLimiter == nil {
		return nil, errors.New("missing options")
	}

	mf := new(MediaFetcher)
	mf.ScanLogger = opts.ScanLogger

	opts.Logger.Debug().
		Any("enhanced", opts.Enhanced).
		Msg("media fetcher: Creating media fetcher")

	if mf.ScanLogger != nil {
		mf.ScanLogger.LogMediaFetcher(zerolog.InfoLevel).
			Msg("Creating media fetcher")
	}

	// Invoke ScanMediaFetcherStarted hook
	event := &ScanMediaFetcherStartedEvent{
		Enhanced:                   opts.Enhanced,
		EnhanceWithOfflineDatabase: opts.EnhanceWithOfflineDatabase,
		DisableAnimeCollection:     opts.DisableAnimeCollection,
	}
	_ = hook.GlobalHookManager.OnScanMediaFetcherStarted().Trigger(event)
	opts.Enhanced = event.Enhanced

	// +---------------------+
	// |     All media       |
	// +---------------------+

	// Fetch latest user's AniList collection
	animeCollectionWithRelations, err := opts.PlatformRef.Get().GetAnimeCollectionWithRelations(ctx)
	if err != nil {
		return nil, err
	}

	mf.AnimeCollectionWithRelations = animeCollectionWithRelations

	// Temporary slice to hold CompleteAnime before conversion
	allCompleteAnime := make([]*medialib.CompleteAnime, 0)

	if !opts.DisableAnimeCollection {
		// For each collection entry, append the media to AllMedia
		for _, list := range animeCollectionWithRelations.GetMediaListCollection().GetLists() {
			for _, entry := range list.GetEntries() {
				allCompleteAnime = append(allCompleteAnime, entry.GetMedia())

				// +---------------------+
				// |        Cache        |
				// +---------------------+
				// We assume the CompleteAnimeCache is empty. Add media to cache.
				opts.CompleteAnimeCache.Set(entry.GetMedia().ID, entry.GetMedia())
			}
		}
		// Handle custom sources
		// Devnote: For now we just get them from opts.AnimeCollection but in the future we could introduce a new method for custom sources to return many CompleteAnime at once
		// right now custom source media wont have any relations data
		if opts.OptionalAnimeCollection != nil {
			for _, list := range opts.OptionalAnimeCollection.GetMediaListCollection().GetLists() {
				if list == nil {
					continue
				}
				for _, entry := range list.GetEntries() {
					if entry == nil || entry.GetMedia() == nil || !customsource.IsExtensionId(entry.GetMedia().GetID()) {
						continue
					}
					allCompleteAnime = append(allCompleteAnime, entry.GetMedia().ToCompleteAnime())
				}
			}
		}
	}

	if mf.ScanLogger != nil {
		mf.ScanLogger.LogMediaFetcher(zerolog.DebugLevel).
			Int("count", len(allCompleteAnime)).
			Msg("Fetched media from AniList collection")
	}

	//--------------------------------------------

	// Get the media IDs from the collection
	mf.CollectionMediaIds = lop.Map(allCompleteAnime, func(m *medialib.CompleteAnime, index int) int {
		return m.ID
	})

	//--------------------------------------------

	// +---------------------+
	// |  Enhanced (Legacy)  |
	// +---------------------+

	// If enhancing (legacy) is on, scan media from local files and get their relations
	if opts.Enhanced && !opts.EnhanceWithOfflineDatabase {

		_, ok := FetchMediaFromLocalFiles(
			ctx,
			opts.PlatformRef.Get(),
			opts.LocalFiles,
			opts.CompleteAnimeCache, // CompleteAnimeCache will be populated on success
			opts.MetadataProviderRef.Get(),
			opts.AnilistRateLimiter,
			mf.ScanLogger,
		)
		if ok {
			// We assume the CompleteAnimeCache is populated.
			// Safe to overwrite allCompleteAnime with the cache content
			// because the cache will contain all media from the user's collection AND scanned ones
			allCompleteAnime = make([]*medialib.CompleteAnime, 0)
			opts.CompleteAnimeCache.Range(func(key int, value *medialib.CompleteAnime) bool {
				allCompleteAnime = append(allCompleteAnime, value)
				return true
			})
		}
	}

	mf.AllMedia = NormalizedMediaFromAnilistComplete(allCompleteAnime)

	// +-------------------------+
	// |  Enhanced (Offline DB)  |
	// +-------------------------+
	// When enhanced mode is on, fetch anime-offline-database to provide more matching candidates

	if opts.Enhanced && opts.EnhanceWithOfflineDatabase {
		if mf.ScanLogger != nil {
			mf.ScanLogger.LogMediaFetcher(zerolog.DebugLevel).
				Msg("Fetching anime-offline-database for enhanced matching")
		}

		// build existing media IDs map for filtering
		existingMediaIDs := make(map[int]bool, len(mf.AllMedia))
		for _, m := range mf.AllMedia {
			existingMediaIDs[m.ID] = true
		}

		offlineMedia, err := animeofflinedb.FetchAndConvertDatabase(existingMediaIDs)
		if err != nil {
			if mf.ScanLogger != nil {
				mf.ScanLogger.LogMediaFetcher(zerolog.WarnLevel).
					Err(err).
					Msg("Failed to fetch anime-offline-database, continuing without it")
			}
		} else {
			if mf.ScanLogger != nil {
				mf.ScanLogger.LogMediaFetcher(zerolog.DebugLevel).
					Int("offlineMediaCount", len(offlineMedia)).
					Msg("Added media from anime-offline-database")
			}

			// Append offline media to AllMedia
			mf.AllMedia = append(mf.AllMedia, offlineMedia...)
		}
	}

	// +---------------------+
	// |   Unknown media     |
	// +---------------------+
	// Media that are not in the user's collection

	// Get the media that are not in the user's collection
	unknownMedia := lo.Filter(mf.AllMedia, func(m *anime.NormalizedMedia, _ int) bool {
		return !lo.Contains(mf.CollectionMediaIds, m.ID)
	})
	// Get the media IDs that are not in the user's collection
	mf.UnknownMediaIds = lop.Map(unknownMedia, func(m *anime.NormalizedMedia, _ int) int {
		return m.ID
	})

	if mf.ScanLogger != nil {
		mf.ScanLogger.LogMediaFetcher(zerolog.DebugLevel).
			Int("unknownMediaCount", len(mf.UnknownMediaIds)).
			Int("allMediaCount", len(mf.AllMedia)).
			Msg("Finished creating media fetcher")
	}

	// Invoke ScanMediaFetcherCompleted hook
	completedEvent := &ScanMediaFetcherCompletedEvent{
		AllMedia:        mf.AllMedia,
		UnknownMediaIds: mf.UnknownMediaIds,
	}
	_ = hook.GlobalHookManager.OnScanMediaFetcherCompleted().Trigger(completedEvent)
	mf.AllMedia = completedEvent.AllMedia
	mf.UnknownMediaIds = completedEvent.UnknownMediaIds

	return mf, nil
}

func NormalizedMediaFromAnilistComplete(c []*medialib.CompleteAnime) []*anime.NormalizedMedia {
	normalizedMediaMap := make(map[int]*anime.NormalizedMedia)

	// Convert CompleteAnime to NormalizedMedia and flatten relations
	for _, m := range c {
		if _, found := normalizedMediaMap[m.ID]; !found {
			normalizedMediaMap[m.ID] = anime.NewNormalizedMedia(m.ToAnime())
		}

		// Process relations
		if m.Relations != nil && m.Relations.Edges != nil && len(m.Relations.Edges) > 0 {
			for _, edgeM := range m.Relations.Edges {
				if edgeM.Node == nil || edgeM.Node.Format == nil || edgeM.RelationType == nil {
					continue
				}
				if *edgeM.Node.Format != medialib.MediaFormatMovie &&
					*edgeM.Node.Format != medialib.MediaFormatOva &&
					*edgeM.Node.Format != medialib.MediaFormatSpecial &&
					*edgeM.Node.Format != medialib.MediaFormatTv {
					continue
				}
				if *edgeM.RelationType != medialib.MediaRelationPrequel &&
					*edgeM.RelationType != medialib.MediaRelationSequel &&
					*edgeM.RelationType != medialib.MediaRelationSpinOff &&
					*edgeM.RelationType != medialib.MediaRelationAlternative &&
					*edgeM.RelationType != medialib.MediaRelationParent {
					continue
				}
				// Make sure we don't overwrite the original media in the map
				if _, found := normalizedMediaMap[edgeM.Node.ID]; !found {
					normalizedMediaMap[edgeM.Node.ID] = anime.NewNormalizedMedia(edgeM.Node)
				}
			}
		}
	}

	ret := make([]*anime.NormalizedMedia, 0, len(normalizedMediaMap))

	for _, m := range normalizedMediaMap {
		ret = append(ret, m)
	}

	return ret
}

//----------------------------------------------------------------------------------------------------------------------

// FetchMediaFromLocalFiles gets media and their relations from local file titles.
// It retrieves unique titles from local files,
// fetches mal.SearchResultAnime from MAL,
// uses these search results to get AniList IDs using metadata.AnimeMetadata mappings,
// queries AniList to retrieve all medialib.Anime using media.GetBaseAnimeById and their relations using medialib.FetchMediaTree.
// It does not return an error if one of the steps fails.
// It returns the scanned media and a boolean indicating whether the process was successful.
func FetchMediaFromLocalFiles(
	ctx context.Context,
	platform platform.Platform,
	localFiles []*anime.LocalFile,
	completeAnime *medialib.CompleteAnimeCache,
	metadataProvider metadata_provider.Provider,
	anilistRateLimiter *limiter.Limiter,
	scanLogger *ScanLogger,
) (ret []*medialib.CompleteAnime, ok bool) {
	defer util.HandlePanicInModuleThen("library/scanner/FetchMediaFromLocalFiles", func() {
		ok = false
	})

	if scanLogger != nil {
		scanLogger.LogMediaFetcher(zerolog.DebugLevel).
			Str("module", "Enhanced").
			Msg("Fetching media from local files")
	}

	rateLimiter := limiter.NewLimiter(time.Second, 20)
	rateLimiter2 := limiter.NewLimiter(time.Second, 20)

	// Get titles
	titles := anime.GetUniqueAnimeTitlesFromLocalFiles(localFiles)

	if scanLogger != nil {
		scanLogger.LogMediaFetcher(zerolog.DebugLevel).
			Str("module", "Enhanced").
			Str("context", spew.Sprint(titles)).
			Msg("Parsed titles from local files")
	}

	// +---------------------+
	// |     MyAnimeList     |
	// +---------------------+

	// Get MAL media from titles
	malSR := parallel.NewSettledResults[string, *mal.SearchResultAnime](titles)
	malSR.AllSettled(func(title string, index int) (*mal.SearchResultAnime, error) {
		rateLimiter.Wait()
		return mal.AdvancedSearchWithMAL(title)
	})
	malRes, ok := malSR.GetFulfilledResults()
	if !ok {
		return nil, false
	}

	// Get duplicate-free version of MAL media
	malMedia := lo.UniqBy(*malRes, func(res *mal.SearchResultAnime) int { return res.ID })
	// Get the MAL media IDs
	malIds := lop.Map(malMedia, func(n *mal.SearchResultAnime, index int) int { return n.ID })

	if scanLogger != nil {
		scanLogger.LogMediaFetcher(zerolog.DebugLevel).
			Str("module", "Enhanced").
			Str("context", spew.Sprint(lo.Map(malMedia, func(n *mal.SearchResultAnime, _ int) string {
				return n.Name
			}))).
			Msg("Fetched MAL media from titles")
	}

	// +---------------------+
	// |       Animap        |
	// +---------------------+

	// Get Animap mappings for each MAL ID and store them in `metadataProvider`
	// This step is necessary because MAL doesn't provide AniList IDs and some MAL media don't exist on AniList
	lop.ForEach(malIds, func(id int, index int) {
		rateLimiter2.Wait()
		//_, _ = metadataProvider.GetAnimeMetadata(metadata.MalPlatform, id)
		_, _ = metadataProvider.GetCache().GetOrSet(metadata_provider.GetAnimeMetadataCacheKey(metadata.MalPlatform, id), func() (*metadata.AnimeMetadata, error) {
			res, err := metadataProvider.GetAnimeMetadata(metadata.MalPlatform, id)
			return res, err
		})
	})

	// +---------------------+
	// |       AniList       |
	// +---------------------+

	// Retrieve the AniList IDs from the Animap mappings stored in the cache
	anilistIds := make([]int, 0)
	metadataProvider.GetCache().Range(func(key string, value *metadata.AnimeMetadata) bool {
		if value != nil {
			anilistIds = append(anilistIds, value.GetMappings().AnilistId)
		}
		return true
	})

	// Fetch all media from the AniList IDs
	anilistMedia := make([]*medialib.CompleteAnime, 0)
	lop.ForEach(anilistIds, func(id int, index int) {
		anilistRateLimiter.Wait()
		var media *medialib.CompleteAnime
		var err error
		media, err = platform.GetAnimeWithRelations(ctx, id)
		if err != nil {
			baseMedia, lErr := platform.GetAnime(ctx, id)
			if lErr == nil {
				media = baseMedia.ToCompleteAnime()
				err = nil
			}
		}
		if err == nil {
			anilistMedia = append(anilistMedia, media)
			if scanLogger != nil {
				scanLogger.LogMediaFetcher(zerolog.DebugLevel).
					Str("module", "Enhanced").
					Str("title", media.GetTitleSafe()).
					Msg("Fetched Anilist media")
			}
		} else {
			if scanLogger != nil {
				scanLogger.LogMediaFetcher(zerolog.WarnLevel).
					Str("module", "Enhanced").
					Int("id", id).
					Msg("Failed to fetch Anilist media")
			}
		}
	})

	if scanLogger != nil {
		scanLogger.LogMediaFetcher(zerolog.DebugLevel).
			Str("module", "Enhanced").
			Str("context", spew.Sprint(lo.Map(anilistMedia, func(n *medialib.CompleteAnime, _ int) string {
				return n.GetTitleSafe()
			}))).
			Msg("Fetched Anilist media from MAL ids")
	}

	// +---------------------+
	// |     MediaTree       |
	// +---------------------+

	// Create a new tree that will hold the fetched relations
	// /!\ This is redundant because we already have a cache, but `FetchMediaTree` needs its
	tree := medialib.NewCompleteAnimeRelationTree()

	start := time.Now()
	// For each media, fetch its relations
	// The relations are fetched in parallel and added to `completeAnime`
	lop.ForEach(anilistMedia, func(m *medialib.CompleteAnime, index int) {
		// We ignore errors because we want to continue even if one of the media fails
		// Bangumi 锚点：CompleteAnime 不再自带递归遍历；GetAnimeWithRelations 返回的一层关系
		// 直接入树，关系节点按需补全为完整条目。
		if m == nil {
			return
		}
		tree.Set(m.ID, m)
		completeAnime.Set(m.ID, m)
		if m.Relations == nil || m.Relations.Edges == nil {
			return
		}
		for _, edge := range m.Relations.Edges {
			if edge == nil || edge.Node == nil || edge.Node.ID == 0 || tree.Has(edge.Node.ID) {
				continue
			}
			anilistRateLimiter.Wait()
			complete, err := platform.GetAnimeWithRelations(context.Background(), edge.Node.ID)
			if err != nil || complete == nil {
				continue
			}
			tree.Set(complete.ID, complete)
			completeAnime.Set(complete.ID, complete)
		}
	})

	// +---------------------+
	// |        Cache        |
	// +---------------------+

	// Retrieve all media from the cache
	scanned := make([]*medialib.CompleteAnime, 0)
	completeAnime.Range(func(key int, value *medialib.CompleteAnime) bool {
		scanned = append(scanned, value)
		return true
	})

	if scanLogger != nil {
		scanLogger.LogMediaFetcher(zerolog.InfoLevel).
			Str("module", "Enhanced").
			Int("ms", int(time.Since(start).Milliseconds())).
			Int("count", len(scanned)).
			Str("context", spew.Sprint(lo.Map(scanned, func(n *medialib.CompleteAnime, _ int) string {
				return n.GetTitleSafe()
			}))).
			Msg("Finished fetching media from local files")
	}

	return scanned, true
}
