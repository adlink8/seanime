package handlers

import (
	"testing"

	"seanime/internal/mediaplayers/mediaplayer"

	"github.com/stretchr/testify/assert"
)

// TestAssembleAsmrPlaybackStatus covers D3: full branch matrix for the status
// assembly pure function. Degrade cases (non-mpv / nil status / non-file) must
// yield running:false and zeroed values; the normal mpv/file case maps fields
// and honors the externally-derived running signal (so a finished playback
// reports running:false while still exposing duration/filepath for the
// auto-complete heuristic, §0.4 / D8).
func TestAssembleAsmrPlaybackStatus(t *testing.T) {
	fileStatus := &mediaplayer.PlaybackStatus{
		Playing:              true,
		CurrentTimeInSeconds: 12.5,
		DurationInSeconds:    120.0,
		Filepath:             "/a/b.flac",
		Filename:             "b.flac",
		PlaybackType:         mediaplayer.PlaybackTypeFile,
	}

	t.Run("normal mpv file playback maps fields", func(t *testing.T) {
		dto := assembleAsmrPlaybackStatus(fileStatus, "mpv", true)
		assert.True(t, dto.Running)
		assert.True(t, dto.Playing)
		assert.Equal(t, 12.5, dto.CurrentTime)
		assert.Equal(t, 120.0, dto.Duration)
		assert.Equal(t, "/a/b.flac", dto.Filepath)
		assert.Equal(t, "b.flac", dto.Filename)
	})

	t.Run("non-mpv default player degrades", func(t *testing.T) {
		dto := assembleAsmrPlaybackStatus(fileStatus, "vlc", true)
		assert.False(t, dto.Running)
		assert.False(t, dto.Playing)
		assert.Zero(t, dto.CurrentTime)
		assert.Zero(t, dto.Duration)
		assert.Empty(t, dto.Filepath)
		assert.Empty(t, dto.Filename)
	})

	t.Run("nil status degrades", func(t *testing.T) {
		dto := assembleAsmrPlaybackStatus(nil, "mpv", true)
		assert.False(t, dto.Running)
		assert.False(t, dto.Playing)
		assert.Zero(t, dto.CurrentTime)
		assert.Empty(t, dto.Filepath)
	})

	t.Run("non-file playback type degrades", func(t *testing.T) {
		streamStatus := &mediaplayer.PlaybackStatus{
			PlaybackType: mediaplayer.PlaybackTypeStream,
		}
		dto := assembleAsmrPlaybackStatus(streamStatus, "mpv", true)
		assert.False(t, dto.Running)
		assert.False(t, dto.Playing)
		assert.Empty(t, dto.Filepath)
	})

	t.Run("finished playback keeps values but running false", func(t *testing.T) {
		dto := assembleAsmrPlaybackStatus(fileStatus, "mpv", false)
		assert.False(t, dto.Running)
		// values retained so the auto-complete heuristic (D8) still works
		assert.Equal(t, 120.0, dto.Duration)
		assert.Equal(t, "/a/b.flac", dto.Filepath)
	})
}

// TestAsmrPlaybackValidatePosition covers the seek position sanity guard
// (D2/D4): 0 <= position < 86400*7.
func TestAsmrPlaybackValidatePosition(t *testing.T) {
	assert.True(t, asmrPlaybackValidatePosition(0))
	assert.True(t, asmrPlaybackValidatePosition(100.5))
	assert.True(t, asmrPlaybackValidatePosition(86400*7-1))
	assert.False(t, asmrPlaybackValidatePosition(-1))
	assert.False(t, asmrPlaybackValidatePosition(86400*7))
	assert.False(t, asmrPlaybackValidatePosition(86400*7+1))
}
