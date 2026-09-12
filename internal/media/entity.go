package media

// entity.go —— AniList BaseAnime/BaseManga 及关联领域实体的镜像层（Wave A）。
//
// 镜像规则：
//   - 字段名与 JSON tag 与 anilist 包逐字段一致（硬约束：前端序列化形状不能破）；
//   - 仅两处改名：anilist.BaseAnime→media.Anime、anilist.BaseManga→media.Manga；
//     嵌套类型按同一规则去查询包装前缀：BaseAnime_Title→Anime_Title、
//     AnimeDetailsById_Media→AnimeDetails、MangaDetailsById_Media→MangaDetails；
//   - anilist 包中 Title/CoverImage/FuzzyDate/Trailer/NextAiringEpisode 等嵌套结构
//     在各查询响应里形状完全相同，此处只定义一份规范结构（Title/CoverImage/
//     FuzzyDate/Trailer/AiringEpisode），旧名一律用类型别名承接（如
//     Anime_Title = Title、AnimeSchedule_Previous_Nodes = AiringEpisode），
//     保证 Wave C 机械替换时旧类型名可用，同时避免数千行重复定义；
//   - Get* 方法语义照搬：nil 接收者返回零值（genqlient 风格）。
//
// Bangumi 扩展（additive，前端可选消费，不影响既有形状）：
//   - Title.Chinese（json:"nameCN"）：中文名，对应 Bangumi name_cn；
//   - Anime.NameCN / Manga.NameCN：顶层便捷字段；
//   - 评分换算：Bangumi rating 为 0-10，AniList meanScore 为 0-100，
//     映射时 MeanScore = round(rating × 10)；NSFW 复用 IsAdult 字段。
//
// 裁剪清单（相对 anilist 包，均为 AniList GraphQL 专属、换锚后无对应物）：
//   - FetchMediaTree 遍历函数（依赖 AnilistClient GraphQL 逐级查询；
//     Bangumi 用 GET /v0/subjects/{id}/subjects 一次取全，Wave B 重建）；
//   - ListAnimeM/ListMangaM/ListMissedSequels/ListRecentAiringAnimeM 及
//     GraphQL Document 常量（同上，query 助手不适用于 Bangumi）；
//   - 未使用的枚举（MediaListSort/MediaSource/MediaTrendSort 等，全仓 grep 零引用）；
//   - Ensure*CollectionEntry 系列（依赖 AnilistClient 拉取缺失条目，
//     media 包无客户端概念，Patch* 原语已保留）。
//   - 复合结构裁剪说明：AnimeDetails/MangaDetails 的推荐位节点
//     （MediaRecommendation）形状是 BaseAnime/BaseManga 的子集，
//     直接复用 *Anime/*Manga（多余字段 omitempty，JSON 形状不变）。

import (
	"fmt"
	"time"

	"github.com/samber/lo"
	"seanime/internal/util/comparison"
	"seanime/internal/util/result"
)

//----------------------------------------------------------------------------------------------------------------------
// 共享嵌套结构（anilist 包中各查询响应的同类嵌套只定义一份）

// Title 标题（anilist.BaseAnime_Title 等的合并镜像）。
// Chinese 为 Bangumi 扩展字段：承载 name_cn 中文名。
type Title struct {
	English       *string `json:"english,omitempty"`
	Native        *string `json:"native,omitempty"`
	Romaji        *string `json:"romaji,omitempty"`
	UserPreferred *string `json:"userPreferred,omitempty"`
	Chinese       *string `json:"nameCN,omitempty"`
}

// 旧名别名：Wave C 消费方按 anilist.Xxx → media.Xxx 机械替换
type (
	Anime_Title         = Title
	Manga_Title         = Title
	CompleteAnime_Title = Title
	AnimeDetails_Title  = Title
	MangaDetails_Title  = Title
)

func (t *Title) GetEnglish() *string {
	if t == nil {
		t = &Title{}
	}
	return t.English
}
func (t *Title) GetNative() *string {
	if t == nil {
		t = &Title{}
	}
	return t.Native
}
func (t *Title) GetRomaji() *string {
	if t == nil {
		t = &Title{}
	}
	return t.Romaji
}
func (t *Title) GetUserPreferred() *string {
	if t == nil {
		t = &Title{}
	}
	return t.UserPreferred
}
func (t *Title) GetChinese() *string {
	if t == nil {
		t = &Title{}
	}
	return t.Chinese
}

// CoverImage 封面图（anilist.BaseAnime_CoverImage 等的合并镜像）
type CoverImage struct {
	Color      *string `json:"color,omitempty"`
	ExtraLarge *string `json:"extraLarge,omitempty"`
	Large      *string `json:"large,omitempty"`
	Medium     *string `json:"medium,omitempty"`
}

type (
	Anime_CoverImage         = CoverImage
	Manga_CoverImage         = CoverImage
	CompleteAnime_CoverImage = CoverImage
)

func (t *CoverImage) GetColor() *string {
	if t == nil {
		t = &CoverImage{}
	}
	return t.Color
}
func (t *CoverImage) GetExtraLarge() *string {
	if t == nil {
		t = &CoverImage{}
	}
	return t.ExtraLarge
}
func (t *CoverImage) GetLarge() *string {
	if t == nil {
		t = &CoverImage{}
	}
	return t.Large
}
func (t *CoverImage) GetMedium() *string {
	if t == nil {
		t = &CoverImage{}
	}
	return t.Medium
}

// FuzzyDate 模糊日期（年/月/日可缺省），对应 anilist.FuzzyDate、
// BaseAnime_StartDate/EndDate、收藏条目的 StartedAt/CompletedAt 等同形结构。
type FuzzyDate struct {
	Year  *int `json:"year,omitempty"`
	Month *int `json:"month,omitempty"`
	Day   *int `json:"day,omitempty"`
}

type (
	Anime_StartDate         = FuzzyDate
	Anime_EndDate           = FuzzyDate
	Manga_StartDate         = FuzzyDate
	Manga_EndDate           = FuzzyDate
	CompleteAnime_StartDate = FuzzyDate
	CompleteAnime_EndDate   = FuzzyDate
	AnimeDetails_StartDate  = FuzzyDate
	AnimeDetails_EndDate    = FuzzyDate
	MangaDetails_StartDate  = FuzzyDate
	MangaDetails_EndDate    = FuzzyDate

	// FuzzyDateInput 变更入参（与 FuzzyDate 同形）
	FuzzyDateInput = FuzzyDate

	// 收藏条目日期旧名别名
	AnimeCollection_MediaListCollection_Lists_Entries_StartedAt   = FuzzyDate
	AnimeCollection_MediaListCollection_Lists_Entries_CompletedAt = FuzzyDate
	MangaCollection_MediaListCollection_Lists_Entries_StartedAt   = FuzzyDate
	MangaCollection_MediaListCollection_Lists_Entries_CompletedAt = FuzzyDate
)

func (t *FuzzyDate) GetYear() *int {
	if t == nil {
		t = &FuzzyDate{}
	}
	return t.Year
}
func (t *FuzzyDate) GetMonth() *int {
	if t == nil {
		t = &FuzzyDate{}
	}
	return t.Month
}
func (t *FuzzyDate) GetDay() *int {
	if t == nil {
		t = &FuzzyDate{}
	}
	return t.Day
}

// AiringEpisode 放送集（anilist.BaseAnime_NextAiringEpisode 及放送表节点同形结构）
type AiringEpisode struct {
	AiringAt        int `json:"airingAt"`
	Episode         int `json:"episode"`
	TimeUntilAiring int `json:"timeUntilAiring"`
}

type (
	Anime_NextAiringEpisode         = AiringEpisode
	CompleteAnime_NextAiringEpisode = AiringEpisode
	AnimeSchedule_Previous_Nodes    = AiringEpisode
	AnimeSchedule_Upcoming_Nodes    = AiringEpisode
)

func (t *AiringEpisode) GetAiringAt() int {
	if t == nil {
		t = &AiringEpisode{}
	}
	return t.AiringAt
}
func (t *AiringEpisode) GetEpisode() int {
	if t == nil {
		t = &AiringEpisode{}
	}
	return t.Episode
}
func (t *AiringEpisode) GetTimeUntilAiring() int {
	if t == nil {
		t = &AiringEpisode{}
	}
	return t.TimeUntilAiring
}

// Trailer 预告片（anilist.BaseAnime_Trailer 等同形结构）
type Trailer struct {
	ID        *string `json:"id,omitempty"`
	Site      *string `json:"site,omitempty"`
	Thumbnail *string `json:"thumbnail,omitempty"`
}

type (
	Anime_Trailer         = Trailer
	Manga_Trailer         = Trailer
	CompleteAnime_Trailer = Trailer
	AnimeDetails_Trailer  = Trailer
	MangaDetails_Trailer  = Trailer
)

func (t *Trailer) GetID() *string {
	if t == nil {
		t = &Trailer{}
	}
	return t.ID
}
func (t *Trailer) GetSite() *string {
	if t == nil {
		t = &Trailer{}
	}
	return t.Site
}
func (t *Trailer) GetThumbnail() *string {
	if t == nil {
		t = &Trailer{}
	}
	return t.Thumbnail
}

//----------------------------------------------------------------------------------------------------------------------
// Anime（≈ anilist.BaseAnime）

