import type { Dictionary } from "../../types"

import { commonDictionary } from "./common"
import { libraryDictionary } from "./library"
import { navigationDictionary } from "./navigation"
import { playerDictionary } from "./player"
import { searchDictionary } from "./search"
import { settingsDictionary } from "./settings"

/** 分模块字典清单（合并顺序即文档展示顺序） */
const moduleDictionaries: readonly [string, Dictionary][] = [
    ["common", commonDictionary],
    ["navigation", navigationDictionary],
    ["library", libraryDictionary],
    ["search", searchDictionary],
    ["settings", settingsDictionary],
    ["player", playerDictionary],
]

/** 合并后的 zh-CN 单一字典 */
export const zhCN = {
    ...commonDictionary,
    ...navigationDictionary,
    ...libraryDictionary,
    ...searchDictionary,
    ...settingsDictionary,
    ...playerDictionary,
} satisfies Dictionary

/** 全部合法词条 key */
export type TranslationKey = keyof typeof zhCN

// 开发期重复 key 检测（构建期零成本，仅警告不抛错）
const expectedKeyCount = moduleDictionaries.reduce(
    (total, [, dictionary]) => total + Object.keys(dictionary).length,
    0,
)
const actualKeyCount = Object.keys(zhCN).length

if (actualKeyCount !== expectedKeyCount) {
    const seen = new Set<string>()
    const duplicates = new Set<string>()
    for (const [moduleName, dictionary] of moduleDictionaries) {
        for (const key of Object.keys(dictionary)) {
            if (seen.has(key)) {
                duplicates.add(`${key} (${moduleName})`)
            } else {
                seen.add(key)
            }
        }
    }
    console.warn(
        `[i18n] 检测到重复词条 key：期望 ${expectedKeyCount} 条，实际 ${actualKeyCount} 条。冲突项：${Array.from(duplicates).join(", ")}`,
    )
}
