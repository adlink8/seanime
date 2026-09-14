package bangumi

// adapter.go —— Bangumi 响应 → media 领域类型的适配层（Wave A）。
//
// 依赖方向：bangumi → media（单向，media 不感知 bangumi）。
//
// 收藏状态映射（Bangumi CollectionType → AniList MediaListStatus）：
//   1 想看 → PLANNING
//   2 看过 → COMPLETED
//   3 在看 → CURRENT
//   4 搁置 → PAUSED
//   5 抛弃 → DROPPED
//
// 评分映射：Bangumi rate 0-10 → AniList score 0-100（×10，与
// media.AnimeFromSubject 的 meanScore 换算规则一致）。
//
// 进度映射：EpStatus（动画为已看集数、书籍为已看话数）→ entry.Progress；
// VolStatus（仅书籍）在 AniList 收藏条目形状中无对应字段，暂不承载。
//
// 条目 ID：AniList 收藏条目 ID 为列表条目 ID，Bangumi 无此概念，
// 直接使用 subject_id（与 media_id = Bangumi subject ID 的库层语义一致）。

import (
	"encoding/json"

	"seanime/internal/media"
)

// SubjectToMedia 将 bangumi.Subject 转换为 media.Subject 镜像结构。
// 两个结构字段高度重合，仅形状略有差异（指针 vs 值）。
func SubjectToMedia(s *Subject) *media.Subject {
	if s == nil {
		return nil
	}

	out := &media.Subject{
		ID:       s.ID,
		Type:     s.Type,
		Name:     s.Name,
		NameCN:   s.NameCN,
		Summary:  s.Summary,
		Date:     s.Date,
		Platform: s.Platform,
		Eps:      s.Eps,
		NSFW:     s.Nsfw,
	}

	if s.Images != (SubjectImages{}) {
		out.Images = &media.SubjectImages{
			Large:  s.Images.Large,
			Common: s.Images.Common,
			Medium: s.Images.Medium,
			Small:  s.Images.Small,
			Grid:   s.Images.Grid,
		}
	}
	if s.Rating != nil {
		out.Rating = &media.SubjectRating{
			Rank:  s.Rating.Rank,
			Total: s.Rating.Total,
			Score: s.Rating.Score,
			Count: s.Rating.Count,
		}
	}
	if len(s.Tags) > 0 {
		out.Tags = make([]media.SubjectTag, 0, len(s.Tags))
		for _, t := range s.Tags {
			// Spoiler（3.9d additive）：v0 tags 元素自带
			out.Tags = append(out.Tags, media.SubjectTag{Name: t.Name, Count: t.Count, Spoiler: t.Spoiler})
		}
	}
	// Infobox（3.9d additive）：信息箱条目原样透传（Value 为原始 JSON）
	if len(s.Infobox) > 0 {
		out.Infobox = make([]media.BangumiInfoboxEntry, 0, len(s.Infobox))
		for _, item := range s.Infobox {
			out.Infobox = append(out.Infobox, media.BangumiInfoboxEntry{Key: item.Key, Value: item.Value})
		}
	}

	return out
}

// LegacySubjectToMedia 将 legacy 检索条目（GET /search/subject/{kw}）转换为 media.Subject 镜像结构。
//
// 与 SubjectToMedia 的差异来自契约 §0.1 的实测地面真值：
//   - legacy 条目**没有 `platform` 字段** → Platform 恒为空
//     （故走它推导的 format 只能落到动画 `TV` / 书籍 `BOOK` 兜底值）；
//   - legacy 条目**没有 `tags` 字段** → Tags 恒为空
//     （故 tag / meta_tags 筛选在 legacy 通路无法本地生效，见 handlers 侧降级注释）；
//   - 日期字段名为 `air_date`（v0 为 `date`），映射进 media.Subject.Date。
func LegacySubjectToMedia(s *LegacySubject) *media.Subject {
	if s == nil {
		return nil
	}

	out := &media.Subject{
		ID:      s.ID,
		Type:    s.Type,
		Name:    s.Name,
		NameCN:  s.NameCN,
		Summary: s.Summary,
		Date:    s.AirDate,
		Eps:     s.Eps,
		// Platform / Tags / NSFW：legacy 响应无对应字段，保持零值（不猜测）。
	}

	if s.Images != (SubjectImages{}) {
		out.Images = &media.SubjectImages{
			Large:  s.Images.Large,
			Common: s.Images.Common,
			Medium: s.Images.Medium,
			Small:  s.Images.Small,
			Grid:   s.Images.Grid,
		}
	}
	if s.Rating != nil {
		out.Rating = &media.SubjectRating{
			Rank:  s.Rating.Rank,
			Total: s.Rating.Total,
			Score: s.Rating.Score,
			Count: s.Rating.Count,
		}
	}

	return out
}

