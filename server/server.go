package server

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/martin-reznik/jinglemanager/lib"
	"html/template"
	"io/ioutil"
	"net/http"
)

// HTTPHandler - handler for serving httpPages
type HTTPHandler struct {
	Context *lib.Context
}

// IndexData - nope yet
type IndexData struct {
	Title          string
	TournamentName string
}

// Index - will serve index page
func (i *HTTPHandler) Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	var t *template.Template
	var err error
	if i.Context.Tournament.Name != "" {
		t, err = template.ParseFiles("static/html/index.html")
	} else {
		t, err = template.ParseFiles("static/html/new_tournament.html")
	}
	if err != nil {
		i.Context.Log.Error(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	t.Execute(w, &IndexData{"Jingle Manager", i.Context.Tournament.Name})
}

// Start - will serve index page
func (i *HTTPHandler) Start(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	t, err := template.ParseFiles("static/html/new_tournament.html")
	if err != nil {
		i.Context.Log.Error(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	t.Execute(w, &IndexData{"Jingle Manager", i.Context.Tournament.Name})
}

// NewTournament - will serve index page
func (i *HTTPHandler) NewTournament(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := r.FormValue("name")
	i.Context.Save()
	i.Context.NewTournament(name)
	http.Redirect(w, r, "/", 302)
}

// FileProxyHandler - struct handling file returns from server
type FileProxyHandler struct {
	Context *lib.Context
}

// Static - handler for static content
func (f *FileProxyHandler) Static(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	path := fmt.Sprintf("static%v", r.URL.Path)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		f.Context.Log.Error(err.Error())
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, string(data))
}
