package core

import (
	"path/filepath"
	"time"

	"seanime/internal/api/animap"
	"seanime/internal/api/animap/dataset"
	"seanime/internal/media"
	"seanime/internal/platforms/bangumi_platform"
	"seanime/internal/platforms/platform"
	"seanime/internal/util"
)

// Phase 4：bangumi→anidb ID 映射服务的生产装配（M1-04 检账）。
//
// 设计要点：
//   - 异步加载：offline 数据集 ~几十 MB 下载/解析不能阻塞启动，
//     就绪前 mappingService 为 nil，平台层 nil 守卫自然降级；
//   - 数据落点：数据集缓存进 Cache.Dir（可再生），队列与人工覆盖表
//     进 Data.AppDataDir（用户数据，不清理）；
//   - 失败降级：下载失败（离线/断网）仅 warn，映射服务缺位不影响其他功能；
//   - 重登装配：AuthenticateToBangumi 会新建平台实例，须重新触发装配。

// MappingDatasetStatus 映射数据集装载状态（mapping/status 端点消费）。
type MappingDatasetStatus struct {
	Loaded    bool      `json:"loaded"`
	Entries   int       `json:"entries"`
	FetchedAt time.Time `json:"fetchedAt"`
	FromCache bool      `json:"fromCache"`
	Error     string    `json:"error,omitempty"`
}

// initMappingService 异步装配映射服务并注入平台（非阻塞，幂等性由调用方保证）。
// 非 Bangumi 平台实现（离线/模拟平台）静默跳过。
func (a *App) initMappingService(p platform.Platform) {
	bp, ok := p.(*bangumi_platform.BangumiPlatform)
	if !ok {
		return
	}
	go func() {
		defer util.HandlePanicThen(func() {})

		start := time.Now()
		entries, fromCache, err := a.loadAnimapDataset()
		if err != nil {
			a.Logger.Warn().Err(err).Msg("mapping: dataset load failed, bangumi→anidb mapping disabled")
			a.MappingMu.Lock()
			a.MappingStatus = MappingDatasetStatus{Error: err.Error()}
			a.MappingMu.Unlock()
			return
		}

		dataDir := a.Config.Data.AppDataDir
		queuePath := filepath.Join(dataDir, "mapping-queue.json")
		overridesPath := filepath.Join(dataDir, "mapping-overrides.json")

		resolver := media.NewAnimapResolver(entries, media.WithQueue(media.NewFileQueue(queuePath)))
		ov := media.NewFileOverrides(overridesPath)
		resolver.SetOverrides(ov.All())

		a.MappingOverridesPath = overridesPath
		a.MappingQueuePath = queuePath
		bp.SetMappingService(resolver)
		a.MappingServiceRef.Set(resolver)

		a.MappingMu.Lock()
		a.MappingStatus = MappingDatasetStatus{
			Loaded:    true,
			Entries:   len(entries),
			FetchedAt: time.Now(),
			FromCache: fromCache,
		}
		a.MappingMu.Unlock()

		a.Logger.Info().
			Int("entries", len(entries)).
			Bool("fromCache", fromCache).
			Dur("took", time.Since(start)).
			Msg("mapping: bangumi→anidb mapping service ready")
	}()
}

// loadAnimapDataset 装载 offline 数据集（带 7 天 TTL 缓存；走配置代理）。
func (a *App) loadAnimapDataset() (entries []animap.Anime, fromCache bool, err error) {
	return dataset.LoadWithStatus(dataset.Options{
		CachePath: filepath.Join(a.Config.Cache.Dir, "animap-dataset.jsonl"),
		ProxyURL:  a.Config.Server.ProxyURL,
	})
}
