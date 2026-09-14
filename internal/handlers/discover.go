package handlers

import (
	"context"
	"math/rand"
	"sort"
	"time"

	"github.com/labstack/echo/v4"
	"seanime/internal/api/asmr"
	"seanime/internal/api/bangumi"
	"seanime/internal/media"
	"seanime/internal/util/result"
)

// discover.go —— 探索页「每日推荐」（契约 03.9c）。
//
// 数据流（每域独立）：
//  1. 库抽样：从本地库取 N=5 条种子（收藏 / 最近播放 / 最近添加语义的确定性近似）；
//     - anime → anime collection（换锚后 media.ID = Bangumi subject ID）
//     - manga / novel → manga collection（按 Format 分流：NOVEL / MANGA 或空）
//     - asmr → 收藏 RJ + 最近播放 RJ（AsmrWorkState / AsmrTrackState）
//  2. 关联聚合：anime/manga/novel → Bangumi GetRelatedSubjects（v0 /subjects/{id}/subjects）；
//     asmr → item-neighbors（recommender，03.8 已有）。
//  3. 滤已入库：各域在库 id 集合内的一律剔除。
//  4. 日期种子洗牌：当天 YYYYMMDD 作 rand 源（Fisher–Yates），当日多次请求结果稳定一致，
//     跨日自然变化；取前 12 条。
//  5. 响应复用各域现有条目 DTO（media.ListAnime / media.ListManga / asmr.Asmr_SearchResult），
//     前端直接渲染卡片。
//
// 降级：库为空、上游失败、无候选 → 一律返回空列表（200），推荐是尽力而为能力，不报错。
// 缓存：进程内 10min（与 anilist list 系列 handler 缓存同级）。

var (
	discoverDailyCache = result.NewCache[string, any]()
)

const (
	// discoverDailySampleSize 每域库抽样种子数（契约 ≈5）
	discoverDailySampleSize = 5
	// discoverDailyCount 最终返回条数（契约 12 条）
	discoverDailyCount = 12
	// discoverDailyMaxSubjectLookups 聚合阶段 GetSubject 补详情的尝试上限（防上游慢查询拖垮响应）
	discoverDailyMaxSubjectLookups = 30
	// discoverDailyCacheTTL handler 缓存 TTL
	discoverDailyCacheTTL = 10 * time.Minute
)

