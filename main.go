package main

import (
	"github.com/martin-reznik/jinglemanager/server"
	"github.com/martin-reznik/logger"
	"net/http"
)

func main() {
	log := logger.NewLog(func(line *logger.LogLine) { line.Print() })
	log.LogSeverity[logger.DEBUG] = true
	defer log.Close()

	httpHandler := server.HTTPHandler{Logger: log}
	fileHandler := server.FileProxyHandler{Logger: log}
	playerHandler := server.PlayerHandler{Logger: log}

	http.HandleFunc("/", httpHandler.Index)

	http.HandleFunc("/css/", fileHandler.Static)
	http.HandleFunc("/js/", fileHandler.Static)
	http.HandleFunc("/images/", fileHandler.Static)

	http.HandleFunc("/play", playerHandler.Play)
	http.HandleFunc("/stop", playerHandler.Stop)
	http.HandleFunc("/list", playerHandler.List)

	log.Info("Server is up and running, open 'http://localhost:8080' in your browser")

	http.ListenAndServe(":8080", nil)
}
