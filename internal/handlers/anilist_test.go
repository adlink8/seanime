package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"seanime/internal/api/bangumi"
	"seanime/internal/core"
	"seanime/internal/media"
	"seanime/internal/platforms/platform"
	"seanime/internal/testmocks"
	"seanime/internal/util"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

// newListTestHandler 构造最小可运行的 Handler：
// 仅填充 AnilistPlatformRef（系统边界依赖），用 testmocks.FakePlatform 提供平台，
// 其 GetBangumiClient() 返回 nil —— 因此 handler 在走到 Bangumi 调用时会干净返回错误响应，
// 而不是 panic。这样测试焦点集中在「缺 page/perPage 不应 panic」这一回归上。
func newListTestHandler() *Handler {
	return &Handler{
		App: &core.App{
			AnilistPlatformRef: util.NewRef[platform.Platform](testmocks.NewFakePlatformBuilder().Build()),
		},
	}
}

// TestHandleAnilistListAnime_NoPageDoesNotPanic 是 03.6-DIAGNOSIS 故障 2 的回归测试：
// HandleAnilistListAnime 在 body 缺 page/perPage 时，旧代码判空后仍裸解引用 *p.Page / *p.PerPage 导致 panic。
// 从 HTTP 边界发 POST（body 为 {}，不含 page/perPage），断言 handler 不 panic。
func TestHandleAnilistListAnime_NoPageDoesNotPanic(t *testing.T) {
	h := newListTestHandler()
	e := echo.New()
	e.POST("/api/v1/anilist/list-anime", h.HandleAnilistListAnime)

	// body 不含 page / perPage（空对象，使 Bind 成功走到判空逻辑）。
	req := httptest.NewRequest(http.MethodPost, "/api/v1/anilist/list-anime", strings.NewReader("{}"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	require.NotPanics(t, func() {
		e.ServeHTTP(rec, req)
	}, "HandleAnilistListAnime must not panic when page/perPage are absent from the request body")

	// 不 panic 即核心回归通过。FakePlatform 无 Bangumi client（离线测试无法提供），
	// 因此 handler 走到干净的错误分支：返回 500 + JSON 错误体，而不是崩溃。
	// 这两个断言有区分度：修复前 panic 会中断 ServeHTTP，Recorder 不会被写入
	// （Code 停留在默认 200、body 为空）。
	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.NotEmpty(t, rec.Body.String())
}

// TestHandleAnilistListRecentAiringAnime_NoPageDoesNotPanic 是故障 2 在 list-recent-anime 上的同源回归：
// HandleAnilistListRecentAiringAnime 同样在缺 page/perPage 时判空后仍解引用而 panic。
// 修复后该 handler 不依赖 Bangumi client，缺省直接返回 200 空列表。
func TestHandleAnilistListRecentAiringAnime_NoPageDoesNotPanic(t *testing.T) {
	h := newListTestHandler()
	e := echo.New()
	e.POST("/api/v1/anilist/list-recent-anime", h.HandleAnilistListRecentAiringAnime)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/anilist/list-recent-anime", strings.NewReader("{}"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	require.NotPanics(t, func() {
		e.ServeHTTP(rec, req)
	}, "HandleAnilistListRecentAiringAnime must not panic when page/perPage are absent from the request body")

	// 修复后该端点缺省返回 200 空列表（结构合法）。
	require.Equal(t, http.StatusOK, rec.Code)
	require.NotEmpty(t, rec.Body.String())
}

///////////////////////////////////////////////////////////////////////////////
// 契约 §2 关键词双通路（03.6b）：关键词非空 → legacy；为空 → v0
///////////////////////////////////////////////////////////////////////////////

// requestRecorder 记录 mock Bangumi 服务端收到的最近一次请求。
// 这是系统边界（外部 API）的观测点：断言 handler 实际打到了哪条通路。
type requestRecorder struct {
	mu     sync.Mutex
	method string
	path   string
	query  url.Values
	body   []byte
	calls  int
}

func (rr *requestRecorder) record(r *http.Request) {
	var body []byte
	if r.Body != nil {
		body, _ = io.ReadAll(r.Body)
	}
	rr.mu.Lock()
	defer rr.mu.Unlock()
	rr.method = r.Method
	rr.path = r.URL.Path
	rr.query = r.URL.Query()
	rr.body = body
	rr.calls++
}

func (rr *requestRecorder) snapshot() (method, path string, query url.Values) {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	return rr.method, rr.path, rr.query
}

func (rr *requestRecorder) callCount() int {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	return rr.calls
}

