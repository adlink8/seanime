package handlers

import (
	"context"
	"errors"
	"fmt"
	"math"
	"seanime/internal/api/bangumi"
	"seanime/internal/media"
	"seanime/internal/platforms/bangumi_platform"
	"seanime/internal/platforms/shared_platform"
	"seanime/internal/util/result"
	"sort"
	"strconv"
	"strings"
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
	return mapAnilistSortToBangumiSort(sorts)
}

// mapAnilistSortToBangumiSort 仅做值映射，不掺入「有关键词就 match」的服务端语义。
// legacy 通路需要它：legacy 无 sort 参数，本地排序要按用户真实选择（score/rank）而非 match。
func mapAnilistSortToBangumiSort(sorts []*media.MediaSort) string {
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

// useLegacySearch 关键词双通路判据（契约 §2 D1）：关键词 trim 后非空 → legacy 通路。
//
// 依据 03.6b-CONTRACT.md §0.1 实测地面真值：v0 `POST /v0/search/subjects` 对 CJK
// 关键词恒 `total=0`（进击的巨人/火影/海贼王全覆盖），只有 legacy
// `GET /search/subject/{kw}` 能返回中文标题结果。故非空关键词必须改走 legacy。
func useLegacySearch(keyword string) bool {
	return strings.TrimSpace(keyword) != ""
}

// serverPageLimit v0 分页硬钳（契约 §0.2 / D8）：服务端 `limit` 上限为 20，
// 请求更大值会被静默钳到 20；若 offset 仍按前端 perPage(48) 步进就会漏条目。
const serverPageLimit = 20

// legacyBookUpstreamStep 书籍分区（type=1）legacy 检索的分页步长
// （契约 03.7 §1 A3：上游 start 从 (page-1)*20 起，步长 20）。
const legacyBookUpstreamStep = serverPageLimit

// legacyBookMaxUpstreamCalls 单次前端请求允许的上游检索次数上限（契约 03.7 §1 A3：
// 上限 10 次 ≈ 200 条上游条目，超出即停并把 hasNextPage 置真，避免深页打爆上游）。
const legacyBookMaxUpstreamCalls = 10

// clampServerPageLimit 将请求的 perPage 收敛到服务端硬上限内。
// 契约 D8：perPage 固定 20。
func clampServerPageLimit(perPage int) int {
	if perPage <= 0 || perPage > serverPageLimit {
		return serverPageLimit
	}
	return perPage
}

// isNovelPlatform 判定 platform 是否属于「轻小说」（契约 03.7 §1 A2 / D2）。
// 比较去首尾空白且大小写不敏感：上游可能返回 `Web` / `web`（契约 §0.1 与
// adapter.go 的 `Web` 写法），而规范值为 `WEB`。
// 仅接受 {小说, WEB} 两个字面值——`轻小说` 不是 Bangumi 的 platform 取值，
// 故不额外收编（避免把漫画分区的条目误判成小说）。
func isNovelPlatform(platform string) bool {
	switch strings.ToLower(strings.TrimSpace(platform)) {
	case "小说", "web":
		return true
	}
	return false
}

// keepPlatformFor 判定某条目是否属于目标通路（契约 03.7 §1 A2 / D2 双向分流）。
//
//	novel=true （轻小说 tab）：保留 platform ∈ {小说, WEB}
//	novel=false（漫画 tab）  ：保留其余（含空 platform）
//
// D4：platform 缺失（空串 / 详情查询失败）归入**漫画**、不归入轻小说——
// 宁可在漫画 tab 多留，也不要在轻小说 tab 误收非小说。
// 故 novel=false 分支对空串返回 true（保留），novel=true 分支对空串返回 false（排除）。
func keepPlatformFor(novel bool, platform string) bool {
	isNovel := isNovelPlatform(platform)
	if novel {
		return isNovel
	}
	return !isNovel
}

// legacySubjectsToAnime 将 legacy 检索条目映射为 media.Anime。
func legacySubjectsToAnime(list []bangumi.LegacySubject) []*media.Anime {
	out := make([]*media.Anime, 0, len(list))
	for i := range list {
		if a := media.AnimeFromSubject(bangumi.LegacySubjectToMedia(&list[i])); a != nil {
			out = append(out, a)
		}
	}
	return out
}

// legacySubjectIDs 提取 legacy 条目的 id 列表（供 A1 批量补查 platform）。
func legacySubjectIDs(list []bangumi.LegacySubject) []int {
	ids := make([]int, 0, len(list))
	for i := range list {
		ids = append(ids, list[i].ID)
	}
	return ids
}

// legacySubjectsToManga 将 legacy 检索条目映射为 media.Manga。
//
// platforms 为 A1 补查到的 id → platform（契约 03.7 §1 A1/A4）：legacy 条目本身
// **没有 platform 字段**（契约 §0.1），这里只把补查结果**写回 media.Subject** 再交给
// 既有适配函数，条目其余字段仍全部来自 legacy 映射（A4）。
// 未解析到（空串）时 Format 落到 BOOK 兜底——该条目按 D4 仍留在漫画通路。
func legacySubjectsToManga(list []bangumi.LegacySubject, platforms map[int]string) []*media.Manga {
	out := make([]*media.Manga, 0, len(list))
	for i := range list {
		m := bangumi.LegacySubjectToMedia(&list[i])
		if m != nil {
			m.Platform = platforms[list[i].ID] // A4：写回 platform，驱动 Format 推导
		}
		if mm := media.MangaFromSubject(m); mm != nil {
			out = append(out, mm)
		}
	}
	return out
}

// legacySubjectsToListAnime 将 legacy 书籍条目映射为「书籍语义」的 media.Anime，
// 供 list-novel 响应使用（响应类型仍是 media.ListAnime，JSON 字段名不变，只改字段值）。
//
// 复用 MangaSubjectAsListAnime 保证 Type=MANGA + Format ∈ {NOVEL, BOOK}（契约 §4 / D9）。
// platforms 为 A1 补查到的 id → platform（契约 03.7 §1 A1/A4）：
// legacy 条目本身**没有 platform 字段**（契约 §0.1），这里只把补查到的 platform
// **写回 media.Subject** 再交给既有适配函数，条目其余字段仍全部来自 legacy 映射
// （A4：不改成用 v0 Subject 数据替换条目内容）。
// 未解析到（空串）时 Format 落到 BOOK 兜底——按 D4 该条目本就不会出现在 novel 通路。
func legacySubjectsToListAnime(list []bangumi.LegacySubject, platforms map[int]string) []*media.Anime {
	out := make([]*media.Anime, 0, len(list))
	for i := range list {
		m := bangumi.LegacySubjectToMedia(&list[i])
		if m != nil {
			m.Platform = platforms[list[i].ID] // A4：写回 platform，驱动 Format 推导
		}
		if a := media.MangaSubjectAsListAnime(m); a != nil {
			out = append(out, a)
		}
	}
	return out
}

// animeMetaTagsFromFormat 将动画 format 映射为 Bangumi `filter.meta_tags`（契约 §1 D4）。
//
// 实测地面真值（03.6b-CONTRACT.md §0.2）：meta_tags 仅 type=2 生效，且仅
// `TV` / `WEB` / `OVA` 三个值有结果；`Movie`/`剧场版`/`原创`/`漫画改` 全 0。
// 故 `MOVIE` / `TV_SHORT` / `SPECIAL` 不再映射（前端已按 D4 移除这些选项），
// 返回 nil 表示「不追加 meta_tags 过滤」——而不是发出一个恒 0 的无效值。
func animeMetaTagsFromFormat(format *media.MediaFormat) []string {
	if format == nil {
		return nil
	}
	switch *format {
	case media.MediaFormatTv:
		return []string{"TV"}
	case media.MediaFormatOna:
		return []string{"WEB"} // AniList 的 ONA ↔ Bangumi 的 WEB
	case media.MediaFormatOva:
		return []string{"OVA"}
	}
	return nil
}

// legacySubjectScore 取 legacy 条目的评分（0–10）；无 rating 时按 0 处理。
func legacySubjectScore(s bangumi.LegacySubject) float64 {
	if s.Rating == nil {
		return 0
	}
	return s.Rating.Score
}

// legacySubjectRank 取 legacy 条目的全局排名；rank<=0 表示上游未给出排名，
// 排序时置于末位（避免「无排名」被当成「第 0 名」抢到最前）。
func legacySubjectRank(s bangumi.LegacySubject) int {
	if s.Rank <= 0 {
		return math.MaxInt
	}
	return s.Rank
}

// legacyScoreMeetsMinimum 本地评分下限过滤（契约 §2 通路 1）。
// 请求参数 averageScore_greater 为 0–100 刻度（D6），条目 rating.score 为 0–10；
// 比较时把条目分数换算成「十分位整数」再比，与 media.AnimeFromSubject 的
// meanScore = int(score*10+0.5) 换算保持一致，避免浮点边界抖动。
func legacyScoreMeetsMinimum(score float64, averageScoreGreater *int) bool {
	if averageScoreGreater == nil {
		return true
	}
	return int(score*10+0.5) >= *averageScoreGreater
}

// legacyMonthsInSeason 季度 → 月份区间，与 seasonAirDateRange（v0 通路）的
// 日历区间保持一致：WINTER 01–03、SPRING 04–06、SUMMER 07–09、FALL 10–12。
func legacyMonthsInSeason(s media.MediaSeason) (int, int) {
	switch s {
	case media.MediaSeasonWinter:
		return 1, 3
	case media.MediaSeasonSpring:
		return 4, 6
	case media.MediaSeasonSummer:
		return 7, 9
	case media.MediaSeasonFall:
		return 10, 12
	}
	return 0, 0
}

// legacyAirDateMatches 本地年份/季度过滤（契约 §2 通路 1）。
// 无 seasonYear 时不过滤（无谓的丢弃比漏条目更糟）；条目日期缺失或不可解析时，
// 在设置了年份筛选的前提下保守丢弃（无法证明命中）。
func legacyAirDateMatches(airDate string, season *media.MediaSeason, seasonYear *int) bool {
	if seasonYear == nil {
		return true
	}
	if airDate == "" {
		return false
	}
	d, err := time.Parse("2006-01-02", airDate)
	if err != nil || d.Year() != *seasonYear {
		return false
	}
	if season == nil || !season.IsValid() {
		return true
	}
	lo, hi := legacyMonthsInSeason(*season)
	if lo == 0 {
		return true
	}
	return int(d.Month()) >= lo && int(d.Month()) <= hi
}

// filterLegacySubjects 对 legacy 结果集做本地过滤（契约 §2 通路 1）。
//
// 为什么需要本地过滤：legacy 端点不支持任何 filter（契约 §0.1），
// 但条目自带 rating.score（0–10）与 air_date，故「评分下限」「年份/季度」可本地等价实现。
func filterLegacySubjects(list []bangumi.LegacySubject, averageScoreGreater *int, season *media.MediaSeason, seasonYear *int) []bangumi.LegacySubject {
	out := make([]bangumi.LegacySubject, 0, len(list))
	for _, s := range list {
		if !legacyScoreMeetsMinimum(legacySubjectScore(s), averageScoreGreater) {
			continue
		}
		if !legacyAirDateMatches(s.AirDate, season, seasonYear) {
			continue
		}
		out = append(out, s)
	}
	return out
}

// sortLegacySubjects 对 legacy 结果集做本地排序（契约 §2 通路 1）。
//   - score → rating.score 降序
//   - rank  → rank 升序（rank 小者优，未排名者置末）
//   - heat / match / 无可映射值 → 保持上游返回顺序（legacy 默认按相关度）
//
// 用 SliceStable 保证同分条目维持上游相对顺序，结果可复现。
func sortLegacySubjects(list []bangumi.LegacySubject, sorts []*media.MediaSort) {
	switch mapAnilistSortToBangumiSort(sorts) {
	case "score":
		sort.SliceStable(list, func(i, j int) bool {
			return legacySubjectScore(list[i]) > legacySubjectScore(list[j])
		})
	case "rank":
		sort.SliceStable(list, func(i, j int) bool {
			return legacySubjectRank(list[i]) < legacySubjectRank(list[j])
		})
	default:
		// heat / match / 无 → 不排序，保持上游相关度顺序
	}
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
		Page                *int               `json:"page,omitempty"`
		Search              *string            `json:"search,omitempty"`
		PerPage             *int               `json:"perPage,omitempty"`
		Sort                []*media.MediaSort `json:"sort,omitempty"`
		Genres              []*string          `json:"genres,omitempty"`
		Tags                []*string          `json:"tags,omitempty"`
		AverageScoreGreater *int               `json:"averageScore_greater,omitempty"`
		Season              *media.MediaSeason `json:"season,omitempty"`
		SeasonYear          *int               `json:"seasonYear,omitempty"`
		Format              *media.MediaFormat `json:"format,omitempty"`
		IsAdult             *bool              `json:"isAdult,omitempty"`
		// status / countryOfOrigin 已按契约 D2/D3 删除：
		// 上游 Bangumi 无 status 与 country 维度（platform 过滤实测无效，见 §0.2），
		// 此前「接收却不消费」（仅进 cacheKey）属于静默忽略，现从接收端一并移除。
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
		p.Genres,
		p.Tags,
		p.AverageScoreGreater,
		p.Season,
		p.SeasonYear,
		p.Format,
		isAdult,
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
	// genres 经 AniList→Bangumi tag 词表映射后并入 filter.tag（03.9a）。
	// Bangumi 无 genres 概念，且英文 genre 直接作 tag 仅欧美条目偶然携带、命中极低，
	// 故先过 mapGenresToBangumiTags 翻译为中文 tag（未命中词表的 genre 丢弃，绝不透传英文值）。
	genreStrs := make([]string, 0, len(p.Genres))
	for _, g := range p.Genres {
		if g != nil && *g != "" {
			genreStrs = append(genreStrs, *g)
		}
	}
	for _, tag := range mapGenresToBangumiTags(genreStrs) {
		filter.Tag = append(filter.Tag, tag)
	}
	// 动画 format → filter.meta_tags（契约 D4：仅 TV / WEB / OVA 三值）
	filter.MetaTags = animeMetaTagsFromFormat(p.Format)
	if isAdult != nil {
		filter.Nsfw = isAdult
	}
	filter.AirDate = seasonAirDateRange(p.Season, p.SeasonYear)
	filter.Rating = ratingFilterFromAverageScore(p.AverageScoreGreater)

	page := 1
	if p.Page != nil {
		page = *p.Page
	}
	// 契约 D8：v0 `limit` 硬钳 20，offset 必须按 20 步进。
	// 前端曾传 perPage=48：服务端把 limit 钳到 20、offset 仍按 48 递增 → 中间 28 条漏掉。
	perPage := serverPageLimit
	if p.PerPage != nil {
		perPage = clampServerPageLimit(*p.PerPage)
	}
	keyword := ""
	if p.Search != nil {
		keyword = *p.Search
	}

	// 契约 §2 D1 双通路：关键词非空 → legacy（v0 对 CJK 关键词恒 total=0，见 §0.1）。
	//
	// ⚠ 明确降级（禁止静默忽略，契约 §2 通路 1 末条）：
	// legacy 条目**不含 tag、也不含 platform/tags 字段**（契约 §0.1 实测地面真值），
	// 故 `tags` 与 `format`(meta_tags) 这两项筛选在 legacy 通路**无法本地过滤，降级为忽略**。
	// 这是能力缺口而非疏漏：上游 legacy 端点既不支持 filter 参数，返回体里也没有可判据的字段。
	// 若后续需要，只能改为「v0 空关键词 + filter」通路另行实现，本阶段不做。
	//
	// 契约 §2 分页补偿：legacy `max_results` 非严格保证（请求 20 可能只回 17），
	// 此处选择「照实返回并在响应中如实反映条数」（不 over-fetch），
	// hasNextPage 按上游实际返回条数（过滤前）计算，保证翻页游标不因本地过滤而错位。
	if useLegacySearch(keyword) {
		start := (page - 1) * perPage
		res, err := client.SearchSubjectsLegacy(
			c.Request().Context(),
			strings.TrimSpace(keyword),
			bangumi.SubjectAnime, // 动画分区
			start,
			perPage,
		)
		if err != nil {
			return h.RespondWithError(c, err)
		}

		// 本地过滤（评分下限按 rating.score 0–10；年份/季度按 air_date）+ 本地排序（§2 通路 1）
		filtered := filterLegacySubjects(res.List, p.AverageScoreGreater, p.Season, p.SeasonYear)
		sortLegacySubjects(filtered, p.Sort)
		mediaList := legacySubjectsToAnime(filtered)

		hasNextPage := start+len(res.List) < res.Results
		total := res.Results
		pi := perPage
		ret := &media.ListAnime{Page: &media.ListAnime_Page{
			Media:    mediaList,
			PageInfo: &media.PageInfo{CurrentPage: &page, PerPage: &pi, Total: &total, HasNextPage: &hasNextPage},
		}}
		anilistListAnimeCache.SetT(cacheKey, ret, time.Minute*10)
		return h.RespondWithData(c, ret)
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

// legacyBookFilter 书籍分区（type=1）legacy 通路的 platform 分流参数（契约 03.7 §1 A2）。
type legacyBookFilter struct {
	// novel=true → 保留 platform ∈ {小说, WEB}（轻小说 tab）；
	// novel=false → 保留其余（漫画 tab，含空 platform，见 D4）。
	novel bool
	// 本地过滤 / 本地排序沿用契约 §2 通路 1 的既有实现（legacy 无 filter / sort 参数）
	averageScoreGreater *int
	season              *media.MediaSeason
	seasonYear          *int
	sorts               []*media.MediaSort
}

// legacyBookPage legacy 书籍检索 + platform 分流后的分页结果。
type legacyBookPage struct {
	subjects    []bangumi.LegacySubject // 本页窗口内的条目（已 platform 分流、已本地排序）
	platforms   map[int]string          // A1 补查结果（id → platform），供 A4 写回 media.Subject
	hasNextPage bool
	total       int
	perPage     int // 实际生效的每页条数（供调用方回填 PageInfo.PerPage，响应形状不变）
}

// fetchLegacyBookPage 执行「legacy 检索 → 补查 platform → 双向分流 → over-fetch 补足」，
// 返回第 page 页。契约 03.7 §1 A1/A2/A3 落地。
//
// A3 算法：上游 start 从 (page-1)*20 起、步长 20 逐页取；每页取回后补查 platform → 分流 → 累积；
// 停止条件（任一）：累积数 >= page*perPage / 上游耗尽 / 上游请求次数达上限 10。
//
// 坐标说明（A3 步骤 4 的窗口换算）：累积序列从上游 start=(page-1)*20 开始，
// 而 perPage 被钳到 20（== A3 的步长），故契约里的全局窗口
// [(page-1)*perPage, page*perPage) 换算到该累积序列即 [0, perPage)；
// 同理步骤 5 的 `累积数 > page*perPage` 在此坐标下即 `累积数 > perPage`。
func fetchLegacyBookPage(ctx context.Context, client *bangumi.Client, keyword string, page, perPage int, f legacyBookFilter) (*legacyBookPage, error) {
	if page < 1 {
		// 非法页码（前端理论上不发送）：按第 1 页处理，避免负 start 触发上游异常语义
		page = 1
	}
	// A3 的步长恒为 20，窗口/目标条数也只在此口径下自洽（契约 D8：perPage 固定 20）。
	perPage = clampServerPageLimit(perPage)
	firstStart := (page - 1) * legacyBookUpstreamStep
	target := page * perPage // A3 停止条件 1

	var (
		accumulated = make([]bangumi.LegacySubject, 0, perPage)
		platforms   = make(map[int]string)
		total       = 0
		exhausted   = false
	)

	for upstreamCalls := 0; len(accumulated) < target && !exhausted && upstreamCalls < legacyBookMaxUpstreamCalls; upstreamCalls++ {
		start := firstStart + upstreamCalls*legacyBookUpstreamStep
		res, err := client.SearchSubjectsLegacy(ctx, keyword, bangumi.SubjectBook, start, legacyBookUpstreamStep)
		if err != nil {
			if upstreamCalls == 0 {
				return nil, err // 首页失败：无任何内容可返回，如实报错
			}
			// 后续页失败：降级为返回已累积部分（不静默忽略——此处即降级点）。
			// 上游为 legacy 端点且 `max_results` 本非严格保证（契约 §0.1），
			// 深层补页失败时宁可少给条目也不让整次检索失败。
			break
		}
		total = res.Results

		// 本地过滤（评分下限 + 年份/季度，契约 §2 通路 1）。
		filtered := filterLegacySubjects(res.List, f.averageScoreGreater, f.season, f.seasonYear)

		// A1：批量解析 platform（失败 → 空串，按 D4 归漫画，不中断整次检索）。
		resolved := client.GetSubjectPlatforms(ctx, legacySubjectIDs(filtered))
		for id, p := range resolved {
			platforms[id] = p
		}

		// A2：按 platform 双向分流后累积。
		for i := range filtered {
			if keepPlatformFor(f.novel, resolved[filtered[i].ID]) {
				accumulated = append(accumulated, filtered[i])
			}
		}

		// A3 停止条件 2：上游耗尽（本页不足一页，或已越过命中总数）。
		if len(res.List) < legacyBookUpstreamStep || start+legacyBookUpstreamStep >= res.Results {
			exhausted = true
		}
	}

	// A3 步骤 4：本页窗口（换算见函数头注释）。
	windowEnd := perPage
	if windowEnd > len(accumulated) {
		windowEnd = len(accumulated)
	}
	window := accumulated[:windowEnd]

	// 本地排序只作用于本页窗口（legacy 无 sort 参数，契约 §2 通路 1）。
	sortLegacySubjects(window, f.sorts)

	// A3 步骤 5：累积数超出本页窗口 → 还有下一页；或上游尚未耗尽。
	hasNextPage := len(accumulated) > perPage || !exhausted

	return &legacyBookPage{
		subjects:    window,
		platforms:   platforms,
		hasNextPage: hasNextPage,
		total:       total,
		perPage:     perPage,
	}, nil
}

// HandleAnilistListNovel
//
//	@summary returns a list of novels (轻小说) based on the search parameters.
//	@desc Bangumi 锚点下没有独立的轻小说分类，内部走契约 03.6b §2 双通路（type=1 书籍分区）：
//	@desc 关键词非空 → legacy 检索（`GET /search/subject/{kw}?type=1`，CJK 唯一可行通路），
//	@desc   对命中 id 补查 v0 详情拿 platform（03.7 §1 A1），再按 platform ∈ {"小说","WEB"} 分流（D2），
//	@desc   并按 03.7 §1 A3 over-fetch 补足条数（platform 缺失者按 D4 归漫画，本通路排除）。
//	@desc 关键词为空 → v0 `POST /v0/search/subjects`（type=1），再按条目 platform ∈ {"小说","WEB"} 过滤。
//	@desc 分页采用 over-fetch 近似：v0 通路每次请求 limit=perPage*3、offset=(page-1)*perPage*3，
//	@desc 过滤后不足 perPage 如实返回，超出截断。平台过滤会使 total 偏大（total 为书籍总数），
//	@desc hasNextPage 以「过滤后条目数 >= perPage」近似判断而非按 total 估算，分页边界略偏，属可接受近似。
//	@desc 请求/响应形状与 list-anime 完全一致（复用 media.ListAnime）。
//	@route /api/v1/anilist/list-novel [POST]
//	@returns media.ListAnime
func (h *Handler) HandleAnilistListNovel(c echo.Context) error {

	type body struct {
		Page                *int               `json:"page,omitempty"`
		Search              *string            `json:"search,omitempty"`
		PerPage             *int               `json:"perPage,omitempty"`
		Sort                []*media.MediaSort `json:"sort,omitempty"`
		Genres              []*string          `json:"genres,omitempty"`
		Tags                []*string          `json:"tags,omitempty"`
		AverageScoreGreater *int               `json:"averageScore_greater,omitempty"`
		Season              *media.MediaSeason `json:"season,omitempty"`
		SeasonYear          *int               `json:"seasonYear,omitempty"`
		Format              *media.MediaFormat `json:"format,omitempty"`
		IsAdult             *bool              `json:"isAdult,omitempty"`
		// status / countryOfOrigin 已按契约 D2/D3 删除：
		// 上游 Bangumi 无 status 与 country 维度（platform 过滤实测无效，见 §0.2），
		// 此前「接收却不消费」（仅进 cacheKey）属于静默忽略，现从接收端一并移除。
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

	// 契约 §2 D1 双通路：关键词非空 → legacy（v0 对 CJK 关键词恒 total=0，见 §0.1）。
	// 轻小说搜索框输入中文标题（如「龙族」）此前恒空，根因即缺这一分支。
	//
	// 契约 03.7 §1 A1/A2（D2/D4）：legacy 条目**不含 platform 字段**（§0.1 实测地面真值），
	// 故对命中 id 补查 v0 详情读 platform，再按 platform ∈ {"小说","WEB"} 过滤；
	// 解析失败/空 platform 按 D4 归入漫画（即本通路排除），绝不丢弃整次检索。
	// 明确降级（禁止静默忽略）：tag 在 legacy 通路仍无法本地过滤（legacy 返回体无 tags
	// 字段，契约 §0.1），故 p.Tags / p.Genres 在此降级为忽略——上游能力缺口，非遗漏。
	//
	// 契约 §2 分页补偿：legacy `max_results` 非严格保证（请求 20 可能只回 17），
	// 且 platform 过滤会进一步削减条数，故按 03.7 §1 A3 做 over-fetch 补足（见 fetchLegacyBookPage）。
	if useLegacySearch(keyword) {
		pg, err := fetchLegacyBookPage(c.Request().Context(), client, strings.TrimSpace(keyword), page, perPage, legacyBookFilter{
			novel:               true,
			averageScoreGreater: p.AverageScoreGreater,
			season:              p.Season,
			seasonYear:          p.SeasonYear,
			sorts:               p.Sort,
		})
		if err != nil {
			return h.RespondWithError(c, err)
		}

		mediaList := legacySubjectsToListAnime(pg.subjects, pg.platforms)

		total := pg.total
		pi := pg.perPage
		ret := &media.ListAnime{Page: &media.ListAnime_Page{
			Media:    mediaList,
			PageInfo: &media.PageInfo{CurrentPage: &page, PerPage: &pi, Total: &total, HasNextPage: &pg.hasNextPage},
		}}
		anilistListNovelCache.SetT(cacheKey, ret, time.Minute*10)
		return h.RespondWithData(c, ret)
	}

	filter := bangumi.SearchFilter{Type: []int{1}} // 1=书籍
	// ⚠ 禁止对书籍分区（type=1）追加 tag：实测 filter.tag 对 type=1 恒 total=0（契约 §0.2）。
	// p.Tags / p.Genres 在此有意忽略（前端已移除书籍标签筛选），非遗漏。
	// 另：此处**不得**为空关键词追加 Bangumi tag「轻小说」——一旦追加，
	// 轻小说探索会从「有结果」直接变「全空」（03.6-DIAGNOSIS §四 高风险提示）。
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
			// D9 / §4：书籍条目必须输出 type=MANGA + format ∈ {NOVEL, BOOK}
			// （修复前误用 AnimeFromSubject → type:ANIME + format:TV，见 §0.5）。
			if a := media.MangaSubjectAsListAnime(bangumi.SubjectToMedia(&res.Data[i])); a != nil {
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
