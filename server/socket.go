package server

import (
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/martin-reznik/jinglemanager/lib"
	"net/http"
)

// SocketHandler - struct for WebSocket handling
type SocketHandler struct {
	Context  *lib.Context
	Upgrader *websocket.Upgrader
}

type ping struct{}

// HandleChangeSocket - will handle sockets
func (h *SocketHandler) HandleChangeSocket(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.Context.Log.Error("upgrade: " + err.Error())
		return
	}
	defer c.Close()

	changeListener, deferFunc := lib.ChannelChange.Subscribe()
	defer deferFunc()

	for {
		select {
		case m := <-changeListener:
			err := c.WriteJSON(m)
			if err != nil {
				h.Context.Log.Error("Write error closing sock: " + err.Error())
				return
			}
		}
	}
}

// HandleLogSocket - will handle sockets
func (h *SocketHandler) HandleLogSocket(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.Context.Log.Error("upgrade: " + err.Error())
		return
	}
	defer c.Close()

	changeListener, deferFunc := lib.ChannelLog.Subscribe()
	defer deferFunc()

	for {
		select {
		case m := <-changeListener:
			err := c.WriteJSON(m)
			if err != nil {
				h.Context.Log.Error("Write error closing sock: " + err.Error())
				return
			}
		}
	}
}
