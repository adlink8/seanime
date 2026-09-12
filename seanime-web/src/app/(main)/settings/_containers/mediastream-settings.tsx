import { useGetMediastreamSettings, useSaveMediastreamSettings } from "@/api/hooks/mediastream.hooks"
import { useWebsocketMessageListener } from "@/app/(main)/_hooks/handle-websockets.ts"
import { useServerStatus } from "@/app/(main)/_hooks/use-server-status"
import { useMediastreamActiveOnDevice } from "@/app/(main)/mediastream/_lib/mediastream.atoms"
import { SettingsCard, SettingsPageHeader } from "@/app/(main)/settings/_components/settings-card"
import { SettingsIsDirty, SettingsSubmitButton } from "@/app/(main)/settings/_components/settings-submit-button"
import { defineSchema, Field, Form } from "@/components/ui/form"
import { LoadingSpinner } from "@/components/ui/loading-spinner"
import { WSEvents } from "@/lib/server/ws-events.ts"
import { t } from "@/lib/i18n"
import React from "react"
import { UseFormReturn } from "react-hook-form"
import { LuTabletSmartphone } from "react-icons/lu"

const mediastreamSchema = defineSchema(({ z }) => z.object({
    transcodeEnabled: z.boolean(),
    transcodeHwAccel: z.string(),
    transcodePreset: z.string().min(2),
    // transcodeThreads: z.number(),
    // preTranscodeEnabled: z.boolean(),
    // preTranscodeLibraryDir: z.string(),
    disableAutoSwitchToDirectPlay: z.boolean(),
    directPlayOnly: z.boolean(),
    ffmpegPath: z.string().min(0),
    ffprobePath: z.string().min(0),
    transcodeHwAccelCustomSettings: z.string().min(0),
}))

const MEDIASTREAM_HW_ACCEL_OPTIONS = [
    { label: t("settings.mediastream.hwaccel_cpu"), value: "cpu" },
    { label: "NVIDIA (NVENC)", value: "nvidia" },
    { label: "Intel (QSV)", value: "qsv" },
    { label: "Intel (QSV Low-Power)", value: "qsv-low-power" },
    { label: "VAAPI", value: "vaapi" },
    { label: "Apple VideoToolbox", value: "videotoolbox" },
    { label: t("settings.ui.banner_type_custom"), value: "custom" },
]

const MEDIASTREAM_PRESET_OPTIONS = [
    { label: "Ultrafast", value: "ultrafast" },
    { label: "Superfast", value: "superfast" },
    { label: "Veryfast", value: "veryfast" },
    { label: "Fast", value: "fast" },
    { label: "Medium", value: "medium" },
]

type MediastreamSettingsProps = {
    children?: React.ReactNode
}

