package bangumi

import (
	"errors"
	"fmt"
	"net/http"
)

// 语义化哨兵错误，通过 errors.Is(err, ErrXxx) 判定。
// *APIError 实现了 Is 方法，按状态码归类。
var (
	ErrUnauthorized = errors.New("bangumi: 鉴权失败，token 无效或过期 (401)")
	ErrNotFound     = errors.New("bangumi: 资源不存在 (404)")
	ErrRateLimited  = errors.New("bangumi: 请求过于频繁 (429)")
	ErrServer       = errors.New("bangumi: Bangumi 服务端错误 (5xx)")
)

// APIError 表示 Bangumi API 返回的非 2xx 响应。
type APIError struct {
	Code   int    // HTTP 状态码
	Method string // 请求方法
	URL    string // 请求 URL
	Body   string // 响应体（截断到 512 字符，用于排查）
}

func (e *APIError) Error() string {
	return fmt.Sprintf("bangumi: %s %s 返回 %d: %s", e.Method, e.URL, e.Code, e.Body)
}

// Is 让 errors.Is 按状态码将 APIError 归类到语义化哨兵错误。
func (e *APIError) Is(target error) bool {
	switch target {
	case ErrUnauthorized:
		return e.Code == http.StatusUnauthorized
	case ErrNotFound:
		return e.Code == http.StatusNotFound
	case ErrRateLimited:
		return e.Code == http.StatusTooManyRequests
	case ErrServer:
		return e.Code >= 500
	}
	return false
}

// StatusCodeOf 从错误中提取 HTTP 状态码；非 API 错误返回 0。
func StatusCodeOf(err error) int {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Code
	}
	return 0
}
