import { useAsmrLocalWork, useAsmrDownload, useAsmrTrackProgress, useAsmrWork } from "@/api/hooks/asmr.hooks"
import { usePlaybackPlayVideo } from "@/api/hooks/playback_manager.hooks"
import { Asmr_Track, Asmr_Work } from "@/api/generated/types"
import { deriveCompletedPaths } from "@/app/(main)/_features/asmr/_lib/asmr-completed-paths"
import { Badge } from "@/components/ui/badge"
import { cn } from "@/components/ui/core/styling"
import { LoadingSpinner } from "@/components/ui/loading-spinner"
import { Modal } from "@/components/ui/modal"
import { SeaImage } from "@/components/shared/sea-image"
import { t } from "@/lib/i18n"
import React, { useEffect, useState } from "react"
import { LuChevronDown, LuCheck, LuDownload, LuFolder, LuMusic, LuPlay } from "react-icons/lu"

// Phase 2.5：ASMR（音声）卡片与详情 modal，供搜索页与探索页共享（M3 前不做音频播放，音轨列表只读展示）
// Phase 3.1：本地库模式复用同一 modal（localMode），叶子音轨加播放/完听；顶部下载按钮（在线作品）。

type AsmrWorkDetailModalProps = {
    work: Asmr_Work
    open: boolean
    onOpenChange: (open: boolean) => void
    /** 本地库模式：拉取 /asmr/library/work/{rjId} 而非在线 /asmr/work/{id} */
    localMode?: boolean
}

