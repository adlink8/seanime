import { AL_MediaSeason } from "@/api/generated/types"
import { useAnilistListNovel } from "@/api/hooks/anilist.hooks"
import { useAsmrPopular, useAsmrSearch } from "@/api/hooks/asmr.hooks"
import { useAnilistListAnime } from "@/api/hooks/anilist.hooks"
import { __discover_isAdultAtom } from "@/app/(main)/discover/_lib/discover.atoms"
import { atom } from "jotai"
import { useAtomValue } from "jotai/react"
import { useInView } from "motion/react"

export const __discover_trendingGenresAtom = atom<string[]>([])
export const __discover_currentSeasonGenresAtom = atom<string[]>([])
export const __discover_pastSeasonGenresAtom = atom<string[]>([])

// 契约 03.9c：探索页 18+ 开关——开启时各域列表请求透传 isAdult: true（后端按 EnableAdultContent 门控）
export function useDiscoverIsAdult(): boolean {
    return useAtomValue(__discover_isAdultAtom)
}

export function useDiscoverTrendingAnime() {
    const genres = useAtomValue(__discover_trendingGenresAtom)
    const isAdult = useDiscoverIsAdult()

    return useAnilistListAnime({
        page: 1,
        perPage: 20,
        sort: ["TRENDING_DESC"],
        genres: genres.length > 0 ? genres : undefined,
        isAdult: isAdult || undefined,
    }, true)

}

export function useDiscoverCurrentSeasonAnime(ref: any) {
    const genres = useAtomValue(__discover_currentSeasonGenresAtom)
    const isAdult = useDiscoverIsAdult()
    const isInView = useInView(ref, { once: true })
    const currentMonth = new Date().getMonth() + 1
    let currentYear = new Date().getFullYear()
    let season: AL_MediaSeason = "SUMMER"
    switch (currentMonth) {
        case 1:
        case 2:
        case 3:
            season = "WINTER"
            break
        case 4:
        case 5:
        case 6:
            season = "SPRING"
            break
        case 7:
        case 8:
        case 9:
            season = "SUMMER"
            break
        case 10:
        case 11:
        case 12:
            season = "FALL"
            break
    }


    return useAnilistListAnime({
        page: 1,
        perPage: 20,
        sort: ["SCORE_DESC"],
        season: season,
        seasonYear: currentYear,
        genres: genres.length > 0 ? genres : undefined,
        isAdult: isAdult || undefined,
    }, isInView)
}

export function useDiscoverPastSeasonAnime(ref: any) {
    const genres = useAtomValue(__discover_pastSeasonGenresAtom)
    const isAdult = useDiscoverIsAdult()
    const isInView = useInView(ref, { once: true })
    const currentMonth = new Date().getMonth() + 1
    const currentYear = new Date().getFullYear()
    let season: AL_MediaSeason = "SUMMER"
    switch (currentMonth) {
        case 1:
        case 2:
        case 3:
            season = "WINTER"
            break
        case 4:
        case 5:
        case 6:
            season = "SPRING"
            break
        case 7:
        case 8:
        case 9:
            season = "SUMMER"
            break
        case 10:
        case 11:
        case 12:
            season = "FALL"
            break
    }
    const pastSeason = season === "WINTER" ? "FALL" : season === "SPRING" ? "WINTER" : season === "SUMMER" ? "SPRING" : "SUMMER"
    const pastYear = season === "WINTER" ? currentYear - 1 : currentYear

    return useAnilistListAnime({
        page: 1,
        perPage: 20,
        sort: ["SCORE_DESC"],
        season: pastSeason,
        seasonYear: pastYear,
        genres: genres.length > 0 ? genres : undefined,
        isAdult: isAdult || undefined,
    }, isInView)
}

export function useDiscoverUpcomingAnime(ref: any) {
    const isInView = useInView(ref, { once: true })
    return useAnilistListAnime({
        page: 1,
        perPage: 20,
        sort: ["TRENDING_DESC"],
        status: ["NOT_YET_RELEASED"],
    }, isInView)
}

export function useDiscoverPopularAnime(ref: any) {
    const isInView = useInView(ref, { once: true })
    return useAnilistListAnime({
        page: 1,
        perPage: 20,
        sort: ["POPULARITY_DESC"],
    }, isInView)
}

export function useDiscoverTrendingMovies(ref: any) {
    const isInView = useInView(ref, { once: true })
    return useAnilistListAnime({
        page: 1,
        perPage: 20,
        format: "MOVIE",
        sort: ["TRENDING_DESC"],
        status: ["RELEASING", "FINISHED"],
    }, isInView)
}

// //////////////////////////////////////////////////////////////////////////////////////
// Phase 2.5：轻小说 + 音声（ASMR）探索区块
// //////////////////////////////////////////////////////////////////////////////////////

/** 由当前月份计算当前季度（与上方 anime 区块的季度逻辑一致） */
function getCurrentSeason(): AL_MediaSeason {
    const currentMonth = new Date().getMonth() + 1
    switch (currentMonth) {
        case 1:
        case 2:
        case 3:
            return "WINTER"
        case 4:
        case 5:
        case 6:
            return "SPRING"
        case 7:
        case 8:
        case 9:
            return "SUMMER"
        default:
            return "FALL"
    }
}

/** 「热门轻小说」区块 */
export function useDiscoverPopularNovel(ref: any) {
    const isInView = useInView(ref, { once: true })
    return useAnilistListNovel({
        page: 1,
        perPage: 24,
        sort: ["POPULARITY_DESC"],
    }, isInView)
}

/** 「本季轻小说」区块（前端计算当前年份月份→season） */
export function useDiscoverThisSeasonNovel(ref: any) {
    const isInView = useInView(ref, { once: true })
    const season = getCurrentSeason()
    const seasonYear = new Date().getFullYear()
    return useAnilistListNovel({
        page: 1,
        perPage: 24,
        sort: ["SCORE_DESC"],
        season,
        seasonYear,
    }, isInView)
}

/** 「热门音声」区块 */
export function useDiscoverPopularAsmr(ref: any) {
    const isInView = useInView(ref, { once: true })
    return useAsmrPopular(isInView)
}

/** 「最新音声」区块 */
export function useDiscoverLatestAsmr(ref: any) {
    const isInView = useInView(ref, { once: true })
    return useAsmrSearch({
        keyword: "",
        order: "dd",
        page: 1,
        perPage: 24,
    }, isInView)
}
