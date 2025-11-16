package repositories

import (
	"context"

	"search_engine/internal/domain/entities"
)

type ContentMetricsRepository interface {
	Create(ctx context.Context, m *entities.ContentMetrics) error
	UpdateByContentID(ctx context.Context, contentID int64, m *entities.ContentMetrics) error
	GetByContentID(ctx context.Context, contentID int64) (*entities.ContentMetrics, error)
	BulkUpsert(ctx context.Context, metrics []entities.ContentMetrics) error
}
