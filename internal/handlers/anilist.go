package handlers

import (
	"errors"
	"fmt"
	"seanime/internal/media"
	"seanime/internal/platforms/bangumi_platform"
	"seanime/internal/platforms/shared_platform"
	"seanime/internal/util/result"
	"seanime/internal/api/bangumi"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// HandleGetAnimeCollection
//
//	@summary returns the user's AniList anime collection.
//	@desc Calling GET will return the cached anime collection.
//	@desc The manga collection is also refreshed in the background, and upon completion, a WebSocket event is sent.
//	@desc Calling POST will refetch both the anime and manga collections.
//	@returns media.AnimeCollection
//	@route /api/v1/anilist/collection [GET,POST]
func (h *Handler) HandleGetAnimeCollection(c echo.Context) error {

	bypassCache := c.Request().Method == "POST"

	if !bypassCache {
		// Get the user's anilist collection
		animeCollection, err := h.App.GetAnimeCollection(false)
		if err != nil {
			return h.RespondWithError(c, err)
		}
		return h.RespondWithData(c, animeCollection)
	}

	animeCollection, err := h.App.RefreshAnimeCollection()
	if err != nil {
		return h.RespondWithError(c, err)
	}

	go func() {
		_, _ = h.App.RefreshMangaCollection()
	}()

	return h.RespondWithData(c, animeCollection)
}

// HandleGetRawAnimeCollection
//
//	@summary returns the user's AniList anime collection without filtering out custom lists.
//	@desc Calling GET will return the cached anime collection.
//	@returns media.AnimeCollection
//	@route /api/v1/anilist/collection/raw [GET,POST]
func (h *Handler) HandleGetRawAnimeCollection(c echo.Context) error {

	bypassCache := c.Request().Method == "POST"

	// Get the user's anilist collection
	animeCollection, err := h.App.GetRawAnimeCollection(bypassCache)
	if err != nil {
		return h.RespondWithError(c, err)
	}

	return h.RespondWithData(c, animeCollection)
}

var tagsCache *media.MediaTagMap

// HandleGetRawAnimeCollectionTags
//
//	@summary returns the AniList tags for the user's raw anime collection.
//	@desc This runs a dedicated AniList tags query used by the lists page filters.
//	@returns media.MediaTagMap
//	@route /api/v1/anilist/collection/raw/tags [GET]
func (h *Handler) HandleGetRawAnimeCollectionTags(c echo.Context) error {
	h.App.OnRefreshAnilistCollectionFuncs.Set("HandleGetRawAnimeCollectionTags", func() {
		tagsCache = nil
	})

	if tagsCache != nil {
		return h.RespondWithData(c, *tagsCache)
	}

	if h.App.GetUser().IsSimulated {
		return h.RespondWithData(c, media.MediaTagMap{})
	}

	// Bangumi 锚点：无用户标签聚合查询，改为分页拉取收藏并聚合各条目的用户标签。
	// 注意：GET 收藏接口的路径参数是 username 句柄（GetMe 解析），不是昵称；
	// App.GetUsername() 返回的是展示昵称，直接传会导致 404。
	client := h.App.AnilistPlatformRef.Get().GetBangumiClient()
	if client == nil {
		return h.RespondWithData(c, media.MediaTagMap{})
	}
	bp, ok := h.App.AnilistPlatformRef.Get().(*bangumi_platform.BangumiPlatform)
	if !ok {
		// 非 Bangumi 平台（offline/simulated）：无标签可聚合
		return h.RespondWithData(c, media.MediaTagMap{})
	}
	userName := bp.Username(c.Request().Context())
	if userName == "" {
		return h.RespondWithData(c, media.MediaTagMap{})
	}

	tags := make(media.MediaTagMap)
	limit := 50
	for offset := 0; ; offset += limit {
		res, err := client.GetUserCollectionsByUser(c.Request().Context(), userName, bangumi.UserCollectionsOpts{Limit: limit, Offset: offset})
		if err != nil {
			return h.RespondWithError(c, err)
		}
		for _, uc := range res.Data {
			if len(uc.Tags) > 0 {
				tags[uc.SubjectID] = uc.Tags
			}
		}
		if len(res.Data) < limit || offset+limit >= res.Total {
			break
		}
	}
	tagsCache = &tags

	return h.RespondWithData(c, tags)
}

// HandleEditAnilistListEntry
//
//	@summary updates the user's list entry on Anilist.
//	@desc This is used to edit an entry on AniList.
//	@desc The "type" field is used to determine if the entry is an anime or manga and refreshes the collection accordingly.
//	@desc The client should refetch collection-dependent queries after this mutation.
//	@returns true
//	@route /api/v1/anilist/list-entry [POST]
func (h *Handler) HandleEditAnilistListEntry(c echo.Context) error {

	type body struct {
		MediaId   *int                   `json:"mediaId"`
		Status    *media.MediaListStatus `json:"status"`
		Score     *int                   `json:"score"`
		Progress  *int                   `json:"progress"`
		StartDate *media.FuzzyDateInput  `json:"startedAt"`
		EndDate   *media.FuzzyDateInput  `json:"completedAt"`
		Type      string                 `json:"type"`
	}

	p := new(body)
	if err := c.Bind(p); err != nil {
		return h.RespondWithError(c, err)
	}

	err := h.App.AnilistPlatformRef.Get().UpdateEntry(
		c.Request().Context(),
		*p.MediaId,
		p.Status,
		p.Score,
		p.Progress,
		p.StartDate,
		p.EndDate,
	)
	if err != nil {
		return h.RespondWithError(c, err)
	}

	switch p.Type {
	case "anime":
		_, _ = h.App.RefreshAnimeCollection()
	case "manga":
		_, _ = h.App.RefreshMangaCollection()
	default:
		_, _ = h.App.RefreshAnimeCollection()
		_, _ = h.App.RefreshMangaCollection()
	}

	return h.RespondWithData(c, true)
}

//----------------------------------------------------------------------------------------------------------------------------------------------------

var (
	detailsCache = result.NewCache[int, *media.AnimeDetails]()
)

// HandleGetAnilistAnimeDetails
//
//	@summary returns more details about an AniList anime entry.
//	@desc This fetches more fields omitted from the base queries.
//	@param id - int - true - "The AniList anime ID"
//	@returns media.AnimeDetails
//	@route /api/v1/anilist/media-details/{id} [GET]
func (h *Handler) HandleGetAnilistAnimeDetails(c echo.Context) error {

	mId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.RespondWithError(c, err)
	}

	if details, ok := detailsCache.Get(mId); ok {
		return h.RespondWithData(c, details)
	}
	details, err := h.App.AnilistPlatformRef.Get().GetAnimeDetails(c.Request().Context(), mId)
	if err != nil {
		return h.RespondWithError(c, err)
	}
	detailsCache.Set(mId, details)

	return h.RespondWithData(c, details)
}

