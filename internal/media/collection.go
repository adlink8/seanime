package media

// collection.go —— AniList 收藏域（MediaListCollection 树）与用户统计的镜像层。
//
// 镜像规则与裁剪说明见 entity.go 头注释。要点：
//   - AnimeCollection/MangaCollection/AnimeCollectionWithRelations 及其
//     MediaListCollection→Lists→Entries 链逐字段镜像（JSON 形状不变）；
//   - Entries.Media 字段分别指向 *Anime / *Manga / *CompleteAnime；
//   - AnimeListEntry/AnimeList/MangaListEntry/MangaList 为别名（照搬 anilist 包）；
//   - 收藏日期（StartedAt/CompletedAt）为 FuzzyDate 别名（见 entity.go）；
//   - PatchAnimeCollectionEntry 等本地补丁原语照搬（不依赖远端客户端）；
//   - Ensure* 系列被裁剪（原实现依赖 AnilistClient 拉取缺失条目）。

import (
	"fmt"
	"time"

	"github.com/goccy/go-json"
)

//----------------------------------------------------------------------------------------------------------------------
// AnimeCollection（≈ anilist.AnimeCollection）

type AnimeCollection struct {
	MediaListCollection *AnimeCollection_MediaListCollection `json:"MediaListCollection,omitempty"`
}

type AnimeCollection_MediaListCollection struct {
	Lists []*AnimeCollection_MediaListCollection_Lists `json:"lists,omitempty"`
}

type AnimeCollection_MediaListCollection_Lists struct {
	Entries      []*AnimeCollection_MediaListCollection_Lists_Entries `json:"entries,omitempty"`
	IsCustomList *bool                                                `json:"isCustomList,omitempty"`
	Name         *string                                              `json:"name,omitempty"`
	Status       *MediaListStatus                                     `json:"status,omitempty"`
}

type AnimeCollection_MediaListCollection_Lists_Entries struct {
	CompletedAt *AnimeCollection_MediaListCollection_Lists_Entries_CompletedAt `json:"completedAt,omitempty"`
	ID          int                                                            `json:"id"`
	Media       *Anime                                                         `json:"media,omitempty"`
	Notes       *string                                                        `json:"notes,omitempty"`
	Private     *bool                                                          `json:"private,omitempty"`
	Progress    *int                                                           `json:"progress,omitempty"`
	Repeat      *int                                                           `json:"repeat,omitempty"`
	Score       *float64                                                       `json:"score,omitempty"`
	StartedAt   *AnimeCollection_MediaListCollection_Lists_Entries_StartedAt    `json:"startedAt,omitempty"`
	Status      *MediaListStatus                                               `json:"status,omitempty"`
}

func (t *AnimeCollection) GetMediaListCollection() *AnimeCollection_MediaListCollection {
	if t == nil {
		t = &AnimeCollection{}
	}
	return t.MediaListCollection
}
func (t *AnimeCollection_MediaListCollection) GetLists() []*AnimeCollection_MediaListCollection_Lists {
	if t == nil {
		t = &AnimeCollection_MediaListCollection{}
	}
	return t.Lists
}
func (t *AnimeCollection_MediaListCollection_Lists) GetEntries() []*AnimeCollection_MediaListCollection_Lists_Entries {
	if t == nil {
		t = &AnimeCollection_MediaListCollection_Lists{}
	}
	return t.Entries
}
func (t *AnimeCollection_MediaListCollection_Lists) GetIsCustomList() *bool {
	if t == nil {
		t = &AnimeCollection_MediaListCollection_Lists{}
	}
	return t.IsCustomList
}
func (t *AnimeCollection_MediaListCollection_Lists) GetName() *string {
	if t == nil {
		t = &AnimeCollection_MediaListCollection_Lists{}
	}
	return t.Name
}
func (t *AnimeCollection_MediaListCollection_Lists) GetStatus() *MediaListStatus {
	if t == nil {
		t = &AnimeCollection_MediaListCollection_Lists{}
	}
	return t.Status
}
func (t *AnimeCollection_MediaListCollection_Lists_Entries) GetCompletedAt() *AnimeCollection_MediaListCollection_Lists_Entries_CompletedAt {
	if t == nil {
		t = &AnimeCollection_MediaListCollection_Lists_Entries{}
	}
	return t.CompletedAt
}
func (t *AnimeCollection_MediaListCollection_Lists_Entries) GetID() int {
	if t == nil {
		t = &AnimeCollection_MediaListCollection_Lists_Entries{}
	}
	return t.ID
}
func (t *AnimeCollection_MediaListCollection_Lists_Entries) GetMedia() *Anime {
	if t == nil {
		t = &AnimeCollection_MediaListCollection_Lists_Entries{}
	}
	return t.Media
}
func (t *AnimeCollection_MediaListCollection_Lists_Entries) GetNotes() *string {
	if t == nil {
		t = &AnimeCollection_MediaListCollection_Lists_Entries{}
	}
	return t.Notes
}
func (t *AnimeCollection_MediaListCollection_Lists_Entries) GetPrivate() *bool {
	if t == nil {
		t = &AnimeCollection_MediaListCollection_Lists_Entries{}
	}
	return t.Private
}
func (t *AnimeCollection_MediaListCollection_Lists_Entries) GetProgress() *int {
	if t == nil {
		t = &AnimeCollection_MediaListCollection_Lists_Entries{}
	}
	return t.Progress
}
func (t *AnimeCollection_MediaListCollection_Lists_Entries) GetRepeat() *int {
	if t == nil {
		t = &AnimeCollection_MediaListCollection_Lists_Entries{}
	}
	return t.Repeat
}
func (t *AnimeCollection_MediaListCollection_Lists_Entries) GetScore() *float64 {
	if t == nil {
		t = &AnimeCollection_MediaListCollection_Lists_Entries{}
	}
	return t.Score
}
func (t *AnimeCollection_MediaListCollection_Lists_Entries) GetStartedAt() *AnimeCollection_MediaListCollection_Lists_Entries_StartedAt {
	if t == nil {
		t = &AnimeCollection_MediaListCollection_Lists_Entries{}
	}
	return t.StartedAt
}
func (t *AnimeCollection_MediaListCollection_Lists_Entries) GetStatus() *MediaListStatus {
	if t == nil {
		t = &AnimeCollection_MediaListCollection_Lists_Entries{}
	}
	return t.Status
}

//----------------------------------------------------------------------------------------------------------------------
// MangaCollection（≈ anilist.MangaCollection）

type MangaCollection struct {
	MediaListCollection *MangaCollection_MediaListCollection `json:"MediaListCollection,omitempty"`
}

type MangaCollection_MediaListCollection struct {
	Lists []*MangaCollection_MediaListCollection_Lists `json:"lists,omitempty"`
}

type MangaCollection_MediaListCollection_Lists struct {
	Entries      []*MangaCollection_MediaListCollection_Lists_Entries `json:"entries,omitempty"`
	IsCustomList *bool                                                `json:"isCustomList,omitempty"`
	Name         *string                                              `json:"name,omitempty"`
	Status       *MediaListStatus                                     `json:"status,omitempty"`
}

type MangaCollection_MediaListCollection_Lists_Entries struct {
	CompletedAt *MangaCollection_MediaListCollection_Lists_Entries_CompletedAt `json:"completedAt,omitempty"`
	ID          int                                                            `json:"id"`
	Media       *Manga                                                         `json:"media,omitempty"`
	Notes       *string                                                        `json:"notes,omitempty"`
	Private     *bool                                                          `json:"private,omitempty"`
	Progress    *int                                                           `json:"progress,omitempty"`
	Repeat      *int                                                           `json:"repeat,omitempty"`
	Score       *float64                                                       `json:"score,omitempty"`
	StartedAt   *MangaCollection_MediaListCollection_Lists_Entries_StartedAt    `json:"startedAt,omitempty"`
	Status      *MediaListStatus                                               `json:"status,omitempty"`
}

