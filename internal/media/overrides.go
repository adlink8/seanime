package media

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// 人工映射覆盖表（Phase 4 清偿工具的持久化层）。
//
// 待映射队列里的条目经人工确认 bangumi ↔ anidb 对应关系后写入本表；
// 生产接线时启动加载、运行时 ApplyOverride 即时生效。覆盖表优先级
// 最高（高于名称解析），错误映射由人工负责——这正是保守自动解析
// 与人工兜底的分工。

// Overrides 覆盖表接口（测试可替换）。
type Overrides interface {
	Get(bangumiID int) (anidbID int, ok bool)
	Put(bangumiID, anidbID int) error
	All() map[int]int
}

// FileOverrides 基于 JSON 文件的覆盖表实现。文件形状：{bangumi_id: anidb_id}。
type FileOverrides struct {
	path string
	mu   sync.Mutex
}

var _ Overrides = (*FileOverrides)(nil)

// NewFileOverrides 创建以 path 为持久化文件的覆盖表。文件不存在时按空表处理。
func NewFileOverrides(path string) *FileOverrides {
	return &FileOverrides{path: path}
}

// Get 查询一条覆盖记录。
func (o *FileOverrides) Get(bangumiID int) (int, bool) {
	o.mu.Lock()
	defer o.mu.Unlock()
	m, err := o.readLocked()
	if err != nil {
		return 0, false
	}
	id, ok := m[bangumiID]
	return id, ok
}

// Put 写入并原子落盘。
func (o *FileOverrides) Put(bangumiID, anidbID int) error {
	if bangumiID <= 0 || anidbID <= 0 {
		return nil // 非法值静默忽略
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	m, err := o.readLocked()
	if err != nil {
		return err
	}
	m[bangumiID] = anidbID
	return o.writeLocked(m)
}

// All 返回全表副本。
func (o *FileOverrides) All() map[int]int {
	o.mu.Lock()
	defer o.mu.Unlock()
	m, err := o.readLocked()
	if err != nil {
		return map[int]int{}
	}
	return m
}

// readLocked 读取文件（调用方需持锁）。
func (o *FileOverrides) readLocked() (map[int]int, error) {
	data, err := os.ReadFile(o.path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[int]int{}, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return map[int]int{}, nil
	}
	var m map[int]int
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// writeLocked 写入文件（先临时文件再 rename，避免写一半损坏）。
func (o *FileOverrides) writeLocked(m map[int]int) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(o.path), 0o755); err != nil {
		return err
	}
	tmp := o.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, o.path)
}
