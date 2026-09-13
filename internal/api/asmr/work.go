package asmr

import (
	"context"
	"net/url"
	"strconv"
)

// WorkInfo 获取作品详情（GET /api/workInfo/{id}，id 为数字 work id）。
// 实测 "RJ..." 形式返回 404 {"error":"未找到该音声"}。
func (c *Client) WorkInfo(ctx context.Context, workID int) (*Asmr_Work, error) {
	var raw rawWork
	if err := c.doGet(ctx, "/workInfo/"+strconv.Itoa(workID), url.Values{}, &raw); err != nil {
		return nil, err
	}
	out := convertWork(raw)
	return &out, nil
}

// Tracks 获取作品音轨树（GET /api/tracks/{id}?v=2）。
// 实测根为 JSON 数组，folder 用 children 嵌套；这里转换为契约 Asmr_Track（folder 用 tracks 嵌套）。
func (c *Client) Tracks(ctx context.Context, workID int) ([]Asmr_Track, error) {
	var raw []rawTrackNode
	if err := c.doGet(ctx, "/tracks/"+strconv.Itoa(workID), url.Values{"v": {"2"}}, &raw); err != nil {
		return nil, err
	}
	return convertTracks(raw), nil
}
