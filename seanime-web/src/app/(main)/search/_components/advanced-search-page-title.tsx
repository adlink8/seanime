import { __advancedSearch_getValue, __advancedSearch_paramsAtom } from "@/app/(main)/search/_lib/advanced-search.atoms"
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
        if (__advancedSearch_getValue(params.sorting)?.includes("SCORE_DESC")) str += "最高评分"
        if (__advancedSearch_getValue(params.sorting)?.includes("TRENDING_DESC")) str += "当前热门"
        if (__advancedSearch_getValue(params.sorting)?.includes("POPULARITY_DESC")) str += "最受喜爱"
        if (__advancedSearch_getValue(params.sorting)?.includes("START_DATE_DESC")) str += "最新推出"
        if (__advancedSearch_getValue(params.sorting)?.includes("EPISODES_DESC")) str += "剧集最多"
        if (__advancedSearch_getValue(params.sorting)?.includes("CHAPTERS_DESC")) str += "章节最多"
        if (!!__advancedSearch_getValue(params.genre)) str += ` ${params.genre?.join(", ")}`
        if (!str) str += "最高评分"
        if (params.type === "anime") str += " 动漫"
        else str += " 漫画"
        if (params.season || params.year) str += " •"
        if (params.season) str += ` ${capitalize(params.season)}`
        if (params.year) str += ` ${params.year}`
        if (!!str) return str
        return params.type === "anime" ? "最受喜爱动漫" : "最受喜爱漫画"
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
