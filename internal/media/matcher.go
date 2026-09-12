package media

import (
	"strings"
	"unicode"
)

// 名称归一化与相似度匹配。
//
// 设计意图：bangumi→anidb 映射的唯一可行路径是「名称匹配」（Bangumi ID 不在
// animap 数据集中），因此归一化规则直接决定映射命中率。规则两侧一致应用
// （索引侧与查询侧同走 Normalize），保证相同输入得到相同键。

// unicodeRoman 需要归一化的季号罗马数字查表（ASCII 大写形式 → 阿拉伯数字）。
var unicodeRoman = map[string]string{
	"I": "1", "II": "2", "III": "3", "IV": "4", "V": "5",
	"VI": "6", "VII": "7", "VIII": "8", "IX": "9", "X": "10",
	"XI": "11", "XII": "12",
}

// unicodeRomanChars 单字符 Unicode 罗马数字（U+2160~U+216B 大写、U+2170~U+217B 小写，
// 类别 Nl）→ ASCII 数字串。Nl 类字符既不是 IsLetter 也不是 IsDigit，
// 若不在此处显式转换会被后续步骤丢弃。小写区段必须覆盖：Normalize 先做
// ToLower，大写罗马数字会先变成小写形式。
var unicodeRomanChars = map[rune]string{
	0x2160: "1", 0x2161: "2", 0x2162: "3", 0x2163: "4", 0x2164: "5",
	0x2165: "6", 0x2166: "7", 0x2167: "8", 0x2168: "9", 0x2169: "10",
	0x216A: "11", 0x216B: "12",
	0x2170: "1", 0x2171: "2", 0x2172: "3", 0x2173: "4", 0x2174: "5",
	0x2175: "6", 0x2176: "7", 0x2177: "8", 0x2178: "9", 0x2179: "10",
	0x217A: "11", 0x217B: "12",
}

// 说明：仅处理 I~XII 的独立 token。将独立 token "V" 映射为 "5" 理论上可能误伤
// 某些标题，但两侧规则一致所以索引/查询对称失效，实际风险可接受；
// 换来的收益是 "Series II" 与 "Series 2" 能归一到同一键。

// bracketPairs 需要剥离的括号对（含内容）。
// 动画标题中括号内通常是副标题/季标记/放送信息（如 "化物語 (全15話)"、
// "Steins;Gate (TV)"），这类后缀在 animap 与 Bangumi 两侧常不一致，剥离后
// 能显著提高精确匹配率。全半角以及 CJK 常见括号都覆盖。
var bracketPairs = [][2]rune{
	{'(', ')'}, {'（', '）'}, // 全半角圆括号
	{'[', ']'}, {'【', '】'}, {'〔', '〕'}, // 方括号系
	{'{', '}'}, {'｛', '｝'},
	{'《', '》'}, {'〈', '〉'}, // 书名号/单书名号
	{'「', '」'}, {'『', '』'}, // 直角引号（部分标题用于括注）
}

// Normalize 将名称归一化为匹配键：
//  1. 转小写；
//  2. 全角 ASCII（！～Ａａ１２等）与全角空格转半角；
//  3. 剥离括号段（副标题剥离）；
//  4. 独立 token 的罗马数字季号转阿拉伯数字（"II"→"2"）；
//  5. 移除剩余空白与标点符号，仅保留字母/数字/CJK。
//
// 例：
//
//	"Steins;Gate"          -> "steinsgate"
//	"化物語（全15話）"       -> "化物語全15話"……注：括号段整段剥离，得 "化物語"
//	"Monogatari Series II" -> "monogatariseries2"
func Normalize(name string) string {
	s := strings.ToLower(name)

	// 全角 -> 半角（0xFF01~0xFF5E 对应 ASCII 0x21~0x7E，偏移 0xFEE0）
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case r >= 0xFF01 && r <= 0xFF5E:
			b.WriteRune(r - 0xFEE0)
		case r == 0x3000: // 全角空格
			b.WriteRune(' ')
		default:
			if d, ok := unicodeRomanChars[r]; ok {
				b.WriteString(d)
			} else {
				b.WriteRune(r)
			}
		}
	}
	s = b.String()

	// 剥离括号段（非贪婪：先内后外自然成立，逐字符扫描配对）
	s = stripBrackets(s)

	// 按 token 处理季号罗马数字（此时标点仍在，token 以空白/标点分隔）
	s = normalizeRomanTokens(s)

	// 移除空白与标点，仅保留字母/数字
	var out strings.Builder
	out.Grow(len(s))
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			out.WriteRune(r)
		}
	}
	return out.String()
}

