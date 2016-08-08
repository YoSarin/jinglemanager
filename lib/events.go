package lib

import (
	"sync"
)

// EventType - event types emmited on Channels
type EventType string

// Channel - channel to emmit events to
type Channel struct {
	name    string
	allowed map[EventType]bool
}

const (
	// EventTypeLog - logging event
	EventTypeLog = EventType("log")
)

var (
	// ChannelLog - Channel should contain log events
	ChannelLog = Channel{name: "log", allowed: map[EventType]bool{EventTypeLog: true}}
)

type event struct {
	Type string
	Data interface{}
}

var listeners = make(map[*Channel][]chan interface{})
var listenersMutex = &sync.Mutex{}

// MultiSubscribe - will create subscribe to multiple channels
func MultiSubscribe(list []*Channel) (chan interface{}, func()) {
	ch := make(chan interface{})
	for _, c := range list {
		listenersMutex.Lock()
		listeners[c] = append(listeners[c], ch)
		listenersMutex.Unlock()
	}
	return ch, func() {
		for _, c := range list {
			listenersMutex.Lock()
			for idx, chn := range listeners[c] {
				if chn == ch {
					listeners[c] = append(listeners[c][:idx], listeners[c][idx+1:]...)
				}
			}
			listenersMutex.Unlock()
		}
		close(ch)
	}
}

// Subscribe - Will subscribe new listener and returns him and his defer function
func (c *Channel) Subscribe() (chan interface{}, func()) {
	ch := make(chan interface{})
	listenersMutex.Lock()
	listeners[c] = append(listeners[c], ch)
	listenersMutex.Unlock()

	return ch, func() {
		listenersMutex.Lock()
		for idx, chn := range listeners[c] {
			if chn == ch {
				listeners[c] = append(listeners[c][:idx], listeners[c][idx+1:]...)
				break
			}
		}
		listenersMutex.Unlock()
		close(ch)
	}
}

// Name - will return channels name
func (c *Channel) Name() string {
	return c.name
}

// Emit - Will emit event to all listeners
func (c *Channel) Emit(evType EventType, data interface{}) {
	if c.allowed[evType] {
		ev := struct {
			Type EventType
			Data interface{}
		}{Type: evType, Data: data}

		listenersMutex.Lock()
		for _, ch := range listeners[c] {
			go func(ch chan interface{}, ev interface{}) {
				defer func() {
					if r := recover(); r != nil && r != "send on closed Channel" {
						// Unknown and unexpected error
						// panic(r)
					}
				}()
				ch <- ev
			}(ch, ev)
		}
		listenersMutex.Unlock()
	}
}
