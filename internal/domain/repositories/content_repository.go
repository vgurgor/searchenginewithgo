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
}


