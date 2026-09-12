package testmocks

import "seanime/internal/media"

type BaseAnimeBuilder struct {
	anime *media.Anime
}

func NewBaseAnimeBuilder(id int, title string) *BaseAnimeBuilder {
	return &BaseAnimeBuilder{anime: &media.Anime{
		ID:       id,
		IDMal:    new(501),
		Status:   new(media.MediaStatusFinished),
		Type:     new(media.MediaTypeAnime),
		Format:   new(media.MediaFormatTv),
		Episodes: new(12),
		IsAdult:  new(false),
		Title: &media.Anime_Title{
			English: new(title),
			Romaji:  new(title),
		},
		Synonyms: []*string{new(title), new("Sample Anime Season 1")},
		StartDate: &media.Anime_StartDate{
			Year:  new(2024),
			Month: new(1),
			Day:   new(2),
		},
	}}
}

func NewBaseAnime(id int, title string) *media.Anime {
	return NewBaseAnimeBuilder(id, title).Build()
}

func (b *BaseAnimeBuilder) WithIDMal(idMal int) *BaseAnimeBuilder {
	b.anime.IDMal = new(idMal)
	return b
}

func (b *BaseAnimeBuilder) WithSiteURL(siteURL string) *BaseAnimeBuilder {
	b.anime.SiteURL = new(siteURL)
	return b
}

func (b *BaseAnimeBuilder) WithTitles(english string, romaji string, native string, userPreferred string) *BaseAnimeBuilder {
	ensureAnimeTitle(b.anime)
	b.anime.Title.English = new(english)
	b.anime.Title.Romaji = new(romaji)
	b.anime.Title.Native = new(native)
	b.anime.Title.UserPreferred = new(userPreferred)
	return b
}

func (b *BaseAnimeBuilder) WithEnglishTitle(title string) *BaseAnimeBuilder {
	ensureAnimeTitle(b.anime)
	b.anime.Title.English = new(title)
	return b
}

func (b *BaseAnimeBuilder) WithRomajiTitle(title string) *BaseAnimeBuilder {
	ensureAnimeTitle(b.anime)
	b.anime.Title.Romaji = new(title)
	return b
}

func (b *BaseAnimeBuilder) WithNativeTitle(title string) *BaseAnimeBuilder {
	ensureAnimeTitle(b.anime)
	b.anime.Title.Native = new(title)
	return b
}

func (b *BaseAnimeBuilder) WithUserPreferredTitle(title string) *BaseAnimeBuilder {
	ensureAnimeTitle(b.anime)
	b.anime.Title.UserPreferred = new(title)
	return b
}

func (b *BaseAnimeBuilder) WithStatus(status media.MediaStatus) *BaseAnimeBuilder {
	b.anime.Status = new(status)
	return b
}

func (b *BaseAnimeBuilder) WithFormat(format media.MediaFormat) *BaseAnimeBuilder {
	b.anime.Format = new(format)
	return b
}

func (b *BaseAnimeBuilder) WithEpisodes(episodes int) *BaseAnimeBuilder {
	b.anime.Episodes = new(episodes)
	return b
}

func (b *BaseAnimeBuilder) WithIsAdult(isAdult bool) *BaseAnimeBuilder {
	b.anime.IsAdult = new(isAdult)
	return b
}

func (b *BaseAnimeBuilder) WithSynonyms(synonyms ...string) *BaseAnimeBuilder {
	b.anime.Synonyms = stringPointers(synonyms...)
	return b
}

func (b *BaseAnimeBuilder) WithStartDate(year int, month int, day int) *BaseAnimeBuilder {
	b.anime.StartDate = &media.Anime_StartDate{
		Year:  new(year),
		Month: new(month),
		Day:   new(day),
	}
	return b
}

func (b *BaseAnimeBuilder) WithEndDate(year int, month int, day int) *BaseAnimeBuilder {
	b.anime.EndDate = &media.Anime_EndDate{
		Year:  new(year),
		Month: new(month),
		Day:   new(day),
	}
	return b
}

func (b *BaseAnimeBuilder) WithCoverImage(url string) *BaseAnimeBuilder {
	b.anime.CoverImage = &media.Anime_CoverImage{
		ExtraLarge: new(url),
		Large:      new(url),
		Medium:     new(url),
	}
	return b
}

func (b *BaseAnimeBuilder) WithBannerImage(url string) *BaseAnimeBuilder {
	b.anime.BannerImage = new(url)
	return b
}

func (b *BaseAnimeBuilder) WithNextAiringEpisode(episode int, airingAt int, timeUntilAiring int) *BaseAnimeBuilder {
	b.anime.NextAiringEpisode = &media.Anime_NextAiringEpisode{
		Episode:         episode,
		AiringAt:        airingAt,
		TimeUntilAiring: timeUntilAiring,
	}
	return b
}

func (b *BaseAnimeBuilder) Build() *media.Anime {
	return b.anime
}

