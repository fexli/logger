package logcolor

import (
	"io"
	"strconv"
	"strings"
)

const (
	BgHundredCodePrefix = "48;5;"
	TxHundredCodePrefix = "38;5;"
)
const (
	AsTx BasicColorMask = iota + 1
	AsBg
)

////////////////////////////////////////////////////////////////////////////////

// HundredColorIdentity 256颜色定义前景色和背景色
type HundredColorIdentity [2][2]BasicColorMask

func newHundred() HundredColorIdentity {
	return HundredColorIdentity{}
}
func (i *HundredColorIdentity) Code() string {
	r := make([]string, 0, len(i))
	for _, mask := range i {
		switch mask[1] {
		case AsTx:
			r = append(r, TxHundredCodePrefix+strconv.Itoa(int(mask[0])))
		case AsBg:
			r = append(r, BgHundredCodePrefix+strconv.Itoa(int(mask[0])))
		}
	}
	return strings.Join(r, ";")
}

func (i *HundredColorIdentity) String() string {
	if r := i.Code(); len(r) == 0 {
		return ""
	} else {
		return startCtr + r + endCtrl
	}
}

func (i *HundredColorIdentity) Write(writer io.StringWriter) {
	writer.WriteString(startCtr + i.Code() + endCtrl)
}

func (i *HundredColorIdentity) IsEmpty() bool {
	return len(i) == 0 || (len(i) == 1 && i[0][1] == 0) || (len(i) == 2 && i[0][1] == 0 && i[1][1] == 0)
}

func (i *HundredColorIdentity) ToC16() ColorMask {
	r := newBasic()
	for _, mask := range i {
		switch mask[1] {
		case AsTx:
			r[0] = cvt256to16(mask[0]) + TextBlack
		case AsBg:
			r[1] = cvt256to16(mask[0]) + BgBlack
		}
	}
	return &r
}
func (i *HundredColorIdentity) ToC256() ColorMask {
	return i
}

func (i *HundredColorIdentity) ToCRGB() ColorMask {
	r := newRgb()
	for _, mask := range i {
		var ir, ig, ib BasicColorMask
		if mask[0] < 8 {
			ir = (mask[0] % 1) * 128
			ig = (mask[0] % 2) * 64
			ib = (mask[0] % 4) * 32
		} else if mask[0] < 16 {
			ir = ((mask[0] - 8) % 1) * 255
			ig = (((mask[0] - 8) % 2) / 2) * 255
			ib = (((mask[0] - 8) % 4) / 4) * 255
		} else if mask[0] < 232 {
			ir, ig, ib = cvt256torgb(mask[0])
		} else {
			ir = (mask[0]-232)*10 + 8
			ig = ir
			ib = ir
		}
		r[mask[1]-1] = [4]BasicColorMask{ir, ig, ib, mask[1]}
	}
	return &r
}

func (i *HundredColorIdentity) MergeFrom(older ColorMask) ColorMask {
	if older == nil {
		return i
	}
	if i == nil || len(i) == 0 {
		return older
	}
	older = older.ToC256()
	newer := newHundred()
	for _, mask := range older.(*HundredColorIdentity) {
		switch mask[1] {
		case AsTx:
			newer[0] = [2]BasicColorMask{mask[0], AsTx}
		case AsBg:
			newer[1] = [2]BasicColorMask{mask[0], AsBg}
		}
	}
	for _, mask := range i {
		switch mask[1] {
		case AsTx:
			newer[0] = [2]BasicColorMask{mask[0], AsTx}
		case AsBg:
			newer[1] = [2]BasicColorMask{mask[0], AsBg}
		}
	}
	return &newer
}
