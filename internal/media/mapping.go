package media

import (
	"sort"
	"sync"
	"time"

	"seanime/internal/api/animap"
)

// bangumi→anidb ID 映射服务。
//
// 核心难点（方向问题）：Bangumi subject ID **不在** animap 数据集中
// （anime-offline-database 覆盖 anidb/anilist/mal 等 13 种 ID，唯独没有 bangumi_id），
// 无法做纯 ID 直查。可行路径是：
//
//	Bangumi 名称（中/日/英）→ 匹配 animap 条目的 title/titles → 取该条目的 anidb_id
//
// 因此本服务加载 animap 条目集合后构建「归一化标题 → anidb ID」的**反向索引**，
// 解析时按名称查索引。匹配器（相似度函数）可注入，单测不打网络。

// NameSet 一次解析请求携带的名称集。
// Bangumi 提供日文原名与中文名；Romaji/英文名由编排层从其他来源补充
// （有则命中率更高，缺省不影响解析）。
type NameSet struct {
	Native  string   // 日文原名（Bangumi name）
	Romaji  string   // 罗马音（可选）
	English string   // 英文名（可选）
	Chinese string   // 中文名（Bangumi name_cn）
	Aliases []string // 其他别名（可选）
}

// All 返回全部非空名称。顺序不保证，仅供遍历。
func (n NameSet) All() []string {
	out := make([]string, 0, 4+len(n.Aliases))
	if n.Native != "" {
		out = append(out, n.Native)
	}
	if n.Romaji != "" {
		out = append(out, n.Romaji)
	}
	if n.English != "" {
		out = append(out, n.English)
	}
	if n.Chinese != "" {
		out = append(out, n.Chinese)
	}
	out = append(out, n.Aliases...)
	return out
}

// MappingService 映射服务接口（Phase 2 编排层依赖此接口）。
type MappingService interface {
	// ResolveBangumiToAniDB 尝试将 Bangumi subject ID 解析为 AniDB ID。
	// ok=false 表示无法保守确认——调用方应视该条目「映射未知」，
	// 而不是把返回值当作猜测结果（错误映射比缺映射危害大）。
	ResolveBangumiToAniDB(bangumiID int, names NameSet) (aniDBID int, ok bool)
	// QueueUnresolved 将无法解析的条目加入待映射队列（幂等）。
	QueueUnresolved(bangumiID int, names NameSet)
	// ResolveMalToAniDB 将 MAL ID 直接解析为 AniDB ID（animap 数据自带双 ID）。
	// ok=false 表示该 MAL ID 不在数据集中。
	ResolveMalToAniDB(malID int) (aniDBID int, ok bool)
}

//----------------------------------------------------------------------------------------------------------------------
// animapResolver

// titleNode 反向索引节点。
// conflict=true 表示该归一化标题对应多个不同 anidb ID（重名条目），
// 精确查询时直接拒绝——宁可入队也不猜。
type titleNode struct {
	anidbID  int
	title    string // 原始标题（调试用）
	conflict bool
}

// AnimapResolver 基于 animap 数据的 MappingService 实现。
type AnimapResolver struct {
	index      map[string]*titleNode // 归一化标题 -> 索引节点
	malIndex   map[int]int           // mal ID -> anidb ID（反向查找用）
	similarity SimilarityFunc        // 可注入的相似度函数
	queue      Queue                 // 待映射队列
	resolved   map[int]int           // bangumiID -> anidbID 解析缓存
	mu         sync.RWMutex
}

// 保守阈值说明：
//   - matchThreshold：模糊匹配下限。编辑距离 ratio ≥ 0.90 意味着最长标题
//     里最多允许 ~10% 字符差异，对标题级字符串足够紧，能挡住同名不同作
//     和前传/续作误配。精确匹配（归一化后相等）不受此阈值影响。
//   - ambiguityEpsilon：次优分数与最优分数差距小于该值时视为歧义入队。
//     典型歧义场景：剧场版与 TV 版标题仅差一个词、同名的不同年份重制版。
//   - minFuzzyLen：归一化后长度低于该值的标题只做精确匹配。短标题
//     （如 "日常"）编辑距离区分度太差，模糊匹配极易误配。
const (
	defaultMatchThreshold  = 0.90
	defaultAmbiguityEps    = 0.02
	defaultMinFuzzyRuneLen = 6
)

// ResolverOption 解析器可选配置。
type ResolverOption func(*AnimapResolver)

