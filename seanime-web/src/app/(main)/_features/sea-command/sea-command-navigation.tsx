import { useGetAnimeCollection } from "@/api/hooks/anilist.hooks"
import { useGetMangaCollection } from "@/api/hooks/manga.hooks"
import { useLibraryCollection } from "@/app/(main)/_hooks/anime-library-collection-loader.ts"
import { useServerStatus } from "@/app/(main)/_hooks/use-server-status"
import { CommandGroup, CommandItem, CommandShortcut } from "@/components/ui/command"
import { t } from "@/lib/i18n"
import { useRouter } from "@/lib/navigation"
import React from "react"
import { BiArrowBack } from "react-icons/bi"
import { CommandHelperText, CommandItemMedia } from "./_components/command-utils"
import { useSeaCommandContext } from "./sea-command"
import { seaCommand_compareMediaTitles } from "./utils"

// only rendered when typing "/anime", "/library" or "/manga"
export function SeaCommandUserMediaNavigation() {

    const { input, select, command: { isCommand, command, args }, scrollToTop } = useSeaCommandContext()
    const { data: animeCollection, isLoading: isAnimeLoading } = useGetAnimeCollection() // should be available instantly
    const { data: mangaCollection, isLoading: isMangaLoading } = useGetMangaCollection()
    const animeLibraryCollection = useLibraryCollection()

    const anime = animeCollection?.MediaListCollection?.lists?.flatMap(n => n?.entries)?.filter(Boolean)?.map(n => n.media)?.filter(Boolean) ?? []
    const manga = mangaCollection?.lists?.flatMap(n => n?.entries)?.filter(Boolean)?.map(n => n.media)?.filter(Boolean) ?? []

    const router = useRouter()

    const query = args.join(" ")
    const filteredAnime = (command === "anime" && query.length > 0) ? anime.filter(n => seaCommand_compareMediaTitles(n.title, query)) : []
    const filteredManga = (command === "manga" && query.length > 0) ? manga.filter(n => seaCommand_compareMediaTitles(n.title, query)) : []
    const filteredAnimeLibrary = (command === "library" && query.length > 0) ? animeLibraryCollection?.lists?.flatMap(l => l.entries)
        ?.filter(n => seaCommand_compareMediaTitles(n?.media?.title, query))
        ?.map(n => n?.media)
        ?.filter(Boolean) ?? [] : []

    return (
        <>
            {query.length === 0 && (
                <>
                    <CommandHelperText
                        command="/anime [title]"
                        description={t("misc.sea_command.find_anime_collection")}
                        show={command === "anime"}
                    />
                    <CommandHelperText
                        command="/manga [title]"
                        description={t("misc.sea_command.find_manga_collection")}
                        show={command === "manga"}
                    />
                    <CommandHelperText
                        command="/library [title]"
                        description={t("misc.sea_command.find_anime_library")}
                        show={command === "library"}
                    />
                </>
            )}

            {command === "anime" && filteredAnime.length > 0 && (
                <CommandGroup heading={t("misc.sea_command.my_anime")}>
                    {filteredAnime.map(n => (
                        <CommandItem
                            key={n.id}
                            onSelect={() => {
                                select(() => {
                                    router.push(`/entry?id=${n.id}`)
                                })
                            }}
                        >
                            <CommandItemMedia media={n} type="anime" />
                        </CommandItem>
                    ))}
                </CommandGroup>
            )}

            {command === "library" && filteredAnimeLibrary.length > 0 && (
                <CommandGroup heading="Library anime">
                    {filteredAnimeLibrary.map(n => (
                        <CommandItem
                            key={n.id}
                            onSelect={() => {
                                select(() => {
                                    router.push(`/entry?id=${n.id}`)
                                })
                            }}
                        >
                            <CommandItemMedia media={n} type="anime" />
                        </CommandItem>
                    ))}
                </CommandGroup>
            )}
            {command === "manga" && filteredManga.length > 0 && (
                <CommandGroup heading="My manga">
                    {filteredManga.map(n => (
                        <CommandItem
                            key={n.id}
                            onSelect={() => {
                                select(() => {
                                    router.push(`/manga/entry?id=${n.id}`)
                                })
                            }}
                        >
                            <CommandItemMedia media={n} type="manga" />
                        </CommandItem>
                    ))}
                </CommandGroup>
            )}
        </>
    )
}

