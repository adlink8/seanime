package handlers

import (
	"context"
	"errors"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
	"seanime/internal/api/asmr"
	asmrlib "seanime/internal/asmr"
)

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// asmr client 装配：asmr.one 搜索域无需用户鉴权（token 空串），
// 代理沿用启动时读取的 config server.proxyurl（api.asmr.one 直连被墙，必须显式代理）。
// 与 bangumi 不同（token 随登录变化），asmr client 全生命周期不变，用进程级单例即可。
var (
	asmrClientOnce sync.Once
	asmrClient     *asmr.Client

	asmrDownloaderOnce sync.Once
	asmrDownloader     *asmrlib.Downloader
)

func (h *Handler) getAsmrClient() (*asmr.Client, error) {
	asmrClientOnce.Do(func() {
		opts := []asmr.Option{asmr.WithProxyURL(h.App.Config.Server.ProxyURL)}
		if h.App.FileCacher != nil {
			opts = append(opts, asmr.WithFileCache(h.App.FileCacher))
		}
		asmrClient = asmr.New("", opts...)
	})
	if asmrClient == nil {
		return nil, errors.New("asmr client not available")
	}
	return asmrClient, nil
}

// getAsmrDownloader 进程级下载器单例（依赖 asmr client + ws 事件管理器）。
func (h *Handler) getAsmrDownloader() (*asmrlib.Downloader, error) {
	asmrDownloaderOnce.Do(func() {
		client, err := h.getAsmrClient()
		if err != nil {
			return
		}
		asmrDownloader = asmrlib.NewDownloader(client, h.App.Config.Asmr.LocalDir, h.App.WSEventManager, h.App.Logger)
	})
	if asmrDownloader == nil {
		return nil, errors.New("asmr downloader not available")
	}
	return asmrDownloader, nil
}

// HandleAsmrSearch
//
//	@summary searches asmr.one for voice works (音声作品).
//	@desc asmr.one 搜索域：无需鉴权，走显式代理。请求形状见契约第 4 节。
//	@route /api/v1/asmr/search [POST]
//	@returns asmr.Asmr_SearchResult
func (h *Handler) HandleAsmrSearch(c echo.Context) error {

	type body struct {
		Keyword  string `json:"keyword"`
		Order    string `json:"order,omitempty"`    // "dd"|"dl"|"dc"|"publish_date"
		Page     *int   `json:"page,omitempty"`     // 1 起
		PerPage  *int   `json:"perPage,omitempty"`  // 默认 20
		Subtitle string `json:"subtitle,omitempty"` // "none"|"jp"|"zh"
	}

	p := new(body)
	if err := c.Bind(p); err != nil {
		return h.RespondWithError(c, err)
	}

	page := 1
	if p.Page != nil && *p.Page > 0 {
		page = *p.Page
	}
	perPage := 20
	if p.PerPage != nil && *p.PerPage > 0 {
		perPage = *p.PerPage
	}

	client, err := h.getAsmrClient()
	if err != nil {
		return h.RespondWithError(c, err)
	}

	res, err := client.Search(c.Request().Context(), asmr.SearchParams{
		Keyword:  p.Keyword,
		Order:    p.Order,
		Page:     page,
		PerPage:  perPage,
		Subtitle: p.Subtitle,
	})
	if err != nil {
		return h.RespondWithError(c, err)
	}

	return h.RespondWithData(c, res)
}

// HandleAsmrPopular
//
//	@summary returns popular asmr.one voice works.
//	@desc 近似实现：/api/popular 需鉴权（实测 401），改用无关键词默认搜索第 1 页。
//	@route /api/v1/asmr/popular [GET]
//	@returns asmr.Asmr_SearchResult
func (h *Handler) HandleAsmrPopular(c echo.Context) error {

	client, err := h.getAsmrClient()
	if err != nil {
		return h.RespondWithError(c, err)
	}

	res, err := client.Popular(c.Request().Context(), 20)
	if err != nil {
		return h.RespondWithError(c, err)
	}

	return h.RespondWithData(c, res)
}

// HandleAsmrWork
//
//	@summary returns asmr.one work details with the track tree.
//	@desc 内部组合 /api/workInfo/{id} 与 /api/tracks/{id}?v=2（folder 用 tracks 嵌套）。
//	@desc {id} 为 asmr.one 数字 work id（"RJ..." 会被上游 404）。
//	@route /api/v1/asmr/work/{id} [GET]
//	@returns asmr.Asmr_WorkDetail
func (h *Handler) HandleAsmrWork(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return h.RespondWithError(c, errors.New("invalid work id"))
	}

	client, err := h.getAsmrClient()
	if err != nil {
		return h.RespondWithError(c, err)
	}

	work, err := client.WorkInfo(c.Request().Context(), id)
	if err != nil {
		return h.RespondWithError(c, err)
	}

	tracks, err := client.Tracks(c.Request().Context(), id)
	if err != nil {
		return h.RespondWithError(c, err)
	}

	return h.RespondWithData(c, &asmr.Asmr_WorkDetail{Asmr_Work: *work, Tracks: tracks})
}

