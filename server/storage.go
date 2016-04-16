package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/martin-reznik/jinglemanager/lib"
	"net/http"
)

// StorageHandler - storage handler
type StorageHandler struct {
	Logger       LogI
	SongList     *lib.FileList
	SoundControl *lib.SoundController
}

// Save - will save data
func (s *StorageHandler) Save(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := ps.ByName("name")
	out := lib.Save(s.Logger, s.SongList, s.SoundControl, name)
	w.Header().Set("Content-type", "application/octet-stream")
	w.Header().Set("Content-disposition", "attachment; filename="+name+".yml")
	w.Write(out)
}

// Load - will load data from specified file
func (s *StorageHandler) Load(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	lib.Load(s.Logger, s.SongList, s.SoundControl, ps.ByName("name"))
}
