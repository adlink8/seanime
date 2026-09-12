import { LuffyError } from "@/components/shared/luffy-error"
import { Button } from "@/components/ui/button"
import { t } from "@/lib/i18n"
import { Link } from "@tanstack/react-router"
import React from "react"

export function NotFound() {
    return (
        <div className="p-4 flex flex-col items-center justify-center h-full">
            <LuffyError title={t("common.error.page_not_found")}>
                <p className="text-[--muted] mb-4">
                    {t("common.error.page_not_found_desc")}
                </p>
                <Link to="/">
                    <Button>{t("common.action.go_home")}</Button>
                </Link>
            </LuffyError>
        </div>
    )
}