// HandleDiscoverDaily
//
//	@summary returns daily recommendations for the discover page.
//	@desc 契约 03.9c：库抽样 → 关联聚合（Bangumi GetRelatedSubjects / asmr item-neighbors）
//	@desc → 滤已入库 → 日期种子洗牌（当日稳定、跨日变化）→ 取 12 条。
//	@desc 任何失败降级为空列表，不报错。响应复用各域现有条目 DTO。
//	@route /api/v1/discover/daily?domain=anime|manga|novel|asmr [GET]
//	@returns media.ListAnime | media.ListManga | asmr.Asmr_SearchResult
func (h *Handler) HandleDiscoverDaily(c echo.Context) error {
	domain := c.QueryParam("domain")
	switch domain {
	case "anime", "manga", "novel", "asmr":
	default:
		return h.RespondWithError(c, echo.NewHTTPError(400, "invalid domain, expected anime|manga|novel|asmr"))
	}

	cacheKey := "daily:" + domain
	if cached, ok := discoverDailyCache.Get(cacheKey); ok {
		return h.RespondWithData(c, cached)
	}

	ctx := c.Request().Context()
	seed := dailyDateSeed(time.Now())

	var ret any
	switch domain {
	case "anime":
		ret = h.discoverDailyBangumi(ctx, seed, 2, false)
	case "manga":
		ret = h.discoverDailyBangumi(ctx, seed, 1, false)
	case "novel":
		ret = h.discoverDailyBangumi(ctx, seed, 1, true)
	case "asmr":
		works, err := h.discoverDailyAsmrWorks(ctx, seed)
		if err != nil {
			works = []asmr.Asmr_Work{}
		}
		ret = &asmr.Asmr_SearchResult{Works: works}
	}

	discoverDailyCache.SetT(cacheKey, ret, discoverDailyCacheTTL)
	return h.RespondWithData(c, ret)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Bangumi 域（anime / manga / novel）
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// discoverDailyBangumi 聚合一个 Bangumi 域的每日推荐。
// subjectType：2=动画，1=书籍；novelOnly=true 时按 platform ∈ {小说, WEB} 过滤（novel 域，D2 语义），
// false 表示不按 platform 过滤（anime 域 / manga 域，书籍空 platform 按 D4 归漫画）。
// 失败降级：返回空 media 列表（不报错）。
func (h *Handler) discoverDailyBangumi(ctx context.Context, seed int64, subjectType int, novelOnly bool) any {
	var (
		seedIDs   []int
		inLibrary map[int]struct{}
	)

	switch {
	case subjectType == 2:
		seedIDs, inLibrary = h.sampleAnimeLibrarySeeds()
	case novelOnly:
		seedIDs, inLibrary = h.sampleMangaLibrarySeeds(media.MediaFormatNovel)
	default:
		seedIDs, inLibrary = h.sampleMangaLibrarySeeds(media.MediaFormatManga)
	}

	client := h.App.AnilistPlatformRef.Get().GetBangumiClient()
	if client == nil || len(seedIDs) == 0 {
		return emptyDiscoverList(subjectType)
	}

	// 1. 聚合候选：种子 → 相关条目，滤分区类型 + 已入库 + 去重
	candidates := make([]int, 0, 64)
	seen := make(map[int]struct{})
	for _, id := range seedIDs {
		rels, err := client.GetRelatedSubjects(ctx, id)
		if err != nil {
			continue // 单种子失败不中断
		}
		for _, rel := range rels {
			if rel.SubjectType != subjectType || rel.SubjectID <= 0 {
				continue
			}
			if _, ok := inLibrary[rel.SubjectID]; ok {
				continue
			}
			if _, ok := seen[rel.SubjectID]; ok {
				continue
			}
			seen[rel.SubjectID] = struct{}{}
			candidates = append(candidates, rel.SubjectID)
		}
	}
	if len(candidates) == 0 {
		return emptyDiscoverList(subjectType)
	}

	// 2. 日期种子洗牌（当日稳定、跨日变化）
	shuffleByDateSeed(candidates, seed)

	// 3. 补详情转 DTO（novel platform 过滤 / 格式转换在此完成），
	//    最多尝试 30 次防慢查询；凑满 12 条即停。
	lookups := 0
	if subjectType == 2 {
		out := make([]*media.Anime, 0, discoverDailyCount)
		for _, id := range candidates {
			if len(out) >= discoverDailyCount || lookups >= discoverDailyMaxSubjectLookups {
				break
			}
			lookups++
			subj, err := client.GetSubject(ctx, id)
			if err != nil {
				continue
			}
			if a := media.AnimeFromSubject(bangumi.SubjectToMedia(subj)); a != nil {
				out = append(out, a)
			}
		}
		page := 1
		pi := len(out)
		total := len(out)
		return &media.ListAnime{Page: &media.ListAnime_Page{
			Media:    out,
			PageInfo: &media.PageInfo{CurrentPage: &page, PerPage: &pi, Total: &total},
		}}
	}

	out := make([]*media.Anime, 0, discoverDailyCount)
	for _, id := range candidates {
		if len(out) >= discoverDailyCount || lookups >= discoverDailyMaxSubjectLookups {
			break
		}
		lookups++
		subj, err := client.GetSubject(ctx, id)
		if err != nil {
			continue
		}
		if novelOnly {
			p := subj.Platform
			if p != "小说" && p != "WEB" {
				continue
			}
		}
		// 书籍条目走「书籍语义」转换（type=MANGA + format ∈ {NOVEL, BOOK}，同 list-novel §4）。
		// manga 域空 platform 按 D4 归漫画（MangaFromSubject → BOOK 兜底，同样落入本列表）。
		if a := media.MangaSubjectAsListAnime(bangumi.SubjectToMedia(subj)); a != nil {
			if !novelOnly && a.Format != nil && *a.Format == media.MediaFormatNovel {
				continue // manga 域排除轻小说条目
			}
			out = append(out, a)
		}
	}
	page := 1
	pi := len(out)
	total := len(out)
	return &media.ListAnime{Page: &media.ListAnime_Page{
		Media:    out,
		PageInfo: &media.PageInfo{CurrentPage: &page, PerPage: &pi, Total: &total},
	}}
}

// emptyDiscoverList 与 Bangumi 域响应同构的空列表（降级返回值；
// anime/manga/novel 三域响应形状均为 media.ListAnime，前端按域取 media）。
func emptyDiscoverList(subjectType int) any {
	return &media.ListAnime{Page: &media.ListAnime_Page{Media: []*media.Anime{}}}
}

// sampleAnimeLibrarySeeds anime 域库抽样： Bangumi 收藏（collection）条目，
// 优先「在看/看过」（追番语义），不足补其余状态；按 subject ID 升序取前 N（确定性）。
// 失败降级返回空。
func (h *Handler) sampleAnimeLibrarySeeds() (seedIDs []int, inLibrary map[int]struct{}) {
	seedIDs = []int{}
	inLibrary = map[int]struct{}{}

	col, err := h.App.GetAnimeCollection(false)
	if err != nil || col == nil {
		return seedIDs, inLibrary
	}

	var prioritized, rest []int
	for _, list := range col.GetMediaListCollection().GetLists() {
		if cl := list.GetIsCustomList(); cl != nil && *cl {
			continue
		}
		for _, e := range list.GetEntries() {
			m := e.GetMedia()
			if m == nil || m.GetID() <= 0 {
				continue
			}
			id := m.GetID()
			if _, ok := inLibrary[id]; ok {
				continue
			}
			inLibrary[id] = struct{}{}
			if st := e.GetStatus(); st != nil &&
				(*st == media.MediaListStatusCurrent || *st == media.MediaListStatusCompleted) {
				prioritized = append(prioritized, id)
			} else {
				rest = append(rest, id)
			}
		}
	}

	sort.Ints(prioritized)
	sort.Ints(rest)
	seedIDs = append(prioritized, rest...)
	if len(seedIDs) > discoverDailySampleSize {
		seedIDs = seedIDs[:discoverDailySampleSize]
	}
	return seedIDs, inLibrary
}

// sampleMangaLibrarySeeds 书籍域（manga/novel）库抽样： manga collection 条目，
// wantFormat=NOVEL → 仅取轻小说条目；wantFormat=MANGA → 取漫画条目（Format 空/未知按 D4 归漫画一并纳入）。
// 排序与截断规则同 sampleAnimeLibrarySeeds。
func (h *Handler) sampleMangaLibrarySeeds(wantFormat media.MediaFormat) (seedIDs []int, inLibrary map[int]struct{}) {
	seedIDs = []int{}
	inLibrary = map[int]struct{}{}

	col, err := h.App.GetMangaCollection(false)
	if err != nil || col == nil {
		return seedIDs, inLibrary
	}

	isNovel := wantFormat == media.MediaFormatNovel
	for _, list := range col.GetMediaListCollection().GetLists() {
		if cl := list.GetIsCustomList(); cl != nil && *cl {
			continue
		}
		for _, e := range list.GetEntries() {
			m := e.GetMedia()
			if m == nil || m.GetID() <= 0 {
				continue
			}
			id := m.GetID()
			if _, ok := inLibrary[id]; ok {
				continue
			}
			inLibrary[id] = struct{}{}

			f := m.GetFormat()
			if isNovel {
				if f == nil || *f != media.MediaFormatNovel {
					continue
				}
			} else if f != nil && *f == media.MediaFormatNovel {
				continue // manga 域排除轻小说
			}
			seedIDs = append(seedIDs, id)
		}
	}

	sort.Ints(seedIDs)
	if len(seedIDs) > discoverDailySampleSize {
		seedIDs = seedIDs[:discoverDailySampleSize]
	}
	return seedIDs, inLibrary
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// ASMR 域
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// discoverDailyAsmrWorks asmr 域每日推荐： 收藏/最近播放 RJ → 数字 work id →
// item-neighbors 聚合 → 滤已知本地 RJ → 日期种子洗牌 → 取 12 条。
// 任何失败返回 error，由调用方降级为空列表。
func (h *Handler) discoverDailyAsmrWorks(ctx context.Context, seed int64) ([]asmr.Asmr_Work, error) {
	client, err := h.getAsmrClient()
	if err != nil {
		return nil, err
	}
	if lerr := h.asmrLoginIfNeeded(ctx, client); lerr != nil {
		return nil, lerr
	}

	// 1. 库抽样：收藏优先，最近播放补足
	fav, err := h.App.Database.GetAsmrFavoriteRjIDs()
	if err != nil {
		fav = nil
	}
	recent, err := h.App.Database.GetAsmrRecentlyPlayedRjIDs(discoverDailySampleSize * 2)
	if err != nil {
		recent = nil
	}

	inLibrary := make(map[string]struct{})
	for _, id := range fav {
		inLibrary[id] = struct{}{}
	}
	for _, id := range recent {
		inLibrary[id] = struct{}{}
	}

	seedRjs := make([]string, 0, discoverDailySampleSize)
	seenSeed := make(map[string]struct{})
	appendSeed := func(ids []string) {
		for _, id := range ids {
			if _, ok := seenSeed[id]; ok {
				continue
			}
			seenSeed[id] = struct{}{}
			seedRjs = append(seedRjs, id)
		}
	}
	appendSeed(fav)
	appendSeed(recent)
	if len(seedRjs) > discoverDailySampleSize {
		seedRjs = seedRjs[:discoverDailySampleSize]
	}
	if len(seedRjs) == 0 {
		return []asmr.Asmr_Work{}, nil
	}

	// 2. 聚合：RJ → 数字 work id（复用 resolveWorkID）→ item-neighbors
	candidates := make([]asmr.Asmr_Work, 0, 64)
	seenWork := make(map[string]struct{})
	for _, rj := range seedRjs {
		itemID, err := h.resolveWorkID(ctx, client, rj)
		if err != nil {
			continue // 单种子失败不中断
		}
		res, err := client.GetItemNeighbors(ctx, itemID, 1, false)
		if err != nil {
			continue // 部分作品无推荐数据/上游失败，属预期
		}
		for _, w := range res.Works {
			if w.RjID == "" {
				continue
			}
			if _, ok := inLibrary[w.RjID]; ok {
				continue
			}
			if _, ok := seenWork[w.RjID]; ok {
				continue
			}
			seenWork[w.RjID] = struct{}{}
			candidates = append(candidates, w)
		}
	}
	if len(candidates) == 0 {
		return []asmr.Asmr_Work{}, nil
	}

	// 3. 日期种子洗牌 + 截断
	works := make([]asmr.Asmr_Work, len(candidates))
	copy(works, candidates)
	shuffleByDateSeed(works, seed)
	if len(works) > discoverDailyCount {
		works = works[:discoverDailyCount]
	}
	return works, nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 纯函数（单测覆盖）
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// dailyDateSeed 日期种子：当天 YYYYMMDD（本地时间），保证当日稳定、跨日变化。
func dailyDateSeed(t time.Time) int64 {
	return int64(t.Year())*10000 + int64(t.Month())*100 + int64(t.Day())
}

// shuffleByDateSeed 按日期种子确定性洗牌（Fisher–Yates，本地 rand 源，不污染全局 rand）。
// 同一种子 + 同一输入 → 恒定输出（当日多次请求稳定一致）；跨日种子变化 → 顺序变化。
func shuffleByDateSeed[T any](items []T, seed int64) {
	rng := rand.New(rand.NewSource(seed))
	for i := len(items) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		items[i], items[j] = items[j], items[i]
	}
}
