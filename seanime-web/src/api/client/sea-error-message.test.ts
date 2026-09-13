import { describe, expect, it } from "vitest"
import { handleSeaError } from "./sea-error-message"

/**
 * 服务端错误文案的行为契约锁。
 *
 * 历史故障：后端 panic 后 echo 返回 500 {"message":"Internal Server Error"}，
 * 响应体没有 .error 字段，前端一律兜底成字面量 "Unknown error" —— 用户看到红色 toast
 * 却拿不到任何可定位的信息（见 03.6-DIAGNOSIS.md Q1）。
 */
describe("handleSeaError", () => {
    describe("without an .error field (fallback path)", () => {
        it("reports the HTTP status together with the server message", () => {
            expect(handleSeaError({ message: "Internal Server Error" }, 500)).toBe("HTTP 500: Internal Server Error")
        })

        it("falls back to the bare HTTP status when there is no message", () => {
            expect(handleSeaError({}, 502)).toBe("HTTP 502")
        })

        it("uses the server message when the status is unknown", () => {
            expect(handleSeaError({ message: "oops" })).toBe("Server error: oops")
        })

        it("keeps the generic label only when nothing at all is known", () => {
            expect(handleSeaError({})).toBe("Unknown error")
        })
    })

    describe("with an .error field (existing branches preserved)", () => {
        it("unwraps a GraphQL error payload", () => {
            expect(handleSeaError({ error: '{"graphqlErrors":[{"message":"Bad request"}]}' }, 200))
                .toBe("AniList error: Bad request")
        })

        it("stays silent for cache misses", () => {
            expect(handleSeaError({ error: "no cached data" }, 500)).toBe("")
        })

        it("labels rate limiting", () => {
            expect(handleSeaError({ error: "Too many requests" }, 429))
                .toBe("AniList: Too many requests, please wait a moment and try again.")
        })

        it("passes a plain string body through", () => {
            expect(handleSeaError("boom", 500)).toBe("Server Error: boom")
        })
    })
})
