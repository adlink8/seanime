package bangumi

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"seanime/internal/media"
)

//----------------------------------------------------------------------------------------------------------------------
// 收藏状态映射

func TestListStatusFromBangumiType(t *testing.T) {
	cases := []struct {
		in   int
		want media.MediaListStatus
	}{
		{CollectionWish, media.MediaListStatusPlanning},    // 想看
		{CollectionDoing, media.MediaListStatusCurrent},    // 在看
		{CollectionDone, media.MediaListStatusCompleted},   // 看过
		{CollectionOnHold, media.MediaListStatusPaused},    // 搁置
		{CollectionDropped, media.MediaListStatusDropped},  // 抛弃
	}
	for _, c := range cases {
		got := ListStatusFromBangumiType(c.in)
		require.NotNil(t, got, "类型 %d 不应映射为 nil", c.in)
		require.Equal(t, c.want, *got)
	}
	require.Nil(t, ListStatusFromBangumiType(99), "未知类型应返回 nil（保守不猜测）")
}

//----------------------------------------------------------------------------------------------------------------------
// Subject 桥接

func fixtureBangumiSubject() *Subject {
	return &Subject{
		ID:       94607,
		Type:     SubjectAnime,
		Name:     "葬送のフリーレン",
		NameCN:   "葬送的芙莉莲",
		Date:     "2023-09-29",
		Platform: "TV",
		Summary:  "简介",
		Eps:      28,
		Nsfw:     false,
		Images: SubjectImages{
			Large:  "https://lain.bgm.tv/l.jpg",
			Common: "https://lain.bgm.tv/c.jpg",
		},
		Rating: &SubjectRating{Score: 9.2},
		Tags:   []SubjectTag{{Name: "奇幻"}},
	}
}

func TestSubjectToMedia(t *testing.T) {
	out := SubjectToMedia(fixtureBangumiSubject())
	require.NotNil(t, out)
	require.Equal(t, 94607, out.ID)
	require.Equal(t, SubjectAnime, out.Type)
	require.Equal(t, "葬送のフリーレン", out.Name)
	require.Equal(t, "葬送的芙莉莲", out.NameCN)
	require.Equal(t, "TV", out.Platform)
	require.Equal(t, 28, out.Eps)
	require.False(t, out.NSFW)
	require.Equal(t, "2023-09-29", out.Date)
	require.NotNil(t, out.Images)
	require.Equal(t, "https://lain.bgm.tv/c.jpg", out.Images.Common)
	require.NotNil(t, out.Rating)
	require.InDelta(t, 9.2, out.Rating.Score, 0.001)
	require.Nil(t, SubjectToMedia(nil))
}

//----------------------------------------------------------------------------------------------------------------------
// 收藏 → media.AnimeCollection / MangaCollection

func fixtureUserCollections() *UserCollectionsResult {
	watching := CollectionDoing
	wish := CollectionWish
	return &UserCollectionsResult{
		Total: 2,
		Data: []UserCollection{
			{
				SubjectID: 94607,
				Subject:   fixtureBangumiSubject(),
				Type:      watching,
				Rate:      8,
				EpStatus:  12,
			},
			{
				SubjectID: 10001,
				Subject: &Subject{
					ID:     10001,
					Type:   SubjectAnime,
					Name:   "新世紀エヴァンゲリオン",
					NameCN: "新世纪福音战士",
					Date:   "1995-10-04",
					Images: SubjectImages{Large: "https://lain.bgm.tv/e.jpg"},
				},
				Type: wish,
				Rate: 10,
			},
		},
	}
}

