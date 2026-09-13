import type { AL_MangaCollection_MediaListCollection_Lists } from "@/api/generated/types"

/**
 * Pure filter (no React/jotai deps) that takes a manga collection's `lists`
 * array and returns an isomorphic structure containing only the entries whose
 * `media.format === "NOVEL"`.
 *
 * Guarantees:
 * - `null` / `undefined` input -> `[]`
 * - entries with `undefined` media are dropped (never crash)
 * - non-NOVEL formats (MANGA / ONE_SHOT / ...) are filtered out
 */
export function filterNovelCollection(
    lists: AL_MangaCollection_MediaListCollection_Lists[] | null | undefined,
): AL_MangaCollection_MediaListCollection_Lists[] {
    if (!lists) return []

    return lists
        .filter((list): list is AL_MangaCollection_MediaListCollection_Lists => !!list)
        .map(list => ({
            ...list,
            entries: (list.entries ?? []).filter(entry => entry?.media?.format === "NOVEL"),
        }))
}
