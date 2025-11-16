package services

type MetricsThresholds struct {
	Percent int
	AbsViews int
	AbsLikes int
	AbsReactions int
}

type MetricsSnapshot struct {
	Views       int64
	Likes       int64
	ReadingTime int
	Reactions   int
}

func HasMetricsChanged(oldM, newM MetricsSnapshot, t MetricsThresholds) bool {
	// Reading time: any change matters
	if oldM.ReadingTime != newM.ReadingTime {
		return true
	}
	if significantChange(oldM.Views, newM.Views, t.Percent, int64(t.AbsViews)) {
		return true
	}
	if significantChange(oldM.Likes, newM.Likes, t.Percent, int64(t.AbsLikes)) {
		return true
	}
	if significantChange(int64(oldM.Reactions), int64(newM.Reactions), t.Percent, int64(t.AbsReactions)) {
		return true
	}
	return false
}

func significantChange(oldVal, newVal int64, percent int, abs int64) bool {
	diff := newVal - oldVal
	if diff < 0 {
		diff = -diff
	}
	if diff >= abs {
		return true
	}
	if oldVal == 0 {
		// if old is 0, percentage change is not well-defined; use abs only
		return diff > 0
	}
	p := (float64(diff) / float64(oldVal)) * 100.0
	return p >= float64(percent)
}


