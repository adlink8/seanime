import { createFileRoute } from "@tanstack/react-router"
import { z } from "zod"

const searchSchema = z.object({
    sorting: z.string().optional(),
    genre: z.string().optional(),
    status: z.string().optional(),
    format: z.string().optional(),
    season: z.string().optional(),
    year: z.coerce.number().optional(),
    type: z.string().optional(),
    // Phase 3.9d additive：Bangumi 标签直达（search/page.tsx 已消费该参数，
    // 未声明时 zod 会在导航时剥掉未知键）
    tags: z.string().optional(),
})

export const Route = createFileRoute("/_main/search/")({
    validateSearch: searchSchema,
})
