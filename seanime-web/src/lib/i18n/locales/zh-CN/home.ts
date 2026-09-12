import type { Dictionary } from "../../types"

/**
 * home 词条表
 * 归属：首页（`app/(main)/_features/home/`）
 * 独占：05-02
 *
 * 复用说明（以下 key 为**只读引用**，本 plan 未在其它模块新增或改写任何词条）：
 * - `common.home.*` / `common.action.*`：05-01 已按 home 文件来源预先登记（见 common.ts 行尾注释）
 * - `common.state.*`：AniList 片单状态名（正在观看 / 二刷重温 / …）的唯一权威词条
 * - `library.continue_watching.title` / `library.filter.*` / `library.stats.*`：媒体库既有词条
 * - `search.filter.genre` / `search.season.*` / `search.status.*` / `search.type.*`：高级搜索常量表既有词条
 * - `navigation.item.home` / `navigation.item.lists`：导航项既有词条
 * 依据 05-CONTRACT §6.3「同名文案跨文件复用时必须复用同一 key，不新建重复词条」。
 */
export const homeDictionary = {
    // —— 首页项目：名称与描述（HOME_ITEMS[x].name / .description）——
    "home.item.continue_watching.description": "展示你当前正在追看的剧集列表。", // src/app/(main)/_features/home/home-items.utils.ts
    "home.item.continue_watching_header.name": "继续观看横幅", // src/app/(main)/_features/home/home-items.utils.ts
    "home.item.continue_watching_header.description": "在顶部展示正在追看的番剧轮播大图横幅。", // src/app/(main)/_features/home/home-items.utils.ts
    "home.item.anime_library.name": "本地番剧媒体库", // src/app/(main)/_features/home/home-items.utils.ts, src/app/(main)/_features/home/home-toolbar.tsx
    "home.item.anime_library.description": "按观看状态分类展示你下载并追看的番剧。", // src/app/(main)/_features/home/home-items.utils.ts
    "home.item.centered_title.name": "居中标题", // src/app/(main)/_features/home/home-items.utils.ts
    "home.item.centered_title.description": "显示一段居中的标题文字。", // src/app/(main)/_features/home/home-items.utils.ts
    "home.item.my_lists.description": "按观看状态展示你片单中的媒体。", // src/app/(main)/_features/home/home-items.utils.ts
    "home.item.local_anime_library.name": "本地番剧库", // src/app/(main)/_features/home/home-items.utils.ts
    "home.item.local_anime_library.description": "以完整网格展示你本地媒体库中的番剧。", // src/app/(main)/_features/home/home-items.utils.ts
    "home.item.library_upcoming_episodes.name": "媒体库即将播出剧集", // src/app/(main)/_features/home/home-items.utils.ts
    "home.item.library_upcoming_episodes.description": "以轮播展示你媒体库中番剧的即将播出剧集。", // src/app/(main)/_features/home/home-items.utils.ts
    "home.item.aired_recently.name": "近期已播出（全球）", // src/app/(main)/_features/home/home-items.utils.ts
    "home.item.aired_recently.description": "以轮播展示近期已播出的番剧剧集。", // src/app/(main)/_features/home/home-items.utils.ts
    "home.item.missed_sequels.name": "错过的续作", // src/app/(main)/_features/home/home-items.utils.ts, src/app/(main)/_features/home/home-screen.tsx
    "home.item.missed_sequels.description": "以轮播展示不在你收藏中的续作。", // src/app/(main)/_features/home/home-items.utils.ts
    "home.item.anime_schedule_calendar.name": "番剧放送日历", // src/app/(main)/_features/home/home-items.utils.ts
    "home.item.anime_schedule_calendar.description": "根据放送日程以日历展示番剧剧集。", // src/app/(main)/_features/home/home-items.utils.ts
    "home.item.local_anime_library_stats.name": "本地番剧库统计", // src/app/(main)/_features/home/home-items.utils.ts
    "home.item.local_anime_library_stats.description": "展示本地番剧库的统计数据。", // src/app/(main)/_features/home/home-items.utils.ts
    "home.item.discover_header.name": "探索发现页头图", // src/app/(main)/_features/home/home-items.utils.ts
    "home.item.discover_header.description": "展示带热门番剧轮播的页头。", // src/app/(main)/_features/home/home-items.utils.ts
    "home.item.anime_carousel.name": "番剧轮播", // src/app/(main)/_features/home/home-items.utils.ts, src/app/(main)/_features/home/home-screen.tsx
    "home.item.anime_carousel.description": "根据所选条件以轮播展示番剧。", // src/app/(main)/_features/home/home-items.utils.ts
    "home.item.manga_carousel.name": "漫画轮播", // src/app/(main)/_features/home/home-items.utils.ts, src/app/(main)/_features/home/home-screen.tsx
    "home.item.manga_carousel.description": "根据所选条件以轮播展示漫画。", // src/app/(main)/_features/home/home-items.utils.ts
    "home.item.manga_library.name": "漫画库", // src/app/(main)/_features/home/home-items.utils.ts
    "home.item.manga_library.description": "按观看状态展示你媒体库中的漫画。", // src/app/(main)/_features/home/home-items.utils.ts

    // —— 首页项目：选项字段标签（_carouselOptions / 各 item 的 options）——
    "home.option.name": "名称", // src/app/(main)/_features/home/home-items.utils.ts
    "home.option.sorting": "排序", // src/app/(main)/_features/home/home-items.utils.ts
    "home.option.sort.popular": "热门", // src/app/(main)/_features/home/home-items.utils.ts
    "home.option.sort.trending": "流行趋势", // src/app/(main)/_features/home/home-items.utils.ts
    "home.option.sort.title_romaji_asc": "罗马字标题 (A-Z)", // src/app/(main)/_features/home/home-items.utils.ts
    "home.option.sort.title_romaji_desc": "罗马字标题 (Z-A)", // src/app/(main)/_features/home/home-items.utils.ts
    "home.option.sort.title_english_asc": "英文标题 (A-Z)", // src/app/(main)/_features/home/home-items.utils.ts
    "home.option.sort.title_english_desc": "英文标题 (Z-A)", // src/app/(main)/_features/home/home-items.utils.ts
    "home.option.sort.score_asc": "评分 (0-10)", // src/app/(main)/_features/home/home-items.utils.ts
    "home.option.sort.score_desc": "评分 (10-0)", // src/app/(main)/_features/home/home-items.utils.ts
    "home.option.text": "文本", // src/app/(main)/_features/home/home-items.utils.ts
    "home.option.country_of_origin": "原产国", // src/app/(main)/_features/home/home-items.utils.ts
    "home.option.layout": "布局", // src/app/(main)/_features/home/home-items.utils.ts
    "home.option.layout.grid": "网格", // src/app/(main)/_features/home/home-items.utils.ts
    "home.option.layout.carousel": "轮播", // src/app/(main)/_features/home/home-items.utils.ts
    "home.option.watch_status": "观看状态", // src/app/(main)/_features/home/home-items.utils.ts
    "home.option.type": "类型", // src/app/(main)/_features/home/home-items.utils.ts
    "home.option.type.global": "全球", // src/app/(main)/_features/home/home-items.utils.ts
    "home.option.custom_list_name": "自定义清单名称（可选）", // src/app/(main)/_features/home/home-items.utils.ts
    "home.option.format.one_shot": "短篇", // src/app/(main)/_features/home/home-items.utils.ts

    // —— 首页空状态 / 错误 / 轮播提示 ——
    "home.empty.add_from_discover": "去“探索发现”添加番剧到片单", // src/app/(main)/_features/home/home-screen.tsx
    "home.empty.add_watching_to_library": "添加正在观看的番剧到媒体库", // src/app/(main)/_features/home/home-screen.tsx
    "home.empty.no_watching": "当前暂无正在追看的番剧", // src/app/(main)/_features/home/home-screen.tsx
    "home.empty.add_to_watching_hint": "将番剧添加到“正在观看”列表中即可开始", // src/app/(main)/_features/home/home-screen.tsx
    "home.error.item_not_found": "未找到该条目", // src/app/(main)/_features/home/home-screen.tsx
    "home.error.item_not_found_with_type": "未找到该条目（{type}）", // src/app/(main)/_features/home/home-screen.tsx
    "home.carousel.nothing_fetched": "未获取到内容，请更新你的筛选条件。", // src/app/(main)/_features/home/home-screen.tsx
    "home.carousel.invalid_item": "条目“{name}”缺少必要的配置项，无法展示。", // src/app/(main)/_features/home/home-screen.tsx

    // —— 首页设置弹窗 ——
    "home.settings.title": "首页设置", // home-settings-button.tsx（工具栏按钮 tooltip）
    "home.settings.anime_library": "番剧库", // src/app/(main)/_features/home/home-settings-modal.tsx
    "home.settings.anime_library_hint_local": "仅显示你本地番剧库中的番剧", // src/app/(main)/_features/home/home-settings-modal.tsx
    "home.settings.anime_library_hint_stream": "当前正在观看列表中的所有番剧都会纳入媒体库", // src/app/(main)/_features/home/home-settings-modal.tsx
    "home.settings.local_anime_only": "仅本地番剧", // src/app/(main)/_features/home/home-settings-modal.tsx
    "home.settings.local_anime_streaming": "本地番剧 + 流媒体", // src/app/(main)/_features/home/home-settings-modal.tsx
    "home.settings.home_layout": "首页布局", // src/app/(main)/_features/home/home-settings-modal.tsx
    "home.settings.no_items": "尚未添加任何项目。在下方添加项目即可自定义首页。", // src/app/(main)/_features/home/home-settings-modal.tsx
    "home.settings.available_items": "可用项目", // src/app/(main)/_features/home/home-settings-modal.tsx
    "home.settings.all_items_added": "所有可用项目都已添加到你的首页。", // src/app/(main)/_features/home/home-settings-modal.tsx
    "home.settings.add": "添加", // src/app/(main)/_features/home/home-settings-modal.tsx
    "home.settings.configure_item": "配置 {name}", // src/app/(main)/_features/home/home-settings-modal.tsx
    "home.settings.customize_hint": "自定义该首页项目的设置。", // src/app/(main)/_features/home/home-settings-modal.tsx
    "home.settings.cancel": "取消", // src/app/(main)/_features/home/home-settings-modal.tsx
    "home.settings.save": "保存", // src/app/(main)/_features/home/home-settings-modal.tsx
    "home.settings.placeholder_enter": "输入{label}", // src/app/(main)/_features/home/home-settings-modal.tsx
    "home.settings.placeholder_select": "选择{label}", // src/app/(main)/_features/home/home-settings-modal.tsx
    "home.settings.unsupported_field_type": "不支持的字段类型：{type}", // src/app/(main)/_features/home/home-settings-modal.tsx
    "home.settings.kind_row": "行", // src/app/(main)/_features/home/home-settings-modal.tsx
    "home.settings.kind_header": "头部", // src/app/(main)/_features/home/home-settings-modal.tsx

    // —— 首页工具栏 ——
    "home.toolbar.library_explorer": "媒体库资源管理器", // src/app/(main)/_features/home/home-toolbar.tsx
    "home.toolbar.missing_extension": "未安装种子源扩展。", // src/app/(main)/_features/home/home-toolbar.tsx
    "home.toolbar.search_in_library": "在媒体库中搜索", // src/app/(main)/_features/home/home-toolbar.tsx
    "home.toolbar.open_directory": "打开本地目录", // src/app/(main)/_features/home/home-toolbar.tsx
    "home.toolbar.ignored_files": "已忽略文件", // src/app/(main)/_features/home/home-toolbar.tsx
    "home.toolbar.scan_summaries": "扫描报告", // src/app/(main)/_features/home/home-toolbar.tsx
    "home.toolbar.unmatched_files": "处理未匹配文件 ({count})", // src/app/(main)/_features/home/home-toolbar.tsx
    "home.toolbar.unknown_media": "处理隐藏媒体 ({count})", // src/app/(main)/_features/home/home-toolbar.tsx
    "home.toolbar.nakama_library": "{name} 的媒体库", // src/app/(main)/_features/home/home-toolbar.tsx

    // —— 首页提示（toast）——
    "home.toast.item_added": "已添加首页项目", // src/app/(main)/_features/home/home-settings-modal.tsx
    "home.toast.item_removed": "已移除首页项目", // src/app/(main)/_features/home/home-settings-modal.tsx
    "home.toast.layout_updated": "首页布局已更新", // src/app/(main)/_features/home/home-settings-modal.tsx
} satisfies Dictionary
