package server

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/martin-reznik/jinglemanager/manager"
	"github.com/martin-reznik/logger"
	"net/http"
	"strconv"
)

// PlayerHandler - player handler
type PlayerHandler struct {
	Logger *logger.Log
}

// Add - will add song
func (p *PlayerHandler) Add(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	songFile := r.URL.Query().Get("filename")
	s := manager.FindSongByFile(songFile)
	if s != nil {
		output, _ := json.Marshal(s)
		w.Write(output)
		return
	}
	s, err := manager.NewSong(songFile, p.Logger)
	if err != nil {
		p.Logger.Error(err.Error())
		http.NotFound(w, r)
		return
	}
	output, _ := json.Marshal(s)
	w.Write(output)
}

// Play - will play song
func (p *PlayerHandler) Play(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))

	if err != nil {
		http.NotFound(w, r)
		return
	}

	s, err := manager.FindSong(id)
	if err != nil || s.IsPlaying() {
		p.Logger.Error(err.Error())
		http.NotFound(w, r)
		return
	}
	s.Play()

	output, _ := json.Marshal(s)
	w.Write(output)
}

// List - will list all actually playing songs
func (p *PlayerHandler) List(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	list := manager.GetAllPlaying()
	output, _ := json.Marshal(list)
	w.Write(output)
}

// Stop - will stop song
func (p *PlayerHandler) Stop(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	s, err := manager.FindSong(id)
	if err != nil || !s.IsPlaying() {
		http.NotFound(w, r)
		return
	}
	s.Stop()

	output, _ := json.Marshal(s)
	w.Write(output)
}
