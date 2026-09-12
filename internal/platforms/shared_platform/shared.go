package shared_platform

import (
	"context"
	"errors"
	"seanime/internal/customsource"
	"seanime/internal/database/db"
	"seanime/internal/extension"
	"seanime/internal/hook"
	"seanime/internal/media"
	"seanime/internal/platforms/platform"
	"seanime/internal/util"
	"seanime/internal/util/result"
	"time"

	"github.com/rs/zerolog"
	"github.com/samber/lo"
)

// 换锚说明（Phase 2 Wave B）：PlatformHelper 的缓存与事件结构全部切换为
// internal/media 镜像类型；customsource 已同步迁移至 media 类型。
// 原 BuildAnimeAiringSchedule（依赖 AniList 逐集播出时间戳）被移除——
// Bangumi 无对应数据，放送表由各平台层降级返回空结构。

type PlatformHelper struct {
	logger              *zerolog.Logger
	customSourceManager *customsource.Manager
	baseAnimeCache      *result.BoundedCache[int, *media.Anime]
	baseMangaCache      *result.BoundedCache[int, *media.Manga]
	completeAnimeCache  *result.BoundedCache[int, *media.CompleteAnime]
	extensionBankRef    *util.Ref[*extension.UnifiedBank]
}

func NewPlatformHelper(extensionBankRef *util.Ref[*extension.UnifiedBank], db *db.Database, logger *zerolog.Logger) *PlatformHelper {
	helper := &PlatformHelper{
		logger:              logger,
		baseAnimeCache:      result.NewBoundedCache[int, *media.Anime](50),
		baseMangaCache:      result.NewBoundedCache[int, *media.Manga](50),
		completeAnimeCache:  result.NewBoundedCache[int, *media.CompleteAnime](10),
		extensionBankRef:    extensionBankRef,
		customSourceManager: customsource.NewManager(extensionBankRef, db, logger),
	}

	return helper
}

func (h *PlatformHelper) Close() {
	if h.customSourceManager != nil {
		h.customSourceManager.Close()
	}
}

func (h *PlatformHelper) ClearCache() {
	h.baseAnimeCache.Clear()
	h.baseMangaCache.Clear()
	h.completeAnimeCache.Clear()
}

