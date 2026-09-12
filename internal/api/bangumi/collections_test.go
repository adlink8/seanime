package bangumi

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetUserCollections(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request, _ int) {
		writeJSON(t, w, UserCollectionsResult{
			Total: 2,
			Data: []UserCollection{
				{SubjectID: 100, Type: CollectionDone, Rate: 8},
				{SubjectID: 200, Type: CollectionWish},
			},
		})
	}
	c := newTestClient(t, ts.URL)

	res, err := c.GetUserCollections(context.Background(), ListCollectionsOpts{
		Limit:       20,
		Offset:      0,
		SubjectType: intPtr(SubjectAnime),
	})
	require.NoError(t, err)

	require.Equal(t, "/v0/users/-/collections", rs.lastPath)
	require.Equal(t, "GET", rs.lastMethod)
	require.Equal(t, "20", rs.lastQuery.Get("limit"))
	require.Empty(t, rs.lastQuery.Get("offset"), "offset 为 0 时不应出现在 query")
	require.Equal(t, "2", rs.lastQuery.Get("subject_type"))

	require.Equal(t, 2, res.Total)
	require.Equal(t, CollectionDone, res.Data[0].Type)
	require.Equal(t, 8, res.Data[0].Rate)
}

func TestUpsertCollectionPost(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request, _ int) {
		w.WriteHeader(200)
	}
	c := newTestClient(t, ts.URL)

	err := c.UpsertCollection(context.Background(), 94619, CollectionUpsertBody{
		Type:    intPtr(CollectionDoing),
		Rate:    intPtr(7),
		Comment: strPtr("好看"),
		Private: boolPtr(false),
		Tags:    []string{"tag1"},
	})
	require.NoError(t, err)

	require.Equal(t, "POST", rs.lastMethod)
	require.Equal(t, "/v0/users/-/collections/94619", rs.lastPath)
	require.Equal(t, "application/json; charset=utf-8", rs.lastHeader.Get("Content-Type"))

	var body map[string]any
	require.NoError(t, json.Unmarshal(rs.lastBody, &body))
	require.Equal(t, float64(CollectionDoing), body["type"])
	require.Equal(t, float64(7), body["rate"])
	require.Equal(t, "好看", body["comment"])
	require.Equal(t, false, body["private"])
	require.Len(t, body["tags"].([]any), 1)
}

func TestPatchCollection(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request, _ int) {
		w.WriteHeader(200)
	}
	c, _ := newTestClientWithCache(t, ts.URL)

	// 只更新 rate，验证 PATCH 局部更新且不带其他字段
	err := c.PatchCollection(context.Background(), 94619, CollectionUpsertBody{Rate: intPtr(10)})
	require.NoError(t, err)

	require.Equal(t, "PATCH", rs.lastMethod)
	require.Equal(t, "/v0/users/-/collections/94619", rs.lastPath)

	var body map[string]any
	require.NoError(t, json.Unmarshal(rs.lastBody, &body))
	require.Equal(t, float64(10), body["rate"])
	_, hasType := body["type"]
	require.False(t, hasType, "未提供的字段不应出现在请求体")
}

func TestGetEpisodeCollections(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request, _ int) {
		writeJSON(t, w, EpisodeCollectionsResult{
			SubjectID: 94619,
			Total:     1,
			Data: []EpisodeCollection{
				{ID: 1, Type: 0, Name: "ep1", NameCN: "第一话", Status: ProgressWatched},
			},
		})
	}
	c := newTestClient(t, ts.URL)

	res, err := c.GetEpisodeCollections(context.Background(), 94619)
	require.NoError(t, err)

	require.Equal(t, "/v0/users/-/collections/94619/episodes", rs.lastPath)
	require.Equal(t, 94619, res.SubjectID)
	require.Equal(t, "第一话", res.Data[0].NameCN)
	require.Equal(t, ProgressWatched, res.Data[0].Status)
}

func TestPatchEpisodeCollections(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request, _ int) {
		w.WriteHeader(200)
	}
	c := newTestClient(t, ts.URL)

	err := c.PatchEpisodeCollections(context.Background(), 94619, []int{1, 2, 8}, ProgressWatched)
	require.NoError(t, err)

	require.Equal(t, "PATCH", rs.lastMethod)
	require.Equal(t, "/v0/users/-/collections/94619/episodes", rs.lastPath)

	var body struct {
		EpisodeIDs []int  `json:"episode_id"`
		Type       string `json:"type"`
	}
	require.NoError(t, json.Unmarshal(rs.lastBody, &body))
	require.Equal(t, []int{1, 2, 8}, body.EpisodeIDs)
	require.Equal(t, ProgressWatched, body.Type)
}