// ListStatusFromBangumiType Bangumi 收藏状态枚举 → media.MediaListStatus。
// 未知值返回 nil（保守：不猜测）。
func ListStatusFromBangumiType(t int) *media.MediaListStatus {
	var st media.MediaListStatus
	switch t {
	case CollectionWish:
		st = media.MediaListStatusPlanning
	case CollectionDone:
		st = media.MediaListStatusCompleted
	case CollectionDoing:
		st = media.MediaListStatusCurrent
	case CollectionOnHold:
		st = media.MediaListStatusPaused
	case CollectionDropped:
		st = media.MediaListStatusDropped
	default:
		return nil
	}
	return &st
}

// listMediaFromUserCollection 条目实体转换：subject → Anime/Manga。
func animeFromUserCollection(uc *UserCollection) *media.Anime {
	if uc == nil || uc.Subject == nil {
		return nil
	}
	return media.AnimeFromSubject(SubjectToMedia(uc.Subject))
}

func mangaFromUserCollection(uc *UserCollection) *media.Manga {
	if uc == nil || uc.Subject == nil {
		return nil
	}
	return media.MangaFromSubject(SubjectToMedia(uc.Subject))
}

// scoreToAnilistScale 0-10 → 0-100
func scoreToAnilistScale(rate int) *float64 {
	if rate <= 0 {
		return nil
	}
	s := float64(rate) * 10
	return &s
}

// AnimeCollectionFromUserCollections 将用户动画收藏转换为 media.AnimeCollection。
// 按 Bangumi 收藏状态分组为 lists（与 AniList MediaListCollection 的形状一致）。
// subject 缺失的条目跳过（列表接口正常都会内嵌 subject 精简信息）。
func AnimeCollectionFromUserCollections(res *UserCollectionsResult) *media.AnimeCollection {
	if res == nil {
		return nil
	}

	col := &media.AnimeCollection{
		MediaListCollection: &media.AnimeCollection_MediaListCollection{},
	}
	lists := make(map[media.MediaListStatus]*media.AnimeCollection_MediaListCollection_Lists)

	for i := range res.Data {
		uc := &res.Data[i]
		status := ListStatusFromBangumiType(uc.Type)
		if status == nil {
			continue // 未知状态，保守跳过
		}
		list, ok := lists[*status]
		if !ok {
			name := string(*status)
			isCustom := false
			list = &media.AnimeCollection_MediaListCollection_Lists{
				Status:       status,
				Name:         &name,
				IsCustomList: &isCustom,
			}
			lists[*status] = list
			col.MediaListCollection.Lists = append(col.MediaListCollection.Lists, list)
		}

		entry := &media.AnimeCollection_MediaListCollection_Lists_Entries{
			ID:     uc.SubjectID,
			Status: status,
		}
		if a := animeFromUserCollection(uc); a != nil {
			entry.Media = a
		}
		if uc.EpStatus > 0 {
			progress := uc.EpStatus
			entry.Progress = &progress
		}
		entry.Score = scoreToAnilistScale(uc.Rate)

		list.Entries = append(list.Entries, entry)
	}

	return col
}

// MangaCollectionFromUserCollections 将用户书籍收藏转换为 media.MangaCollection。
// 分组规则同 AnimeCollectionFromUserCollections。
func MangaCollectionFromUserCollections(res *UserCollectionsResult) *media.MangaCollection {
	if res == nil {
		return nil
	}

	col := &media.MangaCollection{
		MediaListCollection: &media.MangaCollection_MediaListCollection{},
	}
	lists := make(map[media.MediaListStatus]*media.MangaCollection_MediaListCollection_Lists)

	for i := range res.Data {
		uc := &res.Data[i]
		status := ListStatusFromBangumiType(uc.Type)
		if status == nil {
			continue
		}
		list, ok := lists[*status]
		if !ok {
			name := string(*status)
			isCustom := false
			list = &media.MangaCollection_MediaListCollection_Lists{
				Status:       status,
				Name:         &name,
				IsCustomList: &isCustom,
			}
			lists[*status] = list
			col.MediaListCollection.Lists = append(col.MediaListCollection.Lists, list)
		}

		entry := &media.MangaCollection_MediaListCollection_Lists_Entries{
			ID:     uc.SubjectID,
			Status: status,
		}
		if m := mangaFromUserCollection(uc); m != nil {
			entry.Media = m
		}
		if uc.EpStatus > 0 {
			progress := uc.EpStatus
			entry.Progress = &progress
		}
		entry.Score = scoreToAnilistScale(uc.Rate)

		list.Entries = append(list.Entries, entry)
	}

	return col
}

