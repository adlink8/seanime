import { buildSeaQuery } from "@/api/client/requests"
import { API_ENDPOINTS } from "@/api/generated/endpoints"
import { AL_ListAnime, AL_ListManga, AL_ListNovel, Asmr_SearchResult } from "@/api/generated/types"
import { serverAuthTokenAtom } from "@/app/(main)/_atoms/server-status.atoms"
import { __advancedSearch_getValue, __advancedSearch_paramsAtom } from "@/app/(main)/search/_lib/advanced-search.atoms"
import { ADVANCED_SEARCH_SUBTITLE_ASMR, mapSortingToAsmrOrder } from "@/app/(main)/search/_lib/advanced-search-constants"
import { useInfiniteQuery } from "@tanstack/react-query"
import { useAtomValue } from "jotai/react"
import React from "react"

export function useAnilistAdvancedSearch() {

    const params = useAtomValue(__advancedSearch_paramsAtom)
    const password = useAtomValue(serverAuthTokenAtom)

    const { isLoading: isLoading1, data: data1, fetchNextPage: fetchNextPage1, hasNextPage: hasNextPage1 } = useInfiniteQuery({
        queryKey: ["advanced-search-anime", params],
        initialPageParam: 1,
        queryFn: async ({ pageParam }) => {
            const variables = {
                page: pageParam,
                perPage: 48,
                format: __advancedSearch_getValue(params.format)?.toUpperCase(),
                search: (params.title === null || params.title === "") ? undefined : params.title,
                genres: __advancedSearch_getValue(params.genre),
                tags: __advancedSearch_getValue(params.tags),
                season: __advancedSearch_getValue(params.season),
                seasonYear: __advancedSearch_getValue(params.year),
                averageScore_greater: __advancedSearch_getValue(params.minScore) !== undefined
                    ? __advancedSearch_getValue(params.minScore)
                    : undefined,
                sort: (params.title?.length && params.title.length > 0) ? ["SEARCH_MATCH",
                    ...(__advancedSearch_getValue(params.sorting) || ["SCORE_DESC"])] : (__advancedSearch_getValue(params.sorting) || ["SCORE_DESC"]),
                status: params.sorting?.includes("START_DATE_DESC") ? (__advancedSearch_getValue(params.status)
                    ?.filter((n: string) => n !== "NOT_YET_RELEASED") || ["FINISHED", "RELEASING"]) : __advancedSearch_getValue(params.status),
                isAdult: params.isAdult,
            }

            return buildSeaQuery<AL_ListAnime>({
                endpoint: API_ENDPOINTS.ANILIST.AnilistListAnime.endpoint,
                method: "POST",
                data: variables,
                password: password,
            })
        },
        getNextPageParam: (lastPage, pages) => {
            const curr = lastPage?.Page?.pageInfo?.currentPage
            const hasNext = lastPage?.Page?.pageInfo?.hasNextPage
            // console.log("lastPage", lastPage, "pages", pages, "curr", curr, "hasNext", hasNext, "nextPage", (!!curr && hasNext) ? pages.length + 1
            // : undefined)
            return (!!curr && hasNext) ? pages.length + 1 : undefined
        },
        enabled: params.active && params.type === "anime",
        refetchOnMount: true,
    })

    const { isLoading: isLoading2, data: data2, fetchNextPage: fetchNextPage2, hasNextPage: hasNextPage2 } = useInfiniteQuery({
        queryKey: ["advanced-search-manga", params],
        initialPageParam: 1,
        queryFn: async ({ pageParam }) => {
            const variables = {
                page: pageParam,
                perPage: 48,
                search: (params.title === null || params.title === "") ? undefined : params.title,
                genres: __advancedSearch_getValue(params.genre),
                tags: __advancedSearch_getValue(params.tags),
                year: __advancedSearch_getValue(params.year),
                format: __advancedSearch_getValue(params.format)?.toUpperCase(),
                averageScore_greater: __advancedSearch_getValue(params.minScore) !== undefined
                    ? __advancedSearch_getValue(params.minScore)
                    : undefined,
                sort: (params.title?.length && params.title.length > 0) ? ["SEARCH_MATCH",
                    ...(__advancedSearch_getValue(params.sorting) || ["SCORE_DESC"])] : (__advancedSearch_getValue(params.sorting) || ["SCORE_DESC"]),
                status: params.sorting?.includes("START_DATE_DESC") ? (__advancedSearch_getValue(params.status)
                    ?.filter((n: string) => n !== "NOT_YET_RELEASED") || ["FINISHED", "RELEASING"]) : __advancedSearch_getValue(params.status),
                countryOfOrigin: __advancedSearch_getValue(params.countryOfOrigin),
                isAdult: params.isAdult,
            }

            return buildSeaQuery<AL_ListManga>({
                endpoint: API_ENDPOINTS.MANGA.AnilistListManga.endpoint,
                method: "POST",
                data: variables,
                password: password,
            })
        },
        getNextPageParam: (lastPage, pages) => {
            const curr = lastPage?.Page?.pageInfo?.currentPage
            const hasNext = lastPage?.Page?.pageInfo?.hasNextPage
            // console.log("lastPage", lastPage, "pages", pages, "curr", curr, "hasNext", hasNext, "nextPage", (!!curr && hasNext) ? pages.length + 1
            // : undefined)
            return (!!curr && hasNext) ? pages.length + 1 : undefined
        },
        enabled: params.active && params.type === "manga",
        refetchOnMount: true,
    })

    // Phase 2.5：轻小说查询，变量构造复刻 anime（后端能映射多少算多少，契约第 3 节：请求形状与 list-anime 一致）
    const { isLoading: isLoading3, data: data3, fetchNextPage: fetchNextPage3, hasNextPage: hasNextPage3 } = useInfiniteQuery({
        queryKey: ["advanced-search-novel", params],
        initialPageParam: 1,
        queryFn: async ({ pageParam }) => {
            const variables = {
                page: pageParam,
                perPage: 48,
                search: (params.title === null || params.title === "") ? undefined : params.title,
                genres: __advancedSearch_getValue(params.genre),
                tags: __advancedSearch_getValue(params.tags),
                season: __advancedSearch_getValue(params.season),
                seasonYear: __advancedSearch_getValue(params.year),
                averageScore_greater: __advancedSearch_getValue(params.minScore) !== undefined
                    ? __advancedSearch_getValue(params.minScore)
                    : undefined,
                sort: (params.title?.length && params.title.length > 0) ? ["SEARCH_MATCH",
                    ...(__advancedSearch_getValue(params.sorting) || ["SCORE_DESC"])] : (__advancedSearch_getValue(params.sorting) || ["SCORE_DESC"]),
                status: params.sorting?.includes("START_DATE_DESC") ? (__advancedSearch_getValue(params.status)
                    ?.filter((n: string) => n !== "NOT_YET_RELEASED") || ["FINISHED", "RELEASING"]) : __advancedSearch_getValue(params.status),
                isAdult: params.isAdult,
            }

            return buildSeaQuery<AL_ListNovel>({
                endpoint: API_ENDPOINTS.NOVEL.AnilistListNovel.endpoint,
                method: "POST",
                data: variables,
                password: password,
            })
        },
        getNextPageParam: (lastPage, pages) => {
            const curr = lastPage?.Page?.pageInfo?.currentPage
            const hasNext = lastPage?.Page?.pageInfo?.hasNextPage
            return (!!curr && hasNext) ? pages.length + 1 : undefined
        },
        enabled: params.active && params.type === "novel",
        refetchOnMount: true,
    })

    // Phase 2.5：音声（ASMR）查询，参数映射契约第 4 节（keyword/order/page/perPage/subtitle）
    const { isLoading: isLoading4, data: data4, fetchNextPage: fetchNextPage4, hasNextPage: hasNextPage4 } = useInfiniteQuery({
        queryKey: ["advanced-search-asmr", params],
        initialPageParam: 1,
        queryFn: async ({ pageParam }) => {
            const subtitle = params.asmrSubtitle
            const variables = {
                keyword: params.title ?? "",
                order: mapSortingToAsmrOrder(params.sorting?.[0]),
                page: pageParam,
                perPage: 48,
                // 仅当所选字幕值是合法枚举时下发（"none" 在契约里表示过滤为无字幕，这里不区分、直接透传合法值）
                subtitle: subtitle && [...ADVANCED_SEARCH_SUBTITLE_ASMR.map(o => o.value), "none"].includes(subtitle)
                    ? subtitle
                    : undefined,
            }

            return buildSeaQuery<Asmr_SearchResult>({
                endpoint: API_ENDPOINTS.ASMR.AsmrSearch.endpoint,
                method: "POST",
                data: variables,
                password: password,
            })
        },
        getNextPageParam: (lastPage, pages) => {
            const curr = lastPage?.pageInfo?.currentPage
            const hasNext = lastPage?.pageInfo?.hasNextPage
            return (!!curr && hasNext) ? pages.length + 1 : undefined
        },
        enabled: params.active && params.type === "asmr",
        refetchOnMount: true,
    })

    const isLoading = React.useMemo(() =>
        params.type === "anime" ? isLoading1
            : params.type === "manga" ? isLoading2
                : params.type === "novel" ? isLoading3
                    : isLoading4,
        [isLoading1, isLoading2, isLoading3, isLoading4, params.type])
    const data = React.useMemo(() =>
        params.type === "anime" ? data1
            : params.type === "manga" ? data2
                : params.type === "novel" ? data3
                    : data4,
        [data1, data2, data3, data4, params.type])
    const fetchNextPage = React.useMemo(() =>
        params.type === "anime" ? fetchNextPage1
            : params.type === "manga" ? fetchNextPage2
                : params.type === "novel" ? fetchNextPage3
                    : fetchNextPage4,
        [fetchNextPage1, fetchNextPage2, fetchNextPage3, fetchNextPage4, params.type])
    const hasNextPage = React.useMemo(() =>
        params.type === "anime" ? hasNextPage1
            : params.type === "manga" ? hasNextPage2
                : params.type === "novel" ? hasNextPage3
                    : hasNextPage4,
        [hasNextPage1, hasNextPage2, hasNextPage3, hasNextPage4, params.type])

    return {
        isLoading,
        data,
        fetchNextPage,
        hasNextPage,
        type: params.type,
    }
}
