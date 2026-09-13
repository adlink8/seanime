package bangumi

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestSearchSubjectsLegacyParsesFixture 用 testdata fixture 验证 legacy 检索响应解析（离线，不触网）。
//
// 契约依据：03.6b-CONTRACT.md §0.1 —— legacy 通路返回 {results, list[]}，
// 条目 type 为**整数**（1书籍 / 2动画），rating.score 为 0–10 浮点，
// air_date 为 "YYYY-MM-DD" 字符串，rank 为整数。
func TestSearchSubjectsLegacyParsesFixture(t *testing.T) {
	fixture, err := os.ReadFile("testdata/search_legacy.json")
	require.NoError(t, err, "fixture 必须随仓库提交（若被 .gitignore 吞掉需加例外规则）")
	require.NotEmpty(t, fixture)

	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request, _ int) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(fixture)
	}
	c := newTestClient(t, ts.URL)

	res, err := c.SearchSubjectsLegacy(context.Background(), "进击的巨人", SubjectAnime, 20, 20)
	require.NoError(t, err)

	// 请求形态（契约 §2：GET /search/subject/{urlencode(kw)}?type=&start=&max_results=）
	require.Equal(t, http.MethodGet, rs.lastMethod)
	require.Equal(t, "/search/subject/进击的巨人", rs.lastPath)
	require.Equal(t, "2", rs.lastQuery.Get("type"))
	require.Equal(t, "20", rs.lastQuery.Get("start"))
	require.Equal(t, "20", rs.lastQuery.Get("max_results"))
	require.Equal(t, UserAgent, rs.lastHeader.Get("User-Agent"))

	// 顶层 results 可读
	require.Equal(t, 127, res.Results)

	// list[] 解析
	require.Len(t, res.List, 2)

	// 条目 type 是整数
	require.Equal(t, SubjectAnime, res.List[0].Type)
	require.Equal(t, SubjectBook, res.List[1].Type)

	// rating.score 可读（0–10 制，非 0–100）
	require.NotNil(t, res.List[0].Rating)
	require.InDelta(t, 8.2, res.List[0].Rating.Score, 0.0001)

	// air_date 可读
	require.Equal(t, "2013-04-06", res.List[0].AirDate)
	require.Equal(t, "2015-09-09", res.List[1].AirDate)

	// 其余条目字段
	require.Equal(t, 486, res.List[0].ID)
	require.Equal(t, "进击的巨人", res.List[0].NameCN)
	require.Equal(t, 25, res.List[0].Eps)
	require.Equal(t, 12, res.List[0].Rank)
}

// TestSearchSubjectsLegacyOmitsZeroStart 起始页 start=0 / limit=0 时不应写入 query
// （服务端默认 start=0；避免缓存键与 URL 噪声）。
func TestSearchSubjectsLegacyOmitsZeroStart(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request, _ int) {
		writeJSON(t, w, LegacySearchResult{})
	}
	c := newTestClient(t, ts.URL)

	_, err := c.SearchSubjectsLegacy(context.Background(), "naruto", SubjectAnime, 0, 0)
	require.NoError(t, err)

	require.Equal(t, "/search/subject/naruto", rs.lastPath)
	require.Equal(t, "2", rs.lastQuery.Get("type"))
	require.Empty(t, rs.lastQuery.Get("start"))
	require.Empty(t, rs.lastQuery.Get("max_results"))
}
