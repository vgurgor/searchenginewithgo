//go:build integration

package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"search_engine/internal/api/handlers"
	"search_engine/internal/api/dto"
	"search_engine/internal/domain/entities"
	"search_engine/internal/domain/repositories"
	"search_engine/internal/infrastructure/services"
	"context"
)

// stubContentRepo implements repositories.ContentRepository with minimal behavior for search.
type stubContentRepo struct{}

func (s *stubContentRepo) Create(_ context.Context, _ *entities.Content) error { return nil }
func (s *stubContentRepo) GetByID(_ context.Context, _ int64) (*entities.Content, error) { return nil, nil }
func (s *stubContentRepo) GetAll(_ context.Context, _ repositories.ContentFilters, _ repositories.Pagination, _ repositories.SortBy) ([]entities.Content, int64, error) {
	return nil, 0, nil
}
func (s *stubContentRepo) Update(_ context.Context, _ *entities.Content) error { return nil }
func (s *stubContentRepo) BulkInsert(_ context.Context, _ []entities.Content) error { return nil }
func (s *stubContentRepo) SearchByKeyword(_ context.Context, _ string, _ repositories.ContentFilters, _ repositories.Pagination, _ repositories.SortBy) ([]entities.Content, int64, error) {
	return nil, 0, nil
}
func (s *stubContentRepo) GetByProviderKey(_ context.Context, _, _ string) (*entities.Content, error) { return nil, nil }
func (s *stubContentRepo) ListIDs(_ context.Context, _, _ int) ([]int64, error) { return nil, nil }
func (s *stubContentRepo) CountAll(_ context.Context) (int64, error) { return 0, nil }
func (s *stubContentRepo) SearchWithFilters(_ context.Context, keyword string, contentType *entities.ContentType, pagination repositories.Pagination, sort repositories.SearchSort) ([]repositories.ContentWithMetrics, int64, error) {
	// Return one predictable item
	now := time.Now().UTC()
	ct := entities.ContentTypeVideo
	if contentType != nil {
		ct = *contentType
	}
	item := repositories.ContentWithMetrics{
		Content: entities.Content{
			ID:          1,
			ProviderID:  "provider1",
			Title:       "Hello " + keyword,
			ContentType: ct,
			PublishedAt: &now,
		},
		Metrics: entities.ContentMetrics{
			FinalScore: 123.45,
		},
	}
	return []repositories.ContentWithMetrics{item}, 1, nil
}
func (s *stubContentRepo) GetDetailByID(_ context.Context, _ int64) (*repositories.ContentWithMetrics, error) { return nil, nil }
func (s *stubContentRepo) CountByType(_ context.Context) (map[entities.ContentType]int64, error) { return map[entities.ContentType]int64{}, nil }
func (s *stubContentRepo) GetAverageScore(_ context.Context) (float64, error) { return 0, nil }
func (s *stubContentRepo) CountByProvider(_ context.Context) (map[string]int64, error) { return map[string]int64{}, nil }
func (s *stubContentRepo) SoftDelete(_ context.Context, _ int64) error { return nil }
func (s *stubContentRepo) ListIDsByType(_ context.Context, _ entities.ContentType, _, _ int) ([]int64, error) { return nil, nil }
func (s *stubContentRepo) GetAverageScoreByProvider(_ context.Context, _ string) (float64, error) { return 0, nil }

func TestContentsSearchEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	svc := &services.ContentSearchService{
		Repo:            &stubContentRepo{},
		DefaultPageSize: 20,
		MaxPageSize:     100,
		CacheEnabled:    false,
	}
	handlers.RegisterContentRoutes(router, svc, 20, 100)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/contents/search?q=test&type=video&sort=score_desc&page=1&page_size=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body=%s", w.Code, w.Body.String())
	}
	var resp dto.SearchResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if !resp.Success {
		t.Fatalf("expected success=true")
	}
	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 item, got %d", len(resp.Data))
	}
	if resp.Data[0].Score <= 0 {
		t.Fatalf("expected positive score, got %v", resp.Data[0].Score)
	}
}


