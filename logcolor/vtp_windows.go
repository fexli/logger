//go:build windows
// +build windows

package logcolor

import (
	"os"
	"syscall"
	"unsafe"

	"github.com/xo/terminfo"
	"golang.org/x/sys/windows"
)

var (
	kernel32           *syscall.LazyDLL
	procGetConsoleMode *syscall.LazyProc
	procSetConsoleMode *syscall.LazyProc
)

// Get the Windows Version and Build Number
var (
	winVersion, _, buildNumber = windows.RtlGetNtVersionNumbers()
)

const (
	EnableVirtualTerminalProcessingMode uint32 = 0x0004
)

func init() {

	if !SupportColor() {
		return
	}
	if !EnableColor {
		return
	}
	tryEnableVTP(needVTP)
}

func tryEnableVTP(enable bool) bool {
	if !enable {
		return false
	}
	initKernel32Proc()
	return tryEnableOnCONOUT() || tryEnableOnStdout()
}

func initKernel32Proc() {
	if kernel32 != nil {
		return
	}
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

	procGetConsoleMode = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode = kernel32.NewProc("SetConsoleMode")
}

func tryEnableOnCONOUT() bool {
	outHandle, err := syscall.Open("CONOUT$", syscall.O_RDWR, 0)
	if err != nil {
		return false
	}
	return EnableVirtualTerminalProcessing(outHandle, true) == nil
}

func tryEnableOnStdout() bool {
	return EnableVirtualTerminalProcessing(syscall.Stdout, true) == nil
}

func detectSpecialTermColor(_ string) (tl terminfo.ColorLevel, needVTP bool) {
	if os.Getenv("ConEmuANSI") == "ON" {
		// ConEmuANSI is "ON" for generic ANSI support
		// but True Color option is enabled by default
		// I am just assuming that people wouldn't have disabled it
		// Even if it is not enabled then ConEmu will auto round off
		// accordingly
		return terminfo.ColorLevelMillions, false
	}
	// Windows10 Build 10586版本前不支持ANSI颜色
	if buildNumber < 10586 || winVersion < 10 {
		if os.Getenv("ANSICON") != "" {
			if conVersion := os.Getenv("ANSICON_VER"); conVersion >= "181" {
				return terminfo.ColorLevelHundreds, false
			}
			return terminfo.ColorLevelBasic, false
		}
		return terminfo.ColorLevelNone, false
	}
	// Windows10 Build 14931版本前只支持8bit颜色
	if buildNumber < 14931 {
		return terminfo.ColorLevelHundreds, true
	}
	return terminfo.ColorLevelMillions, true
}

// EnableVirtualTerminalProcessing 开启虚拟终端处理(virtual terminal processing)模式
//
// ref from github.com/konsorten/go-windows-terminal-sequences
//
// doc https://docs.microsoft.com/zh-cn/windows/console/console-virtual-terminal-sequences#samples
//
// e.g.:
// 	err := EnableVirtualTerminalProcessing(syscall.Stdout, true) // 开启VTP模式
func EnableVirtualTerminalProcessing(stream syscall.Handle, enable bool) error {
	var mode uint32
	err := syscall.GetConsoleMode(stream, &mode)
	if err != nil {
		// fmt.Println("EnableVirtualTerminalProcessing", err)
		return err
	}

	if enable {
		mode |= EnableVirtualTerminalProcessingMode
	} else {
		mode &^= EnableVirtualTerminalProcessingMode
	}
	ret, _, err := procSetConsoleMode.Call(uintptr(stream), uintptr(mode))
	if ret == 0 {
		return err
	}

	return nil
}

// IsTerminal 返回文件描述符是否为Tty终端
//
// e.g.:
// 	fd := os.Stdout.Fd()
// 	fd := uintptr(syscall.Stdout) // windows
// 	IsTerminal(fd)
func IsTerminal(fd uintptr) bool {
	initKernel32Proc()

	var st uint32
	r, _, e := syscall.Syscall(procGetConsoleMode.Addr(), 2, fd, uintptr(unsafe.Pointer(&st)), 0)
	return r != 0 && e == 0
}
