package anime

import (
	"context"
	"seanime/internal/media"
	"seanime/internal/platforms/platform"
	"seanime/internal/util/comparison"
	"seanime/internal/util/limiter"
	"seanime/internal/util/result"

	"github.com/samber/lo"
)

type NormalizedMedia struct {
	ID          int
	IdMal       *int
	Title       *NormalizedMediaTitle
	Synonyms    []*string
	Format      *media.MediaFormat
	Status      *media.MediaStatus
	Season      *media.MediaSeason
	Year        *int
	StartDate   *NormalizedMediaDate
	Episodes    *int
	BannerImage *string
	CoverImage  *NormalizedMediaCoverImage
	//Relations         *media.CompleteAnimeById_Media_CompleteAnime_Relations
	NextAiringEpisode *NormalizedMediaNextAiringEpisode
	// Whether it was fetched from AniList
	fetched bool
}

type NormalizedMediaTitle struct {
	Romaji        *string
	English       *string
	Native        *string
	UserPreferred *string
}

type NormalizedMediaDate struct {
	Year  *int
	Month *int
	Day   *int
}

type NormalizedMediaCoverImage struct {
	ExtraLarge *string
	Large      *string
	Medium     *string
	Color      *string
}

type NormalizedMediaNextAiringEpisode struct {
	AiringAt        int
	TimeUntilAiring int
	Episode         int
}

type NormalizedMediaCache struct {
	*result.Cache[int, *NormalizedMedia]
}

func NewNormalizedMedia(m *media.Anime) *NormalizedMedia {
	var startDate *NormalizedMediaDate
	if m.GetStartDate() != nil {
		startDate = &NormalizedMediaDate{
			Year:  m.GetStartDate().GetYear(),
			Month: m.GetStartDate().GetMonth(),
			Day:   m.GetStartDate().GetDay(),
		}
	}

	var title *NormalizedMediaTitle
	if m.GetTitle() != nil {
		title = &NormalizedMediaTitle{
			Romaji:        m.GetTitle().GetRomaji(),
			English:       m.GetTitle().GetEnglish(),
			Native:        m.GetTitle().GetNative(),
			UserPreferred: m.GetTitle().GetUserPreferred(),
		}
	}

	var coverImage *NormalizedMediaCoverImage
	if m.GetCoverImage() != nil {
		coverImage = &NormalizedMediaCoverImage{
			ExtraLarge: m.GetCoverImage().GetExtraLarge(),
			Large:      m.GetCoverImage().GetLarge(),
			Medium:     m.GetCoverImage().GetMedium(),
			Color:      m.GetCoverImage().GetColor(),
		}
	}

	var nextAiringEpisode *NormalizedMediaNextAiringEpisode
	if m.GetNextAiringEpisode() != nil {
		nextAiringEpisode = &NormalizedMediaNextAiringEpisode{
			AiringAt:        m.GetNextAiringEpisode().GetAiringAt(),
			TimeUntilAiring: m.GetNextAiringEpisode().GetTimeUntilAiring(),
			Episode:         m.GetNextAiringEpisode().GetEpisode(),
		}
	}

	return &NormalizedMedia{
		ID:                m.GetID(),
		IdMal:             m.GetIDMal(),
		Title:             title,
		Synonyms:          m.GetSynonyms(),
		Format:            m.GetFormat(),
		Status:            m.GetStatus(),
		Season:            m.GetSeason(),
		Year:              m.GetSeasonYear(),
		StartDate:         startDate,
		Episodes:          m.GetEpisodes(),
		BannerImage:       m.GetBannerImage(),
		CoverImage:        coverImage,
		NextAiringEpisode: nextAiringEpisode,
		fetched:           true,
	}
}

// NewNormalizedMediaFromOfflineDB creates a NormalizedMedia from the anime-offline-database.
// The media is marked as not fetched (fetched=false) since it lacks some AniList-specific data.
func NewNormalizedMediaFromOfflineDB(
	id int,
	idMal *int,
	title *NormalizedMediaTitle,
	synonyms []*string,
	format *media.MediaFormat,
	status *media.MediaStatus,
	season *media.MediaSeason,
	year *int,
	startDate *NormalizedMediaDate,
	episodes *int,
	coverImage *NormalizedMediaCoverImage,
) *NormalizedMedia {
	return &NormalizedMedia{
		ID:         id,
		IdMal:      idMal,
		Title:      title,
		Synonyms:   synonyms,
		Format:     format,
		Status:     status,
		Season:     season,
		Year:       year,
		StartDate:  startDate,
		Episodes:   episodes,
		CoverImage: coverImage,
		fetched:    false,
	}
}

// FetchNormalizedMedia fetches the complete anime (with relations) via the platform layer
// （Bangumi 锚点：原实现走 AnilistClient.CompleteAnimeByID，现统一走 Platform.GetAnimeWithRelations）.
func FetchNormalizedMedia(platform platform.Platform, l *limiter.Limiter, cache *media.CompleteAnimeCache, m *NormalizedMedia) error {
	if platform == nil || m == nil {
		return nil
	}

	if m.fetched {
		return nil
	}

	if cache != nil {
		if complete, found := cache.Get(m.ID); found {
			*m = *NewNormalizedMedia(complete.ToAnime())
		}
	}

	l.Wait()
	complete, err := platform.GetAnimeWithRelations(context.Background(), m.ID)
	if err != nil {
		return err
	}

	if cache != nil && complete != nil {
		cache.Set(m.ID, complete)
	}
	*m = *NewNormalizedMedia(complete.ToAnime())
	m.fetched = true
	return nil
}

