package media

// subject_adapter.go —— Bangumi Subject → 完整 Anime/Manga 实体转换（Wave A）。
//
// 推导规则（保守优先，错误映射比缺映射危害大）：
//
//	Status：
//	  - 结束日期已知且已过           → FINISHED
//	  - 开始日期未知或在未来         → NOT_YET_RELEASED
//	  - 其余（开始日期已过）         → RELEASING
//	  说明：Bangumi v0 subject 不直接给出结束日期，本包默认拿不到 end，
//	  因此纯 Subject 输入永远推导不出 FINISHED（保守给 RELEASING）；
//	  Wave B 平台层若从 infobox 等渠道取得结束日期，可用 statusFromDates 传入。
//
//	Season / SeasonYear：由开始日期月份推导（12 月归入次年 WINTER，
//	  与 AniList 的跨年冬季习惯一致）；日期未知时两者留空。
//
//	Format：
//	  - 动画(type=2)：platform/tag 命中「剧场版/映画/电影」→ MOVIE；
//	    platform 命中 OVA/ONA/Web → 对应值；其余保守给 TV。
//	  - 书籍(type=1)：platform 命中「漫画」→ MANGA，「小说/轻小说」→ NOVEL；
//	    其余保守给 BOOK（Bangumi 扩展值，见 enums.go）。
//
//	评分换算：Bangumi rating 0-10 → AniList meanScore 0-100（×10）。
//	NSFW 复用 IsAdult 字段。

import (
	"fmt"
	"time"
)

// AnimeFromSubject 将 Bangumi Subject 镜像结构转换为完整 Anime 实体。
// s 为 nil 时返回 nil。
func AnimeFromSubject(s *Subject) *Anime {
	if s == nil {
		return nil
	}

	start, end := subjectDateRange(s)

	a := &Anime{
		ID:          s.ID,
		SiteURL:     strPtrOf(fmt.Sprintf("https://bgm.tv/subject/%d", s.ID), true),
		Status:      statusFromDates(start, end, time.Now()),
		Type:        mediaTypePtrOf(MediaTypeAnime),
		Episodes:    intPtrOf(s.Eps, s.Eps > 0),
		IsAdult:     boolPtrOf(s.NSFW),
		CountryOfOrigin: strPtrOf("JP", true),
		Description: strPtrOf(s.Summary, s.Summary != ""),
		Title:       titleFromSubject(s),
		CoverImage:  coverImageFromSubject(s),
		NameCN:      s.NameCN,
	}

	a.StartDate = fuzzyDateFromTime(start)
	if end != nil {
		a.EndDate = fuzzyDateFromTime(end)
	}
	if start != nil {
		season, year := seasonFromDate(*start)
		a.Season = &season
		a.SeasonYear = &year
	}

	a.Format = animeFormatFromSubject(s)

	// 评分：0-10 → 0-100（×10）
	if s.Rating != nil && s.Rating.Score > 0 {
		score := int(s.Rating.Score*10 + 0.5)
		a.MeanScore = &score
	}

	// 别名：原日文名与中文名不同时，日文名进 synonyms 以便匹配层复用
	if s.Name != "" && s.NameCN != "" && s.Name != s.NameCN {
		name := s.Name
		a.Synonyms = []*string{&name}
	}

	return a
}

// MangaFromSubject 将 Bangumi Subject 镜像结构转换为完整 Manga 实体。
// 书籍类目（漫画/轻小说/音声）统一走 Manga 路径。
func MangaFromSubject(s *Subject) *Manga {
	if s == nil {
		return nil
	}

	start, end := subjectDateRange(s)

	m := &Manga{
		ID:          s.ID,
		SiteURL:     strPtrOf(fmt.Sprintf("https://bgm.tv/subject/%d", s.ID), true),
		Status:      statusFromDates(start, end, time.Now()),
		Type:        mediaTypePtrOf(MediaTypeManga),
		Chapters:    intPtrOf(s.Eps, s.Eps > 0),
		IsAdult:     boolPtrOf(s.NSFW),
		CountryOfOrigin: strPtrOf("JP", true),
		Description: strPtrOf(s.Summary, s.Summary != ""),
		Title:       titleFromSubject(s),
		CoverImage:  coverImageFromSubject(s),
		NameCN:      s.NameCN,
	}

	m.StartDate = fuzzyDateFromTime(start)
	if end != nil {
		m.EndDate = fuzzyDateFromTime(end)
	}

	m.Format = mangaFormatFromSubject(s)

	if s.Rating != nil && s.Rating.Score > 0 {
		score := int(s.Rating.Score*10 + 0.5)
		m.MeanScore = &score
	}

	if s.Name != "" && s.NameCN != "" && s.Name != s.NameCN {
		name := s.Name
		m.Synonyms = []*string{&name}
	}

	return m
}

// titleFromSubject 标题映射：
//   - English  = subject.Name（Bangumi 无英文名字段，保证下游 GetTitleSafe 等逻辑可用）
//   - Native   = subject.Name
//   - Chinese  = subject.NameCN（Bangumi 扩展）
//   - UserPreferred 优先中文名（换锚的核心收益），缺失时回退原名
func titleFromSubject(s *Subject) *Title {
	name := s.Name
	nameCN := s.NameCN

	t := &Title{
		English: strPtrOf(name, name != ""),
		Native:  strPtrOf(name, name != ""),
		Chinese: strPtrOf(nameCN, nameCN != ""),
	}

	preferred := nameCN
	if preferred == "" {
		preferred = name
	}
	t.UserPreferred = strPtrOf(preferred, preferred != "")

	return t
}

