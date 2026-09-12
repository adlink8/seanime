package media

// enums.go —— AniList 领域枚举的镜像层（Wave A 换锚手术）。
//
// 镜像规则（见 .planning/phases/02-library-reanchor/02-CONTRACT.md）：
//   - 类型名与常量名与 internal/api/anilist 中定义一字不差，
//     仅 BaseAnime→Anime、BaseManga→Manga 两个例外；
//   - 字符串值完全一致（"CURRENT"、"TV"、"FINISHED"…），保证前端序列化形状不变；
//   - 省略 gqlgen 生成的 UnmarshalGQL/MarshalGQL 方法：命名字符串类型的
//     默认 JSON 编解码行为与原实现等价（均产出带引号字符串）；
//   - AllXxx 切片与 IsValid() 保留，供消费方遍历/校验。
//
// Bangumi 扩展（additive）：MediaFormatBook = "BOOK"，承载 Bangumi 书籍类目
// 的保守推导结果（AniList 无 BOOK 值； novels 用 NOVEL，漫画用 MANGA）。

// MediaListStatus 收藏/观看状态
type MediaListStatus string

const (
	// MediaListStatusCurrent 在看/在读
	MediaListStatusCurrent MediaListStatus = "CURRENT"
	// MediaListStatusPlanning 想看/想读
	MediaListStatusPlanning MediaListStatus = "PLANNING"
	// MediaListStatusCompleted 看过/读过
	MediaListStatusCompleted MediaListStatus = "COMPLETED"
	// MediaListStatusDropped 抛弃
	MediaListStatusDropped MediaListStatus = "DROPPED"
	// MediaListStatusPaused 搁置
	MediaListStatusPaused MediaListStatus = "PAUSED"
	// MediaListStatusRepeating 重看/重读
	MediaListStatusRepeating MediaListStatus = "REPEATING"
)

var AllMediaListStatus = []MediaListStatus{
	MediaListStatusCurrent,
	MediaListStatusPlanning,
	MediaListStatusCompleted,
	MediaListStatusDropped,
	MediaListStatusPaused,
	MediaListStatusRepeating,
}

func (e MediaListStatus) IsValid() bool {
	switch e {
	case MediaListStatusCurrent, MediaListStatusPlanning, MediaListStatusCompleted,
		MediaListStatusDropped, MediaListStatusPaused, MediaListStatusRepeating:
		return true
	}
	return false
}

func (e MediaListStatus) String() string { return string(e) }

// MediaStatus 放送/连载状态
type MediaStatus string

const (
	// MediaStatusFinished 已完结
	MediaStatusFinished MediaStatus = "FINISHED"
	// MediaStatusReleasing 放送/连载中
	MediaStatusReleasing MediaStatus = "RELEASING"
	// MediaStatusNotYetReleased 未放售
	MediaStatusNotYetReleased MediaStatus = "NOT_YET_RELEASED"
	// MediaStatusCancelled 腰斩
	MediaStatusCancelled MediaStatus = "CANCELLED"
	// MediaStatusHiatus 暂停更新
	MediaStatusHiatus MediaStatus = "HIATUS"
)

var AllMediaStatus = []MediaStatus{
	MediaStatusFinished,
	MediaStatusReleasing,
	MediaStatusNotYetReleased,
	MediaStatusCancelled,
	MediaStatusHiatus,
}

func (e MediaStatus) IsValid() bool {
	switch e {
	case MediaStatusFinished, MediaStatusReleasing, MediaStatusNotYetReleased,
		MediaStatusCancelled, MediaStatusHiatus:
		return true
	}
	return false
}

func (e MediaStatus) String() string { return string(e) }

// MediaFormat 条目形式
type MediaFormat string

