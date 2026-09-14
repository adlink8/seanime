package handlers

// genreToBangumiTag 是 AniList genre → Bangumi 中文 tag 的映射表（03.9a）。
//
// 词表依据 .planning/codebase/TAG-VOCAB-RECON.md §1 实测（动作=687、搞笑=1000+、机战=1000+ …），
// 覆盖 18 个 AniList genre 规范值。Bangumi 无 genres 维度，搜索时 genre 被并入 filter.tag，
// 故这里把英文 genre 翻译为 Bangumi 实际可命中的中文 tag，修复「选类型筛选基本失效」问题
// （此前英文 tag 仅欧美条目偶然携带，命中极低）。
//
// 未命中词表的 genre 一律丢弃（绝不透传英文原始值）。
var genreToBangumiTag = map[string]string{
	"Action":        "动作",
	"Adventure":     "冒险",
	"Comedy":        "搞笑",
	"Drama":         "剧情",
	"Ecchi":         "福利",
	"Fantasy":       "奇幻",
	"Horror":        "恐怖",
	"Mahou Shoujo":  "魔法少女",
	"Mecha":         "机战",
	"Music":         "音乐",
	"Mystery":       "悬疑",
	"Psychological": "心理",
	"Romance":       "恋爱",
	"Sci-Fi":        "科幻",
	"Slice of Life": "日常",
	"Sports":        "运动",
	// Supernatural→超自然 实测仅 6 条命中（RECON §1），属上游词表事实，保留诚实翻译不裁 UI；
	// 若后续实测体验差可再裁。
	"Supernatural": "超自然",
	"Thriller":     "惊悚",
}

// mapGenresToBangumiTags 将 AniList genre 列表映射为 Bangumi 中文 tag 列表。
//
// 规则：
//   - 输入 nil 或空切片返回 nil（不返回空切片，保证下游 filter.Tag 判空一致）；
//   - 逐项查表，未命中词表的 genre 丢弃，不修改、不透传任何英文原始值；
//   - 大小写敏感：仅匹配上面词表的规范大小写（AniList 下发值即规范大小写）。
func mapGenresToBangumiTags(genres []string) []string {
	if len(genres) == 0 {
		return nil
	}
	out := make([]string, 0, len(genres))
	for _, g := range genres {
		if tag, ok := genreToBangumiTag[g]; ok {
			out = append(out, tag)
		}
	}
	// 全部未命中时返回 nil，避免向 filter.Tag 追加空切片（语义更清晰，且与空输入一致）。
	if len(out) == 0 {
		return nil
	}
	return out
}