// lastV0Filter 解析最近一次 v0 请求体里的 filter 字段（不含时返回空 map）。
func (rr *requestRecorder) lastV0Filter(t *testing.T) map[string]json.RawMessage {
	t.Helper()
	rr.mu.Lock()
	body := append([]byte(nil), rr.body...)
	rr.mu.Unlock()

	var payload struct {
		Filter map[string]json.RawMessage `json:"filter"`
	}
	require.NoError(t, json.Unmarshal(body, &payload))
	return payload.Filter
}

// newBangumiMockClient 起一个 httptest mock 作为 Bangumi 上游，返回指向它的 client 与请求记录器。
// 测试全程离线（httptest 在本机 loopback），不发真实网络请求。
func newBangumiMockClient(t *testing.T, respond func(w http.ResponseWriter, r *http.Request)) (*bangumi.Client, *requestRecorder) {
	t.Helper()
	rr := &requestRecorder{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rr.record(r) // 记录请求（含 body），再返回 canned 响应
		w.Header().Set("Content-Type", "application/json")
		respond(w, r)
	}))
	t.Cleanup(ts.Close)

	client := bangumi.New("",
		bangumi.WithBaseURL(ts.URL),
		bangumi.WithThrottleInterval(time.Millisecond),
	)
	t.Cleanup(client.Close)
	return client, rr
}

func newListTestHandlerWithClient(client *bangumi.Client) *Handler {
	return &Handler{
		App: &core.App{
			AnilistPlatformRef: util.NewRef[platform.Platform](
				testmocks.NewFakePlatformBuilder().WithBangumiClient(client).Build(),
			),
		},
	}
}

// listAnimeEnvelope 仅取测试关心的响应字段（响应 JSON 字段名是前端已生成的契约，不得变动）。
type listAnimeEnvelope struct {
	Data struct {
		Page struct {
			Media []struct {
				ID        int     `json:"id"`
				Type      *string `json:"type"`
				Format    *string `json:"format"`
				MeanScore *int    `json:"meanScore"`
				NameCN    string  `json:"nameCN"`
			} `json:"media"`
			PageInfo struct {
				HasNextPage *bool `json:"hasNextPage"`
				Total       *int  `json:"total"`
			} `json:"pageInfo"`
		} `json:"page"`
	} `json:"data"`
}

const legacyAnimeFixture = `{"results":127,"list":[
  {"id":486,"type":2,"name":"進撃の巨人","name_cn":"进击的巨人","air_date":"2013-04-06","eps":25,"rank":12,
   "rating":{"total":12000,"count":{},"score":8.2}}
]}`

const v0AnimeFixture = `{"total":1,"limit":20,"offset":0,"data":[
  {"id":486,"type":2,"name":"進撃の巨人","name_cn":"进击的巨人","date":"2013-04-06","eps":25,"platform":"TV",
   "rating":{"rank":12,"total":12000,"score":8.2}}
]}`

// TestHandleAnilistListAnime_KeywordRoutesToLegacy 断言：关键词非空必须走 legacy 通路
// （契约 §2 D1）。v0 通路对 CJK 关键词恒 total=0，是「中文搜不到」的根因。
func TestHandleAnilistListAnime_KeywordRoutesToLegacy(t *testing.T) {
	client, rr := newBangumiMockClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(legacyAnimeFixture))
	})
	h := newListTestHandlerWithClient(client)

	e := echo.New()
	e.POST("/api/v1/anilist/list-anime", h.HandleAnilistListAnime)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/anilist/list-anime",
		strings.NewReader(`{"search":"进击的巨人","page":1,"perPage":20}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code, "body=%s", rec.Body.String())

	method, path, query := rr.snapshot()
	require.Equal(t, http.MethodGet, method)
	require.Equal(t, "/search/subject/进击的巨人", path, "非空关键词必须打到 legacy 端点（契约 §2 通路 1）")
	require.Equal(t, "2", query.Get("type"), "动画分区 type=2")
	require.Equal(t, "20", query.Get("max_results"))
	require.Empty(t, query.Get("start"), "page=1 → start=0 应省略")

	var env listAnimeEnvelope
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &env))
	require.Len(t, env.Data.Page.Media, 1)
	require.Equal(t, 486, env.Data.Page.Media[0].ID)
	require.NotNil(t, env.Data.Page.Media[0].Type)
	require.Equal(t, "ANIME", *env.Data.Page.Media[0].Type)
	require.NotNil(t, env.Data.Page.Media[0].MeanScore)
	require.Equal(t, 82, *env.Data.Page.Media[0].MeanScore, "legacy rating.score 8.2（0-10）→ meanScore 82（0-100）")
}

// TestHandleAnilistListAnime_NoKeywordRoutesToV0 断言：关键词为空必须走 v0 filter 通路
// （契约 §2 通路 2），且发出的 limit 为 20、offset 按 20 步进（D8）。
func TestHandleAnilistListAnime_NoKeywordRoutesToV0(t *testing.T) {
	client, rr := newBangumiMockClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(v0AnimeFixture))
	})
	h := newListTestHandlerWithClient(client)

	e := echo.New()
	e.POST("/api/v1/anilist/list-anime", h.HandleAnilistListAnime)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/anilist/list-anime",
		strings.NewReader(`{"page":3,"perPage":48}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code, "body=%s", rec.Body.String())

	method, path, query := rr.snapshot()
	require.Equal(t, http.MethodPost, method)
	require.Equal(t, "/v0/search/subjects", path, "空关键词必须走 v0 通路（契约 §2 通路 2）")
	require.Equal(t, "20", query.Get("limit"), "v0 limit 固定 20（D8）")
	require.Equal(t, "40", query.Get("offset"), "offset 按 20 步进：page=3 → 40（D8）")
}

