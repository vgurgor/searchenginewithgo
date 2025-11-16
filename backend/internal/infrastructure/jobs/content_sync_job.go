package jobs

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
	"search_engine/internal/infrastructure/services"
)

type ContentSyncJob struct {
	Logger     *zap.Logger
	Service    *services.ContentSyncService
	Interval   time.Duration
	Enabled    bool
	MaxRetries int
	RetryDelay time.Duration

	mu      sync.Mutex
	running bool
	stopCh  chan struct{}
}

func NewContentSyncJob(logger *zap.Logger, svc *services.ContentSyncService, interval time.Duration, enabled bool, retries int, retryDelay time.Duration) *ContentSyncJob {
	return &ContentSyncJob{
		Logger:     logger,
		Service:    svc,
		Interval:   interval,
		Enabled:    enabled,
		MaxRetries: retries,
		RetryDelay: retryDelay,
		stopCh:     make(chan struct{}),
	}
}

func (j *ContentSyncJob) Start() {
	if !j.Enabled {
		return
	}
	ticker := time.NewTicker(j.Interval)
	go func() {
		j.Logger.Info("content sync job started", zap.Duration("interval", j.Interval))
		// Run immediately once on start to ensure data is ingested without waiting first tick
		j.runOnce()
		defer j.Logger.Info("content sync job stopped")
		for {
			select {
			case <-ticker.C:
				j.runOnce()
			case <-j.stopCh:
				ticker.Stop()
				return
			}
		}
	}()
}

func (j *ContentSyncJob) Stop() {
	close(j.stopCh)
}

func (j *ContentSyncJob) runOnce() {
	j.mu.Lock()
	if j.running {
		j.Logger.Warn("content sync already running; skipping")
		j.mu.Unlock()
		return
	}
	j.running = true
	j.mu.Unlock()
	defer func() {
		j.mu.Lock()
		j.running = false
		j.mu.Unlock()
	}()

	var lastErr error
	for attempt := 0; attempt <= j.MaxRetries; attempt++ {
		_, err := j.Service.SyncAllProviders(context.Background())
		if err == nil {
			return
		}
		lastErr = err
		j.Logger.Warn("content sync attempt failed", zap.Int("attempt", attempt+1), zap.Error(err))
		time.Sleep(j.RetryDelay)
	}
	if lastErr != nil {
		j.Logger.Error("content sync failed after retries", zap.Error(lastErr))
	}
}
