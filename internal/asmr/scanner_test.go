package asmr

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	apiasmr "seanime/internal/api/asmr"
	"seanime/internal/database/db"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// findTrackByPath 在音轨树中按精确 Path 递归查找节点（用于断言回填结果）。
func findTrackByPath(tracks []apiasmr.Asmr_Track, path string) *apiasmr.Asmr_Track {
	for i := range tracks {
		if tracks[i].Path == path {
			return &tracks[i]
		}
		if found := findTrackByPath(tracks[i].Tracks, path); found != nil {
			return found
		}
	}
	return nil
}

func TestNormalizeRJ(t *testing.T) {
	cases := []struct {
		in     string
		want   string
		wantOK bool
	}{
		{"RJ01234567", "RJ01234567", true},
		{"rj289543", "RJ289543", true},
		{"RJ12345", "RJ12345", true},
		{"RJ1234", "", false},      // 4 位数字，不匹配
		{"RJ123456789", "", false}, // 9 位数字，不匹配
		{"notrj", "", false},
		{"RJ12345 标题 with spaces", "RJ12345", true},
		{"rj00012345", "RJ00012345", true},
		{"rj+12345", "", false}, // 数字前非数字
	}
	for _, c := range cases {
		got, ok := NormalizeRJ(c.in)
		require.Equal(t, c.wantOK, ok, "in=%q", c.in)
		if c.wantOK {
			require.Equal(t, c.want, got, "in=%q", c.in)
		}
	}
}

func makeRJDir(t *testing.T, root, name string) string {
	t.Helper()
	dir := filepath.Join(root, name)
	require.NoError(t, os.MkdirAll(dir, 0o755))
	return dir
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
}

func TestScanDegradationAndCounts(t *testing.T) {
	root := t.TempDir()
	rjDir := makeRJDir(t, root, "RJ01234567 某作品")
	writeFile(t, filepath.Join(rjDir, "01.mp3"), "a")
	writeFile(t, filepath.Join(rjDir, "sub", "02.wav"), "bb")
	writeFile(t, filepath.Join(rjDir, "cover.jpg"), "img") // 非音频，忽略
	makeRJDir(t, root, "random-folder")                    // 不匹配
	require.NoError(t, os.WriteFile(filepath.Join(root, "note.txt"), []byte("x"), 0o644))

	s := NewScanner(root, nil, nil, nil)
	entries, err := s.Scan(context.Background())
	require.NoError(t, err)
	require.Len(t, entries, 1)

	e := entries[0]
	require.Equal(t, "RJ01234567", e.RjID)
	require.Equal(t, "RJ01234567 某作品", e.Title, "降级：title=目录名（无 client）")
	require.Equal(t, 2, e.TrackCount, "应统计 .mp3 + .wav（jpg 忽略）")
	require.Equal(t, int64(3), e.TotalSizeBytes, "a(1)+bb(2)")
}

func TestBuildLocalTrackTree(t *testing.T) {
	root := t.TempDir()
	rjDir := makeRJDir(t, root, "RJ01234567")
	writeFile(t, filepath.Join(rjDir, "01.mp3"), "a")
	writeFile(t, filepath.Join(rjDir, "folderA", "02.flac"), "bb")
	writeFile(t, filepath.Join(rjDir, "folderA", "sub", "03.ogg"), "ccc")

	files, err := collectAudioFiles(rjDir)
	require.NoError(t, err)
	require.Len(t, files, 3)

	tree, basenames := buildLocalTrackTree(rjDir, files)
	require.Len(t, tree, 2) // 顶层：01.mp3 + folderA

	var audio01, folderA *apiasmr.Asmr_Track
	for i := range tree {
		if tree[i].Type == "audio" {
			audio01 = &tree[i]
		} else if tree[i].Type == "folder" {
			folderA = &tree[i]
		}
	}
	require.NotNil(t, audio01)
	require.Equal(t, "01", audio01.Title, "title=去扩展")
	require.FileExists(t, audio01.LocalPath)
	require.Equal(t, "01.mp3", audio01.Path, "Path=相对 RJ 目录的 '/' 分隔路径（前端完听上报用）")
	require.NotNil(t, folderA)
	require.Len(t, folderA.Tracks, 2) // folderA 顶层：02.flac + sub
	// 嵌套叶子的 Path 含目录层级
	var nested *apiasmr.Asmr_Track
	for i := range folderA.Tracks {
		if folderA.Tracks[i].Type == "folder" && len(folderA.Tracks[i].Tracks) > 0 {
			nested = &folderA.Tracks[i].Tracks[0]
		}
	}
	require.NotNil(t, nested)
	require.Equal(t, "folderA/sub/03.ogg", nested.Path)

	_, hasOgg := basenames["03.ogg"]
	require.True(t, hasOgg, "basename 集合应含小写含扩展名")
	_, hasJpg := basenames["cover.jpg"]
	require.False(t, hasJpg, "非音频不应计入")
}