// postListAnime 在 handler 边界发一次 list-anime 请求，返回响应与上游请求快照。
func postListAnime(t *testing.T, client *bangumi.Client, body string) (*httptest.ResponseRecorder, *requestRecorder) {
	t.Helper()
	h := newListTestHandlerWithClient(client)
	e := echo.New()
	e.POST("/api/v1/anilist/list-anime", h.HandleAnilistListAnime)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/anilist/list-anime", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code, "body=%s", rec.Body.String())
	return rec, nil
}

func legacyMockIDs(t *testing.T, rec *httptest.ResponseRecorder) []int {
	t.Helper()
	var env listAnimeEnvelope
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &env))
	ids := make([]int, 0, len(env.Data.Page.Media))
	for _, m := range env.Data.Page.Media {
		ids = append(ids, m.ID)
	}
	return ids
}

// TestHandleAnilistListAnime_LegacyLocalRatingAndYearFilter 断言 legacy 通路的本地过滤与本地排序
// （契约 §2 通路 1）：评分下限按 0–10 比 rating.score，年份按 air_date 过滤，
// `score` 排序 → rating.score 降序。
//
// 独立期望值：候选 3 条为 (2013, 7.5) / (2013, 8.2) / (2015, 9.1)；
// 阈值 75 → ≥7.5 全通过；seasonYear=2013 → 丢弃第三条；SCORE_DESC → 8.2 在前。
// 期望 [2, 1] 由上述字面量手工推导，与实现无关。
func TestHandleAnilistListAnime_LegacyLocalRatingAndYearFilter(t *testing.T) {
	client, rr := newBangumiMockClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"results":3,"list":[
		  {"id":1,"type":2,"name":"a","air_date":"2013-04-06","rating":{"total":10,"count":{},"score":7.5}},
		  {"id":2,"type":2,"name":"b","air_date":"2013-10-01","rating":{"total":10,"count":{},"score":8.2}},
		  {"id":3,"type":2,"name":"c","air_date":"2015-01-05","rating":{"total":10,"count":{},"score":9.1}}
		]}`))
	})

	rec, _ := postListAnime(t, client,
		`{"search":"巨人","averageScore_greater":75,"seasonYear":2013,"sort":["SCORE_DESC"],"page":1,"perPage":20}`)

	_, path, _ := rr.snapshot()
	require.Equal(t, "/search/subject/巨人", path)
	require.Equal(t, []int{2, 1}, legacyMockIDs(t, rec))
}

// TestHandleAnilistListAnime_LegacyLocalRankAscending 断言 `rank` 排序 → 按 rank 升序（rank 小者优）。
// 若误用 score 降序会得到 [1, 2]，故该用例对实现有区分度。
func TestHandleAnilistListAnime_LegacyLocalRankAscending(t *testing.T) {
	client, _ := newBangumiMockClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"results":2,"list":[
		  {"id":1,"type":2,"name":"a","rank":30,"rating":{"total":10,"count":{},"score":7.0}},
		  {"id":2,"type":2,"name":"b","rank":5,"rating":{"total":10,"count":{},"score":6.0}}
		]}`))
	})

	rec, _ := postListAnime(t, client, `{"search":"a","sort":["POPULARITY_DESC"],"page":1,"perPage":20}`)
	require.Equal(t, []int{2, 1}, legacyMockIDs(t, rec))
}

