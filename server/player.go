package server

import (
	"encoding/json"
	"github.com/martin-reznik/jinglemanager/manager"
	"github.com/martin-reznik/logger"
	"net/http"
	"strconv"
)

// PlayerHandler - player handler
type PlayerHandler struct {
	Logger *logger.Log
}

// Play - will play song
func (p *PlayerHandler) Play(w http.ResponseWriter, r *http.Request) {
	songFile := r.URL.Query().Get("file")
	s, err := manager.NewSong(songFile, p.Logger)
	if err != nil {
		p.Logger.Error(err.Error())
		http.NotFound(w, r)
		return
	}
	s.Play()
	output, _ := json.Marshal(s)
	w.Write(output)
}

// List - will list all actually playing songs
func (p *PlayerHandler) List(w http.ResponseWriter, r *http.Request) {
	list := manager.GetAllPlaying()
	output, _ := json.Marshal(list)
	w.Write(output)
}

// Stop - will stop song
func (p *PlayerHandler) Stop(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	s, err := manager.FindPlayingSong(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	s.Stop()
}
