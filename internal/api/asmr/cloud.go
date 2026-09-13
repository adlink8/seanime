package asmr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

// 云端账号同步（契约 §6）。
//
// 地面真值（编排层实测，勿重复实测）：
//   - 登录：POST /api/auth/me {name, password} → JWT（契约记为 token）
//   - 读取：GET /api/review?order=updated_at&sort=desc&page=1&filter=marked|listening 实测 200
//   - 写端点：第三方 works_review.py 佐证 fire-and-forget POST；Wave A 以实测为准（见 SaveReview）
//
// 凭据来源：调用方（handler）从 config 注入，本 client 不读文件/环境变量。
// 严禁把凭据/JWT 写入日志、文件、测试或提交（契约硬约束）。

// authResponse /api/auth/me 登录响应（asmr.one 实测返回 token 字段；兼容 access_token）。
type authResponse struct {
	Token       string `json:"token"`
	AccessToken string `json:"access_token"`
}

// rawReview /api/review 列表单条（字段名宽松解析，按实测为准）。
// asmr.one 单条至少含 work_id；source_id（RJxxxx）可能在顶层或嵌套 work 内。
type rawReview struct {
	WorkID   int    `json:"work_id"`
	SourceID string `json:"source_id"`
	Work     *struct {
		SourceID string `json:"source_id"`
	} `json:"work"`
	UserReview *struct {
		Review string `json:"review"`
	} `json:"user_review"`
}

// jwtCache 内存缓存登录态；401 时由调用方触发重登一次。
type jwtCache struct {
	mu       sync.Mutex
	token    string
	name     string
	password string
}

// Login 登录 asmr.one 并缓存 JWT 内存。成功后后续请求自动带 Authorization。
// 凭据仅在本次调用传入并缓存于 client，绝不外泄。
func (c *Client) Login(ctx context.Context, name, password string) (string, error) {
	c.jwt.mu.Lock()
	c.jwt.name = name
	c.jwt.password = password
	c.jwt.mu.Unlock()

	token, err := c.doLogin(ctx, name, password)
	if err != nil {
		return "", err
	}
	c.jwt.mu.Lock()
	c.jwt.token = token
	c.jwt.mu.Unlock()
	c.token = token
	return token, nil
}

func (c *Client) doLogin(ctx context.Context, name, password string) (string, error) {
	body := map[string]string{"name": name, "password": password}
	var out authResponse
	if err := c.doPost(ctx, "/auth/me", body, &out); err != nil {
		return "", err
	}
	if out.Token != "" {
		return out.Token, nil
	}
	if out.AccessToken != "" {
		return out.AccessToken, nil
	}
	return "", fmt.Errorf("asmr: /api/auth/me 响应未包含 token")
}

// doPost 发送 JSON POST；带 Bearer token；遇 401 重登一次后重试。
// POST 永不走缓存（写操作必须打到服务端）。
func (c *Client) doPost(ctx context.Context, path string, body any, out any) error {
	if err := c.limiter.wait(ctx); err != nil {
		return err
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("asmr: 序列化请求体失败: %w", err)
	}

	do := func() (*http.Response, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL.JoinPath(path).String(), bytes.NewReader(payload))
		if err != nil {
			return nil, fmt.Errorf("asmr: 构造请求失败: %w", err)
		}
		req.Header.Set("User-Agent", UserAgent)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		c.jwt.mu.Lock()
		tok := c.jwt.token
		c.jwt.mu.Unlock()
		if tok != "" {
			req.Header.Set("Authorization", "Bearer "+tok)
		}
		return c.httpClient.Do(req)
	}

	resp, err := do()
	if err != nil {
		return fmt.Errorf("asmr: 请求 %s 失败: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		// 401：用缓存凭据重登一次（契约 §6）
		c.jwt.mu.Lock()
		n, p := c.jwt.name, c.jwt.password
		c.jwt.mu.Unlock()
		if n == "" || p == "" {
			return &APIError{Code: resp.StatusCode, URL: c.baseURL.JoinPath(path).String(), Body: "401 无缓存凭据可重登"}
		}
		newTok, lerr := c.doLogin(ctx, n, p)
		if lerr != nil {
			return lerr
		}
		c.jwt.mu.Lock()
		c.jwt.token = newTok
		c.jwt.mu.Unlock()
		c.token = newTok
		resp2, err := do()
		if err != nil {
			return fmt.Errorf("asmr: 重登后重试 %s 失败: %w", path, err)
		}
		defer resp2.Body.Close()
		resp = resp2
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("asmr: 读取响应失败: %w", err)
	}
	if resp.StatusCode >= 400 {
		return &APIError{Code: resp.StatusCode, URL: c.baseURL.JoinPath(path).String(), Body: truncateBody(data)}
	}
	if out != nil && len(data) > 0 {
		if err := json.Unmarshal(data, out); err != nil {
			return fmt.Errorf("asmr: 解析 %s 响应失败: %w", path, err)
		}
	}
	return nil
}

