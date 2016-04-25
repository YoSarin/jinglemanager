package server

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/martin-reznik/jinglemanager/lib"
	"net/http"
)

// PlayerHandler - player handler
type PlayerHandler struct {
	Context *lib.Context
}

// Add - will add song
func (p *PlayerHandler) Add(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	songFile := r.FormValue("filename")
    s, err := lib.NewSong(songFile, p.Context.Log)
    if err != nil {
        p.Context.Log.Error(err.Error())
    } else {
        p.Context.Songs.AddUniq(s, p.Context.Log)
    }
	output, _ := json.Marshal(p.Context.Songs.GetAll())
	w.Write(output)
}

// Play - will play song
func (p *PlayerHandler) Play(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	f, err := p.Context.Songs.Find(id)
	s, ok := f.(SongI)
	if err != nil || !ok {
		p.Context.Log.Error(err.Error())
		http.NotFound(w, r)
		return
	}
	s.Play()

	output, _ := json.Marshal(p.Context.Songs.GetAll())
	w.Write(output)
}

// List - will list all actually playing songs
func (p *PlayerHandler) List(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	list := p.Context.Songs.GetAll()

	output, _ := json.Marshal(list)
	w.Write(output)
}

// Stop - will stop song
func (p *PlayerHandler) Stop(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	f, err := p.Context.Songs.Find(id)
	s, ok := f.(SongI)
	if err != nil || !ok {
		http.NotFound(w, r)
		return
	}
	s.Stop()

	output, _ := json.Marshal(p.Context.Songs.GetAll())
	w.Write(output)
}

// Pause - will stop song
func (p *PlayerHandler) Pause(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	f, err := p.Context.Songs.Find(id)
	s, ok := f.(SongI)
	if err != nil || !ok {
		http.NotFound(w, r)
		return
	}
	s.Pause()

	output, _ := json.Marshal(p.Context.Songs.GetAll())
	w.Write(output)
}

// Delete - will remove song
func (p *PlayerHandler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	f, err := p.Context.Songs.Find(id)
	s, ok := f.(SongI)
	if err != nil || !ok {
		http.NotFound(w, r)
		return
	}

	if s.IsPlaying() {
		s.Stop()
	}

	p.Context.Songs.Delete(id)

	output, _ := json.Marshal(p.Context.Songs.GetAll())
	w.Write(output)
}
