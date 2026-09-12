import { getServerBaseUrl } from "@/api/client/server-url"
import { useImportLocalFiles } from "@/api/hooks/localfiles.hooks"
import { useServerHMACAuth } from "@/app/(main)/_hooks/use-server-status"
import { Button } from "@/components/ui/button"
import { Modal } from "@/components/ui/modal"
import { TextInput } from "@/components/ui/text-input"
import { openTab } from "@/lib/helpers/browser"
import { t } from "@/lib/i18n"
import React from "react"
import { CgImport } from "react-icons/cg"
import { TbDatabaseExport } from "react-icons/tb"
import { toast } from "sonner"

type DataSettingsProps = {
    children?: React.ReactNode
}

export function DataSettings(props: DataSettingsProps) {

    const {
        children,
        ...rest
    } = props

    const { mutate: importLocalFiles, isPending: isImportingLocalFiles } = useImportLocalFiles()
    const [localFileDataPath, setLocalFileDataPath] = React.useState("")

    function handleImportLocalFiles() {
        if (!localFileDataPath) return

        importLocalFiles({ dataFilePath: localFileDataPath }, {
            onSuccess: () => {
                setLocalFileDataPath("")
            },
        })
    }

    const { getHMACTokenQueryParam } = useServerHMACAuth()

    const handleExportLocalFiles = React.useCallback(async () => {
        try {
            const endpoint = "/api/v1/library/local-files/dump"
            const tokenQuery = await getHMACTokenQueryParam(endpoint)
            openTab(`${getServerBaseUrl()}${endpoint}${tokenQuery}`)
        }
        catch (error) {
            toast.error(t("settings.data.toast_export_failed"))
        }
    }, [getHMACTokenQueryParam])

    return (
        <div className="space-y-4">

            <div>
                <h5>{t("settings.data.title")}</h5>

                <p className="text-[--muted]">
                    {t("settings.data.desc")}
                </p>
            </div>

            <div className="flex flex-wrap gap-2">
                <Button
                    intent="primary-subtle"
                    leftIcon={<TbDatabaseExport className="text-xl" />}
                    size="md"
                    disabled={isImportingLocalFiles}
                    onClick={handleExportLocalFiles}
                >
                    {t("settings.data.export")}
                </Button>

                <Modal
                    title={t("settings.data.import_title")}
                    trigger={
                        <Button
                            intent="white-subtle"
                            leftIcon={<CgImport className="text-xl" />}
                            size="md"
                            disabled={isImportingLocalFiles}
                        >
                            {t("settings.data.import_title")}
                        </Button>
                    }
                >

                    <p>
                        {t("settings.data.import_warning")}
                    </p>

                    <TextInput
                        label={t("settings.data.file_path")}
                        help={t("settings.data.file_path_help")}
                        value={localFileDataPath}
                        onValueChange={setLocalFileDataPath}
                    />

                    <Button
                        intent="white"
                        rounded
                        onClick={handleImportLocalFiles}
                        disabled={isImportingLocalFiles}
                    >{t("settings.action.import")}</Button>

                </Modal>
            </div>
        </div>
    )
}
