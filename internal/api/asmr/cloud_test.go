package asmr

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestNormalizeReviewRJID 表驱动验证 review → rjId 提取（顶层 source_id / id 回退 / 小写前缀）。
// 实测结构（2026-09-14）：source_id 顶层必有；id 为数字 work id。
func TestNormalizeReviewRJID(t *testing.T) {
	cases := []struct {
		name string
		in   rawReview
		want string
	}{
		{"顶层 source_id", rawReview{ID: 1657200, SourceID: "RJ01657200"}, "RJ01657200"},
		{"source_id 缺失回退 id", rawReview{ID: 123456}, "RJ123456"},
		{"小写前缀", rawReview{ID: 789012, SourceID: "rj00789012"}, "RJ00789012"},
		{"全空", rawReview{}, ""},
	}
	for _, c := range cases {
		require.Equal(t, c.want, normalizeReviewRJID(c.in), c.name)
	}
}

// TestCloudReviewsParse 用 testdata fixture 验证 /api/review 分页信封解析（离线，不触网）。
// 2026-09-14 实测：真实响应为 {"works": [...], "pagination": {...}}，非裸数组。
func TestCloudReviewsParse(t *testing.T) {
	data, err := os.ReadFile("testdata/cloud_reviews.json")
	require.NoError(t, err)

	var env reviewEnvelope
	require.NoError(t, json.Unmarshal(data, &env))
	require.Len(t, env.Works, 3)
	require.Equal(t, 3, env.Pagination.TotalCount)

	// 模拟 CloudReviews 的提取逻辑
	got := make(map[string]struct{})
	for _, r := range env.Works {
		if rj := normalizeReviewRJID(r); rj != "" {
			got[rj] = struct{}{}
		}
	}
	require.Contains(t, got, "RJ01657200")
	require.Contains(t, got, "RJ123456")
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
// 契约 §0.3 / D5 / AC-04：body 恰好两字段 work_id(int) + progress(枚举字符串)，
// 不含旧实现的 marked/listening/review/rating。服务端必须解析并断言请求体内容。
func TestSaveReviewHTTP(t *testing.T) {
	rs, ts := newRecordingServer(t)
	var gotBody map[string]any
	rs.handler = func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/review", r.URL.Path)
		require.Equal(t, http.MethodPut, r.Method)
		require.NoError(t, json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusOK)
	}
	c := newTestClient(t, ts.URL)

	// favorite=true → ProgressMarked（D7）
	err := c.SaveReview(context.Background(), 1657200, ProgressMarked)
	require.NoError(t, err)

	// 恰好两个字段
	require.Len(t, gotBody, 2)
	// 仅含 work_id 与 progress
	require.Contains(t, gotBody, "work_id")
	require.Contains(t, gotBody, "progress")
	// 绝不含旧实现的多余字段
	require.NotContains(t, gotBody, "marked")
	require.NotContains(t, gotBody, "listening")
	require.NotContains(t, gotBody, "review")
	require.NotContains(t, gotBody, "rating")
	// 字段类型与值：progress 是枚举字符串（非布尔）
	require.Equal(t, float64(1657200), gotBody["work_id"])
	require.Equal(t, ProgressMarked, gotBody["progress"])
	require.IsType(t, "", gotBody["progress"])
}
