import { useServerStatus } from "@/app/(main)/_hooks/use-server-status"
import { SettingsCard, SettingsPageHeader } from "@/app/(main)/settings/_components/settings-card"
import { SettingsSubmitButton } from "@/app/(main)/settings/_components/settings-submit-button"
import { Alert } from "@/components/ui/alert"
import { Badge } from "@/components/ui/badge"
import { cn } from "@/components/ui/core/styling"
import { Field } from "@/components/ui/form"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import React, { useState } from "react"
import { useWatch } from "react-hook-form"
import { MdOutlineConnectWithoutContact } from "react-icons/md"
import { t } from "@/lib/i18n"

type Props = {
    isPending: boolean
    children?: React.ReactNode
}

const tabsRootClass = cn("w-full contents space-y-4")

const tabsTriggerClass = cn(
    "h-10 lg:justify-center px-3 flex-1",
)

const tabsListClass = cn(
    "w-full flex flex-row lg:flex-row flex-wrap h-fit !mt-4",
)

const tabContentClass = cn(
    "space-y-4 animate-in fade-in-0 duration-300",
)

export function NakamaSettings(props: Props) {

    const {
        isPending,
        children,
        ...rest
    } = props

    const serverStatus = useServerStatus()
    const nakamaIsHost = useWatch({ name: "nakamaIsHost" })

    const [tab, setTab] = useState("peer")

    React.useLayoutEffect(() => {
        setTab(serverStatus?.settings?.nakama?.isHost ? "host" : "peer")
    }, [serverStatus?.settings?.nakama?.isHost])


    return (
        <div className="space-y-4">

            <SettingsPageHeader
                title={t("settings.nakama.title")}
                description={t("settings.nakama.desc")}
                icon={MdOutlineConnectWithoutContact}
            />

            <SettingsCard>
                <Field.Switch
                    side="right"
                    name="nakamaEnabled"
                    label={t("settings.nakama.enable")}
                />

                <Field.Text
                    label={t("settings.field.username")}
                    name="nakamaUsername"
                    placeholder={t("settings.field.username")}
                    help={t("settings.nakama.username_help")}
                />
            </SettingsCard>

            <Tabs
                value={tab}
                onValueChange={setTab}
                variant="pill"
                className={tabsRootClass}
                triggerClass={tabsTriggerClass}
                listClass={tabsListClass}
            >
                <TabsList>
                    <TabsTrigger value="peer">{t("settings.nakama.tab_peer")}</TabsTrigger>
                    <TabsTrigger value="host">{t("settings.nakama.tab_host")} {serverStatus?.settings?.nakama?.isHost &&
                        <Badge intent="info" className="ml-3">{t("settings.nakama.currently_hosting")}</Badge>}</TabsTrigger>
                    {/*<TabsTrigger value="browser-client">Rendering</TabsTrigger>*/}
                </TabsList>

                <TabsContent value="host" className={tabContentClass}>

                    {!serverStatus?.serverHasPassword &&
                        <Alert
                            intent="warning"
                            title={t("settings.nakama.reminder_title")}
                            description={t("settings.nakama.reminder_desc")}
                        />}

                    <SettingsCard className="!bg-gray-900 text-sm">
                        <div>
                            <p>
                                {t("settings.nakama.host_mode_desc_1")}
                            </p>
                            <p>
                                {t("settings.nakama.host_mode_desc_2_prefix")} <strong>{t("settings.nakama.host_mode_cloud_rooms")}</strong> {t("settings.nakama.host_mode_desc_2_suffix")}
                            </p>
                        </div>
                    </SettingsCard>

                    <SettingsCard>

                        <Field.Switch
                            side="right"
                            name="nakamaIsHost"
                            label={t("settings.nakama.enable_host_mode")}
                            // moreHelp="Password must be set in the config file"
                            help={t("settings.nakama.enable_host_mode_help")}
                        />

                        <Field.Text
                            label={t("settings.nakama.passcode")}
                            name="nakamaHostPassword"
                            placeholder={t("settings.nakama.passcode")}
                            type="password"
                            help={t("settings.nakama.passcode_help")}
                        />

                        {/*<Field.Switch*/}
                        {/*    side="right"*/}
                        {/*    name="nakamaHostEnablePortForwarding"*/}
                        {/*    label="Enable port forwarding"*/}
                        {/*    moreHelp="This might not work for all networks."*/}
                        {/*    help="If enabled, this server will expose its port to the internet. This might be required for other clients to connect to this server."*/}
                        {/*/>*/}
                    </SettingsCard>

                    {nakamaIsHost && <SettingsCard title={t("player.menu.title")}>

                        <Field.Switch
                            side="right"
                            name="nakamaHostShareLocalAnimeLibrary"
                            label={t("settings.nakama.share_local_library")}
                            help={t("settings.nakama.share_local_library_help")}
                        />

                        <Field.MediaExclusionSelector
                            name="nakamaHostUnsharedAnimeIds"
                            label={t("settings.nakama.exclude_from_sharing")}
                            help={t("settings.nakama.exclude_from_sharing_help")}
                        />
                    </SettingsCard>}
                </TabsContent>

                <TabsContent value="peer" className={tabContentClass}>
                    <SettingsCard>
                        {serverStatus?.settings?.nakama?.isHost && <Alert intent="info" description={t("settings.nakama.cannot_connect_while_host")} />}

                        <div
                            className={cn(
                                "space-y-4",
                                serverStatus?.settings?.nakama?.isHost ? "hidden" : "",
                            )}
                        >

                            <Field.Text
                                label={t("settings.nakama.server_url")}
                                name="nakamaRemoteServerURL"
                                placeholder="https://{address} or room://{id}"
                                help={t("settings.nakama.server_url_help")}
                            />

                            <Field.Text
                                label={t("settings.nakama.remote_passcode")}
                                name="nakamaRemoteServerPassword"
                                placeholder={t("settings.nakama.passcode")}
                                help={t("settings.nakama.remote_passcode_help")}
                                type="password"
                            />
                        </div>
                    </SettingsCard>

                    {!serverStatus?.settings?.nakama?.isHost && <SettingsCard title={t("player.menu.title")}>
                        <Field.Switch
                            side="right"
                            name="includeNakamaAnimeLibrary"
                            label={t("settings.nakama.use_nakama_library")}
                            help={t("settings.nakama.use_nakama_library_help")}
                        />
                    </SettingsCard>}
                </TabsContent>

            </Tabs>

            <SettingsSubmitButton isPending={isPending} />

        </div>
    )
}
