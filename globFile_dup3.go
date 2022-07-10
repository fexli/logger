//go:build linux && (arm64 || arm)
// +build linux
// +build arm64 arm

package logger

import (
	"os"
	"syscall"
)

func initErr() {
	if GlobalFileHandler == nil {
		return
	}
	if err := syscall.Dup3(int(GlobalFileHandler.Fd()), int(os.Stderr.Fd()), 0); err != nil {
		println("SetStdOutHandle failed:", err.Error())
		os.Exit(0xD00)
	}
}