// stripBrackets 剥离成对括号及其内容，支持嵌套。
func stripBrackets(s string) string {
	var out strings.Builder
	out.Grow(len(s))
	var stack []rune // 未闭合的左括号栈
	for _, r := range s {
		opened := false
		for _, p := range bracketPairs {
			if r == p[0] {
				stack = append(stack, r)
				opened = true
				break
			}
		}
		if opened {
			continue
		}
		closed := false
		for _, p := range bracketPairs {
			if r == p[1] {
				// 弹出栈中最近的左括号（容忍不配对：孤立右括号直接丢弃）
				if len(stack) > 0 {
					stack = stack[:len(stack)-1]
				}
				closed = true
				break
			}
		}
		if closed {
			continue
		}
		if len(stack) == 0 {
			out.WriteRune(r)
		}
	}
	return out.String()
}

// normalizeRomanTokens 将独立 token 的罗马数字替换为阿拉伯数字。
// token 以 unicode 空白/标点分隔；纯 ASCII 字母的 token 且命中罗马数字表才替换。
func normalizeRomanTokens(s string) string {
	isSep := func(r rune) bool {
		return unicode.IsSpace(r) || (unicode.IsPunct(r) && r != '-')
	}
	fields := strings.FieldsFunc(s, isSep)
	if len(fields) == 0 {
		return s
	}
	replaced := false
	for i, f := range fields {
		// 查表用大写键（token 已因 ToLower 变小写）
		if d, ok := unicodeRoman[strings.ToUpper(f)]; ok {
			fields[i] = d
			replaced = true
		}
	}
	if !replaced {
		return s
	}
	// 用单空格重组；后续步骤会移除空白，重组形态不影响最终键
	return strings.Join(fields, " ")
}

// SimilarityFunc 计算两个已归一化名称的相似度，取值 [0,1]。
type SimilarityFunc func(a, b string) float64

// Similarity 基于归一化编辑距离（Levenshtein ratio）的相似度：
//
//	sim = 1 - dist(a,b) / max(len(a), len(b))
//
// 选型理由：自用场景、候选量在万级以内，O(n*m) DP 足够；
// 不引入重型相似度库（n-gram/Jaro-Winkler 等），编辑距离对标题的
// 单字符变体（异体字、错字、空格差异）最直观。
func Similarity(a, b string) float64 {
	if a == b {
		return 1
	}
	ra, rb := []rune(a), []rune(b)
	la, lb := len(ra), len(rb)
	if la == 0 || lb == 0 {
		return 0
	}
	maxLen := la
	if lb > maxLen {
		maxLen = lb
	}
	dist := levenshtein(ra, rb)
	return 1 - float64(dist)/float64(maxLen)
}

// levenshtein 标准编辑距离 DP，空间压缩到单行。
func levenshtein(a, b []rune) int {
	prev := make([]int, len(b)+1)
	curr := make([]int, len(b)+1)
	for j := 0; j <= len(b); j++ {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		curr[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			curr[j] = min3(curr[j-1]+1, prev[j]+1, prev[j-1]+cost)
		}
		prev, curr = curr, prev
	}
	return prev[len(b)]
}

func min3(a, b, c int) int {
	m := a
	if b < m {
		m = b
	}
	if c < m {
		m = c
	}
	return m
}
