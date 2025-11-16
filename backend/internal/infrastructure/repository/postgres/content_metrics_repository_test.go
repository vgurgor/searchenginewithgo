package postgres

import (
	"context"
	"testing"
	"time"

	"search_engine/internal/domain/entities"
)

func TestContentMetricsRepository_CRUD(t *testing.T) {
	pool := getTestPool(t)
	defer pool.Close()
	ctx := context.Background()

	// Ensure a content exists
	cRepo := NewContentRepository(pool)
	title := "Metric Content"
	content := &entities.Content{
		ProviderID:        "test",
		ProviderContentID: "m-001",
		Title:             title,
		ContentType:       entities.ContentTypeVideo,
	}
	if err := cRepo.Create(ctx, content); err != nil {
		t.Fatalf("create content: %v", err)
	}

	mRepo := NewContentMetricsRepository(pool)
	now := time.Now()
	m := &entities.ContentMetrics{
		ContentID:      content.ID,
		Views:          123,
		Likes:          45,
		ReadingTime:    0,
		Reactions:      0,
		FinalScore:     78.9,
		RecalculatedAt: &now,
	}
	if err := mRepo.Create(ctx, m); err != nil {
		t.Fatalf("create metrics: %v", err)
	}
	got, err := mRepo.GetByContentID(ctx, content.ID)
	if err != nil {
		t.Fatalf("get metrics: %v", err)
	}
	if got.FinalScore != 78.9 {
		t.Fatalf("expected score 78.9 got %v", got.FinalScore)
	}

	m.Views = 999
	if err := mRepo.UpdateByContentID(ctx, content.ID, m); err != nil {
		t.Fatalf("update metrics: %v", err)
	}
	got2, err := mRepo.GetByContentID(ctx, content.ID)
	if err != nil {
		t.Fatalf("get metrics after update: %v", err)
	}
	if got2.Views != 999 {
		t.Fatalf("expected updated views 999 got %d", got2.Views)
	}
}