// CloudReviews 读取云端收藏/收听列表（filter=marked|listening），返回归一化 rjId 集合。
// 实测 GET /api/review 200；rjId 从 work.source_id 或顶层 source_id 提取。
func (c *Client) CloudReviews(ctx context.Context, filter string) (map[string]struct{}, error) {
	q := map[string]string{
		"order": "updated_at",
		"sort":  "desc",
		"page":  "1",
		"filter": filter,
	}
	raw := c.baseURL.JoinPath("/review")
	vals := raw.Query()
	for k, v := range q {
		vals.Set(k, v)
	}
	raw.RawQuery = vals.Encode()

	if err := c.limiter.wait(ctx); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, raw.String(), bytes.NewReader(nil))
	if err != nil {
		return nil, fmt.Errorf("asmr: 构造请求失败: %w", err)
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "application/json")
	c.jwt.mu.Lock()
	tok := c.jwt.token
	c.jwt.mu.Unlock()
	if tok != "" {
		req.Header.Set("Authorization", "Bearer "+tok)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("asmr: 请求 /api/review 失败: %w", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("asmr: 读取响应失败: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, &APIError{Code: resp.StatusCode, URL: raw.String(), Body: truncateBody(data)}
	}

	var reviews []rawReview
	if err := json.Unmarshal(data, &reviews); err != nil {
		return nil, fmt.Errorf("asmr: 解析 /api/review 响应失败: %w", err)
	}

	res := make(map[string]struct{}, len(reviews))
	for _, r := range reviews {
		if rj := normalizeReviewRJID(r); rj != "" {
			res[rj] = struct{}{}
		}
	}
	return res, nil
}

// CloudFavorites 云端已收藏 rjId 集合（filter=marked）。
func (c *Client) CloudFavorites(ctx context.Context) (map[string]struct{}, error) {
	return c.CloudReviews(ctx, "marked")
}

// CloudListening 云端收听中 rjId 集合（filter=listening）。
func (c *Client) CloudListening(ctx context.Context) (map[string]struct{}, error) {
	return c.CloudReviews(ctx, "listening")
}

// normalizeReviewRJID 从 review 单条提取归一化 rjId（"RJ"+数字）。
// 兼容 source_id 在顶层或嵌套 work 内两种实测形态。
func normalizeReviewRJID(r rawReview) string {
	cand := r.SourceID
	if cand == "" && r.Work != nil {
		cand = r.Work.SourceID
	}
	cand = strings.TrimSpace(cand)
	cand = strings.TrimPrefix(cand, "RJ")
	cand = strings.TrimPrefix(cand, "rj")
	if cand == "" {
		return ""
	}
	if !strings.ContainsFunc(cand, func(r rune) bool { return r >= '0' && r <= '9' }) {
		return ""
	}
	return "RJ" + cand
}

// SaveReview 写云端收藏/收听状态（契约 §6）。
//
// 端点实测（Wave A）：PUT /api/review（第三方 henntaidesu/asmr.one_download 佐证 review() 用 PUT /api/review）。
// 该端点实测返回 200（body 含 work_id/marked/listening/review/rating 均被接受），
// 但实测发现：以 marked:true 写入后，GET /api/review?filter=marked 列表并未包含该作品——
// 即控制云端"已收藏"列表的字段并非 body 的 marked 布尔（progress 枚举也全部 400 拒绝）。
//
// 因此写同步是否真正生效需以 read-back 校验为准（见 handler.syncCloudFavorite）。
// 本方法仅负责发出写请求；写失败（非 2xx）返回 error，调用方据此 syncedToCloud=false。
//
// TODO(M3.1 编排裁决): 云端 marked/listening 列表的真实写入字段/端点未确认（疑似独立开关或 progress 特定枚举），
// 待编排层确认后对齐 SaveReview body；写失败/未确认不影响本地收藏。
func (c *Client) SaveReview(ctx context.Context, workID int, marked, listening bool) error {
	body := map[string]any{
		"work_id":  workID,
		"marked":   marked,
		"listening": listening,
		"review":   "",
		"rating":   0,
	}
	return c.doPut(ctx, "/review", body, nil)
}

// doPut 发送 JSON PUT；带 Bearer token；遇 401 重登一次后重试（与 doPost 同策略）。
func (c *Client) doPut(ctx context.Context, path string, body any, out any) error {
	return c.doMethod(ctx, http.MethodPut, path, body, out)
}

// doMethod 通用 JSON 写方法（POST/PUT），带 Bearer token 与 401 重登一次。
func (c *Client) doMethod(ctx context.Context, method, path string, body any, out any) error {
	if err := c.limiter.wait(ctx); err != nil {
		return err
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("asmr: 序列化请求体失败: %w", err)
	}

	do := func() (*http.Response, error) {
		req, err := http.NewRequestWithContext(ctx, method, c.baseURL.JoinPath(path).String(), bytes.NewReader(payload))
		if err != nil {
			return nil, fmt.Errorf("asmr: 构造请求失败: %w", err)
		}
		req.Header.Set("User-Agent", UserAgent)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		c.jwt.mu.Lock()
		tok := c.jwt.token
		c.jwt.mu.Unlock()
		if tok != "" {
			req.Header.Set("Authorization", "Bearer "+tok)
		}
		return c.httpClient.Do(req)
	}

	resp, err := do()
	if err != nil {
		return fmt.Errorf("asmr: 请求 %s %s 失败: %w", method, path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		c.jwt.mu.Lock()
		n, p := c.jwt.name, c.jwt.password
		c.jwt.mu.Unlock()
		if n == "" || p == "" {
			return &APIError{Code: resp.StatusCode, URL: c.baseURL.JoinPath(path).String(), Body: "401 无缓存凭据可重登"}
		}
		newTok, lerr := c.doLogin(ctx, n, p)
		if lerr != nil {
			return lerr
		}
		c.jwt.mu.Lock()
		c.jwt.token = newTok
		c.jwt.mu.Unlock()
		c.token = newTok
		resp2, err := do()
		if err != nil {
			return fmt.Errorf("asmr: 重登后重试 %s %s 失败: %w", method, path, err)
		}
		defer resp2.Body.Close()
		resp = resp2
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("asmr: 读取响应失败: %w", err)
	}
	if resp.StatusCode >= 400 {
		return &APIError{Code: resp.StatusCode, URL: c.baseURL.JoinPath(path).String(), Body: truncateBody(data)}
	}
	if out != nil && len(data) > 0 {
		if err := json.Unmarshal(data, out); err != nil {
			return fmt.Errorf("asmr: 解析 %s %s 响应失败: %w", method, path, err)
		}
	}
	return nil
}

// DownloadTrack 通过 client 的代理 http client 下载直链音频（mediaDownloadUrl）。
// 走与搜索域相同的代理/超时配置，但不经 filecache（二进制流）。
func (c *Client) DownloadTrack(ctx context.Context, mediaURL string) (io.ReadCloser, int64, error) {
	if err := c.limiter.wait(ctx); err != nil {
		return nil, 0, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, mediaURL, bytes.NewReader(nil))
	if err != nil {
		return nil, 0, fmt.Errorf("asmr: 构造下载请求失败: %w", err)
	}
	req.Header.Set("User-Agent", UserAgent)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("asmr: 下载 %s 失败: %w", mediaURL, err)
	}
	if resp.StatusCode >= 400 {
		_ = resp.Body.Close()
		return nil, 0, &APIError{Code: resp.StatusCode, URL: mediaURL, Body: "下载直链失败"}
	}
	return resp.Body, resp.ContentLength, nil
}
