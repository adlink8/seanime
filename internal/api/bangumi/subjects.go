package bangumi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// GetSubject 获取条目详情（GET /v0/subjects/{id}）。
// 含 name（日文）/name_cn（中文）/summary（中文简介）；服务端本身有 300s 缓存，客户端再加 24h filecache。
func (c *Client) GetSubject(ctx context.Context, subjectID int) (*Subject, error) {
	var out Subject
	if err := c.doGet(ctx, fmt.Sprintf("/v0/subjects/%d", subjectID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EpisodesFilter 章节列表过滤参数
type EpisodesFilter struct {
	// Type EpType 过滤：0本篇 1SP 2OP 3ED 4PV 5MAD 6其他；nil 返回全部
	Type *int
}

// GetSubjectEpisodes 获取章节列表（GET /v0/subjects/{id}/episodes），
// name_cn 为中文章节名。
func (c *Client) GetSubjectEpisodes(ctx context.Context, subjectID int, filter *EpisodesFilter) (*EpisodesResult, error) {
	query := url.Values{}
	if filter != nil && filter.Type != nil {
		query.Set("type", strconv.Itoa(*filter.Type))
	}

	var out EpisodesResult
	if err := c.doGet(ctx, fmt.Sprintf("/v0/subjects/%d/episodes", subjectID), query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
