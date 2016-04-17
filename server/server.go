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
	Title string
}

// Index - will serve index page
func (i *HTTPHandler) Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	t, err := template.ParseFiles("static/html/index.html")
	if err != nil {
		i.Context.Log.Error(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	t.Execute(w, &IndexData{"Jingle Manager"})
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
