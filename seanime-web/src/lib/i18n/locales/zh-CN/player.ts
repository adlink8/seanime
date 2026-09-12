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
    // ===== video-core（07-03）=====
    "player.common.on": "开", // video-core-settings-menu.tsx
    "player.common.off": "关", // video-core-settings-menu.tsx
    "player.common.unknown_error": "未知错误", // video-core-cast.tsx, video-core-events.ts, video-core-anime-4k-manager.ts
    // 偏好设置弹窗 video-core-preferences.tsx
    "player.prefs.title": "偏好设置", // video-core-preferences.tsx, video-core-settings-menu.tsx
    "player.prefs.tab_keybinds": "键盘快捷键", // video-core-preferences.tsx
    "player.prefs.tab_subtitles_audio": "字幕与音频", // video-core-preferences.tsx
    "player.prefs.tab_general": "常规", // video-core-preferences.tsx
    "player.prefs.tab_translation": "翻译", // video-core-preferences.tsx
    "player.prefs.reset_all": "全部重置", // video-core-preferences.tsx
    "player.prefs.extra_skip_chapters": "额外跳过章节规则", // video-core-preferences.tsx
    "player.prefs.extra_skip_chapters_help": "以逗号分隔的正则表达式（不区分大小写）。既有的片头/片尾规则仍然生效。", // video-core-preferences.tsx
    "player.prefs.screenshot_dir": "截图目录", // video-core-preferences.tsx
    "player.prefs.screenshot_dir_help": "配置截图保存的目录", // video-core-preferences.tsx
    "player.prefs.must_be_absolute_path": "必须是绝对路径", // video-core-preferences.tsx
    "player.prefs.press_key": "按下按键...", // video-core-preferences.tsx
    "player.prefs.kb_value_seconds": "秒", // video-core-preferences.tsx
    "player.prefs.kb_value_increment": "增量", // video-core-preferences.tsx
    "player.prefs.kb_value_percent": "百分比", // video-core-preferences.tsx
    "player.prefs.kb_seek_forward_fine": "快进（微调）", // video-core-preferences.tsx
    "player.prefs.kb_seek_backward_fine": "快退（微调）", // video-core-preferences.tsx
    "player.prefs.kb_seek_forward": "快进", // video-core-preferences.tsx
    "player.prefs.kb_seek_backward": "快退", // video-core-preferences.tsx
    "player.prefs.kb_increase_speed": "提升播放速度", // video-core-preferences.tsx
    "player.prefs.kb_decrease_speed": "降低播放速度", // video-core-preferences.tsx
    "player.prefs.kb_next_chapter": "跳到下一章节", // video-core-preferences.tsx
    "player.prefs.kb_previous_chapter": "跳到上一章节", // video-core-preferences.tsx
    "player.prefs.kb_next_episode": "播放下一集", // video-core-preferences.tsx
    "player.prefs.kb_previous_episode": "播放上一集", // video-core-preferences.tsx
    "player.prefs.kb_cycle_subtitles": "切换字幕轨道", // video-core-preferences.tsx
    "player.prefs.kb_fullscreen": "切换全屏", // video-core-preferences.tsx
    "player.prefs.kb_pip": "切换画中画", // video-core-preferences.tsx
    "player.prefs.kb_take_screenshot": "截图", // video-core-preferences.tsx
    "player.prefs.kb_volume_up": "增大音量", // video-core-preferences.tsx
    "player.prefs.kb_volume_down": "减小音量", // video-core-preferences.tsx
    "player.prefs.kb_mute": "切换静音", // video-core-preferences.tsx
    "player.prefs.kb_cycle_audio": "切换音轨", // video-core-preferences.tsx
    "player.prefs.kb_display_characters": "显示角色面板", // video-core-preferences.tsx
    "player.prefs.kb_stats_for_nerds": "切换技术统计面板", // video-core-preferences.tsx
    "player.prefs.defaults": "默认设置", // video-core-preferences.tsx
    "player.prefs.preferred_sub_language": "首选字幕语言", // video-core-preferences.tsx
    "player.prefs.preferred_audio_language": "首选音频语言", // video-core-preferences.tsx
    "player.prefs.ignored_sub_names": "忽略的字幕名称", // video-core-preferences.tsx
    "player.prefs.ignored_sub_names_help": "匹配首选语言时默认不会选中的字幕轨道。多个名称用逗号分隔。", // video-core-preferences.tsx
    "player.prefs.rendering": "渲染", // video-core-preferences.tsx
    "player.prefs.convert_soft_subs": "将软字幕转换为 ASS", // video-core-preferences.tsx
    "player.prefs.convert_soft_subs_help": "播放器会把其他字幕格式（SRT、VTT 等）转换为 ASS。如果你的语言不受支持，可以添加新字体或关闭此功能。更改此设置后需要重新加载播放器。", // video-core-preferences.tsx
    "player.prefs.translation_saved": "翻译设置已保存", // video-core-preferences.tsx
    "player.prefs.enable_translation": "启用翻译", // video-core-preferences.tsx
    "player.prefs.enable_translation_help": "自动将字幕轨道翻译为所选语言", // video-core-preferences.tsx
    "player.prefs.provider": "翻译服务", // video-core-preferences.tsx
    "player.prefs.provider_google_free": "Google 免费", // video-core-preferences.tsx
    "player.prefs.provider_openai_compatible": "OpenAI 兼容", // video-core-preferences.tsx
    "player.prefs.deepl_limit": "DeepL 不支持所有目标语言。", // video-core-preferences.tsx
    "player.prefs.target_language": "目标语言", // video-core-preferences.tsx
    "player.prefs.target_language_help": "选择字幕要翻译成的语言", // video-core-preferences.tsx
    "player.prefs.base_url": "接口地址", // video-core-preferences.tsx
    "player.prefs.base_url_help": "OpenAI 兼容的 /v1 接口，例如 LM Studio 或 Ollama。", // video-core-preferences.tsx
    "player.prefs.model": "模型", // video-core-preferences.tsx
    "player.prefs.api_key": "API 密钥", // video-core-preferences.tsx
    "player.prefs.api_key_optional": "API 密钥（可选）", // video-core-preferences.tsx
    "player.prefs.api_key_placeholder": "输入你的 API 密钥", // video-core-preferences.tsx
    "player.prefs.reload_notice": "仅在切换翻译服务、语言、接口地址、模型或 API 密钥后，才需要重新加载播放器。", // video-core-preferences.tsx
    "player.prefs.lang_en_us": "英语（美国）", // video-core-preferences.tsx
    "player.prefs.lang_en_gb": "英语（英国）", // video-core-preferences.tsx
    "player.prefs.lang_es": "西班牙语", // video-core-preferences.tsx
    "player.prefs.lang_fr": "法语", // video-core-preferences.tsx
    "player.prefs.lang_de": "德语", // video-core-preferences.tsx
    "player.prefs.lang_it": "意大利语", // video-core-preferences.tsx
    "player.prefs.lang_pt_br": "葡萄牙语（巴西）", // video-core-preferences.tsx
    "player.prefs.lang_pt_pt": "葡萄牙语（葡萄牙）", // video-core-preferences.tsx
    "player.prefs.lang_ru": "俄语", // video-core-preferences.tsx
    "player.prefs.lang_ja": "日语", // video-core-preferences.tsx
    "player.prefs.lang_ko": "韩语", // video-core-preferences.tsx
    "player.prefs.lang_zh_hans": "中文（简体）", // video-core-preferences.tsx
    "player.prefs.lang_zh_hant": "中文（繁体）", // video-core-preferences.tsx
    "player.prefs.lang_ar": "阿拉伯语", // video-core-preferences.tsx
    "player.prefs.lang_tr": "土耳其语", // video-core-preferences.tsx
    "player.prefs.lang_pl": "波兰语", // video-core-preferences.tsx
    "player.prefs.lang_nl": "荷兰语", // video-core-preferences.tsx
    "player.prefs.lang_sv": "瑞典语", // video-core-preferences.tsx
    "player.prefs.lang_nb": "挪威语", // video-core-preferences.tsx
    "player.prefs.lang_da": "丹麦语", // video-core-preferences.tsx
    "player.prefs.lang_fi": "芬兰语", // video-core-preferences.tsx
    "player.prefs.lang_el": "希腊语", // video-core-preferences.tsx
    "player.prefs.lang_cs": "捷克语", // video-core-preferences.tsx
    "player.prefs.lang_hu": "匈牙利语", // video-core-preferences.tsx
    "player.prefs.lang_ro": "罗马尼亚语", // video-core-preferences.tsx
    "player.prefs.lang_id": "印尼语", // video-core-preferences.tsx
    "player.prefs.lang_uk": "乌克兰语", // video-core-preferences.tsx
    "player.prefs.lang_bg": "保加利亚语", // video-core-preferences.tsx
    "player.prefs.lang_sk": "斯洛伐克语", // video-core-preferences.tsx
    "player.prefs.lang_sl": "斯洛文尼亚语", // video-core-preferences.tsx
    "player.prefs.lang_et": "爱沙尼亚语", // video-core-preferences.tsx
    "player.prefs.lang_lv": "拉脱维亚语", // video-core-preferences.tsx
    "player.prefs.lang_lt": "立陶宛语", // video-core-preferences.tsx
    "player.prefs.lang_hi": "印地语", // video-core-preferences.tsx
    "player.prefs.lang_bn": "孟加拉语", // video-core-preferences.tsx
    "player.prefs.lang_ta": "泰米尔语", // video-core-preferences.tsx
    "player.prefs.lang_te": "泰卢固语", // video-core-preferences.tsx
    "player.prefs.lang_mr": "马拉地语", // video-core-preferences.tsx
    "player.prefs.lang_kn": "卡纳达语", // video-core-preferences.tsx
    "player.prefs.lang_ml": "马拉雅拉姆语", // video-core-preferences.tsx
    "player.prefs.lang_pa": "旁遮普语", // video-core-preferences.tsx
    "player.prefs.lang_fa": "波斯语", // video-core-preferences.tsx
    "player.prefs.lang_ur": "乌尔都语", // video-core-preferences.tsx
    "player.prefs.lang_sw": "斯瓦希里语", // video-core-preferences.tsx
    "player.prefs.lang_af": "南非荷兰语", // video-core-preferences.tsx
    "player.prefs.lang_ms": "马来语", // video-core-preferences.tsx
    "player.prefs.lang_hr": "克罗地亚语", // video-core-preferences.tsx
    "player.prefs.lang_sr": "塞尔维亚语", // video-core-preferences.tsx
    "player.prefs.lang_he": "希伯来语", // video-core-preferences.tsx
    "player.prefs.lang_th": "泰语", // video-core-preferences.tsx
    "player.prefs.lang_vi": "越南语", // video-core-preferences.tsx
    // 播放器设置菜单 video-core-settings-menu.tsx
    "player.menu.title": "设置", // video-core-settings-menu.tsx
    "player.menu.playback_speed": "倍速", // video-core-settings-menu.tsx
    "player.menu.auto_play": "自动播放", // video-core-settings-menu.tsx
    "player.menu.auto_next": "自动连播", // video-core-settings-menu.tsx
    "player.menu.skip_op_ed": "跳过 OP/ED", // video-core-settings-menu.tsx
    "player.menu.subtitle_delay": "字幕延迟", // video-core-settings-menu.tsx, video-core-subtitle-menu.tsx
    "player.menu.subtitle_styles": "字幕样式", // video-core-settings-menu.tsx, video-core-subtitle-menu.tsx
    "player.menu.caption_styles": "隐藏式字幕样式", // video-core-settings-menu.tsx
    "player.menu.player_appearance": "播放器外观", // video-core-settings-menu.tsx
    "player.menu.options": "选项", // video-core-settings-menu.tsx
    "player.menu.sub_styles_notice": "字幕自定义不会覆盖包含多种样式的 ASS/SSA 轨道。", // video-core-settings-menu.tsx
    "player.menu.caption_notice": "仅对非 ASS 字幕生效。", // video-core-settings-menu.tsx
    "player.menu.subtitle_delay_hint": "正值延迟字幕，负值提前字幕。", // video-core-settings-menu.tsx
    "player.menu.anime4k_desc": "实时锐化，GPU 开销较大。", // video-core-settings-menu.tsx
    "player.menu.performance_heavy": "高负载", // video-core-settings-menu.tsx
    "player.menu.show_chapter_markers": "显示章节标记", // video-core-settings-menu.tsx
    "player.menu.highlight_skipped_chapters": "高亮已跳过章节", // video-core-settings-menu.tsx
    "player.menu.increase_saturation": "提高饱和度", // video-core-settings-menu.tsx
    "player.menu.font": "字体", // video-core-settings-menu.tsx
    "player.menu.default": "默认", // video-core-settings-menu.tsx
    "player.menu.font_size": "字体大小", // video-core-settings-menu.tsx
    "player.menu.text_color": "文字颜色", // video-core-settings-menu.tsx
    "player.menu.outline": "描边", // video-core-settings-menu.tsx
    "player.menu.outline_width": "描边宽度", // video-core-settings-menu.tsx
    "player.menu.outline_color": "描边颜色", // video-core-settings-menu.tsx
    "player.menu.shadow": "阴影", // video-core-settings-menu.tsx
    "player.menu.shadow_depth": "阴影深度", // video-core-settings-menu.tsx
    "player.menu.shadow_opacity": "阴影不透明度", // video-core-settings-menu.tsx
    "player.menu.shadow_color": "阴影颜色", // video-core-settings-menu.tsx
    "player.menu.background": "背景", // video-core-settings-menu.tsx
    "player.menu.background_opacity": "背景不透明度", // video-core-settings-menu.tsx
    "player.menu.background_color": "背景颜色", // video-core-settings-menu.tsx
    "player.menu.text_shadow": "文字阴影", // video-core-settings-menu.tsx
    "player.menu.custom_font": "自定义字体", // video-core-settings-menu.tsx
    "player.menu.custom_font_help_prefix": "将字体文件放到", // video-core-settings-menu.tsx
    "player.menu.assets_dir": "Seanime 资源目录", // video-core-settings-menu.tsx
    "player.menu.custom_font_help_suffix": "。文件名必须与字体名完全一致。", // video-core-settings-menu.tsx
    "player.menu.file_name": "文件名", // video-core-settings-menu.tsx
    "player.menu.file_name_help": "示例：Noto Sans JP.woff2", // video-core-settings-menu.tsx
    "player.menu.on_with_font": "开，字体", // video-core-settings-menu.tsx
    "player.menu.size_small": "小", // video-core-settings-menu.tsx
    "player.menu.size_medium": "中", // video-core-settings-menu.tsx
    "player.menu.size_large": "大", // video-core-settings-menu.tsx
    "player.menu.size_extra_large": "特大", // video-core-settings-menu.tsx
    "player.menu.none": "无", // video-core-settings-menu.tsx
    "player.menu.color_white": "白色", // video-core-settings-menu.tsx
    "player.menu.color_black": "黑色", // video-core-settings-menu.tsx
    "player.menu.color_gray": "灰色", // video-core-settings-menu.tsx
    "player.menu.color_yellow": "黄色", // video-core-settings-menu.tsx
    "player.menu.color_cyan": "青色", // video-core-settings-menu.tsx
    "player.menu.color_pink": "粉色", // video-core-settings-menu.tsx
    "player.menu.color_purple": "紫色", // video-core-settings-menu.tsx
    "player.menu.color_lime": "青柠色", // video-core-settings-menu.tsx
    "player.menu.quality": "画质", // video-core-resolution-menu.tsx
    "player.menu.auto": "自动", // video-core-resolution-menu.tsx
    "player.menu.audio": "音频", // video-core-audio-menu.tsx
    "player.subs.title": "字幕", // video-core-subtitle-menu.tsx
    "player.subs.drop_hint": "可以将字幕文件拖放到播放器上来添加字幕。", // video-core-subtitle-menu.tsx
    "player.subs.adding_file": "正在添加字幕文件...", // video-core-events.ts
    "player.subs.unknown_type": "未知的字幕类型", // video-core-events.ts
    "player.subs.libass_init_error": "libass 渲染器初始化出错：{error}", // video-core-subtitles.ts
    "player.subs.load_failed": "字幕轨道加载失败", // video-core-subtitles.ts
    "player.subs.load_failed_with_error": "字幕轨道加载失败：{error}", // video-core-subtitles.ts
    "player.subs.convert_failed": "字幕轨道转换失败", // video-core-subtitles.ts
    "player.subs.drop_target_hint": "将字幕文件拖放到播放器上即可加载", // video-core-subtitle-menu.tsx
    // 技术统计 video-core-stats.tsx
    "player.stats.title": "技术统计", // video-core-stats.tsx
    "player.stats.file": "文件", // video-core-stats.tsx
    "player.stats.path": "路径", // video-core-stats.tsx
    "player.stats.stream": "流地址", // video-core-stats.tsx
    "player.stats.media_title": "标题", // video-core-stats.tsx
    "player.stats.episode": "剧集", // video-core-stats.tsx
    "player.stats.mime_type": "MIME 类型", // video-core-stats.tsx
    "player.stats.video_codec": "视频编码", // video-core-stats.tsx
    "player.stats.audio_codecs": "音频编码", // video-core-stats.tsx
    "player.stats.display_video": "显示 / 视频", // video-core-stats.tsx
    "player.stats.framerate": "帧率", // video-core-stats.tsx
    "player.stats.frames": "帧数（总计 / 丢弃）", // video-core-stats.tsx
    "player.stats.render_time": "渲染耗时", // video-core-stats.tsx
    "player.stats.buffer_ahead": "前向缓冲", // video-core-stats.tsx
    "player.stats.playback_rate": "播放倍速", // video-core-stats.tsx
    "player.stats.decoded_corrupted": "解码 / 损坏", // video-core-stats.tsx
    "player.stats.network_state": "网络状态", // video-core-stats.tsx
    "player.stats.ready_state": "就绪状态", // video-core-stats.tsx
    "player.stats.time_duration": "时间 / 时长", // video-core-stats.tsx
    "player.stats.anime4k_mode": "Anime4K 模式", // video-core-stats.tsx
    "player.stats.a4k_framerate": "Anime4K 帧率", // video-core-stats.tsx
    "player.stats.a4k_frame_drops": "Anime4K 掉帧数", // video-core-stats.tsx
    // Anime4K video-core-anime-4k.ts / video-core-anime-4k-manager.ts
    "player.a4k.disabled": "已禁用", // video-core-anime-4k.ts
    "player.a4k.mode_a": "模式 A", // video-core-anime-4k.ts
    "player.a4k.mode_b": "模式 B", // video-core-anime-4k.ts
    "player.a4k.mode_c": "模式 C", // video-core-anime-4k.ts
    "player.a4k.mode_aa": "模式 A+A", // video-core-anime-4k.ts
    "player.a4k.mode_bb": "模式 B+B", // video-core-anime-4k.ts
    "player.a4k.mode_ca": "模式 C+A", // video-core-anime-4k.ts
    "player.a4k.desc_mode_a": "先去除压缩伪影再放大", // video-core-anime-4k.ts
    "player.a4k.desc_mode_b": "温和去除伪影后放大", // video-core-anime-4k.ts
    "player.a4k.desc_mode_c": "放大时降噪，然后再次放大", // video-core-anime-4k.ts
    "player.a4k.desc_mode_aa": "增强修复，画质更佳", // video-core-anime-4k.ts
    "player.a4k.desc_mode_bb": "双重柔和修复", // video-core-anime-4k.ts
    "player.a4k.desc_mode_ca": "降噪 + 修复混合", // video-core-anime-4k.ts
    "player.a4k.desc_cnn_2x_m": "速度与画质均衡", // video-core-anime-4k.ts
    "player.a4k.desc_cnn_2x_vl": "高质量神经网络", // video-core-anime-4k.ts
    "player.a4k.desc_denoise_cnn_2x_vl": "放大的同时去除噪点", // video-core-anime-4k.ts
    "player.a4k.desc_cnn_2x_ul": "最高 CNN 画质", // video-core-anime-4k.ts
    "player.a4k.desc_gan_3x_l": "生成对抗网络，感知画质更佳", // video-core-anime-4k.ts
    "player.a4k.desc_gan_4x_ul": "利用 GAN 技术实现最大倍放大", // video-core-anime-4k.ts
    "player.a4k.recommended_reason": "推荐用于 {gpu} 上的 {height}p 视频", // video-core-anime-4k.ts
    "player.a4k.performance_degraded": "性能下降，已关闭 Anime4K。", // video-core-anime-4k-manager.ts
    // 播放覆盖提示（overlay feedback）
    "player.overlay.beginning": "已回到开头", // video-core-preferences.tsx
    "player.overlay.end": "已跳至结尾", // video-core-preferences.tsx
    "player.overlay.prev_frame": "上一帧", // video-core-preferences.tsx
    "player.overlay.next_frame": "下一帧", // video-core-preferences.tsx
    "player.overlay.skipped_opening": "已跳过片头", // video-core-preferences.tsx
    "player.overlay.skipped_ending": "已跳过片尾", // video-core-preferences.tsx
    "player.overlay.skipped": "已跳过 {label}", // video-core-time-range.tsx
    "player.overlay.speed": "倍速：{speed}x", // video-core-preferences.tsx
    "player.overlay.chapter": "章节：{name}", // video-core-preferences.tsx
    "player.overlay.chapter_num": "第 {num} 章", // video-core-preferences.tsx
    "player.overlay.end_of_chapters": "已到达章节末尾", // video-core-preferences.tsx
    "player.overlay.subtitles": "字幕：{track}", // video-core-preferences.tsx
    "player.overlay.subtitles_off": "字幕：关", // video-core-preferences.tsx
    "player.overlay.audio": "音频：{track}", // video-core-preferences.tsx
    "player.overlay.track": "轨道 {num}", // video-core-preferences.tsx, video-core-audio-menu.tsx, video-core-hls.ts
    "player.overlay.no_audio_tracks": "没有其他音轨", // video-core-preferences.tsx
    "player.overlay.no_next_episode": "没有下一集", // video-core-preferences.tsx
    "player.overlay.no_prev_episode": "没有上一集", // video-core-preferences.tsx
    "player.overlay.progress_restored": "进度已恢复", // video-core.tsx
    "player.overlay.restoring_progress": "正在恢复进度", // video-core.tsx
    "player.overlay.subtitle_track_added": "已添加字幕轨道：{label}", // video-core-events.ts
    "player.watch_party.new_message": "新聊天消息（{count}）", // video-core-watch-party-chat.tsx
    // 终止播放确认 video-core.tsx
    "player.terminate.title": "终止播放？", // video-core.tsx
    "player.terminate.desc": "再次按 Esc 或选择“终止播放”以停止播放。", // video-core.tsx
    "player.terminate.confirm": "终止播放", // video-core.tsx
    "player.terminate.keep": "继续播放", // video-core.tsx
    // 投屏 video-core-cast.tsx
    "player.cast.title": "投屏到设备", // video-core-cast.tsx
    "player.cast.casting": "正在投屏", // video-core-cast.tsx
    "player.cast.casting_to": "正在投屏到 {device}", // video-core-cast.tsx
    "player.cast.connected_to": "已连接到", // video-core-cast.tsx
    "player.cast.connected": "已连接到 Chromecast", // video-core-cast.tsx
    "player.cast.disconnected": "已断开 Chromecast 连接", // video-core-cast.tsx
    "player.cast.disconnect": "断开连接", // video-core-cast.tsx
    "player.cast.searching": "正在搜索设备...", // video-core-cast.tsx
    "player.cast.no_devices": "未找到设备", // video-core-cast.tsx
    "player.cast.scan_again": "重新扫描", // video-core-cast.tsx
    "player.cast.stop": "停止", // video-core-cast.tsx
    "player.cast.play": "播放", // video-core-cast.tsx
    "player.cast.pause": "暂停", // video-core-cast.tsx
    "player.cast.error": "投屏出错：{message}", // video-core-cast.tsx
    "player.cast.connect_failed": "连接失败：{message}", // video-core-cast.tsx
    // 画中画 video-core-pip.ts
    "player.pip.exit": "退出画中画", // video-core.tsx
    "player.pip.enter_error": "进入画中画时发生未知错误", // video-core-pip.ts
    "player.pip.exit_failed": "退出画中画失败：{message}", // video-core-pip.ts
    // 错误与播放状态 video-core.tsx / video-core-events.ts / video-core-audio.ts
    "player.playback.error_occurred": "视频播放出错（代码：{code}）", // video-core.tsx
    "player.playback.stalled_while_buffering": "缓冲时播放停滞", // video-core.tsx
    "player.error.aborted": "媒体播放已中止", // video-core-events.ts
    "player.error.network": "发生网络错误：请查看控制台与网络面板了解详情", // video-core-events.ts
    "player.error.decode": "媒体解码错误：编解码器不受支持或文件已损坏", // video-core-events.ts
    "player.error.decode_detail": "这可能是编解码器兼容性问题。", // video-core-events.ts
    "player.error.format_not_supported": "媒体格式不受支持", // video-core-events.ts
    "player.error.format_detail": "不支持该视频编解码器/封装格式。", // video-core-events.ts
    "player.error.unknown_media": "未知媒体错误", // video-core-events.ts
    "player.audio.codec_unsupported": "播放器不支持此音频编解码器。请尝试其他文件或使用外部播放器。", // video-core-audio.ts
    // HLS video-core-hls.ts
    "player.hls.stalled": "HLS 卡顿：{message}", // video-core.tsx
    "player.hls.level": "清晰度 {num}", // video-core-hls.ts
    "player.hls.not_supported": "当前浏览器不支持 HLS 播放", // video-core-hls.ts
    // 截图 video-core-screenshot.ts / video-core-screenshot-prompt.tsx
    "player.screenshot.saved": "截图已保存", // video-core-screenshot.ts, video-core-screenshot-prompt.tsx
    "player.screenshot.failed": "截图失败", // video-core-screenshot.ts
    "player.screenshot.taking": "正在截图...", // video-core-screenshot.ts
    "player.screenshot.save_server_failed": "保存截图到服务器失败", // video-core-screenshot.ts
    "player.screenshot.folder_saved": "截图目录已保存", // video-core-screenshot-prompt.tsx
    "player.screenshot.folder_save_failed": "保存截图目录失败", // video-core-screenshot-prompt.tsx
    // 角色面板 video-core-in-sight.tsx
    "player.insight.spoilers": "可能包含剧透。", // video-core-in-sight.tsx
    "player.insight.search_placeholder": "搜索角色...", // video-core-in-sight.tsx
    "player.insight.no_results": "无结果", // video-core-in-sight.tsx
    // 剧集列表 video-core-inline-helpers.tsx / _components/episode-pills-grid.tsx
    "player.inline.update_progress": "更新进度", // video-core-inline-helpers.tsx
    "player.episodes.episode_num": "第 {num} 集", // _components/episode-pills-grid.tsx
    "player.episodes.filler": "原创剧情集", // _components/episode-pills-grid.tsx
} satisfies Dictionary
