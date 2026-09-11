import { useGetAnilistCacheLayerStatus, useToggleAnilistCacheLayerStatus } from "@/api/hooks/anilist.hooks"
import { useListAnimeEntryEpisodeTabExtensions } from "@/api/hooks/extensions.hooks"
import { useLocalSyncSimulatedDataToAnilist } from "@/api/hooks/local.hooks"
import { __seaCommand_shortcuts } from "@/app/(main)/_features/sea-command/sea-command"
import { SettingsCard } from "@/app/(main)/settings/_components/settings-card"
import { SettingsSubmitButton } from "@/app/(main)/settings/_components/settings-submit-button"
import { ConfirmationDialog, useConfirmationDialog } from "@/components/shared/confirmation-dialog"
import { Alert } from "@/components/ui/alert"
import { Button } from "@/components/ui/button"
import { cn } from "@/components/ui/core/styling"
import { Field } from "@/components/ui/form"
import { Switch } from "@/components/ui/switch"
import { __isElectronDesktop__ } from "@/types/constants"
import { useAtom } from "jotai/react"
import React from "react"
import { useFormContext, useWatch } from "react-hook-form"
import { FaRedo } from "react-icons/fa"
import { LuCircleAlert, LuCloudUpload, LuDatabaseBackup, LuEyeOff, LuImageOff, LuImages, LuShield, LuStarOff, LuUserPen } from "react-icons/lu"
import { MdDownloading } from "react-icons/md"
import { RiMovieAiLine } from "react-icons/ri"
import { TbAlertSquareRoundedOff, TbBrowserShare, TbChecklist, TbClockPlay, TbDownloadOff, TbProgressCheck, TbRating18Plus } from "react-icons/tb"
import { useServerStatus } from "../../_hooks/use-server-status"

type ServerSettingsProps = {
    isPending: boolean
}

