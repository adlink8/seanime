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
} satisfies Dictionary
