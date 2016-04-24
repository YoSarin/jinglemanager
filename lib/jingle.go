package lib

import (
	"time"
)

// MatchPoint - point of match
type MatchPoint string

const (
	// MatchStart - start of match
	MatchStart = MatchPoint("match_start")
	// MatchEnd - end of match
	MatchEnd = MatchPoint("match_end")
)

// Jingle - holds info about jingle (song, timing and so)
type Jingle struct {
	Song            *Song
	Name            string
	TimeBeforePoint time.Duration
	Point           MatchPoint
}

// NewJingle - will create new Jingle
func NewJingle(name string, song *Song, timeBeforePoint time.Duration, point MatchPoint) *Jingle {
	j := &Jingle{
		Name:            name,
		Song:            song,
		TimeBeforePoint: timeBeforePoint,
		Point:           point,
	}

	ChannelChange.Emit(EventTypeJingleAdded, j)

	return j
}

// ID - will play jingle
func (j Jingle) ID() string {
	return j.Song.ID()
}

// Remove - will play jingle
func (j Jingle) Remove() {

}

// Play - will play jingle
func (j *Jingle) Play() {
	j.Song.Play()
}

// JingleList - will convert UniqueList into array of Jingles
func (l *UniqueList) JingleList() []*Jingle {
	out := make([]*Jingle, len(l.list))
	i := 0
	for _, val := range l.list {
		v, ok := val.(*Jingle)
		if ok {
			out[i] = v
			i++
		}
	}
	return out
}
