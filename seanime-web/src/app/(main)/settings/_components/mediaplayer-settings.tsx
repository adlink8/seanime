import { useExternalPlayerLink } from "@/app/(main)/_atoms/playback.atoms"
import { useServerStatus } from "@/app/(main)/_hooks/use-server-status"
import { SettingsCard, SettingsPageHeader } from "@/app/(main)/settings/_components/settings-card"
import { SettingsSubmitButton } from "@/app/(main)/settings/_components/settings-submit-button"
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from "@/components/ui/accordion"
import { Alert } from "@/components/ui/alert"
import { Field } from "@/components/ui/form"
import { Switch } from "@/components/ui/switch"
import { TextInput } from "@/components/ui/text-input"
import { getDefaultIinaSocket } from "@/lib/server/settings"
import { t } from "@/lib/i18n"
import React from "react"
import { useWatch } from "react-hook-form"
import { FcClapperboard, FcVideoCall, FcVlc } from "react-icons/fc"
import { HiPlay } from "react-icons/hi"
import { IoPlayForwardCircleSharp } from "react-icons/io5"
import { LuCircleArrowOutUpRight, LuMonitorPlay } from "react-icons/lu"
import { RiSettings3Fill } from "react-icons/ri"

type MediaplayerSettingsProps = {
    isPending: boolean
}

export function MediaplayerSettings(props: MediaplayerSettingsProps) {

    const {
        isPending,
    } = props

    const serverStatus = useServerStatus()
    const selectedPlayer = useWatch({ name: "defaultPlayer" })

    return (
        <>
            <SettingsPageHeader
                title={t("player.desktop.title")}
                description={t("player.desktop.desc")}
                icon={LuMonitorPlay}
            />

            <SettingsCard>
                <Field.Select
                    name="defaultPlayer"
                    label={t("player.desktop.default_player")}
                    leftIcon={<FcVideoCall />}
                    options={[
                        { label: "MPV", value: "mpv" },
                        { label: "VLC", value: "vlc" },
                        { label: t("settings.mediaplayer.option_mpc"), value: "mpc-hc" },
                        { label: t("settings.mediaplayer.option_iina"), value: "iina" },
                    ]}
                    help={t("player.desktop.default_player_help")}
                />
                {selectedPlayer === "iina" && <Alert
                    intent="info-basic"
                    description={<p>{t("settings.mediaplayer.iina_notice_1")} <strong>{t("settings.mediaplayer.iina_notice_quit")}</strong>{t("settings.mediaplayer.iina_notice_2")} <strong>{t("settings.mediaplayer.iina_notice_keep")}</strong>{t("settings.mediaplayer.iina_notice_3")}</p>}
                />}
            </SettingsCard>

            <SettingsCard title={t("player.autoplay.title")}>
                <Field.Switch
                    side="right"
                    name="autoPlayNextEpisode"
                    label={t("player.autoplay.next_label")}
                    help={t("player.autoplay.next_help")}
                />
            </SettingsCard>

            <SettingsCard title={t("player.desktop.config_title")}>


                <Field.Text
                    name="mediaPlayerHost"
                    label={t("settings.field.host")}
                    help={t("settings.mediaplayer.host_help")}
                />

                <Accordion
                    type="single"
                    className=""
                    triggerClass="text-[--muted] dark:data-[state=open]:text-white px-0 dark:hover:bg-transparent hover:bg-transparent dark:hover:text-white hover:text-black"
                    itemClass=""
                    contentClass="p-4 border rounded-[--radius-md] bg-[--paper]"
                    collapsible
                    defaultValue={serverStatus?.settings?.mediaPlayer?.defaultPlayer}
                >
                    <AccordionItem value="vlc">
                        <AccordionTrigger>
                            <h4 className="flex gap-2 items-center"><FcVlc /> VLC</h4>
                        </AccordionTrigger>
                        <AccordionContent className="space-y-4">
                            <div className="flex flex-col md:flex-row gap-4">
                                <Field.Text
                                    name="vlcUsername"
                                    label={t("settings.field.username")}
                                />
                                <Field.Text
                                    name="vlcPassword"
                                    label={t("settings.field.password")}
                                    type="password"
                                />
                                <Field.Number
                                    name="vlcPort"
                                    label={t("settings.field.port")}
                                    formatOptions={{
                                        useGrouping: false,
                                    }}
                                    hideControls
                                />
                            </div>
                            <Field.Text
                                name="vlcPath"
                                label={t("settings.mediaplayer.app_path")}
                            />
                        </AccordionContent>
                    </AccordionItem>

                    <AccordionItem value="mpc-hc">
                        <AccordionTrigger>
                            <h4 className="flex gap-2 items-center"><FcClapperboard /> MPC-HC</h4>
                        </AccordionTrigger>
                        <AccordionContent>
                            <div className="flex flex-col md:flex-row gap-4">
                                <Field.Number
                                    name="mpcPort"
                                    label={t("settings.field.port")}
                                    formatOptions={{
                                        useGrouping: false,
                                    }}
                                    hideControls
                                />
                                <Field.Text
                                    name="mpcPath"
                                    label={t("settings.mediaplayer.app_path")}
                                />
                            </div>
                        </AccordionContent>
                    </AccordionItem>

                    <AccordionItem value="mpv">
                        <AccordionTrigger>
                            <h4 className="flex gap-2 items-center"><HiPlay className="mr-1 text-purple-100" /> MPV</h4>
                        </AccordionTrigger>
                        <AccordionContent>
                            <div className="flex gap-4">
                                <Field.Text
                                    name="mpvSocket"
                                    label="Socket"
                                    placeholder={t("settings.mediaplayer.socket_placeholder")}
                                    help={t("settings.mediaplayer.socket_help")}
                                />
                                <Field.Text
                                    name="mpvPath"
                                    label={t("settings.mediaplayer.app_path")}
                                    placeholder={serverStatus?.os === "windows" ? "e.g. C:/Program Files/mpv/mpv.exe" : serverStatus?.os === "darwin"
                                        ? "e.g. /Applications/mpv.app/Contents/MacOS/mpv"
                                        : t("settings.mediaplayer.mpv_path_placeholder_default")}
                                    help={t("settings.mediaplayer.cli_leave_empty")}
                                />
                            </div>
                            <div className="mt-4">
                                <Field.Text
                                    name="mpvArgs"
                                    label={t("player.menu.options")}
                                    placeholder="e.g. --no-config --mute=yes"
                                />
                            </div>
                        </AccordionContent>
                    </AccordionItem>

                    <AccordionItem value="iina">
                        <AccordionTrigger>
                            <h4 className="flex gap-2 items-center"><IoPlayForwardCircleSharp className="mr-1 text-purple-100" /> IINA</h4>
                        </AccordionTrigger>
                        <AccordionContent>
                            <div className="flex gap-4">
                                <Field.Text
                                    name="iinaSocket"
                                    label="Socket"
                                    placeholder={t("settings.mediaplayer.socket_default", { value: getDefaultIinaSocket(serverStatus?.os ?? "") })}
                                />
                                <Field.Text
                                    name="iinaPath"
                                    label={t("settings.mediaplayer.cli_path")}
                                    placeholder={t("settings.mediaplayer.iina_cli_placeholder")}
                                    help={t("settings.mediaplayer.cli_leave_empty")}
                                />
                            </div>
                            <div>
                                <Field.Text
                                    name="iinaArgs"
                                    label={t("player.menu.options")}
                                    placeholder="e.g. --mpv-mute=yes"
                                />
                            </div>
                        </AccordionContent>
                    </AccordionItem>
                </Accordion>
            </SettingsCard>

            <SettingsSubmitButton isPending={isPending} />

        </>
    )
}

