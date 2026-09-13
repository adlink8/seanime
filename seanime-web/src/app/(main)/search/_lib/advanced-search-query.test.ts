import { describe, expect, it } from "vitest"
import { ADVANCED_SEARCH_FORMATS_ANIME, ADVANCED_SEARCH_SORTING, ADVANCED_SEARCH_SORTING_MANGA } from "./advanced-search-constants"
import { buildAdvancedSearchVariables, getAdvancedSearchNextPageParam } from "./advanced-search-query"

/**
 * 高级搜索「查询语义」行为锁。
 *
 * 期望值全部来自 `.planning/phases/02-library-reanchor/03.6b-CONTRACT.md`：
 *   §0.2 上游硬限制：`limit` 硬钳 20、`total` 硬钳 1000；`platform`/`series`/`status` 静默忽略
 *   §0.2 `meta_tags` 仅动画且仅 TV/WEB/OVA 有效；书籍分区恒 0
 *   §0.3 评分全链路 0–10（后端收 0–100 再 ÷10）
 *   §1 D2/D3/D5/D6/D7/D8 决策矩阵
 * 字面量一律直接写死（不反向引用实现常量），避免同义反复。
 */

/** 还原「真正发到网线上」的对象：axios 会对 data 做 JSON 序列化，undefined 键被丢弃 */
function wire(variables: Record<string, any>): Record<string, any> {
    return JSON.parse(JSON.stringify(variables))
}

const ANIME_BASE = {
    title: null,
    sorting: null,
    genre: null,
    tags: null,
    format: null,
    season: null,
    year: null,
    minScore: null,
    isAdult: false,
} as const

describe("分页（契约 D8）", () => {
    // 上游 v0 `limit` 硬钳 20：前端若请求 48，第二页 offset=48 会跳过 20–47 之间的条目。
    it("每页请求 20 条（上游 limit 硬钳）", () => {
        expect(buildAdvancedSearchVariables("anime", ANIME_BASE, 1).perPage).toBe(20)
    })

    it("把无限查询的页码原样下发，由后端换算 offset", () => {
        expect(buildAdvancedSearchVariables("anime", ANIME_BASE, 3).page).toBe(3)
    })

    it("有下一页时按已加载页数递增", () => {
        const lastPage = { Page: { pageInfo: { currentPage: 2, hasNextPage: true } } }
        expect(getAdvancedSearchNextPageParam(lastPage, [{}, {}])).toBe(3)
    })

    it("没有下一页时停在最后一页", () => {
        const lastPage = { Page: { pageInfo: { currentPage: 2, hasNextPage: false } } }
        expect(getAdvancedSearchNextPageParam(lastPage, [{}, {}])).toBeUndefined()
    })

    it("响应形状缺失时不继续翻页", () => {
        expect(getAdvancedSearchNextPageParam(undefined, [{}])).toBeUndefined()
        expect(getAdvancedSearchNextPageParam({}, [{}])).toBeUndefined()
    })
})

describe("评分下限（契约 D6 / §0.3）", () => {
    // 后端收 0–100 后 ÷10 再拼 `filter.rating:[">=N"]`；前端若发 1–9，÷10 后恒为 0.x → 过滤失效。
    it("按 0–100 语义原样下发，不做 ×10", () => {
        expect(buildAdvancedSearchVariables("anime", { ...ANIME_BASE, minScore: "60" }, 1).averageScore_greater).toBe(60)
    })

    it("下限 100 原样下发", () => {
        expect(buildAdvancedSearchVariables("anime", { ...ANIME_BASE, minScore: "100" }, 1).averageScore_greater).toBe(100)
    })

    it("未选评分时不下发该字段", () => {
        expect("averageScore_greater" in wire(buildAdvancedSearchVariables("anime", ANIME_BASE, 1))).toBe(false)
        expect("averageScore_greater" in wire(buildAdvancedSearchVariables("anime", { ...ANIME_BASE, minScore: "" }, 1))).toBe(false)
    })

    it("三种 AniList 检索类型使用同一刻度", () => {
        for (const type of ["anime", "manga", "novel"] as const) {
            expect(buildAdvancedSearchVariables(type, { ...ANIME_BASE, minScore: "70" }, 1).averageScore_greater).toBe(70)
        }
    })
})

describe("排序收敛（契约 D7）", () => {
    // 上游排序能力只有 heat / rank / score / match；其余排序值无对应能力。
    it("未选排序时默认按评分", () => {
        expect(buildAdvancedSearchVariables("anime", ANIME_BASE, 1).sort).toEqual(["SCORE_DESC"])
    })

    it("保留受支持的排序值", () => {
        expect(buildAdvancedSearchVariables("anime", { ...ANIME_BASE, sorting: ["TRENDING_DESC"] }, 1).sort)
            .toEqual(["TRENDING_DESC"])
    })

    it("关键词非空时前置相关度排序", () => {
        expect(buildAdvancedSearchVariables("anime", { ...ANIME_BASE, title: "naruto", sorting: ["POPULARITY_DESC"] }, 1).sort)
            .toEqual(["SEARCH_MATCH", "POPULARITY_DESC"])
    })

    it("关键词只有空白时不算关键词", () => {
        expect(buildAdvancedSearchVariables("anime", { ...ANIME_BASE, title: "   " }, 1).sort).toEqual(["SCORE_DESC"])
    })

    // START_DATE_DESC / EPISODES_DESC / CHAPTERS_DESC 上游无对应能力（契约 D7），
    // 旧 URL 书签里可能残留这些值，必须被丢弃而不是原样下发。
    it("丢弃上游无能力的排序值并回落到默认排序", () => {
        for (const legacy of ["START_DATE_DESC", "EPISODES_DESC", "CHAPTERS_DESC"]) {
            expect(buildAdvancedSearchVariables("anime", { ...ANIME_BASE, sorting: [legacy] }, 1).sort).toEqual(["SCORE_DESC"])
            expect(buildAdvancedSearchVariables("manga", { ...ANIME_BASE, sorting: [legacy] }, 1).sort).toEqual(["SCORE_DESC"])
        }
    })
})