export function ServerSettings(props: ServerSettingsProps) {

    const {
        isPending,
        ...rest
    } = props

    const serverStatus = useServerStatus()
    const { data: episodeTabExtensions } = useListAnimeEntryEpisodeTabExtensions()

    const [shortcuts, setShortcuts] = useAtom(__seaCommand_shortcuts)
    const f = useFormContext()
    const defaultPlaybackSource = useWatch({ name: "defaultPlaybackSource" })

    const defaultPlaybackSourceOptions = React.useMemo(() => {
        const pluginOptions = Array.from(new Map((episodeTabExtensions ?? []).map(ext => [
            `ext:${ext.id}`,
            {
                value: `ext:${ext.id}`,
                label: ext.tabName ? `${ext.tabName} (${ext.name})` : ext.name,
            },
        ])).values()).sort((a, b) => a.label.localeCompare(b.label))

        const options = [
            { value: "-", label: "Automatic" },
            { value: "library", label: "Local library" },
            ...(serverStatus?.debridSettings?.enabled ? [{ value: "debridstream", label: "Debrid streaming" }] : []),
            ...(serverStatus?.torrentstreamSettings?.enabled ? [{ value: "torrentstream", label: "Torrent streaming" }] : []),
            ...(serverStatus?.settings?.library?.enableOnlinestream ? [{ value: "onlinestream", label: "Online streaming" }] : []),
            ...pluginOptions,
        ]

        if (!!defaultPlaybackSource && defaultPlaybackSource.startsWith("ext:") && !options.some(option => option.value === defaultPlaybackSource)) {
            options.push({ value: defaultPlaybackSource, label: "Unavailable plugin" })
        }

        return options
    }, [episodeTabExtensions, serverStatus, defaultPlaybackSource])

    const { mutate: upload, isPending: isUploading } = useLocalSyncSimulatedDataToAnilist()

    const { data: isApiWorking, isLoading: isFetchingApiStatus } = useGetAnilistCacheLayerStatus()
    const { mutate: toggleCacheLayer, isPending: isTogglingCacheLayer } = useToggleAnilistCacheLayerStatus()

    const confirmDialog = useConfirmationDialog({
        title: "Upload to AniList",
        description: "This will upload your local Seanime collection to your AniList account. Are you sure you want to proceed?",
        actionText: "Upload",
        actionIntent: "primary",
        onConfirm: async () => {
            if (isUploading) return
            upload()
        },
    })

    return (
        <div className="space-y-8">

            {(!isApiWorking && !isFetchingApiStatus) && (
                <Alert
                    intent="warning-basic"
                    description={<div className="space-y-1">
                        <p>The AniList API is not working. All requests will be served from the cache.</p>
                        <p>You can disable this in the app settings.</p>
                    </div>}
                    className="fixed top-4 right-4 z-[50] hidden lg:block"
                />
            )}

            <SettingsCard title="剧集与播放记录">
                <Field.Switch
                    side="right"
                    name="autoUpdateProgress"
                    label="自动更新观看进度"
                    help="启用后，当观看剧集达到 80% 时，将自动同步更新您的追番进度。"
                    moreHelp="仅适用于桌面播放器或内置播放器。"
                    icon={<TbProgressCheck className="" />}
                />
                <Field.Switch
                    side="right"
                    name="enableWatchContinuity"
                    label="启用播放历史记录"
                    help="启用后，Seanime 将记录您的播放进度并在下次打开时从上次离开的位置继续播放。"
                    moreHelp="仅适用于桌面播放器或内置播放器。"
                    icon={<TbClockPlay className="" />}
                />

                <div data-settings-default-episode-source>
                    <Field.Select
                        name="defaultPlaybackSource"
                        label="默认剧集来源"
                        help="打开动漫详情页时的默认数据播放源。"
                        leftIcon={<RiMovieAiLine />}
                        options={defaultPlaybackSourceOptions}
                    />
                </div>
            </SettingsCard>

            <SettingsCard title="动漫展示与防剧透">
                {/*<p className="text-[--muted]">*/}
                {/*    Only applies to desktop and integrated players.*/}
                {/*</p>*/}

                <div className="space-y-3">
                    <div data-settings-hide-anime-spoilers>
                        <Field.Switch
                            side="right"
                            label="隐藏动漫剧透"
                            help="在继续观看、剧集列表以及缺失剧集中隐藏剧透性的缩略图和文本。"
                            name="hideAnimeSpoilers"
                            icon={<LuEyeOff className="" />}
                        />
                    </div>

                    {f.watch("hideAnimeSpoilers") && (
                        <div className="space-y-1 pl-4 border-l border-[--border] ml-2">
                            <Field.Switch
                                side="right"
                                label="隐藏缩略图"
                                name="hideAnimeSpoilerThumbnails"
                            />

                            <Field.Switch
                                side="right"
                                label="隐藏剧集标题"
                                name="hideAnimeSpoilerTitles"
                            />

                            <Field.Switch
                                side="right"
                                label="隐藏剧集简介"
                                name="hideAnimeSpoilerDescriptions"
                            />

                            <Field.Switch
                                side="right"
                                label="下一集跳过防剧透"
                                help="从下一集之后的剧集开始隐藏剧透。"
                                name="hideAnimeSpoilerSkipNextEpisode"
                            />
                        </div>
                    )}
                </div>

                <Field.Switch
                    side="right"
                    name="hideAudienceScore"
                    label="隐藏观众评分"
                    help="启用后，直到您主动点击查看前，番剧观众评分将被遮挡。"
                    icon={<LuStarOff className="" />}
                />


                <div className="space-y-3">
                    <Field.Switch
                        side="right"
                        name="enableAdultContent"
                        label="启用成人/R18内容"
                        help="关闭后，成人向内容将从搜索结果和媒体库中隐藏。"
                        icon={<TbRating18Plus className="" />}
                    />
                    {f.watch("enableAdultContent") && <div className="space-y-1 pl-4 border-l border-[--border] ml-2">
                        <Field.Switch
                            side="right"
                            name="blurAdultContent"
                            label="对成人内容应用模糊遮罩"
                            fieldClass={cn(
                                !f.watch("enableAdultContent") && "opacity-50",
                            )}
                        />
                    </div>}
                </div>

                <Field.Switch
                    side="right"
                    name="disableAnimeCardTrailers"
                    label="禁用动漫卡片预告片"
                    help=""
                    icon={<LuImageOff className="" />}
                />

            </SettingsCard>

            <SettingsCard title="扩展安全性">
                <div data-settings-enable-extension-secure-mode>
                    <Field.Switch
                        side="right"
                        name="enableExtensionSecureMode"
                        label="启用扩展安全模式"
                        help="启用后，即使已授予权限，扩展尝试执行敏感操作时 Seanime 仍会弹出确认提示。"
                        icon={<LuShield className="" />}
                    />
                </div>
            </SettingsCard>

            <SettingsCard
                title="本地账户与数据同步"
                description="未绑定或离线使用 AniList 账号时使用本地账户。"
            >
                <div className={cn(serverStatus?.user?.isSimulated && "opacity-50 pointer-events-none")}>
                    <Field.Switch
                        side="right"
                        name="autoSyncToLocalAccount"
                        label="Automatically back up AniList lists"
                        help="If enabled, your local lists will be periodically updated by using your AniList data. This will override any local changes you've made since the last sync."
                        icon={<LuUserPen className="" />}
                    />
                </div>
                <Button
                    size="sm"
                    intent="primary-subtle"
                    loading={isUploading}
                    leftIcon={<LuCloudUpload className="size-4" />}
                    onClick={() => {
                        confirmDialog.open()
                    }}
                    disabled={serverStatus?.user?.isSimulated}
                >
                    Upload local lists to AniList
                </Button>
            </SettingsCard>

            <ConfirmationDialog {...confirmDialog} />

            <SettingsCard title="离线模式" description="绑定 AniList 账户时可用。">

                <Field.Switch
                    side="right"
                    name="autoSyncOfflineLocalData"
                    label="自动同步离线媒体"
                    help="如果关闭，您需要在离线模式页面中手动点击“立即同步”来刷新本地元数据。"
                    moreHelp="如果您在离线时进行了修改且尚未同步到 AniList，该操作会暂停。"
                    icon={<MdDownloading className="" />}
                />

                <Field.Switch
                    side="right"
                    name="autoSaveCurrentMediaOffline"
                    label="自动离线缓存正在观看/阅读的作品"
                    help="启用后，Seanime 会自动将您正在追的番剧和漫画数据保存为离线可用。"
                    icon={<TbChecklist className="" />}
                />

            </SettingsCard>

            <SettingsCard title="元数据与缓存提供方">
                <div className="space-y-3">
                    <Field.Switch
                        side="right"
                        name="disableCacheLayer"
                        label="禁用 AniList 请求磁盘缓存"
                        help="启用后，Seanime 将不再把 AniList 的请求缓存到磁盘。"
                        moreHelp="默认情况下，所有向 AniList 发起的请求均会缓存到本地，保证网络中断时仍可离线访问。"
                        icon={<LuDatabaseBackup className="" />}
                    />
                    {!f.watch("disableCacheLayer") && (
                        <div className="space-y-1 pl-4 border-l border-[--border] ml-2">
                            <Switch
                                value={!isApiWorking}
                                onValueChange={v => toggleCacheLayer()}
                                disabled={isTogglingCacheLayer}
                                label="启用仅缓存模式"
                                moreHelp="Seanime 将直接使用本地缓存数据，不再向外网 API 发起请求。"
                            />
                        </div>
                    )}
                </div>
                <Field.Switch
                    side="right"
                    name="useFallbackMetadataProvider"
                    label="使用备用剧集元数据源"
                    help="启用后，Seanime 将尝试从备用数据源拉取剧集标题与缩略图元数据。"
                    icon={<LuImages className="" />}
                />
            </SettingsCard>

            <SettingsCard title="版本更新">

                <Field.Switch
                    side="right"
                    name="disableUpdateCheck"
                    label={__isElectronDesktop__ ? "不获取更新" : "不自动检查更新"}
                    help={__isElectronDesktop__ ? (<span className="flex gap-2 items-center">
                        <LuCircleAlert className="size-4 text-[--blue]" />
                        <span>启用后将不再提示新版本。</span>
                    </span>) : "启用后，Seanime 将不再自动检测 GitHub 新版本发布。"}
                    moreHelp={__isElectronDesktop__ ? "Seanime Denshi 桌面端无法关闭静默自动更新。" : undefined}
                    icon={<TbDownloadOff className="" />}
                />
                <Field.Select
                    label="更新渠道"
                    name="updateChannel"
                    help={__isElectronDesktop__ ? "同时也适用于 Seanime Denshi 桌面客户端自动更新。" : ""}
                    options={[
                        { label: "GitHub (默认)", value: "github" },
                        { label: "Seanime 官方源", value: "seanime" },
                        { label: "Seanime (Canary 测试版)", value: "seanime_nightly" },
                    ]}
                />
                {serverStatus?.settings?.library?.updateChannel === "seanime" && (
                    <Alert intent="info" description="您当前正在使用 Seanime 托管的发布渠道。" />
                )}
                {serverStatus?.settings?.library?.updateChannel === "seanime_nightly" && (
                    <Alert
                        intent="warning"
                        description="您当前正在使用 Canary 金丝雀尝鲜分支，可能会接收到未经全面测试的不稳定更新。"
                    />
                )}
            </SettingsCard>

            <SettingsCard title="服务器与系统通知">
                <Field.Switch
                    side="right"
                    name="openWebURLOnStart"
                    label="服务启动时自动打开网页端"
                    icon={<TbBrowserShare className="" />}
                />
                <div className="space-y-3">
                    <Field.Switch
                        side="right"
                        name="disableNotifications"
                        label="禁用系统桌面通知"
                        moreHelp="由操作系统展示的弹窗通知"
                        icon={<TbAlertSquareRoundedOff className="" />}
                    />

                    {!f.watch("disableNotifications") && (
                        <div className="space-y-1 pl-4 border-l border-[--border] ml-2">
                            <Field.Switch
                                // side="right"
                                size="sm"
                                name="disableAutoDownloaderNotifications"
                                label="禁用自动下载器通知"
                            />
                            <Field.Switch
                                // side="right"
                                size="sm"
                                name="disableAutoScannerNotifications"
                                label="禁用自动扫描器通知"
                            />
                        </div>)}
                </div>
            </SettingsCard>

            <SettingsCard title="快捷键设置">
                <div className="space-y-4">
                    {[
                        {
                            label: "Open command palette",
                            value: "meta+j",
                            altValue: "q",
                        },
                    ].map(item => {
                        return (
                            <div className="flex gap-2 items-center" key={item.label}>
                                <label className="text-[--gray]">
                                    <span className="font-semibold">{item.label}</span>
                                </label>
                                <div className="flex gap-2 items-center">
                                    <Button
                                        onKeyDownCapture={(e) => {
                                            e.preventDefault()
                                            e.stopPropagation()

                                            const specialKeys = ["Control", "Shift", "Meta", "Command", "Alt", "Option"]
                                            if (!specialKeys.includes(e.key)) {
                                                const keyStr = `${e.metaKey ? "meta+" : ""}${e.ctrlKey ? "ctrl+" : ""}${e.altKey
                                                    ? "alt+"
                                                    : ""}${e.shiftKey ? "shift+" : ""}${e.key.toLowerCase()
                                                    .replace("arrow", "")
                                                    .replace("insert", "ins")
                                                    .replace("delete", "del")
                                                    .replace(" ", "space")
                                                    .replace("+", "plus")}`

                                                // Update the first shortcut
                                                setShortcuts(prev => [keyStr, prev[1]])
                                            }
                                        }}
                                        className="focus:ring-2 focus:ring-[--brand] focus:ring-offset-1"
                                        size="sm"
                                        intent="white-subtle"
                                    >
                                        {shortcuts[0]}
                                    </Button>
                                    <span className="text-[--muted]">or</span>
                                    <Button
                                        onKeyDownCapture={(e) => {
                                            e.preventDefault()
                                            e.stopPropagation()

                                            const specialKeys = ["Control", "Shift", "Meta", "Command", "Alt", "Option"]
                                            if (!specialKeys.includes(e.key)) {
                                                const keyStr = `${e.metaKey ? "meta+" : ""}${e.ctrlKey ? "ctrl+" : ""}${e.altKey
                                                    ? "alt+"
                                                    : ""}${e.shiftKey ? "shift+" : ""}${e.key.toLowerCase()
                                                    .replace("arrow", "")
                                                    .replace("insert", "ins")
                                                    .replace("delete", "del")
                                                    .replace(" ", "space")
                                                    .replace("+", "plus")}`

                                                // Update the second shortcut
                                                setShortcuts(prev => [prev[0], keyStr])
                                            }
                                        }}
                                        className="focus:ring-2 focus:ring-[--brand] focus:ring-offset-1"
                                        size="sm"
                                        intent="white-subtle"
                                    >
                                        {shortcuts[1]}
                                    </Button>
                                </div>
                                {(shortcuts[0] !== "meta+j" || shortcuts[1] !== "q") && (
                                    <Button
                                        onClick={() => {
                                            setShortcuts(["meta+j", "q"])
                                        }}
                                        className="rounded-full"
                                        size="sm"
                                        intent="white-basic"
                                        leftIcon={<FaRedo />}
                                    >
                                        Reset
                                    </Button>
                                )}
                            </div>
                        )
                    })}
                </div>
            </SettingsCard>

            {/*<Accordion*/}
            {/*    type="single"*/}
            {/*    collapsible*/}
            {/*    className="border rounded-[--radius-md]"*/}
            {/*    triggerClass="dark:bg-[--paper]"*/}
            {/*    contentClass="!pt-2 dark:bg-[--paper]"*/}
            {/*>*/}
            {/*    <AccordionItem value="more">*/}
            {/*        <AccordionTrigger className="bg-gray-900 rounded-[--radius-md]">*/}
            {/*            Advanced*/}
            {/*        </AccordionTrigger>*/}
            {/*        <AccordionContent className="pt-6 flex flex-col md:flex-row gap-3">*/}
            {/*            */}
            {/*        </AccordionContent>*/}
            {/*    </AccordionItem>*/}
            {/*</Accordion>*/}


            <SettingsSubmitButton isPending={isPending} />

        </div>
    )
}

