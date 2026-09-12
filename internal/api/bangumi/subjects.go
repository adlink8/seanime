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

// RelatedSubject 相关条目（GET /v0/subjects/{id}/subjects 元素）。
// Type 为关系类型：1前传 2续集 3主线故事 4支线故事 5角色 6其他。
type RelatedSubject struct {
	ID          int            `json:"id"`          // 关系记录 ID（非条目 ID）
	Type        int            `json:"type"`        // 关系类型（见上）
	Name        string         `json:"name"`        // 条目原名
	NameCN      string         `json:"name_cn"`     // 条目中文名
	SubjectID   int            `json:"subject_id"`  // 关联条目 ID（真正的锚点 ID）
	SubjectType int            `json:"subject_type"`// 关联条目类型（1书籍 2动画…）
	Images      *SubjectImages `json:"images,omitempty"`
}

// GetRelatedSubjects 获取条目的相关条目（GET /v0/subjects/{id}/subjects），
// 用于建立条目关系树（前传/续作等）。
func (c *Client) GetRelatedSubjects(ctx context.Context, subjectID int) ([]RelatedSubject, error) {
	var out []RelatedSubject
	if err := c.doGet(ctx, fmt.Sprintf("/v0/subjects/%d/subjects", subjectID), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
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
