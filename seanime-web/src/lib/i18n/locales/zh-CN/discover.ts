import type { Dictionary } from "../../types"

/**
 * discover 词条表
 * 归属：探索发现页（`app/(main)/discover/`）
 * 独占：06-02
 *
 * 复用说明（以下 key 为**只读引用**，本 plan 未在其它模块新增或改写任何词条）：
 * - `common.home.trending`：页面板块标题「当前热门番剧」与首页同名文案复用
 * - `navigation.item.schedule`：页签「放送日历」与导航项同名文案复用
 * - `media.action.preview`：右键菜单 / 按钮「预览」
 * - `media.alt.banner_image` / `media.alt.cover_image`：图片 alt 描述
 * - `library.filter.all`：流派选择器「全部」
 * 依据 06-CONTRACT §2「同名文案必须复用既有 key，禁止新建重复词条」。
 */
export const discoverDictionary = {
    // —— 探索发现页：顶部页签（discover/page.tsx）——
    "discover.tab.anime": "动漫", // src/app/(main)/discover/page.tsx
    "discover.tab.manga": "漫画", // src/app/(main)/discover/page.tsx
    "discover.tab.novel": "轻小说", // src/app/(main)/discover/page.tsx（Phase 3.1 顶层分类）
    "discover.tab.asmr": "音声", // src/app/(main)/discover/page.tsx（Phase 3.1 顶层分类）

    // —— 探索发现页：板块标题（discover/page.tsx）——
    "discover.section.this_season": "本季霸权精选", // src/app/(main)/discover/page.tsx
    "discover.section.past_season": "上季高分佳作", // src/app/(main)/discover/page.tsx
    "discover.section.upcoming": "即将播出新番", // src/app/(main)/discover/page.tsx
    "discover.section.trending_movies": "热门剧场版", // src/app/(main)/discover/page.tsx
    "discover.section.manga_trending_jp": "热门日漫", // src/app/(main)/discover/page.tsx
    "discover.section.manga_trending_kr": "热门韩漫", // src/app/(main)/discover/page.tsx
    "discover.section.manga_trending_cn": "热门国漫", // src/app/(main)/discover/page.tsx

    // —— 探索发现页：轻小说 + 音声区块（Phase 2.5，discover/page.tsx）——
    "discover.section.popular_novel": "热门轻小说", // src/app/(main)/discover/page.tsx
    "discover.section.this_season_novel": "本季轻小说", // src/app/(main)/discover/page.tsx
    "discover.section.popular_asmr": "热门音声", // src/app/(main)/discover/page.tsx
    "discover.section.latest_asmr": "最新音声", // src/app/(main)/discover/page.tsx

    // —— 探索发现页：每日推荐区块（Phase 3.9c，discover/page.tsx）——
    "discover.section.daily": "每日推荐", // src/app/(main)/discover/page.tsx（当前 tab 域，空数据隐藏）

    // —— 探索发现页：页头（_components/discover-page-header.tsx）——
    "discover.header.releasing_now": "正在放送", // src/app/(main)/discover/_components/discover-page-header.tsx
    "discover.header.episodes_released": "{count} 集已播出", // src/app/(main)/discover/_components/discover-page-header.tsx
    "discover.header.total_episodes": "共 {count} 集", // src/app/(main)/discover/_components/discover-page-header.tsx
    "discover.header.go_to_slide": "转到第 {index} 张", // src/app/(main)/discover/_components/discover-page-header.tsx

    // —— 探索发现页：放送时间表（_containers/discover-airing-schedule.tsx）——
    "discover.airing.title": "放送时间表", // src/app/(main)/discover/_containers/discover-airing-schedule.tsx
    "discover.airing.today": "今天", // src/app/(main)/discover/_containers/discover-airing-schedule.tsx
    "discover.airing.open_page": "打开页面", // src/app/(main)/discover/_containers/discover-airing-schedule.tsx
    "discover.airing.episode_airing_at": "第 {episode} 集将于 {time} 播出", // src/app/(main)/discover/_containers/discover-airing-schedule.tsx

    // —— 探索发现页：错过续作（_containers/discover-missed-sequels.tsx）——
    "discover.missed_sequels.title": "你可能错过的续作", // src/app/(main)/discover/_containers/discover-missed-sequels.tsx

    // —— 探索发现页：漫画搜索（_containers/discover-trending-manga-all.tsx）——
    "discover.manga.search_placeholder": "搜索漫画...", // src/app/(main)/discover/_containers/discover-trending-manga-all.tsx
} satisfies Dictionary
