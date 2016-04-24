package lib

// Tournament - struct holding data about tournament
type Tournament struct {
	Name string
}

// NewTournament - will create new tournament
func NewTournament(name string) *Tournament {

	t := &Tournament{
		Name: name,
	}

	ChannelChange.Emit(EventTypeTournamentChange, t)
	return t
}

// SetName - will set new name for tournament
func (t *Tournament) SetName(name string) {
	t.Name = name
	ChannelChange.Emit(EventTypeTournamentChange, t)
}
