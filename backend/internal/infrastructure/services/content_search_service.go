package services

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"search_engine/internal/api/dto"
	"search_engine/internal/domain/entities"
	"search_engine/internal/domain/repositories"
	"search_engine/internal/infrastructure/cache"

	"github.com/redis/go-redis/v9"
)

type ContentSearchService struct {
	Repo         repositories.ContentRepository
	HistoryRepo  repositories.SyncHistoryRepository
	DefaultPageSize int
	MaxPageSize     int
	// Optional cache
	CacheClient   *redis.Client
	CacheEnabled  bool
	CacheTTL      time.Duration
}

func (s *ContentSearchService) SearchContents(ctx context.Context, req dto.SearchRequest) ([]dto.ContentSummaryDTO, int64, error) {
	var ct *entities.ContentType
	if req.ContentType != "" {
		t := strings.ToLower(req.ContentType)
		switch t {
		case "video":
			v := entities.ContentTypeVideo
			ct = &v
		case "text":
			v := entities.ContentTypeText
			ct = &v
		default:
			return nil, 0, errors.New("invalid content type")
		}
	}
	var sort repositories.SearchSort = repositories.SearchSortScoreDesc
	switch req.SortBy {
	case "score_desc":
		sort = repositories.SearchSortScoreDesc
	case "score_asc":
		sort = repositories.SearchSortScoreAsc
	case "date_desc":
		sort = repositories.SearchSortDateDesc
	case "date_asc":
		sort = repositories.SearchSortDateAsc
	default:
		// keep default
	}

	// Cache key
	var cached struct {
		Items []dto.ContentSummaryDTO `json:"items"`
		Total int64                   `json:"total"`
	}
	cacheKey := fmt.Sprintf("sc:%s|%s|%s|%d|%d",
		strings.ToLower(strings.TrimSpace(req.Keyword)),
		strings.ToLower(strings.TrimSpace(req.ContentType)),
		string(sort),
		req.Page,
		req.PageSize,
	)
	if s.CacheEnabled && s.CacheClient != nil && s.CacheTTL > 0 {
		if ok, _ := cache.GetJSON(ctx, s.CacheClient, cacheKey, &cached); ok {
			return cached.Items, cached.Total, nil
		}
	}

	items, total, err := s.Repo.SearchWithFilters(ctx, req.Keyword, ct, repositories.Pagination{Page: req.Page, PageSize: req.PageSize}, sort)
	if err != nil {
		return nil, 0, err
	}
	out := make([]dto.ContentSummaryDTO, 0, len(items))
	for _, row := range items {
		desc := truncateOrNil(row.Content.Description, 200)
		score := row.Metrics.FinalScore
		out = append(out, dto.ContentSummaryDTO{
			ID:          row.Content.ID,
			Title:       row.Content.Title,
			ContentType: string(row.Content.ContentType),
			Description: desc,
			URL:         row.Content.URL,
			ThumbnailURL: row.Content.ThumbnailURL,
			Score:       score,
			PublishedAt: row.Content.PublishedAt,
			Provider:    row.Content.ProviderID,
		})
	}
	if s.CacheEnabled && s.CacheClient != nil && s.CacheTTL > 0 {
		_ = cache.SetJSON(ctx, s.CacheClient, cacheKey, struct {
			Items []dto.ContentSummaryDTO `json:"items"`
			Total int64                   `json:"total"`
		}{Items: out, Total: total}, s.CacheTTL)
	}
	return out, total, nil
}

