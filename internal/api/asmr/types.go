package asmr

import "strconv"

// 类型分两层：
//   - rawXxx：asmr.one 原始响应字段（实测地面真值 2026-09-13，字段名保持 API 原样）
//   - Asmr_Xxx：契约暴露给前端/handlers 的形状（.planning/phases/02-library-reanchor/02.5-CONTRACT.md 第 4 节），
//     字段名必须逐字一致；原始字段 → 契约字段的映射全部在 convert 函数里完成

// rawVAS 原始声优条目
type rawVAS struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// rawTag 原始标签条目（i18n 结构已省略，取中文规范化 name）
type rawTag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// rawCircle 原始社团条目
type rawCircle struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	SourceID string `json:"source_id"`
}

// rawWork asmr.one work 原始结构（搜索结果条目与 /api/workInfo/{id} 响应同构）
type rawWork struct {
	ID            int        `json:"id"`
	Title         string     `json:"title"`
	CircleID      int        `json:"circle_id"`
	CircleName    string     `json:"name"` // 顶层 name 是社团名（circle 冗余字段），与 circle.name 一致
	Nsfw          bool       `json:"nsfw"`
	Release       string     `json:"release"` // YYYY-MM-DD
	DlCount       int        `json:"dl_count"`
	Price         int        `json:"price"`
	RateAverage2D float64    `json:"rate_average_2dp"` // 0-5
	// Phase 3.9b P0 徽章字段（additive，上游原始字段名保持）
	Duration  int `json:"duration"`   // 时长（分钟）
	RateCount int `json:"rate_count"` // 评分人数
	HasSubtitle   bool       `json:"has_subtitle"`
	SourceID      string     `json:"source_id"` // "RJ01657200"
	Vas           []rawVAS   `json:"vas"`
	Tags          []rawTag   `json:"tags"`
	Circle        *rawCircle `json:"circle"`
	MainCoverURL  string     `json:"mainCoverUrl"`
}

// rawPagination 搜索响应分页块
type rawPagination struct {
	CurrentPage int `json:"currentPage"`
	PageSize    int `json:"pageSize"`
	TotalCount  int `json:"totalCount"`
}

// rawSearchResult 搜索响应：{"works":[...], "pagination":{...}}
type rawSearchResult struct {
	Works      []rawWork     `json:"works"`
	Pagination rawPagination `json:"pagination"`
}

// rawTrackNode /api/tracks/{id}?v=2 的节点（根为 JSON 数组，folder 用 children 嵌套）
type rawTrackNode struct {
	Type             string         `json:"type"` // folder / audio / text / image
	Title            string         `json:"title"`
	Children         []rawTrackNode `json:"children,omitempty"` // 仅 folder
	MediaStreamURL   string         `json:"mediaStreamUrl,omitempty"`
	MediaDownloadURL string         `json:"mediaDownloadUrl,omitempty"`
}

// ———— 契约暴露类型（字段名逐字对齐契约第 4 节，前端按此 mock，不得改动） ————

// Asmr_Work 音声作品列表项
type Asmr_Work struct {
	ID          string   `json:"id"`   // asmr.one work id（字符串化）
	RjID        string   `json:"rjId"` // "RJ01234567"
	Title       string   `json:"title"`
	Circle      string   `json:"circle"` // 社团名
	Cvs         []string `json:"cvs"`    // 声优
	Tags        []string `json:"tags"`   // 标签名
	Nsfw        bool     `json:"nsfw"`
	CoverURL    string   `json:"coverUrl"`
	ReleaseDate string   `json:"releaseDate"` // YYYY-MM-DD
	Rating      float64  `json:"rating"`      // 0-5，可为 0
	DlCount     int      `json:"dlCount"`
	Price       int      `json:"price"`
	HasSubtitle bool     `json:"hasSubtitle"`
	// Phase 3.9b P0 徽章字段（additive）
	Duration  int `json:"duration,omitempty"`  // 时长（分钟）
	RateCount int `json:"rateCount,omitempty"` // 评分人数
	// Phase 3.9e 降级方案（additive）：社团 DLsite source_id（如 "RG51931"），
	// 详情页社团名外链 asmr.one 社团页 /circle/{source_id}；聚合端点 D9 证伪后无站内聚合能力。
	CircleSourceID string `json:"circleSourceId,omitempty"`
}

// Asmr_PageInfo 分页信息
type Asmr_PageInfo struct {
	CurrentPage int  `json:"currentPage"`
	PerPage     int  `json:"perPage"`
	HasNextPage bool `json:"hasNextPage"`
	Total       int  `json:"total"`
}

