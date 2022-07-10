//go:build !windows
// +build !windows

package logcolor

import (
	"strings"
	"syscall"

	"github.com/xo/terminfo"
)

// detectSpecialTermColor 检测特殊的终端颜色支持
func detectSpecialTermColor(termVal string) (terminfo.ColorLevel, bool) {
	if termVal == "" {
		return terminfo.ColorLevelNone, false
	}
	if termVal == "screen" {
		return terminfo.ColorLevelHundreds, false
	}
	if strings.Contains(termVal, "256color") {
		return terminfo.ColorLevelHundreds, false
	}
	if strings.Contains(termVal, "xterm") {
		return terminfo.ColorLevelHundreds, false
	}
	return terminfo.ColorLevelBasic, false
}

// IsTerminal 检测指定描述符是否是终端
//
// e.g.:
// 	IsTerminal(os.Stdout.Fd())
func IsTerminal(fd uintptr) bool {
	return fd == uintptr(syscall.Stdout) || fd == uintptr(syscall.Stdin) || fd == uintptr(syscall.Stderr)
}
