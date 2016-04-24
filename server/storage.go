package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/martin-reznik/jinglemanager/lib"
	"io/ioutil"
	"net/http"
)

// StorageHandler - storage handler
type StorageHandler struct {
	Context *lib.Context
}

// Save - will save data
func (s *StorageHandler) Save(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := s.Context.Tournament.Name
	out := s.Context.Save()
	w.Header().Set("Content-type", "application/octet-stream")
	w.Header().Set("Content-disposition", "attachment; filename="+name+".yml")
	w.Write(out)
}

// Load - will load data from specified file
func (s *StorageHandler) Load(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	file, _, err := r.FormFile("file")
	if err != nil {
		s.Context.Log.Error("File upload failed: " + err.Error())
		http.Error(w, "upload failed", 500)
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)

	s.Context.Load(data)
	http.Redirect(w, r, "/", 302)
}
