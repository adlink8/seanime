import { useServerMutation, useServerQuery } from "@/api/client/requests"
import {
    AsmrPlaybackPause_Variables,
    AsmrPlaybackSeek_Variables,
    AsmrPlaybackStatus,
    AsmrSearch_Variables,
    AsmrWork_Variables,
} from "@/api/generated/endpoint.types"
import { API_ENDPOINTS } from "@/api/generated/endpoints"
import {
    Asmr_CloudResponse,
    Asmr_DownloadResponse,
    Asmr_FavoriteResponse,
    Asmr_Library,
    Asmr_LocalWorkDetail,
    Asmr_SearchResult,
    Asmr_TrackProgressResponse,
    Asmr_WorkDetail,
    Nullish,
} from "@/api/generated/types"
import { useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"

// Phase 2.5：ASMR（音声）域 hooks，契约见 .planning/phases/02-library-reanchor/02.5-CONTRACT.md 第 4 节
// 无限翻页由搜索页的 useInfiniteQuery（handle-advanced-search.ts）实现，与现有 anime/manga 高级搜索模式一致。
// Phase 3.1：本地库 / 进度 / 收藏 / 下载 / 云端 hooks，契约见 03.1-CONTRACT.md。

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

/**
 * ASMR 本地库（扫描 RJ 目录，GET /api/v1/asmr/library）
 */
export function useAsmrLibrary(enabled?: boolean) {
    return useServerQuery<Asmr_Library>({
        endpoint: API_ENDPOINTS.ASMR.AsmrLibrary.endpoint,
        method: API_ENDPOINTS.ASMR.AsmrLibrary.methods[0],
        queryKey: [API_ENDPOINTS.ASMR.AsmrLibrary.key],
        enabled: enabled ?? true,
    })
}

/**
 * ASMR 本地作品详情（含音轨树，GET /api/v1/asmr/library/work/{rjId}）
 */
export function useAsmrLocalWork(rjId: Nullish<string>, enabled?: boolean) {
    return useServerQuery<Asmr_LocalWorkDetail>({
        endpoint: API_ENDPOINTS.ASMR.AsmrLibraryWork.endpoint.replace("{rjId}", String(rjId)),
        method: API_ENDPOINTS.ASMR.AsmrLibraryWork.methods[0],
        queryKey: [API_ENDPOINTS.ASMR.AsmrLibraryWork.key, String(rjId)],
        enabled: (!!rjId && enabled) ?? true,
    })
}

/**
 * ASMR 音轨完听进度（POST /api/v1/asmr/track/progress）
 * 成功后失效 library + 对应本地作品查询。
 */
export function useAsmrTrackProgress() {
    const queryClient = useQueryClient()

    return useServerMutation<Asmr_TrackProgressResponse, { rjId: string, trackPath: string, completed: boolean }>({
        endpoint: API_ENDPOINTS.ASMR.AsmrTrackProgress.endpoint,
        method: API_ENDPOINTS.ASMR.AsmrTrackProgress.methods[0],
        mutationKey: [API_ENDPOINTS.ASMR.AsmrTrackProgress.key],
        onSuccess: async (_data, variables) => {
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.ASMR.AsmrLibrary.key] })
            await queryClient.invalidateQueries({
                queryKey: [API_ENDPOINTS.ASMR.AsmrLibraryWork.key, variables.rjId],
            })
        },
    })
}

/**
 * ASMR 收藏（POST /api/v1/asmr/library/favorite）
 * syncedToCloud===false 时提示「本地已收藏（云端同步失败/未配置）」，本地状态照常生效。
 */