func (t *MangaCollection) GetMediaListCollection() *MangaCollection_MediaListCollection {
	if t == nil {
		t = &MangaCollection{}
	}
	return t.MediaListCollection
}
func (t *MangaCollection_MediaListCollection) GetLists() []*MangaCollection_MediaListCollection_Lists {
	if t == nil {
		t = &MangaCollection_MediaListCollection{}
	}
	return t.Lists
}
func (t *MangaCollection_MediaListCollection_Lists) GetEntries() []*MangaCollection_MediaListCollection_Lists_Entries {
	if t == nil {
		t = &MangaCollection_MediaListCollection_Lists{}
	}
	return t.Entries
}
func (t *MangaCollection_MediaListCollection_Lists) GetIsCustomList() *bool {
	if t == nil {
		t = &MangaCollection_MediaListCollection_Lists{}
	}
	return t.IsCustomList
}
func (t *MangaCollection_MediaListCollection_Lists) GetName() *string {
	if t == nil {
		t = &MangaCollection_MediaListCollection_Lists{}
	}
	return t.Name
}
func (t *MangaCollection_MediaListCollection_Lists) GetStatus() *MediaListStatus {
	if t == nil {
		t = &MangaCollection_MediaListCollection_Lists{}
	}
	return t.Status
}
func (t *MangaCollection_MediaListCollection_Lists_Entries) GetCompletedAt() *MangaCollection_MediaListCollection_Lists_Entries_CompletedAt {
	if t == nil {
		t = &MangaCollection_MediaListCollection_Lists_Entries{}
	}
	return t.CompletedAt
}
func (t *MangaCollection_MediaListCollection_Lists_Entries) GetID() int {
	if t == nil {
		t = &MangaCollection_MediaListCollection_Lists_Entries{}
	}
	return t.ID
}
func (t *MangaCollection_MediaListCollection_Lists_Entries) GetMedia() *Manga {
	if t == nil {
		t = &MangaCollection_MediaListCollection_Lists_Entries{}
	}
	return t.Media
}
func (t *MangaCollection_MediaListCollection_Lists_Entries) GetNotes() *string {
	if t == nil {
		t = &MangaCollection_MediaListCollection_Lists_Entries{}
	}
	return t.Notes
}
func (t *MangaCollection_MediaListCollection_Lists_Entries) GetPrivate() *bool {
	if t == nil {
		t = &MangaCollection_MediaListCollection_Lists_Entries{}
	}
	return t.Private
}
func (t *MangaCollection_MediaListCollection_Lists_Entries) GetProgress() *int {
	if t == nil {
		t = &MangaCollection_MediaListCollection_Lists_Entries{}
	}
	return t.Progress
}
func (t *MangaCollection_MediaListCollection_Lists_Entries) GetRepeat() *int {
	if t == nil {
		t = &MangaCollection_MediaListCollection_Lists_Entries{}
	}
	return t.Repeat
}
func (t *MangaCollection_MediaListCollection_Lists_Entries) GetScore() *float64 {
	if t == nil {
		t = &MangaCollection_MediaListCollection_Lists_Entries{}
	}
	return t.Score
}
func (t *MangaCollection_MediaListCollection_Lists_Entries) GetStartedAt() *MangaCollection_MediaListCollection_Lists_Entries_StartedAt {
	if t == nil {
		t = &MangaCollection_MediaListCollection_Lists_Entries{}
	}
	return t.StartedAt
}
func (t *MangaCollection_MediaListCollection_Lists_Entries) GetStatus() *MediaListStatus {
	if t == nil {
		t = &MangaCollection_MediaListCollection_Lists_Entries{}
	}
	return t.Status
}

//----------------------------------------------------------------------------------------------------------------------
// AnimeCollectionWithRelations（≈ anilist.AnimeCollectionWithRelations，条目带关系树）

type AnimeCollectionWithRelations struct {
	MediaListCollection *AnimeCollectionWithRelations_MediaListCollection `json:"MediaListCollection,omitempty"`
}

type AnimeCollectionWithRelations_MediaListCollection struct {
	Lists []*AnimeCollectionWithRelations_MediaListCollection_Lists `json:"lists,omitempty"`
}

type AnimeCollectionWithRelations_MediaListCollection_Lists struct {
	Entries      []*AnimeCollectionWithRelations_MediaListCollection_Lists_Entries `json:"entries,omitempty"`
	IsCustomList *bool                                                             `json:"isCustomList,omitempty"`
	Name         *string                                                           `json:"name,omitempty"`
	Status       *MediaListStatus                                                  `json:"status,omitempty"`
}

type AnimeCollectionWithRelations_MediaListCollection_Lists_Entries struct {
	CompletedAt *AnimeCollectionWithRelations_MediaListCollection_Lists_Entries_CompletedAt `json:"completedAt,omitempty"`
	ID          int                                                                         `json:"id"`
	Media       *CompleteAnime                                                              `json:"media,omitempty"`
	Notes       *string                                                                     `json:"notes,omitempty"`
	Private     *bool                                                                       `json:"private,omitempty"`
	Progress    *int                                                                        `json:"progress,omitempty"`
	Repeat      *int                                                                        `json:"repeat,omitempty"`
	Score       *float64                                                                    `json:"score,omitempty"`
	StartedAt   *AnimeCollectionWithRelations_MediaListCollection_Lists_Entries_StartedAt   `json:"startedAt,omitempty"`
	Status      *MediaListStatus                                                            `json:"status,omitempty"`
}

type (
	AnimeCollectionWithRelations_MediaListCollection_Lists_Entries_StartedAt   = FuzzyDate
	AnimeCollectionWithRelations_MediaListCollection_Lists_Entries_CompletedAt = FuzzyDate
)

