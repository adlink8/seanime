import { useServerStatus } from "@/app/(main)/_hooks/use-server-status"
import {
    ADVANCED_SEARCH_FORMATS_ANIME,
    ADVANCED_SEARCH_MEDIA_GENRES,
    ADVANCED_SEARCH_SEASONS,
    ADVANCED_SEARCH_SORTING,
    ADVANCED_SEARCH_SORTING_ASMR,
    ADVANCED_SEARCH_SORTING_MANGA,
    ADVANCED_SEARCH_SUBTITLE_ASMR,
    ADVANCED_SEARCH_TYPE,
    GENRE_TRANSLATIONS,
    SEASON_TRANSLATIONS,
    mapSortingToAsmrOrder,
} from "@/app/(main)/search/_lib/advanced-search-constants"
import { __advancedSearch_paramsAtom } from "@/app/(main)/search/_lib/advanced-search.atoms"
import { AppLayoutStack } from "@/components/ui/app-layout"
import { IconButton } from "@/components/ui/button"
import { Combobox } from "@/components/ui/combobox"
import { cn } from "@/components/ui/core/styling"
import { Select } from "@/components/ui/select"
import { Switch } from "@/components/ui/switch"
import { TextInput } from "@/components/ui/text-input"
import { useDebounce } from "@/hooks/use-debounce"
import { t } from "@/lib/i18n"
import { getYear } from "date-fns"
import { useAtom } from "jotai/react"
import React, { useState } from "react"
import { BiTrash } from "react-icons/bi"
import { FaRegStar, FaSortAmountDown } from "react-icons/fa"
import { FiSearch } from "react-icons/fi"
import { LuCalendar, LuLeaf } from "react-icons/lu"
import { MdPersonalVideo } from "react-icons/md"
import { TbSwords, TbTagsFilled } from "react-icons/tb"
import { useMount } from "react-use"
import { useUpdateEffect } from "react-use"