func TestMergeTrackTrees(t *testing.T) {
	local := []apiasmr.Asmr_Track{
		{Title: "01", Type: "audio", LocalPath: "/x/01.mp3"},
		{Title: "folder", Type: "folder", Tracks: []apiasmr.Asmr_Track{
			{Title: "02", Type: "audio", LocalPath: "/x/folder/02.wav"},
		}},
	}
	localBasenames := map[string]struct{}{
		"01.mp3": {},
		"02.wav": {},
		"03.ogg": {},
	}
	online := []apiasmr.Asmr_Track{
		{Title: "01.mp3", Type: "audio", MediaDownloadURL: "http://d/01.mp3"}, // 本地已有 → 跳过
		{Title: "04.opus", Type: "audio", MediaDownloadURL: "http://d/04.opus"}, // 在线独有 → 并入
		{Title: "folder", Type: "folder", Tracks: []apiasmr.Asmr_Track{
			{Title: "02.wav", Type: "audio", MediaDownloadURL: "http://d/02.wav"}, // 本地已有 → 跳过
			{Title: "05.m4a", Type: "audio", MediaDownloadURL: "http://d/05.m4a"}, // 在线独有 → 并入子目录
		}},
	}

	merged := mergeTrackTrees(local, online, localBasenames)

	require.Len(t, merged, 3) // 顶层：01(本地), folder(本地), 04.opus(在线)
	var found04, foundFolder bool
	for i := range merged {
		if merged[i].Title == "04.opus" {
			found04 = true
			require.Equal(t, "http://d/04.opus", merged[i].MediaDownloadURL)
			require.Empty(t, merged[i].LocalPath, "在线节点无 localPath")
		}
		if merged[i].Type == "folder" {
			foundFolder = true
			require.Len(t, merged[i].Tracks, 2)
		}
	}
	require.True(t, found04, "在线独有音频应并入顶层")
	require.True(t, foundFolder, "本地 folder 应保留")
}

// TestGetWorkBackfillsCompleted 验证 3.2a：GetWork 返回的音轨树按 DB 逐轨状态回填
// Completed 字段（契约 §0.1/D1-D4）。使用真实临时 DB（不 mock），写入 completion 后调 GetWork 断言。
func TestGetWorkBackfillsCompleted(t *testing.T) {
	root := t.TempDir()
	rjID := "RJ01234567"
	rjDir := makeRJDir(t, root, rjID)

	// 真实音频文件：扁平 + 单层 folder + 深层嵌套 folder（不能只测一层）
	writeFile(t, filepath.Join(rjDir, "01.mp3"), "a")
	writeFile(t, filepath.Join(rjDir, "04.mp3"), "d") // 完全无 DB 记录 → 应为 false
	writeFile(t, filepath.Join(rjDir, "folderA", "02.wav"), "bb")
	writeFile(t, filepath.Join(rjDir, "folderA", "sub", "03.ogg"), "ccc") // 深层叶子

	// 真实临时 DB（不 mock）
	logger := zerolog.Nop()
	database, err := db.NewDatabase(t.TempDir(), "test", &logger)
	require.NoError(t, err)
	t.Cleanup(func() { _ = database.Close() })

	// 写入完听状态：01.mp3 与 深层 folderA/sub/03.ogg 为 true；02.wav 显式记 false
	_, err = database.UpsertAsmrTrackState(rjID, "01.mp3", true)
	require.NoError(t, err)
	_, err = database.UpsertAsmrTrackState(rjID, "folderA/sub/03.ogg", true)
	require.NoError(t, err)
	_, err = database.UpsertAsmrTrackState(rjID, "folderA/02.wav", false)
	require.NoError(t, err)
	// 04.mp3 不写任何记录

	s := NewScanner(root, nil, database, &logger)
	work, ok := s.GetWork(context.Background(), rjID)
	require.True(t, ok, "本地目录存在应返回作品")
	require.NotNil(t, work)

	// ① 本地扁平音轨被正确回填为 true
	got01 := findTrackByPath(work.Tracks, "01.mp3")
	require.NotNil(t, got01, "应存在 01.mp3 叶子")
	require.True(t, got01.Completed, "01.mp3 应为 completed=true")

	// ③ 深层嵌套 folder 下的叶子也被回填
	gotDeep := findTrackByPath(work.Tracks, "folderA/sub/03.ogg")
	require.NotNil(t, gotDeep, "应存在 folderA/sub/03.ogg 叶子")
	require.True(t, gotDeep.Completed, "深层叶子 folderA/sub/03.ogg 应为 completed=true")

	// ② 完全未记录的音轨为 false
	got04 := findTrackByPath(work.Tracks, "04.mp3")
	require.NotNil(t, got04, "应存在 04.mp3 叶子")
	require.False(t, got04.Completed, "无 DB 记录的 04.mp3 应为 false")

	// ② 显式记录为 false 的也应为 false
	got02 := findTrackByPath(work.Tracks, "folderA/02.wav")
	require.NotNil(t, got02, "应存在 folderA/02.wav 叶子")
	require.False(t, got02.Completed, "显式 completed=false 的 02.wav 应为 false")

	// D4 回归：聚合 ListenedCount 仍正确（仅 2 条 completed=true）
	require.Equal(t, 2, work.ListenedCount, "ListenedCount 应等于 completed=true 的轨数")
}
