package server

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/martin-reznik/jinglemanager/lib"
	"net/http"
	"strconv"
	"time"
)

// JingleHandler - player handler
type JingleHandler struct {
	Context *lib.Context
}

// Add - will add song
func (h *JingleHandler) Add(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	file, info, err := r.FormFile("file")
	if err != nil {
		h.Context.Log.Error("File upload failed: " + err.Error())
		http.Error(w, "upload failed", 500)
		return
	}
	defer file.Close()

	filename, err := h.Context.SaveSong(file, info.Filename)
    s, err := lib.NewSong(filename, h.Context.Log)
    if err != nil {
        h.Context.Log.Error(err.Error())
    } else {
        h.Context.Songs.AddUniq(s, h.Context.Log)
    }

	offset, err := strconv.Atoi(r.FormValue("minutes"))
	if err != nil {
		offset = 0
	}

	name, point := jingleName(r.FormValue("relative_to"), time.Duration(offset)*time.Minute)
	lib.NewJingle(
		name,
		s,
		time.Duration(offset)*time.Minute,
		point,
	)
}

// Play - will play song
func (h *JingleHandler) Play(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	f, err := h.Context.Songs.Find(id)
	s, ok := f.(SongI)
	if err != nil || !ok {
		h.Context.Log.Error(err.Error())
		http.NotFound(w, r)
		return
	}
	s.Play()

	output, _ := json.Marshal(h.Context.Songs.GetAll())
	w.Write(output)
}

// List - will list all actually playing songs
func (h *JingleHandler) List(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	list := h.Context.Songs.GetAll()

	output, _ := json.Marshal(list)
	w.Write(output)
}

// Delete - will remove song
func (h *JingleHandler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	f, err := h.Context.Songs.Find(id)
	s, ok := f.(SongI)
	if err != nil || !ok {
		http.NotFound(w, r)
		return
	}

	if s.IsPlaying() {
		s.Stop()
	}

	h.Context.Songs.Delete(id)

	output, _ := json.Marshal(h.Context.Songs.GetAll())
	w.Write(output)
}

func jingleName(relativeTo string, dur time.Duration) (out string, point lib.MatchPoint) {
	switch relativeTo {
	case "before_start":
		point = lib.MatchStart
		if dur.Minutes() == 0 {
			out = "začátek zápasu"
		} else {
			out = fmt.Sprintf("%v minut před začátkem", dur.Minutes())
		}
	case "before_end":
		point = lib.MatchEnd
		if dur.Minutes() == 0 {
			out = "konec zápasu"
		} else {
			out = fmt.Sprintf("%v minut do konce", dur.Minutes())
		}
	}
	return
}
