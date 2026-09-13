package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"seanime/internal/core"
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
