package handlers

import (
	"errors"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
	"seanime/internal/api/asmr"
)

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// asmr client 装配：asmr.one 搜索域无需用户鉴权（token 空串），
// 代理沿用启动时读取的 config server.proxyurl（api.asmr.one 直连被墙，必须显式代理）。
// 与 bangumi 不同（token 随登录变化），asmr client 全生命周期不变，用进程级单例即可。
var (
	asmrClientOnce sync.Once
	asmrClient     *asmr.Client
)

func (h *Handler) getAsmrClient() (*asmr.Client, error) {
	asmrClientOnce.Do(func() {
		asmrClient = asmr.New("", asmr.WithProxyURL(h.App.Config.Server.ProxyURL))
	})
	if asmrClient == nil {
		return nil, errors.New("asmr client not available")
	}
	return asmrClient, nil
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
