package lib

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	// EventTypeJingleAdded - event type related to jingle change
	EventTypeJingleAdded = EventType("jingle_added")
	// EventTypeJingleRemoved - event type related to jingle change
	EventTypeJingleRemoved = EventType("jingle_removed")
)

var (
	// ChannelJingle - channel to emit jingle changes
	ChannelJingle = Channel{name: "jingle", allowed: map[EventType]bool{EventTypeJingleAdded: true, EventTypeJingleRemoved: true}}
)

// MatchPoint - point of match
type MatchPoint string

const (
	// MatchStart - start of match
	MatchStart = MatchPoint("match_start")
	// MatchEnd - end of match
	MatchEnd = MatchPoint("match_end")
	// MatchNone - no match related
	MatchNone = MatchPoint("match_none")
)

// Jingle - holds info about jingle (song, timing and so)
type Jingle struct {
	*Song
	Name            string
	TimeBeforePoint time.Duration
	Point           MatchPoint
	Context         *Context
}

// JingleStorage - struct used for storage
type JingleStorage struct {
	Name            string
	File            string
	Point           MatchPoint
	TimeBeforePoint time.Duration
}

// NewJingle - will create new Jingle
func NewJingle(name string, song *Song, timeBeforePoint time.Duration, point MatchPoint, ctx *Context) *Jingle {
	j := &Jingle{
		Song:            song,
		Name:            name,
		TimeBeforePoint: timeBeforePoint,
		Point:           point,
		Context:         ctx,
	}

	ChannelJingle.Emit(EventTypeJingleAdded, j)

	return j
}

// MarshalJSON - will convert song to JSON
func (j *Jingle) MarshalJSON() ([]byte, error) {
	data := struct {
		ID              string
		Song            *Song
		Name            string
		Point           MatchPoint
		TimeBeforePoint time.Duration
	}{
		j.ID(), j.Song, j.Name, j.Point, j.TimeBeforePoint,
	}
	return json.Marshal(data)
}

// ID - will play jingle
func (j Jingle) ID() string {
	return j.Song.ID()
}

// Play - will play jingle
func (j *Jingle) Play() {
	j.Context.Log.Info("Going to play jingle %v", j.Name)
	j.Context.Sound.MuteApps()
	j.Song.Play(func() {
		j.Context.Sound.UnMuteApps()
	})
}

// ToJingleStorage - will play jingle
func (j *Jingle) ToJingleStorage() *JingleStorage {
	return &JingleStorage{
		Name:            j.Name,
		TimeBeforePoint: j.TimeBeforePoint,
		File:            j.Song.FileName(),
		Point:           j.Point,
	}
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

// JingleStorageList - will convert UniqueList into array of JingleStorage
func (l *UniqueList) JingleStorageList() []*JingleStorage {
	out := make([]*JingleStorage, len(l.list))
	i := 0
	for _, val := range l.list {
		v, ok := val.(*Jingle)
		if ok {
			out[i] = v.ToJingleStorage()
			i++
		}
	}
	return out
}

// JingleName - will compose name of jingle
func JingleName(relativeTo string, dur time.Duration) (out string, point MatchPoint) {
	switch relativeTo {
	case "before_start":
		point = MatchStart
		if dur.Minutes() == 0 {
			out = "začátek zápasu"
		} else {
			out = fmt.Sprintf("%v minut před začátkem", dur.Minutes())
		}
	case "before_end":
		point = MatchEnd
		if dur.Minutes() == 0 {
			out = "konec zápasu"
		} else {
			out = fmt.Sprintf("%v minut do konce", dur.Minutes())
		}
	}
	return
}

// OnRemove - callback for removal form list
func (j *Jingle) OnRemove() {
	ChannelJingle.Emit(EventTypeJingleRemoved, j)
}