//----------------------------------------------------------------------------------------------------------------------
// Wave B additive：关系映射与详情/关系树组装
//
// 关系类型映射（Bangumi RelatedSubject.Type → media.MediaRelation）：
//   1 前传     → PREQUEL
//   2 续集     → SEQUEL
//   3 主线故事 → PARENT
//   4 支线故事 → SIDE_STORY
//   5 角色     → CHARACTER
//   6 其他     → OTHER
// 未知值返回 nil（保守：不猜测）。

// RelationTypeFromBangumiType Bangumi 关系类型 → media.MediaRelation。
func RelationTypeFromBangumiType(t int) *media.MediaRelation {
	var r media.MediaRelation
	switch t {
	case 1:
		r = media.MediaRelationPrequel
	case 2:
		r = media.MediaRelationSequel
	case 3:
		r = media.MediaRelationParent
	case 4:
		r = media.MediaRelationSideStory
	case 5:
		r = media.MediaRelationCharacter
	case 6:
		r = media.MediaRelationOther
	default:
		return nil
	}
	return &r
}

// MediaListStatusToBangumiType media.MediaListStatus → Bangumi 收藏状态枚举。
// REPEATING 在 Bangumi 无对应状态（无重看概念），保守映射为「看过(2)」。
func MediaListStatusToBangumiType(s media.MediaListStatus) (int, bool) {
	switch s {
	case media.MediaListStatusPlanning:
		return CollectionWish, true
	case media.MediaListStatusCompleted, media.MediaListStatusRepeating:
		return CollectionDone, true
	case media.MediaListStatusCurrent:
		return CollectionDoing, true
	case media.MediaListStatusPaused:
		return CollectionOnHold, true
	case media.MediaListStatusDropped:
		return CollectionDropped, true
	default:
		return 0, false
	}
}

// AnimeFromRelatedSubject 由相关条目摘要构建轻量 media.Anime 节点。
//
// 说明：RelatedSubject 只含 id/名称/图片，不含集数、状态等完整字段；
// 为避免「一个关系一个请求」的节流失控（集合大时相关条目可达数十个），
// 关系树节点采用轻量形状，扫描器遍历关系树时会再按需取完整条目。
func AnimeFromRelatedSubject(rs *RelatedSubject) *media.Anime {
	if rs == nil {
		return nil
	}
	a := &media.Anime{
		ID:     rs.SubjectID, // 锚点 ID = 关联条目 ID（rs.ID 是关系记录 ID）
		NameCN: rs.NameCN,
	}
	name := rs.Name
	nameCN := rs.NameCN
	a.Title = &media.Title{
		English: &name, // Bangumi 无英文名，English 槽位放原名（与 titleFromSubject 约定一致）
		Native:  &name,
	}
	if nameCN != "" {
		a.Title.Chinese = &nameCN
		a.Title.UserPreferred = &nameCN
	} else {
		a.Title.UserPreferred = &name
	}
	if rs.Images != nil && rs.Images.Common != "" {
		img := rs.Images.Common
		a.CoverImage = &media.CoverImage{Large: &img, Medium: &img}
	}
	return a
}

// CompleteAnimeFromSubject 组装带关系树的完整 anime 实体
// （subject 详情 + GET /v0/subjects/{id}/subjects 关系列表）。
// 关系节点仅保留 subject_type=2（动画）的条目：CompleteAnime_Relations_Edges.Node
// 类型为 *media.Anime，书籍/音乐等无法保真承载，保守跳过。
func CompleteAnimeFromSubject(s *Subject, related []RelatedSubject) *media.CompleteAnime {
	base := media.AnimeFromSubject(SubjectToMedia(s))
	if base == nil {
		return nil
	}

	ca := &media.CompleteAnime{
		ID:                base.ID,
		IDMal:             base.IDMal,
		SiteURL:           base.SiteURL,
		Status:            base.Status,
		Season:            base.Season,
		SeasonYear:        base.SeasonYear,
		Type:              base.Type,
		Format:            base.Format,
		BannerImage:       base.BannerImage,
		Episodes:          base.Episodes,
		Synonyms:          base.Synonyms,
		IsAdult:           base.IsAdult,
		CountryOfOrigin:   base.CountryOfOrigin,
		MeanScore:         base.MeanScore,
		Description:       base.Description,
		Genres:            base.Genres,
		Duration:          base.Duration,
		Trailer:           base.Trailer,
		Title:             base.Title,
		CoverImage:        base.CoverImage,
		StartDate:         base.StartDate,
		EndDate:           base.EndDate,
		NextAiringEpisode: base.NextAiringEpisode,
		NameCN:            base.NameCN,
		Relations:         &media.CompleteAnime_Relations{},
	}

	for i := range related {
		rs := &related[i]
		if rs.SubjectType != SubjectAnime {
			continue
		}
		relationType := RelationTypeFromBangumiType(rs.Type)
		if relationType == nil {
			continue
		}
		ca.Relations.Edges = append(ca.Relations.Edges, &media.CompleteAnime_Relations_Edges{
			Node:         AnimeFromRelatedSubject(rs),
			RelationType: relationType,
		})
	}

	return ca
}

