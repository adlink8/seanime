package bangumi

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetSubject(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request, _ int) {
		writeJSON(t, w, Subject{
			ID:       94619,
			Type:     SubjectAnime,
			Name:     "新しい台本",
			NameCN:   "新的台本",
			Date:     "2025-04-01",
			Summary:  "这是一段中文简介",
			Platform: "TV",
			Nsfw:     true,
			Eps:      12,
			Tags:     []SubjectTag{{Name: "ASMR", Count: 3}},
		})
	}
	c := newTestClient(t, ts.URL)

	subject, err := c.GetSubject(context.Background(), 94619)
	require.NoError(t, err)

	require.Equal(t, "/v0/subjects/94619", rs.lastPath)
	require.Equal(t, "GET", rs.lastMethod)
	require.Equal(t, 94619, subject.ID)
	require.Equal(t, SubjectAnime, subject.Type)
	require.Equal(t, "新的台本", subject.NameCN)
	require.Equal(t, "新的台本", subject.NameCN)
	require.True(t, subject.Nsfw)
	require.Equal(t, 12, subject.Eps)
	require.Len(t, subject.Tags, 1)
	require.Equal(t, "ASMR", subject.Tags[0].Name)
}

func TestGetSubjectEpisodes(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request, _ int) {
		writeJSON(t, w, EpisodesResult{
			Total: 2,
			Data: []Episode{
				{ID: 1, Type: 0, Name: "ep1", NameCN: "第一话"},
				{ID: 2, Type: 1, Name: "OP", NameCN: "片头曲"},
			},
		})
	}
	c := newTestClient(t, ts.URL)

	epType := 0
	res, err := c.GetSubjectEpisodes(context.Background(), 12, &EpisodesFilter{Type: &epType})
	require.NoError(t, err)

	require.Equal(t, "/v0/subjects/12/episodes", rs.lastPath)
	require.Equal(t, "0", rs.lastQuery.Get("type"), "EpType 过滤应走 query")
	require.Equal(t, 2, res.Total)
	require.Equal(t, "第一话", res.Data[0].NameCN)
	require.Equal(t, "片头曲", res.Data[1].NameCN)
}

func TestGetSubjectEpisodesNoFilter(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request, _ int) {
		writeJSON(t, w, EpisodesResult{Total: 0})
	}
	c := newTestClient(t, ts.URL)

	_, err := c.GetSubjectEpisodes(context.Background(), 12, nil)
	require.NoError(t, err)
	require.Empty(t, rs.lastQuery, "无过滤条件时不应带 query")
}
