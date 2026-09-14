import type { Dictionary } from "../../types"

/**
 * media 词条表
 * 归属：媒体卡片与预览（`app/(main)/_features/media/`）
 * 独占：05-04 —— 05-01 只建空骨架，不得添加词条
 *
 * 说明：`Watching` / `Reading` / `Planning` / `Paused` / `Completed` / `Dropped`
 * 六个列表状态下拉标签，语义与既有 `common.state.*` 完全一致，按契约 §2/§3
 * 「同名文案复用同一 key」**只读引用** common.state.*，不在本模块新建重复词条。
 */
export const mediaDictionary = {
    // —— 操作 ——
    "media.action.add_to_list": "加入列表", // anilist-media-entry-modal.tsx
    "media.action.confirm": "确认", // anilist-media-entry-modal.tsx
    "media.action.save": "保存", // anilist-media-entry-modal.tsx
    "media.action.preview": "预览", // media-entry-card.tsx
    "media.action.add_to_playlist": "加入播放列表", // media-entry-card.tsx
    "media.action.open_in_library_explorer": "在媒体库浏览器中打开", // media-entry-card.tsx
    "media.action.continue": "继续", // media-entry-card.tsx
    "media.action.watch": "观看", // media-entry-card.tsx
    "media.action.start_reading": "开始阅读", // media-entry-card.tsx
    "media.action.open_page": "打开页面", // media-preview-modal.tsx
    "media.action.trailer": "预告片", // media-preview-modal.tsx
    "media.action.remove_offline_data": "移除离线数据", // media-sync-track-button.tsx
    "media.action.save_locally": "保存到本地", // media-sync-track-button.tsx

    // —— 表单字段 ——
    "media.field.score": "评分", // anilist-media-entry-modal.tsx
    "media.field.progress": "进度", // anilist-media-entry-modal.tsx
    "media.field.start_date": "开始日期", // anilist-media-entry-modal.tsx
    "media.field.completion_date": "完成日期", // anilist-media-entry-modal.tsx
    "media.field.status": "状态", // anilist-media-entry-modal.tsx
    "media.field.total_rewatches": "重看次数", // anilist-media-entry-modal.tsx
    "media.field.total_rereads": "重读次数", // anilist-media-entry-modal.tsx

    // —— 列表状态（其余 6 项复用 common.state.*）——
    "media.status.repeating": "重温", // anilist-media-entry-modal.tsx
    "media.status.releasing": "放送中", // media-entry-card-components.tsx
    "media.status.not_yet_released": "尚未放送", // media-entry-card-components.tsx

    // —— 区块 / 卡片 ——
    "media.section.characters": "角色", // media-entry-characters-section.tsx
    "media.card.episode_prefix": "第", // media-entry-card-components.tsx
    "media.card.episode_suffix": "集", // media-entry-card-components.tsx

    // —— 角色 ——
    "media.characters.years_old": "{age} 岁", // media-entry-characters-section.tsx

    // —— 元数据 / 榜单 ——
    "media.metadata.show_audience_score": "显示观众评分", // media-entry-metadata-components.tsx
    "media.rankings.highest_rated": "最高评分", // media-entry-metadata-components.tsx
    "media.rankings.popular": "热门", // media-entry-metadata-components.tsx

    // —— Bangumi 详情增强（Phase 3.9d）——
    "media.bangumi.rank_badge": "Bangumi #{rank}", // media-entry-bangumi-info.tsx
    "media.bangumi.tags_title": "标签", // media-entry-bangumi-info.tsx
    "media.bangumi.infobox_title": "信息", // media-entry-bangumi-info.tsx

    // —— 文件路径选择 ——
    "media.file_path.select_all": "全选文件", // filepath-selector.tsx

    // —— 图片替代文本 ——
    "media.alt.episode_image": "剧集图片", // media-entry-characters-section.tsx
    "media.alt.banner_image": "横幅图片", // media-page-header-components.tsx, media-preview-modal.tsx, media-entry-page-small-banner.tsx
    "media.alt.cover_image": "封面图片", // media-page-header-components.tsx

    // —— 确认弹窗 ——
    "media.confirm.remove_offline_data_description": "此操作将移除该媒体条目的离线数据。确定要继续吗？", // media-sync-track-button.tsx

    // —— 剧集信息弹窗 ——
    "media.episode.default_title": "剧集", // media-episode-info-modal.tsx
    "media.episode.unknown_air_date": "未知放送日期", // media-episode-info-modal.tsx
    "media.episode.no_summary": "暂无简介", // media-episode-info-modal.tsx
    "media.episode.minutes": "分钟", // media-episode-info-modal.tsx

    // —— 重看 / 重读 ——
    "media.repeat.rewatch": "重看", // media-page-header-components.tsx
    "media.repeat.reread": "重读", // media-page-header-components.tsx
    "media.repeat.count": "{count} 次{type}", // media-page-header-components.tsx
} satisfies Dictionary
