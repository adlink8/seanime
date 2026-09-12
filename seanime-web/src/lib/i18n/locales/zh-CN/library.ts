import type { Dictionary } from "../../types"

/**
 * library 词条表
 * 归属：媒体库详情页（统计卡、筛选器）与媒体库排序选项
 */
export const libraryDictionary = {
    "library.stats.size": "媒体库大小", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.stats.file_count": "文件数", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.stats.entry_count": "条目数", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.stats.episodes": "剧集", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.stats.movies": "剧场版", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.stats.specials": "特别篇", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.filter.format": "格式", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.filter.format_all": "全部格式", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.filter.status": "状态", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.filter.status_all": "全部状态", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.filter.tags": "标签", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.filter.tags_all": "全部标签", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.filter.tags_empty": "未找到标签", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.filter.season": "季度", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.filter.season_all": "全部季度", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.filter.year": "年份", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.filter.year_any": "不限年份", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.filter.adult": "成人内容 (R18)", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.filter.all": "全部", // src/app/(main)/_features/anime-library/_screens/detailed-library-view.tsx
    "library.continue_watching.title": "继续观看", // src/app/(main)/_features/anime-library/_containers/continue-watching.tsx
    "library.sort.airdate_desc": "最近播出", // src/lib/helpers/filtering.ts
    "library.sort.airdate": "最早播出", // src/lib/helpers/filtering.ts
    "library.sort.episodes_desc": "最高集数", // src/lib/helpers/filtering.ts
    "library.sort.unwatched_desc": "未看集数最多", // src/lib/helpers/filtering.ts
    "library.sort.score_desc": "评分最高", // src/lib/helpers/filtering.ts
    "library.sort.title": "标题 (A-Z)", // src/lib/helpers/filtering.ts
    "library.sort.last_watched_desc": "最近观看", // src/lib/helpers/filtering.ts
    "library.sort.unread_chapters_desc": "未读章节最多", // src/lib/helpers/filtering.ts
    "library.sort.episodes": "最低集数", // src/lib/helpers/filtering.ts
    "library.sort.unwatched": "未看集数最少", // src/lib/helpers/filtering.ts
    "library.sort.score": "评分最低", // src/lib/helpers/filtering.ts
    "library.sort.start_date_desc": "最近开播", // src/lib/helpers/filtering.ts
    "library.sort.start_date": "最早开播", // src/lib/helpers/filtering.ts
    "library.sort.last_watched": "最早观看", // src/lib/helpers/filtering.ts
    "library.sort.title_desc": "标题 (Z-A)", // src/lib/helpers/filtering.ts
    "library.sort.audience_score_desc": "观众评分最高", // src/lib/helpers/filtering.ts
    "library.sort.audience_score": "观众评分最低", // src/lib/helpers/filtering.ts
    "library.sort.progress_desc": "观看进度最高", // src/lib/helpers/filtering.ts
    "library.sort.progress": "观看进度最低", // src/lib/helpers/filtering.ts
    "library.sort.added_desc": "最近添加/开始", // src/lib/helpers/filtering.ts
    "library.sort.added": "最早添加/开始", // src/lib/helpers/filtering.ts
    "library.sort.ended_desc": "最近完结", // src/lib/helpers/filtering.ts
    "library.sort.ended": "最早完结", // src/lib/helpers/filtering.ts
    "library.sort.released_desc": "最新发布", // src/lib/helpers/filtering.ts
    "library.sort.released": "最早发布", // src/lib/helpers/filtering.ts
    "library.sort.unread_chapters": "未读章节最少", // src/lib/helpers/filtering.ts
    "library.sort.airdate_updates_desc": "最近播出且有更新", // src/lib/helpers/filtering.ts
    "library.sort.airdate_updates": "最早播出且有更新", // src/lib/helpers/filtering.ts

    // ===== P5 / 05-03 追加：媒体库与库浏览器 =====

    // —— 跨界面复用 ——
    "library.common.lock_all_files": "锁定全部文件", // bulk-action-modal.tsx, library-explorer.tsx, toggle-lock-files-button.tsx
    "library.common.unlock_all_files": "解锁全部文件", // bulk-action-modal.tsx, library-explorer.tsx, toggle-lock-files-button.tsx
    "library.common.select_all_files": "全选文件", // unmatched-file-manager.tsx, ignored-file-manager.tsx
    "library.common.unmatched_files": "未匹配的文件", // unmatched-file-manager.tsx, library-explorer.tsx
    "library.common.unlocked_files": "已解锁的文件", // library-explorer.tsx
    "library.common.ignored_files": "已忽略的文件", // ignored-file-manager.tsx, library-explorer.tsx
    "library.common.unknown_media": "未知媒体", // library-explorer.tsx
    "library.common.save": "保存", // library-explorer.tsx
    "library.common.delete": "删除", // library-explorer.tsx
    "library.common.cancel": "取消", // library-explorer.tsx
    "library.common.file": "文件", // library-explorer.tsx
    "library.common.folder": "文件夹", // library-explorer.tsx

    // —— 文件类型选项 ——
    "library.file_type.main": "正片", // library-explorer.tsx, library-explorer-super-update-drawer.tsx
    "library.file_type.special": "特别篇", // library-explorer.tsx, library-explorer-super-update-drawer.tsx
    "library.file_type.nc": "NC", // library-explorer-super-update-drawer.tsx
    "library.file_type.nc_other": "NC/其他", // library-explorer.tsx

    // —— 可访问性替代文本 ——
    "library.alt.banner_image": "横幅图片", // _components/library-header.tsx, _containers/continue-watching-header.tsx
    "library.alt.cover_image": "封面图片", // _containers/continue-watching-header.tsx
    "library.alt.episode_image": "剧集图片", // _containers/continue-watching.tsx

    // —— 媒体库扫描器 ——
    "library.scanner.title": "媒体库扫描", // _containers/scanner-modal.tsx
    "library.scanner.refresh_library": "刷新媒体库", // _containers/scanner-modal.tsx
    "library.scanner.local_files": "本地文件", // _containers/scanner-modal.tsx
    "library.scanner.matching_data": "匹配数据", // _containers/scanner-modal.tsx
    "library.scanner.skip_locked_files": "跳过已锁定的文件", // _containers/scanner-modal.tsx
    "library.scanner.skip_ignored_files": "跳过已忽略的文件", // _containers/scanner-modal.tsx
    "library.scanner.anilist_only": "仅我的 AniList 收藏", // _containers/scanner-modal.tsx
    "library.scanner.anilist_only_help": "仅将本地文件与你的 AniList 收藏进行匹配。", // _containers/scanner-modal.tsx
    "library.scanner.anilist_only_more_help": "速度更快，但若你的收藏不包含媒体库中的所有番剧，通常准确度较低。", // _containers/scanner-modal.tsx
    "library.scanner.enhanced_method": "增强匹配方式", // _containers/scanner-modal.tsx
    "library.scanner.use_offline_db": "使用 Anime Offline Database", // _containers/scanner-modal.tsx
    "library.scanner.use_anilist_api": "使用 AniList API", // _containers/scanner-modal.tsx
    "library.scanner.enhanced_help": "将本地文件与整个 AniList 目录进行匹配。扫描会更慢。", // _containers/scanner-modal.tsx
    "library.scanner.slower_for_large": "大型媒体库会更慢", // _containers/scanner-modal.tsx
    "library.scanner.enhanced_more_help": "Seanime 会为媒体库中发现的每个番剧标题各发送一次 API 请求，可能导致速率限制并让扫描变慢。", // _containers/scanner-modal.tsx
    "library.scanner.simulated_warning": "在未登录 AniList 账号的情况下使用此功能并不推荐——如果你的媒体库很大，可能会触发速率限制并导致扫描变慢。建议使用账号以获得更好的体验。", // _containers/scanner-modal.tsx
    "library.scanner.scan": "扫描", // _containers/scanner-modal.tsx
    "library.scanner.scanning": "扫描中...", // _containers/scan-progress-bar.tsx

    // —— 批量操作 ——
    "library.bulk_actions.title": "批量操作", // _containers/bulk-action-modal.tsx
    "library.bulk_actions.remove_empty_dirs": "移除空目录", // _containers/bulk-action-modal.tsx
    "library.bulk_actions.remove_empty_dirs_desc": "此操作将移除媒体库中的所有空目录。确定要继续吗？", // _containers/bulk-action-modal.tsx
    "library.bulk_actions.files_locked": "文件已锁定", // _containers/bulk-action-modal.tsx
    "library.bulk_actions.files_unlocked": "文件已解锁", // _containers/bulk-action-modal.tsx

    // —— 未匹配文件管理 ——
    "library.unmatched.previous": "上一页", // _containers/unmatched-file-manager.tsx
    "library.unmatched.next": "下一页", // _containers/unmatched-file-manager.tsx
    "library.unmatched.match_selection": "匹配所选", // _containers/unmatched-file-manager.tsx
    "library.unmatched.fetch_suggestions": "获取推荐", // _containers/unmatched-file-manager.tsx
    "library.unmatched.search_on_anilist": "在 AniList 上搜索", // _containers/unmatched-file-manager.tsx
    "library.unmatched.ignore_selection": "忽略所选", // _containers/unmatched-file-manager.tsx
    "library.unmatched.select_anime": "选择番剧", // _containers/unmatched-file-manager.tsx
    "library.unmatched.open_on_anilist": "在 AniList 上打开", // _containers/unmatched-file-manager.tsx
    "library.unmatched.no_suggestions": "未找到推荐，请尝试手动搜索", // _containers/unmatched-file-manager.tsx
    "library.unmatched.files_ignored": "文件已忽略", // _containers/unmatched-file-manager.tsx
    "library.unmatched.label_type": "类型：", // _containers/unmatched-file-manager.tsx
    "library.unmatched.label_aired": "播出：", // _containers/unmatched-file-manager.tsx
    "library.unmatched.label_status": "状态：", // _containers/unmatched-file-manager.tsx

    // —— 已忽略文件管理 ——
    "library.ignored.unignore_selection": "取消忽略所选", // _containers/ignored-file-manager.tsx
    "library.ignored.empty": "没有已忽略的文件", // _containers/ignored-file-manager.tsx
    "library.ignored.files_unignored": "已取消忽略文件", // _containers/ignored-file-manager.tsx

    // —— 隐藏媒体 ——
    "library.unknown.title": "隐藏的媒体", // _containers/unknown-media-manager.tsx
    "library.unknown.description": "Seanime 将 {count} 组匹配到了你{source}收藏中不存在的剧集。", // _containers/unknown-media-manager.tsx
    "library.unknown.description_hint": "添加该媒体后即可在媒体库中看到条目；如果不正确，也可以取消匹配。", // _containers/unknown-media-manager.tsx
    "library.unknown.matched_to": "已匹配到", // _containers/unknown-media-manager.tsx
    "library.unknown.add_all_to_anilist": "全部添加到 AniList", // _containers/unknown-media-manager.tsx
    "library.unknown.add_all_to_collection": "全部添加到收藏", // _containers/unknown-media-manager.tsx
    "library.unknown.add_to_anilist": "添加到 AniList", // _containers/unknown-media-manager.tsx
    "library.unknown.add_to_collection": "添加到收藏", // _containers/unknown-media-manager.tsx
    "library.unknown.unmatch": "取消匹配", // _containers/unknown-media-manager.tsx
    "library.unknown.media_unmatched": "媒体已取消匹配", // _containers/unknown-media-manager.tsx

    // —— 继续观看头部 ——
    "library.continue_watching.go_to_episode": "跳转到第 {number} 集", // _containers/continue-watching-header.tsx
    "library.continue_watching.releasing_now": "正在放送", // _containers/continue-watching-header.tsx
    "library.continue_watching.preview": "预览", // _containers/continue-watching-header.tsx
    "library.continue_watching.ep_abbr": "集", // _containers/continue-watching.tsx

    // —— 自动下载器 ——
    "library.auto_downloader.title": "自动下载器", // _containers/anime-auto-downloader-button.tsx
    "library.auto_downloader.new_rule": "新建规则", // _containers/anime-auto-downloader-button.tsx
    "library.auto_downloader.no_rules": "此番剧暂无规则。", // _containers/anime-auto-downloader-button.tsx

    // —— 空媒体库 ——
    "library.empty.title": "媒体库为空", // _screens/empty-library-view.tsx
    "library.empty.scan_library": "扫描你的媒体库", // _screens/empty-library-view.tsx
    "library.empty.your_library_is_empty": "你的媒体库为空", // _screens/empty-library-view.tsx
    "library.empty.set_path": "设置本地媒体库路径并扫描", // _screens/empty-library-view.tsx
    "library.empty.include_online_stream": "在媒体库中包含在线流媒体", // _screens/empty-library-view.tsx
    "library.empty.include_torrent_stream": "在媒体库中包含种子流媒体", // _screens/empty-library-view.tsx
    "library.empty.include_debrid_stream": "在媒体库中包含 Debrid 流媒体", // _screens/empty-library-view.tsx
    "library.empty.trending_this_season": "本季热门", // _screens/empty-library-view.tsx

    // —— 详情页头部 / 筛选标签 ——
    "library.header.home": "首页", // _screens/detailed-library-view.tsx
    "library.header.host_library": "{name} 的媒体库", // _screens/detailed-library-view.tsx
    "library.search.sorting": "排序", // _screens/detailed-library-view.tsx
    "library.list.lists": "列表", // _screens/detailed-library-view.tsx

    // —— 库浏览器：工具条与筛选 ——
    "library.explorer.title": "媒体库浏览器", // library-explorer/library-explorer.tsx
    "library.explorer.filter": "筛选", // library-explorer/library-explorer.tsx
    "library.explorer.select": "选择", // library-explorer/library-explorer.tsx
    "library.explorer.selecting": "选择中", // library-explorer/library-explorer.tsx
    "library.explorer.search_placeholder": "搜索文件和文件夹...", // library-explorer/library-explorer.tsx
    "library.explorer.unscanned_alert": "部分文件尚未扫描。请先扫描媒体库，才能对它们执行操作。", // library-explorer/library-explorer.tsx
    "library.explorer.unmatched_count": "{count} 个未匹配的文件", // library-explorer/library-explorer.tsx
    "library.explorer.hidden_media_count": "{count} 个文件包含隐藏媒体", // library-explorer/library-explorer.tsx
    "library.explorer.lock_all_matched_alert": "锁定所有正确匹配的文件", // library-explorer/library-explorer.tsx
    "library.explorer.end": "已到底部", // library-explorer/library-explorer.tsx
    "library.explorer.small_screen_notice": "媒体库浏览器仅支持在较大屏幕上渲染。", // library-explorer/library-explorer-drawer.tsx

    // —— 库浏览器：右键菜单与批量操作 ——
    "library.explorer.preview_anime": "预览番剧", // library-explorer/library-explorer.tsx
    "library.explorer.resolve_unknown_media": "解析未知媒体", // library-explorer/library-explorer.tsx
    "library.explorer.super_update": "超级更新", // library-explorer/library-explorer.tsx, library-explorer-super-update.tsx
    "library.explorer.unignore_files": "取消忽略文件", // library-explorer/library-explorer.tsx
    "library.explorer.unmatch_files": "取消匹配文件", // library-explorer/library-explorer.tsx
    "library.explorer.match_files": "匹配文件", // library-explorer/library-explorer.tsx
    "library.explorer.ignore_files": "忽略文件", // library-explorer/library-explorer.tsx
    "library.explorer.edit_metadata": "编辑元数据", // library-explorer/library-explorer.tsx
    "library.explorer.delete_files": "删除文件", // library-explorer/library-explorer.tsx
    "library.explorer.open_in_explorer": "在文件管理器中打开", // library-explorer/library-explorer.tsx
    "library.explorer.more": "更多", // library-explorer/library-explorer.tsx
    "library.explorer.match_n_files": "匹配 {count} 个文件", // library-explorer/library-explorer.tsx
    "library.explorer.unmatch_n_files": "取消匹配 {count} 个文件", // library-explorer/library-explorer.tsx
    "library.explorer.ignore_n_files": "忽略 {count} 个文件", // library-explorer/library-explorer.tsx
    "library.explorer.unignore_n_files": "取消忽略 {count} 个文件", // library-explorer/library-explorer.tsx
    "library.explorer.delete_n_files": "删除 {count} 个文件", // library-explorer/library-explorer.tsx

    // —— 库浏览器：树节点状态与锁定 ——
    "library.explorer.root_name": "番剧媒体库", // library-explorer/library-explorer.tsx
    "library.explorer.unknown_anime": "未知番剧", // library-explorer/library-explorer.tsx
    "library.explorer.not_in_collection": "此媒体不在你的收藏中。", // library-explorer/library-explorer.tsx
    "library.explorer.ignored": "已忽略", // library-explorer/library-explorer.tsx
    "library.explorer.matched_files": "已匹配文件", // library-explorer/library-explorer.tsx
    "library.explorer.not_matched": "未匹配", // library-explorer/library-explorer.tsx
    "library.explorer.not_scanned": "未扫描", // library-explorer/library-explorer.tsx
    "library.explorer.matched": "已匹配", // library-explorer/library-explorer.tsx
    "library.explorer.lock": "锁定", // library-explorer/library-explorer.tsx
    "library.explorer.unlock": "解锁", // library-explorer/library-explorer.tsx
    "library.explorer.lock_all_desc": "这将锁定该目录下的所有文件。确定要继续吗？", // library-explorer/library-explorer.tsx
    "library.explorer.unlock_all_desc": "这将解锁该目录下的所有文件。确定要继续吗？", // library-explorer/library-explorer.tsx
    "library.explorer.no_media_found": "未找到媒体", // library-explorer/library-explorer.tsx
    "library.explorer.metadata_saved": "元数据已保存", // library-explorer/library-explorer.tsx

    // —— 库浏览器：元数据对话框 ——
    "library.explorer.file_metadata_title": "文件元数据", // library-explorer/library-explorer.tsx
    "library.explorer.episode_number": "集数", // library-explorer/library-explorer.tsx, library-explorer-super-update-drawer.tsx
    "library.explorer.episode_number_help": "相对集数。如果是剧场版，集数 = 1", // library-explorer/library-explorer.tsx
    "library.explorer.anidb_episode": "AniDB 剧集", // library-explorer/library-explorer.tsx, library-explorer-super-update-drawer.tsx
    "library.explorer.anidb_episode_help": "特别篇通常包含字母 S", // library-explorer/library-explorer.tsx
    "library.explorer.type": "类型", // library-explorer/library-explorer.tsx
    "library.explorer.file_type_label": "文件类型", // library-explorer-super-update-drawer.tsx

    // —— 库浏览器：详情面板 ——
    "library.explorer.select_file_or_folder": "选择一个文件或文件夹以查看详情", // library-explorer/library-explorer.tsx
    "library.explorer.file_suffix": "{ext} 文件", // library-explorer/library-explorer.tsx
    "library.explorer.path": "路径", // library-explorer/library-explorer.tsx
    "library.explorer.size": "大小", // library-explorer/library-explorer.tsx
    "library.explorer.associated_media": "关联媒体", // library-explorer/library-explorer.tsx
    "library.explorer.anime_series_count": "{count} 部番剧", // library-explorer/library-explorer.tsx
    "library.explorer.library_file": "媒体库文件", // library-explorer/library-explorer.tsx
    "library.explorer.anidb_ep": "AniDB 集", // library-explorer/library-explorer.tsx
    "library.explorer.contents": "内容", // library-explorer/library-explorer.tsx
    "library.explorer.contents_count": "{folders} 个文件夹，{files} 个文件", // library-explorer/library-explorer.tsx
    "library.explorer.delete_desc": "此操作无法撤销。", // library-explorer/library-explorer.tsx
    "library.explorer.select_files_to_delete": "选择要删除的文件", // library-explorer/library-explorer.tsx

    // —— 库浏览器：超级更新 ——
    "library.super_update.description": "一次性更新多个文件名和元数据。", // library-explorer-super-update-drawer.tsx
    "library.super_update.search_for": "查找内容", // library-explorer-super-update-drawer.tsx
    "library.super_update.search_placeholder": "输入要查找的文本...", // library-explorer-super-update-drawer.tsx
    "library.super_update.replace_with": "替换为", // library-explorer-super-update-drawer.tsx
    "library.super_update.replace_placeholder": "输入替换文本...", // library-explorer-super-update-drawer.tsx
    "library.super_update.enum_patterns": "枚举模式：", // library-explorer-super-update-drawer.tsx
    "library.super_update.enum_simple": "${} - 简单计数器 (0, 1, 2...)", // library-explorer-super-update-drawer.tsx
    "library.super_update.enum_start": "${start=5} - 从 5 开始 (5, 6, 7...)", // library-explorer-super-update-drawer.tsx
    "library.super_update.enum_increment": "${increment=2} - 每次递增 2 (0, 2, 4...)", // library-explorer-super-update-drawer.tsx
    "library.super_update.enum_padding": "${padding=3} - 用零填充 (000, 001, 002...)", // library-explorer-super-update-drawer.tsx
    "library.super_update.enum_padding_start": "${padding=3;start=10} - 组合 (010, 011, 012...)", // library-explorer-super-update-drawer.tsx
    "library.super_update.enum_padding_increment": "${padding=2;increment=5} - 填充 + 递增 (00, 05, 10...)", // library-explorer-super-update-drawer.tsx
    "library.super_update.enum_all": "${increment=2;start=1;padding=3} - 全部组合 (001, 003, 005...)", // library-explorer-super-update-drawer.tsx
    "library.super_update.use_regex": "使用正则表达式", // library-explorer-super-update-drawer.tsx
    "library.super_update.case_sensitive": "区分大小写", // library-explorer-super-update-drawer.tsx
    "library.super_update.match_all_occurrences": "匹配所有出现位置", // library-explorer-super-update-drawer.tsx
    "library.super_update.enumerate_items": "枚举项目", // library-explorer-super-update-drawer.tsx
    "library.super_update.text_formatting": "文本格式", // library-explorer-super-update-drawer.tsx
    "library.super_update.edit_file_metadata": "编辑文件元数据", // library-explorer-super-update-drawer.tsx
    "library.super_update.add_rule": "添加规则", // library-explorer-super-update-drawer.tsx
    "library.super_update.no_rules": "暂无元数据编辑规则。点击“添加规则”创建一个。", // library-explorer-super-update-drawer.tsx
    "library.super_update.original_count": "原始 ({count})", // library-explorer-super-update-drawer.tsx
    "library.super_update.renamed_count": "已重命名 ({count})", // library-explorer-super-update-drawer.tsx
    "library.super_update.apply_count": "应用 ({count})", // library-explorer-super-update-drawer.tsx
    "library.super_update.no_files_selected": "未选择文件", // library-explorer-super-update-drawer.tsx
    "library.super_update.metadata": "元数据", // library-explorer-super-update-drawer.tsx
    "library.super_update.no_change": "无变化", // library-explorer-super-update-drawer.tsx
    "library.super_update.renamed": "已重命名", // library-explorer-super-update-drawer.tsx
    "library.super_update.label_episode": "集数：", // library-explorer-super-update-drawer.tsx
    "library.super_update.label_anidb": "AniDB：", // library-explorer-super-update-drawer.tsx
    "library.super_update.label_type": "类型：", // library-explorer-super-update-drawer.tsx
    "library.super_update.success_prefix": "更新成功 ", // library-explorer-super-update-drawer.tsx
    "library.super_update.result_filenames": "{count} 个文件名", // library-explorer-super-update-drawer.tsx
    "library.super_update.result_metadata": "{count} 项元数据", // library-explorer-super-update-drawer.tsx
    "library.super_update.result_filenames_metadata": "{files} 个文件名和 {metadata} 项元数据", // library-explorer-super-update-drawer.tsx
    "library.super_update.applying_changes": "正在应用 {applied} 项更改，跳过 {skipped} 项无效更改", // library-explorer-super-update-drawer.tsx
    "library.super_update.no_valid_changes": "没有可应用的更改。请检查你的设置。", // library-explorer-super-update-drawer.tsx
    "library.super_update.failed_rename": "重命名文件失败：", // library-explorer-super-update-drawer.tsx
    "library.super_update.filter_optional": "筛选（可选）", // library-explorer-super-update-drawer.tsx
    "library.super_update.find_text_regex": "查找（文本/正则）", // library-explorer-super-update-drawer.tsx
    "library.super_update.search": "搜索", // library-explorer-super-update-drawer.tsx
    "library.super_update.replace": "替换", // library-explorer-super-update-drawer.tsx
    "library.super_update.select_type": "选择类型", // library-explorer-super-update-drawer.tsx
    "library.super_update.ph_episode_search": "例如：>=1;<=12;!=5;type=main|special;!type=nc", // library-explorer-super-update-drawer.tsx
    "library.super_update.ph_anidb_search": "输入文本或正则表达式", // library-explorer-super-update-drawer.tsx
    "library.super_update.ph_type_search": "例如：>=1;<12;=5;type=main;!type=special", // library-explorer-super-update-drawer.tsx
    "library.super_update.ph_anidb_filter": "例如：anidb>=1;anidb=S1;anidb!=C2;type=special", // library-explorer-super-update-drawer.tsx
    "library.super_update.ph_replace_episode": "例如：increment=1, decrement=1, start=1，或直接填写数值如 5", // library-explorer-super-update-drawer.tsx
    "library.super_update.ph_replace_anidb": "替换文本（支持枚举模式）", // library-explorer-super-update-drawer.tsx
    "library.super_update.help_search_episode": "运算符：>=, <=, >, <, =, !=, ! | 类型：type=main|special|nc, !type=special", // library-explorer-super-update-drawer.tsx
    "library.super_update.help_search_anidb": "支持正则表达式和大小写敏感选项", // library-explorer-super-update-drawer.tsx
    "library.super_update.help_search_type": "与集数使用相同的运算符：>=, <=, >, <, =, !=, ! | 类型：type=main, !type=special", // library-explorer-super-update-drawer.tsx
    "library.super_update.help_anidb_operators": "AniDB 运算符：anidb>=, anidb<=, anidb=, anidb!=, !anidb= | 格式：数字 (1,12) 或带前缀 (S1,C2,T1)", // library-explorer-super-update-drawer.tsx

    // —— 其他 ——
    "library.collection.show_all": "显示全部", // _containers/library-collection.tsx
    "library.collection.show_unwatched_only": "仅显示未观看", // _containers/library-collection.tsx
    "library.play_random": "随机播放番剧", // _containers/play-random-episode-button.tsx

    // —— 资源管理器：详情面板与删除弹窗（library-explorer.tsx）——
    // 注意：本区块词条已全部并入上方 151-218 行的既有 explorer 区块，勿在此重复添加
} satisfies Dictionary
