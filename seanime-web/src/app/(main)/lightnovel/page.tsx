import { MediaEntryPageLoadingDisplay } from "@/app/(main)/_features/media/_components/media-entry-page-loading-display"
import { useHandleMangaCollection } from "@/app/(main)/manga/_lib/handle-manga-collection"
import { MangaLibraryView } from "@/app/(main)/manga/_screens/manga-library-view"
import { Manga_Collection } from "@/api/generated/types"
import React from "react"

/**
 * 轻小说栏（Phase 2 M2-05：纯元数据域——收藏/跟踪/评价，无下载/阅读）
 *
 * 数据复用 manga 收藏管道：Bangumi 书籍分区（subject_type=1）同时包含漫画与轻小说，
 * 后端 adapter 按 subject platform 推导 Format（小说/WEB→NOVEL、漫画→MANGA）。
 * 本页仅展示 Format=NOVEL 的条目，其余与漫画库行为一致。
 */
function filterNovelCollection(collection: Manga_Collection): Manga_Collection
function filterNovelCollection(collection: Manga_Collection | null | undefined): Manga_Collection | undefined
function filterNovelCollection(collection: Manga_Collection | null | undefined): Manga_Collection | undefined {
    if (!collection?.lists) return collection ?? undefined
    return {
        ...collection,
        lists: collection.lists?.map(l => ({
            ...l,
            entries: l.entries?.filter(e => e.media?.format === "NOVEL"),
        })),
    }
}

export default function Page() {
    const {
        mangaCollection,
        filteredMangaCollection,
        mangaCollectionLoading,
        storedFilters,
        storedProviders,
        mangaCollectionGenres,
        hasManga,
    } = useHandleMangaCollection()

    const novelCollection = filterNovelCollection(mangaCollection)
    const filteredNovelCollection = filterNovelCollection(filteredMangaCollection)
    const hasNovel = !!novelCollection?.lists?.some(l => !!l.entries?.length)

    if (!mangaCollection || mangaCollectionLoading) return <MediaEntryPageLoadingDisplay />

    // 加载守卫后 novelCollection 必然有值（mangaCollection 已定义）
    const collection = novelCollection as Manga_Collection

    return (
        <div
            data-lightnovel-page-container
        >
            <MangaLibraryView
                genres={mangaCollectionGenres}
                collection={collection}
                filteredCollection={filteredNovelCollection}
                storedProviders={storedProviders}
                hasManga={hasNovel}
                isMangaPage
            />
        </div>
    )
}
