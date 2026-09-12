import { Extension_Extension, Extension_InvalidExtension } from "@/api/generated/types"
import { useGetAllExtensions, useInstallExternalExtension } from "@/api/hooks/extensions.hooks"
import { AddExtensionModal } from "@/app/(main)/extensions/_containers/add-extension-modal"
import { ExtensionCard } from "@/app/(main)/extensions/_containers/extension-card"
import { InvalidExtensionCard, UnauthorizedExtensionPluginCard } from "@/app/(main)/extensions/_containers/invalid-extension-card"
import { LuffyError } from "@/components/shared/luffy-error"
import { SeaLink } from "@/components/shared/sea-link"
import { AppLayoutStack } from "@/components/ui/app-layout"
import { Button, IconButton } from "@/components/ui/button"
import { Card } from "@/components/ui/card"
import { DropdownMenu, DropdownMenuItem } from "@/components/ui/dropdown-menu"
import { LoadingSpinner } from "@/components/ui/loading-spinner"
import { TextInput } from "@/components/ui/text-input"
import { useRouter } from "@/lib/navigation"
import { t } from "@/lib/i18n"
import { atom, useSetAtom } from "jotai"
import orderBy from "lodash/orderBy"
import React from "react"
import { BiSearch } from "react-icons/bi"
import { BiDotsVerticalRounded } from "react-icons/bi"
import { CgMediaPodcast } from "react-icons/cg"
import { GrInstallOption } from "react-icons/gr"
import { LuBlocks, LuDownload } from "react-icons/lu"
import { MdDataSaverOn } from "react-icons/md"
import { PiBookFill } from "react-icons/pi"
import { RiFolderDownloadFill } from "react-icons/ri"
import { TbReload } from "react-icons/tb"
import { toast } from "sonner"

type ExtensionListProps = {
    children?: React.ReactNode
}

export const __extensions_currentPageAtom = atom<"installed" | "marketplace">("installed")

export const EXTENSION_TYPE = {
    "plugin": t("extensions.type.plugin"),
    "anime-torrent-provider": t("extensions.type.anime_torrent_provider"),
    "manga-provider": t("extensions.type.manga_provider"),
    "onlinestream-provider": t("extensions.type.onlinestream_provider"),
    "custom-source": t("extensions.type.custom_source"),
}

