package services

import (
	"context"
	"time"

	"go.uber.org/zap"

	"search_engine/internal/domain/entities"
	domainp "search_engine/internal/domain/providers"
	"search_engine/internal/domain/repositories"
)

type SyncResult struct {
	ProviderID      string
	TotalFetched    int
	NewContents     int
	UpdatedContents int
	SkippedContents int
	FailedContents  int
	Duration        time.Duration
	Errors          []string
	SyncedAt        time.Time
}

type IContentSyncService interface {
	SyncAllProviders(ctx context.Context) ([]SyncResult, error)
	SyncProvider(ctx context.Context, providerID string) (SyncResult, error)
	GetLastSyncTime(ctx context.Context, providerID string) (time.Time, error)
	GetSyncHistory(ctx context.Context, limit int) ([]entities.SyncHistory, error)
}

type ContentSyncService struct {
	Logger  *zap.Logger
	Factory interface {
		GetAllProviders() []domainp.IContentProvider
		GetProviderByID(id string) (domainp.IContentProvider, error)
	}
	ProviderClient interface {
		FetchFromProvider(ctx context.Context, providerID string) ([]domainp.ProviderContent, error)
	}
	Contents    repositories.ContentRepository
	Metrics     repositories.ContentMetricsRepository
	ScoreCalc   *ScoreCalculatorService
	HistoryRepo repositories.SyncHistoryRepository
	Thresholds  MetricsThresholds
}

func (s *ContentSyncService) SyncAllProviders(ctx context.Context) ([]SyncResult, error) {
	providers := s.Factory.GetAllProviders()
	results := make([]SyncResult, 0, len(providers))
	for _, p := range providers {
		res, _ := s.SyncProvider(ctx, p.GetProviderID())
		results = append(results, res)
	}
	return results, nil
}

func (s *ContentSyncService) SyncProvider(ctx context.Context, providerID string) (SyncResult, error) {
	start := time.Now().UTC()
	res := SyncResult{ProviderID: providerID, SyncedAt: start}
	h := entities.SyncHistory{
		ProviderID: providerID,
		SyncStatus: entities.SyncStatusInProgress,
		StartedAt:  start,
	}
	if err := s.HistoryRepo.Create(ctx, &h); err != nil {
		s.Logger.Warn("failed to create sync history", zap.String("provider", providerID), zap.Error(err))
	}

	items, err := s.ProviderClient.FetchFromProvider(ctx, providerID)
	if err != nil {
		msg := "fetch failed: " + err.Error()
		s.Logger.Error("provider fetch failed", zap.String("provider", providerID), zap.Error(err))
		res.Errors = append(res.Errors, msg)
		res.FailedContents = 0
		res.Duration = time.Since(start)
		res.TotalFetched = 0
		h.SyncStatus = entities.SyncStatusFailed
		now := time.Now().UTC()
		h.CompletedAt = &now
		h.ErrorMessage = &msg
		h.DurationMs = int(res.Duration.Milliseconds())
		s.persistHistory(ctx, &h)
		return res, err
	}
	res.TotalFetched = len(items)

	for _, pc := range items {
		// check existing
		existing, err := s.Contents.GetByProviderKey(ctx, pc.ProviderID, pc.ProviderContentID)
		if err == nil && existing != nil {
			// compare metrics
			oldM, err := s.Metrics.GetByContentID(ctx, existing.ID)
			if err != nil {
				res.FailedContents++
				s.Logger.Warn("metrics fetch failed", zap.Int64("content_id", existing.ID), zap.Error(err))
				continue
			}
			newSnap := MetricsSnapshot{
				Views:       val64(pc.Views),
				Likes:       val64(pc.Likes),
				ReadingTime: valInt(pc.ReadingTime),
				Reactions:   valInt(pc.Reactions),
			}
			oldSnap := MetricsSnapshot{
				Views:       oldM.Views,
				Likes:       oldM.Likes,
				ReadingTime: oldM.ReadingTime,
				Reactions:   oldM.Reactions,
			}
			if HasMetricsChanged(oldSnap, newSnap, s.Thresholds) {
				oldM.Views = newSnap.Views
				oldM.Likes = newSnap.Likes
				oldM.ReadingTime = newSnap.ReadingTime
				oldM.Reactions = newSnap.Reactions
				if err := s.Metrics.UpdateByContentID(ctx, existing.ID, oldM); err != nil {
					res.FailedContents++
					s.Logger.Error("metrics update failed", zap.Int64("content_id", existing.ID), zap.Error(err))
					continue
				}
				// recalc score
				if _, err := s.ScoreCalc.RecalculateScore(ctx, existing.ID); err != nil {
					res.FailedContents++
					s.Logger.Error("score recalc failed", zap.Int64("content_id", existing.ID), zap.Error(err))
					continue
				}
				res.UpdatedContents++
			} else {
				res.SkippedContents++
			}
			continue
		}
		// new content
		if _, _, err := s.ScoreCalc.ProcessNewContent(ctx, &pc); err != nil {
			res.FailedContents++
			s.Logger.Error("new content processing failed", zap.String("provider", providerID), zap.Error(err))
			continue
		}
		res.NewContents++
	}

	res.Duration = time.Since(start)
	status := entities.SyncStatusSuccess
	if len(res.Errors) > 0 {
		if res.NewContents+res.UpdatedContents+res.SkippedContents > 0 {
			status = entities.SyncStatusPartial
		} else {
			status = entities.SyncStatusFailed
		}
	}
	now := time.Now().UTC()
	h.SyncStatus = status
	h.TotalFetched = res.TotalFetched
	h.NewContents = res.NewContents
	h.UpdatedContents = res.UpdatedContents
	h.SkippedContents = res.SkippedContents
	h.FailedContents = res.FailedContents
	h.CompletedAt = &now
	if len(res.Errors) > 0 {
		msg := res.Errors[0]
		h.ErrorMessage = &msg
	}
	h.DurationMs = int(res.Duration.Milliseconds())
	s.persistHistory(ctx, &h)
	s.Logger.Info("sync completed", zap.String("provider", providerID), zap.Int("fetched", res.TotalFetched), zap.Duration("duration", res.Duration))
	return res, nil
}

func (s *ContentSyncService) GetLastSyncTime(ctx context.Context, providerID string) (time.Time, error) {
	h, err := s.HistoryRepo.GetLastSync(ctx, providerID)
	if err != nil || h == nil || h.CompletedAt == nil {
		return time.Time{}, err
	}
	return *h.CompletedAt, nil
}

func (s *ContentSyncService) GetSyncHistory(ctx context.Context, limit int) ([]entities.SyncHistory, error) {
	return s.HistoryRepo.GetAll(ctx, limit)
}

func (s *ContentSyncService) persistHistory(ctx context.Context, h *entities.SyncHistory) {
	if h == nil {
		return
	}
	var err error
	if h.ID == 0 {
		err = s.HistoryRepo.Create(ctx, h)
	} else {
		err = s.HistoryRepo.Update(ctx, h)
	}
	if err != nil {
		s.Logger.Warn("failed to persist sync history", zap.String("provider", h.ProviderID), zap.Error(err))
	}
}
