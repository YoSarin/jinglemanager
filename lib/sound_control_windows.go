// +build windows

package lib

import (
	"fmt"
	"os/exec"
)

func (a *App) setAppVolume(level float32) bool {
	err := exec.Command("nircmdc.exe", "setappvolume", a.Name, fmt.Sprintf("%.2f", level)).Run()
	if err != nil {
		return false
	}
	return true
}