const cardCheckboxStyles = {
    itemContainerClass: cn(
        "block border border-[--border] cursor-pointer transition overflow-hidden w-full",
        "bg-gray-50 hover:bg-[--subtle] dark:bg-gray-950 border-dashed",
        "data-[checked=false]:opacity-30",
        "data-[checked=true]:bg-white dark:data-[checked=true]:bg-gray-950",
        "focus:ring-2 ring-brand-100 dark:ring-brand-900 ring-offset-1 ring-offset-[--background] focus-within:ring-2 transition",
        "data-[checked=true]:border data-[checked=true]:ring-offset-0",
    ),
    itemClass: cn(
        "hidden",
    ),
    // itemLabelClass: cn(
    //     "border-transparent border data-[checked=true]:border-brand dark:bg-transparent dark:data-[state=unchecked]:bg-transparent",
    //     "data-[state=unchecked]:bg-transparent data-[state=unchecked]:hover:bg-transparent dark:data-[state=unchecked]:hover:bg-transparent",
    //     "focus-visible:ring-0 focus-visible:ring-offset-0 focus-visible:ring-offset-transparent",
    // ),
    // itemLabelClass: "font-medium flex flex-col items-center data-[state=checked]:text-[--brand] cursor-pointer",
    stackClass: "flex md:flex-row flex-col space-y-0 gap-4",
}
