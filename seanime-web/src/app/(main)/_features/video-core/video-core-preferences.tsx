import { useSaveMediaPlayerSettings } from "@/api/hooks/settings.hooks"
import { getSkipPatternError } from "@/app/(main)/_features/media-core/media-core-chapters"
import { mediaCoreDefaultPreferences, mediaCorePreferencesAtom } from "@/app/(main)/_features/media-core/media-core-preferences"
import { vc_subtitleManager } from "@/app/(main)/_features/video-core/video-core"
import { vc_mediaCaptionsManager } from "@/app/(main)/_features/video-core/video-core"
import { vc_audioManager } from "@/app/(main)/_features/video-core/video-core"
import { VideoCoreChapterCue } from "@/app/(main)/_features/video-core/video-core"
import { vc_videoElement } from "@/app/(main)/_features/video-core/video-core-atoms"
import { vc_isMuted } from "@/app/(main)/_features/video-core/video-core-atoms"
import { vc_volume } from "@/app/(main)/_features/video-core/video-core-atoms"
import { vc_isFullscreen } from "@/app/(main)/_features/video-core/video-core-atoms"
import { vc_miniPlayer } from "@/app/(main)/_features/video-core/video-core-atoms"
import { vc_containerElement } from "@/app/(main)/_features/video-core/video-core-atoms"
import { vc_fullscreenManager } from "@/app/(main)/_features/video-core/video-core-fullscreen"
import { useVideoCoreInSight } from "@/app/(main)/_features/video-core/video-core-in-sight"
import { useVideoCoreOverlayFeedback } from "@/app/(main)/_features/video-core/video-core-overlay-display"
import { vc_pip } from "@/app/(main)/_features/video-core/video-core-pip"
import { vc_pipManager } from "@/app/(main)/_features/video-core/video-core-pip"
import { useVideoCorePlaylist } from "@/app/(main)/_features/video-core/video-core-playlist"
import {
    vc_defaultKeybindings,
    vc_initialSettings,
    vc_keybindingsAtom,
    vc_settings,
    vc_showStatsForNerdsAtom,
    vc_storedMutedAtom,
    vc_storedVolumeAtom,
    vc_useLibassRendererAtom,
    VideoCoreKeybindings,
} from "@/app/(main)/_features/video-core/video-core.atoms"
import { vc_dispatchAction } from "@/app/(main)/_features/video-core/video-core.utils"
import { DirectorySelector } from "@/components/shared/directory-selector"
import { Button } from "@/components/ui/button"
import { cn } from "@/components/ui/core/styling"
import { defineSchema, Field, Form } from "@/components/ui/form"
import { Modal } from "@/components/ui/modal"
import { NumberInput } from "@/components/ui/number-input"
import { Switch } from "@/components/ui/switch"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { TextInput } from "@/components/ui/text-input"
import { logger } from "@/lib/helpers/debug"
import { upath } from "@/lib/helpers/upath"
import { t } from "@/lib/i18n"
import { atom, useAtom, useAtomValue } from "jotai"
import { useSetAtom } from "jotai/react"
import React, { useCallback, useEffect, useRef, useState } from "react"
import { UseFormReturn } from "react-hook-form"
import { toast } from "sonner"
import { z } from "zod"
import { useServerStatus } from "../../_hooks/use-server-status"
import { useVideoCoreScreenshot } from "./video-core-screenshot"

export const videoCorePreferencesModalAtom = atom(false)

const tabsRootClass = cn("w-full contents space-y-4")

const tabsTriggerClass = cn(
    "text-base px-6 rounded-[--radius-md] w-fit border-none data-[state=active]:bg-[--subtle] data-[state=active]:text-white dark:hover:text-white",
    "h-10 lg:justify-center px-3 flex-1",
)

const tabsListClass = cn(
    "w-full flex flex-row lg:flex-row flex-wrap h-fit !mt-4",
)

const tabContentClass = cn(
    "space-y-4 animate-in fade-in-0 duration-300",
)

function isEditableKeyboardTarget(target: EventTarget | null) {
    if (!(target instanceof Element)) return false
    if (target instanceof HTMLElement && target.isContentEditable) return true

    return !!target.closest("input, textarea, select, [contenteditable='true'], [role='textbox']")
}

const translationSettingsSchema = defineSchema(({ z, presets }) => z.object({
    vcTranslate: z.boolean().default(false),
    vcTranslateProvider: z.string().default("google"),
    vcTranslateTargetLanguage: z.string().default("en"),
    vcTranslateApiKey: z.string().default(""),
    vcTranslateBaseUrl: z.string().default(""),
    vcTranslateModel: z.string().default(""),
}))

const KeybindingValueInput = ({
    actionKey,
    value,
    onValueChange,
}: {
    actionKey: keyof VideoCoreKeybindings
    value: number
    onValueChange: (value: number) => void
}) => {
    return (
        <NumberInput
            value={value}
            onValueChange={onValueChange}
            size="sm"
            fieldClass="w-16"
            hideControls
            min={0}
            step={actionKey.includes("Speed") ? 0.25 : 1}
            // onKeyDown={(e) => e.stopPropagation()}
            // onInput={(e) => e.stopPropagation()}
        />
    )
}

const KeybindingRow = ({
    action,
    description,
    actionKey,
    hasValue = false,
    valueLabel = "",
    editedKeybindings,
    setEditedKeybindings,
    recordingKey,
    handleKeyRecord,
    formatKeyDisplay = (actionKey: keyof VideoCoreKeybindings) => actionKey,
}: {
    action: string
    description: string
    actionKey: keyof VideoCoreKeybindings
    hasValue?: boolean
    valueLabel?: string
    editedKeybindings: VideoCoreKeybindings
    setEditedKeybindings: React.Dispatch<React.SetStateAction<VideoCoreKeybindings>>
    recordingKey: string | null
    handleKeyRecord: (actionKey: keyof VideoCoreKeybindings) => void
    formatKeyDisplay?: (actionKey: keyof VideoCoreKeybindings) => keyof VideoCoreKeybindings | string
}) => (
    <div className="flex items-center justify-between py-2 border rounded-lg px-3 bg-[--paper]">
        <div className="flex-1">
            <div className="font-medium text-sm">{action}</div>
            {hasValue && (
                <div className="flex items-center gap-2 mt-1">
                    <span className="text-xs text-muted-foreground">{valueLabel}:</span>
                    <KeybindingValueInput
                        actionKey={actionKey}
                        value={("value" in editedKeybindings[actionKey]) ? (editedKeybindings[actionKey] as any).value : 0}
                        onValueChange={(value) => {
                            setEditedKeybindings(prev => ({
                                ...prev,
                                [actionKey]: { ...prev[actionKey], value: value || 0 },
                            }))
                        }}
                    />
                </div>
            )}
        </div>
        <div className="flex items-center gap-2">
            <Button
                intent={recordingKey === actionKey ? "white-subtle" : "gray-subtle"}
                size="sm"
                onClick={() => handleKeyRecord(actionKey)}
                className={cn(
                    "h-8 px-3 text-lg font-mono",
                    recordingKey === actionKey && "!text-xs text-white",
                )}
            >
                {recordingKey === actionKey ? t("player.prefs.press_key") : formatKeyDisplay(editedKeybindings?.[actionKey]?.key as any ?? "" as any)}
            </Button>
        </div>
    </div>
)

