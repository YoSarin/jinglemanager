// +build linux unix

package lib

import (
	"fmt"
	"os/exec"
    "strings"
    "regexp"
    "errors"
)

type unixSpecificData struct {
    sink string
}

func (a *App) setAppVolume(level float32) bool {
    if a.specificData == nil {
        return false
    }
    s, ok := a.specificData.(*unixSpecificData)
    if !ok {
        return false
    }
	err := exec.Command("pactl", "set-sink-input-volume", s.sink, fmt.Sprintf("%v%%", int(100*level))).Run()
	if err != nil {
		return false
	}
	return true
}

func (a *App) platformSpecificStuff(){
    sink, err := a._findAppSink()
    if err == nil {
        s := &unixSpecificData{ sink: sink }
        a.specificData = s
    }
}

func (a *App) refresh() {
    a.platformSpecificStuff()
}

func (a *App) _findAppSink() (string, error) {
	out, _ := exec.Command("pactl", "list", "sink-inputs").Output()
	sinks := strings.Split(string(out), "Sink Input")
    appSearch := regexp.MustCompile(".*application\\.icon_name = \"" + a.Name + "\"")
    sinkSearch := regexp.MustCompile("^\\s+#(\\d+)\\s")
    for _, sinkData := range sinks {
        if "" != appSearch.FindString(sinkData) {
            return sinkSearch.FindStringSubmatch(sinkData)[1], nil
        }
    }
    return "", errors.New("App not found")
}
