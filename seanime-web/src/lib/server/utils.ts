import capitalize from "lodash/capitalize"

export function getLibraryCollectionTitle(type?: string) {
    switch (type) {
        case "CURRENT":
            return "正在观看"
        case "COMPLETED":
            return "已看完"
        case "PLANNING":
            return "计划观看"
        case "PAUSED":
            return "暂停搁置"
        case "DROPPED":
            return "已弃番"
        case "REPEATING":
            return "二刷重温"
        default:
            return capitalize(type ?? "")
    }
}

export function getMangaCollectionTitle(type?: string) {
    switch (type) {
        case "CURRENT":
            return "正在阅读"
        case "COMPLETED":
            return "已看完"
        case "PLANNING":
            return "计划阅读"
        case "PAUSED":
            return "暂停搁置"
        case "DROPPED":
            return "已弃坑"
        case "REPEATING":
            return "重温阅读"
        default:
            return capitalize(type ?? "")
    }
}

export function formatDateAndTimeShort(date: string) {
    return new Date(date).toLocaleString(undefined, {
        dateStyle: "short",
        timeStyle: "short",
    })
}

export function isCustomSource(mId: number | null | undefined) {
    if (mId === null || mId === undefined) return false
    return mId >= 0x80000000
}

export function getCustomSourceExtensionId(m: { siteUrl?: string } | null | undefined) {
    if (!m?.siteUrl) return null
    let s = m.siteUrl.replace("ext_custom_source_", "")
    if (s.includes("|END|")) {
        s = s.split("|END|")[0]
    }
    return s
}

export function getCustomSourceMediaSiteUrl(m: { siteUrl?: string } | null | undefined) {
    if (!m?.siteUrl) return null
    let s = m.siteUrl.replace("ext_custom_source_", "")
    if (s.includes("|END|")) {
        s = s.split("|END|")?.[1] ?? null
    }
    return s
}
