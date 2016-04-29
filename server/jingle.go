package server

import (
	"encoding/json"
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

// Add - will add jingle
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
	s, err := lib.NewSong(filename, h.Context)
	if err != nil {
		h.Context.Log.Error(err.Error())
	} else {
		h.Context.Songs.AddUniq(s, h.Context.Log)
	}
	var (
		name   string
		point  lib.MatchPoint
		offset int
	)

	if r.FormValue("play") == "match_related" {
		offset, err = strconv.Atoi(r.FormValue("minutes"))
		if err != nil {
			h.Context.Log.Error(err.Error())
			offset = 0
		}

		name, point = lib.JingleName(r.FormValue("relative_to"), time.Duration(offset)*time.Minute)
	} else {
		name = filename
		point = lib.MatchNone
		offset = 0
	}

	j := lib.NewJingle(
		name,
		s,
		time.Duration(offset)*time.Minute,
		point,
		h.Context,
	)

	h.Context.Jingles.AddUniq(j, h.Context.Log)

	http.Redirect(w, r, "/", 302)
}

// List - will list all jingles
func (h *JingleHandler) List(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	list := h.Context.Jingles.GetAll()

	output, _ := json.Marshal(list)
	w.Write(output)
}
