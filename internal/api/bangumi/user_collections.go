package bangumi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// UserCollectionsOpts GET /v0/users/{username}/collections 查询参数。
// 区分两个维度：
//   - Type：收藏状态（1想看 2看过 3在看 4搁置 5抛弃），对应 query "type"
//   - SubjectType：条目类型（1书籍 2动画 3音乐 4游戏 6三次元），对应 query "subject_type"
//
// 换锚后取库的典型用法：动画库 = SubjectType=2，书籍/漫画库 = SubjectType=1。
type UserCollectionsOpts struct {
	// Type 收藏状态过滤；nil 返回全部状态
	Type *int
	// SubjectType 条目类型过滤；nil 返回全部类型
	SubjectType *int
	Limit       int // 1-50，服务端默认 30
	Offset      int
}

// GetUserCollectionsByUser 获取指定用户的收藏列表
// （GET /v0/users/{username}/collections）。
// username 可传数字 uid 或用户名；传 "-" 查询当前 token 用户本人
// （私有收藏必须本人 token）。
func (c *Client) GetUserCollectionsByUser(ctx context.Context, username string, opts UserCollectionsOpts) (*UserCollectionsResult, error) {
	if username == "" {
		username = "-"
	}
	query := url.Values{}
	if opts.Type != nil {
		query.Set("type", strconv.Itoa(*opts.Type))
	}
	if opts.SubjectType != nil {
		query.Set("subject_type", strconv.Itoa(*opts.SubjectType))
	}
	if opts.Limit > 0 {
		query.Set("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Offset > 0 {
		query.Set("offset", strconv.Itoa(opts.Offset))
	}

	var out UserCollectionsResult
	path := fmt.Sprintf("/v0/users/%s/collections", username)
	if err := c.doGet(ctx, path, query, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
