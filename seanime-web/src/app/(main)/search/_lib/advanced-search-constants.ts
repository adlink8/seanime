import type { AL_MediaFormat } from "@/api/generated/types"

export const ADVANCED_SEARCH_MEDIA_GENRES = [
    "Action",
    "Adventure",
    "Comedy",
    "Drama",
    "Ecchi",
    "Fantasy",
    "Horror",
    "Mahou Shoujo",
    "Mecha",
    "Music",
    "Mystery",
    "Psychological",
    "Romance",
    "Sci-Fi",
    "Slice of Life",
    "Sports",
    "Supernatural",
    "Thriller",
]

export const GENRE_TRANSLATIONS: Record<string, string> = {
    "Action": "动作",
    "Adventure": "冒险",
    "Comedy": "喜剧",
    "Drama": "剧情",
    "Ecchi": "福利",
    "Fantasy": "奇幻",
    "Horror": "恐怖",
    "Mahou Shoujo": "魔法少女",
    "Mecha": "机战",
    "Music": "音乐",
    "Mystery": "悬疑",
    "Psychological": "心理",
    "Romance": "恋爱",
    "Sci-Fi": "科幻",
    "Slice of Life": "日常",
    "Sports": "运动",
    "Supernatural": "超自然",
    "Thriller": "惊悚",
}

export const ADVANCED_SEARCH_SEASONS = [
    "Winter",
    "Spring",
    "Summer",
    "Fall",
]

export const SEASON_TRANSLATIONS: Record<string, string> = {
    "Winter": "冬季",
    "Spring": "春季",
    "Summer": "夏季",
    "Fall": "秋季",
    "WINTER": "冬季",
    "SPRING": "春季",
    "SUMMER": "夏季",
    "FALL": "秋季",
}

export const ADVANCED_SEARCH_FORMATS: { value: AL_MediaFormat, label: string }[] = [
    { value: "TV", label: "TV 动画" },
    { value: "MOVIE", label: "剧场版 (Movie)" },
    { value: "ONA", label: "网络动画 (ONA)" },
    { value: "OVA", label: "原创录像 (OVA)" },
    { value: "TV_SHORT", label: "泡面番 (TV Short)" },
    { value: "SPECIAL", label: "特别篇 (Special)" },
]

/**
 * 高级搜索页的动画格式选项（契约 3.6b §1 D4 / §0.2）。
 * 上游 `meta_tags` 仅有 TV / WEB / OVA 三个值有效（ONA → WEB），
 * 故 MOVIE / TV_SHORT / SPECIAL 不再出现在搜索页。
 * 注意：给本地媒体库用的 `ADVANCED_SEARCH_FORMATS`（含 MOVIE 等）保持不变。
 */
export const ADVANCED_SEARCH_FORMATS_ANIME: { value: AL_MediaFormat, label: string }[] = [
    { value: "TV", label: "TV 动画" },
    { value: "ONA", label: "网络动画 (ONA)" },
    { value: "OVA", label: "原创录像 (OVA)" },
]

export const ADVANCED_SEARCH_FORMATS_MANGA: { value: AL_MediaFormat, label: string }[] = [
    { value: "MANGA", label: "连载漫画" },
    { value: "ONE_SHOT", label: "短篇单行本 (One Shot)" },
]

// 契约 3.6b §1 D7：上游仅支持 heat / rank / score（+ 关键词非空时的 match），
// 故 START_DATE_DESC / EPISODES_DESC / CHAPTERS_DESC 不再作为选项。
export const ADVANCED_SEARCH_SORTING = [
    { value: "TRENDING_DESC", label: "趋势热度" },
    { value: "POPULARITY_DESC", label: "最受欢迎" },
    { value: "SCORE_DESC", label: "最高评分" },
]

export const ADVANCED_SEARCH_SORTING_MANGA = [
    { value: "TRENDING_DESC", label: "趋势热度" },
    { value: "POPULARITY_DESC", label: "最受欢迎" },
    { value: "SCORE_DESC", label: "最高评分" },
]

export const ADVANCED_SEARCH_TYPE = [
    { value: "anime", label: "动漫" },
    { value: "manga", label: "漫画" },
    // Phase 2.5：新增轻小说 / 音声搜索类型
    { value: "novel", label: "轻小说" },
    { value: "asmr", label: "音声" },
]

// Phase 2.5：asmr.one 排序映射（SCORE_DESC→rating 评分、POPULARITY/TRENDING→dl 下载热度、START_DATE→dd 最新）
// Phase 3.9b D9：asmr.one search 端点 order=dc 实测 400，order=rating 合法且默认降序；
// 故「评分最高」值由 dc 改为 rating（dc 仍被后端兼容映射为 rating）。
export const ADVANCED_SEARCH_SORTING_ASMR = [
    { value: "dl", label: "下载热度" },
    { value: "rating", label: "评分最高" },
    { value: "dd", label: "最新发布" },
    { value: "publish_date", label: "发行日期" },
]

// Phase 2.5：asmr 字幕过滤选项（契约第 4 节 subtitle: "none" | "jp" | "zh"）
export const ADVANCED_SEARCH_SUBTITLE_ASMR = [
    { value: "zh", label: "中文字幕" },
    { value: "jp", label: "日文字幕" },
]

/**
 * Phase 2.5：将 AL 排序值映射为 asmr.one 的 order 参数
 * Phase 3.9b D9：SCORE_DESC → "rating"（asmr.one order=rating 合法且默认降序；dc 已弃用，后端兼容映射 dc→rating）
 */
export function mapSortingToAsmrOrder(sorting: string | null | undefined): string {
    switch (sorting) {
        case "POPULARITY_DESC":
        case "TRENDING_DESC":
            return "dl"
        case "START_DATE_DESC":
            return "dd"
        case "SCORE_DESC":
        default:
            return "rating"
    }
}
