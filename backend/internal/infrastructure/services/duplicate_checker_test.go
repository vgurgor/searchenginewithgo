package services

import "testing"

func TestHasMetricsChanged(t *testing.T) {
	th := MetricsThresholds{Percent: 5, AbsViews: 100, AbsLikes: 10, AbsReactions: 5}
	// No change
	old := MetricsSnapshot{Views: 1000, Likes: 100, ReadingTime: 10, Reactions: 50}
	newV := MetricsSnapshot{Views: 1000, Likes: 100, ReadingTime: 10, Reactions: 50}
	if HasMetricsChanged(old, newV, th) {
		t.Fatalf("expected no change")
	}
	// Absolute views change > 100
	newV.Views = 1201
	if !HasMetricsChanged(old, newV, th) {
		t.Fatalf("expected change on views abs")
	}
	// Percentage likes change > 5%
	newV = old
	newV.Likes = 106
	if !HasMetricsChanged(old, newV, th) {
		t.Fatalf("expected change on likes percent")
	}
	// ReadingTime any change
	newV = old
	newV.ReadingTime = 12
	if !HasMetricsChanged(old, newV, th) {
		t.Fatalf("expected change on reading time")
	}
	// Reactions abs threshold
	newV = old
	newV.Reactions = 56
	if !HasMetricsChanged(old, newV, th) {
		t.Fatalf("expected change on reactions")
	}
}