export function MediastreamSettings(props: MediastreamSettingsProps) {

    const {
        children,
        ...rest
    } = props

    const serverStatus = useServerStatus()

    const { data: settings, isLoading, refetch } = useGetMediastreamSettings(true)

    const { mutate, isPending } = useSaveMediastreamSettings()

    const { activeOnDevice, setActiveOnDevice } = useMediastreamActiveOnDevice()

    useWebsocketMessageListener({
        type: WSEvents.SETTINGS_CHANGED,
        onMessage: () => {
            refetch()
        },
    })

    const formRef = React.useRef<UseFormReturn<any>>(null)

    if (!settings) return <LoadingSpinner />

    return (
        <>
            <SettingsPageHeader
                title={t("settings.mediastream.title")}
                description={t("settings.mediastream.desc")}
                icon={LuTabletSmartphone}
            />

            <Form
                key={settings?.updatedAt ?? "mediastream-settings"}
                schema={mediastreamSchema}
                mRef={formRef}
                onSubmit={data => {
                    if (settings) {
                        mutate({
                                settings: {
                                    ...settings,
                                    ...data,
                                    preTranscodeLibraryDir: "",
                                    preTranscodeEnabled: false,
                                    transcodeThreads: 0,
                                },
                            },
                            {
                                onSuccess: () => {
                                    formRef.current?.reset(formRef.current.getValues())
                                },
                            },
                        )
                    }
                }}
                defaultValues={{
                    transcodeEnabled: settings?.transcodeEnabled ?? false,
                    transcodeHwAccel: settings?.transcodeHwAccel === "none" ? "cpu" : settings?.transcodeHwAccel || "cpu",
                    transcodePreset: settings?.transcodePreset || "fast",
                    // transcodeThreads: settings?.transcodeThreads,
                    // preTranscodeEnabled: settings?.preTranscodeEnabled ?? false,
                    // preTranscodeLibraryDir: settings?.preTranscodeLibraryDir,
                    disableAutoSwitchToDirectPlay: settings?.disableAutoSwitchToDirectPlay ?? false,
                    directPlayOnly: settings?.directPlayOnly ?? false,
                    ffmpegPath: settings?.ffmpegPath || "",
                    ffprobePath: settings?.ffprobePath || "",
                    transcodeHwAccelCustomSettings: settings?.transcodeHwAccelCustomSettings || "{\n	\"name\": \"\",\n	\"decodeFlags\": [\n		\"-hwaccel\", \"\",\n		\"-hwaccel_output_format\", \"\",\n	],\n	\"encodeFlags\": [\n		\"-c:v\", \"\",\n		\"-preset\", \"\",\n		\"-pix_fmt\", \"yuv420p\",\n	],\n	\"scaleFilter\": \"scale=%d:%d\"\n}",
                }}
                stackClass="space-y-8"
            >
                {(f) => (
                    <>
                        <SettingsIsDirty />
                        <SettingsCard>
                            <Field.Switch
                                side="right"
                                name="transcodeEnabled"
                                label={t("settings.action.enable")}
                            />
                        </SettingsCard>

                        {/* <SettingsCard title="Client Playback">
                         <div className="flex gap-4 items-center rounded-[--radius-md]">
                         <MdOutlineDevices className="text-4xl" />
                         <div className="space-y-1">
                         <Checkbox
                         value={activeOnDevice ?? false}
                         onValueChange={v => {
                         setActiveOnDevice((prev) => typeof v === "boolean" ? v : prev)
                         if (v) {
                         toast.success("Media streaming is now active on this device.")
                         } else {
                         toast.info("Media streaming is now inactive on this device.")
                         }
                         }}
                         label="Use media streaming on this device"
                         help="Enable this option if you want to use media streaming on this device."
                         />
                         <p className="text-gray-200">
                         Current client: {serverStatus?.clientDevice}, {serverStatus?.clientPlatform}
                         </p>
                         </div>
                         </div>

                         {(f.watch("transcodeEnabled") && activeOnDevice) && (
                         <Alert
                         intent="info" description={<>
                         Media streaming will be used instead of your external player on this device.
                         </>}
                         />
                         )}
                         </SettingsCard> */}

                        <SettingsCard title={t("settings.mediastream.direct_play_title")}>

                            <Field.Switch
                                side="right"
                                name="disableAutoSwitchToDirectPlay"
                                label={t("settings.mediastream.prefer_transcoding")}
                                help={t("settings.mediastream.prefer_transcoding_help")}
                            />

                            <Field.Switch
                                side="right"
                                name="directPlayOnly"
                                label={t("settings.mediastream.direct_play_only")}
                                help={t("settings.mediastream.direct_play_only_help")}
                            />

                        </SettingsCard>

                        <SettingsCard title={t("settings.mediastream.transcoding_title")}>
                            <Field.Select
                                options={MEDIASTREAM_HW_ACCEL_OPTIONS}
                                name="transcodeHwAccel"
                                label={t("settings.mediastream.hardware_acceleration")}
                                help={t("settings.mediastream.hardware_acceleration_help")}
                            />

                            {f.watch("transcodeHwAccel") === "custom" && (
                                <Field.Textarea
                                    name="transcodeHwAccelCustomSettings"
                                    label={t("settings.mediastream.custom_settings_json")}
                                    className="min-h-[400px]"
                                    help={t("settings.mediastream.custom_settings_help")}
                                />
                            )}

                            <Field.Select
                                options={MEDIASTREAM_PRESET_OPTIONS}
                                name="transcodePreset"
                                label={t("settings.mediastream.transcode_preset")}
                                help={t("settings.mediastream.transcode_preset_help")}
                            />
                        </SettingsCard>

                        <SettingsCard title="FFmpeg">

                            <div className="flex gap-3 items-center">
                                <Field.Text
                                    name="ffmpegPath"
                                    label={t("settings.mediastream.ffmpeg_path")}
                                    help={t("settings.mediastream.ffmpeg_path_help")}
                                />

                                <Field.Text
                                    name="ffprobePath"
                                    label={t("settings.mediastream.ffprobe_path")}
                                    help={t("settings.mediastream.ffprobe_path_help")}
                                />
                            </div>
                        </SettingsCard>

                        <SettingsSubmitButton isPending={isPending} />
                    </>
                )}
            </Form>

            {/*<Separator />*/}

            {/*<h2>Cache</h2>*/}

            {/*<div className="space-y-4">*/}
            {/*    <div className="flex gap-2 items-center">*/}
            {/*        <Button intent="white-subtle" size="sm" onClick={() => getTotalSize()} disabled={isFetchingSize}>*/}
            {/*            Show total size*/}
            {/*        </Button>*/}
            {/*        {!!totalSize && (*/}
            {/*            <p>*/}
            {/*                {totalSize}*/}
            {/*            </p>*/}
            {/*        )}*/}
            {/*    </div>*/}
            {/*    <div className="flex gap-2 flex-wrap items-center">*/}
            {/*        <Button intent="alert-subtle" size="sm" onClick={() => clearCache()} disabled={isClearing}>*/}
            {/*            Clear cache*/}
            {/*        </Button>*/}
            {/*    </div>*/}
            {/*</div>*/}
        </>
    )
}