export function AdvancedSearchOptions() {

    const serverStatus = useServerStatus()
    const [params, setParams] = useAtom(__advancedSearch_paramsAtom)

    const highlightTrash = React.useMemo(() => {
        return !(!params.title?.length &&
            (params.sorting === null || params.sorting?.[0] === "SCORE_DESC") &&
            (params.genre === null || !params.genre.length) &&
            (params.tags === null || !params.tags.length) &&
            (params.status === null || !params.status.length) &&
            params.format === null && params.season === null && params.year === null && params.isAdult === false && params.minScore === null &&
            (params.countryOfOrigin === null || params.type === "anime"))
    }, [params])

    return (
        <AppLayoutStack data-advanced-search-options-container className="px-4 xl:px-0 space-y-3">
            <div data-advanced-search-options-header className="flex flex-col md:flex-row xl:flex-col gap-4 lg:gap-3">
                <TitleInput />
                <Select
                    className="w-full"
                    // Phase 3.7 D1：类型控件此前既无 label 也无 placeholder，补可见标签以消除歧义
                    // （相邻控件用 placeholder 提示、label 则以 fieldLabelClass="hidden" 隐藏，本控件两者皆无，故保留标签可见）
                    label={t("search.filter.type")}
                    options={ADVANCED_SEARCH_TYPE}
                    value={params.type}
                    onValueChange={v => setParams(draft => {
                        draft.type = v as "anime" | "manga" | "novel" | "asmr"
                        return
                    })}
                />
                {params.type === "asmr"
                    ? (
                        // Phase 2.5：asmr 类型下复用排序控件，映射为 asmr.one 的 order 参数
                        <Select
                            leftAddon={
                                <FaSortAmountDown className={cn((params.sorting !== null && params.sorting?.[0] !== "SCORE_DESC") && "text-indigo-300 font-bold text-xl")} />}
                            className="w-full"
                            options={ADVANCED_SEARCH_SORTING_ASMR}
                            value={mapSortingToAsmrOrder(params.sorting?.[0])}
                            onValueChange={v => setParams(draft => {
                                // 将 asmr order 反查回 AL 排序值存储（dl→POPULARITY_DESC / rating→SCORE_DESC / dd / publish_date→START_DATE_DESC）
                                const al = v === "dl" ? "POPULARITY_DESC" : v === "dd" ? "START_DATE_DESC" : v === "publish_date" ? "START_DATE_DESC" : "SCORE_DESC"
                                draft.sorting = [al] as any
                                return
                            })}
                        />
                    )
                    : (
                        <Select
                            // label="Sorting"
                            leftAddon={
                                <FaSortAmountDown className={cn((params.sorting !== null && params.sorting?.[0] !== "SCORE_DESC") && "text-indigo-300 font-bold text-xl")} />}
                            className="w-full"
                            options={params.type === "anime" ? ADVANCED_SEARCH_SORTING : params.type === "novel" ? ADVANCED_SEARCH_SORTING : ADVANCED_SEARCH_SORTING_MANGA}
                            value={params.sorting?.[0] || "SCORE_DESC"}
                            onValueChange={v => setParams(draft => {
                                draft.sorting = [v] as any
                                return
                            })}
                        />
                    )}
            </div>
            {/* Phase 2.5：asmr 类型下过滤器简化，仅保留标题关键词（上方常驻）+ 排序（上）+ 字幕可选 */}
            {params.type === "asmr" ? (
                <div
                    data-advanced-search-options-content
                    className="grid grid-cols-2 md:grid-cols-3 xl:grid-cols-1 gap-4 lg:gap-3 items-end xl:items-start"
                >
                    <Select
                        leftAddon={<TbTagsFilled className={cn((params.asmrSubtitle !== null && !!params.asmrSubtitle) && "text-indigo-300 font-bold text-xl")} />}
                        label={t("search.filter.subtitle")} placeholder={t("search.filter.subtitle_any")} className="w-full"
                        options={ADVANCED_SEARCH_SUBTITLE_ASMR}
                        value={params.asmrSubtitle || ""}
                        onValueChange={v => setParams(draft => {
                            draft.asmrSubtitle = v || null
                            return
                        })}
                        fieldLabelClass="hidden"
                    />
                </div>
            ) : (
                <div
                    data-advanced-search-options-content
                    className="grid grid-cols-2 md:grid-cols-3 xl:grid-cols-1 gap-4 lg:gap-3 items-end xl:items-start"
                >
                <Combobox
                    multiple
                    leftAddon={<TbSwords className={cn((params.genre !== null && !!params.genre.length) && "text-indigo-300 font-bold text-xl")} />}
                    emptyMessage={t("search.filter.genre_empty")}
                    label={t("search.filter.genre")} placeholder={t("search.filter.genre_all")} className="w-full"
                    options={ADVANCED_SEARCH_MEDIA_GENRES.map(genre => ({ value: genre, label: GENRE_TRANSLATIONS[genre] || genre, textValue: `${genre} ${GENRE_TRANSLATIONS[genre] || ""}` }))}
                    value={params.genre ? params.genre : []}
                    onValueChange={v => setParams(draft => {
                        draft.genre = v
                        return
                    })}
                    fieldLabelClass="hidden"
                />
                {params.type === "anime" && <Select
                    leftAddon={<MdPersonalVideo className={cn((params.format !== null && !!params.format) && "text-indigo-300 font-bold text-xl")} />}
                    label={t("search.filter.format")} placeholder={t("search.filter.format_all")} className="w-full"
                    // Phase 3.6b D4：上游 meta_tags 仅 TV / WEB(ONA) / OVA 有效，MOVIE / TV_SHORT / SPECIAL 已移除
                    options={ADVANCED_SEARCH_FORMATS_ANIME}
                    value={params.format || ""}
                    onValueChange={v => setParams(draft => {
                        draft.format = v as any
                        return
                    })}
                    fieldLabelClass="hidden"
                />}
                {/* Phase 3.6b D3：漫画的「国家/地区」筛选已移除（上游无 country 字段，platform 过滤被实测证伪） */}
                {/* Phase 3.6b D5：漫画的「格式」筛选已移除（meta_tags 对书籍分区恒 0，上游无能力） */}
                {(params.type === "anime" || params.type === "novel") && <Select
                    leftAddon={<LuLeaf className={cn((params.season !== null && !!params.season) && "text-indigo-300 font-bold text-xl")} />}
                    placeholder={t("search.filter.season_all")} className="w-full"
                    options={ADVANCED_SEARCH_SEASONS.map(season => ({ value: season.toUpperCase(), label: SEASON_TRANSLATIONS[season] || season }))}
                    value={params.season || ""}
                    onValueChange={v => setParams(draft => {
                        draft.season = v as any
                        return
                    })}
                    fieldLabelClass="hidden"
                />}
                <Select
                    leftAddon={<LuCalendar className={cn((params.year !== null && !!params.year) && "text-indigo-300 font-bold text-xl")} />}
                    label={t("search.filter.year")} placeholder={t("search.filter.year_any")} className="w-full"
                    options={[...Array(70)].map((v, idx) => getYear(new Date()) - idx + 2).map(year => ({
                        value: String(year),
                        label: String(year),
                    }))}
                    value={params.year || ""}
                    onValueChange={v => setParams(draft => {
                        draft.year = v as any
                        return
                    })}
                    fieldLabelClass="hidden"
                />
                {/* Phase 3.6b D2：status 筛选已移除（上游静默忽略该字段） */}
                <Select
                    leftAddon={<FaRegStar className={cn((params.minScore !== null && !!params.minScore) && "text-indigo-300 font-bold text-xl")} />}
                    placeholder={t("search.filter.score_all")} className="w-full"
                    // Phase 3.6b D6：与后端一致按 0–100 下发（后端 ÷10 后拼 rating 过滤）
                    options={[...Array(10)].map((v, idx) => (idx + 1) * 10).map(score => ({
                        value: String(score),
                        label: `≥ ${score}`,
                    }))}
                    value={params.minScore || ""}
                    onValueChange={v => setParams(draft => {
                        draft.minScore = v as any
                        return
                    })}
                />
                {serverStatus?.settings?.anilist?.enableAdultContent && <Switch
                    label={t("search.filter.adult")}
                    value={params.isAdult}
                    onValueChange={v => setParams(draft => {
                        draft.isAdult = v
                        return
                    })}
                    fieldLabelClass="hidden"
                />}
                <IconButton
                    icon={<BiTrash />} intent={highlightTrash ? "alert" : "gray-subtle"} className="flex-none" onClick={() => {
                    setParams(prev => ({
                        ...prev,
                        active: true,
                        title: null,
                        sorting: null,
                        status: null,
                        genre: null,
                        tags: null,
                        format: null,
                        season: null,
                        year: null,
                        minScore: null,
                        countryOfOrigin: null,
                        // isAdult: false,
                    }))
                }}
                    disabled={!highlightTrash}
                />
                </div>
            )}

        </AppLayoutStack>
    )
}

function TitleInput() {
    const [inputValue, setInputValue] = useState("")
    const debouncedTitle = useDebounce(inputValue, 500)
    const [params, setParams] = useAtom(__advancedSearch_paramsAtom)

    const ref = React.useRef<HTMLInputElement | null>(null)

    useMount(() => {
        ref.current?.focus()
    })

    useUpdateEffect(() => {
        setParams(draft => {
            draft.title = debouncedTitle
            return
        })
    }, [debouncedTitle])

    useUpdateEffect(() => {
        setInputValue(params.title || "")
    }, [params.title])

    return (
        <TextInput
            ref={ref}
            leftIcon={<FiSearch />} placeholder={t("search.filter.title_placeholder")} className="w-full"
            value={inputValue}
            onValueChange={v => setInputValue(v)}
        />
    )
}
