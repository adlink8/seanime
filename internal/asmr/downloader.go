package asmr

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	apiasmr "seanime/internal/api/asmr"
	"seanime/internal/events"
	"github.com/rs/zerolog"
)

// Downloader 音声单 work 串行下载器（契约 §7）。
// 走 client 的代理 http client 下载 mediaDownloadUrl，落盘到 {localDir}/{rjId}/{相对路径}。
// 同 rjId 进行中幂等忽略。
type Downloader struct {
	client   *apiasmr.Client
	localDir string
	ws       events.WSEventManagerInterface
	logger   *zerolog.Logger

	mu        sync.Mutex
	inProgress map[string]bool
}

// NewDownloader 构造下载器（进程级单例，handler 持有）。
func NewDownloader(client *apiasmr.Client, localDir string, ws events.WSEventManagerInterface, logger *zerolog.Logger) *Downloader {
	l := zerolog.Nop()
	if logger != nil {
		l = *logger
	}
	return &Downloader{
		client:    client,
		localDir:  localDir,
		ws:        ws,
		logger:    &l,
		inProgress: make(map[string]bool),
	}
}

// Download 异步下载单个 work 的全部 audio 音轨。立即返回（handler 在 goroutine 内调用）。
// 同 rjId 已在下载中则幂等忽略。
func (d *Downloader) Download(ctx context.Context, workID int, rjID string) {
	d.mu.Lock()
	if d.inProgress[rjID] {
		d.mu.Unlock()
		d.logger.Debug().Str("rjId", rjID).Msg("asmr: 下载进行中，幂等忽略")
		return
	}
	d.inProgress[rjID] = true
	d.mu.Unlock()

	defer func() {
		d.mu.Lock()
		delete(d.inProgress, rjID)
		d.mu.Unlock()
	}()

	d.run(ctx, workID, rjID)
}

func (d *Downloader) run(ctx context.Context, workID int, rjID string) {
	if d.client == nil {
		d.emitDone(rjID, false, 1)
		return
	}

	tracks, err := d.client.Tracks(ctx, workID)
	if err != nil {
		d.logger.Error().Err(err).Str("rjId", rjID).Msg("asmr: 获取音轨树失败")
		d.emitDone(rjID, false, 0)
		return
	}

	var tasks []downloadTask
	collectDownloadTasks(tracks, "", &tasks)

	total := len(tasks)
	done := 0
	failed := 0
	destRoot := filepath.Join(d.localDir, rjID)

	for _, t := range tasks {
		select {
		case <-ctx.Done():
			d.logger.Warn().Str("rjId", rjID).Msg("asmr: 下载被取消")
			d.emitDone(rjID, false, failed)
			return
		default:
		}

		dest := filepath.Join(destRoot, filepath.FromSlash(t.relPath))
		// 跳过已存在且非空的
		if fi, ferr := os.Stat(dest); ferr == nil && fi.Size() > 0 {
			done++
			d.emitProgress(rjID, done, total, t.relPath)
			continue
		}
		if mkerr := os.MkdirAll(filepath.Dir(dest), 0o755); mkerr != nil {
			failed++
			d.emitProgress(rjID, done, total, t.relPath)
			continue
		}

		rc, _, derr := d.client.DownloadTrack(ctx, t.url)
		if derr != nil {
			d.logger.Warn().Err(derr).Str("url", t.url).Msg("asmr: 下载音轨失败")
			failed++
			d.emitProgress(rjID, done, total, t.relPath)
			continue
		}
		werr := writeFileAtomic(dest, rc)
		_ = rc.Close()
		if werr != nil {
			d.logger.Warn().Err(werr).Str("dest", dest).Msg("asmr: 写盘失败")
			failed++
		} else {
			done++
		}
		d.emitProgress(rjID, done, total, t.relPath)
	}

	d.emitDone(rjID, failed == 0, failed)
}

// downloadTask 单个待下载音轨任务。
type downloadTask struct {
	relPath string
	url     string
}

// collectDownloadTasks 展平在线音轨树中的 audio 节点为下载任务。
// relPath 由 folder 标题 + 文件标题（含扩展）拼接，保留目录树（契约 §7）。
func collectDownloadTasks(nodes []apiasmr.Asmr_Track, prefix string, out *[]downloadTask) {
	for _, n := range nodes {
		switch n.Type {
		case "folder":
			collectDownloadTasks(n.Tracks, joinRel(prefix, n.Title), out)
		case "audio":
			if n.MediaDownloadURL == "" {
				continue
			}
			*out = append(*out, struct{ relPath, url string }{
				relPath: joinRel(prefix, n.Title),
				url:     n.MediaDownloadURL,
			})
		}
	}
}

func joinRel(prefix, name string) string {
	if prefix == "" {
		return name
	}
	return prefix + "/" + name
}

func writeFileAtomic(dest string, rc io.Reader) error {
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, rc)
	return err
}

func (d *Downloader) emitProgress(rjID string, done, total int, currentTrack string) {
	if d.ws == nil {
		return
	}
	d.ws.SendEvent(events.AsmrDownloadProgress, map[string]any{
		"rjId":        rjID,
		"done":        done,
		"total":       total,
		"currentTrack": currentTrack,
	})
}

func (d *Downloader) emitDone(rjID string, ok bool, failedCount int) {
	if d.ws == nil {
		return
	}
	d.ws.SendEvent(events.AsmrDownloadDone, map[string]any{
		"rjId":       rjID,
		"ok":         ok,
		"failedCount": failedCount,
	})
}

// FilterAudioTitle 暴露给测试：返回去扩展的音频标题（与本地树一致）。
func FilterAudioTitle(title string) string {
	return strings.TrimSuffix(title, filepath.Ext(title))
}