type BaseMangaBuilder struct {
	manga *media.Manga
}

func NewBaseMangaBuilder(id int, title string) *BaseMangaBuilder {
	return &BaseMangaBuilder{manga: &media.Manga{
		ID:      id,
		Status:  new(media.MediaStatusFinished),
		Type:    new(media.MediaTypeManga),
		Format:  new(media.MediaFormatManga),
		IsAdult: new(false),
		Title: &media.Manga_Title{
			English: new(title),
			Romaji:  new(title),
		},
		Synonyms: []*string{new(title), new(title + " Alternative")},
		StartDate: &media.Manga_StartDate{
			Year: new(2023),
		},
	}}
}

func NewBaseManga(id int, title string) *media.Manga {
	return NewBaseMangaBuilder(id, title).Build()
}

func (b *BaseMangaBuilder) WithIDMal(idMal int) *BaseMangaBuilder {
	b.manga.IDMal = new(idMal)
	return b
}

func (b *BaseMangaBuilder) WithSiteURL(siteURL string) *BaseMangaBuilder {
	b.manga.SiteURL = new(siteURL)
	return b
}

func (b *BaseMangaBuilder) WithTitles(english string, romaji string, native string, userPreferred string) *BaseMangaBuilder {
	ensureMangaTitle(b.manga)
	b.manga.Title.English = new(english)
	b.manga.Title.Romaji = new(romaji)
	b.manga.Title.Native = new(native)
	b.manga.Title.UserPreferred = new(userPreferred)
	return b
}

func (b *BaseMangaBuilder) WithEnglishTitle(title string) *BaseMangaBuilder {
	ensureMangaTitle(b.manga)
	b.manga.Title.English = new(title)
	return b
}

func (b *BaseMangaBuilder) WithRomajiTitle(title string) *BaseMangaBuilder {
	ensureMangaTitle(b.manga)
	b.manga.Title.Romaji = new(title)
	return b
}

func (b *BaseMangaBuilder) WithNativeTitle(title string) *BaseMangaBuilder {
	ensureMangaTitle(b.manga)
	b.manga.Title.Native = new(title)
	return b
}

func (b *BaseMangaBuilder) WithUserPreferredTitle(title string) *BaseMangaBuilder {
	ensureMangaTitle(b.manga)
	b.manga.Title.UserPreferred = new(title)
	return b
}

func (b *BaseMangaBuilder) WithStatus(status media.MediaStatus) *BaseMangaBuilder {
	b.manga.Status = new(status)
	return b
}

func (b *BaseMangaBuilder) WithFormat(format media.MediaFormat) *BaseMangaBuilder {
	b.manga.Format = new(format)
	return b
}

func (b *BaseMangaBuilder) WithChapters(chapters int) *BaseMangaBuilder {
	b.manga.Chapters = new(chapters)
	return b
}

func (b *BaseMangaBuilder) WithVolumes(volumes int) *BaseMangaBuilder {
	b.manga.Volumes = new(volumes)
	return b
}

func (b *BaseMangaBuilder) WithIsAdult(isAdult bool) *BaseMangaBuilder {
	b.manga.IsAdult = new(isAdult)
	return b
}

func (b *BaseMangaBuilder) WithSynonyms(synonyms ...string) *BaseMangaBuilder {
	b.manga.Synonyms = stringPointers(synonyms...)
	return b
}

func (b *BaseMangaBuilder) WithStartDate(year int, month int, day int) *BaseMangaBuilder {
	b.manga.StartDate = &media.Manga_StartDate{
		Year:  new(year),
		Month: new(month),
		Day:   new(day),
	}
	return b
}

func (b *BaseMangaBuilder) WithEndDate(year int, month int, day int) *BaseMangaBuilder {
	b.manga.EndDate = &media.Manga_EndDate{
		Year:  new(year),
		Month: new(month),
		Day:   new(day),
	}
	return b
}

func (b *BaseMangaBuilder) WithCoverImage(url string) *BaseMangaBuilder {
	b.manga.CoverImage = &media.Manga_CoverImage{
		ExtraLarge: new(url),
		Large:      new(url),
		Medium:     new(url),
	}
	return b
}

func (b *BaseMangaBuilder) WithBannerImage(url string) *BaseMangaBuilder {
	b.manga.BannerImage = new(url)
	return b
}

func (b *BaseMangaBuilder) Build() *media.Manga {
	return b.manga
}

func ensureAnimeTitle(anime *media.Anime) {
	if anime.Title == nil {
		anime.Title = &media.Anime_Title{}
	}
}

func ensureMangaTitle(manga *media.Manga) {
	if manga.Title == nil {
		manga.Title = &media.Manga_Title{}
	}
}

func stringPointers(values ...string) []*string {
	ret := make([]*string, 0, len(values))
	for _, value := range values {
		ret = append(ret, new(value))
	}
	return ret
}
