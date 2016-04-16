package lib

import (
	"time"
)

// SoundController - contains info about sound controll we want to perform
type SoundController struct {
	Logger             LogI
	TargetSilentVolume int
	TargetLoudVolume   int
	steps              int
	appList            map[string]*App
}

// App - contains info about applications which sound should be manipulated
type App struct {
	Name         string
	Volume       float32
	specificData interface{}
}

// NewSoundController - will create new sound controller
func NewSoundController(logger LogI) *SoundController {
	return &SoundController{
		Logger:             logger,
		TargetLoudVolume:   100,
		TargetSilentVolume: 0,
		steps:              50,
		appList:            make(map[string]*App),
	}
}

// MuteApps - Will mute all apps in controller
func (c *SoundController) MuteApps() {
	c.refresh()
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

// UnMuteApps - Will unmute all apps in controller
func (c *SoundController) UnMuteApps() {
	c.refresh()
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

// Add - Will add an application to controller
func (c *SoundController) Add(appname string) {
	a := &App{appname, 1.0, nil}
	a.platformSpecificStuff()
	c.appList[appname] = a
}

// AddUniq - Will add an application to controller
func (c *SoundController) AddUniq(appname string, l LogI) (bool, error) {
	c.Add(appname)
	return true, nil
}

// AppNames - Will return list of app names
func (c *SoundController) AppNames() []string {
	out := make([]string, len(c.appList))
	i := 0
	for _, val := range c.appList {
		out[i] = val.Name
		i++
	}
	return out
}

// Remove - will remove an application from controller
func (c *SoundController) Remove(appname string) {
	delete(c.appList, appname)
}

// List - will return list of all applications in controller
func (c *SoundController) List() map[string]*App {
	return c.appList
}

func (c *SoundController) refresh() {
	for _, app := range c.appList {
		app.refresh()
	}
}
