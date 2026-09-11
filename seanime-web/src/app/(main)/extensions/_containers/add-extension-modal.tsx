import { Extension_Extension } from "@/api/generated/types"
import { useFetchExternalExtensionData, useInstallExternalExtension, useInstallExternalExtensionRepository } from "@/api/hooks/extensions.hooks"
import { ExtensionDetails } from "@/app/(main)/extensions/_components/extension-details"
import { MarketplaceExtensionCard } from "@/app/(main)/extensions/_containers/marketplace-extensions"
import { Button } from "@/components/ui/button"
import { Modal } from "@/components/ui/modal"
import { Separator } from "@/components/ui/separator"
import { TextInput } from "@/components/ui/text-input"
import React from "react"
import { FiDownload } from "react-icons/fi"
import { LuSearch } from "react-icons/lu"
import { toast } from "sonner"

type AddExtensionModalProps = {
    extensions: Extension_Extension[] | undefined
    children?: React.ReactElement
}

export function AddExtensionModal(props: AddExtensionModalProps) {

    const {
        extensions,
        children,
        ...rest
    } = props

    const [open, setOpen] = React.useState(false)
    const [manifestURL, setManifestURL] = React.useState<string>("")
    const [repositoryURL, setRepositoryURL] = React.useState<string>("")

    const { mutate: fetchExtensionData, data: extensionData, isPending, reset } = useFetchExternalExtensionData(null)
    const {
        mutate: installFromRepository,
        data: repositoryData,
        isPending: isInstallingFromRepo,
        reset: resetRepo,
    } = useInstallExternalExtensionRepository()

    const {
        mutate: installExtension,
        data: installResponse,
        isPending: isInstalling,
    } = useInstallExternalExtension()

    React.useEffect(() => {
        if (installResponse) {
            toast.success(installResponse.message)
            setOpen(false)
            reset()
        }
    }, [installResponse])

    function handleFetchExtensionData() {
        if (!manifestURL) {
            toast.warning("Please provide a valid URL.")
            return
        }

        fetchExtensionData({
            manifestUri: manifestURL,
        })
    }

    function handleInstallFromRepository(install: boolean) {
        if (!repositoryURL) {
            toast.warning("Please provide a valid URL.")
            return
        }

        installFromRepository({
            repositoryUri: repositoryURL,
            install: install,
        }, {
            onSuccess: () => {
                if (install) {
                    toast.success("Extensions installed successfully.")
                    setOpen(false)
                    setRepositoryURL("")
                    resetRepo()
                }
            },
        })
    }

    return (
        <>
            <Modal
                open={open}
                onOpenChange={setOpen}
                trigger={children}
                contentClass="max-w-3xl"
                titleClass="text-center pb-4"
                title="添加扩展"
            >
                <div className="flex gap-4 flex-col lg:flex-row">
                    <div className="lg:w-1/3">
                        <h3 className="text-2xl font-bold">通过 URL 安装</h3>
                        <p className="text-[--muted]">输入扩展 manifest.json 链接进行安装。</p>
                    </div>
                    <div className="lg:w-2/3 gap-3 flex flex-col">
                        <TextInput
                            placeholder="https://example.com/extension.json"
                            value={manifestURL}
                            onValueChange={setManifestURL}
                            // label="URL"
                        />
                        <Button
                            leftIcon={<LuSearch />}
                            intent="white"
                            onClick={handleFetchExtensionData}
                            loading={isPending}
                        >查找</Button>
                    </div>
                </div>

                {!!extensionData && (
                    <>
                        <Separator />

                        <ExtensionDetails extension={extensionData} />

                        {extensions?.find(n => n.id === extensionData.id) ? (
                            <p className="text-center">
                                该扩展已经安装过。
                            </p>
                        ) : (
                            <Button
                                intent="white"
                                loading={isInstalling}
                                onClick={() => {
                                    installExtension({
                                        manifestUri: extensionData?.manifestURI,
                                    })
                                }}
                            >安装</Button>
                        )}
                    </>
                )}

                {!extensionData && (
                    <>
                        <Separator />

                        <p className="text-center text-[--muted]">
                            您也可以输入扩展源仓库地址，批量导入并一键安装多个扩展。
                        </p>

                        <div className="flex gap-4 flex-col lg:flex-row-reverse">
                            <div className="lg:w-1/3">
                                <h3 className="text-xl font-bold">从软件源批量导入</h3>
                                <p className="text-[--muted]">输入扩展仓库地址或 JSON 链接导入扩展。</p>
                            </div>
                            <div className="lg:w-2/3 gap-3 flex flex-col">
                                <TextInput
                                    placeholder={"https://example.com/extensions.json 或 { \"urls\": [...] }"}
                                    value={repositoryURL}
                                    onValueChange={setRepositoryURL}
                                    // label="URL"
                                />
                                <Button
                                    leftIcon={<FiDownload />}
                                    intent="gray-outline"
                                    onClick={() => handleInstallFromRepository(false)}
                                    loading={isInstallingFromRepo}
                                >导入全部</Button>
                            </div>
                        </div>

                        {!!repositoryData && (
                            <>
                                {repositoryData.extensions?.toSorted((a, b) => a.id.toLowerCase().localeCompare(b.id.toLowerCase())).map(ext => (
                                    <div key={ext.id} className="">
                                        <MarketplaceExtensionCard extension={ext} isInstalled={false} hideInstallButton showType />
                                    </div>
                                ))}
                                <Button
                                    leftIcon={<FiDownload />}
                                    intent="white"
                                    onClick={() => handleInstallFromRepository(true)}
                                    loading={isInstallingFromRepo}
                                >全部安装</Button>
                            </>
                        )}
                    </>
                )}

            </Modal>
        </>
    )
}
