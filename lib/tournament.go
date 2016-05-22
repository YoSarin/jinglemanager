package lib

import (
	"time"
)

// Tournament - struct holding data about tournament
type Tournament struct {
	Name       string
	MatchSlots []*TournamentMatchSlot
	context    *Context
}

type TournamentMatchSlot struct {
	*MatchSlot
	place int
}

// TournamentConfig - config for creating tournament schedule
type TournamentConfig struct {
	MinimalMatchLenght time.Duration
	FieldCount         int
}

const (
	// EventTypeTournamentChange - event type related to tournament change
	EventTypeTournamentChange = EventType("tournament_change")
	// EventTypeMatchSlotAdded - event type related to match slot added
	EventTypeMatchSlotAdded = EventType("match_slot_added")
	// EventTypeMatchSlotRemoved - event type related to match slot removed
	EventTypeMatchSlotRemoved = EventType("match_slot_removed")
	// EventTypeMatchSlotChange - event type related to match slot change
	EventTypeMatchSlotChange = EventType("match_slot_change")
)

var (
	ChannelTournament = Channel{name: "tournament", allowed: map[EventType]bool{
		EventTypeTournamentChange: true,
		EventTypeMatchSlotAdded:   true,
		EventTypeMatchSlotRemoved: true,
		EventTypeMatchSlotChange:  true,
	}}
)

// NewTournament - will create new tournament
func NewTournament(name string, context *Context) *Tournament {

	t := &Tournament{
		Name:       name,
		context:    context,
		MatchSlots: []*TournamentMatchSlot{},
	}

	ChannelTournament.Emit(EventTypeTournamentChange, t)
	return t
}

// AddMatchSlot - will add an match to Tournament
func (t *Tournament) AddMatchSlot(m *MatchSlot) {
	if m == nil {
		return
	}
	for _, s := range t.MatchSlots {
		if s.Overlaps(*m) {
			return
		}
	}
	s := &TournamentMatchSlot{m, len(t.MatchSlots)}
	t.MatchSlots = append(t.MatchSlots, s)
	for _, j := range t.context.Jingles.JingleList() {
		m.Notify(j.TimeBeforePoint, j.Point, j.Play)
	}
	ChannelTournament.Emit(EventTypeMatchSlotAdded, struct {
		Slot  *TournamentMatchSlot
		Place int
	}{
		Slot:  s,
		Place: s.place,
	})
}

// RemoveMatchSlot - will set new name for tournament
func (t *Tournament) RemoveMatchSlot(place int) {
	t.MatchSlots[place].Cancel()
	t.MatchSlots[place] = nil
	ChannelTournament.Emit(EventTypeMatchSlotRemoved, place)
}

// PlanJingles - will plan jingles
func (t *Tournament) PlanJingles() {
	for _, m := range t.MatchSlots {
		for _, j := range t.context.Jingles.JingleList() {
			m.Cancel()
			m.Notify(j.TimeBeforePoint, j.Point, j.Play)
		}
	}
}

// SetName - will set new name for tournament
func (t *Tournament) SetName(name string) {
	t.Name = name
	ChannelTournament.Emit(EventTypeTournamentChange, t)
}
