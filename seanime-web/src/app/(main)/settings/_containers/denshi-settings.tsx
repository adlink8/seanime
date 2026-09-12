import { SettingsCard } from "@/app/(main)/settings/_components/settings-card"
import { Switch } from "@/components/ui/switch"
import { t } from "@/lib/i18n"
import React from "react"
import { RiSettings3Fill } from "react-icons/ri"

export function DenshiSettings() {

    const [settings, setSettings] = React.useState<DenshiSettings | null>(null)
    const settingsRef = React.useRef<DenshiSettings | null>(null)
    const [loading, setLoading] = React.useState(true)

    React.useEffect(() => {
        if (window.electron?.denshiSettings) {
            window.electron.denshiSettings.get().then((s) => {
                setSettings(s)
                settingsRef.current = s
                setLoading(false)
            })
        }
    }, [])

    function updateSetting(key: keyof DenshiSettings, value: boolean | string) {
        if (!settingsRef.current || !window.electron?.denshiSettings) return

        const newSettings = { ...settingsRef.current, [key]: value }
        settingsRef.current = newSettings
        setSettings(newSettings)
        window.electron.denshiSettings.set(newSettings)
    }

    if (loading || !settings) {
        return null
    }

    return (
        <div className="space-y-4">
            <SettingsCard title={t("settings.denshi.window_title")}>
                <Switch
                    side="right"
                    value={settings.minimizeToTray}
                    onValueChange={(v) => updateSetting("minimizeToTray", v)}
                    label={t("settings.denshi.minimize_to_tray")}
                    help={t("settings.denshi.minimize_to_tray_help")}
                />
                <Switch
                    side="right"
                    value={settings.openInBackground}
                    onValueChange={(v) => updateSetting("openInBackground", v)}
                    label={t("settings.denshi.open_background")}
                    help={t("settings.denshi.open_background_help")}
                />
            </SettingsCard>

            <SettingsCard title={t("settings.denshi.system_title")}>
                <Switch
                    side="right"
                    value={settings.openAtLaunch}
                    onValueChange={(v) => updateSetting("openAtLaunch", v)}
                    label={t("settings.denshi.open_at_launch")}
                    help={window.electron?.platform === "linux"
                        ? t("settings.denshi.open_at_launch_help_linux")
                        : t("settings.denshi.open_at_launch_help")}
                    disabled={window.electron?.platform === "linux"}
                />
            </SettingsCard>

            <div className="flex items-center gap-2 text-sm text-gray-500 bg-gray-50 dark:bg-gray-900/30 rounded-lg p-3 border border-gray-200 dark:border-gray-800 border-dashed">
                <RiSettings3Fill className="text-base" />
                <span>{t("settings.denshi.auto_save_note")}</span>
            </div>
        </div>
    )
}
