package server

import (
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/martin-reznik/jinglemanager/lib"
	"net/http"
	"time"
)

// SocketHandler - struct for WebSocket handling
type SocketHandler struct {
	Context  *lib.Context
	Upgrader *websocket.Upgrader
}

type ping struct{}

// HandleSocket - will handle sockets
func (h *SocketHandler) HandleSocket(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.Context.Log.Error("upgrade: " + err.Error())
		return
	}
	defer c.Close()
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case m := <-h.Context.Changes:
			h.Context.Log.Info("write")
			err := c.WriteJSON(m)
			if err != nil {
				h.Context.Log.Error("Write error: " + err.Error())
				return
			}
		case <-ticker.C:
			err := c.WriteJSON(ping{})
			if err != nil {
				h.Context.Log.Error("Ping error: " + err.Error())
				return
			}
		}
	}
}
