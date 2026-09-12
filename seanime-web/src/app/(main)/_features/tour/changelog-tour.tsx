import { useSeaCommand } from "@/app/(main)/_features/sea-command/sea-command.tsx"
import { SeaImage } from "@/components/shared/sea-image"
import { useRouter } from "@/lib/navigation"
import { t } from "@/lib/i18n"
import { useAtom } from "jotai"
import { atomWithStorage } from "jotai/utils"
import React from "react"
import { useWindowSize } from "react-use"
import { useServerStatus } from "../../_hooks/use-server-status"
import { __settings_tabAtom } from "../../settings/_components/settings-page.atoms"
import { __scanner_modalIsOpen } from "../anime-library/_containers/scanner-modal"
import { tourHelpers, useTour } from "./tour"
import { TourStep } from "./tour"

export const seenChangelogAtom = atomWithStorage<string | null>("sea-seen-changelog", null, undefined, { getOnInit: true })

function useSetupTour(): Record<string, () => TourStep[]> {
    const serverStatus = useServerStatus()
    const router = useRouter()
    const [, openScannerModal] = useAtom(__scanner_modalIsOpen)
    const [, setSettingsTab] = useAtom(__settings_tabAtom)
    const { setSeaCommandOpen, setSeaCommandInput } = useSeaCommand()

    const get3_5_0 = (): TourStep[] => {
        return [
            {
                id: "changelog-1",
                content: (
                    <div>
                        <h4 className="text-xl font-bold text-white">{t("misc.changelog.whats_new_350")}</h4>
                        <p>{t("misc.changelog.intro_features")}</p>
                    </div>
                ),
                route: "/",
                nextLabel: t("manga.action.start"),
                ignoreOutsideClick: true,
            },
            {
                id: "scanner",
                target: "[data-home-toolbar-scan-button]",
                title: t("misc.changelog.new_scanner"),
                content: t("misc.changelog.new_scanner_content"),
                route: "/",
                advanceOnTargetClick: true,
                ignoreOutsideClick: true,
                condition: () => !!serverStatus?.settings?.library?.libraryPath?.length,
                conditionFailBehavior: "modal",
            },
            {
                id: "scanner-2",
                target: "[data-scanner-modal-content]",
                title: t("misc.changelog.new_scanner"),
                content: t("misc.changelog.new_scanner_content2"),
                route: "/",
                prepare: () => {
                    openScannerModal(true)
                },
                advanceOnTargetClick: true,
                ignoreOutsideClick: true,
                condition: () => !!serverStatus?.settings?.library?.libraryPath?.length,
                conditionFailBehavior: "skip",
            },
            {
                id: "scanner-3",
                target: "[data-settings-anime-library='advanced-accordion-trigger']",
                title: t("misc.changelog.scanner_config"),
                content: t("misc.changelog.scanner_config_content"),
                route: "/settings",
                prepare: async () => {
                    setSettingsTab("library")
                    await tourHelpers.waitForSelector("[data-settings-anime-library='advanced-accordion-trigger']")
                    await tourHelpers.click("[data-settings-anime-library='advanced-accordion-trigger']", 200)
                },
                advanceOnTargetClick: false,
                ignoreOutsideClick: true,
            },
            {
                id: "issue-recorder",
                target: "[data-open-issue-recorder-button]",
                title: t("misc.changelog.issue_recorder"),
                // content: "The issue recorder has been improved and will now record the UI.",
                content: <div>
                    <SeaImage
                        src="https://github.com/5rahim/hibike/blob/main/changelog/3_5-issue-recorder.gif?raw=true"
                        alt={t("misc.changelog.issue_recorder")}
                        width="100%"
                        height="auto"
                        className="rounded-md"
                        allowGif
                    />
                    <p className="mt-2">{t("misc.changelog.issue_recorder_content")}</p>
                </div>,
                route: "/settings",
                prepare: async () => {
                    setSettingsTab("seanime")
                },
                advanceOnTargetClick: false,
                ignoreOutsideClick: true,
                popoverWidth: 500,
            },
            {
                id: "transcode-new-player",
                target: "[data-tab-trigger='mediastream']",
                title: t("misc.changelog.transcode_player"),
                content: t("misc.changelog.transcode_player_content"),
                route: "/settings",
                prepare: async () => {
                    setSettingsTab("mediastream")
                },
                advanceOnTargetClick: false,
                ignoreOutsideClick: true,
            },
            {
                id: "search",
                target: "[data-vertical-menu-item='Search']",
                title: t("navigation.item.search"),
                content: t("misc.changelog.search_content"),
                route: "/search",
                advanceOnTargetClick: false,
                ignoreOutsideClick: true,
            },
            {
                id: "entry",
                title: t("misc.changelog.new_player_features"),
                content: <div>
                    <SeaImage
                        src="https://github.com/5rahim/hibike/blob/main/changelog/3_5-videocore-characters.png?raw=true"
                        alt={t("misc.changelog.char_lookup_alt")}
                        width="100%"
                        height="auto"
                        className="rounded-md"
                    />
                    <p className="mt-2">{t("misc.changelog.char_lookup_content")}</p>
                </div>,
                route: "/",
                advanceOnTargetClick: false,
                ignoreOutsideClick: false,
                popoverWidth: 500,
            },
        ]
    }

    const get3_7_0 = (): TourStep[] => {
        return [
            {
                id: "changelog-1",
                content: (
                    <div>
                        <h4 className="text-xl font-bold text-white">{t("misc.changelog.whats_new_370")}</h4>
                        <p>{t("misc.changelog.intro_features")}</p>
                    </div>
                ),
                route: "/",
                nextLabel: t("manga.action.start"),
                ignoreOutsideClick: true,
            },
            {
                id: "security",
                title: t("misc.changelog.security_improvements"),
                content: t("misc.changelog.security_content"),
                route: "/",
                advanceOnTargetClick: true,
                ignoreOutsideClick: true,
            },
            {
                id: "search",
                target: "[data-advanced-search-options-tags='true']",
                title: t("library.filter.tags"),
                content: t("misc.changelog.tags_content"),
                route: "/search",
                advanceOnTargetClick: false,
                ignoreOutsideClick: true,
            },
            {
                id: "search",
                target: ".sea-command-content",
                title: t("misc.changelog.adult_global_search"),
                content: t("misc.changelog.adult_content"),
                route: "/search",
                advanceOnTargetClick: false,
                ignoreOutsideClick: true,
                prepare: async () => {
                    setSeaCommandOpen(true)
                    setTimeout(() => {
                        setSeaCommandInput("/search ")
                    }, 200)
                    // wait 500ms
                    return new Promise(resolve => setTimeout(resolve, 500))
                },
            },
            {
                id: "changelog-2",
                title: t("misc.changelog.bug_fixes"),
                content: t("misc.changelog.bug_fixes_content_37"),
                route: "/",
                ignoreOutsideClick: true,
            },
        ]
    }

    const get3_8_0 = (): TourStep[] => {
        return [
            {
                id: "changelog-1",
                content: (
                    <div>
                        <h4 className="text-xl font-bold text-white">{t("misc.changelog.whats_new_380")}</h4>
                        <p>{t("misc.changelog.intro_biggest")}</p>
                    </div>
                ),
                route: "/",
                nextLabel: t("manga.action.start"),
                ignoreOutsideClick: true,
            },
            {
                id: "torrent-search",
                title: t("misc.changelog.torrent_search_downloads"),
                content: t("misc.changelog.torrent_content"),
                route: "/",
                ignoreOutsideClick: true,
            },
            {
                id: "subtitle-translation",
                title: t("misc.changelog.subtitle_translation"),
                content: t("misc.changelog.subtitle_content"),
                route: "/",
                ignoreOutsideClick: true,
            },
            {
                id: "external-player-link",
                target: "[data-settings-external-player-link-scheme]",
                title: t("misc.changelog.local_subtitle_files"),
                content: t("misc.changelog.local_sub_content"),
                route: "/settings",
                prepare: async () => {
                    setSettingsTab("external-player-link")
                    await tourHelpers.waitForSelector("[data-settings-external-player-link-scheme]")
                },
                ignoreOutsideClick: true,
                popoverWidth: 460,
            },
            {
                id: "spoilers",
                target: "[data-settings-hide-anime-spoilers]",
                title: t("misc.changelog.hide_spoilers"),
                content: t("misc.changelog.spoilers_content"),
                route: "/settings",
                prepare: async () => {
                    setSettingsTab("seanime")
                    await tourHelpers.waitForSelector("[data-settings-hide-anime-spoilers]")
                },
                ignoreOutsideClick: true,
                popoverWidth: 460,
            },
            {
                id: "online-streaming",
                target: "[data-settings-enable-onlinestream]",
                title: t("misc.changelog.online_streaming"),
                content: t("misc.changelog.online_content_38"),
                route: "/settings",
                prepare: async () => {
                    setSettingsTab("onlinestream")
                    await tourHelpers.waitForSelector("[data-settings-enable-onlinestream]")
                },
                ignoreOutsideClick: true,
                popoverWidth: 460,
            },
            {
                id: "default-episode-source",
                target: "[data-settings-default-episode-source]",
                title: t("misc.changelog.default_episode_source"),
                content: t("misc.changelog.episode_source_content"),
                route: "/settings",
                prepare: async () => {
                    setSettingsTab("seanime")
                    await tourHelpers.waitForSelector("[data-settings-default-episode-source]")
                },
                ignoreOutsideClick: true,
                popoverWidth: 460,
            },
            {
                id: "ui-settings-redesign",
                target: "[data-settings-ui-panel-tabs]",
                title: t("misc.changelog.redesigned_ui"),
                content: t("misc.changelog.ui_content"),
                route: "/settings",
                prepare: async () => {
                    setSettingsTab("ui")
                    await tourHelpers.waitForSelector("[data-settings-ui-panel-tabs]")
                },
                ignoreOutsideClick: true,
                popoverWidth: 460,
            },
            {
                id: "ui-settings-redesign2",
                target: ".settings-ui-navigation-preloading",
                title: t("misc.changelog.route_preloading"),
                content: t("misc.changelog.preloading_content"),
                prepare: async () => {
                    window.scrollTo({ top: document.body.scrollHeight, behavior: "smooth" })
                    await new Promise(resolve => setTimeout(resolve, 1000))
                    await tourHelpers.waitForSelector(".settings-ui-navigation-preloading")
                },
                ignoreOutsideClick: true,
                popoverWidth: 460,
            },
            {
                id: "entry-header-redesign",
                target: "[data-media-page-header]",
                title: t("misc.changelog.ui_updates"),
                content: t("misc.changelog.ui_updates_content"),
                prepare: async () => {
                    router.push("/entry?id=21827")
                    await tourHelpers.waitForSelector("[data-media-page-header]")
                },
                ignoreOutsideClick: true,
                popoverWidth: 460,
            },
            {
                id: "extensions",
                title: t("misc.changelog.extensions"),
                content: t("misc.changelog.extensions_content"),
                route: "/extensions",
                ignoreOutsideClick: true,
            },
            {
                id: "extension-secure-mode",
                target: "[data-settings-enable-extension-secure-mode]",
                title: t("misc.changelog.extension_secure_mode"),
                content: t("misc.changelog.secure_mode_content"),
                route: "/settings",
                prepare: async () => {
                    setSettingsTab("seanime")
                    await tourHelpers.waitForSelector("[data-settings-enable-extension-secure-mode]")
                },
                ignoreOutsideClick: true,
                popoverWidth: 460,
            },
            {
                id: "denshi",
                title: t("misc.changelog.denshi_window_state"),
                content: t("misc.changelog.denshi_content"),
                route: "/settings",
                prepare: async () => {
                    setSettingsTab("denshi")
                },
                condition: () => typeof window !== "undefined" && !!window.electron,
                conditionFailBehavior: "skip",
                ignoreOutsideClick: true,
            },
            // {
            //     id: "denshi",
            //     title: "View Transitions",
            //     content: "Seanime Denshi now uses the View Transitions API for native transitions between different screens.",
            //     route: "/schedule",
            //     prepare: async () => {
            //         setSettingsTab("seanime")
            //         await new Promise(resolve => setTimeout(resolve, 650))
            //         router.push("/lists")
            //         await new Promise(resolve => setTimeout(resolve, 650))
            //         router.push("/settings")
            //         await new Promise(resolve => setTimeout(resolve, 650))
            //         // scroll to bottom
            //         window.scrollTo({ top: document.body.scrollHeight, behavior: "smooth" })
            //         await new Promise(resolve => setTimeout(resolve, 650))
            //         router.push("/schedule")
            //         await new Promise(resolve => setTimeout(resolve, 650))
            //         router.push("/lists")
            //         await new Promise(resolve => setTimeout(resolve, 650))
            //         router.push("/settings")
            //         await new Promise(resolve => setTimeout(resolve, 650))
            //         router.push("/")
            //         await new Promise(resolve => setTimeout(resolve, 500))
            //     },
            //     condition: () => typeof window !== "undefined" && !!window.electron,
            //     conditionFailBehavior: "skip",
            //     ignoreOutsideClick: true,
            // },
            {
                id: "changelog-2",
                title: t("misc.changelog.bug_fixes"),
                content: t("misc.changelog.bug_fixes_content_38"),
                route: "/",
                ignoreOutsideClick: true,
            },
        ]
    }

    const get3_9_0 = (): TourStep[] => {
        return [
            {
                id: "changelog-1",
                content: (
                    <div>
                        <h4 className="text-xl font-bold text-white">{t("misc.changelog.whats_new_390")}</h4>
                        <p>{t("misc.changelog.intro_biggest")}</p>
                    </div>
                ),
                route: "/",
                nextLabel: t("manga.action.start"),
                ignoreOutsideClick: true,
            },
            {
                id: "libmpv-player",
                target: "[data-tab-trigger='playback']",
                title: t("misc.changelog.new_built_in_player"),
                content: t("misc.changelog.denshi_player_content"),
                route: "/settings",
                prepare: async () => {
                    setSettingsTab("playback")
                    await tourHelpers.waitForSelector("[data-tab-trigger='playback']")
                },
                ignoreOutsideClick: true,
                popoverWidth: 460,
                condition: () => typeof window !== "undefined" && !!window.electron,
                conditionFailBehavior: "skip",
            },
            {
                id: "torrent-streaming-perf",
                title: t("misc.changelog.faster_torrent_streaming"),
                content: t("misc.changelog.faster_torrent_content"),
                route: "/",
                ignoreOutsideClick: true,
            },
            {
                id: "debrid-streaming-perf",
                title: t("misc.changelog.faster_debrid_streaming"),
                content: t("misc.changelog.faster_debrid_content"),
                route: "/",
                ignoreOutsideClick: true,
            },
            {
                id: "changelog-2",
                title: t("misc.changelog.bug_fixes"),
                content: t("misc.changelog.bug_fixes_content_39"),
                route: "/",
                ignoreOutsideClick: true,
            },
        ]
    }

    const get3_10_0 = (): TourStep[] => {
        return [
            {
                id: "changelog-1",
                content: (
                    <div>
                        <h4 className="text-xl font-bold text-white">{t("misc.changelog.whats_new_3100")}</h4>
                        <p>{t("misc.changelog.intro_biggest")}</p>
                    </div>
                ),
                route: "/",
                nextLabel: t("manga.action.start"),
                ignoreOutsideClick: true,
            },
            {
                id: "manga-source-refresh",
                title: t("misc.changelog.manga_source_refresh"),
                content: t("misc.changelog.manga_refresh_content"),
                route: "/manga",
                ignoreOutsideClick: true,
            },
            {
                id: "torrent-availability",
                target: "[data-settings-show-torrent-availability]",
                title: t("misc.changelog.torrent_availability"),
                content: t("misc.changelog.availability_content"),
                route: "/settings",
                prepare: async () => {
                    setSettingsTab("torrent")
                    await tourHelpers.waitForSelector("[data-settings-show-torrent-availability]")
                },
                ignoreOutsideClick: true,
                popoverWidth: 460,
            },
            {
                id: "online-streaming",
                title: t("misc.changelog.online_streaming"),
                content: t("misc.changelog.online_content_310"),
                route: "/",
                ignoreOutsideClick: true,
            },
            {
                id: "plugin-tray-badges",
                title: t("misc.changelog.plugin_tray_badges"),
                content: t("misc.changelog.tray_badges_content"),
                route: "/",
                ignoreOutsideClick: true,
            },
            {
                id: "mpvcore-logs",
                target: "[data-tab-trigger='playback']",
                title: t("misc.changelog.export_mpvcore_logs"),
                content: t("misc.changelog.mpvcore_logs_content"),
                route: "/settings",
                prepare: async () => {
                    setSettingsTab("playback")
                    await tourHelpers.waitForSelector("[data-tab-trigger='playback']")
                },
                condition: () => typeof window !== "undefined" && !!window.electron,
                conditionFailBehavior: "skip",
                ignoreOutsideClick: true,
                popoverWidth: 460,
            },
            {
                id: "changelog-2",
                title: t("misc.changelog.bug_fixes"),
                content: t("misc.changelog.bug_fixes_content_310"),
                route: "/",
                ignoreOutsideClick: true,
            },
        ]
    }

    return {
        "3.5.0": get3_5_0,
        "3.7.0": get3_7_0,
        "3.8.0": get3_8_0,
        "3.9.0": get3_9_0,
        "3.10.0": get3_10_0,
    }
}