// TestHandleAnilistListAnime_LegacyHeatKeepsUpstreamOrder 断言 `heat`/`match` → 保持上游返回序
// （契约 §2 通路 1：legacy 默认按相关度）。若误按 score 排序会得到 [2, 1]。
func TestHandleAnilistListAnime_LegacyHeatKeepsUpstreamOrder(t *testing.T) {
	client, _ := newBangumiMockClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"results":2,"list":[
		  {"id":1,"type":2,"name":"a","rating":{"total":10,"count":{},"score":6.0}},
		  {"id":2,"type":2,"name":"b","rating":{"total":10,"count":{},"score":9.0}}
		]}`))
	})

	rec, _ := postListAnime(t, client, `{"search":"a","sort":["TRENDING_DESC"],"page":1,"perPage":20}`)
	require.Equal(t, []int{1, 2}, legacyMockIDs(t, rec))
}

// TestHandleAnilistListAnime_LegacyTagFilterDegradedNotDropped 固化「明确降级」行为：
// legacy 条目不含 tag 字段（契约 §0.1），tag/meta_tags 在 legacy 通路无法本地过滤，
// 因此**不得**据此丢弃条目（禁止静默忽略式的全空结果）。
func TestHandleAnilistListAnime_LegacyTagFilterDegradedNotDropped(t *testing.T) {
	client, _ := newBangumiMockClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"results":2,"list":[
		  {"id":1,"type":2,"name":"a","rating":{"total":10,"count":{},"score":6.0}},
		  {"id":2,"type":2,"name":"b","rating":{"total":10,"count":{},"score":9.0}}
		]}`))
	})

	rec, _ := postListAnime(t, client, `{"search":"a","tags":["Action"],"format":"TV","page":1,"perPage":20}`)
	require.Equal(t, []int{1, 2}, legacyMockIDs(t, rec),
		"tag/format 在 legacy 通路降级为忽略（契约 §0.1），不得丢条目")
}

///////////////////////////////////////////////////////////////////////////////
// 契约 §2：list-manga（type=1 书籍）双通路
///////////////////////////////////////////////////////////////////////////////

func postListManga(t *testing.T, client *bangumi.Client, body string) *httptest.ResponseRecorder {
	t.Helper()
	h := newListTestHandlerWithClient(client)
	e := echo.New()
	e.POST("/api/v1/manga/anilist/list", h.HandleAnilistListManga)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/manga/anilist/list", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code, "body=%s", rec.Body.String())
	return rec
}

// TestHandleAnilistListManga_KeywordRoutesToLegacy 断言漫画分区（type=1）关键词非空 → legacy。
func TestHandleAnilistListManga_KeywordRoutesToLegacy(t *testing.T) {
	client, rr := newBangumiMockClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"results":448,"list":[
		  {"id":276764,"type":1,"name":"進撃の巨人","name_cn":"进击的巨人","air_date":"2015-09-09","eps":4,
		   "rating":{"total":300,"count":{},"score":7.1}}
		]}`))
	})

	rec := postListManga(t, client, `{"search":"进击的巨人","page":2,"perPage":20}`)

	_, path, query := rr.snapshot()
	require.Equal(t, "/search/subject/进击的巨人", path, "漫画分区关键词非空必须走 legacy（契约 §2 通路 1）")
	require.Equal(t, "1", query.Get("type"), "书籍分区 type=1")
	require.Equal(t, "20", query.Get("start"), "page=2 → start=20（契约 §2）")
	require.Equal(t, "20", query.Get("max_results"))

	require.Equal(t, []int{276764}, legacyMockIDs(t, rec))
}

// TestHandleAnilistListManga_NoKeywordRoutesToV0 断言漫画分区关键词为空 → v0，且分页收敛 20。
func TestHandleAnilistListManga_NoKeywordRoutesToV0(t *testing.T) {
	client, rr := newBangumiMockClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"total":1,"limit":20,"offset":20,"data":[
		  {"id":276764,"type":1,"name":"進撃の巨人","name_cn":"进击的巨人","date":"2015-09-09","platform":"漫画",
		   "rating":{"rank":0,"total":300,"score":7.1}}
		]}`))
	})

	rec := postListManga(t, client, `{"page":2,"perPage":48}`)

	_, path, query := rr.snapshot()
	require.Equal(t, "/v0/search/subjects", path, "空关键词走 v0（契约 §2 通路 2）")
	require.Equal(t, "20", query.Get("limit"), "v0 limit 固定 20（D8）")
	require.Equal(t, "20", query.Get("offset"), "offset 按 20 步进：page=2 → 20（D8）")

	// 书籍条目经映射后 Type=MANGA、Format 使用书籍语义（D9）
	var env listAnimeEnvelope
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &env))
	require.Len(t, env.Data.Page.Media, 1)
	require.NotNil(t, env.Data.Page.Media[0].Type)
	require.Equal(t, "MANGA", *env.Data.Page.Media[0].Type)
	require.NotNil(t, env.Data.Page.Media[0].Format)
	require.Equal(t, "MANGA", *env.Data.Page.Media[0].Format, "platform=漫画 → format MANGA")
}

