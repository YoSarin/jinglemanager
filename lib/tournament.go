package lib

import (
	"math"
	"time"
)

// Tournament - struct holding data about tournament
type Tournament struct {
	Name            string
	MatchSlots      []*TournamentMatchSlot
	Authorization   map[string]string
	Public          bool
	context         *Context
	jingleCheckStop chan bool
}

// TournamentMatchSlot - marshalable container for match slots
type TournamentMatchSlot struct {
	*MatchSlot
	place int
}

// TournamentConfig - config for creating tournament schedule
type TournamentConfig struct {
	MinimalMatchLenght time.Duration
	FieldCount         int
}

const (
	// EventTypeTournamentChange - event type related to tournament change
	EventTypeTournamentChange = EventType("tournament_change")
	// EventTypeMatchSlotAdded - event type related to match slot added
	EventTypeMatchSlotAdded = EventType("match_slot_added")
	// EventTypeMatchSlotRemoved - event type related to match slot removed
	EventTypeMatchSlotRemoved = EventType("match_slot_removed")
	// EventTypeMatchSlotChange - event type related to match slot change
	EventTypeMatchSlotChange = EventType("match_slot_change")
)

var (
	// ChannelTournament - channel to emmit tournament changes to
	ChannelTournament = Channel{name: "tournament", allowed: map[EventType]bool{
		EventTypeTournamentChange: true,
		EventTypeMatchSlotAdded:   true,
		EventTypeMatchSlotRemoved: true,
		EventTypeMatchSlotChange:  true,
	}}
)

// NewTournament - will create new tournament
func NewTournament(name string, context *Context) *Tournament {

	t := &Tournament{
		Name:            name,
		context:         context,
		MatchSlots:      []*TournamentMatchSlot{},
		jingleCheckStop: make(chan bool),
	}

	var duration = time.Duration(1) * time.Second

	ticker := time.NewTicker(duration)
	go func() {
		for {
			select {
			case <-ticker.C:
				for _, j := range t.CurrentJingles(time.Now(), duration+time.Duration(1)*time.Second) {
					if !j.IsPlaying() {
						j.Play()
					}
				}
			case <-t.jingleCheckStop:
				ticker.Stop()
			}
		}
	}()

	ChannelTournament.Emit(EventTypeTournamentChange, t)
	return t
}

// AddMatchSlot - will add an match to Tournament
func (t *Tournament) AddMatchSlot(m *MatchSlot) {
	if m == nil {
		return
	}
	for _, s := range t.MatchSlots {
		if s.Overlaps(*m) {
			return
		}
	}
	s := &TournamentMatchSlot{m, len(t.MatchSlots)}
	t.MatchSlots = append(t.MatchSlots, s)

	ChannelTournament.Emit(EventTypeMatchSlotAdded, struct {
		Slot  *TournamentMatchSlot
		Place int
	}{
		Slot:  s,
		Place: s.place,
	})
}

// RemoveMatchSlot - will set new name for tournament
func (t *Tournament) RemoveMatchSlot(place int) {
	t.MatchSlots[place] = nil
	ChannelTournament.Emit(EventTypeMatchSlotRemoved, place)
}

// CancelJingles - will cancel jingles
func (t *Tournament) CancelJingles() {
	t.jingleCheckStop <- true
}

// CurrentJingles - will return list of jingles which should run at this very point in time
func (t *Tournament) CurrentJingles(when time.Time, window time.Duration) []*Jingle {
	var out []*Jingle
	for _, m := range t.MatchSlots {
		for _, j := range t.context.Jingles.JingleList() {
			d := -1 * j.TimeBeforePoint

			var timeTo time.Duration
			if j.Point == MatchStart {
				timeTo = -1 * time.Since(m.StartsAt.Add(d))
			} else {
				timeTo = -1 * time.Since(m.StartsAt.Add(d).Add(m.Duration))
			}
			if (time.Duration(math.Abs(timeTo.Seconds())) * time.Second) < window {
				out = append(out, j)
			}
		}
	}
	return out
}

// SetName - will set new name for tournament
func (t *Tournament) SetName(name string) {
	t.Name = name
	ChannelTournament.Emit(EventTypeTournamentChange, t)
}
