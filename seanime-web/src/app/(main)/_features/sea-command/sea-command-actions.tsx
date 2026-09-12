import { useTorrentstreamDropTorrent } from "@/api/hooks/torrentstream.hooks"
import { __issueReport_overlayOpenAtom, __issueReport_recordingAtom } from "@/app/(main)/_features/issue-report/issue-report"
import { useHandleCopyLatestLogs } from "@/app/(main)/_hooks/logs"
import { CommandGroup, CommandItem, CommandShortcut } from "@/components/ui/command"
import { t } from "@/lib/i18n"
import { useSetAtom } from "jotai/react"
import React from "react"
import { useSeaCommandContext } from "./sea-command"

export function SeaCommandActions() {

    const { input, select, command: { isCommand, command, args }, scrollToTop, close } = useSeaCommandContext()

    const setIssueRecorderOpen = useSetAtom(__issueReport_overlayOpenAtom)
    const setIssueRecorderIsRecording = useSetAtom(__issueReport_recordingAtom)

    const { handleCopyLatestLogs } = useHandleCopyLatestLogs()
    const { mutate: dropTorrent, isPending: droppingTorrent } = useTorrentstreamDropTorrent()

    const reloadPage = () => {
        window.location.reload()
    }

    return (
        <>
            {command === "logs" && (
                <CommandGroup heading={t("misc.sea_command.actions")}>
                    <CommandItem
                        value="Logs"
                        onSelect={() => {
                            select(() => {
                                handleCopyLatestLogs()
                            })
                        }}
                    >
                        {t("misc.sea_command.copy_server_logs")}
                        <CommandShortcut>{t("misc.sea_command.enter")}</CommandShortcut>
                    </CommandItem>
                </CommandGroup>
            )}
            {command === "issue" && (
                <CommandGroup heading={t("misc.sea_command.actions")}>
                    <CommandItem
                        value="Issue"
                        onSelect={() => {
                            select(() => {
                                close()
                                React.startTransition(() => {
                                    setIssueRecorderOpen(true)
                                    setTimeout(() => {
                                        setIssueRecorderIsRecording(true)
                                    }, 500)
                                })
                            })
                        }}
                    >
                        {t("misc.sea_command.record_issue")}
                        <CommandShortcut>{t("misc.sea_command.enter")}</CommandShortcut>
                    </CommandItem>
                </CommandGroup>
            )}
            {command === "droptorrent" && (
                <CommandGroup heading={t("misc.sea_command.actions")}>
                    <CommandItem
                        value="Drop Torrent"
                        onSelect={() => {
                            close()
                            select(() => {
                                dropTorrent(undefined, {
                                    onSuccess: () => {
                                    },
                                })
                            })
                        }}
                    >
                        {t("misc.sea_command.drop_all_torrents")}
                        <CommandShortcut>{droppingTorrent ? t("misc.sea_command.dropping") : t("misc.sea_command.enter")}</CommandShortcut>
                    </CommandItem>
                </CommandGroup>
            )}
            {command === "reload" && (
                <CommandGroup heading={t("misc.sea_command.actions")}>
                    <CommandItem
                        value="Reload Page"
                        onSelect={() => {
                            close()
                            select(() => {
                                reloadPage()
                            })
                        }}
                    >
                        {t("misc.sea_command.reload_page")}
                        <CommandShortcut>{t("misc.sea_command.enter")}</CommandShortcut>
                    </CommandItem>
                </CommandGroup>
            )}
        </>
    )
}
