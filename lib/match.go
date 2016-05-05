package lib

import (
	"time"
)

// Match - struct holds info about matches
type Match struct {
	StartsAt time.Time
	Duration time.Duration
	Field    int
	HomeTeam *Team
	AwayTeam *Team
}

// NewMatch - will create new match
func NewMatch(startTime time.Time, duration time.Duration, field int, homeTeam *Team, awayTeam *Team) *Match {
	return &Match{
		StartsAt: startTime,
		Duration: duration,
		Field:    field,
		HomeTeam: homeTeam,
		AwayTeam: awayTeam,
	}
}
