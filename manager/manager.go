package manager

import (
	"github.com/martin-reznik/logger"
	"github.com/mattn/go-soundplayer"
	"net/http"
)

// Player - player handler
type Player struct {
	Logger *logger.Log
}

// ServeHTTP - will serve HTTP
func (p Player) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	err := soundplayer.Play("media/song.mp3")
	if err != nil {
		p.Logger.Error(err.Error())
	}
	err = soundplayer.Play("media/song.wav")
	if err != nil {
		p.Logger.Error(err.Error())
	}
}
