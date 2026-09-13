package playbackmanager

import (
	"runtime"
	"seanime/internal/events"
	"seanime/internal/library/anime"
	medialib "seanime/internal/media"
	"seanime/internal/mediaplayers/mediaplayer"
	"seanime/internal/testmocks"
	"testing"

	"github.com/samber/mo"
	"github.com/stretchr/testify/require"
)

// TestAsmrIsAsmrLocalFile covers the pure detection helper: in-dir hit, sibling-prefix
// rejection, and case/backslash normalization tolerance (Windows semantics).
func TestAsmrIsAsmrLocalFile(t *testing.T) {
	// ① hit: file inside the asmr local dir
	require.True(t, isAsmrLocalFile("/asmr-local", "/asmr-local/RJ01657200/01/track.mp3"))
	// the directory root itself counts as inside
	require.True(t, isAsmrLocalFile("/asmr-local", "/asmr-local"))

	// ② a similar prefix that is NOT a subpath must NOT match (boundary check).
	require.False(t, isAsmrLocalFile("/asmr-local", "/asmr-localevil/RJ/track.mp3"))
	require.False(t, isAsmrLocalFile("/asmr-local", "/asmr-local-sibling/track.mp3"))
	// a path completely outside the dir
	require.False(t, isAsmrLocalFile("/asmr-local", "/Anime/Frieren/1.mkv"))

	// empty dir disables detection entirely
	require.False(t, isAsmrLocalFile("", "/asmr-local/RJ/track.mp3"))

	// ③ case + backslash normalization tolerance. util.NormalizePath lowercases and converts
	// separators to slash ONLY on windows; on other platforms it is a no-op, so the Windows
	// semantics samples are platform-gated (CI runs this package on ubuntu).
	if runtime.GOOS == "windows" {
		require.True(t, isAsmrLocalFile("C:\\ASMR-LOCAL", "c:/asmr-local/RJ/track.mp3"))
		require.True(t, isAsmrLocalFile("/asmr-local", "/ASMR-LOCAL/RJ/track.mp3"))
		require.True(t, isAsmrLocalFile("C:/AsmrLocal", "C:\\AsmrLocal\\RJ\\track.mp3"))
	} else {
		// non-windows: separators are already slashes and case is preserved
		require.True(t, isAsmrLocalFile("/asmr-local", "/asmr-local/RJ/track.mp3"))
		require.False(t, isAsmrLocalFile("/asmr-local", "/ASMR-LOCAL/RJ/track.mp3"))
	}
}

// TestAsmrTrackingStartedDoesNotCancelAndKeepsOptionsNone drives the ASMR branch directly:
// an ASMR-local filepath must route to handleAsmrTrackingStarted, which sets the playback
// status, keeps the three anime-local Options at None, and never triggers Cancel / error toast.
func TestAsmrTrackingStartedDoesNotCancelAndKeepsOptionsNone(t *testing.T) {
	h := newPlaybackManagerTestWrapper(t)
	h.playbackManager.AsmrLocalDir = "/asmr-local"

	// seed a leftover anime session to prove the branch explicitly resets the Options to None.
	h.playbackManager.currentMediaListEntry = mo.Some(&medialib.AnimeListEntry{
		Media:    testmocks.NewBaseAnime(1, "leftover"),
		Progress: new(0),
	})
	h.playbackManager.currentLocalFile = mo.Some(&anime.LocalFile{})
	h.playbackManager.currentLocalFileWrapperEntry = mo.Some(&anime.LocalFileWrapperEntry{})

	subscriber := h.playbackManager.SubscribeToPlaybackStatus("unit-asmr")
	status := &mediaplayer.PlaybackStatus{
		Filename:             "track.mp3",
		Filepath:             "/asmr-local/RJ01657200/01/track.mp3",
		CompletionPercentage: 0.01,
		CurrentTimeInSeconds: 3,
		DurationInSeconds:    291,
		PlaybackType:         mediaplayer.PlaybackTypeFile,
	}

	// If the asmr branch were NOT taken, getLocalFilePlaybackDetails would fail for this
	// non-anime file and the original failure path would send an ErrorToast immediately
	// followed by MediaPlayerRepository.Cancel(). Absence of the ErrorToast therefore proves
	// Cancel was never reached.
	h.playbackManager.handleTrackingStarted(status)

	// status is recorded
	require.Same(t, status, h.playbackManager.currentMediaPlaybackStatus)
	// playback type is local file
	require.Equal(t, LocalFilePlayback, h.playbackManager.currentPlaybackType)
	// the three anime-local Options are explicitly reset to None (no stale anime state)
	require.True(t, h.playbackManager.currentMediaListEntry.IsAbsent())
	require.True(t, h.playbackManager.currentLocalFile.IsAbsent())
	require.True(t, h.playbackManager.currentLocalFileWrapperEntry.IsAbsent())
	// no error toast -> failure path + Cancel not triggered
	require.Equal(t, 0, h.wsEventManager.count(events.ErrorToast))
	// the asmr branch fired its own tracking-started event
	require.Equal(t, 1, h.wsEventManager.count(events.PlaybackManagerProgressTrackingStarted))
	require.Equal(t, events.PlaybackManagerProgressTrackingStarted, h.wsEventManager.lastType())

	// subscribers were notified (status changed + video started)
	changedEvent := expectPlaybackEvent[PlaybackStatusChangedEvent](t, subscriber.EventCh)
	require.Equal(t, "track.mp3", changedEvent.Status.Filename)
	startedEvent := expectPlaybackEvent[VideoStartedEvent](t, subscriber.EventCh)
	require.Equal(t, "track.mp3", startedEvent.Filename)
	require.Equal(t, "/asmr-local/RJ01657200/01/track.mp3", startedEvent.Filepath)
}

// TestAsmrEmptyDirDoesNotRoute confirms the zero-value (empty AsmrLocalDir) leaves detection
// disabled, so isAsmrLocalFile behaves as a no-op regardless of the path. This guarantees the
// two test constructors that omit the field keep unchanged behavior.
func TestAsmrEmptyDirDoesNotRoute(t *testing.T) {
	h := newPlaybackManagerTestWrapper(t)
	require.Empty(t, h.playbackManager.AsmrLocalDir)
	require.False(t, h.playbackManager.isAsmrLocalFile("/asmr-local/RJ/track.mp3"))
	require.False(t, h.playbackManager.isAsmrLocalFile("/Anime/Frieren/1.mkv"))
}
