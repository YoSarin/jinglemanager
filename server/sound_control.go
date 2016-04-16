package server

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/martin-reznik/jinglemanager/lib"
	"net/http"
)

// SoundControlHandler - handler to controll sounds
type SoundControlHandler struct {
	Logger       LogI
	SoundControl *lib.SoundController
}

// Mute - will mute all apps
func (h *SoundControlHandler) Mute(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	h.SoundControl.MuteApps()
}

// UnMute - will unmute all the apps
func (h *SoundControlHandler) UnMute(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	h.SoundControl.UnMuteApps()
}

// List - will list all the apps
func (h *SoundControlHandler) List(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	out, _ := json.Marshal(h.SoundControl.List())
	w.Write(out)
}

// Add - will add new app to list
func (h *SoundControlHandler) Add(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	app := ps.ByName("app")
	h.SoundControl.Add(app)
	out, _ := json.Marshal(h.SoundControl.List())
	w.Write(out)
}

// Delete - wil remove app from list
func (h *SoundControlHandler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	app := ps.ByName("app")
	h.SoundControl.Remove(app)
	out, _ := json.Marshal(h.SoundControl.List())
	w.Write(out)
}
