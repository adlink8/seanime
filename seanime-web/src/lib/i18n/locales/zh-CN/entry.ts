import type { Dictionary } from "../../types"

/**
 * entry 词条表
 * 归属：条目详情页（app/(main)/entry/）
 * 独占：07-04
 *
 * 只读复用：common.* / player.* / media.* 等同名文案按契约复用，不新建重复词条。
 */
export const entryDictionary = {
    // ===== 种子搜索容器（torrent-search-container.tsx）=====
    "entry.torrent_search.select_torrent_to_stream": "选择要串流的种子", // torrent-search-container.tsx
    "entry.torrent_search.provider_none": "无", // torrent-search-container.tsx
    "entry.torrent_search.smart_search": "智能搜索", // torrent-search-container.tsx
    "entry.torrent_search.smart_search_help": "根据给定参数自动搜索。", // torrent-search-container.tsx
    "entry.torrent_search.smart_search_unsupported": "该提供商不支持智能搜索。", // torrent-search-container.tsx
    "entry.torrent_search.search_across_providers": "跨源搜索", // torrent-search-container.tsx
    "entry.torrent_search.search_across_providers_help": "使用已保存的附加提供商运行相同搜索。", // torrent-search-container.tsx
    "entry.torrent_search.no_providers_found": "未找到提供商", // torrent-search-container.tsx
    "entry.torrent_search.add_providers": "添加提供商", // torrent-search-container.tsx
    "entry.torrent_search.additional_providers": "附加提供商", // torrent-search-container.tsx
    "entry.torrent_search.episode_number": "集数", // torrent-search-container.tsx, episode-item.tsx
    "entry.torrent_search.batches": "批量打包", // torrent-search-container.tsx
    "entry.torrent_search.resolution": "分辨率", // torrent-search-container.tsx, torrent-common-helpers.tsx
    "entry.torrent_search.resolution_any": "任意", // torrent-search-container.tsx
    "entry.torrent_search.best_releases": "最佳发布", // torrent-search-container.tsx
    "entry.torrent_search.refine_title": "精炼标题（{title}）", // torrent-search-container.tsx
    "entry.torrent_search.search_placeholder": "搜索", // torrent-search-container.tsx
    "entry.torrent_search.warn_adult_unsupported": "该提供商不支持成人内容", // torrent-search-container.tsx
    "entry.torrent_search.warn_smart_search_unsupported": "该提供商不支持智能搜索", // torrent-search-container.tsx
    "entry.torrent_search.warn_best_release_unsupported": "该提供商不支持最佳发布搜索", // torrent-search-container.tsx
    "entry.torrent_search.search_failed": "搜索失败", // torrent-search-container.tsx
    "entry.torrent_search.fetch_failed": "从提供商获取种子失败。", // torrent-search-container.tsx
    "entry.torrent_search.retry_search": "重试搜索", // torrent-search-container.tsx
    "entry.torrent_search.checking_provider_again": "正在再次检查提供商...", // torrent-search-container.tsx
    "entry.torrent_search.auto_retry_countdown": "未找到匹配的种子，将在 {time} 后重新检查。", // torrent-search-container.tsx
    "entry.torrent_search.no_extension_selected": "未选择扩展", // torrent-search-container.tsx
    "entry.torrent_search.add_extensions": "添加扩展", // torrent-search-container.tsx
    "entry.torrent_search.previous_selection": "上次选择", // torrent-search-container.tsx

    // ===== 倒计时时间（torrent-search-container.tsx formatRetryTime）=====
    "entry.time.moment": "片刻", // torrent-search-container.tsx
    "entry.time.seconds": "{seconds} 秒", // torrent-search-container.tsx
    "entry.time.minutes_seconds": "{minutes} 分 {seconds} 秒", // torrent-search-container.tsx

    // ===== 下载确认弹窗（torrent-download-modal.tsx）=====
    "entry.torrent_download.choose_destination": "选择下载位置", // torrent-download-modal.tsx
    "entry.torrent_download.with_debrid": "使用 Debrid 服务下载", // torrent-download-modal.tsx
    "entry.torrent_download.destination": "下载位置", // torrent-download-modal.tsx, torrent-download-file-selection.tsx
    "entry.torrent_download.download_torrent_files": "下载 .torrent 文件", // torrent-download-modal.tsx
    "entry.torrent_download.download_all": "下载全部", // torrent-download-modal.tsx
    "entry.torrent_download.download": "下载", // torrent-download-modal.tsx, anime-entry-download-files-modal.tsx, torrent-search-button.tsx
    "entry.torrent_download.download_with_torrent_client": "使用种子客户端下载", // torrent-download-modal.tsx
    "entry.torrent_download.choose_files": "选择要下载的文件", // torrent-download-modal.tsx
    "entry.torrent_download.download_missing_episodes": "下载缺失剧集", // torrent-download-modal.tsx
    "entry.torrent_download.stream": "串流", // torrent-download-modal.tsx, debrid/torrent-stream-file-selection-modal.tsx

    // ===== 种子预览条目（torrent-preview-item.tsx / torrent-table.tsx）=====
    "entry.torrent_preview.highest_quality": "最高画质", // torrent-preview-item.tsx
    "entry.torrent_preview.episode_n": "第 {n} 集", // torrent-preview-item.tsx, torrent-table.tsx, playback-play-pill.tsx
    "entry.torrent_preview.batch": "打包", // torrent-preview-item.tsx
    "entry.torrent_preview.part_and": "第 {a} 部和第 {b} 部", // torrent-preview-item.tsx
    "entry.torrent_preview.parts_range": "第 {a} 至第 {b} 部", // torrent-preview-item.tsx
    "entry.torrent_preview.part_n": "第 {a} 部", // torrent-preview-item.tsx
    "entry.torrent_preview.season_and": "第 {a} 季和第 {b} 季", // torrent-preview-item.tsx
    "entry.torrent_preview.seasons_range": "第 {a} 至第 {b} 季", // torrent-preview-item.tsx
    "entry.torrent_preview.season_n": "第 {a} 季", // torrent-preview-item.tsx
    "entry.torrent_preview.episodes_range": "第 {a} 至第 {b} 集", // torrent-preview-item.tsx
    "entry.torrent_preview.season_suffix": "（第 {n} 季）", // torrent-preview-item.tsx
    "entry.torrent_preview.episode_image_alt": "剧集图片", // torrent-preview-item.tsx
    "entry.torrent_preview.unidentified": "无法识别", // torrent-preview-item.tsx
    "entry.torrent_preview.open_in_browser": "在浏览器中打开", // torrent-preview-item.tsx, torrent-download-modal.tsx
    "entry.torrent_table.nothing_found": "未找到任何内容", // torrent-table.tsx, torrent-preview-list.tsx

    // ===== 种子筛选与排序（torrent-common-helpers.tsx / torrent-item-badges.tsx）=====
    "entry.torrent_filter.result_count": "{count} 个结果", // torrent-common-helpers.tsx
    "entry.torrent_filter.filters": "筛选", // torrent-common-helpers.tsx
    "entry.torrent_filter.filters_notice": "筛选基于种子名称，可能会遗漏部分结果。", // torrent-common-helpers.tsx
    "entry.torrent_filter.multi_subs": "多字幕", // torrent-common-helpers.tsx, torrent-item-badges.tsx
    "entry.torrent_filter.dubbed": "配音", // torrent-common-helpers.tsx, torrent-item-badges.tsx
    "entry.torrent_filter.languages": "语言", // torrent-item-badges.tsx
    "entry.torrent_filter.original_and_dub": "原声+配音", // torrent-item-badges.tsx
    "entry.torrent_filter.seeders": "做种数", // torrent-common-helpers.tsx
    "entry.torrent_filter.seeders_none": "无", // torrent-item-badges.tsx
    "entry.torrent_filter.seeders_unit": "做种", // torrent-item-badges.tsx
    "entry.torrent_filter.date": "日期", // torrent-common-helpers.tsx
    "entry.torrent_filter.debrid_instant": "Debrid 服务即时可用", // torrent-item-badges.tsx

    // ===== 下载按钮（torrent-search-button.tsx）=====
    "entry.torrent_search_button.download_episodes": "下载{range}{n} 集", // torrent-search-button.tsx
    "entry.torrent_search_button.range_batch": "批次 / ", // torrent-search-button.tsx
    "entry.torrent_search_button.range_next": "接下来的 ", // torrent-search-button.tsx
    "entry.torrent_search_button.download_movie": "下载电影", // torrent-search-button.tsx

    // ===== 下载文件选择（torrent-download-file-selection.tsx / anime-entry-download-files-modal.tsx）=====
    "entry.download_files.select_files_to_download": "选择要下载的文件", // torrent-download-file-selection.tsx, anime-entry-download-files-modal.tsx
    "entry.download_files.files_selected_count": "已选择 {selected}/{total} 个文件", // torrent-download-file-selection.tsx
    "entry.download_files.download_selected": "下载所选文件", // torrent-download-file-selection.tsx
    "entry.download_files.popup_notice": "Seanime 将为每个下载的文件打开一个新标签页。请确保浏览器允许弹出窗口。", // anime-entry-download-files-modal.tsx

    // ===== Debrid 串流遮罩（debrid-stream-overlay.tsx）=====
    "entry.debrid_stream.sending_to_player": "正在将串流发送至播放器...", // debrid-stream-overlay.tsx
    "entry.debrid_stream.cancel_and_remove_torrent": "取消并移除种子", // debrid-stream-overlay.tsx
    "entry.debrid_stream.cancel_and_remove_torrent_desc": "确定要取消串流并移除该种子吗？", // debrid-stream-overlay.tsx
    "entry.debrid_stream.cancel_stream": "取消串流", // debrid-stream-overlay.tsx
    "entry.debrid_stream.cancel_stream_desc": "确定要取消串流吗？", // debrid-stream-overlay.tsx
    "entry.debrid_stream.awaiting_stream": "正在等待 Debrid 服务的串流", // debrid-stream-overlay.tsx
    "entry.debrid_stream.close_wont_cancel": "关闭此弹窗不会取消串流", // debrid-stream-overlay.tsx

    // ===== 串流页面（debrid-stream-page.tsx / torrent-stream-page.tsx 共用）=====
    "entry.stream_page.disable_previous_torrent": "禁用上次种子", // debrid-stream-page.tsx, torrent-stream-page.tsx
    "entry.stream_page.disable_previous_torrent_desc": "暂时禁用已保存的上次批量下载，或彻底删除保存的历史记录。", // debrid-stream-page.tsx, torrent-stream-page.tsx
    "entry.stream_page.delete_history": "删除历史记录", // debrid-stream-page.tsx, torrent-stream-page.tsx
    "entry.stream_page.disable_only": "仅禁用", // debrid-stream-page.tsx, torrent-stream-page.tsx
    "entry.stream_page.auto_select": "自动选择", // debrid-stream-page.tsx, torrent-stream-page.tsx
    "entry.stream_page.auto_select_file": "自动选择文件", // debrid-stream-page.tsx, torrent-stream-page.tsx
    "entry.stream_page.auto_select_file_help": "将从你选择的批量种子中自动选择剧集文件", // debrid-stream-page.tsx, torrent-stream-page.tsx
    "entry.stream_page.autoselecting_previous": "正在上次的种子中自动选择", // debrid-stream-page.tsx, torrent-stream-page.tsx
    "entry.stream_page.no_metadata_info": "该番剧没有可用的元数据信息。你可能需要手动选择要串流的文件。", // debrid-stream-page.tsx, torrent-stream-page.tsx

    // ===== 串流文件选择（debrid/torrent-stream-file-selection-modal.tsx 共用）=====
    "entry.stream_file_selection.launching": "正在启动串流...", // debrid-stream-file-selection-modal.tsx, torrent-stream-file-selection-modal.tsx
    "entry.stream_file_selection.fetching_info": "正在获取种子信息...", // debrid-stream-file-selection-modal.tsx, torrent-stream-file-selection-modal.tsx

    // ===== 种子串流（playback-play-pill.tsx / torrent-stream-overlay.tsx / torrent-stream-button.tsx）=====
    "entry.torrent_stream.stop_streaming": "停止串流？", // playback-play-pill.tsx
    "entry.torrent_stream.stop_streaming_desc": "确定要停止并关闭串流吗？", // playback-play-pill.tsx
    "entry.torrent_stream.stop_stream": "停止串流", // playback-play-pill.tsx, torrent-stream-overlay.tsx
    "entry.torrent_stream.loading": "加载中...", // playback-play-pill.tsx, torrent-stream-overlay.tsx
    "entry.torrent_stream.selecting_file": "正在选择文件...", // playback-play-pill.tsx, torrent-stream-overlay.tsx
    "entry.torrent_stream.adding_torrent": "正在添加种子“{name}”", // playback-play-pill.tsx, torrent-stream-overlay.tsx
    "entry.torrent_stream.adding_torrent_ellipsis": "正在添加种子...", // playback-play-pill.tsx
    "entry.torrent_stream.checking_torrent": "正在检查种子“{name}”", // playback-play-pill.tsx, torrent-stream-overlay.tsx
    "entry.torrent_stream.checking_torrent_ellipsis": "正在检查种子...", // playback-play-pill.tsx
    "entry.torrent_stream.getting_metadata": "正在获取元数据...", // torrent-stream-overlay.tsx
    "entry.torrent_stream.sending_to_player": "正在将串流发送至播放器...", // playback-play-pill.tsx, torrent-stream-overlay.tsx
    "entry.torrent_stream.active_streaming": "正在串流", // playback-play-pill.tsx
    "entry.torrent_stream.top_candidates": "最佳候选", // playback-play-pill.tsx
    "entry.torrent_stream.close_view": "关闭种子串流", // torrent-stream-button.tsx
    "entry.torrent_stream.view": "种子串流", // torrent-stream-button.tsx

    // ===== 下拉菜单（anime-entry-dropdown-menu.tsx）=====
    "entry.dropdown.open_directory": "打开目录", // anime-entry-dropdown-menu.tsx
    "entry.dropdown.open_on_anidb": "在 AniDB 打开", // anime-entry-dropdown-menu.tsx, episode-item.tsx
    "entry.dropdown.open_on_mal": "在 MAL 打开", // anime-entry-dropdown-menu.tsx
    "entry.dropdown.copy_id": "复制 ID", // anime-entry-dropdown-menu.tsx
    "entry.dropdown.bulk_actions": "批量操作", // anime-entry-dropdown-menu.tsx
    "entry.dropdown.download_some_files": "下载部分文件", // anime-entry-dropdown-menu.tsx
    "entry.dropdown.unmatch_some_files": "取消匹配部分文件", // anime-entry-dropdown-menu.tsx
    "entry.dropdown.delete_some_files": "删除部分文件", // anime-entry-dropdown-menu.tsx

    // ===== 批量弹窗（anime-entry-unmatch-files-modal.tsx / anime-entry-bulk-delete-files-modal.tsx）=====
    "entry.unmatch.select_files_to_unmatch": "选择要取消匹配的文件", // anime-entry-unmatch-files-modal.tsx
    "entry.bulk_delete.cannot_be_undone": "此操作无法撤销。", // anime-entry-bulk-delete-files-modal.tsx

    // ===== 元数据管理器（anime-entry-metadata-manager.tsx）=====
    "entry.metadata.fetching": "正在获取...", // anime-entry-metadata-manager.tsx
    "entry.metadata.fetch_filler_info": "获取 filler 信息", // anime-entry-metadata-manager.tsx
    "entry.metadata.removing": "正在移除...", // anime-entry-metadata-manager.tsx
    "entry.metadata.remove_filler_info": "移除 filler 信息", // anime-entry-metadata-manager.tsx
    "entry.metadata.parent": "元数据父级", // anime-entry-metadata-manager.tsx
    "entry.metadata.parent_help": "如果父系列不包含特别篇元数据，此操作将无效。", // anime-entry-metadata-manager.tsx
    "entry.metadata.parent_desc": "通过将条目关联到父级番剧，为特别篇添加元数据。", // anime-entry-metadata-manager.tsx
    "entry.metadata.anilist_id": "AniList ID", // anime-entry-metadata-manager.tsx
    "entry.metadata.select_parent": "选择父级", // anime-entry-metadata-manager.tsx
    "entry.metadata.special_offset": "特别篇偏移", // anime-entry-metadata-manager.tsx
    "entry.metadata.special_offset_hint": "0 = 第 1 季，1 = 第 2 季…", // anime-entry-metadata-manager.tsx
    "entry.metadata.remove": "移除", // anime-entry-metadata-manager.tsx

    // ===== 剧集条目与剧集区（episode-item.tsx / episode-section.tsx / undownloaded-episode-list.tsx / torrent-stream-episode-section.tsx）=====
    "entry.episode.play_externally": "外部播放", // episode-item.tsx, episode-section.tsx, torrent-stream-episode-section.tsx
    "entry.episode.copy_stream_url": "复制串流地址", // episode-item.tsx
    "entry.episode.stream_url_copied": "串流地址已复制", // episode-item.tsx
    "entry.episode.episode_number_help": "相对集数。若为剧场版，集数 = 1", // episode-item.tsx
    "entry.episode.anidb_episode": "AniDB 剧集", // episode-item.tsx
    "entry.episode.anidb_episode_help": "特别篇通常包含字母 S", // episode-item.tsx
    "entry.episode.type": "类型", // episode-item.tsx
    "entry.episode.type_main": "正篇", // episode-item.tsx
    "entry.episode.type_special": "特别篇", // episode-item.tsx
    "entry.episode.type_nc_other": "NC/其他", // episode-item.tsx
    "entry.episode.update_metadata": "更新元数据", // episode-item.tsx
    "entry.episode.banner_alt": "横幅", // episode-item.tsx
    "entry.episode.minutes": "分钟", // episode-item.tsx
    "entry.episode.no_summary": "暂无简介", // episode-item.tsx
    "entry.episode.anidb_episode_label": "AniDB 剧集：{episode}", // episode-item.tsx
    "entry.episode.not_in_nakama_library": "不在 Nakama 好友的媒体库中", // episode-section.tsx
    "entry.episode.not_in_your_library": "不在你的媒体库中", // episode-section.tsx
    "entry.episode.not_yet_released": "尚未放送", // episode-section.tsx
    "entry.episode.invalid_episodes_notice": "部分剧集无效。请更新元数据以修复。", // episode-section.tsx
    "entry.episode.specials_header": "特别篇 (SP)", // episode-section.tsx（原裸中文回填）
    "entry.episode.others_header": "其他剧集", // episode-section.tsx（原裸中文回填）
    "entry.episode.aired_on": "放送于 {date}", // undownloaded-episode-list.tsx
    "entry.episode.aired": "放送日期", // undownloaded-episode-list.tsx
    "entry.episode.and_more": "更多...", // undownloaded-episode-list.tsx

    // ===== 关联与推荐（relations-recommendations-section.tsx）=====
    "entry.relations.header": "关联作品", // relations-recommendations-section.tsx（原裸中文回填）
    "entry.relations.manga": "漫画", // relations-recommendations-section.tsx
    "entry.relations.movie_relation": "{relation}（剧场版）", // relations-recommendations-section.tsx
    "entry.relations.recommendations": "相关推荐", // relations-recommendations-section.tsx（原裸中文回填）

    // ===== 页面视图切换（anime-entry-page.tsx / 各 view 按钮）=====
    "entry.page.downloaded_episodes": "已下载剧集", // anime-entry-page.tsx
    "entry.page.torrent_streaming": "种子串流", // anime-entry-page.tsx
    "entry.page.debrid_streaming": "Debrid 串流", // anime-entry-page.tsx, debrid-stream-button.tsx
    "entry.page.online_streaming": "在线串流", // anime-entry-page.tsx, anime-onlinestream-button.tsx
    "entry.page.close_online_streaming": "关闭在线串流", // anime-onlinestream-button.tsx
    "entry.page.close_debrid_streaming": "关闭 Debrid 串流", // debrid-stream-button.tsx
    "entry.page.views": "视图", // anime-entry-page.tsx
    "entry.page.download_torrents": "下载种子", // anime-entry-page.tsx
    "entry.page.local_library": "本地媒体库", // anime-entry-page.tsx

    // ===== 元信息区（meta-section.tsx / anime-entry-silence-toggle.tsx）=====
    "entry.meta.custom_source": "自定义来源", // meta-section.tsx
    "entry.meta.open_in_website": "在网站打开", // meta-section.tsx
    "entry.meta.trailer": "预告片", // meta-section.tsx
    "entry.meta.shared_by": "由 {user} 分享", // meta-section.tsx
    "entry.meta.inaccurate_schedule_1": "无法获取该作品的准确放送时间信息。", // meta-section.tsx
    "entry.meta.inaccurate_schedule_2": "请在线查询放送时间表以获取更多信息。", // meta-section.tsx
    "entry.meta.no_metadata_message": "该条目暂不支持获取剧集元数据。", // meta-section.tsx
    "entry.meta.silence_notifications": "静音通知", // anime-entry-silence-toggle.tsx
    "entry.meta.unsilence_notifications": "取消静音通知", // anime-entry-silence-toggle.tsx
} satisfies Dictionary
