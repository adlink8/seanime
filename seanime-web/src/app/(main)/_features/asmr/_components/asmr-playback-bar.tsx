import { useAsmrPlaybackPause, useAsmrPlaybackSeek, useAsmrPlaybackStatus } from "@/api/hooks/asmr.hooks"
import { Asmr_Track } from "@/api/generated/types"
import { cn } from "@/components/ui/core/styling"
import { t } from "@/lib/i18n"
import { deriveAutoCompletedTrack, formatTime, normalizeFilepath } from "@/app/(main)/_features/asmr/_lib/asmr-playback"
import React, { useEffect, useMemo, useRef, useState } from "react"
import { LuPause, LuPlay } from "react-icons/lu"

// Phase 3.1c（Wave B / D5–D8）：音声播放控制条。
// 仅弹窗打开（组件挂载）时轮询 status；status.filepath 归一化匹配某轨 localPath 才渲染。
// 既有播放/完听逻辑（handlePlayLocal / handleToggleCompleted）零改动，本组件只做「轮询展示 + 自动完听」。

type AsmrPlaybackBarProps = {
    /** 扁平化前的音轨树（组件内会递归展平） */
    tracks: Array<Asmr_Track>
    /** 本地已完听音轨 path 集合（防自动完听重发） */
    completedPaths: Set<string>
    /** 命中自动完听时回调（签名对齐 handleToggleCompleted：path, completed=true） */
    onToggleCompleted: (trackPath: string, completed: boolean) => void
}

/** 递归展平音轨树（folder 嵌套），供 filepath 匹配与自动完听判定使用。 */
function flattenTracks(tracks: Array<Asmr_Track>): Array<Asmr_Track> {
    const out: Array<Asmr_Track> = []
    const walk = (list: Array<Asmr_Track>) => {
        for (const tr of list) {
            out.push(tr)
            if (tr.tracks?.length) walk(tr.tracks)
        }
    }
    walk(tracks)
    return out
}

export function AsmrPlaybackBar({ tracks, completedPaths, onToggleCompleted }: AsmrPlaybackBarProps) {
    // 组件挂载即轮询（弹窗关闭 = 卸载，自然停止轮询，契约 D5/D6）
    const status = useAsmrPlaybackStatus(true)
    const pauseMut = useAsmrPlaybackPause()
    const seekMut = useAsmrPlaybackSeek()

    const flatTracks = useMemo(() => flattenTracks(tracks), [tracks])

    // 保存上一轮 status 快照，用于检测 running→stopped 跃迁
    const prevStatusRef = useRef<{
        running: boolean
        currentTime: number
        duration: number
        filepath: string
    } | null>(null)

    // 自动完听：status 变化且发生 running→stopped 跃迁时判定（契约 D8）
    useEffect(() => {
        const data = status.data
        const prev = prevStatusRef.current
        if (data) {
            const next = { running: data.running, filepath: data.filepath }
            const trackPath = deriveAutoCompletedTrack({
                prev,
                next,
                tracks: flatTracks,
                completed: completedPaths,
            })
            if (trackPath) onToggleCompleted(trackPath, true)
        }
        if (data) {
            prevStatusRef.current = {
                running: data.running,
                currentTime: data.currentTime,
                duration: data.duration,
                filepath: data.filepath,
            }
        }
    }, [status.data, flatTracks, completedPaths, onToggleCompleted])

    // status.filepath 归一化匹配某轨 localPath 才渲染（在线作品 localPath 为空自然不匹配，契约 D6）
    const matchedTrack = useMemo(() => {
        const data = status.data
        if (!data) return undefined
        const norm = normalizeFilepath(data.filepath)
        return flatTracks.find(tr => !!tr.localPath && normalizeFilepath(tr.localPath) === norm)
    }, [status.data, flatTracks])

    // 进度条：拖动用本地 state 跟随，松手才提交 seek（避免 IPC 刷屏，契约 D7）。
    // ⚠ 这些 hooks 必须在下面的提前 return 之前调用——违反 rules-of-hooks会在
    // status 从无到有时抛 "Rendered more hooks than during the previous render"。
    const [dragging, setDragging] = useState(false)
    const [localTime, setLocalTime] = useState(0)
    useEffect(() => {
        if (!dragging && status.data) setLocalTime(status.data.currentTime)
    }, [status.data?.currentTime, dragging])

    // 无匹配（无播放 / 非 mpv / 不匹配）则不渲染（契约 AC-06）
    if (!status.data || !matchedTrack) return null

    const { playing, currentTime, duration, filename } = status.data

    const displayTime = dragging ? localTime : currentTime
    const pct = duration > 0 ? (displayTime / duration) * 100 : 0

    const commitSeek = () => {
        seekMut.mutate({ position: localTime })
    }

    return (
        <div
            data-asmr-playback-bar
            className="flex flex-col gap-2 rounded-[--radius-md] border border-[--border-color] bg-[--background] p-3"
        >
            <div className="flex items-center gap-3">
                <button
                    type="button"
                    onClick={() => pauseMut.mutate({ paused: !playing })}
                    disabled={pauseMut.isPending}
                    data-asmr-playback-toggle
                    title={playing ? t("asmr.playback.pause") : t("asmr.playback.resume")}
                    className="flex-none inline-flex items-center justify-center rounded-[--radius-md] border border-[--border-color] px-2 py-1 text-[--brand] hover:bg-[--background] disabled:opacity-50"
                >
                    {playing ? <LuPause /> : <LuPlay />}
                </button>

                <span className="min-w-0 flex-1 truncate text-sm" data-asmr-playback-filename>
                    {filename}
                </span>

                <span className="flex-none text-xs tabular-nums text-[--muted]" data-asmr-playback-time>
                    {formatTime(displayTime)} / {formatTime(duration)}
                </span>
            </div>

            <input
                type="range"
                min={0}
                max={Math.max(duration, 1)}
                step={0.1}
                value={displayTime}
                onChange={e => {
                    setDragging(true)
                    setLocalTime(Number(e.target.value))
                }}
                onPointerUp={() => {
                    setDragging(false)
                    commitSeek()
                }}
                onMouseUp={() => {
                    setDragging(false)
                    commitSeek()
                }}
                data-asmr-playback-seek
                aria-label="Seek position"
                className={cn(
                    "h-1.5 w-full cursor-pointer appearance-none rounded-full bg-[--border-color] accent-[--brand]",
                    "[&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:h-3 [&::-webkit-slider-thumb]:w-3 [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:bg-[--brand]",
                )}
                style={{ background: `linear-gradient(to right, var(--brand) ${pct}%, var(--border-color) ${pct}%)` }}
            />
        </div>
    )
}