export function ExtensionList(props: ExtensionListProps) {

    const {
        children,
        ...rest
    } = props

    const router = useRouter()

    const [checkForUpdates, setCheckForUpdates] = React.useState(false)
    const [searchTerm, setSearchTerm] = React.useState("")

    const { data: allExtensions, isPending: isLoading, refetch } = useGetAllExtensions(checkForUpdates)

    const setPage = useSetAtom(__extensions_currentPageAtom)

    const {
        mutate: installExtension,
        data: installResponse,
        isPending: isInstalling,
    } = useInstallExternalExtension()

    function orderExtensions(extensions: Extension_Extension[] | undefined) {
        return extensions ?
            orderBy(extensions, ["name", "manifestUri"])
            : []
    }

    function isExtensionInstalled(extensionID: string) {
        return !!allExtensions?.extensions?.find(n => n.id === extensionID) ||
            !!allExtensions?.disabledExtensions?.find(n => n.id === extensionID) ||
            !!allExtensions?.invalidExtensions?.find(n => n.id === extensionID)
    }

    const normalizedSearchTerm = searchTerm.trim().toLowerCase()

    function matchesSearchTerm(extension: Extension_Extension) {
        if (!normalizedSearchTerm) return true

        return [extension.name, extension.description, extension.id]
            .some(value => value?.toLowerCase().includes(normalizedSearchTerm))
    }

    function matchesInvalidExtensionSearchTerm(extension: Extension_InvalidExtension) {
        return matchesSearchTerm(extension.extension)
    }

    const installedExtensions = [
        ...(allExtensions?.extensions ?? []),
        ...(allExtensions?.disabledExtensions ?? []),
    ]

    const enabledExtensions = orderExtensions(allExtensions?.extensions ?? []).filter(matchesSearchTerm)
    const disabledExtensions = orderExtensions(allExtensions?.disabledExtensions ?? []).filter(matchesSearchTerm)

    const pluginExtensions = enabledExtensions.filter(n => n.type === "plugin")
    const animeTorrentExtensions = enabledExtensions.filter(n => n.type === "anime-torrent-provider")
    const mangaExtensions = enabledExtensions.filter(n => n.type === "manga-provider")
    const onlinestreamExtensions = enabledExtensions.filter(n => n.type === "onlinestream-provider")
    const customSourceExtensions = enabledExtensions.filter(n => n.type === "custom-source")

    const nonvalidExtensions = (allExtensions?.invalidExtensions ?? []).filter(n => n.code !== "plugin_permissions_not_granted")
        .filter(matchesInvalidExtensionSearchTerm)
        .sort((a, b) => a.id.localeCompare(b.id))
    const pluginPermissionsNotGrantedExtensions = (allExtensions?.invalidExtensions ?? []).filter(n => n.code === "plugin_permissions_not_granted")
        .filter(matchesInvalidExtensionSearchTerm)
        .sort((a, b) => a.id.localeCompare(b.id))

    const hasVisibleResults =
        pluginPermissionsNotGrantedExtensions.length > 0 ||
        nonvalidExtensions.length > 0 ||
        disabledExtensions.length > 0 ||
        pluginExtensions.length > 0 ||
        customSourceExtensions.length > 0 ||
        animeTorrentExtensions.length > 0 ||
        mangaExtensions.length > 0 ||
        onlinestreamExtensions.length > 0

    if (isLoading) return <LoadingSpinner />

    if (!allExtensions) return <LuffyError>
        {t("extensions.list.fetch_error")}
    </LuffyError>

    return (
        <AppLayoutStack className="gap-6">
            <div className="flex items-center gap-2 flex-wrap">
                <div>
                    <h2>
                        {t("extensions.list.title")}
                    </h2>
                    <p className="text-[--muted] text-sm">
                        {t("extensions.list.subtitle")}
                    </p>
                </div>

                <div className="flex flex-1"></div>

                <div className="flex items-center gap-2 flex-wrap">
                    {!!allExtensions?.hasUpdate?.filter(update => !allExtensions.unsafeExtensions?.[update.extensionID]).length && (
                        <Button
                            className="rounded-full animate-pulse"
                            intent="success"
                            leftIcon={<LuDownload className="text-lg" />}
                            loading={isInstalling}
                            onClick={() => {
                                toast.info(t("extensions.list.toast_installing_updates"))
                                allExtensions?.hasUpdate?.forEach(update => {
                                    if (allExtensions.unsafeExtensions?.[update.extensionID]) {
                                        toast.warning(t("extensions.list.toast_skip_unsafe", { id: update.extensionID }))
                                        return
                                    }
                                    installExtension({
                                        manifestUri: update.manifestURI,
                                    })
                                })
                            }}
                        >
                            {t("extensions.list.update_all")}
                        </Button>
                    )}
                    <Button
                        className="rounded-full"
                        intent="gray-basic"
                        leftIcon={<TbReload className="text-lg" />}
                        disabled={isLoading}
                        onClick={() => {
                            setCheckForUpdates(true)
                            // React.startTransition(() => {
                            //     refetch()
                            // })
                        }}
                    >
                        {t("extensions.list.check_updates")}
                    </Button>
                    <AddExtensionModal extensions={installedExtensions}>
                        <Button
                            className="rounded-full"
                            intent="white-subtle"
                            leftIcon={<GrInstallOption className="text-lg" />}
                        >
                            {t("extensions.list.add_extension")}
                        </Button>
                    </AddExtensionModal>

                    <DropdownMenu trigger={<IconButton icon={<BiDotsVerticalRounded />} intent="gray-basic" />}>

                        <DropdownMenuItem
                            onClick={() => {
                                router.push("/extensions/playground")
                            }}
                        >
                            <span>{t("extensions.list.menu_playground")}</span>
                        </DropdownMenuItem>

                        <DropdownMenuItem
                            onClick={() => {
                                setPage("marketplace")
                            }}
                        >
                            <span>{t("extensions.list.menu_marketplace")}</span>
                        </DropdownMenuItem>
                    </DropdownMenu>
                </div>
            </div>

            <TextInput
                placeholder={t("extensions.list.search_installed_placeholder")}
                value={searchTerm}
                onValueChange={setSearchTerm}
                className="pl-10"
                leftIcon={<BiSearch />}
            />

            {!!searchTerm && !hasVisibleResults && (
                <Card className="p-8 text-center">
                    <p className="text-[--muted]">{t("extensions.list.no_results")}</p>
                </Card>
            )}


            {!!pluginPermissionsNotGrantedExtensions?.length && (
                <Card className="p-4 space-y-6">
                    <h3 className="flex gap-3 items-center">{t("extensions.list.section_permissions_required")}</h3>

                    <div className="grid grid-cols-1 lg:grid-cols-3 2xl:grid-cols-4 gap-4">
                        {pluginPermissionsNotGrantedExtensions.map(extension => (
                            <UnauthorizedExtensionPluginCard
                                key={extension.id}
                                extension={extension}
                                isInstalled={isExtensionInstalled(extension.id)}
                                isUnsafe={allExtensions?.unsafeExtensions?.[extension.id] ?? false}
                            />
                        ))}
                    </div>
                </Card>
            )}
            {!!nonvalidExtensions?.length && (
                <Card className="p-4 space-y-6 border-red-800">

                    <h3 className="flex gap-3 items-center">{t("extensions.list.section_invalid")}</h3>

                    <div className="grid grid-cols-1 lg:grid-cols-3 2xl:grid-cols-4 gap-4">
                        {nonvalidExtensions.map(extension => (
                            <InvalidExtensionCard
                                key={extension.id}
                                extension={extension}
                                isInstalled={isExtensionInstalled(extension.id)}
                            />
                        ))}
                    </div>
                </Card>
            )}

            {!!disabledExtensions?.length && (
                <Card className="p-4 space-y-6">
                    <h3 className="flex gap-3 items-center">{t("extensions.list.section_disabled")}</h3>
                    <div className="grid grid-cols-1 lg:grid-cols-3 2xl:grid-cols-4 gap-4">
                        {disabledExtensions.map(extension => (
                            <ExtensionCard
                                key={extension.id}
                                extension={extension}
                                updateData={allExtensions?.hasUpdate?.find(n => n.extensionID === extension.id)}
                                isInstalled={true}
                                isUnsafe={allExtensions?.unsafeExtensions?.[extension.id] ?? false}
                                isDisabled
                            />
                        ))}
                    </div>
                </Card>
            )}

            {/*<Card className="p-4 space-y-6">*/}

            {!!pluginExtensions?.length && (
                <Card className="p-4 space-y-6">
                    <h3 className="flex gap-3 items-center"><LuBlocks /> {t("extensions.list.section_plugins")}</h3>
                    <div className="grid grid-cols-1 lg:grid-cols-3 2xl:grid-cols-4 gap-4">
                        {pluginExtensions.map(extension => (
                            <ExtensionCard
                                key={extension.id}
                                extension={extension}
                                updateData={allExtensions?.hasUpdate?.find(n => n.extensionID === extension.id)}
                                isInstalled={isExtensionInstalled(extension.id)}
                                userConfigError={allExtensions?.invalidUserConfigExtensions?.find(n => n.id == extension.id)}
                                isUnsafe={allExtensions?.unsafeExtensions?.[extension.id] ?? false}
                                allowReload={true}
                            />
                        ))}
                    </div>
                </Card>
            )}

            {!!customSourceExtensions?.length && (
                <Card className="p-4 space-y-6">
                    <div className="flex items-center gap-4">
                        <h3 className="flex gap-3 items-center"><MdDataSaverOn />{t("extensions.type.custom_source")}</h3>
                        <SeaLink href="/custom-sources" className="text-sm underline underline-offset-2 text-[--muted] hover:text-[--foreground]">
                            {t("extensions.list.browse_all_sources")}
                        </SeaLink>
                    </div>
                    <div className="grid grid-cols-1 lg:grid-cols-3 2xl:grid-cols-4 gap-4">
                        {customSourceExtensions.map(extension => (
                            <ExtensionCard
                                key={extension.id}
                                extension={extension}
                                updateData={allExtensions?.hasUpdate?.find(n => n.extensionID === extension.id)}
                                isInstalled={isExtensionInstalled(extension.id)}
                                userConfigError={allExtensions?.invalidUserConfigExtensions?.find(n => n.id == extension.id)}
                                isUnsafe={allExtensions?.unsafeExtensions?.[extension.id] ?? false}
                                allowReload
                            />
                        ))}
                    </div>
                </Card>
            )}

            {!!animeTorrentExtensions?.length && (
                <Card className="p-4 space-y-6">
                    <h3 className="flex gap-3 items-center"><RiFolderDownloadFill />{t("extensions.marketplace.anime_torrents_section")}</h3>
                    <div className="grid grid-cols-1 lg:grid-cols-3 2xl:grid-cols-4 gap-4">
                        {animeTorrentExtensions.map(extension => (
                            <ExtensionCard
                                key={extension.id}
                                extension={extension}
                                updateData={allExtensions?.hasUpdate?.find(n => n.extensionID === extension.id)}
                                isInstalled={isExtensionInstalled(extension.id)}
                                userConfigError={allExtensions?.invalidUserConfigExtensions?.find(n => n.id == extension.id)}
                                isUnsafe={allExtensions?.unsafeExtensions?.[extension.id] ?? false}
                                allowReload
                            />
                        ))}
                    </div>
                </Card>
            )}


            {!!mangaExtensions?.length && (
                <Card className="p-4 space-y-6">
                    <h3 className="flex gap-3 items-center"><PiBookFill />{t("extensions.marketplace.manga_section")}</h3>
                    <div className="grid grid-cols-1 lg:grid-cols-3 2xl:grid-cols-4 gap-4">
                        {mangaExtensions.map(extension => (
                            <ExtensionCard
                                key={extension.id}
                                extension={extension}
                                updateData={allExtensions?.hasUpdate?.find(n => n.extensionID === extension.id)}
                                isInstalled={isExtensionInstalled(extension.id)}
                                userConfigError={allExtensions?.invalidUserConfigExtensions?.find(n => n.id == extension.id)}
                                isUnsafe={allExtensions?.unsafeExtensions?.[extension.id] ?? false}
                                allowReload
                            />
                        ))}
                    </div>
                </Card>
            )}

            {!!onlinestreamExtensions?.length && (
                <Card className="p-4 space-y-6">
                    <h3 className="flex gap-3 items-center"><CgMediaPodcast /> {t("extensions.marketplace.online_streaming_section")}</h3>
                    <div className="grid grid-cols-1 lg:grid-cols-3 2xl:grid-cols-4 gap-4">
                        {onlinestreamExtensions.map(extension => (
                            <ExtensionCard
                                key={extension.id}
                                extension={extension}
                                updateData={allExtensions?.hasUpdate?.find(n => n.extensionID === extension.id)}
                                isInstalled={isExtensionInstalled(extension.id)}
                                userConfigError={allExtensions?.invalidUserConfigExtensions?.find(n => n.id == extension.id)}
                                isUnsafe={allExtensions?.unsafeExtensions?.[extension.id] ?? false}
                                allowReload
                            />
                        ))}
                    </div>
                </Card>
            )}

            {/*</Card>*/}
        </AppLayoutStack>
    )
}