// WithSimilarity 注入自定义相似度函数（默认 Similarity）。
func WithSimilarity(f SimilarityFunc) ResolverOption {
	return func(r *AnimapResolver) { r.similarity = f }
}

// WithQueue 注入待映射队列（默认 nil，解析失败仅返回 false 不落盘；
// 生产环境应注入 FileQueue）。
func WithQueue(q Queue) ResolverOption {
	return func(r *AnimapResolver) { r.queue = q }
}

// NewAnimapResolver 用 animap 条目集合构建解析器。
// entries 由调用方提供（animap 包只提供按 ID 单查的 API，无全量拉取；
// 全量数据的获取与缓存是编排层/Phase 4 清偿工具的职责）。
func NewAnimapResolver(entries []animap.Anime, opts ...ResolverOption) *AnimapResolver {
	r := &AnimapResolver{
		index:      BuildTitleIndex(entries),
		malIndex:   BuildMalIndex(entries),
		similarity: Similarity,
		resolved:   make(map[int]int),
	}
	for _, o := range opts {
		o(r)
	}
	return r
}


// BuildTitleIndex 构建「归一化标题 → anidb ID」反向索引。
//
// 为什么是反向索引：animap 的原生查询方向是 anidb/anilist/mal ID → entry，
// 我们需要反过来按名称找 ID。索引覆盖 entry.Title 与 entry.Titles 的全部值
// （animap titles 的 key 是语言代码，值才是各语言标题，统一取值不做语言区分，
// 因为查询侧 NameSet 也混合多语言）。
//
// 同一归一化标题命中多个不同 anidb ID 时标记 conflict（保守拒绝）。
func BuildTitleIndex(entries []animap.Anime) map[string]*titleNode {
	index := make(map[string]*titleNode)
	add := func(title string, anidbID int) {
		if title == "" || anidbID == 0 {
			return
		}
		key := Normalize(title)
		if key == "" {
			return
		}
		if node, ok := index[key]; ok {
			// 重名：不同 ID 则标记冲突
			if node.anidbID != anidbID {
				node.conflict = true
			}
			return
		}
		index[key] = &titleNode{anidbID: anidbID, title: title}
	}
	for _, e := range entries {
		id := 0
		if e.Mappings != nil {
			id = e.Mappings.AnidbID
		}
		if id == 0 {
			continue
		}
		add(e.Title, id)
		for _, t := range e.Titles {
			add(t, id)
		}
	}
	return index
}

// BuildMalIndex 构建「MAL ID → AniDB ID」直查索引（additive，Wave B 反向查找用）。
// animap 条目的 Mappings 同时携带 anidb_id 与 mal_id，属同一作品的双 ID，
// 可直接建索引；同一 MAL ID 命中多个不同 AniDB ID 时保守丢弃（宁可缺映射）。
func BuildMalIndex(entries []animap.Anime) map[int]int {
	index := make(map[int]int)
	for _, e := range entries {
		if e.Mappings == nil || e.Mappings.MalID == 0 || e.Mappings.AnidbID == 0 {
			continue
		}
		if existing, ok := index[e.Mappings.MalID]; ok && existing != e.Mappings.AnidbID {
			continue // 冲突：保留首个，其余丢弃
		}
		index[e.Mappings.MalID] = e.Mappings.AnidbID
	}
	return index
}

