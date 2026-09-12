package anime_test

import (
	"seanime/internal/customsource"
	"seanime/internal/library/anime"
	"seanime/internal/media"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetScheduleItemsFormatsDeduplicates(t *testing.T) {
	// schedule items are merged from all schedule buckets,
	// deduped by media/episode/time
	h := newAnimeTestWrapper(t)

	patchAnimeCollectionEntry(t, h.animeCollection, 154587, media.AnimeCollectionEntryPatch{
		Status:        new(media.MediaListStatusCurrent),
		AiredEpisodes: new(12),
	})
	patchCollectionEntryEpisodeCount(t, h.animeCollection, 154587, 12)

	patchAnimeCollectionEntry(t, h.animeCollection, 146065, media.AnimeCollectionEntryPatch{
		Status:        new(media.MediaListStatusCurrent),
		AiredEpisodes: new(1),
	})
	patchCollectionEntryEpisodeCount(t, h.animeCollection, 146065, 1)
	movieFormat := media.MediaFormatMovie
	patchCollectionEntryFormat(t, h.animeCollection, 146065, movieFormat)
	movieEntry := findCollectionEntryByMediaID(t, h.animeCollection, 146065)
	fallbackTitle := "movie fallback"
	movieEntry.Media.Title.UserPreferred = nil
	movieEntry.Media.Title.English = &fallbackTitle

	// extension-backed ids should not leak into the schedule list.
	extensionEntry := findCollectionEntryByMediaID(t, h.animeCollection, 21)
	extensionID := customsource.GenerateMediaId(1, 99)
	extensionEntry.Media.ID = extensionID
	extensionEntry.Status = new(media.MediaListStatusCurrent)

	// Bangumi 锚点：media 包将分页包装合并为单对象（Media 由切片变为单个 *AnimeSchedule），
	// 原 Ongoing 中的 extension 条目挪至 Preceding，泄漏验证语义不变。
	animeSchedule := &media.AnimeAiringSchedule{
		Ongoing: &media.AnimeAiringSchedule_Ongoing{
			Media: newAnimeSchedule(154587,
				[]*media.AnimeSchedule_Previous_Nodes{newPreviousScheduleNode(1_700_000_100, 11, -100)},
				[]*media.AnimeSchedule_Upcoming_Nodes{newUpcomingScheduleNode(1_700_000_200, 12, 200)},
			),
		},
		OngoingNext: &media.AnimeAiringSchedule_OngoingNext{
			Media: newAnimeSchedule(154587, nil, []*media.AnimeSchedule_Upcoming_Nodes{newUpcomingScheduleNode(1_700_000_200, 12, 200)}),
		},
		Upcoming: &media.AnimeAiringSchedule_Upcoming{
			Media: newAnimeSchedule(146065, nil, []*media.AnimeSchedule_Upcoming_Nodes{newUpcomingScheduleNode(1_700_000_300, 1, 300)}),
		},
		Preceding: &media.AnimeAiringSchedule_Preceding{
			Media: newAnimeSchedule(extensionID, nil, []*media.AnimeSchedule_Upcoming_Nodes{newUpcomingScheduleNode(1_700_000_050, 1, 50)}),
		},
	}

	items := anime.GetScheduleItems(animeSchedule, h.animeCollection)

	require.Len(t, items, 3)
	require.Len(t, findScheduleItems(items, 154587, 12), 1)
	require.Empty(t, findScheduleItems(items, extensionID, 1))

	previousItem := findScheduleItem(t, items, 154587, 11)
	require.Equal(t, time.Unix(1_700_000_100, 0).UTC(), previousItem.DateTime)
	require.Equal(t, previousItem.DateTime.Format("15:04"), previousItem.Time)
	require.False(t, previousItem.IsSeasonFinale)
	require.False(t, previousItem.IsMovie)

	finaleItem := findScheduleItem(t, items, 154587, 12)
	require.True(t, finaleItem.IsSeasonFinale)

	movieItem := findScheduleItem(t, items, 146065, 1)
	require.Equal(t, fallbackTitle, movieItem.Title)
	require.True(t, movieItem.IsMovie)
	require.True(t, movieItem.IsSeasonFinale)
}

func TestGetScheduleItemsHandlesNilInputs(t *testing.T) {
	// nil inputs should just give the caller an empty slice instead of exploding.
	require.Empty(t, anime.GetScheduleItems(nil, nil))
}

func newAnimeSchedule(mediaID int, previous []*media.AnimeSchedule_Previous_Nodes, upcoming []*media.AnimeSchedule_Upcoming_Nodes) *media.AnimeSchedule {
	ret := &media.AnimeSchedule{ID: mediaID}
	if previous != nil {
		ret.Previous = &media.AnimeSchedule_Previous{Nodes: previous}
	}
	if upcoming != nil {
		ret.Upcoming = &media.AnimeSchedule_Upcoming{Nodes: upcoming}
	}
	return ret
}

func newPreviousScheduleNode(airingAt int, episode int, timeUntilAiring int) *media.AnimeSchedule_Previous_Nodes {
	return &media.AnimeSchedule_Previous_Nodes{
		AiringAt:        airingAt,
		Episode:         episode,
		TimeUntilAiring: timeUntilAiring,
	}
}

func newUpcomingScheduleNode(airingAt int, episode int, timeUntilAiring int) *media.AnimeSchedule_Upcoming_Nodes {
	return &media.AnimeSchedule_Upcoming_Nodes{
		AiringAt:        airingAt,
		Episode:         episode,
		TimeUntilAiring: timeUntilAiring,
	}
}

func findScheduleItem(t *testing.T, items []*anime.ScheduleItem, mediaID int, episodeNumber int) *anime.ScheduleItem {
	t.Helper()
	matching := findScheduleItems(items, mediaID, episodeNumber)
	require.Len(t, matching, 1)
	return matching[0]
}

func findScheduleItems(items []*anime.ScheduleItem, mediaID int, episodeNumber int) []*anime.ScheduleItem {
	ret := make([]*anime.ScheduleItem, 0)
	for _, item := range items {
		if item.MediaId == mediaID && item.EpisodeNumber == episodeNumber {
			ret = append(ret, item)
		}
	}
	return ret
}
