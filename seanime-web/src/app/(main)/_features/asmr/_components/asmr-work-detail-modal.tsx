import { useAsmrWork } from "@/api/hooks/asmr.hooks"
import { Asmr_Track, Asmr_Work } from "@/api/generated/types"
import { Badge } from "@/components/ui/badge"
import { cn } from "@/components/ui/core/styling"
import { LoadingSpinner } from "@/components/ui/loading-spinner"
import { Modal } from "@/components/ui/modal"
import { SeaImage } from "@/components/shared/sea-image"
import { t } from "@/lib/i18n"
import React, { useState } from "react"
import { LuChevronDown, LuFolder, LuMusic } from "react-icons/lu"

// Phase 2.5：ASMR（音声）卡片与详情 modal，供搜索页与探索页共享（M3 前不做音频播放，音轨列表只读展示）

type AsmrWorkDetailModalProps = {
    work: Asmr_Work
    open: boolean
    onOpenChange: (open: boolean) => void
}

export function AsmrWorkDetailModal({ work, open, onOpenChange }: AsmrWorkDetailModalProps) {
    // 仅在 modal 打开时拉取详情（含音轨树）
    const { data: detail, isLoading } = useAsmrWork(open ? work.id : undefined, open)

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
                    </div>
                </div>

                {/* 音轨树（只读展示，M3 前不做播放） */}
                <div className="space-y-2" data-asmr-work-detail-tracks-container>
                    <h4 className="font-semibold text-sm">{t("search.asmr.tracks")}</h4>
                    {isLoading && <LoadingSpinner />}
                    {!isLoading && !!detail?.tracks?.length && (
                        <div className="max-h-72 overflow-y-auto rounded-[--radius-md] border border-[--border-color] p-2">
                            <AsmrTrackList tracks={detail.tracks} />
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

function AsmrTrackList({ tracks, depth = 0 }: { tracks: Array<Asmr_Track>, depth?: number }) {
    return (
        <>
            {tracks.map((track, idx) => track.type === "folder"
                ? <AsmrTrackFolder key={`${depth}-${idx}-${track.title}`} track={track} depth={depth} />
                : (
                    <div
                        key={`${depth}-${idx}-${track.title}`}
                        data-asmr-track-item
                        className="flex items-center gap-2 py-1 text-sm"
                        style={{ paddingLeft: `${depth * 16}px` }}
                    >
                        <LuMusic className="flex-none text-[--muted]" />
                        <span className="truncate">{track.title}</span>
                    </div>
                ))}
        </>
    )
}

function AsmrTrackFolder({ track, depth }: { track: Asmr_Track, depth: number }) {
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
                    <AsmrTrackList tracks={track.tracks} depth={depth + 1} />
                </div>
            )}
        </div>
    )
}
