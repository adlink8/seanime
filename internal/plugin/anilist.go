package plugin

import (
	"context"
	"errors"
	"seanime/internal/api/bangumi"
	"seanime/internal/extension"
	"seanime/internal/extension_repo/prompt"
	"seanime/internal/goja/goja_bindings"
	"seanime/internal/library/anime"
	"seanime/internal/media"
	"seanime/internal/platforms/platform"
	gojautil "seanime/internal/util/goja"

	"github.com/dop251/goja"
	"github.com/rs/zerolog"
)

type Anilist struct {
	ctx       *AppContextImpl
	ext       *extension.Extension
	logger    *zerolog.Logger
	scheduler *gojautil.Scheduler
}

// BindAnilist binds the anilist API to the Goja runtime.
// Permissions need to be checked by the caller.
// Permissions needed: anilist
func (a *AppContextImpl) BindAnilist(vm *goja.Runtime, logger *zerolog.Logger, ext *extension.Extension, scheduler *gojautil.Scheduler) {
	al := &Anilist{
		ctx:       a,
		ext:       ext,
		logger:    new(logger.With().Str("id", ext.ID).Logger()),
		scheduler: scheduler,
	}
	anilistObj := getAnilistObj(vm)
	_ = anilistObj.Set("refreshAnimeCollection", al.RefreshAnimeCollection)
	_ = anilistObj.Set("refreshMangaCollection", al.RefreshMangaCollection)
	// Bangumi 锚点：请求提供者概念随 AniList client 移除，保留方法名返回固定值（插件兼容）。
	_ = anilistObj.Set("getRequestProvider", func() string {
		return "official"
	})

	// Bind anilist platform
	anilistPlatformRef, ok := a.anilistPlatformRef.Get()
	if ok {
		_ = anilistObj.Set("updateEntry", func(mediaID int, status *media.MediaListStatus, scoreRaw *int, progress *int, startedAt *media.FuzzyDateInput, completedAt *media.FuzzyDateInput) error {
			return anilistPlatformRef.Get().UpdateEntry(context.Background(), mediaID, status, scoreRaw, progress, startedAt, completedAt)
		})
		_ = anilistObj.Set("updateEntryProgress", func(mediaID int, progress int, totalEpisodes *int) error {
			return anilistPlatformRef.Get().UpdateEntryProgress(context.Background(), mediaID, progress, totalEpisodes)
		})
		_ = anilistObj.Set("updateEntryRepeat", func(mediaID int, repeat int) error {
			return anilistPlatformRef.Get().UpdateEntryRepeat(context.Background(), mediaID, repeat)
		})
		_ = anilistObj.Set("deleteEntry", func(mediaID int, entryId int) error {
			return anilistPlatformRef.Get().DeleteEntry(context.Background(), mediaID, entryId)
		})
		_ = anilistObj.Set("getAnimeCollection", func(bypassCache bool) (*media.AnimeCollection, error) {
			return anilistPlatformRef.Get().GetAnimeCollection(context.Background(), bypassCache)
		})
		_ = anilistObj.Set("getRawAnimeCollection", func(bypassCache bool) (*media.AnimeCollection, error) {
			return anilistPlatformRef.Get().GetRawAnimeCollection(context.Background(), bypassCache)
		})
		_ = anilistObj.Set("getMangaCollection", func(bypassCache bool) (*media.MangaCollection, error) {
			return anilistPlatformRef.Get().GetMangaCollection(context.Background(), bypassCache)
		})
		_ = anilistObj.Set("getRawMangaCollection", func(bypassCache bool) (*media.MangaCollection, error) {
			return anilistPlatformRef.Get().GetRawMangaCollection(context.Background(), bypassCache)
		})
		_ = anilistObj.Set("getAnime", func(mediaID int) (*media.Anime, error) {
			return anilistPlatformRef.Get().GetAnime(context.Background(), mediaID)
		})
		_ = anilistObj.Set("getManga", func(mediaID int) (*media.Manga, error) {
			return anilistPlatformRef.Get().GetManga(context.Background(), mediaID)
		})
		_ = anilistObj.Set("getAnimeDetails", func(mediaID int) (*media.AnimeDetails, error) {
			return anilistPlatformRef.Get().GetAnimeDetails(context.Background(), mediaID)
		})
		_ = anilistObj.Set("getMangaDetails", func(mediaID int) (*media.MangaDetails, error) {
			return anilistPlatformRef.Get().GetMangaDetails(context.Background(), mediaID)
		})
		_ = anilistObj.Set("getAnimeCollectionWithRelations", func() (*media.AnimeCollectionWithRelations, error) {
			return anilistPlatformRef.Get().GetAnimeCollectionWithRelations(context.Background())
		})
		_ = anilistObj.Set("addMediaToCollection", func(mIds []int) error {
			return anilistPlatformRef.Get().AddMediaToCollection(context.Background(), mIds)
		})
		_ = anilistObj.Set("getStudioDetails", func(studioID int) (*media.StudioDetails, error) {
			return anilistPlatformRef.Get().GetStudioDetails(context.Background(), studioID)
		})
		_ = anilistObj.Set("listAnime", func(page *int, search *string, perPage *int, sort []*media.MediaSort, status []*media.MediaStatus, genres []*string, tags []*string, averageScoreGreater *int, season *media.MediaSeason, seasonYear *int, format *media.MediaFormat, isAdult *bool) (*media.ListAnime, error) {
			// Bangumi 锚点：仅关键词搜索可映射，其余过滤参数忽略（TODO(M4)）。
			return searchAnimeViaBangumi(anilistPlatformRef.Get(), page, search, perPage)
		})
		_ = anilistObj.Set("listManga", func(page *int, search *string, perPage *int, sort []*media.MediaSort, status []*media.MediaStatus, genres []*string, tags []*string, averageScoreGreater *int, startDateGreater *string, startDateLesser *string, format *media.MediaFormat, countryOfOrigin *string, isAdult *bool) (*media.ListManga, error) {
			// Bangumi 锚点：仅关键词搜索可映射，其余过滤参数忽略（TODO(M4)）。
			return searchMangaViaBangumi(anilistPlatformRef.Get(), page, search, perPage)
		})
		_ = anilistObj.Set("listRecentAnime", func(page *int, perPage *int, airingAtGreater *int, airingAtLesser *int, notYetAired *bool) (*media.ListRecentAnime, error) {
			// Bangumi 锚点降级：无逐集放送时间戳端点，返回空列表（M4 补齐）。
			return &media.ListRecentAnime{}, nil
		})
		_ = anilistObj.Set("clearCache", func() {
			anilistPlatformRef.Get().ClearCache()
			anime.ClearEpisodeCollectionCache()
			anime.ClearMissingEpisodesCache()
			anime.ClearScheduleCache()
		})
		_ = anilistObj.Set("customQuery", func(body map[string]interface{}, token string) (interface{}, error) {
			// Bangumi 锚点：无自定义 GraphQL 查询入口，保留方法名并返回错误（插件兼容）。
			return nil, errors.New("customQuery is not supported: the AniList GraphQL endpoint has been replaced by Bangumi")
		})

	}

	_ = vm.Set("$anilist", anilistObj)
}

