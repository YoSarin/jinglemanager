package lib

import (
	"time"
)

// Jingle - holds info about jingle (song, timing and so)
type Jingle struct {
	Song               *Song
	Name               string
	TimeFromMatchStart *time.Duration
}

// NewJingle - will create new Jingle
func NewJingle(name string, song *Song, timeFromStart *time.Duration) *Jingle {
	return &Jingle{
		Name:               name,
		Song:               song,
		TimeFromMatchStart: timeFromStart,
	}
}

// Play - will play jingle
func (j *Jingle) Play() {
	j.Song.Play()
}
