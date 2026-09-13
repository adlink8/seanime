/**
 * 高级搜索「查询语义」组装层。
 *
 * 为什么单独成模块：查询组装是纯逻辑（参数 → 请求体），抽离后可用 vitest 直接验证；
 * vitest 未配置 `@/` 别名，故此文件**不得**引入任何 `@/` 依赖
 * （先例：`src/api/client/sea-error-message.ts`）。
 *
 * 权威契约：`.planning/phases/02-library-reanchor/03.6b-CONTRACT.md`
 * §0.2 上游硬限制、§0.3 评分刻度、§1 决策矩阵 D2–D8。
 */

import { ADVANCED_SEARCH_FORMATS_ANIME } from "./advanced-search-constants"

/**
 * 契约 D8 / §0.2：上游 v0 检索的 `limit` 被服务端硬钳到 20。
 * 前端若请求更多，第二页的 offset 会按请求值跳步，20 之后的条目被整段跳过（漏条目）。
 * 因此前端页码步进必须与 20 对齐。
 *
 * ASMR 通路（asmr.one）不受该硬钳影响，但搜索页各类型统一使用同一页大小。
 */
export const ADVANCED_SEARCH_PER_PAGE = 20

export type AdvancedSearchType = "anime" | "manga" | "novel" | "asmr"

/** 查询组装所需的最小参数面（结构上兼容 `__advancedSearch_paramsAtom`） */
export interface AdvancedSearchQueryParams {
    title: string | null
    sorting: string[] | null
    genre: string[] | null
    tags: string[] | null
    format: string | null
    season: string | null
    year: string | null
    minScore: string | null
    isAdult: boolean
}

/**
 * 分页游标：后端按 `page` 自行换算 offset，因此这里只需按已加载页数步进。
 * 契约 §5.1 明确要求核对此处 —— 页码递增与 `perPage: 20` 必须同刻度。
 */
export function getAdvancedSearchNextPageParam(lastPage: any, pages: any[]): number | undefined {
    const curr = lastPage?.Page?.pageInfo?.currentPage
    const hasNext = lastPage?.Page?.pageInfo?.hasNextPage
    return (!!curr && hasNext) ? pages.length + 1 : undefined
}

/** 去掉首尾空白后的关键词；空串按「无关键词」处理 */
function normalizeKeyword(title: string | null): string | null {
    const trimmed = (title ?? "").trim()
    return trimmed.length ? trimmed : null
}

/** 年份按数值下发（`seasonYear` / `year` 在请求类型里均为 number） */
function toYearValue(year: string | null): number | undefined {
    if (!year) return undefined
    const value = Number(year)
    return Number.isFinite(value) ? value : undefined
}

/** 搜索页动画格式白名单：与选项列表同源，避免 UI 与请求体各写一份（契约 D4） */
const ADVANCED_SEARCH_FORMATS_ANIME_VALUES: string[] = ADVANCED_SEARCH_FORMATS_ANIME.map((option) => option.value)

/**
 * 契约 D6 / §0.3：本字段全程按 **0–100** 语义下发。
 * 后端负责 0–100 → ÷10 → `filter.rating:[">=N"]`（Bangumi 侧为 0–10 制）。
 * 历史故障：前端发 1–9，后端 ÷10 后恒 `>=0.x`，评分下限形同失效。
 */
export function toAverageScoreFilter(minScore: string | null | undefined): number | undefined {
    if (minScore === null || minScore === undefined || String(minScore).trim() === "") return undefined
    const value = Number(minScore)
    return Number.isFinite(value) ? value : undefined
}

/**
 * 契约 D7：上游排序能力只有 `heat` / `rank` / `score`（+ 关键词非空时的 `match`）。
 * `START_DATE_DESC` / `EPISODES_DESC` / `CHAPTERS_DESC` 无对应能力，**一律丢弃**——
 * 旧 URL 书签（`?sorting=...`）可能残留这些值，不能原样下发。
 */
export const ADVANCED_SEARCH_SUPPORTED_SORTS: readonly string[] = ["TRENDING_DESC", "POPULARITY_DESC", "SCORE_DESC"]

export const ADVANCED_SEARCH_DEFAULT_SORT = "SCORE_DESC"

/** 关键词非空时前置的相关度排序（上游 `match`） */
const ADVANCED_SEARCH_KEYWORD_SORT = "SEARCH_MATCH"

export function buildAdvancedSearchSortParam(sorting: string[] | null | undefined, title: string | null | undefined): string[] {
    const requested = (sorting ?? []).filter((value) => ADVANCED_SEARCH_SUPPORTED_SORTS.includes(value))
    const base = requested.length ? requested : [ADVANCED_SEARCH_DEFAULT_SORT]
    return normalizeKeyword(title ?? null) ? [ADVANCED_SEARCH_KEYWORD_SORT, ...base] : base
}

/**
 * 组装某个检索类型的一次请求参数。
 * 各分支的字段裁剪规则见契约 §1：
 *   D2 不再下发 `status`（上游静默忽略）
 *   D3 不再下发 `countryOfOrigin`（上游无该字段，`platform` 过滤已被实测证伪）
 *   D4 动画 `format` 仅允许 TV / ONA / OVA
 *   D5 漫画不再下发 `format`（`meta_tags` 对书籍分区恒 0）
 */
export function buildAdvancedSearchVariables(
    type: Exclude<AdvancedSearchType, "asmr">,
    params: AdvancedSearchQueryParams,
    page: number,
): Record<string, any> {
    const base = {
        page,
        perPage: ADVANCED_SEARCH_PER_PAGE,
        search: normalizeKeyword(params.title) ?? undefined,
        averageScore_greater: toAverageScoreFilter(params.minScore),
        sort: buildAdvancedSearchSortParam(params.sorting, params.title),
        isAdult: params.isAdult,
        genres: params.genre?.length ? params.genre : undefined,
        tags: params.tags?.length ? params.tags : undefined,
    }

    if (type === "anime") {
        return {
            ...base,
            format: ADVANCED_SEARCH_FORMATS_ANIME_VALUES.includes(params.format ?? "") ? params.format : undefined,
            season: params.season || undefined,
            seasonYear: toYearValue(params.year),
        }
    }

    if (type === "manga") {
        // 漫画没有季度概念，年份走 `year`；format / countryOfOrigin 一律不下发（D3 / D5）
        return {
            ...base,
            year: toYearValue(params.year),
        }
    }

    return {
        ...base,
        season: params.season || undefined,
        seasonYear: toYearValue(params.year),
    }
}
