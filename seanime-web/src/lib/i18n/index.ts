/**
 * i18n 公共入口（barrel）
 *
 * 页面侧一律通过 `import { t } from "@/lib/i18n"` 取词。
 * 这里采用显式导出清单而非 `export *`：避免把内部实现（如词条私有结构、调试分支）
 * 意外提升为事实上的公共 API，公共面只保留明确承诺的符号。
 */

export { t } from "./t"
export { zhCN } from "./locales/zh-CN"
export type { TranslationKey } from "./locales/zh-CN"

export { DICTIONARY_MODULES } from "./types"
export type { Dictionary, DictionaryModule, TranslateParams } from "./types"
