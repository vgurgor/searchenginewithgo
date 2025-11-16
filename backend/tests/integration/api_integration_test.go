//go:build integration

package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"search_engine/internal/api/dto"
	"search_engine/internal/api/handlers"
	"search_engine/internal/infrastructure/repository/postgres"
	"search_engine/internal/infrastructure/services"
	"search_engine/tests"
)

func TestHealthEndpoint(t *testing.T) {
	// Setup test environment
	env := tests.SetupTestEnvironment(t)
	defer env.Cleanup()

	// Create health handler
	healthHandler := handlers.NewHealthHandler(env.DB, env.Redis, time.Now(), "1.0.0-test", env.Logger)

	// Create router
	router := gin.New()
	router.GET("/health", healthHandler)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	// Serve request
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "healthy", response["status"])

	// Check services health
	services, ok := response["services"].(map[string]interface{})
	require.True(t, ok, "services should be a map")

	postgres, ok := services["postgres"].(map[string]interface{})
	require.True(t, ok, "postgres should be a map")
	assert.True(t, postgres["healthy"].(bool), "postgres should be healthy")

	redis, ok := services["redis"].(map[string]interface{})
	require.True(t, ok, "redis should be a map")
	assert.True(t, redis["healthy"].(bool), "redis should be healthy")

	assert.Contains(t, response, "timestamp")
	assert.Contains(t, response, "uptime_seconds")
	assert.Equal(t, "1.0.0-test", response["version"])
}

