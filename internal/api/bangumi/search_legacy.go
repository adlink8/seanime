package bangumi

import (
	"context"
	"net/url"
	"strconv"
)

// search_legacy.go —— legacy 检索通路（GET /search/subject/{keyword}）。
//
// 契约依据（03.6b-CONTRACT.md §0.1 / §2）：
//   - v0 `POST /v0/search/subjects` 对 CJK 关键词恒 total=0，legacy 通路是中文标题唯一可行通路；
//   - legacy 响应结构为 `{results, list: []}`，条目字段与 v0 subject 不同：
//     `id` / `type`(整数) / `name` / `name_cn` / `images` / `rating{total,count,score}` /
//     `air_date` / `eps` / `rank` / `summary` / `collection`；
//   - **无 `platform` 字段、无 `tags` 字段**，且不支持任何 filter / sort；
//   - 分页用 `start` + `max_results`（`page` 参数无效，实测 page=1/2/3 返回相同 ids）。
//   - `max_results` 非严格保证：请求 20 可能只返回 17 条。
//
// 本文件复用 Client 既有的节流 / filecache / 显式代理（WithProxyURL）/ UA 机制：
// 走 doGet（GET → 可缓存），与 v0 通路共用同一 limiter 与 http.Client。

// LegacySubject legacy 检索结果条目。
// 与 v0 Subject 刻意分开定义：字段集更窄（无 platform/tags/infobox/nsfw），
// 复用会造成「字段存在但恒为空」的误导。
type LegacySubject struct {
	ID         int            `json:"id"`
	Type       int            `json:"type"` // 1书籍 2动画 3音乐 4游戏 6三次元
	Name       string         `json:"name"`
	NameCN     string         `json:"name_cn"` // 中文名（中文搜索的核心收益）
	Images     SubjectImages  `json:"images"`
	Rating     *SubjectRating `json:"rating"` // rating.score 为 0–10 制
	AirDate    string         `json:"air_date"`
	Eps        int            `json:"eps"`
	Rank       int            `json:"rank"`
	Summary    string         `json:"summary"`
	Collection map[string]int `json:"collection,omitempty"`
}

// LegacySearchResult legacy 检索响应：{results, list[]}。
// Results 是命中总数（分页判据），List 是当前页条目（条数可能少于请求的 max_results）。
type LegacySearchResult struct {
	Results int             `json:"results"`
	List    []LegacySubject `json:"list"`
}

// SearchSubjectsLegacy 调用 legacy 检索端点（GET /search/subject/{keyword}）。
//
// subjectType：动画=SubjectAnime(2)、书籍=SubjectBook(1)。
// start：偏移量（契约 §2 规定 start = 20 * (page-1)）。
// limit：请求的 max_results（契约 §2 规定请求 20）。
//
// 注意：该端点**不支持 filter / sort**，调用方若传了额外筛选条件，
// 必须在本地对结果集过滤/排序，或按契约 §0.1 显式降级（禁止静默忽略）。
func (c *Client) SearchSubjectsLegacy(ctx context.Context, keyword string, subjectType int, start, limit int) (*LegacySearchResult, error) {
	query := url.Values{}
	query.Set("type", strconv.Itoa(subjectType))
	if start > 0 {
		query.Set("start", strconv.Itoa(start))
	}
	if limit > 0 {
		query.Set("max_results", strconv.Itoa(limit))
	}

	// 关键词是路径段：用 PathEscape 保证 CJK / 空格 / '/' 等不破坏路径结构。
	// JoinPath 会保留已转义的 RawPath，故 URL 形态稳定（见 client.go doRequest）。
	path := "/search/subject/" + url.PathEscape(keyword)

	var out LegacySearchResult
	if err := c.doGet(ctx, path, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
