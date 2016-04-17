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
	Awayteam *Team
}
