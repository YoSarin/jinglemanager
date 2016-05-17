package server

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/martin-reznik/jinglemanager/lib"
	"net/http"
	"strconv"
	"time"
)

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
	))
}

func (h *SlotHandler) List(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	output, _ := json.Marshal(h.Context.Tournament.MatchSlots)
	w.Write(output)
}
