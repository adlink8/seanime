import { useGetAutoDownloaderProfiles } from "@/api/hooks/auto_downloader.hooks"
import { useAnimeListTorrentProviderExtensions } from "@/api/hooks/extensions.hooks"
import { Button, CloseButton, IconButton } from "@/components/ui/button"
import { Field } from "@/components/ui/form"
import { TextInput } from "@/components/ui/text-input"
import { t } from "@/lib/i18n"
import React from "react"
import { useFieldArray } from "react-hook-form"
import { BiPlus } from "react-icons/bi"

type TextArrayFieldProps<T extends string | number> = {
    name: string
    control: any
    type?: "text" | "number"
    label?: string
    placeholder?: string
    separatorText?: string
    suggestions?: string[]
    suggestionLabels?: string[]
}

export function TextArrayField<T extends string | number>(props: TextArrayFieldProps<T>) {
    const { fields, append, remove } = useFieldArray({
        control: props.control,
        name: props.name,
    })

    return (
        <div className="space-y-2">
            {props.label && <div className="flex items-center">
                <div className="text-base font-semibold">{props.label}</div>
            </div>}
            {!!props.suggestions?.length && (
                <div className="flex flex-wrap gap-2 mb-2">
                    {props.suggestions.map((suggestion, index) => (
                        <Button
                            key={suggestion}
                            intent="gray-subtle"
                            size="sm"
                            className="rounded-full"
                            onClick={() => append(props.type === "number" ? parseInt(suggestion) : suggestion)}
                        >
                            {props.suggestionLabels?.[index] || suggestion}
                        </Button>
                    ))}
                </div>
            )}
            {fields.map((field, index) => (
                <React.Fragment key={field.id}>
                    <div className="flex gap-2 items-center">
                        {props.type === "text" && (
                            <TextInput
                                {...props.control.register(`${props.name}.${index}`)}
                                placeholder={props.placeholder}
                            />
                        )}
                        {props.type === "number" && (
                            <TextInput
                                type="number"
                                {...props.control.register(`${props.name}.${index}`, {
                                    valueAsNumber: true,
                                    min: 1,
                                    validate: (value: number) => !isNaN(value),
                                })}
                            />
                        )}
                        <CloseButton
                            size="sm"
                            intent="alert-subtle"
                            onClick={() => remove(index)}
                        />
                    </div>
                    {(!!props.separatorText && index < fields.length - 1) && (
                        <p className="text-center text-[--muted]">{props.separatorText}</p>
                    )}
                </React.Fragment>
            ))}
            <IconButton
                intent="success-subtle"
                className="rounded-full"
                onClick={() => append(props.type === "number" ? 1 : "")}
                icon={<BiPlus />}
            />
        </div>
    )
}

type ReleaseGroupsFieldProps = {
    name: string
    control: any
}

export function ReleaseGroupsField(props: ReleaseGroupsFieldProps) {
    const suggestions = [
        "SubsPlease",
        "Erai-raws",
        "VARYG",
        "EMBER",
        "Judas",
        "ASW",
        "Tsundere-Raws",
    ]

    return (
        <div className="border rounded-[--radius] p-4 relative !mt-8 space-y-3">
            <div className="absolute -top-2.5 tracking-wide font-semibold uppercase text-sm left-4 bg-gray-950 px-2">{t("autodownloader.section.release_groups")}</div>
            <p className="text-sm">
                {t("autodownloader.help.release_groups_short")}
            </p>

            <TextArrayField
                name={props.name}
                control={props.control}
                type="text"
                placeholder={t("autodownloader.placeholder.release_group")}
                separatorText={t("settings.common.or")}
                suggestions={suggestions}
            />
        </div>
    )
}

type ResolutionsFieldProps = {
    name: string
    control: any
}

export function ResolutionsField(props: ResolutionsFieldProps) {
    const suggestions = ["1080p", "720p", "2160p", "480p"]

    return (
        <div className="border rounded-[--radius] p-4 relative !mt-8 space-y-3">
            <div className="absolute -top-2.5 tracking-wide font-semibold uppercase text-sm left-4 bg-gray-950 px-2">{t("entry.torrent_search.resolution")}</div>
            <p className="text-sm">
                {t("autodownloader.help.resolutions")}
            </p>

            <TextArrayField
                name={props.name}
                control={props.control}
                type="text"
                placeholder={t("autodownloader.placeholder.resolution")}
                separatorText={t("settings.common.or")}
                suggestions={suggestions}
            />
        </div>
    )
}

type AdditionalTermsFieldProps = {
    name: string
    control: any
    defaultOpen?: boolean
}