// ResolveMalToAniDB 实现 MappingService 的 MAL → AniDB 直查。
func (r *AnimapResolver) ResolveMalToAniDB(malID int) (int, bool) {
	if malID <= 0 {
		return 0, false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.malIndex[malID]
	return id, ok
}

//----------------------------------------------------------------------------------------------------------------------
// 反向 ID 索引（anidb/mal → bangumi）

// ReverseIDIndex 反向 ID 索引（additive，Wave B 平台层 GetAnimeByMalID 用）。
//
// 背景：animap 数据集不含 bangumi ID，无法离线直查 anidb/mal → bangumi。
// 本索引采用「机会式」填充策略：平台层每次通过正向解析（名称 → anidb）
// 或其他渠道确认了 bangumi ↔ anidb/mal 的对应关系后调用 Put 记录，
// 后续 GetAnimeByMalID 先查本索引，命中即零成本直查。
//
// 未命中时调用方应自行走降级路径（如入队等待 Phase 4 清偿工具回填），
// 而不是猜测——错误映射比缺映射危害大。
type ReverseIDIndex struct {
	mu             sync.RWMutex
	anidbToBangumi map[int]int
	malToBangumi   map[int]int
}

// NewReverseIDIndex 返回空的反向索引。
func NewReverseIDIndex() *ReverseIDIndex {
	return &ReverseIDIndex{
		anidbToBangumi: make(map[int]int),
		malToBangumi:   make(map[int]int),
	}
}

// Put 记录一组已确认的 ID 对应关系（任意 ID 为 0 时忽略该方向）。
// 后写覆盖先写（以最近确认为准）。
func (r *ReverseIDIndex) Put(bangumiID, aniDBID, malID int) {
	if bangumiID <= 0 {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if aniDBID > 0 {
		r.anidbToBangumi[aniDBID] = bangumiID
	}
	if malID > 0 {
		r.malToBangumi[malID] = bangumiID
	}
}

// BangumiByAniDB 按 AniDB ID 查 bangumi ID。
func (r *ReverseIDIndex) BangumiByAniDB(aniDBID int) (int, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.anidbToBangumi[aniDBID]
	return id, ok
}

// BangumiByMal 按 MAL ID 查 bangumi ID。
func (r *ReverseIDIndex) BangumiByMal(malID int) (int, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.malToBangumi[malID]
	return id, ok
}

// Len 返回已记录的 bangumi ID 数（诊断用）。
func (r *ReverseIDIndex) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.anidbToBangumi)
}

// ResolveBangumiToAniDB 实现 MappingService。
// 解析策略（保守优先）：
//  1. 已缓存直接返回（同一条目多次解析零开销）；
//  2. 精确匹配：任一名称归一化后命中索引 → 返回（冲突节点拒绝）；
//  3. 模糊匹配：对全索引计算相似度取最优，需同时满足
//     最优分 ≥ threshold 且 领先次优 ≥ epsilon，否则视为歧义；
//  4. 失败：入队（若配置）并返回 false。
func (r *AnimapResolver) ResolveBangumiToAniDB(bangumiID int, names NameSet) (int, bool) {
	r.mu.RLock()
	if id, ok := r.resolved[bangumiID]; ok {
		r.mu.RUnlock()
		return id, true
	}
	r.mu.RUnlock()

	id, ok := r.resolve(bangumiID, names)
	if ok {
		r.mu.Lock()
		r.resolved[bangumiID] = id
		r.mu.Unlock()
	} else {
		r.QueueUnresolved(bangumiID, names)
	}
	return id, ok
}

// resolve 实际解析逻辑（不加缓存锁）。
func (r *AnimapResolver) resolve(bangumiID int, names NameSet) (int, bool) {
	// 精确匹配
	for _, name := range names.All() {
		key := Normalize(name)
		if key == "" {
			continue
		}
		if node, ok := r.index[key]; ok {
			if node.conflict {
				return 0, false // 重名条目，保守拒绝
			}
			return node.anidbID, true
		}
	}

	// 模糊匹配
	type score struct {
		id    int
		score float64
	}
	bestByID := make(map[int]float64)
	for _, name := range names.All() {
		q := Normalize(name)
		if len([]rune(q)) < defaultMinFuzzyRuneLen {
			continue // 短名称不参与模糊匹配
		}
		for key, node := range r.index {
			if node.conflict || len([]rune(key)) < defaultMinFuzzyRuneLen {
				continue
			}
			s := r.similarity(q, key)
			if s > bestByID[node.anidbID] {
				bestByID[node.anidbID] = s
			}
		}
	}
	if len(bestByID) > 0 {
		scores := make([]score, 0, len(bestByID))
		for id, s := range bestByID {
			scores = append(scores, score{id, s})
		}
		sort.Slice(scores, func(i, j int) bool { return scores[i].score > scores[j].score })
		best, second := scores[0].score, 0.0
		if len(scores) > 1 {
			second = scores[1].score
		}
		if best >= defaultMatchThreshold && (best-second) >= defaultAmbiguityEps {
			return scores[0].id, true
		}
		// 分数不够高或存在近似并列：视为歧义，拒绝
		return 0, false
	}

	return 0, false
}

// QueueUnresolved 实现 MappingService。未配置队列时静默跳过。
func (r *AnimapResolver) QueueUnresolved(bangumiID int, names NameSet) {
	if r.queue == nil {
		return
	}
	_ = r.queue.Add(Unresolved{
		BangumiID: bangumiID,
		Names:     names,
		QueuedAt:  time.Now(),
	})
}
