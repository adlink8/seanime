package hibikecustomsource

import (
	"context"
	"seanime/internal/api/metadata"
	"seanime/internal/media"
)

// Custom sources allow users to add custom media, often things not available on AniList, to the app.
// At runtime the extension will be assigned a unique extension identifier (int). Extension IDs use bit-based separation:
// AniList IDs: 0 to 2^31-1, Extension IDs: 2^31+ with embedded extension identifier and local ID.
// A custom source can be identified by its ID and SiteUrl property (e.g. "ext_custom_source_{extId}|END|https://example.com")
//
// Custom source media are fetched in the Platform and Metadata provider. The customsource.Manager is tasked with storing tracking info.

type (
	Settings struct {
		SupportsAnime bool `json:"supportsAnime"`
		SupportsManga bool `json:"supportsManga"`
	}

	ListAnimeResponse struct {
		Media      []*media.Anime `json:"media"`
		Page       int            `json:"page"`
		TotalPages int            `json:"totalPages"`
		Total      int            `json:"total"`
	}

	ListMangaResponse struct {
		Media      []*media.Manga `json:"media"`
		Page       int            `json:"page"`
		TotalPages int            `json:"totalPages"`
		Total      int            `json:"total"`
	}

	Provider interface {
		GetExtensionIdentifier() int
		GetSettings() Settings
		GetAnime(ctx context.Context, id []int) ([]*media.Anime, error)
		ListAnime(ctx context.Context, search string, page int, perPage int) (*ListAnimeResponse, error)
		GetAnimeWithRelations(ctx context.Context, id int) (*media.CompleteAnime, error)
		GetAnimeMetadata(ctx context.Context, id int) (*metadata.AnimeMetadata, error)
		GetAnimeDetails(ctx context.Context, id int) (*media.AnimeDetails, error)
		GetManga(ctx context.Context, id []int) ([]*media.Manga, error)
		ListManga(ctx context.Context, search string, page int, perPage int) (*ListMangaResponse, error)
		GetMangaDetails(ctx context.Context, id int) (*media.MangaDetails, error)
	}
)