export function AdditionalTermsField(props: AdditionalTermsFieldProps) {
    const suggestions = [
        // Video
        { label: t("autodownloader.suggestion.hevc"), value: "H265,H.265,x265,HEVC" },
        { label: t("autodownloader.suggestion.avc"), value: "H264,H.264,x264,AVC" },
        { label: "10-bit", value: "10bit,10-bit,10 bit" },
        { label: "HDR", value: "HDR,HDR10,HDR10+" },
        { label: t("autodownloader.suggestion.dolby_vision"), value: "Dolby Vision,DolbyVision,DoVi" },
        { label: t("autodownloader.suggestion.remux"), value: "Remux" },
        // Audio
        { label: "FLAC", value: "FLAC" },
        { label: t("autodownloader.suggestion.opus"), value: "Opus" },
        { label: "AAC", value: "AAC,AAC2.0" },
        { label: "E-AC3", value: "E-AC3,EAC3,EAC-3,E-AC-3" },
        { label: "TrueHD", value: "TrueHD,True-HD" },
        { label: "DTS", value: "DTS,DTS-HD,DTS-X" },
        // Source
        { label: "BluRay", value: "BluRay,Blu-Ray,BDRip" },
        { label: "WEB-DL", value: "WEB-DL,WEBDL,WEB DL,WebRip,Web-Rip,WEB RIP" },
        { label: "DVD", value: "DVD,DVD-Rip,DVDRip" },
        // Anime
        { label: t("autodownloader.suggestion.dual_audio"), value: "Dual Audio,Dual-Audio,DualAudio" },
        { label: t("settings.autoselect.pref_multiple_audio"), value: "Multi Audio,Multi-Audio,MultiAudio,Dual Audio,Dual-Audio,DualAudio" },
        { label: "Multi-Sub", value: "Multi-Sub,Multi Sub,Multisub,Multisubs,Multi-subs,Multi subs" },
        { label: t("autodownloader.suggestion.dubbed"), value: "Dubbed,Dub" },
    ]

    return (
        <div className="border rounded-[--radius] p-4 relative !mt-8 space-y-3">
            <div className="absolute -top-2.5 tracking-wide font-semibold uppercase text-sm left-4 bg-gray-950 px-2">{t("autodownloader.section.terms")}</div>
            <div>
                <p className="text-sm -top-2 relative"><span className="font-semibold">
                    {t("autodownloader.help.additional_terms")}</span> {t("autodownloader.help.additional_terms_more")}</p>
            </div>

            <TextArrayField
                name={props.name}
                control={props.control}
                type="text"
                placeholder={t("autodownloader.placeholder.terms")}
                separatorText={t("autodownloader.separator.and")}
                suggestions={suggestions.map(s => s.value)}
                suggestionLabels={suggestions.map(s => s.label)}
            />
        </div>
    )
}

export function ExcludeTermsField(props: AdditionalTermsFieldProps) {
    const suggestions = [
        // Video
        { label: t("autodownloader.suggestion.hevc"), value: "H265,H.265,x265,HEVC" },
        { label: t("autodownloader.suggestion.avc"), value: "H264,H.264,x264,AVC" },
        { label: "10-bit", value: "10bit,10-bit,10 bit" },
        { label: "HDR", value: "HDR,HDR10,HDR10+" },
        { label: t("autodownloader.suggestion.dolby_vision"), value: "Dolby Vision,DolbyVision,DoVi" },
        { label: t("autodownloader.suggestion.remux"), value: "Remux" },
        // Audio
        { label: "FLAC", value: "FLAC" },
        // Source
        { label: "BluRay", value: "BluRay,Blu-Ray,BDRip" },
        { label: "DVD", value: "DVD,DVD-Rip,DVDRip" },
        { label: t("autodownloader.suggestion.cam"), value: "Cam,CamRip" },
        // Anime
        { label: t("autodownloader.suggestion.dubbed"), value: "Dubbed,Dub" },
    ]

    return (
        <div className="border rounded-[--radius] p-4 relative !mt-8 space-y-3">
            <div className="absolute -top-2.5 tracking-wide font-semibold uppercase text-sm left-4 bg-gray-950 px-2">{t("settings.autoselect.exclude_terms")}</div>
            <p className="text-sm"><span className="font-semibold">
                {t("autodownloader.help.exclude_terms")}</span> {t("autodownloader.help.exclude_terms_more")}</p>

            <TextArrayField
                name={props.name}
                control={props.control}
                type="text"
                placeholder={t("autodownloader.placeholder.terms")}
                separatorText={t("autodownloader.separator.and")}
                suggestions={suggestions.map(s => s.value)}
                suggestionLabels={suggestions.map(s => s.label)}
            />
        </div>
    )
}

type ProvidersFieldProps = {
    name: string
    control: any
}

export function ProvidersField(props: ProvidersFieldProps) {
    const { data: extensions } = useAnimeListTorrentProviderExtensions()

    return (
        <div className="border rounded-[--radius] p-4 relative !mt-8 space-y-3">
            <div className="absolute -top-2.5 tracking-wide font-semibold uppercase text-sm left-4 bg-gray-950 px-2">{t("settings.debrid.provider")}</div>
            <p className="text-sm">
                {t("autodownloader.help.providers")}
            </p>
            <Field.Combobox
                name={props.name}
                options={extensions?.map(ext => ({
                    label: ext.name,
                    textValue: ext.name,
                    value: ext.id,
                })) ?? []}
                multiple
                label={t("autodownloader.field.select_providers")}
                emptyMessage={t("settings.autoselect.empty_providers")}
            />
        </div>
    )
}

type ProfileSelectFieldProps = {
    name: string
}

export function ProfileSelectField(props: ProfileSelectFieldProps) {
    const { data: profiles } = useGetAutoDownloaderProfiles()

    return (
        <div className="border rounded-[--radius] p-4 relative !mt-8 space-y-3">
            <div className="absolute -top-2.5 tracking-wide font-semibold uppercase text-sm left-4 bg-gray-950 px-2">{t("autodownloader.section.profile")}</div>
            <p className="text-sm">
                {t("autodownloader.help.profile")}
            </p>
            <Field.Combobox
                name={props.name}
                options={[
                    ...(profiles?.map(profile => ({
                        label: profile.name,
                        textValue: profile.name,
                        value: String(profile.dbId),
                    })) ?? []),
                ]}
                label={t("autodownloader.field.select_profile")}
                emptyMessage={t("autodownloader.empty.no_profile")}
            />
        </div>
    )
}
