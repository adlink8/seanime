import type { Dictionary } from "../../types"

/**
 * manga 词条表
 * 归属：漫画模块（app/(main)/manga/）
 * 独占：08-01
 *
 * 只读复用：common.* / media.* / library.* / entry.* 等同名文案按契约复用，不新建重复词条。
 * 已复用（勿重复新增）：
 *   - common.action.retry（Retry）/ common.error.client_side（Client side error）
 *   - common.action.refresh（Refresh）/ common.state.reading（Currently Reading）
 *   - settings.action.reset（Reset）/ library.common.cancel（Cancel）
 *   - media.action.confirm（Confirm）/ media.action.start_reading（Start reading）
 *   - library.filter.all（All）/ library.collection.show_all（Show all）
 *   - entry.meta.custom_source / entry.meta.open_in_website（meta-section.tsx）
 *   - library.unmatched.open_on_anilist（Open on AniList）/ entry.dropdown.copy_id（Copy ID）
 *   - entry.relations.header（Relations）/ entry.relations.recommendations（Recommendations）
 *   - entry.torrent_filter.languages（Language）
 * 注意：_lib/language-map.ts 的 name 字段是死数据（消费点只读 nativeName），冻结不译；
 * nativeName 保持语言本名原样（日本語/中文/繁體…）。
 */
