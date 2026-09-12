package bangumi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// ListCollectionsOpts GET /v0/users/-/collections 分页参数
type ListCollectionsOpts struct {
	Limit       int // 1-50
	Offset      int
	SubjectType *int // 1书籍 2动画 3音乐 4游戏 6三次元；nil 返回全部
}

// GetUserCollections 获取当前用户的收藏列表（GET /v0/users/-/collections），
// "-" 为本人；私有收藏需本人 token。
func (c *Client) GetUserCollections(ctx context.Context, opts ListCollectionsOpts) (*UserCollectionsResult, error) {
	query := url.Values{}
	if opts.Limit > 0 {
		query.Set("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Offset > 0 {
		query.Set("offset", strconv.Itoa(opts.Offset))
	}
	if opts.SubjectType != nil {
		query.Set("subject_type", strconv.Itoa(*opts.SubjectType))
	}

	var out UserCollectionsResult
	if err := c.doGet(ctx, "/v0/users/-/collections", query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpsertCollection 创建或更新收藏（POST /v0/users/-/collections/{subject_id}），
// 收藏不存在则创建，已存在则整体更新。
func (c *Client) UpsertCollection(ctx context.Context, subjectID int, body CollectionUpsertBody) error {
	return c.doWrite(ctx, http.MethodPost,
		fmt.Sprintf("/v0/users/-/collections/%d", subjectID), nil, body, nil)
}

// PatchCollection 局部更新收藏（PATCH /v0/users/-/collections/{subject_id}），
// 只更新提供的字段。
func (c *Client) PatchCollection(ctx context.Context, subjectID int, body CollectionUpsertBody) error {
	return c.doWrite(ctx, http.MethodPatch,
		fmt.Sprintf("/v0/users/-/collections/%d", subjectID), nil, body, nil)
}

// DeleteCollection 删除收藏（DELETE /v0/users/-/collections/{subject_id}），
// 成功时服务端返回 204 空响应体。
// （Wave B additive：平台层 DeleteEntry 使用）
func (c *Client) DeleteCollection(ctx context.Context, subjectID int) error {
	return c.doWrite(ctx, http.MethodDelete,
		fmt.Sprintf("/v0/users/-/collections/%d", subjectID), nil, nil, nil)
}

// GetEpisodeCollections 查询某条目的章节观看进度
// （GET /v0/users/-/collections/{subject_id}/episodes）。
func (c *Client) GetEpisodeCollections(ctx context.Context, subjectID int) (*EpisodeCollectionsResult, error) {
	var out EpisodeCollectionsResult
	if err := c.doGet(ctx, fmt.Sprintf("/v0/users/-/collections/%d/episodes", subjectID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// PatchEpisodeCollections 批量更新章节进度
// （PATCH /v0/users/-/collections/{subject_id}/episodes，请求体 {episode_id: [1,2,8], type}），
// 服务端自动重算完成度。
func (c *Client) PatchEpisodeCollections(ctx context.Context, subjectID int, episodeIDs []int, progressType string) error {
	body := EpisodeProgressBody{EpisodeIDs: episodeIDs, Type: progressType}
	return c.doWrite(ctx, http.MethodPatch,
		fmt.Sprintf("/v0/users/-/collections/%d/episodes", subjectID), nil, body, nil)
}