export function SeaCommandNavigation() {

    const serverStatus = useServerStatus()

    const { input, select, command: { isCommand, command, args } } = useSeaCommandContext()

    const router = useRouter()

    const pages = [
        {
            name: t("navigation.item.home"),
            href: "/",
            flag: "home",
            show: !serverStatus?.isOffline,
        },
        {
            name: t("navigation.item.schedule"),
            href: "/schedule",
            flag: "schedule",
            show: !serverStatus?.isOffline,
        },
        {
            name: t("navigation.sidebar.settings"),
            href: "/settings",
            flag: "settings",
            show: !serverStatus?.isOffline,
        },
        {
            name: t("navigation.item.manga"),
            href: "/manga",
            flag: "manga",
            show: !serverStatus?.isOffline,
        },
        {
            name: t("navigation.item.discover"),
            href: "/discover",
            flag: "discover",
            show: !serverStatus?.isOffline,
        },
        {
            name: t("navigation.item.lists"),
            href: "/lists",
            flag: "lists",
            show: !serverStatus?.isOffline,
        },
        {
            name: t("navigation.sidebar.auto_downloader"),
            href: "/auto-downloader",
            flag: "auto-downloader",
            show: !serverStatus?.isOffline,
        },
        {
            name: t("navigation.sidebar.torrent_list"),
            href: "/torrent-list",
            flag: "torrent-list",
            show: !serverStatus?.isOffline,
        },
        {
            name: t("navigation.sidebar.scan_summaries"),
            href: "/scan-summaries",
            flag: "scan-summaries",
            show: !serverStatus?.isOffline,
        },
        {
            name: t("navigation.sidebar.extensions"),
            href: "/extensions",
            flag: "extensions",
            show: !serverStatus?.isOffline,
        },
        {
            name: t("misc.sea_command.advanced_search"),
            href: "/search",
            flag: "search",
            show: !serverStatus?.isOffline,
        },
    ]

    // If no args, show all pages
    // If args, show pages that match the args
    const filteredPages = pages.filter(page => page.flag.startsWith(command))


    // if (!input.startsWith("/")) return null


    return (
        <>
            {command.startsWith("ba") && (
                <CommandGroup heading={t("misc.sea_command.navigation")}>
                    <CommandItem
                        onSelect={() => {
                            select(() => {
                                router.back()
                            })
                        }}
                    >
                        <BiArrowBack className="mr-2 h-4 w-4" />
                        <span>{t("misc.sea_command.go_back")}</span>
                    </CommandItem>
                </CommandGroup>
            )}
            {command.startsWith("fo") && (
                <CommandGroup heading="Navigation">
                    <CommandItem
                        onSelect={() => {
                            select(() => {
                                router.forward()
                            })
                        }}
                    >
                        <BiArrowBack className="mr-2 h-4 w-4 rotate-180" />
                        <span>{t("misc.sea_command.go_forward")}</span>
                    </CommandItem>
                </CommandGroup>
            )}

            {/*Typing `/library`, `/schedule`, etc. without args*/}
            {isCommand && filteredPages.length > 0 && args.length === 0 && (
                <CommandGroup heading={t("misc.sea_command.screens")}>
                    <>
                        {filteredPages.filter(page => page.show).map(page => (
                            <CommandItem
                                key={page.flag}
                                onSelect={() => {
                                    select(() => {
                                        router.push(page.href)
                                    })
                                }}
                            >
                                <span className="text-sm tracking-wide font-bold text-[--muted]">{t("misc.sea_command.go_to")}&nbsp;</span>{" "}{page.name}
                                {command === page.flag ? <CommandShortcut>{t("misc.sea_command.enter")}</CommandShortcut> : <CommandShortcut>/{page.flag}</CommandShortcut>}
                            </CommandItem>
                        ))}
                    </>
                </CommandGroup>
            )}
            {(command !== "back" && command !== "forward") && (
                <CommandGroup heading="Navigation">
                    {/* {command === "" && ( */}
                    <>
                        <CommandItem
                            onSelect={() => {
                                select(() => {
                                    router.back()
                                })
                            }}
                        >
                            <BiArrowBack className="mr-2 h-4 w-4" />
                            <span>{t("misc.sea_command.go_back")}</span>
                        </CommandItem>
                        <CommandItem
                            onSelect={() => {
                                select(() => {
                                    router.forward()
                                })
                            }}
                        >
                            <BiArrowBack className="mr-2 h-4 w-4 rotate-180" />
                            <span>{t("misc.sea_command.go_forward")}</span>
                        </CommandItem>
                    </>
                    {/* )} */}
                </CommandGroup>
            )}
        </>
    )
}
