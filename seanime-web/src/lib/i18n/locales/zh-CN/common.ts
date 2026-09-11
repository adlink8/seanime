import type { Dictionary } from "../../types"

/**
 * common 词条表
 * 归属：跨界面通用状态名、空状态、通用操作按钮（首页 / Debrid / 服务端工具函数）
 */
export const commonDictionary = {
    "common.state.watching": "正在观看", // src/lib/server/utils.ts
    "common.state.completed": "已看完", // src/lib/server/utils.ts
    "common.state.planning": "计划观看", // src/lib/server/utils.ts
    "common.state.paused": "暂停搁置", // src/lib/server/utils.ts
    "common.state.dropped": "已弃番", // src/lib/server/utils.ts
    "common.state.rewatching": "二刷重温", // src/lib/server/utils.ts
    "common.state.reading": "正在阅读", // src/lib/server/utils.ts
    "common.state.plan_to_read": "计划阅读", // src/lib/server/utils.ts
    "common.state.manga_dropped": "已弃坑", // src/lib/server/utils.ts
    "common.state.rereading": "重温阅读", // src/lib/server/utils.ts
    "common.empty.title": "暂无内容", // src/app/(main)/debrid/page.tsx
    "common.empty.debrid_torrents": "暂无活跃的种子任务", // src/app/(main)/debrid/page.tsx
    "common.action.refresh": "刷新", // src/app/(main)/_features/home/home-toolbar.tsx
    "common.action.scan": "扫描", // src/app/(main)/_features/home/home-toolbar.tsx
    "common.action.refresh_library": "刷新本地番剧", // src/app/(main)/_features/home/home-toolbar.tsx
    "common.action.scan_library": "扫描本地番剧", // src/app/(main)/_features/home/home-toolbar.tsx
    "common.action.install_extensions": "安装所需扩展", // src/app/(main)/_features/home/home-toolbar.tsx
    "common.action.playlists": "播放列表", // src/app/(main)/_features/home/home-toolbar.tsx
    "common.action.batch_operations": "批量操作", // src/app/(main)/_features/home/home-toolbar.tsx
    "common.home.empty_title": "你的首页暂无番剧内容", // src/app/(main)/_features/home/home-screen.tsx
    "common.home.scan_library": "扫描本地番剧媒体库", // src/app/(main)/_features/home/home-screen.tsx
    "common.home.trending": "当前热门番剧", // src/app/(main)/_features/home/home-screen.tsx
} satisfies Dictionary
