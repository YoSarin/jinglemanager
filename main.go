package main

import (
	"github.com/julienschmidt/httprouter"
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

	router := httprouter.New()
	router.GET("/", httpHandler.Index)

	router.GET("/css/*filepath", fileHandler.Static)
	router.GET("/js/*filepath", fileHandler.Static)
	router.GET("/images/*filepath", fileHandler.Static)

	router.POST("/track/play/:id", playerHandler.Play)
	router.POST("/track/add", playerHandler.Add)
	router.POST("/track/stop/:id", playerHandler.Stop)
	router.GET("/track/list", playerHandler.List)

	log.Info("Server is up and running, open 'http://localhost:8080' in your browser")

	http.ListenAndServe(":8080", router)
}
