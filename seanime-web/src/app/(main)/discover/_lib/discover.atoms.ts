import { atom } from "jotai"

export const __discord_pageTypeAtom = atom<"anime" | "schedule" | "manga" | "novel" | "asmr">("anime")

// 契约 03.9c：探索页 18+ 开关（仅 EnableAdultContent 开启时可见，请求带 isAdult: true）
export const __discover_isAdultAtom = atom(false)
