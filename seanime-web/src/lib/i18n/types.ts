/**
 * i18n 类型契约
 *
 * 词条 key 采用点号小写 ASCII：<module>.<section>.<name>
 * 例如 navigation.sidebar.library、common.action.confirm。
 * key 中严禁出现中文或空格。
 */

/** 单个模块的词条表 */
export type Dictionary = Record<string, string>

/** 词条表的分模块粒度（INFRA-04 要求） */
export type DictionaryModule =
    | "common"
    | "navigation"
    | "library"
    | "search"
    | "settings"
    | "player"
    | "home"
    | "media"
    | "discover"
    | "schedule"
    | "settings-media"
    | "entry"
    | "mpv"

/** 模块清单（运行时可用，顺序即文档中的展示顺序） */
export const DICTIONARY_MODULES: readonly DictionaryModule[] = [
    "common",
    "navigation",
    "library",
    "search",
    "settings",
    "player",
    "home",
    "media",
    "discover",
    "schedule",
    "settings-media",
    "entry",
    "mpv",
] as const

/** t(key, params) 的插值参数，占位符写法为 {name} */
export type TranslateParams = Record<string, string | number>
