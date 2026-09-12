// Package bangumi 提供 Bangumi API v0 (https://api.bgm.tv) 的 Go 客户端。
//
// 关键设计（依据 .planning/phases/01-foundation/01-CONTRACT.md 硬约束）：
//   - 节流：全局 0.5s ticker 限流器（每个 interval 放行一个请求，首个请求立即放行），
//     社区安全线 ≈120 req/min，官方无 rate limit 文档；可通过 WithThrottleInterval 调整（<=0 关闭）
//   - 退避：收到 429 时优先读取 Retry-After 头（秒），缺失则按 base*2^n 指数退避，最多重试 3 次
//   - 缓存：GET 请求成功响应落 filecache（TTL 24h），键为 "方法 路径?query"；
//     POST/PATCH 永不缓存（写操作必须打到服务端）
//   - 鉴权：Bearer token 由构造函数注入，不读文件、不读环境变量做默认值
package bangumi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"seanime/internal/util/filecache"
)

const (
	// DefaultBaseURL Bangumi API 官方地址
	DefaultBaseURL = "https://api.bgm.tv"

	// UserAgent 官方要求的 UA 头（契约硬约束，所有请求必带）
	UserAgent = "li/adlink8 (seanime-personal-fork)"

	defaultTimeout          = 15 * time.Second // 出站超时（契约硬约束）
	defaultThrottleInterval = 500 * time.Millisecond
	defaultRetryBaseDelay   = 1 * time.Second
	maxRetries              = 3 // 429 最多重试次数（契约硬约束）

	cacheBucketName = "bangumi_api"
	cacheTTL        = 24 * time.Hour
)

type Client struct {
	token          string
	baseURL        *url.URL
	httpClient     *http.Client
	logger         *zerolog.Logger
	cache          *filecache.Cacher // 为 nil 时禁用缓存
	limiter        *tickerLimiter
	retryBaseDelay time.Duration
	maxRetries     int
}

type options struct {
	httpClient       *http.Client
	logger           *zerolog.Logger
	cache            *filecache.Cacher
	baseURL          string
	throttleInterval time.Duration
	retryBaseDelay   time.Duration
}

type Option func(*options)

// WithHTTPClient 注入自定义 http.Client（测试用 httptest 或需要自定义传输层时）
func WithHTTPClient(hc *http.Client) Option {
	return func(o *options) { o.httpClient = hc }
}

// WithLogger 注入 zerolog logger，默认 Nop
func WithLogger(l *zerolog.Logger) Option {
	return func(o *options) { o.logger = l }
}

// WithFileCache 启用 GET 响应缓存（TTL 24h）
func WithFileCache(c *filecache.Cacher) Option {
	return func(o *options) { o.cache = c }
}

// WithBaseURL 覆盖 API 地址（单测指向 httptest mock server）
func WithBaseURL(u string) Option {
	return func(o *options) { o.baseURL = u }
}

// WithThrottleInterval 覆盖节流间隔，<=0 表示关闭节流（仅测试用，生产保持默认 0.5s）
func WithThrottleInterval(d time.Duration) Option {
	return func(o *options) { o.throttleInterval = d }
}

// WithRetryBaseDelay 覆盖 429 退避基准延迟（生产默认 1s，测试调短加速）
func WithRetryBaseDelay(d time.Duration) Option {
	return func(o *options) { o.retryBaseDelay = d }
}

// New 创建 Client。token 由调用方注入（个人令牌，scope: write:collection），
// token 为空时请求不带 Authorization 头（公开端点仍可用，NSFW 内容会被 404）。
func New(token string, opts ...Option) *Client {
	o := options{
		throttleInterval: defaultThrottleInterval,
		retryBaseDelay:   defaultRetryBaseDelay,
	}
	for _, opt := range opts {
		opt(&o)
	}

	baseURLStr := o.baseURL
	if baseURLStr == "" {
		baseURLStr = DefaultBaseURL
	}
	baseURL, err := url.Parse(baseURLStr)
	if err != nil {
		// 构造函数不返回错误，非法地址回退官方地址
		baseURL, _ = url.Parse(DefaultBaseURL)
	}

	logger := zerolog.Nop()
	if o.logger != nil {
		logger = *o.logger
	}

	httpClient := o.httpClient
	if httpClient == nil {
		// 不设置自定义 Transport：http.DefaultTransport 默认走 ProxyFromEnvironment，
		// 与契约"代理交给用户环境（Clash）"一致
		httpClient = &http.Client{Timeout: defaultTimeout}
	}

	return &Client{
		token:          token,
		baseURL:        baseURL,
		httpClient:     httpClient,
		logger:         &logger,
		cache:          o.cache,
		limiter:        newTickerLimiter(o.throttleInterval),
		retryBaseDelay: o.retryBaseDelay,
		maxRetries:     maxRetries,
	}
}

// Close 停止后台节流 ticker，进程长期运行时避免 goroutine 泄漏。
// Close 后 Client 不可再用。
func (c *Client) Close() {
	c.limiter.stop()
}

