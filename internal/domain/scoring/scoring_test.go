package scoring

import (
	"testing"
	"time"

	"search_engine/internal/domain/entities"
)

func newEngine() *ScoringEngine {
	return &ScoringEngine{
		VideoTypeMultiplier: 1.5,
		TextTypeMultiplier:  1.0,
		Freshness: FreshnessConfig{
			WithinOneWeekScore:    5,
			WithinOneMonthScore:   3,
			WithinThreeMonthsScore: 1,
		},
	}
}

func TestVideoHighEngagement(t *testing.T) {
	engine := newEngine()
	p := time.Now().UTC().Add(-3 * 24 * time.Hour)
	c := entities.Content{ContentType: entities.ContentTypeVideo, PublishedAt: &p}
	m := entities.ContentMetrics{Views: 100000, Likes: 5000}
	score, _ := engine.CalculateScore(c, m)
	if score <= 150 {
		t.Fatalf("expected high score > 150 got %v", score)
	}
}

func TestTextLowEngagementOld(t *testing.T) {
	engine := newEngine()
	p := time.Now().UTC().AddDate(0, -6, 0)
	c := entities.Content{ContentType: entities.ContentTypeText, PublishedAt: &p}
	m := entities.ContentMetrics{ReadingTime: 5, Reactions: 10}
	score, _ := engine.CalculateScore(c, m)
	if score >= 20 {
		t.Fatalf("expected low score < 20 got %v", score)
	}
}

func TestVideoZeroViews(t *testing.T) {
	engine := newEngine()
	now := time.Now().UTC()
	c := entities.Content{ContentType: entities.ContentTypeVideo, PublishedAt: &now}
	m := entities.ContentMetrics{Views: 0, Likes: 0}
	score, _ := engine.CalculateScore(c, m)
	if score < 0 {
		t.Fatalf("expected non-negative base score got %v", score)
	}
}

func TestVeryOldContentFreshnessZero(t *testing.T) {
	engine := newEngine()
	p := time.Now().UTC().AddDate(-2, 0, 0)
	c := entities.Content{ContentType: entities.ContentTypeText, PublishedAt: &p}
	m := entities.ContentMetrics{}
	score, _ := engine.CalculateScore(c, m)
	// freshness should add 0; base and engagement are zero -> final 0
	if score != 0 {
		t.Fatalf("expected 0 score for old content got %v", score)
	}
}