// HandleAsmrLibrary
//
//	@summary scans the local asmr directory and returns library entries (契约 §1).
//	@route /api/v1/asmr/library [GET]
//	@returns asmrlib.LibraryResponse
func (h *Handler) HandleAsmrLibrary(c echo.Context) error {
	client, err := h.getAsmrClient()
	if err != nil {
		return h.RespondWithError(c, err)
	}

	scanner := asmrlib.NewScanner(h.App.Config.Asmr.LocalDir, client, h.App.Database, h.App.Logger)
	entries, err := scanner.Scan(c.Request().Context())
	if err != nil {
		return h.RespondWithError(c, err)
	}

	return h.RespondWithData(c, &asmrlib.LibraryResponse{
		LocalDir: h.App.Config.Asmr.LocalDir,
		Entries:  entries,
	})
}

// HandleAsmrLibraryWork
//
//	@summary returns merged local+online track tree for a single RJ work (契约 §2).
//	@route /api/v1/asmr/library/work/:rjId [GET]
//	@returns asmrlib.LibraryWork
func (h *Handler) HandleAsmrLibraryWork(c echo.Context) error {
	rjID := c.Param("rjId")
	if rjID == "" {
		return h.RespondWithError(c, errors.New("rjId is required"))
	}

	client, err := h.getAsmrClient()
	if err != nil {
		return h.RespondWithError(c, err)
	}

	scanner := asmrlib.NewScanner(h.App.Config.Asmr.LocalDir, client, h.App.Database, h.App.Logger)
	work, ok := scanner.GetWork(c.Request().Context(), rjID)
	if !ok {
		return h.RespondWithError(c, errors.New("work not found locally or online: "+rjID))
	}

	return h.RespondWithData(c, work)
}

// HandleAsmrTrackProgress
//
//	@summary records a track's completed state and returns the listened count (契约 §4).
//	@route /api/v1/asmr/track/progress [POST]
//	@returns { ok, listenedCount }
func (h *Handler) HandleAsmrTrackProgress(c echo.Context) error {
	type body struct {
		RjID      string `json:"rjId"`
		TrackPath string `json:"trackPath"`
		Completed bool   `json:"completed"`
	}
	p := new(body)
	if err := c.Bind(p); err != nil {
		return h.RespondWithError(c, err)
	}
	if p.RjID == "" || p.TrackPath == "" {
		return h.RespondWithError(c, errors.New("rjId and trackPath are required"))
	}

	listenedCount, err := h.App.Database.UpsertAsmrTrackState(p.RjID, p.TrackPath, p.Completed)
	if err != nil {
		return h.RespondWithError(c, err)
	}

	return h.RespondWithData(c, map[string]any{
		"ok":           true,
		"listenedCount": listenedCount,
	})
}

// HandleAsmrLibraryFavorite
//
//	@summary toggles favorite for a work locally; syncs to cloud if account configured (契约 §5).
//	@route /api/v1/asmr/library/favorite [POST]
//	@returns { ok, syncedToCloud }
func (h *Handler) HandleAsmrLibraryFavorite(c echo.Context) error {
	type body struct {
		RjID    string `json:"rjId"`
		Favorite bool   `json:"favorite"`
	}
	p := new(body)
	if err := c.Bind(p); err != nil {
		return h.RespondWithError(c, err)
	}
	if p.RjID == "" {
		return h.RespondWithError(c, errors.New("rjId is required"))
	}

	if err := h.App.Database.UpsertAsmrWorkState(p.RjID, p.Favorite); err != nil {
		return h.RespondWithError(c, err)
	}

	syncedToCloud := false
	if h.App.Config.Asmr.Name != "" && h.App.Config.Asmr.Password != "" {
		if synced, serr := h.syncCloudFavorite(c.Request().Context(), p.RjID, p.Favorite); serr == nil {
			syncedToCloud = synced
		}
		// 云写失败不影响本地结果（契约 §5）
	}

	return h.RespondWithData(c, map[string]any{
		"ok":            true,
		"syncedToCloud": syncedToCloud,
	})
}