// Anime 镜像 anilist.BaseAnime，字段与 JSON tag 逐字段一致。
// ID 语义 = Bangumi subject ID（换锚后主锚点）。
type Anime struct {
	ID                int             `json:"id"`
	IDMal             *int            `json:"idMal,omitempty"`
	SiteURL           *string         `json:"siteUrl,omitempty"`
	Status            *MediaStatus    `json:"status,omitempty"`
	Season            *MediaSeason    `json:"season,omitempty"`
	Type              *MediaType      `json:"type,omitempty"`
	Format            *MediaFormat    `json:"format,omitempty"`
	SeasonYear        *int            `json:"seasonYear,omitempty"`
	BannerImage       *string         `json:"bannerImage,omitempty"`
	Episodes          *int            `json:"episodes,omitempty"`
	Synonyms          []*string       `json:"synonyms,omitempty"`
	IsAdult           *bool           `json:"isAdult,omitempty"`
	CountryOfOrigin   *string         `json:"countryOfOrigin,omitempty"`
	MeanScore         *int            `json:"meanScore,omitempty"`
	Description       *string         `json:"description,omitempty"`
	Genres            []*string       `json:"genres,omitempty"`
	Duration          *int            `json:"duration,omitempty"`
	Trailer           *Anime_Trailer  `json:"trailer,omitempty"`
	Title             *Anime_Title    `json:"title,omitempty"`
	CoverImage        *Anime_CoverImage `json:"coverImage,omitempty"`
	StartDate         *Anime_StartDate `json:"startDate,omitempty"`
	EndDate           *Anime_EndDate  `json:"endDate,omitempty"`
	NextAiringEpisode *Anime_NextAiringEpisode `json:"nextAiringEpisode,omitempty"`

	// NameCN —— Bangumi 扩展便捷字段（中文名，对应 subject.name_cn）
	NameCN string `json:"nameCN,omitempty"`
}

func (t *Anime) GetID() int {
	if t == nil {
		t = &Anime{}
	}
	return t.ID
}
func (t *Anime) GetIDMal() *int {
	if t == nil {
		t = &Anime{}
	}
	return t.IDMal
}
func (t *Anime) GetSiteURL() *string {
	if t == nil {
		t = &Anime{}
	}
	return t.SiteURL
}
func (t *Anime) GetStatus() *MediaStatus {
	if t == nil {
		t = &Anime{}
	}
	return t.Status
}
func (t *Anime) GetSeason() *MediaSeason {
	if t == nil {
		t = &Anime{}
	}
	return t.Season
}
func (t *Anime) GetType() *MediaType {
	if t == nil {
		t = &Anime{}
	}
	return t.Type
}
func (t *Anime) GetFormat() *MediaFormat {
	if t == nil {
		t = &Anime{}
	}
	return t.Format
}
func (t *Anime) GetSeasonYear() *int {
	if t == nil {
		t = &Anime{}
	}
	return t.SeasonYear
}
func (t *Anime) GetBannerImage() *string {
	if t == nil {
		t = &Anime{}
	}
	return t.BannerImage
}
func (t *Anime) GetEpisodes() *int {
	if t == nil {
		t = &Anime{}
	}
	return t.Episodes
}
func (t *Anime) GetSynonyms() []*string {
	if t == nil {
		t = &Anime{}
	}
	return t.Synonyms
}
func (t *Anime) GetIsAdult() *bool {
	if t == nil {
		t = &Anime{}
	}
	return t.IsAdult
}
func (t *Anime) GetCountryOfOrigin() *string {
	if t == nil {
		t = &Anime{}
	}
	return t.CountryOfOrigin
}
func (t *Anime) GetMeanScore() *int {
	if t == nil {
		t = &Anime{}
	}
	return t.MeanScore
}
func (t *Anime) GetDescription() *string {
	if t == nil {
		t = &Anime{}
	}
	return t.Description
}
func (t *Anime) GetGenres() []*string {
	if t == nil {
		t = &Anime{}
	}
	return t.Genres
}
func (t *Anime) GetDuration() *int {
	if t == nil {
		t = &Anime{}
	}
	return t.Duration
}
func (t *Anime) GetTrailer() *Anime_Trailer {
	if t == nil {
		t = &Anime{}
	}
	return t.Trailer
}
func (t *Anime) GetTitle() *Anime_Title {
	if t == nil {
		t = &Anime{}
	}
	return t.Title
}
func (t *Anime) GetCoverImage() *Anime_CoverImage {
	if t == nil {
		t = &Anime{}
	}
	return t.CoverImage
}
func (t *Anime) GetStartDate() *Anime_StartDate {
	if t == nil {
		t = &Anime{}
	}
	return t.StartDate
}
func (t *Anime) GetEndDate() *Anime_EndDate {
	if t == nil {
		t = &Anime{}
	}
	return t.EndDate
}
func (t *Anime) GetNextAiringEpisode() *Anime_NextAiringEpisode {
	if t == nil {
		t = &Anime{}
	}
	return t.NextAiringEpisode
}

// ---------------------------------------------------------------------------------------------------------------------
// Anime 安全辅助（照搬 anilist/media_helper.go 语义）

func (m *Anime) GetTitleSafe() string {
	if m.GetTitle().GetEnglish() != nil {
		return *m.GetTitle().GetEnglish()
	}
	if m.GetTitle().GetRomaji() != nil {
		return *m.GetTitle().GetRomaji()
	}
	// Bangumi 扩展：无英/罗马音时回退中文名，再回退空串
	if m.GetTitle().GetChinese() != nil {
		return *m.GetTitle().GetChinese()
	}
	return ""
}

func (m *Anime) GetEnglishTitleSafe() string {
	if m.GetTitle().GetEnglish() != nil {
		return *m.GetTitle().GetEnglish()
	}
	return ""
}

func (m *Anime) GetRomajiTitleSafe() string {
	if m.GetTitle().GetRomaji() != nil {
		return *m.GetTitle().GetRomaji()
	}
	if m.GetTitle().GetEnglish() != nil {
		return *m.GetTitle().GetEnglish()
	}
	return ""
}

func (m *Anime) GetPreferredTitle() string {
	if m.GetTitle().GetUserPreferred() != nil {
		return *m.GetTitle().GetUserPreferred()
	}
	return m.GetTitleSafe()
}

func (m *Anime) GetCoverImageSafe() string {
	if m.GetCoverImage().GetExtraLarge() != nil {
		return *m.GetCoverImage().GetExtraLarge()
	}
	if m.GetCoverImage().GetLarge() != nil {
		return *m.GetCoverImage().GetLarge()
	}
	if m.GetBannerImage() != nil {
		return *m.GetBannerImage()
	}
	return ""
}

func (m *Anime) GetBannerImageSafe() string {
	if m.GetBannerImage() != nil {
		return *m.GetBannerImage()
	}
	return m.GetCoverImageSafe()
}

func (m *Anime) IsMovieOrSingleEpisode() bool {
	if m == nil {
		return false
	}
	if m.GetTotalEpisodeCount() == 1 {
		return true
	}
	return false
}

func (m *Anime) GetSynonymsDeref() []string {
	if m.Synonyms == nil {
		return nil
	}
	return lo.Map(m.Synonyms, func(s *string, i int) string { return *s })
}

func (m *Anime) GetSynonymsContainingSeason() []string {
	if m.Synonyms == nil {
		return nil
	}
	return lo.Filter(lo.Map(m.Synonyms, func(s *string, i int) string { return *s }), func(s string, i int) bool { return comparison.ValueContainsSeason(s) })
}

func (m *Anime) GetStartYearSafe() int {
	if m == nil || m.StartDate == nil || m.StartDate.Year == nil {
		return 0
	}
	return *m.StartDate.Year
}

func (m *Anime) IsMovie() bool {
	if m == nil {
		return false
	}
	if m.Format == nil {
		return false
	}
	return *m.Format == MediaFormatMovie
}

func (m *Anime) IsFinished() bool {
	if m == nil {
		return false
	}
	if m.Status == nil {
		return false
	}
	return *m.Status == MediaStatusFinished
}

func (m *Anime) GetAllTitles() []*string {
	titles := make([]*string, 0)
	if m.HasRomajiTitle() {
		titles = append(titles, m.Title.Romaji)
	}
	if m.HasEnglishTitle() {
		titles = append(titles, m.Title.English)
	}
	if m.HasChineseTitle() {
		titles = append(titles, m.Title.Chinese)
	}
	if m.HasSynonyms() && len(m.Synonyms) > 1 {
		titles = append(titles, lo.Filter(m.Synonyms, func(s *string, i int) bool { return comparison.ValueContainsSeason(*s) })...)
	}
	return titles
}

func (m *Anime) GetAllTitlesDeref() []string {
	titles := make([]string, 0)
	if m.HasRomajiTitle() {
		titles = append(titles, *m.Title.Romaji)
	}
	if m.HasEnglishTitle() {
		titles = append(titles, *m.Title.English)
	}
	if m.HasChineseTitle() {
		titles = append(titles, *m.Title.Chinese)
	}
	if m.HasSynonyms() && len(m.Synonyms) > 1 {
		syn := lo.Filter(m.Synonyms, func(s *string, i int) bool { return comparison.ValueContainsSeason(*s) })
		for _, s := range syn {
			titles = append(titles, *s)
		}
	}
	return titles
}

func (m *Anime) GetMainTitles() []*string {
	titles := make([]*string, 0)
	if m.HasRomajiTitle() {
		titles = append(titles, m.Title.Romaji)
	}
	if m.HasEnglishTitle() {
		titles = append(titles, m.Title.English)
	}
	return titles
}

func (m *Anime) GetMainTitlesDeref() []string {
	titles := make([]string, 0)
	if m.HasRomajiTitle() {
		titles = append(titles, *m.Title.Romaji)
	}
	if m.HasEnglishTitle() {
		titles = append(titles, *m.Title.English)
	}
	return titles
}

