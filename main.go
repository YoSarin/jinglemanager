package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/martin-reznik/jinglemanager/lib"
	"github.com/martin-reznik/jinglemanager/server"
	"github.com/martin-reznik/logger"
	"github.com/skratchdot/open-golang/open"
	"net/http"
	"sync"
)

func main() {

	log := logger.NewLog(func(line *logger.LogLine) { line.Print() })
	log.LogSeverity[logger.DEBUG] = true
	defer log.Close()

	httpHandler := server.HTTPHandler{Logger: log}
	fileHandler := server.FileProxyHandler{Logger: log}
	playerHandler := server.PlayerHandler{Logger: log, SongList: lib.NewFileList()}
	controlHandler := server.ControlHandler{Logger: log, Player: lib.NewController(log)}

	router := httprouter.New()
	router.GET("/", httpHandler.Index)

	router.GET("/css/*filepath", fileHandler.Static)
	router.GET("/js/*filepath", fileHandler.Static)
	router.GET("/images/*filepath", fileHandler.Static)

	router.POST("/track/add", playerHandler.Add)
	router.GET("/track/list", playerHandler.List)
	router.POST("/track/play/:id", playerHandler.Play)
	router.POST("/track/stop/:id", playerHandler.Stop)
	router.POST("/track/pause/:id", playerHandler.Pause)
	router.DELETE("/track/delete/:id", playerHandler.Delete)

	router.POST("/app/mute", controlHandler.Mute)
	router.POST("/app/unmute", controlHandler.UnMute)
	router.POST("/app/add/:app", controlHandler.Add)
	router.DELETE("/app/remove/:app", controlHandler.Delete)
	router.GET("/app/list", controlHandler.List)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		http.ListenAndServe(":8080", router)
	}()

	log.Info("Server is up and running, open 'http://localhost:8080' in your browser")

	open.Start("http://localhost:8080")
	wg.Wait()
}