export function AsmrWorkDetailModal({ work, open, onOpenChange, localMode }: AsmrWorkDetailModalProps) {
    // 仅在 modal 打开时拉取详情（含音轨树）；本地模式用 rjId 拉本地作品
    const { data: onlineDetail, isLoading: onlineLoading } = useAsmrWork(open && !localMode ? work.id : undefined, open)
    const { data: localDetail, isLoading: localLoading } = useAsmrLocalWork(open && localMode ? work.rjId : undefined, open)

    const detail = localMode ? localDetail : onlineDetail
    const isLoading = localMode ? localLoading : onlineLoading

    const playVideo = usePlaybackPlayVideo()
    const trackProgress = useAsmrTrackProgress()
    const download = useAsmrDownload()

    // 本地已完听音轨集合。保留本地 Set 以支持 handleToggleCompleted 的乐观更新；
    // 服务端真相由下方 effect 在 detail 变化时灌入（Phase 3.2 / D8）。
    const [completedPaths, setCompletedPaths] = useState<Set<string>>(() => new Set())

    // Phase 3.2（D8）：detail 变化时用服务端逐轨 completed 重置本地集合，实现跨会话回显。
    // ⚠ 依赖只放 detail 本身（react-query data 在结构不变时保持引用稳定）；
    // 若把派生结果/新建对象放进依赖，会在每次渲染重置 → 刚勾上的乐观更新被立刻抹掉（勾选闪回）。
    // 勾选路径：乐观 set → POST 成功 → invalidate → refetch 产生新 detail 引用 → 此处以服务端已落库的真相重置，值一致故不闪回。
    useEffect(() => {
        setCompletedPaths(deriveCompletedPaths(detail?.tracks))
    }, [detail])

    const handlePlayLocal = React.useCallback((path: string) => {
        // 仅走 playback-manager（mpv），禁止 directstream/mediastream/nativeplayer
        playVideo.mutate({ path })
    }, [playVideo])

    const handleToggleCompleted = React.useCallback((trackPath: string, completed: boolean) => {
        if (!work.rjId) return
        setCompletedPaths(prev => {
            const next = new Set(prev)
            if (completed) next.add(trackPath)
            else next.delete(trackPath)
            return next
        })
        trackProgress.mutate({ rjId: work.rjId, trackPath, completed })
    }, [work.rjId, trackProgress])

    const handleDownloadWork = React.useCallback(() => {
        if (!work.rjId) return
        download.mutate({ workId: work.id, rjId: work.rjId })
    }, [work.id, work.rjId, download])

    return (
        <Modal
            open={open}
            onOpenChange={onOpenChange}
            title={work.title}
            titleClass="text-lg line-clamp-2"
            contentClass="max-w-3xl overflow-hidden"
        >
            <div data-asmr-work-detail-container className="space-y-4 px-2 py-1">
                <div className="flex flex-col sm:flex-row gap-4">
                    {/* 封面：NSFW 加模糊遮罩，hover 揭示 */}
                    <div
                        data-asmr-work-detail-cover-container
                        className="group relative flex-none w-40 sm:w-48 aspect-[3/4] rounded-[--radius-md] overflow-hidden bg-[--background]"
                    >
                        <SeaImage
                            src={work.coverUrl}
                            alt={work.title}
                            className={cn(
                                "w-full h-full object-cover object-center",
                                work.nsfw && "blur-lg group-hover:blur-none transition-all duration-300",
                            )}
                        />
                        {work.nsfw && (
                            <div
                                data-asmr-work-detail-cover-nsfw-veil
                                className="absolute inset-0 flex items-end justify-center pb-2 group-hover:opacity-0 transition-opacity"
                            >
                                <Badge intent="alert-solid" size="sm">R18</Badge>
                            </div>
                        )}
                    </div>

                    <div className="flex-1 min-w-0 space-y-2">
                        <p className="text-sm text-[--muted]" data-asmr-work-detail-rj>{work.rjId}</p>
                        {work.circle && (
                            <p className="text-sm" data-asmr-work-detail-circle>
                                <span className="text-[--muted]">{t("search.asmr.circle")}：</span>{work.circle}
                            </p>
                        )}
                        {!!work.cvs.length && (
                            <p className="text-sm" data-asmr-work-detail-cvs>
                                <span className="text-[--muted]">{t("search.asmr.cv")}：</span>{work.cvs.join("、")}
                            </p>
                        )}
                        <div className="flex flex-wrap gap-x-4 gap-y-1 text-sm">
                            {work.rating > 0 && (
                                <p data-asmr-work-detail-rating>
                                    <span className="text-[--muted]">{t("search.asmr.rating")}：</span>{work.rating.toFixed(1)}
                                </p>
                            )}
                            {work.releaseDate && (
                                <p data-asmr-work-detail-release-date>
                                    <span className="text-[--muted]">{t("search.asmr.release_date")}：</span>{work.releaseDate}
                                </p>
                            )}
                            {!!work.dlCount && (
                                <p data-asmr-work-detail-dl-count>
                                    <span className="text-[--muted]">{t("search.asmr.dl_count")}：</span>{work.dlCount}
                                </p>
                            )}
                        </div>
                        {!!work.tags.length && (
                            <div className="flex flex-wrap gap-1 pt-1" data-asmr-work-detail-tags-container>
                                {work.tags.map(tag => (
                                    <Badge key={tag} intent="gray" size="sm">{tag}</Badge>
                                ))}
                            </div>
                        )}
                        {/* 下载到本地（在线作品场景；本地库已存在则隐藏） */}
                        {!localMode && (
                            <div className="pt-1" data-asmr-work-detail-download-container>
                                <button
                                    type="button"
                                    onClick={handleDownloadWork}
                                    disabled={download.isPending}
                                    data-asmr-work-detail-download-button
                                    className="inline-flex items-center gap-1.5 rounded-[--radius-md] border border-[--border-color] px-2.5 py-1 text-sm hover:text-[--brand] disabled:opacity-50"
                                >
                                    <LuDownload className="flex-none" />
                                    {t("asmr.detail.download_to_local")}
                                </button>
                            </div>
                        )}
                    </div>
                </div>

                {/* 音轨树（Phase 3.1：本地音轨加播放/完听） */}
                <div className="space-y-2" data-asmr-work-detail-tracks-container>
                    <h4 className="font-semibold text-sm">{t("search.asmr.tracks")}</h4>
                    {isLoading && <LoadingSpinner />}
                    {!isLoading && !!detail?.tracks?.length && (
                        <div className="max-h-72 overflow-y-auto rounded-[--radius-md] border border-[--border-color] p-2">
                            <AsmrTrackList
                                tracks={detail.tracks}
                                rjId={work.rjId}
                                completedPaths={completedPaths}
                                onToggleCompleted={handleToggleCompleted}
                                onPlayLocal={handlePlayLocal}
                                onDownloadWork={handleDownloadWork}
                            />
                        </div>
                    )}
                    {!isLoading && !detail?.tracks?.length && (
                        <p className="text-sm text-[--muted]">{t("search.asmr.tracks_empty")}</p>
                    )}
                </div>
            </div>
        </Modal>
    )
}

type AsmrTrackListProps = {
    tracks: Array<Asmr_Track>
    depth?: number
    rjId: string
    completedPaths: Set<string>
    onToggleCompleted: (trackPath: string, completed: boolean) => void
    onPlayLocal: (path: string) => void
    onDownloadWork: () => void
}