// TestHandleAnilistListManga_LegacyLocalYearAndRatingFilter 断言漫画 legacy 通路的本地过滤与排序。
// 候选：(2015, 7.1) / (2015, 8.4) / (2016, 9.0)；阈值 80 保留后两条的年份 2015 里 8.4；
// year=2015 丢弃 2016 条；SCORE_DESC → 8.4 在前。期望 [12, 11]。
func TestHandleAnilistListManga_LegacyLocalYearAndRatingFilter(t *testing.T) {
	client, _ := newBangumiMockClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"results":3,"list":[
		  {"id":11,"type":1,"name":"a","air_date":"2015-01-05","rating":{"total":10,"count":{},"score":7.1}},
		  {"id":12,"type":1,"name":"b","air_date":"2015-06-05","rating":{"total":10,"count":{},"score":8.4}},
		  {"id":13,"type":1,"name":"c","air_date":"2016-06-05","rating":{"total":10,"count":{},"score":9.0}}
		]}`))
	})

	rec := postListManga(t, client, `{"search":"a","averageScore_greater":80,"year":2015,"sort":["SCORE_DESC"],"page":1,"perPage":20}`)
	require.Equal(t, []int{12}, legacyMockIDs(t, rec))
}

///////////////////////////////////////////////////////////////////////////////
// 契约 §1 D4：format → meta_tags 收敛（仅 TV / WEB / OVA 三值有效）
///////////////////////////////////////////////////////////////////////////////

func TestAnimeMetaTagsFromFormat(t *testing.T) {
	cases := []struct {
		name   string
		format *media.MediaFormat
		want   []string
	}{
		{"TV → TV", formatPtr(media.MediaFormatTv), []string{"TV"}},
		{"ONA → WEB", formatPtr(media.MediaFormatOna), []string{"WEB"}},
		{"OVA → OVA", formatPtr(media.MediaFormatOva), []string{"OVA"}},
		{"MOVIE 不再产生 meta_tags", formatPtr(media.MediaFormatMovie), nil},
		{"TV_SHORT 不再产生 meta_tags", formatPtr(media.MediaFormatTvShort), nil},
		{"SPECIAL 不再产生 meta_tags", formatPtr(media.MediaFormatSpecial), nil},
		{"MUSIC 不产生 meta_tags", formatPtr(media.MediaFormatMusic), nil},
		{"nil 不产生 meta_tags", nil, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, animeMetaTagsFromFormat(tc.format))
		})
	}
}

func formatPtr(f media.MediaFormat) *media.MediaFormat { return &f }

// TestHandleAnilistListAnime_FormatMapsToMetaTags 在 handler 边界确认 format 实际落到 filter.meta_tags，
// 且 `MOVIE`/`TV_SHORT`/`SPECIAL` 不再产生任何 meta_tags（契约 D4）。
func TestHandleAnilistListAnime_FormatMapsToMetaTags(t *testing.T) {
	cases := []struct {
		format string
		want   []string
	}{
		{"TV", []string{"TV"}},
		{"ONA", []string{"WEB"}},
		{"OVA", []string{"OVA"}},
		{"MOVIE", nil},
		{"TV_SHORT", nil},
		{"SPECIAL", nil},
	}
	for _, tc := range cases {
		t.Run("format="+tc.format, func(t *testing.T) {
			client, rr := newBangumiMockClient(t, func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte(v0AnimeFixture))
			})
			postListAnime(t, client, fmt.Sprintf(`{"page":1,"perPage":20,"format":%q}`, tc.format))

			filter := rr.lastV0Filter(t)
			raw, has := filter["meta_tags"]
			if tc.want == nil {
				require.False(t, has, "format=%s 不应产生 meta_tags（契约 D4）", tc.format)
				return
			}
			require.True(t, has, "format=%s 应产生 meta_tags", tc.format)
			var got []string
			require.NoError(t, json.Unmarshal(raw, &got))
			require.Equal(t, tc.want, got)
		})
	}
}

// TestHandleAnilistListManga_EmptyKeywordOmitsBookTagFilter 断言书籍分区（type=1）的 v0 请求
// **禁止携带 tag / meta_tags**（实测恒 total=0，契约 §0.2 / §2 通路 2）。
func TestHandleAnilistListManga_EmptyKeywordOmitsBookTagFilter(t *testing.T) {
	client, rr := newBangumiMockClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"total":0,"limit":20,"offset":0,"data":[]}`))
	})
	postListManga(t, client, `{"page":1,"perPage":20,"tags":["Action"],"genres":["Comedy"],"format":"MANGA"}`)

	filter := rr.lastV0Filter(t)
	_, hasTag := filter["tag"]
	require.False(t, hasTag, "书籍分区禁止 tag（type=1 恒 total=0，契约 §0.2）")
	_, hasMeta := filter["meta_tags"]
	require.False(t, hasMeta, "书籍分区禁止 meta_tags（type=1 恒 0，契约 §0.2）")

	rawType, hasType := filter["type"]
	require.True(t, hasType)
	var gotType []int
	require.NoError(t, json.Unmarshal(rawType, &gotType))
	require.Equal(t, []int{1}, gotType, "书籍分区 type=1")
}

