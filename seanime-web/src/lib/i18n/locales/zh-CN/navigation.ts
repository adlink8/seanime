import type { Dictionary } from "../../types"

/**
 * navigation 词条表
 * 归属：顶部菜单、主侧边栏、离线侧边栏的导航项与导航相关操作
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
} satisfies Dictionary
