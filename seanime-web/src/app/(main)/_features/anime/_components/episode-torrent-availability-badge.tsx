import type { Anime_EpisodeTorrentAvailability } from "@/api/generated/types"
import { Badge } from "@/components/ui/badge"
import { t } from "@/lib/i18n"
import { LuCircleCheck, LuCircleHelp, LuClock3, LuLoaderCircle } from "react-icons/lu"

export function EpisodeTorrentAvailabilityBadge({ status }: { status?: Anime_EpisodeTorrentAvailability }) {
    if (status === "available") {
        return <Badge size="sm" intent="success-solid" leftIcon={<LuCircleCheck />} title={t("misc.anime.torrent_found_title")}>
            {t("misc.anime.torrent_available")}
        </Badge>
    }
    if (status === "checking") {
        return <Badge
            size="sm"
            intent="primary-solid"
            leftIcon={<LuLoaderCircle className="animate-spin" />}
            title={t("misc.anime.checking_provider_title")}
        >
            {t("misc.anime.checking_torrents")}
        </Badge>
    }
    if (status === "waiting") {
        return <Badge size="sm" intent="warning-solid" leftIcon={<LuClock3 />} title={t("misc.anime.no_torrent_title")}>
            {t("misc.anime.waiting_torrent")}
        </Badge>
    }
    if (status === "unknown") {
        return <Badge size="sm" intent="gray-solid" leftIcon={<LuCircleHelp />} title={t("misc.anime.provider_failed_title")}>
            {t("misc.anime.availability_unknown")}
        </Badge>
    }
    return null
}
