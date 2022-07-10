package logcolor

import (
	"io"
	"strconv"
	"strings"
)

const (
	defaultColor16Length = 8 // \x1b[3x;4xm
	highLightExtra       = 60
	textBase             = 3
	textLightBase        = 9
	backgroundBase       = 4
	backgroundLightBase  = 10
)
const (
	TextBlack BasicColorMask = iota + 30
	TextRed
	TextGreen
	TextYellow
	TextBlue
	TextMagenta
	TextCyan
	TextWhite
)
const (
	TextDarkGray BasicColorMask = iota + 90
	TextLightRed
	TextLightGreen
	TextLightYellow
	TextLightBlue
	TextLightMagenta
	TextLightCyan
	TextLightWhite
)

const (
	BgBlack BasicColorMask = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

const (
	BgDarkGray BasicColorMask = iota + 100
	BgLightRed
	BgLightGreen
	BgLightYellow
	BgLightBlue
	BgLightMagenta
	BgLightCyan
	BgLightWhite
)

// BasicColorMask 是最基本的颜色控制字符串
type BasicColorMask uint8

func (c BasicColorMask) Code() string {
	if c < TextBlack || c > TextLightWhite {
		return ""
	}
	return strconv.Itoa(int(c))
}

func (c BasicColorMask) String() string {
	if r := c.Code(); len(r) == 0 {
		return ""
	} else {
		return startCtr + r + endCtrl
	}
}

func (c BasicColorMask) Write(writer io.StringWriter) {
	writer.WriteString(startCtr + c.Code() + endCtrl)
}

func (c BasicColorMask) IsEmpty() bool {
	return c == 0
}

func (c BasicColorMask) ToC16() ColorMask {
	return &BasicColorIdentity{c}
}

func (c BasicColorMask) ToC256() ColorMask {
	identity := BasicColorIdentity{c}
	return identity.ToC256()
}

func (c BasicColorMask) ToCRGB() ColorMask {
	identity := BasicColorIdentity{c}
	return identity.ToCRGB()
}

func (c BasicColorMask) MergeFrom(older ColorMask) ColorMask {
	identity := BasicColorIdentity{c}
	return identity.MergeFrom(older)
}

////////////////////////////////////////////////////////////////////////////////

// BasicColorIdentity 基础颜色定义前景色和背景色
type BasicColorIdentity [2]BasicColorMask

func newBasic() BasicColorIdentity {
	return BasicColorIdentity{}
}

func (b *BasicColorIdentity) Code() string {
	r := make([]string, 0)
	for _, opt := range b {
		if opt == 0 {
			continue
		}
		r = append(r, strconv.Itoa(int(opt)))
	}
	return strings.Join(r, ";")
}

func (b *BasicColorIdentity) String() string {
	if r := b.Code(); len(r) == 0 {
		return ""
	} else {
		return startCtr + r + endCtrl
	}
}

func (b *BasicColorIdentity) Write(writer io.StringWriter) {
	if r := b.Code(); len(r) == 0 {
		return
	} else {
		writer.WriteString(startCtr + r + endCtrl)
	}
}

func (b *BasicColorIdentity) IsEmpty() bool {
	return b[0] == 0 && b[1] == 0
}

func (b *BasicColorIdentity) ToC16() ColorMask {
	return b
}
func (b *BasicColorIdentity) ToC256() ColorMask {
	h := newHundred()
	for _, mask := range *b {
		switch mask / 10 {
		case textBase:
			h[0] = [2]BasicColorMask{mask % 10, AsTx}
		case backgroundBase:
			h[1] = [2]BasicColorMask{mask % 10, AsBg}
		case textLightBase:
			h[0] = [2]BasicColorMask{mask%10 + 8, AsTx}
		case backgroundLightBase:
			h[1] = [2]BasicColorMask{mask%10 + 8, AsBg}
		}
	}
	return &h
}
func (b *BasicColorIdentity) ToCRGB() ColorMask {
	r := newRgb()
	for _, mask := range *b {
		col := mask % 10
		switch mask / 10 {
		case textBase:
			r[0] = [4]BasicColorMask{128 * ((col) % 1), 128 * ((col) % 2), 128 * ((col) % 4), AsTx}
		case backgroundBase:
			r[1] = [4]BasicColorMask{128 * ((col) % 1), 128 * ((col) % 2), 128 * ((col) % 4), AsBg}
		case textLightBase:
			r[0] = [4]BasicColorMask{255 * ((col) % 1), 255 * ((col) % 2), 255 * ((col) % 4), AsTx}
		case backgroundLightBase:
			r[1] = [4]BasicColorMask{255 * ((col) % 1), 255 * ((col) % 2), 255 * ((col) % 4), AsBg}
		}
	}
	return &r
}

func (b *BasicColorIdentity) MergeFrom(older ColorMask) ColorMask {
	if older == nil {
		return b
	}
	if b == nil || len(b) == 0 {
		return older
	}
	older = older.ToC16()
	newer := newBasic()
	for _, mask := range older.(*BasicColorIdentity) {
		switch mask / 10 {
		case textBase, textLightBase:
			newer[0] = mask
		case backgroundBase, backgroundLightBase:
			newer[1] = mask
		}
	}
	for _, mask := range b {
		switch mask / 10 {
		case textBase, textLightBase:
			newer[0] = mask
		case backgroundBase, backgroundLightBase:
			newer[1] = mask
		}
	}
	return &newer
}