const (
	// MediaFormatTv 电视动画
	MediaFormatTv MediaFormat = "TV"
	// MediaFormatTvShort 15 分钟以内的电视动画
	MediaFormatTvShort MediaFormat = "TV_SHORT"
	// MediaFormatMovie 剧场版
	MediaFormatMovie MediaFormat = "MOVIE"
	// MediaFormatSpecial 特别篇
	MediaFormatSpecial MediaFormat = "SPECIAL"
	// MediaFormatOva OVA
	MediaFormatOva MediaFormat = "OVA"
	// MediaFormatOna ONA
	MediaFormatOna MediaFormat = "ONA"
	// MediaFormatMusic MV
	MediaFormatMusic MediaFormat = "MUSIC"
	// MediaFormatManga 漫画
	MediaFormatManga MediaFormat = "MANGA"
	// MediaFormatNovel 轻小说/小说
	MediaFormatNovel MediaFormat = "NOVEL"
	// MediaFormatOneShot 短篇
	MediaFormatOneShot MediaFormat = "ONE_SHOT"

	// MediaFormatBook —— Bangumi 扩展（additive）：书籍类目的保守推导值。
	// AniList 无 "BOOK"；当书籍类目无法进一步区分为漫画/小说时使用。
	MediaFormatBook MediaFormat = "BOOK"
)

var AllMediaFormat = []MediaFormat{
	MediaFormatTv,
	MediaFormatTvShort,
	MediaFormatMovie,
	MediaFormatSpecial,
	MediaFormatOva,
	MediaFormatOna,
	MediaFormatMusic,
	MediaFormatManga,
	MediaFormatNovel,
	MediaFormatOneShot,
	MediaFormatBook,
}

func (e MediaFormat) IsValid() bool {
	switch e {
	case MediaFormatTv, MediaFormatTvShort, MediaFormatMovie, MediaFormatSpecial,
		MediaFormatOva, MediaFormatOna, MediaFormatMusic, MediaFormatManga,
		MediaFormatNovel, MediaFormatOneShot, MediaFormatBook:
		return true
	}
	return false
}

func (e MediaFormat) String() string { return string(e) }

// MediaSeason 放送季度
type MediaSeason string

const (
	// MediaSeasonWinter 12 月-2 月
	MediaSeasonWinter MediaSeason = "WINTER"
	// MediaSeasonSpring 3 月-5 月
	MediaSeasonSpring MediaSeason = "SPRING"
	// MediaSeasonSummer 6 月-8 月
	MediaSeasonSummer MediaSeason = "SUMMER"
	// MediaSeasonFall 9 月-11 月
	MediaSeasonFall MediaSeason = "FALL"
)

var AllMediaSeason = []MediaSeason{MediaSeasonWinter, MediaSeasonSpring, MediaSeasonSummer, MediaSeasonFall}

func (e MediaSeason) IsValid() bool {
	switch e {
	case MediaSeasonWinter, MediaSeasonSpring, MediaSeasonSummer, MediaSeasonFall:
		return true
	}
	return false
}

func (e MediaSeason) String() string { return string(e) }

// MediaSort 排序方式
type MediaSort string

