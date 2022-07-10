package logcolor

import (
	"io"
	"os"
	"regexp"

	"github.com/xo/terminfo"
)

const (
	// ColorCodeRegExp 颜色代码正则，用于匹配颜色代码
	ColorCodeRegExp = `\033\[[\d;?]+m`
)

var (
	// EnableColor 切换是否打开颜色渲染
	EnableColor = !(os.Getenv("TERM") == "dumb") && !(os.Getenv("NO_COLOR") == "1")
	// the color support level for current terminal
	// needVTP - need enable VTP, only for windows OS
	colorLevel, needVTP = detectTermColorLevel()
	// match color codes
	codeRegex = regexp.MustCompile(ColorCodeRegExp)
)

// TermColorLevel value on current ENV
func TermColorLevel() terminfo.ColorLevel {
	return colorLevel
}

// SupportColor on the current ENV
func SupportColor() bool {
	return colorLevel > terminfo.ColorLevelNone
}

// Support256Color on the current ENV
func Support256Color() bool {
	return colorLevel > terminfo.ColorLevelBasic
}

// SupportTrueColor on the current ENV
func SupportTrueColor() bool {
	return colorLevel > terminfo.ColorLevelHundreds
}

// ForceSetColorLevel force open color render
func ForceSetColorLevel(level terminfo.ColorLevel) terminfo.ColorLevel {
	oldLevelVal := colorLevel
	colorLevel = level
	return oldLevelVal
}

// ForceColor force open color render
func ForceColor() terminfo.ColorLevel {
	return ForceSetColor(terminfo.ColorLevelMillions)
}

// ForceSetColor force open color render
func ForceSetColor(color terminfo.ColorLevel) terminfo.ColorLevel {
	return ForceSetColorLevel(color)
}

// ClearColorCode 清除文本中颜色控制符
func ClearColorCode(str string) string {
	return codeRegex.ReplaceAllString(str, "")
}

// SetTitle 设置终端标题
func SetTitle(std io.StringWriter, title string) {
	std.WriteString("\033]0;" + title + "")
}
