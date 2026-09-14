// Package asmr 提供 asmr.one API (https://api.asmr.one/api) 的 Go 客户端（音声作品搜索域）。
//
// 关键设计（依据 .planning/phases/02-library-reanchor/02.5-CONTRACT.md + 2026-09-13 实测）：
//   - 鉴权：搜索/详情/音轨端点均无需 token（实测 /api/popular 需要，已用无关键词搜索替代）；
//     New(token) 预留 token 参数，非空时携带 Authorization: Bearer（供未来解锁 popular 等端点）
//   - 代理：api.asmr.one 直连 TLS 被断（被墙），必须显式注入代理（config server.proxyurl），
//     实现方式与 bangumi.WithProxyURL 相同
//   - 节流：全局 300ms ticker 限流器（asmr.one 无官方 rate limit 文档，比 bangumi 放宽）；
//     无 429 退避重试（实测未遇到 429，保持简单）
//   - 缓存：无 filecache（asmr 数据时效性强、体量大），缓存策略交给 handler 层
//   - 礼仪：UA 同 bangumi client
package asmr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog"
	"seanime/internal/util/filecache"
)

const (
	// DefaultBaseURL asmr.one API 地址
	DefaultBaseURL = "https://api.asmr.one/api"

	// UserAgent 与 bangumi client 一致的 UA 礼仪
	UserAgent = "li/adlink8 (seanime-personal-fork)"

	defaultTimeout          = 15 * time.Second
	defaultThrottleInterval = 300 * time.Millisecond

	// cacheBucketName / cacheTTL 照抄 bangumi：GET 响应落 filecache 24h（契约 §1 client GET filecache 24h）
	cacheBucketName = "asmr_api"
	cacheTTL        = 24 * time.Hour
)

type Client struct {
	token      string
	baseURL    *url.URL
	httpClient *http.Client
	logger     *zerolog.Logger
	limiter    *tickerLimiter
	cache      *filecache.Cacher // 为 nil 时禁用缓存
	jwt        jwtCache          // 登录态内存缓存
}

type options struct {
	httpClient       *http.Client
	logger           *zerolog.Logger
	cache            *filecache.Cacher
	baseURL          string
	throttleInterval time.Duration
	proxyURL         string
}

type Option func(*options)

// WithHTTPClient 注入自定义 http.Client（测试用 httptest）
func WithHTTPClient(hc *http.Client) Option {
	return func(o *options) { o.httpClient = hc }
}

// WithLogger 注入 zerolog logger，默认 Nop
func WithLogger(l *zerolog.Logger) Option {
	return func(o *options) { o.logger = l }
}

// WithBaseURL 覆盖 API 地址（单测指向 httptest mock server）
func WithBaseURL(u string) Option {
	return func(o *options) { o.baseURL = u }
}

// WithThrottleInterval 覆盖节流间隔，<=0 关闭节流（仅测试用，生产保持默认 300ms）
func WithThrottleInterval(d time.Duration) Option {
	return func(o *options) { o.throttleInterval = d }
}

// WithProxyURL 显式设置出站代理（如 Clash "http://127.0.0.1:7897"）。
// api.asmr.one 直连被墙，代理配置必须随 config.toml 显式下发。
func WithProxyURL(raw string) Option {
	return func(o *options) { o.proxyURL = raw }
}

// WithFileCache 启用 GET 响应缓存（TTL 24h，键含 URL）。照抄 bangumi client 模式。
func WithFileCache(c *filecache.Cacher) Option {
	return func(o *options) { o.cache = c }
}

// New 创建 Client。token 可为空串（搜索域无需鉴权）；非空时请求带 Authorization 头。
func New(token string, opts ...Option) *Client {
	o := options{
		throttleInterval: defaultThrottleInterval,
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
		if o.proxyURL != "" {
			if pu, perr := url.Parse(o.proxyURL); perr == nil && pu.Scheme != "" && pu.Host != "" {
				httpClient = &http.Client{
					Timeout:   defaultTimeout,
					Transport: &http.Transport{Proxy: http.ProxyURL(pu)},
				}
			} else {
				logger.Warn().Str("proxyURL", o.proxyURL).Msg("asmr: 非法代理地址，忽略并直连")
				httpClient = &http.Client{Timeout: defaultTimeout}
			}
		} else {
			httpClient = &http.Client{Timeout: defaultTimeout}
		}
	}

	return &Client{
		token:      token,
		baseURL:    baseURL,
		httpClient: httpClient,
		logger:     &logger,
		limiter:    newTickerLimiter(o.throttleInterval),
		cache:      o.cache,
	}
}

// Close 停止后台节流 ticker，避免 goroutine 泄漏。Close 后 Client 不可再用。
func (c *Client) Close() {
	c.limiter.stop()
}

// doGet 发送 GET 请求并解析 JSON 响应。带 filecache 层（命中则不触网，照抄 bangumi）。
func (c *Client) doGet(ctx context.Context, path string, query url.Values, out any) error {
	return c.doGetWithCache(ctx, path, query, out, true)
}

// doGetNoCache 直连请求，读写均绕过 filecache。
// 契约 03.9c：探索页「最新音声」时效性强，24h 缓存会让最新入库长时间不刷新，该路径直连。
func (c *Client) doGetNoCache(ctx context.Context, path string, query url.Values, out any) error {
	return c.doGetWithCache(ctx, path, query, out, false)
}

func (c *Client) doGetWithCache(ctx context.Context, path string, query url.Values, out any, useCache bool) error {
	u := c.baseURL.JoinPath(path)
	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}

	key := c.cacheKey(http.MethodGet, u.String())

	if useCache && c.cache != nil {
		var cached string
		if ok, err := c.cache.Get(filecache.NewBucket(cacheBucketName, cacheTTL), key, &cached); err == nil && ok {
			c.logger.Debug().Str("key", key).Msg("asmr: 缓存命中")
			return json.Unmarshal([]byte(cached), out)
		}
	}

	if err := c.limiter.wait(ctx); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), bytes.NewReader(nil))
	if err != nil {
		return fmt.Errorf("asmr: 构造请求失败: %w", err)
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("asmr: 请求 %s 失败: %w", u.String(), err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("asmr: 读取响应失败: %w", err)
	}

	if resp.StatusCode >= 400 {
		return &APIError{
			Code: resp.StatusCode,
			URL:  u.String(),
			Body: truncateBody(data),
		}
	}

	if len(data) == 0 {
		return fmt.Errorf("asmr: %s 返回空响应", u.String())
	}

	// 仅成功响应写缓存
	if useCache && c.cache != nil {
		_ = c.cache.Set(filecache.NewBucket(cacheBucketName, cacheTTL), key, string(data))
	}

	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("asmr: 解析 %s 响应失败: %w", u.String(), err)
	}
	return nil
}

// cacheKey 生成缓存键（方法 + 完整 URL，含 query）。
func (c *Client) cacheKey(method, fullURL string) string {
	return method + " " + fullURL
}

// APIError asmr.one 非 2xx 响应错误
type APIError struct {
	Code int
	URL  string
	Body string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("asmr: %s 返回 %d: %s", e.URL, e.Code, e.Body)
}

func truncateBody(data []byte) string {
	const maxLen = 512
	s := string(bytes.TrimSpace(data))
	if len(s) > maxLen {
		return s[:maxLen] + "..."
	}
	return s
}

// tickerLimiter 基于 ticker 的全局限流器（与 bangumi client 同构，间隔独立）
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
			default:
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
