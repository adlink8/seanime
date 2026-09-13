package asmr

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	apiasmr "seanime/internal/api/asmr"
	"github.com/stretchr/testify/require"
)

// TestDownloaderLandingAndSkip 验证下载落盘 + 已存在非空文件跳过（契约 §7）。
// 全程 httptest/离线，不触网。
func TestDownloaderLandingAndSkip(t *testing.T) {
	fileBody := []byte("audio-bytes-content")
	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/tracks/123" && r.URL.Query().Get("v") == "2":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[{"type":"audio","title":"01.mp3","mediaDownloadUrl":"` + ts.URL + `/dl/01.mp3"}]`))
		case r.URL.Path == "/dl/01.mp3":
			_, _ = w.Write(fileBody)
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client := apiasmr.New("", apiasmr.WithBaseURL(ts.URL), apiasmr.WithThrottleInterval(time.Millisecond))
	defer client.Close()

	localDir := t.TempDir()
	dl := NewDownloader(client, localDir, nil, nil)

	dl.Download(context.Background(), 123, "RJ12345")

	dest := filepath.Join(localDir, "RJ12345", "01.mp3")
	got, err := os.ReadFile(dest)
	require.NoError(t, err, "下载应落盘")
	require.Equal(t, fileBody, got)

	// 再次下载：已存在非空文件应跳过（不重下，内容不变）
	dl.Download(context.Background(), 123, "RJ12345")
	got2, err := os.ReadFile(dest)
	require.NoError(t, err)
	require.Equal(t, fileBody, got2, "跳过已存在文件，内容不变")
}

// TestDownloaderSkipEmptyFile 验证 size==0 的空文件会被覆盖重下（契约 §7 跳过"非空"）。
func TestDownloaderSkipEmptyFile(t *testing.T) {
	fileBody := []byte("audio-bytes-content")
	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/tracks/9" && r.URL.Query().Get("v") == "2":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[{"type":"audio","title":"x.mp3","mediaDownloadUrl":"` + ts.URL + `/dl/x.mp3"}]`))
		case r.URL.Path == "/dl/x.mp3":
			_, _ = w.Write(fileBody)
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client := apiasmr.New("", apiasmr.WithBaseURL(ts.URL), apiasmr.WithThrottleInterval(time.Millisecond))
	defer client.Close()

	localDir := t.TempDir()
	dest := filepath.Join(localDir, "RJ9", "x.mp3")
	require.NoError(t, os.MkdirAll(filepath.Dir(dest), 0o755))
	require.NoError(t, os.WriteFile(dest, nil, 0o644)) // 空文件

	dl := NewDownloader(client, localDir, nil, nil)
	dl.Download(context.Background(), 9, "RJ9")

	got, err := os.ReadFile(dest)
	require.NoError(t, err)
	require.Equal(t, fileBody, got, "空文件应被重下覆盖")
}
