package bangumi

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearchSubjects(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request, _ int) {
		writeJSON(t, w, SearchResult{
			Total: 1,
			Limit: 20,
			Data: []Subject{
				{ID: 94619, NameCN: "音声作品", Type: SubjectBook, Nsfw: true},
			},
		})
	}
	c := newTestClient(t, ts.URL)

	res, err := c.SearchSubjects(context.Background(), SearchSubjectsOpts{
		Keyword: "RJ123456",
		Filter: SearchFilter{
			Type:    []int{SubjectBook},
			Tag:     []string{"ASMR"},
			Nsfw:    boolPtr(true),
			AirDate: []string{">2026-01-01", "<2026-04-01"},
			Rating:  []string{">=8"},
		},
		Limit:  20,
		Offset: 0,
	})
	require.NoError(t, err)

	// 请求形态断言
	require.Equal(t, "POST", rs.lastMethod)
	require.Equal(t, "/v0/search/subjects", rs.lastPath)
	require.Equal(t, "20", rs.lastQuery.Get("limit"))
	require.Empty(t, rs.lastQuery.Get("offset"), "offset 为 0 时不应出现在 query")

	// 请求体断言
	var body struct {
		Keyword string `json:"keyword"`
		Filter  struct {
			Type    []int    `json:"type"`
			Tag     []string `json:"tag"`
			Nsfw    *bool    `json:"nsfw"`
			AirDate []string `json:"air_date"`
			Rating  []string `json:"rating"`
		} `json:"filter"`
	}
	require.NoError(t, json.Unmarshal(rs.lastBody, &body))
	require.Equal(t, "RJ123456", body.Keyword)
	require.Equal(t, []int{SubjectBook}, body.Filter.Type)
	require.Equal(t, []string{"ASMR"}, body.Filter.Tag)
	require.NotNil(t, body.Filter.Nsfw)
	require.True(t, *body.Filter.Nsfw)
	// air_date 必须是操作符+全日期字符串数组（实测 ">2026-07" 这类年月格式会被 400 拒绝）
	require.Equal(t, []string{">2026-01-01", "<2026-04-01"}, body.Filter.AirDate)
	require.Equal(t, []string{">=8"}, body.Filter.Rating)

	// 响应解析断言
	require.Equal(t, 1, res.Total)
	require.Equal(t, "音声作品", res.Data[0].NameCN)
	require.Equal(t, SubjectBook, res.Data[0].Type)
}

func TestSearchSubjectsOmitsEmptyFields(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = func(w http.ResponseWriter, _ *http.Request, _ int) {
		writeJSON(t, w, SearchResult{Total: 0})
	}
	c := newTestClient(t, ts.URL)

	_, err := c.SearchSubjects(context.Background(), SearchSubjectsOpts{Keyword: "无过滤"})
	require.NoError(t, err)

	require.Empty(t, rs.lastQuery, "limit/offset 为 0 时不应出现在 query")

	var body map[string]any
	require.NoError(t, json.Unmarshal(rs.lastBody, &body))
	require.Equal(t, "无过滤", body["keyword"])
	_, hasFilter := body["filter"]
	require.False(t, hasFilter, "空 filter 应被 omitempty 省略")
}
