package bangumi

// 手写最小响应类型集，字段对齐 Bangumi API v0 OpenAPI 规范。
// 自用客户端，不做全套代码生成（YAGNI）；服务端新增字段会被 json 忽略。

import "encoding/json"

// SubjectType 枚举（无 5）
const (
	SubjectBook  = 1 // 书籍（音声作品归此分区下的「音声」子分类）
	SubjectAnime = 2 // 动画
	SubjectMusic = 3 // 音乐
	SubjectGame  = 4 // 游戏
	SubjectReal  = 6 // 三次元
)

// CollectionType 收藏类型枚举
const (
	CollectionWish    = 1 // 想看
	CollectionDone    = 2 // 看过
	CollectionDoing   = 3 // 在看
	CollectionOnHold  = 4 // 搁置
	CollectionDropped = 5 // 抛弃
)

// 章节进度类型（PATCH /v0/users/-/collections/{subject_id}/episodes 的 type 字段）
const (
	ProgressWatched = "watched" // 看过
	ProgressQueue   = "queue"   // 想看
	ProgressDropped = "dropped" // 抛弃
)

// InfoboxItem 条目信息箱条目（GET /v0/subjects/{id} 的 infobox 字段，Wave B additive）。
// Value 形态不固定：单值时为 JSON 字符串（如 "TV"），
// 多值时为 [{"v":"京都动画"}] 对象数组，统一用 json.RawMessage 延迟解析。
type InfoboxItem struct {
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
}

// Subject 条目详情（GET /v0/subjects/{id} / 搜索结果元素）
type Subject struct {
	ID       int            `json:"id"`
	Type     int            `json:"type"`     // 1书籍 2动画 3音乐 4游戏 6三次元
	Name     string         `json:"name"`     // 日文原名
	NameCN   string         `json:"name_cn"`  // 中文名（换锚的核心收益）
	Date     string         `json:"date"`     // 放送/发售日期
	Summary  string         `json:"summary"`  // 中文简介
	Platform string         `json:"platform"` // TV / Web / OVA / 游戏...
	Images   SubjectImages  `json:"images"`
	Nsfw     bool           `json:"nsfw"`
	Eps      int            `json:"eps"` // 章节数
	Rating   *SubjectRating `json:"rating,omitempty"`
	Tags     []SubjectTag   `json:"tags,omitempty"`
	Infobox  []InfoboxItem  `json:"infobox,omitempty"` // 信息箱（制作公司等，Wave B additive）
}

type SubjectImages struct {
	Large  string `json:"large"`
	Common string `json:"common"`
	Medium string `json:"medium"`
	Small  string `json:"small"`
	Grid   string `json:"grid"`
}

type SubjectRating struct {
	Rank  int            `json:"rank"`
	Total int            `json:"total"`
	Score float64        `json:"score"`
	Count map[string]int `json:"count"`
}

type SubjectTag struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// Episode 章节（EpType：0本篇 1SP 2OP 3ED 4PV 5MAD 6其他）
type Episode struct {
	ID      int    `json:"id"`
	Type    int    `json:"type"`
	Name    string `json:"name"`
	NameCN  string `json:"name_cn"` // 中文章节名
	AirDate string `json:"airdate"`
	Desc    string `json:"desc"`
}

type EpisodesResult struct {
	Total int       `json:"total"`
	Data  []Episode `json:"data"`
}

// SearchResult POST /v0/search/subjects 响应
type SearchResult struct {
	Total  int       `json:"total"`
	Limit  int       `json:"limit"`
	Offset int       `json:"offset"`
	Data   []Subject `json:"data"`
}

// UserCollection 用户收藏条目
type UserCollection struct {
	SubjectID int      `json:"subject_id"`
	Subject   *Subject `json:"subject,omitempty"` // 列表响应内嵌条目详情
	Type      int      `json:"type"`              // 1想看 2看过 3在看 4搁置 5抛弃
	Rate      int      `json:"rate"`              // 0-10
	Comment   string   `json:"comment"`
	Private   bool     `json:"private"`
	Tags      []string `json:"tags"`
	EpStatus  int      `json:"ep_status"`  // 已观看章节数
	VolStatus int      `json:"vol_status"` // 已观看卷数（仅书籍类）
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type UserCollectionsResult struct {
	Total  int              `json:"total"`
	Limit  int              `json:"limit"`
	Offset int              `json:"offset"`
	Data   []UserCollection `json:"data"`
}

// CollectionUpsertBody POST/PATCH /v0/users/-/collections/{subject_id} 请求体。
// 指针字段区分「未提供」与零值；POST 不存在则创建，PATCH 局部更新。
type CollectionUpsertBody struct {
	Type      *int     `json:"type,omitempty"` // 1-5
	Rate      *int     `json:"rate,omitempty"` // 0-10
	Comment   *string  `json:"comment,omitempty"`
	Private   *bool    `json:"private,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	EpStatus  *int     `json:"ep_status,omitempty"`  // 仅书籍类
	VolStatus *int     `json:"vol_status,omitempty"` // 仅书籍类
}

// EpisodeCollection 章节观看进度（GET /v0/users/-/collections/{subject_id}/episodes）
type EpisodeCollection struct {
	ID     int    `json:"id"`
	Type   int    `json:"type"` // EpType
	Name   string `json:"name"`
	NameCN string `json:"name_cn"`
	Status string `json:"status"` // watched / queue / dropped
}

type EpisodeCollectionsResult struct {
	SubjectID int                 `json:"subject_id"`
	Total     int                 `json:"total"`
	Data      []EpisodeCollection `json:"data"`
}

// EpisodeProgressBody PATCH 章节进度请求体；服务端自动重算完成度
type EpisodeProgressBody struct {
	EpisodeIDs []int  `json:"episode_id"`
	Type       string `json:"type"` // watched / queue / dropped
}

// User 当前 token 对应用户（GET /v0/me）
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Sign     string `json:"sign"`
}