export function useAsmrFavorite() {
    const queryClient = useQueryClient()

    return useServerMutation<Asmr_FavoriteResponse, { rjId: string, favorite: boolean }>({
        endpoint: API_ENDPOINTS.ASMR.AsmrLibraryFavorite.endpoint,
        method: API_ENDPOINTS.ASMR.AsmrLibraryFavorite.methods[0],
        mutationKey: [API_ENDPOINTS.ASMR.AsmrLibraryFavorite.key],
        onSuccess: async (data) => {
            await queryClient.invalidateQueries({ queryKey: [API_ENDPOINTS.ASMR.AsmrLibrary.key] })
            if (data?.syncedToCloud === false) {
                toast.warning("本地已收藏（云端同步失败/未配置）")
            }
        },
    })
}

/**
 * ASMR 下载（POST /api/v1/asmr/download，异步）
 * 成功 toast「开始下载」。
 */
export function useAsmrDownload() {
    return useServerMutation<Asmr_DownloadResponse, { workId: string, rjId: string }>({
        endpoint: API_ENDPOINTS.ASMR.AsmrDownload.endpoint,
        method: API_ENDPOINTS.ASMR.AsmrDownload.methods[0],
        mutationKey: [API_ENDPOINTS.ASMR.AsmrDownload.key],
        onSuccess: async () => {
            toast.success("开始下载")
        },
    })
}

/**
 * ASMR 云端同步状态（GET /api/v1/asmr/cloud，可选）
 */
export function useAsmrCloud(enabled?: boolean) {
    return useServerQuery<Asmr_CloudResponse>({
        endpoint: API_ENDPOINTS.ASMR.AsmrCloud.endpoint,
        method: API_ENDPOINTS.ASMR.AsmrCloud.methods[0],
        queryKey: [API_ENDPOINTS.ASMR.AsmrCloud.key],
        enabled: enabled ?? true,
    })
}

/**
 * ASMR 播放控制条状态轮询（GET /api/v1/asmr/playback/status，Phase 3.1c 新增）
 * 仅弹窗打开（组件挂载）时 enabled，refetchInterval 1000ms（契约 D5）。
 * 无播放/非 mpv 时后端返回 running:false + 零值，前端据此不渲染控制条（契约 D3/D6）。
 */
export function useAsmrPlaybackStatus(enabled?: boolean) {
    return useServerQuery<AsmrPlaybackStatus>({
        endpoint: API_ENDPOINTS.ASMR.AsmrPlaybackStatus.endpoint,
        method: API_ENDPOINTS.ASMR.AsmrPlaybackStatus.methods[0],
        queryKey: [API_ENDPOINTS.ASMR.AsmrPlaybackStatus.key],
        refetchInterval: 1000,
        enabled: enabled ?? true,
    })
}

/**
 * ASMR 播放暂停/继续（POST /api/v1/asmr/playback/pause，Phase 3.1c 新增）
 * body { paused: bool }。无活跃播放时后端返 409，onError 静默吞掉（契约 D4）。
 */
export function useAsmrPlaybackPause() {
    return useServerMutation<boolean, AsmrPlaybackPause_Variables>({
        endpoint: API_ENDPOINTS.ASMR.AsmrPlaybackPause.endpoint,
        method: API_ENDPOINTS.ASMR.AsmrPlaybackPause.methods[0],
        mutationKey: [API_ENDPOINTS.ASMR.AsmrPlaybackPause.key],
        onError: () => {
            // 静默：无活跃播放（409）等错误不提示
        },
    })
}

/**
 * ASMR 播放跳转（POST /api/v1/asmr/playback/seek，Phase 3.1c 新增）
 * body { position: number }（绝对秒）。无活跃播放时后端返 409，onError 静默吞掉（契约 D4）。
 */
export function useAsmrPlaybackSeek() {
    return useServerMutation<boolean, AsmrPlaybackSeek_Variables>({
        endpoint: API_ENDPOINTS.ASMR.AsmrPlaybackSeek.endpoint,
        method: API_ENDPOINTS.ASMR.AsmrPlaybackSeek.methods[0],
        mutationKey: [API_ENDPOINTS.ASMR.AsmrPlaybackSeek.key],
        onError: () => {
            // 静默：无活跃播放（409）等错误不提示
        },
    })
}