// BindAnilistCustomClient binds runtime AniList client swap APIs to the Goja runtime.
// Permissions need to be checked by the caller.
// Permissions needed: custom-client
func (a *AppContextImpl) BindAnilistCustomClient(vm *goja.Runtime, logger *zerolog.Logger, ext *extension.Extension, scheduler *gojautil.Scheduler) {
	al := &Anilist{
		ctx:       a,
		ext:       ext,
		logger:    new(logger.With().Str("id", ext.ID).Logger()),
		scheduler: scheduler,
	}
	anilistObj := getAnilistObj(vm)
	// Bangumi 锚点：请求提供者概念随 AniList client 移除，保留方法名返回固定值（插件兼容）。
	_ = anilistObj.Set("getRequestProvider", func() string {
		return "official"
	})
	_ = anilistObj.Set("useOfficialApi", func() goja.Value {
		return al.runAction(vm, func() error {
			if err := al.ctx.ask(al.ext, prompt.Options{
				Kind:       "custom-client",
				Action:     "restore official AniList client",
				Resource:   "AniList client",
				Message:    "Allow \"" + al.ext.Name + "\" to restore the official AniList client?",
				AllowLabel: "Restore",
			}); err != nil {
				return err
			}

			// Bangumi 锚点：AniList client 运行时切换已停用，保留方法名并返回错误（插件兼容）。
			return errors.New("useOfficialApi is not supported: the AniList client has been replaced by Bangumi")
		})
	})
	_ = anilistObj.Set("useCustomApi", func(value goja.Value) goja.Value {
		return al.runAction(vm, func() error {
			config, err := readCustomClientConfig(vm, value)
			if err != nil {
				return err
			}
			if err := al.ctx.ask(al.ext, customClientPromptOptions(al.ext, config)); err != nil {
				return err
			}

			// Bangumi 锚点：AniList client 运行时切换已停用，保留方法名并返回错误（插件兼容）。
			_ = config
			return errors.New("useCustomApi is not supported: the AniList client has been replaced by Bangumi")
		})
	})

	_ = vm.Set("$anilist", anilistObj)
}

