package handlers

import (
	"reflect"
	"testing"
)

// TestMapGenresToBangumiTags 是 03.9a genre→Bangumi tag 词表映射的表驱动测试。
// 覆盖：命中映射、未命中丢弃、空/nil 返回 nil、混合输入、顺序保持、大小写敏感。
func TestMapGenresToBangumiTags(t *testing.T) {
	cases := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name:  "nil 返回 nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "空切片返回 nil",
			input: []string{},
			want:  nil,
		},
		{
			name:  "单项命中",
			input: []string{"Comedy"},
			want:  []string{"搞笑"},
		},
		{
			name:  "多典型项命中且保持顺序",
			input: []string{"Action", "Comedy", "Mecha", "Romance"},
			want:  []string{"动作", "搞笑", "机战", "恋爱"},
		},
		{
			name:  "未命中项被丢弃且不透传英文值",
			input: []string{"Comedy", "SomeMadeUpGenre", "Sci-Fi"},
			want:  []string{"搞笑", "科幻"},
		},
		{
			name:  "全部未命中返回 nil",
			input: []string{"Foo", "Bar", "Baz"},
			want:  nil,
		},
		{
			name:  "大小写敏感：错误大小写视为未命中",
			input: []string{"comedy", "ACTION"},
			want:  nil,
		},
		{
			name:  "Supernatural 低频但保留映射",
			input: []string{"Supernatural"},
			want:  []string{"超自然"},
		},
		{
			name:  "含空白项混入",
			input: []string{"", "Drama", "   "},
			want:  []string{"剧情"},
		},
		{
			name:  "全部 18 项映射闭环",
			input: []string{"Action", "Adventure", "Comedy", "Drama", "Ecchi", "Fantasy", "Horror", "Mahou Shoujo", "Mecha", "Music", "Mystery", "Psychological", "Romance", "Sci-Fi", "Slice of Life", "Sports", "Supernatural", "Thriller"},
			want:  []string{"动作", "冒险", "搞笑", "剧情", "福利", "奇幻", "恐怖", "魔法少女", "机战", "音乐", "悬疑", "心理", "恋爱", "科幻", "日常", "运动", "超自然", "惊悚"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := mapGenresToBangumiTags(tc.input)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("mapGenresToBangumiTags(%v) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}
