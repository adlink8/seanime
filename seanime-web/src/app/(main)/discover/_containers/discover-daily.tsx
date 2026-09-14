import { AsmrWorkCard } from "@/app/(main)/_features/asmr/_components/asmr-work-card"
import { MediaEntryCard } from "@/app/(main)/_features/media/_components/media-entry-card"
import { MediaEntryCardSkeleton } from "@/app/(main)/_features/media/_components/media-entry-card-skeleton"
import { useDiscoverDaily, type DiscoverDailyDomain } from "@/api/hooks/discover.hooks"
import { AL_ListAnime, AL_ListManga, Asmr_SearchResult } from "@/api/generated/types"
import { Carousel, CarouselContent, CarouselDotButtons } from "@/components/ui/carousel"
import React from "react"

// Phase 3.9c：探索页「每日推荐」区块（当前 tab 域）。
// 数据流见后端 internal/handlers/discover.go：库抽样 → 关联聚合 → 滤已入库 → 日期种子洗牌 → 12 条。
// 空数据（库为空/无候选）时整块隐藏。

export function DiscoverDaily({ domain }: { domain: DiscoverDailyDomain }) {

    const ref = React.useRef<HTMLDivElement>(null)
    const { data, isLoading } = useDiscoverDaily(domain)

    // 按域取条目列表（响应复用各域现有 DTO）
    let items: { key: string, node: React.ReactNode }[] = []
    if (domain === "asmr") {
        const works = (data as Asmr_SearchResult | undefined)?.works?.filter(Boolean) ?? []
        items = works.map(work => ({
            key: work.id,
            node: (
                <AsmrWorkCard
                    work={work}
                    containerClassName="w-auto basis-[200px] md:basis-[250px] mx-2 mt-8 mb-0"
                />
            ),
        }))
    } else if (domain === "manga") {
        const media = (data as AL_ListManga | undefined)?.Page?.media?.filter(Boolean) ?? []
        items = media.map(m => ({
            key: String(m.id),
            node: (
                <MediaEntryCard
                    media={m}
                    containerClassName="basis-[200px] md:basis-[250px] mx-2 mt-8 mb-0"
                    type="manga"
                />
            ),
        }))
    } else if (domain === "novel") {
        // novel 响应 = AL_ListAnime（书籍语义条目，复用 manga 卡片跳 /manga/entry，同 novel 容器惯例）
        const media = (data as AL_ListAnime | undefined)?.Page?.media?.filter(Boolean) ?? []
        items = media.map(m => ({
            key: String(m.id),
            node: (
                <MediaEntryCard
                    media={m}
                    containerClassName="basis-[200px] md:basis-[250px] mx-2 mt-8 mb-0"
                    type="manga"
                    hideReleasingBadge
                />
            ),
        }))
    } else {
        const media = (data as AL_ListAnime | undefined)?.Page?.media?.filter(Boolean) ?? []
        items = media.map(m => ({
            key: String(m.id),
            node: (
                <MediaEntryCard
                    media={m}
                    containerClassName="basis-[200px] md:basis-[250px] mx-2 mt-8 mb-0"
                    type="anime"
                />
            ),
        }))
    }

    // 空数据隐藏（加载中先渲染骨架，避免区块跳动）
    if (!isLoading && items.length === 0) return null

    return (
        <Carousel
            className="w-full max-w-full"
            gap="xl"
            opts={{
                align: "start",
                dragFree: true,
            }}
            autoScroll
        >
            <CarouselDotButtons />
            <CarouselContent className="px-6" ref={ref}>
                {!isLoading ? items.map(item => (
                    <React.Fragment key={item.key}>{item.node}</React.Fragment>
                )) : [...Array(10).keys()].map((v, idx) => <MediaEntryCardSkeleton key={idx} />)}
            </CarouselContent>
        </Carousel>
    )
}
