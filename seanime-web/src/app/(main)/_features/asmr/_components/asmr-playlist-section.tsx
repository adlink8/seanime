import { useAsmrPlaylistList, useAsmrPlaylistWorks } from "@/api/hooks/asmr.hooks"
import { AsmrWorkCard } from "@/app/(main)/_features/asmr/_components/asmr-work-card"
import { MediaEntryCardSkeleton } from "@/app/(main)/_features/media/_components/media-entry-card-skeleton"
import { cn } from "@/components/ui/core/styling"
import { LoadingSpinner } from "@/components/ui/loading-spinner"
import { t } from "@/lib/i18n"
import React, { useState } from "react"
import { LuListMusic } from "react-icons/lu"

// Phase 3.8：/asmr 页「云端播放列表」区块（契约 03.8-CONTRACT.md §3 Wave C）
// GET /asmr/playlist/list → 未配置凭据（configured:false）时整块不渲染；
// 点选列表 → GET /asmr/playlist/works?id=uuid → 复用 AsmrWorkCard 网格展示。

export function AsmrPlaylistSection() {
    const { data, isLoading, isError, refetch } = useAsmrPlaylistList()
    const [selectedId, setSelectedId] = useState<string | null>(null)

    // 未配置 asmr.one 账号：区块整体隐藏（契约 D6）
    if (data && !data.configured) return null

    const playlists = data?.playlists ?? []

    // 选中但列表刷新后消失：回落到第一个
    const activeId = selectedId && playlists.some(p => p.id === selectedId)
        ? selectedId
        : playlists[0]?.id ?? null

    return (
        <div data-asmr-playlist-section className="space-y-4">
            <h3 data-asmr-playlist-title className="flex items-center gap-2 text-lg font-semibold">
                <LuListMusic className="text-[--brand]" />
                {t("asmr.playlist.title")}
            </h3>

            {isLoading && (
                <div className="flex justify-center py-8"><LoadingSpinner /></div>
            )}

            {isError && (
                <div className="py-4 text-center space-y-2" data-asmr-playlist-error>
                    <p className="text-sm text-[--muted]">{t("asmr.playlist.load_error")}</p>
                    <button
                        type="button"
                        onClick={() => refetch()}
                        className="rounded-[--radius-md] border border-[--border-color] px-3 py-1.5 text-sm hover:text-[--brand]"
                    >
                        {t("asmr.library.retry")}
                    </button>
                </div>
            )}

            {!isLoading && !isError && !playlists.length && (
                <p data-asmr-playlist-empty className="text-sm text-[--muted]">{t("asmr.playlist.empty")}</p>
            )}

            {!!playlists.length && (
                <>
                    <div data-asmr-playlist-chips className="flex flex-wrap gap-2">
                        {playlists.map(playlist => (
                            <button
                                key={playlist.id}
                                type="button"
                                data-asmr-playlist-chip
                                data-asmr-playlist-id={playlist.id}
                                onClick={() => setSelectedId(playlist.id)}
                                className={cn(
                                    "inline-flex items-center gap-1.5 rounded-full border px-3 py-1.5 text-sm transition-colors",
                                    playlist.id === activeId
                                        ? "border-[--brand] bg-[--brand] text-white"
                                        : "border-[--border-color] hover:text-[--brand]",
                                )}
                            >
                                <span className="max-w-[160px] truncate">{playlist.name}</span>
                                <span className="text-xs opacity-80">{t("asmr.playlist.work_count", { count: playlist.worksCount })}</span>
                            </button>
                        ))}
                    </div>

                    {activeId && <AsmrPlaylistWorksGrid playlistId={activeId} />}
                </>
            )}
        </div>
    )
}

function AsmrPlaylistWorksGrid({ playlistId }: { playlistId: string }) {
    const { data, isLoading } = useAsmrPlaylistWorks(playlistId, true)
    const works = data?.works?.filter(Boolean) ?? []

    if (isLoading) {
        return (
            <div data-asmr-playlist-works-loading className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-4">
                {[...Array(6).keys()].map(v => <MediaEntryCardSkeleton key={v} />)}
            </div>
        )
    }

    if (!works.length) {
        return <p data-asmr-playlist-works-empty className="text-sm text-[--muted]">{t("asmr.playlist.empty")}</p>
    }

    return (
        <div data-asmr-playlist-works-grid className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-4">
            {works.map(work => (
                <AsmrWorkCard key={work.id} work={work} />
            ))}
        </div>
    )
}