// coverImageFromSubject 封面映射：Bangumi large → extraLarge/large，medium → medium。
func coverImageFromSubject(s *Subject) *CoverImage {
	if s.Images == nil {
		return nil
	}
	ci := &CoverImage{}
	if s.Images.Large != "" {
		ci.ExtraLarge = strPtrOf(s.Images.Large, true)
		ci.Large = strPtrOf(s.Images.Large, true)
	}
	if s.Images.Medium != "" {
		ci.Medium = strPtrOf(s.Images.Medium, true)
	}
	if ci.ExtraLarge == nil && ci.Large == nil && ci.Medium == nil {
		return nil
	}
	return ci
}

// statusFromDates 状态推导，规则见文件头注释。
// end 为 nil（Bangumi v0 常态）时不会返回 FINISHED。
func statusFromDates(start, end *time.Time, now time.Time) *MediaStatus {
	var st MediaStatus
	switch {
	case end != nil && end.Before(now):
		st = MediaStatusFinished
	case start == nil || start.After(now):
		st = MediaStatusNotYetReleased
	default:
		st = MediaStatusReleasing
	}
	return &st
}

// seasonFromDate 由日期推导季度与年份；12 月归入次年 WINTER。
// day 未定时按月初处理（仅月份参与推导）。
func seasonFromDate(d time.Time) (MediaSeason, int) {
	year := d.Year()
	season := MediaSeasonWinter
	switch m := d.Month(); {
	case m >= 3 && m <= 5:
		season = MediaSeasonSpring
	case m >= 6 && m <= 8:
		season = MediaSeasonSummer
	case m >= 9 && m <= 11:
		season = MediaSeasonFall
	case m == 12:
		// 12 月放送的作品按 AniList 习惯归入次年冬季档
		season = MediaSeasonWinter
		year++
	}
	return season, year
}

// animeFormatFromSubject 动画形式推导（规则见文件头注释）。
// 剧场版证据最具体，优先级最高。
func animeFormatFromSubject(s *Subject) *MediaFormat {
	platform := subjectPlatform(s)
	f := MediaFormatTv // 保守默认

	switch {
	case containsAny([]string{platform}, movieKeywords()) || hasTag(s, movieKeywords()):
		f = MediaFormatMovie
	case containsAny([]string{platform}, []string{"OVA"}):
		f = MediaFormatOva
	case containsAny([]string{platform}, []string{"ONA", "Web", "网络"}):
		f = MediaFormatOna
	}

	return &f
}

// mangaFormatFromSubject 书籍形式推导。
func mangaFormatFromSubject(s *Subject) *MediaFormat {
	platform := subjectPlatform(s)
	switch {
	case containsAny([]string{platform}, []string{"漫画", "Comic", "Manga"}):
		f := MediaFormatManga
		return &f
	case containsAny([]string{platform}, []string{"小说", "轻小说", "Novel", "Light"}) ||
		hasTag(s, []string{"小说", "轻小说"}):
		f := MediaFormatNovel
		return &f
	default:
		f := MediaFormatBook
		return &f
	}
}

// subjectDateRange 解析 Subject 的开始日期；Bangumi v0 无结束日期，
// 此处恒返回 end=nil（保留参数位以便 Wave B 扩展）。
func subjectDateRange(s *Subject) (start, end *time.Time) {
	start = parseBangumiDate(s.Date)
	return start, nil
}

// subjectPlatform 返回 Subject 的平台信息。
// Platform 为 Wave A 新增的 additive 字段（镜像 bangumi.Subject.Platform）。
func subjectPlatform(s *Subject) string { return s.Platform }

func movieKeywords() []string { return []string{"剧场版", "劇場版", "映画", "电影", "Movie"} }

func hasTag(s *Subject, keywords []string) bool {
	for _, t := range s.Tags {
		for _, k := range keywords {
			if t.Name == k {
				return true
			}
		}
	}
	return false
}

func containsAny(sources []string, keywords []string) bool {
	for _, src := range sources {
		if src == "" {
			continue
		}
		for _, k := range keywords {
			if len(k) > 0 && contains(src, k) {
				return true
			}
		}
	}
	return false
}

func contains(s, sub string) bool {
	return len(sub) > 0 && len(s) >= len(sub) && indexOf(s, sub) >= 0
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func fuzzyDateFromTime(t *time.Time) *FuzzyDate {
	if t == nil {
		return nil
	}
	y, m, d := t.Date()
	return &FuzzyDate{Year: &y, Month: (*int)(&m), Day: &d}
}

func strPtrOf(s string, cond bool) *string {
	if !cond {
		return nil
	}
	return &s
}

func intPtrOf(i int, cond bool) *int {
	if !cond {
		return nil
	}
	return &i
}

func boolPtrOf(b bool) *bool { return &b }

func mediaTypePtrOf(t MediaType) *MediaType { return &t }