func (s *ContentSearchService) GetContentByID(ctx context.Context, id int64) (*dto.ContentDetailDTO, error) {
	cacheKey := fmt.Sprintf("content:detail:%d", id)

	// Try cache first if enabled
	if s.CacheEnabled && s.CacheClient != nil {
		var cached dto.ContentDetailDTO
		if err := cache.GetJSON(s.CacheClient, cacheKey, &cached); err == nil {
			return &cached, nil
		}
	}

	row, err := s.Repo.GetDetailByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, nil
	}
	var viewsPtr, likesPtr *int64
	var rtPtr, reacPtr *int
	if row.Metrics.Views != 0 {
		v := row.Metrics.Views
		viewsPtr = &v
	}
	if row.Metrics.Likes != 0 {
		v := row.Metrics.Likes
		likesPtr = &v
	}
	if row.Metrics.ReadingTime != 0 {
		v := row.Metrics.ReadingTime
		rtPtr = &v
	}
	if row.Metrics.Reactions != 0 {
		v := row.Metrics.Reactions
		reacPtr = &v
	}
	desc := row.Content.Description
	result := &dto.ContentDetailDTO{
		ContentSummaryDTO: dto.ContentSummaryDTO{
			ID:          row.Content.ID,
			Title:       row.Content.Title,
			ContentType: string(row.Content.ContentType),
			Description: desc,
			URL:         row.Content.URL,
			ThumbnailURL: row.Content.ThumbnailURL,
			Score:       row.Metrics.FinalScore,
			PublishedAt: row.Content.PublishedAt,
			Provider:    row.Content.ProviderID,
		},
		Metrics: dto.MetricsDTO{
			Views: viewsPtr,
			Likes: likesPtr,
			ReadingTime: rtPtr,
			Reactions: reacPtr,
			RecalculatedAt: row.Metrics.RecalculatedAt,
		},
	}

	// Cache the result if enabled
	if s.CacheEnabled && s.CacheClient != nil {
		cache.SetJSON(s.CacheClient, cacheKey, result, s.CacheTTL)
	}

	return result, nil
}

// InvalidateContentCache clears the cache for a specific content item
func (s *ContentSearchService) InvalidateContentCache(id int64) error {
	if !s.CacheEnabled || s.CacheClient == nil {
		return nil
	}

	cacheKey := fmt.Sprintf("content:detail:%d", id)
	return s.CacheClient.Del(ctx, cacheKey).Err()
}

// InvalidateSearchCache clears search result caches (pattern-based invalidation)
func (s *ContentSearchService) InvalidateSearchCache() error {
	if !s.CacheEnabled || s.CacheClient == nil {
		return nil
	}

	// Note: Redis doesn't support pattern deletion efficiently in production
	// This is a simplified implementation for demo purposes
	// In production, consider using cache tags or separate cache keys with TTL

	// For now, we'll just clear some common search patterns
	patterns := []string{
		"search:*",
		"stats:*",
	}

	for _, pattern := range patterns {
		keys, err := s.CacheClient.Keys(ctx, pattern).Result()
		if err != nil {
			continue
		}
		if len(keys) > 0 {
			s.CacheClient.Del(ctx, keys...)
		}
	}

	return nil
}

func (s *ContentSearchService) GetStats(ctx context.Context) (dto.StatsDTO, error) {
	total, err := s.Repo.CountAll(ctx)
	if err != nil {
		return dto.StatsDTO{}, err
	}
	byType, err := s.Repo.CountByType(ctx)
	if err != nil {
		return dto.StatsDTO{}, err
	}
	avg, err := s.Repo.GetAverageScore(ctx)
	if err != nil {
		return dto.StatsDTO{}, err
	}
	byProvider, err := s.Repo.CountByProvider(ctx)
	if err != nil {
		return dto.StatsDTO{}, err
	}
	providers := make([]dto.StatsProviderDTO, 0, len(byProvider))
	var lastSync *int64
	for pid, cnt := range byProvider {
		entry := dto.StatsProviderDTO{ProviderID: pid, ContentCount: cnt}
		if s.HistoryRepo != nil {
			if h, err := s.HistoryRepo.GetLastSync(ctx, pid); err == nil && h != nil && h.CompletedAt != nil {
				entry.LastSync = h.CompletedAt
				if lastSync == nil {
					v := h.CompletedAt.Unix()
					lastSync = &v
				} else if h.CompletedAt.Unix() > *lastSync {
					v := h.CompletedAt.Unix()
					lastSync = &v
				}
			}
		}
		providers = append(providers, entry)
	}
	var lastSyncTime *int64 = lastSync
	var lastSyncPtr *time.Time
	if lastSyncTime != nil {
		t := time.Unix(*lastSyncTime, 0).UTC()
		lastSyncPtr = &t
	}
	return dto.StatsDTO{
		TotalContents: total,
		TotalVideos:   byType[entities.ContentTypeVideo],
		TotalTexts:    byType[entities.ContentTypeText],
		AverageScore:  round2(avg),
		LastSync:      lastSyncPtr,
		Providers:     providers,
	}, nil
}

func truncateOrNil(s *string, max int) *string {
	if s == nil {
		return nil
	}
	r := []rune(*s)
	if len(r) <= max {
		return s
	}
	tr := string(r[:max])
	return &tr
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

