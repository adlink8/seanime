import type { Dictionary } from "../../types"

/**
 * navigation 词条表
 * 归属：顶部菜单、主侧边栏、离线侧边栏的导航项与导航相关操作
 *
 * 导航项收敛：`main-sidebar.tsx` 与 `top-menu.tsx` 的重复导航项
 * 统一登记为 `navigation.item.*`，两处共用同一批 key。
 */
export const navigationDictionary = {
    "navigation.top_menu.home": "首页", // src/app/(main)/_features/navigation/top-menu.tsx
    "navigation.top_menu.schedule": "放送日历", // src/app/(main)/_features/navigation/top-menu.tsx
    "navigation.top_menu.manga": "漫画", // src/app/(main)/_features/navigation/top-menu.tsx
    "navigation.top_menu.lists": "我的片单", // src/app/(main)/_features/navigation/top-menu.tsx
    "navigation.top_menu.discover": "探索发现", // src/app/(main)/_features/navigation/top-menu.tsx
    "navigation.sidebar.home": "首页", // src/app/(main)/_features/navigation/main-sidebar.tsx
    "navigation.sidebar.search": "搜索", // src/app/(main)/_features/navigation/main-sidebar.tsx
    "navigation.sidebar.torrent_list": "种子下载列表", // src/app/(main)/_features/navigation/main-sidebar.tsx
    "navigation.sidebar.seeding": "做种中", // src/app/(main)/_features/navigation/main-sidebar.tsx
    "navigation.sidebar.debrid": "Debrid 云端", // src/app/(main)/_features/navigation/main-sidebar.tsx
    "navigation.sidebar.scan_summaries": "本地扫描报告", // src/app/(main)/_features/navigation/main-sidebar.tsx
    "navigation.sidebar.auto_downloader": "自动追番下载", // src/app/(main)/_features/navigation/main-sidebar.tsx
    "navigation.sidebar.more": "更多", // src/app/(main)/_features/navigation/main-sidebar.tsx
    "navigation.sidebar.refresh_data": "刷新追番数据", // src/app/(main)/_features/navigation/main-sidebar.tsx
    "navigation.sidebar.nakama": "同好联动 (Nakama)", // src/app/(main)/_features/navigation/main-sidebar.tsx
    "navigation.sidebar.extensions": "扩展中心", // src/app/(main)/_features/navigation/main-sidebar.tsx
    "navigation.sidebar.offline_sync": "离线同步", // src/app/(main)/_features/navigation/main-sidebar.tsx
    "navigation.sidebar.settings": "系统设置", // src/app/(main)/_features/navigation/main-sidebar.tsx
    "navigation.action.login": "登录账号", // src/app/(main)/_features/navigation/main-sidebar.tsx
    "navigation.action.logout": "退出登录", // src/app/(main)/_features/navigation/main-sidebar.tsx
    "navigation.action.logout_confirm": "确定要退出登录吗？", // src/app/(main)/_features/navigation/main-sidebar.tsx

    // —— 导航项：main-sidebar 与 top-menu 共用（收敛后的唯一 key）——
    "navigation.item.home": "首页", // src/app/(main)/_features/navigation/main-sidebar.tsx, src/app/(main)/_features/navigation/top-menu.tsx
    "navigation.item.schedule": "放送日历", // src/app/(main)/_features/navigation/main-sidebar.tsx, src/app/(main)/_features/navigation/top-menu.tsx
    "navigation.item.manga": "漫画", // src/app/(main)/_features/navigation/main-sidebar.tsx, src/app/(main)/_features/navigation/top-menu.tsx, src/app/(main)/_features/navigation/offline-sidebar.tsx
    "navigation.item.lists": "我的片单", // src/app/(main)/_features/navigation/main-sidebar.tsx, src/app/(main)/_features/navigation/top-menu.tsx
    "navigation.item.discover": "探索发现", // src/app/(main)/_features/navigation/main-sidebar.tsx, src/app/(main)/_features/navigation/top-menu.tsx
    "navigation.item.search": "搜索", // src/app/(main)/_features/navigation/main-sidebar.tsx

    // —— 主侧栏专属导航项 ——
    "navigation.sidebar.torrent_list_with_seeding": "种子下载列表 ({count} 做种中)", // src/app/(main)/_features/navigation/main-sidebar.tsx

    // —— 主侧栏登录弹窗 / 退出确认 ——
    "navigation.main_sidebar.login_title": "使用 AniList 登录", // src/app/(main)/_features/navigation/main-sidebar.tsx
    "navigation.main_sidebar.login_description": "推荐使用 AniList 账号登录。", // src/app/(main)/_features/navigation/main-sidebar.tsx
    "navigation.main_sidebar.get_token": "获取 AniList Token", // src/app/(main)/_features/navigation/main-sidebar.tsx
    "navigation.main_sidebar.continue": "继续", // src/app/(main)/_features/navigation/main-sidebar.tsx
    "navigation.main_sidebar.enter_token": "输入 Token", // src/app/(main)/_features/navigation/main-sidebar.tsx

    // —— 离线侧栏 ——
    "navigation.offline.anime_library": "离线番剧库", // src/app/(main)/_features/navigation/offline-sidebar.tsx
    "navigation.offline.disable": "关闭离线模式", // src/app/(main)/_features/navigation/offline-sidebar.tsx
    "navigation.offline.disable_confirm": "确定要关闭离线模式吗？", // src/app/(main)/_features/navigation/offline-sidebar.tsx
    "navigation.offline.label": "离线模式", // src/app/(main)/_features/navigation/offline-sidebar.tsx
} satisfies Dictionary
