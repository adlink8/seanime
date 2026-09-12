import { usePlaybackPlayRandomVideo } from "@/api/hooks/playback_manager.hooks"
import { DropdownMenuItem } from "@/components/ui/dropdown-menu"
import { t } from "@/lib/i18n"
import React from "react"
import { LiaRandomSolid } from "react-icons/lia"

type PlayRandomEpisodeButtonProps = {
    children?: React.ReactNode
}

export function PlayRandomEpisodeButton(props: PlayRandomEpisodeButtonProps) {

    const {
        children,
        ...rest
    } = props

    const { mutate: playRandom, isPending } = usePlaybackPlayRandomVideo()

    return (
        <>
            <DropdownMenuItem>
                <LiaRandomSolid className="text-2xl" />
                <span>{t("library.play_random")}</span>
            </DropdownMenuItem>
        </>
    )
}
