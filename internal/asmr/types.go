package asmr

import (
	apiasmr "seanime/internal/api/asmr"
)

// LibraryEntry 音声本地库单个 RJ 条目的契约响应形状（契约 §1）。
// JSON 字段名逐字对齐契约；既有 api/asmr 类型字段不改动，这里单独定义子集。
type LibraryEntry struct {
	RjID           string   `json:"rjId"`
	Title          string   `json:"title"`
	Circle         string   `json:"circle"`
	CoverURL       string   `json:"coverUrl"`
	Nsfw           bool     `json:"nsfw"`
	Rating         float64  `json:"rating"`
	ReleaseDate    string   `json:"releaseDate"`
	Cvs            []string `json:"cvs"`
	Tags           []string `json:"tags"`
	TrackCount     int      `json:"trackCount"`
	TotalSizeBytes int64    `json:"totalSizeBytes"`
	IsFavorite     bool     `json:"isFavorite"`
	ListenedCount  int      `json:"listenedCount"`
}

// LibraryWork 音声作品详情（条目 + 合并后的音轨树，契约 §2）。
type LibraryWork struct {
	LibraryEntry
	Tracks []apiasmr.Asmr_Track `json:"tracks"`
}

// LibraryResponse GET /asmr/library 响应根对象。
type LibraryResponse struct {
	LocalDir string         `json:"localDir"`
	Entries  []LibraryEntry `json:"entries"`
}

// audioExtensions 音声本地库音频扩展白名单（契约 §1，不得用 util.IsValidVideoExtension）。
var audioExtensions = map[string]struct{}{
	".mp3":  {},
	".wav":  {},
	".flac": {},
	".m4a":  {},
	".ogg":  {},
	".opus": {},
}

// isAudioExtension 判断扩展名（小写，含点）是否为音频。
func isAudioExtension(ext string) bool {
	_, ok := audioExtensions[ext]
	return ok
}
