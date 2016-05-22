// +build windows

package lib

import (
	"fmt"
	"os/exec"
)

func (a *App) setAppVolume(level float32) bool {
	cmd := exec.Command("nircmdc.exe", "setappvolume", a.Name, fmt.Sprintf("%.2f", level))
	err := cmd.Start()
	if err != nil {
		return false
	}
	err = cmd.Wait()
	if err != nil {
		return false
	}
	return true
}

func (a *App) platformSpecificStuff() {
	// Nothing to do, bcs we have nircmd
}

func (a *App) refresh() {
	// nope
}
