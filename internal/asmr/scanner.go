package asmr

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"

	apiasmr "seanime/internal/api/asmr"
	"seanime/internal/database/db"
	"github.com/rs/zerolog"
)

// rjDirPattern 一级子目录名匹配：大小写不敏感 RJ + 5~8 位数字（契约 §1）。
//
// 注意：Bash 工具会把 "\\" 折叠成 "\"，故正则全部用字符类而非反斜杠转义写法。
var rjDirPattern = regexp.MustCompile(`(?i)^RJ(\d{5,8})`)

// NormalizeRJ 从目录名提取归一化 rjId（"RJ"+原始数字，不补零）。
// 例：RJ01234567 → "RJ01234567"；rj289543 → "RJ289543"。
// 数字段严格 5~8 位：若 RJ 后紧跟 9 位及以上数字（如 RJ123456789），视为不匹配（契约 §1 范围）。
// 不匹配返回 ("", false)。
func NormalizeRJ(dirName string) (string, bool) {
	m := rjDirPattern.FindStringIndex(dirName)
	if m == nil {
		return "", false
	}
	matchedEnd := m[1]
	// 拒绝 RJ 后跟随 9+ 位数字：匹配到的数字段之后若仍是数字，则超出 5~8 位范围
	if matchedEnd < len(dirName) {
		next := dirName[matchedEnd]
		if next >= '0' && next <= '9' {
			return "", false
		}
	}
	// matched = "RJ" + 数字，取数字部分
	digits := dirName[m[0]+2 : m[1]]
	return "RJ" + digits, true
}

// Scanner 音声本地库扫描器（契约 §1）。即时扫描、不建大表；状态查 DB。
type Scanner struct {
	localDir string
	client   *apiasmr.Client
	database *db.Database
	logger   *zerolog.Logger
}

// NewScanner 构造扫描器。client 为 nil 时在线元数据全部降级（title=目录名）。
func NewScanner(localDir string, client *apiasmr.Client, database *db.Database, logger *zerolog.Logger) *Scanner {
	l := zerolog.Nop()
	if logger != nil {
		l = *logger
	}
	// 首扫自动创建扫描根目录（契约 §配置：localDir 允许不存在）
	_ = os.MkdirAll(localDir, 0o755)
	return &Scanner{
		localDir: localDir,
		client:   client,
		database: database,
		logger:   &l,
	}
}

// Scan 遍历 localDir 一级子目录，生成库条目（契约 §1）。
// 在线元数据解析并发 ≤4；单个 RJ 失败不影响整体。
func (s *Scanner) Scan(ctx context.Context) ([]LibraryEntry, error) {
	entries, err := os.ReadDir(s.localDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []LibraryEntry{}, nil
		}
		return nil, err
	}

	type job struct {
		dirName string
	}
	jobs := make([]job, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			if _, ok := NormalizeRJ(e.Name()); ok {
				jobs = append(jobs, job{dirName: e.Name()})
			}
		}
	}
	if len(jobs) == 0 {
		return []LibraryEntry{}, nil
	}

	// 并发 ≤4
	sem := make(chan struct{}, 4)
	var mu sync.Mutex
	out := make([]LibraryEntry, 0, len(jobs))
	var wg sync.WaitGroup

	for _, j := range jobs {
		wg.Add(1)
		go func(j job) {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-ctx.Done():
				return
			}
			entry, err := s.scanOne(ctx, j.dirName)
			if err != nil {
				s.logger.Warn().Err(err).Str("dir", j.dirName).Msg("asmr: 扫描单条失败，跳过")
				return
			}
			mu.Lock()
			out = append(out, entry)
			mu.Unlock()
		}(j)
	}
	wg.Wait()

	return out, nil
}

// scanOne 扫描单个 RJ 目录：在线元数据 + 本地音轨统计 + DB 状态。
func (s *Scanner) scanOne(ctx context.Context, dirName string) (LibraryEntry, error) {
	rjID, _ := NormalizeRJ(dirName)
	rjDir := filepath.Join(s.localDir, dirName)

	files, err := collectAudioFiles(rjDir)
	if err != nil {
		return LibraryEntry{}, err
	}
	var totalSize int64
	for _, f := range files {
		totalSize += f.size
	}

	entry := LibraryEntry{
		RjID:           rjID,
		Title:          dirName, // 降级默认：title=目录名
		Cvs:            []string{},
		Tags:           []string{},
		TrackCount:     len(files),
		TotalSizeBytes: totalSize,
	}

	// 在线元数据（失败降级，不影响本地结果）
	if s.client != nil {
		if work, werr := s.fetchWork(ctx, rjID); werr == nil && work != nil {
			entry.Title = work.Title
			entry.Circle = work.Circle
			entry.CoverURL = work.CoverURL
			entry.Nsfw = work.Nsfw
			entry.Rating = work.Rating
			entry.ReleaseDate = work.ReleaseDate
			entry.Cvs = work.Cvs
			entry.Tags = work.Tags
		}
	}

	// DB 状态
	if s.database != nil {
		if ws, ok := s.database.GetAsmrWorkState(rjID); ok {
			entry.IsFavorite = ws.Favorite
		}
		if states, serr := s.database.GetAsmrTrackStates(rjID); serr == nil {
			n := 0
			for _, st := range states {
				if st.Completed {
					n++
				}
			}
			entry.ListenedCount = n
		}
	}

	return entry, nil
}

