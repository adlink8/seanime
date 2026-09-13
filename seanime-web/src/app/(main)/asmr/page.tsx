import { AsmrLibraryView } from "@/app/(main)/asmr/_components/asmr-library-view"
import { t } from "@/lib/i18n"
import React from "react"

// Phase 3.1：本地音声库页（/asmr 路由，TanStack Router 文件式自动注册）

export default function Page() {
    return (
        <div data-asmr-page-container className="space-y-4 p-4">
            <div className="flex items-center justify-between">
                <h1 className="text-xl font-semibold" data-asmr-page-title>{t("asmr.page.title")}</h1>
            </div>
            <AsmrLibraryView />
        </div>
    )
}
