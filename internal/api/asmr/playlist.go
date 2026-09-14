package asmr

import (
	"context"
	"net/url"
	"strconv"
	"strings"
)

// 播放列表域（契约 .planning/phases/02-library-reanchor/03.8-CONTRACT.md §0/§1）。
//
// 地面真值（2026-09-14 编排层实测 + Yuro 客户端源码佐证，勿重复实测）：
//   - GET /api/playlist/get-playlists?page=N → {"playlists":[{id(uuid), name, works_count,
//     latestWorkID, mainCoverUrl, updated_at, ...}]}
//   - GET /api/playlist/get-playlist-works?id=<uuid>&page=N&pageSize=M → WorksResponse 信封
//     （works 条目与 rawWork 同构，复用 convertSearchResult）
//   - POST /api/playlist/add-works-to-playlist / remove-works-from-playlist，
//     body 恰为 {id: uuid, works: [int]}
//   - 猜测路径 /api/playlists 等均 404（实测）——勿再探测。
//   - 系统播放列表（__SYS_PLAYLIST_MARKED / __SYS_PLAYLIST_LIKED）与 review 域重叠，
//     D2 裁决：client 层过滤不透出。

// rawPlaylist get-playlists 单条原始结构（字段名保持 API 原样）。
type rawPlaylist struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Privacy       int    `json:"privacy"`
	PlaybackCount int    `json:"playback_count"`
	WorksCount    int    `json:"works_count"`
	LatestWorkID  int    `json:"latestWorkID"`
	MainCoverURL  string `json:"mainCoverUrl"`
	UpdatedAt     string `json:"updated_at"`
}

// rawPlaylists get-playlists 响应信封。
type rawPlaylists struct {
	Playlists []rawPlaylist `json:"playlists"`
}

// Asmr_Playlist 契约播放列表（系统列表已过滤）。
type Asmr_Playlist struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	WorksCount   int    `json:"worksCount"`
	LatestWorkID string `json:"latestWorkId"` // 字符串化，与 Asmr_Work.ID 一致
	CoverURL     string `json:"coverUrl"`
	UpdatedAt    string `json:"updatedAt"`
}

// convertPlaylist 原始 → 契约；系统列表（__SYS_ 前缀）返回 false 表示过滤。
func convertPlaylist(p rawPlaylist) (Asmr_Playlist, bool) {
	if strings.HasPrefix(p.Name, "__SYS_") {
		return Asmr_Playlist{}, false
	}
	return Asmr_Playlist{
		ID:           p.ID,
		Name:         p.Name,
		WorksCount:   p.WorksCount,
		LatestWorkID: strconv.Itoa(p.LatestWorkID),
		CoverURL:     p.MainCoverURL,
		UpdatedAt:    p.UpdatedAt,
	}, true
}

// GetPlaylists 拉取当前账号播放列表（系统列表已过滤）。
func (c *Client) GetPlaylists(ctx context.Context, page int) ([]Asmr_Playlist, error) {
	var raw rawPlaylists
	if err := c.doGet(ctx, "/playlist/get-playlists", url.Values{"page": {strconv.Itoa(page)}}, &raw); err != nil {
		return nil, err
	}
	out := make([]Asmr_Playlist, 0, len(raw.Playlists))
	for _, p := range raw.Playlists {
		if conv, ok := convertPlaylist(p); ok {
			out = append(out, conv)
		}
	}
	return out, nil
}

// GetPlaylistWorks 拉取指定播放列表内的作品（WorksResponse 信封，复用搜索转换）。
func (c *Client) GetPlaylistWorks(ctx context.Context, playlistID string, page, pageSize int) (Asmr_SearchResult, error) {
	var raw rawSearchResult
	q := url.Values{
		"id":       {playlistID},
		"page":     {strconv.Itoa(page)},
		"pageSize": {strconv.Itoa(pageSize)},
	}
	if err := c.doGet(ctx, "/playlist/get-playlist-works", q, &raw); err != nil {
		return Asmr_SearchResult{}, err
	}
	return convertSearchResult(raw), nil
}

// AddWorksToPlaylist 把作品加入播放列表（body 恰 {id, works:[int]}）。
func (c *Client) AddWorksToPlaylist(ctx context.Context, playlistID string, workIDs []int) error {
	return c.doPost(ctx, "/playlist/add-works-to-playlist", map[string]any{"id": playlistID, "works": workIDs}, nil)
}

// RemoveWorksFromPlaylist 把作品移出播放列表（body 恰 {id, works:[int]}）。
func (c *Client) RemoveWorksFromPlaylist(ctx context.Context, playlistID string, workIDs []int) error {
	return c.doPost(ctx, "/playlist/remove-works-from-playlist", map[string]any{"id": playlistID, "works": workIDs}, nil)
}
