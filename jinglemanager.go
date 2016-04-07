package main

import (
	"fmt"
	"github.com/martin-reznik/logger"
	"io/ioutil"
	"net/http"
)

// Index - index handler
type Index struct {
	logger *logger.Log
}

// ServeHTTP - will serve HTTP
func (i Index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	path := "static/html/index.html"
	data, err := ioutil.ReadFile(path)
	if err != nil {
		i.logger.Error(err.Error())
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, string(data))
}

// Static - handler for static content
type Static struct {
	logger *logger.Log
}

func (s Static) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := fmt.Sprintf("static%v", r.URL.Path)
	s.logger.Debug(fmt.Sprintf("Path to fetch: %v", path))
	data, err := ioutil.ReadFile(path)
	if err != nil {
		s.logger.Error(err.Error())
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, string(data))
}
