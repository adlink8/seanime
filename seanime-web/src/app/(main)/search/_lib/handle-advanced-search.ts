import { buildSeaQuery } from "@/api/client/requests"
import { API_ENDPOINTS } from "@/api/generated/endpoints"
import { AL_ListAnime, AL_ListManga, AL_ListNovel, Asmr_SearchResult } from "@/api/generated/types"
import { serverAuthTokenAtom } from "@/app/(main)/_atoms/server-status.atoms"
import { __advancedSearch_paramsAtom } from "@/app/(main)/search/_lib/advanced-search.atoms"
import { ADVANCED_SEARCH_SUBTITLE_ASMR, mapSortingToAsmrOrder } from "@/app/(main)/search/_lib/advanced-search-constants"
// Phase 3.6b：查询参数组装抽到纯模块（无 `@/` 依赖，可被 vitest 直接覆盖）
import {
    ADVANCED_SEARCH_PER_PAGE,
    buildAdvancedSearchVariables,
    getAdvancedSearchNextPageParam,
} from "@/app/(main)/search/_lib/advanced-search-query"
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
            return buildSeaQuery<AL_ListAnime>({
                endpoint: API_ENDPOINTS.ANILIST.AnilistListAnime.endpoint,
                method: "POST",
                data: buildAdvancedSearchVariables("anime", params, pageParam),
                password: password,
            })
        },
        getNextPageParam: getAdvancedSearchNextPageParam,
        enabled: params.active && params.type === "anime",
        refetchOnMount: true,
    })

    const { isLoading: isLoading2, data: data2, fetchNextPage: fetchNextPage2, hasNextPage: hasNextPage2 } = useInfiniteQuery({
        queryKey: ["advanced-search-manga", params],
        initialPageParam: 1,
        queryFn: async ({ pageParam }) => {
            return buildSeaQuery<AL_ListManga>({
                endpoint: API_ENDPOINTS.MANGA.AnilistListManga.endpoint,
                method: "POST",
                data: buildAdvancedSearchVariables("manga", params, pageParam),
                password: password,
            })
        },
        getNextPageParam: getAdvancedSearchNextPageParam,
        enabled: params.active && params.type === "manga",
        refetchOnMount: true,
    })

    // Phase 2.5：轻小说查询，变量构造复刻 anime（后端能映射多少算多少，契约第 3 节：请求形状与 list-anime 一致）
    const { isLoading: isLoading3, data: data3, fetchNextPage: fetchNextPage3, hasNextPage: hasNextPage3 } = useInfiniteQuery({
        queryKey: ["advanced-search-novel", params],
        initialPageParam: 1,
        queryFn: async ({ pageParam }) => {
            return buildSeaQuery<AL_ListNovel>({
                endpoint: API_ENDPOINTS.NOVEL.AnilistListNovel.endpoint,
                method: "POST",
                data: buildAdvancedSearchVariables("novel", params, pageParam),
                password: password,
            })
        },
        getNextPageParam: getAdvancedSearchNextPageParam,
        enabled: params.active && params.type === "novel",
        refetchOnMount: true,
    })

    // Phase 2.5：音声（ASMR）查询，参数映射契约第 4 节（keyword/order/page/perPage/subtitle）
    // 注意：asmr.one 是另一条上游通路，不受 AniList v0 的 limit 硬钳 20 影响（契约 §0.2），
    // 这里沿用统一的页大小以保持列表行为一致。
    const { isLoading: isLoading4, data: data4, fetchNextPage: fetchNextPage4, hasNextPage: hasNextPage4 } = useInfiniteQuery({
        queryKey: ["advanced-search-asmr", params],
        initialPageParam: 1,
        queryFn: async ({ pageParam }) => {
            const subtitle = params.asmrSubtitle
            const variables = {
                keyword: params.title ?? "",
                order: mapSortingToAsmrOrder(params.sorting?.[0]),
                page: pageParam,
                perPage: ADVANCED_SEARCH_PER_PAGE,
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
