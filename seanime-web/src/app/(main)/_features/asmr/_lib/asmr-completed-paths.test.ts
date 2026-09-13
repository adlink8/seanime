import { describe, expect, it } from "vitest"
import { deriveCompletedPaths } from "./asmr-completed-paths"
import type { Asmr_Track } from "../../../../../api/generated/types"

/**
 * 完听勾选态派生的行为契约锁（Phase 3.2 / D8-D9）。
 *
 * 历史故障：弹窗勾选态是组件本地 useState<Set<string>>，初始恒为空且从不读服务端
 * → 刷新页面即全部丢失（关闭再打开同一卡片看着"保持"只是因为组件没卸载）。
 * 本纯函数把服务端 detail.tracks 里的逐轨 completed 派生为合集，供 effect 重置本地态。
 */
describe("deriveCompletedPaths", () => {
    it("returns an empty set for undefined tracks", () => {
        expect(deriveCompletedPaths(undefined).size).toBe(0)
    })

    it("collects completed leaf paths from a flat track tree", () => {
        const tracks: Array<Asmr_Track> = [
            { title: "01.mp3", type: "audio", path: "01.mp3", completed: true },
            { title: "02.mp3", type: "audio", path: "02.mp3", completed: false },
            { title: "03.mp3", type: "audio", path: "03.mp3", completed: true },
        ]

        const result = deriveCompletedPaths(tracks)

        expect(result.size).toBe(2)
        expect(result.has("01.mp3")).toBe(true)
        expect(result.has("02.mp3")).toBe(false)
        expect(result.has("03.mp3")).toBe(true)
    })

    it("recurses into multi-level nested folders", () => {
        const tracks: Array<Asmr_Track> = [
            {
                title: "disc1",
                type: "folder",
                tracks: [
                    { title: "a.mp3", type: "audio", path: "disc1/a.mp3", completed: true },
                    {
                        title: "bonus",
                        type: "folder",
                        tracks: [
                            { title: "b.mp3", type: "audio", path: "disc1/bonus/b.mp3", completed: true },
                            { title: "c.mp3", type: "audio", path: "disc1/bonus/c.mp3", completed: false },
                        ],
                    },
                ],
            },
        ]

        const result = deriveCompletedPaths(tracks)

        expect(result.size).toBe(2)
        expect(result.has("disc1/a.mp3")).toBe(true)
        expect(result.has("disc1/bonus/b.mp3")).toBe(true)
        expect(result.has("disc1/bonus/c.mp3")).toBe(false)
    })

    it("excludes nodes whose completed flag is missing or false", () => {
        const tracks: Array<Asmr_Track> = [
            { title: "no-flag.mp3", type: "audio", path: "no-flag.mp3" },
            { title: "explicit-false.mp3", type: "audio", path: "explicit-false.mp3", completed: false },
        ]

        expect(deriveCompletedPaths(tracks).size).toBe(0)
    })

    it("excludes online tracks that have no path", () => {
        const tracks: Array<Asmr_Track> = [
            { title: "online", type: "audio", mediaStreamUrl: "https://example.com/a.mp3", completed: true },
        ]

        expect(deriveCompletedPaths(tracks).size).toBe(0)
    })
})
