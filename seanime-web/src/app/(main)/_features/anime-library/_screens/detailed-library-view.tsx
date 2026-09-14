import { Anime_Episode, Anime_LibraryCollectionEntry, Anime_LibraryCollectionList } from "@/api/generated/types"
import {
    __library_debouncedSearchInputAtom,
    __library_paramsAtom,
    __library_selectedListAtom,
    DETAILED_LIBRARY_DEFAULT_PARAMS,
    useHandleDetailedLibraryCollection,
} from "@/app/(main)/_features/anime-library/_lib/handle-detailed-library-collection.ts"
import { __home_currentView } from "@/app/(main)/_features/home/home-screen"
import { MediaCardLazyGrid } from "@/app/(main)/_features/media/_components/media-card-grid"
import { MediaEntryCard } from "@/app/(main)/_features/media/_components/media-entry-card"
import { MediaGenreSelector } from "@/app/(main)/_features/media/_components/media-genre-selector"
import { useNakamaStatus } from "@/app/(main)/_features/nakama/nakama-manager"
import { useServerStatus } from "@/app/(main)/_hooks/use-server-status"
import {
    ADVANCED_SEARCH_FORMATS,
    ADVANCED_SEARCH_SEASONS,
    GENRE_TRANSLATIONS,
    SEASON_TRANSLATIONS,
} from "@/app/(main)/search/_lib/advanced-search-constants"
import { PageWrapper } from "@/components/shared/page-wrapper"
import { AppLayoutStack } from "@/components/ui/app-layout"
import { IconButton } from "@/components/ui/button"
import { Carousel, CarouselContent, CarouselDotButtons } from "@/components/ui/carousel"
import { Combobox } from "@/components/ui/combobox"
import { cn } from "@/components/ui/core/styling"
import { LoadingSpinner } from "@/components/ui/loading-spinner"
import { Select } from "@/components/ui/select"
import { Switch } from "@/components/ui/switch"
import { StaticTabs } from "@/components/ui/tabs"
import { TextInput } from "@/components/ui/text-input"
import { useDebounce } from "@/hooks/use-debounce"
import { ANIME_COLLECTION_SORTING_OPTIONS } from "@/lib/helpers/filtering"
import { t } from "@/lib/i18n"
import { getLibraryCollectionTitle } from "@/lib/server/utils"
import { useThemeSettings } from "@/lib/theme/theme-hooks"
import { getYear } from "date-fns"
import { useSetAtom } from "jotai"
import { useAtom } from "jotai/react"
import React from "react"
import { AiOutlineArrowLeft } from "react-icons/ai"
import { BiTrash } from "react-icons/bi"
import { FaSortAmountDown } from "react-icons/fa"
import { FiSearch } from "react-icons/fi"
import { LuCalendar, LuLeaf } from "react-icons/lu"
import { MdPersonalVideo } from "react-icons/md"

type LibraryViewProps = {
    collectionList: Anime_LibraryCollectionList[]
    continueWatchingList: Anime_Episode[]
    isLoading: boolean
    hasEntries: boolean
    streamingMediaIds: number[]
    isNakamaLibrary: boolean
    isHomeItem?: boolean
    type?: "carousel" | "grid"
}

