import { useAsmrFavorite } from "@/api/hooks/asmr.hooks"
import { Asmr_LibraryEntry } from "@/api/generated/types"
import { usePlaybackPlayVideo } from "@/api/hooks/playback_manager.hooks"
import { AsmrWorkDetailModal } from "@/app/(main)/_features/asmr/_components/asmr-work-detail-modal"
import { Badge } from "@/components/ui/badge"
import { cn } from "@/components/ui/core/styling"
import { SeaImage } from "@/components/shared/sea-image"
import { t } from "@/lib/i18n"
import React, { useState } from "react"
import { LuMusic, LuStar } from "react-icons/lu"

// Phase 3.1：本地音声库卡片（GET /api/v1/asmr/library 条目）
// 复用封面模糊、circle/cvs/tags/rating，并展示 trackCount/totalSizeBytes/已听进度/收藏星标。

function formatBytes(bytes: number): string {
    if (!bytes || bytes < 0) return "0 B"
    const units = ["B", "KB", "MB", "GB", "TB"]
    const i = Math.min(units.length - 1, Math.floor(Math.log(bytes) / Math.log(1024)))
    const value = bytes / Math.pow(1024, i)
    return `${value.toFixed(i === 0 ? 0 : 1)} ${units[i]}`
}

type AsmrLibraryCardProps = {
    entry: Asmr_LibraryEntry
    containerClassName?: string
}

export function AsmrLibraryCard({ entry, containerClassName }: AsmrLibraryCardProps) {
    const [modalOpen, setModalOpen] = useState(false)
    const favorite = useAsmrFavorite()
    const [isFavorite, setIsFavorite] = useState(entry.isFavorite)

    const progressText = entry.trackCount > 0
        ? t("asmr.library.heard_progress", { listened: entry.listenedCount, total: entry.trackCount })
        : ""

    const toggleFavorite = React.useCallback((e: React.MouseEvent) => {
        e.stopPropagation()
        const next = !isFavorite
        setIsFavorite(next)
        favorite.mutate({ rjId: entry.rjId, favorite: next })
    }, [isFavorite, favorite, entry.rjId])

    // 本地库详情 modal 需要 Asmr_Work 形态，从库条目合成（id 用 rjId 占位，本地模式不依赖在线 id）
    const workForModal = React.useMemo(() => ({
        id: entry.rjId,
        rjId: entry.rjId,
        title: entry.title,
        circle: entry.circle,
        cvs: entry.cvs,
        tags: entry.tags,
        nsfw: entry.nsfw,
        coverUrl: entry.coverUrl,
        releaseDate: entry.releaseDate,
        rating: entry.rating,
        dlCount: 0,
        price: 0,
        hasSubtitle: false,
    }), [entry])

    return (
        <>
            <div
                data-asmr-library-card-container
                data-asmr-rj-id={entry.rjId}
                className={cn("group relative w-[200px] md:w-[220px] flex-none cursor-pointer", containerClassName)}
                onClick={() => setModalOpen(true)}
            >
                <div
                    data-asmr-library-card-cover-container
                    className="relative w-full aspect-[3/4] rounded-[--radius-md] overflow-hidden bg-[--background] border border-[--border-color]"
                >
                    <SeaImage
                        src={entry.coverUrl}
                        alt={entry.title}
                        className={cn(
                            "w-full h-full object-cover object-center transition-all duration-300",
                            entry.nsfw && "blur-lg group-hover:blur-none",
                        )}
                    />
                    {entry.nsfw && (
                        <div
                            data-asmr-library-card-nsfw-veil
                            className="absolute inset-0 flex items-end justify-start p-1 group-hover:opacity-0 transition-opacity"
                        >
                            <Badge intent="alert-solid" size="sm">{t("search.asmr.r18")}</Badge>
                        </div>
                    )}
                    {entry.rating > 0 && (
                        <div data-asmr-library-card-rating-badge className="absolute bottom-1 right-1 z-[5]">
                            <Badge intent="gray-solid" size="sm" className="!bg-gray-950 !bg-opacity-90">
                                ★ {entry.rating.toFixed(1)}
                            </Badge>
                        </div>
                    )}
                    {/* 收藏星标 */}
                    <button
                        type="button"
                        onClick={toggleFavorite}
                        title={t("asmr.library.favorite")}
                        data-asmr-library-card-favorite-button
                        className={cn(
                            "absolute top-1 right-1 z-[5] flex-none inline-flex items-center justify-center rounded-full p-1",
                            isFavorite ? "text-yellow-400" : "text-[--muted] hover:text-yellow-400",
                        )}
                    >
                        <LuStar className={cn("flex-none", isFavorite && "fill-current")} />
                    </button>
                </div>

                <div className="mt-2 space-y-1 w-full">
                    <p data-asmr-library-card-title className="text-sm font-medium line-clamp-2 leading-snug">{entry.title}</p>
                    {entry.circle && (
                        <p data-asmr-library-card-circle className="text-xs text-[--muted] line-clamp-1">{entry.circle}</p>
                    )}
                    <p data-asmr-library-card-rj className="text-xs text-[--muted]">{entry.rjId}</p>
                    <p data-asmr-library-card-meta className="text-xs text-[--muted] flex items-center gap-1">
                        <LuMusic className="flex-none" />
                        {t("asmr.library.track_count", { count: entry.trackCount })}
                        {" · "}
                        {formatBytes(entry.totalSizeBytes)}
                    </p>
                    {!!progressText && (
                        <p data-asmr-library-card-progress className="text-xs text-[--muted]">{progressText}</p>
                    )}
                    {!!entry.tags.length && (
                        <div className="flex flex-wrap gap-1 pt-0.5" data-asmr-library-card-tags-container>
                            {entry.tags.slice(0, 3).map(tag => (
                                <Badge key={tag} intent="gray" size="sm">{tag}</Badge>
                            ))}
                        </div>
                    )}
                </div>
            </div>

            <AsmrWorkDetailModal
                work={workForModal}
                open={modalOpen}
                onOpenChange={setModalOpen}
                localMode
            />
        </>
    )
}
