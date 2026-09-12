import type { Dictionary } from "../../types"

/**
 * extensions 词条表
 * 归属：扩展（extensions/）
 * 独占：09-03
 *
 * 只读复用：common.* / library.* 等同名文案按契约复用，不新建重复词条。
 */
export const extensionsDictionary = {
    // ===== 扩展中心页（page.tsx）=====
    "extensions.page.tab_installed": "已安装插件", // extensions/page.tsx
    "extensions.page.tab_marketplace": "插件市场", // extensions/page.tsx

    // ===== 扩展列表（extension-list.tsx）=====
    "extensions.list.title": "扩展中心", // extension-list.tsx
    "extensions.list.subtitle": "管理已安装的插件与第三方内容提供商。", // extension-list.tsx
    "extensions.list.fetch_error": "无法获取扩展。", // extension-list.tsx
    "extensions.list.update_all": "全部更新", // extension-list.tsx
    "extensions.list.toast_installing_updates": "正在安装更新...", // extension-list.tsx
    "extensions.list.toast_skip_unsafe": "已跳过「{id}」：该扩展使用了不安全标志，请手动更新。", // extension-list.tsx
    "extensions.list.check_updates": "检查更新", // extension-list.tsx（extension-card.tsx 复用）
    "extensions.list.add_extension": "添加扩展", // extension-list.tsx
    "extensions.list.menu_playground": "调试沙箱", // extension-list.tsx（playground/page 复用）
    "extensions.list.menu_marketplace": "插件市场", // extension-list.tsx
    "extensions.list.search_installed_placeholder": "搜索已安装扩展...", // extension-list.tsx
    "extensions.list.no_results": "没有符合搜索条件的扩展。", // extension-list.tsx
    "extensions.list.section_permissions_required": "需要授权权限", // extension-list.tsx（invalid-extension-card.tsx 弹窗标题复用）
    "extensions.list.section_invalid": "无效扩展", // extension-list.tsx
    "extensions.list.section_disabled": "已禁用", // extension-list.tsx
    "extensions.list.section_plugins": "核心插件", // extension-list.tsx
    "extensions.list.browse_all_sources": "浏览全部数据源", // extension-list.tsx

    // ===== 扩展类型展示名（extension-list.tsx EXTENSION_TYPE；marketplace 分组标题、playground 下拉复用）=====
    "extensions.type.plugin": "插件", // extension-list.tsx
    "extensions.type.anime_torrent_provider": "动漫种子提供商", // extension-list.tsx
    "extensions.type.manga_provider": "漫画提供商", // extension-list.tsx
    "extensions.type.onlinestream_provider": "在线流媒体提供商", // extension-list.tsx
    "extensions.type.custom_source": "自定义数据源", // extension-list.tsx

    // ===== 插件市场（marketplace-extensions.tsx）=====
    "extensions.marketplace.modal_title": "仓库 URL", // marketplace-extensions.tsx
    "extensions.marketplace.url_hint": "请输入仓库 JSON 文件的 URL。", // marketplace-extensions.tsx
    "extensions.marketplace.url_label": "市场 URL", // marketplace-extensions.tsx
    "extensions.marketplace.url_placeholder": "输入市场 URL", // marketplace-extensions.tsx
    "extensions.marketplace.url_invalid": "请输入有效的 URL", // marketplace-extensions.tsx
    "extensions.marketplace.apply_default": "恢复默认", // marketplace-extensions.tsx
    "extensions.marketplace.title": "插件市场", // marketplace-extensions.tsx
    "extensions.marketplace.desc": "浏览并安装来自仓库的扩展。", // marketplace-extensions.tsx
    "extensions.marketplace.source_label": "来源：", // marketplace-extensions.tsx
    "extensions.marketplace.official_repository": "官方仓库", // marketplace-extensions.tsx
    "extensions.marketplace.toast_refreshed": "已刷新", // marketplace-extensions.tsx
    "extensions.marketplace.toast_url_updated": "市场 URL 已更新", // marketplace-extensions.tsx
    "extensions.marketplace.toast_fetch_failed": "无法从所提供的 URL 获取扩展", // marketplace-extensions.tsx
    "extensions.marketplace.toast_reset_default": "已重置为默认市场 URL", // marketplace-extensions.tsx
    "extensions.marketplace.toast_fetch_default_failed": "无法从默认 URL 获取扩展", // marketplace-extensions.tsx
    "extensions.marketplace.change_repository": "更换仓库", // marketplace-extensions.tsx
    "extensions.marketplace.no_providers_title": "暂无内容提供商", // marketplace-extensions.tsx
    "extensions.marketplace.no_providers_desc": "Seanime 默认插件市场已不再收录内容提供商。请在网上寻找新的仓库地址并添加。", // marketplace-extensions.tsx
    "extensions.marketplace.add_new_repository": "添加新仓库", // marketplace-extensions.tsx
    "extensions.marketplace.tab_all_types": "全部类型", // marketplace-extensions.tsx
    "extensions.marketplace.tab_anime_torrents": "动漫种子", // marketplace-extensions.tsx
    "extensions.marketplace.tab_online_streaming": "在线流媒体", // marketplace-extensions.tsx
    "extensions.marketplace.tab_custom_sources": "自定义源", // marketplace-extensions.tsx
    "extensions.marketplace.anime_torrents_section": "动漫种子源", // marketplace-extensions.tsx（extension-list.tsx 复用）
    "extensions.marketplace.manga_section": "漫画源", // marketplace-extensions.tsx（extension-list.tsx 复用）
    "extensions.marketplace.online_streaming_section": "在线流媒体源", // marketplace-extensions.tsx（extension-list.tsx 复用）
    "extensions.marketplace.custom_sources_popover": "自定义数据源不提供任何流媒体功能，需要种子或在线流媒体提供商才能实现。", // marketplace-extensions.tsx
    "extensions.marketplace.custom_sources_desc": "自定义数据源让你可以浏览 AniList 之外的媒体。", // marketplace-extensions.tsx
    "extensions.marketplace.fetch_error": "无法获取市场扩展。", // marketplace-extensions.tsx
    "extensions.marketplace.no_results": "没有符合筛选条件的扩展。", // marketplace-extensions.tsx
    "extensions.marketplace.search_placeholder": "搜索扩展...", // marketplace-extensions.tsx
    "extensions.marketplace.all_languages": "全部语言", // marketplace-extensions.tsx
    "extensions.marketplace.missing_type_online_streaming": "在线流媒体", // marketplace-extensions.tsx
    "extensions.marketplace.missing_type_torrent_streaming": "种子流媒体", // marketplace-extensions.tsx
    "extensions.common.update_available": "有可用更新", // marketplace-extensions.tsx 徽章（extension-card.tsx toast 复用）

    // ===== 添加扩展弹窗（add-extension-modal.tsx）=====
    "extensions.add.install_via_url": "通过 URL 安装", // add-extension-modal.tsx
    "extensions.add.install_via_url_desc": "输入扩展 manifest.json 链接进行安装。", // add-extension-modal.tsx
    "extensions.add.find": "查找", // add-extension-modal.tsx
    "extensions.add.already_installed": "该扩展已经安装过。", // add-extension-modal.tsx
    "extensions.add.install": "安装", // add-extension-modal.tsx
    "extensions.add.repo_hint": "您也可以输入扩展源仓库地址，批量导入并一键安装多个扩展。", // add-extension-modal.tsx
    "extensions.add.import_from_repo": "从软件源批量导入", // add-extension-modal.tsx
    "extensions.add.import_from_repo_desc": "输入扩展仓库地址或 JSON 链接导入扩展。", // add-extension-modal.tsx
    "extensions.add.repo_placeholder": "https://example.com/extensions.json 或 { \"urls\": [...] }", // add-extension-modal.tsx（示例 URL 冻结）
    "extensions.add.import_all": "导入全部", // add-extension-modal.tsx
    "extensions.add.install_all": "全部安装", // add-extension-modal.tsx
    "extensions.add.toast_invalid_url": "请提供有效的 URL。", // add-extension-modal.tsx
    "extensions.add.toast_install_success": "扩展安装成功。", // add-extension-modal.tsx

    // ===== 扩展卡片（extension-card.tsx）=====
    "extensions.card.alt_icon": "扩展图标", // extension-card.tsx（marketplace/invalid-extension-card/extension-details 复用）
    "extensions.card.preferences": "偏好设置", // extension-card.tsx（extension-user-config.tsx 弹窗标题复用）
    "extensions.card.info": "详情", // extension-card.tsx
    "extensions.card.documentation": "文档", // extension-card.tsx（extension-details.tsx 复用）
    "extensions.card.code": "代码", // extension-card.tsx（extension-code.tsx 弹窗标题复用）
    "extensions.card.builtin": "内置", // extension-card.tsx（extension-details.tsx 复用）
    "extensions.card.disabled": "已禁用", // extension-card.tsx
    "extensions.card.disable": "禁用", // extension-card.tsx
    "extensions.card.uninstall": "卸载", // extension-card.tsx
    "extensions.card.remove_title": "移除 {name}", // extension-card.tsx
    "extensions.card.remove_undone": "此操作无法撤销。", // extension-card.tsx
    "extensions.card.update_available_prefix": "有可用更新：", // extension-card.tsx
    "extensions.card.view_updated_code": "查看更新后代码", // extension-card.tsx
    "extensions.card.install_update": "安装更新", // extension-card.tsx
    "extensions.card.toast_up_to_date": "扩展已是最新版本", // extension-card.tsx
    "extensions.card.toast_no_id": "扩展缺少 ID", // extension-card.tsx（invalid-extension-card.tsx 复用）

    // ===== 无效扩展卡片（invalid-extension-card.tsx）=====
    "extensions.invalid.error_details_title": "错误详情", // invalid-extension-card.tsx
    "extensions.invalid.load_failed": "Seanime 加载该扩展失败。如果不确定这意味着什么，请联系扩展作者。", // invalid-extension-card.tsx
    "extensions.invalid.code_label": "错误码：", // invalid-extension-card.tsx
    "extensions.invalid.manifest_error": "Manifest 错误", // invalid-extension-card.tsx
    "extensions.invalid.incompatible": "与当前版本的 Seanime 不兼容", // invalid-extension-card.tsx
    "extensions.invalid.invalid_code": "代码无效或不兼容", // invalid-extension-card.tsx
    "extensions.invalid.invalid_id": "无效 ID", // invalid-extension-card.tsx
    "extensions.invalid.grant": "授权", // invalid-extension-card.tsx
    "extensions.invalid.permission_request": "插件 {name} 正在请求以下权限：", // invalid-extension-card.tsx
    "extensions.invalid.unsafe_warning": "该插件依赖不安全标志才能运行，Seanime 无法保证其安全性。", // invalid-extension-card.tsx
    "extensions.invalid.view_code": "查看代码", // invalid-extension-card.tsx
    "extensions.invalid.grant_permissions": "授予权限", // invalid-extension-card.tsx

    // ===== 扩展详情（extension-details.tsx）=====
    "extensions.details.website": "网站", // extension-details.tsx
    "extensions.details.id_label": "ID：", // extension-details.tsx
    "extensions.details.author_label": "作者：", // extension-details.tsx
    "extensions.details.language_label": "语言：", // extension-details.tsx
    "extensions.details.manifest_url": "Manifest URL：", // extension-details.tsx
    "extensions.details.notes_label": "备注：", // extension-details.tsx

    // ===== 扩展代码弹窗（extension-code.tsx）=====
    "extensions.code.edit_hint": "你可以在此处编辑该扩展的代码。", // extension-code.tsx

    // ===== 扩展偏好设置弹窗（extension-user-config.tsx）=====
    "extensions.user_config.edit_hint": "你可以在此处编辑该扩展的偏好设置。", // extension-user-config.tsx
    "extensions.user_config.fill_required": "请填写必填的配置项。", // extension-user-config.tsx
    "extensions.user_config.default_value": "默认值：{value}", // extension-user-config.tsx
    "extensions.user_config.na": "暂无", // extension-user-config.tsx

    // ===== 调试沙箱（playground/_containers/extension-playground.tsx）=====
    "extensions.playground.oops": "哎呀！", // extension-playground.tsx
    "extensions.playground.screen_too_small": "屏幕尺寸太小，无法显示。", // extension-playground.tsx
    "extensions.playground.output": "输出", // extension-playground.tsx
    "extensions.playground.run": "运行", // extension-playground.tsx
    "extensions.playground.running": "运行中...", // extension-playground.tsx
    "extensions.playground.toast_invalid_function": "选择的函数无效。", // extension-playground.tsx
    "extensions.playground.toast_no_output": "没有可复制的输出", // extension-playground.tsx
    "extensions.playground.method": "方法", // extension-playground.tsx
    "extensions.playground.media_id": "媒体 ID", // extension-playground.tsx
    "extensions.playground.query": "查询关键词", // extension-playground.tsx
    "extensions.playground.episode_number": "剧集编号", // extension-playground.tsx
    "extensions.playground.resolution": "分辨率", // extension-playground.tsx
    "extensions.playground.torrent_json": "种子 JSON", // extension-playground.tsx
    "extensions.playground.manga_id": "漫画 ID", // extension-playground.tsx
    "extensions.playground.chapter_id": "章节 ID", // extension-playground.tsx
    "extensions.playground.episode_id": "剧集 ID", // extension-playground.tsx
    "extensions.playground.episode_json": "剧集 JSON", // extension-playground.tsx
    "extensions.playground.server": "服务器", // extension-playground.tsx
    "extensions.playground.batch": "批量", // extension-playground.tsx
    "extensions.playground.best_releases": "最佳发布", // extension-playground.tsx
    "extensions.playground.manga_match_hint": "Seanime 将根据漫画标题自动选择最佳匹配。", // extension-playground.tsx
    "extensions.playground.anime_match_hint": "Seanime 将根据番剧标题自动选择最佳匹配。", // extension-playground.tsx
    "extensions.playground.lang_typescript": "TypeScript", // extension-playground.tsx
    "extensions.playground.lang_javascript": "JavaScript", // extension-playground.tsx

    // ===== 09-05 收官（plugin）=====
    "extensions.plugin_tray.title": "插件托盘", // plugin-sidebar-tray.tsx Tooltip
    "extensions.plugin_tray.count": "插件托盘 ({n})", // plugin-sidebar-tray.tsx Tooltip（带角标数）
    "extensions.plugin_tray.pin": "固定", // plugin-sidebar-tray.tsx Tooltip
    "extensions.plugin_tray.unpin": "取消固定", // plugin-sidebar-tray.tsx Tooltip
    "extensions.plugin_tray.empty": "暂无托盘插件", // plugin-sidebar-tray.tsx 空态
    "extensions.plugin_tray.debug": "调试", // plugin-sidebar-tray.tsx 调试区标题
    "extensions.plugin_tray.dev_mode_notice": "这些扩展以开发模式加载。", // plugin-sidebar-tray.tsx
    "extensions.plugin_tray.debug_logs": "调试日志", // plugin-sidebar-tray.tsx Tooltip
    "extensions.plugin.unknown_dom_action": "未知的 DOM 操作：{action}", // dom-manager.ts console.warn
} satisfies Dictionary
