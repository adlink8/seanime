package scanner

import (
	"seanime/internal/api/metadata_provider"
	"seanime/internal/database/db"
	"seanime/internal/events"
	"seanime/internal/library/anime"
	"seanime/internal/media"
	"seanime/internal/platforms/platform"
	"seanime/internal/testmocks"
	"seanime/internal/testutil"
	"seanime/internal/util/limiter"
	"testing"

	"github.com/rs/zerolog"
)

const scannerTestLibraryDir = "E:/Anime"

type scannerTestWrapper struct {
	Env                *testutil.TestEnv
	Config             *testutil.Config
	Logger             *zerolog.Logger
	Database           *db.Database
	Platform           platform.Platform
	MetadataProvider   metadata_provider.Provider
	CompleteAnimeCache *media.CompleteAnimeCache
	AnilistRateLimiter *limiter.Limiter
	WSEventManager     events.WSEventManagerInterface
	LibraryDir         string
}

func newScannerFixtureWrapper(t testing.TB) *scannerTestWrapper {
	t.Helper()

	env := testutil.NewTestEnv(t)
	return newScannerWrapper(t, env, "")
}

func newScannerLiveWrapper(t testing.TB) *scannerTestWrapper {
	t.Helper()

	env := testutil.NewTestEnv(t, testutil.Anilist())

	return newScannerWrapper(t, env, "")
}

func newScannerWrapper(t testing.TB, env *testutil.TestEnv, username string) *scannerTestWrapper {
	t.Helper()

	logger := env.Logger()
	database := env.MustNewDatabase(logger)
	// Bangumi 锚点：测试基建改用 FakePlatform（原 AniList fixture client 已随包裁剪）。
	platform := testmocks.NewFakePlatformBuilder().Build()
	if username != "" {
		platform.SetUsername(username)
	}

	return &scannerTestWrapper{
		Env:                env,
		Config:             env.Config(),
		Logger:             logger,
		Database:           database,
		Platform:           platform,
		MetadataProvider:   metadata_provider.NewTestProviderWithEnv(env, database),
		CompleteAnimeCache: media.NewCompleteAnimeCache(),
		AnilistRateLimiter: limiter.NewAnilistLimiter(),
		WSEventManager:     events.NewMockWSEventManager(logger),
		LibraryDir:         scannerTestLibraryDir,
	}
}

func (h *scannerTestWrapper) LocalFiles(paths ...string) []*anime.LocalFile {
	localFiles := make([]*anime.LocalFile, 0, len(paths))
	for _, path := range paths {
		localFiles = append(localFiles, anime.NewLocalFile(path, h.LibraryDir))
	}

	return localFiles
}
