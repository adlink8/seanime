package handlers

import (
	"seanime/internal/core"
	"seanime/internal/customsource"
	"seanime/internal/media"
	"seanime/internal/platforms/bangumi_platform"

	"github.com/labstack/echo/v4"
)

// Phase 4 映射队列清偿工具（M4-01）：
// 提供映射服务状态、待映射队列查看、人工确认清偿、存量库覆盖率实测。
// 自用管理端点——前端无 UI，curl/浏览器直接消费。

// HandleGetMappingStatus
//
//	@summary mapping service status
//	@route /api/v1/mapping/status [GET]
func (h *Handler) HandleGetMappingStatus(c echo.Context) error {
	type statusResponse struct {
		Ready           bool                              `json:"ready"`   // 映射服务是否已装配
		Dataset         *core.MappingDatasetStatus        `json:"dataset"` // 数据集装载状态
		QueueSize       int                               `json:"queueSize"`
		QueueAvailable  bool                              `json:"queueAvailable"`
		Overrides       map[int]int                       `json:"overrides"`
		ReverseMappings int                               `json:"reverseMappings"` // 反向索引已确认条数
	}

	resp := statusResponse{Overrides: map[int]int{}}

	if svc := h.App.MappingServiceRef.Get(); svc != nil {
		resp.Ready = true
		if r, ok := svc.(*media.AnimapResolver); ok {
			if items, ok2 := r.QueueSnapshot(); ok2 {
				resp.QueueSize = len(items)
				resp.QueueAvailable = true
			}
			resp.Overrides = r.OverridesSnapshot()
		}
	}

	h.App.MappingMu.Lock()
	st := h.App.MappingStatus
	h.App.MappingMu.Unlock()
	resp.Dataset = &st

	// 反向索引条数（平台未离线时）
	if bp, ok := h.App.AnilistPlatformRef.Get().(*bangumi_platform.BangumiPlatform); ok && bp != nil {
		if idx := bp.ReverseIDIndex(); idx != nil {
			resp.ReverseMappings = idx.Len()
		}
	}

	return h.RespondWithData(c, resp)
}

// HandleGetMappingQueue
//
//	@summary list unresolved mapping queue items
//	@route /api/v1/mapping/queue [GET]
func (h *Handler) HandleGetMappingQueue(c echo.Context) error {
	svc := h.App.MappingServiceRef.Get()
	if svc == nil {
		return h.RespondWithData(c, []media.Unresolved{})
	}
	r, ok := svc.(*media.AnimapResolver)
	if !ok {
		return h.RespondWithData(c, []media.Unresolved{})
	}
	items, ok2 := r.QueueSnapshot()
	if !ok2 {
		return h.RespondWithError(c, echo.NewHTTPError(400, "mapping queue not available"))
	}
	if items == nil {
		items = []media.Unresolved{}
	}
	return h.RespondWithData(c, items)
}

type mappingResolveBody struct {
	BangumiID int `json:"bangumiId"`
	AniDBID   int `json:"anidbId"`
}

// HandlePostMappingResolve
//
//	@summary manually confirm a bangumi→anidb mapping (clears queue entry)
//	@route /api/v1/mapping/resolve [POST]
func (h *Handler) HandlePostMappingResolve(c echo.Context) error {
	var body mappingResolveBody
	if err := c.Bind(&body); err != nil {
		return h.RespondWithError(c, err)
	}
	if body.BangumiID <= 0 || body.AniDBID <= 0 {
		return h.RespondWithError(c, echo.NewHTTPError(400, "bangumiId and anidbId must be positive"))
	}

	svc := h.App.MappingServiceRef.Get()
	if svc == nil {
		return h.RespondWithError(c, echo.NewHTTPError(400, "mapping service not ready"))
	}
	r, ok := svc.(*media.AnimapResolver)
	if !ok {
		return h.RespondWithError(c, echo.NewHTTPError(400, "mapping service not available"))
	}

	// 1. 持久化覆盖表（重启后仍生效）
	ov := media.NewFileOverrides(h.App.MappingOverridesPath)
	if err := ov.Put(body.BangumiID, body.AniDBID); err != nil {
		return h.RespondWithError(c, err)
	}
	// 2. 运行时立即生效
	r.ApplyOverride(body.BangumiID, body.AniDBID)
	// 3. 清偿出队
	r.QueueRemove(body.BangumiID)
	// 4. 反向索引回填
	if bp, ok := h.App.AnilistPlatformRef.Get().(*bangumi_platform.BangumiPlatform); ok && bp != nil {
		bp.ConfirmMapping(body.BangumiID, body.AniDBID)
	}

	return h.RespondWithData(c, map[string]bool{"ok": true})
}

