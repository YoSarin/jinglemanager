package lib

import (
	"time"
)

type Match struct {
	StartsAt time.Time
	Duration time.Duration
	Field    int
	HomeTeam *Team
	Awayteam *Team
}
