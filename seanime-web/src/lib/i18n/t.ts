/**
 * 统一取词函数 t(key, params?)
 *
 * 设计要点（后续阶段依赖，勿擅改）：
 *
 * 1. 形参 key 用 `string` 而非 `TranslationKey` 联合类型。
 *    原因：i18n 是一个持续增量的基建，后续阶段会不断新增词条；若收窄为联合类型，
 *    则「传入一个尚未登记的 key」会变成编译错误，与第 2 点的运行时回退语义直接冲突。
 *    取词失败必须是运行时可观测的降级行为，而不是把整个构建卡死的类型错误。
 *
 * 2. 未命中 key 时原样返回 key 本身，不抛异常、不打日志。
 *    原因：叙事页面在词条未覆盖时应当显示原始 key 而不是白屏；同时未命中是设计内的
 *    正常状态（渐进式本地化），高频 console 噪音会淹没真正的错误。
 *
 * 3. 缺插值参数时保留 `{name}` 占位符原样，不抛异常。
 *    原因：与第 2 点同理，缺参是渐进式覆盖的中间态，抛出会连带炸掉整棵组件树。
 *
 * 4. 本函数为纯函数、框架无关：不 import React、不依赖任何浏览器 API。
 *    原因：应用内 `_lib` 目录下的非组件模块（服务端工具、数据转换）同样需要取词。
 */

import { zhCN } from "./locales/zh-CN"
import type { Dictionary, TranslateParams } from "./types"

/**
 * 以索引签名视角读取词条表。
 * `zhCN` 经 `satisfies` 后类型仍是 193 个字面量 key 的精确对象类型，没有索引签名，
 * 无法用 `string` 下标访问；这里仅在模块内部做一次结构升级，不改变对外类型契约。
 */
const dictionary: Dictionary = zhCN

/** 匹配 {name} 形式的插值占位符，name 为单词字符 */
const PLACEHOLDER_PATTERN = /\{(\w+)\}/g

export function t(key: string, params?: TranslateParams): string {
    const template = dictionary[key]
    if (template === undefined) return key
    if (!params) return template

    return template.replace(
        PLACEHOLDER_PATTERN,
        (placeholder: string, name: string): string =>
            params[name] === undefined ? placeholder : String(params[name]),
    )
}
