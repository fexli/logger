package logcolor

import (
	"io"
	"strings"
)

// Color 基于ColorMask定义的任意颜色
type Color struct {
	Identity ColorMask    `json:"identity"`
	Options  ColorOptions `json:"options"`
}

func (b *Color) Code() string {
	if b == nil {
		return ""
	}
	r := make([]string, 0, 2)
	if co := b.Options.Code(); len(co) > 0 {
		r = append(r, co)
	}
	if b.Identity != nil {
		if id := b.Identity.Code(); len(id) > 0 {
			r = append(r, id)
		}
	}
	return strings.Join(r, ";")
}

func (b *Color) String() string {
	if r := b.Code(); len(r) == 0 {
		return ""
	} else {
		return startCtr + r + endCtrl
	}
}

func (b *Color) WriteStart(writer io.StringWriter) {
	writer.WriteString(b.String())
}

func (b *Color) WriteEnd(writer io.StringWriter) {
	if b == nil || (b.Options&opMask == 0 && (b.Identity == nil || b.Identity.IsEmpty())) {
		return
	}
	writer.WriteString(resetCtr)
}

func (b *Color) MergeTo(newer *Color) *Color {
	if newer == nil {
		return b
	}
	if b == nil {
		return newer
	}
	mergedColor := new(Color)
	if newer.Identity != nil {
		mergedColor.Identity = newer.Identity.MergeFrom(b.Identity)
	} else {
		mergedColor.Identity = b.Identity
	}
	mergedColor.Options = newer.Options.MergeFrom(b.Options)
	return mergedColor
}

func NewColor(color ColorMask, options ...ColorOptions) *Color {
	return &Color{
		Identity: color,
		Options:  NewColorOptions(options...),
	}
}

////////////////////////////////////////////////////////////////////////////////
// Basic Color Creator

func BlackString(info string) *LogTextCtx {
	return &LogTextCtx{
		Log:   info,
		Color: NewColor(&BasicColorIdentity{TextBlack}),
	}
}

func RedString(info string) *LogTextCtx {
	return &LogTextCtx{
		Log:   info,
		Color: NewColor(&BasicColorIdentity{TextRed}),
	}
}

func GreenString(info string) *LogTextCtx {
	return &LogTextCtx{
		Log:   info,
		Color: NewColor(&BasicColorIdentity{TextGreen}),
	}
}

func YellowString(info string) *LogTextCtx {
	return &LogTextCtx{
		Log:   info,
		Color: NewColor(&BasicColorIdentity{TextYellow}),
	}
}

func BlueString(info string) *LogTextCtx {
	return &LogTextCtx{
		Log:   info,
		Color: NewColor(&BasicColorIdentity{TextBlue}),
	}
}

func MagentaString(info string) *LogTextCtx {
	return &LogTextCtx{
		Log:   info,
		Color: NewColor(&BasicColorIdentity{TextMagenta}),
	}
}

func CyanString(info string) *LogTextCtx {
	return &LogTextCtx{
		Log:   info,
		Color: NewColor(&BasicColorIdentity{TextCyan}),
	}
}

func WhiteString(info string) *LogTextCtx {
	return &LogTextCtx{
		Log:   info,
		Color: NewColor(&BasicColorIdentity{TextWhite}),
	}
}

func DarkGrayString(info string) *LogTextCtx {
	return &LogTextCtx{
		Log:   info,
		Color: NewColor(&BasicColorIdentity{TextDarkGray}),
	}
}

func LightRedString(info string) *LogTextCtx {
	return &LogTextCtx{
		Log:   info,
		Color: NewColor(&BasicColorIdentity{TextLightRed}),
	}
}
func LightGreenString(info string) *LogTextCtx {
	return &LogTextCtx{
		Log:   info,
		Color: NewColor(&BasicColorIdentity{TextLightGreen}),
	}
}

func LightYellowString(info string) *LogTextCtx {
	return &LogTextCtx{
		Log:   info,
		Color: NewColor(&BasicColorIdentity{TextLightYellow}),
	}
}

func LightBlueString(info string) *LogTextCtx {
	return &LogTextCtx{
		Log:   info,
		Color: NewColor(&BasicColorIdentity{TextLightBlue}),
	}
}

func LightMagentaString(info string) *LogTextCtx {
	return &LogTextCtx{
		Log:   info,
		Color: NewColor(&BasicColorIdentity{TextLightMagenta}),
	}
}

func LightCyanString(info string) *LogTextCtx {
	return &LogTextCtx{
		Log:   info,
		Color: NewColor(&BasicColorIdentity{TextLightCyan}),
	}
}

func LightWhiteString(info string) *LogTextCtx {
	return &LogTextCtx{
		Log:   info,
		Color: NewColor(&BasicColorIdentity{TextLightWhite}),
	}
}