function AsmrTrackList({
    tracks,
    depth = 0,
    rjId,
    completedPaths,
    onToggleCompleted,
    onPlayLocal,
    onDownloadWork,
}: AsmrTrackListProps) {
    return (
        <>
            {tracks.map((track, idx) => track.type === "folder"
                ? <AsmrTrackFolder
                    key={`${depth}-${idx}-${track.title}`}
                    track={track}
                    depth={depth}
                    rjId={rjId}
                    completedPaths={completedPaths}
                    onToggleCompleted={onToggleCompleted}
                    onPlayLocal={onPlayLocal}
                    onDownloadWork={onDownloadWork}
                />
                : (
                    <div
                        key={`${depth}-${idx}-${track.title}`}
                        data-asmr-track-item
                        className="flex items-center gap-2 py-1 text-sm"
                        style={{ paddingLeft: `${depth * 16}px` }}
                    >
                        <LuMusic className="flex-none text-[--muted]" />
                        <span className="truncate">{track.title}</span>

                        {/* 本地音轨：播放按钮（localPath 非空时显示） */}
                        {!!track.localPath && (
                            <button
                                type="button"
                                onClick={() => onPlayLocal(track.localPath as string)}
                                title={t("asmr.track.play")}
                                data-asmr-track-play-button
                                className="ml-auto flex-none inline-flex items-center justify-center rounded-[--radius-md] px-1.5 py-0.5 text-[--brand] hover:bg-[--background]"
                            >
                                <LuPlay />
                            </button>
                        )}

                        {/* 本地音轨：完听勾选框（path 相对路径存在时显示） */}
                        {!!track.path && (
                            <label
                                className="ml-auto flex-none inline-flex items-center gap-1 cursor-pointer text-xs text-[--muted]"
                                data-asmr-track-completed-label
                            >
                                <input
                                    type="checkbox"
                                    checked={completedPaths.has(track.path)}
                                    onChange={e => onToggleCompleted(track.path as string, e.target.checked)}
                                    data-asmr-track-completed-checkbox
                                    className="accent-[--brand]"
                                />
                                <LuCheck className={cn("flex-none", completedPaths.has(track.path) && "text-[--brand]")} />
                            </label>
                        )}

                        {/* 在线音轨（无 localPath）：显示「下载到本地」 */}
                        {!track.localPath && (
                            <button
                                type="button"
                                onClick={onDownloadWork}
                                title={t("asmr.detail.download_to_local")}
                                data-asmr-track-download-button
                                className="ml-auto flex-none inline-flex items-center gap-1 rounded-[--radius-md] border border-[--border-color] px-1.5 py-0.5 text-xs hover:text-[--brand]"
                            >
                                <LuDownload className="flex-none" />
                                {t("asmr.detail.download_to_local_short")}
                            </button>
                        )}
                    </div>
                ))}
        </>
    )
}

type AsmrTrackFolderProps = {
    track: Asmr_Track
    depth: number
    rjId: string
    completedPaths: Set<string>
    onToggleCompleted: (trackPath: string, completed: boolean) => void
    onPlayLocal: (path: string) => void
    onDownloadWork: () => void
}

function AsmrTrackFolder({
    track,
    depth,
    rjId,
    completedPaths,
    onToggleCompleted,
    onPlayLocal,
    onDownloadWork,
}: AsmrTrackFolderProps) {
    const [expanded, setExpanded] = useState(depth === 0)

    return (
        <div data-asmr-track-folder>
            <button
                type="button"
                className="flex items-center gap-2 py-1 text-sm w-full text-left hover:text-[--brand]"
                style={{ paddingLeft: `${depth * 16}px` }}
                onClick={() => setExpanded(prev => !prev)}
            >
                <LuChevronDown className={cn("flex-none transition-transform", !expanded && "-rotate-90")} />
                <LuFolder className="flex-none text-[--muted]" />
                <span className="truncate">{track.title}</span>
            </button>
            {expanded && !!track.tracks?.length && (
                <div data-asmr-track-folder-content>
                    <AsmrTrackList
                        tracks={track.tracks}
                        depth={depth + 1}
                        rjId={rjId}
                        completedPaths={completedPaths}
                        onToggleCompleted={onToggleCompleted}
                        onPlayLocal={onPlayLocal}
                        onDownloadWork={onDownloadWork}
                    />
                </div>
            )}
        </div>
    )
}
