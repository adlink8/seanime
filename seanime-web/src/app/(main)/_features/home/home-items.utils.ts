import { Models_HomeItem, Nullish } from "@/api/generated/types"
import { ADVANCED_SEARCH_COUNTRIES_MANGA, ADVANCED_SEARCH_MEDIA_GENRES, GENRE_TRANSLATIONS } from "@/app/(main)/search/_lib/advanced-search-constants"
import { t } from "@/lib/i18n"

export const DEFAULT_HOME_ITEMS: Models_HomeItem[] = [
    {
        id: "anime-continue-watching",
        type: "anime-continue-watching",
        schemaVersion: 1,
    },
    {
        id: "anime-library",
        type: "anime-library",
        schemaVersion: 1,
        options: {
            statuses: ["CURRENT", "PAUSED", "PLANNING", "COMPLETED", "DROPPED"],
            layout: "grid",
        },
    },
]

export function isAnimeLibraryItemsOnly(items: Nullish<Models_HomeItem[]>) {
    if (!items) return true

    for (const item of items) {
        if (![
            "anime-continue-watching",
            "anime-library",
            "anime-continue-watching-header",
            "local-anime-library",
            "local-anime-library-stats",
            "library-upcoming-episodes",
        ].includes(item.type)) {
            return false
        }
    }
    return true
}

type HomeItemSchema = {
    name: string
    kind: ("row" | "header")[]
    options?: { label: string, name: string, type: string, options?: any[] }[]
    schemaVersion: number
    description?: string
}

const _carouselOptions = [
    {
        label: t("home.option.name"),
        type: "text",
        name: "name",
    },
    {
        label: t("home.option.sorting"),
        type: "select",
        name: "sorting",
        options: [
            {
                label: t("home.option.sort.popular"),
                value: "POPULARITY_DESC",
            },
            {
                label: t("home.option.sort.trending"),
                value: "TRENDING_DESC",
            },
            {
                label: t("home.option.sort.title_romaji_asc"),
                value: "TITLE_ROMAJI_ASC",
            },
            {
                label: t("home.option.sort.title_romaji_desc"),
                value: "TITLE_ROMAJI_DESC",
            },
            {
                label: t("home.option.sort.title_english_asc"),
                value: "TITLE_ENGLISH_ASC",
            },
            {
                label: t("home.option.sort.title_english_desc"),
                value: "TITLE_ENGLISH_DESC",
            },
            {
                label: t("home.option.sort.score_asc"),
                value: "SCORE",
            },
            {
                label: t("home.option.sort.score_desc"),
                value: "SCORE_DESC",
            },
        ],
    },
    {
        label: t("library.filter.status"),
        type: "multi-select",
        name: "status",
        options: [
            {
                label: t("search.status.releasing"),
                value: "RELEASING",
            },
            {
                label: t("search.status.finished"),
                value: "FINISHED",
            },
            {
                label: t("search.status.not_yet_released"),
                value: "NOT_YET_RELEASED",
            },
        ],
    },
    {
        label: t("library.filter.format"),
        type: "select",
        name: "format",
        options: [
            {
                label: "TV",
                value: "TV",
            },
            {
                label: t("library.stats.movies"),
                value: "MOVIE",
            },
            {
                label: "OVA",
                value: "OVA",
            },
            {
                label: "ONA",
                value: "ONA",
            },
            {
                label: t("library.stats.specials"),
                value: "SPECIAL",
            },
        ],
    },
    {
        label: t("search.filter.genre"),
        type: "multi-select",
        options: ADVANCED_SEARCH_MEDIA_GENRES.map(n => ({ value: n, label: GENRE_TRANSLATIONS[n] || n })),
        name: "genres",
    },
    {
        label: t("library.filter.season"),
        type: "select",
        name: "season",
        options: [
            { value: "WINTER", label: t("search.season.winter") },
            { value: "SPRING", label: t("search.season.spring") },
            { value: "SUMMER", label: t("search.season.summer") },
            { value: "FALL", label: t("search.season.fall") },
        ],
    },
    {
        label: t("library.filter.year"),
        type: "number",
        name: "year",
        min: 0,
        max: 2100,
    },
    {
        label: t("home.option.country_of_origin"),
        type: "select",
        name: "countryOfOrigin",
        options: ADVANCED_SEARCH_COUNTRIES_MANGA,
    },
]

