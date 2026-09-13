package bangumi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetSubject(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request, _ int) {
		writeJSON(t, w, Subject{
			ID:       94619,
			Type:     SubjectAnime,
			Name:     "新しい台本",
			NameCN:   "新的台本",
			Date:     "2025-04-01",
			Summary:  "这是一段中文简介",
			Platform: "TV",
			Nsfw:     true,
			Eps:      12,
			Tags:     []SubjectTag{{Name: "ASMR", Count: 3}},
		})
	}
	c := newTestClient(t, ts.URL)

	subject, err := c.GetSubject(context.Background(), 94619)
	require.NoError(t, err)

	require.Equal(t, "/v0/subjects/94619", rs.lastPath)
	require.Equal(t, "GET", rs.lastMethod)
	require.Equal(t, 94619, subject.ID)
	require.Equal(t, SubjectAnime, subject.Type)
	require.Equal(t, "新的台本", subject.NameCN)
	require.Equal(t, "新的台本", subject.NameCN)
	require.True(t, subject.Nsfw)
	require.Equal(t, 12, subject.Eps)
	require.Len(t, subject.Tags, 1)
	require.Equal(t, "ASMR", subject.Tags[0].Name)
}

func TestGetSubjectEpisodes(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request, _ int) {
		writeJSON(t, w, EpisodesResult{
			Total: 2,
			Data: []Episode{
				{ID: 1, Type: 0, Name: "ep1", NameCN: "第一话"},
				{ID: 2, Type: 1, Name: "OP", NameCN: "片头曲"},
			},
		})
	}
	c := newTestClient(t, ts.URL)

	epType := 0
	res, err := c.GetSubjectEpisodes(context.Background(), 12, &EpisodesFilter{Type: &epType})
	require.NoError(t, err)

	require.Equal(t, "/v0/subjects/12/episodes", rs.lastPath)
	require.Equal(t, "0", rs.lastQuery.Get("type"), "EpType 过滤应走 query")
	require.Equal(t, 2, res.Total)
	require.Equal(t, "第一话", res.Data[0].NameCN)
	require.Equal(t, "片头曲", res.Data[1].NameCN)
}

func TestGetSubjectEpisodesNoFilter(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request, _ int) {
		writeJSON(t, w, EpisodesResult{Total: 0})
	}
	c := newTestClient(t, ts.URL)

	_, err := c.GetSubjectEpisodes(context.Background(), 12, nil)
	require.NoError(t, err)
	require.Empty(t, rs.lastQuery, "无过滤条件时不应带 query")
}