export function DetailedLibraryView(props: LibraryViewProps) {

    const {
        // collectionList: _collectionList,
        continueWatchingList,
        isLoading,
        hasEntries,
        streamingMediaIds,
        isNakamaLibrary,
        isHomeItem,
        type = "grid",
        ...rest
    } = props

    const ts = useThemeSettings()
    const setView = useSetAtom(__home_currentView)
    const nakamaStatus = useNakamaStatus()

    const {
        stats,
        libraryCollectionList,
        libraryGenres,
        libraryEntries,
    } = useHandleDetailedLibraryCollection()

    const [selectedList, setSelectedList] = useAtom(__library_selectedListAtom)

    React.useLayoutEffect(() => {
        if (selectedList !== "-" && selectedList !== "all") {
            setSelectedList("-")
        }
    }, [])

    if (isLoading) return <LoadingSpinner />

    if (!hasEntries) return null

    return (
        <PageWrapper className="px-4 space-y-8 relative z-[4]" data-detailed-library-view-container>

            {/* <div
             className={cn(
             "absolute top-[-20rem] left-0 w-full h-[30rem] bg-gradient-to-t from-[--background] to-transparent z-[-1]",
             TRANSPARENT_SIDEBAR_BANNER_IMG_STYLE,
             )}
             /> */}

            {!isHomeItem && <div className="flex flex-col md:flex-row gap-4 justify-between" data-detailed-library-view-header-container>
                <div className="flex gap-4 items-center relative w-fit">
                    <IconButton
                        icon={<AiOutlineArrowLeft />}
                        rounded
                        intent="white-outline"
                        size="sm"
                        onClick={() => setView("base")}
                    />
                    {!isNakamaLibrary && <h3 className="text-ellipsis truncate">{t("library.header.home")}</h3>}
                    {isNakamaLibrary &&
                        <h3 className="text-ellipsis truncate">{t("library.header.host_library", { name: nakamaStatus?.hostConnectionStatus?.username || "Host" })}</h3>}
                </div>

                <SearchInput />
            </div>}

            {(!isHomeItem) && <div
                className={cn(
                    "grid grid-cols-3 lg:grid-cols-6 gap-4 [&>div]:text-center [&>div>p]:text-[--muted]",
                    isNakamaLibrary && "lg:grid-cols-5",
                )}
                data-detailed-library-view-stats-container
            >
                {!isNakamaLibrary && <div>
                    <h3>{stats?.totalSize}</h3>
                    <p>{t("library.stats.size")}</p>
                </div>}
                <div>
                    <h3>{stats?.totalFiles}</h3>
                    <p>{t("library.stats.file_count")}</p>
                </div>
                <div>
                    <h3>{stats?.totalEntries}</h3>
                    <p>{t("library.stats.entry_count")}</p>
                </div>
                <div>
                    <h3>{stats?.totalShows}</h3>
                    <p>{t("library.stats.episodes")}</p>
                </div>
                <div>
                    <h3>{stats?.totalMovies}</h3>
                    <p>{t("library.stats.movies")}</p>
                </div>
                <div>
                    <h3>{stats?.totalSpecials}</h3>
                    <p>{t("library.stats.specials")}</p>
                </div>
            </div>}

            <SearchOptions />

            <GenreSelector genres={libraryGenres} />

            {selectedList !== "all" && libraryCollectionList.map(collection => {
                if (!collection.entries?.length) return null
                return <LibraryCollectionListItem key={collection.type} list={collection} streamingMediaIds={streamingMediaIds} type={type} />
            })}

            {selectedList === "all" && <MergedLibraryCollectionList entries={libraryEntries} streamingMediaIds={streamingMediaIds} type={type} />}
        </PageWrapper>
    )
}

const LibraryCollectionListItem = React.memo(({ list, streamingMediaIds, type }: {
    list: Anime_LibraryCollectionList,
    streamingMediaIds: number[],
    type: "grid" | "carousel"
}) => {

    const [selectedList, setSelectedList] = useAtom(__library_selectedListAtom)

    if (selectedList !== "-" && selectedList !== list.type) return null

    return (
        <React.Fragment key={list.type}>
            <h2>{getLibraryCollectionTitle(list.type)} <span className="text-[--muted] font-medium ml-3">{list?.entries?.length ?? 0}</span></h2>
            {type === "grid" && <MediaCardLazyGrid itemCount={list?.entries?.length || 0}>
                {list.entries?.map(entry => {
                    return <LibraryCollectionEntryItem key={entry.mediaId} entry={entry} streamingMediaIds={streamingMediaIds} type={type} />
                })}
            </MediaCardLazyGrid>}

            {type === "carousel" && <Carousel
                className={cn("w-full max-w-full !mt-0")}
                gap="xl"
                opts={{
                    align: "start",
                    dragFree: true,
                }}
                autoScroll={false}
            >
                <CarouselDotButtons className="-top-2" />
                <CarouselContent className="px-6">
                    {list.entries?.filter(Boolean)?.map(entry => {
                        return <LibraryCollectionEntryItem key={entry.mediaId} entry={entry} streamingMediaIds={streamingMediaIds} type={type} />
                    })}
                </CarouselContent>
            </Carousel>}
        </React.Fragment>
    )
})

