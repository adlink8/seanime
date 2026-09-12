import { AL_AnimeDetailsById_Media, Anime_Entry } from "@/api/generated/types"
import {
    useDeleteMediaMetadataParent,
    useGetMediaMetadataParent,
    usePopulateFillerData,
    useRemoveFillerData,
    useSaveMediaMetadataParent,
} from "@/api/hooks/metadata.hooks"
import { Help } from "@/components/shared/help"
import { Button } from "@/components/ui/button"
import { LoadingSpinner } from "@/components/ui/loading-spinner"
import { Modal } from "@/components/ui/modal"
import { NumberInput } from "@/components/ui/number-input"
import { Separator } from "@/components/ui/separator"
import { t } from "@/lib/i18n"
import { atom } from "jotai"
import { useAtom } from "jotai/react"
import React from "react"
import { IoSearchCircle } from "react-icons/io5"

type AnimeEntryMetadataManagerProps = {
    entry: Anime_Entry
    details?: AL_AnimeDetailsById_Media
}

export const __metadataManager_isOpenAtom = atom(false)

export function AnimeEntryMetadataManager(props: AnimeEntryMetadataManagerProps) {

    const { entry, details } = props

    const [isOpen, setOpen] = useAtom(__metadataManager_isOpenAtom)

    const { mutate: filler_populate, isPending: filler_isPopulating } = usePopulateFillerData()
    const { mutate: filler_remove, isPending: filler_isRemoving } = useRemoveFillerData()

    const cannotAddMetadata = entry.media?.format !== "TV" && entry.media?.format !== "TV_SHORT" && entry.media?.format !== "ONA"

    const { data: metadataParentData, isLoading } = useGetMediaMetadataParent(entry.mediaId)
    const { mutate: saveMetadataParent, isPending: isSavingMetadataParent } = useSaveMediaMetadataParent()
    const { mutate: deleteMetadataParent, isPending: isDeletingMetadataParent } = useDeleteMediaMetadataParent()

    const [metadataParentId, setMetadataParentId] = React.useState<number | null>(null)
    const [specialOffset, setSpecialOffset] = React.useState<number | null>(null)

    React.useEffect(() => {
        if (!metadataParentData?.id) return
        setMetadataParentId(metadataParentData.parentId)
        setSpecialOffset(metadataParentData.specialOffset)
    }, [metadataParentData])

    const parentId = details?.relations?.edges?.find(e => e.relationType === "PARENT")?.node?.id

    return (
        <Modal
            open={isOpen}
            onOpenChange={setOpen}
            title={t("library.super_update.metadata")}
            contentClass="max-w-2xl"
            titleClass=""
        >
            <h3 className="text-center">AnimeFillerList</h3>

            <div className="flex lg:flex-row flex-col gap-2">
                <Button
                    className="w-full"
                    intent="primary-subtle"
                    leftIcon={<IoSearchCircle className="text-xl" />}
                    onClick={() => filler_populate({ mediaId: entry.mediaId })}
                    loading={filler_isPopulating}
                    disabled={filler_isPopulating || filler_isRemoving || cannotAddMetadata}
                >
                    {filler_isPopulating ? t("entry.metadata.fetching") : t("entry.metadata.fetch_filler_info")}
                </Button>
                <Button
                    className="w-full"
                    intent="gray-subtle"
                    onClick={() => filler_remove({ mediaId: entry.mediaId })}
                    loading={filler_isRemoving}
                    disabled={filler_isPopulating || filler_isRemoving || cannotAddMetadata}
                >
                    {filler_isRemoving ? t("entry.metadata.removing") : t("entry.metadata.remove_filler_info")}
                </Button>
            </div>

            {!entry.anidbId && (
                <>

                    <Separator />

                    <div>
                        <h3 className="text-center flex gap-2 items-center justify-center">{t("entry.metadata.parent")}
                                                                                           <Help content={t("entry.metadata.parent_help")} />
                        </h3>
                        <p className="text-sm text-[--muted] text-center">
                            {t("entry.metadata.parent_desc")}
                        </p>
                    </div>

                    {isLoading ? <LoadingSpinner /> : (
                        <>
                            <div className="flex gap-2 flex-col lg:flex-row">
                                <NumberInput
                                    leftAddon={t("entry.metadata.anilist_id")}
                                    addonClass="justify-center text-center font-semibold"
                                    hideControls
                                    value={metadataParentId ?? ""}
                                    onValueChange={setMetadataParentId}
                                    formatOptions={{ useGrouping: false }}
                                    rightAddon={(!!parentId && metadataParentId !== parentId) ? <Button
                                        size="xs" intent="gray-link"
                                        onClick={() => setMetadataParentId(parentId)}
                                    >
                                        {t("entry.metadata.select_parent")}
                                    </Button> : undefined}
                                />

                                <NumberInput
                                    leftAddon={t("entry.metadata.special_offset")}
                                    value={specialOffset ?? ""}
                                    onValueChange={setSpecialOffset}
                                    addonClass="text-center font-semibold"
                                    hideControls
                                    placeholder={t("entry.metadata.special_offset_hint")}
                                />
                            </div>

                            <div className="flex gap-2">
                                <Button
                                    className="w-full"
                                    intent="primary-subtle"
                                    loading={isSavingMetadataParent}
                                    onClick={() => {
                                        if (!metadataParentId) return
                                        saveMetadataParent({
                                            mediaId: entry.mediaId,
                                            parentId: metadataParentId,
                                            specialOffset: specialOffset || 0,
                                        })
                                    }}
                                    disabled={isSavingMetadataParent || isDeletingMetadataParent}
                                >
                                    {t("media.action.save")}
                                </Button>
                                {!!metadataParentData?.id && <Button
                                    className="w-full"
                                    intent="gray-subtle"
                                    loading={isDeletingMetadataParent || isSavingMetadataParent}
                                    onClick={() => deleteMetadataParent({ mediaId: entry.mediaId })}
                                >
                                    {t("entry.metadata.remove")}
                                </Button>}
                            </div>
                        </>
                    )}

                </>
            )}

        </Modal>
    )
}
