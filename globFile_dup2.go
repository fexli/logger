//go:build (aix && ppc64) || (darwin && amd64) || (darwin && arm64) || (dragonfly && amd64) || (freebsd && 386) || (freebsd && amd64) || (freebsd && arm) || (freebsd && arm64) || (linux && 386) || (linux && amd64) || (linux && mips) || (linux && mips64) || (linux && mips64le) || (linux && ppc64) || (linux && ppc64le) || (linux && s390x) || (netbsd && 386) || (netbsd && amd64) || (netbsd && arm) || (netbsd && arm64) || (openbsd && 386) || (openbsd && amd64) || (openbsd && arm) || (openbsd && arm64) || (openbsd && mips64)
// +build aix,ppc64 darwin,amd64 darwin,arm64 dragonfly,amd64 freebsd,386 freebsd,amd64 freebsd,arm freebsd,arm64 linux,386 linux,amd64 linux,mips linux,mips64 linux,mips64le linux,ppc64 linux,ppc64le linux,s390x netbsd,386 netbsd,amd64 netbsd,arm netbsd,arm64 openbsd,386 openbsd,amd64 openbsd,arm openbsd,arm64 openbsd,mips64

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
		println("SetStdOutHandle[2] failed:", err.Error())
		os.Exit(0xD00)
	}
}
