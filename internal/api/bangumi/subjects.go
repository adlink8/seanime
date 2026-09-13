package bangumi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"sync"
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

// platformLookupConcurrency platform 补查的并发上限（契约 03.7 §1 A1：≤4）。
// 用于防止瞬时打爆上游，仍保留（即使补查已不受共享 ticker 约束，见 fetchSubjectForPlatform）。
const platformLookupConcurrency = 4

// fetchSubjectForPlatform 读取条目详情（GET /v0/subjects/{id}）供 platform 补查使用。
//
// 为什么这条通路可以不受全局节流（Phase 3.7 A2-2）：
//   - 依据 .planning/codebase/BANGUMI-API-RECON.md:25，「OpenAPI 规范无任何 rate limit
//     头/说明」——client.go 的 defaultThrottleInterval=500ms 是客户端当初的**保守自选值**，
//     不是上游硬要求；且详情端点有 300s 服务端缓存，滥用风险低；
//   - 一次搜索要补查 20（关键词精准）~60+（小说占比低）个 id，若串行经共享 ticker
//     需 10~30s，用户会在验收时立刻发现；
//   - 安全网仍在：本通路与其他请求共用同一条 doRequest 管道，429 退避（Retry-After /
//     指数退避，最多 3 次）完整生效；并发由调用方的 platformLookupConcurrency=4 限制。
//
// 不变量：公开方法 GetSubject 行为不变、仍受共享 ticker 约束；只有此内部通路绕过节流。
func (c *Client) fetchSubjectForPlatform(ctx context.Context, subjectID int) (*Subject, error) {
	var out Subject
	if err := c.doGetWithThrottle(ctx, fmt.Sprintf("/v0/subjects/%d", subjectID), nil, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetSubjectPlatforms 批量读取条目的 platform（契约 03.7 §1 A1）。
//
// 用途：legacy 检索条目**不含 platform 字段**（契约 §0.1 实测地面真值），
// 书籍分区要按 platform 分流（D2）只能对命中 id 补查 v0 详情。
//
// 返回 map[int]string（id → platform）。语义与降级：
//   - 单条失败（网络错误 / >=400 / ctx 取消）→ 该 id 的 platform 视为**空串**，
//     **不中断整次检索**（契约 A1）；空串按 D4 归入漫画，绝不丢弃条目；
//   - 同 id 去重，避免重复请求；
//   - 缓存：Client 级内存缓存（c.platformCache，见 client.go 字段注释）。
//     生产 client 未传 WithFileCache（client.go 的 doGet 缓存分支落空），
//     故这层内存缓存是 platform 补查不重复打上游的关键。
//   - 补查请求走 fetchSubjectForPlatform（不受共享 ticker 约束，见其注释）。
func (c *Client) GetSubjectPlatforms(ctx context.Context, ids []int) map[int]string {
	out := make(map[int]string, len(ids))
	if len(ids) == 0 {
		return out
	}

	uniq := make([]int, 0, len(ids))
	seen := make(map[int]struct{}, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		uniq = append(uniq, id)
	}

	var (
		mu  sync.Mutex
		wg  sync.WaitGroup
		sem = make(chan struct{}, platformLookupConcurrency)
	)
	for _, id := range uniq {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// 缓存命中：直接回填，不发请求。
			c.platformMu.Lock()
			cached, hit := c.platformCache[id]
			c.platformMu.Unlock()
			if hit {
				mu.Lock()
				out[id] = cached
				mu.Unlock()
				return
			}

			sem <- struct{}{}
			defer func() { <-sem }()

			s, err := c.fetchSubjectForPlatform(ctx, id)
			mu.Lock()
			if err != nil || s == nil {
				// 解析失败（网络错误 / >=400 / ctx 取消）→ 本次结果为空串
				// （A1：不中断整次检索；D4：空 platform 归漫画）。
				//
				// **失败绝不写入缓存**：platform 缓存的生命周期是整个 Client，
				// 若把一次网络抖动产生的空串缓存下来，该 id 会被永久污染成
				// 「非小说」（D4 归漫画），且后续搜索再也拿不到正确值。
				// 不缓存则下次可重试；空串仍如实返回给本次调用方。
				out[id] = ""
			} else {
				out[id] = s.Platform
				// 仅成功解析才入缓存（platform 是不可变的上游属性，可安全长驻）。
				// 注意：成功但 platform 为空串也入缓存——那是上游确实没有该字段，
				// 与「请求失败」不同，重试也不会有别的结果。
				c.platformMu.Lock()
				c.platformCache[id] = s.Platform
				c.platformMu.Unlock()
			}
			mu.Unlock()
		}(id)
	}
	wg.Wait()

	return out
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
