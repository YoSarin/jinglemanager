package server

import (
	"fmt"
	"github.com/martin-reznik/logger"
	"html/template"
	"io/ioutil"
	"net/http"
)

// HTTPHandler - handler for serving httpPages
type HTTPHandler struct {
	Logger *logger.Log
}

// IndexData - nope yet
type IndexData struct {
	Title string
	Body  string
}

// Index - will serve index page
func (i *HTTPHandler) Index(w http.ResponseWriter, r *http.Request) {
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

// FileProxyHandler - struct handling file returns from server
type FileProxyHandler struct {
	Logger *logger.Log
}

// Static - handler for static content
func (f *FileProxyHandler) Static(w http.ResponseWriter, r *http.Request) {
	path := fmt.Sprintf("static%v", r.URL.Path)
	f.Logger.Debug(fmt.Sprintf("Path to fetch: %v", path))
	data, err := ioutil.ReadFile(path)
	if err != nil {
		f.Logger.Error(err.Error())
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, string(data))
}
