//go:build integration

// 可选集成验证：使用真实 animap 数据（网络拉取）验证反向索引与名称匹配。
//
// 门禁（CI 上必须默认跳过）：
//   1. 构建标签 integration（go test -tags=./internal/media/...）；
//   2. 环境变量 SEANIME_MEDIA_INTEGRATION=1。
// 两项同时满足才执行，任一缺失则 t.Skip。
//
// 数据获取方式：internal/api/animap 只提供按 ID 单查（anilist/mal/anidb ID → entry），
// 无全量拉取端点。本测试用已知 anilist ID 逐条拉取真实条目，构建反向索引，
// 再用条目自身标题作为查询验证「名称 → anidb」的往返解析。
package media

import (
	"os"
	"testing"

	"seanime/internal/api/animap"
)

// 真实动画样本的 anilist ID（公开可查）
var integrationAnilistIDs = []int{
	5081,   // Bakemonogatari
	16498,  // Shingeki no Kyojin
	154587, // Sousou no Frieren
	9253,   // Steins;Gate
	4381,   // Suzumiya Haruhi no Yuuutsu
	2167,   // CLANNAD
	30,     // Neon Genesis Evangelion
	12189,  // Koukaku Kidoutai: Stand Alone Complex
	21827,  // Violet Evergarden
	4224,   // Toradora
}

func TestIntegration_RealAnimapResolve(t *testing.T) {
	if os.Getenv("SEANIME_MEDIA_INTEGRATION") != "1" {
		t.Skip("需要环境变量 SEANIME_MEDIA_INTEGRATION=1（真实网络集成测试，CI 默认跳过）")
	}

	// 逐条拉取真实 animap 条目
	var entries []animap.Anime
	for _, id := range integrationAnilistIDs {
		e, err := animap.FetchAnimapMedia("anilist", id)
		if err != nil {
			// 单条失败不致命（ID 变动/网络抖动），记录并继续
			t.Logf("anilist_id=%d 拉取失败: %v", id, err)
			continue
		}
		if e.Mappings == nil || e.Mappings.AnidbID == 0 {
			t.Logf("anilist_id=%d 无 anidb 映射，跳过", id)
			continue
		}
		entries = append(entries, *e)
	}
	if len(entries) < 5 {
		t.Skipf("真实数据拉取不足（成功 %d 条），跳过验证", len(entries))
	}

	r := NewAnimapResolver(entries)

	// 往返验证：用条目自身标题集查询，应解析回同一 anidb ID
	success := 0
	for _, e := range entries {
		names := NameSet{Native: e.Title}
		for _, v := range e.Titles {
			names.Aliases = append(names.Aliases, v)
		}
		got, ok := r.ResolveBangumiToAniDB(e.Mappings.AnidbID, names)
		if !ok || got != e.Mappings.AnidbID {
			t.Errorf("%q 往返解析失败: ok=%v got=%d want=%d", e.Title, ok, got, e.Mappings.AnidbID)
			continue
		}
		success++
	}
	t.Logf("真实数据往返解析成功率: %d/%d", success, len(entries))
}
