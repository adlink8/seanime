package media

import (
	"path/filepath"
	"testing"
	"time"

	"seanime/internal/api/animap"
)

//----------------------------------------------------------------------------------------------------------------------
// Mock animap 数据：手工编写的 12 条 animap 形状条目，覆盖真实动画样本。
// anidb ID 参考公开资料的常见值（测试只关心映射机制的确定性，不校验 ID 的真实权威性）。

func mockAnimapEntries() []animap.Anime {
	m := func(anidb int) *animap.AnimeMapping { return &animap.AnimeMapping{AnidbID: anidb} }
	return []animap.Anime{
		{Title: "Bakemonogatari", Mappings: m(6327), Titles: map[string]string{"ja": "化物語", "zh-Hans": "化物语"}},
		{Title: "Shingeki no Kyojin", Mappings: m(8679), Titles: map[string]string{"ja": "進撃の巨人", "zh-Hans": "进击的巨人"}},
		{Title: "Sousou no Frieren", Mappings: m(18139), Titles: map[string]string{"ja": "葬送のフリーレン", "zh-Hans": "葬送的芙莉莲"}},
		{Title: "Steins;Gate", Mappings: m(7052), Titles: map[string]string{"ja": "シュタインズ・ゲート", "zh-Hans": "命运石之门"}},
		{Title: "Suzumiya Haruhi no Yuuutsu", Mappings: m(4550), Titles: map[string]string{"ja": "涼宮ハルヒの憂鬱", "zh-Hans": "凉宫春日的忧郁"}},
		{Title: "CLANNAD", Mappings: m(4586), Titles: map[string]string{"ja": "CLANNAD"}},
		{Title: "Shin Seiki Evangelion", Mappings: m(289), Titles: map[string]string{"ja": "新世紀エヴァンゲリオン", "en": "Neon Genesis Evangelion", "zh-Hans": "新世纪福音战士"}},
		{Title: "Koukaku Kidoutai: Stand Alone Complex", Mappings: m(819), Titles: map[string]string{"ja": "攻殻機動隊 STAND ALONE COMPLEX", "zh-Hans": "攻壳机动队"}},
		{Title: "Violet Evergarden", Mappings: m(13484), Titles: map[string]string{"ja": "ヴァイオレット・エヴァーガーデン", "zh-Hans": "紫罗兰永恒花园"}},
		{Title: "Toradora!", Mappings: m(5340), Titles: map[string]string{"ja": "とらドラ!", "zh-Hans": "龙与虎"}},
		// 近似匹配样本：查询侧带单字符变体时依赖编辑距离
		{Title: "Zankyou no Terror", Mappings: m(9891), Titles: map[string]string{"ja": "残響のテロル", "zh-Hans": "恐怖残响"}},
		// 括号变体样本：animap 标题带副标题后缀，查询侧是纯净名
		{Title: "Ansatsu Kyoushitsu (TV)", Mappings: m(10734), Titles: map[string]string{"ja": "暗殺教室", "zh-Hans": "暗杀教室"}},
	}
}

// knownSamples 已知动画样本（≥10 部）：模拟 Bangumi 侧的 NameSet。
// bangumiID 使用 1~N 的测试序号即可。
var knownSamples = []struct {
	name      string
	bangumiID int
	names     NameSet
	wantAniDB int
}{
	{"化物语", 1, NameSet{Native: "化物語", Chinese: "化物语"}, 6327},
	{"进击的巨人", 2, NameSet{Native: "進撃の巨人", Chinese: "进击的巨人"}, 8679},
	{"葬送的芙莉莲", 3, NameSet{Native: "葬送のフリーレン", Chinese: "葬送的芙莉莲"}, 18139},
	{"命运石之门", 4, NameSet{Native: "シュタインズ・ゲート", Romaji: "Steins;Gate", Chinese: "命运石之门"}, 7052},
	{"凉宫春日的忧郁", 5, NameSet{Native: "涼宮ハルヒの憂鬱", Chinese: "凉宫春日的忧郁"}, 4550},
	{"CLANNAD", 6, NameSet{English: "CLANNAD", Chinese: "团子大家族"}, 4586},
	{"新世纪福音战士", 7, NameSet{Native: "新世紀エヴァンゲリオン", English: "Neon Genesis Evangelion", Chinese: "新世纪福音战士"}, 289},
	{"攻壳机动队SAC", 8, NameSet{Native: "攻殻機動隊 STAND ALONE COMPLEX", Chinese: "攻壳机动队"}, 819},
	{"紫罗兰永恒花园", 9, NameSet{Native: "ヴァイオレット・エヴァーガーデン", Chinese: "紫罗兰永恒花园"}, 13484},
	{"龙与虎", 10, NameSet{Native: "とらドラ!", Romaji: "Toradora", Chinese: "龙与虎"}, 5340},
	{"暗杀教室-括号变体", 11, NameSet{Native: "暗殺教室", Romaji: "Ansatsu Kyoushitsu", Chinese: "暗杀教室"}, 10734},
}

