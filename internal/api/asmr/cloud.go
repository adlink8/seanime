package asmr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// 云端账号同步（契约 §6）。
//
// 地面真值（编排层实测，勿重复实测；2026-09-14 真实网络复测）：
//   - 登录：POST /api/auth/me {name, password} → JWT（契约记为 token）
//   - 读取：GET /api/review?order=updated_at&sort=desc&page=1&filter=marked|listening 实测 200。
//     **响应是分页信封 {"works": [...], "pagination": {...}}，非裸数组**（2026-09-14 实测，
//     旧版按 []rawReview 解析对真实响应必然失败——已修）。
//   - 单条为**扁平结构**：{"id": 390787, "source_id": "RJ390787", "progress": "listening", ...}，
//     无 work_id / 嵌套 work / user_review 字段（2026-09-14 实测）。
//   - 写端点：PUT /api/review，body 仅 {work_id, progress}（progress 为枚举字符串，非布尔）。
//     **2026-09-14 真实账号实测通过**：PUT marked → 条目移入 marked filter 且 progress 回读一致；
//     恢复 listening 同样生效。
//   - 取消收藏语义（实测坐实）：progress 是唯一状态模型，条目只会在 filter 间移动、不会消失；
//     无 DELETE 端点、无 "clear" 枚举值 → 云端无「取消收藏」，handler favorite=false 放弃同步是终态。
//
// 第三方源码佐证（body 形状两路独立）：
//       https://raw.githubusercontent.com/henntaidesu/asmr.one_download/master/src/asmr_api/works_review.py
//       https://raw.githubusercontent.com/asmroneapp/Yuro/main/lib/data/services/api_service.dart （updateWorkMarkStatus / convertMarkStatusToApi）
//
// 凭据来源：调用方（handler）从 config 注入，本 client 不读文件/环境变量。
// 严禁把凭据/JWT 写入日志、文件、测试或提交（契约硬约束）。

// authResponse /api/auth/me 登录响应（asmr.one 实测返回 token 字段；兼容 access_token）。
type authResponse struct {
	Token       string `json:"token"`
	AccessToken string `json:"access_token"`
}

// rawReview /api/review 列表单条（2026-09-14 实测扁平结构，见文件头地面真值）。
type rawReview struct {
	ID       int    `json:"id"`
	SourceID string `json:"source_id"`
	Progress string `json:"progress"`
}

// reviewEnvelope /api/review 分页信封（2026-09-14 实测：{"works": [...], "pagination": {...}}）。
type reviewEnvelope struct {
	Works      []rawReview `json:"works"`
	Pagination struct {
		TotalCount int `json:"totalCount"`
	} `json:"pagination"`
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

	var env reviewEnvelope
	if err := json.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("asmr: 解析 /api/review 响应失败: %w", err)
	}

	res := make(map[string]struct{}, len(env.Works))
	for _, r := range env.Works {
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
// 实测结构（2026-09-14）：source_id 顶层必有（如 "RJ390787"）；缺失时回退用数字 id 拼 "RJ{id}"。
func normalizeReviewRJID(r rawReview) string {
	cand := strings.TrimSpace(r.SourceID)
	if cand == "" {
		if r.ID <= 0 {
			return ""
		}
		cand = strconv.Itoa(r.ID)
	}
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

// Progress* 云端收藏/收听状态枚举（asmr.one review.progress 取值，非布尔）。
// 证据：
//   - https://raw.githubusercontent.com/henntaidesu/asmr.one_download/master/src/asmr_api/works_review.py
//   - https://raw.githubusercontent.com/asmroneapp/Yuro/main/lib/data/services/api_service.dart （updateWorkMarkStatus / convertMarkStatusToApi）
// 禁止在调用处散落裸字符串（契约 D6）。
const (
	ProgressMarked    = "marked"    // 想听
	ProgressListening = "listening" // 在听
	ProgressListened  = "listened"  // 听过
	ProgressReplay    = "replay"    // 重听
	ProgressPostponed = "postponed" // 搁置
)

// SaveReview 写云端收藏/收听状态（契约 §0.3 / D5 / Wave C）。
//
// 端点是 PUT /api/review，HTTP 方法与路径均与旧实现一致（本就没错），错的是旧实现发出的 body 形状。
// 正确 body 仅两个字段：{work_id:int, progress:string}，progress 为枚举字符串（见 Progress* 常量），
// 绝不能用布尔 marked/listening 表示。
//
// 证据（两路独立第三方源码，均发 PUT /api/review + body {work_id, progress}）：
//   - https://raw.githubusercontent.com/henntaidesu/asmr.one_download/master/src/asmr_api/works_review.py
//   - https://raw.githubusercontent.com/asmroneapp/Yuro/main/lib/data/services/api_service.dart （updateWorkMarkStatus / convertMarkStatusToApi）
//
// 取消收藏（favorite=false）语义：经三路调研（asmr.one_download 脚本、asmr.one 官方前端 henntaidesu/asmr.one、
// Yuro 客户端）均未找到"清除/取消"的 progress 值，也未找到 DELETE /api/review 之类端点。
// progress 枚举仅 5 个值（marked/listening/listened/replay/postponed），无 "unmark"/"clear"。
// 故 handler 在 favorite=false 时显式放弃云端写同步（syncedToCloud=false），见其内部 TODO。
//
// 写失败（非 2xx）返回 error，调用方据此 syncedToCloud=false。
func (c *Client) SaveReview(ctx context.Context, workID int, progress string) error {
	body := map[string]any{
		"work_id":  workID,
		"progress": progress,
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
