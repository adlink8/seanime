package local

import (
	"seanime/internal/api/metadata_provider"
	"seanime/internal/database/db"
	"seanime/internal/database/db_bridge"
	"seanime/internal/database/models"
	"seanime/internal/events"
	"seanime/internal/library/anime"
	"seanime/internal/manga"
	"seanime/internal/platforms/platform"
	"seanime/internal/testmocks"
	"seanime/internal/testutil"
	"seanime/internal/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func NewTestManager(t *testing.T, db *db.Database) Manager {
	env := testutil.NewTestEnv(t)

	logger := env.Logger()
	metadataProvider := metadata_provider.NewTestProviderWithEnv(env, db)
	metadataProviderRef := util.NewRef(metadataProvider)
	mangaRepository := manga.NewTestRepositoryWithEnv(env, db)

	wsEventManager := events.NewMockWSEventManager(logger)
	// Bangumi 锚点：测试基建改用 FakePlatform（原 AniList fixture client 已随包裁剪）。
	fakePlatform := testmocks.NewFakePlatformBuilder().Build()
	anilistPlatformRef := util.NewRef[platform.Platform](fakePlatform)

	localDir := env.MustMkdirData("offline")
	assetsDir := env.MustMkdirData("offline", "assets")

	var localFilesCount int64
	err := db.Gorm().Model(&models.LocalFiles{}).Count(&localFilesCount).Error
	require.NoError(t, err)
	if localFilesCount == 0 {
		_, err = db_bridge.InsertLocalFiles(db, make([]*anime.LocalFile, 0))
		require.NoError(t, err)
	}

	m, err := NewManager(&NewManagerOptions{
		LocalDir:            localDir,
		AssetDir:            assetsDir,
		Logger:              logger,
		MetadataProviderRef: metadataProviderRef,
		MangaRepository:     mangaRepository,
		Database:            db,
		WSEventManager:      wsEventManager,
		AnilistPlatformRef:  anilistPlatformRef,
		IsOffline:           false,
	})
	require.NoError(t, err)

	return m
}
