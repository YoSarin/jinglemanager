package manager

import (
	"github.com/martin-reznik/logger"
	"net/http"
)

// PlayerHandler - player handler
type PlayerHandler struct {
	Logger *logger.Log
}

// Play - will play song
func (p *PlayerHandler) Play(w http.ResponseWriter, r *http.Request) {
	songFile := r.URL.Query().Get("file")
	s := NewSong(songFile, p.Logger)
	s.Play()
}

// Stop - will stop song
func (p *PlayerHandler) Stop(w http.ResponseWriter, r *http.Request) {
	s, err := FindPlayingSong(1)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	s.Stop()
}