// fetchWork 搜索并精确匹配 rjId → 返回 WorkInfo（nil 表示降级）。
func (s *Scanner) fetchWork(ctx context.Context, rjID string) (*apiasmr.Asmr_Work, error) {
	res, err := s.client.Search(ctx, apiasmr.SearchParams{Keyword: rjID, Page: 1, PerPage: 20})
	if err != nil {
		return nil, err
	}
	for _, w := range res.Works {
		if w.RjID == rjID {
			return &w, nil
		}
	}
	return nil, nil
}

// findWorkDir 在 localDir 下找到归一化 rjId 匹配的实际目录名。
// 用于 /asmr/library/work/{rjId} 定位本地树根。未找到返回 ("", false)。
func (s *Scanner) findWorkDir(rjID string) (string, bool) {
	entries, err := os.ReadDir(s.localDir)
	if err != nil {
		return "", false
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if norm, ok := NormalizeRJ(e.Name()); ok && norm == rjID {
			return e.Name(), true
		}
	}
	return "", false
}

// GetWork 构造单个作品的合并详情（契约 §2）：本地树优先 + 在线树补充。
// 本地目录不存在且在线无数据时返回 (nil, false)（handler 据此 404）。
func (s *Scanner) GetWork(ctx context.Context, rjID string) (*LibraryWork, bool) {
	dirName, hasLocal := s.findWorkDir(rjID)
	rjDir := ""
	if hasLocal {
		rjDir = filepath.Join(s.localDir, dirName)
	}

	// 本地树
	var localTree []apiasmr.Asmr_Track
	localBasenames := map[string]struct{}{}
	if hasLocal {
		files, _ := collectAudioFiles(rjDir)
		localTree, localBasenames = buildLocalTrackTree(rjDir, files)
	}

	// 在线元数据 + 在线音轨树
	var entry LibraryEntry
	entry.RjID = rjID
	entry.Title = rjID
	entry.Cvs = []string{}
	entry.Tags = []string{}

	// 本地文件计数与体积（仅本地存在时有意义）
	if hasLocal {
		if files, _ := collectAudioFiles(rjDir); true {
			var sz int64
			for _, f := range files {
				sz += f.size
			}
			entry.TrackCount = len(files)
			entry.TotalSizeBytes = sz
		}
	}

	var onlineTree []apiasmr.Asmr_Track
	if s.client != nil {
		if work, _ := s.fetchWork(ctx, rjID); work != nil {
			entry.Title = work.Title
			entry.Circle = work.Circle
			entry.CoverURL = work.CoverURL
			entry.Nsfw = work.Nsfw
			entry.Rating = work.Rating
			entry.ReleaseDate = work.ReleaseDate
			entry.Cvs = work.Cvs
			entry.Tags = work.Tags
			if wid, perr := strconv.Atoi(work.ID); perr == nil && wid > 0 {
				if tr, terr := s.client.Tracks(ctx, wid); terr == nil {
					onlineTree = tr
				}
			}
		}
	}

	// 合并：本地树优先，在线仅在本地没有对应音频时补充（契约 §2）
	tracks := mergeTrackTrees(localTree, onlineTree, localBasenames)

	// DB 状态
	if s.database != nil {
		if ws, ok := s.database.GetAsmrWorkState(rjID); ok {
			entry.IsFavorite = ws.Favorite
		}
		if states, serr := s.database.GetAsmrTrackStates(rjID); serr == nil {
			n := 0
			for _, st := range states {
				if st.Completed {
					n++
				}
			}
			entry.ListenedCount = n
		}
	}

	if !hasLocal && len(onlineTree) == 0 {
		// 既无本地也无在线：无法构造有效作品
		return nil, false
	}

	return &LibraryWork{LibraryEntry: entry, Tracks: tracks}, true
}

// collectAudioFiles 递归收集目录下全部音频文件，返回相对路径（/ 分隔）与绝对路径、大小。
func collectAudioFiles(rjDir string) ([]audioFile, error) {
	var files []audioFile
	err := filepath.WalkDir(rjDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if !isAudioExtension(ext) {
			return nil
		}
		info, ierr := d.Info()
		if ierr != nil {
			return nil
		}
		rel, rerr := filepath.Rel(rjDir, path)
		if rerr != nil {
			rel = path
		}
		files = append(files, audioFile{
			relPath:  filepath.ToSlash(rel),
			absPath:  filepath.Clean(path),
			baseName: filepath.Base(path),
			size:     info.Size(),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

type audioFile struct {
	relPath  string
	absPath  string
	baseName string
	size     int64
}
