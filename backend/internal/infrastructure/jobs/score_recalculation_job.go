package jobs

import (
	"context"
	"time"

	"go.uber.org/zap"

	"search_engine/internal/domain/repositories"
	"search_engine/internal/infrastructure/services"
)

type ScoreRecalculationJob struct {
	Logger    *zap.Logger
	Repo      repositories.ContentRepository
	Service   *services.ScoreCalculatorService
	BatchSize int
	Interval  time.Duration
	stopCh    chan struct{}
}

func NewScoreRecalculationJob(logger *zap.Logger, repo repositories.ContentRepository, svc *services.ScoreCalculatorService, batchSize int, interval time.Duration) *ScoreRecalculationJob {
	return &ScoreRecalculationJob{
		Logger:    logger,
		Repo:      repo,
		Service:   svc,
		BatchSize: batchSize,
		Interval:  interval,
		stopCh:    make(chan struct{}),
	}
}

func (j *ScoreRecalculationJob) Start() {
	ticker := time.NewTicker(j.Interval)
	go func() {
		j.Logger.Info("score recalculation job started", zap.Duration("interval", j.Interval))
		defer j.Logger.Info("score recalculation job stopped")
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

func (j *ScoreRecalculationJob) Stop() {
	close(j.stopCh)
}

func (j *ScoreRecalculationJob) runOnce() {
	ctx := context.Background()
	total, err := j.Repo.CountAll(ctx)
	if err != nil {
		j.Logger.Error("count contents failed", zap.Error(err))
		return
	}
	start := time.Now()
	var processed int64
	for offset := 0; offset < int(total); offset += j.BatchSize {
		ids, err := j.Repo.ListIDs(ctx, offset, j.BatchSize)
		if err != nil {
			j.Logger.Error("list ids failed", zap.Error(err))
			break
		}
		for _, id := range ids {
			if _, err := j.Service.RecalculateScore(ctx, id); err != nil {
				j.Logger.Warn("recalc score failed", zap.Int64("content_id", id), zap.Error(err))
			}
			processed++
		}
	}
	j.Logger.Info("recalculation batch completed", zap.Int64("processed", processed), zap.Duration("duration", time.Since(start)))
}