func getAnilistObj(vm *goja.Runtime) *goja.Object {
	value := vm.Get("$anilist")
	if value != nil && !goja.IsUndefined(value) && !goja.IsNull(value) {
		return value.ToObject(vm)
	}

	obj := vm.NewObject()
	_ = vm.Set("$anilist", obj)
	return obj
}

// CustomClientConfig AniList 自定义端点配置。
// Bangumi 锚点下仅保留结构以兼容插件 API 形状，运行时切换已停用。
type CustomClientConfig struct {
	Name          string
	Endpoint      string
	Token         string
	Authenticated bool
	Headers       map[string]string
}

func customClientPromptOptions(ext *extension.Extension, config CustomClientConfig) prompt.Options {
	name := config.Name
	if name == "" {
		name = "custom"
	}

	details := []string{"Endpoint: " + config.Endpoint}
	if config.Authenticated || config.Token != "" || len(config.Headers) > 0 {
		details = append(details, "May authenticate requests")
	}

	return prompt.Options{
		Kind:       "custom-client",
		Action:     "switch AniList client to \"" + name + "\"",
		Resource:   "AniList client",
		Message:    "Allow \"" + ext.Name + "\" to switch Seanime's AniList client to \"" + name + "\"?",
		Details:    details,
		AllowLabel: "Switch",
	}
}

func readCustomClientConfig(vm *goja.Runtime, value goja.Value) (CustomClientConfig, error) {
	if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
		return CustomClientConfig{}, errors.New("anilist custom client options are required")
	}

	obj := value.ToObject(vm)
	config := CustomClientConfig{
		Name:          readString(obj, "name"),
		Endpoint:      readString(obj, "endpoint"),
		Token:         readString(obj, "token"),
		Authenticated: readBool(obj, "authenticated"),
		Headers:       readStringMap(vm, obj.Get("headers")),
	}

	return config, nil
}

func readString(obj *goja.Object, key string) string {
	value := obj.Get(key)
	if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
		return ""
	}

	return value.String()
}

func readBool(obj *goja.Object, key string) bool {
	value := obj.Get(key)
	if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
		return false
	}

	return value.ToBoolean()
}

func readStringMap(vm *goja.Runtime, value goja.Value) map[string]string {
	if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
		return nil
	}

	obj := value.ToObject(vm)
	ret := make(map[string]string)
	for _, key := range obj.Keys() {
		ret[key] = readString(obj, key)
	}

	return ret
}

