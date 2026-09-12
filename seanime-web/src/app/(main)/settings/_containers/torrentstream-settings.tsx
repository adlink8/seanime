import { Models_TorrentstreamSettings } from "@/api/generated/types"
import { useGetTorrentstreamSettings } from "@/api/hooks/torrentstream.hooks"
import { useSaveTorrentstreamSettings, useTorrentstreamDropTorrent } from "@/api/hooks/torrentstream.hooks"
import { useWebsocketMessageListener } from "@/app/(main)/_hooks/handle-websockets.ts"
import { AutoSelectProfileButton } from "@/app/(main)/settings/_components/autoselect-profile-form"
import { SettingsCard } from "@/app/(main)/settings/_components/settings-card"
import { SettingsIsDirty, SettingsSubmitButton } from "@/app/(main)/settings/_components/settings-submit-button"
import { ExperimentalBadge } from "@/components/shared/beta-badge"
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from "@/components/ui/accordion"
import { Alert } from "@/components/ui/alert"
import { Button } from "@/components/ui/button"
import { defineSchema, Field, Form } from "@/components/ui/form"
import { WSEvents } from "@/lib/server/ws-events.ts"
import { t } from "@/lib/i18n"
import React from "react"
import { UseFormReturn } from "react-hook-form"
import { FcFolder } from "react-icons/fc"
import { SiBittorrent } from "react-icons/si"
import { toast } from "sonner"

const torrentstreamSchema = defineSchema(({ z }) => z.object({
    enabled: z.boolean(),
    downloadDir: z.string(),
    autoSelect: z.boolean(),
    disableIPV6: z.boolean(),
    addToLibrary: z.boolean(),
    // streamingServerPort: z.number(),
    // streamingServerHost: z.string(),
    torrentClientHost: z.string().optional().default(""),
    torrentClientPort: z.number(),
    preferredResolution: z.string(),
    includeInLibrary: z.boolean(),
    streamUrlAddress: z.string().optional().default(""),
    slowSeeding: z.boolean().optional().default(false),
    preloadNextStream: z.boolean().optional().default(false),
    disableAcceleratedStartup: z.boolean().optional().default(false),
}))


type TorrentstreamSettingsProps = {
    children?: React.ReactNode
    settings: Models_TorrentstreamSettings | undefined
}

