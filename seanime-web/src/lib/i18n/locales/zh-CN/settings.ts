import type { Dictionary } from "../../types"

/**
 * settings 词条表
 * 归属：系统设置页的标签页、设置卡片标题与字段标签
 */
export const settingsDictionary = {
    "settings.page.title": "系统设置", // src/app/(main)/settings/page.tsx
    "settings.tab.general": "常规设置", // src/app/(main)/settings/page.tsx
    "settings.tab.appearance": "界面外观", // src/app/(main)/settings/page.tsx
    "settings.tab.library": "本地动漫库", // src/app/(main)/settings/page.tsx
    "settings.tab.video_playback": "视频播放", // src/app/(main)/settings/page.tsx
    "settings.tab.desktop_player": "桌面播放器", // src/app/(main)/settings/page.tsx
    "settings.tab.external_player": "外部播放器关联", // src/app/(main)/settings/page.tsx
    "settings.tab.transcode": "转码 / 直播流", // src/app/(main)/settings/page.tsx
    "settings.tab.torrent_provider": "种子源提供商", // src/app/(main)/settings/page.tsx
    "settings.tab.download_client": "下载客户端", // src/app/(main)/settings/page.tsx
    "settings.tab.downloader": "边下边播", // src/app/(main)/settings/page.tsx
    "settings.tab.debrid": "Debrid 云端服务", // src/app/(main)/settings/page.tsx
    "settings.tab.onlinestream": "在线流媒体", // src/app/(main)/settings/page.tsx
    "settings.tab.manga": "漫画设置", // src/app/(main)/settings/page.tsx
    "settings.tab.nakama": "Nakama 共享", // src/app/(main)/settings/page.tsx
    "settings.tab.discord": "Discord 状态", // src/app/(main)/settings/page.tsx
    "settings.tab.denshi": "Denshi 桌面端", // src/app/(main)/settings/page.tsx
    "settings.tab.logs_cache": "日志与缓存", // src/app/(main)/settings/page.tsx
    "settings.action.sponsor": "赞助支持", // src/app/(main)/settings/page.tsx
    "settings.action.open_data_dir": "打开数据目录", // src/app/(main)/settings/page.tsx
    "settings.action.report_issue": "反馈问题", // src/app/(main)/settings/page.tsx
    "settings.action.check_update": "检查更新", // src/app/(main)/settings/page.tsx
    "settings.toast.up_to_date": "当前已是最新版本", // src/app/(main)/settings/page.tsx
    "settings.card.general_title": "常规设置", // src/app/(main)/settings/page.tsx
    "settings.card.general_desc": "通用应用与服务配置", // src/app/(main)/settings/page.tsx
    "settings.card.library_title": "本地动漫库", // src/app/(main)/settings/page.tsx
    "settings.card.library_desc": "管理本地动漫媒体库与扫描配置", // src/app/(main)/settings/page.tsx
    "settings.card.onlinestream_title": "在线流媒体", // src/app/(main)/settings/page.tsx
    "settings.card.onlinestream_desc": "配置第三方在线动漫播放与流媒体选项", // src/app/(main)/settings/page.tsx
    "settings.card.torrent_provider_title": "种子源提供商", // src/app/(main)/settings/page.tsx
    "settings.card.torrent_provider_desc": "配置番剧种子检索与默认提供商扩展", // src/app/(main)/settings/page.tsx
    "settings.card.download_client_title": "下载客户端", // src/app/(main)/settings/page.tsx
    "settings.card.download_client_desc": "配置外部或内置下载客户端", // src/app/(main)/settings/page.tsx
    "settings.card.client_integration_title": "客户端集成", // src/app/(main)/settings/page.tsx
    "settings.field.default_download_client": "默认下载客户端", // src/app/(main)/settings/page.tsx
    "settings.field.builtin_client": "内置客户端", // src/app/(main)/settings/page.tsx
    "settings.field.max_connections_per_torrent": "每个种子最大连接数", // src/app/(main)/settings/page.tsx
    "settings.field.max_active_downloads": "最大同时下载数", // src/app/(main)/settings/page.tsx
    "settings.field.download_speed_limit": "下载限速 (KB/s)", // src/app/(main)/settings/page.tsx
    "settings.field.upload_speed_limit": "上传限速 (KB/s)", // src/app/(main)/settings/page.tsx
    "settings.field.show_active_downloads": "在侧边栏显示活跃下载数", // src/app/(main)/settings/page.tsx
    "settings.field.open_client_on_startup": "开机启动时打开下载客户端", // src/app/(main)/settings/page.tsx
    "settings.card.appearance_desc": "自定义主题风格与界面展示选项", // src/app/(main)/settings/page.tsx
    "settings.card.streaming_desc": "配置种子在线缓冲与边下边播", // src/app/(main)/settings/page.tsx
    "settings.card.logs_desc": "查看系统运行日志", // src/app/(main)/settings/page.tsx
    "settings.card.cache_desc": "管理本地文件与数据缓存", // src/app/(main)/settings/page.tsx
    "settings.card.library_dir_title": "本地媒体库目录", // src/app/(main)/settings/_containers/anime-library-settings.tsx
    "settings.field.main_library_path": "主动漫库路径", // src/app/(main)/settings/_containers/anime-library-settings.tsx
    "settings.field.extended_library_path": "扩展动漫库路径", // src/app/(main)/settings/_containers/anime-library-settings.tsx
    "settings.card.library_scan_title": "媒体库扫描", // src/app/(main)/settings/_containers/anime-library-settings.tsx
    "settings.field.auto_refresh_library": "自动定时刷新媒体库", // src/app/(main)/settings/_containers/anime-library-settings.tsx
    "settings.field.scan_on_startup": "服务启动时自动扫描刷新", // src/app/(main)/settings/_containers/anime-library-settings.tsx
    "settings.card.episodes_title": "剧集与播放记录", // src/app/(main)/settings/_containers/server-settings.tsx
    "settings.field.auto_update_progress": "自动更新观看进度", // src/app/(main)/settings/_containers/server-settings.tsx
    "settings.field.enable_playback_history": "启用播放历史记录", // src/app/(main)/settings/_containers/server-settings.tsx
    "settings.field.default_episode_source": "默认剧集来源", // src/app/(main)/settings/_containers/server-settings.tsx
    "settings.card.spoilers_title": "动漫展示与防剧透", // src/app/(main)/settings/_containers/server-settings.tsx
    "settings.field.hide_spoilers": "隐藏动漫剧透", // src/app/(main)/settings/_containers/server-settings.tsx

    // —— 07-01 追加：通用动作 / 选项（settings 域内多文件复用）——
    "settings.action.save": "保存", // _components/settings-submit-button.tsx, _containers/ui-settings.tsx
    "settings.action.reset": "重置", // _components/settings-submit-button.tsx, _containers/server-settings.tsx
    "settings.action.enable": "启用", // _containers/debrid-settings.tsx, _containers/discord-rich-presence-settings.tsx
    "settings.action.upload": "上传", // 确认弹窗 actionText：server-settings / local-settings / anilist-settings
    "settings.action.import": "导入", // _containers/data-settings.tsx
    "settings.common.or": "或", // _containers/server-settings.tsx
    "settings.option.none": "无", // page.tsx 与 debrid-settings.tsx 的「None」选项
    "settings.option.highest": "最高", // debrid-settings.tsx（torrentstream 同名词段复用）
    "settings.option.automatic": "自动", // server-settings.tsx 默认剧集来源选项
    "settings.option.local_library": "本地媒体库", // server-settings.tsx
    "settings.option.debrid_streaming": "Debrid 流媒体", // server-settings.tsx
    "settings.option.torrent_streaming": "种子流媒体", // server-settings.tsx
    "settings.option.online_streaming": "在线流媒体", // server-settings.tsx
    "settings.option.unavailable_plugin": "不可用的扩展", // server-settings.tsx

    // —— 07-01 追加：通用字段（page.tsx / 容器复用）——
    "settings.field.host": "主机", // page.tsx qBittorrent/Transmission 区块
    "settings.field.username": "用户名", // page.tsx
    "settings.field.password": "密码", // page.tsx
    "settings.field.port": "端口", // page.tsx
    "settings.field.executable_path": "执行文件路径", // page.tsx
    "settings.field.tags": "标签", // page.tsx
    "settings.field.category": "分类", // page.tsx
    "settings.field.listen_port": "监听端口", // page.tsx 内置客户端
    "settings.field.default_provider": "默认提供商", // page.tsx 种子源提供商

    // —— 07-01 追加：page.tsx 回填 ——
    "settings.help.qbittorrent_tags": "添加到下载任务的标签，逗号分隔，如 seanime,anime", // page.tsx
    "settings.help.qbittorrent_category": "添加到下载任务的分类目录名称。", // page.tsx
    "settings.help.speed_limit_zero": "设置为 0 表示不限速。", // page.tsx 下载/上传限速共用
    "settings.help.show_active_downloads": "在左侧导航栏的种子图标旁显示当前活跃下载数量。(较消耗内存)", // page.tsx
    "settings.card.discord_title": "Discord", // page.tsx Discord 页头
    "settings.card.discord_desc": "配置 Discord 正在播放状态展示", // page.tsx
    "settings.field.enable_onlinestream": "启用在线流媒体", // page.tsx
    "settings.help.enable_onlinestream": "直接从第三方在线流媒体源观看动漫剧集。", // page.tsx
    "settings.card.home_display_title": "首页展示", // page.tsx 在线流媒体卡片
    "settings.field.include_streaming_in_library": "在追番列表中包含流媒体作品", // page.tsx
    "settings.help.include_streaming_in_library": "在您的动漫片单中展示正在观看的在线流媒体动漫。", // page.tsx
    "settings.field.show_torrent_availability": "在最近剧集中显示种子可用性", // page.tsx
    "settings.help.show_torrent_availability": "为媒体库缺失的最近剧集，以及使用种子或 Debrid 流媒体播放时的「继续观看」添加角标。", // page.tsx
    "settings.help.torrent_provider": "供搜索引擎使用。如不需要种子支持可选择「无」。", // page.tsx
    "settings.help.torrent_availability_missing": "选择一个种子源提供商以检查剧集可用性。", // page.tsx
    "settings.card.denshi_title": "Denshi", // page.tsx Denshi 页头
    "settings.card.denshi_desc": "桌面客户端设置", // page.tsx
    "settings.client.transmission": "Transmission", // page.tsx 下载客户端标题与选项（专有名词保留）
    "settings.card.logs_title": "运行日志", // page.tsx 日志页头
    "settings.card.cache_title": "缓存管理", // page.tsx 缓存页头

    // —— 07-01 追加：server-settings.tsx ——
    "settings.alert.anilist_api_down": "AniList API 无法访问，所有请求将由缓存提供。", // server-settings.tsx
    "settings.alert.anilist_api_down_hint": "您可以在应用设置中禁用此功能。", // server-settings.tsx
    "settings.help.auto_update_progress": "启用后，当观看剧集达到 80% 时，将自动同步更新您的追番进度。", // server-settings.tsx
    "settings.help.auto_update_progress_more": "仅适用于桌面播放器或内置播放器。", // server-settings.tsx（两处共用）
    "settings.help.playback_history": "启用后，Seanime 将记录您的播放进度并在下次打开时从上次离开的位置继续播放。", // server-settings.tsx
    "settings.help.default_episode_source": "打开动漫详情页时的默认数据播放源。", // server-settings.tsx
    "settings.field.hide_thumbnails": "隐藏缩略图", // server-settings.tsx
    "settings.field.hide_episode_titles": "隐藏剧集标题", // server-settings.tsx
    "settings.field.hide_episode_descriptions": "隐藏剧集简介", // server-settings.tsx
    "settings.field.skip_spoilers_next_episode": "下一集跳过防剧透", // server-settings.tsx
    "settings.help.hide_spoilers": "在继续观看、剧集列表以及缺失剧集中隐藏剧透性的缩略图和文本。", // server-settings.tsx
    "settings.help.skip_spoilers_next_episode": "从下一集之后的剧集开始隐藏剧透。", // server-settings.tsx
    "settings.field.hide_audience_score": "隐藏观众评分", // server-settings.tsx
    "settings.help.hide_audience_score": "启用后，直到您主动点击查看前，番剧观众评分将被遮挡。", // server-settings.tsx
    "settings.field.enable_adult_content": "启用成人/R18内容", // server-settings.tsx
    "settings.help.enable_adult_content": "关闭后，成人向内容将从搜索结果和媒体库中隐藏。", // server-settings.tsx
    "settings.field.blur_adult_content": "对成人内容应用模糊遮罩", // server-settings.tsx
    "settings.field.disable_card_trailers": "禁用动漫卡片预告片", // server-settings.tsx
    "settings.card.extension_security_title": "扩展安全性", // server-settings.tsx
    "settings.field.enable_extension_secure_mode": "启用扩展安全模式", // server-settings.tsx
    "settings.help.enable_extension_secure_mode": "启用后，即使已授予权限，扩展尝试执行敏感操作时 Seanime 仍会弹出确认提示。", // server-settings.tsx
    "settings.card.local_sync_title": "本地账户与数据同步", // server-settings.tsx
    "settings.card.local_sync_desc": "未绑定或离线使用 AniList 账号时使用本地账户。", // server-settings.tsx
    "settings.field.auto_backup_anilist": "自动备份 AniList 列表", // server-settings.tsx
    "settings.help.auto_backup_anilist": "启用后，将定期使用您的 AniList 数据更新本地列表。这会覆盖上次同步后您在本地所做的修改。", // server-settings.tsx
    "settings.action.upload_local_lists": "将本地列表上传至 AniList", // server-settings.tsx
    "settings.card.offline_title": "离线模式", // server-settings.tsx
    "settings.card.offline_desc": "绑定 AniList 账户时可用。", // server-settings.tsx
    "settings.field.auto_sync_offline_media": "自动同步离线媒体", // server-settings.tsx
    "settings.help.auto_sync_offline_media": "如果关闭，您需要在离线模式页面中手动点击“立即同步”来刷新本地元数据。", // server-settings.tsx
    "settings.help.auto_sync_offline_media_more": "如果您在离线时进行了修改且尚未同步到 AniList，该操作会暂停。", // server-settings.tsx
    "settings.field.auto_offline_cache_watching": "自动离线缓存正在观看/阅读的作品", // server-settings.tsx
    "settings.help.auto_offline_cache_watching": "启用后，Seanime 会自动将您正在追的番剧和漫画数据保存为离线可用。", // server-settings.tsx
    "settings.card.metadata_cache_title": "元数据与缓存提供方", // server-settings.tsx
    "settings.field.disable_anilist_cache": "禁用 AniList 请求磁盘缓存", // server-settings.tsx
    "settings.help.disable_anilist_cache": "启用后，Seanime 将不再把 AniList 的请求缓存到磁盘。", // server-settings.tsx
    "settings.help.disable_anilist_cache_more": "默认情况下，所有向 AniList 发起的请求均会缓存到本地，保证网络中断时仍可离线访问。", // server-settings.tsx
    "settings.field.cache_only_mode": "启用仅缓存模式", // server-settings.tsx
    "settings.help.cache_only_mode": "Seanime 将直接使用本地缓存数据，不再向外网 API 发起请求。", // server-settings.tsx
    "settings.field.fallback_metadata_provider": "使用备用剧集元数据源", // server-settings.tsx
    "settings.help.fallback_metadata_provider": "启用后，Seanime 将尝试从备用数据源拉取剧集标题与缩略图元数据。", // server-settings.tsx
    "settings.card.update_title": "版本更新", // server-settings.tsx
    "settings.field.no_update_desktop": "不获取更新", // server-settings.tsx 桌面端
    "settings.field.no_update_check": "不自动检查更新", // server-settings.tsx 网页端
    "settings.help.no_update_desktop": "启用后将不再提示新版本。", // server-settings.tsx
    "settings.help.no_update_check": "启用后，Seanime 将不再自动检测 GitHub 新版本发布。", // server-settings.tsx
    "settings.help.no_update_desktop_more": "Seanime Denshi 桌面端无法关闭静默自动更新。", // server-settings.tsx
    "settings.field.update_channel": "更新渠道", // server-settings.tsx
    "settings.help.update_channel_desktop": "同时也适用于 Seanime Denshi 桌面客户端自动更新。", // server-settings.tsx
    "settings.option.update_channel_github": "GitHub (默认)", // server-settings.tsx
    "settings.option.update_channel_seanime": "Seanime 官方源", // server-settings.tsx
    "settings.option.update_channel_canary": "Seanime (Canary 测试版)", // server-settings.tsx
    "settings.toast.update_channel_seanime": "您当前正在使用 Seanime 托管的发布渠道。", // server-settings.tsx
    "settings.toast.update_channel_canary": "您当前正在使用 Canary 金丝雀尝鲜分支，可能会接收到未经全面测试的不稳定更新。", // server-settings.tsx
    "settings.card.notifications_title": "服务器与系统通知", // server-settings.tsx
    "settings.field.open_web_on_start": "服务启动时自动打开网页端", // server-settings.tsx
    "settings.field.disable_desktop_notifications": "禁用系统桌面通知", // server-settings.tsx
    "settings.help.disable_desktop_notifications": "由操作系统展示的弹窗通知", // server-settings.tsx
    "settings.field.disable_downloader_notifications": "禁用自动下载器通知", // server-settings.tsx
    "settings.field.disable_scanner_notifications": "禁用自动扫描器通知", // server-settings.tsx
    "settings.card.shortcuts_title": "快捷键设置", // server-settings.tsx
    "settings.field.open_command_palette": "打开命令面板", // server-settings.tsx

    // —— 07-01 追加：确认弹窗（server/local/anilist 共用）——
    "settings.confirm.upload_to_anilist": "上传到 AniList", // server-settings / local-settings / anilist-settings
    "settings.confirm.upload_to_anilist_desc": "这将把您的本地 Seanime 收藏上传到您的 AniList 账户。确定要继续吗？", // 同上

    // —— 07-01 追加：ui-settings.tsx ——
    "settings.ui.tab_general": "常规", // ui-settings.tsx 子标签
    "settings.ui.apply_to_client": "应用到当前客户端", // ui-settings.tsx
    "settings.ui.custom_css_note": "自定义 CSS 将保存在服务端，需要手动应用到每个客户端。", // ui-settings.tsx
    "settings.ui.custom_css_note_devtools": "如果 CSS 导致界面无法正常显示，可随时通过开发者工具从本地存储中将其移除。", // ui-settings.tsx
    "settings.field.custom_css": "自定义 CSS", // ui-settings.tsx
    "settings.placeholder.custom_css": "自定义 CSS", // ui-settings.tsx
    "settings.field.mobile_custom_css": "移动端自定义 CSS", // ui-settings.tsx
    "settings.help.custom_css_desktop": "应用于 1024px 以上屏幕。", // ui-settings.tsx
    "settings.help.custom_css_mobile": "应用于 1024px 以下屏幕。", // ui-settings.tsx
    "settings.card.sorting_title": "排序", // ui-settings.tsx
    "settings.alert.watch_continuity_sorting": "需要启用播放历史记录才能使用「上次观看」排序选项。", // ui-settings.tsx
    "settings.field.continue_watching_sorting": "继续观看排序", // ui-settings.tsx
    "settings.field.anime_library_sorting": "动漫媒体库排序", // ui-settings.tsx
    "settings.field.manga_library_sorting": "漫画媒体库排序", // ui-settings.tsx
    "settings.card.theme_title": "主题", // ui-settings.tsx
    "settings.field.enable_color_settings": "启用颜色设置", // ui-settings.tsx
    "settings.field.live_preview": "实时预览", // ui-settings.tsx
    "settings.help.live_preview": "关闭将重新加载页面，未应用的更改将丢失。", // ui-settings.tsx
    "settings.field.background_color": "背景颜色", // ui-settings.tsx
    "settings.help.background_color": "默认：#070707", // ui-settings.tsx
    "settings.field.accent_color": "强调颜色", // ui-settings.tsx
    "settings.help.accent_color": "默认：#6152df", // ui-settings.tsx
    "settings.card.banners_title": "横幅与背景", // ui-settings.tsx
    "settings.field.background_image_path": "背景图片路径", // ui-settings.tsx
    "settings.placeholder.image_path": "例如 image.png", // ui-settings.tsx
    "settings.help.background_image": "所有页面的背景图片。非媒体库页面会变暗显示。", // ui-settings.tsx
    "settings.field.background_opacity": "背景图片不透明度", // ui-settings.tsx
    "settings.placeholder.default_10": "默认：10", // ui-settings.tsx 背景与横幅不透明度共用
    "settings.field.banner_image_path": "横幅图片路径", // ui-settings.tsx
    "settings.help.banner_image": "所有页面的横幅图片。", // ui-settings.tsx
    "settings.field.banner_position": "横幅位置", // ui-settings.tsx
    "settings.placeholder.default_position": "默认：50% 50%", // ui-settings.tsx
    "settings.field.banner_opacity": "横幅不透明度", // ui-settings.tsx
    "settings.field.home_banner_type": "首页横幅类型", // ui-settings.tsx
    "settings.field.media_banner_image": "媒体页横幅图片", // ui-settings.tsx
    "settings.field.media_banner_size": "媒体页横幅尺寸", // ui-settings.tsx
    "settings.ui.banner_type_dynamic": "动态", // ui-settings.tsx 横幅类型选项
    "settings.ui.banner_type_custom": "自定义", // ui-settings.tsx
    "settings.ui.open_assets_dir": "打开素材目录", // ui-settings.tsx
    "settings.help.home_banner_type": "在所有媒体库页面使用横幅图片。", // ui-settings.tsx
    "settings.card.tweaks_title": "微调", // ui-settings.tsx
    "settings.field.navigation_preloading": "导航预加载", // ui-settings.tsx
    "settings.ui.preload_disabled": "已禁用", // ui-settings.tsx 预加载模式选项
    "settings.ui.preload_disabled_desc": "不预加载", // ui-settings.tsx
    "settings.ui.preload_intent": "意图", // ui-settings.tsx
    "settings.ui.preload_intent_desc": "悬停时预加载", // ui-settings.tsx
    "settings.ui.preload_faster_intent": "快速意图", // ui-settings.tsx
    "settings.ui.preload_faster_intent_desc": "更激进地预加载", // ui-settings.tsx
    "settings.ui.preload_viewport": "视口", // ui-settings.tsx
    "settings.ui.preload_viewport_desc": "进入可视区域时预加载", // ui-settings.tsx
    "settings.help.navigation_preloading": "应用于当前客户端的媒体页。预加载可能使您更快达到请求频率限制。", // ui-settings.tsx
    "settings.alert.preloading_simulated": "由于请求频率限制，未登录 AniList 账号时导航预加载已禁用。", // ui-settings.tsx
    "settings.field.remove_genre_selector": "移除流派选择器", // ui-settings.tsx
    "settings.field.enable_blurring": "启用模糊效果", // ui-settings.tsx
    "settings.help.blurring_effects": "可能会影响部分设备的性能。", // ui-settings.tsx
    "settings.field.media_blurred_bg": "媒体页模糊背景", // ui-settings.tsx
    "settings.help.media_blurred_bg": "可能导致性能问题。", // ui-settings.tsx
    "settings.field.anime_unwatched_count": "动漫卡片未看计数", // ui-settings.tsx
    "settings.field.manga_unread_count": "漫画卡片未读计数", // ui-settings.tsx
    "settings.field.media_glassy_bg": "媒体卡片毛玻璃背景", // ui-settings.tsx
    "settings.field.episode_legacy_layout": "剧集卡片：传统布局", // ui-settings.tsx
    "settings.field.episode_hide_summary": "剧集条目：隐藏简介", // ui-settings.tsx
    "settings.field.episode_hide_filename": "剧集条目：隐藏文件名", // ui-settings.tsx
    "settings.field.disable_carousel_autoscroll": "禁用轮播自动滚动", // ui-settings.tsx
    "settings.field.smaller_carousel_cards": "更小的轮播剧集卡片", // ui-settings.tsx
    "settings.field.expand_sidebar_hover": "悬停时展开侧边栏", // ui-settings.tsx
    "settings.help.expand_sidebar": "可能导致插件托盘出现视觉故障。", // ui-settings.tsx
    "settings.field.disable_sidebar_transparency": "禁用侧边栏透明", // ui-settings.tsx
    "settings.field.hide_top_navbar_web": "隐藏顶部导航栏（网页端）", // ui-settings.tsx
    "settings.field.hide_top_navbar": "隐藏顶部导航栏", // ui-settings.tsx
    "settings.help.hide_top_navbar": "切换为仅侧边栏模式。", // ui-settings.tsx
    "settings.field.unpinned_menu_items": "未固定的菜单项", // ui-settings.tsx
    "settings.ui.no_items_selected": "未选择任何项", // ui-settings.tsx Combobox 空状态
    "settings.toast.light_theme_unsupported": "Seanime 不支持浅色主题", // ui-settings.tsx
    "settings.nav.auto_downloader_menu": "自动下载器", // ui-settings.tsx 菜单固定选项（非侧栏导航文案，navigation.sidebar.auto_downloader 语义不同）

    // —— 07-01 追加：logs-settings.tsx ——
    "settings.logs.column_name": "名称", // logs-settings.tsx 表格列头
    "settings.logs.copy_server_logs": "复制当前服务器日志", // logs-settings.tsx
    "settings.logs.open_logs_dir": "打开日志目录", // logs-settings.tsx
    "settings.logs.record_issue": "记录问题", // logs-settings.tsx
    "settings.logs.delete_selected": "删除所选", // logs-settings.tsx
    "settings.logs.copy_to_clipboard": "复制到剪贴板", // logs-settings.tsx
    "settings.toast.copied_to_clipboard": "已复制到剪贴板", // logs-settings.tsx
    "settings.logs.filter_server": "服务器", // logs-settings.tsx 筛选选项
    "settings.logs.filter_scanner": "扫描器", // logs-settings.tsx
    "settings.logs.profiling_title": "性能分析", // logs-settings.tsx
    "settings.logs.memory_statistics": "内存统计", // logs-settings.tsx
    "settings.logs.force_gc": "强制垃圾回收", // logs-settings.tsx
    "settings.logs.heap_allocated": "堆已分配", // logs-settings.tsx
    "settings.logs.heap_in_use": "堆使用中", // logs-settings.tsx
    "settings.logs.heap_system": "堆系统内存", // logs-settings.tsx
    "settings.logs.total_allocated": "累计分配", // logs-settings.tsx
    "settings.logs.goroutines": "Goroutine 数量", // logs-settings.tsx
    "settings.logs.gc_cycles": "GC 次数", // logs-settings.tsx
    "settings.logs.click_refresh_hint": "点击「刷新」加载内存统计", // logs-settings.tsx
    "settings.logs.memory": "内存", // logs-settings.tsx
    "settings.logs.heap_profile": "堆内存 Profile", // logs-settings.tsx
    "settings.logs.allocations_profile": "内存分配 Profile", // logs-settings.tsx
    "settings.logs.goroutine_profile": "Goroutine Profile", // logs-settings.tsx
    "settings.logs.duration_seconds": "持续时间（秒）", // logs-settings.tsx
    "settings.logs.download_cpu_profile": "下载 CPU Profile", // logs-settings.tsx
    "settings.logs.cpu_hint": "CPU 分析将运行指定的时长（1-300 秒）", // logs-settings.tsx

    // —— 07-01 追加：debrid-settings.tsx ——
    "settings.debrid.title": "Debrid 服务", // debrid-settings.tsx
    "settings.debrid.desc": "配置您的 Debrid 服务集成", // debrid-settings.tsx
    "settings.debrid.auto_downloader_alert_title": "自动下载器未使用 Debrid", // debrid-settings.tsx
    "settings.debrid.auto_downloader_alert_before": "自动下载器已启用但未使用 Debrid。请在", // debrid-settings.tsx
    "settings.debrid.auto_downloader_alert_link": "自动下载器设置", // debrid-settings.tsx
    "settings.debrid.auto_downloader_alert_after": "中配置使用您的 Debrid 服务。", // debrid-settings.tsx
    "settings.debrid.provider": "提供商", // debrid-settings.tsx 卡片标题与字段标签共用
    "settings.debrid.provider_torbox": "TorBox", // debrid-settings.tsx（专有名词保留）
    "settings.debrid.provider_realdebrid": "Real-Debrid", // debrid-settings.tsx（专有名词保留）
    "settings.debrid.provider_alldebrid": "AllDebrid", // debrid-settings.tsx（专有名词保留）
    "settings.debrid.provider_premiumize": "Premiumize", // debrid-settings.tsx（专有名词保留）
    "settings.debrid.provider_dummy": "Dummy Debrid", // debrid-settings.tsx（专有名词保留）
    "settings.debrid.api_key": "API 密钥", // debrid-settings.tsx
    "settings.debrid.streaming_title": "Debrid 流媒体", // debrid-settings.tsx
    "settings.debrid.streaming_desc": "配置如何从您的 Debrid 服务串流播放剧集", // debrid-settings.tsx
    "settings.debrid.home_screen_title": "首页", // debrid-settings.tsx
    "settings.debrid.include_streaming": "在动漫列表中包含流媒体", // debrid-settings.tsx（torrentstream 同名词段复用）
    "settings.debrid.include_streaming_help": "在您的动漫列表中展示正在观看的流媒体标题。", // debrid-settings.tsx
    "settings.debrid.autoselect_title": "自动选择", // debrid-settings.tsx
    "settings.debrid.autoselect_help": "让 Seanime 基于缓存与分辨率自动寻找最佳种子。", // debrid-settings.tsx
    "settings.debrid.preferred_resolution": "首选分辨率", // debrid-settings.tsx（torrentstream 同名词段复用）
    "settings.debrid.preferred_resolution_help": "启用自动选择后，Seanime 将尝试寻找此分辨率的种子。", // debrid-settings.tsx
    "settings.debrid.error_expected_array": "应为一个 JSON 数组", // debrid-settings.tsx 校验消息
    "settings.debrid.error_invalid_json": "JSON 格式无效", // debrid-settings.tsx 校验消息
    "settings.toast.settings_saved": "设置已保存", // debrid-settings.tsx
    "settings.toast.dummy_profile_saved": "Dummy 配置文件已保存", // debrid-settings.tsx
    "settings.debrid.profile_title": "Dummy Debrid 配置文件", // debrid-settings.tsx
    "settings.debrid.enable_profile": "启用配置文件", // debrid-settings.tsx
    "settings.debrid.profile_name": "配置文件名称", // debrid-settings.tsx
    "settings.debrid.fallback_mkv_path": "备用 MKV 路径", // debrid-settings.tsx
    "settings.debrid.cache_available": "缓存可用", // debrid-settings.tsx
    "settings.debrid.files": "文件", // debrid-settings.tsx
    "settings.debrid.network_title": "Dummy 网络", // debrid-settings.tsx
    "settings.debrid.ready_delay": "就绪延迟 (ms)", // debrid-settings.tsx
    "settings.debrid.progress_interval": "进度间隔 (ms)", // debrid-settings.tsx
    "settings.debrid.first_byte_delay": "首字节延迟 (ms)", // debrid-settings.tsx
    "settings.debrid.bandwidth": "带宽 (B/s)", // debrid-settings.tsx
    "settings.debrid.chunk_size": "分块大小 (字节)", // debrid-settings.tsx
    "settings.debrid.jitter": "抖动 (ms)", // debrid-settings.tsx

    // —— 07-01 追加：local-settings.tsx / anilist-settings.tsx ——
    "settings.local.title": "本地账户", // local-settings.tsx
    "settings.local.desc": "由 Seanime 管理的本地动漫与漫画列表", // local-settings.tsx
    "settings.local.auto_sync": "从 AniList 自动同步", // local-settings.tsx
    "settings.local.auto_sync_help": "定期使用您的 AniList 数据更新本地收藏。", // local-settings.tsx
    "settings.anilist.desc": "管理您的 AniList 账户", // anilist-settings.tsx

    // —— 07-01 追加：data-settings.tsx ——
    "settings.data.title": "本地文件", // data-settings.tsx
    "settings.data.desc": "已扫描的本地文件数据。", // data-settings.tsx
    "settings.data.export": "导出本地文件数据", // data-settings.tsx
    "settings.data.toast_export_failed": "生成导出令牌失败", // data-settings.tsx
    "settings.data.import_title": "导入本地文件", // data-settings.tsx 弹窗标题与按钮共用
    "settings.data.import_warning": "这将覆盖您现有的媒体库数据，请确保已有备份。", // data-settings.tsx
    "settings.data.file_path": "数据文件路径", // data-settings.tsx
    "settings.data.file_path_help": "包含本地文件数据的 JSON 文件路径。", // data-settings.tsx

    // —— 07-01 追加：filecache-settings.tsx ——
    "settings.filecache.show_total": "显示总大小", // filecache-settings.tsx
    "settings.filecache.clear_manga": "清除漫画缓存", // filecache-settings.tsx
    "settings.filecache.clear_mediastream": "清除媒体流缓存", // filecache-settings.tsx
    "settings.filecache.clear_onlinestream": "清除在线流媒体缓存", // filecache-settings.tsx

    // —— 07-01 追加：denshi-settings.tsx ——
    "settings.denshi.window_title": "窗口", // denshi-settings.tsx
    "settings.denshi.minimize_to_tray": "关闭时最小化到托盘", // denshi-settings.tsx
    "settings.denshi.minimize_to_tray_help": "启用后，关闭窗口将把应用最小化到系统托盘而不是退出。", // denshi-settings.tsx
    "settings.denshi.open_background": "后台启动", // denshi-settings.tsx
    "settings.denshi.open_background_help": "启用后，应用启动时隐藏，可从系统托盘显示。", // denshi-settings.tsx
    "settings.denshi.system_title": "系统", // denshi-settings.tsx
    "settings.denshi.open_at_launch": "开机自启", // denshi-settings.tsx
    "settings.denshi.open_at_launch_help_linux": "此功能在 Linux 上不受支持。", // denshi-settings.tsx
    "settings.denshi.open_at_launch_help": "启用后，登录电脑时应用将自动启动。", // denshi-settings.tsx
    "settings.denshi.auto_save_note": "设置会自动保存，并在重启后生效", // denshi-settings.tsx

    // —— 07-01 追加：discord-rich-presence-settings.tsx ——
    "settings.discord.rich_presence": "状态展示", // discord-rich-presence-settings.tsx（Rich Presence）
    "settings.discord.rich_presence_desc": "在 Discord 中展示您正在观看或阅读的内容。", // discord-rich-presence-settings.tsx
    "settings.discord.hide_repo_button": "隐藏 Seanime 仓库按钮", // discord-rich-presence-settings.tsx
    "settings.discord.show_profile_button": "显示 AniList 资料页按钮", // discord-rich-presence-settings.tsx
    "settings.discord.show_profile_button_help": "显示一个打开您 AniList 资料页的按钮。", // discord-rich-presence-settings.tsx

    // —— 07-01 追加：settings-submit-button.tsx ——
    "settings.toast.unsaved_changes": "您有未保存的更改。", // settings-submit-button.tsx
} satisfies Dictionary
