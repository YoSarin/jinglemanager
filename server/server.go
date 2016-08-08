package server

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/martin-reznik/jinglemanager/lib"
	"html/template"
	"net/http"
	"strconv"
	"time"
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
		http.Redirect(w, r, "/start", 302)
		return
	}
	if err != nil {
		i.Context.Log.Error(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	t.Execute(w, &IndexData{"Jingle Manager", i.Context.Tournament.Name})
}

// Start - will serve creator page
func (i *HTTPHandler) Start(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	t, err := template.ParseFiles("static/html/new_tournament.html")
	if err != nil {
		i.Context.Log.Error(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	t.Execute(w, &struct {
		Title          string
		TournamentName string
		List           []string
	}{
		Title:          "Jingle Manager",
		TournamentName: i.Context.Tournament.Name,
		List:           i.Context.ListTournaments(),
	})
}

// NewTournament - will serve index page
func (i *HTTPHandler) NewTournament(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := r.FormValue("name")
	loc := time.Now().Location()
	parsed, err := time.ParseInLocation("15:04", r.FormValue("start"), loc)
	if err != nil {
		i.Context.Log.Error("Fuj: ", err)
		http.Redirect(w, r, "/start", 302)
		return
	}
	n := time.Now()
	start := time.Date(n.Year(), n.Month(), n.Day(), parsed.Hour(), parsed.Minute(), parsed.Second(), parsed.Nanosecond(), parsed.Location())
	matchDuration, err := strconv.Atoi(r.FormValue("duration"))
	if err != nil {
		i.Context.Log.Error("Fuj: ", err)
		http.Redirect(w, r, "/start", 302)
		return
	}
	breakDuration, err := strconv.Atoi(r.FormValue("break"))
	if err != nil {
		i.Context.Log.Error("Fuj: ", err)
		http.Redirect(w, r, "/start", 302)
		return
	}
	count, err := strconv.Atoi(r.FormValue("count"))
	if err != nil {
		i.Context.Log.Error("Fuj: ", err)
		http.Redirect(w, r, "/start", 302)
		return
	}

	i.Context.Save()
	i.Context.NewTournament(name)

	for n := 0; n < count; n++ {
		matchStart := start.Add(time.Duration(n*(matchDuration+breakDuration)) * time.Minute)
		m := lib.NewMatchSlot(matchStart, time.Duration(matchDuration)*time.Minute, i.Context)
		i.Context.Tournament.AddMatchSlot(m)
	}
	http.Redirect(w, r, "/", 302)
}

func (i *HTTPHandler) Health(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	output, _ := json.Marshal(map[string]string{"status": "ok"})
	w.Write(output)
}

// FileProxyHandler - struct handling file returns from server
type FileProxyHandler struct {
	Context *lib.Context
}

// Static - handler for static content
func (f *FileProxyHandler) Static(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	path := fmt.Sprintf("static%v", r.URL.Path)
	http.ServeFile(w, r, path)
}
