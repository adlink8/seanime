import { describe, expect, it } from "vitest"
import {
    deriveAutoCompletedTrack,
    formatTime,
    normalizeFilepath,
    type AsmrPlaybackTrack,
} from "./asmr-playback"

// Phase 3.1c（Wave B / D8）纯函数 seam 行为契约锁。
// 历史故障：自动完听若不在前端做保守判定，会在手动暂停/切歌时误标「已听」，污染进度。
// 故 deriveAutoCompletedTrack 必须仅在 running→stopped 跃迁且比例 >= 0.9 且 filepath 归一化匹配时触发。

const tracks: AsmrPlaybackTrack[] = [
    { path: "RJ/01/file.wav", localPath: "C:/RJ/01/file.wav" },
    { path: "RJ/02/other.wav", localPath: "C:/RJ/02/other.wav" },
]

describe("normalizeFilepath", () => {
    it("converts backslashes to forward slashes and lowercases", () => {
        expect(normalizeFilepath("C:\\RJ\\vob\\TRACK01.WAV")).toBe("c:/rj/vob/track01.wav")
    })
})

describe("formatTime", () => {
    it("formats whole minutes as m:ss", () => {
        expect(formatTime(65)).toBe("1:05")
    })

    it("zero-pads seconds", () => {
        expect(formatTime(5)).toBe("0:05")
    })

    it("treats negative as 0:00", () => {
        expect(formatTime(-1)).toBe("0:00")
    })

    it("treats NaN / Infinity as 0:00", () => {
        expect(formatTime(NaN)).toBe("0:00")
        expect(formatTime(Infinity)).toBe("0:00")
    })
})

describe("deriveAutoCompletedTrack", () => {
    it("returns the track path when naturally finished (ratio >= 0.9)", () => {
        const result = deriveAutoCompletedTrack({
            prev: { running: true, currentTime: 95, duration: 100, filepath: "C:/RJ/01/file.wav" },
            next: { running: false, filepath: "C:/RJ/01/file.wav" },
            tracks,
            completed: new Set(),
        })
        expect(result).toBe("RJ/01/file.wav")
    })

    it("does not trigger when ratio is below 0.9", () => {
        const result = deriveAutoCompletedTrack({
            prev: { running: true, currentTime: 89, duration: 100, filepath: "C:/RJ/01/file.wav" },
            next: { running: false, filepath: "C:/RJ/01/file.wav" },
            tracks,
            completed: new Set(),
        })
        expect(result).toBeNull()
    })

    it("does not match when filepath differs", () => {
        const result = deriveAutoCompletedTrack({
            prev: { running: true, currentTime: 100, duration: 100, filepath: "C:/RJ/99/nope.wav" },
            next: { running: false, filepath: "C:/RJ/99/nope.wav" },
            tracks,
            completed: new Set(),
        })
        expect(result).toBeNull()
    })

    it("matches via normalization (backslash + case-insensitive) against localPath", () => {
        const result = deriveAutoCompletedTrack({
            prev: { running: true, currentTime: 100, duration: 100, filepath: "c:\\rj\\01\\FILE.WAV" },
            next: { running: false, filepath: "c:\\rj\\01\\FILE.WAV" },
            tracks,
            completed: new Set(),
        })
        expect(result).toBe("RJ/01/file.wav")
    })

    it("does not re-fire for an already completed track", () => {
        const result = deriveAutoCompletedTrack({
            prev: { running: true, currentTime: 100, duration: 100, filepath: "C:/RJ/01/file.wav" },
            next: { running: false, filepath: "C:/RJ/01/file.wav" },
            tracks,
            completed: new Set(["RJ/01/file.wav"]),
        })
        expect(result).toBeNull()
    })

    it("does nothing when prev is null", () => {
        const result = deriveAutoCompletedTrack({
            prev: null,
            next: { running: false, filepath: "C:/RJ/01/file.wav" },
            tracks,
            completed: new Set(),
        })
        expect(result).toBeNull()
    })

    it("does not trigger when still running (no running→stopped transition)", () => {
        const result = deriveAutoCompletedTrack({
            prev: { running: true, currentTime: 100, duration: 100, filepath: "C:/RJ/01/file.wav" },
            next: { running: true, filepath: "C:/RJ/01/file.wav" },
            tracks,
            completed: new Set(),
        })
        expect(result).toBeNull()
    })
})