// searchAnimeViaBangumi 通过 Bangumi 搜索接口实现插件 listAnime。
func searchAnimeViaBangumi(p platform.Platform, page *int, search *string, perPage *int) (*media.ListAnime, error) {
	client := p.GetBangumiClient()
	if client == nil {
		return nil, errors.New("bangumi client not available")
	}
	pageN, perPageN := 1, 20
	if page != nil {
		pageN = *page
	}
	if perPage != nil {
		perPageN = *perPage
	}
	keyword := ""
	if search != nil {
		keyword = *search
	}
	res, err := client.SearchSubjects(context.Background(), bangumi.SearchSubjectsOpts{
		Keyword: keyword,
		Sort:    "match",
		Filter:  bangumi.SearchFilter{Type: []int{2}},
		Limit:   perPageN,
		Offset:  (pageN - 1) * perPageN,
	})
	if err != nil {
		return nil, err
	}
	mediaList := make([]*media.Anime, 0, len(res.Data))
	for i := range res.Data {
		if an := media.AnimeFromSubject(bangumi.SubjectToMedia(&res.Data[i])); an != nil {
			mediaList = append(mediaList, an)
		}
	}
	hasNextPage := res.Offset+len(res.Data) < res.Total
	total := res.Total
	pi := perPageN
	return &media.ListAnime{Page: &media.ListAnime_Page{
		Media:    mediaList,
		PageInfo: &media.PageInfo{CurrentPage: &pageN, PerPage: &pi, Total: &total, HasNextPage: &hasNextPage},
	}}, nil
}

// searchMangaViaBangumi 通过 Bangumi 搜索接口实现插件 listManga。
func searchMangaViaBangumi(p platform.Platform, page *int, search *string, perPage *int) (*media.ListManga, error) {
	client := p.GetBangumiClient()
	if client == nil {
		return nil, errors.New("bangumi client not available")
	}
	pageN, perPageN := 1, 20
	if page != nil {
		pageN = *page
	}
	if perPage != nil {
		perPageN = *perPage
	}
	keyword := ""
	if search != nil {
		keyword = *search
	}
	res, err := client.SearchSubjects(context.Background(), bangumi.SearchSubjectsOpts{
		Keyword: keyword,
		Sort:    "match",
		Filter:  bangumi.SearchFilter{Type: []int{1}},
		Limit:   perPageN,
		Offset:  (pageN - 1) * perPageN,
	})
	if err != nil {
		return nil, err
	}
	mediaList := make([]*media.Manga, 0, len(res.Data))
	for i := range res.Data {
		if m := media.MangaFromSubject(bangumi.SubjectToMedia(&res.Data[i])); m != nil {
			mediaList = append(mediaList, m)
		}
	}
	hasNextPage := res.Offset+len(res.Data) < res.Total
	total := res.Total
	pi := perPageN
	return &media.ListManga{Page: &media.ListManga_Page{
		Media:    mediaList,
		PageInfo: &media.PageInfo{CurrentPage: &pageN, PerPage: &pi, Total: &total, HasNextPage: &hasNextPage},
	}}, nil
}

func (a *Anilist) RefreshAnimeCollection() {
	a.logger.Trace().Msg("plugin: Refreshing anime collection")
	onRefreshAnilistAnimeCollection, ok := a.ctx.onRefreshAnilistAnimeCollection.Get()
	if !ok {
		return
	}

	onRefreshAnilistAnimeCollection()
}

func (a *Anilist) RefreshMangaCollection() {
	a.logger.Trace().Msg("plugin: Refreshing manga collection")
	onRefreshAnilistMangaCollection, ok := a.ctx.onRefreshAnilistMangaCollection.Get()
	if !ok {
		return
	}

	onRefreshAnilistMangaCollection()
}

func (a *Anilist) runAction(vm *goja.Runtime, run func() error) goja.Value {
	promise, resolve, reject := vm.NewPromise()

	go func() {
		err := run()

		a.scheduler.ScheduleAsync(func() error {
			if err != nil {
				reject(goja_bindings.NewErrorString(vm, err.Error()))
				return nil
			}

			resolve(goja.Undefined())
			return nil
		})
	}()

	return vm.ToValue(promise)
}
