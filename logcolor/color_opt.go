package logcolor

import (
	"io"
	"strings"
)

type ColorOptions uint32

const (
	OpBold       ColorOptions = 1 << iota // 加粗
	OpFaint                               // 模糊
	OpItalic                              // 斜体
	OpUnderline                           // 下划线
	OpBlinkSlow                           // 慢闪烁
	OpBlinkFast                           // 快闪烁
	OpInverse                             // 反显（反色）
	OpConceal                             // 隐藏（暗色）
	OpCrossedOut                          // 删除线

	OpNoBold
	OpNoFaint
	OpNoItalic
	OpNoUnderline
	OpNoBlinkSlow
	OpNoBlinkFast
	OpNoInverse
	OpNoConceal
	OpNoCrossedOut

	OpReset

	opMask   ColorOptions = OpBold | OpFaint | OpItalic | OpUnderline | OpBlinkSlow | OpBlinkFast | OpInverse | OpConceal | OpCrossedOut
	opNoMask ColorOptions = OpNoBold | OpNoFaint | OpNoItalic | OpNoUnderline | OpNoBlinkSlow | OpNoBlinkFast | OpNoInverse | OpNoConceal | OpNoCrossedOut
)

func (c ColorOptions) Code() string {
	if c == 0 {
		return ""
	}
	if c&OpReset != 0 {
		return "0"
	}
	r := make([]string, 0)
	if c&OpBold != 0 && c&OpNoBold == 0 {
		r = append(r, "1")
	}
	if c&OpFaint != 0 && c&OpNoFaint == 0 {
		r = append(r, "2")
	}
	if c&OpItalic != 0 && c&OpNoItalic == 0 {
		r = append(r, "3")
	}
	if c&OpUnderline != 0 && c&OpNoUnderline == 0 {
		r = append(r, "4")
	}
	if c&OpBlinkSlow != 0 && c&OpNoBlinkSlow == 0 {
		r = append(r, "5")
	}
	if c&OpBlinkFast != 0 && c&OpNoBlinkFast == 0 {
		r = append(r, "6")
	}
	if c&OpInverse != 0 && c&OpNoInverse == 0 {
		r = append(r, "7")
	}
	if c&OpConceal != 0 && c&OpNoConceal == 0 {
		r = append(r, "8")
	}
	if c&OpCrossedOut != 0 && c&OpNoCrossedOut == 0 {
		r = append(r, "9")
	}
	return strings.Join(r, ";")
}
func (c ColorOptions) String() string {
	if r := c.Code(); len(r) == 0 {
		return ""
	} else {
		return startCtr + r + endCtrl
	}
}
func (c ColorOptions) Write(writer io.StringWriter) {
	writer.WriteString(startCtr + c.Code() + endCtrl)
}

func (c ColorOptions) MergeFrom(older ColorOptions) ColorOptions {
	if older == 0 {
		return c
	}
	var merged ColorOptions
	merged |= c & opMask
	if older&OpBold != 0 && older&OpNoBold == 0 {
		merged |= OpBold | (c & OpBold)
	}
	if older&OpFaint != 0 && older&OpNoFaint == 0 {
		merged |= OpFaint | (c & OpFaint)
	}
	if older&OpItalic != 0 && older&OpNoItalic == 0 {
		merged |= OpItalic | (c & OpItalic)
	}
	if older&OpUnderline != 0 && older&OpNoUnderline == 0 {
		merged |= OpUnderline | (c & OpUnderline)
	}
	if older&OpBlinkSlow != 0 && older&OpNoBlinkSlow == 0 {
		merged |= OpBlinkSlow | (c & OpBlinkSlow)
	}
	if older&OpBlinkFast != 0 && older&OpNoBlinkFast == 0 {
		merged |= OpBlinkFast | (c & OpBlinkFast)
	}
	if older&OpInverse != 0 && older&OpNoInverse == 0 {
		merged |= OpInverse | (c & OpInverse)
	}
	if older&OpConceal != 0 && older&OpNoConceal == 0 {
		merged |= OpConceal | (c & OpConceal)
	}
	if older&OpCrossedOut != 0 && older&OpNoCrossedOut == 0 {
		merged |= OpCrossedOut | (c & OpCrossedOut)
	}
	return merged
}

func NewColorOptions(options ...ColorOptions) ColorOptions {
	var merged ColorOptions
	for _, option := range options {
		merged |= option
	}
	return merged
}
