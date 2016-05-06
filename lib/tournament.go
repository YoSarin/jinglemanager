package lib

import (
	"time"
)

// Tournament - struct holding data about tournament
type Tournament struct {
	Name       string
	MatchSlots []*MatchSlot
	Context    *Context
}

// TournamentConfig - config for creating tournament schedule
type TournamentConfig struct {
	MinimalMatchLenght time.Duration
	FieldCount         int
}

// NewTournament - will create new tournament
func NewTournament(name string, context *Context) *Tournament {

	t := &Tournament{
		Name:       name,
		Context:    context,
		MatchSlots: []*MatchSlot{},
	}

	ChannelChange.Emit(EventTypeTournamentChange, t)
	return t
}

// AddMatchSlot - will add an match to Tournament
func (t *Tournament) AddMatchSlot(m *MatchSlot) {
	for _, s := range t.MatchSlots {
		if s.Overlaps(*m) {
			return
		}
	}
	t.MatchSlots = append(t.MatchSlots, m)
	for _, j := range t.Context.Jingles.JingleList() {
		m.Notify(j.TimeBeforePoint, j.Point, func() {
			t.Context.Log.Info("Kaboom!")
		})
	}
}

// SetName - will set new name for tournament
func (t *Tournament) SetName(name string) {
	t.Name = name
	ChannelChange.Emit(EventTypeTournamentChange, t)
}