const (
	MediaSortID               MediaSort = "ID"
	MediaSortIDDesc           MediaSort = "ID_DESC"
	MediaSortTitleRomaji      MediaSort = "TITLE_ROMAJI"
	MediaSortTitleRomajiDesc  MediaSort = "TITLE_ROMAJI_DESC"
	MediaSortTitleEnglish     MediaSort = "TITLE_ENGLISH"
	MediaSortTitleEnglishDesc MediaSort = "TITLE_ENGLISH_DESC"
	MediaSortTitleNative      MediaSort = "TITLE_NATIVE"
	MediaSortTitleNativeDesc  MediaSort = "TITLE_NATIVE_DESC"
	MediaSortType             MediaSort = "TYPE"
	MediaSortTypeDesc         MediaSort = "TYPE_DESC"
	MediaSortFormat           MediaSort = "FORMAT"
	MediaSortFormatDesc       MediaSort = "FORMAT_DESC"
	MediaSortStartDate        MediaSort = "START_DATE"
	MediaSortStartDateDesc    MediaSort = "START_DATE_DESC"
	MediaSortEndDate          MediaSort = "END_DATE"
	MediaSortEndDateDesc      MediaSort = "END_DATE_DESC"
	MediaSortScore            MediaSort = "SCORE"
	MediaSortScoreDesc        MediaSort = "SCORE_DESC"
	MediaSortPopularity       MediaSort = "POPULARITY"
	MediaSortPopularityDesc   MediaSort = "POPULARITY_DESC"
	MediaSortTrending         MediaSort = "TRENDING"
	MediaSortTrendingDesc     MediaSort = "TRENDING_DESC"
	MediaSortEpisodes         MediaSort = "EPISODES"
	MediaSortEpisodesDesc     MediaSort = "EPISODES_DESC"
	MediaSortDuration         MediaSort = "DURATION"
	MediaSortDurationDesc     MediaSort = "DURATION_DESC"
	MediaSortStatus           MediaSort = "STATUS"
	MediaSortStatusDesc       MediaSort = "STATUS_DESC"
	MediaSortChapters         MediaSort = "CHAPTERS"
	MediaSortChaptersDesc     MediaSort = "CHAPTERS_DESC"
	MediaSortVolumes          MediaSort = "VOLUMES"
	MediaSortVolumesDesc      MediaSort = "VOLUMES_DESC"
	MediaSortUpdatedAt        MediaSort = "UPDATED_AT"
	MediaSortUpdatedAtDesc    MediaSort = "UPDATED_AT_DESC"
	MediaSortSearchMatch      MediaSort = "SEARCH_MATCH"
	MediaSortFavourites       MediaSort = "FAVOURITES"
	MediaSortFavouritesDesc   MediaSort = "FAVOURITES_DESC"
)

var AllMediaSort = []MediaSort{
	MediaSortID, MediaSortIDDesc,
	MediaSortTitleRomaji, MediaSortTitleRomajiDesc,
	MediaSortTitleEnglish, MediaSortTitleEnglishDesc,
	MediaSortTitleNative, MediaSortTitleNativeDesc,
	MediaSortType, MediaSortTypeDesc,
	MediaSortFormat, MediaSortFormatDesc,
	MediaSortStartDate, MediaSortStartDateDesc,
	MediaSortEndDate, MediaSortEndDateDesc,
	MediaSortScore, MediaSortScoreDesc,
	MediaSortPopularity, MediaSortPopularityDesc,
	MediaSortTrending, MediaSortTrendingDesc,
	MediaSortEpisodes, MediaSortEpisodesDesc,
	MediaSortDuration, MediaSortDurationDesc,
	MediaSortStatus, MediaSortStatusDesc,
	MediaSortChapters, MediaSortChaptersDesc,
	MediaSortVolumes, MediaSortVolumesDesc,
	MediaSortUpdatedAt, MediaSortUpdatedAtDesc,
	MediaSortSearchMatch,
	MediaSortFavourites, MediaSortFavouritesDesc,
}

func (e MediaSort) IsValid() bool {
	for _, v := range AllMediaSort {
		if e == v {
			return true
		}
	}
	return false
}

func (e MediaSort) String() string { return string(e) }

// MediaType 条目大类（动画/漫画）
type MediaType string

const (
	// MediaTypeAnime 动画
	MediaTypeAnime MediaType = "ANIME"
	// MediaTypeManga 漫画（书籍类目统称）
	MediaTypeManga MediaType = "MANGA"
)

var AllMediaType = []MediaType{MediaTypeAnime, MediaTypeManga}

func (e MediaType) IsValid() bool {
	switch e {
	case MediaTypeAnime, MediaTypeManga:
		return true
	}
	return false
}

func (e MediaType) String() string { return string(e) }

// MediaRelation 条目关系类型
type MediaRelation string