const MergedLibraryCollectionList = React.memo(({ entries, streamingMediaIds, type }: {
    entries: Anime_LibraryCollectionEntry[],
    streamingMediaIds: number[]
    type: "grid" | "carousel"
}) => {

    return (
        <React.Fragment>
            {type === "grid" && <MediaCardLazyGrid itemCount={entries?.length || 0}>
                {entries?.map(entry => {
                    return <LibraryCollectionEntryItem key={entry.mediaId} entry={entry} streamingMediaIds={streamingMediaIds} type={type} />
                })}
            </MediaCardLazyGrid>}

            {type === "carousel" && <Carousel
                className={cn("w-full max-w-full !mt-0")}
                gap="xl"
                opts={{
                    align: "start",
                    dragFree: true,
                }}
                autoScroll={false}
            >
                <CarouselDotButtons className="-top-2" />
                <CarouselContent className="px-6">
                    {entries?.filter(Boolean)?.map(entry => {
                        return <LibraryCollectionEntryItem key={entry.mediaId} entry={entry} streamingMediaIds={streamingMediaIds} type={type} />
                    })}
                </CarouselContent>
            </Carousel>}
        </React.Fragment>
    )
})

const LibraryCollectionEntryItem = React.memo(({ entry, streamingMediaIds, type }: {
    entry: Anime_LibraryCollectionEntry,
    streamingMediaIds: number[]
    type: "grid" | "carousel"
}) => {
    return (
        <MediaEntryCard
            media={entry.media!}
            listData={entry.listData}
            libraryData={entry.libraryData}
            nakamaLibraryData={entry.nakamaLibraryData}
            showListDataButton
            withAudienceScore={false}
            type="anime"
            containerClassName={type === "carousel" ? "basis-[200px] md:basis-[250px] mx-2 mt-8 mb-0" : undefined}
            // showLibraryBadge={!!streamingMediaIds?.length && !streamingMediaIds.includes(entry.mediaId)}
        />
    )
})

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

const SearchInput = () => {

    const [inputValue, setInputValue] = React.useState("")
    const setDebouncedInput = useSetAtom(__library_debouncedSearchInputAtom)
    const debouncedInput = useDebounce(inputValue, 500)

    React.useEffect(() => {
        setDebouncedInput(inputValue)
    }, [debouncedInput])


    return (
        <div className="w-full md:w-[300px]">
            <TextInput
                leftIcon={<FiSearch />}
                value={inputValue}
                onValueChange={v => {
                    setInputValue(v)
                }}
                className="rounded-full bg-gray-900/50"
            />
        </div>
    )
}

