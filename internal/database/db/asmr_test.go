package db

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func newTestDB(t *testing.T) *Database {
	t.Helper()
	logger := zerolog.Nop()
	database, err := NewDatabase(t.TempDir(), "test", &logger)
	require.NoError(t, err)
	t.Cleanup(func() { _ = database.Close() })
	return database
}

func TestAsmrWorkStateUpsert(t *testing.T) {
	database := newTestDB(t)

	ok, found := database.GetAsmrWorkState("RJ12345")
	require.False(t, found, "初始不应存在")
	require.Nil(t, ok)

	require.NoError(t, database.UpsertAsmrWorkState("RJ12345", true))
	ws, found := database.GetAsmrWorkState("RJ12345")
	require.True(t, found)
	require.True(t, ws.Favorite)

	// 更新为 false（唯一键不新增行）
	require.NoError(t, database.UpsertAsmrWorkState("RJ12345", false))
	ws, _ = database.GetAsmrWorkState("RJ12345")
	require.False(t, ws.Favorite)
}

func TestAsmrTrackStateUpsertAndCount(t *testing.T) {
	database := newTestDB(t)

	count, err := database.UpsertAsmrTrackState("RJ12345", "01.mp3", true)
	require.NoError(t, err)
	require.Equal(t, 1, count)

	count, err = database.UpsertAsmrTrackState("RJ12345", "02.mp3", true)
	require.NoError(t, err)
	require.Equal(t, 2, count)

	// 同音轨标记未完成 → listenedCount 回退到 1
	count, err = database.UpsertAsmrTrackState("RJ12345", "01.mp3", false)
	require.NoError(t, err)
	require.Equal(t, 1, count)

	// 幂等：重复 completed 不增计数
	count, err = database.UpsertAsmrTrackState("RJ12345", "02.mp3", true)
	require.NoError(t, err)
	require.Equal(t, 1, count)

	states, err := database.GetAsmrTrackStates("RJ12345")
	require.NoError(t, err)
	require.Len(t, states, 2)
}
