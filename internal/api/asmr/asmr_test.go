package asmr

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// recordingServer 包装 httptest.Server，记录最近一次请求详情
type recordingServer struct {
	mu         sync.Mutex
	lastMethod string
	lastPath   string
	lastQuery  url.Values
	lastHeader http.Header
	handler    http.HandlerFunc
}

func newRecordingServer(t *testing.T) (*recordingServer, *httptest.Server) {
	t.Helper()
	rs := &recordingServer{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rs.mu.Lock()
		rs.lastMethod = r.Method
		rs.lastPath = r.URL.Path
		rs.lastQuery = r.URL.Query()
		rs.lastHeader = r.Header.Clone()
		h := rs.handler
		rs.mu.Unlock()
		if h != nil {
			h(w, r)
		}
	}))
	t.Cleanup(ts.Close)
	return rs, ts
}

func newTestClient(t *testing.T, serverURL string) *Client {
	t.Helper()
	c := New("", WithBaseURL(serverURL), WithThrottleInterval(time.Millisecond))
	t.Cleanup(c.Close)
	return c
}

func serveFixture(t *testing.T, name string) http.HandlerFunc {
	t.Helper()
	data, err := os.ReadFile("testdata/" + name)
	require.NoError(t, err)
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(data)
	}
}

func TestSearch(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = serveFixture(t, "search.json")
	c := newTestClient(t, ts.URL)

	res, err := c.Search(context.Background(), SearchParams{
		Keyword:  "RJ01234567",
		Order:    "dl",
		Page:     2,
		PerPage:  5,
		Subtitle: "zh",
	})
	require.NoError(t, err)

	// 请求形态断言：order/subtitle 走 query，page 走路径
	require.Equal(t, "GET", rs.lastMethod)
	require.Equal(t, "/search/2", rs.lastPath)
	require.Equal(t, "RJ01234567", rs.lastQuery.Get("keyword"))
	require.Equal(t, "dl_count", rs.lastQuery.Get("order"), "契约 order=dl 应映射为实测合法枚举 dl_count")
	require.Equal(t, "1", rs.lastQuery.Get("subtitle"), "契约 subtitle=zh 应映射为实测合法值 1")
	require.Equal(t, "5", rs.lastQuery.Get("pageSize"))
	require.NotEmpty(t, rs.lastHeader.Get("User-Agent"))

	// 响应解析/转换断言（契约字段名）
	require.Equal(t, 2, len(res.Works))
	w := res.Works[0]
	require.Equal(t, "1657200", w.ID, "work id 应字符串化")
	require.Equal(t, "RJ01657200", w.RjID)
	require.Equal(t, "巨乳ギャルJKとだらだらしつつ元気を取り戻していくお話 ～どすけべお姉ちゃん、派遣します・EXTrack～", w.Title)
	require.Equal(t, "青春×フェティシズム", w.Circle, "circle 取自顶层 name")
	require.Empty(t, w.Cvs)
	require.Empty(t, w.Tags)
	require.True(t, w.Nsfw)
	require.Equal(t, "https://api.asmr.one/api/cover/1657200.jpg?type=main", w.CoverURL)
	require.Equal(t, "2026-06-27", w.ReleaseDate)
	require.InDelta(t, 4.9, w.Rating, 0.001)
	require.Equal(t, 0, w.DlCount)
	require.Equal(t, 0, w.Price)
	require.True(t, w.HasSubtitle)

	// 第二条：有声优/标签
	w2 := res.Works[1]
	require.Equal(t, []string{"浅木式"}, w2.Cvs)
	require.Equal(t, []string{"纯爱/甜蜜", "ASMR"}, w2.Tags)
	require.Equal(t, 350, w2.DlCount)
	require.Equal(t, 1980, w2.Price)

	// 分页信息
	require.Equal(t, 1, res.PageInfo.CurrentPage)
	require.Equal(t, 20, res.PageInfo.PerPage)
	require.Equal(t, 62441, res.PageInfo.Total)
	require.True(t, res.PageInfo.HasNextPage, "1*20 < 62441 应有下一页")
}

func TestSearchOrderMapping(t *testing.T) {
	cases := map[string]string{
		"dd":           "create_date",
		"dl":           "dl_count",
		"dc":           "dl_count",
		"publish_date": "release",
		"":             "", // 默认：不传 order
		"unknown":      "",
	}
	for in, want := range cases {
		rs, ts := newRecordingServer(t)
		rs.handler = func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"works":[],"pagination":{"currentPage":1,"pageSize":20,"totalCount":0}}`))
		}
		c := newTestClient(t, ts.URL)
		_, err := c.Search(context.Background(), SearchParams{Keyword: "k", Order: in, Page: 1})
		require.NoError(t, err, "order=%q", in)
		require.Equal(t, want, rs.lastQuery.Get("order"), "order=%q", in)
	}
}

func TestSearchHasNextPageFalse(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"works":[],"pagination":{"currentPage":2,"pageSize":20,"totalCount":30}}`))
	}
	c := newTestClient(t, ts.URL)
	res, err := c.Search(context.Background(), SearchParams{Page: 2})
	require.NoError(t, err)
	require.False(t, res.PageInfo.HasNextPage, "2*20 >= 30 应无下一页")
}