// GetCurrentEpisodeCount 返回当前已放送集数；无集数且无下一集信息时返回 -1。
func (m *Anime) GetCurrentEpisodeCount() int {
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

func (m *Anime) GetCurrentEpisodeCountOrNil() *int {
	n := m.GetCurrentEpisodeCount()
	if n == -1 {
		return nil
	}
	return &n
}

// GetTotalEpisodeCount 返回总集数，未知时返回 -1。
func (m *Anime) GetTotalEpisodeCount() int {
	ceil := -1
	if m.Episodes != nil {
		ceil = *m.Episodes
	}
	return ceil
}

func (m *Anime) GetTotalEpisodeCountOrNil() *int {
	return m.Episodes
}

// GetPossibleSeasonNumber 从别名/标题中推断季度号，未知返回 -1。
func (m *Anime) GetPossibleSeasonNumber() int {
	if m == nil || m.Synonyms == nil || len(m.Synonyms) == 0 {
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

func (m *Anime) HasEnglishTitle() bool {
	return m.Title.English != nil
}

func (m *Anime) HasRomajiTitle() bool {
	return m.Title.Romaji != nil
}

// HasChineseTitle —— Bangumi 扩展判断
func (m *Anime) HasChineseTitle() bool {
	return m.Title.Chinese != nil
}

func (m *Anime) HasSynonyms() bool {
	return m.Synonyms != nil
}

//----------------------------------------------------------------------------------------------------------------------
// Manga（≈ anilist.BaseManga）

// Manga 镜像 anilist.BaseManga。ID 语义 = Bangumi subject ID。
type Manga struct {
	ID              int                `json:"id"`
	IDMal           *int               `json:"idMal,omitempty"`
	SiteURL         *string            `json:"siteUrl,omitempty"`
	Status          *MediaStatus       `json:"status,omitempty"`
	Season          *MediaSeason       `json:"season,omitempty"`
	Type            *MediaType         `json:"type,omitempty"`
	Format          *MediaFormat       `json:"format,omitempty"`
	BannerImage     *string            `json:"bannerImage,omitempty"`
	Chapters        *int               `json:"chapters,omitempty"`
	Volumes         *int               `json:"volumes,omitempty"`
	Synonyms        []*string          `json:"synonyms,omitempty"`
	IsAdult         *bool              `json:"isAdult,omitempty"`
	CountryOfOrigin *string            `json:"countryOfOrigin,omitempty"`
	MeanScore       *int               `json:"meanScore,omitempty"`
	Description     *string            `json:"description,omitempty"`
	Genres          []*string          `json:"genres,omitempty"`
	Title           *Manga_Title       `json:"title,omitempty"`
	CoverImage      *Manga_CoverImage  `json:"coverImage,omitempty"`
	StartDate       *Manga_StartDate   `json:"startDate,omitempty"`
	EndDate         *Manga_EndDate     `json:"endDate,omitempty"`

	// NameCN —— Bangumi 扩展便捷字段
	NameCN string `json:"nameCN,omitempty"`
}

func (t *Manga) GetID() int {
	if t == nil {
		t = &Manga{}
	}
	return t.ID
}
func (t *Manga) GetIDMal() *int {
	if t == nil {
		t = &Manga{}
	}
	return t.IDMal
}
func (t *Manga) GetSiteURL() *string {
	if t == nil {
		t = &Manga{}
	}
	return t.SiteURL
}
func (t *Manga) GetStatus() *MediaStatus {
	if t == nil {
		t = &Manga{}
	}
	return t.Status
}
func (t *Manga) GetSeason() *MediaSeason {
	if t == nil {
		t = &Manga{}
	}
	return t.Season
}
func (t *Manga) GetType() *MediaType {
	if t == nil {
		t = &Manga{}
	}
	return t.Type
}
func (t *Manga) GetFormat() *MediaFormat {
	if t == nil {
		t = &Manga{}
	}
	return t.Format
}
func (t *Manga) GetBannerImage() *string {
	if t == nil {
		t = &Manga{}
	}
	return t.BannerImage
}
func (t *Manga) GetChapters() *int {
	if t == nil {
		t = &Manga{}
	}
	return t.Chapters
}
func (t *Manga) GetVolumes() *int {
	if t == nil {
		t = &Manga{}
	}
	return t.Volumes
}
func (t *Manga) GetSynonyms() []*string {
	if t == nil {
		t = &Manga{}
	}
	return t.Synonyms
}
func (t *Manga) GetIsAdult() *bool {
	if t == nil {
		t = &Manga{}
	}
	return t.IsAdult
}
func (t *Manga) GetCountryOfOrigin() *string {
	if t == nil {
		t = &Manga{}
	}
	return t.CountryOfOrigin
}
func (t *Manga) GetMeanScore() *int {
	if t == nil {
		t = &Manga{}
	}
	return t.MeanScore
}
func (t *Manga) GetDescription() *string {
	if t == nil {
		t = &Manga{}
	}
	return t.Description
}
func (t *Manga) GetGenres() []*string {
	if t == nil {
		t = &Manga{}
	}
	return t.Genres
}
func (t *Manga) GetTitle() *Manga_Title {
	if t == nil {
		t = &Manga{}
	}
	return t.Title
}
func (t *Manga) GetCoverImage() *Manga_CoverImage {
	if t == nil {
		t = &Manga{}
	}
	return t.CoverImage
}
func (t *Manga) GetStartDate() *Manga_StartDate {
	if t == nil {
		t = &Manga{}
	}
	return t.StartDate
}
func (t *Manga) GetEndDate() *Manga_EndDate {
	if t == nil {
		t = &Manga{}
	}
	return t.EndDate
}

func (m *Manga) GetTitleSafe() string {
	if m.GetTitle().GetEnglish() != nil {
		return *m.GetTitle().GetEnglish()
	}
	if m.GetTitle().GetRomaji() != nil {
		return *m.GetTitle().GetRomaji()
	}
	if m.GetTitle().GetChinese() != nil {
		return *m.GetTitle().GetChinese()
	}
	return ""
}
func (m *Manga) GetRomajiTitleSafe() string {
	if m.GetTitle().GetRomaji() != nil {
		return *m.GetTitle().GetRomaji()
	}
	if m.GetTitle().GetEnglish() != nil {
		return *m.GetTitle().GetEnglish()
	}
	return ""
}

func (m *Manga) GetPreferredTitle() string {
	if m.GetTitle().GetUserPreferred() != nil {
		return *m.GetTitle().GetUserPreferred()
	}
	return m.GetTitleSafe()
}

func (m *Manga) GetCoverImageSafe() string {
	if m.GetCoverImage().GetExtraLarge() != nil {
		return *m.GetCoverImage().GetExtraLarge()
	}
	if m.GetCoverImage().GetLarge() != nil {
		return *m.GetCoverImage().GetLarge()
	}
	if m.GetBannerImage() != nil {
		return *m.GetBannerImage()
	}
	return ""
}
func (m *Manga) GetBannerImageSafe() string {
	if m.GetBannerImage() != nil {
		return *m.GetBannerImage()
	}
	return m.GetCoverImageSafe()
}

func (m *Manga) GetAllTitles() []*string {
	titles := make([]*string, 0)
	if m.HasRomajiTitle() {
		titles = append(titles, m.Title.Romaji)
	}
	if m.HasEnglishTitle() {
		titles = append(titles, m.Title.English)
	}
	if m.HasChineseTitle() {
		titles = append(titles, m.Title.Chinese)
	}
	if m.HasSynonyms() && len(m.Synonyms) > 1 {
		titles = append(titles, m.Synonyms...)
	}
	return titles
}

func (m *Manga) GetMainTitlesDeref() []string {
	titles := make([]string, 0)
	if m.HasRomajiTitle() {
		titles = append(titles, *m.Title.Romaji)
	}
	if m.HasEnglishTitle() {
		titles = append(titles, *m.Title.English)
	}
	return titles
}

func (m *Manga) HasEnglishTitle() bool {
	return m.Title.English != nil
}
func (m *Manga) HasRomajiTitle() bool {
	return m.Title.Romaji != nil
}
func (m *Manga) HasChineseTitle() bool {
	return m.Title.Chinese != nil
}
func (m *Manga) HasSynonyms() bool {
	return m.Synonyms != nil
}

func (m *Manga) GetStartYearSafe() int {
	if m.GetStartDate() != nil && m.GetStartDate().GetYear() != nil {
		return *m.GetStartDate().GetYear()
	}
	return 0
}

//----------------------------------------------------------------------------------------------------------------------
// CompleteAnime（镜像 anilist.CompleteAnime，含 relations，scanner 建关系树用）

type CompleteAnime struct {
	ID                int                        `json:"id"`
	IDMal             *int                       `json:"idMal,omitempty"`
	SiteURL           *string                    `json:"siteUrl,omitempty"`
	Status            *MediaStatus               `json:"status,omitempty"`
	Season            *MediaSeason               `json:"season,omitempty"`
	SeasonYear        *int                       `json:"seasonYear,omitempty"`
	Type              *MediaType                 `json:"type,omitempty"`
	Format            *MediaFormat               `json:"format,omitempty"`
	BannerImage       *string                    `json:"bannerImage,omitempty"`
	Episodes          *int                       `json:"episodes,omitempty"`
	Synonyms          []*string                  `json:"synonyms,omitempty"`
	IsAdult           *bool                      `json:"isAdult,omitempty"`
	CountryOfOrigin   *string                    `json:"countryOfOrigin,omitempty"`
	MeanScore         *int                       `json:"meanScore,omitempty"`
	Description       *string                    `json:"description,omitempty"`
	Genres            []*string                  `json:"genres,omitempty"`
	Duration          *int                       `json:"duration,omitempty"`
	Trailer           *CompleteAnime_Trailer     `json:"trailer,omitempty"`
	Title             *CompleteAnime_Title       `json:"title,omitempty"`
	CoverImage        *CompleteAnime_CoverImage  `json:"coverImage,omitempty"`
	StartDate         *CompleteAnime_StartDate   `json:"startDate,omitempty"`
	EndDate           *CompleteAnime_EndDate     `json:"endDate,omitempty"`
	NextAiringEpisode *CompleteAnime_NextAiringEpisode `json:"nextAiringEpisode,omitempty"`
	Relations         *CompleteAnime_Relations   `json:"relations,omitempty"`

	// NameCN —— Bangumi 扩展便捷字段
	NameCN string `json:"nameCN,omitempty"`
}

func (t *CompleteAnime) GetID() int {
	if t == nil {
		t = &CompleteAnime{}
	}
	return t.ID
}
func (t *CompleteAnime) GetIDMal() *int {
	if t == nil {
		t = &CompleteAnime{}
	}
	return t.IDMal
}
func (t *CompleteAnime) GetSiteURL() *string {
	if t == nil {
		t = &CompleteAnime{}
	}
	return t.SiteURL
}
func (t *CompleteAnime) GetStatus() *MediaStatus {
	if t == nil {
		t = &CompleteAnime{}
	}
	return t.Status
}
func (t *CompleteAnime) GetSeason() *MediaSeason {
	if t == nil {
		t = &CompleteAnime{}
	}
	return t.Season
}
func (t *CompleteAnime) GetSeasonYear() *int {
	if t == nil {
		t = &CompleteAnime{}
	}
	return t.SeasonYear
}
func (t *CompleteAnime) GetType() *MediaType {
	if t == nil {
		t = &CompleteAnime{}
	}
	return t.Type
}
func (t *CompleteAnime) GetFormat() *MediaFormat {
	if t == nil {
		t = &CompleteAnime{}
	}
	return t.Format
}
func (t *CompleteAnime) GetBannerImage() *string {
	if t == nil {
		t = &CompleteAnime{}
	}
	return t.BannerImage
}
func (t *CompleteAnime) GetEpisodes() *int {
	if t == nil {
		t = &CompleteAnime{}
	}
	return t.Episodes
}
func (t *CompleteAnime) GetSynonyms() []*string {
	if t == nil {
		t = &CompleteAnime{}
	}
	return t.Synonyms
}
func (t *CompleteAnime) GetIsAdult() *bool {
	if t == nil {
		t = &CompleteAnime{}
	}
	return t.IsAdult
}
func (t *CompleteAnime) GetCountryOfOrigin() *string {
	if t == nil {
		t = &CompleteAnime{}
	}
	return t.CountryOfOrigin
}
func (t *CompleteAnime) GetMeanScore() *int {
	if t == nil {
		t = &CompleteAnime{}
	}
	return t.MeanScore
}
func (t *CompleteAnime) GetDescription() *string {
	if t == nil {
		t = &CompleteAnime{}
	}
	return t.Description
}
func (t *CompleteAnime) GetGenres() []*string {
	if t == nil {
		t = &CompleteAnime{}
	}
	return t.Genres
}
func (t *CompleteAnime) GetDuration() *int {
	if t == nil {
		t = &CompleteAnime{}
	}
	return t.Duration
}
func (t *CompleteAnime) GetTrailer() *CompleteAnime_Trailer {
	if t == nil {
		t = &CompleteAnime{}
	}
	return t.Trailer
}
func (t *CompleteAnime) GetTitle() *CompleteAnime_Title {
	if t == nil {
		t = &CompleteAnime{}
	}
	return t.Title
}
func (t *CompleteAnime) GetCoverImage() *CompleteAnime_CoverImage {
	if t == nil {
		t = &CompleteAnime{}
	}
	return t.CoverImage
}
func (t *CompleteAnime) GetStartDate() *CompleteAnime_StartDate {
	if t == nil {
		t = &CompleteAnime{}
	}
	return t.StartDate
}
func (t *CompleteAnime) GetEndDate() *CompleteAnime_EndDate {
	if t == nil {
		t = &CompleteAnime{}
	}
	return t.EndDate
}
func (t *CompleteAnime) GetNextAiringEpisode() *CompleteAnime_NextAiringEpisode {
	if t == nil {
		t = &CompleteAnime{}
	}
	return t.NextAiringEpisode
}
func (t *CompleteAnime) GetRelations() *CompleteAnime_Relations {
	if t == nil {
		t = &CompleteAnime{}
	}
	return t.Relations
}

// CompleteAnime_Relations 关系边集（节点为 BaseAnime 镜像 Anime）
type CompleteAnime_Relations struct {
	Edges []*CompleteAnime_Relations_Edges `json:"edges,omitempty"`
}

type CompleteAnime_Relations_Edges struct {
	Node         *Anime         `json:"node,omitempty"`
	RelationType *MediaRelation `json:"relationType,omitempty"`
}

func (t *CompleteAnime_Relations) GetEdges() []*CompleteAnime_Relations_Edges {
	if t == nil {
		t = &CompleteAnime_Relations{}
	}
	return t.Edges
}
func (t *CompleteAnime_Relations_Edges) GetNode() *Anime {
	if t == nil {
		t = &CompleteAnime_Relations_Edges{}
	}
	return t.Node
}
func (t *CompleteAnime_Relations_Edges) GetRelationType() *MediaRelation {
	if t == nil {
		t = &CompleteAnime_Relations_Edges{}
	}
	return t.RelationType
}

// ---------------------------------------------------------------------------------------------------------------------
// CompleteAnime 安全辅助（照搬语义）

func (m *CompleteAnime) GetTitleSafe() string {
	if m.GetTitle().GetEnglish() != nil {
		return *m.GetTitle().GetEnglish()
	}
	if m.GetTitle().GetRomaji() != nil {
		return *m.GetTitle().GetRomaji()
	}
	if m.GetTitle().GetChinese() != nil {
		return *m.GetTitle().GetChinese()
	}
	return "N/A"
}
func (m *CompleteAnime) GetRomajiTitleSafe() string {
	if m.GetTitle().GetRomaji() != nil {
		return *m.GetTitle().GetRomaji()
	}
	if m.GetTitle().GetEnglish() != nil {
		return *m.GetTitle().GetEnglish()
	}
	return "N/A"
}

func (m *CompleteAnime) GetPreferredTitle() string {
	if m.GetTitle().GetUserPreferred() != nil {
		return *m.GetTitle().GetUserPreferred()
	}
	return m.GetTitleSafe()
}

func (m *CompleteAnime) GetCoverImageSafe() string {
	if m.GetCoverImage().GetExtraLarge() != nil {
		return *m.GetCoverImage().GetExtraLarge()
	}
	if m.GetCoverImage().GetLarge() != nil {
		return *m.GetCoverImage().GetLarge()
	}
	if m.GetBannerImage() != nil {
		return *m.GetBannerImage()
	}
	return ""
}

func (m *CompleteAnime) GetBannerImageSafe() string {
	if m.GetBannerImage() != nil {
		return *m.GetBannerImage()
	}
	return m.GetCoverImageSafe()
}

func (m *CompleteAnime) IsMovieOrSingleEpisode() bool {
	if m == nil {
		return false
	}
	if m.GetTotalEpisodeCount() == 1 {
		return true
	}
	return false
}

func (m *CompleteAnime) IsMovie() bool {
	if m == nil {
		return false
	}
	if m.Format == nil {
		return false
	}
	return *m.Format == MediaFormatMovie
}

func (m *CompleteAnime) IsFinished() bool {
	if m == nil {
		return false
	}
	if m.Status == nil {
		return false
	}
	return *m.Status == MediaStatusFinished
}

func (m *CompleteAnime) GetAllTitles() []*string {
	titles := make([]*string, 0)
	if m.HasRomajiTitle() {
		titles = append(titles, m.Title.Romaji)
	}
	if m.HasEnglishTitle() {
		titles = append(titles, m.Title.English)
	}
	if m.HasChineseTitle() {
		titles = append(titles, m.Title.Chinese)
	}
	if m.HasSynonyms() && len(m.Synonyms) > 1 {
		titles = append(titles, lo.Filter(m.Synonyms, func(s *string, i int) bool { return comparison.ValueContainsSeason(*s) })...)
	}
	return titles
}

func (m *CompleteAnime) GetAllTitlesDeref() []string {
	titles := make([]string, 0)
	if m.HasRomajiTitle() {
		titles = append(titles, *m.Title.Romaji)
	}
	if m.HasEnglishTitle() {
		titles = append(titles, *m.Title.English)
	}
	if m.HasChineseTitle() {
		titles = append(titles, *m.Title.Chinese)
	}
	if m.HasSynonyms() && len(m.Synonyms) > 1 {
		syn := lo.Filter(m.Synonyms, func(s *string, i int) bool { return comparison.ValueContainsSeason(*s) })
		for _, s := range syn {
			titles = append(titles, *s)
		}
	}
	return titles
}

// GetCurrentEpisodeCount 语义同 Anime.GetCurrentEpisodeCount。
func (m *CompleteAnime) GetCurrentEpisodeCount() int {
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

// GetTotalEpisodeCount 语义同 Anime.GetTotalEpisodeCount。
func (m *CompleteAnime) GetTotalEpisodeCount() int {
	ceil := -1
	if m.Episodes != nil {
		ceil = *m.Episodes
	}
	return ceil
}

func (m *CompleteAnime) GetPossibleSeasonNumber() int {
	if m == nil || m.Synonyms == nil || len(m.Synonyms) == 0 {
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

func (m *CompleteAnime) HasEnglishTitle() bool {
	return m.Title.English != nil
}

func (m *CompleteAnime) HasRomajiTitle() bool {
	return m.Title.Romaji != nil
}

func (m *CompleteAnime) HasChineseTitle() bool {
	return m.Title.Chinese != nil
}

func (m *CompleteAnime) HasSynonyms() bool {
	return m.Synonyms != nil
}

// EdgeNarrowFormats / EdgeBroaderFormats 关系边匹配用的形式集合
var (
	EdgeNarrowFormats  = []MediaFormat{MediaFormatTv, MediaFormatTvShort}
	EdgeBroaderFormats = []MediaFormat{MediaFormatTv, MediaFormatTvShort, MediaFormatOna, MediaFormatOva, MediaFormatMovie, MediaFormatSpecial}
)

func (m *CompleteAnime) FindEdge(relation string, formats []MediaFormat) (*Anime, bool) {
	if m.GetRelations() == nil {
		return nil, false
	}
	for _, edge := range m.GetRelations().GetEdges() {
		if edge.GetRelationType().String() == relation {
			for _, fm := range formats {
				if fm.String() == edge.GetNode().GetFormat().String() {
					return edge.GetNode(), true
				}
			}
		}
	}
	return nil, false
}

func (e *CompleteAnime_Relations_Edges) IsBroadRelationFormat() bool {
	return edgeHasFormat(e, EdgeBroaderFormats)
}
func (e *CompleteAnime_Relations_Edges) IsNarrowRelationFormat() bool {
	return edgeHasFormat(e, EdgeNarrowFormats)
}

func edgeHasFormat(e *CompleteAnime_Relations_Edges, formats []MediaFormat) bool {
	if e.GetNode() == nil {
		return false
	}
	if e.GetNode().GetFormat() == nil {
		return false
	}
	for _, fm := range formats {
		if fm.String() == e.GetNode().GetFormat().String() {
			return true
		}
	}
	return false
}

// ToAnime（原 anilist.CompleteAnime.ToBaseAnime）：共享嵌套结构后为纯字段拷贝。
func (m *CompleteAnime) ToAnime() *Anime {
	if m == nil {
		return nil
	}
	return &Anime{
		ID:                m.ID,
		IDMal:             m.IDMal,
		SiteURL:           m.SiteURL,
		Status:            m.Status,
		Season:            m.Season,
		Type:              m.Type,
		Format:            m.Format,
		SeasonYear:        m.SeasonYear,
		BannerImage:       m.BannerImage,
		Episodes:          m.Episodes,
		Synonyms:          m.Synonyms,
		IsAdult:           m.IsAdult,
		CountryOfOrigin:   m.CountryOfOrigin,
		MeanScore:         m.MeanScore,
		Description:       m.Description,
		Genres:            m.Genres,
		Duration:          m.Duration,
		Trailer:           m.Trailer,
		Title:             m.Title,
		CoverImage:        m.CoverImage,
		StartDate:         m.StartDate,
		EndDate:           m.EndDate,
		NextAiringEpisode: m.NextAiringEpisode,
		NameCN:            m.NameCN,
	}
}

// ToCompleteAnime（原 anilist.BaseAnime.ToCompleteAnime）：无关系边的 CompleteAnime。
func (m *Anime) ToCompleteAnime() *CompleteAnime {
	if m == nil {
		return nil
	}
	return &CompleteAnime{
		ID:                m.ID,
		IDMal:             m.IDMal,
		SiteURL:           m.SiteURL,
		Format:            m.Format,
		Episodes:          m.Episodes,
		Status:            m.Status,
		Synonyms:          m.Synonyms,
		BannerImage:       m.BannerImage,
		Season:            m.Season,
		SeasonYear:        m.SeasonYear,
		Type:              m.Type,
		IsAdult:           m.IsAdult,
		CountryOfOrigin:   m.CountryOfOrigin,
		Genres:            m.Genres,
		Duration:          m.Duration,
		Description:       m.Description,
		MeanScore:         m.MeanScore,
		Trailer:           m.Trailer,
		Title:             m.Title,
		CoverImage:        m.CoverImage,
		StartDate:         m.StartDate,
		EndDate:           m.EndDate,
		NextAiringEpisode: m.NextAiringEpisode,
		Relations: &CompleteAnime_Relations{
			Edges: make([]*CompleteAnime_Relations_Edges, 0),
		},
		NameCN: m.NameCN,
	}
}

//----------------------------------------------------------------------------------------------------------------------
// BaseCharacter（详情页角色，镜像 anilist.BaseCharacter）

type BaseCharacter struct {
	ID          int                       `json:"id"`
	IsFavourite bool                      `json:"isFavourite"`
	Gender      *string                   `json:"gender,omitempty"`
	Age         *string                   `json:"age,omitempty"`
	DateOfBirth *BaseCharacter_DateOfBirth `json:"dateOfBirth,omitempty"`
	Name        *BaseCharacter_Name       `json:"name,omitempty"`
	Image       *BaseCharacter_Image      `json:"image,omitempty"`
	Description *string                   `json:"description,omitempty"`
	SiteURL     *string                   `json:"siteUrl,omitempty"`
}

type (
	BaseCharacter_DateOfBirth = FuzzyDate
)

type BaseCharacter_Name struct {
	Alternative []*string `json:"alternative,omitempty"`
	Full        *string   `json:"full,omitempty"`
	Native      *string   `json:"native,omitempty"`
}

type BaseCharacter_Image struct {
	Large *string `json:"large,omitempty"`
}

func (t *BaseCharacter) GetID() int {
	if t == nil {
		t = &BaseCharacter{}
	}
	return t.ID
}
func (t *BaseCharacter) GetIsFavourite() bool {
	if t == nil {
		t = &BaseCharacter{}
	}
	return t.IsFavourite
}
func (t *BaseCharacter) GetGender() *string {
	if t == nil {
		t = &BaseCharacter{}
	}
	return t.Gender
}
func (t *BaseCharacter) GetAge() *string {
	if t == nil {
		t = &BaseCharacter{}
	}
	return t.Age
}
func (t *BaseCharacter) GetDateOfBirth() *BaseCharacter_DateOfBirth {
	if t == nil {
		t = &BaseCharacter{}
	}
	return t.DateOfBirth
}
func (t *BaseCharacter) GetName() *BaseCharacter_Name {
	if t == nil {
		t = &BaseCharacter{}
	}
	return t.Name
}
func (t *BaseCharacter) GetImage() *BaseCharacter_Image {
	if t == nil {
		t = &BaseCharacter{}
	}
	return t.Image
}
func (t *BaseCharacter) GetDescription() *string {
	if t == nil {
		t = &BaseCharacter{}
	}
	return t.Description
}
func (t *BaseCharacter) GetSiteURL() *string {
	if t == nil {
		t = &BaseCharacter{}
	}
	return t.SiteURL
}

func (t *BaseCharacter_Name) GetAlternative() []*string {
	if t == nil {
		t = &BaseCharacter_Name{}
	}
	return t.Alternative
}
func (t *BaseCharacter_Name) GetFull() *string {
	if t == nil {
		t = &BaseCharacter_Name{}
	}
	return t.Full
}
func (t *BaseCharacter_Name) GetNative() *string {
	if t == nil {
		t = &BaseCharacter_Name{}
	}
	return t.Native
}

func (t *BaseCharacter_Image) GetLarge() *string {
	if t == nil {
		t = &BaseCharacter_Image{}
	}
	return t.Large
}

//----------------------------------------------------------------------------------------------------------------------
// AnimeDetails（≈ anilist.AnimeDetailsById_Media）与 MangaDetails（≈ MangaDetailsById_Media）

// AnimeDetails 条目详情（含角色/Staff/制作组/排名/推荐/关系）
type AnimeDetails struct {
	AverageScore    *int                             `json:"averageScore,omitempty"`
	Characters      *AnimeDetails_Characters         `json:"characters,omitempty"`
	Description     *string                          `json:"description,omitempty"`
	Duration        *int                             `json:"duration,omitempty"`
	EndDate         *AnimeDetails_EndDate            `json:"endDate,omitempty"`
	Genres          []*string                        `json:"genres,omitempty"`
	ID              int                              `json:"id"`
	MeanScore       *int                             `json:"meanScore,omitempty"`
	Popularity      *int                             `json:"popularity,omitempty"`
	Rankings        []*AnimeDetails_Rankings         `json:"rankings,omitempty"`
	Recommendations *AnimeDetails_Recommendations    `json:"recommendations,omitempty"`
	Relations       *AnimeDetails_Relations          `json:"relations,omitempty"`
	SiteURL         *string                          `json:"siteUrl,omitempty"`
	Staff           *AnimeDetails_Staff              `json:"staff,omitempty"`
	StartDate       *AnimeDetails_StartDate          `json:"startDate,omitempty"`
	Studios         *AnimeDetails_Studios            `json:"studios,omitempty"`
	Trailer         *AnimeDetails_Trailer            `json:"trailer,omitempty"`
}

type AnimeDetails_Characters struct {
	Edges []*AnimeDetails_Characters_Edges `json:"edges,omitempty"`
}

type AnimeDetails_Characters_Edges struct {
	ID   *int           `json:"id,omitempty"`
	Name *string        `json:"name,omitempty"`
	Node *BaseCharacter `json:"node,omitempty"`
	Role *CharacterRole `json:"role,omitempty"`
}

type AnimeDetails_Staff struct {
	Edges []*AnimeDetails_Staff_Edges `json:"edges,omitempty"`
}

type AnimeDetails_Staff_Edges struct {
	Node *AnimeDetails_Staff_Edges_Node `json:"node,omitempty"`
	Role *string                        `json:"role,omitempty"`
}

type AnimeDetails_Staff_Edges_Node struct {
	ID   int                                `json:"id"`
	Name *AnimeDetails_Staff_Edges_Node_Name `json:"name,omitempty"`
}

type AnimeDetails_Staff_Edges_Node_Name struct {
	Full *string `json:"full,omitempty"`
}

type AnimeDetails_Rankings struct {
	AllTime *bool         `json:"allTime,omitempty"`
	Context string        `json:"context"`
	Format  MediaFormat   `json:"format"`
	Rank    int           `json:"rank"`
	Season  *MediaSeason  `json:"season,omitempty"`
	Type    MediaRankType `json:"type"`
	Year    *int          `json:"year,omitempty"`
}

type AnimeDetails_Recommendations struct {
	Edges []*AnimeDetails_Recommendations_Edges `json:"edges,omitempty"`
}

type AnimeDetails_Recommendations_Edges struct {
	Node *AnimeDetails_Recommendations_Edges_Node `json:"node,omitempty"`
}

type AnimeDetails_Recommendations_Edges_Node struct {
	// MediaRecommendation 简化说明：原类型为 BaseAnime 字段子集，直接复用 *Anime
	MediaRecommendation *Anime `json:"mediaRecommendation,omitempty"`
}

type AnimeDetails_Relations struct {
	Edges []*AnimeDetails_Relations_Edges `json:"edges,omitempty"`
}

type AnimeDetails_Relations_Edges struct {
	Node         *Anime         `json:"node,omitempty"`
	RelationType *MediaRelation `json:"relationType,omitempty"`
}

type AnimeDetails_Studios struct {
	Nodes []*AnimeDetails_Studios_Nodes `json:"nodes,omitempty"`
}

type AnimeDetails_Studios_Nodes struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (t *AnimeDetails) GetAverageScore() *int {
	if t == nil {
		t = &AnimeDetails{}
	}
	return t.AverageScore
}
func (t *AnimeDetails) GetCharacters() *AnimeDetails_Characters {
	if t == nil {
		t = &AnimeDetails{}
	}
	return t.Characters
}
func (t *AnimeDetails) GetDescription() *string {
	if t == nil {
		t = &AnimeDetails{}
	}
	return t.Description
}
func (t *AnimeDetails) GetDuration() *int {
	if t == nil {
		t = &AnimeDetails{}
	}
	return t.Duration
}
func (t *AnimeDetails) GetEndDate() *AnimeDetails_EndDate {
	if t == nil {
		t = &AnimeDetails{}
	}
	return t.EndDate
}
func (t *AnimeDetails) GetGenres() []*string {
	if t == nil {
		t = &AnimeDetails{}
	}
	return t.Genres
}
func (t *AnimeDetails) GetID() int {
	if t == nil {
		t = &AnimeDetails{}
	}
	return t.ID
}
func (t *AnimeDetails) GetMeanScore() *int {
	if t == nil {
		t = &AnimeDetails{}
	}
	return t.MeanScore
}
func (t *AnimeDetails) GetPopularity() *int {
	if t == nil {
		t = &AnimeDetails{}
	}
	return t.Popularity
}
func (t *AnimeDetails) GetRankings() []*AnimeDetails_Rankings {
	if t == nil {
		t = &AnimeDetails{}
	}
	return t.Rankings
}
func (t *AnimeDetails) GetRecommendations() *AnimeDetails_Recommendations {
	if t == nil {
		t = &AnimeDetails{}
	}
	return t.Recommendations
}
func (t *AnimeDetails) GetRelations() *AnimeDetails_Relations {
	if t == nil {
		t = &AnimeDetails{}
	}
	return t.Relations
}
func (t *AnimeDetails) GetSiteURL() *string {
	if t == nil {
		t = &AnimeDetails{}
	}
	return t.SiteURL
}
func (t *AnimeDetails) GetStaff() *AnimeDetails_Staff {
	if t == nil {
		t = &AnimeDetails{}
	}
	return t.Staff
}
func (t *AnimeDetails) GetStartDate() *AnimeDetails_StartDate {
	if t == nil {
		t = &AnimeDetails{}
	}
	return t.StartDate
}
func (t *AnimeDetails) GetStudios() *AnimeDetails_Studios {
	if t == nil {
		t = &AnimeDetails{}
	}
	return t.Studios
}
func (t *AnimeDetails) GetTrailer() *AnimeDetails_Trailer {
	if t == nil {
		t = &AnimeDetails{}
	}
	return t.Trailer
}

func (t *AnimeDetails_Relations) GetEdges() []*AnimeDetails_Relations_Edges {
	if t == nil {
		t = &AnimeDetails_Relations{}
	}
	return t.Edges
}
func (t *AnimeDetails_Relations_Edges) GetNode() *Anime {
	if t == nil {
		t = &AnimeDetails_Relations_Edges{}
	}
	return t.Node
}
func (t *AnimeDetails_Relations_Edges) GetRelationType() *MediaRelation {
	if t == nil {
		t = &AnimeDetails_Relations_Edges{}
	}
	return t.RelationType
}

func (t *AnimeDetails_Studios) GetNodes() []*AnimeDetails_Studios_Nodes {
	if t == nil {
		t = &AnimeDetails_Studios{}
	}
	return t.Nodes
}
func (t *AnimeDetails_Studios_Nodes) GetID() int {
	if t == nil {
		t = &AnimeDetails_Studios_Nodes{}
	}
	return t.ID
}
func (t *AnimeDetails_Studios_Nodes) GetName() string {
	if t == nil {
		t = &AnimeDetails_Studios_Nodes{}
	}
	return t.Name
}

// MangaDetails 条目详情（书籍类目：无 Staff/Studios/Trailer）
type MangaDetails struct {
	Characters      *MangaDetails_Characters        `json:"characters,omitempty"`
	Duration        *int                            `json:"duration,omitempty"`
	Genres          []*string                       `json:"genres,omitempty"`
	ID              int                             `json:"id"`
	Rankings        []*MangaDetails_Rankings        `json:"rankings,omitempty"`
	Recommendations *MangaDetails_Recommendations   `json:"recommendations,omitempty"`
	Relations       *MangaDetails_Relations         `json:"relations,omitempty"`
	SiteURL         *string                         `json:"siteUrl,omitempty"`
}

type MangaDetails_Characters struct {
	Edges []*MangaDetails_Characters_Edges `json:"edges,omitempty"`
}

type MangaDetails_Characters_Edges struct {
	ID   *int           `json:"id,omitempty"`
	Name *string        `json:"name,omitempty"`
	Node *BaseCharacter `json:"node,omitempty"`
	Role *CharacterRole `json:"role,omitempty"`
}

type MangaDetails_Rankings struct {
	AllTime *bool         `json:"allTime,omitempty"`
	Context string        `json:"context"`
	Format  MediaFormat   `json:"format"`
	Rank    int           `json:"rank"`
	Season  *MediaSeason  `json:"season,omitempty"`
	Type    MediaRankType `json:"type"`
	Year    *int          `json:"year,omitempty"`
}

type MangaDetails_Recommendations struct {
	Edges []*MangaDetails_Recommendations_Edges `json:"edges,omitempty"`
}

type MangaDetails_Recommendations_Edges struct {
	Node *MangaDetails_Recommendations_Edges_Node `json:"node,omitempty"`
}

type MangaDetails_Recommendations_Edges_Node struct {
	// MediaRecommendation 简化说明：原类型为 BaseManga 字段子集，直接复用 *Manga
	MediaRecommendation *Manga `json:"mediaRecommendation,omitempty"`
}

type MangaDetails_Relations struct {
	Edges []*MangaDetails_Relations_Edges `json:"edges,omitempty"`
}

type MangaDetails_Relations_Edges struct {
	Node         *Manga         `json:"node,omitempty"`
	RelationType *MediaRelation `json:"relationType,omitempty"`
}

func (t *MangaDetails) GetCharacters() *MangaDetails_Characters {
	if t == nil {
		t = &MangaDetails{}
	}
	return t.Characters
}
func (t *MangaDetails) GetDuration() *int {
	if t == nil {
		t = &MangaDetails{}
	}
	return t.Duration
}
func (t *MangaDetails) GetGenres() []*string {
	if t == nil {
		t = &MangaDetails{}
	}
	return t.Genres
}
func (t *MangaDetails) GetID() int {
	if t == nil {
		t = &MangaDetails{}
	}
	return t.ID
}
func (t *MangaDetails) GetRankings() []*MangaDetails_Rankings {
	if t == nil {
		t = &MangaDetails{}
	}
	return t.Rankings
}
func (t *MangaDetails) GetRecommendations() *MangaDetails_Recommendations {
	if t == nil {
		t = &MangaDetails{}
	}
	return t.Recommendations
}
func (t *MangaDetails) GetRelations() *MangaDetails_Relations {
	if t == nil {
		t = &MangaDetails{}
	}
	return t.Relations
}
func (t *MangaDetails) GetSiteURL() *string {
	if t == nil {
		t = &MangaDetails{}
	}
	return t.SiteURL
}

func (t *MangaDetails_Relations) GetEdges() []*MangaDetails_Relations_Edges {
	if t == nil {
		t = &MangaDetails_Relations{}
	}
	return t.Edges
}
func (t *MangaDetails_Relations_Edges) GetNode() *Manga {
	if t == nil {
		t = &MangaDetails_Relations_Edges{}
	}
	return t.Node
}
func (t *MangaDetails_Relations_Edges) GetRelationType() *MediaRelation {
	if t == nil {
		t = &MangaDetails_Relations_Edges{}
	}
	return t.RelationType
}

//----------------------------------------------------------------------------------------------------------------------
// 查询结果包装（镜像 anilist 的 ByID / Page 包装层）

type (
	BaseAnimeByID     struct{ Media *Anime `json:"Media,omitempty"` }
	BaseMangaByID     struct{ Media *Manga `json:"Media,omitempty"` }
	BaseAnimeByMalID  struct{ Media *Anime `json:"Media,omitempty"` }
	CompleteAnimeByID struct{ Media *CompleteAnime `json:"Media,omitempty"` }
	AnimeDetailsByID  struct{ Media *AnimeDetails `json:"Media,omitempty"` }
	MangaDetailsByID  struct{ Media *MangaDetails `json:"Media,omitempty"` }
)

func (t *BaseAnimeByID) GetMedia() *Anime {
	if t == nil {
		t = &BaseAnimeByID{}
	}
	return t.Media
}
func (t *BaseMangaByID) GetMedia() *Manga {
	if t == nil {
		t = &BaseMangaByID{}
	}
	return t.Media
}
func (t *BaseAnimeByMalID) GetMedia() *Anime {
	if t == nil {
		t = &BaseAnimeByMalID{}
	}
	return t.Media
}
func (t *CompleteAnimeByID) GetMedia() *CompleteAnime {
	if t == nil {
		t = &CompleteAnimeByID{}
	}
	return t.Media
}
func (t *AnimeDetailsByID) GetMedia() *AnimeDetails {
	if t == nil {
		t = &AnimeDetailsByID{}
	}
	return t.Media
}
func (t *MangaDetailsByID) GetMedia() *MangaDetails {
	if t == nil {
		t = &MangaDetailsByID{}
	}
	return t.Media
}

// PageInfo 分页信息（ListAnime/ListManga 等 Page 结构共用）
type PageInfo struct {
	CurrentPage *int  `json:"currentPage,omitempty"`
	HasNextPage *bool `json:"hasNextPage,omitempty"`
	LastPage    *int  `json:"lastPage,omitempty"`
	PerPage     *int  `json:"perPage,omitempty"`
	Total       *int  `json:"total,omitempty"`
}

func (t *PageInfo) GetCurrentPage() *int {
	if t == nil {
		t = &PageInfo{}
	}
	return t.CurrentPage
}
func (t *PageInfo) GetHasNextPage() *bool {
	if t == nil {
		t = &PageInfo{}
	}
	return t.HasNextPage
}
func (t *PageInfo) GetLastPage() *int {
	if t == nil {
		t = &PageInfo{}
	}
	return t.LastPage
}
func (t *PageInfo) GetPerPage() *int {
	if t == nil {
		t = &PageInfo{}
	}
	return t.PerPage
}
func (t *PageInfo) GetTotal() *int {
	if t == nil {
		t = &PageInfo{}
	}
	return t.Total
}

type (
	// ListAnime / ListManga 列表搜索结果；ListAnimeAll/ListMangaAll 形状相同直接别名
	ListAnime     struct{ Page *ListAnime_Page `json:"Page,omitempty"` }
	ListAnimeAll  = ListAnime
	ListManga     struct{ Page *ListManga_Page `json:"Page,omitempty"` }
	ListMangaAll  = ListManga
)

type ListAnime_Page struct {
	Media    []*Anime  `json:"media,omitempty"`
	PageInfo *PageInfo `json:"pageInfo,omitempty"`
}

type ListManga_Page struct {
	Media    []*Manga  `json:"media,omitempty"`
	PageInfo *PageInfo `json:"pageInfo,omitempty"`
}

func (t *ListAnime) GetPage() *ListAnime_Page {
	if t == nil {
		t = &ListAnime{}
	}
	return t.Page
}
func (t *ListManga) GetPage() *ListManga_Page {
	if t == nil {
		t = &ListManga{}
	}
	return t.Page
}
func (t *ListAnime_Page) GetMedia() []*Anime {
	if t == nil {
		t = &ListAnime_Page{}
	}
	return t.Media
}
func (t *ListManga_Page) GetMedia() []*Manga {
	if t == nil {
		t = &ListManga_Page{}
	}
	return t.Media
}
func (t *ListAnime_Page) GetPageInfo() *PageInfo {
	if t == nil {
		t = &ListAnime_Page{}
	}
	return t.PageInfo
}
func (t *ListManga_Page) GetPageInfo() *PageInfo {
	if t == nil {
		t = &ListManga_Page{}
	}
	return t.PageInfo
}

// SearchBaseAnimeByIds / SearchBaseManga 按 ID 批量搜索
type SearchBaseAnimeByIds struct {
	Page *SearchBaseAnimeByIds_Page `json:"Page,omitempty"`
}

type SearchBaseAnimeByIds_Page struct {
	Media    []*Anime `json:"media,omitempty"`
	PageInfo *SearchBaseAnimeByIds_Page_PageInfo `json:"pageInfo,omitempty"`
}

type SearchBaseAnimeByIds_Page_PageInfo struct {
	HasNextPage *bool `json:"hasNextPage,omitempty"`
}

type SearchBaseManga struct {
	Page *SearchBaseManga_Page `json:"Page,omitempty"`
}

type SearchBaseManga_Page struct {
	Media    []*Manga `json:"media,omitempty"`
	PageInfo *SearchBaseAnimeByIds_Page_PageInfo `json:"pageInfo,omitempty"`
}

func (t *SearchBaseAnimeByIds) GetPage() *SearchBaseAnimeByIds_Page {
	if t == nil {
		t = &SearchBaseAnimeByIds{}
	}
	return t.Page
}
func (t *SearchBaseAnimeByIds_Page) GetMedia() []*Anime {
	if t == nil {
		t = &SearchBaseAnimeByIds_Page{}
	}
	return t.Media
}
func (t *SearchBaseAnimeByIds_Page) GetPageInfo() *SearchBaseAnimeByIds_Page_PageInfo {
	if t == nil {
		t = &SearchBaseAnimeByIds_Page{}
	}
	return t.PageInfo
}
func (t *SearchBaseAnimeByIds_Page_PageInfo) GetHasNextPage() *bool {
	if t == nil {
		t = &SearchBaseAnimeByIds_Page_PageInfo{}
	}
	return t.HasNextPage
}
func (t *SearchBaseManga) GetPage() *SearchBaseManga_Page {
	if t == nil {
		t = &SearchBaseManga{}
	}
	return t.Page
}
func (t *SearchBaseManga_Page) GetMedia() []*Manga {
	if t == nil {
		t = &SearchBaseManga_Page{}
	}
	return t.Media
}
func (t *SearchBaseManga_Page) GetPageInfo() *SearchBaseAnimeByIds_Page_PageInfo {
	if t == nil {
		t = &SearchBaseManga_Page{}
	}
	return t.PageInfo
}

// ListRecentAnime 近期放送（含逐集放送计划）
type ListRecentAnime struct {
	Page *ListRecentAnime_Page `json:"Page,omitempty"`
}

type ListRecentAnime_Page struct {
	AiringSchedules []*ListRecentAnime_Page_AiringSchedules `json:"airingSchedules,omitempty"`
	PageInfo        *PageInfo                               `json:"pageInfo,omitempty"`
}

type ListRecentAnime_Page_AiringSchedules struct {
	AiringAt        int    `json:"airingAt"`
	Episode         int    `json:"episode"`
	ID              int    `json:"id"`
	Media           *Anime `json:"media,omitempty"`
	TimeUntilAiring int    `json:"timeUntilAiring"`
}

func (t *ListRecentAnime) GetPage() *ListRecentAnime_Page {
	if t == nil {
		t = &ListRecentAnime{}
	}
	return t.Page
}
func (t *ListRecentAnime_Page) GetAiringSchedules() []*ListRecentAnime_Page_AiringSchedules {
	if t == nil {
		t = &ListRecentAnime_Page{}
	}
	return t.AiringSchedules
}
func (t *ListRecentAnime_Page) GetPageInfo() *PageInfo {
	if t == nil {
		t = &ListRecentAnime_Page{}
	}
	return t.PageInfo
}

// 变更响应包装（镜像 anilist UpdateMediaListEntry 系列）
type (
	UpdateMediaListEntry struct {
		SaveMediaListEntry *UpdateMediaListEntry_SaveMediaListEntry `json:"SaveMediaListEntry,omitempty"`
	}
	UpdateMediaListEntry_SaveMediaListEntry struct {
		ID int `json:"id"`
	}
	UpdateMediaListEntryProgress struct {
		SaveMediaListEntry *UpdateMediaListEntryProgress_SaveMediaListEntry `json:"SaveMediaListEntry,omitempty"`
	}
	UpdateMediaListEntryProgress_SaveMediaListEntry struct {
		ID int `json:"id"`
	}
	UpdateMediaListEntryRepeat struct {
		SaveMediaListEntry *UpdateMediaListEntryRepeat_SaveMediaListEntry `json:"SaveMediaListEntry,omitempty"`
	}
	UpdateMediaListEntryRepeat_SaveMediaListEntry struct {
		ID int `json:"id"`
	}
	DeleteEntry struct {
		DeleteMediaListEntry *DeleteEntry_DeleteMediaListEntry `json:"DeleteMediaListEntry,omitempty"`
	}
	DeleteEntry_DeleteMediaListEntry struct {
		Deleted *bool `json:"deleted,omitempty"`
	}
)

func (t *UpdateMediaListEntry) GetSaveMediaListEntry() *UpdateMediaListEntry_SaveMediaListEntry {
	if t == nil {
		t = &UpdateMediaListEntry{}
	}
	return t.SaveMediaListEntry
}
func (t *UpdateMediaListEntryProgress) GetSaveMediaListEntry() *UpdateMediaListEntryProgress_SaveMediaListEntry {
	if t == nil {
		t = &UpdateMediaListEntryProgress{}
	}
	return t.SaveMediaListEntry
}
func (t *UpdateMediaListEntryRepeat) GetSaveMediaListEntry() *UpdateMediaListEntryRepeat_SaveMediaListEntry {
	if t == nil {
		t = &UpdateMediaListEntryRepeat{}
	}
	return t.SaveMediaListEntry
}
func (t *DeleteEntry) GetDeleteMediaListEntry() *DeleteEntry_DeleteMediaListEntry {
	if t == nil {
		t = &DeleteEntry{}
	}
	return t.DeleteMediaListEntry
}

//----------------------------------------------------------------------------------------------------------------------
// 放送时间表（镜像 anilist.AnimeSchedule 与 AnimeAiringSchedule 系列）

// AnimeSchedule 单条目的前后放送节点
type AnimeSchedule struct {
	ID       int                     `json:"id"`
	IDMal    *int                    `json:"idMal,omitempty"`
	Previous *AnimeSchedule_Previous `json:"previous,omitempty"`
	Upcoming *AnimeSchedule_Upcoming `json:"upcoming,omitempty"`
}

type AnimeSchedule_Previous struct {
	Nodes []*AnimeSchedule_Previous_Nodes `json:"nodes,omitempty"`
}

type AnimeSchedule_Upcoming struct {
	Nodes []*AnimeSchedule_Upcoming_Nodes `json:"nodes,omitempty"`
}

func (t *AnimeSchedule) GetID() int {
	if t == nil {
		t = &AnimeSchedule{}
	}
	return t.ID
}
func (t *AnimeSchedule) GetIDMal() *int {
	if t == nil {
		t = &AnimeSchedule{}
	}
	return t.IDMal
}
func (t *AnimeSchedule) GetPrevious() *AnimeSchedule_Previous {
	if t == nil {
		t = &AnimeSchedule{}
	}
	return t.Previous
}
func (t *AnimeSchedule) GetUpcoming() *AnimeSchedule_Upcoming {
	if t == nil {
		t = &AnimeSchedule{}
	}
	return t.Upcoming
}
func (t *AnimeSchedule_Previous) GetNodes() []*AnimeSchedule_Previous_Nodes {
	if t == nil {
		t = &AnimeSchedule_Previous{}
	}
	return t.Nodes
}
func (t *AnimeSchedule_Upcoming) GetNodes() []*AnimeSchedule_Upcoming_Nodes {
	if t == nil {
		t = &AnimeSchedule_Upcoming{}
	}
	return t.Nodes
}

// schedulePage 放送时间表分页包装的公共形状：{media: AnimeSchedule}
type schedulePage struct {
	Media *AnimeSchedule `json:"media,omitempty"`
}

func (t *schedulePage) GetMedia() *AnimeSchedule {
	if t == nil {
		t = &schedulePage{}
	}
	return t.Media
}

// AnimeAiringSchedule 分季放送时间表。
// Ongoing/OngoingNext/Upcoming/UpcomingNext/Preceding 原为五个同形包装类型
// （AnimeAiringSchedule_Ongoing 等），此处合并为 schedulePage 别名。
type AnimeAiringSchedule struct {
	Ongoing      *schedulePage `json:"ongoing,omitempty"`
	OngoingNext  *schedulePage `json:"ongoingNext,omitempty"`
	Upcoming     *schedulePage `json:"upcoming,omitempty"`
	UpcomingNext *schedulePage `json:"upcomingNext,omitempty"`
	Preceding    *schedulePage `json:"preceding,omitempty"`
}

type (
	AnimeAiringSchedule_Ongoing      = schedulePage
	AnimeAiringSchedule_OngoingNext  = schedulePage
	AnimeAiringSchedule_Upcoming     = schedulePage
	AnimeAiringSchedule_UpcomingNext = schedulePage
	AnimeAiringSchedule_Preceding    = schedulePage
)

func (t *AnimeAiringSchedule) GetOngoing() *schedulePage {
	if t == nil {
		t = &AnimeAiringSchedule{}
	}
	return t.Ongoing
}
func (t *AnimeAiringSchedule) GetOngoingNext() *schedulePage {
	if t == nil {
		t = &AnimeAiringSchedule{}
	}
	return t.OngoingNext
}
func (t *AnimeAiringSchedule) GetUpcoming() *schedulePage {
	if t == nil {
		t = &AnimeAiringSchedule{}
	}
	return t.Upcoming
}
func (t *AnimeAiringSchedule) GetUpcomingNext() *schedulePage {
	if t == nil {
		t = &AnimeAiringSchedule{}
	}
	return t.UpcomingNext
}
func (t *AnimeAiringSchedule) GetPreceding() *schedulePage {
	if t == nil {
		t = &AnimeAiringSchedule{}
	}
	return t.Preceding
}

type AnimeAiringScheduleRaw struct {
	Page *AnimeAiringScheduleRaw_Page `json:"Page,omitempty"`
}

type AnimeAiringScheduleRaw_Page struct {
	Media []*AnimeSchedule `json:"media,omitempty"`
}

func (t *AnimeAiringScheduleRaw) GetPage() *AnimeAiringScheduleRaw_Page {
	if t == nil {
		t = &AnimeAiringScheduleRaw{}
	}
	return t.Page
}
func (t *AnimeAiringScheduleRaw_Page) GetMedia() []*AnimeSchedule {
	if t == nil {
		t = &AnimeAiringScheduleRaw_Page{}
	}
	return t.Media
}

//----------------------------------------------------------------------------------------------------------------------
// 缓存与关系树（镜像 anilist/media.go、media_tree.go 中的类型部分）

// BaseAnimeCache / CompleteAnimeCache FetchMediaTree 结果临时缓存
type (
	BaseAnimeCache struct{ *result.Cache[int, *Anime] }

	CompleteAnimeCache struct{ *result.Cache[int, *CompleteAnime] }
)

// NewBaseAnimeCache returns a new result.Cache[int, *Anime].
func NewBaseAnimeCache() *BaseAnimeCache {
	return &BaseAnimeCache{result.NewCache[int, *Anime]()}
}

// NewCompleteAnimeCache returns a new result.Cache[int, *CompleteAnime].
func NewCompleteAnimeCache() *CompleteAnimeCache {
	return &CompleteAnimeCache{result.NewCache[int, *CompleteAnime]()}
}

// CompleteAnimeRelationTree 关系树（scanner 等用于存放跨作品遍历结果）
type CompleteAnimeRelationTree struct {
	*result.Map[int, *CompleteAnime]
}

// NewCompleteAnimeRelationTree returns a new result.Map[int, *CompleteAnime].
func NewCompleteAnimeRelationTree() *CompleteAnimeRelationTree {
	return &CompleteAnimeRelationTree{result.NewMap[int, *CompleteAnime]()}
}

// FetchMediaTreeRelation 关系树遍历方向。
// 注意：FetchMediaTree 逐级抓取函数本身被裁剪（依赖 AniList GraphQL 客户端），
// Bangumi 侧由 Wave B 用 GET /v0/subjects/{id}/subjects 一次性建树。
type FetchMediaTreeRelation = string

const (
	FetchMediaTreeSequels  FetchMediaTreeRelation = "sequels"
	FetchMediaTreePrequels FetchMediaTreeRelation = "prequels"
	FetchMediaTreeAll      FetchMediaTreeRelation = "all"
)

//----------------------------------------------------------------------------------------------------------------------
// 季度工具（照搬 anilist/utils.go GetSeasonInfo）

type GetSeasonKind int

const (
	GetSeasonKindCurrent GetSeasonKind = iota
	GetSeasonKindNext
	GetSeasonKindPrevious
)

// GetSeasonInfo 返回给定时间的季度与年份。
func GetSeasonInfo(now time.Time, kind GetSeasonKind) (MediaSeason, int) {
	month, year := now.Month(), now.Year()

	getSeasonIndex := func(m time.Month) int {
		switch {
		case m >= 3 && m <= 5: // spring: 3, 4, 5
			return 1
		case m >= 6 && m <= 8: // summer: 6, 7, 8
			return 2
		case m >= 9 && m <= 11: // fall: 9, 10, 11
			return 3
		default: // winter: 12, 1, 2
			return 0
		}
	}

	seasons := []MediaSeason{MediaSeasonWinter, MediaSeasonSpring, MediaSeasonSummer, MediaSeasonFall}
	var index int

	switch kind {
	case GetSeasonKindCurrent:
		index = getSeasonIndex(month)

	case GetSeasonKindNext:
		nextMonth := month + 3
		nextYear := year
		if nextMonth > 12 {
			nextMonth -= 12
			nextYear++
		}
		index = getSeasonIndex(nextMonth)
		year = nextYear

	case GetSeasonKindPrevious:
		prevMonth := month - 3
		prevYear := year
		if prevMonth <= 0 {
			prevMonth += 12
			prevYear--
		}
		index = getSeasonIndex(prevMonth)
		year = prevYear
	}

	return seasons[index], year
}

// ListAnimeCacheKey / ListMangaCacheKey 照搬 anilist/list.go（离线缓存键）
func ListAnimeCacheKey(
	Page *int,
	Search *string,
	PerPage *int,
	Sort []*MediaSort,
	Status []*MediaStatus,
	Genres []*string,
	Tags []*string,
	AverageScoreGreater *int,
	Season *MediaSeason,
	SeasonYear *int,
	Format *MediaFormat,
	IsAdult *bool,
	CountryOfOrigin *string,
) string {
	key := "ListAnime"
	if Page != nil {
		key += fmt.Sprintf("_%d", *Page)
	}
	if Search != nil {
		key += fmt.Sprintf("_%s", *Search)
	}
	if PerPage != nil {
		key += fmt.Sprintf("_%d", *PerPage)
	}
	if Sort != nil {
		key += fmt.Sprintf("_%v", Sort)
	}
	if Status != nil {
		key += fmt.Sprintf("_%v", Status)
	}
	if Genres != nil {
		key += fmt.Sprintf("_%v", Genres)
	}
	if Tags != nil {
		key += fmt.Sprintf("_%v", Tags)
	}
	if AverageScoreGreater != nil {
		key += fmt.Sprintf("_%d", *AverageScoreGreater)
	}
	if Season != nil {
		key += fmt.Sprintf("_%s", *Season)
	}
	if SeasonYear != nil {
		key += fmt.Sprintf("_%d", *SeasonYear)
	}
	if Format != nil {
		key += fmt.Sprintf("_%s", *Format)
	}
	if IsAdult != nil {
		key += fmt.Sprintf("_%t", *IsAdult)
	}
	if CountryOfOrigin != nil {
		key += fmt.Sprintf("_%s", *CountryOfOrigin)
	}
	return key
}

func ListMangaCacheKey(
	Page *int,
	Search *string,
	PerPage *int,
	Sort []*MediaSort,
	Status []*MediaStatus,
	Genres []*string,
	Tags []*string,
	AverageScoreGreater *int,
	Season *MediaSeason,
	SeasonYear *int,
	Format *MediaFormat,
	CountryOfOrigin *string,
	IsAdult *bool,
) string {
	key := "ListManga"
	if Page != nil {
		key += fmt.Sprintf("_%d", *Page)
	}
	if Search != nil {
		key += fmt.Sprintf("_%s", *Search)
	}
	if PerPage != nil {
		key += fmt.Sprintf("_%d", *PerPage)
	}
	if Sort != nil {
		key += fmt.Sprintf("_%v", Sort)
	}
	if Status != nil {
		key += fmt.Sprintf("_%v", Status)
	}
	if Genres != nil {
		key += fmt.Sprintf("_%v", Genres)
	}
	if Tags != nil {
		key += fmt.Sprintf("_%v", Tags)
	}
	if AverageScoreGreater != nil {
		key += fmt.Sprintf("_%d", *AverageScoreGreater)
	}
	if Season != nil {
		key += fmt.Sprintf("_%s", *Season)
	}
	if SeasonYear != nil {
		key += fmt.Sprintf("_%d", *SeasonYear)
	}
	if Format != nil {
		key += fmt.Sprintf("_%s", *Format)
	}
	if CountryOfOrigin != nil {
		key += fmt.Sprintf("_%s", *CountryOfOrigin)
	}
	if IsAdult != nil {
		key += fmt.Sprintf("_%t", *IsAdult)
	}
	return key
}
