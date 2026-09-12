package media

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------
// Subject → Anime / Manga 转换

// tvSubject 常规 TV 动画样本（已放送）
func tvSubject() *Subject {
	return &Subject{
		ID:       94607,
		Type:     2,
		Name:     "葬送のフリーレン",
		NameCN:   "葬送的芙莉莲",
		Summary:  "简介文本",
		Date:     "2023-09-29",
		Platform: "TV",
		Eps:      28,
		NSFW:     false,
		Images: &SubjectImages{
			Large:  "https://lain.bgm.tv/pic/cover/l/xx.jpg",
			Common: "https://lain.bgm.tv/pic/cover/c/xx.jpg",
			Medium: "https://lain.bgm.tv/pic/cover/m/xx.jpg",
		},
		Tags:   []SubjectTag{{Name: "奇幻", Count: 200}, {Name: "TV", Count: 100}},
		Rating: &SubjectRating{Score: 9.2, Rank: 5, Total: 1000},
	}
}

func TestAnimeFromSubject_TV(t *testing.T) {
	a := AnimeFromSubject(tvSubject())
	if a == nil {
		t.Fatal("转换结果为 nil")
	}

	if a.ID != 94607 {
		t.Errorf("ID = %d, want 94607", a.ID)
	}
	if a.Type == nil || *a.Type != MediaTypeAnime {
		t.Errorf("Type 应为 ANIME")
	}
	if a.Format == nil || *a.Format != MediaFormatTv {
		t.Errorf("常规 TV 动画 Format 应保守推导为 TV，实际 %v", a.Format)
	}
	if a.Status == nil || *a.Status != MediaStatusReleasing {
		t.Errorf("开始日期已过且无结束日期时应为 RELEASING，实际 %v", a.Status)
	}
	if a.Season == nil || *a.Season != MediaSeasonFall {
		t.Errorf("9 月放送应为 FALL，实际 %v", a.Season)
	}
	if a.SeasonYear == nil || *a.SeasonYear != 2023 {
		t.Errorf("SeasonYear = %v, want 2023", a.SeasonYear)
	}
	if a.Episodes == nil || *a.Episodes != 28 {
		t.Errorf("Episodes = %v, want 28", a.Episodes)
	}
	if a.NameCN != "葬送的芙莉莲" {
		t.Errorf("NameCN = %q", a.NameCN)
	}
	// 评分换算：9.2 × 10 = 92
	if a.MeanScore == nil || *a.MeanScore != 92 {
		t.Errorf("MeanScore 应为 92（Bangumi 9.2 ×10），实际 %v", a.MeanScore)
	}
	if a.IsAdult == nil || *a.IsAdult {
		t.Errorf("NSFW=false 应映射 IsAdult=false")
	}
	if a.SiteURL == nil || *a.SiteURL != "https://bgm.tv/subject/94607" {
		t.Errorf("SiteURL = %v", a.SiteURL)
	}

	// 标题映射：English=Native=原名，Chinese=中文名，UserPreferred 优先中文
	if a.Title == nil {
		t.Fatal("Title 为 nil")
	}
	if a.Title.English == nil || *a.Title.English != "葬送のフリーレン" {
		t.Errorf("Title.English 应为原日文名")
	}
	if a.Title.Chinese == nil || *a.Title.Chinese != "葬送的芙莉莲" {
		t.Errorf("Title.Chinese 应为中文名（Bangumi 扩展字段）")
	}
	if a.Title.UserPreferred == nil || *a.Title.UserPreferred != "葬送的芙莉莲" {
		t.Errorf("UserPreferred 应优先中文名")
	}

	// 封面：large → extraLarge/large
	if a.CoverImage == nil || a.CoverImage.ExtraLarge == nil {
		t.Errorf("CoverImage.ExtraLarge 应取 Bangumi large 图")
	}

	// 原名进 synonyms（中日名不同）
	if len(a.Synonyms) != 1 || *a.Synonyms[0] != "葬送のフリーレン" {
		t.Errorf("Synonyms 应包含日文原名")
	}
}

func TestAnimeFromSubject_NilSafe(t *testing.T) {
	if AnimeFromSubject(nil) != nil {
		t.Error("nil 输入应返回 nil")
	}
	var a *Anime
	if a.GetID() != 0 {
		t.Error("nil 接收者的 GetID 应返回零值")
	}
	// nil 接收者链式调用不 panic，且返回零值
	if v := a.GetTitle().GetEnglish(); v != nil {
		t.Error("nil 接收者的链式 Get* 应返回 nil")
	}
	if v := a.GetCoverImageSafe(); v != "" {
		t.Errorf("nil 接收者 GetCoverImageSafe 应返回空串，实际 %q", v)
	}
}

func TestAnimeFromSubject_Movie(t *testing.T) {
	s := tvSubject()
	s.Platform = "剧场版"
	a := AnimeFromSubject(s)
	if a.Format == nil || *a.Format != MediaFormatMovie {
		t.Errorf("platform=剧场版 应推导 MOVIE，实际 %v", a.Format)
	}

	s2 := tvSubject()
	s2.Platform = ""
	s2.Tags = append(s2.Tags, SubjectTag{Name: "剧场版"})
	a2 := AnimeFromSubject(s2)
	if a2.Format == nil || *a2.Format != MediaFormatMovie {
		t.Errorf("tag=剧场版 应推导 MOVIE，实际 %v", a2.Format)
	}
}

