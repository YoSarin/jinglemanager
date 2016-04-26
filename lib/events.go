package lib

type eventType string
type channel struct {
	name    string
	allowed map[eventType]bool
}

const (
	// EventTypeLog - logging event
	EventTypeLog = eventType("log")
	// EventTypeCleanup - event type related to total cleanup
	EventTypeCleanup = eventType("cleanup")
	// EventTypeVolumeChange - event type related to apps volume change
	EventTypeVolumeChange = eventType("volume_changed")
	// EventTypeSongChange - event type related to song change
	EventTypeSongChange = eventType("song_changed")
	// EventTypeSongAdded - event type related to song change
	EventTypeSongAdded = eventType("song_added")
	// EventTypeSongRemoved - event type related to song change
	EventTypeSongRemoved = eventType("song_removed")
	// EventTypeAppAdded - event type related to app list change
	EventTypeAppAdded = eventType("app_added")
	// EventTypeAppRemoved - event type related to app list change
	EventTypeAppRemoved = eventType("app_removed")
	// EventTypeTournamentChange - event type related to tournament change
	EventTypeTournamentChange = eventType("tournament_change")
	// EventTypeJingleAdded - event type related to jingle change
	EventTypeJingleAdded = eventType("jingle_added")
	// EventTypeJingleRemoved - event type related to jingle change
	EventTypeJingleRemoved = eventType("jingle_removed")
)

var (
	// ChannelChange - Channel should contain change events
	ChannelChange = channel{name: "change", allowed: map[eventType]bool{
		EventTypeAppAdded:         true,
		EventTypeAppRemoved:       true,
		EventTypeCleanup:          true,
		EventTypeSongAdded:        true,
		EventTypeSongChange:       true,
		EventTypeSongRemoved:      true,
		EventTypeVolumeChange:     true,
        EventTypeTournamentChange: true,
        EventTypeJingleAdded:      true,
        EventTypeJingleRemoved:    true,
	}}
	// ChannelLog - Channel should contain log events
	ChannelLog = channel{name: "log", allowed: map[eventType]bool{EventTypeLog: true}}
)

type event struct {
	Type string
	Data interface{}
}

var listeners = make(map[*channel][]chan interface{})

// Subscribe - Will subscribe new listener and returns him and his defer function
func (c *channel) Subscribe() (chan interface{}, func()) {
	ch := make(chan interface{})
	listeners[c] = append(listeners[c], ch)

	return ch, func() {
		for idx, chn := range listeners[c] {
			if chn == ch {
				listeners[c] = append(listeners[c][:idx], listeners[c][idx+1:]...)
				break
			}
		}
		close(ch)
	}
}

// Emit - Will emit event to all listeners
func (c *channel) Emit(evType eventType, data interface{}) {
	if c.allowed[evType] {
		ev := struct {
			Type eventType
			Data interface{}
		}{Type: evType, Data: data}
		for _, ch := range listeners[c] {
			go func(ch chan interface{}, ev interface{}) {
				defer func() {
					if r := recover(); r != nil && r != "send on closed channel" {
						// Unknown and unexpected error
						panic(r)
					}
				}()
				ch <- ev
			}(ch, ev)
		}
	}
}
