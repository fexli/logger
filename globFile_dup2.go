//go:build !windows && !(linux && (arm64 || arm)) && !aix && !darwin && !dragonfly && !freebsd && !netbsd && !hurd && !ios && !js && !linux && !nacl && !plan9 && !solaris && !zos
// +build !windows
// +build !linux !arm64,!arm
// +build !aix
// +build !darwin
// +build !dragonfly
// +build !freebsd
// +build !netbsd
// +build !hurd
// +build !ios
// +build !js
// +build !linux
// +build !nacl
// +build !plan9
// +build !solaris
// +build !zos

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
