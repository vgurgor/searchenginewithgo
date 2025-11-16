package services

import (
	"context"
	"time"

	"go.uber.org/zap"

	"search_engine/internal/domain/entities"
	"search_engine/internal/domain/providers"
	"search_engine/internal/domain/repositories"
	"search_engine/internal/domain/scoring"
)

type ScoreCalculatorService struct {
	Contents repositories.ContentRepository
	Metrics  repositories.ContentMetricsRepository
	Engine   scoring.IScoringService
	Logger   *zap.Logger
}

func (s *ScoreCalculatorService) ProcessNewContent(ctx context.Context, pc *providers.ProviderContent) (int64, float64, error) {
	c := entities.Content{
		ProviderID:        pc.ProviderID,
		ProviderContentID: pc.ProviderContentID,
		Title:             pc.Title,
		ContentType:       mapContentType(pc.ContentType),
		URL:               strPtrOrNil(pc.URL),
		ThumbnailURL:      strPtrOrNil(pc.ThumbnailURL),
		Description:       strPtrOrNil(pc.Description),
		PublishedAt:       timePtrOrNil(pc.PublishedAt),
	}
	if err := s.Contents.Create(ctx, &c); err != nil {
		return 0, 0, err
	}
	m := entities.ContentMetrics{
		ContentID:   c.ID,
		Views:       val64(pc.Views),
		Likes:       val64(pc.Likes),
		ReadingTime: valInt(pc.ReadingTime),
		Reactions:   valInt(pc.Reactions),
	}
	score, err := s.Engine.CalculateScore(&c, &m)
	if err != nil {
		return c.ID, 0, err
	}
	m.FinalScore = score
	now := time.Now().UTC()
	m.RecalculatedAt = &now
	if err := s.Metrics.Create(ctx, &m); err != nil {
		return c.ID, 0, err
	}
	s.Logger.Info("score calculated", zap.Int64("content_id", c.ID), zap.Float64("score", score))
	return c.ID, score, nil
}

func (s *ScoreCalculatorService) RecalculateScore(ctx context.Context, contentID int64) (float64, error) {
	c, err := s.Contents.GetByID(ctx, contentID)
	if err != nil {
		return 0, err
	}
	m, err := s.Metrics.GetByContentID(ctx, contentID)
	if err != nil {
		return 0, err
	}
	score, err := s.Engine.CalculateScore(c, m)
	if err != nil {
		return 0, err
	}
	m.FinalScore = score
	now := time.Now().UTC()
	m.RecalculatedAt = &now
	if err := s.Metrics.UpdateByContentID(ctx, contentID, m); err != nil {
		return 0, err
	}
	s.Logger.Info("score recalculated", zap.Int64("content_id", contentID), zap.Float64("score", score))
	return score, nil
}

func mapContentType(t string) entities.ContentType {
	if t == "video" || t == "Video" {
		return entities.ContentTypeVideo
	}
	return entities.ContentTypeText
}

func strPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
func timePtrOrNil(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	tt := t.UTC()
	return &tt
}
func val64(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}
func valInt(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}
