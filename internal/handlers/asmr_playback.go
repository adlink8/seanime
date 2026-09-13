package handlers

import (
	"net/http"

	"seanime/internal/mediaplayers/mediaplayer"

	"github.com/labstack/echo/v4"
)

// AsmrPlaybackStatusDTO is the D3 response shape for the asmr playback status
// endpoint. It intentionally mirrors only the fields the front-end control bar
// needs and never reaches into tracking/WS semantics.
type AsmrPlaybackStatusDTO struct {
	Running     bool    `json:"running"`
	Playing     bool    `json:"playing"`
	CurrentTime float64 `json:"currentTime"`
	Duration    float64 `json:"duration"`
	Filepath    string  `json:"filepath"`
	Filename    string  `json:"filename"`
}

// asmrPlaybackMaxPosition is the sanity upper bound for seek targets (seconds).
// 7 days of seconds; guards against absurd payloads (D2).
const asmrPlaybackMaxPosition = 86400 * 7

// assembleAsmrPlaybackStatus is a pure, testable mapping of a Repository
// PlaybackStatus into the D3 DTO. Degradation rules (D3): default player is not
// mpv, status is nil, or PlaybackType is not "file" all yield running:false with
// zeroed fields and a 200 (no 4xx), so the front-end polls degrade silently.
//
// The `running` signal is supplied by the caller (Repository.IsRunning()) so the
// function stays free of Repository coupling and unit-testable; it captures the
// §0.4 behavior where a finished mpv file playback reports running:false while
// the other fields remain populated for the D8 auto-complete heuristic.
func assembleAsmrPlaybackStatus(status *mediaplayer.PlaybackStatus, defaultPlayer string, running bool) AsmrPlaybackStatusDTO {
	if status == nil || defaultPlayer != "mpv" || status.PlaybackType != mediaplayer.PlaybackTypeFile {
		return AsmrPlaybackStatusDTO{Running: false}
	}
	return AsmrPlaybackStatusDTO{
		Running:     running,
		Playing:     status.Playing,
		CurrentTime: status.CurrentTimeInSeconds,
		Duration:    status.DurationInSeconds,
		Filepath:    status.Filepath,
		Filename:    status.Filename,
	}
}

// asmrPlaybackValidatePosition validates the seek target (D2/D4): the position
// must be a non-negative absolute second below the sanity ceiling.
func asmrPlaybackValidatePosition(position float64) bool {
	return position >= 0 && position < asmrPlaybackMaxPosition
}

// HandleAsmrPlaybackStatus
//
//	@summary returns the current asmr (file) playback status from the default media player.
//	@desc Pure forward of Repository.PullStatus. Degrades to running:false (HTTP 200) when the default player is not mpv, there is no status, or the playback type is not "file". Never touches tracking/WS.
//	@route /api/v1/asmr/playback/status [GET]
//	@returns AsmrPlaybackStatusDTO
func (h *Handler) HandleAsmrPlaybackStatus(c echo.Context) error {
	repo := h.App.MediaPlayerRepository
	if repo == nil {
		return h.RespondWithData(c, assembleAsmrPlaybackStatus(nil, "", false))
	}

	status, _ := repo.PullStatus()
	defaultPlayer := repo.GetDefault()
	running := repo.IsRunning()

	return h.RespondWithData(c, assembleAsmrPlaybackStatus(status, defaultPlayer, running))
}

// HandleAsmrPlaybackPause
//
//	@summary pauses or resumes the current asmr (file) playback.
//	@desc Pure forward of Repository.Pause/Resume. Returns 409 when there is no active playback. Never touches tracking/WS.
//	@route /api/v1/asmr/playback/pause [POST]
//	@returns bool
func (h *Handler) HandleAsmrPlaybackPause(c echo.Context) error {
	type body struct {
		Paused bool `json:"paused"`
	}
	b := new(body)
	if err := c.Bind(b); err != nil {
		return h.RespondWithError(c, err)
	}

	repo := h.App.MediaPlayerRepository
	if repo == nil || !repo.IsRunning() {
		return echo.NewHTTPError(http.StatusConflict, "no active playback")
	}

	var err error
	if b.Paused {
		err = repo.Pause()
	} else {
		err = repo.Resume()
	}
	if err != nil {
		return h.RespondWithError(c, err)
	}

	return h.RespondWithData(c, true)
}

// HandleAsmrPlaybackSeek
//
//	@summary seeks the current asmr (file) playback to an absolute position in seconds.
//	@desc Pure forward of Repository.SeekTo. Validates the position (>=0, < 7 days) and returns 409 when there is no active playback. Never touches tracking/WS.
//	@route /api/v1/asmr/playback/seek [POST]
//	@returns bool
func (h *Handler) HandleAsmrPlaybackSeek(c echo.Context) error {
	type body struct {
		Position float64 `json:"position"`
	}
	b := new(body)
	if err := c.Bind(b); err != nil {
		return h.RespondWithError(c, err)
	}

	if !asmrPlaybackValidatePosition(b.Position) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid position")
	}

	repo := h.App.MediaPlayerRepository
	if repo == nil || !repo.IsRunning() {
		return echo.NewHTTPError(http.StatusConflict, "no active playback")
	}

	if err := repo.SeekTo(b.Position); err != nil {
		return h.RespondWithError(c, err)
	}

	return h.RespondWithData(c, true)
}
