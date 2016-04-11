package server

import (
	"fmt"
	"syscall"
	"unsafe"
)

func Test() {
	var mod = syscall.NewLazyDLL("nircmd.dll")
	var proc = mod.NewProc("DoNirCmd")

	ret, _, _ := proc.Call(
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("mutesysvolume 1"))),
	)
	fmt.Printf("Return: %d\n", ret)

}
