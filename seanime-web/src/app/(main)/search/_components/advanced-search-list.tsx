import { AsmrWorkCard } from "@/app/(main)/_features/asmr/_components/asmr-work-card"
import { MediaCardLazyGrid } from "@/app/(main)/_features/media/_components/media-card-grid"
import { MediaEntryCard } from "@/app/(main)/_features/media/_components/media-entry-card"
import { useAnilistAdvancedSearch } from "@/app/(main)/search/_lib/handle-advanced-search"
import { cn } from "@/components/ui/core/styling"
import { LoadingSpinner } from "@/components/ui/loading-spinner"
import { t } from "@/lib/i18n"
import React from "react"
import { AiOutlinePlusCircle } from "react-icons/ai"

export function AdvancedSearchList() {

    const { isLoading, data, fetchNextPage, hasNextPage, type } = useAnilistAdvancedSearch()

    // Phase 2.5：novel 结果与 anime/manga 同构（AL_ListNovel = AL_ListAnime），asmr 结果走 works
    const items = data?.pages.filter(Boolean).flatMap((n: any) => (type === "asmr" ? n.works : n.Page?.media)).filter(Boolean)

    return <>
        {!isLoading && <MediaCardLazyGrid itemCount={items?.length ?? 0}>
            {type === "asmr"
                ? items?.map(work => (
                    <AsmrWorkCard
                        key={`${work.id}`}
                        work={work as any}
                        containerClassName="w-full"
                    />
                ))
                : items?.map(media => (
                    <MediaEntryCard
                        key={`${media.id}`}
                        media={media}
                        showLibraryBadge={true}
                        showTrailer
                        // novel 结果字段与 manga 卡片同构，复用 manga 卡片（点击跳 /manga/entry）
                        type={type === "novel" ? "manga" : (type as "anime" | "manga")}
                    />
                ))}
        </MediaCardLazyGrid>}
        {isLoading && <LoadingSpinner />}
        {((data?.pages.filter(Boolean).flatMap((n: any) => (type === "asmr" ? n.works : n.Page?.media)).filter(Boolean) || []).length > 0 && hasNextPage) &&
            <div
                data-advanced-search-list-load-more-container
                className={cn(
                    "relative flex flex-col rounded-[--radius-md] animate-none",
                    "cursor-pointer border border-none text-[--muted] hover:text-white pt-24 items-center gap-2 transition",
                )}
                onClick={() => fetchNextPage()}
            >
                <AiOutlinePlusCircle className="text-4xl" />
                <p className="text-lg font-medium">{t("search.list.load_more")}</p>
            </div>}
    </>
}