// syncCloudFavorite 登录并写云端收藏标记，再以 read-back 校验是否真正生效。
// 返回 (synced, nil)：synced 表示云端 marked 列表状态与期望一致（写端点/字段未确认时为 false）。
func (h *Handler) syncCloudFavorite(ctx context.Context, rjID string, favorite bool) (bool, error) {
	client, err := h.getAsmrClient()
	if err != nil {
		return false, err
	}
	if _, err := client.Login(ctx, h.App.Config.Asmr.Name, h.App.Config.Asmr.Password); err != nil {
		return false, err
	}
	workID, werr := h.resolveWorkID(ctx, client, rjID)
	if werr != nil {
		return false, werr
	}
	if err := client.SaveReview(ctx, workID, favorite, false); err != nil {
		// 写失败：降级，syncedToCloud=false（不影响本地收藏）
		return false, nil
	}
	// read-back 校验：云端 marked 列表是否反映期望状态。
	fav, ferr := client.CloudFavorites(ctx)
	if ferr != nil {
		return false, nil
	}
	return containsKey(fav, rjID) == favorite, nil
}

func containsKey(m map[string]struct{}, k string) bool {
	_, ok := m[k]
	return ok
}

// resolveWorkID 通过 Search 精确匹配 rjId 取得数字 work id。
func (h *Handler) resolveWorkID(ctx context.Context, client *asmr.Client, rjID string) (int, error) {
	res, err := client.Search(ctx, asmr.SearchParams{Keyword: rjID, Page: 1, PerPage: 20})
	if err != nil {
		return 0, err
	}
	for _, w := range res.Works {
		if w.RjID == rjID {
			return strconv.Atoi(w.ID)
		}
	}
	return 0, errors.New("asmr: 未找到匹配 work: " + rjID)
}

// HandleAsmrCloud
//
//	@summary returns cloud favorites/listening rjId sets (契约 §6).
//	@route /api/v1/asmr/cloud [GET]
//	@returns { favorites, listening, configured }
func (h *Handler) HandleAsmrCloud(c echo.Context) error {
	if h.App.Config.Asmr.Name == "" || h.App.Config.Asmr.Password == "" {
		return h.RespondWithData(c, map[string]any{
			"favorites":  []string{},
			"listening":  []string{},
			"configured": false,
		})
	}

	client, err := h.getAsmrClient()
	if err != nil {
		return h.RespondWithError(c, err)
	}
	if _, err := client.Login(c.Request().Context(), h.App.Config.Asmr.Name, h.App.Config.Asmr.Password); err != nil {
		return h.RespondWithError(c, err)
	}

	fav, err := client.CloudFavorites(c.Request().Context())
	if err != nil {
		return h.RespondWithError(c, err)
	}
	lis, err := client.CloudListening(c.Request().Context())
	if err != nil {
		return h.RespondWithError(c, err)
	}

	return h.RespondWithData(c, map[string]any{
		"favorites":  sortedKeys(fav),
		"listening":  sortedKeys(lis),
		"configured": true,
	})
}

// HandleAsmrDownload
//
//	@summary downloads a work's audio tracks asynchronously (契约 §7).
//	@route /api/v1/asmr/download [POST]
//	@returns { ok }
func (h *Handler) HandleAsmrDownload(c echo.Context) error {
	type body struct {
		WorkID string `json:"workId"`
		RjID   string `json:"rjId"`
	}
	p := new(body)
	if err := c.Bind(p); err != nil {
		return h.RespondWithError(c, err)
	}
	if p.WorkID == "" || p.RjID == "" {
		return h.RespondWithError(c, errors.New("workId and rjId are required"))
	}
	workID, err := strconv.Atoi(p.WorkID)
	if err != nil || workID <= 0 {
		return h.RespondWithError(c, errors.New("invalid workId"))
	}

	dl, err := h.getAsmrDownloader()
	if err != nil {
		return h.RespondWithError(c, err)
	}

	// 异步执行，立即返回（202 语义）。
	// 注意：不能用 c.Request().Context()——handler 返回后该 context 即被取消，
	// 后台下载的第一个 HTTP 请求就会 context canceled。
	// ctx 的 cancel 由 goroutine 退出时调用（Download 返回后），handler 里不能 defer cancel，
	// 否则 handler 一返回就会把后台下载取消掉。
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		dl.Download(ctx, workID, p.RjID)
	}()

	return h.RespondWithData(c, map[string]any{"ok": true})
}

// sortedKeys 将 map[rjId]struct{} 转为排序后的字符串切片（稳定响应）。
func sortedKeys(m map[string]struct{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	// 稳定排序（string 默认字典序）
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j] < out[i] {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
}
