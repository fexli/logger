//go:build aix || darwin || dragonfly || freebsd || netbsd || hurd || ios || js || (linux && !arm && !arm64) || nacl || plan9 || solaris || zos
// +build aix darwin dragonfly freebsd netbsd hurd ios js linux,!arm,!arm64 nacl plan9 solaris zos

package logger

func initErr() {
	println("SetStdOutHandle failed: armbe architecture is not supported")
}
