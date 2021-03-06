package server

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/martin-reznik/jinglemanager/lib"
	"net/http"
	"strconv"
	"time"
)

// SlotHandler - handling slot API events
type SlotHandler struct {
	Context *lib.Context
}

// Add - will add new match slot
func (h *SlotHandler) Add(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	dur, _ := strconv.Atoi(r.FormValue("duration"))
	start, _ := time.Parse("2006-01-02 15:04:05", r.FormValue("start"))
	h.Context.Tournament.AddMatchSlot(lib.NewMatchSlot(
		start,
		time.Duration(dur)*time.Minute,
		h.Context,
	))
}

// Postpone - will postpone all match slots by given duration in minutes
func (h *SlotHandler) Postpone(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	dur, _ := strconv.Atoi(r.FormValue("postpone"))
	for _, slot := range h.Context.Tournament.MatchSlots {
		slot.StartsAt = slot.StartsAt.Add(time.Duration(dur) * time.Minute)
	}
	h.List(w, r, ps)
}

// Remove - will remove match slot
func (h *SlotHandler) Remove(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	dur, _ := strconv.Atoi(r.FormValue("duration"))
	start, _ := time.Parse("2006-01-02 15:04:05", r.FormValue("start"))
	h.Context.Tournament.AddMatchSlot(lib.NewMatchSlot(
		start,
		time.Duration(dur)*time.Minute,
		h.Context,
	))
}

// List - will return list of match slots
func (h *SlotHandler) List(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	output, _ := json.Marshal(h.Context.Tournament.MatchSlots)
	w.Write(output)
}
