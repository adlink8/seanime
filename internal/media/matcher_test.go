package media

import (
	"testing"
)

// 匹配器单测：归一化与相似度的边界，全部离线。

func TestNormalize_FullWidthToHalfWidth(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"Ｓｔｅｉｎｓ；Ｇａｔｅ", "steinsgate"},       // 全角字母+全角分号
		{"化物語　全１５話", "化物語全15話"},             // 全角空格 + 全角数字
		{"ＧＯＳＩＣＫ!", "gosick"},                // 全角字母 + 半角标点
		{"食戟のソーマ（ＴＶ版）", "食戟のソーマ"},            // 全角括号剥离
	}
	for _, c := range cases {
		if got := Normalize(c.in); got != c.want {
			t.Errorf("Normalize(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestNormalize_BracketSubtitleStrip(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"化物語（全15話）", "化物語"},                            // 全角括号副标题
		{"Steins;Gate (TV)", "steinsgate"},              // 半角括号
		{"Shingeki no Kyojin [2013]", "shingekinokyojin"}, // 方括号年份
		{"〈物語〉セカンドシーズン", "セカンドシーズン"},                     // 书名号剥离
		{"CLANNAD ～AFTER STORY～", "clannadafterstory"},  // 非括号波浪线保留（仅标点被移除）
		{"孤立的右括号 ) test", "孤立的右括号test"},                 // 不配对括号容忍
	}
	for _, c := range cases {
		if got := Normalize(c.in); got != c.want {
			t.Errorf("Normalize(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestNormalize_PunctuationAndSeason(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"Steins;Gate", "steinsgate"},
		{"Toradora!", "toradora"},
		{"K-ON!", "kon"},
		{"Monogatari Series II", "monogatariseries2"},   // 罗马数字季号
		{"Attack on Titan III", "attackontitan3"},
		{"Overlord Ⅳ", "overlord4"}, // 单字符 Unicode 罗马数字（U+2163）显式转换
		{"", ""},
		{"　", ""}, // 纯全角空格
	}
	for _, c := range cases {
		if got := Normalize(c.in); got != c.want {
			t.Errorf("Normalize(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestNormalize_ConsistencyBetweenIndexAndQuery(t *testing.T) {
	// 索引侧与查询侧对"同义异形"标题必须得到同一键
	pairs := [][2]string{
		{"Steins;Gate", "Ｓｔｅｉｎｓ；Ｇａｔｅ"},
		{"Toradora!", "Toradora"},
		{"化物語 (全15話)", "化物語"},
		{"Violet Evergarden", "violet  evergarden"}, // 多空格
	}
	for _, p := range pairs {
		a, b := Normalize(p[0]), Normalize(p[1])
		if a != b {
			t.Errorf("Normalize(%q)=%q != Normalize(%q)=%q", p[0], a, p[1], b)
		}
	}
}

func TestSimilarity_Basics(t *testing.T) {
	// 完全相同/不同
	if got := Similarity("abc", "abc"); got != 1 {
		t.Errorf("Similarity identical = %v, want 1", got)
	}
	if got := Similarity("abcd", "wxyz"); got != 0 {
		t.Errorf("Similarity disjoint = %v, want 0", got)
	}
	// 单字符错字："shingekinokyojin"(16) vs "shingekinokyojim"(16) → dist 1 → 1-1/16 = 0.9375
	if got := Similarity("shingekinokyojin", "shingekinokyojim"); got != 0.9375 {
		t.Errorf("Similarity single-char typo = %v, want 0.9375", got)
	}
	// 前缀差异大的两个词
	if got := Similarity("frieren", "sousounofrieren"); got > 0.6 {
		t.Errorf("Similarity prefix-extended = %v, want < 0.6", got)
	}
}

func TestSimilarity_ThresholdBehavior(t *testing.T) {
	// 阈值 0.90 的语义校验：10% 以内的字符差异可通过，超过则拒绝
	long := "suzumiyaharuheinoyuuutsu"
	oneOff := "suzumiyaharuheinoyuuutss" // 尾部错 1 字符
	if s := Similarity(oneOff, long); s < 0.9 {
		t.Errorf("1-char typo on len %d = %v, want >= 0.9", len([]rune(long)), s)
	}
	heavily := "suzumiya" // 截断过多
	if s := Similarity(heavily, long); s >= 0.9 {
		t.Errorf("truncated title = %v, want < 0.9", s)
	}
}
