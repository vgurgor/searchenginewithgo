package entities

import "time"

type SyncStatus string

const (
	SyncStatusSuccess    SyncStatus = "success"
	SyncStatusPartial    SyncStatus = "partial"
	SyncStatusFailed     SyncStatus = "failed"
	SyncStatusInProgress SyncStatus = "in_progress"
	SyncStatusSkipped    SyncStatus = "skipped"
)

type SyncHistory struct {
	ID              int64
	ProviderID      string
	SyncStatus      SyncStatus
	TotalFetched    int
	NewContents     int
	UpdatedContents int
	SkippedContents int
	FailedContents  int
	ErrorMessage    *string
	StartedAt       time.Time
	CompletedAt     *time.Time
	DurationMs      int
}