func TestAnimeCollectionFromUserCollections(t *testing.T) {
	col := AnimeCollectionFromUserCollections(fixtureUserCollections())
	require.NotNil(t, col)
	require.NotNil(t, col.MediaListCollection)
	require.Len(t, col.MediaListCollection.Lists, 2)

	// 按 status 分组：CURRENT(在看) 与 PLANNING(想看)
	statuses := map[media.MediaListStatus]int{}
	for _, l := range col.MediaListCollection.Lists {
		statuses[*l.Status] = len(l.Entries)
		require.NotNil(t, l.Name)
		require.Equal(t, string(*l.Status), *l.Name)
	}
	require.Equal(t, 1, statuses[media.MediaListStatusCurrent])
	require.Equal(t, 1, statuses[media.MediaListStatusPlanning])

	// 在看条目：进度 12、评分 8×10=80、media 实体已填充
	entry := col.MediaListCollection.Lists[0].Entries[0]
	require.Equal(t, 94607, entry.ID)
	require.Equal(t, 12, *entry.Progress)
	require.InDelta(t, 80.0, *entry.Score, 0.001)
	require.NotNil(t, entry.Media)
	require.Equal(t, 94607, entry.Media.ID)
	require.Equal(t, "葬送的芙莉莲", entry.Media.NameCN)

	// 想看条目：无进度、评分 100
	entry2 := col.MediaListCollection.Lists[1].Entries[0]
	require.Nil(t, entry2.Progress)
	require.InDelta(t, 100.0, *entry2.Score, 0.001)

	// 序列化形状：MediaListCollection 首字母大写（GraphQL 包装层习惯）+ lists/entries 小写
	data, err := json.Marshal(col)
	require.NoError(t, err)
	var probe map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(data, &probe))
	require.Contains(t, probe, "MediaListCollection")
	var mlc map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(probe["MediaListCollection"], &mlc))
	require.Contains(t, mlc, "lists")
}

func TestMangaCollectionFromUserCollections(t *testing.T) {
	res := fixtureUserCollections()
	res.Data[0].SubjectID = 20001
	res.Data[0].Subject = &Subject{
		ID:   20001,
		Type: SubjectBook,
		Name: "よふかしのうた",
		NameCN: "彻夜之歌",
		Eps:  100,
	}
	col := MangaCollectionFromUserCollections(res)
	require.NotNil(t, col)
	require.Len(t, col.MediaListCollection.Lists, 2)

	entry := col.MediaListCollection.Lists[0].Entries[0]
	require.Equal(t, 20001, entry.ID)
	require.NotNil(t, entry.Media)
	require.Equal(t, "彻夜之歌", entry.Media.NameCN)
	require.Equal(t, 12, *entry.Progress)

	require.Nil(t, MangaCollectionFromUserCollections(nil))
}

func TestAnimeCollectionFromUserCollections_UnknownStatusSkipped(t *testing.T) {
	res := &UserCollectionsResult{
		Data: []UserCollection{{SubjectID: 1, Type: 99}},
	}
	col := AnimeCollectionFromUserCollections(res)
	require.NotNil(t, col)
	require.Empty(t, col.MediaListCollection.Lists)
}

//----------------------------------------------------------------------------------------------------------------------
// 新端点：按用户取收藏 + 相关条目

func TestGetUserCollectionsByUser(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, r *http.Request, call int) {
		writeJSON(t, w, UserCollectionsResult{Total: 0, Data: []UserCollection{}})
	}
	c := newTestClient(t, ts.URL)

	two := 2
	_, err := c.GetUserCollectionsByUser(context.Background(), "12345", UserCollectionsOpts{
		SubjectType: &two,
		Limit:       50,
	})
	require.NoError(t, err)
	require.Equal(t, http.MethodGet, rs.lastMethod)
	require.Equal(t, "/v0/users/12345/collections", rs.lastPath)
	require.Equal(t, "2", rs.lastQuery.Get("subject_type"))
	require.Equal(t, "50", rs.lastQuery.Get("limit"))

	// type 过滤参数
	wish := CollectionWish
	_, err = c.GetUserCollectionsByUser(context.Background(), "-", UserCollectionsOpts{Type: &wish})
	require.NoError(t, err)
	require.Equal(t, "/v0/users/-/collections", rs.lastPath)
	require.Equal(t, "1", rs.lastQuery.Get("type"))
}

func TestGetRelatedSubjects(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, r *http.Request, call int) {
		writeJSON(t, w, []RelatedSubject{
			{
				ID:          1,
				Type:        2, // 续集
				Name:        "葬送のフリーレン",
				NameCN:      "葬送的芙莉莲",
				SubjectID:   400602,
				SubjectType: SubjectAnime,
			},
		})
	}
	c := newTestClient(t, ts.URL)

	rels, err := c.GetRelatedSubjects(context.Background(), 94607)
	require.NoError(t, err)
	require.Len(t, rels, 1)
	require.Equal(t, 2, rels[0].Type)
	require.Equal(t, 400602, rels[0].SubjectID)
	require.Equal(t, "/v0/subjects/94607/subjects", rs.lastPath)
}