// Asmr_SearchResult 搜索/popular 响应
type Asmr_SearchResult struct {
	Works    []Asmr_Work   `json:"works"`
	PageInfo Asmr_PageInfo `json:"pageInfo"`
}

// Asmr_Track 音轨节点（folder 递归嵌套）
type Asmr_Track struct {
	Title            string       `json:"title"`
	Type             string       `json:"type"` // "audio"|"folder"|"text"|"image"
	MediaStreamURL   string       `json:"mediaStreamUrl,omitempty"`
	MediaDownloadURL string       `json:"mediaDownloadUrl,omitempty"`
	// LocalPath 音声本地库扩展：本地存在的音频文件绝对路径（契约 §2 合并规则）。
	// additive 新增可选字段，既有字段名一律不动。
	LocalPath string `json:"localPath,omitempty"`
	// Path 音声本地库扩展：相对 RJ 目录的 '/' 分隔路径（仅本地音轨），
	// 前端用于 POST /asmr/track/progress 的 trackPath。
	Path      string       `json:"path,omitempty"`
	Tracks    []Asmr_Track `json:"tracks,omitempty"` // folder 嵌套
	// Completed 音声本地库扩展：单条音轨的逐轨完听状态（契约 3.2a / D1）。
	// additive 新增可选字段，既有字段名一律不动；结束时省略（omitempty）表示未完成或本地库路径不回填。
	// 注意：该字段仅由本地库路径 Scanner.GetWork 回填；线上搜索路径 Asmr_WorkDetail 永不回填
	// （在线音轨无 Path，天然不会被匹配）——这是预期行为，非 bug（契约 §8）。
	Completed bool `json:"completed,omitempty"`
}

// Asmr_WorkDetail 作品详情（Asmr_Work 全字段 + tracks）。
// 该结构被线上搜索详情复用：其 tracks 来自 convertTracks，节点无 Path，
// 故 Completed 永为 omitted（契约 §8，预期行为，非 bug）。
type Asmr_WorkDetail struct {
	Asmr_Work
	Tracks []Asmr_Track `json:"tracks"`
}

// convertWork 原始 work → 契约 Asmr_Work
func convertWork(w rawWork) Asmr_Work {
	out := Asmr_Work{
		ID:          strconv.Itoa(w.ID),
		RjID:        w.SourceID,
		Title:       w.Title,
		Circle:      w.CircleName,
		Nsfw:        w.Nsfw,
		CoverURL:    w.MainCoverURL,
		ReleaseDate: w.Release,
		Rating:      w.RateAverage2D,
		DlCount:     w.DlCount,
		Price:       w.Price,
		HasSubtitle: w.HasSubtitle,
		Duration:    w.Duration,
		RateCount:   w.RateCount,
	}
	if w.Circle != nil && w.Circle.Name != "" {
		// 顶层 name 缺失时回退 circle.name
		out.Circle = w.Circle.Name
		out.CircleSourceID = w.Circle.SourceID
	}
	out.Cvs = make([]string, 0, len(w.Vas))
	for _, v := range w.Vas {
		out.Cvs = append(out.Cvs, v.Name)
	}
	out.Tags = make([]string, 0, len(w.Tags))
	for _, t := range w.Tags {
		out.Tags = append(out.Tags, t.Name)
	}
	return out
}

// convertSearchResult 原始搜索响应 → 契约 Asmr_SearchResult
func convertSearchResult(raw rawSearchResult) Asmr_SearchResult {
	out := Asmr_SearchResult{
		Works: make([]Asmr_Work, 0, len(raw.Works)),
		PageInfo: Asmr_PageInfo{
			CurrentPage: raw.Pagination.CurrentPage,
			PerPage:     raw.Pagination.PageSize,
			Total:       raw.Pagination.TotalCount,
		},
	}
	for _, w := range raw.Works {
		out.Works = append(out.Works, convertWork(w))
	}
	// hasNextPage：当前页已取满且未达总数
	out.PageInfo.HasNextPage = raw.Pagination.CurrentPage*raw.Pagination.PageSize < raw.Pagination.TotalCount
	return out
}

// convertTracks 原始音轨节点树 → 契约 Asmr_Track 树。
// 原始 folder 用 children 嵌套、叶子字段扁平，契约统一为 type/tracks 嵌套形状。
func convertTracks(nodes []rawTrackNode) []Asmr_Track {
	out := make([]Asmr_Track, 0, len(nodes))
	for _, n := range nodes {
		t := Asmr_Track{
			Title:            n.Title,
			Type:             n.Type,
			MediaStreamURL:   n.MediaStreamURL,
			MediaDownloadURL: n.MediaDownloadURL,
		}
		if n.Type == "folder" {
			t.Tracks = convertTracks(n.Children)
		}
		out = append(out, t)
	}
	return out
}
