import { PluginWebviewSlot } from "@/app/(main)/_features/plugin/webview/plugin-webviews"
import { useServerStatus } from "@/app/(main)/_hooks/use-server-status"
import { DiscoverPageHeader } from "@/app/(main)/discover/_components/discover-page-header"
import { DiscoverAiringSchedule } from "@/app/(main)/discover/_containers/discover-airing-schedule"
import { DiscoverDaily } from "@/app/(main)/discover/_containers/discover-daily"
import { DiscoverLatestAsmr, DiscoverPopularAsmr } from "@/app/(main)/discover/_containers/discover-asmr"
import { DiscoverMissedSequelsSection } from "@/app/(main)/discover/_containers/discover-missed-sequels"
import { DiscoverPopularNovel, DiscoverThisSeasonNovel } from "@/app/(main)/discover/_containers/discover-novel"
import { DiscoverPastSeason, DiscoverThisSeason } from "@/app/(main)/discover/_containers/discover-popular"
import { DiscoverTrending } from "@/app/(main)/discover/_containers/discover-trending"
import { DiscoverTrendingCountry } from "@/app/(main)/discover/_containers/discover-trending-country"
import { DiscoverTrendingMovies } from "@/app/(main)/discover/_containers/discover-trending-movies"
import { DiscoverUpcoming } from "@/app/(main)/discover/_containers/discover-upcoming"
import { __discord_pageTypeAtom, __discover_isAdultAtom } from "@/app/(main)/discover/_lib/discover.atoms"
import { RecentReleases } from "@/app/(main)/schedule/_containers/recent-releases"
import { PageWrapper } from "@/components/shared/page-wrapper"
import { StaticTabs } from "@/components/ui/tabs"
import { Switch } from "@/components/ui/switch"
import { t } from "@/lib/i18n"
import { useRouter, useSearchParams } from "@/lib/navigation"
import { useAtom } from "jotai/react"
import { AnimatePresence, motion } from "motion/react"
import React from "react"


