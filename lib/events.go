package lib

type channel string
type eventType string

var (
	// ChannelChange - Channel should contain or change events
	ChannelChange = channel("change")

	// EventTypeVolumeChange - event type related to apps volume change
	EventTypeVolumeChange = eventType("volume_change")
	// EventTypeSongChange - event type related to song change
	EventTypeSongChange = eventType("song_change")
)

type event struct {
	Type string
	Data interface{}
}

var listeners = make(map[channel][]chan interface{})

// Subscribe - Will subscribe new listener and returns him and his defer function
func (c channel) Subscribe() (chan interface{}, func()) {
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
func (c channel) Emit(evType eventType, data interface{}) {
	ev := struct {
		Type eventType
		Data interface{}
	}{Type: evType, Data: data}
	for _, ch := range listeners[c] {
		ch <- ev
	}
}
