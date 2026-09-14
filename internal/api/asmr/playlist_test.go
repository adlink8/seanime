package asmr

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestGetPlaylists 契约 AC-01：路径/query 正确、系统列表过滤、字段映射（03.8 §1）。
func TestGetPlaylists(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = serveFixture(t, "playlists.json")
	c := newTestClient(t, ts.URL)

	lists, err := c.GetPlaylists(context.Background(), 2)
	require.NoError(t, err)

	require.Equal(t, "GET", rs.lastMethod)
	require.Equal(t, "/playlist/get-playlists", rs.lastPath)
	require.Equal(t, "2", rs.lastQuery.Get("page"))

	// fixture 含 1 个系统列表 → 过滤后剩 2
	require.Len(t, lists, 2)
	require.Equal(t, "a2ceeaaa-27c7-41ff-9855-6c929dceca1c", lists[0].ID)
	require.Equal(t, "姐系", lists[0].Name)
	require.Equal(t, 4, lists[0].WorksCount)
	require.Equal(t, "1168990", lists[0].LatestWorkID)
	require.Equal(t, "https://example.com/cover.jpg?type=main", lists[0].CoverURL)
	require.Equal(t, "2025-03-08 00:25:23", lists[0].UpdatedAt)
}

// TestGetPlaylistWorks 契约 AC-02：信封响应复用 convertSearchResult（零新转换）。
func TestGetPlaylistWorks(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = serveFixture(t, "cloud_reviews_works.json")
	c := newTestClient(t, ts.URL)

	res, err := c.GetPlaylistWorks(context.Background(), "a2ceeaaa-27c7-41ff-9855-6c929dceca1c", 1, 12)
	require.NoError(t, err)

	require.Equal(t, "GET", rs.lastMethod)
	require.Equal(t, "/playlist/get-playlist-works", rs.lastPath)
	require.Equal(t, "a2ceeaaa-27c7-41ff-9855-6c929dceca1c", rs.lastQuery.Get("id"))
	require.Equal(t, "1", rs.lastQuery.Get("page"))
	require.Equal(t, "12", rs.lastQuery.Get("pageSize"))

	require.Len(t, res.Works, 2)
	require.Equal(t, "RJ01168990", res.Works[0].RjID)
	require.Equal(t, "義理あね。2", res.Works[0].Title)
	require.Equal(t, 2, res.PageInfo.Total)
}

// TestPlaylistMutateBody 契约 AC-03：add/remove 请求体恰 {id, works:[int]}。
func TestPlaylistMutateBody(t *testing.T) {
	for _, tc := range []struct {
		name   string
		call   func(c *Client) error
		method string
		path   string
	}{
		{"add", func(c *Client) error { return c.AddWorksToPlaylist(context.Background(), "pid-1", []int{390787}) }, "POST", "/playlist/add-works-to-playlist"},
		{"remove", func(c *Client) error { return c.RemoveWorksFromPlaylist(context.Background(), "pid-1", []int{390787}) }, "POST", "/playlist/remove-works-from-playlist"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rs, ts := newRecordingServer(t)
			var gotBody map[string]any
			rs.handler = func(w http.ResponseWriter, r *http.Request) {
				require.NoError(t, json.NewDecoder(r.Body).Decode(&gotBody))
				w.WriteHeader(http.StatusOK)
			}
			c := newTestClient(t, ts.URL)

			require.NoError(t, tc.call(c))
			require.Equal(t, tc.method, rs.lastMethod)
			require.Equal(t, tc.path, rs.lastPath)
			require.Len(t, gotBody, 2)
			require.Equal(t, "pid-1", gotBody["id"])
			require.Equal(t, []any{float64(390787)}, gotBody["works"])
		})
	}
}

// TestGetPopularWorks 契约：POST body 形状 {keyword,page,subtitle,localSubtitled}。
func TestGetPopularWorks(t *testing.T) {
	rs, ts := newRecordingServer(t)
	var gotBody map[string]any
	rs.handler = func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, json.NewDecoder(r.Body).Decode(&gotBody))
		data, _ := os.ReadFile("testdata/cloud_reviews_works.json")
		_, _ = w.Write(data)
	}
	c := newTestClient(t, ts.URL)

	res, err := c.GetPopularWorks(context.Background(), 3, true)
	require.NoError(t, err)

	require.Equal(t, "POST", rs.lastMethod)
	require.Equal(t, "/recommender/popular", rs.lastPath)
	require.Len(t, gotBody, 4)
	require.Equal(t, " ", gotBody["keyword"])
	require.Equal(t, float64(3), gotBody["page"])
	require.Equal(t, float64(1), gotBody["subtitle"])
	require.Equal(t, float64(0), gotBody["localSubtitled"])
	require.Len(t, res.Works, 2)
}

// TestGetItemNeighbors 契约：POST body 形状 {keyword:"",itemId,page,subtitle,localSubtitled}。
func TestGetItemNeighbors(t *testing.T) {
	rs, ts := newRecordingServer(t)
	var gotBody map[string]any
	rs.handler = func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, json.NewDecoder(r.Body).Decode(&gotBody))
		data, _ := os.ReadFile("testdata/cloud_reviews_works.json")
		_, _ = w.Write(data)
	}
	c := newTestClient(t, ts.URL)

	_, err := c.GetItemNeighbors(context.Background(), 390787, 1, false)
	require.NoError(t, err)

	require.Equal(t, "/recommender/item-neighbors", rs.lastPath)
	require.Len(t, gotBody, 5)
	require.Equal(t, "", gotBody["keyword"])
	require.Equal(t, float64(390787), gotBody["itemId"])
	require.Equal(t, float64(0), gotBody["subtitle"])
}

// TestListWorks 契约：GET /works query 形状。
func TestListWorks(t *testing.T) {
	rs, ts := newRecordingServer(t)
	rs.handler = serveFixture(t, "cloud_reviews_works.json")
	c := newTestClient(t, ts.URL)

	_, err := c.ListWorks(context.Background(), "create_date", "desc", 2, false)
	require.NoError(t, err)

	require.Equal(t, "GET", rs.lastMethod)
	require.Equal(t, "/works", rs.lastPath)
	require.Equal(t, "create_date", rs.lastQuery.Get("order"))
	require.Equal(t, "desc", rs.lastQuery.Get("sort"))
	require.Equal(t, "2", rs.lastQuery.Get("page"))
	require.Equal(t, "0", rs.lastQuery.Get("subtitle"))
}
