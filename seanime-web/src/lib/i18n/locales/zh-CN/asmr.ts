import type { Dictionary } from "../../types"

/**
 * asmr 词条表（Phase 3.1：本地音声库 /asmr 页、库卡片、详情 modal 播放/完听/下载）
 * 搜索/探索页已复用的 asmr.* 词条见 search.ts。
 */
export const asmrDictionary = {
    "asmr.page.title": "本地音声库", // src/app/(main)/asmr/page.tsx

    "asmr.library.heard_progress": "已听 {listened}/{total}", // src/app/(main)/_features/asmr/_components/asmr-library-card.tsx
    "asmr.library.favorite": "收藏", // src/app/(main)/_features/asmr/_components/asmr-library-card.tsx
    "asmr.library.track_count": "{count} 音轨", // src/app/(main)/_features/asmr/_components/asmr-library-card.tsx
    "asmr.library.load_error": "加载本地音声库失败", // src/app/(main)/asmr/_components/asmr-library-view.tsx
    "asmr.library.retry": "重试", // src/app/(main)/asmr/_components/asmr-library-view.tsx
    "asmr.library.empty_title": "本地音声库为空", // src/app/(main)/asmr/_components/asmr-library-view.tsx
    "asmr.library.empty_desc": "把以 RJ 号命名的文件夹（如 RJ01234567）放进扫描目录后刷新，即可在此看到作品。", // src/app/(main)/asmr/_components/asmr-library-view.tsx
    "asmr.library.empty_dir": "默认扫描目录：", // src/app/(main)/asmr/_components/asmr-library-view.tsx
    "asmr.library.refresh": "刷新", // src/app/(main)/asmr/_components/asmr-library-view.tsx

    "asmr.detail.download_to_local": "下载到本地", // src/app/(main)/_features/asmr/_components/asmr-work-detail-modal.tsx
    "asmr.detail.download_to_local_short": "下载", // src/app/(main)/_features/asmr/_components/asmr-work-detail-modal.tsx

    "asmr.track.play": "播放", // src/app/(main)/_features/asmr/_components/asmr-work-detail-modal.tsx

    "asmr.playback.pause": "暂停", // src/app/(main)/_features/asmr/_components/asmr-playback-bar.tsx
    "asmr.playback.resume": "继续", // src/app/(main)/_features/asmr/_components/asmr-playback-bar.tsx
}
