package scanner

import (
	"seanime/internal/media"
)

// ensureAnimeCollectionWithRelationsEntry 确保集合中存在指定条目并应用补丁。
// Bangumi 锚点：原 media.EnsureAnimeCollectionWithRelationsEntry 依赖 AnilistClient 拉取缺失条目，
// 该概念已随 AniList client 裁剪；测试版改为在条目缺失时插入仅含 ID 的最小条目。
// 依赖 test/data AniList fixture 的用例属 CI 跳过名单。
func ensureAnimeCollectionWithRelationsEntry(collection *media.AnimeCollectionWithRelations, mediaID int, patch media.AnimeCollectionEntryPatch) {
	if collection == nil {
		return
	}
	if _, found := collection.FindAnime(mediaID); found {
		media.PatchAnimeCollectionWithRelationsEntry(collection, mediaID, patch)
		return
	}
	if collection.MediaListCollection == nil {
		collection.MediaListCollection = &media.AnimeCollectionWithRelations_MediaListCollection{}
	}
	if len(collection.MediaListCollection.Lists) == 0 {
		collection.MediaListCollection.Lists = append(collection.MediaListCollection.Lists, &media.AnimeCollectionWithRelations_MediaListCollection_Lists{})
	}
	list := collection.MediaListCollection.Lists[0]
	entry := &media.AnimeCollectionWithRelations_MediaListCollection_Lists_Entries{
		ID:    mediaID,
		Media: &media.CompleteAnime{ID: mediaID},
	}
	if patch.Status != nil {
		entry.Status = patch.Status
	}
	list.Entries = append(list.Entries, entry)
}
