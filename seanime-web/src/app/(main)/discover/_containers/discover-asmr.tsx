import { AsmrWorkCard } from "@/app/(main)/_features/asmr/_components/asmr-work-card"
import { MediaEntryCardSkeleton } from "@/app/(main)/_features/media/_components/media-entry-card-skeleton"
import {
    useDiscoverLatestAsmr,
    useDiscoverPopularAsmr,
} from "@/app/(main)/discover/_lib/handle-discover-queries"
import { Carousel, CarouselContent, CarouselDotButtons } from "@/components/ui/carousel"
import React from "react"

// Phase 2.5：探索页音声（ASMR）区块，卡片点击打开详情 modal（AsmrWorkCard 内置）

export function DiscoverPopularAsmr() {

    const ref = React.useRef<HTMLDivElement>(null)
    const { data, isLoading } = useDiscoverPopularAsmr(ref)

    return (
        <Carousel
            className="w-full max-w-full"
            gap="md"
            opts={{
                align: "start",
                dragFree: true,
            }}
            autoScroll
        >
            <CarouselDotButtons />
            <CarouselContent className="px-6" ref={ref}>
                {!isLoading ? data?.works?.filter(Boolean).map(work => {
                    return (
                        <AsmrWorkCard
                            key={work.id}
                            work={work}
                            containerClassName="w-auto basis-[200px] md:basis-[250px] mx-2 mt-8 mb-0"
                        />
                    )
                }) : [...Array(10).keys()].map((v, idx) => <MediaEntryCardSkeleton key={idx} />)}
            </CarouselContent>
        </Carousel>
    )
}

export function DiscoverLatestAsmr() {

    const ref = React.useRef<HTMLDivElement>(null)
    const { data, isLoading } = useDiscoverLatestAsmr(ref)

    return (
        <Carousel
            className="w-full max-w-full"
            gap="md"
            opts={{
                align: "start",
                dragFree: true,
            }}
            autoScroll
        >
            <CarouselDotButtons />
            <CarouselContent className="px-6" ref={ref}>
                {!isLoading ? data?.works?.filter(Boolean).map(work => {
                    return (
                        <AsmrWorkCard
                            key={work.id}
                            work={work}
                            containerClassName="w-auto basis-[200px] md:basis-[250px] mx-2 mt-8 mb-0"
                        />
                    )
                }) : [...Array(10).keys()].map((v, idx) => <MediaEntryCardSkeleton key={idx} />)}
            </CarouselContent>
        </Carousel>
    )
}
