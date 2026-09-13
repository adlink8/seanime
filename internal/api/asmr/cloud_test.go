package asmr

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestNormalizeReviewRJID 表驱动验证 review → rjId 提取（顶层/嵌套/小写前缀）。
func TestNormalizeReviewRJID(t *testing.T) {
	cases := []struct {
		name string
		in   rawReview
		want string
	}{
		{"顶层 source_id", rawReview{SourceID: "RJ01657200"}, "RJ01657200"},
		{"嵌套 work.source_id", rawReview{Work: &struct {
			SourceID string `json:"source_id"`
		}{SourceID: "RJ00123456"}}, "RJ00123456"},
		{"小写前缀", rawReview{SourceID: "rj00789012"}, "RJ00789012"},
		{"无 source_id", rawReview{WorkID: 5}, ""},
	}
	for _, c := range cases {
		require.Equal(t, c.want, normalizeReviewRJID(c.in), c.name)
	}
}

// TestCloudReviewsParse 用 testdata fixture 验证 /api/review 解析（离线，不触网）。
func TestCloudReviewsParse(t *testing.T) {
	data, err := os.ReadFile("testdata/cloud_reviews.json")
	require.NoError(t, err)

	var reviews []rawReview
	require.NoError(t, json.Unmarshal(data, &reviews))
	require.Len(t, reviews, 3)

	// 模拟 CloudReviews 的提取逻辑
	got := make(map[string]struct{})
	for _, r := range reviews {
		if rj := normalizeReviewRJID(r); rj != "" {
			got[rj] = struct{}{}
		}
	}
	require.Contains(t, got, "RJ01657200")
	require.Contains(t, got, "RJ00123456")
	require.Contains(t, got, "RJ00789012")
	require.Len(t, got, 3)
}

// TestCloudReviewsHTTP 通过 httptest 验证 GET /api/review 走代理/鉴权头并解析。
func TestCloudReviewsHTTP(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/review", r.URL.Path)
		require.Equal(t, "marked", r.URL.Query().Get("filter"))
		// 返回 fixture
		data, _ := os.ReadFile("testdata/cloud_reviews.json")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(data)
	}
	c := newTestClient(t, ts.URL)

	res, err := c.CloudReviews(context.Background(), "marked")
	require.NoError(t, err)
	require.Contains(t, res, "RJ01657200")
	require.Len(t, res, 3)
}

// TestSaveReviewHTTP 验证写端点请求体形状（离线 httptest，不触真实网络）。
// 写端点为实测候选；此处仅确认请求可达且 body 含 work_id/marked/listening。
func TestSaveReviewHTTP(t *testing.T) {
	rs, ts := newRecordingServer(t)
	var gotBody map[string]any
	rs.handler = func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/review", r.URL.Path)
		require.Equal(t, http.MethodPut, r.Method)
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(http.StatusOK)
	}
	c := newTestClient(t, ts.URL)

	err := c.SaveReview(context.Background(), 1657200, true, false)
	require.NoError(t, err)
	require.Equal(t, float64(1657200), gotBody["work_id"])
	require.Equal(t, true, gotBody["marked"])
	require.Equal(t, false, gotBody["listening"])
}