export function ExternalPlayerLinkSettings() {

    const { externalPlayerLink, setExternalPlayerLink, encodePath, setEncodePath } = useExternalPlayerLink()

    return (
        <>
            <SettingsPageHeader
                title={t("player.external.title")}
                description={t("player.external.desc")}
                icon={LuCircleArrowOutUpRight}
            />

            <Alert
                intent="info" description={<>
                    {t("player.external.device_only")}
                </>}
            />

            <SettingsCard>
                <div data-settings-external-player-link-scheme>
                    <TextInput
                        label={t("player.external.scheme_label")}
                        placeholder={t("player.external.scheme_placeholder")}
                        value={externalPlayerLink}
                        onValueChange={setExternalPlayerLink}
                    />
                </div>
            </SettingsCard>

            <SettingsCard>
                <Switch
                    side="right"
                    name="encodePath"
                    label={t("settings.mediaplayer.encode_path")}
                    help={t("settings.mediaplayer.encode_path_help")}
                    value={encodePath}
                    onValueChange={setEncodePath}
                />
            </SettingsCard>

            <div className="flex items-center gap-2 text-sm text-gray-500 bg-gray-50 dark:bg-gray-900/30 rounded-lg p-3 border border-gray-200 dark:border-gray-800 border-dashed">
                <RiSettings3Fill className="text-base" />
                <span>{t("settings.common.saved_automatically")}</span>
            </div>
        </>
    )
}
