package bangumi

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"seanime/internal/util/filecache"
)

// recordingServer 包装 httptest.Server，记录请求计数与最近一次请求详情，
// handler 可按调用序号返回不同响应（模拟 429 后恢复等场景）。
type recordingServer struct {
	mu         sync.Mutex
	count      int
	lastMethod string
	lastPath   string
	lastQuery  url.Values
	lastHeader http.Header
	lastBody   []byte
	handler    func(w http.ResponseWriter, r *http.Request, call int)
}

func newRecordingServer(t *testing.T) (*recordingServer, *httptest.Server) {
	t.Helper()
	rs := &recordingServer{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body []byte
		if r.Body != nil {
			body, _ = io.ReadAll(r.Body)
		}
		rs.mu.Lock()
		call := rs.count + 1
		rs.count = call
		rs.lastMethod = r.Method
		rs.lastPath = r.URL.Path
		rs.lastQuery = r.URL.Query()
		rs.lastHeader = r.Header.Clone()
		rs.lastBody = body
		h := rs.handler
		rs.mu.Unlock()
		if h != nil {
			h(w, r, call)
		}
	}))
	return rs, ts
}

func (rs *recordingServer) callCount() int {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	return rs.count
}

// newTestClient 创建指向 mock server 的客户端，节流/退避调短以加速测试
func newTestClient(t *testing.T, serverURL string, opts ...Option) *Client {
	t.Helper()
	all := append([]Option{
		WithBaseURL(serverURL),
		WithThrottleInterval(time.Millisecond),
		WithRetryBaseDelay(5 * time.Millisecond),
	}, opts...)
	c := New("test-token", all...)
	t.Cleanup(c.Close)
	return c
}

func newTestClientWithCache(t *testing.T, serverURL string) (*Client, *filecache.Cacher) {
	t.Helper()
	cache, err := filecache.NewCacher(t.TempDir())
	require.NoError(t, err)
	return newTestClient(t, serverURL, WithFileCache(cache)), cache
}

func writeJSON(t *testing.T, w http.ResponseWriter, v any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func intPtr(i int) *int       { return &i }
func strPtr(s string) *string { return &s }
func boolPtr(b bool) *bool    { return &b }
