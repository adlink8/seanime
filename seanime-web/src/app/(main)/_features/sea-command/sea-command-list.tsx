import { CommandGroup, CommandItem, CommandShortcut } from "@/components/ui/command"
import { t } from "@/lib/i18n"
import { usePathname, useSearchParams } from "@/lib/navigation"
import { useSeaCommandContext } from "./sea-command"

// renders when "/" is typed
export function SeaCommandList() {

    const pathname = usePathname()
    const searchParams = useSearchParams()
    const mediaId = Number(searchParams.get("id"))
    const isAnimePage = (pathname === "/entry" || pathname === "/offline/entry/anime") && Number.isFinite(mediaId) && mediaId > 0

    const { input, setInput, select, command: { isCommand, command, args }, scrollToTop } = useSeaCommandContext()

    const commands = [
        {
            command: "anime",
            description: t("misc.sea_command.find_in_collection"),
            show: true,
        },
        {
            command: "manga",
            description: t("misc.sea_command.find_in_collection"),
            show: true,
        },
        {
            command: "library",
            description: t("misc.sea_command.find_in_library"),
            show: true,
        },
        {
            command: "search",
            description: t("misc.sea_command.search_anilist"),
            show: true,
        },
        {
            command: "magnet",
            description: t("misc.sea_command.stream_or_download"),
            show: true,
        },
        {
            command: "logs",
            description: t("misc.sea_command.copy_logs_desc"),
            show: true,
        },
        {
            command: "issue",
            description: t("misc.sea_command.record_issue"),
            show: true,
        },
        {
            command: "droptorrent",
            description: t("misc.sea_command.drop_torrent_desc"),
            show: input.startsWith("/d"),
        },
        {
            command: "reload",
            description: t("misc.sea_command.reload_page"),
            show: input.startsWith("/r"),
        },
        {
            command: "spoilers",
            description: t("misc.sea_command.toggle_spoilers"),
            show: isAnimePage,
        },
    ]

    const filtered = commands.filter(n => n.show && n.command.startsWith(command) && n.command != command)

    if (!filtered?.length) return null

    return (
        <>
            <CommandGroup heading="Autocomplete">
                {filtered.map(command => (
                    <CommandItem
                        key={command.command}
                        onSelect={() => {
                            setInput(`/${command.command}`)
                        }}
                    >
                        <span className="tracking-widest text-sm">/{command.command}</span>
                        <CommandShortcut className="text-[--muted]">{command.description}</CommandShortcut>
                    </CommandItem>
                ))}
            </CommandGroup>
        </>
    )
}
