package repositories

import (
	"context"
	"search_engine/internal/domain/entities"
)

type SyncHistoryRepository interface {
	Create(ctx context.Context, h *entities.SyncHistory) error
	GetByProviderID(ctx context.Context, providerID string, limit int) ([]entities.SyncHistory, error)
	GetLastSync(ctx context.Context, providerID string) (*entities.SyncHistory, error)
	GetAll(ctx context.Context, limit int) ([]entities.SyncHistory, error)
}