export function useChangelogTourListener() {
    const serverStatus = useServerStatus()
    const [seenChangelog, setSeenChangelog] = useAtom(seenChangelogAtom)
    const { start } = useTour()
    const tours = useSetupTour()
    const { width } = useWindowSize()
    const isMobile = width < 768

    const toursRef = React.useRef(tours)
    toursRef.current = tours

    const started = React.useRef(false)
    const timeout = React.useRef<NodeJS.Timeout | null>(null)

    React.useEffect(() => {
        if (!serverStatus?.showChangelogTour) return
        if (serverStatus.isOffline) return
        if (isMobile) return
        if (started.current) return

        if (seenChangelog === serverStatus.showChangelogTour) return

        started.current = true

        const tourId = serverStatus.showChangelogTour

        if (timeout.current) clearTimeout(timeout.current)
        timeout.current = setTimeout(() => {
            const getSteps = toursRef.current[tourId]
            if (getSteps) {
                start(getSteps(), tourId, () => {
                    console.log("tour completed")
                    setSeenChangelog(tourId)
                })
            }
        }, 1000)

        return () => {
            if (timeout.current) clearTimeout(timeout.current)
        }
    }, [serverStatus?.showChangelogTour, serverStatus?.isOffline, seenChangelog, start, setSeenChangelog, isMobile])

    return null
}
