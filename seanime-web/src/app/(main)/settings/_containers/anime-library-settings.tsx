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
import { t } from "@/lib/i18n"

type LibrarySettingsProps = {
    isPending: boolean
}

export function AnimeLibrarySettings(props: LibrarySettingsProps) {

    const { isPending } = props

    const useLegacyMatching = useWatch({ name: "scannerUseLegacyMatching" })


    return (
        <div className="space-y-8">

            <SettingsCard title={t("settings.card.library_dir_title")}>
                <Field.DirectorySelector
                    name="libraryPath"
                    label={t("settings.field.main_library_path")}
                    leftIcon={<FcFolder />}
                    help={t("settings.anime_library.main_path_help")}
                    shouldExist
                />

                <Field.MultiDirectorySelector
                    name="libraryPaths"
                    label={t("settings.field.extended_library_path")}
                    leftIcon={<FcFolder />}
                    help={t("settings.anime_library.extended_path_help")}
                    shouldExist
                />
            </SettingsCard>

            <SettingsCard title={t("settings.card.library_scan_title")}>

                <Field.Switch
                    side="right"
                    name="autoScan"
                    label={t("settings.field.auto_refresh_library")}
                    moreHelp={<p>
                        {t("settings.anime_library.auto_refresh_more_help")}
                    </p>}
                />

                <Field.Switch
                    side="right"
                    name="refreshLibraryOnStart"
                    label={t("settings.field.scan_on_startup")}
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
                        {t("settings.anime_library.advanced_rules")}
                    </AccordionTrigger>
                    <AccordionContent className="space-y-4">
                        {!useLegacyMatching && <div className="space-y-4">
                            <div>
                                <p className="font-semibold text-lg mb-2">{t("settings.anime_library.scanner_rules_title")}</p>
                                <p className="text-sm text-[--muted] mb-4">
                                    {t("settings.anime_library.scanner_rules_desc")}
                                </p>
                            </div>
                            <ScannerConfigEditor />
                        </div>}

                        <>
                            <Field.Switch
                                name="scannerUseLegacyMatching"
                                label={t("settings.anime_library.use_legacy_matching")}
                                help={t("settings.anime_library.use_legacy_matching_help")}
                                moreHelp={t("settings.anime_library.use_legacy_matching_more_help")}
                            />
                        </>

                        {useLegacyMatching && <div className="flex flex-col md:flex-row gap-3">
                            <Field.Select
                                options={[
                                    { value: "-", label: t("settings.anime_library.algorithm_default") },
                                    { value: "sorensen-dice", label: "Sorensen-Dice" },
                                    { value: "jaccard", label: "Jaccard" },
                                ]}
                                name="scannerMatchingAlgorithm"
                                label={t("settings.anime_library.matching_algorithm")}
                                help={t("settings.anime_library.matching_algorithm_help")}
                            />
                            <Field.Number
                                name="scannerMatchingThreshold"
                                label={t("settings.anime_library.matching_threshold")}
                                placeholder="0.5"
                                help={t("settings.anime_library.matching_threshold_help")}
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

