import type { Dictionary } from "../../types"

/**
 * library 词条表
 * 归属：媒体库详情页（统计卡、筛选器）与媒体库排序选项
 */
export const libraryDictionary = {
    "library.stats.size": "媒体库大小", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.stats.file_count": "文件数", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.stats.entry_count": "条目数", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.stats.episodes": "剧集", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.stats.movies": "剧场版", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.stats.specials": "特别篇", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.filter.format": "格式", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.filter.format_all": "全部格式", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.filter.status": "状态", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.filter.status_all": "全部状态", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.filter.tags": "标签", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.filter.tags_all": "全部标签", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.filter.tags_empty": "未找到标签", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.filter.season": "季度", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.filter.season_all": "全部季度", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.filter.year": "年份", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.filter.year_any": "不限年份", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.filter.adult": "成人内容 (R18)", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.filter.all": "全部", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.continue_watching.title": "继续观看", // src/app/(main)/_features/anime-library/_containers/continue-watching.tsx
    "library.sort.airdate_desc": "最近播出", // src/lib/helpers/filtering.ts
    "library.sort.airdate": "最早播出", // src/lib/helpers/filtering.ts
    "library.sort.episodes_desc": "最高集数", // src/lib/helpers/filtering.ts
    "library.sort.unwatched_desc": "未看集数最多", // src/lib/helpers/filtering.ts
    "library.sort.score_desc": "评分最高", // src/lib/helpers/filtering.ts
    "library.sort.title": "标题 (A-Z)", // src/lib/helpers/filtering.ts
    "library.sort.last_watched_desc": "最近观看", // src/lib/helpers/filtering.ts
    "library.sort.unread_chapters_desc": "未读章节最多", // src/lib/helpers/filtering.ts
} satisfies Dictionary
