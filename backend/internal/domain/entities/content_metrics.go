package entities

import "time"

type ContentMetrics struct {
	ID             int64      `json:"id"`
	ContentID      int64      `json:"contentId"`
	Views          int64      `json:"views"`
	Likes          int64      `json:"likes"`
	ReadingTime    int        `json:"readingTime"`
	Reactions      int        `json:"reactions"`
	FinalScore     float64    `json:"finalScore"`
	RecalculatedAt *time.Time `json:"recalculatedAt,omitempty"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}
