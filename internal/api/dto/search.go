package dto

import (
	"time"
	"strings"
)

type SearchRequest struct {
	Keyword    string
	ContentType string // "video" | "text" | ""
	SortBy     string  // "score_desc" | "score_asc" | "date_desc" | "date_asc"
	Page       int
	PageSize   int
}

func (r *SearchRequest) Normalize(defaultPage, defaultPageSize, maxPageSize int) {
	r.Keyword = strings.TrimSpace(r.Keyword)
	if r.Page <= 0 {
		r.Page = defaultPage
	}
	if r.PageSize <= 0 || r.PageSize > maxPageSize {
		r.PageSize = defaultPageSize
	}
	if r.SortBy == "" {
		r.SortBy = "score_desc"
	}
}

type ErrorDTO struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

type PaginationDTO struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

type ContentSummaryDTO struct {
	ID           int64      `json:"id"`
	Title        string     `json:"title"`
	ContentType  string     `json:"content_type"`
	Description  *string    `json:"description,omitempty"`
	URL          *string    `json:"url,omitempty"`
	ThumbnailURL *string    `json:"thumbnail_url,omitempty"`
	Score        float64    `json:"score"`
	PublishedAt  *time.Time `json:"published_at,omitempty"`
	Provider     string     `json:"provider"`
}

type MetricsDTO struct {
	Views         *int64     `json:"views,omitempty"`
	Likes         *int64     `json:"likes,omitempty"`
	ReadingTime   *int       `json:"reading_time,omitempty"`
	Reactions     *int       `json:"reactions,omitempty"`
	RecalculatedAt *time.Time `json:"recalculated_at,omitempty"`
}

type ContentDetailDTO struct {
	ContentSummaryDTO
	Metrics MetricsDTO `json:"metrics"`
}

type SearchResponse struct {
	Success    bool               `json:"success"`
	Data       []ContentSummaryDTO`json:"data"`
	Pagination PaginationDTO      `json:"pagination"`
	Error      *ErrorDTO          `json:"error,omitempty"`
}

type ContentResponse struct {
	Success bool             `json:"success"`
	Data    *ContentDetailDTO`json:"data,omitempty"`
	Error   *ErrorDTO        `json:"error,omitempty"`
}

type StatsProviderDTO struct {
	ProviderID  string    `json:"provider_id"`
	ContentCount int64    `json:"content_count"`
	LastSync    *time.Time`json:"last_sync,omitempty"`
}

type StatsDTO struct {
	TotalContents int64             `json:"total_contents"`
	TotalVideos   int64             `json:"total_videos"`
	TotalTexts    int64             `json:"total_texts"`
	AverageScore  float64           `json:"average_score"`
	LastSync      *time.Time        `json:"last_sync,omitempty"`
	Providers     []StatsProviderDTO`json:"providers"`
}

type StatsResponse struct {
	Success bool     `json:"success"`
	Data    StatsDTO `json:"data"`
}


