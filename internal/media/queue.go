package media

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// 待映射队列：ResolveBangumiToAniDB 无法保守确认的条目（名称未命中/相似度不足/
// 存在歧义）落盘为 JSON 文件，供 Phase 4 清偿工具（人工确认/补充数据源）消费。
//
// 设计要点：
//   - 路径可注入（NewFileQueue），单测写临时目录，生产由编排层给配置路径；
//   - 按 bangumiID 去重，同一条目重复解析失败不产生重复记录；
//   - 读-改-写全程持锁，进程内并发安全（自用单机场景，无跨进程锁需求）。

// Unresolved 一条待映射记录。
type Unresolved struct {
	BangumiID int      `json:"bangumi_id"` // Bangumi subject ID
	Names     NameSet  `json:"names"`      // 解析时的名称集，供清偿工具人工核对
	QueuedAt  time.Time `json:"queued_at"` // 入队时间
}

// Queue 待映射队列接口。animapResolver 只依赖该接口，单测可替换实现。
type Queue interface {
	Add(u Unresolved) error
	Load() ([]Unresolved, error)
}

// FileQueue 基于 JSON 文件的待映射队列实现。
type FileQueue struct {
	path string
	mu   sync.Mutex
}

// 编译期断言：FileQueue 满足 Queue 接口。
var _ Queue = (*FileQueue)(nil)

// NewFileQueue 创建以 path 为持久化文件的队列。文件不存在时按空队列处理。
func NewFileQueue(path string) *FileQueue {
	return &FileQueue{path: path}
}

// Add 入队一条记录（按 BangumiID 去重），随后原子落盘。
func (q *FileQueue) Add(u Unresolved) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	items, err := q.readLocked()
	if err != nil {
		return err
	}
	for _, it := range items {
		if it.BangumiID == u.BangumiID {
			return nil // 已存在，去重
		}
	}
	items = append(items, u)
	return q.writeLocked(items)
}

// Load 读取全部待映射记录。
func (q *FileQueue) Load() ([]Unresolved, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.readLocked()
}

// Remove 清偿后移除指定条目（幂等：不存在不报错）。
func (q *FileQueue) Remove(bangumiIDs ...int) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	items, err := q.readLocked()
	if err != nil {
		return err
	}
	drop := make(map[int]bool, len(bangumiIDs))
	for _, id := range bangumiIDs {
		drop[id] = true
	}
	kept := items[:0]
	for _, it := range items {
		if !drop[it.BangumiID] {
			kept = append(kept, it)
		}
	}
	if len(kept) == len(items) {
		return nil // 无变化不落盘
	}
	return q.writeLocked(kept)
}

// readLocked 读取文件（调用方需持锁）。
func (q *FileQueue) readLocked() ([]Unresolved, error) {
	data, err := os.ReadFile(q.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}
	var items []Unresolved
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// writeLocked 写入文件（先写临时文件再 rename，避免写一半损坏）。
func (q *FileQueue) writeLocked(items []Unresolved) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(q.path), 0o755); err != nil {
		return err
	}
	tmp := q.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, q.path)
}
