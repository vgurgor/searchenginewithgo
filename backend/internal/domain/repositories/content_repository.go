package repositories

import (
	"context"

	"search_engine/internal/domain/entities"
)

type ContentFilters struct {
	ContentType *entities.ContentType
	ProviderID  *string
}

type Pagination struct {
	Page     int
	PageSize int
}

type SortBy string

const (
	SortByPopularity SortBy = "popularity"
	SortByRelevance  SortBy = "relevance"
)

type ContentRepository interface {
	Create(ctx context.Context, c *entities.Content) error
	GetByID(ctx context.Context, id int64) (*entities.Content, error)
	GetAll(ctx context.Context, filters ContentFilters, pagination Pagination, sort SortBy) ([]entities.Content, int64, error)
	Update(ctx context.Context, c *entities.Content) error
	BulkInsert(ctx context.Context, contents []entities.Content) error
	SearchByKeyword(ctx context.Context, keyword string, filters ContentFilters, pagination Pagination, sort SortBy) ([]entities.Content, int64, error)
	// Extended
	GetByProviderKey(ctx context.Context, providerID, providerContentID string) (*entities.Content, error)
	ListIDs(ctx context.Context, offset, limit int) ([]int64, error)
	CountAll(ctx context.Context) (int64, error)
	// Search & Stats
	SearchWithFilters(ctx context.Context, keyword string, contentType *entities.ContentType, pagination Pagination, sort SearchSort) ([]ContentWithMetrics, int64, error)
	GetDetailByID(ctx context.Context, id int64) (*ContentWithMetrics, error)
	CountByType(ctx context.Context) (map[entities.ContentType]int64, error)
	GetAverageScore(ctx context.Context) (float64, error)
	CountByProvider(ctx context.Context) (map[string]int64, error)
	SoftDelete(ctx context.Context, id int64) error
	ListIDsByType(ctx context.Context, t entities.ContentType, offset, limit int) ([]int64, error)
	GetAverageScoreByProvider(ctx context.Context, providerID string) (float64, error)
}

type SearchSort string

const (
	SearchSortScoreDesc SearchSort = "score_desc"
	SearchSortScoreAsc  SearchSort = "score_asc"
	SearchSortDateDesc  SearchSort = "date_desc"
	SearchSortDateAsc   SearchSort = "date_asc"
)

type ContentWithMetrics struct {
	Content entities.Content
	Metrics entities.ContentMetrics
}