//----------------------------------------------------------------------------------------------------------------------
// 详情组装（infobox 尽力映射）

// infoboxValue 取信息箱条目的首个字符串值。
// 兼容两种形态：单值字符串 "TV"；对象数组 [{"v":"京都动画"}]。
// 解析失败返回空串（缺字段留空，不猜测）。
func infoboxValue(items []InfoboxItem, key string) []string {
	for _, item := range items {
		if item.Key != key {
			continue
		}
		// 形态 1：纯字符串
		var s string
		if err := json.Unmarshal(item.Value, &s); err == nil {
			return []string{s}
		}
		// 形态 2：[{"v":"..."}]
		var arr []struct {
			V string `json:"v"`
		}
		if err := json.Unmarshal(item.Value, &arr); err == nil {
			out := make([]string, 0, len(arr))
			for _, a := range arr {
				if a.V != "" {
					out = append(out, a.V)
				}
			}
			return out
		}
	}
	return nil
}

// StudioNamesFromInfobox 从信息箱「制作」字段提取制作公司名列表。
// Bangumi 无制作公司 ID 体系，仅返回名称（平台层负责生成稳定的伪 ID）。
func StudioNamesFromInfobox(s *Subject) []string {
	if s == nil {
		return nil
	}
	return infoboxValue(s.Infobox, "制作")
}

// AnimeDetailsFromSubject 组装 anime 详情（anime 页专用）。
//
// 尽力映射说明（Bangumi v0 缺失的字段留空，不伪造）：
//   - Studios：infobox「制作」；Bangumi 无公司 ID，节点 ID 由调用方
//     （平台层）生成伪 ID，此处填 0 占位；
//   - Characters/Staff/Recommendations/Rankings/Trailer/Duration：
//     v0 无对应端点/字段，留空（TODO: 后续按需接入 /v0/subjects/{id}/characters）。
func AnimeDetailsFromSubject(s *Subject) *media.AnimeDetails {
	if s == nil {
		return nil
	}
	m := SubjectToMedia(s)
	base := media.AnimeFromSubject(m)

	d := &media.AnimeDetails{
		ID:          base.ID,
		SiteURL:     base.SiteURL,
		Description: base.Description,
		Genres:      base.Genres,
		MeanScore:   base.MeanScore,
		StartDate:   base.StartDate,
		EndDate:     base.EndDate,
	}

	// 平均分：直接取 Bangumi 评分（0-10 → 0-100）
	if s.Rating != nil && s.Rating.Score > 0 {
		avg := int(s.Rating.Score*10 + 0.5)
		d.AverageScore = &avg
	}

	// 制作公司（infobox「制作」）；伪 ID 由平台层回填
	names := StudioNamesFromInfobox(s)
	if len(names) > 0 {
		d.Studios = &media.AnimeDetails_Studios{}
		for _, name := range names {
			d.Studios.Nodes = append(d.Studios.Nodes, &media.AnimeDetails_Studios_Nodes{
				ID:   0, // Bangumi 无公司 ID，平台层用名称哈希回填
				Name: name,
			})
		}
	}

	return d
}

// MangaDetailsFromSubject 组装 manga 详情（manga 页专用）。
// 映射策略同 AnimeDetailsFromSubject。
// 注意：media.MangaDetails 形状较窄（无 Description/MeanScore 等字段），
// 仅承载其已有字段。
func MangaDetailsFromSubject(s *Subject) *media.MangaDetails {
	if s == nil {
		return nil
	}
	base := media.MangaFromSubject(SubjectToMedia(s))

	return &media.MangaDetails{
		ID:      base.ID,
		SiteURL: base.SiteURL,
		Genres:  base.Genres,
		// 3.9d additive：标签云 / 排名 / 信息箱
		BangumiTags: base.BangumiTags,
		BangumiRank: base.BangumiRank,
		Infobox:     base.Infobox,
	}
}
