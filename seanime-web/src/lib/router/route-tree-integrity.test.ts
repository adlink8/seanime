import * as fs from "node:fs"
import path from "node:path"
import { fileURLToPath } from "node:url"
import { describe, expect, it } from "vitest"

/**
 * 路由完整性门禁（Phase 3.6 故障 3 的回归锁）
 *
 * 历史故障：/lightnovel 的页面组件写好了、侧栏与顶栏入口也链上了，但
 * `src/routes/_main/lightnovel/` 路由定义从未创建 —— 生成器忠实工作，
 * routeTree.gen.ts 里没有它，客户端路由于是给出 404（后端对未知路径返 SPA 回落 200）。
 *
 * 本测试断言：导航中硬编码的每一个站内 href，都必须能在路由树里找到对应路由。
 *
 * 期望值来源（独立来源，不用被测代码重算）：
 *   routeTree.gen.ts 中由 TanStack Router 生成的真实路由 path/id 字面量。
 */

const projectRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "../../..")

const ROUTE_TREE_FILE = path.join(projectRoot, "src/routeTree.gen.ts")

/** 承载站内导航入口的文件（主侧栏 / 顶栏 / 离线侧栏 / 离线顶栏） */
const NAV_FILES = [
    path.join(projectRoot, "src/app/(main)/_features/navigation/main-sidebar.tsx"),
    path.join(projectRoot, "src/app/(main)/_features/navigation/top-menu.tsx"),
    path.join(projectRoot, "src/app/(main)/_features/navigation/offline-sidebar.tsx"),
    path.join(projectRoot, "src/app/(main)/_features/offline/_components/offline-top-menu.tsx"),
]

/** 归一化为可比较的 URL 路径：去尾斜杠，根路径保留 "/"，空串视作根 */
function normalizeUrlPath(raw: string): string {
    if (raw === "" || raw === "/") return "/"
    const withoutTrailingSlash = raw.endsWith("/") ? raw.slice(0, -1) : raw
    return withoutTrailingSlash.startsWith("/") ? withoutTrailingSlash : `/${withoutTrailingSlash}`
}

/** 从路由树生成产物里取出全部真实路由路径（path 与 id 均为路由字面量） */
function readRoutePaths(): Set<string> {
    const source = fs.readFileSync(ROUTE_TREE_FILE, "utf8")
    const paths = new Set<string>()
    const pattern = /^\s*(?:id|path):\s*(["'])([^"']*)\1,?\s*$/gm

    for (const match of source.matchAll(pattern)) {
        paths.add(normalizeUrlPath(match[2]))
    }

    return paths
}

/** 从导航源码里取出硬编码的站内 href 字面量（`href: "/x"`） */
function readNavigationHrefs(): string[] {
    const hrefs = new Set<string>()

    for (const file of NAV_FILES) {
        const source = fs.readFileSync(file, "utf8")
        for (const match of source.matchAll(/\bhref:\s*(["'])(\/[^"']*)\1/g)) {
            const raw = match[2]
            // 跳过协议相对链接、带参数的动态路由（导航不硬编码这些）
            if (raw.startsWith("//")) continue
            if (raw.includes(":")) continue
            hrefs.add(normalizeUrlPath(raw))
        }
    }

    return [...hrefs].sort()
}

const routePaths = readRoutePaths()
const navigationHrefs = readNavigationHrefs()

describe("route tree integrity", () => {
    it("extracts navigation hrefs from the sidebar and top menu", () => {
        // 防止导航文件被移动/改写后测试静默空转
        expect(navigationHrefs.length).toBeGreaterThan(0)
    })

    it("extracts real route paths from the generated route tree", () => {
        // 已知会在路由树中的站内路径（来自诊断报告核对过的既有路由）
        expect(routePaths.has("/asmr")).toBe(true)
        expect(routePaths.has("/schedule")).toBe(true)
        expect(routePaths.has("/")).toBe(true)
    })

    it("has a route for every in-app navigation href", () => {
        const missing = navigationHrefs.filter(href => !routePaths.has(href))

        expect(missing).toEqual([])
    })
})