func TestWorkInfo(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = serveFixture(t, "workinfo.json")
	c := newTestClient(t, ts.URL)

	w, err := c.WorkInfo(context.Background(), 1657200)
	require.NoError(t, err)

	require.Equal(t, "GET", rs.lastMethod)
	require.Equal(t, "/workInfo/1657200", rs.lastPath)

	require.Equal(t, "1657200", w.ID)
	require.Equal(t, "RJ01657200", w.RjID)
	require.Equal(t, "青春×フェティシズム", w.Circle)
	require.Equal(t, 4.9, w.Rating)
	require.Equal(t, "2026-06-27", w.ReleaseDate)
}

func TestTracks(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = serveFixture(t, "tracks.json")
	c := newTestClient(t, ts.URL)

	tracks, err := c.Tracks(context.Background(), 1657200)
	require.NoError(t, err)

	require.Equal(t, "GET", rs.lastMethod)
	require.Equal(t, "/tracks/1657200", rs.lastPath)
	require.Equal(t, "2", rs.lastQuery.Get("v"))

	// 根为两层：folder + 顶层 audio 叶子
	require.Len(t, tracks, 2)

	// folder 节点：type=folder，children → tracks 递归转换
	folder := tracks[0]
	require.Equal(t, "folder", folder.Type)
	require.Equal(t, "02英語", folder.Title)
	require.Empty(t, folder.MediaStreamURL)
	require.Len(t, folder.Tracks, 2, "folder 应有 2 个子节点（子 folder + lrc 叶子）")
	subFolder := folder.Tracks[0]
	require.Equal(t, "folder", subFolder.Type)
	require.Len(t, subFolder.Tracks, 2, "嵌套 folder 应有 2 个叶子")
	leaf := subFolder.Tracks[0]
	require.Equal(t, "audio", leaf.Type)
	require.Equal(t, "TrackEX_Shiki Asagi's special free talk.mp3", leaf.Title)
	require.Equal(t, "https://raw.kiko-play-niptan.one/media/stream/daily/2026-08-07/RJ01657200/01.mp3", leaf.MediaStreamURL)
	require.Equal(t, "https://raw.kiko-play-niptan.one/media/download/daily/2026-08-07/RJ01657200/01.mp3", leaf.MediaDownloadURL)
	// folder 直接子叶子（lrc 文本）
	require.Equal(t, "text", folder.Tracks[1].Type)
	require.Equal(t, "TrackEX_Shiki Asagi's special free talk.lrc", folder.Tracks[1].Title)

	// 顶层 audio 叶子
	topLeaf := tracks[1]
	require.Equal(t, "audio", topLeaf.Type)
	require.Equal(t, "01トラック.mp3", topLeaf.Title)
	require.Nil(t, topLeaf.Tracks, "叶子不应有 tracks 字段（omitempty）")
}

func TestAPIError(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"未找到该音声"}`))
	}
	c := newTestClient(t, ts.URL)

	_, err := c.WorkInfo(context.Background(), 1)
	require.Error(t, err)
	var apiErr *APIError
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, http.StatusNotFound, apiErr.Code)
	require.Contains(t, apiErr.Body, "未找到该音声")
}

// 防御：确保契约 JSON 形状不被意外改动（关键键名快照）
func TestContractShape(t *testing.T) {
	w := Asmr_Work{Cvs: []string{}, Tags: []string{}}
	b, err := json.Marshal(w)
	require.NoError(t, err)
	require.JSONEq(t, `{"id":"","rjId":"","title":"","circle":"","cvs":[],"tags":[],"nsfw":false,"coverUrl":"","releaseDate":"","rating":0,"dlCount":0,"price":0,"hasSubtitle":false}`, string(b))

	tr := Asmr_Track{Title: "x", Type: "audio", MediaStreamURL: "s", MediaDownloadURL: "d"}
	b, err = json.Marshal(tr)
	require.NoError(t, err)
	require.JSONEq(t, `{"title":"x","type":"audio","mediaStreamUrl":"s","mediaDownloadUrl":"d"}`, string(b))

	folder := Asmr_Track{Title: "f", Type: "folder", Tracks: []Asmr_Track{{Title: "c", Type: "audio"}}}
	b, err = json.Marshal(folder)
	require.NoError(t, err)
	// folder 有子节点时 tracks 必须序列化；tracks 为空 slice 时被 omitempty 省略，
	// 与契约 "tracks?"（可选）语义一致
	require.JSONEq(t, `{"title":"f","type":"folder","tracks":[{"title":"c","type":"audio"}]}`, string(b))
}
