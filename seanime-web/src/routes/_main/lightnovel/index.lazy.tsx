import Page from "@/app/(main)/lightnovel/page"
import { createLazyFileRoute } from "@tanstack/react-router"

export const Route = createLazyFileRoute("/_main/lightnovel/")({
    component: Page,
})
