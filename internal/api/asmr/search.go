package asmr

import (
	"context"
	"net/url"
	"strconv"
)

// 实测地面真值（2026-09-13，走代理 curl 实测 api.asmr.one）：
//
//   - 搜索：GET /api/search/{page}?keyword=...&order=...&subtitle=...&pageSize=...（无需鉴权）
//     注意 order/subtitle 是 query 参数，不是路径段：/api/search/dd/1 返回 404 "Cannot GET"。
//   - order 合法枚举（express-validator 白名单实测）：空（服务端默认）、create_date、dl_count、
//     price、release、id、rating；dc/publish_date/rate 等一律 400 {"error":"order: Invalid value"}。
//     Phase 3.9b D9：order=rating 实测 200 且默认降序（评分最高在前）；order=dc 实测 400，故前端
//     SORTING_ASMR 的「评分最高」值由 dc 改为 rating（见 mapOrder）。
//   - subtitle 合法值：0（全部）、1（有字幕）；jp/2/3 等一律 400。API 不支持按字幕语言过滤。
//   - pageSize 生效（默认 20），响应 pagination.pageSize 回显。
//   - sort 参数（契约 03.2c / D9，2026-09-14 经代理实测）：合法；order=release 时服务端**默认即 desc**
//     （最新在前），sort=desc 同默认，sort=asc 生效（最旧在前）。SearchParams 故不设 Sort 字段，
//     tracker 依赖默认 desc 即可。
//   - 响应形状：{"works":[...], "pagination":{"currentPage","pageSize","totalCount"}}。
//   - popular：GET /api/popular 实测 401 {"error":"No authorization token was found"}（需登录态）；
//     /api/search/popular/{page} 不存在（404）。Popular() 以无关键词默认搜索（page 1）近似替代，
//     服务端默认排序实测返回热门/最新混合结果，属可接受近似。
//   - workInfo：GET /api/workInfo/{id}，{id} 为数字 work id（"RJ..." 返回 404 {"error":"未找到该音声"}）。
//   - tracks：GET /api/tracks/{id}?v=2，根为 JSON 数组，folder 节点用 children 嵌套，
//     叶子含 mediaStreamUrl/mediaDownloadUrl，type ∈ {folder, audio, text, image}。

// SearchParams 搜索参数（契约第 4 节 POST /api/v1/asmr/search 请求体的后端视角）
type SearchParams struct {
	Keyword string
	// Order 契约前端枚举："dd"|"dl"|"rating"|"publish_date"，空走默认。
	// 映射到实测合法枚举（见 mapOrder）。注意 "dc" 为历史别名，同映射为 rating。
	Order string
	// Page 从 1 开始
	Page int
	// PerPage 默认 20（服务端 pageSize）
	PerPage int
	// Subtitle 契约前端枚举："none"|"jp"|"zh"，空不过滤。
	// API 仅支持 0/1（全部/有字幕），"jp"/"zh" 均映射为 1（不区分语言的近似），
	// "none" 视为不过滤（API 无法按"无字幕"过滤）。
	Subtitle string
}

// mapOrder 契约 order 枚举 → asmr.one 实测合法枚举。
// dd（日期降序）→ create_date（最新入库在前）；dl（下载量）→ dl_count；
// dc（历史别名，前端已改名为 rating）→ rating（评分最高，服务端默认降序）；
// rating → rating（评分最高）；publish_date → release；其余/空 → 不传（服务端默认）。
// 注意：asmr.one search/works 端点 order=dc 实测 400（契约 03.9b D9），故 dc 不再映射为 dl_count。
func mapOrder(o string) string {
	switch o {
	case "dd":
		return "create_date"
	case "dl":
		return "dl_count"
	case "dc", "rating":
		return "rating"
	case "publish_date":
		return "release"
	default:
		return ""
	}
}

// mapSubtitle 契约 subtitle 枚举 → asmr.one 实测合法值
func mapSubtitle(s string) string {
	switch s {
	case "jp", "zh":
		return "1"
	default:
		return ""
	}
}

// Search 搜索音声作品（GET /api/search/{page}）
func (c *Client) Search(ctx context.Context, p SearchParams) (*Asmr_SearchResult, error) {
	page := p.Page
	if page < 1 {
		page = 1
	}
	q := url.Values{}
	if p.Keyword != "" {
		q.Set("keyword", p.Keyword)
	}
	if o := mapOrder(p.Order); o != "" {
		q.Set("order", o)
	}
	if s := mapSubtitle(p.Subtitle); s != "" {
		q.Set("subtitle", s)
	}
	if p.PerPage > 0 {
		q.Set("pageSize", strconv.Itoa(p.PerPage))
	}

	var raw rawSearchResult
	if err := c.doGet(ctx, "/search/"+strconv.Itoa(page), q, &raw); err != nil {
		return nil, err
	}
	out := convertSearchResult(raw)
	return &out, nil
}

// Popular 获取热门音声作品。
// 近似实现：/api/popular 需鉴权（实测 401），改用无关键词默认搜索第 1 页，
// 结果为服务端默认排序的热门/最新混合，与"热门"语义近似。
func (c *Client) Popular(ctx context.Context, perPage int) (*Asmr_SearchResult, error) {
	return c.Search(ctx, SearchParams{Page: 1, PerPage: perPage})
}
