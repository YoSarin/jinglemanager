package lib

import (
	"time"
)

type Jingle struct {
	Song               *Song
	Name               string
	TimeFromMatchStart *time.Duration
}

func NewJingle(name string, song *Song, timeFromStart *time.Duration) *Jingle {
	return &Jingle{
		Name:               name,
		Song:               song,
		TimeFromMatchStart: timeFromStart,
	}
}

func (j *Jingle) Play() {
	j.Song.Play()
}
