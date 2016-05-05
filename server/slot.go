package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/martin-reznik/jinglemanager/lib"
	"net/http"
	"time"
)

type SlotHandler struct {
	Context *lib.Context
}

// Add - will add new match slot
func (h *SlotHandler) Add(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	h.Context.Tournament.AddMatchSlot(lib.NewMatchSlot(
		time.Now().Add(10*time.Second),
		5*time.Minute,
	))
}