// TestGetSubjectPlatforms_BoundedConcurrencyAndFailureIsolation 固化契约 03.7 §1 A1：
//   - 并发上限 ≤ 4（在 HTTP 边界观测真实在途请求数，而非读实现常量）；
//   - 单条失败（>=400）→ 该 id 的 platform 视为空串（D4 归漫画），**不中断**整次批量；
//   - 返回 map[int]string（id → platform），重复 id 去重（只请求一次）。
func TestGetSubjectPlatforms_BoundedConcurrencyAndFailureIsolation(t *testing.T) {
	var (
		mu        sync.Mutex
		inFlight  int
		maxFlight int
		perID     = map[int]int{}
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		inFlight++
		if inFlight > maxFlight {
			maxFlight = inFlight
		}
		mu.Unlock()

		time.Sleep(20 * time.Millisecond) // 制造重叠窗口，使并发度可被观测

		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/v0/subjects/"))

		mu.Lock()
		inFlight--
		if err == nil {
			perID[id]++
		}
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		switch {
		case err != nil:
			w.WriteHeader(http.StatusNotFound)
		case id == 2:
			// 单条失败：契约 A1 要求不中断整次检索
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"detail unavailable"}`))
		default:
			writeJSON(t, w, Subject{ID: id, Platform: "小说"})
		}
	}))
	t.Cleanup(ts.Close)

	// 关闭节流：否则 0.5s/请求的串行化会掩盖并发上限（生产的节流由 client 内部 limiter 保证）
	c := newTestClient(t, ts.URL, WithThrottleInterval(0))

	got := c.GetSubjectPlatforms(context.Background(), []int{1, 2, 3, 1, 4, 5, 6, 7, 8})

	mu.Lock()
	defer mu.Unlock()

	require.LessOrEqual(t, maxFlight, 4, "并发上限必须 ≤ 4（契约 A1）")
	require.GreaterOrEqual(t, maxFlight, 2,
		"实现必须确实并发（否则本用例无法区分串行实现）")

	require.Equal(t, 1, perID[1], "重复 id 应去重，只请求一次")
	total := 0
	for _, n := range perID {
		total += n
	}
	require.Equal(t, 8, total, "去重后应恰好请求 8 个 id")

	require.Len(t, got, 8, "返回 map[int]string（id → platform）")
	require.Equal(t, "小说", got[1])
	require.Equal(t, "小说", got[8])
	require.Empty(t, got[2], "详情失败 → platform 视为空串（A1 / D4），且不中断整次检索")
	require.Contains(t, got, 3, "单条失败不得影响其他 id")
}

// TestGetSubjectPlatforms_BypassesSharedThrottle 是 A2-3 RED 1（核心性能回归）。
//
// 背景（Wave A 遗留的阻塞级性能问题）：生产 client 是 New(token, WithProxyURL(...))，
// **未传 WithFileCache** → platform 补查在 doGet 完全落空缓存；而 client 的
// defaultThrottleInterval=500ms 是单个令牌的全局限流器（client.go:354），
// 所有请求串行每 500ms 放行一个。于是「并发 ≤4」只影响排队深度、不影响吞吐：
// 20 个未缓存 id 的补查 ≈ 19 × 500ms ≈ 9.5s，用户搜索会卡 10s 量级。
//
// 断言：20 个未缓存 id 在 500ms 节流下必须在 3s 内完成（远小于 20×500ms=10s）。
// 依据 .planning/codebase/BANGUMI-API-RECON.md:25 —— 上游 OpenAPI **无任何 rate limit
// 头/说明**，500ms 是客户端保守自选值而非上游硬要求，故 platform 补查可走不受
// 共享 ticker 约束的专用通路，仍保留并发上限与 429 退避作为安全网。
//
// 修复前本用例必然失败（实测约 9.5s）——这是本任务最重要的 RED 证据。
func TestGetSubjectPlatforms_BypassesSharedThrottle(t *testing.T) {
	const n = 20

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/v0/subjects/"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		writeJSON(t, w, Subject{ID: id, Platform: "小说"})
	}))
	t.Cleanup(ts.Close)

	// 显式使用生产的 500ms 节流：本用例就是要证明 platform 补查不再受它串行约束。
	c := newTestClient(t, ts.URL, WithThrottleInterval(500*time.Millisecond))

	ids := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		ids = append(ids, i)
	}

	start := time.Now()
	got := c.GetSubjectPlatforms(context.Background(), ids)
	elapsed := time.Since(start)
	t.Logf("20 个未缓存 id 的补查实测耗时: %v", elapsed)

	require.Len(t, got, n)
	require.Equal(t, "小说", got[1])
	require.Equal(t, "小说", got[n])
	require.Less(t, elapsed, 3*time.Second,
		"platform 补查不得被 500ms 全局串行 ticker 约束（20 id 串行约 10s；数据依据 BANGUMI-API-RECON.md:25）")
}

// TestGetSubjectPlatforms_CacheHitDoesNotHitServer 是 A2-3 RED 2。
//
// platform 是不可变的上游属性，客户端级内存缓存应在 Client 生命周期内生效：
// 同一 id 第二次调用必须命中缓存、**不发新请求**（服务端计数不增）。
// 修复前无客户端缓存 → 第二次仍触网 → callCount=2，失败。
func TestGetSubjectPlatforms_CacheHitDoesNotHitServer(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request, _ int) {
		writeJSON(t, w, Subject{ID: 94619, Platform: "漫画"})
	}
	c := newTestClient(t, ts.URL)

	first := c.GetSubjectPlatforms(context.Background(), []int{94619})
	second := c.GetSubjectPlatforms(context.Background(), []int{94619})

	require.Equal(t, "漫画", first[94619])
	require.Equal(t, "漫画", second[94619])
	require.Equal(t, 1, rs.callCount(), "第二次应命中客户端内存缓存，不触网")
}

// TestGetSubjectPlatforms_FailureNotCached 是 A2-3 RED 3。
//
// 「解析失败」不得写入缓存：一次网络抖动（此处模拟为 500）若把空串写进缓存，
// 会在 Client 生命周期内永久污染该 id，导致该条目永远被误判为漫画（D4）。
// 断言：首次失败返回空串，第二次调用能重试并拿到正确 platform（服务端计数=2）。
//
// 注意：若实现「无条件缓存」（含失败结果），本用例第二次将直出空串 → 失败。
func TestGetSubjectPlatforms_FailureNotCached(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request, call int) {
		if call == 1 {
			// 首次：模拟详情端点抖动（>=400，不触发 429 重试）
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"detail unavailable"}`))
			return
		}
		writeJSON(t, w, Subject{ID: 7, Platform: "小说"})
	}
	c := newTestClient(t, ts.URL)

	failed := c.GetSubjectPlatforms(context.Background(), []int{7})
	require.Empty(t, failed[7], "首次失败 → platform 视为空串（A1 / D4）")

	retried := c.GetSubjectPlatforms(context.Background(), []int{7})
	require.Equal(t, "小说", retried[7],
		"解析失败的 id 不得写入缓存，第二次必须重试并拿到正确 platform")
	require.Equal(t, 2, rs.callCount(), "失败不入缓存 → 第二次确实重新请求")
}
