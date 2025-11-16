package scoring

import (
	"math"
	"time"

	"search_engine/internal/domain/entities"
)

type FreshnessConfig struct {
	WithinOneWeekScore     float64
	WithinOneMonthScore    float64
	WithinThreeMonthsScore float64
}

type ScoringEngine struct {
	VideoTypeMultiplier float64
	TextTypeMultiplier  float64
	Freshness           FreshnessConfig
}

type IScoringService interface {
	CalculateScore(content *entities.Content, metrics *entities.ContentMetrics) (float64, error)
}

func (s *ScoringEngine) CalculateScore(content *entities.Content, metrics *entities.ContentMetrics) (float64, error) {
	base := s.calculateBaseScore(content.ContentType, metrics)
	typeMul := s.getTypeMultiplier(content.ContentType)
	fresh := s.calculateFreshnessScore(content.PublishedAt)
	eng := s.calculateEngagementScore(content.ContentType, metrics)
	final := (base * typeMul) + fresh + eng
	return round2(final), nil
}

func (s *ScoringEngine) calculateBaseScore(ct entities.ContentType, m *entities.ContentMetrics) float64 {
	switch ct {
	case entities.ContentTypeVideo:
		views := float64(max64(m.Views, 0))
		likes := float64(max64(m.Likes, 0))
		return (views / 1000.0) + (likes / 100.0)
	default:
		rt := float64(maxInt(m.ReadingTime, 0))
		reac := float64(maxInt(m.Reactions, 0))
		return rt + (reac / 50.0)
	}
}

func (s *ScoringEngine) calculateFreshnessScore(p *time.Time) float64 {
	if p == nil {
		return 0
	}
	days := time.Since(p.UTC()).Hours() / 24
	switch {
	case days <= 7:
		return s.Freshness.WithinOneWeekScore
	case days <= 30:
		return s.Freshness.WithinOneMonthScore
	case days <= 90:
		return s.Freshness.WithinThreeMonthsScore
	default:
		return 0
	}
}

func (s *ScoringEngine) calculateEngagementScore(ct entities.ContentType, m *entities.ContentMetrics) float64 {
	switch ct {
	case entities.ContentTypeVideo:
		if m.Views > 0 {
			return (float64(max64(m.Likes, 0)) / float64(m.Views)) * 10.0
		}
		return 0
	default:
		if m.ReadingTime > 0 {
			return (float64(maxInt(m.Reactions, 0)) / float64(m.ReadingTime)) * 5.0
		}
		return 0
	}
}

func (s *ScoringEngine) getTypeMultiplier(ct entities.ContentType) float64 {
	if ct == entities.ContentTypeVideo {
		return s.VideoTypeMultiplier
	}
	return s.TextTypeMultiplier
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func max64(v int64, minValue int64) int64 {
	if v < minValue {
		return minValue
	}
	return v
}

func maxInt(v int, minValue int) int {
	if v < minValue {
		return minValue
	}
	return v
}
