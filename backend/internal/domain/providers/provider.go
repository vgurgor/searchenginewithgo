package providers

import "time"

type RateLimit struct {
	RequestsPerMinute int
}

type ProviderContent struct {
	ProviderID        string
	ProviderContentID string
	Title             string
	ContentType       string // "video" | "text"
	Description       string
	URL               string
	ThumbnailURL      string
	Views             *int64
	Likes             *int64
	ReadingTime       *int
	Reactions         *int
	PublishedAt       time.Time
}

type IContentProvider interface {
	FetchContents() ([]ProviderContent, error)
	GetProviderID() string
	GetRateLimit() RateLimit
}
