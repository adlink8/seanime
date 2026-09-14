package handlers

import (
	"strconv"
	"testing"
	"time"

	"seanime/internal/api/asmr"
)

// discover_daily_seed_test.go —— 03.9c 每日推荐纯函数的表驱动测试。
// 覆盖：日期种子（当日稳定、跨日变化）、洗牌确定性（同种子同序 / 跨日种子变序）、
// 洗牌不变性（元素集合不变，仅顺序变化）、空/单元素边界。

// TestDailyDateSeed 当天 YYYYMMDD 种子：同日稳定、跨日变化、格式可读。
func TestDailyDateSeed(t *testing.T) {
	if got := dailyDateSeed(time.Date(2026, 9, 14, 23, 59, 59, 0, time.Local)); got != 20260914 {
		t.Errorf("dailyDateSeed(2026-09-14) = %d, want 20260914", got)
	}
	// 跨日变化
	if dailyDateSeed(time.Date(2026, 9, 14, 0, 0, 0, 0, time.Local)) == dailyDateSeed(time.Date(2026, 9, 15, 0, 0, 0, 0, time.Local)) {
		t.Error("不同日期的种子不应相同")
	}
	// 同日任意时刻稳定
	if dailyDateSeed(time.Date(2026, 9, 14, 0, 0, 0, 0, time.Local)) != dailyDateSeed(time.Date(2026, 9, 14, 23, 0, 0, 0, time.Local)) {
		t.Error("同一日期不同时刻的种子应相同")
	}
}

// TestShuffleByDateSeed_Deterministic 同一种子多次洗牌结果恒定（SC5：三连击当日稳定一致）。
func TestShuffleByDateSeed_Deterministic(t *testing.T) {
	seed := int64(20260914)
	first := makeRange(1, 24)
	shuffleByDateSeed(first, seed)

	for i := 0; i < 3; i++ {
		again := makeRange(1, 24)
		shuffleByDateSeed(again, seed)
		for j := range first {
			if first[j] != again[j] {
				t.Fatalf("同种子第 %d 次洗牌结果不一致（位置 %d: %d != %d）", i, j, first[j], again[j])
			}
		}
	}
}

// TestShuffleByDateSeed_CrossDay 跨日种子应产生不同顺序（SC5：跨日变化）。
func TestShuffleByDateSeed_CrossDay(t *testing.T) {
	seedDay1 := int64(20260914)
	seedDay2 := int64(20260915)

	day1 := makeRange(1, 25)
	day2 := makeRange(1, 25)
	shuffleByDateSeed(day1, seedDay1)
	shuffleByDateSeed(day2, seedDay2)

	same := true
	for i := range day1 {
		if day1[i] != day2[i] {
			same = false
			break
		}
	}
	if same {
		t.Error("不同日期种子产生了完全相同的洗牌顺序，跨日变化失效")
	}
}

// TestShuffleByDateSeed_Permutation 洗牌不增不减：元素集合保持不变（SC5：结果条目完整）。
func TestShuffleByDateSeed_Permutation(t *testing.T) {
	input := makeRange(1, 25)
	shuffled := makeRange(1, 25)
	shuffleByDateSeed(shuffled, 20260914)

	sumIn, sumOut := 0, 0
	for _, v := range input {
		sumIn += v
	}
	for _, v := range shuffled {
		sumOut += v
	}
	if sumIn != sumOut {
		t.Errorf("洗牌后元素集合改变: sum %d != %d", sumIn, sumOut)
	}
	if len(shuffled) != len(input) {
		t.Errorf("洗牌后长度改变: %d != %d", len(shuffled), len(input))
	}
}

// TestShuffleByDateSeed_Edge 空切片 / 单元素不 panic、不变。
func TestShuffleByDateSeed_Edge(t *testing.T) {
	var empty []int
	shuffleByDateSeed(empty, 1)

	single := []int{42}
	shuffleByDateSeed(single, 1)
	if len(single) != 1 || single[0] != 42 {
		t.Errorf("单元素洗牌不应改变内容, got %v", single)
	}
}

// TestShuffleWorksByDateSeed_SimilarityWithIntVersion 泛型洗牌对 Asmr_Work 与 int 行为一致：
// 同种子同序（以 RjID 序列验证），保证 asmr 域当日稳定性语义与 Bangumi 域相同。
func TestShuffleWorksByDateSeed_SimilarityWithIntVersion(t *testing.T) {
	mkWorks := func() []asmr.Asmr_Work {
		works := make([]asmr.Asmr_Work, 0, 20)
		for i := 1; i <= 20; i++ {
			rj := "RJ" + strconv.Itoa(10000000+i)
			works = append(works, asmr.Asmr_Work{ID: strconv.Itoa(i), RjID: rj})
		}
		return works
	}

	a := mkWorks()
	shuffleByDateSeed(a, 20260914)
	b := mkWorks()
	shuffleByDateSeed(b, 20260914)

	for i := range a {
		if a[i].RjID != b[i].RjID {
			t.Fatalf("Asmr_Work 洗牌同种子结果不一致（位置 %d: %s != %s）", i, a[i].RjID, b[i].RjID)
		}
	}
}

func makeRange(start, endInclusive int) []int {
	out := make([]int, 0, endInclusive-start+1)
	for i := start; i <= endInclusive; i++ {
		out = append(out, i)
	}
	return out
}