const (
	// MediaRelationAdaptation 不同形式的改编
	MediaRelationAdaptation MediaRelation = "ADAPTATION"
	// MediaRelationPrequel 前传
	MediaRelationPrequel MediaRelation = "PREQUEL"
	// MediaRelationSequel 续作
	MediaRelationSequel MediaRelation = "SEQUEL"
	// MediaRelationParent 主线
	MediaRelationParent MediaRelation = "PARENT"
	// MediaRelationSideStory 支线
	MediaRelationSideStory MediaRelation = "SIDE_STORY"
	// MediaRelationCharacter 共享角色
	MediaRelationCharacter MediaRelation = "CHARACTER"
	// MediaRelationSummary 精简版
	MediaRelationSummary MediaRelation = "SUMMARY"
	// MediaRelationAlternative 替代版本
	MediaRelationAlternative MediaRelation = "ALTERNATIVE"
	// MediaRelationSpinOff 衍生作
	MediaRelationSpinOff MediaRelation = "SPIN_OFF"
	// MediaRelationOther 其他
	MediaRelationOther MediaRelation = "OTHER"
	// MediaRelationSource 原作
	MediaRelationSource MediaRelation = "SOURCE"
	// MediaRelationCompilation 合集
	MediaRelationCompilation MediaRelation = "COMPILATION"
	// MediaRelationContains 包含
	MediaRelationContains MediaRelation = "CONTAINS"
)

var AllMediaRelation = []MediaRelation{
	MediaRelationAdaptation, MediaRelationPrequel, MediaRelationSequel, MediaRelationParent,
	MediaRelationSideStory, MediaRelationCharacter, MediaRelationSummary, MediaRelationAlternative,
	MediaRelationSpinOff, MediaRelationOther, MediaRelationSource, MediaRelationCompilation,
	MediaRelationContains,
}

func (e MediaRelation) IsValid() bool {
	for _, v := range AllMediaRelation {
		if e == v {
			return true
		}
	}
	return false
}

func (e MediaRelation) String() string { return string(e) }

// MediaRankType 排名依据
type MediaRankType string

const (
	// MediaRankTypeRated 按评分
	MediaRankTypeRated MediaRankType = "RATED"
	// MediaRankTypePopular 按热度
	MediaRankTypePopular MediaRankType = "POPULAR"
)

var AllMediaRankType = []MediaRankType{MediaRankTypeRated, MediaRankTypePopular}

func (e MediaRankType) IsValid() bool {
	switch e {
	case MediaRankTypeRated, MediaRankTypePopular:
		return true
	}
	return false
}

func (e MediaRankType) String() string { return string(e) }

// CharacterRole 角色定位
type CharacterRole string

const (
	CharacterRoleMain       CharacterRole = "MAIN"
	CharacterRoleSupporting CharacterRole = "SUPPORTING"
	CharacterRoleBackground CharacterRole = "BACKGROUND"
)

var AllCharacterRole = []CharacterRole{CharacterRoleMain, CharacterRoleSupporting, CharacterRoleBackground}

func (e CharacterRole) IsValid() bool {
	switch e {
	case CharacterRoleMain, CharacterRoleSupporting, CharacterRoleBackground:
		return true
	}
	return false
}

func (e CharacterRole) String() string { return string(e) }

// AiringSort 放送时间表排序
type AiringSort string

const (
	AiringSortID          AiringSort = "ID"
	AiringSortIDDesc      AiringSort = "ID_DESC"
	AiringSortMediaID     AiringSort = "MEDIA_ID"
	AiringSortMediaIDDesc AiringSort = "MEDIA_ID_DESC"
	AiringSortTime        AiringSort = "TIME"
	AiringSortTimeDesc    AiringSort = "TIME_DESC"
	AiringSortEpisode     AiringSort = "EPISODE"
	AiringSortEpisodeDesc AiringSort = "EPISODE_DESC"
)

var AllAiringSort = []AiringSort{
	AiringSortID, AiringSortIDDesc,
	AiringSortMediaID, AiringSortMediaIDDesc,
	AiringSortTime, AiringSortTimeDesc,
	AiringSortEpisode, AiringSortEpisodeDesc,
}

func (e AiringSort) IsValid() bool {
	for _, v := range AllAiringSort {
		if e == v {
			return true
		}
	}
	return false
}

func (e AiringSort) String() string { return string(e) }
