import { Anime_AutoDownloaderRule } from "@/api/generated/types"
import {
    useDeleteAutoDownloaderRule,
    useGetAutoDownloaderItems,
    useGetAutoDownloaderProfiles,
    useGetAutoDownloaderRules,
    useRunAutoDownloader,
} from "@/api/hooks/auto_downloader.hooks"
import { useAnimeListTorrentProviderExtensions } from "@/api/hooks/extensions.hooks"
import { useSaveAutoDownloaderSettings } from "@/api/hooks/settings.hooks"
import { __anilist_userAnimeMediaAtom } from "@/app/(main)/_atoms/anilist.atoms"
import { useServerStatus } from "@/app/(main)/_hooks/use-server-status"
import { AutoDownloaderRuleItem } from "@/app/(main)/auto-downloader/_components/autodownloader-rule-item"
import { AutoDownloaderBatchRuleForm } from "@/app/(main)/auto-downloader/_containers/autodownloader-batch-rule-form"
import { AutoDownloaderProfiles } from "@/app/(main)/auto-downloader/_containers/autodownloader-profiles"
import { AutodownloaderQueue } from "@/app/(main)/auto-downloader/_containers/autodownloader-queue"
import { AutoDownloaderRuleForm } from "@/app/(main)/auto-downloader/_containers/autodownloader-rule-form"
import { SettingsCard } from "@/app/(main)/settings/_components/settings-card"
import { ConfirmationDialog, useConfirmationDialog } from "@/components/shared/confirmation-dialog"
import { Alert } from "@/components/ui/alert"
import { Badge } from "@/components/ui/badge"
import { Button, IconButton } from "@/components/ui/button"
import { Card } from "@/components/ui/card"
import { cn } from "@/components/ui/core/styling"
import { Drawer } from "@/components/ui/drawer"
import { DropdownMenu, DropdownMenuItem } from "@/components/ui/dropdown-menu"
import { defineSchema, Field, Form } from "@/components/ui/form"
import { LoadingSpinner } from "@/components/ui/loading-spinner"
import { Modal } from "@/components/ui/modal"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { useBoolean } from "@/hooks/use-disclosure"
import { t } from "@/lib/i18n"
import { useAtomValue } from "jotai/react"
import React from "react"
import { BiDotsVerticalRounded } from "react-icons/bi"
import { FaSquareRss } from "react-icons/fa6"
import { LuTrash } from "react-icons/lu"
import { MdOutlineAdd } from "react-icons/md"
import { toast } from "sonner"
import { z } from "zod"

const tabContentClass = cn(
    "space-y-4 animate-in fade-in-0 duration-300",
)

const settingsSchema = defineSchema(({ z, presets }) => z.object({
    provider: presets.multiSelect,
    interval: z.number().transform(n => {
        if (n < 15) {
            toast.info(t("autodownloader.toast.interval_min"))
            return 15
        }
        return n
    }),
    enabled: z.boolean(),
    downloadAutomatically: z.boolean(),
    enableEnhancedQueries: z.boolean(),
    enableSeasonCheck: z.boolean(),
    useDebrid: z.boolean(),
}))

