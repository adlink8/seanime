import type { Dictionary } from "../../types"

/**
 * schedule 词条表
 * 归属：时间表页（`app/(main)/schedule/`）
 * 独占：06-02
 *
 * 复用说明（以下 key 为**只读引用**，本 plan 未在其它模块新增或改写任何词条）：
 * - `common.state.watching/planning/completed/paused`：日历设置里的片单状态选项，
 *   为 AniList 状态名的唯一权威词条
 * - `media.alt.banner_image`：日历背景图 alt 描述
 * - `library.stats.movies`：近期已播出的「Movie」格式角标
 * 依据 06-CONTRACT §2「同名文案必须复用既有 key，禁止新建重复词条」。
 */
export const scheduleDictionary = {
    // —— 时间表页：标题与副标题（schedule/page.tsx、_containers/upcoming-episodes.tsx）——
    "schedule.title": "放送日历 / 追番排期", // src/app/(main)/schedule/page.tsx
    "schedule.personal_list_hint": "基于你的个人片单", // src/app/(main)/schedule/page.tsx, src/app/(main)/schedule/_containers/upcoming-episodes.tsx
    "schedule.upcoming.title": "即将更新剧集", // src/app/(main)/schedule/_containers/upcoming-episodes.tsx
    "schedule.missing.title": "媒体库缺失的剧集", // src/app/(main)/schedule/_components/missing-episodes.tsx
    "schedule.missing.silenced": "已静音的剧集", // src/app/(main)/schedule/_components/missing-episodes.tsx
    "schedule.recent.title": "近期已播出", // src/app/(main)/schedule/_containers/recent-releases.tsx

    // —— 时间表页：剧集标题（_containers/upcoming-episodes.tsx、_components/schedule-calendar.tsx）——
    "schedule.episode_title": "第 {ep} 集", // src/app/(main)/schedule/_containers/upcoming-episodes.tsx, src/app/(main)/schedule/_components/schedule-calendar.tsx

    // —— 时间表页：放送日历设置（_components/schedule-calendar.tsx）——
    "schedule.calendar.week_start": "一周起始日", // src/app/(main)/schedule/_components/schedule-calendar.tsx
    "schedule.calendar.monday": "周一", // src/app/(main)/schedule/_components/schedule-calendar.tsx
    "schedule.calendar.sunday": "周日", // src/app/(main)/schedule/_components/schedule-calendar.tsx
    "schedule.calendar.status": "状态", // src/app/(main)/schedule/_components/schedule-calendar.tsx
    "schedule.calendar.indicate_watched": "标记已观看的剧集", // src/app/(main)/schedule/_components/schedule-calendar.tsx
    "schedule.calendar.disable_transitions": "关闭图片切换动画", // src/app/(main)/schedule/_components/schedule-calendar.tsx

    // —— 时间表页：放送日历空状态与计数（_components/schedule-calendar.tsx）——
    "schedule.calendar.empty_month": "本月暂无排播剧集", // src/app/(main)/schedule/_components/schedule-calendar.tsx
    "schedule.calendar.empty_today": "今天暂无排播剧集", // src/app/(main)/schedule/_components/schedule-calendar.tsx
    "schedule.calendar.empty_day": "这一天暂无排播剧集。", // src/app/(main)/schedule/_components/schedule-calendar.tsx
    "schedule.calendar.finale": "最终集", // src/app/(main)/schedule/_components/schedule-calendar.tsx
    "schedule.calendar.scheduled_count": "{count} 个排播剧集", // src/app/(main)/schedule/_components/schedule-calendar.tsx
    "schedule.calendar.episode_count": "{count} 集", // src/app/(main)/schedule/_components/schedule-calendar.tsx
    "schedule.calendar.more": "还有 {count} 个", // src/app/(main)/schedule/_components/schedule-calendar.tsx
} satisfies Dictionary