// newTestResolver 组装带临时文件队列的解析器，返回解析器与队列文件路径。
func newTestResolver(t *testing.T, entries []animap.Anime) (*AnimapResolver, string) {
	t.Helper()
	qPath := filepath.Join(t.TempDir(), "unresolved.json")
	r := NewAnimapResolver(entries, WithQueue(NewFileQueue(qPath)))
	return r, qPath
}

// TestResolve_KnownSamples 核心验收：≥10 部已知动画全部映射成功。
func TestResolve_KnownSamples(t *testing.T) {
	r, _ := newTestResolver(t, mockAnimapEntries())

	success := 0
	for _, s := range knownSamples {
		got, ok := r.ResolveBangumiToAniDB(s.bangumiID, s.names)
		if !ok {
			t.Errorf("[样本 %d:%s] 期望映射到 anidb %d，实际未解析", s.bangumiID, s.name, s.wantAniDB)
			continue
		}
		if got != s.wantAniDB {
			t.Errorf("[样本 %d:%s] 期望 anidb %d，实际 %d", s.bangumiID, s.name, s.wantAniDB, got)
			continue
		}
		success++
	}
	if success < 10 {
		t.Errorf("样本映射成功 %d/%d，低于契约要求的 10 部", success, len(knownSamples))
	}
	t.Logf("样本映射成功率: %d/%d = %.0f%%", success, len(knownSamples), float64(success)/float64(len(knownSamples))*100)
}

// TestResolve_ExactMatchesStillWorkOnSecondCall 验证解析缓存：第二次调用直接命中。
func TestResolve_Caching(t *testing.T) {
	r, _ := newTestResolver(t, mockAnimapEntries())
	s := knownSamples[0]
	for i := 0; i < 3; i++ {
		got, ok := r.ResolveBangumiToAniDB(s.bangumiID, s.names)
		if !ok || got != s.wantAniDB {
			t.Fatalf("第 %d 次调用解析失败: ok=%v id=%d", i+1, ok, got)
		}
	}
}

// TestResolve_BracketedQuery 查询侧带括号副标题（如来自其他数据源的脏数据）仍应精确命中。
func TestResolve_BracketedQuery(t *testing.T) {
	r, _ := newTestResolver(t, mockAnimapEntries())
	got, ok := r.ResolveBangumiToAniDB(100, NameSet{Native: "化物語（全15話）"})
	if !ok || got != 6327 {
		t.Errorf("括号变体查询失败: ok=%v id=%d", ok, got)
	}
}

// TestResolve_FuzzyVariant 近似匹配：查询侧单字符变体（模拟名称差异）应命中。
func TestResolve_FuzzyVariant(t *testing.T) {
	r, _ := newTestResolver(t, mockAnimapEntries())
	// "Zankyou no Terrer"（1 字符错字）→ "Zankyou no Terror"（anidb 9891）
	got, ok := r.ResolveBangumiToAniDB(200, NameSet{Romaji: "Zankyou no Terrer"})
	if !ok || got != 9891 {
		t.Errorf("近似匹配失败: ok=%v id=%d", ok, got)
	}
}

// TestResolve_AmbiguousTieQueues 歧义保护：两个候选分数近似并列时必须拒绝并入队。
// 两个 mock 条目 "Example Anime Show"/"Example Anime Shows" 与查询
// "Example Anime Showa" 的分数差 < 0.02，宁可入队不硬猜。
func TestResolve_AmbiguousTieQueues(t *testing.T) {
	entries := []animap.Anime{
		{Title: "Example Anime Show", Mappings: &animap.AnimeMapping{AnidbID: 100}},
		{Title: "Example Anime Shows", Mappings: &animap.AnimeMapping{AnidbID: 200}},
	}
	r, qPath := newTestResolver(t, entries)

	got, ok := r.ResolveBangumiToAniDB(300, NameSet{English: "Example Anime Showa"})
	if ok {
		t.Fatalf("歧义候选不应解析，实际返回 %d", got)
	}
	assertQueued(t, qPath, 300)
}

