import { useListMangaProviderExtensions } from "@/api/hooks/extensions.hooks"
import { useServerStatus } from "@/app/(main)/_hooks/use-server-status"
import { useStoredMangaProviders } from "@/app/(main)/manga/_lib/handle-manga-selected-provider"
import { SettingsCard, SettingsPageHeader } from "@/app/(main)/settings/_components/settings-card"
import { SettingsSubmitButton } from "@/app/(main)/settings/_components/settings-submit-button"
import { ConfirmationDialog, useConfirmationDialog } from "@/components/shared/confirmation-dialog"
import { Button } from "@/components/ui/button"
import { Field } from "@/components/ui/form"
import { atom } from "jotai"
import { useAtom } from "jotai/react"
import React from "react"
import { useFormContext } from "react-hook-form"
import { LuBookOpen } from "react-icons/lu"
import { toast } from "sonner"
import { t } from "@/lib/i18n"

type MangaSettingsProps = {
    isPending: boolean
}

const __manga_storedProvidersHistoryAtom = atom<Record<string, string> | null>(null)

export function MangaSettings(props: MangaSettingsProps) {

    const {
        isPending,
        ...rest
    } = props

    const serverStatus = useServerStatus()
    const f = useFormContext()

    const { data: extensions } = useListMangaProviderExtensions()

    const { storedProviders, overwriteStoredProviders, overwriteStoredProvidersWith } = useStoredMangaProviders(extensions)
    const [storedProvidersHistory, setStoredProvidersHistory] = useAtom(__manga_storedProvidersHistoryAtom)

    const options = React.useMemo(() => {
        return [
            { label: t("settings.option.automatic"), value: "-" },
            ...(extensions?.map(provider => ({
                label: provider.name,
                value: provider.id,
            })) ?? []).sort((a, b) => a.label.localeCompare(b.label)),
        ]
    }, [extensions])

    const defaultProviderExt = extensions?.find(e => e.id === serverStatus?.settings?.manga?.defaultMangaProvider)

    const confirmDialog = useConfirmationDialog({
        title: t("settings.manga.overwrite_confirm_title"),
        description: t("settings.manga.overwrite_confirm_desc"),
        actionText: t("settings.manga.overwrite_action"),
        actionIntent: "warning",
        onConfirm: async () => {
            if (!defaultProviderExt) return
            const oldProviders = structuredClone(storedProviders)
            overwriteStoredProvidersWith(defaultProviderExt.id)
            toast.success(t("settings.manga.overwritten_toast"))
            setTimeout(() => {
                setStoredProvidersHistory(oldProviders)
            }, 500)
        },
    })

    return (
        <>
            <SettingsPageHeader
                title={t("settings.manga.title")}
                description={t("settings.manga.desc")}
                icon={LuBookOpen}
            />

            <SettingsCard>
                <Field.Switch
                    side="right"
                    name="enableManga"
                    label={<span className="flex gap-1 items-center">{t("settings.action.enable")}</span>}
                    help={t("settings.manga.enable_help")}
                />
                <Field.Switch
                    side="right"
                    name="mangaAutoUpdateProgress"
                    label={t("settings.manga.auto_update_progress")}
                    help={t("settings.manga.auto_update_progress_help")}
                />
            </SettingsCard>

            <SettingsCard title={t("settings.manga.provider_title")}>
                <Field.Select
                    name="defaultMangaProvider"
                    label={t("settings.manga.default_provider")}
                    help={t("settings.manga.default_provider_help")}
                    options={options}
                />
                {(!!defaultProviderExt && f.watch("defaultMangaProvider") === serverStatus?.settings?.manga?.defaultMangaProvider) && (
                    <div className="flex w-full space-x-4 flex-wrap">
                        <Button className="px-0 py-1" intent="warning-link" onClick={() => confirmDialog.open()}>
                            {t("settings.manga.overwrite_with", { name: defaultProviderExt.name })}
                        </Button>
                        {!!storedProvidersHistory && (
                            <Button
                                className="px-0 py-1" intent="gray-link" onClick={() => {
                                overwriteStoredProviders(storedProvidersHistory)
                                toast.success(t("settings.manga.restored_toast"))
                                setStoredProvidersHistory(null)
                            }}
                            >
                                {t("settings.manga.undo")}
                            </Button>
                        )}
                    </div>
                )}
            </SettingsCard>

            <SettingsCard title={t("settings.manga.local_provider_title")} description={t("settings.manga.local_provider_desc")}>

                <Field.DirectorySelector
                    name="mangaLocalSourceDirectory"
                    label={t("settings.manga.local_source_directory")}
                    help={t("settings.manga.local_source_directory_help")}
                />
            </SettingsCard>

            <ConfirmationDialog {...confirmDialog} />

            <SettingsSubmitButton isPending={isPending} />
        </>
    )
}
