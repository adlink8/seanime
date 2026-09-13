import { MediaEntryCard } from "@/app/(main)/_features/media/_components/media-entry-card"
import { MediaEntryCardSkeleton } from "@/app/(main)/_features/media/_components/media-entry-card-skeleton"
import {
    useDiscoverPopularNovel,
    useDiscoverThisSeasonNovel,
} from "@/app/(main)/discover/_lib/handle-discover-queries"
import { Carousel, CarouselContent, CarouselDotButtons } from "@/components/ui/carousel"
import React from "react"

// Phase 2.5：探索页轻小说区块
// novel 结果响应形状与 anime/manga 同构（AL_ListNovel = AL_ListAnime），复用 manga 卡片（点击跳 /manga/entry）

export function DiscoverPopularNovel() {

    const ref = React.useRef<HTMLDivElement>(null)
    const { data, isLoading } = useDiscoverPopularNovel(ref)

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
                {!isLoading ? data?.Page?.media?.filter(Boolean).map(media => {
                    return (
                        <MediaEntryCard
                            key={media.id}
                            media={media}
                            containerClassName="basis-[200px] md:basis-[250px] mx-2 mt-8 mb-0"
                            type="manga"
                            hideReleasingBadge
                        />
                    )
                }) : [...Array(10).keys()].map((v, idx) => <MediaEntryCardSkeleton key={idx} />)}
            </CarouselContent>
        </Carousel>
    )
}

export function DiscoverThisSeasonNovel() {

    const ref = React.useRef<HTMLDivElement>(null)
    const { data, isLoading } = useDiscoverThisSeasonNovel(ref)

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
                {!isLoading ? data?.Page?.media?.filter(Boolean).map(media => {
                    return (
                        <MediaEntryCard
                            key={media.id}
                            media={media}
                            containerClassName="basis-[200px] md:basis-[250px] mx-2 mt-8 mb-0"
                            type="manga"
                            hideReleasingBadge
                        />
                    )
                }) : [...Array(10).keys()].map((v, idx) => <MediaEntryCardSkeleton key={idx} />)}
            </CarouselContent>
        </Carousel>
    )
}
