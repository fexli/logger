//go:build windows
// +build windows

package logger

import (
	"os"
	"syscall"
)

const kernel32 = "kernel32.dll"

func initErr() {
	if GlobalFileHandler == nil {
		return
	}
	kernal := syscall.NewLazyDLL(kernel32)
	setStdHandle := kernal.NewProc("SetStdHandle")
	sh := syscall.STD_ERROR_HANDLE
	if v, _, err := setStdHandle.Call(uintptr(sh), uintptr(GlobalFileHandler.Fd())); v == 0 {
		println("SetStdOutHandle failed:", err.Error())
		os.Exit(0xD00)
	}
	_ = syscall.FreeLibrary(syscall.Handle(kernal.Handle()))
}
