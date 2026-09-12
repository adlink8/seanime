package bangumi

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// SearchFilter 搜索过滤条件（POST /v0/search/subjects 请求体的 filter 字段）
type SearchFilter struct {
	Type []int    `json:"type,omitempty"` // 条目类型：1书籍 2动画 3音乐 4游戏 6三次元
	Tag  []string `json:"tag,omitempty"`  // 标签过滤（如 "ASMR"）
	// Nsfw 三态：nil 不过滤；true 只返回 R18（音声域依赖此项）；false 排除 R18
	Nsfw *bool `json:"nsfw,omitempty"`
}

// SearchSubjectsOpts 搜索参数
type SearchSubjectsOpts struct {
	Keyword string
	// Sort 排序：match / heat / rank / score；留空走服务端默认
	Sort   string
	Filter SearchFilter
	Limit  int
	Offset int
}

// searchSubjectsBody POST 请求体（limit/offset 走 query，不在 body 里）。
// Filter 用指针：结构体类型 omitempty 不生效，需要整体判空后决定是否携带。
type searchSubjectsBody struct {
	Keyword string        `json:"keyword"`
	Sort    string        `json:"sort,omitempty"`
	Filter  *SearchFilter `json:"filter,omitempty"`
}

// SearchSubjects 搜索条目（POST /v0/search/subjects，实验性端点）。
// keyword + filter 走请求体，limit/offset 走 query。写方法路径，不缓存。
func (c *Client) SearchSubjects(ctx context.Context, opts SearchSubjectsOpts) (*SearchResult, error) {
	query := url.Values{}
	if opts.Limit > 0 {
		query.Set("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Offset > 0 {
		query.Set("offset", strconv.Itoa(opts.Offset))
	}

	var filter *SearchFilter
	if len(opts.Filter.Type) > 0 || len(opts.Filter.Tag) > 0 || opts.Filter.Nsfw != nil {
		filter = &opts.Filter
	}

	var out SearchResult
	err := c.doWrite(ctx, http.MethodPost, "/v0/search/subjects", query, searchSubjectsBody{
		Keyword: opts.Keyword,
		Sort:    opts.Sort,
		Filter:  filter,
	}, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
