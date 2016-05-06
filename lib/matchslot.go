package lib

import (
	"time"
)

// MatchSlot - struct holds info about matches
type MatchSlot struct {
	StartsAt time.Time
	Duration time.Duration
}

// NewMatchSlot - will create new match
func NewMatchSlot(startTime time.Time, duration time.Duration) *MatchSlot {
	slot := &MatchSlot{
		StartsAt: startTime,
		Duration: duration,
	}
	return slot
}

func (m *MatchSlot) Overlaps(m2 MatchSlot) bool {
	if m.StartsAt.After(m2.StartsAt) && m.StartsAt.Before(m2.StartsAt.Add(m2.Duration)) {
		return true
	}
	if m2.StartsAt.After(m.StartsAt) && m2.StartsAt.Before(m.StartsAt.Add(m.Duration)) {
		return true
	}
	if m.StartsAt.Equal(m2.StartsAt) {
		return true
	}
	return false
}

func (m *MatchSlot) Notify(d time.Duration, p MatchPoint, notifier func()) {
	var dur time.Duration
	if p == MatchStart {
		dur = -1 * time.Since(m.StartsAt.Add(d))
	} else {
		dur = -1 * time.Since(m.StartsAt.Add(d).Add(m.Duration))
	}

	time.AfterFunc(dur, notifier)
}
