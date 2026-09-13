// Phase 3.1c（Wave B / D8）：音声播放控制条纯函数 seam。
// 历史故障：自动完听若在前端不做保守判定，会在手动暂停/切歌时误标「已听」，污染进度。
// 故 deriveAutoCompletedTrack 必须仅在 running→stopped 跃迁且比例 >= 0.9 且 filepath 归一化匹配时触发。
// 本文件零 React 依赖，可被 vitest 直接单测（相对导入 + import type）。

/**
 * 归一化文件路径：反斜杠→正斜杠 + 小写（契约 D6/D8 匹配规则，跨平台大小写不敏感）。
 */
export function normalizeFilepath(p: string): string {
    return p.replace(/\\/g, "/").toLowerCase()
}

/**
 * 秒数格式化为 m:ss（契约 D7）。负数 / NaN / 非有限 → "0:00"。
 */
export function formatTime(sec: number): string {
    if (!Number.isFinite(sec) || sec < 0) return "0:00"
    const total = Math.floor(sec)
    const m = Math.floor(total / 60)
    const s = total % 60
    return `${m}:${s.toString().padStart(2, "0")}`
}

export type AsmrPlaybackTrack = {
    /** 相对 RJ 目录的 '/' 分隔路径（POST /asmr/track/progress 用） */
    path?: string
    /** 本地音频绝对路径（匹配 status.filepath 用） */
    localPath?: string
}

export type DeriveAutoCompletedTrackInput = {
    /** 上一轮轮询的 status（首次/无上轮为 null） */
    prev: {
        running: boolean
        currentTime: number
        duration: number
        filepath: string
    } | null
    /** 本轮轮询的 status（至少含 running / filepath） */
    next: {
        running: boolean
        filepath: string
    }
    /** 扁平化后的音轨树（含 path/localPath） */
    tracks: AsmrPlaybackTrack[]
    /** 已完听音轨 path 集合（去重，避免对已标记轨重发） */
    completed: Set<string>
}

/**
 * 自动完听判定（契约 D8）：轮询中发现 prev.running===true && next.running===false
 * 且 prev.filepath 归一化匹配某轨 localPath，且 prev.duration>0 且 prev.currentTime/prev.duration >= 0.9，
 * 且该轨 path 尚未在 completed 集合中 → 返回该轨 path；否则 null。
 *
 * 防误触发：仅捕获「正在播 → 停止」的跃迁（自然播完 mpv 退出 = running:false），
 * 手动暂停（running 仍 true，playing=false）不触发；切歌（prev.filepath 不匹配）不触发；
 * 比例不足 / 已 completed / prev 缺失 → 不动作。
 */
export function deriveAutoCompletedTrack(input: DeriveAutoCompletedTrackInput): string | null {
    const { prev, next, tracks, completed } = input
    if (!prev || !prev.running) return null
    if (next.running) return null
    if (prev.duration <= 0) return null
    if (prev.currentTime / prev.duration < 0.9) return null

    const normPrev = normalizeFilepath(prev.filepath)
    for (const track of tracks) {
        if (!track.localPath) continue
        if (normalizeFilepath(track.localPath) !== normPrev) continue
        const path = track.path ?? null
        if (!path) continue
        if (completed.has(path)) continue
        return path
    }
    return null
}
