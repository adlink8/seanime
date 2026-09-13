package asmr

import (
	"path/filepath"
	"strings"

	apiasmr "seanime/internal/api/asmr"
)

// buildLocalTrackTree 由音频文件列表构建本地音轨树（folder 递归嵌套）。
// 每个音频叶子：title=文件名(去扩展)，type="audio"，localPath=绝对路径（契约 §2）。
// 返回顶层节点切片与本地音频 basename（小写，含扩展）集合，供在线合并去重。
func buildLocalTrackTree(rjDir string, files []audioFile) ([]apiasmr.Asmr_Track, map[string]struct{}) {
	// 稳定排序：按相对路径，使目录结构有序
	roots := make([]apiasmr.Asmr_Track, 0, 8)
	basenames := make(map[string]struct{}, len(files))

	for _, f := range files {
		basenames[strings.ToLower(f.baseName)] = struct{}{}
		parts := strings.Split(f.relPath, "/")
		insertLocalTrack(&roots, parts, f.absPath, f.relPath)
	}
	return roots, basenames
}

// insertLocalTrack 将单个音频文件按相对路径 parts 插入树（parts 末段为文件名）。
// relPath 为 '/' 分隔的完整相对路径，写入叶子的 Path 字段（前端完听上报用）。
func insertLocalTrack(nodes *[]apiasmr.Asmr_Track, parts []string, absPath, relPath string) {
	if len(parts) == 1 {
		base := parts[0]
		title := strings.TrimSuffix(base, filepath.Ext(base))
		*nodes = append(*nodes, apiasmr.Asmr_Track{
			Title:     title,
			Type:      "audio",
			LocalPath: absPath,
			Path:      relPath,
		})
		return
	}
	folderTitle := parts[0]
	var folder *apiasmr.Asmr_Track
	for i := range *nodes {
		if (*nodes)[i].Type == "folder" && (*nodes)[i].Title == folderTitle {
			folder = &(*nodes)[i]
			break
		}
	}
	if folder == nil {
		*nodes = append(*nodes, apiasmr.Asmr_Track{Title: folderTitle, Type: "folder"})
		folder = &(*nodes)[len(*nodes)-1]
	}
	insertLocalTrack(&folder.Tracks, parts[1:], absPath, relPath)
}

// mergeTrackTrees 合并本地树与在线树（契约 §2：本地优先）。
// 本地树原样保留；在线音频节点若其 basename（小写，含扩展）不在本地集合里，则按 folder 结构并入结果。
//
// TODO(M3.1 编排裁决): 合并键当前用"basename（文件名去扩展后比对含扩展的 basename）"近似匹配，
// 未做完整相对路径（目录层级）精确匹配——在线 folder 标题与本地目录名往往不同，basename 去重已覆盖绝大多数"已下载"判定。
// 若需严格按相对路径合并，待编排层确认在线树相对路径构造规则后调整。
func mergeTrackTrees(local, online []apiasmr.Asmr_Track, localBasenames map[string]struct{}) []apiasmr.Asmr_Track {
	result := append([]apiasmr.Asmr_Track{}, local...)
	mergeOnlineInto(&result, online, localBasenames)
	return result
}

// markTracksCompleted 递归回填逐轨完听状态（契约 3.2a / D2-D4）。
// completedPaths 为 DB 中 completed=true 的 track_path 集合。按节点 Path 精确匹配
// completedPaths（D3，不做 basename 近似）。仅本地音轨有 Path，故在线音轨 / 线上搜索路径
// 天然不会被回填（Asmr_WorkDetail 永不回填，预期行为，非 bug，见 types.go §8）。
func markTracksCompleted(nodes []apiasmr.Asmr_Track, completedPaths map[string]struct{}) {
	for i := range nodes {
		if nodes[i].Path != "" {
			if _, ok := completedPaths[nodes[i].Path]; ok {
				nodes[i].Completed = true
			}
		}
		if len(nodes[i].Tracks) > 0 {
			markTracksCompleted(nodes[i].Tracks, completedPaths)
		}
	}
}

// mergeOnlineInto 将在线树中"本地没有对应音频"的节点并入 target（保留 folder 嵌套）。
func mergeOnlineInto(target *[]apiasmr.Asmr_Track, nodes []apiasmr.Asmr_Track, localBasenames map[string]struct{}) {
	for _, n := range nodes {
		switch n.Type {
		case "folder":
			var f *apiasmr.Asmr_Track
			for i := range *target {
				if (*target)[i].Type == "folder" && (*target)[i].Title == n.Title {
					f = &(*target)[i]
					break
				}
			}
			if f == nil {
				*target = append(*target, apiasmr.Asmr_Track{Title: n.Title, Type: "folder"})
				f = &(*target)[len(*target)-1]
			}
			mergeOnlineInto(&f.Tracks, n.Tracks, localBasenames)
		case "audio":
			base := strings.ToLower(filepath.Base(n.Title))
			if _, ok := localBasenames[base]; ok {
				continue // 本地优先：已存在则跳过在线节点
			}
			*target = append(*target, n)
		}
		// text/image 节点忽略（不参与合并）
	}
}
