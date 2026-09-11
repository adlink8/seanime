import { SettingsCard } from "@/app/(main)/settings/_components/settings-card"
import { SettingsSubmitButton } from "@/app/(main)/settings/_components/settings-submit-button"
import { DataSettings } from "@/app/(main)/settings/_containers/data-settings"
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from "@/components/ui/accordion"
import { Field } from "@/components/ui/form"
import { Separator } from "@/components/ui/separator"
import { javascript } from "@codemirror/lang-javascript"
import { vscodeDark } from "@uiw/codemirror-theme-vscode"
import CodeMirror from "@uiw/react-codemirror"
import React from "react"
import { useFormContext, useWatch } from "react-hook-form"
import { FcFolder } from "react-icons/fc"

type LibrarySettingsProps = {
    isPending: boolean
}

export function AnimeLibrarySettings(props: LibrarySettingsProps) {

    const { isPending } = props

    const useLegacyMatching = useWatch({ name: "scannerUseLegacyMatching" })


    return (
        <div className="space-y-8">

            <SettingsCard title="本地媒体库目录">
                <Field.DirectorySelector
                    name="libraryPath"
                    label="主动漫库路径"
                    leftIcon={<FcFolder />}
                    help="存放动漫视频文件的根目录路径。（请保持路径字母大小写一致）"
                    shouldExist
                />

                <Field.MultiDirectorySelector
                    name="libraryPaths"
                    label="扩展动漫库路径"
                    leftIcon={<FcFolder />}
                    help="如果您的动漫存放在多个不同磁盘或目录下，可以在此处添加其他路径。"
                    shouldExist
                />
            </SettingsCard>

            <SettingsCard title="媒体库扫描">

                <Field.Switch
                    side="right"
                    name="autoScan"
                    label="自动定时刷新媒体库"
                    moreHelp={<p>
                        批量添加新番时，可能需要等待扫描器识别。
                    </p>}
                />

                <Field.Switch
                    side="right"
                    name="refreshLibraryOnStart"
                    label="服务启动时自动扫描刷新"
                />
            </SettingsCard>

            {/*<SettingsCard title="Advanced">*/}

            <Accordion
                type="single"
                collapsible
                className="border rounded-[--radius-md]"
                triggerClass="dark:bg-[--paper]"
                contentClass="!pt-2 dark:bg-[--paper]"
                defaultValue={(useLegacyMatching) ? "more" : undefined}
            >
                <AccordionItem value="more">
                    <AccordionTrigger className="bg-gray-900 rounded-[--radius-md]" data-settings-anime-library="advanced-accordion-trigger">
                        高级扫描规则
                    </AccordionTrigger>
                    <AccordionContent className="space-y-4">
                        {!useLegacyMatching && <div className="space-y-4">
                            <div>
                                <p className="font-semibold text-lg mb-2">扫描器规则配置 (JSON)</p>
                                <p className="text-sm text-[--muted] mb-4">
                                    以 JSON 格式配置高级扫描与番剧识别规则，支持自定义匹配正则和媒体库清洗。
                                </p>
                            </div>
                            <ScannerConfigEditor />
                        </div>}

                        <>
                            <Field.Switch
                                name="scannerUseLegacyMatching"
                                label="Use legacy matching algorithm"
                                help="Enable to use the legacy matching algorithms. (Versions 3.4 and below)"
                                moreHelp="The legacy matching algorithm uses simpler methods which may be less accurate."
                            />
                        </>

                        {useLegacyMatching && <div className="flex flex-col md:flex-row gap-3">
                            <Field.Select
                                options={[
                                    { value: "-", label: "Levenshtein + Sorensen-Dice (Default)" },
                                    { value: "sorensen-dice", label: "Sorensen-Dice" },
                                    { value: "jaccard", label: "Jaccard" },
                                ]}
                                name="scannerMatchingAlgorithm"
                                label="Matching algorithm"
                                help="Choose the algorithm used to match files to AniList entries."
                            />
                            <Field.Number
                                name="scannerMatchingThreshold"
                                label="Matching threshold"
                                placeholder="0.5"
                                help="The minimum score required for a file to be matched to an AniList entry. Default is 0.5."
                                formatOptions={{
                                    minimumFractionDigits: 1,
                                    maximumFractionDigits: 1,
                                }}
                                max={1.0}
                                step={0.1}
                            />
                        </div>}

                        <Separator />

                        <DataSettings />
                    </AccordionContent>
                </AccordionItem>
            </Accordion>

            {/*</SettingsCard>*/}

            <SettingsSubmitButton isPending={isPending} />

        </div>
    )
}

function ScannerConfigEditor() {
    const { setValue } = useFormContext()
    const scannerConfig = useWatch({ name: "scannerConfig" })

    const [value, setLocalValue] = React.useState(scannerConfig || "")

    React.useEffect(() => {
        setLocalValue(scannerConfig || "")
    }, [scannerConfig])

    const handleChange = React.useCallback((val: string) => {
        setLocalValue(val)
        setValue("scannerConfig", val, { shouldDirty: true })
    }, [setValue])

    return (
        <div className="overflow-hidden rounded-[--radius-md] border">
            <CodeMirror
                value={value}
                height="400px"
                theme={vscodeDark}
                extensions={[javascript()]}
                onChange={handleChange}
                basicSetup={{
                    lineNumbers: true,
                    foldGutter: true,
                    bracketMatching: true,
                    syntaxHighlighting: true,
                    highlightActiveLine: true,
                }}
                placeholder={`{
  "matching": {
    "rules": []
  },
  "hydration": {
    "rules": []
  },
  "logs": {
    "verbose": false
  }
}`}
            />
        </div>
    )
}

