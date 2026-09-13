import Page from "@/app/(main)/asmr/page"
import { createLazyFileRoute } from "@tanstack/react-router"

export const Route = createLazyFileRoute("/_main/asmr/")({
    component: Page,
})
