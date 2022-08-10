//go:build aix || darwin || dragonfly || freebsd || netbsd || hurd || ios || js || (linux && !arm && !arm64) || nacl || plan9 || solaris || zos
// +build aix darwin dragonfly freebsd netbsd hurd ios js linux,!arm,!arm64 nacl plan9 solaris zos

package logger

import "runtime"

func initErr() {
	println("SetStdOutHandle failed: " + runtime.GOARCH + " architecture with " + runtime.GOOS + " system is not supported")
}