export function TorrentstreamSettings(props: TorrentstreamSettingsProps) {

    const {
        children,
        settings,
        ...rest
    } = props

    const { mutate, isPending } = useSaveTorrentstreamSettings()
    const { refetch } = useGetTorrentstreamSettings()

    const { mutate: dropTorrent, isPending: droppingTorrent } = useTorrentstreamDropTorrent()

    useWebsocketMessageListener({
        type: WSEvents.SETTINGS_CHANGED,
        onMessage: () => {
            refetch()
        },
    })

    const formRef = React.useRef<UseFormReturn<any>>(null)

    if (!settings) return null

    return (
        <>
            <Form
                key={settings?.updatedAt ?? "torrentstream-settings"}
                schema={torrentstreamSchema}
                mRef={formRef}
                onSubmit={data => {
                    if (settings) {
                        mutate({
                                settings: {
                                    ...settings,
                                    ...data,
                                    preferredResolution: data.preferredResolution === "-" ? "" : data.preferredResolution,
                                },
                            },
                            {
                                onSuccess: () => {
                                    formRef.current?.reset(formRef.current.getValues())
                                    toast.success(t("settings.toast.settings_saved"))
                                },
                            },
                        )
                    }
                }}
                defaultValues={{
                    enabled: settings.enabled,
                    autoSelect: settings.autoSelect,
                    downloadDir: settings.downloadDir || "",
                    disableIPV6: settings.disableIPV6,
                    addToLibrary: settings.addToLibrary,
                    // streamingServerPort: settings.streamingServerPort,
                    // streamingServerHost: settings.streamingServerHost || "",
                    torrentClientHost: settings.torrentClientHost || "",
                    torrentClientPort: settings.torrentClientPort,
                    preferredResolution: settings.preferredResolution || "-",
                    includeInLibrary: settings.includeInLibrary,
                    streamUrlAddress: settings.streamUrlAddress || "",
                    slowSeeding: settings.slowSeeding,
                    preloadNextStream: settings.preloadNextStream,
                    disableAcceleratedStartup: settings.disableAcceleratedStartup,
                }}
                stackClass="space-y-8"
            >
                {(f) => (
                    <>
                        <SettingsIsDirty />
                        <SettingsCard>
                            <Field.Switch
                                side="right"
                                name="enabled"
                                label={t("settings.action.enable")}
                            />
                        </SettingsCard>

                        <SettingsCard title={t("settings.debrid.home_screen_title")}>
                            <Field.Switch
                                side="right"
                                name="includeInLibrary"
                                label={t("settings.debrid.include_streaming")}
                                help={t("settings.debrid.include_streaming_help")}
                            />
                        </SettingsCard>

                        <SettingsCard title={t("settings.debrid.autoselect_title")}>
                            <Field.Switch
                                side="right"
                                name="autoSelect"
                                label={t("settings.action.enable")}
                                help={t("settings.torrentstream.autoselect_help")}
                            />

                            <Field.Select
                                name="preferredResolution"
                                label={t("settings.debrid.preferred_resolution")}
                                help={t("settings.debrid.preferred_resolution_help")}
                                options={[
                                    { label: t("settings.option.highest"), value: "-" },
                                    { label: "480p", value: "480" },
                                    { label: "720p", value: "720" },
                                    { label: "1080p", value: "1080" },
                                ]}
                            />

                            <div className="pt-2">
                                <AutoSelectProfileButton />
                            </div>

                            <Field.Switch
                                side="right"
                                name="preloadNextStream"
                                label={<span>{t("settings.torrentstream.preload_next")} <ExperimentalBadge title={t("settings.common.unstable")} /></span>}
                                help={t("settings.torrentstream.preload_next_help")}
                                moreHelp={t("settings.torrentstream.preload_next_more_help")}
                            />
                        </SettingsCard>


                        {/*<Field.Switch
                         side="right"*/}
                        {/*    name="addToLibrary"*/}
                        {/*    label="Add to library"*/}
                        {/*    help="Keep completely downloaded files in corresponding library entries."*/}
                        {/*/>*/}

                        {/* <SettingsCard title="Torrent Client" description="Seanime uses a built-in torrent client to download torrents.">

                         </SettingsCard> */}

                        <Accordion
                            type="single"
                            collapsible
                            className="border rounded-[--radius-md]"
                            triggerClass="dark:bg-[--paper]"
                            contentClass="!pt-2 dark:bg-[--paper]"
                        >
                            <AccordionItem value="more">
                                <AccordionTrigger className="bg-gray-900 rounded-[--radius-md]">
                                    {t("settings.torrentstream.torrent_client")}
                                </AccordionTrigger>
                                <AccordionContent className="space-y-4">
                                    <div className="flex items-center gap-3">

                                        <Field.Text
                                            name="torrentClientHost"
                                            label={t("settings.field.host")}
                                            help={t("settings.torrentstream.host_help")}
                                        />

                                        <Field.Number
                                            name="torrentClientPort"
                                            label={t("settings.field.port")}
                                            formatOptions={{
                                                useGrouping: false,
                                            }}
                                            help={t("settings.torrentstream.port_help")}
                                        />

                                    </div>

                                    <Field.Switch
                                        side="right"
                                        name="disableIPv6"
                                        label={t("settings.torrentstream.disable_ipv6")}
                                    />

                                    <Field.Switch
                                        side="right"
                                        name="slowSeeding"
                                        label={t("settings.torrentstream.slow_seeding")}
                                        moreHelp={t("settings.torrentstream.slow_seeding_more_help")}
                                    />

                                    <Field.Switch
                                        side="right"
                                        name="disableAcceleratedStartup"
                                        label={t("settings.torrentstream.disable_accelerated_startup")}
                                        disabled={f.watch("slowSeeding")}
                                        moreHelp={t("settings.torrentstream.disable_accelerated_startup_more_help")}
                                    />
                                </AccordionContent>
                            </AccordionItem>
                        </Accordion>

                        <Accordion
                            type="single"
                            collapsible
                            className="border rounded-[--radius-md]"
                            triggerClass="dark:bg-[--paper]"
                            contentClass="!pt-2 dark:bg-[--paper]"
                        >
                            <AccordionItem value="more">
                                <AccordionTrigger className="bg-gray-900 rounded-[--radius-md]">
                                    {t("settings.common.advanced")}
                                </AccordionTrigger>
                                <AccordionContent className="pt-6 space-y-4">
                                    <Field.Text
                                        name="streamUrlAddress"
                                        label={t("settings.torrentstream.stream_url_address")}
                                        placeholder={t("settings.torrentstream.stream_url_placeholder")}
                                        help={t("settings.torrentstream.stream_url_help")}
                                    />

                                    <Field.DirectorySelector
                                        name="downloadDir"
                                        label={t("settings.torrentstream.cache_directory")}
                                        leftIcon={<FcFolder />}
                                        help={t("settings.torrentstream.cache_directory_help")}
                                        shouldExist
                                    />
                                    <Alert
                                        intent="warning"
                                        description={t("settings.torrentstream.cache_directory_warning")}
                                    />
                                </AccordionContent>
                            </AccordionItem>
                        </Accordion>


                        <div className="flex w-full items-center">
                            <SettingsSubmitButton isPending={isPending} />
                            <div className="flex flex-1"></div>
                            <Button
                                leftIcon={<SiBittorrent />} intent="alert-subtle" onClick={() => dropTorrent()}
                                disabled={droppingTorrent}
                            >
                                {t("settings.torrentstream.drop_torrent")}
                            </Button>
                        </div>
                    </>
                )}
            </Form>
        </>
    )
}
