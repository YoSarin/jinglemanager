package server

import (
	"fmt"
	"github.com/martin-reznik/logger"
	"html/template"
	"io/ioutil"
	"net/http"
)

// Index - index handler
type Index struct {
	Logger *logger.Log
}

type IndexData struct {
	Title string
	Body  string
}

// ServeHTTP - will serve HTTP
func (i Index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	t, err := template.ParseFiles("static/html/index.html")
	if err != nil {
		i.Logger.Error(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	t.Execute(w, &IndexData{"Titulek", "Nazdar"})
}

// Static - handler for static content
type Static struct {
	Logger *logger.Log
}

func (s Static) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := fmt.Sprintf("static%v", r.URL.Path)
	s.Logger.Debug(fmt.Sprintf("Path to fetch: %v", path))
	data, err := ioutil.ReadFile(path)
	if err != nil {
		s.Logger.Error(err.Error())
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, string(data))
}
