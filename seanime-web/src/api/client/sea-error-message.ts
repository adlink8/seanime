/**
 * 服务端错误响应 → 用户可见文案。
 *
 * 纯函数、无传输层依赖，便于单测（vitest 无法解析 `@/` 别名，故独立成模块）。
 */
export function handleSeaError(data: unknown, status?: number): string {
    if (typeof data === "string") return "Server Error: " + data

    const err = (data as any)?.error as string

    if (!err) {
        // 兜底：后端未提供 .error 时（例如 panic 后 echo 的 500 {"message":"Internal Server Error"}），
        // 退回到 HTTP 状态码 + message，避免一律显示无信息量的 "Unknown error"
        const message = typeof (data as any)?.message === "string" ? (data as any).message.trim() : ""
        if (status) {
            return message ? `HTTP ${status}: ${message}` : `HTTP ${status}`
        }
        return message ? `Server error: ${message}` : "Unknown error"
    }

    if (err.includes("Too many requests"))
        return "AniList: Too many requests, please wait a moment and try again."

    try {
        const graphqlErr = JSON.parse(err) as any
        console.log("AniList error", graphqlErr)
        if (graphqlErr.graphqlErrors && graphqlErr.graphqlErrors.length > 0 && !!graphqlErr.graphqlErrors[0]?.message) {
            return "AniList error: " + graphqlErr.graphqlErrors[0]?.message
        }
        return "AniList error"
    }
    catch (e) {
        if (err.includes("no cached data") || err.includes("cache lookup failed")) {
            return ""
        }
        return "Error: " + err
    }
}
