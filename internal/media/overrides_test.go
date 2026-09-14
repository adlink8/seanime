package media

import (
	"path/filepath"
	"testing"

	"seanime/internal/api/animap"
)

// Phase 4 清偿工具支撑层测试：覆盖表持久化 + 队列移除 + resolver 覆盖优先级。

func TestFileOverrides_PutGetPersist(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", "overrides.json")
	o := NewFileOverrides(path)

	if err := o.Put(1001, 42); err != nil {
		t.Fatalf("Put: %v", err)
	}
	if id, ok := o.Get(1001); !ok || id != 42 {
		t.Fatalf("Get after Put: id=%d ok=%v", id, ok)
	}
	// 非法值静默忽略
	_ = o.Put(0, 42)
	_ = o.Put(1002, 0)
	if _, ok := o.Get(1002); ok {
		t.Fatal("invalid Put should be ignored")
	}
	// 重新打开应读到持久化数据
	o2 := NewFileOverrides(path)
	if id, ok := o2.Get(1001); !ok || id != 42 {
		t.Fatalf("persisted: id=%d ok=%v", id, ok)
	}
}

func TestFileQueue_Remove(t *testing.T) {
	path := filepath.Join(t.TempDir(), "queue.json")
	q := NewFileQueue(path)
	_ = q.Add(Unresolved{BangumiID: 1, Names: NameSet{Native: "a"}})
	_ = q.Add(Unresolved{BangumiID: 2, Names: NameSet{Native: "b"}})
	_ = q.Add(Unresolved{BangumiID: 3})

	if err := q.Remove(2, 999); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	items, _ := q.Load()
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	// 幂等：重复移除不报错
	if err := q.Remove(2); err != nil {
		t.Fatalf("idempotent Remove: %v", err)
	}
}

func TestAnimapResolver_OverridePriority(t *testing.T) {
	// 独一无二的标题，避免与 mockAnimapEntries 既有条目冲突
	entries := []animap.Anime{
		{Title: "Unique Auto Title", Mappings: &animap.AnimeMapping{AnidbID: 1}},
	}
	r := NewAnimapResolver(entries)

	// 自动解析先确认基线（精确命中 anidb=1）
	id, ok := r.ResolveBangumiToAniDB(777, NameSet{Native: "Unique Auto Title"})
	if !ok || id != 1 {
		t.Fatalf("baseline: id=%d ok=%v", id, ok)
	}

	// 未命中条目 → 覆盖表兜底
	r.ApplyOverride(888, 99)
	if id, ok := r.ResolveBangumiToAniDB(888, NameSet{Native: "什么都没有"}); !ok || id != 99 {
		t.Fatalf("override: id=%d ok=%v", id, ok)
	}

	// 覆盖表优先于自动解析（即使自动命中也以人工为准）
	r.ApplyOverride(777, 55)
	if id, ok := r.ResolveBangumiToAniDB(777, NameSet{Native: "Unique Auto Title"}); !ok || id != 55 {
		t.Fatalf("override priority: id=%d ok=%v", id, ok)
	}
}
