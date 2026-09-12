import { usePatchSetting } from "@/api/hooks/settings.hooks"
import { MediaCoreControlButtonIcon } from "@/app/(main)/_features/media-core/media-core-control-bar"
import {
    MediaCoreMenu,
    MediaCoreMenuOption,
    MediaCoreMenuSectionBody,
    MediaCoreMenuSubmenuBody,
    MediaCoreMenuSubOption,
    MediaCoreMenuSubSubmenuBody,
    MediaCoreMenuTitle,
    MediaCoreSettingSelect,
    MediaCoreSettingTextInput,
} from "@/app/(main)/_features/media-core/media-core-menu"
import { useServerStatus } from "@/app/(main)/_hooks/use-server-status"
import { Button } from "@/components/ui/button"
import { Switch } from "@/components/ui/switch"
import React from "react"
import { HiFastForward } from "react-icons/hi"
import { ImFileText } from "react-icons/im"
import { IoCaretForwardCircleOutline } from "react-icons/io5"
import { LuChevronUp, LuHeading, LuPaintbrush, LuPalette, LuSettings2, LuSparkles, LuTvMinimalPlay } from "react-icons/lu"
import { MdOutlineAccessTime, MdOutlineSubtitles, MdSpeed } from "react-icons/md"
import { RiShadowLine } from "react-icons/ri"
import { TbArrowForwardUp } from "react-icons/tb"
import { upath } from "@/lib/helpers/upath"
import { t } from "@/lib/i18n"
import { mc_parseCustomMpvConfig, mc_resolveAnime4KProfile } from "./mpv-core"
import type { MpvCoreAnime4KQuality, MpvCoreSettings, MpvCoreShaderMode, MpvCoreShaderSettings } from "./mpv-core.atoms"

const mpvSubtitleFontSizeOptions = [
    { label: t("player.menu.size_small"), value: 28 },
    { label: t("player.menu.size_medium"), value: 38 },
    { label: t("player.menu.size_large"), value: 48 },
    { label: t("player.menu.size_extra_large"), value: 58 },
]
const mpvSubtitleColorOptions = [
    { label: t("player.menu.color_white"), value: "#FFFFFF" },
    { label: t("player.menu.color_black"), value: "#000000" },
    { label: t("player.menu.color_gray"), value: "#808080" },
    { label: t("player.menu.color_yellow"), value: "#FFD700" },
    { label: t("player.menu.color_cyan"), value: "#00FFFF" },
    { label: t("player.menu.color_pink"), value: "#FF69B4" },
    { label: t("player.menu.color_purple"), value: "#9370DB" },
    { label: t("player.menu.color_lime"), value: "#00FF00" },
]
const mpvSubtitleOutlineOptions = [
    { label: t("player.menu.none"), value: 0 },
    { label: t("player.menu.size_small"), value: 2 },
    { label: t("player.menu.size_medium"), value: 3 },
    { label: t("player.menu.size_large"), value: 4 },
]
const mpvSubtitleShadowOptions = [
    { label: t("player.menu.none"), value: 0 },
    { label: t("player.menu.size_small"), value: 1 },
    { label: t("player.menu.size_medium"), value: 2 },
    { label: t("player.menu.size_large"), value: 3 },
]
const mpvSubtitleOpacityOptions = [
    { label: "100%", value: 1 },
    { label: "80%", value: 0.8 },
    { label: "70%", value: 0.7 },
    { label: "50%", value: 0.5 },
    { label: "25%", value: 0.25 },
    { label: "0%", value: 0 },
]