export function SearchOptions() {

    const serverStatus = useServerStatus()
    const [params, setParams] = useAtom(__library_paramsAtom)
    const [selectedIndex, setSelectedIndex] = useAtom(__library_selectedListAtom)

    return (
        <AppLayoutStack className="px-4 xl:px-0" data-detailed-library-view-search-options-container>
            <div className="flex w-full justify-center">
                <StaticTabs
                    className="w-fit mb-6"
                    triggerClass="px-4 py-1"
                    items={[
                        { name: t("library.list.lists"), isCurrent: selectedIndex === "-", onClick: () => setSelectedIndex("-") },
                        { name: t("library.filter.all"), isCurrent: selectedIndex === "all", onClick: () => setSelectedIndex("all") },
                        { name: t("common.state.watching"), isCurrent: selectedIndex === "CURRENT", onClick: () => setSelectedIndex("CURRENT") },
                        { name: t("common.state.planning"), isCurrent: selectedIndex === "PLANNING", onClick: () => setSelectedIndex("PLANNING") },
                        { name: t("common.state.paused"), isCurrent: selectedIndex === "PAUSED", onClick: () => setSelectedIndex("PAUSED") },
                        { name: t("common.state.completed"), isCurrent: selectedIndex === "COMPLETED", onClick: () => setSelectedIndex("COMPLETED") },
                        { name: t("common.state.dropped"), isCurrent: selectedIndex === "DROPPED", onClick: () => setSelectedIndex("DROPPED") },
                    ]}
                />
            </div>
            <div
                className="grid grid-cols-2 md:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-[1fr_1fr_1fr_1fr_1fr_1fr_auto_auto] gap-4 items-start"
                data-detailed-library-view-search-options-grid
            >
                <Select
                    label={t("library.search.sorting")}
                    leftAddon={<FaSortAmountDown className={cn(params.sorting !== "TITLE" && "text-indigo-300 font-bold text-xl")} />}
                    className="w-full"
                    fieldClass="flex items-center"
                    inputContainerClass="w-full"
                    options={ANIME_COLLECTION_SORTING_OPTIONS}
                    value={params.sorting || "TITLE"}
                    onValueChange={v => setParams(draft => {
                        draft.sorting = v as any
                        return
                    })}
                    fieldLabelClass="hidden"
                    // disabled={!!params.title && params.title.length > 0}
                />
                <Select
                    leftAddon={
                        <MdPersonalVideo className={cn((params.format as any) !== null && (params.format as any) !== "" && "text-indigo-300 font-bold text-xl")} />}
                    label={t("library.filter.format")} placeholder={t("library.filter.format_all")}
                    className="w-full"
                    fieldClass="w-full"
                    options={ADVANCED_SEARCH_FORMATS}
                    value={params.format || ""}
                    onValueChange={v => setParams(draft => {
                        draft.format = v as any
                        return
                    })}
                    fieldLabelClass="hidden"
                />
                <Select
                    leftAddon={
                        <LuLeaf className={cn((params.season as any) !== null && (params.season as any) !== "" && "text-indigo-300 font-bold text-xl")} />}
                    label={t("library.filter.season")}
                    placeholder={t("library.filter.season_all")}
                    className="w-full"
                    fieldClass="w-full flex items-center"
                    inputContainerClass="w-full"
                    options={ADVANCED_SEARCH_SEASONS.map(season => ({ value: season.toUpperCase(), label: SEASON_TRANSLATIONS[season] || season }))}
                    value={params.season || ""}
                    onValueChange={v => setParams(draft => {
                        draft.season = v as any
                        return
                    })}
                    fieldLabelClass="hidden"
                />
                <Select
                    leftAddon={<LuCalendar className={cn((params.year !== null && params.year !== "") && "text-indigo-300 font-bold text-xl")} />}
                    label={t("library.filter.year")} placeholder={t("library.filter.year_any")}
                    className="w-full"
                    fieldClass="w-full"
                    options={[...Array(70)].map((v, idx) => getYear(new Date()) - idx).map(year => ({
                        value: String(year),
                        label: String(year),
                    }))}
                    value={params.year || ""}
                    onValueChange={v => setParams(draft => {
                        draft.year = v as any
                        return
                    })}
                    fieldLabelClass="hidden"
                />
                <div className="flex gap-4 items-center w-full">
                    <IconButton
                        icon={<BiTrash />} intent="alert-subtle" className="flex-none" onClick={() => {
                        setParams(DETAILED_LIBRARY_DEFAULT_PARAMS)
                    }}
                    />
                </div>
                {serverStatus?.settings?.anilist?.enableAdultContent && <div className="flex h-full items-center">
                    <Switch
                        label={t("library.filter.adult")}
                        value={params.isAdult}
                        onValueChange={v => setParams(draft => {
                            draft.isAdult = v
                            return
                        })}
                        fieldLabelClass="hidden"
                    />
                </div>}
            </div>

        </AppLayoutStack>
    )
}

function GenreSelector({ genres }: { genres: string[] }) {
    const [params, setParams] = useAtom(__library_paramsAtom)
    return (
        <MediaGenreSelector
            items={[
                {
                    name: t("library.filter.all"),
                    isCurrent: !params!.genre?.length,
                    onClick: () => setParams(draft => {
                        draft.genre = []
                        return
                    }),
                },
                ...genres.map(genre => ({
                    name: GENRE_TRANSLATIONS[genre] || genre,
                    isCurrent: params!.genre?.includes(genre) ?? false,
                    onClick: () => setParams(draft => {
                        if (draft.genre?.includes(genre)) {
                            draft.genre = draft.genre?.filter(g => g !== genre)
                        } else {
                            draft.genre = [...(draft.genre || []), genre]
                        }
                        return
                    }),
                })),
            ]}
        />
    )
}
