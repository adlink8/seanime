import { useLocalAddTrackedMedia, useLocalGetIsMediaTracked, useLocalRemoveTrackedMedia } from "@/api/hooks/local.hooks"
import { ConfirmationDialog, useConfirmationDialog } from "@/components/shared/confirmation-dialog"
import { IconButton } from "@/components/ui/button"
import { Tooltip } from "@/components/ui/tooltip"
import { t } from "@/lib/i18n"
import React from "react"
import { MdOutlineDownloadForOffline, MdOutlineOfflinePin } from "react-icons/md"

type MediaSyncTrackButtonProps = {
    mediaId: number
    type: "anime" | "manga"
    size?: "sm" | "md" | "lg"
}

export function MediaSyncTrackButton(props: MediaSyncTrackButtonProps) {

    const {
        mediaId,
        type,
        size,
        ...rest
    } = props

    const { data: isTracked, isLoading } = useLocalGetIsMediaTracked(mediaId, type)
    const { mutate: addMedia, isPending: isAdding } = useLocalAddTrackedMedia()
    const { mutate: removeMedia, isPending: isRemoving } = useLocalRemoveTrackedMedia()

    function handleToggle() {
        if (isTracked) {
            removeMedia({ mediaId, type })
        } else {
            addMedia({
                media: [{
                    mediaId: mediaId,
                    type: type,
                }],
            })
        }
    }

    const confirmUntrack = useConfirmationDialog({
        title: t("media.action.remove_offline_data"),
        description: t("media.confirm.remove_offline_data_description"),
        onConfirm: () => {
            handleToggle()
        },
    })

    return (
        <>
            <Tooltip
                trigger={<IconButton
                    icon={isTracked ? <MdOutlineOfflinePin /> : <MdOutlineDownloadForOffline />}
                    onClick={() => isTracked ? confirmUntrack.open() : handleToggle()}
                    loading={isLoading || isAdding || isRemoving}
                    intent={isTracked ? "primary-subtle" : "gray-subtle"}
                    size={size}
                    {...rest}
                />}
            >
                {isTracked ? t("media.action.remove_offline_data") : t("media.action.save_locally")}
            </Tooltip>

            <ConfirmationDialog {...confirmUntrack} />
        </>
    )
}