export interface MpvCoreSettingsMenuProps {
    openMenu: string | null
    openSection: string | null
    openSubSection: string | null
    setOpenMenu: (value: string | null) => void
    setOpenSection: (value: string | null) => void
    setOpenSubSection: (value: string | null) => void
    isFullscreen: boolean
    containerElement: HTMLElement | null
    speed: number
    changeSpeed: (value: number) => Promise<void>
    autoPlay: boolean
    setAutoPlay: (value: boolean) => void
    autoNext: boolean
    setAutoNext: (value: boolean) => void
    autoSkip: boolean
    setAutoSkip: (value: boolean) => void
    subtitleDelay: number
    setSubtitleDelay: (value: number) => void
    showChapterMarkers: boolean
    setChapterMarkers: (value: boolean) => void
    highlightOPEDChapters: boolean
    setHighlightOPEDChapters: (value: boolean) => void
    showStats: boolean
    setShowStats: (value: boolean) => void
    mpvSettings: MpvCoreSettings
    setMpvSettings: (
        update: MpvCoreSettings | ((current: MpvCoreSettings) => MpvCoreSettings)
    ) => void
    shaderSettings: MpvCoreShaderSettings
    setShaderSettings: (
        update: MpvCoreShaderSettings | ((current: MpvCoreShaderSettings) => MpvCoreShaderSettings)
    ) => void
    anime4kDirectory: MpvCoreAnime4KDirectory | null
    anime4kError: string | null
    onRefreshAnime4K: () => void
    onOpenPreferences: () => void
}