func NewNormalizedMediaCache() *NormalizedMediaCache {
	return &NormalizedMediaCache{result.NewCache[int, *NormalizedMedia]()}
}

// Helper methods

func (m *NormalizedMedia) GetTitleSafe() string {
	if m.Title == nil {
		return ""
	}
	if m.Title.UserPreferred != nil {
		return *m.Title.UserPreferred
	}
	if m.Title.English != nil {
		return *m.Title.English
	}
	if m.Title.Romaji != nil {
		return *m.Title.Romaji
	}
	if m.Title.Native != nil {
		return *m.Title.Native
	}
	return ""
}

func (m *NormalizedMedia) HasRomajiTitle() bool {
	return m.Title != nil && m.Title.Romaji != nil
}

func (m *NormalizedMedia) HasEnglishTitle() bool {
	return m.Title != nil && m.Title.English != nil
}

func (m *NormalizedMedia) HasSynonyms() bool {
	return len(m.Synonyms) > 0
}

func (m *NormalizedMedia) GetAllTitles() []*string {
	titles := make([]*string, 0)
	if m.Title == nil {
		return titles
	}
	if m.Title.Romaji != nil {
		titles = append(titles, m.Title.Romaji)
	}
	if m.Title.English != nil {
		titles = append(titles, m.Title.English)
	}
	if m.Title.Native != nil {
		titles = append(titles, m.Title.Native)
	}
	if m.Title.UserPreferred != nil {
		titles = append(titles, m.Title.UserPreferred)
	}
	titles = append(titles, m.Synonyms...)
	return titles
}

// GetPossibleSeasonNumber returns the possible season number for that media and -1 if it doesn't have one.
// It looks at the synonyms and returns the highest season number found.
func (m *NormalizedMedia) GetPossibleSeasonNumber() int {
	if m == nil || len(m.Synonyms) == 0 {
		return -1
	}
	titles := lo.Filter(m.Synonyms, func(s *string, i int) bool { return comparison.ValueContainsSeason(*s) })
	if m.HasEnglishTitle() {
		titles = append(titles, m.Title.English)
	}
	if m.HasRomajiTitle() {
		titles = append(titles, m.Title.Romaji)
	}
	seasons := lo.Map(titles, func(s *string, i int) int { return comparison.ExtractSeasonNumber(*s) })
	return lo.Max(seasons)
}

// FetchMediaTree populates the relation tree with the given media's sequels/prequels via the platform layer.
// Bangumi 的 GetAnimeWithRelations 一次返回主题及其一层相关条目，因此用 BFS 逐层扩展；
// 第二层起返回的 relations 与首层相同，依赖 tree.Has 去重自然收敛。
// rel 为 all 时两个方向都收集（原 AniList 实现按边方向收敛，此处语义稍宽，保守可接受）。
func (m *NormalizedMedia) FetchMediaTree(
	rel media.FetchMediaTreeRelation,
	platform platform.Platform,
	rl *limiter.Limiter,
	tree *media.CompleteAnimeRelationTree,
	cache *media.CompleteAnimeCache,
) error {
	if m == nil || platform == nil {
		return nil
	}

	rl.Wait()
	root, err := platform.GetAnimeWithRelations(context.Background(), m.ID)
	if err != nil {
		return err
	}
	if root == nil {
		return nil
	}

	queue := []*media.CompleteAnime{root}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur == nil {
			continue
		}
		if tree.Has(cur.ID) {
			cache.Set(cur.ID, cur)
			continue
		}
		cache.Set(cur.ID, cur)
		tree.Set(cur.ID, cur)

		edges := cur.GetRelations().GetEdges()
		for _, edge := range edges {
			if edge == nil || edge.GetNode() == nil || edge.GetRelationType() == nil {
				continue
			}
			node := edge.GetNode()
			rt := *edge.GetRelationType()
			// 方向过滤：sequels 只追续作，prequels 只追前传，all 收两者
			if rel == media.FetchMediaTreeSequels && rt != media.MediaRelationSequel {
				continue
			}
			if rel == media.FetchMediaTreePrequels && rt != media.MediaRelationPrequel {
				continue
			}
			if rt != media.MediaRelationSequel && rt != media.MediaRelationPrequel {
				continue
			}
			if status := node.GetStatus(); status == nil || *status == media.MediaStatusNotYetReleased {
				continue
			}
			if !edge.IsBroadRelationFormat() || tree.Has(node.ID) {
				continue
			}
			rl.Wait()
			next, err := platform.GetAnimeWithRelations(context.Background(), node.ID)
			if err != nil || next == nil {
				continue
			}
			queue = append(queue, next)
		}
	}
	return nil
}

// GetCurrentEpisodeCount returns the current episode number for that media and -1 if it doesn't have one.
// i.e. -1 is returned if the media has no episodes AND the next airing episode is not set.
func (m *NormalizedMedia) GetCurrentEpisodeCount() int {
	ceil := -1
	if m.Episodes != nil {
		ceil = *m.Episodes
	}
	if m.NextAiringEpisode != nil {
		if m.NextAiringEpisode.Episode > 0 {
			ceil = m.NextAiringEpisode.Episode - 1
		}
	}
	return ceil
}

// GetTotalEpisodeCount returns the total episode number for that media and -1 if it doesn't have one
func (m *NormalizedMedia) GetTotalEpisodeCount() int {
	ceil := -1
	if m.Episodes != nil {
		ceil = *m.Episodes
	}
	return ceil
}