export const HOME_ITEMS = {
    "centered-title": {
        name: t("home.item.centered_title.name"),
        kind: ["row"],
        schemaVersion: 1,
        description: t("home.item.centered_title.description"),
        options: [{
            label: t("home.option.text"),
            type: "text",
            name: "text",
        }],
    },
    "anime-continue-watching": {
        name: t("library.continue_watching.title"),
        kind: ["row", "header"],
        schemaVersion: 1,
        description: t("home.item.continue_watching.description"),
    },
    "anime-continue-watching-header": {
        name: t("home.item.continue_watching_header.name"),
        kind: ["header"],
        schemaVersion: 1,
        description: t("home.item.continue_watching_header.description"),
    },
    "anime-library": {
        name: t("home.item.anime_library.name"),
        kind: ["row"],
        schemaVersion: 2,
        description: t("home.item.anime_library.description"),
        options: [
            {
                label: t("home.option.watch_status"),
                name: "statuses",
                type: "multi-select",
                options: [
                    {
                        value: "CURRENT",
                        label: t("common.state.watching"),
                    },
                    {
                        value: "PAUSED",
                        label: t("common.state.paused"),
                    },
                    {
                        value: "PLANNING",
                        label: t("common.state.planning"),
                    },
                    {
                        value: "COMPLETED",
                        label: t("common.state.completed"),
                    },
                    {
                        value: "DROPPED",
                        label: t("common.state.dropped"),
                    },
                ],
            },
            {
                label: t("home.option.layout"),
                name: "layout",
                type: "select",
                options: [
                    {
                        label: t("home.option.layout.grid"),
                        value: "grid",
                    },
                    {
                        label: t("home.option.layout.carousel"),
                        value: "carousel",
                    },
                ],
            },
        ],
    },
    "my-lists": {
        name: t("navigation.item.lists"),
        kind: ["row"],
        schemaVersion: 1,
        description: t("home.item.my_lists.description"),
        options: [
            {
                label: t("library.filter.status"),
                name: "statuses",
                type: "multi-select",
                options: [
                    {
                        value: "CURRENT",
                        label: t("common.state.watching"),
                    },
                    {
                        value: "REPEATING",
                        label: t("common.state.rewatching"),
                    },
                    {
                        value: "PAUSED",
                        label: t("common.state.paused"),
                    },
                    {
                        value: "PLANNING",
                        label: t("common.state.planning"),
                    },
                    {
                        value: "COMPLETED",
                        label: t("common.state.completed"),
                    },
                    {
                        value: "DROPPED",
                        label: t("common.state.dropped"),
                    },
                ],
            },
            {
                label: t("home.option.layout"),
                name: "layout",
                type: "select",
                options: [
                    {
                        label: t("home.option.layout.grid"),
                        value: "grid",
                    },
                    {
                        label: t("home.option.layout.carousel"),
                        value: "carousel",
                    },
                ],
            },
            {
                label: t("home.option.type"),
                name: "type",
                type: "select",
                options: [
                    {
                        label: t("search.type.anime"),
                        value: "anime",
                    },
                    {
                        label: t("search.type.manga"),
                        value: "manga",
                    },
                ],
            },
            {
                label: t("home.option.custom_list_name"),
                type: "text",
                name: "customListName",
            },
        ],
    },
    "local-anime-library": {
        name: t("home.item.local_anime_library.name"),
        kind: ["row"],
        schemaVersion: 2,
        description: t("home.item.local_anime_library.description"),
        options: [
            {
                label: t("home.option.layout"),
                name: "layout",
                type: "select",
                options: [
                    {
                        label: t("home.option.layout.grid"),
                        value: "grid",
                    },
                    {
                        label: t("home.option.layout.carousel"),
                        value: "carousel",
                    },
                ],
            },
        ],
    },
    "library-upcoming-episodes": {
        name: t("home.item.library_upcoming_episodes.name"),
        kind: ["row"],
        schemaVersion: 1,
        description: t("home.item.library_upcoming_episodes.description"),
    },
    "aired-recently": {
        name: t("home.item.aired_recently.name"),
        kind: ["row"],
        schemaVersion: 1,
        description: t("home.item.aired_recently.description"),
    },
    "missed-sequels": {
        name: t("home.item.missed_sequels.name"),
        kind: ["row"],
        schemaVersion: 1,
        description: t("home.item.missed_sequels.description"),
    },
    "anime-schedule-calendar": {
        name: t("home.item.anime_schedule_calendar.name"),
        kind: ["row"],
        schemaVersion: 2,
        description: t("home.item.anime_schedule_calendar.description"),
        options: [
            {
                label: t("home.option.type"),
                name: "type",
                type: "select",
                options: [
                    {
                        label: t("navigation.item.lists"),
                        value: "my-lists",
                    },
                    {
                        label: t("home.option.type.global"),
                        value: "global",
                    },
                ],
            },
        ],
    },
    "local-anime-library-stats": {
        name: t("home.item.local_anime_library_stats.name"),
        kind: ["row"],
        schemaVersion: 1,
        description: t("home.item.local_anime_library_stats.description"),
    },
    "discover-header": {
        name: t("home.item.discover_header.name"),
        kind: ["header"],
        schemaVersion: 1,
        description: t("home.item.discover_header.description"),
    },
    "anime-carousel": {
        name: t("home.item.anime_carousel.name"),
        kind: ["row"],
        schemaVersion: 3,
        options: _carouselOptions,
        description: t("home.item.anime_carousel.description"),
    },
    "manga-carousel": {
        name: t("home.item.manga_carousel.name"),
        kind: ["row"],
        schemaVersion: 1,
        description: t("home.item.manga_carousel.description"),
        options: _carouselOptions.map(n => {
            if (n.name === "format") {
                return {
                    ...n,
                    options: [
                        {
                            label: t("search.type.manga"),
                            value: "MANGA",
                        },
                        {
                            label: t("home.option.format.one_shot"),
                            value: "ONE_SHOT",
                        },
                    ],
                }
            }
            return n
        }),
    },
    "manga-library": {
        name: t("home.item.manga_library.name"),
        kind: ["row", "header"],
        schemaVersion: 2,
        description: t("home.item.manga_library.description"),
        options: [
            {
                label: t("library.filter.status"),
                name: "statuses",
                type: "multi-select",
                options: [
                    {
                        value: "CURRENT",
                        label: t("common.state.reading"),
                    },
                    {
                        value: "PAUSED",
                        label: t("common.state.paused"),
                    },
                ],
            },
            {
                label: t("home.option.layout"),
                name: "layout",
                type: "select",
                options: [
                    {
                        label: t("home.option.layout.grid"),
                        value: "grid",
                    },
                    {
                        label: t("home.option.layout.carousel"),
                        value: "carousel",
                    },
                ],
            },
        ],
    },
} as Record<string, HomeItemSchema>

export const HOME_ITEM_IDS = Object.keys(HOME_ITEMS) as (keyof typeof HOME_ITEMS)[]

// export type HomeItemID = (keyof typeof HOME_ITEMS)
