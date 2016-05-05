package lib

import (
	"time"
)

// Tournament - struct holding data about tournament
type Tournament struct {
	Name       string
	log        LogI
	MatchSlots []*MatchSlot
}

// TournamentConfig - config for creating tournament schedule
type TournamentConfig struct {
	MinimalMatchLenght time.Duration
	FieldCount         int
}

// NewTournament - will create new tournament
func NewTournament(name string, log LogI) *Tournament {

	t := &Tournament{
		Name:       name,
		log:        log,
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
}

// SetName - will set new name for tournament
func (t *Tournament) SetName(name string) {
	t.Name = name
	ChannelChange.Emit(EventTypeTournamentChange, t)
}