// TestHandleAnilistListNovel_EmptyKeywordOmitsBookTagFilter 同源约束：轻小说分区（书籍）禁止 tag。
func TestHandleAnilistListNovel_EmptyKeywordOmitsBookTagFilter(t *testing.T) {
	client, rr := newBangumiMockClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"total":0,"limit":20,"offset":0,"data":[]}`))
	})
	h := newListTestHandlerWithClient(client)
	e := echo.New()
	e.POST("/api/v1/anilist/list-novel", h.HandleAnilistListNovel)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/anilist/list-novel",
		strings.NewReader(`{"page":1,"perPage":20,"tags":["轻小说"],"genres":["Romance"]}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code, "body=%s", rec.Body.String())

	filter := rr.lastV0Filter(t)
	_, hasTag := filter["tag"]
	require.False(t, hasTag, "书籍分区禁止 tag（契约 §0.2；曾因追加「轻小说」tag 导致探索页全空）")
}

// TestHandleAnilistListAnime_StatusAndCountryOfOriginNotInCacheKey 断言 status / countryOfOrigin
// 已完全退出语义与缓存键（契约 D2/D3）：两次请求仅这两项不同 → 只打一次上游（第二次命中 handler 缓存）。
func TestHandleAnilistListAnime_StatusAndCountryOfOriginNotInCacheKey(t *testing.T) {
	client, rr := newBangumiMockClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(v0AnimeFixture))
	})
	h := newListTestHandlerWithClient(client)
	e := echo.New()
	e.POST("/api/v1/anilist/list-anime", h.HandleAnilistListAnime)

	send := func(body string) {
		t.Helper()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/anilist/list-anime", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code, "body=%s", rec.Body.String())
	}

	// seasonYear=2099 让缓存键与其他用例隔离
	send(`{"page":1,"perPage":20,"seasonYear":2099,"status":["FINISHED"],"countryOfOrigin":"JP"}`)
	send(`{"page":1,"perPage":20,"seasonYear":2099,"status":["RELEASING"],"countryOfOrigin":"KR"}`)

	require.Equal(t, 1, rr.callCount(),
		"status/countryOfOrigin 不再参与语义与缓存键（D2/D3），第二次请求应命中缓存")
}

///////////////////////////////////////////////////////////////////////////////
// 契约 §2：list-novel（type=1 书籍，轻小说）双通路
//
// Wave A 已给 anime/manga 接上「关键词非空 → legacy」双通路，唯独漏了 novel：
// 轻小说搜索框输中文标题（如「龙族」）仍返回空。以下用例固化 novel 的同款契约。
///////////////////////////////////////////////////////////////////////////////

// postListNovel 在 handler 边界发一次 list-novel 请求（断言 200）。
func postListNovel(t *testing.T, client *bangumi.Client, body string) *httptest.ResponseRecorder {
	t.Helper()
	h := newListTestHandlerWithClient(client)
	e := echo.New()
	e.POST("/api/v1/anilist/list-novel", h.HandleAnilistListNovel)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/anilist/list-novel", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code, "body=%s", rec.Body.String())
	return rec
}