export function MpvCoreSettingsMenu(props: MpvCoreSettingsMenuProps) {
    const {
        openMenu,
        openSection,
        openSubSection,
        setOpenMenu,
        setOpenSection,
        setOpenSubSection,
        isFullscreen,
        containerElement,
        speed,
        changeSpeed,
        autoPlay,
        setAutoPlay,
        autoNext,
        setAutoNext,
        autoSkip,
        setAutoSkip,
        subtitleDelay,
        setSubtitleDelay,
        showChapterMarkers,
        setChapterMarkers,
        highlightOPEDChapters,
        setHighlightOPEDChapters,
        showStats,
        setShowStats,
        mpvSettings,
        setMpvSettings,
        shaderSettings,
        setShaderSettings,
        anime4kDirectory,
        anime4kError,
        onRefreshAnime4K,
        onOpenPreferences,
    } = props
    const serverStatus = useServerStatus()
    const { mutate: patchSetting } = usePatchSetting()

    const { parsed: parsedCustomConfig } = React.useMemo(() => mc_parseCustomMpvConfig(mpvSettings.customMpvConfig), [mpvSettings.customMpvConfig])
    const hasCustomDeband = "deband" in parsedCustomConfig
    const customDebandEnabled = hasCustomDeband
        ? (parsedCustomConfig["deband"] !== "no" && parsedCustomConfig["deband"] !== "false")
        : false
    const debandActive = hasCustomDeband ? customDebandEnabled : mpvSettings.deband

    const [subFontName, setSubFontName] = React.useState(mpvSettings.subtitleCustomization.fontName)

    React.useEffect(() => {
        if (openSection === t("player.menu.subtitle_styles")) {
            setSubFontName(mpvSettings.subtitleCustomization.fontName)
        }
    }, [openSection, mpvSettings.subtitleCustomization.fontName])

    const updateSubtitleStyle = <Key extends keyof MpvCoreSettings["subtitleCustomization"]>(
        key: Key,
        value: MpvCoreSettings["subtitleCustomization"][Key],
    ) => {
        setMpvSettings(current => ({
            ...current,
            subtitleCustomization: {
                ...current.subtitleCustomization,
                [key]: value,
            },
        }))
    }

    return (
        <MediaCoreMenu
            name="settings"
            openMenu={openMenu}
            onOpenMenuChange={setOpenMenu}
            onOpenSectionChange={setOpenSection}
            onOpenSubSectionChange={setOpenSubSection}
            isFullscreen={isFullscreen}
            containerElement={containerElement}
            trigger={
                <MediaCoreControlButtonIcon
                    icons={[["default", LuChevronUp]]}
                    state="default"
                    onClick={() => { }}
                    isMobile={false}
                    isMiniPlayer={false}
                />
            }
        >
            <MediaCoreMenuSectionBody show={!openSection}>
                <MediaCoreMenuTitle>{t("player.menu.title")}</MediaCoreMenuTitle>
                <MediaCoreMenuOption
                    title={t("player.menu.playback_speed")}
                    icon={MdSpeed}
                    value={`${speed.toFixed(2)}x`}
                    openSection={openSection}
                    onOpenSectionChange={setOpenSection}
                />
                <MediaCoreMenuOption
                    title={t("player.menu.auto_play")}
                    icon={IoCaretForwardCircleOutline}
                    value={autoPlay ? "On" : "Off"}
                    openSection={openSection}
                    onOpenSectionChange={setOpenSection}
                />
                <MediaCoreMenuOption
                    title={t("player.menu.auto_next")}
                    icon={HiFastForward}
                    value={autoNext ? "On" : "Off"}
                    openSection={openSection}
                    onOpenSectionChange={setOpenSection}
                />
                <MediaCoreMenuOption
                    title={t("player.menu.skip_op_ed")}
                    icon={TbArrowForwardUp}
                    value={autoSkip ? "On" : "Off"}
                    openSection={openSection}
                    onOpenSectionChange={setOpenSection}
                />
                <MediaCoreMenuOption
                    title={t("mpv.shader.title")}
                    icon={LuSparkles}
                    value={[
                        debandActive && "Deband",
                        shaderSettings.mode !== "off" && (
                            shaderSettings.mode === "custom"
                                ? "Custom"
                                : `${shaderSettings.anime4kMode.replace("mode-", "").toUpperCase()} (${shaderSettings.anime4kQuality.toUpperCase()})`
                        ),
                    ].filter(Boolean).join(", ") || "Off"}
                    openSection={openSection}
                    onOpenSectionChange={setOpenSection}
                />
                <MediaCoreMenuOption
                    title={t("player.menu.subtitle_delay")}
                    icon={MdOutlineAccessTime}
                    value={`${subtitleDelay.toFixed(1)}s`}
                    openSection={openSection}
                    onOpenSectionChange={setOpenSection}
                />
                <MediaCoreMenuOption
                    title={t("player.menu.subtitle_styles")}
                    icon={MdOutlineSubtitles}
                    value={mpvSettings.subtitleCustomization.enabled ? `On${!!mpvSettings.subtitleCustomization.fontName ? ", Font" : ""}` : "Off"}
                    openSection={openSection}
                    onOpenSectionChange={setOpenSection}
                />
                <MediaCoreMenuOption
                    title={t("player.menu.player_appearance")}
                    icon={LuTvMinimalPlay}
                    openSection={openSection}
                    onOpenSectionChange={setOpenSection}
                />
                <MediaCoreMenuOption
                    title={t("player.prefs.title")}
                    icon={LuSettings2}
                    openSection={openSection}
                    onOpenSectionChange={setOpenSection}
                    onClick={onOpenPreferences}
                />
            </MediaCoreMenuSectionBody>

            <MediaCoreMenuSubmenuBody show={!!openSection && !openSubSection}>
                <MediaCoreMenuOption
                    title={t("player.menu.subtitle_styles")}
                    icon={MdOutlineSubtitles}
                    openSection={openSection}
                    onOpenSectionChange={setOpenSection}
                >
                    <MediaCoreSettingSelect
                        options={[
                            { label: t("player.common.on"), value: 1 },
                            { label: t("player.common.off"), value: 0 },
                        ]}
                        onValueChange={value => updateSubtitleStyle("enabled", value === 1)}
                        value={mpvSettings.subtitleCustomization.enabled ? 1 : 0}
                    />
                    {mpvSettings.subtitleCustomization.enabled && (
                        <>
                            <p className="text-[--muted] text-sm my-2">{t("player.menu.options")}</p>
                            <MediaCoreMenuSubOption
                                title={t("player.menu.font")}
                                icon={LuHeading}
                                parentId={t("player.menu.subtitle_styles")}
                                value={!mpvSettings.subtitleCustomization.fontName ? "Default" : mpvSettings.subtitleCustomization.fontName?.slice(0,
                                    11) + (!!mpvSettings.subtitleCustomization.fontName?.length && mpvSettings.subtitleCustomization.fontName?.length > 10
                                    ? "..."
                                    : "")}
                                openSection={openSection}
                                openSubSection={openSubSection}
                                onOpenSubSectionChange={setOpenSubSection}
                            />
                            <MediaCoreMenuSubOption
                                title={t("player.menu.font_size")}
                                icon={LuHeading}
                                parentId={t("player.menu.subtitle_styles")}
                                value={`${mpvSettings.subtitleCustomization.fontSize}px`}
                                openSection={openSection}
                                openSubSection={openSubSection}
                                onOpenSubSectionChange={setOpenSubSection}
                            />
                            <MediaCoreMenuSubOption
                                title={t("player.menu.text_color")}
                                icon={LuPalette}
                                parentId={t("player.menu.subtitle_styles")}
                                value={mpvSubtitleColorOptions.find(option => option.value === mpvSettings.subtitleCustomization.primaryColor)?.label}
                                openSection={openSection}
                                openSubSection={openSubSection}
                                onOpenSubSectionChange={setOpenSubSection}
                            />
                            <MediaCoreMenuSubOption
                                title={t("player.menu.outline")}
                                icon={ImFileText}
                                parentId={t("player.menu.subtitle_styles")}
                                value={`${mpvSubtitleOutlineOptions.find(option => option.value === mpvSettings.subtitleCustomization.outline)?.label}, ${mpvSubtitleColorOptions.find(option => option.value === mpvSettings.subtitleCustomization.outlineColor)?.label}`}
                                openSection={openSection}
                                openSubSection={openSubSection}
                                onOpenSubSectionChange={setOpenSubSection}
                            />
                            <MediaCoreMenuSubOption
                                title={t("player.menu.shadow")}
                                icon={RiShadowLine}
                                parentId={t("player.menu.subtitle_styles")}
                                value={`${mpvSubtitleShadowOptions.find(option => option.value === mpvSettings.subtitleCustomization.shadow)?.label}, ${mpvSubtitleColorOptions.find(option => option.value === mpvSettings.subtitleCustomization.backColor)?.label}`}
                                openSection={openSection}
                                openSubSection={openSubSection}
                                onOpenSubSectionChange={setOpenSubSection}
                            />
                        </>
                    )}
                </MediaCoreMenuOption>
                <MediaCoreMenuOption
                    title={t("player.menu.playback_speed")}
                    icon={MdSpeed}
                    openSection={openSection}
                    onOpenSectionChange={setOpenSection}
                >
                    <MediaCoreSettingSelect
                        options={[0.5, 0.9, 1, 1.1, 1.5, 2, 3, 4].map(value => ({
                            label: `${value}x`,
                            value,
                        }))}
                        value={speed}
                        onValueChange={value => changeSpeed(Number(value))}
                        isFullscreen={isFullscreen}
                        containerElement={containerElement}
                    />
                </MediaCoreMenuOption>
                <MediaCoreMenuOption
                    title={t("player.menu.auto_play")}
                    icon={IoCaretForwardCircleOutline}
                    openSection={openSection}
                    onOpenSectionChange={setOpenSection}
                >
                    <MediaCoreSettingSelect
                        options={[
                            { label: t("player.common.on"), value: 1 },
                            { label: t("player.common.off"), value: 0 },
                        ]}
                        value={autoPlay ? 1 : 0}
                        onValueChange={value => setAutoPlay(Boolean(value))}
                    />
                </MediaCoreMenuOption>
                <MediaCoreMenuOption
                    title={t("player.menu.auto_next")}
                    icon={HiFastForward}
                    openSection={openSection}
                    onOpenSectionChange={setOpenSection}
                >
                    <MediaCoreSettingSelect
                        options={[
                            { label: t("player.common.on"), value: 1 },
                            { label: t("player.common.off"), value: 0 },
                        ]}
                        value={autoNext ? 1 : 0}
                        onValueChange={value => setAutoNext(Boolean(value))}
                    />
                </MediaCoreMenuOption>
                <MediaCoreMenuOption
                    title={t("player.menu.skip_op_ed")}
                    icon={TbArrowForwardUp}
                    openSection={openSection}
                    onOpenSectionChange={setOpenSection}
                >
                    <MediaCoreSettingSelect
                        options={[
                            { label: t("player.common.on"), value: 1 },
                            { label: t("player.common.off"), value: 0 },
                        ]}
                        value={autoSkip ? 1 : 0}
                        onValueChange={value => setAutoSkip(Boolean(value))}
                    />
                </MediaCoreMenuOption>
                <MediaCoreMenuOption
                    title={t("player.menu.subtitle_delay")}
                    icon={MdOutlineAccessTime}
                    openSection={openSection}
                    onOpenSectionChange={setOpenSection}
                >
                    <p className="text-sm text-[--muted] mb-2">{t("mpv.menu.subtitle_delay_hint")}</p>
                    <div className="flex gap-1.5 items-center mt-3">
                        {[-0.5, -0.1].map(delta => (
                            <Button
                                key={delta}
                                className="px-1 !text-xs flex-1"
                                intent="gray-subtle"
                                size="sm"
                                onClick={() => setSubtitleDelay(Number((subtitleDelay + delta).toFixed(1)))}
                            >
                                {delta}
                            </Button>
                        ))}
                        <span className="text-sm text-center text-[--muted] px-1 flex-1">
                            {subtitleDelay.toFixed(1)}s
                        </span>
                        {[0.1, 0.5].map(delta => (
                            <Button
                                key={delta}
                                className="px-1 !text-xs flex-1"
                                intent="gray-subtle"
                                size="sm"
                                onClick={() => setSubtitleDelay(Number((subtitleDelay + delta).toFixed(1)))}
                            >
                                +{delta}
                            </Button>
                        ))}
                    </div>
                    <MediaCoreSettingSelect
                        options={[-2, -1, -0.5, 0, 0.5, 1, 2].map(value => ({
                            label: `${value}s`,
                            value,
                        }))}
                        value={[-2, -1, -0.5, 0, 0.5, 1, 2].includes(subtitleDelay) ? subtitleDelay : null}
                        onValueChange={value => setSubtitleDelay(Number(value))}
                    />
                </MediaCoreMenuOption>
                <MediaCoreMenuOption
                    title={t("mpv.shader.title")}
                    icon={LuSparkles}
                    openSection={openSection}
                    onOpenSectionChange={setOpenSection}
                >
                    <div className="border-b border-[--border] pb-3 mb-3">
                        <Switch
                            label={t("mpv.shader.debanding")}
                            side="right"
                            fieldClass="hover:bg-transparent hover:border-transparent px-0 ml-0 w-full"
                            size="sm"
                            value={debandActive}
                            disabled={hasCustomDeband}
                            help={hasCustomDeband ? t("mpv.shader.debanding_help") : undefined}
                            onValueChange={checked => setMpvSettings(current => ({ ...current, deband: checked }))}
                        />
                    </div>
                    <p className="text-[--muted] text-sm mb-2">
                        {t("player.menu.anime4k_desc")}
                    </p>
                    <MediaCoreSettingSelect
                        options={[
                            { label: t("player.common.off"), value: "off" },
                            { label: t("mpv.shader.anime4k_preset"), value: "anime4k", description: t("mpv.shader.anime4k_preset_desc") },
                            { label: t("mpv.shader.custom_shaders"), value: "custom", description: t("mpv.shader.custom_shaders_desc") },
                        ]}
                        value={shaderSettings.mode}
                        onValueChange={value => setShaderSettings(current => ({
                            ...current,
                            mode: value as MpvCoreShaderMode,
                        }))}
                        isFullscreen={isFullscreen}
                        containerElement={containerElement}
                    />

                    {shaderSettings.mode === "anime4k" && (
                        <>
                            <p className="text-[--muted] text-sm my-2">{t("mpv.menu.preset")}</p>
                            <MediaCoreSettingSelect
                                options={[
                                    { label: t("player.a4k.mode_a"), value: "mode-a", description: t("player.a4k.desc_mode_a") },
                                    { label: t("player.a4k.mode_b"), value: "mode-b", description: t("player.a4k.desc_mode_b") },
                                    { label: t("player.a4k.mode_c"), value: "mode-c", description: t("player.a4k.desc_mode_c") },
                                    { label: t("player.a4k.mode_aa"), value: "mode-aa", description: t("player.a4k.desc_mode_aa") },
                                    { label: t("player.a4k.mode_bb"), value: "mode-bb", description: t("player.a4k.desc_mode_bb") },
                                    { label: t("player.a4k.mode_ca"), value: "mode-ca", description: t("player.a4k.desc_mode_ca") },
                                    { label: t("mpv.shader.a4k_label_cnn_2x_medium"), value: "cnn-2x-medium", description: t("player.a4k.desc_cnn_2x_m") },
                                    { label: t("mpv.shader.a4k_label_cnn_2x_very_large"), value: "cnn-2x-very-large", description: t("player.a4k.desc_cnn_2x_vl") },
                                    { label: t("mpv.shader.a4k_label_denoise_cnn_2x_vl"), value: "denoise-cnn-2x-very-large", description: t("player.a4k.desc_denoise_cnn_2x_vl") },
                                    { label: t("mpv.shader.a4k_label_cnn_2x_ultra_large"), value: "cnn-2x-ultra-large", description: t("player.a4k.desc_cnn_2x_ul") },
                                ]}
                                value={shaderSettings.anime4kMode}
                                onValueChange={value => setShaderSettings(current => ({
                                    ...current,
                                    anime4kMode: String(value),
                                }))}
                                isFullscreen={isFullscreen}
                                containerElement={containerElement}
                            />
                            <p className="text-[--muted] text-sm my-2">{t("player.menu.quality")}</p>
                            <MediaCoreSettingSelect
                                options={[
                                    { label: t("mpv.shader.quality_fast"), value: "fast" },
                                    { label: t("mpv.shader.quality_hq"), value: "hq" },
                                ]}
                                value={shaderSettings.anime4kQuality}
                                onValueChange={value => setShaderSettings(current => ({
                                    ...current,
                                    anime4kQuality: value as MpvCoreAnime4KQuality,
                                }))}
                                isFullscreen={isFullscreen}
                                containerElement={containerElement}
                            />
                        </>
                    )}

                    {shaderSettings.mode === "custom" && (
                        <div className="mt-4 border-t border-gray-800 pt-4 max-h-48 overflow-y-auto space-y-2 pr-1 custom-scrollbar">
                            <p className="text-[--muted] text-xs font-semibold uppercase tracking-wider mb-2">{t("mpv.shader.select_shaders")}</p>
                            {anime4kDirectory?.shaders.length ? (
                                anime4kDirectory.shaders.map(shader => {
                                    const isChecked = (shaderSettings.customShaders || []).includes(shader.name)
                                    return (
                                        <div key={shader.name} className="flex items-center justify-between text-sm py-0.5">
                                            <span className="truncate text-gray-300 mr-2" title={shader.name}>
                                                {shader.name.split("/").pop() || shader.name}
                                            </span>
                                            <Switch
                                                size="sm"
                                                value={isChecked}
                                                onValueChange={checked => {
                                                    setShaderSettings(current => {
                                                        const list = current.customShaders || []
                                                        const nextList = checked
                                                            ? [...list, shader.name]
                                                            : list.filter(name => name !== shader.name)
                                                        return { ...current, customShaders: nextList }
                                                    })
                                                }}
                                            />
                                        </div>
                                    )
                                })
                            ) : (
                                <p className="text-[--muted] text-sm italic">{t("mpv.shader.none_found")}</p>
                            )}
                        </div>
                    )}

                    {/*<p className="text-[--muted] text-sm my-2 break-all">*/}
                    {/*    {anime4kDirectory?.directory || shaderSettings.directory || "No shader folder selected"}*/}
                    {/*</p>*/}
                    <p className="text-[--muted] text-sm my-2">
                        {t("mpv.shader.detected_count", { count: anime4kDirectory?.shaders.length ?? 0 })}
                    </p>
                    {shaderSettings.mode === "anime4k" && mc_resolveAnime4KProfile(anime4kDirectory, shaderSettings.anime4kMode, shaderSettings.anime4kQuality).missing.length > 0 && (
                        <p className="text-red-300 text-sm mb-2 break-words">
                            {t("mpv.shader.missing_files", { files: mc_resolveAnime4KProfile(anime4kDirectory, shaderSettings.anime4kMode, shaderSettings.anime4kQuality).missing.join(", ") })}
                        </p>
                    )}
                    {anime4kError && <p className="text-red-300 text-sm mb-2 break-words">{anime4kError}</p>}
                    <div className="flex flex-wrap gap-2">
                        <Button
                            size="sm"
                            intent="gray-subtle"
                            onClick={() => window.electron?.mpvCore.openAnime4KDirectory(
                                anime4kDirectory?.directory || shaderSettings.directory,
                            )}
                        >
                            {t("mpv.shader.open_folder")}
                        </Button>
                        <Button size="sm" intent="gray-subtle" onClick={onRefreshAnime4K}>
                            {t("common.action.refresh")}
                        </Button>
                    </div>
                </MediaCoreMenuOption>
                <MediaCoreMenuOption
                    title={t("player.menu.player_appearance")}
                    icon={LuPaintbrush}
                    openSection={openSection}
                    onOpenSectionChange={setOpenSection}
                >
                    <Switch
                        label={t("player.menu.show_chapter_markers")}
                        side="right"
                        fieldClass="hover:bg-transparent hover:border-transparent px-0 ml-0 w-full"
                        size="sm"
                        value={showChapterMarkers}
                        onValueChange={setChapterMarkers}
                    />
                    <Switch
                        label={t("player.menu.highlight_skipped_chapters")}
                        side="right"
                        fieldClass="hover:bg-transparent hover:border-transparent px-0 ml-0 w-full"
                        size="sm"
                        value={highlightOPEDChapters}
                        onValueChange={setHighlightOPEDChapters}
                    />
                </MediaCoreMenuOption>
            </MediaCoreMenuSubmenuBody>
            <MediaCoreMenuSubSubmenuBody show={!!openSubSection}>
                <MediaCoreMenuSubOption
                    title={t("player.menu.font")}
                    icon={LuHeading}
                    parentId={t("player.menu.subtitle_styles")}
                    openSection={openSection}
                    openSubSection={openSubSection}
                    onOpenSubSectionChange={setOpenSubSection}
                >
                    <p className="text-sm mb-2">{t("player.menu.font")}</p>
                    <MediaCoreSettingTextInput
                        label={t("mpv.menu.font_name")}
                        value={subFontName}
                        onValueChange={setSubFontName}
                        help={t("mpv.menu.font_name_help")}
                    />
                    <div className="flex w-full mt-2">
                        <Button size="sm" intent="gray-subtle" onClick={() => updateSubtitleStyle("fontName", subFontName)}>
                            {t("media.action.save")}
                        </Button>
                    </div>
                </MediaCoreMenuSubOption>
                <MediaCoreMenuSubOption
                    title={t("player.menu.font_size")}
                    icon={LuHeading}
                    parentId={t("player.menu.subtitle_styles")}
                    openSection={openSection}
                    openSubSection={openSubSection}
                    onOpenSubSectionChange={setOpenSubSection}
                >
                    <MediaCoreSettingSelect
                        options={mpvSubtitleFontSizeOptions}
                        value={mpvSettings.subtitleCustomization.fontSize}
                        onValueChange={value => updateSubtitleStyle("fontSize", Number(value))}
                    />
                </MediaCoreMenuSubOption>
                <MediaCoreMenuSubOption
                    title={t("player.menu.text_color")}
                    icon={LuPalette}
                    parentId={t("player.menu.subtitle_styles")}
                    openSection={openSection}
                    openSubSection={openSubSection}
                    onOpenSubSectionChange={setOpenSubSection}
                >
                    <MediaCoreSettingSelect
                        options={mpvSubtitleColorOptions}
                        value={mpvSettings.subtitleCustomization.primaryColor}
                        onValueChange={value => updateSubtitleStyle("primaryColor", String(value))}
                    />
                </MediaCoreMenuSubOption>
                <MediaCoreMenuSubOption
                    title={t("player.menu.outline")}
                    icon={ImFileText}
                    parentId={t("player.menu.subtitle_styles")}
                    openSection={openSection}
                    openSubSection={openSubSection}
                    onOpenSubSectionChange={setOpenSubSection}
                >
                    <p className="text-[--muted] text-sm mb-2">{t("player.menu.outline_width")}</p>
                    <MediaCoreSettingSelect
                        options={mpvSubtitleOutlineOptions}
                        value={mpvSettings.subtitleCustomization.outline}
                        onValueChange={value => updateSubtitleStyle("outline", Number(value))}
                    />
                    <p className="text-[--muted] text-sm my-2">{t("player.menu.outline_color")}</p>
                    <MediaCoreSettingSelect
                        options={mpvSubtitleColorOptions}
                        value={mpvSettings.subtitleCustomization.outlineColor}
                        onValueChange={value => updateSubtitleStyle("outlineColor", String(value))}
                    />
                </MediaCoreMenuSubOption>
                <MediaCoreMenuSubOption
                    title={t("player.menu.shadow")}
                    icon={RiShadowLine}
                    parentId={t("player.menu.subtitle_styles")}
                    openSection={openSection}
                    openSubSection={openSubSection}
                    onOpenSubSectionChange={setOpenSubSection}
                >
                    <p className="text-[--muted] text-sm mb-2">{t("player.menu.shadow_depth")}</p>
                    <MediaCoreSettingSelect
                        options={mpvSubtitleShadowOptions}
                        value={mpvSettings.subtitleCustomization.shadow}
                        onValueChange={value => updateSubtitleStyle("shadow", Number(value))}
                    />
                    <p className="text-[--muted] text-sm my-2">{t("player.menu.shadow_opacity")}</p>
                    <MediaCoreSettingSelect
                        options={mpvSubtitleOpacityOptions}
                        value={mpvSettings.subtitleCustomization.backColorOpacity}
                        onValueChange={value => updateSubtitleStyle("backColorOpacity", Number(value))}
                    />
                    <p className="text-[--muted] text-sm my-2">{t("player.menu.shadow_color")}</p>
                    <MediaCoreSettingSelect
                        options={mpvSubtitleColorOptions}
                        value={mpvSettings.subtitleCustomization.backColor}
                        onValueChange={value => updateSubtitleStyle("backColor", String(value))}
                    />
                </MediaCoreMenuSubOption>
            </MediaCoreMenuSubSubmenuBody>
        </MediaCoreMenu>
    )
}
