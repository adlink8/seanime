import { __advancedSearch_getValue, __advancedSearch_paramsAtom } from "@/app/(main)/search/_lib/advanced-search.atoms"
import { t } from "@/lib/i18n"
import { useAtomValue } from "jotai/react"
import capitalize from "lodash/capitalize"
import startCase from "lodash/startCase"
import React from "react"

export function AdvancedSearchPageTitle() {

    const params = useAtomValue(__advancedSearch_paramsAtom)

    const title = React.useMemo(() => {
        let str = ""
        if (params.title && params.title.length > 0) {
            str += startCase(params.title)
            return str
        }
        // if (!!__advancedSearch_getValue(params.genre)) str += params.genre?.join(", ") || ""
        if (__advancedSearch_getValue(params.sorting)?.includes("SCORE_DESC")) str += t("search.sort.highest_score")
        if (__advancedSearch_getValue(params.sorting)?.includes("TRENDING_DESC")) str += t("search.sort.trending_now")
        if (__advancedSearch_getValue(params.sorting)?.includes("POPULARITY_DESC")) str += t("search.sort.most_loved")
        if (__advancedSearch_getValue(params.sorting)?.includes("START_DATE_DESC")) str += t("search.sort.latest_released")
        if (__advancedSearch_getValue(params.sorting)?.includes("EPISODES_DESC")) str += t("search.sort.most_episodes")
        if (__advancedSearch_getValue(params.sorting)?.includes("CHAPTERS_DESC")) str += t("search.sort.most_chapters")
        if (!!__advancedSearch_getValue(params.genre)) str += ` ${params.genre?.join(", ")}`
        if (!str) str += t("search.sort.highest_score")
        if (params.type === "anime") str += ` ${t("search.type.anime")}`
        else str += ` ${t("search.type.manga")}`
        if (params.season || params.year) str += " •"
        if (params.season) str += ` ${capitalize(params.season)}`
        if (params.year) str += ` ${params.year}`
        if (!!str) return str
        return params.type === "anime" ? t("search.title.most_loved_anime") : t("search.title.most_loved_manga")
    }, [params.title, params.genre, params.sorting, params.type, params.season, params.year])

    // const secondaryTitle = React.useMemo(() => {
    //     let str = ""
    //     if (params.season) str += ` ${capitalize(params.season)}`
    //     if (params.year) str += ` ${params.year}`
    //     return str || null
    // }, [params.genre, params.season, params.year])

    return (
        <div data-advanced-search-page-title-container>
            <h2 data-advanced-search-page-title className="line-clamp-2">{title}</h2>
            {/*{secondaryTitle && <p className="text-xl line-clamp-1">{secondaryTitle}</p>}*/}
        </div>
    )
}
