import { useAsmrLibrary } from "@/api/hooks/asmr.hooks"
import { AsmrLibraryCard } from "@/app/(main)/_features/asmr/_components/asmr-library-card"
import { LoadingSpinner } from "@/components/ui/loading-spinner"
import { t } from "@/lib/i18n"
import React from "react"
import { LuFolderOpen } from "react-icons/lu"

// Phase 3.1：本地音声库视图（/asmr 页 container）
// 调 GET /api/v1/asmr/library 渲染网格；空库显示引导。

export function AsmrLibraryView() {
    const { data, isLoading, isError, refetch } = useAsmrLibrary()

    if (isLoading) {
        return <div className="flex justify-center py-16"><LoadingSpinner /></div>
    }

    const entries = data?.entries ?? []

    if (isError) {
        return (
            <div className="py-16 text-center space-y-3" data-asmr-library-error>
                <p className="text-sm text-[--muted]">{t("asmr.library.load_error")}</p>
                <button
                    type="button"
                    onClick={() => refetch()}
                    className="rounded-[--radius-md] border border-[--border-color] px-3 py-1.5 text-sm hover:text-[--brand]"
                >
                    {t("asmr.library.retry")}
                </button>
            </div>
        )
    }

    if (!entries.length) {
        return (
            <div
                data-asmr-library-empty
                className="mx-auto max-w-xl py-16 px-4 text-center space-y-4"
            >
                <LuFolderOpen className="mx-auto text-4xl text-[--muted]" />
                <h3 className="text-lg font-semibold">{t("asmr.library.empty_title")}</h3>
                <p className="text-sm text-[--muted] leading-relaxed">
                    {t("asmr.library.empty_desc")}
                </p>
                <p className="text-sm text-[--muted]">
                    {t("asmr.library.empty_dir")}{" "}
                    <code className="rounded bg-[--background] px-1.5 py-0.5 text-[--brand]">
                        {data?.localDir || "$SEANIME_DATA_DIR/asmr-local"}
                    </code>
                </p>
                <button
                    type="button"
                    onClick={() => refetch()}
                    className="rounded-[--radius-md] border border-[--border-color] px-3 py-1.5 text-sm hover:text-[--brand]"
                >
                    {t("asmr.library.refresh")}
                </button>
            </div>
        )
    }

    return (
        <div data-asmr-library-grid className="flex flex-wrap gap-4">
            {entries.map(entry => (
                <AsmrLibraryCard key={entry.rjId} entry={entry} />
            ))}
        </div>
    )
}
