package entities

import "time"

type ContentType string

const (
	ContentTypeVideo ContentType = "video"
	ContentTypeText  ContentType = "text"
)

type Content struct {
	ID                int64       `json:"id"`
	ProviderID        string      `json:"providerId"`
	ProviderContentID string      `json:"providerContentId"`
	Title             string      `json:"title"`
	ContentType       ContentType `json:"contentType"`
	Description       *string     `json:"description,omitempty"`
	URL               *string     `json:"url,omitempty"`
	ThumbnailURL      *string     `json:"thumbnailUrl,omitempty"`
	PublishedAt       *time.Time  `json:"publishedAt,omitempty"`
	CreatedAt         time.Time   `json:"createdAt"`
	UpdatedAt         time.Time   `json:"updatedAt"`

	// Relation
	Metrics *ContentMetrics `json:"metrics,omitempty"`
}
