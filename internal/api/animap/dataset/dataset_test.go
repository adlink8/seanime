package dataset

import (
	"testing"
)

// Convert：manami 条目 → animap.Anime 的转换规则测试。

func TestConvert_WithAnidbAndMal(t *testing.T) {
	e := &manamiEntry{
		Sources: []string{
			"https://anidb.net/anime/1234",
			"https://myanimelist.net/anime/5678",
			"https://anilist.co/anime/99",
		},
		Title: "Test Anime",
		Synonyms: []string{"Test Anime", "テストアニメ", "测试动画"},
	}
	a, ok := Convert(e)
	if !ok {
		t.Fatal("expected ok=true")
	}
	if a.Mappings == nil || a.Mappings.AnidbID != 1234 || a.Mappings.MalID != 5678 {
		t.Fatalf("unexpected mappings: %+v", a.Mappings)
	}
	if a.Title != "Test Anime" {
		t.Fatalf("unexpected title: %q", a.Title)
	}
	// Title 不重复出现在 Titles；synonyms 全部收录
	if _, dup := a.Titles["Test Anime"]; dup {
		t.Fatalf("title should not duplicate in Titles map")
	}
	if a.Titles["テストアニメ"] != "テストアニメ" || a.Titles["测试动画"] != "测试动画" {
		t.Fatalf("synonyms missing: %+v", a.Titles)
	}
}

func TestConvert_NoAnidbSkipped(t *testing.T) {
	e := &manamiEntry{
		Sources: []string{"https://anilist.co/anime/99"},
		Title:   "No AniDB",
	}
	if _, ok := Convert(e); ok {
		t.Fatal("expected ok=false when no anidb source")
	}
}

func TestConvert_TrailingSlashAndInvalid(t *testing.T) {
	// 带尾斜杠的 ID 应可提取
	e := &manamiEntry{
		Sources: []string{"https://anidb.net/anime/42/"},
		Title:   "Trailing",
	}
	a, ok := Convert(e)
	if !ok || a.Mappings.AnidbID != 42 {
		t.Fatalf("expected anidb=42, got ok=%v id=%d", ok, a.Mappings.AnidbID)
	}
	// 非 ID source 应被忽略
	e2 := &manamiEntry{
		Sources: []string{"https://anidb.net/anime/notanumber"},
		Title:   "Bad",
	}
	if _, ok := Convert(e2); ok {
		t.Fatal("expected ok=false for non-numeric id")
	}
}
