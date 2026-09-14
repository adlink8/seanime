// Package media 提供 Seanime 自有的领域模型与 bangumi→anidb ID 映射服务。
//
// 设计背景：Seanime 现有库条目直接复用 AniList 的 anilist.BaseAnime/BaseManga，
// 数据结构与数据源强耦合。换锚 Bangumi 前必须先建立一层自有模型，
// 使上游（库管理、元数据 provider 等）与具体数据源解耦。
//
// 本包与 internal/api/bangumi 并行开发（Phase 1 契约），故不 import 该包：
// 这里定义的 Subject 结构体按 Bangumi v0 API 的 JSON 形状镜像其关键字段，
// Phase 2 编排层可直接将 client 响应反序列化进 Subject（或做薄适配），
// 两个包在接口层面汇合，互不依赖。
package media

import (
	"encoding/json"
	"time"
)

// MediaKind 条目类型。
// 命名说明：Phase 1 原名 MediaType，Wave A 镜像 anilist 领域类型时
// 与 AniList 枚举 MediaType("ANIME"/"MANGA") 撞名，故改名为 MediaKind。
// JSON 值（"anime"/"manga"/…）不变。
// ASMR 类目在 Bangumi 无独立 type（归书籍分区），asmr 锚点使用 asmr.one 的 workID，
// 由 Media.ASMRWorkID 承载，Phase 2/3 编排层负责填充。
type MediaKind string

const (
	MediaKindAnime   MediaKind = "anime"
	MediaKindManga   MediaKind = "manga"
	MediaKindNovel   MediaKind = "novel"
	MediaKindASMR    MediaKind = "asmr"
	MediaKindUnknown MediaKind = ""
)

// Media 是 Seanime 的统一领域模型。
// ID 字段为 Bangumi subject ID（项目换锚后的主锚点）。
type Media struct {
	ID   int       `json:"id"`   // Bangumi subject ID
	Type MediaKind `json:"type"` // 条目类型

	Name   string `json:"name"`    // 原名（日文）
	NameCN string `json:"nameCN"`  // 中文名
	Summary string `json:"summary"` // 简介（Bangumi 返回中文）
	ImageURL string `json:"imageUrl"`

	Tags []string `json:"tags,omitempty"`

	AirDate *time.Time `json:"airDate,omitempty"` // 放送/连载开始日期
	EndDate *time.Time `json:"endDate,omitempty"` // 结束日期（Bangumi v0 无直接字段，Phase 2 可由 aired 补充）

	Eps int `json:"eps"` // 集数（书籍类为话数，音声类为 0）

	NSFW   bool    `json:"nsfw"`
	Rating float64 `json:"rating"` // Bangumi 评分（0-10）

	// ASMRWorkID 保留字段：asmr.one 的 work ID。
	// asmr 锚点场景下 Media.ID 仍为 bangumi subject ID，workID 通过映射层关联。
	ASMRWorkID int `json:"asmrWorkId,omitempty"`
}

//----------------------------------------------------------------------------------------------------------------------
// Bangumi v0 Subject 镜像结构（鸭子类型适配）

// Subject 按 Bangumi v0 API `GET /v0/subjects/{id}` 响应的关键字段形状定义。
// 这里的字段名/JSON tag 与 bangumi API 对齐，但类型是本包自有定义——
// Phase 2 编排层可直接 json.Unmarshal 进该结构，或编写薄适配器，
// 从而保持本包对 internal/api/bangumi 零依赖。
type Subject struct {
	ID       int            `json:"id"`               // subject ID
	Type     int            `json:"type"`             // 1书籍/2动画/3音乐/4游戏/6三次元
	Name     string         `json:"name"`             // 原名（日文）
	NameCN   string         `json:"name_cn"`          // 中文名
	Summary  string         `json:"summary"`          // 简介（中文）
	Date     string         `json:"date,omitempty"`   // 放送/出版日期，如 2009-04-03
	Platform string         `json:"platform,omitempty"` // 放送/出版平台（TV / 剧场版 / 漫画 / 小说…，Wave A 新增）
	Eps      int            `json:"eps"`              // 话数
	NSFW     bool           `json:"nsfw"`             // 是否 R18（查 R18 内容必须带 token）
	Images   *SubjectImages `json:"images,omitempty"` // 封面图（5 尺寸）
	Tags     []SubjectTag   `json:"tags,omitempty"`   // 标签
	Rating   *SubjectRating `json:"rating,omitempty"` // 评分
	// Infobox（3.9d additive）：Bangumi v0 信息箱，Value 形态不固定
	//（单值字符串或 [{"v":"..."}] 数组），延迟解析。
	Infobox []BangumiInfoboxEntry `json:"infobox,omitempty"`
}

