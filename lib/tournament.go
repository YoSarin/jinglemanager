package lib

// Tournament - struct holding data about tournament
type Tournament struct {
	Name string
}

// NewTournament - will create new tournament
func NewTournament(name string) *Tournament {
	return &Tournament{
		Name: name,
	}
}