// TestResolve_ConflictTitleQueues 重名保护：同一标题对应多个 anidb ID 时拒绝。
func TestResolve_ConflictTitleQueues(t *testing.T) {
	entries := []animap.Anime{
		{Title: "Re:Zero kara Hajimeru Isekai Seikatsu", Mappings: &animap.AnimeMapping{AnidbID: 21226}},
		{Title: "Re:Zero kara Hajimeru Isekai Seikatsu", Mappings: &animap.AnimeMapping{AnidbID: 21519}},
	}
	r, qPath := newTestResolver(t, entries)

	got, ok := r.ResolveBangumiToAniDB(400, NameSet{English: "Re:Zero kara Hajimeru Isekai Seikatsu"})
	if ok {
		t.Fatalf("冲突标题不应解析，实际返回 %d", got)
	}
	assertQueued(t, qPath, 400)
}

// TestResolve_BelowThresholdQueues 相似度不足时拒绝并入队（不硬猜）。
func TestResolve_BelowThresholdQueues(t *testing.T) {
	r, qPath := newTestResolver(t, mockAnimapEntries())

	got, ok := r.ResolveBangumiToAniDB(500, NameSet{Chinese: "完全无关的原创动画作品"})
	if ok {
		t.Fatalf("无关名称不应解析，实际返回 %d", got)
	}
	assertQueued(t, qPath, 500)
}

// TestResolve_ShortTitleExactOnly 短标题只做精确匹配："龙与虎"截断变体不应命中。
func TestResolve_ShortTitleExactOnly(t *testing.T) {
	r, _ := newTestResolver(t, mockAnimapEntries())
	// "龙与"（3 字符 < 6）不参与模糊匹配，索引中无此精确键
	got, ok := r.ResolveBangumiToAniDB(600, NameSet{Chinese: "龙与"})
	if ok {
		t.Errorf("短标题模糊匹配被意外放行: id=%d", got)
	}
}

// TestQueueUnresolved_Explicit 幂等入队：同一条目重复入队只落一条。
func TestQueueUnresolved_Explicit(t *testing.T) {
	r, qPath := newTestResolver(t, mockAnimapEntries())
	names := NameSet{Chinese: "未知作品A"}
	r.QueueUnresolved(700, names)
	r.QueueUnresolved(700, names)
	r.QueueUnresolved(701, NameSet{Chinese: "未知作品B"})

	q := NewFileQueue(qPath)
	items, err := q.Load()
	if err != nil {
		t.Fatalf("读取队列失败: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("期望队列 2 条（去重后），实际 %d", len(items))
	}
	if items[0].BangumiID != 700 || items[1].BangumiID != 701 {
		t.Errorf("队列记录顺序/内容不符: %+v", items)
	}
	if items[0].Names.Chinese != "未知作品A" {
		t.Errorf("队列应携带名称集供清偿工具核对: %+v", items[0].Names)
	}
}

// TestFileQueue_Persistence 空队列/落盘/重载行为。
func TestFileQueue_Persistence(t *testing.T) {
	qPath := filepath.Join(t.TempDir(), "sub", "queue.json") // 目录不存在时应自动创建
	q := NewFileQueue(qPath)

	// 空队列
	items, err := q.Load()
	if err != nil || len(items) != 0 {
		t.Fatalf("新队列应为空: items=%v err=%v", items, err)
	}

	u := Unresolved{BangumiID: 42, Names: NameSet{Native: "テスト", Chinese: "测试"}, QueuedAt: time.Now()}
	if err := q.Add(u); err != nil {
		t.Fatalf("Add 失败: %v", err)
	}
	// 重新实例化（模拟重启后重载）
	q2 := NewFileQueue(qPath)
	items, err = q2.Load()
	if err != nil {
		t.Fatalf("重载失败: %v", err)
	}
	if len(items) != 1 || items[0].BangumiID != 42 || items[0].Names.Native != "テスト" {
		t.Fatalf("重载数据不符: %+v", items)
	}
}

// assertQueued 断言队列文件中存在指定 bangumiID 的记录。
func assertQueued(t *testing.T, qPath string, bangumiID int) {
	t.Helper()
	items, err := NewFileQueue(qPath).Load()
	if err != nil {
		t.Fatalf("读取队列失败: %v", err)
	}
	for _, it := range items {
		if it.BangumiID == bangumiID {
			return
		}
	}
	t.Errorf("bangumiID %d 未入队，队列内容: %+v", bangumiID, items)
}
