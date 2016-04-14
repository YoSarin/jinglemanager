// +build windows

package lib

import (
	"time"
)

type Controller struct {
	Logger             LogI
	TargetSilentVolume int
	TargetLoudVolume   int
	steps              int
	appList            map[string]*App
}

type App struct {
	Name   string
	Volume float32
}

func NewController(logger LogI) *Controller {
	return &Controller{
		Logger:             logger,
		TargetLoudVolume:   100,
		TargetSilentVolume: 0,
		steps:              50,
		appList:            make(map[string]*App),
	}
}

func (c *Controller) MuteApps() {
	for i := 0; i < c.steps; i++ {
		level := float32(float32(c.TargetLoudVolume)-(float32(i)*float32(c.TargetLoudVolume-c.TargetSilentVolume)/float32(c.steps))) / 100.0
		for _, app := range c.appList {
			go func(app *App) {
				app.setAppVolume(level)
				app.Volume = level
			}(app)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (c *Controller) UnMuteApps() {
	for i := 0; i < c.steps; i++ {
		level := float32(float32(c.TargetSilentVolume)+(float32(i)*float32(c.TargetLoudVolume-c.TargetSilentVolume)/float32(c.steps))) / 100.0
		for _, app := range c.appList {
			go func(app *App) {
				app.setAppVolume(level)
				app.Volume = level
			}(app)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (c *Controller) Add(appname string) {
	c.appList[appname] = &App{appname, 1.0}
}

func (c *Controller) Remove(appname string) {
	delete(c.appList, appname)
}

func (c *Controller) List() map[string]*App {
	return c.appList
}
