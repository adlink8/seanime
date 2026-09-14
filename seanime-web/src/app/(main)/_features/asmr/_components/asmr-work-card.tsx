import { Asmr_Work } from "@/api/generated/types"
import { AsmrWorkDetailModal } from "@/app/(main)/_features/asmr/_components/asmr-work-detail-modal"
import { Badge } from "@/components/ui/badge"
import { cn } from "@/components/ui/core/styling"
import { SeaImage } from "@/components/shared/sea-image"
import { t } from "@/lib/i18n"
import React, { useState } from "react"

// Phase 2.5：ASMR（音声）结果卡片，供搜索页与探索页共享
// NSFW 封面复用项目的模糊遮罩思路：blur + hover 揭示（参照 MediaEntryCardAdultVeil 行为）

// Phase 3.9b P0：时长格式化（分钟 → "Xh Ym" / "Ym"；缺省或 ≤0 返回 null 不显示）
function formatDuration(min?: number): string | null {
    if (!min || min <= 0) return null
    if (min >= 60) return `${Math.floor(min / 60)}h ${min % 60}m`
    return `${min}m`
}

type AsmrWorkCardProps = {
    work: Asmr_Work
    containerClassName?: string
    onClick?: () => void
}

export function AsmrWorkCard({ work, containerClassName, onClick }: AsmrWorkCardProps) {
    const [modalOpen, setModalOpen] = useState(false)

    const handleOpen = React.useCallback(() => {
        if (onClick) {
            onClick()
            return
        }
        setModalOpen(true)
    }, [onClick])

    return (
        <>
            <div
                data-asmr-work-card-container
                data-asmr-work-id={work.id}
                className={cn("group relative w-[200px] md:w-[220px] flex-none cursor-pointer", containerClassName)}
                onClick={handleOpen}
            >
                <div
                    data-asmr-work-card-cover-container
                    className="relative w-full aspect-[3/4] rounded-[--radius-md] overflow-hidden bg-[--background] border border-[--border-color]"
                >
                    <SeaImage
                        src={work.coverUrl}
                        alt={work.title}
                        className={cn(
                            "w-full h-full object-cover object-center transition-all duration-300",
                            work.nsfw && "blur-lg group-hover:blur-none",
                        )}
                    />
                    {work.nsfw && (
                        <div
                            data-asmr-work-card-nsfw-veil
                            className="absolute inset-0 flex items-end justify-start p-1 group-hover:opacity-0 transition-opacity"
                        >
                            <Badge intent="alert-solid" size="sm">{t("search.asmr.r18")}</Badge>
                        </div>
                    )}
                    {work.hasSubtitle && (
                        <div data-asmr-work-card-subtitle-badge className="absolute top-1 right-1 z-[5]">
                            <Badge intent="primary-solid" size="sm">{t("search.asmr.subtitle_badge")}</Badge>
                        </div>
                    )}
                    {formatDuration(work.duration) && (
                        <div data-asmr-work-card-duration-badge className="absolute top-1 left-1 z-[5]">
                            <Badge intent="gray-solid" size="sm" className="!bg-gray-950 !bg-opacity-90">
                                {formatDuration(work.duration)}
                            </Badge>
                        </div>
                    )}
                    {work.rating > 0 && (
                        <div data-asmr-work-card-rating-badge className="absolute bottom-1 right-1 z-[5]">
                            {/* Phase 3.9b P0：rate_count<10 时灰显防误导 */}
                            <Badge
                                intent="gray-solid"
                                size="sm"
                                className={cn("!bg-gray-950 !bg-opacity-90", work.rateCount != null && work.rateCount < 10 && "opacity-50")}
                            >
                                ★ {work.rating.toFixed(1)}
                            </Badge>
                        </div>
                    )}
                </div>

                <div className="mt-2 space-y-1 w-full">
                    <p data-asmr-work-card-title className="text-sm font-medium line-clamp-2 leading-snug">{work.title}</p>
                    {work.circle && (
                        <p data-asmr-work-card-circle className="text-xs text-[--muted] line-clamp-1">{work.circle}</p>
                    )}
                    <p data-asmr-work-card-rj className="text-xs text-[--muted]">{work.rjId}</p>
                    {!!work.cvs.length && (
                        <p data-asmr-work-card-cvs className="text-xs text-[--muted] line-clamp-1">{work.cvs.join("、")}</p>
                    )}
                    {!!work.tags.length && (
                        <div className="flex flex-wrap gap-1 pt-0.5" data-asmr-work-card-tags-container>
                            {work.tags.slice(0, 3).map(tag => (
                                <Badge key={tag} intent="gray" size="sm">{tag}</Badge>
                            ))}
                        </div>
                    )}
                </div>
            </div>

            <AsmrWorkDetailModal work={work} open={modalOpen} onOpenChange={setModalOpen} />
        </>
    )
}
