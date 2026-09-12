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

    // —— 加载层 ——
    "common.loading.alt": "正在加载...", // src/components/shared/loading-overlay-with-logo.tsx
    "common.loading.default": "正 在 加 载 中 . . .", // src/components/shared/loading-overlay-with-logo.tsx
    "common.loading.server_data": "服 务 加 载 中 . . .", // src/app/(main)/server-data-wrapper.tsx
    "common.action.reload": "重新加载", // src/components/shared/loading-overlay-with-logo.tsx

    // —— 错误边界 / 错误页 ——
    "common.error.client_side": "客户端错误", // src/app/(main)/error.tsx, src/components/shared/app-error-boundary.tsx
    "common.error.unexpected": "发生了意料之外的错误。", // src/app/(main)/error.tsx, src/components/shared/app-error-boundary.tsx
    "common.error.title": "出错了！", // src/components/shared/luffy-error.tsx
    "common.error.transcode_disabled": "转码功能尚未启用", // src/app/(main)/server-data-wrapper.tsx
    "common.action.retry": "重试", // src/components/shared/luffy-error.tsx, src/app/(main)/error.tsx

    // —— 服务数据包装器（服务端未就绪 / 未登录引导）——
    "common.status.updating": "Seanime 正在更新中。更新完成并恢复连接后请刷新页面。", // src/app/(main)/server-data-wrapper.tsx
    "common.welcome.title": "欢迎使用 Seanime！", // src/app/(main)/server-data-wrapper.tsx
    "common.action.get_anilist_token": "获取 AniList 授权 Token", // src/app/(main)/server-data-wrapper.tsx
    "common.action.login_anilist_account": "登录 AniList 账号", // src/app/(main)/server-data-wrapper.tsx
    "common.action.continue_and_login": "继续并登录", // src/app/(main)/server-data-wrapper.tsx
    "common.form.token_required": "必须输入 Token", // src/app/(main)/server-data-wrapper.tsx, src/app/(main)/_features/navigation/main-sidebar.tsx
    "common.form.token_label": "输入授权 Token", // src/app/(main)/server-data-wrapper.tsx

    // —— 通用确认 ——
    "common.action.yes": "是", // src/app/(main)/_features/navigation/offline-sidebar.tsx
} satisfies Dictionary
