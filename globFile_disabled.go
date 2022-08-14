//go:build !(aix && ppc64) && !(darwin && amd64) && !(darwin && arm64) && !(dragonfly && amd64) && !(freebsd && 386) && !(freebsd && amd64) && !(freebsd && arm) && !(freebsd && arm64) && !(linux && 386) && !(linux && amd64) && !(linux && mips) && !(linux && mips64) && !(linux && mips64le) && !(linux && mipsle) && !(linux && ppc64) && !(linux && ppc64le) && !(linux && s390x) && !(netbsd && 386) && !(netbsd && amd64) && !(netbsd && arm) && !(netbsd && arm64) && !(openbsd && 386) && !(openbsd && amd64) && !(openbsd && arm) && !(openbsd && arm64) && !(openbsd && mips64) && !(linux && arm) && !(linux && arm64) && !windows
// +build !aix !ppc64
// +build !darwin !amd64
// +build !darwin !arm64
// +build !dragonfly !amd64
// +build !freebsd !386
// +build !freebsd !amd64
// +build !freebsd !arm
// +build !freebsd !arm64
// +build !linux !386
// +build !linux !amd64
// +build !linux !mips
// +build !linux !mips64
// +build !linux !mips64le
// +build !linux !mipsle
// +build !linux !ppc64
// +build !linux !ppc64le
// +build !linux !s390x
// +build !netbsd !386
// +build !netbsd !amd64
// +build !netbsd !arm
// +build !netbsd !arm64
// +build !openbsd !386
// +build !openbsd !amd64
// +build !openbsd !arm
// +build !openbsd !arm64
// +build !openbsd !mips64
// +build !linux !arm
// +build !linux !arm64
// +build !windows

package logger

import "runtime"

func initErr() {
	println("SetStdOutHandle failed: " + runtime.GOARCH + " architecture with " + runtime.GOOS + " system is not supported")
}