// GetMe 获取当前 token 对应用户（GET /v0/me），用于 token 有效性校验。
func (c *Client) GetMe(ctx context.Context) (*User, error) {
	var out User
	if err := c.doGet(ctx, "/v0/me", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// doGet 发送 GET 请求，带缓存层（命中则不触网）。
func (c *Client) doGet(ctx context.Context, path string, query url.Values, out any) error {
	key := c.cacheKey(http.MethodGet, path, query)
	if c.cache != nil {
		var cached string
		if ok, err := c.cache.Get(filecache.NewBucket(cacheBucketName, cacheTTL), key, &cached); err == nil && ok {
			c.logger.Debug().Str("key", key).Msg("bangumi: 缓存命中")
			return json.Unmarshal([]byte(cached), out)
		}
	}

	respBody, err := c.doRequest(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return err
	}

	// 仅成功响应写缓存（doRequest 对 >=400 已返回错误）
	if c.cache != nil {
		_ = c.cache.Set(filecache.NewBucket(cacheBucketName, cacheTTL), key, string(respBody))
	}

	return json.Unmarshal(respBody, out)
}

// doWrite 发送 POST/PATCH 请求，永不缓存。
func (c *Client) doWrite(ctx context.Context, method, path string, query url.Values, body any, out any) error {
	var payload []byte
	if body != nil {
		var err error
		if payload, err = json.Marshal(body); err != nil {
			return fmt.Errorf("bangumi: 序列化请求体失败: %w", err)
		}
	}

	respBody, err := c.doRequest(ctx, method, path, query, payload)
	if err != nil {
		return err
	}

	if out != nil && len(respBody) > 0 {
		return json.Unmarshal(respBody, out)
	}
	return nil
}

// doRequest 核心请求管道：节流 -> 发送 -> 429 退避重试（<=3 次）。
func (c *Client) doRequest(ctx context.Context, method, path string, query url.Values, body []byte) ([]byte, error) {
	u := c.baseURL.JoinPath(path)
	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}

	for attempt := 0; ; attempt++ {
		// 节流：每个请求（含重试）都先等一个令牌
		if err := c.limiter.wait(ctx); err != nil {
			return nil, err
		}

		var req *http.Request
		var err error
		if body != nil {
			req, err = http.NewRequestWithContext(ctx, method, u.String(), bytes.NewReader(body))
		} else {
			req, err = http.NewRequestWithContext(ctx, method, u.String(), nil)
		}
		if err != nil {
			return nil, fmt.Errorf("bangumi: 构造请求失败: %w", err)
		}

		req.Header.Set("User-Agent", UserAgent)
		req.Header.Set("Accept", "application/json")
		if c.token != "" {
			req.Header.Set("Authorization", "Bearer "+c.token)
		}
		if body != nil {
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			// 网络层错误不重试（超时由 http.Client.Timeout 兜底）
			return nil, fmt.Errorf("bangumi: 请求 %s %s 失败: %w", method, u.String(), err)
		}

		data, readErr := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if readErr != nil {
			return nil, fmt.Errorf("bangumi: 读取响应失败: %w", readErr)
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			apiErr := &APIError{
				Code:   resp.StatusCode,
				Method: method,
				URL:    u.String(),
				Body:   truncateBody(data),
			}
			if attempt < c.maxRetries {
				delay := c.backoffDelay(resp.Header.Get("Retry-After"), attempt)
				c.logger.Warn().
					Int("attempt", attempt+1).
					Dur("delay", delay).
					Msg("bangumi: 收到 429，退避后重试")
				if sleepErr := sleepCtx(ctx, delay); sleepErr != nil {
					return nil, sleepErr
				}
				continue
			}
			return nil, apiErr
		}

		if resp.StatusCode >= 400 {
			return nil, &APIError{
				Code:   resp.StatusCode,
				Method: method,
				URL:    u.String(),
				Body:   truncateBody(data),
			}
		}

		return data, nil
	}
}

// backoffDelay 计算退避时长：优先使用服务端 Retry-After（秒，0 或非法值忽略），
// 否则按 base * 2^attempt 指数增长（1x, 2x, 4x...）。
func (c *Client) backoffDelay(retryAfter string, attempt int) time.Duration {
	if sec, err := strconv.Atoi(strings.TrimSpace(retryAfter)); err == nil && sec > 0 {
		return time.Duration(sec) * time.Second
	}
	return c.retryBaseDelay << attempt
}

// cacheKey 缓存键 = 方法 + 请求路径 + 规范化 query。
// 同一端点不同参数互不污染；写操作不经过此函数（不缓存）。
func (c *Client) cacheKey(method, path string, query url.Values) string {
	if len(query) == 0 {
		return method + " " + path
	}
	return method + " " + path + "?" + query.Encode()
}

func truncateBody(data []byte) string {
	const maxLen = 512
	s := strings.TrimSpace(string(data))
	if len(s) > maxLen {
		return s[:maxLen] + "..."
	}
	return s
}

// sleepCtx 可被 context 取消的 sleep
func sleepCtx(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

// tickerLimiter 基于 ticker 的全局限流器：首个请求立即放行，
// 之后每个 interval 放行一个令牌，保证请求间隔 >= interval。
type tickerLimiter struct {
	disabled bool
	tokens   chan struct{}
	ticker   *time.Ticker
}

func newTickerLimiter(interval time.Duration) *tickerLimiter {
	l := &tickerLimiter{tokens: make(chan struct{}, 1)}
	if interval <= 0 {
		l.disabled = true
		return l
	}
	l.ticker = time.NewTicker(interval)
	l.tokens <- struct{}{} // 预填一个令牌，首个请求无延迟
	go func() {
		for range l.ticker.C {
			select {
			case l.tokens <- struct{}{}:
			default: // 令牌满则丢弃，保持恒定速率
			}
		}
	}()
	return l
}

func (l *tickerLimiter) wait(ctx context.Context) error {
	if l.disabled {
		return nil
	}
	select {
	case <-l.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (l *tickerLimiter) stop() {
	if l.ticker != nil {
		l.ticker.Stop()
	}
}
