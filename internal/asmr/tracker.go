package asmr

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	apiasmr "seanime/internal/api/asmr"
	"seanime/internal/database/db"
	"seanime/internal/events"
	"seanime/internal/util"

	"github.com/rs/zerolog"
)

// Tracker 新作跟踪自动拉取（契约 03.2c / M3-04）。
//
// 定位：与动漫 auto-downloader 并列的轻量 HTTP 队列——只借其调度结构模式，不复用其实现。
// v1 裁决（契约 D1–D11）：
//   - 规则 = config 白名单社团名列表（D1），命中判定为纯函数精确匹配
//   - 参数存 config.toml [asmr.tracker]（D2），钳制在配置加载层完成，本模块假定已合法
//   - 去重双层：本地目录已存在（D3①）+ DB asmr_tracker_seen（D3②）
//   - 首跑观察窗：seen 表为空的首轮只记录不下载（D5）
//   - 单轮上限 max_per_run + 磁盘剩余 min_free_gb 保护（D6）
//   - 下载只调既有 Downloader.Download——该方法本身是**同步阻塞**的（downloader.go run() 内联执行），
//     tracker 串行调用即可精确控制 max_per_run，无需完成事件（§5.2 裁决）
//   - v1 不登录、不读 asmr.one 凭据（D10）
//
// Search 走 client 全局 500ms 单令牌 ticker：每轮 ≤ len(circles) 个请求，与用户手动搜索共用限流，可接受。
type Tracker struct {
	cfg        TrackerConfig
	client     *apiasmr.Client
	downloader *Downloader
	database   *db.Database
	ws         events.WSEventManagerInterface
	logger     *zerolog.Logger

	mu      sync.Mutex
	running bool
	lastRun *TrackerRunResult

	stopCh   chan struct{}
	stopOnce sync.Once
}

// TrackerConfig 跟踪器参数（config [asmr.tracker] 映射，加载层已钳制）。
type TrackerConfig struct {
	Enabled        bool
	IntervalMinutes int   // 默认 120，最小 30
	Circles        []string
	MaxPerRun      int   // 默认 2，最小 1
	MinFreeGB      int   // 默认 10
	BackfillDays   int   // 默认 7；0 = 永远只记录不下载（契约 D2）
}

// TrackerRunResult 单轮运行结果（status 端点返回）。
type TrackerRunResult struct {
	RanAt        time.Time `json:"ranAt"`
	FirstRun     bool      `json:"firstRun"`
	Evaluated    int       `json:"evaluated"`
	Invalid      int       `json:"invalid"`
	SkippedLocal int       `json:"skippedLocal"`
	SkippedSeen  int       `json:"skippedSeen"`
	Recorded     int       `json:"recorded"`
	Downloaded   int       `json:"downloaded"`
	DiskStopped  bool      `json:"diskStopped"`
	CircleErrors []string  `json:"circleErrors,omitempty"`
}

// NewTracker 构造跟踪器。依赖与 downloader 同源（进程级单例，handler 持有）。
func NewTracker(cfg TrackerConfig, client *apiasmr.Client, downloader *Downloader, database *db.Database, ws events.WSEventManagerInterface, logger *zerolog.Logger) *Tracker {
	l := zerolog.Nop()
	if logger != nil {
		l = *logger
	}
	return &Tracker{
		cfg:        cfg,
		client:     client,
		downloader: downloader,
		database:   database,
		ws:         ws,
		logger:     &l,
		stopCh:     make(chan struct{}),
	}
}