func (h *PlatformHelper) GetCustomSourceManager() *customsource.Manager {
	return h.customSourceManager
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Custom Source
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (h *PlatformHelper) HandleCustomSourceAnime(ctx context.Context, mediaID int) (*media.Anime, bool, error) {
	if h.customSourceManager == nil {
		return nil, false, nil
	}

	if customSource, localId, isCustom, hasExtension := h.customSourceManager.GetProviderFromId(mediaID); isCustom {
		if !hasExtension {
			return nil, true, errors.New("custom source does not exist or identifier has changed")
		}
		ret, err := customSource.GetProvider().GetAnime(ctx, []int{localId})
		if err != nil {
			return nil, true, err
		}
		if len(ret) == 0 {
			return nil, true, errors.New("no anime found")
		}
		media := ret[0]
		customsource.NormalizeMedia(customSource.GetExtensionIdentifier(), customSource.GetID(), media)
		return media, true, nil
	}

	return nil, false, nil
}

func (h *PlatformHelper) HandleCustomSourceAnimeDetails(ctx context.Context, mediaID int) (*media.AnimeDetails, bool, error) {
	if h.customSourceManager == nil {
		return nil, false, nil
	}

	if customSource, localId, isCustom, hasExtension := h.customSourceManager.GetProviderFromId(mediaID); isCustom {
		if !hasExtension {
			return nil, true, errors.New("custom source does not exist or identifier has changed")
		}
		ret, err := customSource.GetProvider().GetAnimeDetails(ctx, localId)
		if err != nil {
			return nil, true, err
		}
		customsource.NormalizeMedia(customSource.GetExtensionIdentifier(), customSource.GetID(), ret)
		return ret, true, nil
	}

	return nil, false, nil
}

func (h *PlatformHelper) HandleCustomSourceAnimeWithRelations(ctx context.Context, mediaID int) (*media.CompleteAnime, bool, error) {
	if h.customSourceManager == nil {
		return nil, false, nil
	}

	if customSource, localId, isCustom, hasExtension := h.customSourceManager.GetProviderFromId(mediaID); isCustom {
		if !hasExtension {
			return nil, true, errors.New("custom source does not exist or identifier has changed")
		}
		ret, err := customSource.GetProvider().GetAnimeWithRelations(ctx, localId)
		if err != nil {
			return nil, true, err
		}
		customsource.NormalizeMedia(customSource.GetExtensionIdentifier(), customSource.GetID(), ret)
		return ret, true, nil
	}

	return nil, false, nil
}

func (h *PlatformHelper) HandleCustomSourceManga(ctx context.Context, mediaID int) (*media.Manga, bool, error) {
	if h.customSourceManager == nil {
		return nil, false, nil
	}

	if customSource, localId, isCustom, hasExtension := h.customSourceManager.GetProviderFromId(mediaID); isCustom {
		if !hasExtension {
			return nil, true, errors.New("custom source does not exist or identifier has changed")
		}
		ret, err := customSource.GetProvider().GetManga(ctx, []int{localId})
		if err != nil {
			return nil, true, err
		}
		if len(ret) == 0 {
			return nil, true, errors.New("no manga found")
		}
		media := ret[0]
		customsource.NormalizeMedia(customSource.GetExtensionIdentifier(), customSource.GetID(), media)
		return media, true, nil
	}

	return nil, false, nil
}

func (h *PlatformHelper) HandleCustomSourceMangaDetails(ctx context.Context, mediaID int) (*media.MangaDetails, bool, error) {
	if h.customSourceManager == nil {
		return nil, false, nil
	}

	if customSource, localId, isCustom, hasExtension := h.customSourceManager.GetProviderFromId(mediaID); isCustom {
		if !hasExtension {
			return nil, true, errors.New("custom source does not exist or identifier has changed")
		}
		ret, err := customSource.GetProvider().GetMangaDetails(ctx, localId)
		if err != nil {
			return nil, true, err
		}
		customsource.NormalizeMedia(customSource.GetExtensionIdentifier(), customSource.GetID(), ret)
		return ret, true, nil
	}

	return nil, false, nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Cache
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (h *PlatformHelper) GetCachedBaseAnime(mediaID int) (*media.Anime, bool) {
	return h.baseAnimeCache.Get(mediaID)
}

func (h *PlatformHelper) SetCachedBaseAnime(mediaID int, anime *media.Anime) {
	h.baseAnimeCache.SetT(mediaID, anime, time.Minute*30)
}

func (h *PlatformHelper) GetCachedBaseManga(mediaID int) (*media.Manga, bool) {
	return h.baseMangaCache.Get(mediaID)
}

func (h *PlatformHelper) SetCachedBaseManga(mediaID int, manga *media.Manga) {
	h.baseMangaCache.SetT(mediaID, manga, time.Minute*30)
}

func (h *PlatformHelper) GetCachedCompleteAnime(mediaID int) (*media.CompleteAnime, bool) {
	return h.completeAnimeCache.Get(mediaID)
}

func (h *PlatformHelper) SetCachedCompleteAnime(mediaID int, anime *media.CompleteAnime) {
	h.completeAnimeCache.SetT(mediaID, anime, 4*time.Hour)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Hook Events
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (h *PlatformHelper) TriggerGetAnimeEvent(anime *media.Anime) (*media.Anime, error) {
	event := new(platform.GetAnimeEvent)
	event.Anime = anime
	err := hook.GlobalHookManager.OnGetAnime().Trigger(event)
	if err != nil {
		return nil, err
	}
	return event.Anime, nil
}

func (h *PlatformHelper) TriggerGetAnimeDetailsEvent(anime *media.AnimeDetails) (*media.AnimeDetails, error) {
	event := new(platform.GetAnimeDetailsEvent)
	event.Anime = anime
	err := hook.GlobalHookManager.OnGetAnimeDetails().Trigger(event)
	if err != nil {
		return nil, err
	}
	return event.Anime, nil
}

func (h *PlatformHelper) TriggerGetMangaEvent(manga *media.Manga) (*media.Manga, error) {
	event := new(platform.GetMangaEvent)
	event.Manga = manga
	err := hook.GlobalHookManager.OnGetManga().Trigger(event)
	if err != nil {
		return nil, err
	}
	return event.Manga, nil
}

func (h *PlatformHelper) TriggerGetStudioDetailsEvent(studio *media.StudioDetails) (*media.StudioDetails, error) {
	event := new(platform.GetStudioDetailsEvent)
	event.Studio = studio
	err := hook.GlobalHookManager.OnGetStudioDetails().Trigger(event)
	if err != nil {
		return nil, err
	}
	return event.Studio, nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Custom Source
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (h *PlatformHelper) MergeCustomSourceAnimeEntries(collection *media.AnimeCollection) {
	if h.customSourceManager != nil {
		h.customSourceManager.MergeAnimeEntries(collection)
	}
}

func (h *PlatformHelper) MergeCustomSourceMangaEntries(collection *media.MangaCollection) {
	if h.customSourceManager != nil {
		h.customSourceManager.MergeMangaEntries(collection)
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Update
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (h *PlatformHelper) HandleCustomSourceUpdateEntry(ctx context.Context, mediaID int, status *media.MediaListStatus, scoreRaw *int, progress *int, startedAt *media.FuzzyDateInput, completedAt *media.FuzzyDateInput) (bool, error) {
	if h.customSourceManager != nil && customsource.IsExtensionId(mediaID) {
		err := h.customSourceManager.UpdateEntry(ctx, mediaID, status, scoreRaw, progress, startedAt, completedAt)
		return true, err
	}
	return false, nil
}

func (h *PlatformHelper) HandleCustomSourceUpdateEntryProgress(ctx context.Context, mediaID int, progress int, totalCount *int) (bool, error) {
	if h.customSourceManager != nil && customsource.IsExtensionId(mediaID) {
		err := h.customSourceManager.UpdateEntryProgress(ctx, mediaID, progress, totalCount)
		return true, err
	}
	return false, nil
}

func (h *PlatformHelper) HandleCustomSourceUpdateEntryRepeat(ctx context.Context, mediaID int, repeat int) (bool, error) {
	if h.customSourceManager != nil && customsource.IsExtensionId(mediaID) {
		err := h.customSourceManager.UpdateEntryRepeat(ctx, mediaID, repeat)
		return true, err
	}
	return false, nil
}

func (h *PlatformHelper) HandleCustomSourceDeleteEntry(ctx context.Context, mediaID int, entryId int) (bool, error) {
	if h.customSourceManager != nil && customsource.IsExtensionId(mediaID) {
		err := h.customSourceManager.DeleteEntry(ctx, mediaID, entryId)
		return true, err
	}
	return false, nil
}

func (h *PlatformHelper) TriggerUpdateEntryHooks(ctx context.Context, mediaID int, status *media.MediaListStatus, scoreRaw *int, progress *int, startedAt *media.FuzzyDateInput, completedAt *media.FuzzyDateInput, updateFunc func(event *platform.PreUpdateEntryEvent) error) error {
	// Trigger pre-update hook
	event := new(platform.PreUpdateEntryEvent)
	event.MediaID = &mediaID
	event.Status = status
	event.ScoreRaw = scoreRaw
	event.Progress = progress
	event.StartedAt = startedAt
	event.CompletedAt = completedAt

	err := hook.GlobalHookManager.OnPreUpdateEntry().Trigger(event)
	if err != nil {
		return err
	}

	if event.DefaultPrevented {
		return nil
	}

	// Execute the update
	err = updateFunc(event)
	if err != nil {
		return err
	}

	// Trigger post-update hook
	postEvent := new(platform.PostUpdateEntryEvent)
	postEvent.MediaID = &mediaID
	err = hook.GlobalHookManager.OnPostUpdateEntry().Trigger(postEvent)
	return err
}

// TriggerUpdateEntryProgressHooks triggers pre and post update entry progress hooks
func (h *PlatformHelper) TriggerUpdateEntryProgressHooks(ctx context.Context, mediaID int, progress int, totalCount *int, updateFunc func(event *platform.PreUpdateEntryProgressEvent) error) error {
	// Trigger pre-update hook
	event := new(platform.PreUpdateEntryProgressEvent)
	event.MediaID = &mediaID
	event.Progress = &progress
	event.TotalCount = totalCount
	currentStatus := media.MediaListStatusCurrent
	event.Status = &currentStatus

	_ = hook.GlobalHookManager.OnPreUpdateEntryProgress().Trigger(event)

	if event.DefaultPrevented {
		return nil
	}

	// Execute the update
	err := updateFunc(event)
	if err != nil {
		return err
	}

	// Trigger post-update hook
	postEvent := new(platform.PostUpdateEntryProgressEvent)
	postEvent.MediaID = &mediaID
	_ = hook.GlobalHookManager.OnPostUpdateEntryProgress().Trigger(postEvent)
	return err
}

func (h *PlatformHelper) TriggerUpdateEntryRepeatHooks(ctx context.Context, mediaID int, repeat int, updateFunc func(event *platform.PreUpdateEntryRepeatEvent) error) error {
	// Trigger pre-update hook
	event := new(platform.PreUpdateEntryRepeatEvent)
	event.MediaID = &mediaID
	event.Repeat = &repeat

	err := hook.GlobalHookManager.OnPreUpdateEntryRepeat().Trigger(event)
	if err != nil {
		return err
	}

	if event.DefaultPrevented {
		return nil
	}

	// Execute the update
	err = updateFunc(event)
	if err != nil {
		return err
	}

	// Trigger post-update hook
	postEvent := new(platform.PostUpdateEntryRepeatEvent)
	postEvent.MediaID = &mediaID
	err = hook.GlobalHookManager.OnPostUpdateEntryRepeat().Trigger(postEvent)
	return err
}

func (h *PlatformHelper) TriggerDeleteEntryHooks(ctx context.Context, mediaID int, entryId int, deleteFunc func(event *platform.PreDeleteEntryEvent) error) error {
	// Trigger pre-delete hook
	event := new(platform.PreDeleteEntryEvent)
	event.MediaID = &mediaID
	event.EntryID = &entryId

	err := hook.GlobalHookManager.OnPreDeleteEntry().Trigger(event)
	if err != nil {
		return err
	}

	if event.DefaultPrevented {
		return nil
	}

	// Execute the deletion
	err = deleteFunc(event)
	if err != nil {
		return err
	}

	// Trigger post-delete hook
	postEvent := new(platform.PostDeleteEntryEvent)
	postEvent.MediaID = &mediaID
	postEvent.EntryID = &entryId
	err = hook.GlobalHookManager.OnPostDeleteEntry().Trigger(postEvent)
	return err
}

func (h *PlatformHelper) FilterOutCustomAnimeLists(lists []*media.AnimeCollection_MediaListCollection_Lists) []*media.AnimeCollection_MediaListCollection_Lists {
	return lo.Filter(lists, func(list *media.AnimeCollection_MediaListCollection_Lists, _ int) bool {
		return list.Status != nil
	})
}

func (h *PlatformHelper) FilterOutCustomMangaLists(lists []*media.MangaCollection_MediaListCollection_Lists) []*media.MangaCollection_MediaListCollection_Lists {
	return lo.Filter(lists, func(list *media.MangaCollection_MediaListCollection_Lists, _ int) bool {
		return list.Status != nil
	})
}
