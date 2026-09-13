package handlers

import (
	"errors"
	"fmt"
	"seanime/internal/api/bangumi"
	"seanime/internal/media"
	"seanime/internal/platforms/bangumi_platform"
	"seanime/internal/platforms/shared_platform"
	"seanime/internal/util/result"
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
	anilistListNovelCache       = result.NewCache[string, *media.ListAnime]()
)

// resolveBangumiSort 将 AniList MediaSort 映射为 Bangumi 搜索 sort 值。
// sort 是数组，取第一个可映射值；有 search 词时优先 sort=match；
// 无可映射值时返回空串，走 Bangumi 服务端默认排序（实测合法）。
// 实测地面真值（2026-09-13）：sort ∈ {match, heat, rank, score} 均可用。
func resolveBangumiSort(sorts []*media.MediaSort, keyword string) string {
	if keyword != "" {
		return "match"
	}
	for _, s := range sorts {
		if s == nil {
			continue
		}
		switch *s {
		case media.MediaSortTrendingDesc:
			return "heat"
		case media.MediaSortPopularityDesc:
			return "rank"
		case media.MediaSortScoreDesc:
			return "score"
		}
	}
	return ""
}

// seasonAirDateRange 将 AniList season/seasonYear 映射为 Bangumi air_date 全日期区间。
// 实测地面真值：air_date 元素必须是操作符+全日期（如 ">2026-01-01"），年月格式会被 400 拒绝。
// 区间端点按契约：WINTER=01-01~03-31、SPRING=04-01~06-30、SUMMER=07-01~09-30、FALL=10-01~12-31。
// 上界取下一节点首日（如 WINTER 上界 "<2026-04-01"，等效覆盖 03-31 全天）；
// FALL 上界取下一年 01-01，避免漏掉 12-31 当天放送。仅 season 无 year 时忽略。
func seasonAirDateRange(season *media.MediaSeason, year *int) []string {
	if year == nil {
		return nil
	}
	if season != nil && season.IsValid() {
		var start, end string
		switch *season {
		case media.MediaSeasonWinter:
			start, end = "01-01", fmt.Sprintf("%d-04-01", *year)
		case media.MediaSeasonSpring:
			start, end = "04-01", fmt.Sprintf("%d-07-01", *year)
		case media.MediaSeasonSummer:
			start, end = "07-01", fmt.Sprintf("%d-10-01", *year)
		case media.MediaSeasonFall:
			start, end = "10-01", fmt.Sprintf("%d-01-01", *year+1)
		default:
			return nil
		}
		return []string{fmt.Sprintf(">%d-%s", *year, start), "<" + end}
	}
	// 仅年份：全年区间（契约规定 ">YYYY-01-01","<YYYY-12-31"）
	return []string{fmt.Sprintf(">%d-01-01", *year), fmt.Sprintf("<%d-12-31", *year)}
}

