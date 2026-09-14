import { AL_BangumiInfoboxEntry, AL_BangumiTagInfo } from "@/api/generated/types"
import { Badge } from "@/components/ui/badge"
import { cn } from "@/components/ui/core/styling"
import { t } from "@/lib/i18n"
import { useRouter } from "@/lib/navigation"
import React from "react"
import { LuTrophy } from "react-icons/lu"
import { MdSell } from "react-icons/md"

type BangumiRankBadgeProps = {
    rank?: number
    className?: string
}

/**
 * Bangumi 排名徽章（Phase 3.9d）：bangumiRank > 0 时显示「Bangumi #N」。
 */
export function BangumiRankBadge(props: BangumiRankBadgeProps) {

    const {
        rank,
        className,
        ...rest
    } = props

    if (!rank || rank <= 0) return null

    return (
        <Badge
            data-bangumi-rank-badge
            size="lg"
            intent="gray"
            leftIcon={<LuTrophy className="text-sm" />}
            iconClass="text-yellow-500/70"
            className={cn("transition-all hover:opacity-100 rounded-full bg-transparent dark:text-[--muted] border-transparent px-0 hover:bg-transparent dark:hover:text-[--foreground]", className)}
            {...rest}
        >
            {t("media.bangumi.rank_badge", { rank: String(rank) })}
        </Badge>
    )
}

type BangumiTagsCloudProps = {
    tags?: Array<AL_BangumiTagInfo>
    type: "anime" | "manga"
    className?: string
}

/**
 * Bangumi 标签云（Phase 3.9d）：按 count 降序最多 20 个 chip；
 * spoiler 标签直接不显示。点击 chip 跳转搜索页按标签筛选。
 */
export function BangumiTagsCloud(props: BangumiTagsCloudProps) {

    const {
        tags,
        type,
        className,
        ...rest
    } = props

    const router = useRouter()

    const visibleTags = React.useMemo(() => {
        return (tags || [])
            .filter(tag => !tag.spoiler && !!tag.name)
            .sort((a, b) => b.count - a.count)
            .slice(0, 20)
    }, [tags])

    if (!visibleTags.length) return null

    return (
        <div
            data-bangumi-tags-cloud
            className={cn("flex flex-wrap items-center gap-2", className)}
            {...rest}
        >
            <Badge
                size="lg"
                intent="basic"
                className="bg-transparent border-transparent px-0 dark:text-[--muted]"
                leftIcon={<MdSell className="text-sm" />}
            >
                {t("media.bangumi.tags_title")}
            </Badge>
            {visibleTags.map(tag => (
                <Badge
                    key={tag.name}
                    data-bangumi-tags-cloud-item
                    size="md"
                    intent="gray"
                    role="button"
                    tabIndex={0}
                    onClick={() => router.push(`/search?tags=${encodeURIComponent(tag.name)}&type=${type}`)}
                    className="cursor-pointer transition-all hover:opacity-80"
                >
                    {tag.name}
                    {!!tag.count && <span className="ml-1 opacity-60">{tag.count}</span>}
                </Badge>
            ))}
        </div>
    )
}

type BangumiInfoboxListProps = {
    entries?: Array<AL_BangumiInfoboxEntry>
    className?: string
}

/**
 * Bangumi 信息箱（Phase 3.9d）：紧凑键值列表，最多显示前 8 条。
 * value 兼容单值字符串与 [{ v }] 对象数组两种形态。
 */
export function BangumiInfoboxList(props: BangumiInfoboxListProps) {

    const {
        entries,
        className,
        ...rest
    } = props

    const visibleEntries = React.useMemo(() => {
        return (entries || [])
            .filter(entry => !!entry?.key)
            .slice(0, 8)
    }, [entries])

    const formatValue = React.useCallback((value: AL_BangumiInfoboxEntry["value"]): string => {
        if (typeof value === "string") return value
        if (Array.isArray(value)) {
            return value.map(item => item?.v || item?.k || "").filter(Boolean).join(", ")
        }
        return ""
    }, [])

    if (!visibleEntries.length) return null

    return (
        <div
            data-bangumi-infobox
            className={cn("flex flex-col gap-1", className)}
            {...rest}
        >
            <Badge
                size="lg"
                intent="basic"
                className="bg-transparent border-transparent px-0 dark:text-[--muted]"
            >
                {t("media.bangumi.infobox_title")}
            </Badge>
            <div className="flex flex-col gap-1 rounded-[--radius] border border-[--border-color] divide-y divide-[--border-color] bg-[--bg] overflow-hidden">
                {visibleEntries.map((entry, index) => {
                    const value = formatValue(entry.value)
                    if (!value) return null
                    return (
                        <div key={`${entry.key}-${index}`} className="flex gap-3 px-3 py-1.5 text-sm">
                            <span className="text-[--muted] flex-none w-24 truncate">{entry.key}</span>
                            <span className="min-w-0 break-words">{value}</span>
                        </div>
                    )
                })}
            </div>
        </div>
    )
}
