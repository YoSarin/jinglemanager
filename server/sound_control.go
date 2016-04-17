package server

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/martin-reznik/jinglemanager/lib"
	"net/http"
)

// SoundControlHandler - handler to controll sounds
type SoundControlHandler struct {
	Context *lib.Context
}

// Mute - will mute all apps
func (h *SoundControlHandler) Mute(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	h.Context.Sound.MuteApps()
}

// UnMute - will unmute all the apps
func (h *SoundControlHandler) UnMute(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	h.Context.Sound.UnMuteApps()
}

// List - will list all the apps
func (h *SoundControlHandler) List(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	out, _ := json.Marshal(h.Context.Sound.List())
	w.Write(out)
}

// Add - will add new app to list
func (h *SoundControlHandler) Add(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	app := r.FormValue("name")
	h.Context.Sound.Add(app)
	out, _ := json.Marshal(h.Context.Sound.List())
	w.Write(out)
}

// Delete - wil remove app from list
func (h *SoundControlHandler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	h.Context.Sound.Remove(id)
	out, _ := json.Marshal(h.Context.Sound.List())
	w.Write(out)
}
