package dto

import (
	"time"

	"search_engine/internal/domain/entities"
)

type SearchFilters struct {
	Keyword     *string               `json:"keyword,omitempty"`
	ContentType *entities.ContentType `json:"contentType,omitempty"`
	SortBy      string                `json:"sortBy,omitempty"` // "popularity" | "relevance"
	Page        int                   `json:"page,omitempty"`
	PageSize    int                   `json:"pageSize,omitempty"`
}

type ContentResponse struct {
	ID           int64                `json:"id"`
	Title        string               `json:"title"`
	ContentType  entities.ContentType `json:"contentType"`
	Description  *string              `json:"description,omitempty"`
	URL          *string              `json:"url,omitempty"`
	ThumbnailURL *string              `json:"thumbnailUrl,omitempty"`
	Score        *float64             `json:"score,omitempty"`
	PublishedAt  *time.Time           `json:"publishedAt,omitempty"`
}

type PaginatedResponse[T any] struct {
	Data       []T   `json:"data"`
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalCount int64 `json:"totalCount"`
	TotalPages int   `json:"totalPages"`
}