func (t *AnimeCollectionWithRelations) GetMediaListCollection() *AnimeCollectionWithRelations_MediaListCollection {
	if t == nil {
		t = &AnimeCollectionWithRelations{}
	}
	return t.MediaListCollection
}
func (t *AnimeCollectionWithRelations_MediaListCollection) GetLists() []*AnimeCollectionWithRelations_MediaListCollection_Lists {
	if t == nil {
		t = &AnimeCollectionWithRelations_MediaListCollection{}
	}
	return t.Lists
}
func (t *AnimeCollectionWithRelations_MediaListCollection_Lists) GetEntries() []*AnimeCollectionWithRelations_MediaListCollection_Lists_Entries {
	if t == nil {
		t = &AnimeCollectionWithRelations_MediaListCollection_Lists{}
	}
	return t.Entries
}
func (t *AnimeCollectionWithRelations_MediaListCollection_Lists) GetIsCustomList() *bool {
	if t == nil {
		t = &AnimeCollectionWithRelations_MediaListCollection_Lists{}
	}
	return t.IsCustomList
}
func (t *AnimeCollectionWithRelations_MediaListCollection_Lists) GetName() *string {
	if t == nil {
		t = &AnimeCollectionWithRelations_MediaListCollection_Lists{}
	}
	return t.Name
}
func (t *AnimeCollectionWithRelations_MediaListCollection_Lists) GetStatus() *MediaListStatus {
	if t == nil {
		t = &AnimeCollectionWithRelations_MediaListCollection_Lists{}
	}
	return t.Status
}
func (t *AnimeCollectionWithRelations_MediaListCollection_Lists_Entries) GetCompletedAt() *AnimeCollectionWithRelations_MediaListCollection_Lists_Entries_CompletedAt {
	if t == nil {
		t = &AnimeCollectionWithRelations_MediaListCollection_Lists_Entries{}
	}
	return t.CompletedAt
}
func (t *AnimeCollectionWithRelations_MediaListCollection_Lists_Entries) GetID() int {
	if t == nil {
		t = &AnimeCollectionWithRelations_MediaListCollection_Lists_Entries{}
	}
	return t.ID
}
func (t *AnimeCollectionWithRelations_MediaListCollection_Lists_Entries) GetMedia() *CompleteAnime {
	if t == nil {
		t = &AnimeCollectionWithRelations_MediaListCollection_Lists_Entries{}
	}
	return t.Media
}
func (t *AnimeCollectionWithRelations_MediaListCollection_Lists_Entries) GetNotes() *string {
	if t == nil {
		t = &AnimeCollectionWithRelations_MediaListCollection_Lists_Entries{}
	}
	return t.Notes
}
func (t *AnimeCollectionWithRelations_MediaListCollection_Lists_Entries) GetPrivate() *bool {
	if t == nil {
		t = &AnimeCollectionWithRelations_MediaListCollection_Lists_Entries{}
	}
	return t.Private
}
func (t *AnimeCollectionWithRelations_MediaListCollection_Lists_Entries) GetProgress() *int {
	if t == nil {
		t = &AnimeCollectionWithRelations_MediaListCollection_Lists_Entries{}
	}
	return t.Progress
}
func (t *AnimeCollectionWithRelations_MediaListCollection_Lists_Entries) GetRepeat() *int {
	if t == nil {
		t = &AnimeCollectionWithRelations_MediaListCollection_Lists_Entries{}
	}
	return t.Repeat
}
func (t *AnimeCollectionWithRelations_MediaListCollection_Lists_Entries) GetScore() *float64 {
	if t == nil {
		t = &AnimeCollectionWithRelations_MediaListCollection_Lists_Entries{}
	}
	return t.Score
}
func (t *AnimeCollectionWithRelations_MediaListCollection_Lists_Entries) GetStartedAt() *AnimeCollectionWithRelations_MediaListCollection_Lists_Entries_StartedAt {
	if t == nil {
		t = &AnimeCollectionWithRelations_MediaListCollection_Lists_Entries{}
	}
	return t.StartedAt
}
func (t *AnimeCollectionWithRelations_MediaListCollection_Lists_Entries) GetStatus() *MediaListStatus {
	if t == nil {
		t = &AnimeCollectionWithRelations_MediaListCollection_Lists_Entries{}
	}
	return t.Status
}

//----------------------------------------------------------------------------------------------------------------------
// 收藏条目别名与查询辅助（照搬 anilist/collection_helper.go、entries.go、manga.go）

type (
	// AnimeListEntry 动画收藏条目；AnimeList 动画收藏列表
	AnimeListEntry = AnimeCollection_MediaListCollection_Lists_Entries
	AnimeList      = AnimeCollection_MediaListCollection_Lists

	// MangaListEntry 漫画收藏条目；MangaList 漫画收藏列表
	MangaListEntry = MangaCollection_MediaListCollection_Lists_Entries
	MangaList      = MangaCollection_MediaListCollection_Lists

	EntryDate struct {
		Year  *int `json:"year,omitempty"`
		Month *int `json:"month,omitempty"`
		Day   *int `json:"day,omitempty"`
	}
)

func (m *AnimeListEntry) GetProgressSafe() int {
	if m == nil {
		return 0
	}
	if m.Progress == nil {
		return 0
	}
	return *m.Progress
}

func (m *AnimeListEntry) GetScoreSafe() float64 {
	if m == nil {
		return 0
	}
	if m.Score == nil {
		return 0
	}
	return *m.Score
}

func (m *AnimeListEntry) GetRepeatSafe() int {
	if m == nil {
		return 0
	}
	if m.Repeat == nil {
		return 0
	}
	return *m.Repeat
}

func (m *AnimeListEntry) GetStatusSafe() MediaListStatus {
	if m == nil {
		return ""
	}
	if m.Status == nil {
		return ""
	}
	return *m.Status
}

func (m *MangaListEntry) GetRepeatSafe() int {
	if m == nil {
		return 0
	}
	if m.Repeat == nil {
		return 0
	}
	return *m.Repeat
}

// GetListEntryFromAnimeId 按 media ID 查找动画收藏条目。
func (ac *AnimeCollection) GetListEntryFromAnimeId(id int) (*AnimeListEntry, bool) {
	if ac == nil || ac.MediaListCollection == nil {
		return nil, false
	}

	var entry *AnimeCollection_MediaListCollection_Lists_Entries
	for _, l := range ac.MediaListCollection.Lists {
		if l.Entries == nil || len(l.Entries) == 0 {
			continue
		}
		for _, e := range l.Entries {
			if e.Media.ID == id {
				entry = e
				break
			}
		}
	}
	if entry == nil {
		return nil, false
	}

	return entry, true
}

// GetAllAnime 返回收藏中全部动画（按 ID 去重）。
func (ac *AnimeCollection) GetAllAnime() []*Anime {
	if ac == nil {
		return make([]*Anime, 0)
	}

	var ret []*Anime
	addedId := make(map[int]bool)
	for _, l := range ac.MediaListCollection.Lists {
		if l.Entries == nil || len(l.Entries) == 0 {
			continue
		}
		for _, e := range l.Entries {
			if _, ok := addedId[e.Media.ID]; !ok {
				ret = append(ret, e.Media)
				addedId[e.Media.ID] = true
			}
		}
	}
	return ret
}

// FindAnime 按 media ID 查找收藏中的动画实体。
func (ac *AnimeCollection) FindAnime(mediaId int) (*Anime, bool) {
	if ac == nil {
		return nil, false
	}

	for _, l := range ac.MediaListCollection.Lists {
		if l.Entries == nil || len(l.Entries) == 0 {
			continue
		}
		for _, e := range l.Entries {
			if e.Media.ID == mediaId {
				return e.Media, true
			}
		}
	}
	return nil, false
}

// GetListEntryFromMediaId 同 GetListEntryFromAnimeId（带关系版本）。
func (ac *AnimeCollectionWithRelations) GetListEntryFromMediaId(id int) (*AnimeCollectionWithRelations_MediaListCollection_Lists_Entries, bool) {
	if ac == nil || ac.MediaListCollection == nil {
		return nil, false
	}

	var entry *AnimeCollectionWithRelations_MediaListCollection_Lists_Entries
	for _, l := range ac.MediaListCollection.Lists {
		if l.Entries == nil || len(l.Entries) == 0 {
			continue
		}
		for _, e := range l.Entries {
			if e.Media.ID == id {
				entry = e
				break
			}
		}
	}
	if entry == nil {
		return nil, false
	}

	return entry, true
}

// GetAllAnime 返回带关系收藏中的全部条目（按 ID 去重）。
func (ac *AnimeCollectionWithRelations) GetAllAnime() []*CompleteAnime {
	var ret []*CompleteAnime
	addedId := make(map[int]bool)
	for _, l := range ac.MediaListCollection.Lists {
		if l.Entries == nil || len(l.Entries) == 0 {
			continue
		}
		for _, e := range l.Entries {
			if _, ok := addedId[e.Media.ID]; !ok {
				ret = append(ret, e.Media)
				addedId[e.Media.ID] = true
			}
		}
	}
	return ret
}

// FindAnime 按 media ID 查找带关系收藏中的条目。
func (ac *AnimeCollectionWithRelations) FindAnime(mediaId int) (*CompleteAnime, bool) {
	for _, l := range ac.MediaListCollection.Lists {
		if l.Entries == nil || len(l.Entries) == 0 {
			continue
		}
		for _, e := range l.Entries {
			if e.Media.ID == mediaId {
				return e.Media, true
			}
		}
	}
	return nil, false
}