// TestHandleAnilistListNovel_KeywordRoutesToLegacy 断言轻小说分区关键词非空 → legacy
// （契约 §2 D1）。v0 对 CJK 关键词恒 total=0，是「中文轻小说搜不到」的根因。
// 若仍打 v0（本用例会看到 POST /v0/search/subjects），中文搜索必空。
func TestHandleAnilistListNovel_KeywordRoutesToLegacy(t *testing.T) {
	client, rr := newBangumiMockClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"results":448,"list":[
		  {"id":331401,"type":1,"name":"龍族","name_cn":"龙族","air_date":"2010-04-01","eps":10,
		   "rating":{"total":120,"count":{},"score":8.0}}
		]}`))
	})

	rec := postListNovel(t, client, `{"search":"龙族","page":2,"perPage":20}`)

	method, path, query := rr.snapshot()
	require.Equal(t, http.MethodGet, method, "非空关键词必须走 legacy GET（契约 §2 通路 1）")
	require.Equal(t, "/search/subject/龙族", path, "必须打到 legacy 端点，而不是 /v0/search/subjects")
	require.Equal(t, "1", query.Get("type"), "轻小说属书籍分区 type=1（契约 §4）")
	require.Equal(t, "20", query.Get("start"), "page=2 → start=20（契约 §2）")
	require.Equal(t, "20", query.Get("max_results"))

	require.Equal(t, []int{331401}, legacyMockIDs(t, rec))
}

// TestHandleAnilistListNovel_NoKeywordRoutesToV0 断言轻小说分区关键词为空 → 维持 v0 通路
// （契约 §2 通路 2），且书籍分区禁止 tag（type=1 恒 total=0，契约 §0.2）。
func TestHandleAnilistListNovel_NoKeywordRoutesToV0(t *testing.T) {
	client, rr := newBangumiMockClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"total":0,"limit":20,"offset":0,"data":[]}`))
	})

	rec := postListNovel(t, client, `{"page":2,"perPage":20}`)

	method, path, _ := rr.snapshot()
	require.Equal(t, http.MethodPost, method, "空关键词仍走 v0 POST（契约 §2 通路 2）")
	require.Equal(t, "/v0/search/subjects", path, "空关键词不得走 legacy")
	require.Equal(t, http.StatusOK, rec.Code)

	filter := rr.lastV0Filter(t)
	rawType, hasType := filter["type"]
	require.True(t, hasType)
	var gotType []int
	require.NoError(t, json.Unmarshal(rawType, &gotType))
	require.Equal(t, []int{1}, gotType, "书籍分区 type=1")
	_, hasTag := filter["tag"]
	require.False(t, hasTag, "书籍分区禁止 tag（契约 §0.2）")
}

// TestHandleAnilistListNovel_LegacyNoPlatformFilter 固化 legacy 通路的「明确降级」：
// legacy 条目**没有 `platform` 字段**（契约 §0.1 实测地面真值，`responseGroup=large` 亦然），
// 故轻小说的 platform ∈ {小说,WEB} 过滤在 legacy 通路**无法本地执行，降级为忽略**。
// 这是上游能力缺口而非疏漏——**不得**因为拿不到 platform 就把条目全部丢弃（结果全空）。
func TestHandleAnilistListNovel_LegacyNoPlatformFilter(t *testing.T) {
	// 两条书目均无 platform 字段（与 legacy 上游返回体一致）。
	client, rr := newBangumiMockClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"results":2,"list":[
		  {"id":1,"type":1,"name":"a","name_cn":"甲","air_date":"2010-04-01","rating":{"total":10,"count":{},"score":7.0}},
		  {"id":2,"type":1,"name":"b","name_cn":"乙","air_date":"2011-04-01","rating":{"total":10,"count":{},"score":8.0}}
		]}`))
	})

	rec := postListNovel(t, client, `{"search":"幻想","page":1,"perPage":20}`)

	_, path, _ := rr.snapshot()
	require.Equal(t, "/search/subject/幻想", path)
	require.Equal(t, []int{1, 2}, legacyMockIDs(t, rec),
		"legacy 无 platform 字段 → platform 过滤降级为忽略，不得据此清空条目")

	// 书籍语义仍须保持（D9）：type=MANGA、format ∈ {NOVEL, BOOK}。
	var env listAnimeEnvelope
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &env))
	require.Len(t, env.Data.Page.Media, 2)
	for _, m := range env.Data.Page.Media {
		require.NotNil(t, m.Type)
		require.Equal(t, "MANGA", *m.Type, "书籍条目必须输出 type=MANGA（D9）")
		require.NotNil(t, m.Format)
		require.Contains(t, []string{"NOVEL", "BOOK"}, *m.Format, "format 必须 ∈ {NOVEL, BOOK}（D9）")
	}
}

// TestHandleAnilistListNovel_LegacyLocalRatingAndYearFilter 断言 novel legacy 通路确实接上了
// 本地过滤与本地排序（契约 §2 通路 1），与 anime/manga 通路一致：评分下限按 0–10 比 rating.score，
// 年份按 air_date 过滤，SCORE_DESC → rating.score 降序。
//
// 候选 (2014, 7.0) / (2014, 8.6) / (2016, 9.0)；阈值 80 → 丢弃 7.0；seasonYear=2014 → 丢弃 2016 条；
// SCORE_DESC → 8.6 在前。期望 [12] 由字面量手工推导，与实现无关。
func TestHandleAnilistListNovel_LegacyLocalRatingAndYearFilter(t *testing.T) {
	client, rr := newBangumiMockClient(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"results":3,"list":[
		  {"id":11,"type":1,"name":"a","air_date":"2014-01-05","rating":{"total":10,"count":{},"score":7.0}},
		  {"id":12,"type":1,"name":"b","air_date":"2014-06-05","rating":{"total":10,"count":{},"score":8.6}},
		  {"id":13,"type":1,"name":"c","air_date":"2016-06-05","rating":{"total":10,"count":{},"score":9.0}}
		]}`))
	})

	rec := postListNovel(t, client,
		`{"search":"幻想","averageScore_greater":80,"seasonYear":2014,"sort":["SCORE_DESC"],"page":1,"perPage":20}`)

	_, path, _ := rr.snapshot()
	require.Equal(t, "/search/subject/幻想", path)
	require.Equal(t, []int{12}, legacyMockIDs(t, rec))
}

