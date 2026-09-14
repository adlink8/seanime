package asmr

import "context"

// 推荐域（契约 03.8 §0/§1）。
//
// 地面真值（2026-09-14 编排层实测，勿重复实测）：
//   - POST /recommender/popular body {keyword:" ", page, subtitle:0|1, localSubtitled:0}
//     → WorksResponse 信封（实测 200 / 20 条）
//   - POST /recommender/item-neighbors body {keyword:"", itemId:int, page, subtitle:0|1,
//     localSubtitled:0} → WorksResponse（实测 200；部分作品 works 空属预期，无推荐数据）
//   - POST /recommender/recommend-for-user 需 userId(uuid)，登录响应 {token,user} 未见
//     稳定 uuid 字段 → v1 暂缓（需求文档裁决，不编造参数来源）。
//   - 响应条目与 rawWork 同构 → 全部复用 convertSearchResult（D1，零新转换）。

// subtitleInt bool → 0|1（Yuro 源码实测 body 用整数）。
func subtitleInt(subtitle bool) int {
	if subtitle {
		return 1
	}
	return 0
}

// GetPopularWorks asmr.one 热门作品（发现页热门区块数据源）。
func (c *Client) GetPopularWorks(ctx context.Context, page int, subtitle bool) (Asmr_SearchResult, error) {
	var raw rawSearchResult
	body := map[string]any{
		"keyword":       " ",
		"page":          page,
		"subtitle":      subtitleInt(subtitle),
		"localSubtitled": 0,
	}
	if err := c.doPost(ctx, "/recommender/popular", body, &raw); err != nil {
		return Asmr_SearchResult{}, err
	}
	return convertSearchResult(raw), nil
}

// GetItemNeighbors 相似作品推荐（部分作品可能返回空 works，属预期）。
func (c *Client) GetItemNeighbors(ctx context.Context, itemID int, page int, subtitle bool) (Asmr_SearchResult, error) {
	var raw rawSearchResult
	body := map[string]any{
		"keyword":        "",
		"itemId":         itemID,
		"page":           page,
		"subtitle":       subtitleInt(subtitle),
		"localSubtitled": 0,
	}
	if err := c.doPost(ctx, "/recommender/item-neighbors", body, &raw); err != nil {
		return Asmr_SearchResult{}, err
	}
	return convertSearchResult(raw), nil
}
