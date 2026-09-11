import { describe, expect, it } from "vitest"

import { t } from "./t"
import { zhCN } from "./locales/zh-CN"

/**
 * t() 的行为契约锁：
 * - 命中 → 返回中文词条
 * - 未命中 → 原样返回 key，且不抛异常（INFRA-02 的核心，禁止口头声明）
 * - 插值 → 值中的 {name} 被 params.name 替换；缺参数时占位符原样保留
 *
 * 注：plan 的用例举例用 `{count}`，但正式词条表中现成的插值词条是
 * `player.external.scheme_placeholder`（占位符为 `{url}`）。t() 的替换分支是
 * 按 `\{(\w+)\}` 通用匹配的，`{url}` 与 `{count}` 走同一段代码，故直接复用真实
 * 词条进行断言——不为了测试通过而往 locales/zh-CN 下塞占位词条。
 */

describe("t", () => {
    it("已存在的 key 返回词条表中的中文", () => {
        const firstKey = Object.keys(zhCN)[0]!
        expect(t(firstKey)).toBe(zhCN[firstKey])
    })

    it("不存在的 key 原样返回 key 本身且不抛异常", () => {
        const missing = "__definitely_missing_key__"
        expect(() => t(missing)).not.toThrow()
        expect(t(missing)).toBe(missing)
    })

    it("对真实插值词条，{url} 被 params.url 替换（含数字参数被字符串化）", () => {
        const key = "player.external.scheme_placeholder"
        const template = zhCN[key]!
        // 前置断言：确认该词条确实是带占位符的模板
        expect(template).toContain("{url}")

        const rendered = t(key, { url: "https://example.com/stream.m3u8" })
        expect(rendered).toContain("https://example.com/stream.m3u8")
        expect(rendered).not.toContain("{url}")

        expect(t(key, { url: 123 })).toContain("123")
    })

    it("缺少插值参数时占位符原样保留且不抛异常", () => {
        const key = "player.external.scheme_placeholder"
        const template = zhCN[key]!

        // 完全不传 params
        expect(() => t(key)).not.toThrow()
        expect(t(key)).toBe(template)

        // 传了 params 但不含对应占位符名
        expect(() => t(key, { other: "x" })).not.toThrow()
        expect(t(key, { other: "x" })).toBe(template)
    })

    it("空字符串 key 返回空字符串且不抛异常", () => {
        expect(() => t("")).not.toThrow()
        expect(t("")).toBe("")
    })
})