export const mangaDictionary = {
    // —— 章节通用（chapter-list / drawer / reader-bar 等）——
    "manga.chapter.x": "第 {n} 章", // chapter-list.tsx, chapter-downloads-drawer.tsx, downloaded-chapter-list.tsx
    "manga.chapter_list.name": "名称", // chapter-list.tsx 表头
    "manga.chapter_list.scanlator": "扫译者", // chapter-list.tsx 表头/筛选
    "manga.chapter_list.number": "编号", // chapter-list.tsx, downloaded-chapter-list.tsx 表头
    "manga.chapter_list.provider": "提供商", // downloaded-chapter-list.tsx 表头
    "manga.chapter_list.chapter_id": "章节 ID", // downloaded-chapter-list.tsx 表头
    "manga.chapter_list.queued": "排队中", // chapter-list.tsx, downloaded-chapter-list.tsx
    "manga.chapter_list.source": "来源", // chapter-list.tsx leftAddon
    "manga.chapter_list.manual_match": "手动匹配", // chapter-list.tsx
    "manga.chapter_list.no_chapters": "未找到章节", // chapter-list.tsx
    "manga.chapter_list.try_another_source": "尝试其他来源", // chapter-list.tsx
    "manga.chapter_list.chapters": "章节", // chapter-list.tsx h2, downloaded-chapter-list.tsx 表头
    "manga.reader.continue_reading": "继续阅读", // chapter-list.tsx
    "manga.chapter_list.show_unread": "仅显示未读", // chapter-list.tsx
    "manga.chapter_list.show_downloaded": "显示已下载", // chapter-list.tsx
    "manga.chapter_list.download_unread_title": "下载未读章节", // chapter-list.tsx
    "manga.chapter_list.add_to_queue": "加入队列", // chapter-list.tsx
    "manga.chapter_list.add_unread_confirm": "将 {count} 个未读章节加入下载队列？", // chapter-list.tsx
    "manga.chapter_list.next_chapter_heading": "下一章", // chapter-list.tsx 命令面板
    "manga.chapter_list.upcoming_chapters_heading": "后续章节", // chapter-list.tsx 命令面板
    "manga.chapter_list.refresh_source": "刷新来源", // chapter-list.tsx
    "manga.chapter_list.saved_source_unavailable": "已保存的来源不可用", // chapter-list.tsx
    "manga.chapter_list.no_source_available": "无可用漫画来源", // chapter-list.tsx
    "manga.chapter_list.saved_source_uninstalled": "已保存的来源「{provider}」当前未安装。请重新安装或选择其他来源以继续。", // chapter-list.tsx
    "manga.chapter_list.install_source_hint": "请安装或选择漫画来源以加载章节。", // chapter-list.tsx
    "manga.chapter_list.filters_error": "无法加载章节筛选器。请联系扩展开发者：\"{error}\"", // chapter-list.tsx
    "manga.chapter_list.delete_selected": "删除所选章节（{count}）", // downloaded-chapter-list.tsx
    "manga.chapter_list.show_queued": "显示排队中", // downloaded-chapter-list.tsx
    "manga.refresh.started_toast": "来源刷新已开始", // chapter-list.tsx
    "manga.refresh.not_eligible_title": "仅在看与重温中的漫画可刷新来源", // chapter-list.tsx title
    "manga.refresh.refresh_running": "刷新进行中", // chapter-list.tsx

    // —— 来源刷新弹窗（manga-source-refresh-modal.tsx）——
    "manga.refresh.title": "刷新漫画来源", // MangaSourceRefreshModal 标题
    "manga.refresh.action_close": "关闭", // 底部按钮
    "manga.refresh.stop": "停止刷新", // 底部按钮
    "manga.refresh.stopping": "正在停止...", // 底部按钮
    "manga.refresh.stopping_after_current": "将在当前请求完成后停止", // 状态文案
    "manga.refresh.refreshing_selected": "正在刷新已选来源", // 状态文案
    "manga.refresh.searching_installed": "正在搜索已安装来源", // 状态文案
    "manga.refresh.progress_of": "{current} / {total}", // 进度计数
    "manga.refresh.progress_aria": "来源刷新进度", // ProgressBar aria-label
    "manga.refresh.background_notice": "可以关闭此弹窗，刷新将在后台继续。", // 运行中提示
    "manga.refresh.complete": "来源刷新完成", // 终态
    "manga.refresh.stopped": "来源刷新已停止", // 终态
    "manga.refresh.failed": "来源刷新失败", // 终态
    "manga.refresh.review_issues": "查看问题（{count}）", // 问题折叠面板
    "manga.refresh.issue_not_found_saved": "已保存的来源未返回章节。", // 问题明细
    "manga.refresh.issue_not_found": "未找到匹配的来源。", // 问题明细
    "manga.refresh.issue_provider_error": "一个或多个提供商出现错误。", // 问题明细
    "manga.refresh.providers_installed": "已安装 {count} 个提供商", // 安装计数
    "manga.refresh.current_and_rereading_only": "仅限在看与重温中的漫画", // 范围说明
    "manga.refresh.find_alternatives": "查找替代来源", // 底部按钮
    "manga.refresh.retry_failed": "重试失败项", // 底部按钮
    "manga.refresh.run_again": "再次运行", // 底部按钮
    "manga.refresh.done": "完成", // 底部按钮
    "manga.refresh.start": "开始刷新", // 底部按钮
    "manga.refresh.mode_refresh_selected_title": "刷新已选来源", // 模式单选
    "manga.refresh.mode_refresh_selected_desc": "更新已保存来源的漫画。", // 模式单选
    "manga.refresh.mode_find_missing_title": "查找缺失来源", // 模式单选
    "manga.refresh.mode_find_missing_desc": "在所有已安装提供商中搜索没有来源的漫画。", // 模式单选
    "manga.refresh.mode_refresh_and_find_title": "刷新并查找缺失来源", // 模式单选
    "manga.refresh.mode_refresh_and_find_desc": "先更新已保存的来源，再搜索缺失来源。", // 模式单选
    "manga.refresh.mode_reevaluate_all_title": "重新评估全部来源", // 模式单选
    "manga.refresh.mode_reevaluate_all_desc": "对比所有已安装提供商，并允许更换已保存的来源。", // 模式单选
    "manga.refresh.alert_reevaluate": "当其他提供商拥有更多独立章节时，现有来源选择可能被替换。", // Alert
    "manga.refresh.alert_waiting_sync": "正在等待服务端漫画偏好设置同步完成。", // Alert
    "manga.refresh.alert_install_provider": "请安装漫画提供商以搜索缺失或替代来源。", // Alert
    "manga.refresh.summary_refreshed": "{count} 个已刷新", // formatRefreshSummary
    "manga.refresh.summary_found": "{count} 个已找到", // formatRefreshSummary
    "manga.refresh.summary_replaced": "{count} 个已更换", // formatRefreshSummary
    "manga.refresh.summary_not_found": "{count} 个未找到", // formatRefreshSummary
    "manga.refresh.summary_failed": "{count} 个失败", // formatRefreshSummary

    // —— 阅读器设置抽屉（chapter-reader-settings.tsx）——
    "manga.settings.title": "设置", // Drawer title
    "manga.settings.open": "打开设置", // 下拉菜单
    "manga.settings.toggle_fullscreen": "切换全屏", // 下拉菜单
    "manga.settings.show_bar": "显示工具栏", // 下拉菜单
    "manga.settings.hide_bar": "隐藏工具栏", // 下拉菜单
    "manga.settings.reading_mode": "阅读模式", // RadioGroup label
    "manga.settings.mode_long_strip": "条带模式", // 模式选项
    "manga.settings.mode_single_page": "单页", // 模式选项
    "manga.settings.mode_double_page": "双页", // 模式选项
    "manga.settings.offset": "偏移量", // NumberInput label
    "manga.settings.page_fit": "页面适配", // RadioGroup label
    "manga.settings.fit_contain": "适应高度", // 适配选项
    "manga.settings.fit_overflow": "高度溢出", // 适配选项
    "manga.settings.fit_cover": "适应宽度", // 适配选项
    "manga.settings.fit_true_size": "原始尺寸", // 适配选项
    "manga.settings.zoom": "缩放", // NumberInput label
    "manga.settings.page_container_width": "页面容器宽度", // NumberInput label
    "manga.settings.page_stretch": "页面拉伸", // RadioGroup label
    "manga.settings.page_stretch_help": "「拉伸」在条带模式下会强制所有页面与容器同宽。", // help
    "manga.settings.reset_defaults_for": "重置 {mode} 的默认设置", // 重置按钮
    "manga.settings.page_gap": "页面间距", // Switch label
    "manga.settings.page_gap_shadow": "页面间距阴影", // Switch label
    "manga.settings.reading_direction": "阅读方向", // RadioGroup label
    "manga.settings.dir_ltr": "从左到右", // 方向选项
    "manga.settings.dir_rtl": "从右到左", // 方向选项
    "manga.settings.progress_bar": "进度条", // Switch label
    "manga.settings.stretch_none": "无", // 拉伸选项
    "manga.settings.stretch_stretch": "拉伸", // 拉伸选项
    "manga.settings.editable_keybinds": "可编辑快捷键", // h4
    "manga.settings.click_to_edit": "点击编辑", // 提示
    "manga.settings.prev_chapter": "上一章", // 快捷键行
    "manga.settings.next_chapter": "下一章", // 快捷键行
    "manga.settings.prev_page": "上一页", // 快捷键行
    "manga.settings.next_page": "下一页", // 快捷键行
    "manga.settings.keyboard_shortcuts": "键盘快捷键", // h4
    "manga.settings.kbs_update_progress": "更新进度并前往下一章", // 快捷键说明
    "manga.settings.kbs_toggle_bar": "切换底部工具栏显示", // 快捷键说明
    "manga.settings.kbs_switch_mode": "切换阅读模式", // 快捷键说明
    "manga.settings.kbs_switch_direction": "切换阅读方向", // 快捷键说明
    "manga.settings.kbs_switch_fit": "切换页面适配", // 快捷键说明
    "manga.settings.kbs_switch_stretch": "切换页面拉伸", // 快捷键说明
    "manga.settings.kbs_inc_offset": "增加双页偏移", // 快捷键说明
    "manga.settings.kbs_dec_offset": "减少双页偏移", // 快捷键说明
    "manga.settings.toast_double_page_small_screen": "小屏幕设备不支持双页模式。", // toast.error

    // —— 手动匹配（manga-manual-mapping-modal.tsx）——
    "manga.manual_match.title": "手动匹配", // Modal/确认弹窗标题
    "manga.manual_match.description": "将此漫画匹配到搜索结果", // Modal description
    "manga.manual_match.confirm_match": "确认将此条目匹配到「{title}」？", // 确认弹窗
    "manga.manual_match.distinct_chapters": "独立章节数", // 预览字段
    "manga.manual_match.latest_chapter": "最新章节", // 预览字段
    "manga.manual_match.unknown": "未知", // 预览字段
    "manga.manual_match.languages_label": "语言：", // 预览字段
    "manga.manual_match.scanlators_label": "扫译者：", // 预览字段
    "manga.manual_match.review_notice": "保存匹配前请先核对提供商的结果。", // 确认弹窗描述
    "manga.manual_match.current_mapping": "当前匹配：", // 现有映射
    "manga.manual_match.remove_mapping": "移除匹配", // 按钮
    "manga.manual_match.no_match": "暂无手动匹配", // 空状态
    "manga.manual_match.search_placeholder": "输入标题...", // 输入框
    "manga.manual_match.no_chapters_toast": "该结果未找到章节", // toast.error

    // —— 下载队列（chapter-downloads-drawer.tsx / chapter-downloads-button.tsx）——
    "manga.action.stop": "停止", // chapter-downloads-drawer.tsx
    "manga.action.clear_all": "清空全部", // chapter-downloads-drawer.tsx
    "manga.action.start": "开始", // chapter-downloads-drawer.tsx
    "manga.downloads.button": "下载", // ChapterDownloadsButton
    "manga.downloads.title": "已下载章节", // 抽屉标题 / h3
    "manga.downloads.queue": "队列", // h3
    "manga.downloads.reset_errored": "重置出错章节", // 按钮
    "manga.downloads.fetch_queue_failed": "无法获取下载队列", // LuffyError
    "manga.downloads.oops": "糟糕！", // LuffyError title
    "manga.downloads.errored": "出错", // 状态
    "manga.downloads.empty_queue": "队列中暂无任务", // 空状态
    "manga.downloads.downloaded_header": "已下载", // h3
    "manga.downloads.media_id": "媒体 {id}", // 列表项
    "manga.downloads.n_chapters": "{count} 个章节", // 列表项
    "manga.downloads.not_in_anilist": "不在你的 AniList 收藏中", // 列表项
    "manga.downloads.n_chapters_badge": "共 {count} 章", // 卡片角标
    "manga.downloads.no_chapters": "尚未下载章节", // 空状态

    // —— 阅读器（chapter-reader-drawer.tsx / chapter-page.tsx / manga-reader-bar.tsx / handle-chapter-reader.ts）——
    "manga.reader.toast_no_page_dimensions": "无法从该提供商获取页面尺寸，正在切换到分页模式。", // toast.error
    "manga.reader.close_reader": "关闭阅读器", // 命令面板
    "manga.reader.heading": "阅读器", // 命令面板分组
    "manga.reader.update_progress": "更新进度（{current} / {total}）", // 进度按钮
    "manga.reader.failed_to_load_pages": "页面加载失败", // LuffyError title
    "manga.reader.load_pages_error": "加载本章页面时发生错误。", // 错误描述
    "manga.reader.load_pages_error_hint": "请重新加载页面、重新加载来源或更换来源。", // 错误提示
    "manga.reader.page_alt": "第 {n} 页", // chapter-page alt
    "manga.reader.chapter_prefix": "第", // manga-reader-bar 章节号前缀
    "manga.reader.chapter_suffix": "章", // manga-reader-bar 章节号后缀
    "manga.reader.zoom_flash": "缩放：{value}%", // handle-chapter-reader 闪现提示
    "manga.reader.offset_flash": "双页偏移：{value}", // handle-chapter-reader 闪现提示

    // —— 漫画媒体库（manga-library-view.tsx / library-header.tsx）——
    "manga.library.not_found_title": "未找到漫画", // 回填：原裸中文
    "manga.library.empty_hint": "您的媒体库中尚未添加任何漫画。", // 回填：原裸中文
    "manga.library.browse_manga": "浏览漫画", // 回填：原裸中文
    "manga.library.aria_refresh_running": "来源刷新进行中", // aria-label
    "manga.library.view_source_refresh": "查看漫画来源刷新", // Tooltip
    "manga.library.refresh_manga_sources": "刷新漫画来源", // Tooltip
    "manga.library.aria_list_actions": "漫画列表操作", // aria-label
    "manga.library.view_refresh_short": "查看来源刷新", // 下拉菜单
    "manga.library.refresh_sources_short": "刷新来源", // 下拉菜单
    "manga.library.unread_only": "仅显示未读章节", // 下拉菜单
    "manga.library.hide_unread_counts": "隐藏未读数", // 下拉菜单
    "manga.library.show_unread_counts": "显示未读数", // 下拉菜单
    "manga.library.banner_alt": "横幅图片", // library-header alt

    // —— 关联与推荐（manga-recommendations.tsx）——
    "manga.relations.movie_suffix": "（剧场版）", // 关联作品角标后缀
} satisfies Dictionary
