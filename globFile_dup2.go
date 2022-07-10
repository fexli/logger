//go:build !windows && !(linux && (arm64 || arm)) && !aix && !darwin && !dragonfly && !freebsd && !netbsd && !hurd && !ios && !js && !linux && !nacl && !plan9 && !solaris && !zos

package logger

import (
	"os"
	"syscall"
)

func initErr() {
	if GlobalFileHandler == nil {
		return
	}
	if err := syscall.Dup2(int(GlobalFileHandler.Fd()), int(os.Stderr.Fd())); err != nil {
		println("SetStdOutHandle failed:", err.Error())
		os.Exit(0xD00)
	}
}
