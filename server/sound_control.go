package server

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/martin-reznik/jinglemanager/lib"
	"net/http"
)

type ControlHandler struct {
	Logger LogI
	Player *lib.Controller
}

func (h *ControlHandler) Mute(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	h.Player.MuteApps()
}

func (h *ControlHandler) UnMute(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	h.Player.UnMuteApps()
}

func (h *ControlHandler) List(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	out, _ := json.Marshal(h.Player.List())
	w.Write(out)
}

func (h *ControlHandler) Add(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	app := ps.ByName("app")
	h.Player.Add(app)
	out, _ := json.Marshal(h.Player.List())
	w.Write(out)
}

func (h *ControlHandler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	app := ps.ByName("app")
	h.Player.Remove(app)
	out, _ := json.Marshal(h.Player.List())
	w.Write(out)
}