type SubjectImages struct {
	Large  string `json:"large,omitempty"`
	Common string `json:"common,omitempty"`
	Medium string `json:"medium,omitempty"`
	Small  string `json:"small,omitempty"`
	Grid   string `json:"grid,omitempty"`
}

type SubjectTag struct {
	Name    string `json:"name"`
	Count   int    `json:"count"`
	Spoiler bool   `json:"spoiler,omitempty"` // 3.9d additive：剧透标签
}

// BangumiTagInfo（3.9d additive）：带热度与剧透标记的标签，
// 由 Subject.Tags 映射而来，随 Anime/Manga/AnimeDetails 输出到前端。
type BangumiTagInfo struct {
	Name    string `json:"name"`
	Count   int    `json:"count"`
	Spoiler bool   `json:"spoiler,omitempty"`
}

// BangumiInfoboxEntry（3.9d additive）：信息箱条目。
// Value 兼容两种形态：单值字符串（如 "TV"）与对象数组（如 [{"v":"京都动画"}]），
// 用 json.RawMessage 延迟解析，序列化时原样透传给前端。
type BangumiInfoboxEntry struct {
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
}

type SubjectRating struct {
	Rank  int     `json:"rank"`
	Total int     `json:"total"`
	Count map[string]int `json:"count"`
	Score float64 `json:"score"`
}

// FromSubject 将 Bangumi Subject 镜像结构转换为领域模型 Media。
//
// 转换规则：
//   - Type 映射见 mediaTypeFromBangumi：书籍(1) 默认归为 novel，
//     漫画/音声的细分依赖 tags（如 "ASMR" 标签），由 Phase 2 编排层二次修正；
//   - Date 解析失败不报错（置 nil），日期格式兼容 "2006-01-02" 与 RFC3339；
//   - NameCN 为空时回退到 Name。
func FromSubject(s *Subject) *Media {
	if s == nil {
		return nil
	}

	m := &Media{
		ID:      s.ID,
		Type:    mediaTypeFromBangumi(s.Type),
		Name:    s.Name,
		NameCN:  s.NameCN,
		Summary: s.Summary,
		Eps:     s.Eps,
		NSFW:    s.NSFW,
		Tags:    make([]string, 0, len(s.Tags)),
	}

	if m.NameCN == "" {
		m.NameCN = s.Name
	}

	for _, t := range s.Tags {
		m.Tags = append(m.Tags, t.Name)
	}

	if s.Images != nil {
		// 优先 common（页面标准尺寸），其次 large
		m.ImageURL = s.Images.Common
		if m.ImageURL == "" {
			m.ImageURL = s.Images.Large
		}
	}

	if s.Rating != nil {
		m.Rating = s.Rating.Score
	}

	m.AirDate = parseBangumiDate(s.Date)

	// ASMR 标签嗅探：命中则标记为音声类目（Phase 2 实测确认 tag 方案后保留/调整）
	for _, t := range m.Tags {
		if t == "ASMR" {
			m.Type = MediaKindASMR
			break
		}
	}

	return m
}

// mediaTypeFromBangumi 将 Bangumi SubjectType 枚举映射为 Media 类型。
// 1书籍 / 2动画 / 3音乐 / 4游戏 / 6三次元（无 5）。
// 书籍分区下混有漫画/轻小说/音声，此处保守地默认归为 novel，
// 细分（manga/asmr）交由编排层根据 subType/tags 修正。
func mediaTypeFromBangumi(t int) MediaKind {
	switch t {
	case 1:
		return MediaKindNovel
	case 2:
		return MediaKindAnime
	default:
		return MediaKindUnknown
	}
}

// parseBangumiDate 解析 Bangumi 日期字符串，失败返回 nil（容忍脏数据）。
func parseBangumiDate(s string) *time.Time {
	if s == "" {
		return nil
	}
	for _, layout := range []string{"2006-01-02", time.RFC3339, "2006-01"} {
		if t, err := time.Parse(layout, s); err == nil {
			return &t
		}
	}
	return nil
}

// NameSetNames 供映射层使用的便捷方法：将 Media 的名称集转为 NameSet。
// Bangumi 无罗马音字段，Romaji 留空——映射层按名称匹配时以中文名与原名为主。
func (m *Media) NameSet() NameSet {
	return NameSet{
		Native:  m.Name,
		Chinese: m.NameCN,
	}
}