// Start 启动周期调度 goroutine（Enabled=false 时为 no-op）。
// 失败退避：单轮出现 circle 错误则间隔 ×2，上限 24h；成功轮恢复基准间隔（D4）。
// 配置在进程启动时读取一次，改动 config 需重启生效（v1 裁决，轻量优先）。
func (t *Tracker) Start(ctx context.Context) {
	if t == nil || !t.cfg.Enabled || len(t.cfg.Circles) == 0 {
		t.logger.Info().Msg("asmr tracker: 未启用或无追更社团，调度未启动")
		return
	}
	base := time.Duration(t.cfg.IntervalMinutes) * time.Minute
	go func() {
		defer util.HandlePanicInModuleThen("asmr/tracker.Start", func() {})
		interval := 2 * time.Minute // 首轮避让启动期
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.stopCh:
				return
			case <-time.After(interval):
			}
			res := t.Run(ctx)
			if res != nil && len(res.CircleErrors) > 0 {
				interval = min24h(interval * 2)
			} else {
				interval = base
			}
		}
	}()
	t.logger.Info().Int("circles", len(t.cfg.Circles)).Int("intervalMinutes", t.cfg.IntervalMinutes).Msg("asmr tracker: 调度已启动")
}

// Stop 停止周期调度（幂等）。
func (t *Tracker) Stop() {
	t.stopOnce.Do(func() { close(t.stopCh) })
}

func min24h(d time.Duration) time.Duration {
	const max = 24 * time.Hour
	if d > max {
		return max
	}
	return d
}

// Run 手动/定时触发一轮。进行中重复触发返回 nil（调用方提示幂等忽略，D8）。
func (t *Tracker) Run(ctx context.Context) *TrackerRunResult {
	t.mu.Lock()
	if t.running {
		t.mu.Unlock()
		return nil
	}
	t.running = true
	t.mu.Unlock()
	defer func() {
		t.mu.Lock()
		t.running = false
		t.mu.Unlock()
	}()

	res := &TrackerRunResult{RanAt: time.Now(), CircleErrors: make([]string, 0)}
	if t.client == nil || t.downloader == nil || t.database == nil {
		res.CircleErrors = append(res.CircleErrors, "tracker dependencies not available")
		t.lastRun = res
		return res
	}

	// 首跑观察窗判定（D5）
	seenCount, err := t.database.CountAsmrTrackerSeen()
	if err != nil {
		res.CircleErrors = append(res.CircleErrors, fmt.Sprintf("seen 计数失败: %v", err))
		t.lastRun = res
		return res
	}
	res.FirstRun = seenCount == 0

	// 下载窗口（BackfillDays=0 → 永远只记录，契约 D2）
	recordOnly := res.FirstRun || t.cfg.BackfillDays <= 0
	windowStart := time.Now().AddDate(0, 0, -t.cfg.BackfillDays)

	for _, circle := range t.cfg.Circles {
		select {
		case <-ctx.Done():
			res.CircleErrors = append(res.CircleErrors, "context canceled")
			t.lastRun = res
			return res
		case <-t.stopCh:
			res.CircleErrors = append(res.CircleErrors, "tracker stopped")
			t.lastRun = res
			return res
		default:
		}

		search, err := t.client.Search(ctx, apiasmr.SearchParams{
			Keyword: circle,
			Order:   "publish_date", // → 上游 release，默认最新在前（D9 实测 2026-09-14）
			Page:    1,
			PerPage: 20,
		})
		if err != nil {
			res.CircleErrors = append(res.CircleErrors, fmt.Sprintf("circle %q: search: %v", circle, err))
			continue // 单 circle 失败不阻塞其他 circle（AC-3）
		}

		matched, invalid := filterTrackerCandidates(search.Works, circle, windowStart)
		res.Evaluated += len(search.Works)
		res.Invalid += invalid

		for _, w := range matched {
			if res.Downloaded >= t.cfg.MaxPerRun {
				break // 单轮上限（D6）
			}
			rjID := w.RjID

			// 去重①：本地目录已存在且含文件 → 视为已入档（D3①）
			local, err := hasLocalFiles(filepath.Join(t.downloader.localDir, rjID))
			if err == nil && local {
				res.SkippedLocal++
				_ = t.database.MarkAsmrWorkSeen(rjID) // 顺手记录，防重复评估
				continue
			}

			// 去重②：seen 表（D3②）
			seen, err := t.database.IsAsmrWorkSeen(rjID)
			if err != nil {
				res.CircleErrors = append(res.CircleErrors, fmt.Sprintf("rj %s: seen 查询: %v", rjID, err))
				continue
			}
			if seen {
				res.SkippedSeen++
				continue
			}

			if recordOnly {
				// 首跑观察窗 / backfill=0：只记录不下载（D5）
				if merr := t.database.MarkAsmrWorkSeen(rjID); merr != nil {
					res.CircleErrors = append(res.CircleErrors, fmt.Sprintf("rj %s: seen 写入: %v", rjID, merr))
					continue
				}
				res.Recorded++
				continue
			}

			// 磁盘保护（D6）
			free, err := freeDiskBytes(t.downloader.localDir)
			if err == nil && free < uint64(t.cfg.MinFreeGB)*1024*1024*1024 {
				t.logger.Warn().Str("freeGB", strconv.FormatUint(free/1024/1024/1024, 10)).Msg("asmr tracker: 磁盘剩余低于阈值，本轮停止下发")
				res.DiskStopped = true
				t.emit(res)
				t.lastRun = res
				return res
			}

			// 下载：Downloader.Download 同步阻塞（见类型注释 §5.2），串行保证 max_per_run 精确
			workID, perr := strconv.Atoi(w.ID)
			if perr != nil {
				res.CircleErrors = append(res.CircleErrors, fmt.Sprintf("rj %s: 非法 work id %q", rjID, w.ID))
				continue
			}
			t.logger.Info().Str("rjId", rjID).Str("title", w.Title).Msg("asmr tracker: 命中追更，开始自动拉取")
			t.downloader.Download(ctx, workID, rjID)
			_ = t.database.MarkAsmrWorkSeen(rjID)
			res.Downloaded++
		}
	}

	t.emit(res)
	t.lastRun = res
	t.logger.Info().Int("downloaded", res.Downloaded).Int("recorded", res.Recorded).Int("skippedLocal", res.SkippedLocal).Int("skippedSeen", res.SkippedSeen).Msg("asmr tracker: 单轮完成")
	return res
}

