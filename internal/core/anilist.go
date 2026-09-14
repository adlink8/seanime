package core

import (
	"context"
	"errors"
	"seanime/internal/api/bangumi"
	"seanime/internal/database/models"
	"seanime/internal/events"
	"seanime/internal/media"
	"seanime/internal/platforms/bangumi_platform"
	"seanime/internal/platforms/platform"
	"seanime/internal/platforms/simulated_platform"
	"seanime/internal/plugin"
	"seanime/internal/user"
	"seanime/internal/util"
	"time"

	"github.com/goccy/go-json"
)

// GetUser returns the currently logged-in user or a simulated one.
func (a *App) GetUser() *user.User {
	if a.user == nil {
		return user.NewSimulatedUser()
	}
	return a.user
}

// GetUsername returns the username of the currently logged-in user
func (a *App) GetUsername() string {
	if a.user == nil {
		return ""
	}
	if a.user.Viewer == nil {
		return ""
	}
	return a.user.Viewer.GetName()
}

func (a *App) GetUserAnilistToken() string {
	if a.user == nil || a.user.Token == user.SimulatedUserToken {
		return ""
	}

	return a.user.Token
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// UpdatePlatform changes the current platform to the provided one.
func (a *App) UpdatePlatform(platform platform.Platform) {
	if a.AnilistPlatformRef.IsPresent() {
		a.AnilistPlatformRef.Get().Close()
	}
	a.AnilistPlatformRef.Set(platform)
	a.AddOnRefreshAnilistCollectionFunc("anilist-platform", func() {
		a.AnilistPlatformRef.Get().ClearCache()
	})
}

// UpdateAnilistClientToken will update the Bangumi client token.
// This function should be called when a user logs in.
// Bangumi 锚点：token 沿用 account 表存储，函数名保留（Phase 4 统一清理命名）。
// 代理配置沿用启动时读取的 server.proxyURL。
func (a *App) UpdateAnilistClientToken(token string) {
	a.BangumiClientRef.Set(bangumi.New(token, bangumi.WithProxyURL(a.Config.Server.ProxyURL)))
}

// UseOfficialAnilistClient Bangumi 锚点下 AniList client 运行时切换已停用。
// 保留函数以满足插件 API 装配（plugin.AnilistActions），调用时返回错误。
func (a *App) UseOfficialAnilistClient() error {
	return errors.New("AniList client runtime switching is not supported: the data provider is now Bangumi")
}

// UseCustomAnimeClient Bangumi 锚点下 AniList client 运行时切换已停用。
// 保留函数以满足插件 API 装配（plugin.AnilistActions），调用时返回错误。
func (a *App) UseCustomAnimeClient(config plugin.CustomClientConfig) error {
	_ = config
	return errors.New("AniList client runtime switching is not supported: the data provider is now Bangumi")
}

func (a *App) LoginToAnilist(token string) error {
	if token == "" {
		return errors.New("token is empty")
	}

	a.UpdateAnilistClientToken(token)

	// Bangumi 锚点：原 AniList GetViewer 校验改为 Bangumi GetMe。
	me, err := a.BangumiClientRef.Get().GetMe(context.Background())
	if err != nil {
		a.Logger.Error().Msg("Could not authenticate to Bangumi")
		return err
	}

	if len(me.Username) == 0 {
		return errors.New("could not find user")
	}

	// user.User.Viewer 形状沿用 media.GetViewer_Viewer（前端契约不破），Bangumi 昵称填入 Name。
	viewer := &media.GetViewer_Viewer{
		Name: me.Username,
	}

	bytes, err := json.Marshal(viewer)
	if err != nil {
		a.Logger.Err(err).Msg("scan: could not save local files")
	}

	_, err = a.Database.UpsertAccount(&models.Account{
		BaseModel: models.BaseModel{
			ID:        1,
			UpdatedAt: time.Now(),
		},
		Username: me.Username,
		Token:    token,
		Viewer:   bytes,
	})
	if err != nil {
		return err
	}

	a.Logger.Info().Msg("app: Authenticated to Bangumi")

	bangumiPlatform := bangumi_platform.NewBangumiPlatform(a.BangumiClientRef.Get(), a.BangumiCacheDir, a.ExtensionBankRef, a.Logger, a.Database, a.LogoutFromAnilist)
	a.UpdatePlatform(bangumiPlatform)
	a.initMappingService(bangumiPlatform) // Phase 4：异步装配 bangumi→anidb 映射服务

	a.InitOrRefreshAnilistData()
	a.InitOrRefreshModules()

	go func() {
		defer util.HandlePanicThen(func() {})
		a.InitOrRefreshTorrentstreamSettings()
		a.InitOrRefreshMediastreamSettings()
		a.InitOrRefreshDebridSettings()
	}()

	return nil
}

// LogoutFromAnilist clears the AniList token and switches to the simulated platform.
// This is called internally when the token is detected as invalid.
func (a *App) LogoutFromAnilist() {
	// prevent multiple concurrent calls (e.g. from parallel failing requests)
	if !a.logoutInProgress.CompareAndSwap(false, true) {
		return
	}
	defer a.logoutInProgress.Store(false)

	a.UpdateAnilistClientToken("")

	simulatedPlatform, err := simulated_platform.NewSimulatedPlatform(a.LocalManager, a.BangumiClientRef.Get(), a.BangumiCacheDir, a.ExtensionBankRef, a.Logger, a.Database)
	if err != nil {
		a.Logger.Error().Err(err).Msg("app: Failed to create simulated platform during auto-logout")
	} else {
		a.UpdatePlatform(simulatedPlatform)
	}

	_, _ = a.Database.UpsertAccount(&models.Account{
		BaseModel: models.BaseModel{
			ID:        1,
			UpdatedAt: time.Now(),
		},
		Username: "",
		Token:    "",
		Viewer:   nil,
	})

	a.Logger.Debug().Msg("app: Logged out from AniList, switched to simulated platform")

	a.InitOrRefreshModules()
	a.InitOrRefreshAnilistData()
}

// GetAnimeCollection returns the user's Anilist collection if it in the cache, otherwise it queries Anilist for the user's collection.
// When bypassCache is true, it will always query Anilist for the user's collection
func (a *App) GetAnimeCollection(bypassCache bool) (*media.AnimeCollection, error) {
	return a.AnilistPlatformRef.Get().GetAnimeCollection(context.Background(), bypassCache)
}

// GetRawAnimeCollection is the same as GetAnimeCollection but returns the raw collection that includes custom lists
func (a *App) GetRawAnimeCollection(bypassCache bool) (*media.AnimeCollection, error) {
	return a.AnilistPlatformRef.Get().GetRawAnimeCollection(context.Background(), bypassCache)
}

func (a *App) SyncAnilistToSimulatedCollection() {
	if a.LocalManager != nil &&
		!a.GetUser().IsSimulated &&
		a.Settings != nil &&
		a.Settings.Library != nil &&
		a.Settings.Library.AutoSyncToLocalAccount {
		_ = a.LocalManager.SynchronizeAnilistToSimulatedCollection()
	}
}

// RefreshAnimeCollection queries Anilist for the user's collection
func (a *App) RefreshAnimeCollection() (*media.AnimeCollection, error) {
	go func() {
		a.OnRefreshAnilistCollectionFuncs.Range(func(key string, f func()) bool {
			go f()
			return true
		})
	}()

	ret, err := a.AnilistPlatformRef.Get().RefreshAnimeCollection(context.Background())

	if err != nil {
		return nil, err
	}

	// Save the collection to PlaybackManager
	a.PlaybackManager.SetAnimeCollection(ret)

	// Save the collection to AutoDownloader
	a.AutoDownloader.SetAnimeCollection(ret)

	// Save the collection to LocalManager
	a.LocalManager.SetAnimeCollection(ret)

	// Save the collection to DirectStreamManager
	a.DirectStreamManager.SetAnimeCollection(ret)

	// Save the collection to LibraryExplorer
	a.LibraryExplorer.SetAnimeCollection(ret)

	a.AutoScanner.SetAnimeCollection(ret)

	//a.SyncAnilistToSimulatedCollection()

	a.WSEventManager.SendEvent(events.RefreshedAnilistAnimeCollection, nil)

	return ret, nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// GetMangaCollection is the same as GetAnimeCollection but for manga
func (a *App) GetMangaCollection(bypassCache bool) (*media.MangaCollection, error) {
	return a.AnilistPlatformRef.Get().GetMangaCollection(context.Background(), bypassCache)
}

// GetRawMangaCollection does not exclude custom lists
func (a *App) GetRawMangaCollection(bypassCache bool) (*media.MangaCollection, error) {
	return a.AnilistPlatformRef.Get().GetRawMangaCollection(context.Background(), bypassCache)
}

// RefreshMangaCollection queries Anilist for the user's manga collection
func (a *App) RefreshMangaCollection() (*media.MangaCollection, error) {
	mc, err := a.AnilistPlatformRef.Get().RefreshMangaCollection(context.Background())

	if err != nil {
		return nil, err
	}

	a.LocalManager.SetMangaCollection(mc)

	a.WSEventManager.SendEvent(events.RefreshedAnilistMangaCollection, nil)

	return mc, nil
}