//----------------------------------------------------------------------------------------------------------------------------------------------------

var studioDetailsMap = result.NewMap[int, *media.StudioDetails]()

// HandleGetAnilistStudioDetails
//
//	@summary returns details about a studio.
//	@desc This fetches media produced by the studio.
//	@param id - int - true - "The AniList studio ID"
//	@returns media.StudioDetails
//	@route /api/v1/anilist/studio-details/{id} [GET]
func (h *Handler) HandleGetAnilistStudioDetails(c echo.Context) error {

	mId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.RespondWithError(c, err)
	}

	if details, ok := studioDetailsMap.Get(mId); ok {
		return h.RespondWithData(c, details)
	}
	details, err := h.App.AnilistPlatformRef.Get().GetStudioDetails(c.Request().Context(), mId)
	if err != nil {
		return h.RespondWithError(c, err)
	}

	go func() {
		if details != nil {
			studioDetailsMap.Set(mId, details)
		}
	}()

	return h.RespondWithData(c, details)
}

//----------------------------------------------------------------------------------------------------------------------------------------------------

// HandleDeleteAnilistListEntry
//
//	@summary deletes an entry from the user's AniList list.
//	@desc This is used to delete an entry on AniList.
//	@desc The "type" field is used to determine if the entry is an anime or manga and refreshes the collection accordingly.
//	@desc The client should refetch collection-dependent queries after this mutation.
//	@route /api/v1/anilist/list-entry [DELETE]
//	@returns bool
func (h *Handler) HandleDeleteAnilistListEntry(c echo.Context) error {

	type body struct {
		MediaId *int    `json:"mediaId"`
		Type    *string `json:"type"`
	}

	p := new(body)
	if err := c.Bind(p); err != nil {
		return h.RespondWithError(c, err)
	}

	if p.Type == nil || p.MediaId == nil {
		return h.RespondWithError(c, errors.New("missing parameters"))
	}

	var listEntryID int

	switch *p.Type {
	case "anime":
		// Get the list entry ID
		animeCollection, err := h.App.GetAnimeCollection(false)
		if err != nil {
			return h.RespondWithError(c, err)
		}

		listEntry, found := animeCollection.GetListEntryFromAnimeId(*p.MediaId)
		if !found {
			return h.RespondWithError(c, errors.New("list entry not found"))
		}
		listEntryID = listEntry.ID
	case "manga":
		// Get the list entry ID
		mangaCollection, err := h.App.GetMangaCollection(false)
		if err != nil {
			return h.RespondWithError(c, err)
		}

		listEntry, found := mangaCollection.GetListEntryFromMangaId(*p.MediaId)
		if !found {
			return h.RespondWithError(c, errors.New("list entry not found"))
		}
		listEntryID = listEntry.ID
	}

	// Delete the list entry
	err := h.App.AnilistPlatformRef.Get().DeleteEntry(c.Request().Context(), *p.MediaId, listEntryID)
	if err != nil {
		return h.RespondWithError(c, err)
	}

	switch *p.Type {
	case "anime":
		_, _ = h.App.RefreshAnimeCollection()
	case "manga":
		_, _ = h.App.RefreshMangaCollection()
	}

	return h.RespondWithData(c, true)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var (
	anilistListAnimeCache       = result.NewCache[string, *media.ListAnime]()
	anilistListRecentAnimeCache = result.NewCache[string, *media.ListRecentAnime]() // holds 1 value
)

// HandleAnilistListAnime
//
//	@summary returns a list of anime based on the search parameters.
//	@desc This is used by the "Discover" and "Advanced Search".
//	@route /api/v1/anilist/list-anime [POST]
//	@returns media.ListAnime
func (h *Handler) HandleAnilistListAnime(c echo.Context) error {

	type body struct {
		Page                *int                 `json:"page,omitempty"`
		Search              *string              `json:"search,omitempty"`
		PerPage             *int                 `json:"perPage,omitempty"`
		Sort                []*media.MediaSort   `json:"sort,omitempty"`
		Status              []*media.MediaStatus `json:"status,omitempty"`
		Genres              []*string            `json:"genres,omitempty"`
		Tags                []*string            `json:"tags,omitempty"`
		AverageScoreGreater *int                 `json:"averageScore_greater,omitempty"`
		Season              *media.MediaSeason   `json:"season,omitempty"`
		SeasonYear          *int                 `json:"seasonYear,omitempty"`
		Format              *media.MediaFormat   `json:"format,omitempty"`
		IsAdult             *bool                `json:"isAdult,omitempty"`
		CountryOfOrigin     *string              `json:"countryOfOrigin,omitempty"`
	}

	p := new(body)
	if err := c.Bind(p); err != nil {
		return h.RespondWithError(c, err)
	}

	if p.Page == nil || p.PerPage == nil {
		*p.Page = 1
		*p.PerPage = 20
	}

	var isAdult *bool = nil
	if p.IsAdult != nil {
		isAdult = new(*p.IsAdult && h.App.Settings.GetAnilist().EnableAdultContent)
	}

	cacheKey := media.ListAnimeCacheKey(
		p.Page,
		p.Search,
		p.PerPage,
		p.Sort,
		p.Status,
		p.Genres,
		p.Tags,
		p.AverageScoreGreater,
		p.Season,
		p.SeasonYear,
		p.Format,
		isAdult,
		p.CountryOfOrigin,
	)

	cached, ok := anilistListAnimeCache.Get(cacheKey)
	if ok {
		return h.RespondWithData(c, cached)
	}

	// Bangumi 锚点：原 AniList 复杂过滤搜索改为 SearchSubjects。
	// 可映射：关键词/标签/成人内容过滤/分页；sort、status、genres、averageScore、
	// season、seasonYear、format、countryOfOrigin 无对应过滤条件，忽略（TODO(M4)：前端过滤选项同步裁剪）。
	client := h.App.AnilistPlatformRef.Get().GetBangumiClient()
	if client == nil {
		return h.RespondWithError(c, errors.New("bangumi client not available"))
	}

	filter := bangumi.SearchFilter{Type: []int{2}} // 2=动画
	for _, t := range p.Tags {
		if t != nil && *t != "" {
			filter.Tag = append(filter.Tag, *t)
		}
	}
	if isAdult != nil {
		filter.Nsfw = isAdult
	}

	page := 1
	if p.Page != nil {
		page = *p.Page
	}
	perPage := 20
	if p.PerPage != nil {
		perPage = *p.PerPage
	}
	keyword := ""
	if p.Search != nil {
		keyword = *p.Search
	}

	res, err := client.SearchSubjects(c.Request().Context(), bangumi.SearchSubjectsOpts{
		Keyword: keyword,
		Sort:    "match",
		Filter:  filter,
		Limit:   perPage,
		Offset:  (page - 1) * perPage,
	})
	if err != nil {
		return h.RespondWithError(c, err)
	}

	mediaList := make([]*media.Anime, 0, len(res.Data))
	for i := range res.Data {
		if a := media.AnimeFromSubject(bangumi.SubjectToMedia(&res.Data[i])); a != nil {
			mediaList = append(mediaList, a)
		}
	}
	hasNextPage := res.Offset+len(res.Data) < res.Total
	total := res.Total
	pi := perPage
	ret := &media.ListAnime{Page: &media.ListAnime_Page{
		Media:    mediaList,
		PageInfo: &media.PageInfo{CurrentPage: &page, PerPage: &pi, Total: &total, HasNextPage: &hasNextPage},
	}}

	if ret != nil {
		anilistListAnimeCache.SetT(cacheKey, ret, time.Minute*10)
	}

	return h.RespondWithData(c, ret)
}

// HandleAnilistListRecentAiringAnime
//
//	@summary returns a list of recently aired anime.
//	@desc This is used by the "Schedule" page to display recently aired anime.
//	@route /api/v1/anilist/list-recent-anime [POST]
//	@returns media.ListRecentAnime
func (h *Handler) HandleAnilistListRecentAiringAnime(c echo.Context) error {

	type body struct {
		Page            *int                `json:"page,omitempty"`
		Search          *string             `json:"search,omitempty"`
		PerPage         *int                `json:"perPage,omitempty"`
		AiringAtGreater *int                `json:"airingAt_greater,omitempty"`
		AiringAtLesser  *int                `json:"airingAt_lesser,omitempty"`
		NotYetAired     *bool               `json:"notYetAired,omitempty"`
		Sort            []*media.AiringSort `json:"sort,omitempty"`
	}

	p := new(body)
	if err := c.Bind(p); err != nil {
		return h.RespondWithError(c, err)
	}

	if p.Page == nil || p.PerPage == nil {
		*p.Page = 1
		*p.PerPage = 50
	}

	cacheKey := fmt.Sprintf("%v-%v-%v-%v-%v-%v-%v", p.Page, p.Search, p.PerPage, p.AiringAtGreater, p.AiringAtLesser, p.NotYetAired, p.Sort)

	cached, ok := anilistListRecentAnimeCache.Get(cacheKey)
	if ok {
		return h.RespondWithData(c, cached)
	}

	// Bangumi 锚点降级：无逐集放送时间戳端点（见契约 Wave B GetAnimeAiringSchedule 降级），
	// 返回空列表，Schedule 页留待 M4 补齐。
	ret := &media.ListRecentAnime{}

	anilistListRecentAnimeCache.SetT(cacheKey, ret, time.Hour*1)

	return h.RespondWithData(c, ret)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var anilistMissedSequelsCache = result.NewCache[int, []*media.Anime]()

// HandleAnilistListMissedSequels
//
//	@summary returns a list of sequels not in the user's list.
//	@desc This is used by the "Discover" page to display sequels the user may have missed.
//	@route /api/v1/anilist/list-missed-sequels [GET]
//	@returns []media.Anime
func (h *Handler) HandleAnilistListMissedSequels(c echo.Context) error {

	cached, ok := anilistMissedSequelsCache.Get(1)
	if ok {
		return h.RespondWithData(c, cached)
	}

	// Get complete anime collection
	animeCollection, err := h.App.AnilistPlatformRef.Get().GetAnimeCollectionWithRelations(c.Request().Context())
	if err != nil {
		return h.RespondWithError(c, err)
	}

	// Bangumi 锚点：原 media.ListMissedSequels（AnilistClient 批量查询）改为本地遍历
	// GetAnimeCollectionWithRelations 的关系边——续作节点随集合返回，无需二次请求。
	ret := make([]*media.Anime, 0)
	seen := make(map[int]struct{})
	for _, list := range animeCollection.GetMediaListCollection().GetLists() {
		if list.Status == nil || !(*list.Status == media.MediaListStatusCompleted || *list.Status == media.MediaListStatusRepeating || *list.Status == media.MediaListStatusPaused) || list.Entries == nil {
			continue
		}
		for _, entry := range list.Entries {
			if entry == nil || entry.GetMedia() == nil {
				continue
			}
			for _, edge := range entry.GetMedia().GetRelations().GetEdges() {
				if edge == nil || edge.GetRelationType() == nil || *edge.GetRelationType() != media.MediaRelationSequel {
					continue
				}
				sequel := edge.GetNode()
				if sequel == nil {
					continue
				}
				if _, found := animeCollection.FindAnime(sequel.GetID()); found {
					continue
				}
				if status := sequel.GetStatus(); status == nil || (*status != media.MediaStatusFinished && *status != media.MediaStatusReleasing) {
					continue
				}
				if _, ok := seen[sequel.GetID()]; ok {
					continue
				}
				seen[sequel.GetID()] = struct{}{}
				ret = append(ret, sequel)
			}
		}
	}
	if len(ret) > 10 {
		ret = ret[:10]
	}

	anilistMissedSequelsCache.SetT(1, ret, time.Hour*4)

	return h.RespondWithData(c, ret)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var anilistStatsCache = result.NewCache[int, *media.Stats]()

// HandleGetAniListStats
//
//	@summary returns the anilist stats.
//	@desc This returns the AniList stats for the user.
//	@route /api/v1/anilist/stats [GET]
//	@returns media.Stats
func (h *Handler) HandleGetAniListStats(c echo.Context) error {
	cached, ok := anilistStatsCache.Get(0)
	if ok {
		return h.RespondWithData(c, cached)
	}

	stats, err := h.App.AnilistPlatformRef.Get().GetViewerStats(c.Request().Context())
	if err != nil {
		return h.RespondWithError(c, err)
	}

	ret := media.GetStats(stats)

	anilistStatsCache.SetT(0, ret, time.Hour*1)

	return h.RespondWithData(c, ret)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// HandleGetAnilistCacheLayerStatus
//
//	@summary returns the status of the AniList cache layer.
//	@desc This returns the status of the AniList cache layer.
//	@route /api/v1/anilist/cache-layer/status [GET]
//	@returns bool
func (h *Handler) HandleGetAnilistCacheLayerStatus(c echo.Context) error {
	return h.RespondWithData(c, shared_platform.IsWorking.Load())
}

// HandleToggleAnilistCacheLayerStatus
//
//	@summary toggles the status of the AniList cache layer.
//	@desc This toggles the status of the AniList cache layer.
//	@route /api/v1/anilist/cache-layer/status [POST]
//	@returns bool
func (h *Handler) HandleToggleAnilistCacheLayerStatus(c echo.Context) error {
	shared_platform.IsWorking.Store(!shared_platform.IsWorking.Load())
	return h.RespondWithData(c, shared_platform.IsWorking.Load())
}
