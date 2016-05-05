package lib

import (
	"time"
)

// Tournament - struct holding data about tournament
type Tournament struct {
	Name    string
	Matches []*Match
}

// TournamentConfig - config for creating tournament schedule
type TournamentConfig struct {
	MinimalMatchLenght time.Duration
	FieldCount         int
}

// NewTournament - will create new tournament
func NewTournament(name string) *Tournament {

	t := &Tournament{
		Name: name,
	}

	ChannelChange.Emit(EventTypeTournamentChange, t)
	return t
}

// AddMatch - will add an match to Tournament
func (t *Tournament) AddMatch(m *Match) {
	t.Matches = append(t.Matches, m)
}

// SetName - will set new name for tournament
func (t *Tournament) SetName(name string) {
	t.Name = name
	ChannelChange.Emit(EventTypeTournamentChange, t)
}

func (t *Tournament) generateSchedule(startDates []time.Time, teamsBySeeding []*Team, config TournamentConfig) {
	panic("Not implemented yet")
}