// GetListEntryFromMangaId 按 media ID 查找漫画收藏条目。
func (ac *MangaCollection) GetListEntryFromMangaId(id int) (*MangaListEntry, bool) {
	if ac == nil || ac.MediaListCollection == nil {
		return nil, false
	}

	var entry *MangaCollection_MediaListCollection_Lists_Entries
	for _, l := range ac.MediaListCollection.Lists {
		if l.Entries == nil || len(l.Entries) == 0 {
			continue
		}
		for _, e := range l.Entries {
			if e.Media.ID == id {
				entry = e
				break
			}
		}
	}
	if entry == nil {
		return nil, false
	}

	return entry, true
}

// IFuzzyDate 日期抽象（照搬）
type IFuzzyDate interface {
	GetYear() *int
	GetMonth() *int
	GetDay() *int
}

// FuzzyDateToString 模糊日期 → RFC3339 字符串（无年份时返回空串）。
func FuzzyDateToString(d IFuzzyDate) string {
	if d == nil {
		return ""
	}
	return fuzzyDateToString(d.GetYear(), d.GetMonth(), d.GetDay())
}

func ToEntryStartDate(d *AnimeCollection_MediaListCollection_Lists_Entries_StartedAt) string {
	if d == nil {
		return ""
	}
	return fuzzyDateToString(d.GetYear(), d.GetMonth(), d.GetDay())
}

func ToEntryCompletionDate(d *AnimeCollection_MediaListCollection_Lists_Entries_CompletedAt) string {
	if d == nil {
		return ""
	}
	return fuzzyDateToString(d.GetYear(), d.GetMonth(), d.GetDay())
}