export function VideoCorePreferencesModal({ isWebPlayer }: { isWebPlayer: boolean }) {
    const isFullscreen = useAtomValue(vc_isFullscreen)
    const containerElement = useAtomValue(vc_containerElement)
    const [open, setOpen] = useAtom(videoCorePreferencesModalAtom)
    const [keybindings, setKeybindings] = useAtom(vc_keybindingsAtom)
    const [editedKeybindings, setEditedKeybindings] = useState<VideoCoreKeybindings>(keybindings)
    const [useLibassRenderer, setUseLibassRenderer] = useAtom(vc_useLibassRendererAtom)
    const [editedUseLibassRenderer, setEditedUseLibassRenderer] = useState(useLibassRenderer)

    const [recordingKey, setRecordingKey] = useState<string | null>(null)

    const [tab, setTab] = useState("keybinds")
    const { mutate: saveMediaPlayerSettings } = useSaveMediaPlayerSettings()
    const serverStatus = useServerStatus()
    const mediaPlayerSettings = serverStatus?.settings?.mediaPlayer
    const translateProvider = mediaPlayerSettings?.vcTranslateProvider || "google"
    const translationFormRef = useRef<UseFormReturn<any>>(null)

    const [settings, setSettings] = useAtom(vc_settings)
    const [editedSubLanguage, setEditedSubLanguage] = useState(settings.preferredSubtitleLanguage)
    const [editedAudioLanguage, setEditedAudioLanguage] = useState(settings.preferredAudioLanguage)
    const [editedSubsBlacklist, setEditedSubsBlacklist] = useState(settings.preferredSubtitleBlacklist)
    const [editedSubtitleDelay, setEditedSubtitleDelay] = useState(settings.subtitleDelay ?? 0)
    const [editedScreenshotDir, setEditedScreenshotDir] = useState(mediaPlayerSettings?.screenshotDir ?? "")
    const [preferences, setPreferences] = useAtom(mediaCorePreferencesAtom)
    const [editedSkipPatterns, setEditedSkipPatterns] = useState(preferences.skipPatterns)
    const skipPatternError = React.useMemo(() => getSkipPatternError(editedSkipPatterns), [editedSkipPatterns])

    const isAbsolute = React.useMemo(() => {
        if (!editedScreenshotDir) return true
        return upath.isAbsolute(editedScreenshotDir)
    }, [editedScreenshotDir])
    // const [editedSubCustomization, setEditedSubCustomization] = useState<VideoCoreSettings["subtitleCustomization"]>(
    //     settings.subtitleCustomization || vc_initialSettings.subtitleCustomization
    // )
    const subtitleManager = useAtomValue(vc_subtitleManager)
    const mediaCaptionsManager = useAtomValue(vc_mediaCaptionsManager)

    // Reset edited keybindings and language preferences when modal opens
    useEffect(() => {
        if (open) {
            setEditedKeybindings(keybindings)
            setEditedSubLanguage(settings.preferredSubtitleLanguage)
            setEditedAudioLanguage(settings.preferredAudioLanguage)
            setEditedSubsBlacklist(settings.preferredSubtitleBlacklist)
            setEditedSubtitleDelay(settings.subtitleDelay ?? 0)
            setEditedUseLibassRenderer(useLibassRenderer)
            setEditedScreenshotDir(mediaPlayerSettings?.screenshotDir ?? "")
            setEditedSkipPatterns(preferences.skipPatterns)
            // setEditedSubCustomization(settings.subtitleCustomization || vc_initialSettings.subtitleCustomization)
        }
    }, [open, keybindings, settings, useLibassRenderer, mediaPlayerSettings, preferences.skipPatterns])

    const handleKeyRecord = (actionKey: keyof VideoCoreKeybindings) => {
        setRecordingKey(actionKey)

        const handleKeyDown = (e: KeyboardEvent) => {
            e.preventDefault()
            e.stopPropagation()
            e.stopImmediatePropagation()

            setEditedKeybindings(prev => ({
                ...prev,
                [actionKey]: {
                    ...prev[actionKey],
                    key: e.code,
                },
            }))

            setRecordingKey(null)
            document.removeEventListener("keydown", handleKeyDown, true)
        }

        document.addEventListener("keydown", handleKeyDown, true)
    }

    const handleSave = () => {
        if (skipPatternError) return
        setKeybindings(editedKeybindings)
        const newSettings = {
            ...settings,
            preferredSubtitleLanguage: editedSubLanguage,
            preferredAudioLanguage: editedAudioLanguage,
            preferredSubtitleBlacklist: editedSubsBlacklist,
            subtitleDelay: editedSubtitleDelay,
            // subtitleCustomization: editedSubCustomization,
        }
        setSettings(newSettings)
        setPreferences(current => ({ ...current, skipPatterns: editedSkipPatterns.trim() }))
        setUseLibassRenderer(editedUseLibassRenderer)

        const currentMediaPlayer = serverStatus?.settings?.mediaPlayer
        if (currentMediaPlayer) {
            saveMediaPlayerSettings({
                mediaPlayer: {
                    ...currentMediaPlayer,
                    screenshotDir: editedScreenshotDir,
                },
            })
        }

        // Update subtitle manager with new settings
        subtitleManager?.updateSettings(newSettings)
        mediaCaptionsManager?.updateSettings(newSettings)
        setOpen(false)
    }

    const handleReset = () => {
        setEditedKeybindings(vc_defaultKeybindings)
        setEditedSubLanguage(vc_initialSettings.preferredSubtitleLanguage)
        setEditedAudioLanguage(vc_initialSettings.preferredAudioLanguage)
        setEditedSubsBlacklist(vc_initialSettings.preferredSubtitleBlacklist)
        setEditedSubtitleDelay(vc_initialSettings.subtitleDelay)
        setEditedUseLibassRenderer(true)
        setEditedScreenshotDir(mediaPlayerSettings?.screenshotDir ?? "")
        setEditedSkipPatterns(mediaCoreDefaultPreferences.skipPatterns)
        // setEditedSubCustomization(vc_initialSettings.subtitleCustomization)
    }

    const formatKeyDisplay = (keyCode: string) => {
        const keyMap: Record<string, string> = {
            "KeyA": "A", "KeyB": "B", "KeyC": "C", "KeyD": "D", "KeyE": "E", "KeyF": "F",
            "KeyG": "G", "KeyH": "H", "KeyI": "I", "KeyJ": "J", "KeyK": "K", "KeyL": "L",
            "KeyM": "M", "KeyN": "N", "KeyO": "O", "KeyP": "P", "KeyQ": "Q", "KeyR": "R",
            "KeyS": "S", "KeyT": "T", "KeyU": "U", "KeyV": "V", "KeyW": "W", "KeyX": "X",
            "KeyY": "Y", "KeyZ": "Z",
            "ArrowUp": "↑", "ArrowDown": "↓", "ArrowLeft": "←", "ArrowRight": "→",
            "BracketLeft": "[", "BracketRight": "]",
            "Space": "⎵",
        }
        return keyMap[keyCode] || keyCode
    }

    function handleSaveTranslationSettings(data: z.infer<typeof translationSettingsSchema>) {
        const currentMediaPlayer = serverStatus?.settings?.mediaPlayer!
        const translateModel = data.vcTranslateProvider === "openai" && (!data.vcTranslateModel || data.vcTranslateModel === "local-model")
            ? "gpt-4o-mini"
            : data.vcTranslateModel

        saveMediaPlayerSettings({
            mediaPlayer: {
                ...currentMediaPlayer,
                vcTranslate: data.vcTranslate,
                vcTranslateTargetLanguage: data.vcTranslateTargetLanguage,
                vcTranslateProvider: data.vcTranslateProvider,
                vcTranslateApiKey: data.vcTranslateApiKey,
                vcTranslateBaseUrl: data.vcTranslateBaseUrl,
                vcTranslateModel: translateModel,
            },
        }, {
            onSuccess: () => {
                toast.success(t("player.prefs.translation_saved"))
                translationFormRef.current?.reset(translationFormRef.current.getValues())

                subtitleManager?.updateShouldTranslate(data.vcTranslate ? data.vcTranslateTargetLanguage : null)
                mediaCaptionsManager?.updateShouldTranslate(data.vcTranslate ? data.vcTranslateTargetLanguage : null)
            },
        })
    }

    return (
        <Modal
            title={t("player.prefs.title")}
            open={open}
            onOpenChange={setOpen}
            contentClass="max-w-5xl focus:outline-none focus-visible:outline-none outline-none bg-[--background] backdrop-blur-sm z-[101]"
            overlayClass="z-[150] bg-black/50"
            portalContainer={isFullscreen ? containerElement || undefined : undefined}
        >

            <Tabs
                value={tab}
                onValueChange={setTab}
                className={tabsRootClass}
                triggerClass={tabsTriggerClass}
                listClass={tabsListClass}
            >
                <TabsList className="flex-wrap max-w-full bg-[--paper] p-2 border rounded-xl">
                    <TabsTrigger value="keybinds">{t("player.prefs.tab_keybinds")}</TabsTrigger>
                    <TabsTrigger value="subtitles">{t("player.prefs.tab_subtitles_audio")}</TabsTrigger>
                    <TabsTrigger value="general">{t("player.prefs.tab_general")}</TabsTrigger>
                    <TabsTrigger value="translation">{t("player.prefs.tab_translation")}</TabsTrigger>
                    {/*<TabsTrigger value="browser-client">Rendering</TabsTrigger>*/}
                </TabsList>

                <TabsContent value="general" className={tabContentClass}>
                    <div className="space-y-4">
                        <TextInput
                            label={t("player.prefs.extra_skip_chapters")}
                            value={editedSkipPatterns}
                            onValueChange={setEditedSkipPatterns}
                            placeholder="^Intro$,^Outro$,^Preview$"
                            help={t("player.prefs.extra_skip_chapters_help")}
                            error={skipPatternError}
                            onKeyDown={event => event.stopPropagation()}
                            onInput={event => event.stopPropagation()}
                        />
                        <DirectorySelector
                            value={editedScreenshotDir}
                            onSelect={setEditedScreenshotDir}
                            label={t("player.prefs.screenshot_dir")}
                            help={t("player.prefs.screenshot_dir_help")}
                            error={!isAbsolute ? t("player.prefs.must_be_absolute_path") : ""}
                        />

                        <div className="flex items-center justify-between pt-6">
                            <Button
                                intent="gray-outline"
                                onClick={handleReset}
                            >
                                {t("player.prefs.reset_all")}
                            </Button>
                            <div className="flex gap-2">
                                <Button
                                    intent="gray-outline"
                                    onClick={() => setOpen(false)}
                                >
                                    {t("library.common.cancel")}
                                </Button>
                                <Button
                                    intent="primary"
                                    onClick={handleSave}
                                    disabled={!isAbsolute || !!skipPatternError}
                                >
                                    {t("media.action.save")}
                                </Button>
                            </div>
                        </div>
                    </div>
                </TabsContent>

                <TabsContent value="keybinds" className={tabContentClass}>
                    <div className="space-y-3 hidden lg:block">
                        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
                            <div>
                                {/* <h3 className="text-lg font-semibold mb-4 text-white">Playback</h3> */}
                                <div className="space-y-3">
                                    <KeybindingRow
                                        action={t("player.prefs.kb_seek_forward_fine")}
                                        description={t("player.prefs.kb_seek_forward_fine")}
                                        actionKey="seekForwardFine"
                                        editedKeybindings={editedKeybindings}
                                        setEditedKeybindings={setEditedKeybindings}
                                        recordingKey={recordingKey}
                                        handleKeyRecord={handleKeyRecord}
                                        formatKeyDisplay={formatKeyDisplay}
                                        hasValue={true}
                                        valueLabel={t("player.prefs.kb_value_seconds")}
                                    />
                                    <KeybindingRow
                                        action={t("player.prefs.kb_seek_backward_fine")}
                                        description={t("player.prefs.kb_seek_backward_fine")}
                                        actionKey="seekBackwardFine"
                                        editedKeybindings={editedKeybindings}
                                        setEditedKeybindings={setEditedKeybindings}
                                        recordingKey={recordingKey}
                                        handleKeyRecord={handleKeyRecord}
                                        formatKeyDisplay={formatKeyDisplay}
                                        hasValue={true}
                                        valueLabel={t("player.prefs.kb_value_seconds")}
                                    />
                                    <KeybindingRow
                                        action={t("player.prefs.kb_seek_forward")}
                                        description={t("player.prefs.kb_seek_forward")}
                                        actionKey="seekForward"
                                        editedKeybindings={editedKeybindings}
                                        setEditedKeybindings={setEditedKeybindings}
                                        recordingKey={recordingKey}
                                        handleKeyRecord={handleKeyRecord}
                                        formatKeyDisplay={formatKeyDisplay}
                                        hasValue={true}
                                        valueLabel={t("player.prefs.kb_value_seconds")}
                                    />
                                    <KeybindingRow
                                        action={t("player.prefs.kb_seek_backward")}
                                        description={t("player.prefs.kb_seek_backward")}
                                        actionKey="seekBackward"
                                        editedKeybindings={editedKeybindings}
                                        setEditedKeybindings={setEditedKeybindings}
                                        recordingKey={recordingKey}
                                        handleKeyRecord={handleKeyRecord}
                                        formatKeyDisplay={formatKeyDisplay}
                                        hasValue={true}
                                        valueLabel={t("player.prefs.kb_value_seconds")}
                                    />
                                    <KeybindingRow
                                        action={t("player.prefs.kb_increase_speed")}
                                        description={t("player.prefs.kb_increase_speed")}
                                        actionKey="increaseSpeed"
                                        editedKeybindings={editedKeybindings}
                                        setEditedKeybindings={setEditedKeybindings}
                                        recordingKey={recordingKey}
                                        handleKeyRecord={handleKeyRecord}
                                        formatKeyDisplay={formatKeyDisplay}
                                        hasValue={true}
                                        valueLabel={t("player.prefs.kb_value_increment")}
                                    />
                                    <KeybindingRow
                                        action={t("player.prefs.kb_decrease_speed")}
                                        description={t("player.prefs.kb_decrease_speed")}
                                        actionKey="decreaseSpeed"
                                        editedKeybindings={editedKeybindings}
                                        setEditedKeybindings={setEditedKeybindings}
                                        recordingKey={recordingKey}
                                        handleKeyRecord={handleKeyRecord}
                                        formatKeyDisplay={formatKeyDisplay}
                                        hasValue={true}
                                        valueLabel={t("player.prefs.kb_value_increment")}
                                    />
                                </div>
                            </div>

                            <div>
                                {/* <h3 className="text-lg font-semibold mb-4 text-white">Navigation</h3> */}
                                <div className="space-y-3">
                                    <KeybindingRow
                                        action={t("player.prefs.kb_next_chapter")}
                                        description={t("player.prefs.kb_next_chapter")}
                                        actionKey="nextChapter"
                                        editedKeybindings={editedKeybindings}
                                        setEditedKeybindings={setEditedKeybindings}
                                        recordingKey={recordingKey}
                                        handleKeyRecord={handleKeyRecord}
                                        formatKeyDisplay={formatKeyDisplay}
                                    />
                                    <KeybindingRow
                                        action={t("player.prefs.kb_previous_chapter")}
                                        description={t("player.prefs.kb_previous_chapter")}
                                        actionKey="previousChapter"
                                        editedKeybindings={editedKeybindings}
                                        setEditedKeybindings={setEditedKeybindings}
                                        recordingKey={recordingKey}
                                        handleKeyRecord={handleKeyRecord}
                                        formatKeyDisplay={formatKeyDisplay}
                                    />
                                    <KeybindingRow
                                        action={t("player.prefs.kb_next_episode")}
                                        description={t("player.prefs.kb_next_episode")}
                                        actionKey="nextEpisode"
                                        editedKeybindings={editedKeybindings}
                                        setEditedKeybindings={setEditedKeybindings}
                                        recordingKey={recordingKey}
                                        handleKeyRecord={handleKeyRecord}
                                        formatKeyDisplay={formatKeyDisplay}
                                    />
                                    <KeybindingRow
                                        action={t("player.prefs.kb_previous_episode")}
                                        description={t("player.prefs.kb_previous_episode")}
                                        actionKey="previousEpisode"
                                        editedKeybindings={editedKeybindings}
                                        setEditedKeybindings={setEditedKeybindings}
                                        recordingKey={recordingKey}
                                        handleKeyRecord={handleKeyRecord}
                                        formatKeyDisplay={formatKeyDisplay}
                                    />
                                    <KeybindingRow
                                        action={t("player.prefs.kb_cycle_subtitles")}
                                        description={t("player.prefs.kb_cycle_subtitles")}
                                        actionKey="cycleSubtitles"
                                        editedKeybindings={editedKeybindings}
                                        setEditedKeybindings={setEditedKeybindings}
                                        recordingKey={recordingKey}
                                        handleKeyRecord={handleKeyRecord}
                                        formatKeyDisplay={formatKeyDisplay}
                                    />
                                    <KeybindingRow
                                        action={t("player.prefs.kb_fullscreen")}
                                        description={t("player.prefs.kb_fullscreen")}
                                        actionKey="fullscreen"
                                        editedKeybindings={editedKeybindings}
                                        setEditedKeybindings={setEditedKeybindings}
                                        recordingKey={recordingKey}
                                        handleKeyRecord={handleKeyRecord}
                                        formatKeyDisplay={formatKeyDisplay}
                                    />
                                    <KeybindingRow
                                        action={t("player.prefs.kb_pip")}
                                        description={t("player.prefs.kb_pip")}
                                        actionKey="pictureInPicture"
                                        editedKeybindings={editedKeybindings}
                                        setEditedKeybindings={setEditedKeybindings}
                                        recordingKey={recordingKey}
                                        handleKeyRecord={handleKeyRecord}
                                        formatKeyDisplay={formatKeyDisplay}
                                    />
                                    <KeybindingRow
                                        action={t("player.prefs.kb_take_screenshot")}
                                        description={t("player.prefs.kb_take_screenshot")}
                                        actionKey="takeScreenshot"
                                        editedKeybindings={editedKeybindings}
                                        setEditedKeybindings={setEditedKeybindings}
                                        recordingKey={recordingKey}
                                        handleKeyRecord={handleKeyRecord}
                                        formatKeyDisplay={formatKeyDisplay}
                                    />
                                </div>
                            </div>

                            <div>
                                {/* <h3 className="text-lg font-semibold mb-4 text-white">Audio</h3> */}
                                <div className="space-y-3">
                                    <KeybindingRow
                                        action={t("player.prefs.kb_volume_up")}
                                        description={t("player.prefs.kb_volume_up")}
                                        actionKey="volumeUp"
                                        editedKeybindings={editedKeybindings}
                                        setEditedKeybindings={setEditedKeybindings}
                                        recordingKey={recordingKey}
                                        handleKeyRecord={handleKeyRecord}
                                        formatKeyDisplay={formatKeyDisplay}
                                        hasValue={true}
                                        valueLabel={t("player.prefs.kb_value_percent")}
                                    />
                                    <KeybindingRow
                                        action={t("player.prefs.kb_volume_down")}
                                        description={t("player.prefs.kb_volume_down")}
                                        actionKey="volumeDown"
                                        editedKeybindings={editedKeybindings}
                                        setEditedKeybindings={setEditedKeybindings}
                                        recordingKey={recordingKey}
                                        handleKeyRecord={handleKeyRecord}
                                        formatKeyDisplay={formatKeyDisplay}
                                        hasValue={true}
                                        valueLabel={t("player.prefs.kb_value_percent")}
                                    />
                                    <KeybindingRow
                                        action={t("player.prefs.kb_mute")}
                                        description={t("player.prefs.kb_mute")}
                                        actionKey="mute"
                                        editedKeybindings={editedKeybindings}
                                        setEditedKeybindings={setEditedKeybindings}
                                        recordingKey={recordingKey}
                                        handleKeyRecord={handleKeyRecord}
                                        formatKeyDisplay={formatKeyDisplay}
                                    />
                                    <KeybindingRow
                                        action={t("player.prefs.kb_cycle_audio")}
                                        description={t("player.prefs.kb_cycle_audio")}
                                        actionKey="cycleAudio"
                                        editedKeybindings={editedKeybindings}
                                        setEditedKeybindings={setEditedKeybindings}
                                        recordingKey={recordingKey}
                                        handleKeyRecord={handleKeyRecord}
                                        formatKeyDisplay={formatKeyDisplay}
                                    />
                                    <KeybindingRow
                                        action={t("player.prefs.kb_display_characters")}
                                        description={t("player.prefs.kb_display_characters")}
                                        actionKey="openInSight"
                                        editedKeybindings={editedKeybindings}
                                        setEditedKeybindings={setEditedKeybindings}
                                        recordingKey={recordingKey}
                                        handleKeyRecord={handleKeyRecord}
                                        formatKeyDisplay={formatKeyDisplay}
                                    />
                                    <KeybindingRow
                                        action={t("player.prefs.kb_stats_for_nerds")}
                                        description={t("player.prefs.kb_stats_for_nerds")}
                                        actionKey="statsForNerds"
                                        editedKeybindings={editedKeybindings}
                                        setEditedKeybindings={setEditedKeybindings}
                                        recordingKey={recordingKey}
                                        handleKeyRecord={handleKeyRecord}
                                        formatKeyDisplay={formatKeyDisplay}
                                    />
                                </div>
                            </div>
                        </div>

                        <div className="flex items-center justify-between pt-6">
                            <Button
                                intent="gray-outline"
                                onClick={handleReset}
                            >
                                {t("player.prefs.reset_all")}
                            </Button>
                            <div className="flex gap-2">
                                <Button
                                    intent="gray-outline"
                                    onClick={() => setOpen(false)}
                                >
                                    {t("library.common.cancel")}
                                </Button>
                                <Button
                                    intent="primary"
                                    onClick={handleSave}
                                    disabled={!!skipPatternError}
                                >
                                    {t("media.action.save")}
                                </Button>
                            </div>
                        </div>
                    </div>
                </TabsContent>
                <TabsContent value="subtitles" className={tabContentClass}>
                    <div className="space-y-3">
                        <h3 className="text-lg font-semibold text-white">{t("player.prefs.defaults")}</h3>
                        <div className="grid grid-cols-2 gap-4">
                            <div className="space-y-2">
                                <label className="text-sm font-medium text-muted-foreground">
                                    {t("player.prefs.preferred_sub_language")}
                                </label>
                                <TextInput
                                    value={editedSubLanguage}
                                    onValueChange={setEditedSubLanguage}
                                    placeholder="eng,jpn,spa"
                                    onKeyDown={(e) => e.stopPropagation()}
                                    onInput={(e) => e.stopPropagation()}
                                />
                            </div>
                            <div className="space-y-2">
                                <label className="text-sm font-medium text-muted-foreground">
                                    {t("player.prefs.preferred_audio_language")}
                                </label>
                                <TextInput
                                    value={editedAudioLanguage}
                                    onValueChange={setEditedAudioLanguage}
                                    placeholder="jpn,eng,kor"
                                    onKeyDown={(e) => e.stopPropagation()}
                                    onInput={(e) => e.stopPropagation()}
                                />
                            </div>
                        </div>
                        <div className="space-y-2">
                            <label className="text-sm font-medium text-muted-foreground">
                                {t("player.prefs.ignored_sub_names")}
                            </label>
                            <TextInput
                                value={editedSubsBlacklist}
                                onValueChange={setEditedSubsBlacklist}
                                placeholder="e.g. signs & songs,signs/songs"
                                onKeyDown={(e) => e.stopPropagation()}
                                onInput={(e) => e.stopPropagation()}
                                help={t("player.prefs.ignored_sub_names_help")}
                            />
                        </div>
                    </div>

                    {isWebPlayer && <div className="space-y-3">
                        <h3 className="text-lg font-semibold text-white">{t("player.prefs.rendering")}</h3>
                        <div className="space-y-2">
                            <Switch
                                side="right"
                                label={t("player.prefs.convert_soft_subs")}
                                value={editedUseLibassRenderer}
                                onValueChange={setEditedUseLibassRenderer}
                                help={t("player.prefs.convert_soft_subs_help")}
                            />
                        </div>
                    </div>}

                    <div className="flex items-center justify-between pt-6">
                        <Button
                            intent="gray-outline"
                            onClick={handleReset}
                        >
                            {t("player.prefs.reset_all")}
                        </Button>
                        <div className="flex gap-2">
                            <Button
                                intent="gray-outline"
                                onClick={() => setOpen(false)}
                            >
                                {t("library.common.cancel")}
                            </Button>
                            <Button
                                intent="primary"
                                onClick={handleSave}
                                disabled={!!skipPatternError}
                            >
                                {t("media.action.save")}
                            </Button>
                        </div>
                    </div>
                </TabsContent>
                <TabsContent value="translation" className={tabContentClass}>
                    <Form
                        schema={translationSettingsSchema}
                        onSubmit={handleSaveTranslationSettings}
                        defaultValues={{
                            vcTranslate: mediaPlayerSettings?.vcTranslate ?? false,
                            vcTranslateProvider: translateProvider,
                            vcTranslateTargetLanguage: mediaPlayerSettings?.vcTranslateTargetLanguage || "en",
                            vcTranslateApiKey: mediaPlayerSettings?.vcTranslateApiKey || "",
                            vcTranslateBaseUrl: mediaPlayerSettings?.vcTranslateBaseUrl || "http://localhost:1234/v1",
                            vcTranslateModel: mediaPlayerSettings?.vcTranslateModel || (translateProvider === "openai"
                                ? "gpt-4o-mini"
                                : "local-model"),
                        }}
                        stackClass="space-y-4 relative"
                        mRef={translationFormRef}
                    >
                        {(f) => {
                            const provider = f.watch("vcTranslateProvider")
                            const usesOpenAIProvider = provider === "openai" || provider === "openai-compatible"

                            return (
                                <div className="space-y-4">
                                    <div className="space-y-4">
                                        <Field.Switch
                                            name="vcTranslate"
                                            side="right"
                                            label={t("player.prefs.enable_translation")}
                                            help={t("player.prefs.enable_translation_help")}
                                        />
                                        <div className="space-y-2">
                                            <Field.Select
                                                label={t("player.prefs.provider")}
                                                name="vcTranslateProvider"
                                                options={[
                                                    { value: "google", label: t("player.prefs.provider_google_free") },
                                                    { value: "deepl", label: "DeepL" },
                                                    { value: "openai", label: "OpenAI" },
                                                    { value: "openai-compatible", label: t("player.prefs.provider_openai_compatible") },
                                                ]}
                                                contentClass="z-[999]"
                                            />
                                        </div>

                                        {provider === "deepl" && (
                                            <p>
                                                {t("player.prefs.deepl_limit")}
                                            </p>
                                        )}

                                        <div className="space-y-2">
                                            <Field.Select
                                                label={t("player.prefs.target_language")}
                                                name="vcTranslateTargetLanguage"
                                                options={[
                                                    // DeepL
                                                    { value: "en-US", label: t("player.prefs.lang_en_us") },
                                                    { value: "en-GB", label: t("player.prefs.lang_en_gb") },
                                                    { value: "es", label: t("player.prefs.lang_es") },
                                                    { value: "fr", label: t("player.prefs.lang_fr") },
                                                    { value: "de", label: t("player.prefs.lang_de") },
                                                    { value: "it", label: t("player.prefs.lang_it") },
                                                    { value: "pt-BR", label: t("player.prefs.lang_pt_br") },
                                                    { value: "pt-PT", label: t("player.prefs.lang_pt_pt") },
                                                    { value: "ru", label: t("player.prefs.lang_ru") },
                                                    { value: "ja", label: t("player.prefs.lang_ja") },
                                                    { value: "ko", label: t("player.prefs.lang_ko") },
                                                    { value: "zh-hans", label: t("player.prefs.lang_zh_hans") },
                                                    { value: "zh-hant", label: t("player.prefs.lang_zh_hant") },
                                                    { value: "ar", label: t("player.prefs.lang_ar") },
                                                    { value: "tr", label: t("player.prefs.lang_tr") },
                                                    { value: "pl", label: t("player.prefs.lang_pl") },
                                                    { value: "nl", label: t("player.prefs.lang_nl") },
                                                    { value: "sv", label: t("player.prefs.lang_sv") },
                                                    { value: "nb", label: t("player.prefs.lang_nb") },
                                                    { value: "da", label: t("player.prefs.lang_da") },
                                                    { value: "fi", label: t("player.prefs.lang_fi") },
                                                    { value: "el", label: t("player.prefs.lang_el") },
                                                    { value: "cs", label: t("player.prefs.lang_cs") },
                                                    { value: "hu", label: t("player.prefs.lang_hu") },
                                                    { value: "ro", label: t("player.prefs.lang_ro") },
                                                    { value: "id", label: t("player.prefs.lang_id") },
                                                    { value: "uk", label: t("player.prefs.lang_uk") },
                                                    { value: "bg", label: t("player.prefs.lang_bg") },
                                                    { value: "sk", label: t("player.prefs.lang_sk") },
                                                    { value: "sl", label: t("player.prefs.lang_sl") },
                                                    { value: "et", label: t("player.prefs.lang_et") },
                                                    { value: "lv", label: t("player.prefs.lang_lv") },
                                                    { value: "lt", label: t("player.prefs.lang_lt") },
                                                    // Not currently supported by DeepL
                                                    { value: "hi", label: t("player.prefs.lang_hi") },
                                                    { value: "bn", label: t("player.prefs.lang_bn") },
                                                    { value: "ta", label: t("player.prefs.lang_ta") },
                                                    { value: "te", label: t("player.prefs.lang_te") },
                                                    { value: "mr", label: t("player.prefs.lang_mr") },
                                                    { value: "kn", label: t("player.prefs.lang_kn") },
                                                    { value: "ml", label: t("player.prefs.lang_ml") },
                                                    { value: "pa", label: t("player.prefs.lang_pa") },
                                                    { value: "fa", label: t("player.prefs.lang_fa") },
                                                    { value: "ur", label: t("player.prefs.lang_ur") },
                                                    { value: "sw", label: t("player.prefs.lang_sw") },
                                                    { value: "af", label: t("player.prefs.lang_af") },
                                                    { value: "ms", label: t("player.prefs.lang_ms") },
                                                    { value: "hr", label: t("player.prefs.lang_hr") },
                                                    { value: "sr", label: t("player.prefs.lang_sr") },
                                                    { value: "he", label: t("player.prefs.lang_he") },
                                                    { value: "th", label: t("player.prefs.lang_th") },
                                                    { value: "vi", label: t("player.prefs.lang_vi") },
                                                ]}
                                                contentClass="z-[999]"
                                                help={t("player.prefs.target_language_help")}
                                            />
                                        </div>

                                        <div className="space-y-2">
                                            {provider === "openai-compatible" && (
                                                <Field.Text
                                                    label={t("player.prefs.base_url")}
                                                    name="vcTranslateBaseUrl"
                                                    placeholder="http://localhost:1234/v1"
                                                    onKeyDown={(e) => e.stopPropagation()}
                                                    onInput={(e) => e.stopPropagation()}
                                                    help={t("player.prefs.base_url_help")}
                                                />
                                            )}
                                            {usesOpenAIProvider && (
                                                <Field.Text
                                                    label={t("player.prefs.model")}
                                                    name="vcTranslateModel"
                                                    placeholder={provider === "openai-compatible" ? "local-model" : "gpt-4o-mini"}
                                                    onKeyDown={(e) => e.stopPropagation()}
                                                    onInput={(e) => e.stopPropagation()}
                                                />
                                            )}
                                            {provider !== "google" && (
                                                <Field.Text
                                                    label={provider === "openai-compatible" ? t("player.prefs.api_key_optional") : t("player.prefs.api_key")}
                                                    name="vcTranslateApiKey"
                                                    placeholder={t("player.prefs.api_key_placeholder")}
                                                    onKeyDown={(e) => e.stopPropagation()}
                                                    onInput={(e) => e.stopPropagation()}
                                                    type="password"
                                                />
                                            )}
                                        </div>
                                    </div>

                                    <p className="text-[--muted]">
                                        {t("player.prefs.reload_notice")}
                                    </p>

                                    <div className="flex items-center justify-end pt-6">
                                        <div className="flex gap-2">
                                            <Button
                                                type="button"
                                                intent="gray-outline"
                                                onClick={() => setOpen(false)}
                                            >
                                                {t("library.common.cancel")}
                                            </Button>
                                            <Button
                                                type="submit"
                                                intent="primary"
                                            >
                                                {t("media.action.save")}
                                            </Button>
                                        </div>
                                    </div>
                                </div>
                            )
                        }}
                    </Form>
                </TabsContent>
            </Tabs>

        </Modal>
    )
}