func TestAnimeFromSubject_NotYetReleased(t *testing.T) {
	s := tvSubject()
	s.Date = time.Now().AddDate(1, 0, 0).Format("2006-01-02")
	a := AnimeFromSubject(s)
	if a.Status == nil || *a.Status != MediaStatusNotYetReleased {
		t.Errorf("未来放送应为 NOT_YET_RELEASED，实际 %v", a.Status)
	}
}

func TestStatusFromDates(t *testing.T) {
	now := time.Now()
	past := now.AddDate(-1, 0, 0)
	future := now.AddDate(1, 0, 0)

	// 结束日期已过 → FINISHED（Wave B 从 infobox 等取得 end 时的路径）
	if st := statusFromDates(&past, &past, now); *st != MediaStatusFinished {
		t.Errorf("结束日期已过应为 FINISHED")
	}
	// 无结束日期、开始已过 → RELEASING（保守）
	if st := statusFromDates(&past, nil, now); *st != MediaStatusReleasing {
		t.Errorf("无结束日期应为 RELEASING")
	}
	// 开始在未来 → NOT_YET_RELEASED
	if st := statusFromDates(&future, nil, now); *st != MediaStatusNotYetReleased {
		t.Errorf("开始在未来应为 NOT_YET_RELEASED")
	}
	// 无日期 → NOT_YET_RELEASED
	if st := statusFromDates(nil, nil, now); *st != MediaStatusNotYetReleased {
		t.Errorf("无日期应为 NOT_YET_RELEASED")
	}
}

func TestSeasonFromDate(t *testing.T) {
	cases := []struct {
		month  time.Month
		want   MediaSeason
		year   int
		wantYr int
	}{
		{time.January, MediaSeasonWinter, 2023, 2023},
		{time.March, MediaSeasonSpring, 2023, 2023},
		{time.July, MediaSeasonSummer, 2023, 2023},
		{time.October, MediaSeasonFall, 2023, 2023},
		{time.December, MediaSeasonWinter, 2023, 2024}, // 12 月归次年冬季
	}
	for _, c := range cases {
		got, yr := seasonFromDate(time.Date(c.year, c.month, 15, 0, 0, 0, 0, time.UTC))
		if got != c.want || yr != c.wantYr {
			t.Errorf("seasonFromDate(%s) = (%s, %d), want (%s, %d)", c.month, got, yr, c.want, c.wantYr)
		}
	}
}

func TestMangaFromSubject(t *testing.T) {
	mk := func(platform string, tags ...string) *Subject {
		s := tvSubject()
		s.Type = 1
		s.Platform = platform
		s.Tags = nil
		for _, tg := range tags {
			s.Tags = append(s.Tags, SubjectTag{Name: tg})
		}
		return s
	}

	if m := MangaFromSubject(mk("漫画")); m.Format == nil || *m.Format != MediaFormatManga {
		t.Errorf("platform=漫画 应推导 MANGA，实际 %v", m.Format)
	}
	if m := MangaFromSubject(mk("轻小说")); m.Format == nil || *m.Format != MediaFormatNovel {
		t.Errorf("platform=轻小说 应推导 NOVEL，实际 %v", m.Format)
	}
	// 保守默认：BOOK（Bangumi 扩展值）
	m := MangaFromSubject(mk(""))
	if m.Format == nil || *m.Format != MediaFormatBook {
		t.Errorf("书籍默认应推导 BOOK，实际 %v", m.Format)
	}
	if m.Type == nil || *m.Type != MediaTypeManga {
		t.Errorf("Type 应为 MANGA")
	}
	if m.Chapters == nil || *m.Chapters != 28 {
		t.Errorf("Chapters 应取 Eps，实际 %v", m.Chapters)
	}
	if MangaFromSubject(nil) != nil {
		t.Error("nil 输入应返回 nil")
	}
}

// TestAnimeFromSubject_JSONShape 硬约束验证：JSON 序列化字段名与 anilist 一致
// （换锚不换形——前端序列化形状不能破），Bangumi 扩展字段为 additive。
func TestAnimeFromSubject_JSONShape(t *testing.T) {
	data, err := json.Marshal(AnimeFromSubject(tvSubject()))
	if err != nil {
		t.Fatal(err)
	}
	var keys map[string]json.RawMessage
	if err := json.Unmarshal(data, &keys); err != nil {
		t.Fatal(err)
	}

	wantKeys := []string{
		"id", "siteUrl", "status", "season", "type", "format",
		"seasonYear", "episodes", "synonyms", "isAdult",
		"countryOfOrigin", "meanScore", "description",
		"title", "coverImage", "startDate", "nameCN",
	}
	for _, k := range wantKeys {
		if _, ok := keys[k]; !ok {
			t.Errorf("序列化结果缺少字段 %q", k)
		}
	}

	// title 内部形状（romaji 为 nil 时被 omitempty 省略，属预期行为）
	var title map[string]json.RawMessage
	if err := json.Unmarshal(keys["title"], &title); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"english", "native", "userPreferred", "nameCN"} {
		if _, ok := title[k]; !ok {
			t.Errorf("title 缺少字段 %q", k)
		}
	}

	// 枚举值字符串保持 anilist 形状
	s := string(data)
	if !strings.Contains(s, `"status":"RELEASING"`) {
		t.Errorf("status 序列化值应为 RELEASING：%s", s)
	}
	if !strings.Contains(s, `"type":"ANIME"`) {
		t.Errorf("type 序列化值应为 ANIME：%s", s)
	}
	if !strings.Contains(s, `"format":"TV"`) {
		t.Errorf("format 序列化值应为 TV：%s", s)
	}
}
