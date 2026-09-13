package asmr

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	apiasmr "seanime/internal/api/asmr"

	"github.com/stretchr/testify/require"
)

// 契约 03.2c Wave A：filterTrackerCandidates 纯函数 TDD。
// 覆盖：① 命中 ② 社团不符 ③ 窗口外 ④ Release 非法/空串弃用并计数（不 panic）。
func TestFilterTrackerCandidates(t *testing.T) {
	windowStart := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	works := []apiasmr.Asmr_Work{
		{RjID: "RJ001", Circle: "Whirlpool", ReleaseDate: "2026-09-10"},  // ① 命中
		{RjID: "RJ002", Circle: "OtherCircle", ReleaseDate: "2026-09-10"}, // ② 社团不符
		{RjID: "RJ003", Circle: "Whirlpool", ReleaseDate: "2020-01-01"},   // ③ 窗口外
		{RjID: "RJ004", Circle: "Whirlpool", ReleaseDate: ""},             // ④ 空串
		{RjID: "RJ005", Circle: "Whirlpool", ReleaseDate: "not-a-date"},   // ④ 非法日期
		{RjID: "RJ006", Circle: "Whirlpool", ReleaseDate: "2026-09-01"},   // ① 边界：恰好等于窗口起点
	}

	matched, invalid := filterTrackerCandidates(works, "Whirlpool", windowStart)

	require.Len(t, matched, 2, "应命中 RJ001 与 RJ006（窗口起点含边界）")
	require.Equal(t, "RJ001", matched[0].RjID)
	require.Equal(t, "RJ006", matched[1].RjID)
	require.Equal(t, 2, invalid, "空串与非法日期应计入弃用")

	// 输入不被修改（纯函数）
	require.Equal(t, "Whirlpool", works[0].Circle)
}

// hasLocalFiles：目录不存在 / 空目录 → false；任意层级含文件 → true（D3①）。
func TestHasLocalFiles(t *testing.T) {
	root := t.TempDir()

	ok, err := hasLocalFiles(filepath.Join(root, "RJ99999999"))
	require.NoError(t, err)
	require.False(t, ok, "目录不存在 → false")

	empty := filepath.Join(root, "RJ88888888")
	require.NoError(t, os.MkdirAll(empty, 0o755))
	ok, err = hasLocalFiles(empty)
	require.NoError(t, err)
	require.False(t, ok, "空目录 → false")

	nested := filepath.Join(root, "RJ77777777", "folder1")
	require.NoError(t, os.MkdirAll(nested, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(nested, "01.mp3"), []byte("x"), 0o644))
	ok, err = hasLocalFiles(filepath.Join(root, "RJ77777777"))
	require.NoError(t, err)
	require.True(t, ok, "嵌套子目录含文件 → true")
}
