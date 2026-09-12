import type { Dictionary } from "../../types"

/**
 * settings-media 词条表（settings 域第二文件，与 settings.ts 按 key 判重互补）
 * 归属：设置页·媒体与流（torrentstream / mediastream / manga / nakama /
 *       anime-library / playback / mediaplayer / autoselect-profile）
 * 独占：07-02
 *
 * 只读复用（不新建重复词条）：
 * - 播放器既有词条见 `player.ts`（playback/mediaplayer-settings 的 Phase 2 种子）
 * - `common.*` / `library.*` / `manga.*`（P8 建）等同名文案按契约复用
 */
export const settingsMediaDictionary = {
    // ===== playback-settings.tsx（新增部分；页头/模式卡复用 player.* 种子词条）=====
    "settings.playback.updated": "播放设置已更新", // playback-settings.tsx toast
    "settings.playback.no_scheme": "尚未设置外部播放器自定义协议", // playback-settings.tsx Alert
    "settings.playback.denshi_title": "Seanime Denshi", // playback-settings.tsx 卡片标题（专有名词保留）
    "settings.playback.use_builtin": "使用内置播放器", // playback-settings.tsx
    "settings.playback.use_builtin_help": "启用后，所有媒体播放将使用内置播放器（覆盖下方设置）", // playback-settings.tsx
    "settings.playback.engine_title": "内置播放器引擎", // playback-settings.tsx
    "settings.playback.engine_desc": "选择 Denshi 集成播放所使用的渲染器。", // playback-settings.tsx
    "settings.playback.engine_videocore_desc": "基于 Chromium 视频处理的 HTML5 播放器。", // playback-settings.tsx VideoCore 选项描述
    "settings.playback.engine_mpvcore_desc": "基于 libmpv 的原生播放器，支持更多编解码器。", // playback-settings.tsx MpvCore 选项描述
    "settings.playback.enable_logging": "启用日志记录", // playback-settings.tsx
    "settings.playback.logging_help": "启用后，调试日志将写入 Denshi 数据目录。", // playback-settings.tsx
    "settings.playback.export_logs": "导出日志", // playback-settings.tsx
    "settings.playback.logs_exported": "MpvCore 日志已导出", // playback-settings.tsx toast
    "settings.playback.logs_export_failed": "导出 MpvCore 日志失败", // playback-settings.tsx toast 失败文案
    "settings.playback.custom_mpv_options": "自定义 MPV 选项", // playback-settings.tsx
    "settings.playback.custom_mpv_desc_prefix": "添加自定义", // playback-settings.tsx 描述前半（<code>mpv.conf</code> 夹中间）
    "settings.playback.custom_mpv_desc_suffix": "选项。", // playback-settings.tsx 描述后半
    "settings.playback.custom_mpv_placeholder": "# 在此处添加自定义配置", // playback-settings.tsx 文本域占位

    // ===== mediaplayer-settings.tsx（新增部分；页头/连播/外部关联复用 player.* 种子词条）=====
    "settings.mediaplayer.iina_notice_1": "为了让 IINA 正常配合 Seanime 工作，请确保在 IINA 常规设置中勾选", // mediaplayer-settings.tsx IINA 提示前半
    "settings.mediaplayer.iina_notice_quit": "所有窗口关闭后退出", // mediaplayer-settings.tsx IINA 提示强调项 1
    "settings.mediaplayer.iina_notice_2": "，并取消勾选", // mediaplayer-settings.tsx IINA 提示中段
    "settings.mediaplayer.iina_notice_keep": "播放完毕后保持窗口开启", // mediaplayer-settings.tsx IINA 提示强调项 2
    "settings.mediaplayer.iina_notice_3": "。", // mediaplayer-settings.tsx IINA 提示结尾
    "settings.mediaplayer.host_help": "VLC/MPC-HC", // mediaplayer-settings.tsx Host 字段帮助（专有名词保留）
    "settings.mediaplayer.option_mpc": "MPC-HC (Windows)", // mediaplayer-settings.tsx 默认播放器选项（专有名词保留）
    "settings.mediaplayer.option_iina": "IINA (macOS)", // mediaplayer-settings.tsx 默认播放器选项（专有名词保留）
    "settings.mediaplayer.app_path": "应用程序路径", // mediaplayer-settings.tsx VLC/MPC-HC 路径字段共用
    "settings.mediaplayer.socket_placeholder": "留空则使用自动生成的 socket", // mediaplayer-settings.tsx
    "settings.mediaplayer.socket_help": "仅当您想让 Seanime 附加到特定的 mpv IPC socket 时才需设置。", // mediaplayer-settings.tsx
    "settings.mediaplayer.mpv_path_placeholder_default": "默认使用 CLI", // mediaplayer-settings.tsx Linux 占位
    "settings.mediaplayer.cli_leave_empty": "留空则使用命令行。", // mediaplayer-settings.tsx mpv/iina 路径字段共用
    "settings.mediaplayer.cli_path": "CLI 路径", // mediaplayer-settings.tsx
    "settings.mediaplayer.iina_cli_placeholder": "IINA CLI 的路径", // mediaplayer-settings.tsx
    "settings.mediaplayer.socket_default": "默认：'{value}'", // mediaplayer-settings.tsx IINA socket 占位（含插值）
    "settings.mediaplayer.encode_path": "在 URL 中编码文件路径（仅媒体库）", // mediaplayer-settings.tsx 外部播放器关联
    "settings.mediaplayer.encode_path_help": "启用后，文件路径将在 URL 中进行 base64 编码，以避免特殊字符引起的问题。", // mediaplayer-settings.tsx
    "settings.common.saved_automatically": "设置会自动保存", // playback-settings / mediaplayer-settings 页尾提示共用

    // ===== torrentstream-settings.tsx =====
    "settings.torrentstream.autoselect_help": "让 Seanime 自动寻找最佳种子。", // torrentstream-settings.tsx
    "settings.torrentstream.preload_next": "预加载下一集", // torrentstream-settings.tsx
    "settings.torrentstream.preload_next_help": "在后台开始下载下一集。", // torrentstream-settings.tsx
    "settings.torrentstream.preload_next_more_help": "此功能仅部分实现，请不要依赖其正常工作。", // torrentstream-settings.tsx
    "settings.common.unstable": "不稳定", // torrentstream-settings.tsx ExperimentalBadge
    "settings.torrentstream.torrent_client": "种子客户端", // torrentstream-settings.tsx 手风琴标题
    "settings.torrentstream.host_help": "留空则使用默认值。此主机用于监听新的 uTP 和 TCP BitTorrent 连接。", // torrentstream-settings.tsx
    "settings.torrentstream.port_help": "留空则使用默认值。默认为 43213。", // torrentstream-settings.tsx
    "settings.torrentstream.disable_ipv6": "禁用 IPv6", // torrentstream-settings.tsx
    "settings.torrentstream.slow_seeding": "慢速做种", // torrentstream-settings.tsx
    "settings.torrentstream.slow_seeding_more_help": "这有助于避免网络问题。注意：慢速做种会显著延迟启动。", // torrentstream-settings.tsx
    "settings.torrentstream.disable_accelerated_startup": "禁用加速启动", // torrentstream-settings.tsx
    "settings.torrentstream.disable_accelerated_startup_more_help": "开启后将禁用启动期间激进的节点发现与连接数限制。", // torrentstream-settings.tsx
    "settings.common.advanced": "高级", // torrentstream-settings.tsx 手风琴标题
    "settings.torrentstream.stream_url_address": "流 URL 地址", // torrentstream-settings.tsx
    "settings.torrentstream.stream_url_placeholder": "例如: 0.0.0.0:43211", // torrentstream-settings.tsx
    "settings.torrentstream.stream_url_help": "修改流 URL 的格式。留空则使用默认值。", // torrentstream-settings.tsx
    "settings.torrentstream.cache_directory": "缓存目录", // torrentstream-settings.tsx
    "settings.torrentstream.cache_directory_help": "流式播放时种子的下载位置。留空则使用默认缓存目录。", // torrentstream-settings.tsx
    "settings.torrentstream.cache_directory_warning": "请选择一个空目录，以免丢失数据。", // torrentstream-settings.tsx Alert
    "settings.torrentstream.drop_torrent": "删除种子", // torrentstream-settings.tsx 按钮

    // ===== mediastream-settings.tsx（NVIDIA/Intel/VAAPI/Apple 与 FFmpeg preset 名为技术参数，保留）=====
    "settings.mediastream.title": "转码 / 直连播放", // mediastream-settings.tsx 页头
    "settings.mediastream.desc": "管理转码与直连播放设置", // mediastream-settings.tsx 页头
    "settings.mediastream.hwaccel_cpu": "CPU（已禁用）", // mediastream-settings.tsx 硬件加速选项
    "settings.mediastream.direct_play_title": "直连播放", // mediastream-settings.tsx 卡片标题
    "settings.mediastream.prefer_transcoding": "优先转码", // mediastream-settings.tsx
    "settings.mediastream.prefer_transcoding_help": "启用后，即使客户端支持媒体编解码器，Seanime 也不会自动切换到直连播放。", // mediastream-settings.tsx
    "settings.mediastream.direct_play_only": "仅直连播放", // mediastream-settings.tsx
    "settings.mediastream.direct_play_only_help": "仅允许直连播放，永远不会启动转码。", // mediastream-settings.tsx
    "settings.mediastream.transcoding_title": "转码", // mediastream-settings.tsx 卡片标题
    "settings.mediastream.hardware_acceleration": "硬件加速", // mediastream-settings.tsx
    "settings.mediastream.hardware_acceleration_help": "强烈建议启用硬件加速，以获得更流畅的转码体验。", // mediastream-settings.tsx
    "settings.mediastream.custom_settings_json": "自定义设置 (JSON)", // mediastream-settings.tsx
    "settings.mediastream.custom_settings_help": "仅作用于视频流，scaleFilter = -vf，-map、-bufsize、-b:v、-maxrate 会自动应用。", // mediastream-settings.tsx
    "settings.mediastream.transcode_preset": "转码预设", // mediastream-settings.tsx
    "settings.mediastream.transcode_preset_help": "推荐使用“Fast”。VAAPI 不支持预设。", // mediastream-settings.tsx
    "settings.mediastream.ffmpeg_path": "FFmpeg 路径", // mediastream-settings.tsx
    "settings.mediastream.ffmpeg_path_help": "FFmpeg 可执行文件的路径。如果已在 PATH 中可留空。", // mediastream-settings.tsx
    "settings.mediastream.ffprobe_path": "FFprobe 路径", // mediastream-settings.tsx
    "settings.mediastream.ffprobe_path_help": "FFprobe 可执行文件的路径。如果已在 PATH 中可留空。", // mediastream-settings.tsx

    // ===== manga-settings.tsx =====
    "settings.manga.title": "漫画", // manga-settings.tsx 页头
    "settings.manga.desc": "管理您的漫画媒体库", // manga-settings.tsx 页头
    "settings.manga.enable_help": "阅读漫画系列、下载章节并追踪进度。", // manga-settings.tsx
    "settings.manga.auto_update_progress": "自动更新进度", // manga-settings.tsx
    "settings.manga.auto_update_progress_help": "启用后，读到章节末尾时将自动更新您的进度。", // manga-settings.tsx
    "settings.manga.provider_title": "漫画源", // manga-settings.tsx 卡片标题（漫画域习惯译法，与 Debrid「提供商」区分）
    "settings.manga.default_provider": "默认漫画源", // manga-settings.tsx
    "settings.manga.default_provider_help": "打开新漫画系列时默认选中的漫画源。", // manga-settings.tsx
    "settings.manga.overwrite_confirm_title": "覆盖所有来源", // manga-settings.tsx 确认弹窗
    "settings.manga.overwrite_confirm_desc": "这将用默认漫画源覆盖您已打开的所有漫画系列的来源选择。确定要继续吗？", // manga-settings.tsx 确认弹窗
    "settings.manga.overwrite_action": "覆盖", // manga-settings.tsx 确认弹窗按钮
    "settings.manga.overwritten_toast": "所有来源选择已被覆盖。", // manga-settings.tsx toast
    "settings.manga.overwrite_with": "用 {name} 覆盖所有漫画来源", // manga-settings.tsx 按钮（含插值）
    "settings.manga.restored_toast": "已恢复之前的来源选择。", // manga-settings.tsx toast
    "settings.manga.undo": "撤销", // manga-settings.tsx 按钮
    "settings.manga.local_provider_title": "本地漫画源", // manga-settings.tsx 卡片标题
    "settings.manga.local_provider_desc": "从您的本地目录阅读漫画系列。", // manga-settings.tsx
    "settings.manga.local_source_directory": "本地漫画目录", // manga-settings.tsx
    "settings.manga.local_source_directory_help": "漫画存放目录。仅由本地漫画源使用。", // manga-settings.tsx

    // ===== nakama-settings.tsx =====
    "settings.nakama.title": "Nakama", // nakama-settings.tsx 页头（专有名词保留）
    "settings.nakama.desc": "与其他 Seanime 实例通信", // nakama-settings.tsx 页头
    "settings.nakama.enable": "启用 Nakama", // nakama-settings.tsx
    "settings.nakama.username_help": "用于向其他实例标识此服务器的用户名。留空将分配一个随机 ID。", // nakama-settings.tsx
    "settings.nakama.tab_peer": "以节点身份连接", // nakama-settings.tsx 标签页
    "settings.nakama.tab_host": "作为主机", // nakama-settings.tsx 标签页
    "settings.nakama.currently_hosting": "正在作为主机运行", // nakama-settings.tsx 徽标
    "settings.nakama.reminder_title": "提醒", // nakama-settings.tsx Alert
    "settings.nakama.reminder_desc": "将服务器暴露到互联网之前，请先在配置文件中设置密码。", // nakama-settings.tsx Alert
    "settings.nakama.host_mode_desc_1": "主机模式适用于可通过互联网访问的 Seanime 服务器。", // nakama-settings.tsx
    "settings.nakama.host_mode_desc_2_prefix": "不过，您可以使用", // nakama-settings.tsx 第二段前半
    "settings.nakama.host_mode_cloud_rooms": "Cloud Rooms", // nakama-settings.tsx 强调项（专有名词保留）
    "settings.nakama.host_mode_desc_2_suffix": "来在不暴露服务器的情况下举办一起观看派对。", // nakama-settings.tsx 第二段后半
    "settings.nakama.enable_host_mode": "启用主机模式", // nakama-settings.tsx
    "settings.nakama.enable_host_mode_help": "启用后，此服务器将作为其他客户端的主机。需要先设置主机密码。", // nakama-settings.tsx
    "settings.nakama.passcode": "通行码", // nakama-settings.tsx 标签与占位共用
    "settings.nakama.passcode_help": "设置通行码以保护您的主机模式与房间。此通行码应与服务器密码不同。", // nakama-settings.tsx
    "settings.nakama.share_local_library": "共享本地动漫媒体库", // nakama-settings.tsx
    "settings.nakama.share_local_library_help": "启用后，此服务器将向其他客户端共享其本地动漫媒体库。", // nakama-settings.tsx
    "settings.nakama.exclude_from_sharing": "排除共享的动漫", // nakama-settings.tsx
    "settings.nakama.exclude_from_sharing_help": "选择您不想与其他客户端共享的动漫。", // nakama-settings.tsx
    "settings.nakama.cannot_connect_while_host": "主机模式下无法连接到其他主机。", // nakama-settings.tsx Alert
    "settings.nakama.server_url": "Nakama 服务器 URL", // nakama-settings.tsx
    "settings.nakama.server_url_help": "要连接的 Nakama 主机的 URL。", // nakama-settings.tsx
    "settings.nakama.remote_passcode": "Nakama 通行码", // nakama-settings.tsx
    "settings.nakama.remote_passcode_help": "连接 Nakama 主机所需的通行码。", // nakama-settings.tsx
    "settings.nakama.use_nakama_library": "使用 Nakama 的动漫媒体库", // nakama-settings.tsx
    "settings.nakama.use_nakama_library_help": "启用后，如果对方正在共享，将使用 Nakama 的动漫媒体库作为您的媒体库。", // nakama-settings.tsx

    // ===== anime-library-settings.tsx（目录/扫描标签复用 settings.ts 既有词条）=====
    "settings.anime_library.main_path_help": "存放动漫视频文件的根目录路径。（请保持路径字母大小写一致）", // anime-library-settings.tsx
    "settings.anime_library.extended_path_help": "如果您的动漫存放在多个不同磁盘或目录下，可以在此处添加其他路径。", // anime-library-settings.tsx
    "settings.anime_library.auto_refresh_more_help": "批量添加新番时，可能需要等待扫描器识别。", // anime-library-settings.tsx
    "settings.anime_library.advanced_rules": "高级扫描规则", // anime-library-settings.tsx 手风琴标题
    "settings.anime_library.scanner_rules_title": "扫描器规则配置 (JSON)", // anime-library-settings.tsx
    "settings.anime_library.scanner_rules_desc": "以 JSON 格式配置高级扫描与番剧识别规则，支持自定义匹配正则和媒体库清洗。", // anime-library-settings.tsx
    "settings.anime_library.use_legacy_matching": "使用传统匹配算法", // anime-library-settings.tsx
    "settings.anime_library.use_legacy_matching_help": "启用以使用传统匹配算法。（3.4 及以下版本）", // anime-library-settings.tsx
    "settings.anime_library.use_legacy_matching_more_help": "传统匹配算法使用更简单的方法，准确性可能较低。", // anime-library-settings.tsx
    "settings.anime_library.matching_algorithm": "匹配算法", // anime-library-settings.tsx
    "settings.anime_library.matching_algorithm_help": "选择用于将文件匹配到 AniList 条目的算法。", // anime-library-settings.tsx
    "settings.anime_library.matching_threshold": "匹配阈值", // anime-library-settings.tsx
    "settings.anime_library.matching_threshold_help": "文件匹配到 AniList 条目所需的最低分数。默认为 0.5。", // anime-library-settings.tsx
    "settings.anime_library.algorithm_default": "Levenshtein + Sorensen-Dice（默认）", // anime-library-settings.tsx 选项（算法专名保留）

    // ===== autoselect-profile-form.tsx（profile 字段名/条件表达式/suggestions 数据值冻结）=====
    "settings.autoselect.customize": "自定义自动选择", // autoselect-profile-form.tsx 按钮
    "settings.autoselect.summary_resolutions": "分辨率：", // autoselect-profile-form.tsx 概要标签
    "settings.autoselect.summary_groups": "分组：", // autoselect-profile-form.tsx 概要标签
    "settings.autoselect.summary_providers": "提供商：", // autoselect-profile-form.tsx 概要标签
    "settings.autoselect.summary_codecs": "首选编码：", // autoselect-profile-form.tsx 概要标签
    "settings.autoselect.summary_sources": "首选来源：", // autoselect-profile-form.tsx 概要标签
    "settings.autoselect.summary_languages": "首选语言：", // autoselect-profile-form.tsx 概要标签
    "settings.autoselect.summary_multi_audio": "多音轨：", // autoselect-profile-form.tsx 概要标签
    "settings.autoselect.summary_multi_subs": "多字幕：", // autoselect-profile-form.tsx 概要标签
    "settings.autoselect.summary_best_releases": "最佳发布：", // autoselect-profile-form.tsx 概要标签
    "settings.autoselect.summary_batches": "批量合集：", // autoselect-profile-form.tsx 概要标签
    "settings.autoselect.release_prefs": "发布偏好", // autoselect-profile-form.tsx 分组标题
    "settings.autoselect.metadata_prefs": "元数据偏好", // autoselect-profile-form.tsx 分组标题
    "settings.autoselect.special_prefs": "特殊偏好", // autoselect-profile-form.tsx 分组标题
    "settings.autoselect.thresholds": "阈值", // autoselect-profile-form.tsx 分组标题
    "settings.autoselect.min_seeders": "最小做种数", // autoselect-profile-form.tsx
    "settings.autoselect.min_size": "最小体积", // autoselect-profile-form.tsx
    "settings.autoselect.max_size": "最大体积", // autoselect-profile-form.tsx
    "settings.autoselect.min_size_placeholder": "例如: 100MB", // autoselect-profile-form.tsx
    "settings.autoselect.max_size_placeholder": "例如: 2GB 或 10GiB", // autoselect-profile-form.tsx
    "settings.autoselect.providers_ordered": "提供商（按顺序，最多 3 个）", // autoselect-profile-form.tsx
    "settings.autoselect.providers_hint": "按优先级顺序最多选择 3 个提供商", // autoselect-profile-form.tsx
    "settings.autoselect.empty_providers": "未找到提供商", // autoselect-profile-form.tsx Combobox 空状态
    "settings.autoselect.release_groups_ordered": "发布组（按顺序）", // autoselect-profile-form.tsx
    "settings.autoselect.release_groups_hint": "按优先级顺序的首选发布组", // autoselect-profile-form.tsx
    "settings.autoselect.resolutions_ordered": "分辨率（按顺序）", // autoselect-profile-form.tsx
    "settings.autoselect.resolutions_hint": "按优先级顺序的首选分辨率", // autoselect-profile-form.tsx
    "settings.autoselect.exclude_terms": "排除关键词", // autoselect-profile-form.tsx
    "settings.autoselect.exclude_terms_hint": "排除包含这些关键词的种子", // autoselect-profile-form.tsx
    "settings.autoselect.exclude_terms_placeholder": "例如: CamRip, Cam RIP", // autoselect-profile-form.tsx
    "settings.autoselect.languages_ordered": "首选语言（按顺序）", // autoselect-profile-form.tsx
    "settings.autoselect.languages_hint": "按优先级排序的首选语言列表", // autoselect-profile-form.tsx
    "settings.autoselect.codecs_ordered": "首选编码（按顺序）", // autoselect-profile-form.tsx
    "settings.autoselect.codecs_hint": "按优先级排序的首选编码列表（逗号分隔多个备选）", // autoselect-profile-form.tsx
    "settings.autoselect.sources_ordered": "首选来源（按顺序）", // autoselect-profile-form.tsx
    "settings.autoselect.sources_hint": "按优先级排序的首选来源列表（逗号分隔多个备选）", // autoselect-profile-form.tsx
    "settings.autoselect.require_language": "要求语言匹配", // autoselect-profile-form.tsx
    "settings.autoselect.require_language_help": "未找到首选语言时拒绝", // autoselect-profile-form.tsx
    "settings.autoselect.require_codec": "要求编码匹配", // autoselect-profile-form.tsx
    "settings.autoselect.require_codec_help": "未找到首选编码时拒绝", // autoselect-profile-form.tsx
    "settings.autoselect.require_source": "要求来源匹配", // autoselect-profile-form.tsx
    "settings.autoselect.require_source_help": "未找到首选来源时拒绝", // autoselect-profile-form.tsx
    "settings.autoselect.pref_multiple_audio": "多音轨", // autoselect-profile-form.tsx PreferenceField 标签
    "settings.autoselect.pref_multiple_subtitles": "多字幕", // autoselect-profile-form.tsx
    "settings.autoselect.pref_batches": "批量合集", // autoselect-profile-form.tsx
    "settings.autoselect.pref_best_releases": "最佳发布", // autoselect-profile-form.tsx
    "settings.autoselect.pref_neutral": "中立", // autoselect-profile-form.tsx PreferenceField 选项
    "settings.autoselect.pref_prefer": "优先", // autoselect-profile-form.tsx
    "settings.autoselect.pref_avoid": "规避", // autoselect-profile-form.tsx
    "settings.autoselect.pref_only": "仅限", // autoselect-profile-form.tsx
    "settings.autoselect.pref_never": "绝不", // autoselect-profile-form.tsx
    "settings.mediaplayer.socket_label": "套接字（Socket）", // mediaplayer-settings.tsx Field.Text label（09-05 残留）
} satisfies Dictionary