describe("移除上游无能力的筛选（契约 D2 / D3 / D5）", () => {
    // D2：status 上游静默忽略；D3：countryOfOrigin 上游无该字段（platform 过滤已被实测证伪）
    it("三种检索类型都不再下发 status", () => {
        for (const type of ["anime", "manga", "novel"] as const) {
            const vars = buildAdvancedSearchVariables(type, { ...ANIME_BASE, ...{ status: ["FINISHED"] } as any }, 1)
            expect("status" in wire(vars)).toBe(false)
        }
    })

    it("漫画与轻小说都不再下发 countryOfOrigin", () => {
        for (const type of ["manga", "novel"] as const) {
            const vars = buildAdvancedSearchVariables(type, { ...ANIME_BASE, ...{ countryOfOrigin: "JP" } as any }, 1)
            expect("countryOfOrigin" in wire(vars)).toBe(false)
        }
    })

    // D5：`meta_tags` 对书籍分区（type=1）恒 0，漫画侧下发 format 只会得到空结果
    it("漫画不再下发 format", () => {
        const vars = buildAdvancedSearchVariables("manga", { ...ANIME_BASE, format: "MANGA" }, 1)
        expect("format" in wire(vars)).toBe(false)
    })
})

describe("动画 format 收敛（契约 D4 / §0.2）", () => {
    // 上游 `meta_tags` 仅 TV / WEB(ONA) / OVA 三个值有效
    it("TV / ONA / OVA 原样下发", () => {
        for (const format of ["TV", "ONA", "OVA"]) {
            expect(buildAdvancedSearchVariables("anime", { ...ANIME_BASE, format }, 1).format).toBe(format)
        }
    })

    // MOVIE / TV_SHORT / SPECIAL 实测无有效 meta_tag（`Movie`/`剧场版` 全 0），不能再下发
    it("MOVIE / TV_SHORT / SPECIAL 不下发", () => {
        for (const format of ["MOVIE", "TV_SHORT", "SPECIAL"]) {
            expect("format" in wire(buildAdvancedSearchVariables("anime", { ...ANIME_BASE, format }, 1))).toBe(false)
        }
    })
})

describe("保留字段回归锁（契约 §0.2：tag 仅动画可用、air_date 双分区可用）", () => {
    it("动画下发流派 / 标签 / 季度 / 年份", () => {
        const vars = buildAdvancedSearchVariables("anime", {
            ...ANIME_BASE,
            genre: ["Action"],
            tags: ["Isekai"],
            season: "FALL",
            year: "2024",
        }, 2)
        expect(vars.genres).toEqual(["Action"])
        expect(vars.tags).toEqual(["Isekai"])
        expect(vars.season).toBe("FALL")
        expect(vars.seasonYear).toBe(2024)
        expect(vars.page).toBe(2)
    })

    it("漫画下发流派 / 标签 / 年份（无季度字段）", () => {
        const vars = buildAdvancedSearchVariables("manga", {
            ...ANIME_BASE,
            genre: ["Action"],
            tags: ["Isekai"],
            year: "2024",
        }, 1)
        expect(vars.genres).toEqual(["Action"])
        expect(vars.tags).toEqual(["Isekai"])
        expect(vars.year).toBe(2024)
        expect("season" in wire(vars)).toBe(false)
    })

    it("未选中的多选项不下发", () => {
        const vars = wire(buildAdvancedSearchVariables("anime", { ...ANIME_BASE, genre: [], tags: [] }, 1))
        expect("genres" in vars).toBe(false)
        expect("tags" in vars).toBe(false)
        expect("season" in vars).toBe(false)
    })
})

describe("筛选选项列表收敛（契约 D4 / D7 / §5.3）", () => {
    const values = (options: { value: string }[]) => options.map((option) => option.value)

    it("动画格式只剩 TV / ONA / OVA", () => {
        expect(values(ADVANCED_SEARCH_FORMATS_ANIME)).toEqual(["TV", "ONA", "OVA"])
    })

    it("动画排序只剩趋势 / 热度 / 评分", () => {
        expect(values(ADVANCED_SEARCH_SORTING)).toEqual(["TRENDING_DESC", "POPULARITY_DESC", "SCORE_DESC"])
    })

    it("漫画排序与动画同集合（上游排序能力相同）", () => {
        expect(values(ADVANCED_SEARCH_SORTING_MANGA)).toEqual(["TRENDING_DESC", "POPULARITY_DESC", "SCORE_DESC"])
    })
})
