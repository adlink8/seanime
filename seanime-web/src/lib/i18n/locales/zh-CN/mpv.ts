import type { Dictionary } from "../../types"

/**
 * mpv 词条表
 * 归属：mpv 播放器（app/(main)/_features/mpv-core/）
 * 独占：07-05
 *
 * 只读复用：common.* / player.* / media.* 等同名文案按契约复用，不新建重复词条。
 */
export const mpvDictionary = {
    // ===== 设置菜单 mpv-core-settings-menu.tsx =====
    "mpv.menu.preset": "预设", // mpv-core-settings-menu.tsx（Anime4K 预设分区标题）
    "mpv.menu.font_name": "字体名称", // mpv-core-settings-menu.tsx
    "mpv.menu.font_name_help": "示例：Noto Sans JP", // mpv-core-settings-menu.tsx
    "mpv.menu.subtitle_delay_hint": "调整字幕相对视频的时间。字幕出现太晚时请使用负值。", // mpv-core-settings-menu.tsx
    // ===== 着色器 mpv-core-settings-menu.tsx =====
    "mpv.shader.title": "着色器", // mpv-core-settings-menu.tsx
    "mpv.shader.anime4k_preset": "Anime4K 预设", // mpv-core-settings-menu.tsx
    "mpv.shader.anime4k_preset_desc": "使用 Anime4K 放大预设", // mpv-core-settings-menu.tsx
    "mpv.shader.custom_shaders": "自定义着色器", // mpv-core-settings-menu.tsx
    "mpv.shader.custom_shaders_desc": "启用单个自定义着色器", // mpv-core-settings-menu.tsx
    "mpv.shader.debanding": "去色带", // mpv-core-settings-menu.tsx
    "mpv.shader.debanding_help": "已在 MPV 配置中定义", // mpv-core-settings-menu.tsx
    "mpv.shader.select_shaders": "选择着色器", // mpv-core-settings-menu.tsx
    "mpv.shader.none_found": "目录中未找到着色器文件。", // mpv-core-settings-menu.tsx
    "mpv.shader.detected_count": "检测到 {count} 个着色器文件。", // mpv-core-settings-menu.tsx
    "mpv.shader.missing_files": "缺少：{files}", // mpv-core-settings-menu.tsx, mpv-core-player-inner.tsx
    "mpv.shader.open_folder": "打开文件夹", // mpv-core-settings-menu.tsx
    "mpv.shader.a4k_label_cnn_2x_medium": "CNN 2x 中等", // mpv-core-settings-menu.tsx
    "mpv.shader.a4k_label_cnn_2x_very_large": "CNN 2x 特大", // mpv-core-settings-menu.tsx
    "mpv.shader.a4k_label_denoise_cnn_2x_vl": "降噪 CNN 2x 特大", // mpv-core-settings-menu.tsx
    "mpv.shader.a4k_label_cnn_2x_ultra_large": "CNN 2x 超大", // mpv-core-settings-menu.tsx
    "mpv.shader.quality_fast": "快速（GPU 占用低）", // mpv-core-settings-menu.tsx
    "mpv.shader.quality_hq": "高画质（高负载）", // mpv-core-settings-menu.tsx
    // ===== MPV 配置 mpv-core-player-inner.tsx =====
    "mpv.config.unavailable": "当前 Denshi 构建中 MPV 配置文件不可用", // mpv-core-player-inner.tsx
    "mpv.config.write_failed": "写入 MPV 配置失败", // mpv-core-player-inner.tsx
    // ===== 字幕 mpv-core-player-inner.tsx =====
    "mpv.subs.apply_failed": "应用字幕设置失败", // mpv-core-player-inner.tsx
    "mpv.subs.unsupported_format": "不支持的字幕格式", // mpv-core-player-inner.tsx
    "mpv.subs.loaded_file": "已加载字幕 {name}", // mpv-core-player-inner.tsx
    // ===== 画中画 mpv-core-player-inner.tsx =====
    "mpv.pip.toggle_failed": "切换画中画失败", // mpv-core-player-inner.tsx
    // ===== 截图 mpv-core-player-inner.tsx / mpv-core-screenshot-prompt.tsx =====
    "mpv.screenshot.saved_to": "截图已保存至 {path}", // mpv-core-player-inner.tsx, mpv-core-screenshot-prompt.tsx
    // ===== 投屏 mpv-core-cast-button.tsx =====
    "mpv.cast.connected": "已连接", // mpv-core-cast-button.tsx
    // ===== 技术统计 mpv-core-stats.tsx =====
    "mpv.stats.source": "来源", // mpv-core-stats.tsx
    "mpv.stats.video": "视频", // mpv-core-stats.tsx
    "mpv.stats.color_format": "颜色 / 格式", // mpv-core-stats.tsx
    "mpv.stats.frame_drops": "掉帧（输出 / 解码器）", // mpv-core-stats.tsx
    "mpv.stats.presenter_drops": "呈现器掉帧（队列 / 浏览器）", // mpv-core-stats.tsx
    "mpv.stats.av_sync": "A/V 同步", // mpv-core-stats.tsx
    "mpv.stats.hw_decode": "硬件解码", // mpv-core-stats.tsx
    "mpv.stats.container": "封装格式", // mpv-core-stats.tsx
    // ===== 09-05 收官（player/mpv）=====
    "mpv.stats.custom_shaders_active": "自定义着色器（{count} 个启用）", // mpv-core-stats.tsx
} satisfies Dictionary
