import { createFileRoute } from "@tanstack/react-router"
import { t } from "@/lib/i18n"

export const Route = createFileRoute("/_main/error-test")({
    component: ErrorTest,
})

function ErrorTest() {
    throw new Error("This is a test error")

    return (
        <div className="p-4">
            <h1>{t("misc.error.test_page")}</h1>
        </div>
    )
}