// LastRun 返回上次运行结果（status 端点）。
func (t *Tracker) LastRun() *TrackerRunResult {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.lastRun
}

// Config 返回当前配置（status 端点）。
func (t *Tracker) Config() TrackerConfig {
	return t.cfg
}

// filterTrackerCandidates 命中判定纯函数（D1）：Circle trim 后精确匹配 + Release ≥ windowStart。
// Release 空串/非法日期的条目弃用并计入 invalid，绝不 panic（AC-1）。
// 返回 (matched, invalidCount)。窗口判定含边界（Release == windowStart 命中）。
func filterTrackerCandidates(works []apiasmr.Asmr_Work, circle string, windowStart time.Time) ([]apiasmr.Asmr_Work, int) {
	matched := make([]apiasmr.Asmr_Work, 0)
	invalid := 0
	target := strings.TrimSpace(circle)
	for _, w := range works {
		if strings.TrimSpace(w.Circle) != target {
			continue
		}
		rel, err := time.Parse("2006-01-02", strings.TrimSpace(w.ReleaseDate))
		if err != nil {
			invalid++
			continue
		}
		if rel.Before(windowStart) {
			continue
		}
		matched = append(matched, w)
	}
	return matched, invalid
}

// hasLocalFiles 判断 {localDir}/{rjID} 目录是否存在且任意层级含 ≥1 文件（D3①）。
// 目录不存在返回 (false, nil)（os.Stat 与 ReadDir 的 NotFound 都算"无"）。
func hasLocalFiles(dir string) (bool, error) {
	fi, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if !fi.IsDir() {
		return false, errors.New("asmr tracker: 本地作品路径不是目录")
	}
	found := false
	werr := filepath.WalkDir(dir, func(_ string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			found = true
			return filepath.SkipAll
		}
		return nil
	})
	if werr != nil {
		return false, werr
	}
	return found, nil
}

func (t *Tracker) emit(res *TrackerRunResult) {
	if t.ws == nil {
		return
	}
	t.ws.SendEvent(events.AsmrTrackerEvent, map[string]any{
		"ranAt":       res.RanAt.Format(time.RFC3339),
		"downloaded":  res.Downloaded,
		"recorded":    res.Recorded,
		"diskStopped": res.DiskStopped,
		"errors":      len(res.CircleErrors),
	})
}
