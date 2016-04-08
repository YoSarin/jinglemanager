package main

import (
	"github.com/martin-reznik/jinglemanager/manager"
	"github.com/martin-reznik/jinglemanager/server"
	"github.com/martin-reznik/logger"
	"net/http"
)

func main() {
	log := logger.NewLog(func(line *logger.LogLine) { line.Print() })
	log.LogSeverity[logger.DEBUG] = true
	defer log.Close()

	http.Handle("/", server.Index{Logger: log})

	http.Handle("/css/", server.Static{Logger: log})
	http.Handle("/js/", server.Static{Logger: log})
	http.Handle("/images/", server.Static{Logger: log})

	http.Handle("/play", manager.Player{Logger: log})

	log.Info("Server is up and running, open 'http://localhost:8080' in your browser")

	http.ListenAndServe(":8080", nil)
}