export default function Page() {

    const serverStatus = useServerStatus()
    const router = useRouter()
    const [pageType, setPageType] = useAtom(__discord_pageTypeAtom)
    const [isAdult, setIsAdult] = useAtom(__discover_isAdultAtom)
    const searchParams = useSearchParams()
    const searchType = searchParams.get("type")


    React.useEffect(() => {
        if (searchType) {
            setPageType(searchType as any)
        }
    }, [searchParams])

    return (
        <>
            <DiscoverPageHeader />
            <motion.div
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ duration: 0.5, delay: 0.6 }}
                className="p-4 sm:p-8 space-y-10 pb-10 relative z-[4]"
                data-discover-page-container
            >

                <div
                    className="lg:absolute w-full lg:-top-10 left-0 flex gap-4 p-4 items-center justify-center flex-wrap"
                    data-discover-page-header-tabs-container
                >
                    <div className="max-w-fit border rounded-full overflow-hidden" data-discover-page-header-tabs-inner-container>
                        <StaticTabs
                            className="h-full py-0"
                            triggerClass="px-4 py-2 h-full rounded-full border-transparent"
                            pillClass="rounded-full border-transparent"
                            items={[
                                { name: t("discover.tab.anime"), isCurrent: pageType === "anime", onClick: () => setPageType("anime") },
                                { name: t("navigation.item.schedule"), isCurrent: pageType === "schedule", onClick: () => setPageType("schedule") },
                                ...(serverStatus?.settings?.library?.enableManga ? [{
                                    name: t("discover.tab.manga"),
                                    isCurrent: pageType === "manga",
                                    onClick: () => setPageType("manga"),
                                }] : []),
                                ...(serverStatus?.settings?.library?.enableManga ? [{
                                    name: t("discover.tab.novel"),
                                    isCurrent: pageType === "novel",
                                    onClick: () => setPageType("novel"),
                                }] : []),
                                {
                                    name: t("discover.tab.asmr"),
                                    isCurrent: pageType === "asmr",
                                    onClick: () => setPageType("asmr"),
                                },
                            ]}
                        />
                    </div>
                    {/* 契约 03.9c：18+ 开关——仅在设置开启 EnableAdultContent 时可见（复用 library.filter.adult 文案） */}
                    {serverStatus?.settings?.anilist?.enableAdultContent && (
                        <div className="max-w-fit border rounded-full px-4 py-2 flex items-center" data-discover-page-adult-toggle>
                            <Switch
                                label={t("library.filter.adult")}
                                value={isAdult}
                                onValueChange={setIsAdult}
                                fieldLabelClass="hidden"
                            />
                        </div>
                    )}
                    {/*{!!customSources?.length && <div data-discover-page-header-custom-source-container>*/}
                    {/*    <SeaLink href="/custom-sources">*/}
                    {/*        <Button*/}
                    {/*            leftIcon={<MdDataSaverOn className="text-lg" />}*/}
                    {/*            intent="gray-outline"*/}
                    {/*            // size="lg"*/}
                    {/*            className="rounded-full"*/}
                    {/*            onClick={() => router.push("/search")}*/}
                    {/*        >*/}
                    {/*            Custom sources*/}
                    {/*        </Button>*/}
                    {/*    </SeaLink>*/}
                    {/*</div>}*/}
                    {/*<div data-discover-page-header-advanced-search-container>*/}
                    {/*    <Button*/}
                    {/*        leftIcon={<FaSearch />}*/}
                    {/*        intent="gray-outline"*/}
                    {/*        // size="lg"*/}
                    {/*        className="rounded-full"*/}
                    {/*        onClick={() => router.push("/search")}*/}
                    {/*    >*/}
                    {/*        Advanced search*/}
                    {/*    </Button>*/}
                    {/*</div>*/}
                </div>

                <PluginWebviewSlot slot="after-discover-screen-header" />

                <AnimatePresence mode="wait" initial={false}>
                    {pageType === "anime" && <PageWrapper
                        key="anime"
                        className="relative 2xl:order-first pb-10 pt-4 space-y-8"
                        {...{
                            initial: { opacity: 0, y: 60 },
                            animate: { opacity: 1, y: 0 },
                            exit: { opacity: 0, scale: 0.99 },
                            transition: {
                                duration: 0.35,
                            },
                        }}
                        data-discover-page-anime-container
                    >
                        {/* Phase 3.9c：每日推荐（当前 tab 域，空数据隐藏） */}
                        <div className="space-y-2 z-[5] relative" data-discover-page-anime-daily-container>
                            <h2>{t("discover.section.daily")}</h2>
                            <DiscoverDaily domain="anime" />
                        </div>
                        <div className="space-y-2 z-[5] relative" data-discover-page-anime-trending-container>
                            <h2>{t("common.home.trending")}</h2>
                            <DiscoverTrending />
                        </div>
                        <RecentReleases />
                        <div className="space-y-2 z-[5] relative" data-discover-page-anime-highest-rated-container>
                            <h2>{t("discover.section.this_season")}</h2>
                            <DiscoverThisSeason />
                        </div>
                        <div className="space-y-2 z-[5] relative" data-discover-page-anime-highest-rated-container>
                            <h2>{t("discover.section.past_season")}</h2>
                            <DiscoverPastSeason />
                        </div>
                        <DiscoverMissedSequelsSection />
                        <div className="space-y-2 z-[5] relative" data-discover-page-anime-upcoming-container>
                            <h2>{t("discover.section.upcoming")}</h2>
                            <DiscoverUpcoming />
                        </div>
                        <div className="space-y-2 z-[5] relative" data-discover-page-anime-trending-movies-container>
                            <h2>{t("discover.section.trending_movies")}</h2>
                            <DiscoverTrendingMovies />
                        </div>
                        {/*<div className="space-y-2 z-[5] relative">*/}
                        {/*    <h2>Popular shows</h2>*/}
                        {/*    <DiscoverPopular />*/}
                        {/*</div>*/}
                    </PageWrapper>}
                    {pageType === "schedule" && <PageWrapper
                        key="schedule"
                        className="relative 2xl:order-first pb-10 pt-4"
                        data-discover-page-schedule-container
                        {...{
                            initial: { opacity: 0, y: 60 },
                            animate: { opacity: 1, y: 0 },
                            exit: { opacity: 0, scale: 0.99 },
                            transition: {
                                duration: 0.35,
                            },
                        }}
                    >
                        <DiscoverAiringSchedule />
                    </PageWrapper>}
                    {pageType === "manga" && <PageWrapper
                        key="manga"
                        className="relative 2xl:order-first pb-10 pt-4"
                        data-discover-page-manga-container
                        {...{
                            initial: { opacity: 0, y: 60 },
                            animate: { opacity: 1, y: 0 },
                            exit: { opacity: 0, scale: 0.99 },
                            transition: {
                                duration: 0.35,
                            },
                        }}
                    >
                        {/*<div className="space-y-2 z-[5] relative">*/}
                        {/*    <h2>Trending right now</h2>*/}
                        {/*    <DiscoverTrendingMangaAll />*/}
                        {/*</div>*/}
                        {/* Phase 3.9c：每日推荐（当前 tab 域，空数据隐藏） */}
                        <div className="space-y-2 z-[5] relative" data-discover-page-manga-daily-container>
                            <h2>{t("discover.section.daily")}</h2>
                            <DiscoverDaily domain="manga" />
                        </div>
                        <div className="space-y-2 z-[5] relative" data-discover-page-manga-trending-container>
                            <h2>{t("discover.section.manga_trending_jp")}</h2>
                            <DiscoverTrendingCountry country="JP" forDiscoverHeader />
                        </div>
                        <div className="space-y-2 z-[5] relative" data-discover-page-manga-trending-manhwa-container>
                            <h2>{t("discover.section.manga_trending_kr")}</h2>
                            <DiscoverTrendingCountry country="KR" />
                        </div>
                        <div className="space-y-2 z-[5] relative" data-discover-page-manga-trending-manhua-container>
                            <h2>{t("discover.section.manga_trending_cn")}</h2>
                            <DiscoverTrendingCountry country="CN" />
                        </div>
                        {/*<div className="space-y-2 z-[5] relative">*/}
                        {/*    <DiscoverMangaSearchBar />*/}
                        {/*</div>*/}
                    </PageWrapper>}
                    {/* Phase 3.1：轻小说升为顶层分类 */}
                    {pageType === "novel" && <PageWrapper
                        key="novel"
                        className="relative 2xl:order-first pb-10 pt-4"
                        data-discover-page-novel-container
                        {...{
                            initial: { opacity: 0, y: 60 },
                            animate: { opacity: 1, y: 0 },
                            exit: { opacity: 0, scale: 0.99 },
                            transition: {
                                duration: 0.35,
                            },
                        }}
                    >
                        {/* Phase 3.9c：每日推荐（当前 tab 域，空数据隐藏） */}
                        <div className="space-y-2 z-[5] relative" data-discover-page-novel-daily-container>
                            <h2>{t("discover.section.daily")}</h2>
                            <DiscoverDaily domain="novel" />
                        </div>
                        <div className="space-y-2 z-[5] relative" data-discover-page-novel-popular-container>
                            <h2>{t("discover.section.popular_novel")}</h2>
                            <DiscoverPopularNovel />
                        </div>
                        <div className="space-y-2 z-[5] relative" data-discover-page-novel-this-season-container>
                            <h2>{t("discover.section.this_season_novel")}</h2>
                            <DiscoverThisSeasonNovel />
                        </div>
                    </PageWrapper>}
                    {/* Phase 3.1：音声升为顶层分类 */}
                    {pageType === "asmr" && <PageWrapper
                        key="asmr"
                        className="relative 2xl:order-first pb-10 pt-4"
                        data-discover-page-asmr-container
                        {...{
                            initial: { opacity: 0, y: 60 },
                            animate: { opacity: 1, y: 0 },
                            exit: { opacity: 0, scale: 0.99 },
                            transition: {
                                duration: 0.35,
                            },
                        }}
                    >
                        <div className="space-y-2 z-[5] relative" data-discover-page-asmr-popular-container>
                            <h2>{t("discover.section.popular_asmr")}</h2>
                            <DiscoverPopularAsmr />
                        </div>
                        <div className="space-y-2 z-[5] relative" data-discover-page-asmr-latest-container>
                            <h2>{t("discover.section.latest_asmr")}</h2>
                            <DiscoverLatestAsmr />
                        </div>
                    </PageWrapper>}
                </AnimatePresence>

            </motion.div>
        </>
    )
}
