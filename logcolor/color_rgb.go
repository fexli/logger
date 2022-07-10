package logcolor

import (
	"io"
	"strconv"
	"strings"
)

const (
	BgRgbCodePrefix = "48;2;"
	TxRgbCodePrefix = "38;2;"
)

////////////////////////////////////////////////////////////////////////////////

type RGBColorIdentity [2][4]BasicColorMask

func newRgb() RGBColorIdentity {
	return RGBColorIdentity{}
}

func (c *RGBColorIdentity) Code() string {
	r := make([]string, 0, len(c))
	for _, mask := range c {
		if mask[0] == 0 && mask[1] == 0 && mask[2] == 0 {
			continue
		}
		switch mask[3] {
		case AsTx:
			r = append(r, TxRgbCodePrefix+strconv.Itoa(int(mask[0]))+";"+strconv.Itoa(int(mask[1]))+";"+strconv.Itoa(int(mask[2])))
		case AsBg:
			r = append(r, BgRgbCodePrefix+strconv.Itoa(int(mask[0]))+";"+strconv.Itoa(int(mask[1]))+";"+strconv.Itoa(int(mask[2])))
		}
	}
	return strings.Join(r, ";")
}

func (c *RGBColorIdentity) String() string {
	if r := c.Code(); len(r) == 0 {
		return ""
	} else {
		return startCtr + r + endCtrl
	}
}

func (c *RGBColorIdentity) Write(writer io.StringWriter) {
	writer.WriteString(startCtr + c.Code() + endCtrl)
}

func (c *RGBColorIdentity) IsEmpty() bool {
	return c == nil || len(c) == 0 || (c[0][3] == 0 && c[1][3] == 0)
}

func (c *RGBColorIdentity) ToC16() ColorMask {
	r := newBasic()
	for _, masks := range c {
		switch masks[3] {
		case AsTx:
			r[0] = cvt256to16(masks[0]) + TextBlack
		case AsBg:
			r[1] = cvt256to16(masks[0]) + BgBlack
		}
	}
	return &r
}
func (c *RGBColorIdentity) ToC256() ColorMask {
	r := newHundred()
	for _, masks := range c {
		switch masks[3] {
		case AsTx:
			r[0] = [2]BasicColorMask{cvtrgbto256(masks[0], masks[1], masks[2]), AsTx}
		case AsBg:
			r[1] = [2]BasicColorMask{cvtrgbto256(masks[0], masks[1], masks[2]), AsBg}
		}
	}
	return &r
}
func (c *RGBColorIdentity) ToCRGB() ColorMask {
	return c
}

func (c *RGBColorIdentity) MergeFrom(older ColorMask) ColorMask {
	if older == nil {
		return c
	}
	if c == nil || len(c) == 0 {
		return older
	}
	older = older.ToCRGB()
	newer := newRgb()
	for _, mask := range older.(*RGBColorIdentity) {
		switch mask[3] {
		case AsTx:
			newer[0] = [4]BasicColorMask{mask[0], mask[1], mask[2], AsTx}
		case AsBg:
			newer[1] = [4]BasicColorMask{mask[0], mask[1], mask[2], AsBg}
		}
	}
	for _, mask := range c {
		switch mask[3] {
		case AsTx:
			newer[0] = [4]BasicColorMask{mask[0], mask[1], mask[2], AsTx}
		case AsBg:
			newer[1] = [4]BasicColorMask{mask[0], mask[1], mask[2], AsBg}
		}
	}
	return &newer
}

func RGB(r, g, b uint8, isBg ...bool) ColorMask {
	if len(isBg) != 0 && isBg[0] {
		return &RGBColorIdentity{1: {BasicColorMask(r), BasicColorMask(g), BasicColorMask(b), AsBg}}
	} else {
		return &RGBColorIdentity{0: {BasicColorMask(r), BasicColorMask(g), BasicColorMask(b), AsTx}}
	}
}

func NewRGBColor(r, g, b uint8, isBg ...bool) *Color {
	return &Color{
		Identity: RGB(r, g, b, isBg...),
	}
}

func (c *RGBColorIdentity) And(r, g, b uint8, isBg ...bool) *RGBColorIdentity {
	if len(isBg) != 0 && isBg[0] {
		c[1] = [4]BasicColorMask{BasicColorMask(r), BasicColorMask(g), BasicColorMask(b), AsBg}
	} else {
		c[0] = [4]BasicColorMask{BasicColorMask(r), BasicColorMask(g), BasicColorMask(b), AsTx}
	}
	return c
}