// ratingFilterFromAverageScore 将 AniList averageScore_greater（0-100）映射为
// Bangumi rating 过滤（0-10 制），N=分数÷10 向下取整；0-10 保留 0。
func ratingFilterFromAverageScore(score *int) []string {
	if score == nil || *score < 0 {
		return nil
	}
	n := *score / 10
	if n > 10 {
		n = 10
	}
	return []string{fmt.Sprintf(">=%d", n)}
}

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

	// 缺省补默认值：两个字段独立判空，nil 时先 new 出对象再赋值（见 03.6-DIAGNOSIS 故障 2）。
	if p.Page == nil {
		p.Page = new(int)
		*p.Page = 1
	}
	if p.PerPage == nil {
		p.PerPage = new(int)
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
	// 已映射：关键词/标签(genres+tags)/成人内容/评分/季度/排序/分页；
	// status、format、countryOfOrigin 无对应过滤条件，忽略（TODO(M4)：前端过滤选项同步裁剪）。
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
	// genres 逐个并入 filter.tag（Bangumi 无 genres 概念，标签是最接近的过滤维度）
	for _, g := range p.Genres {
		if g != nil && *g != "" {
			filter.Tag = append(filter.Tag, *g)
		}
	}
	if isAdult != nil {
		filter.Nsfw = isAdult
	}
	filter.AirDate = seasonAirDateRange(p.Season, p.SeasonYear)
	filter.Rating = ratingFilterFromAverageScore(p.AverageScoreGreater)

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
		Sort:    resolveBangumiSort(p.Sort, keyword),
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

// HandleAnilistListNovel
//
//	@summary returns a list of novels (轻小说) based on the search parameters.
//	@desc Bangumi 锚点下没有独立的轻小说分类：内部用 SearchSubjects type=1（书籍）搜索，
//	@desc 再按条目 platform ∈ {"小说","WEB"} 过滤（platform 为空值的书籍条目会被丢弃）。
//	@desc 分页采用 over-fetch 近似：每次请求 limit=perPage*3、offset=(page-1)*perPage*3，
//	@desc 过滤后不足 perPage 如实返回，超出截断。平台过滤会使 total 偏大（total 为书籍总数），
//	@desc hasNextPage 以「过滤后条目数 >= perPage」近似判断而非按 total 估算，分页边界略偏，属可接受近似。
//	@desc 请求/响应形状与 list-anime 完全一致（复用 media.ListAnime）。
//	@route /api/v1/anilist/list-novel [POST]
//	@returns media.ListAnime
func (h *Handler) HandleAnilistListNovel(c echo.Context) error {

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

	page := 1
	if p.Page != nil && *p.Page > 0 {
		page = *p.Page
	}
	perPage := 20
	if p.PerPage != nil && *p.PerPage > 0 {
		perPage = *p.PerPage
	}
	keyword := ""
	if p.Search != nil {
		keyword = *p.Search
	}
	isAdult := false
	if p.IsAdult != nil {
		isAdult = *p.IsAdult && h.App.Settings.GetAnilist().EnableAdultContent
	}

	cacheKey := fmt.Sprintf("novel:%d:%d:%s:%v:%v:%v:%d:%v:%v:%v", page, perPage, keyword, p.Sort, p.Tags, p.Genres, p.AverageScoreGreater, p.Season, p.SeasonYear, isAdult)
	if cached, ok := anilistListNovelCache.Get(cacheKey); ok {
		return h.RespondWithData(c, cached)
	}

	client := h.App.AnilistPlatformRef.Get().GetBangumiClient()
	if client == nil {
		return h.RespondWithError(c, errors.New("bangumi client not available"))
	}

	filter := bangumi.SearchFilter{Type: []int{1}} // 1=书籍
	for _, t := range p.Tags {
		if t != nil && *t != "" {
			filter.Tag = append(filter.Tag, *t)
		}
	}
	for _, g := range p.Genres {
		if g != nil && *g != "" {
			filter.Tag = append(filter.Tag, *g)
		}
	}
	if p.IsAdult != nil {
		v := *p.IsAdult && h.App.Settings.GetAnilist().EnableAdultContent
		filter.Nsfw = &v
	}
	filter.AirDate = seasonAirDateRange(p.Season, p.SeasonYear)
	filter.Rating = ratingFilterFromAverageScore(p.AverageScoreGreater)

	// 注意：此处**不得**为空关键词再追加 Bangumi tag「轻小说」。
	// 实测 filter.tag 对 type=1（书籍分区）恒返 total=0（漫画/轻小说/小说均如此），
	// 一旦追加，轻小说探索会从「有结果」直接变「全空」。轻小说子集只靠下方
	// platform ∈ {小说, WEB} 过滤（实测有效：books 分区前 100 条中 29 条 platform=小说）。

	// over-fetch：书籍搜索结果中轻小说占比低（平台过滤会大量丢弃），且服务端
	// limit 被钳制在 20（实测请求 50 仅返回 20），单次请求不够——改为多页循环
	// 拉取（最多 5 页），凑够 perPage 条轻小说或拉完即停。
	// 注意 limit/offset 上限是 20，先用 limit=20 探明服务端实际返回数再循环。
	const serverPageLimit = 20
	novelPlatforms := map[string]bool{"小说": true, "WEB": true}
	mediaList := make([]*media.Anime, 0, perPage)
	reachedEnd := false
	for pageIdx := 0; pageIdx < 5 && len(mediaList) < perPage && !reachedEnd; pageIdx++ {
		offset := (page - 1 + pageIdx) * serverPageLimit
		res, err := client.SearchSubjects(c.Request().Context(), bangumi.SearchSubjectsOpts{
			Keyword: keyword,
			Sort:    resolveBangumiSort(p.Sort, keyword),
			Filter:  filter,
			Limit:   serverPageLimit,
			Offset:  offset,
		})
		if err != nil {
			if pageIdx == 0 {
				return h.RespondWithError(c, err)
			}
			break // 后续页失败降级为返回已获取部分
		}
		for i := range res.Data {
			if !novelPlatforms[res.Data[i].Platform] {
				continue
			}
			if a := media.AnimeFromSubject(bangumi.SubjectToMedia(&res.Data[i])); a != nil {
				mediaList = append(mediaList, a)
			}
		}
		if len(res.Data) < serverPageLimit || offset+serverPageLimit >= res.Total {
			reachedEnd = true
		}
	}
	if len(mediaList) > perPage {
		mediaList = mediaList[:perPage]
	}

	// hasNextPage 近似：过滤后条目数 >= perPage 视为还有下一页。
	// 不用原始 total 除系数估算——轻小说占比未知，系数估算不可靠。
	hasNextPage := len(mediaList) >= perPage && !reachedEnd
	total := len(mediaList)
	pi := perPage
	ret := &media.ListAnime{Page: &media.ListAnime_Page{
		Media:    mediaList,
		PageInfo: &media.PageInfo{CurrentPage: &page, PerPage: &pi, Total: &total, HasNextPage: &hasNextPage},
	}}

	anilistListNovelCache.SetT(cacheKey, ret, time.Minute*10)

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

	// 缺省补默认值：两个字段独立判空，nil 时先 new 出对象再赋值（见 03.6-DIAGNOSIS 故障 2）。
	if p.Page == nil {
		p.Page = new(int)
		*p.Page = 1
	}
	if p.PerPage == nil {
		p.PerPage = new(int)
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
