import { useServerQuery } from "@/api/client/requests"
import { DiscoverDaily_Variables } from "@/api/generated/endpoint.types"
import { API_ENDPOINTS } from "@/api/generated/endpoints"
import { AL_ListAnime, AL_ListManga, Asmr_SearchResult } from "@/api/generated/types"

// Phase 3.9c：探索页「每日推荐」hooks。
// 后端 GET /api/v1/discover/daily?domain=… 按域复用现有条目 DTO，前端按当前 tab 域解析：
//   anime/novel → AL_ListAnime（Page.media，novel 为书籍语义的 type=MANGA 条目）
//   manga       → AL_ListManga（Page.media）
//   asmr        → Asmr_SearchResult（works）
// 失败/空数据时后端返回空列表，前端隐藏区块。

export type DiscoverDailyDomain = DiscoverDaily_Variables["domain"]

export function useDiscoverDaily(domain: DiscoverDailyDomain) {
    return useServerQuery<AL_ListAnime | AL_ListManga | Asmr_SearchResult, DiscoverDaily_Variables>({
        endpoint: API_ENDPOINTS.DISCOVER.DiscoverDaily.endpoint,
        method: API_ENDPOINTS.DISCOVER.DiscoverDaily.methods[0],
        queryKey: [API_ENDPOINTS.DISCOVER.DiscoverDaily.key, domain],
        data: { domain },
        enabled: true,
    })
}