// HandlePostMappingCoverage
//
//	@summary measure bangumi→anidb mapping coverage over the current library
//	@desc resolves every library entry (offline: MAL direct lookup + name matching + overrides)
//	@route /api/v1/mapping/coverage [POST]
func (h *Handler) HandlePostMappingCoverage(c echo.Context) error {
	svc := h.App.MappingServiceRef.Get()
	if svc == nil {
		return h.RespondWithError(c, echo.NewHTTPError(400, "mapping service not ready"))
	}

	collection, err := h.App.GetAnimeCollection(false)
	if err != nil {
		return h.RespondWithError(c, err)
	}
	if collection == nil {
		return h.RespondWithError(c, echo.NewHTTPError(400, "anime collection not loaded"))
	}

	type unresolvedItem struct {
		BangumiID int    `json:"bangumiId"`
		Name      string `json:"name"`
	}
	type coverageResponse struct {
		Total      int              `json:"total"`
		Resolved   int              `json:"resolved"`
		Coverage   float64          `json:"coverage"` // 0-100
		Unresolved []unresolvedItem `json:"unresolved"`
	}

	seen := map[int]bool{}
	resp := coverageResponse{Unresolved: []unresolvedItem{}}

	for _, list := range collection.GetMediaListCollection().GetLists() {
		for _, entry := range list.GetEntries() {
			m := entry.GetMedia()
			if m == nil {
				continue
			}
			id := m.GetID()
			if id <= 0 || seen[id] || customsource.IsExtensionId(id) {
				continue // 自定义源条目无 Bangumi 锚点，跳过
			}
			seen[id] = true

			// 名称集：Bangumi 提供日文原名 + 中文名；synonyms 为补充别名
			names := media.NameSet{Chinese: m.NameCN}
			if t := m.Title; t != nil {
				names.Native = derefStr(t.Native)
				names.English = derefStr(t.English)
				names.Romaji = derefStr(t.Romaji)
			}
			for _, syn := range m.Synonyms {
				if syn != nil && *syn != "" {
					names.Aliases = append(names.Aliases, *syn)
				}
			}

			resp.Total++
			// MAL 直查优先（零成本、零误配）；失败回退名称匹配
			if m.IDMal != nil && *m.IDMal > 0 {
				if _, ok := svc.ResolveMalToAniDB(*m.IDMal); ok {
					resp.Resolved++
					continue
				}
			}
			if _, ok := svc.ResolveBangumiToAniDB(id, names); ok {
				resp.Resolved++
				continue
			}
			resp.Unresolved = append(resp.Unresolved, unresolvedItem{
				BangumiID: id,
				Name:      displayName(names, m.NameCN),
			})
		}
	}

	if resp.Total > 0 {
		resp.Coverage = float64(resp.Resolved) / float64(resp.Total) * 100
	}
	return h.RespondWithData(c, resp)
}

// displayName 覆盖率报告里的显示名（中文名优先，退回日文原名）。
func displayName(names media.NameSet, nameCN string) string {
	if nameCN != "" {
		return nameCN
	}
	if names.Native != "" {
		return names.Native
	}
	if names.English != "" {
		return names.English
	}
	return ""
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
