package manga

import (
	"testing"

	"seanime/internal/api/bangumi"
	"seanime/internal/media"
	"seanime/internal/platforms/platform"
	"seanime/internal/testmocks"
	"seanime/internal/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewCollectionFromBangumiNilScoreProgress 是 03.6-DIAGNOSIS 故障 1 的回归测试：
// Bangumi 适配器对「未评分 / 无进度」条目返回 nil 的 Progress/Score，
// 旧 NewCollection 在 collection.go:84-85 裸解引用 *entry.Progress / *entry.Score 导致 panic。
//
// Seam：用真实 Bangumi 适配器（MangaCollectionFromUserCollections）产出 collection 喂 NewCollection，
// 不手搓被 mock 的 collection。构造一条 Rate=0、EpStatus=0 的收藏（合法且常见状态），
// 断言 NewCollection 不 panic，且结果条目字段可读（nil → 零值，而非崩溃）。
func TestNewCollectionFromBangumiNilScoreProgress(t *testing.T) {
	// 真实适配器的入参：一条「未评分 / 无进度」的书籍收藏。
	res := &bangumi.UserCollectionsResult{
		Total: 1,
		Data: []bangumi.UserCollection{
			{
				SubjectID: 1,
				Type:      bangumi.CollectionDoing, // → Current
				Rate:      0,                       // → Score = nil
				EpStatus:  0,                       // → Progress = nil
				Subject: &bangumi.Subject{
					ID:     1,
					Type:   1, // 书籍
					Name:   "Test Manga",
					NameCN: "测试漫画",
				},
			},
		},
	}

	// 走真实适配器（不是手搓的 collection）。
	col := bangumi.MangaCollectionFromUserCollections(res)
	require.NotNil(t, col)

	var capturedPanic interface{}
	func() {
		defer func() {
			capturedPanic = recover()
		}()

		coll, err := NewCollection(&NewCollectionOptions{
			MangaCollection: col,
			// PlatformRef 仅为满足 NewCollection 的入参契约（系统边界依赖），
			// 与本次要测的「适配器→collection 转换」无关。
			PlatformRef: util.NewRef[platform.Platform](testmocks.NewFakePlatformBuilder().Build()),
		})
		require.NoError(t, err)
		require.NotNil(t, coll)

		// 断言：未评分/无进度的条目被安全转换，字段可读（零值而非 panic）。
		require.Len(t, coll.Lists, 1)
		require.Len(t, coll.Lists[0].Entries, 1)
		entry := coll.Lists[0].Entries[0]
		require.NotNil(t, entry.EntryListData)
		assert.Equal(t, 1, entry.MediaId)
		// nil 指针应安全转为对应值类型的零值。
		assert.Equal(t, 0, entry.EntryListData.Progress)  // nil → 0
		assert.Equal(t, 0.0, entry.EntryListData.Score)   // nil → 0.0
		assert.Equal(t, media.MediaListStatusCurrent, *entry.EntryListData.Status)
	}()

	require.Nil(t, capturedPanic, "NewCollection must not panic on nil Progress/Score produced by the Bangumi adapter")
}
