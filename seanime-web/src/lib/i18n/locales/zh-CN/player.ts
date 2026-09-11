import type { Dictionary } from "../../types"

/**
 * player 词条表
 * 归属：播放偏好、桌面播放器 / 外部播放器关联、连播与流式播放设置
 * 说明：video-core / mpv-core / onlinestream 目录当前 0 中文，
 *       故 seed 词条取自设置页中真实的播放器配置文案。
 */
export const playerDictionary = {
    "player.playback.title": "视频播放偏好", // src/app/(main)/settings/_components/playback-settings.tsx
    "player.playback.desc": "配置本设备的视频播放模式与引擎", // src/app/(main)/settings/_components/playback-settings.tsx
    "player.playback.device": "设备:", // src/app/(main)/settings/_components/playback-settings.tsx
    "player.local_media.title": "已下载媒体播放", // src/app/(main)/settings/_components/playback-settings.tsx
    "player.local_media.desc": "选择播放本地存储动漫文件的方式。", // src/app/(main)/settings/_components/playback-settings.tsx
    "player.mode.desktop_player": "桌面媒体播放器", // src/app/(main)/settings/_components/playback-settings.tsx
    "player.mode.desktop_player_desc": "使用配置的桌面播放器打开流并自动同步播放进度。", // src/app/(main)/settings/_components/playback-settings.tsx
    "player.mode.transcode": "转码 / 网页直连播放", // src/app/(main)/settings/_components/playback-settings.tsx
    "player.mode.transcode_desc": "通过网页端内嵌的 HTML5 播放器直接观看本地文件。", // src/app/(main)/settings/_components/playback-settings.tsx
    "player.mode.transcode_disabled_desc": "需先开启转码流服务方可使用网页播放器。", // src/app/(main)/settings/_components/playback-settings.tsx
    "player.mode.external_protocol": "外部播放器协议关联", // src/app/(main)/settings/_components/playback-settings.tsx
    "player.mode.external_protocol_desc": "使用自定义 URL Scheme 协议将流地址发送至外部播放器应用。", // src/app/(main)/settings/_components/playback-settings.tsx
    "player.status.disabled": "已禁用", // src/app/(main)/settings/_components/playback-settings.tsx
    "player.streaming.title": "种子与 Debrid 在线流式播放", // src/app/(main)/settings/_components/playback-settings.tsx
    "player.streaming.desc": "选择播放种子边下边播或 Debrid 云端流媒体的方式。", // src/app/(main)/settings/_components/playback-settings.tsx
    "player.desktop.title": "桌面播放器", // src/app/(main)/settings/_components/mediaplayer-settings.tsx
    "player.desktop.desc": "Seanime 内置支持与 MPV、VLC、IINA 以及 MPC-HC 播放器深度联动与进度跟踪。", // src/app/(main)/settings/_components/mediaplayer-settings.tsx
    "player.desktop.default_player": "默认播放器", // src/app/(main)/settings/_components/mediaplayer-settings.tsx
    "player.desktop.default_player_help": "用于播放视频文件并自动同步追番进度的播放器。", // src/app/(main)/settings/_components/mediaplayer-settings.tsx
    "player.autoplay.title": "连播设置", // src/app/(main)/settings/_components/mediaplayer-settings.tsx
    "player.autoplay.next_label": "自动连续播放下一集", // src/app/(main)/settings/_components/mediaplayer-settings.tsx
    "player.autoplay.next_help": "启用后，当当前剧集播放完毕后，Seanime 将在简短倒数后自动调起下一集。", // src/app/(main)/settings/_components/mediaplayer-settings.tsx
    "player.desktop.config_title": "播放器配置", // src/app/(main)/settings/_components/mediaplayer-settings.tsx
    "player.external.title": "外部播放器关联 (URL Scheme)", // src/app/(main)/settings/_components/mediaplayer-settings.tsx
    "player.external.desc": "通过自定义协议链接将流媒体发送至本机的第三方播放器。", // src/app/(main)/settings/_components/mediaplayer-settings.tsx
    "player.external.device_only": "仅适用于当前设备。", // src/app/(main)/settings/_components/mediaplayer-settings.tsx
    "player.external.scheme_label": "自定义协议 Scheme", // src/app/(main)/settings/_components/mediaplayer-settings.tsx
    "player.external.scheme_placeholder": "例如: outplayer://{url} 或 iina://weblink?url={url}", // src/app/(main)/settings/_components/mediaplayer-settings.tsx
} satisfies Dictionary