func fuzzyDateToString(year *int, month *int, day *int) string {
	_year := 0
	if year != nil {
		_year = *year
	}
	if _year == 0 {
		return ""
	}
	_month := 0
	if month != nil {
		_month = *month
	}
	_day := 0
	if day != nil {
		_day = *day
	}
	return time.Date(_year, time.Month(_month), _day, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
}

// AddEntryToList 将条目按状态加入对应列表，无对应列表时新建（照搬语义）。
func (mc *AnimeCollection_MediaListCollection) AddEntryToList(entry *AnimeCollection_MediaListCollection_Lists_Entries, status MediaListStatus) {
	if mc == nil || entry == nil {
		return
	}

	if mc.Lists == nil {
		mc.Lists = make([]*AnimeCollection_MediaListCollection_Lists, 0)
	}

	for _, list := range mc.Lists {
		if list.Status != nil && *list.Status == status {
			if list.Entries == nil {
				list.Entries = make([]*AnimeCollection_MediaListCollection_Lists_Entries, 0)
			}
			list.Entries = append(list.Entries, entry)
			return
		}
	}

	newList := &AnimeCollection_MediaListCollection_Lists{
		Status:  &status,
		Entries: []*AnimeCollection_MediaListCollection_Lists_Entries{entry},
	}
	mc.Lists = append(mc.Lists, newList)
}

// Copy 深拷贝（JSON 往返）。
func (ac *AnimeCollection) Copy() *AnimeCollection {
	if ac == nil {
		return nil
	}
	marshaled, err := json.Marshal(ac)
	if err != nil {
		return nil
	}
	var copy AnimeCollection
	err = json.Unmarshal(marshaled, &copy)
	if err != nil {
		return nil
	}
	return &copy
}

func (ac *AnimeList) CopyT() *AnimeCollection_MediaListCollection_Lists {
	if ac == nil {
		return nil
	}
	marshaled, err := json.Marshal(ac)
	if err != nil {
		return nil
	}
	var copy AnimeCollection_MediaListCollection_Lists
	err = json.Unmarshal(marshaled, &copy)
	if err != nil {
		return nil
	}
	return &copy
}

//----------------------------------------------------------------------------------------------------------------------
// 收藏条目补丁原语（照搬 anilist/test_collection_helpers.go 中不依赖客户端的部分）

// AnimeCollectionEntryPatch 收藏条目补丁。
type AnimeCollectionEntryPatch struct {
	Status            *MediaListStatus
	Progress          *int
	Score             *float64
	Repeat            *int
	AiredEpisodes     *int
	NextAiringEpisode *Anime_NextAiringEpisode
}

// PatchAnimeCollectionEntry 原地修改收藏条目（找不到条目时 panic，照搬语义）。
func PatchAnimeCollectionEntry(collection *AnimeCollection, mediaID int, patch AnimeCollectionEntryPatch) *AnimeCollection {
	if collection == nil {
		panic("media: anime collection is nil")
	}

	entry, currentList := findAnimeCollectionEntry(collection, mediaID)
	if entry == nil {
		panic(fmt.Sprintf("media: anime %d not found in collection", mediaID))
	}

	if patch.Status != nil {
		currentList = moveAnimeEntryToStatus(collection, currentList, entry, *patch.Status)
		_ = currentList
	}

	applyAnimeCollectionEntryPatch(entry, patch)
	return collection
}

// PatchAnimeCollectionWithRelationsEntry 带关系版本的补丁。
func PatchAnimeCollectionWithRelationsEntry(collection *AnimeCollectionWithRelations, mediaID int, patch AnimeCollectionEntryPatch) *AnimeCollectionWithRelations {
	if collection == nil {
		panic("media: anime collection with relations is nil")
	}

	entry, currentList := findAnimeCollectionWithRelationsEntry(collection, mediaID)
	if entry == nil {
		panic(fmt.Sprintf("media: anime %d not found in relation collection", mediaID))
	}

	if patch.Status != nil {
		currentList = moveAnimeRelationsEntryToStatus(collection, currentList, entry, *patch.Status)
		_ = currentList
	}

	applyAnimeCollectionWithRelationsEntryPatch(entry, patch)
	return collection
}

// EnsureAnimeCollectionWithRelationsEntry 被裁剪：原实现依赖 AnilistClient
// 拉取缺失条目；media 包无客户端概念，由 Wave B 在平台层基于 Patch* 原语重建。

func applyAnimeCollectionEntryPatch(entry *AnimeCollection_MediaListCollection_Lists_Entries, patch AnimeCollectionEntryPatch) {
	if patch.Status != nil {
		entry.Status = patch.Status
	}
	if patch.Progress != nil {
		entry.Progress = patch.Progress
	}
	if patch.Score != nil {
		entry.Score = patch.Score
	}
	if patch.Repeat != nil {
		entry.Repeat = patch.Repeat
	}
	if patch.AiredEpisodes != nil {
		entry.Media.Episodes = patch.AiredEpisodes
	}
	if patch.NextAiringEpisode != nil {
		entry.Media.NextAiringEpisode = patch.NextAiringEpisode
	}
}

func applyAnimeCollectionWithRelationsEntryPatch(entry *AnimeCollectionWithRelations_MediaListCollection_Lists_Entries, patch AnimeCollectionEntryPatch) {
	if patch.Status != nil {
		entry.Status = patch.Status
	}
	if patch.Progress != nil {
		entry.Progress = patch.Progress
	}
	if patch.Score != nil {
		entry.Score = patch.Score
	}
	if patch.Repeat != nil {
		entry.Repeat = patch.Repeat
	}
	if patch.AiredEpisodes != nil {
		entry.Media.Episodes = patch.AiredEpisodes
	}
}

func findAnimeCollectionEntry(collection *AnimeCollection, mediaID int) (*AnimeCollection_MediaListCollection_Lists_Entries, *AnimeCollection_MediaListCollection_Lists) {
	if collection == nil || collection.MediaListCollection == nil {
		return nil, nil
	}

	for _, list := range collection.MediaListCollection.Lists {
		if list == nil || list.Entries == nil {
			continue
		}
		for _, entry := range list.Entries {
			if entry != nil && entry.GetMedia().GetID() == mediaID {
				return entry, list
			}
		}
	}

	return nil, nil
}

func findAnimeCollectionWithRelationsEntry(collection *AnimeCollectionWithRelations, mediaID int) (*AnimeCollectionWithRelations_MediaListCollection_Lists_Entries, *AnimeCollectionWithRelations_MediaListCollection_Lists) {
	if collection == nil || collection.MediaListCollection == nil {
		return nil, nil
	}

	for _, list := range collection.MediaListCollection.Lists {
		if list == nil || list.Entries == nil {
			continue
		}
		for _, entry := range list.Entries {
			if entry != nil && entry.GetMedia().GetID() == mediaID {
				return entry, list
			}
		}
	}

	return nil, nil
}

func moveAnimeEntryToStatus(collection *AnimeCollection, currentList *AnimeCollection_MediaListCollection_Lists, entry *AnimeCollection_MediaListCollection_Lists_Entries, status MediaListStatus) *AnimeCollection_MediaListCollection_Lists {
	if currentList != nil && currentList.Status != nil && *currentList.Status == status {
		return currentList
	}
	if currentList != nil {
		removeAnimeEntry(currentList, entry.GetMedia().GetID())
	}

	target := ensureAnimeStatusList(collection, status)
	target.Entries = append(target.Entries, entry)
	return target
}

func moveAnimeRelationsEntryToStatus(collection *AnimeCollectionWithRelations, currentList *AnimeCollectionWithRelations_MediaListCollection_Lists, entry *AnimeCollectionWithRelations_MediaListCollection_Lists_Entries, status MediaListStatus) *AnimeCollectionWithRelations_MediaListCollection_Lists {
	if currentList != nil && currentList.Status != nil && *currentList.Status == status {
		return currentList
	}
	if currentList != nil {
		removeAnimeRelationsEntry(currentList, entry.GetMedia().GetID())
	}

	target := ensureAnimeRelationsStatusList(collection, status)
	target.Entries = append(target.Entries, entry)
	return target
}

func ensureAnimeStatusList(collection *AnimeCollection, status MediaListStatus) *AnimeCollection_MediaListCollection_Lists {
	if collection.MediaListCollection == nil {
		collection.MediaListCollection = &AnimeCollection_MediaListCollection{}
	}

	for _, list := range collection.MediaListCollection.Lists {
		if list != nil && list.Status != nil && *list.Status == status {
			if list.Entries == nil {
				list.Entries = []*AnimeCollection_MediaListCollection_Lists_Entries{}
			}
			return list
		}
	}

	name := string(status)
	isCustomList := false
	list := &AnimeCollection_MediaListCollection_Lists{
		Status:       ptrOf(status),
		Name:         &name,
		IsCustomList: &isCustomList,
		Entries:      []*AnimeCollection_MediaListCollection_Lists_Entries{},
	}
	collection.MediaListCollection.Lists = append(collection.MediaListCollection.Lists, list)
	return list
}

func ensureAnimeRelationsStatusList(collection *AnimeCollectionWithRelations, status MediaListStatus) *AnimeCollectionWithRelations_MediaListCollection_Lists {
	if collection.MediaListCollection == nil {
		collection.MediaListCollection = &AnimeCollectionWithRelations_MediaListCollection{}
	}

	for _, list := range collection.MediaListCollection.Lists {
		if list != nil && list.Status != nil && *list.Status == status {
			if list.Entries == nil {
				list.Entries = []*AnimeCollectionWithRelations_MediaListCollection_Lists_Entries{}
			}
			return list
		}
	}

	name := string(status)
	isCustomList := false
	list := &AnimeCollectionWithRelations_MediaListCollection_Lists{
		Status:       ptrOf(status),
		Name:         &name,
		IsCustomList: &isCustomList,
		Entries:      []*AnimeCollectionWithRelations_MediaListCollection_Lists_Entries{},
	}
	collection.MediaListCollection.Lists = append(collection.MediaListCollection.Lists, list)
	return list
}

func removeAnimeEntry(list *AnimeCollection_MediaListCollection_Lists, mediaID int) {
	for idx, entry := range list.GetEntries() {
		if entry != nil && entry.GetMedia().GetID() == mediaID {
			list.Entries = append(list.Entries[:idx], list.Entries[idx+1:]...)
			return
		}
	}
}

func removeAnimeRelationsEntry(list *AnimeCollectionWithRelations_MediaListCollection_Lists, mediaID int) {
	for idx, entry := range list.GetEntries() {
		if entry != nil && entry.GetMedia().GetID() == mediaID {
			list.Entries = append(list.Entries[:idx], list.Entries[idx+1:]...)
			return
		}
	}
}

func ptrOf[T any](value T) *T {
	return &value
}

//----------------------------------------------------------------------------------------------------------------------
// 收藏标签查询（≈ anilist.AnimeCollectionTags / MangaCollectionTags）

type AnimeCollectionTags struct {
	MediaListCollection *AnimeCollectionTags_MediaListCollection `json:"MediaListCollection,omitempty"`
}

type AnimeCollectionTags_MediaListCollection struct {
	Lists []*AnimeCollectionTags_MediaListCollection_Lists `json:"lists,omitempty"`
}

type AnimeCollectionTags_MediaListCollection_Lists struct {
	Entries []*AnimeCollectionTags_MediaListCollection_Lists_Entries `json:"entries,omitempty"`
}

type AnimeCollectionTags_MediaListCollection_Lists_Entries struct {
	ID    int                                                  `json:"id"`
	Media *AnimeCollectionTags_MediaListCollection_Lists_Entries_Media `json:"media,omitempty"`
}

type AnimeCollectionTags_MediaListCollection_Lists_Entries_Media struct {
	ID   int                                                          `json:"id"`
	Tags []*AnimeCollectionTags_MediaListCollection_Lists_Entries_Media_Tags `json:"tags,omitempty"`
}

type AnimeCollectionTags_MediaListCollection_Lists_Entries_Media_Tags struct {
	Name string `json:"name"`
}

type MangaCollectionTags struct {
	MediaListCollection *MangaCollectionTags_MediaListCollection `json:"MediaListCollection,omitempty"`
}

type MangaCollectionTags_MediaListCollection struct {
	Lists []*MangaCollectionTags_MediaListCollection_Lists `json:"lists,omitempty"`
}

type MangaCollectionTags_MediaListCollection_Lists struct {
	Entries []*MangaCollectionTags_MediaListCollection_Lists_Entries `json:"entries,omitempty"`
}

type MangaCollectionTags_MediaListCollection_Lists_Entries struct {
	ID    int                                                  `json:"id"`
	Media *MangaCollectionTags_MediaListCollection_Lists_Entries_Media `json:"media,omitempty"`
}

type MangaCollectionTags_MediaListCollection_Lists_Entries_Media struct {
	ID   int                                                          `json:"id"`
	Tags []*MangaCollectionTags_MediaListCollection_Lists_Entries_Media_Tags `json:"tags,omitempty"`
}

type MangaCollectionTags_MediaListCollection_Lists_Entries_Media_Tags struct {
	Name string `json:"name"`
}

func (t *AnimeCollectionTags) GetMediaListCollection() *AnimeCollectionTags_MediaListCollection {
	if t == nil {
		t = &AnimeCollectionTags{}
	}
	return t.MediaListCollection
}
func (t *AnimeCollectionTags_MediaListCollection) GetLists() []*AnimeCollectionTags_MediaListCollection_Lists {
	if t == nil {
		t = &AnimeCollectionTags_MediaListCollection{}
	}
	return t.Lists
}
func (t *AnimeCollectionTags_MediaListCollection_Lists) GetEntries() []*AnimeCollectionTags_MediaListCollection_Lists_Entries {
	if t == nil {
		t = &AnimeCollectionTags_MediaListCollection_Lists{}
	}
	return t.Entries
}
func (t *AnimeCollectionTags_MediaListCollection_Lists_Entries) GetMedia() *AnimeCollectionTags_MediaListCollection_Lists_Entries_Media {
	if t == nil {
		t = &AnimeCollectionTags_MediaListCollection_Lists_Entries{}
	}
	return t.Media
}
func (t *AnimeCollectionTags_MediaListCollection_Lists_Entries_Media) GetID() int {
	if t == nil {
		t = &AnimeCollectionTags_MediaListCollection_Lists_Entries_Media{}
	}
	return t.ID
}
func (t *AnimeCollectionTags_MediaListCollection_Lists_Entries_Media) GetTags() []*AnimeCollectionTags_MediaListCollection_Lists_Entries_Media_Tags {
	if t == nil {
		t = &AnimeCollectionTags_MediaListCollection_Lists_Entries_Media{}
	}
	return t.Tags
}
func (t *AnimeCollectionTags_MediaListCollection_Lists_Entries_Media_Tags) GetName() string {
	if t == nil {
		t = &AnimeCollectionTags_MediaListCollection_Lists_Entries_Media_Tags{}
	}
	return t.Name
}

func (t *MangaCollectionTags) GetMediaListCollection() *MangaCollectionTags_MediaListCollection {
	if t == nil {
		t = &MangaCollectionTags{}
	}
	return t.MediaListCollection
}
func (t *MangaCollectionTags_MediaListCollection) GetLists() []*MangaCollectionTags_MediaListCollection_Lists {
	if t == nil {
		t = &MangaCollectionTags_MediaListCollection{}
	}
	return t.Lists
}
func (t *MangaCollectionTags_MediaListCollection_Lists) GetEntries() []*MangaCollectionTags_MediaListCollection_Lists_Entries {
	if t == nil {
		t = &MangaCollectionTags_MediaListCollection_Lists{}
	}
	return t.Entries
}
func (t *MangaCollectionTags_MediaListCollection_Lists_Entries) GetMedia() *MangaCollectionTags_MediaListCollection_Lists_Entries_Media {
	if t == nil {
		t = &MangaCollectionTags_MediaListCollection_Lists_Entries{}
	}
	return t.Media
}
func (t *MangaCollectionTags_MediaListCollection_Lists_Entries_Media) GetID() int {
	if t == nil {
		t = &MangaCollectionTags_MediaListCollection_Lists_Entries_Media{}
	}
	return t.ID
}
func (t *MangaCollectionTags_MediaListCollection_Lists_Entries_Media) GetTags() []*MangaCollectionTags_MediaListCollection_Lists_Entries_Media_Tags {
	if t == nil {
		t = &MangaCollectionTags_MediaListCollection_Lists_Entries_Media{}
	}
	return t.Tags
}
func (t *MangaCollectionTags_MediaListCollection_Lists_Entries_Media_Tags) GetName() string {
	if t == nil {
		t = &MangaCollectionTags_MediaListCollection_Lists_Entries_Media_Tags{}
	}
	return t.Name
}

//----------------------------------------------------------------------------------------------------------------------
// MediaTagMap（照搬 anilist/tags.go）

// MediaTagMap media ID → 标签名列表
type MediaTagMap map[int][]string

func MediaTagMapFromAnimeCollectionTags(data *AnimeCollectionTags) MediaTagMap {
	ret := make(MediaTagMap)
	if data == nil || data.GetMediaListCollection() == nil {
		return ret
	}

	for _, list := range data.GetMediaListCollection().GetLists() {
		if list == nil {
			continue
		}
		for _, entry := range list.GetEntries() {
			if entry == nil || entry.GetMedia() == nil {
				continue
			}
			for _, tag := range entry.GetMedia().GetTags() {
				if tag == nil {
					continue
				}
				ret.add(entry.GetMedia().GetID(), tag.GetName())
			}
		}
	}

	return ret
}

func MediaTagMapFromMangaCollectionTags(data *MangaCollectionTags) MediaTagMap {
	ret := make(MediaTagMap)
	if data == nil || data.GetMediaListCollection() == nil {
		return ret
	}

	for _, list := range data.GetMediaListCollection().GetLists() {
		if list == nil {
			continue
		}
		for _, entry := range list.GetEntries() {
			if entry == nil || entry.GetMedia() == nil {
				continue
			}
			for _, tag := range entry.GetMedia().GetTags() {
				if tag == nil {
					continue
				}
				ret.add(entry.GetMedia().GetID(), tag.GetName())
			}
		}
	}

	return ret
}

func (m MediaTagMap) add(mediaID int, tagName string) {
	if tagName == "" {
		return
	}

	existing := m[mediaID]
	for _, current := range existing {
		if current == tagName {
			return
		}
	}

	m[mediaID] = append(existing, tagName)
}

//----------------------------------------------------------------------------------------------------------------------
// 用户统计（≈ anilist.ViewerStats / Stats）

type ViewerStats struct {
	Viewer *ViewerStats_Viewer `json:"Viewer,omitempty"`
}

type ViewerStats_Viewer struct {
	Statistics *ViewerStats_Viewer_Statistics `json:"statistics,omitempty"`
}

type ViewerStats_Viewer_Statistics struct {
	Anime *ViewerStats_Viewer_Statistics_Anime `json:"anime,omitempty"`
	Manga *ViewerStats_Viewer_Statistics_Manga `json:"manga,omitempty"`
}

type ViewerStats_Viewer_Statistics_Anime struct {
	Count           int                     `json:"count"`
	EpisodesWatched int                     `json:"episodesWatched"`
	Formats         []*UserFormatStats      `json:"formats,omitempty"`
	Genres          []*UserGenreStats       `json:"genres,omitempty"`
	MeanScore       float64                 `json:"meanScore"`
	MinutesWatched  int                     `json:"minutesWatched"`
	ReleaseYears    []*UserReleaseYearStats `json:"releaseYears,omitempty"`
	Scores          []*UserScoreStats       `json:"scores,omitempty"`
	StartYears      []*UserStartYearStats   `json:"startYears,omitempty"`
	Statuses        []*UserStatusStats      `json:"statuses,omitempty"`
	Studios         []*UserStudioStats      `json:"studios,omitempty"`
}

type ViewerStats_Viewer_Statistics_Manga struct {
	ChaptersRead int                     `json:"chaptersRead"`
	Count        int                     `json:"count"`
	Formats      []*UserFormatStats      `json:"formats,omitempty"`
	Genres       []*UserGenreStats       `json:"genres,omitempty"`
	MeanScore    float64                 `json:"meanScore"`
	ReleaseYears []*UserReleaseYearStats `json:"releaseYears,omitempty"`
	Scores       []*UserScoreStats       `json:"scores,omitempty"`
	StartYears   []*UserStartYearStats   `json:"startYears,omitempty"`
	Statuses     []*UserStatusStats      `json:"statuses,omitempty"`
	Studios      []*UserStudioStats      `json:"studios,omitempty"`
}

type UserFormatStats struct {
	Format         *MediaFormat `json:"format,omitempty"`
	MeanScore      float64      `json:"meanScore"`
	Count          int          `json:"count"`
	MinutesWatched int          `json:"minutesWatched"`
	MediaIds       []*int       `json:"mediaIds"`
	ChaptersRead   int          `json:"chaptersRead"`
}

type UserGenreStats struct {
	Genre          *string `json:"genre,omitempty"`
	MeanScore      float64 `json:"meanScore"`
	Count          int     `json:"count"`
	MinutesWatched int     `json:"minutesWatched"`
	MediaIds       []*int  `json:"mediaIds"`
	ChaptersRead   int     `json:"chaptersRead"`
}

type UserStatusStats struct {
	Status         *MediaListStatus `json:"status,omitempty"`
	MeanScore      float64          `json:"meanScore"`
	Count          int              `json:"count"`
	MinutesWatched int              `json:"minutesWatched"`
	MediaIds       []*int           `json:"mediaIds"`
	ChaptersRead   int              `json:"chaptersRead"`
}

type UserScoreStats struct {
	Score          *int    `json:"score,omitempty"`
	MeanScore      float64 `json:"meanScore"`
	Count          int     `json:"count"`
	MinutesWatched int     `json:"minutesWatched"`
	MediaIds       []*int  `json:"mediaIds"`
	ChaptersRead   int     `json:"chaptersRead"`
}

type UserStudioStats struct {
	Studio         *UserStudioStats_Studio `json:"studio,omitempty"`
	MeanScore      float64                 `json:"meanScore"`
	Count          int                     `json:"count"`
	MinutesWatched int                     `json:"minutesWatched"`
	MediaIds       []*int                  `json:"mediaIds"`
	ChaptersRead   int                     `json:"chaptersRead"`
}

type UserStudioStats_Studio struct {
	ID                int    `json:"id"`
	IsAnimationStudio bool   `json:"isAnimationStudio"`
	Name              string `json:"name"`
}

type UserStartYearStats struct {
	StartYear      *int    `json:"startYear,omitempty"`
	MeanScore      float64 `json:"meanScore"`
	Count          int     `json:"count"`
	MinutesWatched int     `json:"minutesWatched"`
	MediaIds       []*int  `json:"mediaIds"`
	ChaptersRead   int     `json:"chaptersRead"`
}

type UserReleaseYearStats struct {
	ReleaseYear    *int    `json:"releaseYear,omitempty"`
	MeanScore      float64 `json:"meanScore"`
	Count          int     `json:"count"`
	MinutesWatched int     `json:"minutesWatched"`
	MediaIds       []*int  `json:"mediaIds"`
	ChaptersRead   int     `json:"chaptersRead"`
}

// ViewerStats_Viewer_Statistics_Anime_Studios_UserStudioStats_Studio 为
// anilist 包中的重复定义，形状与 UserStudioStats_Studio 完全一致，别名承接。
type ViewerStats_Viewer_Statistics_Anime_Studios_UserStudioStats_Studio  = UserStudioStats_Studio
type ViewerStats_Viewer_Statistics_Manga_Studios_UserStudioStats_Studio  = UserStudioStats_Studio

func (t *ViewerStats) GetViewer() *ViewerStats_Viewer {
	if t == nil {
		t = &ViewerStats{}
	}
	return t.Viewer
}
func (t *ViewerStats_Viewer) GetStatistics() *ViewerStats_Viewer_Statistics {
	if t == nil {
		t = &ViewerStats_Viewer{}
	}
	return t.Statistics
}
func (t *ViewerStats_Viewer_Statistics) GetAnime() *ViewerStats_Viewer_Statistics_Anime {
	if t == nil {
		t = &ViewerStats_Viewer_Statistics{}
	}
	return t.Anime
}
func (t *ViewerStats_Viewer_Statistics) GetManga() *ViewerStats_Viewer_Statistics_Manga {
	if t == nil {
		t = &ViewerStats_Viewer_Statistics{}
	}
	return t.Manga
}
func (t *ViewerStats_Viewer_Statistics_Anime) GetCount() int {
	if t == nil {
		t = &ViewerStats_Viewer_Statistics_Anime{}
	}
	return t.Count
}
func (t *ViewerStats_Viewer_Statistics_Anime) GetEpisodesWatched() int {
	if t == nil {
		t = &ViewerStats_Viewer_Statistics_Anime{}
	}
	return t.EpisodesWatched
}
func (t *ViewerStats_Viewer_Statistics_Anime) GetFormats() []*UserFormatStats {
	if t == nil {
		t = &ViewerStats_Viewer_Statistics_Anime{}
	}
	return t.Formats
}
func (t *ViewerStats_Viewer_Statistics_Anime) GetGenres() []*UserGenreStats {
	if t == nil {
		t = &ViewerStats_Viewer_Statistics_Anime{}
	}
	return t.Genres
}
func (t *ViewerStats_Viewer_Statistics_Anime) GetMeanScore() float64 {
	if t == nil {
		t = &ViewerStats_Viewer_Statistics_Anime{}
	}
	return t.MeanScore
}
func (t *ViewerStats_Viewer_Statistics_Anime) GetMinutesWatched() int {
	if t == nil {
		t = &ViewerStats_Viewer_Statistics_Anime{}
	}
	return t.MinutesWatched
}
func (t *ViewerStats_Viewer_Statistics_Anime) GetReleaseYears() []*UserReleaseYearStats {
	if t == nil {
		t = &ViewerStats_Viewer_Statistics_Anime{}
	}
	return t.ReleaseYears
}
func (t *ViewerStats_Viewer_Statistics_Anime) GetScores() []*UserScoreStats {
	if t == nil {
		t = &ViewerStats_Viewer_Statistics_Anime{}
	}
	return t.Scores
}
func (t *ViewerStats_Viewer_Statistics_Anime) GetStartYears() []*UserStartYearStats {
	if t == nil {
		t = &ViewerStats_Viewer_Statistics_Anime{}
	}
	return t.StartYears
}
func (t *ViewerStats_Viewer_Statistics_Anime) GetStatuses() []*UserStatusStats {
	if t == nil {
		t = &ViewerStats_Viewer_Statistics_Anime{}
	}
	return t.Statuses
}
func (t *ViewerStats_Viewer_Statistics_Anime) GetStudios() []*UserStudioStats {
	if t == nil {
		t = &ViewerStats_Viewer_Statistics_Anime{}
	}
	return t.Studios
}
func (t *ViewerStats_Viewer_Statistics_Manga) GetChaptersRead() int {
	if t == nil {
		t = &ViewerStats_Viewer_Statistics_Manga{}
	}
	return t.ChaptersRead
}
func (t *ViewerStats_Viewer_Statistics_Manga) GetCount() int {
	if t == nil {
		t = &ViewerStats_Viewer_Statistics_Manga{}
	}
	return t.Count
}
func (t *ViewerStats_Viewer_Statistics_Manga) GetFormats() []*UserFormatStats {
	if t == nil {
		t = &ViewerStats_Viewer_Statistics_Manga{}
	}
	return t.Formats
}
func (t *ViewerStats_Viewer_Statistics_Manga) GetGenres() []*UserGenreStats {
	if t == nil {
		t = &ViewerStats_Viewer_Statistics_Manga{}
	}
	return t.Genres
}
func (t *ViewerStats_Viewer_Statistics_Manga) GetMeanScore() float64 {
	if t == nil {
		t = &ViewerStats_Viewer_Statistics_Manga{}
	}
	return t.MeanScore
}
func (t *ViewerStats_Viewer_Statistics_Manga) GetReleaseYears() []*UserReleaseYearStats {
	if t == nil {
		t = &ViewerStats_Viewer_Statistics_Manga{}
	}
	return t.ReleaseYears
}
func (t *ViewerStats_Viewer_Statistics_Manga) GetScores() []*UserScoreStats {
	if t == nil {
		t = &ViewerStats_Viewer_Statistics_Manga{}
	}
	return t.Scores
}
func (t *ViewerStats_Viewer_Statistics_Manga) GetStartYears() []*UserStartYearStats {
	if t == nil {
		t = &ViewerStats_Viewer_Statistics_Manga{}
	}
	return t.StartYears
}
func (t *ViewerStats_Viewer_Statistics_Manga) GetStatuses() []*UserStatusStats {
	if t == nil {
		t = &ViewerStats_Viewer_Statistics_Manga{}
	}
	return t.Statuses
}
func (t *ViewerStats_Viewer_Statistics_Manga) GetStudios() []*UserStudioStats {
	if t == nil {
		t = &ViewerStats_Viewer_Statistics_Manga{}
	}
	return t.Studios
}

// Stats 汇总统计（照搬 anilist/stats.go）
type (
	Stats struct {
		AnimeStats *AnimeStats `json:"animeStats"`
		MangaStats *MangaStats `json:"mangaStats"`
	}

	AnimeStats struct {
		Count           int                     `json:"count"`
		MinutesWatched  int                     `json:"minutesWatched"`
		EpisodesWatched int                     `json:"episodesWatched"`
		MeanScore       float64                 `json:"meanScore"`
		Genres          []*UserGenreStats       `json:"genres"`
		Formats         []*UserFormatStats      `json:"formats"`
		Statuses        []*UserStatusStats      `json:"statuses"`
		Studios         []*UserStudioStats      `json:"studios"`
		Scores          []*UserScoreStats       `json:"scores"`
		StartYears      []*UserStartYearStats   `json:"startYears"`
		ReleaseYears    []*UserReleaseYearStats `json:"releaseYears"`
	}

	MangaStats struct {
		Count        int                     `json:"count"`
		ChaptersRead int                     `json:"chaptersRead"`
		MeanScore    float64                 `json:"meanScore"`
		Genres       []*UserGenreStats       `json:"genres"`
		Statuses     []*UserStatusStats      `json:"statuses"`
		Scores       []*UserScoreStats       `json:"scores"`
		StartYears   []*UserStartYearStats   `json:"startYears"`
		ReleaseYears []*UserReleaseYearStats `json:"releaseYears"`
	}
)

// GetStats 将 ViewerStats 摊平为 Stats（照搬语义，无 panic 保护——
// 原 util.HandlePanicInModuleWithError 属于 anilist 包调用点习惯，此处纯转换无 panic 面）。
func GetStats(stats *ViewerStats) *Stats {
	allStats := stats.GetViewer().GetStatistics()

	return &Stats{
		AnimeStats: &AnimeStats{
			Count:           allStats.GetAnime().GetCount(),
			MinutesWatched:  allStats.GetAnime().GetMinutesWatched(),
			EpisodesWatched: allStats.GetAnime().GetEpisodesWatched(),
			MeanScore:       allStats.GetAnime().GetMeanScore(),
			Genres:          allStats.GetAnime().GetGenres(),
			Formats:         allStats.GetAnime().GetFormats(),
			Statuses:        allStats.GetAnime().GetStatuses(),
			Studios:         allStats.GetAnime().GetStudios(),
			Scores:          allStats.GetAnime().GetScores(),
			StartYears:      allStats.GetAnime().GetStartYears(),
			ReleaseYears:    allStats.GetAnime().GetReleaseYears(),
		},
		MangaStats: &MangaStats{
			Count:        allStats.GetManga().GetCount(),
			ChaptersRead: allStats.GetManga().GetChaptersRead(),
			MeanScore:    allStats.GetManga().GetMeanScore(),
			Genres:       allStats.GetManga().GetGenres(),
			Statuses:     allStats.GetManga().GetStatuses(),
			Scores:       allStats.GetManga().GetScores(),
			StartYears:   allStats.GetManga().GetStartYears(),
			ReleaseYears: allStats.GetManga().GetReleaseYears(),
		},
	}
}

//----------------------------------------------------------------------------------------------------------------------
// StudioDetails（≈ anilist.StudioDetails）

type StudioDetails struct {
	Studio *StudioDetails_Studio `json:"Studio,omitempty"`
}

type StudioDetails_Studio struct {
	ID                int                     `json:"id"`
	IsAnimationStudio bool                    `json:"isAnimationStudio"`
	Media             *StudioDetails_Studio_Media `json:"media,omitempty"`
	Name              string                  `json:"name"`
}

type StudioDetails_Studio_Media struct {
	Nodes []*Anime `json:"nodes,omitempty"`
}

func (t *StudioDetails) GetStudio() *StudioDetails_Studio {
	if t == nil {
		t = &StudioDetails{}
	}
	return t.Studio
}
func (t *StudioDetails_Studio) GetID() int {
	if t == nil {
		t = &StudioDetails_Studio{}
	}
	return t.ID
}
func (t *StudioDetails_Studio) GetIsAnimationStudio() bool {
	if t == nil {
		t = &StudioDetails_Studio{}
	}
	return t.IsAnimationStudio
}
func (t *StudioDetails_Studio) GetMedia() *StudioDetails_Studio_Media {
	if t == nil {
		t = &StudioDetails_Studio{}
	}
	return t.Media
}
func (t *StudioDetails_Studio) GetName() string {
	if t == nil {
		t = &StudioDetails_Studio{}
	}
	return t.Name
}
func (t *StudioDetails_Studio_Media) GetNodes() []*Anime {
	if t == nil {
		t = &StudioDetails_Studio_Media{}
	}
	return t.Nodes
}

//----------------------------------------------------------------------------------------------------------------------
// GetViewer（≈ anilist.GetViewer，当前登录用户）

type GetViewer struct {
	Viewer *GetViewer_Viewer `json:"Viewer,omitempty"`
}

type GetViewer_Viewer struct {
	Avatar      *GetViewer_Viewer_Avatar  `json:"avatar,omitempty"`
	BannerImage *string                   `json:"bannerImage,omitempty"`
	IsBlocked   *bool                     `json:"isBlocked,omitempty"`
	Name        string                    `json:"name"`
	Options     *GetViewer_Viewer_Options `json:"options,omitempty"`
}

type GetViewer_Viewer_Avatar struct {
	Large  *string `json:"large,omitempty"`
	Medium *string `json:"medium,omitempty"`
}

type GetViewer_Viewer_Options struct {
	AiringNotifications *bool   `json:"airingNotifications,omitempty"`
	DisplayAdultContent *bool   `json:"displayAdultContent,omitempty"`
	ProfileColor        *string `json:"profileColor,omitempty"`
}

func (t *GetViewer) GetViewer() *GetViewer_Viewer {
	if t == nil {
		t = &GetViewer{}
	}
	return t.Viewer
}
func (t *GetViewer_Viewer) GetAvatar() *GetViewer_Viewer_Avatar {
	if t == nil {
		t = &GetViewer_Viewer{}
	}
	return t.Avatar
}
func (t *GetViewer_Viewer) GetBannerImage() *string {
	if t == nil {
		t = &GetViewer_Viewer{}
	}
	return t.BannerImage
}
func (t *GetViewer_Viewer) GetIsBlocked() *bool {
	if t == nil {
		t = &GetViewer_Viewer{}
	}
	return t.IsBlocked
}
func (t *GetViewer_Viewer) GetName() string {
	if t == nil {
		t = &GetViewer_Viewer{}
	}
	return t.Name
}
func (t *GetViewer_Viewer) GetOptions() *GetViewer_Viewer_Options {
	if t == nil {
		t = &GetViewer_Viewer{}
	}
	return t.Options
}
func (t *GetViewer_Viewer_Avatar) GetLarge() *string {
	if t == nil {
		t = &GetViewer_Viewer_Avatar{}
	}
	return t.Large
}
func (t *GetViewer_Viewer_Avatar) GetMedium() *string {
	if t == nil {
		t = &GetViewer_Viewer_Avatar{}
	}
	return t.Medium
}
func (t *GetViewer_Viewer_Options) GetAiringNotifications() *bool {
	if t == nil {
		t = &GetViewer_Viewer_Options{}
	}
	return t.AiringNotifications
}
func (t *GetViewer_Viewer_Options) GetDisplayAdultContent() *bool {
	if t == nil {
		t = &GetViewer_Viewer_Options{}
	}
	return t.DisplayAdultContent
}
func (t *GetViewer_Viewer_Options) GetProfileColor() *string {
	if t == nil {
		t = &GetViewer_Viewer_Options{}
	}
	return t.ProfileColor
}
