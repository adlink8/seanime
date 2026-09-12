import { useLocalSyncSimulatedDataToAnilist } from "@/api/hooks/local.hooks"
import { SettingsPageHeader } from "@/app/(main)/settings/_components/settings-card"
import { SettingsSubmitButton } from "@/app/(main)/settings/_components/settings-submit-button"
import { ConfirmationDialog, useConfirmationDialog } from "@/components/shared/confirmation-dialog"
import { t } from "@/lib/i18n"
import React from "react"
import { SiAnilist } from "react-icons/si"

type Props = {
    isPending: boolean
    children?: React.ReactNode
}

export function AnilistSettings(props: Props) {

    const {
        isPending,
        children,
        ...rest
    } = props

    const { mutate: upload, isPending: isUploading } = useLocalSyncSimulatedDataToAnilist()

    const confirmDialog = useConfirmationDialog({
        title: t("settings.confirm.upload_to_anilist"),
        description: t("settings.confirm.upload_to_anilist_desc"),
        actionText: t("settings.action.upload"),
        actionIntent: "primary",
        onConfirm: async () => {
            if (isUploading) return
            upload()
        },
    })

    return (
        <div className="space-y-4">

            <SettingsPageHeader
                title="AniList"
                description={t("settings.anilist.desc")}
                icon={SiAnilist}
            />


            <SettingsSubmitButton isPending={isPending} />

            <ConfirmationDialog {...confirmDialog} />

        </div>
    )
}
