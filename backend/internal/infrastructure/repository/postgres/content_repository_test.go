package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"search_engine/internal/domain/entities"
	"search_engine/internal/domain/repositories"
)

func getTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL not set; skipping repo tests")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("pgxpool.New: %v", err)
	}
	return pool
}

func TestContentRepository_CRUD(t *testing.T) {
	pool := getTestPool(t)
	defer pool.Close()
	ctx := context.Background()

	repo := NewContentRepository(pool)

	// Create
	title := "Test Title"
	desc := "Test Description"
	now := time.Now().Add(-time.Hour)
	c := &entities.Content{
		ProviderID:        "test",
		ProviderContentID: "c-001",
		Title:             title,
		ContentType:       entities.ContentTypeText,
		Description:       &desc,
		PublishedAt:       &now,
	}
	if err := repo.Create(ctx, c); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if c.ID == 0 {
		t.Fatalf("expected ID set after create")
	}

	// GetByID
	got, err := repo.GetByID(ctx, c.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.Title != title {
		t.Fatalf("expected title %q got %q", title, got.Title)
	}

	// Update
	newTitle := "Updated Title"
	c.Title = newTitle
	if err := repo.Update(ctx, c); err != nil {
		t.Fatalf("Update: %v", err)
	}
	got2, err := repo.GetByID(ctx, c.ID)
	if err != nil {
		t.Fatalf("GetByID after update: %v", err)
	}
	if got2.Title != newTitle {
		t.Fatalf("expected updated title %q got %q", newTitle, got2.Title)
	}

	// GetAll with pagination
	items, total, err := repo.GetAll(ctx, repositories.ContentFilters{}, repositories.Pagination{Page: 1, PageSize: 10}, repositories.SortByRelevance)
	if err != nil {
		t.Fatalf("GetAll: %v", err)
	}
	if total == 0 || len(items) == 0 {
		t.Fatalf("expected some items in list")
	}

	// SearchByKeyword
	kwItems, kwTotal, err := repo.SearchByKeyword(ctx, "Title", repositories.ContentFilters{}, repositories.Pagination{Page: 1, PageSize: 5}, repositories.SortByPopularity)
	if err != nil {
		t.Fatalf("SearchByKeyword: %v", err)
	}
	if kwTotal == 0 || len(kwItems) == 0 {
		t.Fatalf("expected search results")
	}
}


