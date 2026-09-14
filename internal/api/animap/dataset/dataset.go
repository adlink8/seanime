package dataset

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"seanime/internal/api/animap"

	"github.com/goccy/go-json"
)

// Phase 4 映射服务数据源：manami-project/anime-offline-database JSONL。
//
// 说明：animap 包原生只提供按 ID 单查（内部元数据 API），无全量拉取；
// 而 AnimapResolver 需要全量条目构建反向索引。本装载器复用与
// animeofflinedb 相同的上游数据集（JSONL，每行一个条目），但按映射
// 服务的需要转换：从 sources 提取 anidb/mal ID，title+synonyms 作为
// 多语言标题集。
//
// 缓存策略：下载结果原子落盘 cachePath，TTL 内直接读缓存（数据集
// 每日更新，7 天 TTL 足够自用场景）。

const (
	// DefaultDatabaseURL 上游数据集地址（与 animeofflinedb.DatabaseURL 同源）。
	DefaultDatabaseURL = "https://github.com/manami-project/anime-offline-database/releases/download/latest/anime-offline-database.jsonl"
	// DefaultTTL 缓存有效期。
	DefaultTTL = 7 * 24 * time.Hour

	anidbPrefix = "https://anidb.net/anime/"
	malPrefix   = "https://myanimelist.net/anime/"
)

// Options 装载选项。
type Options struct {
	URL       string        // 数据集 URL（空 = DefaultDatabaseURL）
	CachePath string        // 缓存文件路径（必填）
	ProxyURL  string        // 出站代理（空 = 直连）
	TTL       time.Duration // 缓存有效期（0 = DefaultTTL）
	Client    *http.Client  // 注入自定义 client（测试用）
}

// manamiEntry 上游 JSONL 单行形状（仅取映射所需字段）。
type manamiEntry struct {
	Sources  []string `json:"sources"`
	Title    string   `json:"title"`
	Synonyms []string `json:"synonyms"`
}

// Load 装载全量 animap 条目：优先读缓存，过期/缺失则下载并回写缓存。
// 返回的条目保证 Mappings.AnidbID > 0。
func Load(opts Options) ([]animap.Anime, error) {
	entries, _, err := LoadWithStatus(opts)
	return entries, err
}

// LoadWithStatus 同 Load，额外返回是否命中缓存（status 端点用）。
func LoadWithStatus(opts Options) (entries []animap.Anime, fromCache bool, err error) {
	if opts.CachePath == "" {
		return nil, false, errors.New("animap dataset: cache path is required")
	}
	ttl := opts.TTL
	if ttl == 0 {
		ttl = DefaultTTL
	}
	u := opts.URL
	if u == "" {
		u = DefaultDatabaseURL
	}

	// 缓存新鲜 → 直接读
	if st, statErr := os.Stat(opts.CachePath); statErr == nil && time.Since(st.ModTime()) < ttl {
		if es, parseErr := parseFile(opts.CachePath); parseErr == nil {
			return es, true, nil
		}
		// 缓存损坏 → 重新下载
	}

	tmp := opts.CachePath + ".dl"
	if err := download(u, tmp, opts); err != nil {
		return nil, false, fmt.Errorf("animap dataset: download failed: %w", err)
	}
	entries, err = parseFile(tmp)
	if err != nil {
		os.Remove(tmp)
		return nil, false, err
	}
	if err := os.MkdirAll(filepath.Dir(opts.CachePath), 0o755); err != nil {
		os.Remove(tmp)
		return nil, false, err
	}
	if err := os.Rename(tmp, opts.CachePath); err != nil {
		os.Remove(tmp)
		return nil, false, err
	}
	return entries, false, nil
}

// download 下载 JSONL 到 path（尊重代理配置）。
func download(u, path string, opts Options) error {
	client := opts.Client
	if client == nil {
		tr := http.DefaultTransport.(*http.Transport).Clone()
		if opts.ProxyURL != "" {
			pu, err := url.Parse(opts.ProxyURL)
			if err != nil {
				return err
			}
			tr.Proxy = http.ProxyURL(pu)
		}
		client = &http.Client{Transport: tr, Timeout: 10 * time.Minute}
	}
	resp, err := client.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("status " + resp.Status)
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

// parseFile 解析 JSONL 文件为 animap 条目集合。
func parseFile(path string) ([]animap.Anime, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	entries := make([]animap.Anime, 0, 30000)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		// 首行为数据集元信息（非条目）
		if lineNum == 1 {
			var meta map[string]any
			if json.Unmarshal(scanner.Bytes(), &meta) == nil {
				if _, isEntry := meta["sources"]; !isEntry {
					continue
				}
			}
		}
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var e manamiEntry
		if err := json.Unmarshal(line, &e); err != nil {
			continue // 坏行跳过，不中断
		}
		if a, ok := Convert(&e); ok {
			entries = append(entries, a)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

// Convert 将上游条目转换为 animap.Anime。
// ok=false 表示条目无 anidb ID（映射服务只消费带 anidb 的条目）。
func Convert(e *manamiEntry) (animap.Anime, bool) {
	anidbID := extractID(e.Sources, anidbPrefix)
	if anidbID == 0 {
		return animap.Anime{}, false
	}
	malID := extractID(e.Sources, malPrefix)

	// Titles：title + synonyms 去重（BuildTitleIndex 对 Titles 全值建索引）
	titles := make(map[string]string, 1+len(e.Synonyms))
	addTitle := func(t string) {
		t = strings.TrimSpace(t)
		if t == "" {
			return
		}
		if _, ok := titles[t]; !ok {
			titles[t] = t
		}
	}
	addTitle(e.Title)
	for _, s := range e.Synonyms {
		addTitle(s)
	}
	delete(titles, e.Title) // Title 字段单独存，Titles 里不重复

	return animap.Anime{
		Title:    e.Title,
		Titles:   titles,
		Mappings: &animap.AnimeMapping{AnidbID: anidbID, MalID: malID},
	}, true
}

// extractID 从 sources 前缀中提取数字 ID。
func extractID(sources []string, prefix string) int {
	for _, s := range sources {
		if !strings.HasPrefix(s, prefix) {
			continue
		}
		idStr := s[len(prefix):]
		if idx := strings.IndexAny(idStr, "/?"); idx != -1 {
			idStr = idStr[:idx]
		}
		if id, err := strconv.Atoi(idStr); err == nil {
			return id
		}
	}
	return 0
}
