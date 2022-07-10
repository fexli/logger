package logcolor

import (
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"github.com/xo/terminfo"
)

// DetectColorLevel 检测当前终端的颜色支持级别，返回`terminfo.ColorLevel`作为支持级别
// 支持等级共分为3级：
//  - `terminfo.ColorLevelNone`：不支持颜色
//  - `terminfo.ColorLevelBasic`：支持基本颜色（16级）
//  - `terminfo.ColorLevelHundreds`：支持百级颜色（256级）
//  - `terminfo.ColorLevelMillions`：支持千万颜色（真彩）
func DetectColorLevel() terminfo.ColorLevel {
	level, _ := detectTermColorLevel()
	return level
}

// detect terminal color support level
//
// refer https://github.com/Delta456/box-cli-maker
func detectTermColorLevel() (level terminfo.ColorLevel, needVTP bool) {
	if val := os.Getenv("WSL_DISTRO_NAME"); val != "" {
		// WSL support true-color
		if detectWSL() {
			return terminfo.ColorLevelMillions, false
		}
	}

	isWin := runtime.GOOS == "windows"
	termVal := os.Getenv("TERM")

	// on TERM=screen: not support true-color
	if termVal != "screen" {
		// JetBrains IDE: JediTerm Support True Color
		if val := os.Getenv("TERMINAL_EMULATOR"); val == "JetBrains-JediTerm" {
			return terminfo.ColorLevelMillions, isWin
		}
	}

	// level, err = terminfo.ColorLevelFromEnv()
	level = detectColorLevelFromEnv(termVal, isWin)

	// fallback: simple detect by TERM value string.
	if level == terminfo.ColorLevelNone {
		// on Windows: enable VTP as it has True Color support
		level, needVTP = detectSpecialTermColor(termVal)
	}
	return
}

// detectColorLevelFromEnv 在`terminfo.ColorLevelFromEnv`的基础上，增加了部分检测
func detectColorLevelFromEnv(termVal string, isWin bool) terminfo.ColorLevel {
	colorTerm, termProg, forceColor := os.Getenv("COLORTERM"), os.Getenv("TERM_PROGRAM"), os.Getenv("FORCE_COLOR")
	switch {
	case strings.Contains(colorTerm, "truecolor") || strings.Contains(colorTerm, "24bit"):
		if termVal == "screen" { // on TERM=screen: not support true-color
			return terminfo.ColorLevelHundreds
		}
		return terminfo.ColorLevelMillions
	case colorTerm != "" || forceColor != "":
		return terminfo.ColorLevelBasic
	case termProg == "Apple_Terminal":
		return terminfo.ColorLevelHundreds
	case termProg == "Terminus" || termProg == "Hyper":
		if termVal == "screen" { // on TERM=screen: not support true-color
			return terminfo.ColorLevelHundreds
		}
		return terminfo.ColorLevelMillions
	case termProg == "iTerm.app":
		if termVal == "screen" { // on TERM=screen: not support true-color
			return terminfo.ColorLevelHundreds
		}
		// check iTerm version
		ver := os.Getenv("TERM_PROGRAM_VERSION")
		if ver != "" {
			i, err := strconv.Atoi(strings.Split(ver, ".")[0])
			if err != nil {
				return terminfo.ColorLevelHundreds
			}
			if i == 3 {
				return terminfo.ColorLevelMillions
			}
		}
		return terminfo.ColorLevelHundreds
	}
	if !isWin && termVal != "" {
		ti, err := terminfo.Load(termVal)
		if err != nil {
			return terminfo.ColorLevelNone
		}
		v, ok := ti.Nums[terminfo.MaxColors]
		switch {
		case !ok || v <= 16:
			return terminfo.ColorLevelNone
		case ok && v >= 256:
			return terminfo.ColorLevelHundreds
		}
		return terminfo.ColorLevelBasic
	}
	return terminfo.ColorLevelNone
}

var detectedWSL bool
var wslContents string

// detectWSL 检测当前系统是否是WSL环境
func detectWSL() bool {
	if !detectedWSL {
		detectedWSL = true

		b := make([]byte, 1024)
		f, err := os.Open("/proc/version")
		if err == nil {
			_, _ = f.Read(b) // ignore error
			if err = f.Close(); err != nil {
			}
			wslContents = string(b)
			return strings.Contains(wslContents, "Microsoft")
		}
	}
	return false
}

// isWSL 是否是WSL环境
func isWSL() bool {
	if val := os.Getenv("WSL_DISTRO_NAME"); val == "" {
		return false
	}
	wsl, err := ioutil.ReadFile("/proc/sys/kernel/osrelease")
	if err != nil {
		return false
	}
	content := strings.ToLower(string(wsl))
	return strings.Contains(content, "microsoft")
}

// IsWindows 是否是Windows系统
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// IsConsole 是否是控制台终端
func IsConsole(w io.Writer) bool {
	o, ok := w.(*os.File)
	if !ok {
		return false
	}
	fd := o.Fd()
	return fd == uintptr(syscall.Stdout) || fd == uintptr(syscall.Stdin) || fd == uintptr(syscall.Stderr)
}

// IsMSys 是否是msys(MINGW64)环境
func IsMSys() bool {
	// like "MSYSTEM=MINGW64"
	if len(os.Getenv("MSYSTEM")) > 0 {
		return true
	}

	return false
}

// IsSupportColor 当前终端是否支持任意颜色
func IsSupportColor() bool {
	return IsSupport16Color()
}

// IsSupport16Color 当前终端是否支持16级颜色
func IsSupport16Color() bool {
	level, _ := detectTermColorLevel()
	return level > terminfo.ColorLevelNone
}

// IsSupport256Color 当前终端是否支持256级颜色
func IsSupport256Color() bool {
	level, _ := detectTermColorLevel()
	return level > terminfo.ColorLevelBasic
}

// IsSupportRGBColor 当前终端是否支持RGB颜色，重定向自`IsSupportTrueColor`
func IsSupportRGBColor() bool {
	return IsSupportTrueColor()
}

// IsSupportTrueColor 当前终端是否支持TrueColor颜色
func IsSupportTrueColor() bool {
	level, _ := detectTermColorLevel()
	return level > terminfo.ColorLevelHundreds
}
