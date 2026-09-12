import { useScanLocalFiles } from "@/api/hooks/scan.hooks"
import { __anilist_userAnimeMediaAtom } from "@/app/(main)/_atoms/anilist.atoms"

import { useSeaCommandInject } from "@/app/(main)/_features/sea-command/use-inject"
import { useServerStatus } from "@/app/(main)/_hooks/use-server-status"
import { AppLayoutStack } from "@/components/ui/app-layout"
import { Button } from "@/components/ui/button"
import { cn } from "@/components/ui/core/styling"
import { Modal } from "@/components/ui/modal"
import { RadioGroup } from "@/components/ui/radio-group"
import { Separator } from "@/components/ui/separator"
import { Switch } from "@/components/ui/switch"
import { useBoolean } from "@/hooks/use-disclosure"
import { t } from "@/lib/i18n"
import { useThemeSettings } from "@/lib/theme/theme-hooks"
import { atom } from "jotai"
import { useAtom } from "jotai/react"
import React from "react"
import { FiSearch } from "react-icons/fi"

export const __scanner_modalIsOpen = atom(false)
export const __scanner_isScanningAtom = atom(false)


export function ScannerModal() {
    const serverStatus = useServerStatus()
    const ts = useThemeSettings()
    const [isOpen, setOpen] = useAtom(__scanner_modalIsOpen)
    const [, setScannerIsScanning] = useAtom(__scanner_isScanningAtom)
    const [userMedia] = useAtom(__anilist_userAnimeMediaAtom)
    const anilistDataOnly = useBoolean(true)
    const skipLockedFiles = useBoolean(true)
    const skipIgnoredFiles = useBoolean(true)
    const enhanceWithOfflineDatabase = useBoolean(true)

    const { mutate: scanLibrary, isPending: isScanning } = useScanLocalFiles(() => {
        setOpen(false)
    })

    React.useEffect(() => {
        if (!userMedia?.length) anilistDataOnly.off()
        else anilistDataOnly.on()
    }, [userMedia])

    React.useEffect(() => {
        setScannerIsScanning(isScanning)
    }, [isScanning])

    function handleScan() {
        scanLibrary({
            enhanced: !anilistDataOnly.active,
            skipLockedFiles: skipLockedFiles.active,
            skipIgnoredFiles: skipIgnoredFiles.active,
            enhanceWithOfflineDatabase: enhanceWithOfflineDatabase.active,
        })
        setOpen(false)
    }

    const { inject, remove } = useSeaCommandInject()
    React.useEffect(() => {
        inject("scanner-controls", {
            priority: 1,
            items: [{
                id: "refresh",
                value: "refresh",
                heading: t("library.command.heading"),
                render: () => (
                    <p>{t("library.scanner.refresh_library")}</p>
                ),
                onSelect: ({ ctx }) => {
                    ctx.close()
                    setTimeout(() => {
                        handleScan()
                    }, 500)
                },
                showBasedOnInput: "startsWith",
            }],
            filter: ({ item, input }) => {
                if (!input) return true
                return item.value.toLowerCase().includes(input.toLowerCase())
            },
            shouldShow: ({ ctx }) => ctx.router.pathname === "/",
            showBasedOnInput: "startsWith",
        })

        return () => remove("scanner-controls")
    }, [])

    return (
        <>
            <Modal
                data-scanner-modal
                open={isOpen}
                onOpenChange={o => {
                    // if (!isScanning) {
                    //     setOpen(o)
                    // }
                    setOpen(o)
                }}
                title={t("library.scanner.title")}
                titleClass="text-center"
                contentClass={cn(
                    "space-y-4 max-w-2xl rounded-xl",
                    // ts.enableBlurringEffects && "bg-gray-950 bg-opacity-80 backdrop-blur-sm firefox:bg-opacity-100 firefox:backdrop-blur-none",
                )}
                overlayClass={cn(ts.enableBlurringEffects && "bg-gray-950/70 backdrop-blur-sm")}
            >
                {/*<GlowingEffect*/}
                {/*    spread={50}*/}
                {/*    glow={true}*/}
                {/*    disabled={false}*/}
                {/*    proximity={100}*/}
                {/*    inactiveZone={0.01}*/}
                {/*    // movementDuration={4}*/}
                {/*    className="!mt-0 opacity-30"*/}
                {/*/>*/}



                {serverStatus?.user?.isSimulated && <div className="border border-dashed rounded-md py-2 px-4 !mt-5">
                    {t("library.scanner.simulated_warning")}
                </div>}

                <div className="space-y-4" data-scanner-modal-content>

                    <AppLayoutStack className="space-y-2">
                        <h5 className="text-[--muted]">{t("library.scanner.local_files")}</h5>
                        <Switch
                            side="right"
                            label={t("library.scanner.skip_locked_files")}
                            value={skipLockedFiles.active}
                            onValueChange={v => skipLockedFiles.set(v as boolean)}
                            // size="lg"
                        />
                        <Switch
                            side="right"
                            label={t("library.scanner.skip_ignored_files")}
                            value={skipIgnoredFiles.active}
                            onValueChange={v => skipIgnoredFiles.set(v as boolean)}
                            // size="lg"
                        />

                        <Separator />

                        <AppLayoutStack className="space-y-2">
                            <h5 className="text-[--muted]">{t("library.scanner.matching_data")}</h5>
                            <Switch
                                side="right"
                                label={t("library.scanner.anilist_only")}
                                moreHelp={t("library.scanner.anilist_only_more_help")}
                                help={anilistDataOnly.active ? t("library.scanner.anilist_only_help") : ""}
                                value={anilistDataOnly.active}
                                onValueChange={v => anilistDataOnly.set(v as boolean)}
                                // className="data-[state=checked]:bg-amber-700 dark:data-[state=checked]:bg-amber-700"
                                // size="lg"

                                disabled={!userMedia?.length}
                            />
                            {!anilistDataOnly.active && <RadioGroup
                                label={t("library.scanner.enhanced_method")}
                                options={[
                                    { value: "database", label: t("library.scanner.use_offline_db") },
                                    { value: "anilist", label: t("library.scanner.use_anilist_api") },
                                ]}
                                size="lg"
                                stackClass="space-y-2 py-1"
                                value={enhanceWithOfflineDatabase.active ? "database" : "anilist"}
                                onValueChange={v => enhanceWithOfflineDatabase.set(v === "database")}
                                help={enhanceWithOfflineDatabase.active
                                    ? <span>{t("library.scanner.enhanced_help")}</span>
                                    : <span><span className="text-[--orange]">{t("library.scanner.slower_for_large")}</span>.
                                        {t("library.scanner.enhanced_more_help")}</span>}
                            />}
                        </AppLayoutStack>

                    </AppLayoutStack>
                </div>
                <Button
                    onClick={handleScan}
                    intent="primary"
                    leftIcon={<FiSearch />}
                    loading={isScanning}
                    className="w-full"
                    disabled={!serverStatus?.settings?.library?.libraryPath}
                >
                    {t("library.scanner.scan")}
                </Button>
            </Modal>
        </>
    )

}
