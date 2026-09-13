import { describe, expect, it } from "vitest"
import type { AL_MangaCollection_MediaListCollection_Lists } from "@/api/generated/types"
import { filterNovelCollection } from "./novel-collection-filter"

const makeEntry = (format: "NOVEL" | "MANGA" | "ONE_SHOT" | undefined, id: number) =>
    ({
        id,
        media: format === undefined ? undefined : { format },
    }) as AL_MangaCollection_MediaListCollection_Lists["entries"] extends (infer E)[]
        ? E
        : never

describe("filterNovelCollection", () => {
    it("① keeps only NOVEL entries when a list mixes formats", () => {
        const lists: AL_MangaCollection_MediaListCollection_Lists[] = [
            {
                name: "CURRENT",
                status: "CURRENT",
                entries: [makeEntry("NOVEL", 1), makeEntry("MANGA", 2), makeEntry("NOVEL", 3)],
            },
        ]

        const result = filterNovelCollection(lists)

        expect(result).toHaveLength(1)
        expect(result[0].entries).toHaveLength(2)
        expect(result[0].entries?.every(e => e?.media?.format === "NOVEL")).toBe(true)
        // list metadata is preserved (isomorphic structure)
        expect(result[0].name).toBe("CURRENT")
        expect(result[0].status).toBe("CURRENT")
    })

    it("② is safe with empty arrays and undefined/null input", () => {
        expect(filterNovelCollection([])).toEqual([])
        expect(filterNovelCollection(undefined)).toEqual([])
        expect(filterNovelCollection(null)).toEqual([])
    })

    it("③ does not crash when media is undefined", () => {
        const lists: AL_MangaCollection_MediaListCollection_Lists[] = [
            {
                name: "X",
                entries: [makeEntry(undefined, 1), makeEntry("NOVEL", 2)],
            },
        ]

        const result = filterNovelCollection(lists)

        expect(result).toHaveLength(1)
        expect(result[0].entries).toHaveLength(1)
        expect(result[0].entries?.[0]?.media?.format).toBe("NOVEL")
    })

    it("④ drops every non-NOVEL entry (MANGA / ONE_SHOT)", () => {
        const lists: AL_MangaCollection_MediaListCollection_Lists[] = [
            { name: "A", entries: [makeEntry("MANGA", 1), makeEntry("ONE_SHOT", 2)] },
            { name: "B", entries: [makeEntry("NOVEL", 3)] },
        ]

        const result = filterNovelCollection(lists)

        expect(result[0].entries).toHaveLength(0)
        expect(result[1].entries).toHaveLength(1)
        expect(
            result.every(l => (l.entries ?? []).every(e => e?.media?.format === "NOVEL")),
        ).toBe(true)
    })
})
