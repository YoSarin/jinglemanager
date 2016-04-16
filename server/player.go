package server

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/martin-reznik/jinglemanager/lib"
	"net/http"
)

// PlayerHandler - player handler
type PlayerHandler struct {
	Logger   LogI
	SongList *lib.FileList
}

// SongI - interface used to play songs
type SongI interface {
	Play()
	Pause()
	Stop()
	IsPlaying() bool
	Position() float64
	FileName() string
}

// Add - will add song
func (p *PlayerHandler) Add(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	songFile := r.URL.Query().Get("filename")
	p.SongList.AddUniq(songFile, p.Logger)
	output, _ := json.Marshal(p.SongList.GetAll())
	w.Write(output)
}

// Play - will play song
func (p *PlayerHandler) Play(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	f, err := p.SongList.Find(id)
	s, ok := f.(SongI)
	if err != nil || !ok || s.IsPlaying() {
		p.Logger.Error(err.Error())
		http.NotFound(w, r)
		return
	}
	s.Play()

	output, _ := json.Marshal(p.SongList.GetAll())
	w.Write(output)
}

// List - will list all actually playing songs
func (p *PlayerHandler) List(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	list := p.SongList.GetAll()

	output, _ := json.Marshal(list)
	w.Write(output)
}

// Stop - will stop song
func (p *PlayerHandler) Stop(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	f, err := p.SongList.Find(id)
	s, ok := f.(SongI)
	if err != nil || !ok || !s.IsPlaying() {
		http.NotFound(w, r)
		return
	}
	s.Stop()

	output, _ := json.Marshal(p.SongList.GetAll())
	w.Write(output)
}

// Pause - will stop song
func (p *PlayerHandler) Pause(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	f, err := p.SongList.Find(id)
	s, ok := f.(SongI)
	if err != nil || !ok || !s.IsPlaying() {
		http.NotFound(w, r)
		return
	}
	s.Pause()

	output, _ := json.Marshal(p.SongList.GetAll())
	w.Write(output)
}

// Delete - will remove song
func (p *PlayerHandler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	f, err := p.SongList.Find(id)
	s, ok := f.(SongI)
	if err != nil || !ok {
		http.NotFound(w, r)
		return
	}

	if s.IsPlaying() {
		s.Stop()
	}

	p.SongList.Delete(id)

	output, _ := json.Marshal(p.SongList.GetAll())
	w.Write(output)
}