export function AutoDownloaderPage() {
    const serverStatus = useServerStatus()
    const userMedia = useAtomValue(__anilist_userAnimeMediaAtom)
    const { data: extensions, isLoading: isLoadingExtensions } = useAnimeListTorrentProviderExtensions()

    const [tab, setTab] = React.useState("rules")

    const createRuleModal = useBoolean(false)
    const createBatchRuleModal = useBoolean(false)

    const { mutate: runAutoDownloader, isPending: isRunning } = useRunAutoDownloader()

    const { mutate: updateSettings, isPending } = useSaveAutoDownloaderSettings()

    const { data, isLoading } = useGetAutoDownloaderRules()

    const { data: items, isLoading: itemsLoading } = useGetAutoDownloaderItems()
    const { data: profiles } = useGetAutoDownloaderProfiles()

    const { mutate: deleteNoLongerAiring, isPending: deletingRule } = useDeleteAutoDownloaderRule(-1)

    const confirmDeleteNoLongerAiring = useConfirmationDialog({
        title: t("autodownloader.action.remove_finished"),
        description: t("autodownloader.confirm_remove_finished_desc"),
        onConfirm: () => {
            deleteNoLongerAiring()
        },
    })

    function handleSaveSettings(data: z.infer<typeof settingsSchema>) {
        updateSettings({
            ...data,
            provider: !!data.provider?.length ? data.provider[0] : "",
        })
    }

    function sortRules(a: Anime_AutoDownloaderRule, b: Anime_AutoDownloaderRule) {
        const mediaA = userMedia?.find(m => m.id === a.mediaId)
        const mediaB = userMedia?.find(m => m.id === b.mediaId)
        if (mediaA && !mediaB) return 1
        if (!mediaA && mediaB) return -1
        if (!mediaA && !mediaB) return 0
        if (mediaA?.status !== mediaB?.status) {
            if (mediaA?.status === "RELEASING") return -1
            if (mediaB?.status === "RELEASING") return 1
            if (mediaA?.status === "FINISHED") return 1
            if (mediaB?.status === "FINISHED") return -1
            if (mediaA?.status === "NOT_YET_RELEASED") return 1
            if (mediaB?.status === "NOT_YET_RELEASED") return -1
        }
        return mediaA?.title?.userPreferred?.localeCompare(mediaB?.title?.userPreferred ?? "") ?? 0
    }

    return (
        <div className="space-y-4">
            <ConfirmationDialog {...confirmDeleteNoLongerAiring} />

            <Tabs
                variant="pill"
                value={tab}
                onValueChange={setTab}
            >
                <TabsList>
                    <TabsTrigger value="rules">{t("autodownloader.tab.rules")}</TabsTrigger>
                    <TabsTrigger value="profiles">{t("autodownloader.tab.profiles")}</TabsTrigger>
                    <TabsTrigger value="queue">
                        {t("manga.downloads.queue")}
                        {!!items?.length && (
                            <Badge className="ml-2 font-bold" intent="alert" size="sm">
                                {items.length}
                            </Badge>
                        )}
                    </TabsTrigger>
                    <TabsTrigger value="settings">{t("player.menu.title")}</TabsTrigger>
                </TabsList>
                <TabsContent value="rules" className={tabContentClass}>
                    <div className="pt-4">
                        {(isLoading && isLoadingExtensions) && <LoadingSpinner />}
                        {(!isLoading && !isLoadingExtensions) && (
                            <div className="space-y-4">

                                {!serverStatus?.settings?.autoDownloader?.enabled && (
                                    <Alert
                                        intent="warning"
                                        description={<p>
                                            {t("autodownloader.alert.disabled")} <Button
                                            className="py-0 h-auto"
                                            intent="white-link"
                                            onClick={() => setTab("settings")}
                                        >{t("autodownloader.alert.enable_here")}</Button>
                                        </p>}
                                    />
                                )}

                                <Card className="p-4 space-y-4">
                                    <ul className="text-base text-[--muted]">
                                        <li>{t("autodownloader.rules_desc")}
                                        </li>
                                    </ul>

                                    <div className="w-full flex items-center gap-2">
                                        <DropdownMenu
                                            trigger={<Button
                                                className="rounded-full"
                                                intent="white-subtle"
                                                leftIcon={<MdOutlineAdd className="text-lg" />}

                                            >
                                                {t("library.auto_downloader.new_rule")}
                                            </Button>}
                                        >
                                            <DropdownMenuItem onClick={createRuleModal.on}>
                                                {t("autodownloader.action.one_series")}
                                            </DropdownMenuItem>
                                            <DropdownMenuItem onClick={createBatchRuleModal.on}>
                                                {t("autodownloader.action.multiple_series")}
                                            </DropdownMenuItem>
                                        </DropdownMenu>
                                        <div className="flex flex-1"></div>
                                        <Button
                                            className=""
                                            intent="gray-basic"
                                            leftIcon={<FaSquareRss />}
                                            onClick={() => {
                                                runAutoDownloader()
                                            }}
                                            loading={isRunning}
                                            disabled={!serverStatus?.settings?.autoDownloader?.enabled}
                                        >
                                            {t("autodownloader.action.check_rss")}
                                        </Button>
                                        <DropdownMenu
                                            trigger={<IconButton
                                                className=""
                                                intent="gray-basic"
                                                icon={<BiDotsVerticalRounded className="text-lg" />}
                                            />}
                                        >
                                            <DropdownMenuItem
                                                onClick={confirmDeleteNoLongerAiring.open}
                                                className="text-[--red]"
                                                disabled={deletingRule}
                                            >
                                                <LuTrash /> {t("autodownloader.action.remove_finished")}
                                            </DropdownMenuItem>
                                        </DropdownMenu>
                                    </div>

                                    {(!data?.length) && <div className="p-4 text-[--muted] text-center">{t("autodownloader.empty.no_rules")}</div>}
                                    {(!!data?.length) && <div className="space-y-2">
                                        {data?.toSorted(sortRules)?.map(rule => (
                                            <AutoDownloaderRuleItem
                                                key={rule.dbId}
                                                rule={rule}
                                                userMedia={userMedia}
                                                profiles={profiles ?? []}
                                                extensions={extensions ?? []}
                                            />
                                        ))}
                                    </div>}
                                </Card>
                            </div>
                        )}
                    </div>
                </TabsContent>

                <TabsContent value="profiles" className={tabContentClass}>
                    <AutoDownloaderProfiles />
                </TabsContent>

                <TabsContent value="queue" className={tabContentClass}>

                    <div className="pt-4">
                        <AutodownloaderQueue items={items} isLoading={itemsLoading} />
                    </div>

                </TabsContent>

                <TabsContent value="settings" className={tabContentClass}>
                    <div className="pt-4">
                        <Form
                            schema={settingsSchema}
                            onSubmit={handleSaveSettings}
                            defaultValues={{
                                provider: serverStatus?.settings?.autoDownloader?.provider ? [serverStatus.settings.autoDownloader.provider] : [],
                                enabled: serverStatus?.settings?.autoDownloader?.enabled ?? false,
                                interval: serverStatus?.settings?.autoDownloader?.interval || 15,
                                downloadAutomatically: serverStatus?.settings?.autoDownloader?.downloadAutomatically ?? false,
                                enableEnhancedQueries: serverStatus?.settings?.autoDownloader?.enableEnhancedQueries ?? false,
                                enableSeasonCheck: serverStatus?.settings?.autoDownloader?.enableSeasonCheck ?? false,
                                useDebrid: serverStatus?.settings?.autoDownloader?.useDebrid ?? false,
                            }}
                            stackClass="space-y-4"
                        >
                            {(f) => (
                                <>
                                    <SettingsCard>
                                        <Field.Switch
                                            side="right"
                                            label={t("settings.action.enable")}
                                            name="enabled"
                                        />

                                        <Field.Switch
                                            side="right"
                                            label={t("autodownloader.field.use_debrid")}
                                            name="useDebrid"
                                        />

                                        {f.watch("useDebrid") && !(serverStatus?.debridSettings?.enabled && !!serverStatus?.debridSettings?.provider) && (
                                            <Alert
                                                intent="alert"
                                                title={t("autodownloader.alert.debrid_title")}
                                                description={t("autodownloader.alert.debrid_desc")}
                                            />
                                        )}

                                        <Field.Combobox
                                            name="provider"
                                            options={extensions?.toSorted((a, b) => a.id.localeCompare(b.id))?.map(ext => ({
                                                label: ext.name,
                                                textValue: ext.name,
                                                value: ext.id,
                                            })) ?? []}
                                            label={t("settings.field.default_provider")}
                                            emptyMessage={t("autodownloader.empty.no_extensions")}
                                        />
                                    </SettingsCard>

                                    <SettingsCard
                                        className={cn(
                                            !f.watch("enabled") && "pointer-events-none opacity-50",
                                        )}
                                    >
                                        {/*<Field.Switch*/}
                                        {/*    side="right"*/}
                                        {/*    label="Use smart search queries"*/}
                                        {/*    name="enableEnhancedQueries"*/}
                                        {/*    help="Seanime will use smart search queries for more targeted results, if the extension allows it."*/}
                                        {/*/>*/}
                                        <Field.Switch
                                            side="right"
                                            label={t("autodownloader.field.download_immediately")}
                                            name="downloadAutomatically"
                                            help={t("autodownloader.help.download_immediately")}
                                        />
                                        <Field.Number
                                            label={t("autodownloader.field.interval")}
                                            help={t("autodownloader.help.interval")}
                                            name="interval"
                                            leftAddon={t("autodownloader.field.every")}
                                            rightAddon={t("entry.episode.minutes")}
                                            size="sm"
                                            className="text-center w-20"
                                            min={15}
                                        />
                                    </SettingsCard>

                                    <SettingsCard
                                        className={cn(
                                            !f.watch("enabled") && "pointer-events-none opacity-50",
                                        )}
                                    >
                                        <Field.Switch
                                            side="right"
                                            label={t("autodownloader.field.strict_season_check")}
                                            name="enableSeasonCheck"
                                            help={t("autodownloader.help.strict_season_check")}
                                        />
                                    </SettingsCard>

                                    <Field.Submit role="save" loading={isPending}>{t("settings.action.save")}</Field.Submit>
                                </>
                            )}
                        </Form>
                    </div>
                </TabsContent>

            </Tabs>


            <Modal
                open={createRuleModal.active}
                onOpenChange={createRuleModal.off}
                title={t("library.auto_downloader.new_rule")}
                contentClass="max-w-4xl"
            >
                <AutoDownloaderRuleForm type="create" onRuleCreatedOrDeleted={() => createRuleModal.off()} />
            </Modal>


            <Drawer
                open={createBatchRuleModal.active}
                onOpenChange={createBatchRuleModal.off}
                title={t("autodownloader.modal.create_rules")}
                size="xl"
            >
                <p className="text-[--muted] py-4">
                    {t("autodownloader.batch_intro")}
                </p>
                <AutoDownloaderBatchRuleForm onRuleCreated={() => createBatchRuleModal.off()} rules={data ?? []} />
            </Drawer>
        </div>
    )

}


