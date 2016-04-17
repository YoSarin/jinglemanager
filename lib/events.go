package lib

import (
	"fmt"
)

type channel string
type eventType string

var (
	// ChannelChange - Channel should contain or change events
	ChannelChange = channel("change")

	// EventTypeVolumeChange - event type related to apps volume change
	EventTypeVolumeChange = eventType("volume_change")
	// EventTypeSongChange - event type related to song change
	EventTypeSongChange = eventType("song_change")
	// EventTypeSongAdded - event type related to song change
	EventTypeSongAdded = eventType("song_added")
	// EventTypeSongRemoved - event type related to song change
	EventTypeSongRemoved = eventType("song_removed")
	// EventTypeAppAdded - event type related to app list change
	EventTypeAppAdded = eventType("app_added")
	// EventTypeAppRemoved - event type related to app list change
	EventTypeAppRemoved = eventType("app_removed")
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
		go func(ch chan interface{}, ev interface{}) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Listener gone", r)
				}
			}()
			ch <- ev
		}(ch, ev)
	}
}
