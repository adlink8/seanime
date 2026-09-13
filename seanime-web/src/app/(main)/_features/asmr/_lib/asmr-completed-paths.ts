import type { Asmr_Track } from "@/api/generated/types"

/**
 * 从服务端音轨树派生「已完听音轨路径」集合。
 *
 * 纯函数、无 React/传输层依赖，便于单测（本项目无 vitest.config，`@/` 别名在测试中
 * 解析不了，故测试用相对导入，与 src/api/client/sea-error-message.test.ts 先例一致）。
 *
 * 规则（契约 D8/D9）：
 * - 递归 folder 嵌套（音轨树多层 folder）；
 * - 只收 `path` 存在**且** `completed === true` 的节点；
 * - 在线音轨无 `path` → 天然不入集合（后端也只回填本地音轨）。
 */
export function deriveCompletedPaths(tracks: Asmr_Track[] | undefined): Set<string> {
    const completed = new Set<string>()

    const walk = (nodes: Asmr_Track[] | undefined) => {
        if (!nodes) return
        for (const node of nodes) {
            if (node.path && node.completed === true) {
                completed.add(node.path)
            }
            if (node.tracks?.length) {
                walk(node.tracks)
            }
        }
    }

    walk(tracks)
    return completed
}
