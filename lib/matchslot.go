package lib

import (
	"time"
)

// MatchSlot - struct holds info about matches
type MatchSlot struct {
	StartsAt  time.Time
	Duration  time.Duration
	notifiers []*time.Timer
	context   *Context
}

// NewMatchSlot - will create new match
func NewMatchSlot(startTime time.Time, duration time.Duration, ctx *Context) *MatchSlot {
	slot := &MatchSlot{
		StartsAt: startTime,
		Duration: duration,
		context:  ctx,
	}
	return slot
}

// Overlaps - return true if those two match slots are overlapsing
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

// Notify - will schedule jingles for this matchslot
func (m *MatchSlot) Notify(j *Jingle) {
	d := -1 * j.TimeBeforePoint

	var dur time.Duration
	if j.Point == MatchStart {
		dur = -1 * time.Since(m.StartsAt.Add(d))
	} else {
		dur = -1 * time.Since(m.StartsAt.Add(d).Add(m.Duration))
	}
	if dur > 0 {
		m.context.Log.Info("Scheduling jingle: %v -> '%v'", time.Now().Add(dur), j.Name)
		m.notifiers = append(m.notifiers, time.AfterFunc(dur, j.Play))
	}
}

// Cancel - will cancel all scheduled jingles for this match slot
func (m *MatchSlot) Cancel() {
	for _, t := range m.notifiers {
		t.Stop()
	}
}
