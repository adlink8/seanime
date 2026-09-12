import { SettingsCard } from "@/app/(main)/settings/_components/settings-card"
import { cn } from "@/components/ui/core/styling"
import { Field } from "@/components/ui/form"
import { t } from "@/lib/i18n"
import React from "react"
import { useFormContext } from "react-hook-form"

type DiscordRichPresenceSettingsProps = {
    children?: React.ReactNode
}

export function DiscordRichPresenceSettings(props: DiscordRichPresenceSettingsProps) {

    const {
        children,
        ...rest
    } = props

    const { watch } = useFormContext()

    const enableRichPresence = watch("enableRichPresence")

    return (
        <>
            <SettingsCard title={t("settings.discord.rich_presence")} description={t("settings.discord.rich_presence_desc")}>
                <div className="space-y-3">
                    <Field.Switch
                        side="right"
                        name="enableRichPresence"
                        label={<span className="flex gap-1 items-center">{t("settings.action.enable")}</span>}
                    />
                    <div
                        className={cn(
                            "flex gap-4 items-center flex-col md:flex-row !mt-3",
                            enableRichPresence ? "opacity-100" : "opacity-50 pointer-events-none",
                        )}
                    >
                        <Field.Checkbox
                            name="enableAnimeRichPresence"
                            label={t("discover.tab.anime")}
                            fieldClass="w-fit"
                        />
                        <Field.Checkbox
                            name="enableMangaRichPresence"
                            label={t("navigation.item.manga")}
                            fieldClass="w-fit"
                        />
                    </div>
                </div>

                <Field.Switch
                    side="right"
                    name="richPresenceHideSeanimeRepositoryButton"
                    label={t("settings.discord.hide_repo_button")}
                />

                {/*<Field.Switch*/}
                {/*    side="right"*/}
                {/*    name="richPresenceShowAniListMediaButton"*/}
                {/*    label="Show AniList Media Button"*/}
                {/*    help="Show a button to open the media page on AniList."*/}
                {/*/>*/}

                <Field.Switch
                    side="right"
                    name="richPresenceShowAniListProfileButton"
                    label={t("settings.discord.show_profile_button")}
                    help={t("settings.discord.show_profile_button_help")}
                />

                {/*<Field.Switch*/}
                {/*    side="right"*/}
                {/*    name="richPresenceUseMediaTitleStatus"*/}
                {/*    label={<span className="flex gap-2 items-center">Use Media Title as Status <LuTriangleAlert className="text-[--orange]" /></span>}*/}
                {/*    moreHelp="Does not work with the default Discord Desktop Client."*/}
                {/*    help="Replace 'Seanime' with the media title in the activity status. Only works if you use a discord client that utilizes arRPC."*/}
                {/*/>*/}
            </SettingsCard>
        </>
    )
}