export function VideoCoreKeybindingController(props: {
    active: boolean
    videoRef: React.RefObject<HTMLVideoElement | null>,
    chapterCues: VideoCoreChapterCue[],
    introEndTime: number | undefined,
    introStartTime: number | undefined
    endingEndTime: number | undefined,
    endingStartTime: number | undefined
}) {
    const {
        active,
        chapterCues,
        introEndTime,
        introStartTime,
        endingEndTime,
        endingStartTime,
    } = props

    const [keybindings] = useAtom(vc_keybindingsAtom)
    const isKeybindingsModalOpen = useAtomValue(videoCorePreferencesModalAtom)
    const fullscreen = useAtomValue(vc_isFullscreen)
    const isMiniPlayer = useAtomValue(vc_miniPlayer)
    const pip = useAtomValue(vc_pip)
    const volume = useAtomValue(vc_volume)
    const setMiniPlayer = useSetAtom(vc_miniPlayer)
    const setVolume = useSetAtom(vc_storedVolumeAtom)
    const muted = useAtomValue(vc_isMuted)
    const setMuted = useSetAtom(vc_storedMutedAtom)
    const videoElement = useAtomValue(vc_videoElement)
    const { toggleOpen: toggleInSight } = useVideoCoreInSight()
    const { showOverlayFeedback } = useVideoCoreOverlayFeedback()

    const setShowStats = useSetAtom(vc_showStatsForNerdsAtom)

    const action = useSetAtom(vc_dispatchAction)

    const subtitleManager = useAtomValue(vc_subtitleManager)
    const mediaCaptionsManager = useAtomValue(vc_mediaCaptionsManager)
    const audioManager = useAtomValue(vc_audioManager)
    const fullscreenManager = useAtomValue(vc_fullscreenManager)
    const pipManager = useAtomValue(vc_pipManager)

    const { playEpisode, hasNextEpisode, hasPreviousEpisode } = useVideoCorePlaylist()

    // Rate limiting for seeking operations
    const lastSeekTime = useRef(0)
    const SEEK_THROTTLE_MS = 100 // Minimum time between seek operations

    function seek(seconds: number) {
        const isPaused = videoElement?.paused
        if (!isPaused) {
            videoElement?.pause()
        }
        action({ type: "seek", payload: { time: seconds, flashTime: true } })
        if (!isPaused) {
            videoElement?.play()?.catch()
        }
    }

    function seekTo(to: number) {
        const isPaused = videoElement?.paused
        if (!isPaused) {
            videoElement?.pause()
        }
        action({ type: "seekTo", payload: { time: to, flashTime: true } })
        if (!isPaused) {
            videoElement?.play()?.catch()
        }
    }

    const { takeScreenshot } = useVideoCoreScreenshot()

    //
    // Keyboard shortcuts
    //

    const handleKeyboardShortcuts = useCallback(async (e: KeyboardEvent) => {
            // Don't handle player shortcuts while typing or while keybindings modal is open.
            if (
                e.defaultPrevented ||
                isKeybindingsModalOpen ||
                isEditableKeyboardTarget(e.target) ||
                isEditableKeyboardTarget(document.activeElement)
            ) {
                return
            }

            // Ignore combinations with modifier keys
            if (e.ctrlKey || e.shiftKey || e.altKey || e.metaKey) {
                return
            }

            if (!videoElement || !active) {
                return
            }

            const video = videoElement

            if (isMiniPlayer) {
                if (e.code === keybindings.fullscreen.key) {
                    e.preventDefault()
                    setMiniPlayer(false)
                    window.requestAnimationFrame(() => {
                        window.requestAnimationFrame(() => {
                            handleToggleFullscreen()
                        })
                    })
                }
                return
            }

            if (e.code === "Space" || e.code === "Enter") {
                e.preventDefault()
                if (video.paused) {
                    await video.play()
                    showOverlayFeedback({ message: "PLAY", type: "icon" })
                } else {
                    video.pause()
                    showOverlayFeedback({ message: "PAUSE", type: "icon" })
                }
                return
            }

            // Home, go to beginning
            if (e.code === "Home") {
                e.preventDefault()
                seekTo(0)
                showOverlayFeedback({ message: t("player.overlay.beginning") })
                return
            }

            // End, go to end
            if (e.code === "End") {
                e.preventDefault()
                seekTo(video.duration)
                showOverlayFeedback({ message: t("player.overlay.end") })
                return
            }

            // Escape - Exit fullscreen
            if (e.code === "Escape" && fullscreen) {
                e.preventDefault()
                fullscreenManager?.exitFullscreen()
                return
            }

            // Number keys 0-9, seek to percentage (0%, 10%, 20%, ..., 90%)
            if (e.code.startsWith("Digit") && e.code.length === 6) {
                e.preventDefault()
                const digit = parseInt(e.code.slice(-1))
                const percentage = digit * 10
                const seekTime = Math.max(0, Math.min(video.duration, (video.duration * percentage) / 100))
                seekTo(seekTime)
                // showOverlayFeedback({ message: `${percentage}%` })
                return
            }

            // frame-by-frame seeking, assuming 24fps
            if (e.code === "Comma") {
                e.preventDefault()
                seek(-1 / 24)
                showOverlayFeedback({ message: t("player.overlay.prev_frame") })
                return
            }

            if (e.code === "Period") {
                e.preventDefault()
                seek(1 / 24)
                showOverlayFeedback({ message: t("player.overlay.next_frame") })
                return
            }

            // Helper function to check if seeking is rate limited
            const canSeek = () => {
                const now = Date.now()
                if (now - lastSeekTime.current < SEEK_THROTTLE_MS) {
                    return false
                }
                lastSeekTime.current = now
                return true
            }

            // Check which shortcut was pressed
            if (e.code === keybindings.seekForward.key) {
                e.preventDefault()
                if (!canSeek()) return

                if (props.introEndTime && props.introStartTime && video.currentTime < props.introEndTime && video.currentTime >= props.introStartTime) {
                    seekTo(props.introEndTime)
                    showOverlayFeedback({ message: t("player.overlay.skipped_opening") })
                    return
                }
                if (props.endingEndTime && props.endingStartTime && video.currentTime < props.endingEndTime && video.currentTime >= props.endingStartTime) {
                    seekTo(props.endingEndTime)
                    showOverlayFeedback({ message: t("player.overlay.skipped_ending") })
                    return
                }
                seek(keybindings.seekForward.value)
                video.dispatchEvent(new Event("seeked"))
            } else if (e.code === keybindings.seekBackward.key) {
                e.preventDefault()
                if (!canSeek()) return
                seek(-keybindings.seekBackward.value)
                video.dispatchEvent(new Event("seeked"))
            } else if (e.code === keybindings.seekForwardFine.key) {
                e.preventDefault()
                // if (!canSeek()) return
                video.dispatchEvent(new Event("seeking"))
                seek(keybindings.seekForwardFine.value)
                video.dispatchEvent(new Event("seeked"))
            } else if (e.code === keybindings.seekBackwardFine.key) {
                e.preventDefault()
                // if (!canSeek()) return
                video.dispatchEvent(new Event("seeking"))
                seek(-keybindings.seekBackwardFine.value)
                video.dispatchEvent(new Event("seeked"))
            } else if (e.code === keybindings.nextChapter.key) {
                e.preventDefault()
                handleNextChapter()
            } else if (e.code === keybindings.previousChapter.key) {
                e.preventDefault()
                handlePreviousChapter()
            } else if (e.code === keybindings.volumeUp.key) {
                e.preventDefault()
                const newVolume = Math.min(1, volume + keybindings.volumeUp.value / 100)
                setVolume(newVolume)
            } else if (e.code === keybindings.volumeDown.key) {
                e.preventDefault()
                const newVolume = Math.max(0, volume - keybindings.volumeDown.value / 100)
                setVolume(newVolume)
            } else if (e.code === keybindings.mute.key) {
                e.preventDefault()
                setMuted(!muted)
            } else if (e.code === keybindings.cycleSubtitles.key) {
                e.preventDefault()
                handleCycleSubtitles()
            } else if (e.code === keybindings.cycleAudio.key) {
                e.preventDefault()
                handleCycleAudio()
            } else if (e.code === keybindings.nextEpisode.key) {
                e.preventDefault()
                handleNextEpisode()
            } else if (e.code === keybindings.previousEpisode.key) {
                e.preventDefault()
                handlePreviousEpisode()
            } else if (e.code === keybindings.fullscreen.key) {
                e.preventDefault()
                handleToggleFullscreen()
            } else if (e.code === keybindings.pictureInPicture.key) {
                e.preventDefault()
                handleTogglePictureInPicture()
            } else if (e.code === keybindings.takeScreenshot.key) {
                e.preventDefault()
                takeScreenshot()
            } else if (e.code === keybindings.openInSight.key) {
                e.preventDefault()
                toggleInSight()
            } else if (e.code === keybindings.statsForNerds.key) {
                e.preventDefault()
                setShowStats(prev => !prev)
            } else if (e.code === keybindings.increaseSpeed.key) {
                e.preventDefault()
                const newRate = Math.min(8, video.playbackRate + keybindings.increaseSpeed.value)
                video.playbackRate = newRate
                showOverlayFeedback({ message: t("player.overlay.speed", { speed: newRate.toFixed(2) }) })
            } else if (e.code === keybindings.decreaseSpeed.key) {
                e.preventDefault()
                const newRate = Math.max(0.20, video.playbackRate - keybindings.decreaseSpeed.value)
                video.playbackRate = newRate
                showOverlayFeedback({ message: t("player.overlay.speed", { speed: newRate.toFixed(2) }) })
            }
        },
        [keybindings, volume, muted, seek, active, fullscreen, pip, showOverlayFeedback, introEndTime, introStartTime, isKeybindingsModalOpen,
            toggleInSight, videoElement, isMiniPlayer, setMiniPlayer])

    // Keyboard shortcut handlers
    const handleNextChapter = useCallback(() => {
        if (!videoElement || !chapterCues) return

        const currentTime = videoElement.currentTime

        // Sort chapters by start time to ensure proper order
        const sortedChapters = [...chapterCues].sort((a, b) => a.startTime - b.startTime)

        // Find the next chapter (with a small buffer to avoid edge cases)
        const nextChapter = sortedChapters.find(chapter => chapter.startTime > currentTime + 1)
        if (nextChapter) {
            seekTo(nextChapter.startTime)
            // Try to get chapter name from video track cues
            const chapterName = nextChapter.text
            showOverlayFeedback({ message: chapterName ? t("player.overlay.chapter", { name: chapterName }) : t("player.overlay.chapter_num", { num: sortedChapters.indexOf(nextChapter) + 1 }) })
        } else {
            // If no next chapter, go to the end
            const lastChapter = sortedChapters[sortedChapters.length - 1]
            if (lastChapter && lastChapter.endTime) {
                seekTo(lastChapter.endTime)
                showOverlayFeedback({ message: t("player.overlay.end_of_chapters") })
            }
        }
    }, [chapterCues, seekTo, showOverlayFeedback])

    const handlePreviousChapter = useCallback(() => {
        if (!videoElement || !chapterCues) return

        const currentTime = videoElement.currentTime

        // Sort chapters by start time to ensure proper order
        const sortedChapters = [...chapterCues].sort((a, b) => a.startTime - b.startTime)

        // Find the current chapter first
        const currentChapterIndex = sortedChapters.findIndex((chapter, index) => {
            const nextChapter = sortedChapters[index + 1]
            return chapter.startTime <= currentTime && (!nextChapter || currentTime < nextChapter.startTime)
        })

        if (currentChapterIndex > 0) {
            // Go to previous chapter
            const previousChapter = sortedChapters[currentChapterIndex - 1]
            seekTo(previousChapter.startTime)
            const chapterName = previousChapter.text
            showOverlayFeedback({ message: chapterName ? t("player.overlay.chapter", { name: chapterName }) : t("player.overlay.chapter_num", { num: currentChapterIndex }) })
        } else if (currentChapterIndex === 0) {
            // Already in first chapter, go to the beginning
            seekTo(0)
            const firstChapter = sortedChapters[0]
            const chapterName = firstChapter.text
            showOverlayFeedback({ message: chapterName ? t("player.overlay.chapter", { name: chapterName }) : t("player.overlay.chapter_num", { num: 1 }) })
        } else {
            // If we can't determine current chapter, just go to the beginning
            seekTo(0)
            showOverlayFeedback({ message: "Beginning" })
        }
    }, [chapterCues, seekTo, showOverlayFeedback])


    const handleCycleSubtitles = useCallback(() => {
        if (!videoElement) return
        // TODO: make it work when both types are combined
        let found = false
        if (subtitleManager) {
            // Cycle to next track or disable if we're at the end
            const nextTrackNumber = subtitleManager.getNextTrackNumber(subtitleManager.getSelectedTrackNumberOrNull())

            // Enable next track if available
            if (nextTrackNumber > -1) {
                subtitleManager?.selectTrack(nextTrackNumber)
                const trackName = subtitleManager.getTrack(nextTrackNumber)?.label || t("player.overlay.track", { num: nextTrackNumber })
                showOverlayFeedback({ message: t("player.overlay.subtitles", { track: trackName }) })
                found = true
            }
        }
        if (mediaCaptionsManager) {
            const currentTrackIdx = mediaCaptionsManager.getSelectedTrackIndexOrNull() ?? -1
            const nextTrackIdx = currentTrackIdx + 1
            const nextTrack = mediaCaptionsManager.getTrack(nextTrackIdx)

            // Enable next track if available
            if (nextTrack) {
                mediaCaptionsManager?.selectTrack(nextTrackIdx)
                const trackName = mediaCaptionsManager.getTrack(nextTrackIdx)?.label || t("player.overlay.track", { num: nextTrackIdx })
                showOverlayFeedback({ message: t("player.overlay.subtitles", { track: trackName }) })
                found = true
            }
        }

        if (!found) {
            showOverlayFeedback({ message: t("player.overlay.subtitles_off") })
            subtitleManager?.setNoTrack()
            mediaCaptionsManager?.setNoTrack()
        }
    }, [subtitleManager, mediaCaptionsManager])

    const handleCycleAudio = useCallback(() => {
        if (!videoElement || !audioManager) return

        // HLS stream
        if (audioManager.isHlsStream()) {
            const currentTrackNumber = audioManager.getSelectedTrackNumberOrNull()
            if (currentTrackNumber === null) {
                showOverlayFeedback({ message: t("player.overlay.no_audio_tracks") })
                return
            }
            const audioTracks = audioManager.getHlsAudioTracks()

            const nextTrackNumber = (currentTrackNumber + 1) % (audioTracks.length)

            const nextTrack = audioTracks.find(n => n.id === nextTrackNumber)
            if (nextTrack) {
                const trackName = nextTrack.name || nextTrack.language || t("player.overlay.track", { num: nextTrack.id + 1 })
                showOverlayFeedback({ message: t("player.overlay.audio", { track: trackName }) })
                audioManager.selectTrack(nextTrackNumber)
            }

            return
        }

        const audioTracks = videoElement.audioTracks
        if (!audioTracks || audioTracks.length <= 1) {
            showOverlayFeedback({ message: t("player.overlay.no_audio_tracks") })
            return
        }

        // Find currently enabled track
        let currentTrackIndex = -1
        for (let i = 0; i < audioTracks.length; i++) {
            if (audioTracks[i].enabled) {
                currentTrackIndex = i
                break
            }
        }

        // Cycle to next track
        const nextIndex = (currentTrackIndex + 1) % audioTracks.length

        // Disable all tracks first
        for (let i = 0; i < audioTracks.length; i++) {
            audioTracks[i].enabled = false
        }

        // Enable next track
        audioTracks[nextIndex].enabled = true
        audioTracks.dispatchEvent?.(new Event("change"))
        audioManager.syncSelectedTrack()

        const trackName = audioTracks[nextIndex].label || audioTracks[nextIndex].language || t("player.overlay.track", { num: nextIndex + 1 })
        showOverlayFeedback({ message: t("player.overlay.audio", { track: trackName }) })
    }, [audioManager])

    const log = logger("VideoCoreKeybindings")

    const handleNextEpisode = useCallback(() => {
        if (!hasNextEpisode) {
            showOverlayFeedback({ message: t("player.overlay.no_next_episode") })
            log.info("No next episode available")
            return
        }

        log.info("Playing next episode")
        playEpisode("next")
    }, [hasNextEpisode, playEpisode, showOverlayFeedback])

    const handlePreviousEpisode = useCallback(() => {
        if (!hasPreviousEpisode) {
            showOverlayFeedback({ message: t("player.overlay.no_prev_episode") })
            log.info("No previous episode available")
            return
        }

        log.info("Playing previous episode")
        playEpisode("previous")
    }, [hasPreviousEpisode, playEpisode, showOverlayFeedback])

    const handleToggleFullscreen = useCallback(() => {
        fullscreenManager?.toggleFullscreen()

        React.startTransition(() => {
            setTimeout(() => {
                videoElement?.focus()
            }, 100)
        })
    }, [fullscreenManager])

    const handleTogglePictureInPicture = useCallback(() => {
        pipManager?.enterPip()

        React.startTransition(() => {
            setTimeout(() => {
                videoElement?.focus()
            }, 100)
        })
    }, [pip, pipManager])

    // Add keyboard event listeners
    useEffect(() => {
        if (!active || !videoElement) return

        document.addEventListener("keydown", handleKeyboardShortcuts)

        return () => {
            document.removeEventListener("keydown", handleKeyboardShortcuts)
        }
    }, [handleKeyboardShortcuts, active, videoElement])

    // Handle fullscreen state changes to ensure video gets focused
    useEffect(() => {
        if (!active || !videoElement) return

        const handleFullscreenChange = () => {
            // Small delay to ensure fullscreen transition is complete
            setTimeout(() => {
                if (document.fullscreenElement && videoElement) {
                    videoElement.focus()
                }
            }, 100)
        }

        document.addEventListener("fullscreenchange", handleFullscreenChange)

        return () => {
            document.removeEventListener("fullscreenchange", handleFullscreenChange)
        }
    }, [active, videoElement])

    return null
}
