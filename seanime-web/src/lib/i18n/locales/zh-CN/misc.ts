import type { Dictionary } from "../../types"

/**
 * misc 词条表
 * 归属：长尾杂项（lists/sync/debrid/scan-summaries/mal/mediastream/onlinestream/custom-sources/auth/webview/error）
 * 独占：09-04
 *
 * 只读复用：common.* / library.* 等同名文案按契约复用，不新建重复词条。
 * 已复用：search.type.anime/manga、common.state.*、common.action.retry/refresh、
 *   common.empty.title/debrid_torrents、library.filter.*、library.search.sorting、
 *   library.explorer.end/delete_desc/episode_number、search.filter.genre*、
 *   media.action.save_locally/remove_offline_data/open_in_library_explorer、
 *   media.confirm.remove_offline_data_description、torrent.stats.downloading/seeding、
 *   torrent.filter.downloading、torrent.client.remove_single/per_page、
 *   torrent.qbittorrent.open_in_browser、search.page.custom_source、navigation.item.search
 */
export const miscDictionary = {
    // ===== lists（anilist-collection-lists.tsx）=====
    "misc.lists.tab_stats": "数据统计", // anilist-collection-lists.tsx StaticTabs
    "misc.lists.all_lists": "全部列表", // anilist-collection-lists.tsx Select label

    // ===== lists（anilist-stats.tsx）=====
    "misc.stats.format_tv_short": "TV 短篇", // anilist-stats.tsx formatName
    "misc.stats.format_movie": "剧场版", // anilist-stats.tsx formatName
    "misc.stats.format_special": "特别篇", // anilist-stats.tsx formatName
    "misc.stats.format_music": "音乐", // anilist-stats.tsx formatName
    "misc.stats.total_anime": "动漫总数", // anilist-stats.tsx MetricCard label
    "misc.stats.episodes": "观看集数", // anilist-stats.tsx MetricCard label
    "misc.stats.watch_time": "观看时长", // anilist-stats.tsx MetricCard label
    "misc.stats.mean_score": "平均评分", // anilist-stats.tsx MetricCard label
    "misc.stats.started_this_year": "今年开始", // anilist-stats.tsx MetricCard label
    "misc.stats.started_last_year": "去年开始", // anilist-stats.tsx MetricCard label
    "misc.stats.no_activity": "暂无动态", // anilist-stats.tsx MetricCard sub
    "misc.stats.total_manga": "漫画总数", // anilist-stats.tsx MetricCard label
    "misc.stats.chapters": "已读章节", // anilist-stats.tsx MetricCard label
    "misc.stats.n_hours": "{count} 小时", // anilist-stats.tsx 观看时长 sub
    "misc.stats.n_days": "{count} 天", // anilist-stats.tsx daysWatched
    "misc.stats.score_distribution": "评分分布", // anilist-stats.tsx ChartSection title
    "misc.stats.by_score": "按评分统计条目数", // anilist-stats.tsx ChartSection description
    "misc.stats.list_status": "列表状态", // anilist-stats.tsx StatusChart title
    "misc.stats.status_description": "按当前列表状态的条目数", // anilist-stats.tsx StatusChart description
    "misc.stats.format_mix": "按条目数的格式占比", // anilist-stats.tsx ChartSection description
    "misc.stats.top_genres": "热门流派", // anilist-stats.tsx ChartSection title
    "misc.stats.most_watched_genres": "最常观看的流派", // anilist-stats.tsx ChartSection description
    "misc.stats.most_read_genres": "最常阅读的流派", // anilist-stats.tsx ChartSection description
    "misc.stats.started_by_year": "按开始年份统计", // anilist-stats.tsx ChartSection title
    "misc.stats.titles_started_each_year": "每年开始观看的条目数", // anilist-stats.tsx ChartSection description
    "misc.stats.release_years": "发布年份", // anilist-stats.tsx ChartSection title
    "misc.stats.by_release_year": "按原始发布年份分组的条目数", // anilist-stats.tsx ChartSection description
    "misc.stats.top_studios": "热门制作公司", // anilist-stats.tsx ChartSection title
    "misc.stats.studios_most_watched": "观看条目最多的制作公司", // anilist-stats.tsx ChartSection description
    "misc.stats.titles": "条目数", // anilist-stats.tsx 图表 categories
    "misc.stats.most_used_score": "最常用评分", // anilist-stats.tsx Highlight label
    "misc.stats.completion": "完结率", // anilist-stats.tsx Highlight label
    "misc.stats.top_genre": "主要流派", // anilist-stats.tsx Highlight label
    "misc.stats.peak_start_year": "开始高峰年份", // anilist-stats.tsx Highlight label
    "misc.stats.completion_detail": "已看完 {completed} 部 / 已开始 {started} 部", // anilist-stats.tsx Highlight detail
    "misc.stats.n_started": "{count} 部开始观看", // anilist-stats.tsx Highlight detail
    "misc.stats.n_titles": "{count} 部", // anilist-stats.tsx valueLabel
    "misc.stats.n_hours_watched": "观看 {hours} 小时", // anilist-stats.tsx meta
    "misc.stats.score_avg": "平均评分 {score}", // anilist-stats.tsx formatAvg
    "misc.stats.unit_titles": "部", // anilist-stats.tsx unit 参数
    "misc.stats.unit_chapters": "章", // anilist-stats.tsx secondaryLabel
    "misc.stats.unknown": "未知", // anilist-stats.tsx 兜底显示名

    // ===== debrid（debrid/page.tsx）=====
    "misc.debrid.not_enabled_title": "Debrid 未启用", // debrid/page.tsx LuffyError title
    "misc.debrid.not_enabled": "Debrid 服务未启用或未配置", // debrid/page.tsx LuffyError 内容
    "misc.debrid.failed_title": "连接失败", // debrid/page.tsx LuffyError title
    "misc.debrid.failed": "无法连接到 Debrid 服务，请检查你的设置。", // debrid/page.tsx 错误提示
    "misc.debrid.see_torrents": "查看你的 Debrid 服务种子任务", // debrid/page.tsx 副标题
    "misc.debrid.refreshed": "已刷新", // debrid/page.tsx toast
    "misc.debrid.dashboard": "控制面板", // debrid/page.tsx 按钮
    "misc.debrid.queued": "排队中", // debrid/page.tsx Badge label
    "misc.debrid.waiting": "等待中...", // debrid/page.tsx Badge label
    "misc.debrid.downloading_locally": "正在下载到本地", // debrid/page.tsx Tooltip
    "misc.debrid.preparing": "正在准备本地文件", // debrid/page.tsx 提示
    "misc.debrid.cancel_download": "取消下载", // debrid/page.tsx Tooltip
    "misc.debrid.download": "下载", // debrid/page.tsx Modal title + 按钮
    "misc.debrid.destination": "保存位置", // debrid/page.tsx DirectorySelector label
    "misc.debrid.destination_help": "种子的保存位置", // debrid/page.tsx DirectorySelector help

    // ===== sync（sync/page.tsx + sync-add-media-modal.tsx）=====
    "misc.sync.not_logged_in": "未登录账号", // sync/page.tsx LuffyError title
    "misc.sync.login_required": "该功能仅对已登录认证的用户开放。", // sync/page.tsx LuffyError 内容
    "misc.sync.enable_offline": "启用离线模式", // sync/page.tsx 按钮
    "misc.sync.disable_offline": "停用离线模式", // sync/page.tsx 按钮
    "misc.sync.title": "离线媒体", // sync/page.tsx h2
    "misc.sync.description": "查看你已保存到本地供离线使用的媒体。", // sync/page.tsx 副标题
    "misc.sync.modal_title": "同步", // sync/page.tsx Modal title
    "misc.sync.sync_now": "立即同步", // sync/page.tsx 按钮
    "misc.sync.update_local": "更新本地数据", // sync/page.tsx 按钮
    "misc.sync.update_local_desc_head": "用 AniList 的数据更新你的本地快照。这将覆盖你的离线更改。你可以在 ", // sync/page.tsx 说明前段
    "misc.sync.settings_path": "设置 > 应用 > 离线模式", // sync/page.tsx kbd 内容
    "misc.sync.update_local_desc_tail": " 中自动执行此操作。", // sync/page.tsx 说明后段
    "misc.sync.upload_to_anilist": "上传本地更改到 AniList", // sync/page.tsx 按钮
    "misc.sync.upload_desc": "用本地快照的数据更新你的 AniList 列表。应在离线做出更改后执行。", // sync/page.tsx 说明
    "misc.sync.irreversible": "更改不可撤销。", // sync/page.tsx Alert description
    "misc.sync.local_storage_size": "本地存储大小：", // sync/page.tsx 提示
    "misc.sync.unsaved_alert": "你有 {items} 尚未保存到本地，这些内容你当前正在{status}。", // sync/page.tsx Alert description
    "misc.sync.n_anime": "{count} 部动漫", // sync/page.tsx 拼接段
    "misc.sync.n_manga": "{count} 部漫画", // sync/page.tsx 拼接段
    "misc.sync.and": " 和 ", // sync/page.tsx 连接词
    "misc.sync.watching": "观看", // sync/page.tsx 拼接段
    "misc.sync.reading": "阅读", // sync/page.tsx 拼接段
    "misc.sync.unsynced_warning": "你有尚未同步到 AniList 的本地更改。", // sync/page.tsx Alert
    "misc.sync.autorefresh_paused": "离线数据的自动刷新将会暂停。", // sync/page.tsx Alert
    "misc.sync.upload_local_changes": "上传本地更改", // sync/page.tsx 按钮
    "misc.sync.delete_local_changes": "删除本地更改", // sync/page.tsx 按钮
    "misc.sync.no_tracked_media": "暂无已保存的媒体", // sync/page.tsx LuffyError title
    "misc.sync.saved_anime": "已保存的动漫", // sync/page.tsx h3
    "misc.sync.saved_manga": "已保存的漫画", // sync/page.tsx h3
    "misc.sync.syncing": "同步中", // sync/page.tsx Badge
    "misc.sync.local_changes_ignored": "已忽略本地更改。", // sync/page.tsx toast
    "misc.sync.saved_media": "已保存的媒体", // sync-add-media-modal.tsx Modal title
    "misc.sync.select_media": "选择要保存的媒体", // sync-add-media-modal.tsx 按钮
    "misc.sync.select_media_desc": "选择你想保存到本地的媒体。点击已保存的媒体可将其从本地存储中移除。", // sync-add-media-modal.tsx 说明

    // ===== scan-summaries（scan-summaries/page.tsx）=====
    "misc.scan.title": "扫描报告概览", // scan-summaries/page.tsx h2
    "misc.scan.description": "查看媒体库最新扫描的详细记录与匹配日志", // scan-summaries/page.tsx 副标题
    "misc.scan.empty": "暂无扫描报告", // scan-summaries/page.tsx 空状态
    "misc.scan.search_placeholder": "搜索文件名...", // scan-summaries/page.tsx placeholder
    "misc.scan.scanned_n_media": "Seanime 成功扫描了 {count} 个媒体", // scan-summaries/page.tsx 统计
    "misc.scan.n_matches": "({count} 个匹配)", // scan-summaries/page.tsx 统计
    "misc.scan.n_unmatched": "{count} 个文件未匹配成功", // scan-summaries/page.tsx 统计
    "misc.scan.unmatched_files": "未匹配文件", // scan-summaries/page.tsx h5
    "misc.scan.scanned_media": "已扫描媒体", // scan-summaries/page.tsx h5
    "misc.scan.collection_status": "该动漫{status}你的 AniList 收藏中", // scan-summaries/page.tsx 收藏状态
    "misc.scan.present": "已存在于", // scan-summaries/page.tsx 收藏状态
    "misc.scan.not_present": "不存在于", // scan-summaries/page.tsx 收藏状态
    "misc.scan.n_files_scanned": "已扫描 {count} 个文件", // scan-summaries/page.tsx 统计
    "misc.scan.errors_found": "发现错误", // scan-summaries/page.tsx 提示
    "misc.scan.warnings_found": "发现警告", // scan-summaries/page.tsx 提示
    "misc.scan.view_logs": "查看扫描日志", // scan-summaries/page.tsx AccordionTrigger
    "misc.scan.label_title": "标题", // scan-summaries/page.tsx 解析字段
    "misc.scan.label_part": "部分", // scan-summaries/page.tsx 解析字段
    "misc.scan.label_episode_title": "剧集标题", // scan-summaries/page.tsx 解析字段
    "misc.scan.report_github": "请在 GitHub 仓库上报此问题", // scan-summaries/page.tsx PANIC 提示

    // ===== mediastream（mediastream/page.tsx + test/page.tsx）=====
    "misc.mediastream.track": "轨道 {index}", // mediastream/page.tsx 字幕轨道
    "misc.mediastream.playback_error": "播放错误", // mediastream/page.tsx LuffyError title
    "misc.mediastream.could_not_load": "无法加载媒体容器。", // mediastream/page.tsx 错误兜底
    "misc.mediastream.complete_movie": "完整剧场版", // mediastream/page.tsx 剧集标题
    "misc.mediastream.episode_n": "第 {number} 集", // mediastream/page.tsx、onlinestream-page.tsx 剧集标题
    "misc.mediastream.file_path_missing": "未找到该剧集的文件路径", // mediastream/page.tsx toast
    "misc.mediastream.playback": "播放", // mediastream/page.tsx Modal title
    "misc.mediastream.playback_info": "播放信息", // mediastream/page.tsx 按钮
    "misc.mediastream.direct_play_ok": "文件的音视频编码与此客户端兼容。推荐直接播放。", // mediastream/page.tsx Alert
    "misc.mediastream.transcode_needed": "文件的音视频编码与此客户端不兼容。需要转码。", // mediastream/page.tsx Alert
    "misc.mediastream.stream_type": "流类型：", // mediastream/page.tsx 字段
    "misc.mediastream.video_codec": "视频编码：", // mediastream/page.tsx 字段
    "misc.mediastream.audio_codec": "音频编码：", // mediastream/page.tsx 字段
    "misc.mediastream.container_data": "媒体容器数据", // mediastream/page.tsx Modal title
    "misc.mediastream.more_data": "更多数据", // mediastream/page.tsx 按钮
    "misc.mediastream.switch_to_transcoding": "切换到转码", // mediastream/page.tsx 按钮
    "misc.mediastream.switch_hint": "若想切换到转码，请在媒体流设置中启用「优先转码」", // mediastream/page.tsx 提示
    "misc.mediastream.switch_to_direct": "切换到直接播放", // mediastream/page.tsx 按钮
    "misc.mediastream.switch_to_grid": "切换到网格视图", // mediastream/page.tsx、onlinestream-page.tsx title
    "misc.mediastream.switch_to_list": "切换到列表视图", // mediastream/page.tsx、onlinestream-page.tsx title
    "misc.mediastream.test_go_away": "请离开。", // mediastream/test/page.tsx 测试页文案

    // ===== onlinestream（onlinestream-page.tsx + manual-matching + auto-provider-cycler）=====
    "misc.onlinestream.cancel_trying": "取消尝试", // onlinestream-page.tsx 按钮
    "misc.onlinestream.try_all": "尝试所有可用提供者", // onlinestream-page.tsx 按钮
    "misc.onlinestream.find_other_providers": "查找其他提供者", // onlinestream-page.tsx Select label
    "misc.onlinestream.select_provider": "选择提供者", // onlinestream-page.tsx placeholder
    "misc.onlinestream.stream": "流", // onlinestream-page.tsx Popover title
    "misc.onlinestream.cache": "缓存", // onlinestream-page.tsx 按钮
    "misc.onlinestream.cache_desc": "如果流播放出现问题，请清空缓存。", // onlinestream-page.tsx 说明
    "misc.onlinestream.empty_cache": "清空流缓存", // onlinestream-page.tsx 按钮
    "misc.onlinestream.manual_match": "手动匹配", // onlinestream-page.tsx、onlinestream-manual-matching.tsx
    "misc.onlinestream.no_provider": "未选择提供者", // onlinestream-page.tsx LuffyError title
    "misc.onlinestream.provider_error": "提供者错误", // onlinestream-page.tsx LuffyError title
    "misc.onlinestream.episode_list_error": "无法从提供者获取剧集列表。", // onlinestream-page.tsx 错误提示
    "misc.onlinestream.switch_to_subs": "切换到字幕", // onlinestream-page.tsx 按钮
    "misc.onlinestream.switch_to_dub": "切换到配音", // onlinestream-page.tsx 按钮
    "misc.onlinestream.hls_error": "HLS 错误：{message}", // onlinestream-page.tsx 错误前缀
    "misc.onlinestream.no_working_providers": "未找到可用的提供者", // use-onlinestream-auto-provider-cycler.ts
    "misc.onlinestream.no_dubbed_providers": "没有支持配音的提供者", // use-onlinestream-auto-provider-cycler.ts toast
    "misc.onlinestream.no_providers": "没有可用的提供者", // use-onlinestream-auto-provider-cycler.ts toast
    "misc.onlinestream.stopped_trying": "已停止尝试提供者", // use-onlinestream-auto-provider-cycler.ts toast
    "misc.onlinestream.manual_match_desc": "将此动漫匹配到提供者的搜索结果。", // onlinestream-manual-matching.tsx Modal description
    "misc.onlinestream.current_mapping": "当前匹配：", // onlinestream-manual-matching.tsx 提示
    "misc.onlinestream.remove_mapping": "移除匹配", // onlinestream-manual-matching.tsx 按钮
    "misc.onlinestream.no_manual_match": "无手动匹配", // onlinestream-manual-matching.tsx 空状态
    "misc.onlinestream.search_placeholder": "输入标题...", // onlinestream-manual-matching.tsx placeholder
    "misc.onlinestream.look_for_dubs": "查找配音", // onlinestream-manual-matching.tsx Switch label
    "misc.onlinestream.dubs_help": "仅对搜索结果支持配音的提供者生效。", // onlinestream-manual-matching.tsx moreHelp
    "misc.onlinestream.confirm_match_desc": "确定要将此动漫匹配到该搜索结果吗？", // onlinestream-manual-matching.tsx 确认弹窗

    // ===== custom-sources（custom-sources/page.tsx）=====
    "misc.custom_sources.select_source": "选择源", // custom-sources/page.tsx placeholder
    "misc.custom_sources.search_titles": "搜索标题...", // custom-sources/page.tsx placeholder
    "misc.custom_sources.select_to_view": "选择一个源以查看其内容", // custom-sources/page.tsx 空状态
    "misc.custom_sources.load_failed": "加载内容失败", // custom-sources/page.tsx LuffyError title
    "misc.custom_sources.load_error": "从 {provider} 加载内容时出错", // custom-sources/page.tsx 错误提示
    "misc.custom_sources.no_results": "未找到结果", // custom-sources/page.tsx LuffyError title
    "misc.custom_sources.no_results_desc": "未找到符合搜索条件的{type}", // custom-sources/page.tsx 错误提示
    "misc.custom_sources.n_results": "{count} 条结果", // custom-sources/page.tsx 统计
    "misc.custom_sources.page_of": "第 {current} 页，共 {total} 页", // custom-sources/page.tsx 分页

    // ===== mal（mal/auth/callback/_page.tsx；mal/_page.tsx 命中均为注释误报）=====
    "misc.mal.invalid_url": "无效的 URL 或 Challenge", // mal/auth/callback/_page.tsx
    "misc.mal.authenticating": "正在认证...", // mal/auth/callback/_page.tsx

    // ===== 扫描日志查看器（app/scan-log-viewer/，09-05）=====
    "misc.scan_log.title": "扫描日志分析器", // page.tsx h1
    "misc.scan_log.clear_log": "清除日志", // page.tsx 按钮
    "misc.scan_log.loading_saved": "正在加载已保存的日志...", // page.tsx 加载态
    "misc.scan_log.drop_file": "拖入日志文件", // page.tsx 拖放提示
    "misc.scan_log.load_another": "加载其他文件", // page.tsx 按钮
    "misc.scan_log.load_file": "加载扫描日志文件", // page.tsx 按钮
    "misc.scan_log.save_failed_storage": "保存日志到浏览器存储失败", // page.tsx toast
    "misc.scan_log.restored": "已恢复上次保存的扫描日志", // page.tsx toast
    "misc.scan_log.saving": "正在保存日志...", // page.tsx toast
    "misc.scan_log.saved": "日志已保存，可在下次会话中使用", // page.tsx toast
    "misc.scan_log.save_failed": "保存日志失败", // page.tsx toast
    "misc.scan_log.cleared": "已清除保存的日志", // page.tsx toast
    "misc.scan_log.empty_hint": "加载扫描日志文件以开始分析", // scan-log-viewer.tsx 空状态
    "misc.scan_log.back": "← 返回", // scan-log-viewer.tsx 按钮
    "misc.scan_log.no_logs_file": "未找到该文件的日志。", // scan-log-viewer.tsx 空状态
    "misc.scan_log.no_logs_any_phase": "所有阶段均未找到该文件的日志。", // scan-log-viewer.tsx 空状态
    "misc.scan_log.tab_overview": "概览", // scan-log-viewer.tsx 页签
    "misc.scan_log.tab_parsed": "已解析文件", // scan-log-viewer.tsx 页签
    "misc.scan_log.tab_matcher": "匹配器", // scan-log-viewer.tsx 页签/流程标题
    "misc.scan_log.tab_hydrator": "补全器", // scan-log-viewer.tsx 页签/流程标题
    "misc.scan_log.tab_issues": "问题（{count}）", // scan-log-viewer.tsx 页签
    "misc.scan_log.parsing": "解析", // scan-log-viewer.tsx 流程标题
    "misc.scan_log.total_files": "文件总数", // scan-log-viewer.tsx StatCard
    "misc.scan_log.issues": "问题", // scan-log-viewer.tsx StatCard/筛选按钮
    "misc.scan_log.n_matched": "已匹配（{count}）", // scan-log-viewer.tsx 筛选按钮
    "misc.scan_log.n_unmatched": "未匹配（{count}）", // scan-log-viewer.tsx 筛选按钮
    "misc.scan_log.pipeline": "流水线", // scan-log-viewer.tsx 标题
    "misc.scan_log.file_discovery": "文件发现", // scan-log-viewer.tsx 流水线步骤
    "misc.scan_log.media_fetch": "媒体获取", // scan-log-viewer.tsx 流水线步骤
    "misc.scan_log.token_index": "分词索引", // scan-log-viewer.tsx 流水线步骤
    "misc.scan_log.n_files": "{count} 个文件", // scan-log-viewer.tsx 统计
    "misc.scan_log.n_media_new": "{count} 个媒体（{count2} 个新增）", // scan-log-viewer.tsx 统计
    "misc.scan_log.n_tokens": "{count} 个分词", // scan-log-viewer.tsx 统计
    "misc.scan_log.events": "事件", // scan-log-viewer.tsx 标题
    "misc.scan_log.n_issues": "{count} 个问题", // scan-log-viewer.tsx 统计
    "misc.scan_log.search_parsed": "搜索已解析文件...", // scan-log-viewer.tsx placeholder
    "misc.scan_log.view_full_flow": "查看完整流程", // scan-log-viewer.tsx 按钮
    "misc.scan_log.parsed_data": "解析数据", // scan-log-viewer.tsx 标题
    "misc.scan_log.search_match": "按文件名或匹配结果搜索...", // scan-log-viewer.tsx placeholder
    "misc.scan_log.search_media": "按文件名或媒体 ID 搜索...", // scan-log-viewer.tsx placeholder
    "misc.scan_log.search_issues": "搜索问题...", // scan-log-viewer.tsx placeholder
    "misc.scan_log.no_issues": "未发现问题", // scan-log-viewer.tsx 空状态
    "misc.scan_log.candidates_found": "{count} 个候选", // scan-log-viewer.tsx 统计

    // ===== 问题报告查看器（app/issue-report/page.tsx，09-05）=====
    "misc.issue_report.title": "问题报告", // page.tsx h1
    "misc.issue_report.loading_saved": "正在加载已保存的报告...", // page.tsx 加载态
    "misc.issue_report.drop_file": "拖入报告文件", // page.tsx 拖放提示
    "misc.issue_report.load_another": "加载其他文件", // page.tsx 按钮
    "misc.issue_report.load_file": "加载报告文件", // page.tsx 按钮
    "misc.issue_report.clear_report": "清除报告", // page.tsx 按钮
    "misc.issue_report.cleared": "已清除报告", // page.tsx toast
    "misc.issue_report.empty_hint": "加载问题报告 JSON 或 ZIP 文件", // page.tsx 空状态
    "misc.issue_report.save_failed_storage": "保存报告到浏览器存储失败", // page.tsx toast
    "misc.issue_report.restored": "已恢复上次保存的问题报告", // page.tsx toast
    "misc.issue_report.decompressing": "正在解压报告...", // page.tsx toast
    "misc.issue_report.loaded": "报告已加载", // page.tsx toast
    "misc.issue_report.decompress_failed": "解压报告失败", // page.tsx toast
    "misc.issue_report.parse_failed": "解析报告失败", // page.tsx toast
    "misc.issue_report.tab_replay": "会话回放", // page.tsx 页签
    "misc.issue_report.tab_timeline": "时间线", // page.tsx 页签
    "misc.issue_report.tab_network": "网络（{count}）", // page.tsx 页签
    "misc.issue_report.tab_console": "控制台（{count}）", // page.tsx 页签
    "misc.issue_report.tab_clicks": "点击（{count}）", // page.tsx 页签
    "misc.issue_report.tab_websocket": "WebSocket（{count}）", // page.tsx 页签
    "misc.issue_report.tab_server_logs": "服务端日志", // page.tsx 页签
    "misc.issue_report.tab_screenshots": "截图（{count}）", // page.tsx 页签
    "misc.issue_report.tab_scan_logs": "扫描日志", // page.tsx 页签
    "misc.issue_report.report_info": "报告信息", // page.tsx 标题
    "misc.issue_report.version": "版本", // page.tsx 明细
    "misc.issue_report.os": "操作系统", // page.tsx 明细
    "misc.issue_report.user_agent": "用户代理", // page.tsx 明细
    "misc.issue_report.created": "创建时间", // page.tsx 明细
    "misc.issue_report.viewport": "视口", // page.tsx 明细
    "misc.issue_report.user_description": "用户描述", // page.tsx 标题
    "misc.issue_report.total_events": "事件总数", // page.tsx StatCard
    "misc.issue_report.errors": "错误", // page.tsx StatCard/筛选
    "misc.issue_report.n_errors": "错误（{count}）", // page.tsx 筛选
    "misc.issue_report.n_warnings": "警告（{count}）", // page.tsx 筛选
    "misc.issue_report.net_errors": "网络错误", // page.tsx StatCard
    "misc.issue_report.recording": "录制", // page.tsx StatCard
    "misc.issue_report.clicks": "点击", // page.tsx 流水线/筛选
    "misc.issue_report.network": "网络", // page.tsx 流水线/筛选
    "misc.issue_report.console": "控制台", // page.tsx 流水线/筛选
    "misc.issue_report.queries": "查询", // page.tsx 流水线/筛选
    "misc.issue_report.navigations": "页面导航", // page.tsx 流水线
    "misc.issue_report.screenshots": "截图", // page.tsx 流水线/筛选
    "misc.issue_report.nav_short": "导航", // page.tsx 筛选
    "misc.issue_report.server": "服务端", // page.tsx 筛选
    "misc.issue_report.n_req": "{count} 个请求", // page.tsx 流水线
    "misc.issue_report.server_status": "服务端状态", // page.tsx 标题
    "misc.issue_report.unlocked_local_files": "已解锁的本地文件（{count}）", // page.tsx 标题
    "misc.issue_report.request_body": "请求体", // page.tsx 标题
    "misc.issue_report.response": "响应", // page.tsx 标题
    "misc.issue_report.no_screenshots": "此报告中没有截图", // page.tsx 空状态
    "misc.issue_report.replay_load_failed": "会话回放播放器加载失败", // page.tsx 回放错误
    "misc.issue_report.errors_only": "仅显示错误", // page.tsx Switch
    "misc.issue_report.include_server_logs": "包含服务端日志", // page.tsx Switch
    "misc.issue_report.logs": "日志", // page.tsx 筛选
    "misc.issue_report.incoming": "传入", // page.tsx WebSocket 筛选
    "misc.issue_report.outgoing": "传出", // page.tsx WebSocket 筛选
    "misc.issue_report.search_network": "搜索网络请求...", // page.tsx placeholder
    "misc.issue_report.search_console": "搜索控制台日志...", // page.tsx placeholder
    "misc.issue_report.search_clicks": "搜索点击事件...", // page.tsx placeholder
    "misc.issue_report.search_server": "搜索服务端日志...", // page.tsx placeholder
    "misc.issue_report.search_websocket": "搜索 WebSocket 事件...", // page.tsx placeholder
    "misc.issue_report.n_dom_events": "{count} 个 DOM 事件", // page.tsx 回放统计

    // ===== API 文档页（app/docs/page.tsx，09-05）=====
    "misc.docs.url_params": "URL 参数", // page.tsx h5
    "misc.docs.returns": "返回类型", // page.tsx 说明

    // ===== 服务端认证（app/public/auth/server-auth.tsx，09-05）=====
    "misc.auth.title": "需要密码", // server-auth.tsx Modal title
    "misc.auth.description": "此 Seanime 服务器需要身份验证。", // server-auth.tsx Modal description
    "misc.auth.password_required": "必须输入密码", // server-auth.tsx 表单校验
    "misc.auth.enter_password": "输入密码", // server-auth.tsx Field label

    // ===== 崩溃页（app/splashscreen/crash/page.tsx，09-05）=====
    "misc.crash.title": "出错了", // crash/page.tsx LuffyError title

    // ===== 组件调试页（app/(main)/test/page.tsx，09-05；Lorem* 为占位假文，有意保留）=====
    "misc.test.sections": "分区", // page.tsx 侧栏标题
    "misc.test.section_buttons": "按钮与操作", // page.tsx SECTIONS/h2
    "misc.test.section_forms": "表单与控件", // page.tsx SECTIONS/h2
    "misc.test.section_dialogs": "对话框与浮层", // page.tsx SECTIONS/h2
    "misc.test.section_feedback": "反馈与标签页", // page.tsx SECTIONS/h2
    "misc.test.section_layout": "布局与骨架屏", // page.tsx SECTIONS/h2
    "misc.test.loading": "加载中", // page.tsx Switch
    "misc.test.disable_buttons": "禁用按钮", // page.tsx Switch
    "misc.test.rounded": "圆角", // page.tsx Switch

    // ===== 更新日志引导（app/(main)/_features/tour/changelog-tour.tsx，09-05）=====
    "misc.changelog.whats_new_350": "3.5.0 有什么新功能？", // 3.5.0 标题
    "misc.changelog.whats_new_370": "3.7.0 有什么新功能？", // 3.7.0 标题
    "misc.changelog.whats_new_380": "3.8.0 有什么新功能？", // 3.8.0 标题
    "misc.changelog.whats_new_390": "3.9.0 有什么新功能？", // 3.9.0 标题
    "misc.changelog.whats_new_3100": "3.10.0 有什么新功能？", // 3.10.0 标题
    "misc.changelog.intro_features": "来看看这次的新功能。", // 3.5/3.7 引言
    "misc.changelog.intro_biggest": "来看看这个版本最重要的新内容。", // 3.8/3.9/3.10 引言
    "misc.changelog.new_scanner": "全新扫描器", // 3.5 步骤标题
    "misc.changelog.new_scanner_content": "扫描器的内部逻辑已彻底重构，现在使用更智能的上下文感知算法，匹配更准确。", // 3.5 内容
    "misc.changelog.new_scanner_content2": "扫描器现已支持使用 Anime Offline Database 匹配数据。", // 3.5 内容
    "misc.changelog.scanner_config": "扫描器配置", // 3.5 步骤标题
    "misc.changelog.scanner_config_content": "你现在可以微调扫描器的匹配行为。详情请查看文档。", // 3.5 内容
    "misc.changelog.issue_recorder": "问题记录器", // 3.5 步骤标题/图片 alt
    "misc.changelog.issue_recorder_content": "问题记录器已改进，现在可以录制 UI，让 Bug 报告更有参考价值。", // 3.5 内容
    "misc.changelog.search_content": "搜索菜单项现在会打开搜索页面。你仍然可以在任意页面按 'S' 快速搜索。", // 3.5 内容（标题复用 navigation.item.search）
    "misc.changelog.transcode_player": "转码播放器", // 3.5 步骤标题
    "misc.changelog.transcode_player_content": "转码/Direct Play 现在使用 Seanime Denshi 与在线流媒体相同的自定义播放器。", // 3.5 内容
    "misc.changelog.char_lookup_alt": "角色查找", // 3.5 图片 alt
    "misc.changelog.new_player_features": "新播放器功能", // 3.5 步骤标题
    "misc.changelog.char_lookup_content": "观看时按 'H' 可快速查找角色。按 'Z' 可切换 Stats for Nerds。", // 3.5 内容
    "misc.changelog.security_improvements": "安全性改进", // 3.7 步骤标题
    "misc.changelog.security_content": "3.7.0 包含多项安全性改进，包括安全模式。详情请查看文档。", // 3.7 内容
    "misc.changelog.tags_content": "搜索页面现在支持按标签搜索。", // 3.7 内容（标题复用 library.filter.tags）
    "misc.changelog.adult_global_search": "全局搜索中的成人内容", // 3.7 步骤标题
    "misc.changelog.adult_content": "启用成人内容后，全局搜索不再过滤成人条目。（提示：按 's' 打开全局搜索）", // 3.7 内容
    "misc.changelog.bug_fixes": "Bug 修复", // 步骤标题（3.7/3.8/3.9/3.10 共用）
    "misc.changelog.bug_fixes_content_37": "本版本修复了多个 Bug，包括与 Seanime Denshi 和插件相关的问题。详情请阅读完整更新日志。", // 3.7 内容
    "misc.changelog.bug_fixes_content_38": "本版本修复了多个 Bug，包括与内置播放器相关的问题。详情请阅读完整更新日志。", // 3.8 内容
    "misc.changelog.bug_fixes_content_39": "本版本修复了 MPV/IINA 进度追踪、漫画图片代理等问题，并将 Seanime Denshi 的 Electron 升级到 42.4.0。详情请阅读完整更新日志。", // 3.9 内容
    "misc.changelog.bug_fixes_content_310": "本版本还修复了 VideoCore 中字幕卡顿或缺失、MPV/IINA 观看进度、媒体容器检测等问题。详情请阅读完整更新日志。", // 3.10 内容
    "misc.changelog.torrent_search_downloads": "种子搜索与下载", // 3.8 步骤标题
    "misc.changelog.torrent_content": "种子搜索现在可以同时向多个提供者发出请求。本版本还修复了若干 Debrid 下载的边缘问题。", // 3.8 内容
    "misc.changelog.subtitle_translation": "字幕翻译", // 3.8 步骤标题
    "misc.changelog.subtitle_content": "字幕翻译器现已支持 OpenAI 兼容的本地 LLM，可使用 LM Studio、Ollama 等作为本地翻译后端。", // 3.8 内容
    "misc.changelog.local_subtitle_files": "本地字幕文件", // 3.8 步骤标题
    "misc.changelog.local_sub_content": "现在会自动从视频文件夹中读取本地字幕文件，外部播放器链接可使用新的 '{subtitleUrl}' 占位符引用这些字幕。", // 3.8 内容
    "misc.changelog.hide_spoilers": "隐藏剧透", // 3.8 步骤标题
    "misc.changelog.spoilers_content": "现在可以在整个应用中隐藏剧透。在动漫页面使用新的 '/spoilers' 命令可为该动漫切换剧透隐藏。", // 3.8 内容
    "misc.changelog.online_streaming": "在线流媒体", // 3.8/3.10 步骤标题
    "misc.changelog.online_content_38": "在线流媒体改用基于 HTTP/1 的新代理，并可自动轮换提供者直到找到可用的一项。", // 3.8 内容
    "misc.changelog.online_content_310": "在线流媒体现在可以代理字幕、刷新过期的源 URL，并更可靠地记住你的提供者、服务端、画质、音轨与字幕选择。", // 3.10 内容
    "misc.changelog.default_episode_source": "默认剧集来源", // 3.8 步骤标题
    "misc.changelog.episode_source_content": "选择进入动漫页面时默认打开的剧集来源。", // 3.8 内容
    "misc.changelog.redesigned_ui": "重新设计的界面设置", // 3.8 步骤标题
    "misc.changelog.ui_content": "用户界面设置面板已重新设计，导航更轻松。", // 3.8 内容
    "misc.changelog.route_preloading": "路由预加载", // 3.8 步骤标题
    "misc.changelog.preloading_content": "Seanime 现在可以在后台预加载路由，让导航更加流畅。你可以在新的界面设置面板中调整预加载行为。", // 3.8 内容
    "misc.changelog.ui_updates": "界面更新", // 3.8 步骤标题
    "misc.changelog.ui_updates_content": "媒体页头部经过了小幅重新设计，并新增了动画与过渡效果，体验更流畅。", // 3.8 内容
    "misc.changelog.extensions": "扩展", // 3.8 步骤标题
    "misc.changelog.extensions_content": "扩展现在可以在不卸载的情况下停用，插件也新增了设置、认证与扩展管理相关 API。", // 3.8 内容
    "misc.changelog.extension_secure_mode": "扩展安全模式", // 3.8 步骤标题
    "misc.changelog.secure_mode_content": "启用扩展安全模式后，每当扩展尝试执行敏感操作都会弹出确认提示。", // 3.8 内容
    "misc.changelog.denshi_window_state": "Denshi 窗口状态", // 3.8 步骤标题
    "misc.changelog.denshi_content": "Seanime Denshi 现在会记住窗口位置和大小，重新打开应用即可恢复之前的桌面布局。", // 3.8 内容
    "misc.changelog.new_built_in_player": "全新内置播放器（Denshi）", // 3.9 步骤标题
    "misc.changelog.denshi_player_content": "Denshi 现在内置了基于 libmpv 的播放器，可在应用视口内进行硬件加速渲染，编解码与字幕支持完善，并支持 mpv.conf 选项和着色器。", // 3.9 内容
    "misc.changelog.faster_torrent_streaming": "更快的种子流媒体", // 3.9 步骤标题
    "misc.changelog.faster_torrent_content": "种子流媒体启动速度最高提升 20%（取决于做种情况），下载进度报告更准确，并修复了批量选择问题。", // 3.9 内容
    "misc.changelog.faster_debrid_streaming": "更快的 Debrid 流媒体", // 3.9 步骤标题
    "misc.changelog.faster_debrid_content": "已缓存流的 Debrid 流媒体启动时间最高缩短 5 秒。", // 3.9 内容
    "misc.changelog.manga_source_refresh": "漫画源刷新", // 3.10 步骤标题
    "misc.changelog.manga_refresh_content": "漫画源刷新已改进。除更新已保存的源外，现在还可以为无源的漫画指定源，或重新评估所有源。", // 3.10 内容
    "misc.changelog.torrent_availability": "种子可用性", // 3.10 步骤标题
    "misc.changelog.availability_content": "启用后，在使用种子或 Debrid 流媒体时，会在媒体库缺失的近期剧集和「继续观看」中显示可用性徽章。", // 3.10 内容
    "misc.changelog.plugin_tray_badges": "插件托盘角标", // 3.10 步骤标题
    "misc.changelog.tray_badges_content": "即使插件未固定，托盘插件的角标数量现在也会显示在主托盘图标上。", // 3.10 内容
    "misc.changelog.export_mpvcore_logs": "导出 MpvCore 日志", // 3.10 步骤标题
    "misc.changelog.mpvcore_logs_content": "启用 MpvCore 日志后，现在可以从视频播放设置中导出日志，方便更快提交 Bug 报告。", // 3.10 内容

    // ===== 初始设置向导（app/(main)/_features/getting-started/getting-started-page.tsx，09-05）=====
    "misc.getting_started.alt_logo": "Seanime 标志", // Logo alt
    "misc.getting_started.library_title": "番剧媒体库", // h2
    "misc.getting_started.library_desc": "选择存放动漫文件的文件夹，Seanime 将在此扫描你的收藏。", // h2 说明
    "misc.getting_started.library_path_label": "番剧媒体库路径", // Field label
    "misc.getting_started.library_path_help": "选择存放动漫收藏的主文件夹，之后可以添加更多文件夹。", // Field help
    "misc.getting_started.denshi_note": "Seanime Denshi 内置了默认启用的媒体播放器，你仍可配置外部媒体播放器。", // Denshi 提示
    "misc.getting_started.media_player_title": "{external}媒体播放器", // h2
    "misc.getting_started.external_prefix": "外部 ", // media_player_title 的 external 插值（桌面端）
    "misc.getting_started.media_player_desc": "配置你想用于观看动漫并自动记录进度的外部媒体播放器。", // h2 说明
    "misc.getting_started.desktop_player_label": "桌面媒体播放器", // Field label（playback.atoms.tsx 复用）
    "misc.getting_started.desktop_player_or_transcode": "桌面媒体播放器或转码 / Direct Play", // playback.atoms.tsx
    "misc.getting_started.external_player_link": "外部播放器链接", // playback.atoms.tsx
    "misc.getting_started.mpv_help_win": "推荐使用 MPV 以获得更好的字幕渲染与种子串流体验。", // Field help
    "misc.getting_started.mpv_help_mac": "在 macOS 上推荐使用 MPV 或 IINA。", // Field help
    "misc.getting_started.mpv_recommended": "MPV（推荐）", // 下拉选项
    "misc.getting_started.mpv_install_hint": "在 Windows 上可用 Scoop 或 Chocolatey 轻松安装 MPV。在 macOS 上可用 Homebrew 安装 MPV。", // 提示
    "misc.getting_started.iina_alert": "要让 IINA 与 Seanime 正常配合，请在 IINA 通用设置中确保「所有窗口关闭后退出」被勾选，且「播放结束后保持窗口打开」被取消勾选。", // IINA 提示
    "misc.getting_started.vlc_config": "VLC 配置", // 配置标题
    "misc.getting_started.mpc_config": "MPC-HC 配置", // 配置标题
    "misc.getting_started.host": "主机", // Field label
    "misc.getting_started.port": "端口", // Field label
    "misc.getting_started.username": "用户名", // Field label
    "misc.getting_started.password": "密码", // Field label
    "misc.getting_started.vlc_path_label": "VLC 可执行文件路径", // Field label
    "misc.getting_started.mpc_path_label": "MPC-HC 可执行文件路径", // Field label
    "misc.getting_started.exe_path": "可执行文件路径", // Field label
    "misc.getting_started.downloading_title": "下载", // h2
    "misc.getting_started.downloading_desc": "配置用于下载的种子客户端。", // h2 说明
    "misc.getting_started.torrent_client": "种子客户端", // 标题
    "misc.getting_started.torrent_client_desc": "用于下载动漫种子的客户端", // 说明
    "misc.getting_started.client_label": "客户端", // Field label
    "misc.getting_started.qb_settings": "qBittorrent 设置", // 配置标题
    "misc.getting_started.transmission_settings": "Transmission 设置", // 配置标题
    "misc.getting_started.debrid_service": "Debrid 服务", // h2 / Field label
    "misc.getting_started.debrid_desc": "Debrid 服务提供更快的下载速度和云端即时串流。", // h2 说明
    "misc.getting_started.api_key_help": "由 Debrid 服务提供的 API 密钥。", // Field help（文案复用 player.prefs.api_key）
    "misc.getting_started.additional_features": "附加功能", // h2
    "misc.getting_started.additional_features_desc": "选择要启用的附加功能，之后可在设置中随时更改。", // h2 说明
    "misc.getting_started.feature_manga_desc": "阅读和下载漫画章节", // 功能说明
    "misc.getting_started.feature_torrent_streaming_desc": "无需等待下载即可串流种子", // 功能说明
    "misc.getting_started.feature_nsfw": "成人内容", // 功能标题
    "misc.getting_started.feature_nsfw_desc": "在媒体库和搜索中显示成人内容", // 功能说明
    "misc.getting_started.feature_online_streaming_desc": "观看在线来源的动漫", // 功能说明
    "misc.getting_started.feature_discord": "Discord Rich Presence", // 功能标题
    "misc.getting_started.feature_discord_desc": "在 Discord 上展示你正在观看的内容", // 功能说明
    "misc.getting_started.feature_transcode": "转码 / Direct Play", // 功能标题
    "misc.getting_started.feature_transcode_desc": "在其他设备上串流已下载的文件", // 功能说明

    // ===== Nakama 观影（_features/nakama/nakama-manager.tsx、nakama-watch-party-chat.tsx，09-05）=====
    "misc.nakama.reconnect_initiated": "已发起重新连接", // toast
    "misc.nakama.reconnect_failed": "重新连接失败：{message}", // toast
    "misc.nakama.cleanup_success": "已清理失效连接", // toast
    "misc.nakama.cleanup_failed": "清理失败：{message}", // toast
    "misc.nakama.party_created": "观影会已创建", // toast
    "misc.nakama.party_create_failed": "创建观影会失败：{message}", // toast
    "misc.nakama.joining_party": "正在加入观影会", // toast
    "misc.nakama.leaving_party": "正在退出观影会", // toast
    "misc.nakama.room_created": "房间创建成功", // toast
    "misc.nakama.disconnected": "已断开房间连接", // toast
    "misc.nakama.disconnect_failed": "断开房间连接失败：{message}", // toast
    "misc.nakama.currently_hosting": "正在主持", // Badge
    "misc.nakama.cleaning_up": "正在清理...", // 按钮
    "misc.nakama.remove_stale": "移除失效连接", // 按钮
    "misc.nakama.cloud_room": "云端房间", // 标题/Badge
    "misc.nakama.disconnecting": "正在断开...", // 按钮
    "misc.nakama.disconnect": "断开连接", // 按钮
    "misc.nakama.cloud_room_limitation": "云端房间不支持本地文件与 Debrid 播放。", // 说明
    "misc.nakama.host_url_passcode": "Nakama 主机地址与通行码", // 提示
    "misc.nakama.host_url": "主机地址", // TextInput leftAddon
    "misc.nakama.creating": "正在创建...", // 按钮
    "misc.nakama.create_cloud_room": "创建云端房间", // 按钮/确认标题
    "misc.nakama.create_cloud_room_desc": "继续操作即表示你同意在房间有效期间，通过 Seanime 的服务器广播你的播放状态以与同伴同步。每天限创建 10 个房间，每个房间限 4 名同伴（限制可能调整）。", // 确认描述
    "misc.nakama.auto_join_hint": "你将自动加入房间。", // Tooltip
    "misc.nakama.direct_connections": "直连（{count}）", // 标题
    "misc.nakama.no_connected_peers": "暂无已连接的同伴", // 空状态
    "misc.nakama.host_connection": "主机连接", // 标题
    "misc.nakama.connection_mode": "连接模式", // 标签
    "misc.nakama.direct": "直连", // Badge
    "misc.nakama.not_active": "Nakama 未启用", // 空状态
    "misc.nakama.not_active_hint": "在设置中配置 Nakama 以连接主机或开始主持", // 空状态
    "misc.nakama.watch_party": "观影会", // 标题
    "misc.nakama.create_watch_party": "创建观影会", // 按钮
    "misc.nakama.party_active_hint": "有一个进行中的观影会！加入即可一起同步观看。", // 说明
    "misc.nakama.joining": "正在加入...", // 按钮
    "misc.nakama.join_watch_party": "加入观影会", // 按钮
    "misc.nakama.join_hint": "连接到主机即可加入观影会", // 空状态
    "misc.nakama.no_active_party": "当前没有进行中的观影会", // 空状态
    "misc.nakama.enable_relay": "启用中继模式", // Tooltip
    "misc.nakama.leaving": "正在退出...", // 按钮
    "misc.nakama.leave": "退出", // 按钮
    "misc.nakama.participants": "参与者（{count}）", // 标题
    "misc.nakama.me": "（我）", // 参与者后缀
    "misc.nakama.relay": "中继", // Badge
    "misc.nakama.origin": "源", // Badge
    "misc.nakama.promote_origin": "提升为源", // 按钮
    "misc.nakama.buffer": "缓冲", // 标签
    "misc.nakama.buffer_health": "同步缓冲健康度", // Tooltip
    "misc.nakama.buffering": "缓冲中", // Badge
    "misc.nakama.chat_title": "观影会聊天", // nakama-watch-party-chat.tsx 标题
    "misc.nakama.no_messages": "暂无消息", // nakama-watch-party-chat.tsx 空状态
    "misc.nakama.chat_placeholder": "输入消息...", // nakama-watch-party-chat.tsx placeholder

    // ===== Sea Command（_features/sea-command/，09-05）=====
    "misc.sea_command.find_anime_collection": "在你的收藏中查找动漫", // sea-command-navigation.tsx
    "misc.sea_command.find_manga_collection": "在你的收藏中查找漫画", // sea-command-navigation.tsx
    "misc.sea_command.find_anime_library": "在你的媒体库中查找动漫", // sea-command-navigation.tsx
    "misc.sea_command.my_anime": "我的动漫", // sea-command-navigation.tsx heading
    "misc.sea_command.navigation": "导航", // sea-command-navigation.tsx heading
    "misc.sea_command.screens": "页面", // sea-command-navigation.tsx heading
    "misc.sea_command.go_to": "前往：", // sea-command-navigation.tsx
    "misc.sea_command.go_back": "返回", // sea-command-navigation.tsx
    "misc.sea_command.go_forward": "前进", // sea-command-navigation.tsx
    "misc.sea_command.enter": "回车", // CommandShortcut
    "misc.sea_command.advanced_search": "高级搜索", // sea-command-navigation.tsx 页面名
    "misc.sea_command.magnet_name": "磁力链接", // sea-command-torrent-magnet.tsx name
    "misc.sea_command.magnet_desc": "粘贴磁力链接以开始串流或下载。", // sea-command-torrent-magnet.tsx description
    "misc.sea_command.magnet_heading": "粘贴磁力链接", // sea-command-torrent-magnet.tsx heading
    "misc.sea_command.select_anime": "选择一个动漫", // sea-command-torrent-magnet.tsx heading
    "misc.sea_command.actions": "操作", // sea-command-actions.tsx heading
    "misc.sea_command.dropping": "正在丢弃...", // sea-command-actions.tsx
    "misc.sea_command.stream": "串流", // sea-command-torrent-magnet.tsx
    "misc.sea_command.stream_pick_file": "串流并选择文件", // sea-command-torrent-magnet.tsx
    "misc.sea_command.stream_debrid": "使用 Debrid 串流", // sea-command-torrent-magnet.tsx
    "misc.sea_command.stream_debrid_pick_file": "使用 Debrid 串流并选择文件", // sea-command-torrent-magnet.tsx
    "misc.sea_command.magnet_continue_hint": "粘贴有效的磁力链接以继续。", // sea-command-torrent-magnet.tsx
    "misc.sea_command.magnet_link_label": "磁力链接：", // sea-command-torrent-magnet.tsx
    "misc.sea_command.back": "返回", // sea-command-torrent-magnet.tsx 按钮
    "misc.sea_command.loading_anime": "正在加载你的动漫...", // sea-command-torrent-magnet.tsx
    "misc.sea_command.loading_episodes": "正在加载剧集...", // sea-command-torrent-magnet.tsx
    "misc.sea_command.no_episodes": "未找到剧集。", // sea-command-torrent-magnet.tsx
    "misc.sea_command.stream_unavailable": "该剧集无法串流，因为缺少 AniDB 映射。", // sea-command-torrent-magnet.tsx
    "misc.sea_command.download_unavailable": "无法下载，因为未能解析默认下载目录。", // sea-command-torrent-magnet.tsx
    "misc.sea_command.find_in_collection": "在你的收藏中查找", // sea-command-list.tsx
    "misc.sea_command.find_in_library": "在你的动漫媒体库中查找", // sea-command-list.tsx
    "misc.sea_command.search_anilist": "在 AniList 上搜索", // sea-command-list.tsx
    "misc.sea_command.stream_or_download": "通过磁力链接串流或下载", // sea-command-list.tsx
    "misc.sea_command.copy_logs_desc": "复制当前日志", // sea-command-list.tsx
    "misc.sea_command.record_issue": "记录一个问题", // sea-command-list.tsx、sea-command-actions.tsx
    "misc.sea_command.drop_torrent_desc": "丢弃当前的种子串流任务", // sea-command-list.tsx
    "misc.sea_command.reload_page": "重新加载页面", // sea-command-list.tsx、sea-command-actions.tsx
    "misc.sea_command.toggle_spoilers": "切换此动漫的剧透显示", // sea-command-list.tsx
    "misc.sea_command.copy_server_logs": "复制当前服务端日志", // sea-command-actions.tsx
    "misc.sea_command.drop_all_torrents": "丢弃种子串流客户端中的全部种子", // sea-command-actions.tsx
    "misc.sea_command.search_anime": "搜索动漫", // sea-command-search.tsx
    "misc.sea_command.search_manga": "搜索漫画", // sea-command-search.tsx
    "misc.sea_command.nothing_found": "未找到任何结果", // sea-command-search.tsx
    "misc.sea_command.nothing_found_desc": "没有找到与该名称匹配的内容，请重试。", // sea-command-search.tsx
    "misc.sea_command.unhide_spoilers": "取消隐藏当前动漫的剧透", // sea-command-spoilers.tsx
    "misc.sea_command.rehide_spoilers": "重新隐藏当前动漫的剧透", // sea-command-spoilers.tsx
    "misc.sea_command.input_placeholder": "输入命令或内容...", // sea-command.tsx placeholder
    "misc.anime.episode_image_alt": "剧集图片", // command-utils.tsx 等图片 alt
    "misc.anime.cover_image_alt": "封面图片", // playback 相关图片 alt
    "misc.anime.banner_image_alt": "横幅图片", // library-header.tsx 图片 alt
    "misc.anime.radial_shadow_alt": "径向阴影", // continue-watching-header.tsx 图片 alt

    // ===== 播放列表（_features/playlists/，09-05）=====
    "misc.playlist.add_episodes": "添加剧集", // playlist-editor.tsx 按钮
    "misc.playlist.episode_n": "第 {count} 集", // playlist-editor.tsx / global-playlist-manager.tsx
    "misc.playlist.no_episodes": "未找到剧集", // playlist-editor.tsx 空状态
    "misc.playlist.debrid_streaming": "Debrid 串流", // playlist-editor.tsx / debrid-stream-page.tsx
    "misc.playlist.search_placeholder": "搜索", // playlist-editor.tsx placeholder（复用 entry.torrent_search.search_placeholder 文案）
    "misc.playlist.update_progress_title": "更新进度？", // global-playlist-manager.tsx
    "misc.playlist.update_progress_desc": "要更新当前剧集的观看进度吗？", // global-playlist-manager.tsx
    "misc.playlist.no": "否", // global-playlist-manager.tsx 按钮
    "misc.playlist.reopen_episode": "重新打开剧集", // global-playlist-manager.tsx
    "misc.playlist.stop_playlist": "停止播放列表", // global-playlist-manager.tsx
    "misc.playlist.playing_episode": "正在播放《{anime}》第 {episode} 集", // global-playlist-manager.tsx toast
    "misc.playlist.name_label": "名称", // playlist-editor-modal.tsx（复用 home.option.name 文案）
    "misc.playlist.delete_playlist": "删除播放列表", // playlist-editor-modal.tsx
    "misc.playlist.add_to_new_playlist": "添加到新播放列表", // playlist-list-modal.tsx
    "misc.playlist.create_playlist": "创建播放列表", // playlist-list-modal.tsx
    "misc.playlist.no_playlists": "暂无播放列表", // playlist-list-modal.tsx 空状态

    // ===== 进度追踪（_features/progress-tracking/，09-05）=====
    "misc.progress.progress_auto_updated": "你的进度将自动更新", // playback-manager-progress-tracking.tsx
    "misc.progress.progress_updated": "进度已更新", // playback-manager-progress-tracking.tsx toast/提示
    "misc.progress.update_progress_now": "立即更新进度", // playback-manager / manual-progress-tracking.tsx
    "misc.progress.play_next_episode": "播放下一集", // 确认标题/按钮
    "misc.progress.play_next_confirm_desc": "确定要播放下一集吗？", // 确认描述
    "misc.progress.confirm": "确认", // 确认按钮
    "misc.progress.stop_playlist_confirm_desc": "确定要停止播放列表吗？列表将被删除。", // 确认描述
    "misc.progress.playlist_remaining": "此后还有 {count} 集", // 播放列表剩余提示
    "misc.progress.playing_externally": "正在外部播放", // manual-progress-tracking.tsx
    "misc.progress.open_in_external_player": "在外部播放器中打开", // manual-progress-tracking.tsx
    "misc.progress.next_episode": "下一集", // autoplay-countdown-modal.tsx
    "misc.progress.playing_next_in": "即将播放下一集", // autoplay-countdown-modal.tsx
    "misc.progress.play_now": "立即播放", // autoplay-countdown-modal.tsx

    // ===== 更新（_features/update/ 与 _electron/，09-05）=====
    "misc.update.available_toast_title": "有可用更新！", // update-modal.tsx / electron-update-modal.tsx h3
    "misc.update.available_menu": "有可用更新", // VerticalMenu item name
    "misc.update.update_now": "立即更新", // update-modal.tsx 按钮
    "misc.update.update_seanime": "更新 Seanime", // update-modal.tsx Modal 标题
    "misc.update.from_desktop": "请从桌面应用更新 Seanime。", // update-modal.tsx Alert
    "misc.update.update_note": "Seanime 将通过下载并替换现有文件来执行更新。详情请参阅文档。", // update-modal.tsx 说明
    "misc.update.download_and_install": "下载并安装", // update-modal.tsx / electron-update-modal.tsx 按钮
    "misc.update.see_on_github": "在 GitHub 上查看", // 按钮
    "misc.update.download_new_release": "下载新版本", // 下载弹窗标题
    "misc.update.select_destination": "选择保存位置", // 下载弹窗 label
    "misc.update.critical_patch": "此更新包含关键补丁", // update-helper.tsx
    "misc.update.whats_new": "更新内容", // update-helper.tsx
    "misc.update.other_updates": "你错过的其他更新", // update-helper.tsx
    "misc.update.desktop_unavailable": "此更新尚未支持桌面客户端。请稍等片刻或前往 GitHub 页面了解详情。", // electron-update-modal.tsx Alert
    "misc.update.install_now": "立即安装", // electron-update-modal.tsx
    "misc.update.installing": "正在安装...", // electron-update-modal.tsx
    "misc.update.downloading_pct": "正在下载... {count}%", // electron-update-modal.tsx
    "misc.update.update_error": "更新出错：{message}", // electron-update-modal.tsx toast
    "misc.update.install_failed": "安装更新失败：{message}", // electron-update-modal.tsx toast
    "misc.update.installed_closing": "更新已安装！应用即将关闭...", // electron-update-modal.tsx toast
    "misc.update.installed_restarting": "更新已安装。正在重启应用...", // electron-update-modal.tsx toast
    "misc.update.installed_restart_app": "更新已安装，请重启应用。", // electron-update-modal.tsx 提示
    "misc.electron.restart_server_desc": "后台服务进程已停止响应，请重启以继续。", // electron-restart-server-prompt.tsx
    "misc.electron.restart_server": "重启服务", // electron-restart-server-prompt.tsx 按钮
    "misc.electron.restart_server_hint": "如果多次重试后仍出现此提示，请重新启动应用程序。", // electron-restart-server-prompt.tsx

    // ===== 剧集组件（_features/anime/_components/，09-05）=====
    "misc.anime.torrent_available": "有可用种子", // episode-torrent-availability-badge.tsx
    "misc.anime.checking_torrents": "正在检查种子", // episode-torrent-availability-badge.tsx
    "misc.anime.waiting_torrent": "等待种子", // episode-torrent-availability-badge.tsx
    "misc.anime.availability_unknown": "可用性未知", // episode-torrent-availability-badge.tsx
    "misc.anime.torrent_found_title": "已找到匹配的种子", // episode-torrent-availability-badge.tsx
    "misc.anime.checking_provider_title": "正在检查所选种子提供者", // episode-torrent-availability-badge.tsx
    "misc.anime.no_torrent_title": "尚未找到匹配的种子", // episode-torrent-availability-badge.tsx
    "misc.anime.provider_failed_title": "无法检查种子提供者", // episode-torrent-availability-badge.tsx
    "misc.anime.filler": "原创剧情", // episode-grid-item.tsx
    "misc.anime.unidentified": "未识别", // episode-grid-item.tsx
    "misc.anime.no_metadata": "未找到元数据", // episode-grid-item.tsx
    "misc.anime.open_page": "打开页面", // episode-card.tsx
    "misc.anime.add_to_playlist": "添加到播放列表", // episode-card.tsx
    "misc.anime.no_trailer": "未找到预告片", // trailer-modal.tsx
    "misc.anime.youtube_video": "YouTube 视频", // trailer-modal.tsx
    "misc.anime.next_episode_prefix": "下一集：", // media-entry-card-components.tsx

    // ===== 播放器（video-core / mpv-core / media-core，09-05）=====
    "misc.player.font_size_help": "字体大小占视频高度的百分比", // video-core-settings-menu.tsx
    "misc.player.caption_font_help": "字幕使用的字体", // video-core-settings-menu.tsx
    "misc.player.font_family_title": "字体", // video-core-settings-menu.tsx
    "misc.player.outline_title": "描边", // video-core-settings-menu.tsx
    "misc.player.rendering": "渲染", // video-core-preferences.tsx 区块标题
    "misc.player.playback_section": "播放", // video-core-preferences.tsx 区块标题
    "misc.player.navigation_section": "导航", // video-core-preferences.tsx 区块标题
    "misc.player.audio_section": "音频", // video-core-preferences.tsx 区块标题
    "misc.player.beginning": "开头", // video-core-preferences.tsx overlay message
    "misc.player.capture_display": "捕获 / 显示", // video-core-stats.tsx
    "misc.player.prologue": "序幕", // video-core.utils.ts 章节类型
    "misc.player.opening": "片头", // video-core.utils.ts / mpv-core.tsx 章节类型
    "misc.player.ending": "片尾", // video-core.utils.ts / mpv-core.tsx 章节类型
    "misc.player.mistimed_delayed": "时序错位 / 延迟", // mpv-core-stats.tsx
    "misc.player.anime4k_stats": "Anime4K（{mode}）- {quality}", // mpv-core-stats.tsx
    "misc.player.custom_shaders_stats": "自定义（{count} 个启用）", // mpv-core-stats.tsx
    "misc.player.screenshot_folder": "截图文件夹", // screenshot-dir-modal.tsx
    "misc.player.screenshot_dir_desc": "选择保存视频截图的文件夹。", // screenshot-dir-modal.tsx
    "misc.player.video_playback_time": "视频播放时间", // media-core-control-bar.tsx aria-label
    "misc.player.seek_chapter_end": "跳转到第 {count} 章结尾", // media-core-control-bar.tsx
    "misc.player.playback_error": "播放错误", // media-core-overlays.tsx
    "misc.player.close_player": "关闭播放器", // media-core-overlays.tsx
    "misc.player.drawer_a11y_title": "抽屉面板", // media-core-drawer.tsx 隐藏标题
    "misc.player.caption_english": "英语", // video-core-media-captions.ts
    "misc.player.caption_spanish": "西班牙语", // video-core-media-captions.ts
    "misc.player.absolute_path_required": "必须是绝对路径", // screenshot-dir-modal.tsx
    "misc.player.preview_alt": "预览", // media-core-control-bar.tsx 缩略图 alt
    "misc.player.stream_error_fallback": "播放视频流时发生错误，请稍后重试。", // media-core-overlays.tsx
    "misc.tour.skip": "跳过", // tour-overlay.tsx
    "misc.tour.back": "上一步", // tour-overlay.tsx
    "misc.tour.next": "下一步", // tour-overlay.tsx
    "misc.tour.done": "完成", // tour-overlay.tsx
    "misc.entry.episode_section_specials": "特别篇", // episode-section.tsx
    "misc.entry.episode_section_others": "其他", // episode-section.tsx
    "misc.entry.episode_section_episodes": "剧集", // episode-section.tsx（Next Episode 复用 misc.progress.next_episode）
    "misc.entry.undownloaded_not_downloaded": "以下剧集尚未下载：", // undownloaded-episode-list.tsx
    "misc.entry.undownloaded_not_watched": "以下剧集既未观看也未下载：", // undownloaded-episode-list.tsx
    "misc.entry.undownloaded_not_in_library": "以下剧集不在你的媒体库中：", // undownloaded-episode-list.tsx
    "misc.entry.untitled": "未命名", // torrent-download-file-selection.tsx
    "misc.hooks.logs_deleted": "日志已删除", // status.hooks.ts
    "misc.hooks.logs_fetch_failed": "无法获取日志", // status.hooks.ts
    "misc.hooks.logs_copy_failed": "复制日志失败：{message}", // status.hooks.ts
    "misc.hooks.gc_completed": "垃圾回收完成", // status.hooks.ts
    "misc.exclusion.modal_title": "选择不参与共享的动漫", // media-exclusion-selector.tsx
    "misc.exclusion.select_all": "全选", // media-exclusion-selector.tsx
    "misc.exclusion.deselect_all": "取消全选", // media-exclusion-selector.tsx
    "misc.exclusion.select_adult": "选择成人内容", // media-exclusion-selector.tsx
    "misc.exclusion.selected_count": "已选 {count} 个（不会共享）", // media-exclusion-selector.tsx

    // ===== 其他散点（09-05）=====
    "misc.common.discover": "发现", // search/page.tsx、custom-sources/page.tsx
    "misc.common.library_watcher": "媒体库监视器", // library-watcher.tsx
    "misc.common.library_watcher_desc": "检测到媒体库有变更，请刷新你的条目。", // library-watcher.tsx
    "misc.common.scan_your_library": "扫描你的媒体库", // library-watcher.tsx
    "misc.common.library_scanned": "媒体库扫描完成", // library-watcher.tsx toast
    "misc.common.test_page": "测试", // main-sidebar.tsx 页面名
    "misc.common.command_palette": "命令面板", // main-sidebar.tsx
    "misc.common.manga_downloads": "漫画下载", // top-navbar.tsx
    "misc.common.library_explorer": "媒体库浏览器", // home-toolbar.tsx
    "misc.common.show_unread": "只显示未读", // offline-chapter-list.tsx
    "misc.common.refresh_anilist": "刷新 AniList", // refresh-anilist-button.tsx
    "misc.common.bulk_actions_note": "这些操作不会影响已忽略的文件。", // bulk-action-modal.tsx
    "misc.common.open_media_external": "在外部播放器中打开媒体", // external-player-link-button.tsx
    "misc.common.popular_shows": "热门节目", // discover/page.tsx
    "misc.common.trending_now": "当前热门", // discover/page.tsx
    "misc.common.choose_destination": "选择下载目录", // torrent-download-modal.tsx
    "misc.common.adult_filter": "成人", // torrent-search-container.tsx
    "misc.common.awaiting_stream": "等待串流", // debrid-stream-overlay.tsx
    "misc.common.downloading_torrent": "正在下载种子...", // debrid-stream-overlay.tsx
    "misc.common.add_to_library_hint": "将此动漫添加到你的媒体库以查看其剧集", // episode-section.tsx
    "misc.common.programming_language": "编程语言：", // extension-details.tsx
    "misc.common.manifest_url": "Manifest 地址：", // extension-details.tsx
    "misc.common.syncing": "同步进行中", // sync/page.tsx
    "misc.hooks.profile_generating": "正在生成 {profile} 性能剖析...", // status.hooks.ts
    "misc.hooks.http_error": "HTTP 错误，状态码：{status}", // status.hooks.ts
    "misc.hooks.profile_downloaded": "性能剖析「{profile}」已下载", // status.hooks.ts
    "misc.hooks.profile_download_failed": "下载 {profile} 性能剖析失败", // status.hooks.ts
    "misc.hooks.cpu_generating": "正在生成 CPU 性能剖析（{count} 秒）...", // status.hooks.ts
    "misc.hooks.cpu_downloaded": "CPU 性能剖析（{count} 秒）已下载", // status.hooks.ts
    "misc.hooks.cpu_download_failed": "下载 CPU 性能剖析失败", // status.hooks.ts
    "misc.issue.recorder_title": "问题记录器", // issue-report.tsx
    "misc.issue.start_recording": "开始录制", // issue-report.tsx
    "misc.issue.recording": "录制中", // issue-report.tsx
    "misc.issue.attach_screenshot": "附加截图", // issue-report.tsx
    "misc.issue.description_placeholder": "描述你遇到的问题...", // issue-report.tsx
    "misc.issue.cancel_recording": "取消录制", // issue-report.tsx
    "misc.issue.download_failed_msg": "下载问题报告失败", // issue-report.tsx
    "misc.issue.include_scanner_logs": "包含媒体库扫描日志", // issue-report.tsx
    "misc.issue.dom_recording_failed": "启动 DOM 录制失败", // issue-report.tsx
    "misc.issue.screenshot_capture_failed": "无法截取屏幕截图", // issue-report.tsx
    "misc.issue.caption_prompt": "为此截图添加说明（可选）：", // issue-report.tsx
    "misc.issue.screenshot_added": "截图已添加到报告", // issue-report.tsx
    "misc.issue.report_saved": "问题报告已保存", // issue-report.tsx
    "misc.issue.edit_note": "编辑备注", // issue-report.tsx
    "misc.issue.add_note": "添加备注", // issue-report.tsx
    "misc.issue.description_tooltip": "添加你遇到问题的描述", // issue-report.tsx
    "misc.issue.saving": "保存中...", // issue-report.tsx
    "misc.issue.stop_and_save": "停止并保存", // issue-report.tsx
    "misc.issue.events_count": "{count} 个事件", // issue-report.tsx
    "misc.issue.images_count": "{count} 张图片", // issue-report.tsx
    "misc.plugin.title": "插件", // plugin-sidebar-tray.tsx
    "misc.plugin.copy": "复制", // plugin-debug-window.tsx
    "misc.plugin.clear": "清空", // plugin-debug-window.tsx
    "misc.plugin.close": "关闭", // plugin-debug-window.tsx（复用 manga.refresh.action_close 文案）
    "misc.plugin.no_logs": "暂无日志", // plugin-debug-window.tsx
    "misc.plugin.search_logs": "搜索日志...", // plugin-debug-window.tsx
    "misc.plugin.something_went_wrong": "出错了：", // registry.tsx
    "misc.plugin.allow_extension": "允许扩展「{name}」执行以下操作：{action}？", // extension-prompt.tsx
    "misc.plugin.on_resource": "作用于「{resource}」", // extension-prompt.tsx
    "misc.plugin.no_metadata_hint": "此动漫没有可用的元数据信息，剧集映射可能不完整。", // plugin-entry-episode-tabs.tsx
    "misc.plugin.anonymous_extension": "某个扩展", // extension-prompt.tsx
    "misc.plugin.would_like_action": "想要执行以下操作", // extension-prompt.tsx
    "misc.plugin.dont_allow": "不允许", // extension-prompt.tsx
    "misc.plugin.allow": "允许", // extension-prompt.tsx

    // ===== 09-05 收官（misc）=====
    "misc.getting_started.step_library": "本地番剧库", // getting-started-page.tsx STEPS
    "misc.getting_started.step_library_desc": "选择你存放本地动漫文件的目录", // getting-started-page.tsx STEPS
    "misc.getting_started.step_player": "播放器配置", // getting-started-page.tsx STEPS
    "misc.getting_started.step_player_desc": "配置你的本地视频播放器", // getting-started-page.tsx STEPS
    "misc.getting_started.step_downloading": "下载设置", // getting-started-page.tsx STEPS
    "misc.getting_started.step_downloading_desc": "配置种子下载器与来源", // getting-started-page.tsx STEPS
    "misc.getting_started.step_debrid": "Debrid 云解析", // getting-started-page.tsx STEPS
    "misc.getting_started.step_debrid_desc": "可选的高级云端流媒体服务", // getting-started-page.tsx STEPS
    "misc.getting_started.step_features": "功能选项", // getting-started-page.tsx STEPS
    "misc.getting_started.step_features_desc": "启用扩展功能与模式", // getting-started-page.tsx STEPS
    "misc.getting_started.settings_note": "这些设置后续可随时在【系统设置】中更改", // getting-started-page.tsx 提示
    "misc.getting_started.launch_btn": "启动并进入 Seanime", // getting-started-page.tsx 提交按钮
    "misc.getting_started.by_author": "作者：5rahim", // getting-started-page.tsx 页脚署名
    "misc.getting_started.none_option": "无", // getting-started-page.tsx 下拉「None」选项
    "misc.nakama.reconnect": "重新连接", // nakama-manager.tsx 按钮
    "misc.nakama.reconnecting": "正在重新连接...", // nakama-manager.tsx 按钮
    "misc.nakama.passcode": "通行码", // nakama-manager.tsx TextInput leftAddon
    "misc.nakama.no_password_set": "未设置密码", // nakama-manager.tsx 空值兜底
    "misc.nakama.expires": "过期时间：", // nakama-manager.tsx 字段
    "misc.nakama.cloud_rooms_desc": "云端房间通过 Seanime 的 API 主持观影会，无需将你的服务器暴露到互联网。", // nakama-manager.tsx 说明
    "misc.common.library_refreshing": "正在刷新你的媒体库...", // library-watcher.tsx 扫描进度
    "misc.issue_report.filtered_events": "已过滤 {count} 个事件", // issue-report/page.tsx 时间线统计
    "misc.issue_report.screenshot_alt": "截图", // issue-report/page.tsx 图片 alt
    "misc.issue_report.entries": "{count} 条记录", // issue-report/page.tsx 控制台统计
    "misc.issue_report.screenshots_count": "截图（{count}）", // issue-report/page.tsx 面板标题
    "misc.issue_report.screenshot_n": "截图 {count}", // issue-report/page.tsx 截图编号
    "misc.issue_report.full_screenshot_alt": "完整截图", // issue-report/page.tsx 图片 alt
    "misc.issue_report.no_replay_data": "此报告中没有会话回放数据", // issue-report/page.tsx 空状态
    "misc.issue_report.timeline_count": "时间线（{count}）", // issue-report/page.tsx 回放时间线
    "misc.issue_report.waiting_events": "等待事件...", // issue-report/page.tsx 空状态
    "misc.issue_report.scan_logs_count": "扫描日志（{count}）", // issue-report/page.tsx 面板标题
    "misc.issue_report.first_is_earliest": "越靠前 = 越早", // issue-report/page.tsx 说明
    "misc.issue_report.scan_log_n": "扫描日志 {count}", // issue-report/page.tsx 抽屉标题/按钮
    "misc.docs.body": "请求体", // docs/page.tsx h5
    "misc.error.client_title": "客户端错误", // main.tsx LuffyError title
    "misc.error.unexpected": "Seanime 遇到了意外错误，请重试。", // main.tsx 错误页正文
    "misc.error.test_page": "错误测试", // routes/_main/error-test.tsx h1
    "misc.tour.loading": "加载中...", // tour-overlay.tsx 加载提示
} satisfies Dictionary
