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
    // —— 通知（api/hooks/anilist.hooks.ts）——
    "common.toast.anilist_synced": "AniList 追番数据已同步至最新", // api/hooks/anilist.hooks.ts
    "common.toast.entry_updated": "追番条目已更新", // api/hooks/anilist.hooks.ts
    "common.toast.entry_deleted": "条目已删除", // api/hooks/anilist.hooks.ts
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

    // ===== 09-05 收官（common）=====
    // —— 媒体排除选择器（components/shared/media-exclusion-selector.tsx）——
    "common.media_exclusion.title": "选择要从分享中排除的番剧", // media-exclusion-selector.tsx Modal title
    "common.media_exclusion.description": "选择你不想与其他客户端分享的番剧。被选中的番剧将对已连接的客户端不可见。", // media-exclusion-selector.tsx 说明文本
    "common.media_exclusion.excluded_count": "{count} 部番剧已排除分享", // media-exclusion-selector.tsx 状态行
    "common.media_exclusion.not_visible_hint": "（对其他客户端不可见）", // media-exclusion-selector.tsx 提示
    "common.media_exclusion.more": "还有 {count} 部", // media-exclusion-selector.tsx 溢出提示
    "common.media_exclusion.edit_selection": "编辑选择", // media-exclusion-selector.tsx 触发按钮
    "common.media_exclusion.select_anime": "选择番剧", // media-exclusion-selector.tsx 触发按钮
    "common.media_exclusion.select_all": "全选", // media-exclusion-selector.tsx
    "common.media_exclusion.deselect_all": "取消全选", // media-exclusion-selector.tsx
    "common.media_exclusion.select_adult": "选择成人内容", // media-exclusion-selector.tsx
    "common.media_exclusion.selected_count": "已选 {count} 部（不会分享）", // media-exclusion-selector.tsx
    "common.media_exclusion.section_all": "全部", // media-exclusion-selector.tsx 分区标题
    "common.media_exclusion.done": "完成（已选 {count} 部）", // media-exclusion-selector.tsx 底部按钮
    // —— 徽章（components/shared/beta-badge.tsx）——
    "common.badge.beta": "公测版", // beta-badge.tsx
    "common.badge.alpha": "内测版", // beta-badge.tsx
    "common.badge.experimental": "实验性功能", // beta-badge.tsx title 默认值
    // —— 通用确认弹窗（components/shared/confirmation-dialog.tsx）——
    "common.dialog.confirm_description": "确定要继续吗？", // confirmation-dialog.tsx description 默认值
    "common.action.confirm": "确认", // confirmation-dialog.tsx actionText 默认值
    "common.action.cancel": "取消", // confirmation-dialog.tsx cancelText 默认值
    // —— 404 页（components/shared/not-found.tsx）——
    "common.error.page_not_found": "页面未找到", // not-found.tsx title
    "common.error.page_not_found_desc": "你访问的页面不存在。", // not-found.tsx
    "common.action.go_home": "返回首页", // not-found.tsx
    // —— 文件 / 目录选择器（components/shared/directory-selector.tsx, file-tree-selector.tsx, file-selector.tsx, _file-selector.tsx）——
    "common.action.change_library": "更换媒体库", // directory-selector.tsx
    "common.file.select_directory": "选择目录", // directory-selector.tsx Modal title
    "common.file.select_file_or_directory": "选择文件或目录", // file-selector.tsx / _file-selector.tsx Modal title
    "common.file.likely_match": "疑似匹配", // file-tree-selector.tsx
    "common.file.search_files": "搜索文件...", // file-tree-selector.tsx placeholder
    // —— 无集数元数据提示（components/shared/slider-episode-item.tsx）——
    "common.empty.no_metadata": "未找到元数据", // slider-episode-item.tsx
    // —— 表单（components/ui/form/fields.tsx, schema-presets.ts）——
    "common.form.directory_label": "目录", // fields.tsx label
    "common.form.expected_file": "应为文件类型", // schema-presets.ts Zod 校验消息（表单错误展示）
    // —— 无障碍文案（components/ui/*）——
    "common.a11y.open_main_menu": "打开主菜单", // navigation-menu.tsx, app-sidebar.tsx sr-only
    "common.a11y.drawer": "抽屉", // drawer.tsx, vaul/index.tsx VisuallyHidden Title
    "common.a11y.dialog": "对话框", // modal.tsx VisuallyHidden Title
    // —— 通用操作（components/ui/date-picker/date-picker.tsx）——
    "common.action.clear": "清除", // date-picker.tsx
    "common.action.refresh_anilist": "刷新 AniList", // refresh-anilist-button.tsx Tooltip
    // —— 提示（_features/rate-limit-loader.tsx, api/hooks/status.hooks.ts）——
    "common.toast.anilist_rate_limit": "AniList 速率受限：{seconds} 秒后重试", // rate-limit-loader.tsx
    "common.toast.generating_goroutine_profile": "正在生成 goroutine 性能剖析...", // status.hooks.ts toast.info
} satisfies Dictionary
