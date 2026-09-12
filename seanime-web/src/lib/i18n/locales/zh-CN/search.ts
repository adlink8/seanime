import type { Dictionary } from "../../types"

/**
 * search 词条表
 * 归属：探索发现（搜索）页、高级搜索筛选项、流派 / 季度 / 状态 / 格式 / 国家常量表
 */
export const searchDictionary = {
    "search.page.title": "探索发现", // src/app/(main)/search/page.tsx
    "search.page.custom_source": "自定义源", // src/app/(main)/search/page.tsx
    "search.filter.genre": "流派", // src/app/(main)/search/_components/advanced-search-options.tsx
    "search.filter.genre_all": "全部流派", // src/app/(main)/search/_components/advanced-search-options.tsx
    "search.filter.genre_empty": "未找到流派", // src/app/(main)/search/_components/advanced-search-options.tsx
    "search.filter.tags_empty": "未找到标签", // src/app/(main)/search/_components/advanced-search-options.tsx
    "search.filter.country": "国家/地区", // src/app/(main)/search/_components/advanced-search-options.tsx
    "search.filter.country_all": "全部国家/地区", // src/app/(main)/search/_components/advanced-search-options.tsx
    "search.filter.score_all": "全部评分", // src/app/(main)/search/_components/advanced-search-options.tsx
    "search.filter.title_placeholder": "搜索标题...", // src/app/(main)/search/_components/advanced-search-options.tsx
    "search.filter.tags": "标签", // src/app/(main)/search/_components/advanced-search-options.tsx
    "search.filter.tags_all": "全部标签", // src/app/(main)/search/_components/advanced-search-options.tsx
    "search.filter.format": "格式", // src/app/(main)/search/_components/advanced-search-options.tsx
    "search.filter.format_all": "全部格式", // src/app/(main)/search/_components/advanced-search-options.tsx
    "search.filter.season_all": "全部季度", // src/app/(main)/search/_components/advanced-search-options.tsx
    "search.filter.year": "年份", // src/app/(main)/search/_components/advanced-search-options.tsx
    "search.filter.year_any": "不限年份", // src/app/(main)/search/_components/advanced-search-options.tsx
    "search.filter.status": "状态", // src/app/(main)/search/_components/advanced-search-options.tsx
    "search.filter.status_all": "全部状态", // src/app/(main)/search/_components/advanced-search-options.tsx
    "search.filter.adult": "成人内容 (R18)", // src/app/(main)/search/_components/advanced-search-options.tsx
    "search.sort.trending_now": "当前热门", // src/app/(main)/search/_components/advanced-search-page-title.tsx
    "search.sort.most_loved": "最受喜爱", // src/app/(main)/search/_components/advanced-search-page-title.tsx
    "search.sort.latest_released": "最新推出", // src/app/(main)/search/_components/advanced-search-page-title.tsx
    "search.sort.most_episodes": "剧集最多", // src/app/(main)/search/_components/advanced-search-page-title.tsx
    "search.sort.most_chapters": "章节最多", // src/app/(main)/search/_components/advanced-search-page-title.tsx
    "search.title.most_loved_manga": "最受喜爱漫画", // src/app/(main)/search/_components/advanced-search-page-title.tsx
    "search.list.load_more": "加载更多", // src/app/(main)/search/_components/advanced-search-list.tsx
    "search.type.anime": "动漫", // src/app/(main)/search/_lib/advanced-search-constants.ts
    "search.type.manga": "漫画", // src/app/(main)/search/_lib/advanced-search-constants.ts
    "search.genre.action": "动作", // src/app/(main)/search/_lib/advanced-search-constants.ts
    "search.genre.adventure": "冒险", // src/app/(main)/search/_lib/advanced-search-constants.ts
    "search.genre.comedy": "喜剧", // src/app/(main)/search/_lib/advanced-search-constants.ts
    "search.genre.drama": "剧情", // src/app/(main)/search/_lib/advanced-search-constants.ts
    "search.genre.fantasy": "奇幻", // src/app/(main)/search/_lib/advanced-search-constants.ts
    "search.genre.mystery": "悬疑", // src/app/(main)/search/_lib/advanced-search-constants.ts
    "search.genre.romance": "恋爱", // src/app/(main)/search/_lib/advanced-search-constants.ts
    "search.genre.sci_fi": "科幻", // src/app/(main)/search/_lib/advanced-search-constants.ts
    "search.genre.slice_of_life": "日常", // src/app/(main)/search/_lib/advanced-search-constants.ts
    "search.season.winter": "冬季", // src/app/(main)/search/_lib/advanced-search-constants.ts
    "search.season.spring": "春季", // src/app/(main)/search/_lib/advanced-search-constants.ts
    "search.season.summer": "夏季", // src/app/(main)/search/_lib/advanced-search-constants.ts
    "search.season.fall": "秋季", // src/app/(main)/search/_lib/advanced-search-constants.ts
    "search.status.finished": "已完结", // src/app/(main)/search/_lib/advanced-search-constants.ts
    "search.status.releasing": "连载中", // src/app/(main)/search/_lib/advanced-search-constants.ts
    "search.status.not_yet_released": "即将上线", // src/app/(main)/search/_lib/advanced-search-constants.ts
    "search.format.tv": "TV 动画", // src/app/(main)/search/_lib/advanced-search-constants.ts
    "search.format.movie": "剧场版 (Movie)", // src/app/(main)/search/_lib/advanced-search-constants.ts
    "search.country.jp": "日本", // src/app/(main)/search/_lib/advanced-search-constants.ts
    "search.country.kr": "韩国", // src/app/(main)/search/_lib/advanced-search-constants.ts
    "search.country.cn": "中国大陆", // src/app/(main)/search/_lib/advanced-search-constants.ts
    "search.country.tw": "中国台湾", // src/app/(main)/search/_lib/advanced-search-constants.ts
    "search.title.most_loved_anime": "最受喜爱动漫", // src/app/(main)/search/_components/advanced-search-page-title.tsx
    "search.sort.highest_score": "最高评分", // src/app/(main)/search/_components/advanced-search-page-title.tsx
} satisfies Dictionary
