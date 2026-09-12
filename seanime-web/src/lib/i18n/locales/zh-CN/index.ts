import type { Dictionary } from "../../types"

import { commonDictionary } from "./common"
import { discoverDictionary } from "./discover"
import { homeDictionary } from "./home"
import { libraryDictionary } from "./library"
import { mediaDictionary } from "./media"
import { navigationDictionary } from "./navigation"
import { playerDictionary } from "./player"
import { scheduleDictionary } from "./schedule"
import { mangaDictionary } from "./manga"
import { torrentDictionary } from "./torrent"
import { autodownloaderDictionary } from "./autodownloader"
import { extensionsDictionary } from "./extensions"
import { miscDictionary } from "./misc"
import { searchDictionary } from "./search"
import { entryDictionary } from "./entry"
import { mpvDictionary } from "./mpv"
import { settingsDictionary } from "./settings"
import { settingsMediaDictionary } from "./settings-media"

/** 分模块字典清单（合并顺序即文档展示顺序） */
const moduleDictionaries: readonly [string, Dictionary][] = [
    ["common", commonDictionary],
    ["navigation", navigationDictionary],
    ["library", libraryDictionary],
    ["search", searchDictionary],
    ["settings", settingsDictionary],
    ["player", playerDictionary],
    ["home", homeDictionary],
    ["media", mediaDictionary],
    ["settings-media", settingsMediaDictionary],
    ["entry", entryDictionary],
    ["mpv", mpvDictionary],
    ["discover", discoverDictionary],
    ["schedule", scheduleDictionary],
    ["manga", mangaDictionary],
    ["torrent", torrentDictionary],
    ["autodownloader", autodownloaderDictionary],
    ["extensions", extensionsDictionary],
    ["misc", miscDictionary],
]

/** 合并后的 zh-CN 单一字典 */
export const zhCN = {
    ...commonDictionary,
    ...navigationDictionary,
    ...libraryDictionary,
    ...searchDictionary,
    ...settingsDictionary,
    ...playerDictionary,
    ...homeDictionary,
    ...mediaDictionary,
    ...discoverDictionary,
    ...scheduleDictionary,
    ...mangaDictionary,
    ...settingsMediaDictionary,
    ...entryDictionary,
    ...mpvDictionary,
    ...torrentDictionary,
    ...autodownloaderDictionary,
    ...extensionsDictionary,
    ...miscDictionary,
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
