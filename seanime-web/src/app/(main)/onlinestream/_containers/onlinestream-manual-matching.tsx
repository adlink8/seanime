import { Anime_Entry } from "@/api/generated/types"
import {
    useGetOnlinestreamMapping,
    useOnlinestreamManualMapping,
    useOnlinestreamManualSearch,
    useRemoveOnlinestreamMapping,
} from "@/api/hooks/onlinestream.hooks"
import { ConfirmationDialog, useConfirmationDialog } from "@/components/shared/confirmation-dialog"
import { SeaLink } from "@/components/shared/sea-link"
import { AppLayoutStack } from "@/components/ui/app-layout"
import { Button, IconButton } from "@/components/ui/button"
import { cn } from "@/components/ui/core/styling"
import { defineSchema, Field, Form, InferType } from "@/components/ui/form"
import { LoadingSpinner } from "@/components/ui/loading-spinner"
import { Modal } from "@/components/ui/modal"
import { Separator } from "@/components/ui/separator"
import { Tooltip } from "@/components/ui/tooltip"
import { t } from "@/lib/i18n"
import React from "react"
import { BiLinkExternal } from "react-icons/bi"
import { FiSearch } from "react-icons/fi"

type OnlinestreamManualMappingModalProps = {
    entry: Anime_Entry
    provider: string
    children: React.ReactElement
}

export function OnlinestreamManualMappingModal(props: OnlinestreamManualMappingModalProps) {

    const {
        children,
        entry,
        provider,
        ...rest
    } = props

    return (
        <>
            <Modal
                title={t("misc.onlinestream.manual_match")}
                description={t("misc.onlinestream.manual_match_desc")}
                trigger={children}
                contentClass="max-w-4xl"
            >
                <Content entry={entry} provider={provider} />
            </Modal>
        </>
    )
}

const searchSchema = defineSchema(({ z }) => z.object({
    query: z.string().min(1),
    dubbed: z.boolean().default(false),
}))

function Content({ entry, provider }: { entry: Anime_Entry, provider: string }) {
    const selectedProvider = provider

    // Get current mapping
    const { data: existingMapping, isLoading: mappingLoading } = useGetOnlinestreamMapping({
        provider: selectedProvider || undefined,
        mediaId: entry.mediaId,
    })

    // Search
    const { mutate: search, data: searchResults, isPending: searchLoading, reset } = useOnlinestreamManualSearch(entry.mediaId, selectedProvider)

    function handleSearch(data: InferType<typeof searchSchema>) {
        if (selectedProvider) {
            search({
                provider: selectedProvider,
                query: data.query,
                dubbed: data.dubbed,
            })
        }
    }

    // Match
    const { mutate: match, isPending: isMatching } = useOnlinestreamManualMapping()

    // Unmatch
    const { mutate: unmatch, isPending: isUnmatching } = useRemoveOnlinestreamMapping()

    const [animeId, setAnimeId] = React.useState<string | null>(null)
    const confirmMatch = useConfirmationDialog({
        title: t("misc.onlinestream.manual_match"),
        description: t("misc.onlinestream.confirm_match_desc"),
        actionText: t("media.action.confirm"),
        actionIntent: "success",
        onConfirm: () => {
            if (animeId && selectedProvider) {
                match({
                    provider: selectedProvider,
                    mediaId: entry.mediaId,
                    animeId: animeId,
                })
                reset()
                setAnimeId(null)
            }
        },
    })

    return (
        <>
            {mappingLoading ? (
                <LoadingSpinner />
            ) : (
                <AppLayoutStack>
                    <div className="text-center">
                        {!!existingMapping?.animeId ? (
                            <AppLayoutStack>
                                <p>
                                    {t("misc.onlinestream.current_mapping")} <span>{existingMapping.animeId}</span>
                                </p>
                                <Button
                                    intent="alert-subtle" loading={isUnmatching} onClick={() => {
                                    if (selectedProvider) {
                                        unmatch({
                                            provider: selectedProvider,
                                            mediaId: entry.mediaId,
                                        })
                                    }
                                }}
                                >
                                    {t("misc.onlinestream.remove_mapping")}
                                </Button>
                            </AppLayoutStack>
                        ) : (
                            <p className="text-[--muted] italic">{t("misc.onlinestream.no_manual_match")}</p>
                        )}
                    </div>

                    <Separator />

                    <Form schema={searchSchema} onSubmit={handleSearch}>
                        <div className="space-y-2">
                            <Field.Text
                                name="query"
                                placeholder={t("misc.onlinestream.search_placeholder")}
                                leftIcon={<FiSearch className="text-xl text-[--muted]" />}
                                fieldClass="w-full"
                            />

                            <Field.Switch
                                name="dubbed"
                                label={t("misc.onlinestream.look_for_dubs")}
                                side="right"
                                moreHelp={t("misc.onlinestream.dubs_help")}
                            />

                            <Field.Submit intent="white" loading={isMatching || searchLoading || mappingLoading} className="">{t("navigation.item.search")}</Field.Submit>
                        </div>
                    </Form>

                    {searchLoading ? <LoadingSpinner /> : (
                        <>
                            <div className="space-y-2">
                                {searchResults?.map(item => (
                                    <div
                                        key={item.id}
                                        className={cn(
                                            "flex justify-between items-center",
                                        )}
                                    >
                                        <p
                                            onClick={() => {
                                                setAnimeId(item.id)
                                                React.startTransition(() => {
                                                    confirmMatch.open()
                                                })
                                            }}
                                            className="cursor-pointer hover:underline"
                                        >
                                            {item.title}
                                        </p>
                                        <div>
                                            <SeaLink href={item.url} target="_blank">
                                                <Tooltip
                                                    trigger={<IconButton
                                                        icon={<BiLinkExternal />}
                                                        intent="primary-basic"
                                                        size="xs"
                                                    />}
                                                >
                                                    {t("torrent.qbittorrent.open_in_browser")}
                                                </Tooltip>
                                            </SeaLink>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </>
                    )}

                </AppLayoutStack>
            )}

            <ConfirmationDialog {...confirmMatch} />
        </>
    )
}
