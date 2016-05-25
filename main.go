package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"github.com/martin-reznik/jinglemanager/lib"
	"github.com/martin-reznik/jinglemanager/router"
	"github.com/martin-reznik/jinglemanager/server"
	"github.com/martin-reznik/logger"
	"github.com/skratchdot/open-golang/open"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"
)

func main() {
	flagOpenBrowser := flag.Bool("browser", false, "do open browser")
	flag.Parse()
	log := logger.NewLog(func(line *logger.LogLine) {
		lib.ChannelLog.Emit(lib.EventTypeLog, line)
		line.Print()
	}, &logger.Config{GoRoutinesLogTicker: 5 * time.Second})

	log.LogSeverity[logger.DEBUG] = true

	Ctx := lib.NewContext(log)
	Ctx.LoadCurrent()

	defer Ctx.AppClosed()

	/*
		Ctx.Tournament.AddMatchSlot(lib.NewMatchSlot(
			time.Now().Add(30*time.Second),
			5*time.Minute,
			Ctx,
		))
	*/

	httpHandler := server.HTTPHandler{Context: Ctx}
	fileHandler := server.FileProxyHandler{Context: Ctx}
	playerHandler := server.PlayerHandler{Context: Ctx}
	jingleHandler := server.JingleHandler{Context: Ctx}
	controlHandler := server.SoundControlHandler{Context: Ctx}
	slotHandler := server.SlotHandler{Context: Ctx}
	storageHandler := server.StorageHandler{Context: Ctx}
	socketHandler := server.SocketHandler{Context: Ctx, Upgrader: &websocket.Upgrader{}}

	web := router.NewRouter(log)

	web.AddMiddleware(router.NewAuthMiddleware(log))

	web.GET("/", httpHandler.Index)
	web.GET("/start", httpHandler.Start)
	web.POST("/tournament/new", httpHandler.NewTournament)

	web.GET("/css/*filepath", fileHandler.Static)
	web.GET("/js/*filepath", fileHandler.Static)
	web.GET("/images/*filepath", fileHandler.Static)

	web.POST("/track/add", playerHandler.Add)
	web.GET("/track/list", playerHandler.List)
	web.POST("/track/play/:id", playerHandler.Play)
	web.POST("/track/stop/:id", playerHandler.Stop)
	web.POST("/track/pause/:id", playerHandler.Pause)
	web.DELETE("/track/delete/:id", playerHandler.Delete)

	web.POST("/jingle/add", jingleHandler.Add)
	web.GET("/jingle/list", jingleHandler.List)

	web.POST("/app/mute", controlHandler.Mute)
	web.POST("/app/unmute", controlHandler.UnMute)
	web.POST("/app/add", controlHandler.Add)
	web.DELETE("/app/delete/:id", controlHandler.Delete)
	web.GET("/app/list", controlHandler.List)

	web.POST("/slot/add", slotHandler.Add)
	web.GET("/slot/list", slotHandler.List)

	web.GET("/save", storageHandler.Save)
	web.POST("/load", storageHandler.Load)
	web.POST("/open/:filename", storageHandler.Open)

	web.GET("/changes", socketHandler.HandleChangeSocket)
	web.GET("/logs", socketHandler.HandleLogSocket)

	wg := sync.WaitGroup{}
	go func() {
		// running server
		http.ListenAndServe(":8080", web)
	}()

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-ticker.C:
				log.Info("GC run")
				runtime.GC()
			}
		}
	}()

	wg.Add(1)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		// listening for interrupt to save progress
		defer wg.Done()
		for signal := range c {
			switch signal {
			case os.Interrupt:
				return
			}
		}
	}()

	log.Info("Server is up and running, open 'http://localhost:8080' in your browser")

	if *flagOpenBrowser {
		open.Start("http://localhost:8080")
	}
	wg.Wait()

	log.Info("Done")
}
