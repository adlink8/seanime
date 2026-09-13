import { useServerQuery } from "@/api/client/requests"
import { AsmrSearch_Variables, AsmrWork_Variables } from "@/api/generated/endpoint.types"
import { API_ENDPOINTS } from "@/api/generated/endpoints"
import { Asmr_SearchResult, Asmr_WorkDetail, Nullish } from "@/api/generated/types"

// Phase 2.5：ASMR（音声）域 hooks，契约见 .planning/phases/02-library-reanchor/02.5-CONTRACT.md 第 4 节
// 无限翻页由搜索页的 useInfiniteQuery（handle-advanced-search.ts）实现，与现有 anime/manga 高级搜索模式一致。

/**
 * ASMR 搜索（单页查询，page 由 variables.page 指定，从 1 开始）
 */
export function useAsmrSearch(variables: AsmrSearch_Variables, enabled?: boolean) {
    return useServerQuery<Asmr_SearchResult, AsmrSearch_Variables>({
        endpoint: API_ENDPOINTS.ASMR.AsmrSearch.endpoint,
        method: API_ENDPOINTS.ASMR.AsmrSearch.methods[0],
        queryKey: [API_ENDPOINTS.ASMR.AsmrSearch.key, variables],
        data: variables,
        enabled: enabled ?? true,
    })
}

/**
 * 热门 ASMR（无请求参数）
 */
export function useAsmrPopular(enabled?: boolean) {
    return useServerQuery<Asmr_SearchResult>({
        endpoint: API_ENDPOINTS.ASMR.AsmrPopular.endpoint,
        method: API_ENDPOINTS.ASMR.AsmrPopular.methods[0],
        queryKey: [API_ENDPOINTS.ASMR.AsmrPopular.key],
        enabled: enabled ?? true,
    })
}

/**
 * ASMR 作品详情（含音轨树）
 */
export function useAsmrWork(id: Nullish<string>, enabled?: boolean) {
    return useServerQuery<Asmr_WorkDetail>({
        endpoint: API_ENDPOINTS.ASMR.AsmrWork.endpoint.replace("{id}", String(id)),
        method: API_ENDPOINTS.ASMR.AsmrWork.methods[0],
        queryKey: [API_ENDPOINTS.ASMR.AsmrWork.key, String(id)],
        enabled: (!!id && enabled) ?? true,
    })
}