///////////////////////////////////////////////////////////////////////////////
// 契约 §4 / D9：list-novel 书籍类型映射（Phase 3.6c）
///////////////////////////////////////////////////////////////////////////////

// TestHandleAnilistListNovel_BookTypeMapping 断言 list-novel 返回的书籍条目为书籍语义：
// type=MANGA、format ∈ {NOVEL, BOOK}。修复前用 AnimeFromSubject → type=ANIME + format=TV（§0.5）。
//
// 注意：platform→format 的推导只存在于 v0 通路（legacy 条目**无 platform 字段**，契约 §0.1），
// 故本用例必须以**空关键词**发出请求，走 §2 通路 2 的 v0 分支（非空关键词现已改走 legacy）。
// 各子用例用不同 seasonYear 隔离 handler 级缓存（缓存键含 seasonYear）。
func TestHandleAnilistListNovel_BookTypeMapping(t *testing.T) {
	cases := []struct {
		platform   string
		wantFormat string
		seasonYear int
	}{
		{"小说", "NOVEL", 2014},
		{"WEB", "BOOK", 2015},
	}

	for _, tc := range cases {
		t.Run("platform="+tc.platform, func(t *testing.T) {
			client, _ := newBangumiMockClient(t, func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte(`{"total":1,"limit":20,"offset":0,"data":[
				  {"id":276764,"type":1,"name":"無職転生","name_cn":"无职转生","date":"2014-01-24","platform":"` + tc.platform + `",
				   "eps":28,"rating":{"rank":0,"total":300,"score":9.2}}
				]}`))
			})

			h := newListTestHandlerWithClient(client)
			e := echo.New()
			e.POST("/api/v1/anilist/list-novel", h.HandleAnilistListNovel)

			// 空关键词 → v0 通路（platform 字段存在，才谈得上 platform→format 推导）。
			req := httptest.NewRequest(http.MethodPost, "/api/v1/anilist/list-novel",
				strings.NewReader(fmt.Sprintf(`{"page":1,"perPage":20,"seasonYear":%d}`, tc.seasonYear)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			require.Equal(t, http.StatusOK, rec.Code, "body=%s", rec.Body.String())

			var env listAnimeEnvelope
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &env))
			require.Len(t, env.Data.Page.Media, 1)

			m := env.Data.Page.Media[0]
			require.NotNil(t, m.Type)
			require.Equal(t, "MANGA", *m.Type, "书籍条目必须输出 type=MANGA（D9）")
			require.NotNil(t, m.Format)
			require.Equal(t, tc.wantFormat, *m.Format, "format 必须 ∈ {NOVEL, BOOK}（D9）")
			require.Equal(t, "无职转生", m.NameCN)
		})
	}
}
