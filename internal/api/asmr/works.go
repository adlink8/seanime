package asmr

import (
	"context"
	"net/url"
	"strconv"
)

// 全量作品浏览域（契约 03.8 §0/§1）。
//
// 地面真值（2026-09-14 实测）：GET /api/works?page=&subtitle=0&order=&sort=
// → WorksResponse 信封（实测 200 / 20 条；Yuro 源码默认 order=create_date, sort=desc）。
// 条目与 rawWork 同构 → 复用 convertSearchResult。tracker 的「最新作」拉取可复用本方法
// （order=release&sort=desc 与 search.go 的 order=release 默认方向结论一致，见 search.go 注释）。

// ListWorks 按排序浏览 asmr.one 全量作品。
func (c *Client) ListWorks(ctx context.Context, order, sort string, page int, subtitle bool) (Asmr_SearchResult, error) {
	var raw rawSearchResult
	q := url.Values{
		"page":     {strconv.Itoa(page)},
		"subtitle": {strconv.Itoa(subtitleInt(subtitle))},
		"order":    {order},
		"sort":     {sort},
	}
	if err := c.doGet(ctx, "/works", q, &raw); err != nil {
		return Asmr_SearchResult{}, err
	}
	return convertSearchResult(raw), nil
}