func TestContentsSearchIntegration(t *testing.T) {
	// Setup test environment
	env := tests.SetupTestEnvironment(t)
	defer env.Cleanup()

	ctx := context.Background()

	// Insert test data
	testContents := []struct {
		title       string
		contentType string
		score       float64
		provider    string
	}{
		{"Go Programming Tutorial", "video", 95.5, "provider1"},
		{"Python Basics", "text", 85.2, "provider2"},
		{"Advanced Go Patterns", "video", 88.7, "provider1"},
		{"Web Development Guide", "text", 78.3, "provider2"},
	}

	for i, tc := range testContents {
		contentID := int64(i + 1)
		contentType := "video"
		if tc.contentType == "text" {
			contentType = "text"
		}

		_, err := env.DB.Exec(ctx, `
			INSERT INTO contents (id, provider_id, provider_content_id, title, content_type, published_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, contentID, tc.provider, fmt.Sprintf("test-%d", contentID), tc.title, contentType, time.Now())
		require.NoError(t, err)

		_, err = env.DB.Exec(ctx, `
			INSERT INTO content_metrics (content_id, final_score)
			VALUES ($1, $2)
		`, contentID, tc.score)
		require.NoError(t, err)
	}

	// Create services
	repo := postgres.NewContentRepository(env.DB)
	searchSvc := &services.ContentSearchService{
		Repo:            repo,
		DefaultPageSize: 20,
		MaxPageSize:     100,
		CacheClient:     env.Redis,
		CacheEnabled:    false,
	}

	// Create router
	router := gin.New()
	handlers.RegisterContentRoutes(router, searchSvc, 20, 100)

	// Test search without filters
	t.Run("search_all", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/contents/search?page_size=10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.SearchResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response.Success)
		assert.Len(t, response.Data, 4)
		assert.Equal(t, int64(4), response.Pagination.TotalItems)
	})

	// Test search with keyword
	t.Run("search_with_keyword", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/contents/search?q=Go&page_size=10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var response dto.SearchResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Debug: Print response if not successful
		if !response.Success {
			t.Logf("Search failed - Code: %d, Response: %+v", w.Code, response)
		}

		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, response.Success)
		assert.Len(t, response.Data, 2) // Should find "Go Programming Tutorial" and "Advanced Go Patterns"
	})

	// Test search with content type filter
	t.Run("search_with_type_filter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/contents/search?type=video&page_size=10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.SearchResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response.Success)
		assert.Len(t, response.Data, 2) // Should find 2 videos
		for _, item := range response.Data {
			assert.Equal(t, "video", item.ContentType)
		}
	})

	// Test pagination
	t.Run("search_pagination", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/contents/search?page=1&page_size=2", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.SearchResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response.Success)
		assert.Len(t, response.Data, 2)
		assert.Equal(t, 1, response.Pagination.Page)
		assert.Equal(t, 2, response.Pagination.PageSize)
		assert.Equal(t, int64(4), response.Pagination.TotalItems)
		assert.Equal(t, 2, response.Pagination.TotalPages)
	})
}

func TestContentsStatsIntegration(t *testing.T) {
	// Setup test environment
	env := tests.SetupTestEnvironment(t)
	defer env.Cleanup()

	ctx := context.Background()

	// Insert test data
	_, err := env.DB.Exec(ctx, `
		INSERT INTO contents (id, provider_id, provider_content_id, title, content_type, published_at)
		VALUES
			(1, 'provider1', 'test-1', 'Video 1', 'video', $1),
			(2, 'provider1', 'test-2', 'Video 2', 'video', $1),
			(3, 'provider2', 'test-3', 'Text 1', 'text', $1),
			(4, 'provider2', 'test-4', 'Text 2', 'text', $1),
			(5, 'provider2', 'test-5', 'Text 3', 'text', $1)
	`, time.Now())
	require.NoError(t, err)

	_, err = env.DB.Exec(ctx, `
		INSERT INTO content_metrics (content_id, final_score)
		VALUES
			(1, 90.0), (2, 85.0), (3, 80.0), (4, 75.0), (5, 70.0)
	`)
	require.NoError(t, err)

	// Create services
	repo := postgres.NewContentRepository(env.DB)
	searchSvc := &services.ContentSearchService{
		Repo:            repo,
		DefaultPageSize: 20,
		MaxPageSize:     100,
		CacheClient:     env.Redis,
		CacheEnabled:    false,
	}

	// Create router
	router := gin.New()
	handlers.RegisterContentRoutes(router, searchSvc, 20, 100)

	// Test stats endpoint
	req := httptest.NewRequest(http.MethodGet, "/api/v1/contents/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.StatsResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.True(t, response.Success)
	assert.Equal(t, int64(5), response.Data.TotalContents)
	assert.Equal(t, int64(2), response.Data.TotalVideos)
	assert.Equal(t, int64(3), response.Data.TotalTexts)
	assert.Greater(t, response.Data.AverageScore, 0.0)
	assert.Len(t, response.Data.Providers, 2)
}

func TestContentsDetailIntegration(t *testing.T) {
	// Setup test environment
	env := tests.SetupTestEnvironment(t)
	defer env.Cleanup()

	ctx := context.Background()

	// Insert test content
	now := time.Now()
	_, err := env.DB.Exec(ctx, `
		INSERT INTO contents (id, provider_id, provider_content_id, title, content_type, description, url, published_at)
		VALUES (1, 'provider1', 'test-1', 'Test Video', 'video', 'Test description', 'https://example.com/video', $1)
	`, now)
	require.NoError(t, err)

	_, err = env.DB.Exec(ctx, `
		INSERT INTO content_metrics (content_id, views, likes, reading_time, reactions, final_score)
		VALUES (1, 1000, 50, 0, 0, 85.5)
	`)
	require.NoError(t, err)

	// Create services
	repo := postgres.NewContentRepository(env.DB)
	searchSvc := &services.ContentSearchService{
		Repo:            repo,
		DefaultPageSize: 20,
		MaxPageSize:     100,
		CacheClient:     env.Redis,
		CacheEnabled:    false,
	}

	// Create router
	router := gin.New()
	handlers.RegisterContentRoutes(router, searchSvc, 20, 100)

	// Test detail endpoint
	req := httptest.NewRequest(http.MethodGet, "/api/v1/contents/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.APIContentResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.True(t, response.Success)
	assert.NotNil(t, response.Data)
	assert.Equal(t, int64(1), response.Data.ID)
	assert.Equal(t, "Test Video", response.Data.Title)
	assert.Equal(t, "video", response.Data.ContentType)
	assert.Equal(t, "Test description", *response.Data.Description)
	assert.Equal(t, "https://example.com/video", *response.Data.URL)
	assert.Equal(t, int64(1000), *response.Data.Metrics.Views)
	assert.Equal(t, int64(50), *response.Data.Metrics.Likes)
	assert.Equal(t, 85.5, response.Data.Score)
}

func TestErrorHandlingIntegration(t *testing.T) {
	// Setup test environment
	env := tests.SetupTestEnvironment(t)
	defer env.Cleanup()

	// Create services
	repo := postgres.NewContentRepository(env.DB)
	searchSvc := &services.ContentSearchService{
		Repo:            repo,
		DefaultPageSize: 20,
		MaxPageSize:     100,
		CacheClient:     env.Redis,
		CacheEnabled:    false,
	}

	// Create router
	router := gin.New()
	handlers.RegisterContentRoutes(router, searchSvc, 20, 100)

	// Test invalid page size
	t.Run("invalid_page_size", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/contents/search?page_size=200", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var response dto.SearchResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Debug: Print response
		if w.Code != http.StatusBadRequest {
			t.Logf("Expected 400, got %d - Response: %+v", w.Code, response)
		}

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.False(t, response.Success)
		assert.NotNil(t, response.Error)
		assert.Equal(t, "INVALID_PARAMETER", response.Error.Code)
	})

	// Test content not found
	t.Run("content_not_found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/contents/999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response dto.APIContentResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.False(t, response.Success)
		assert.NotNil(t, response.Error)
		assert.Equal(t, "CONTENT_NOT_FOUND", response.Error.Code)
	})
}
