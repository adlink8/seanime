package bangumi

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRequestHeaders(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request, _ int) {
		writeJSON(t, w, map[string]any{"ok": true})
	}
	c := newTestClient(t, ts.URL)

	_, err := c.GetMe(context.Background())
	require.NoError(t, err)

	require.Equal(t, UserAgent, rs.lastHeader.Get("User-Agent"))
	require.Equal(t, "Bearer test-token", rs.lastHeader.Get("Authorization"))
	require.Equal(t, "application/json", rs.lastHeader.Get("Accept"))
}

func TestEmptyTokenOmitsAuthHeader(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request, _ int) {
		writeJSON(t, w, map[string]any{"id": 1})
	}
	c := New("", WithBaseURL(ts.URL), WithThrottleInterval(time.Millisecond), WithRetryBaseDelay(5*time.Millisecond))
	t.Cleanup(c.Close)

	_, err := c.GetMe(context.Background())
	require.NoError(t, err)
	require.Empty(t, rs.lastHeader.Get("Authorization"))
}

func TestUnauthorized401(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request, _ int) {
		w.WriteHeader(401)
	}
	c := newTestClient(t, ts.URL)

	_, err := c.GetMe(context.Background())
	require.ErrorIs(t, err, ErrUnauthorized)
	require.Equal(t, 401, StatusCodeOf(err))
}

func TestNotFound404(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request, _ int) {
		w.WriteHeader(404)
	}
	c := newTestClient(t, ts.URL)

	_, err := c.GetSubject(context.Background(), 999999)
	require.ErrorIs(t, err, ErrNotFound)
	require.Equal(t, 404, StatusCodeOf(err))
}

func TestServerError500(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request, _ int) {
		w.WriteHeader(500)
	}
	c := newTestClient(t, ts.URL)

	_, err := c.GetSubject(context.Background(), 1)
	require.ErrorIs(t, err, ErrServer)
	require.Equal(t, 500, StatusCodeOf(err))
}

func TestRetryOn429WithRetryAfter(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request, call int) {
		if call < 3 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(429)
			return
		}
		writeJSON(t, w, map[string]any{"id": 1, "name": "ok"})
	}
	c := newTestClient(t, ts.URL)

	subject, err := c.GetSubject(context.Background(), 1)
	require.NoError(t, err)
	require.Equal(t, 1, subject.ID)
	require.Equal(t, 3, rs.callCount(), "前 2 次 429 + 第 3 次成功")
}

func TestRetryOn429ExponentialBackoff(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request, call int) {
		if call < 3 {
			// 不带 Retry-After 头，走指数退避
			w.WriteHeader(429)
			return
		}
		writeJSON(t, w, map[string]any{"id": 1})
	}
	c := newTestClient(t, ts.URL, WithRetryBaseDelay(40*time.Millisecond))

	start := time.Now()
	_, err := c.GetSubject(context.Background(), 1)
	require.NoError(t, err)

	elapsed := time.Since(start)
	// 退避应为 40ms + 80ms = 120ms；放宽到 100ms 容忍调度误差
	require.GreaterOrEqual(t, elapsed.Milliseconds(), int64(100))
	require.Less(t, elapsed, 5*time.Second)
	require.Equal(t, 3, rs.callCount())
}

func TestRetry429Exhausted(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request, _ int) {
		w.Header().Set("Retry-After", "0")
		w.WriteHeader(429)
	}
	c := newTestClient(t, ts.URL)

	_, err := c.GetSubject(context.Background(), 1)
	require.ErrorIs(t, err, ErrRateLimited)
	require.Equal(t, maxRetries+1, rs.callCount(), "1 次首发 + 3 次重试")
}

func TestThrottleInterval(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request, _ int) {
		writeJSON(t, w, map[string]any{"id": 1})
	}
	c := newTestClient(t, ts.URL, WithThrottleInterval(80*time.Millisecond))

	start := time.Now()
	for i := 0; i < 3; i++ {
		_, err := c.GetSubject(context.Background(), 1)
		require.NoError(t, err)
	}
	elapsed := time.Since(start)

	// 首个请求立即放行，后两个各等 ~80ms，总耗时 >= 160ms；放宽到 150ms
	require.GreaterOrEqual(t, elapsed.Milliseconds(), int64(150))
	require.Less(t, elapsed, 5*time.Second)
	require.Equal(t, 3, rs.callCount())
}

func TestThrottleDisabled(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request, _ int) {
		writeJSON(t, w, map[string]any{"id": 1})
	}
	c := newTestClient(t, ts.URL, WithThrottleInterval(0))

	start := time.Now()
	for i := 0; i < 5; i++ {
		_, err := c.GetSubject(context.Background(), 1)
		require.NoError(t, err)
	}
	require.Less(t, time.Since(start).Milliseconds(), int64(100), "关闭节流后应无间隔")
}

func TestCacheHitDoesNotHitServer(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request, _ int) {
		writeJSON(t, w, map[string]any{"id": 94619, "name_cn": "缓存测试"})
	}
	c, _ := newTestClientWithCache(t, ts.URL)

	s1, err := c.GetSubject(context.Background(), 94619)
	require.NoError(t, err)
	s2, err := c.GetSubject(context.Background(), 94619)
	require.NoError(t, err)

	require.Equal(t, 1, rs.callCount(), "第二次应命中缓存，不触网")
	require.Equal(t, s1.NameCN, s2.NameCN)
	require.Equal(t, s1.ID, s2.ID)
}

func TestWriteOperationsNeverCached(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request, _ int) {
		w.WriteHeader(200)
	}
	c, _ := newTestClientWithCache(t, ts.URL)

	for i := 0; i < 2; i++ {
		err := c.UpsertCollection(context.Background(), 94619, CollectionUpsertBody{Type: intPtr(CollectionDone)})
		require.NoError(t, err)
	}
	require.Equal(t, 2, rs.callCount(), "POST 不缓存，两次都触网")
}

func TestTimeoutIsApplied(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request, _ int) {
		time.Sleep(50 * time.Millisecond)
		writeJSON(t, w, map[string]any{"id": 1})
	}
	// 注入带短超时的 http.Client 验证超时生效路径（生产默认 15s）
	hc := &http.Client{Timeout: 10 * time.Millisecond}
	c := newTestClient(t, ts.URL, WithHTTPClient(hc))

	_, err := c.GetSubject(context.Background(), 1)
	require.Error(t, err, "10ms 超时应打断 50ms 的慢响应")
}
