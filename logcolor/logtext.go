package logcolor

import (
	"bytes"
	"io"
	"unsafe"
)

type LogTextCtx struct {
	Log      string  `json:"log"`
	Color    *Color  `json:"color"`
	InnerLog LogText `json:"inner"`
}

type ColorMask interface {
	Code() string                        // 转换为不带前后缀的颜色控制字符串
	String() string                      // 转换为颜色控制字符串
	Write(writer io.StringWriter)        // 将控制串写入writer
	IsEmpty() bool                       // 是否为空
	ToC16() ColorMask                    // 转换为16色`BasicColorIdentity`
	ToC256() ColorMask                   // 转换为128色`HundredColorIdentity`
	ToCRGB() ColorMask                   // 转换为RGB色`RGBColorIdentity`
	MergeFrom(older ColorMask) ColorMask // 合并两个颜色段
}
type LogText []*LogTextCtx

// WriteBytes write the text sequence control with
// `\x1b[` colored control sequence into io.StringWriter.
func (t *LogTextCtx) WriteBytes(buffer io.StringWriter, prevMask *Color) {
	if t == nil || buffer == nil {
		return
	}
	// check inner text
	if t.InnerLog != nil && len(t.InnerLog) != 0 {
		for _, ctx := range t.InnerLog {
			//mark := prevMask.MergeTo(ctx.Color)
			//mark.WriteStart(buffer)
			ctx.WriteBytes(buffer, prevMask)
			//mark.WriteEnd(buffer)
			//buffer.WriteString("[RIX]")
		}
	} else {
		if t.Log != "" {
			mask := prevMask.MergeTo(t.Color)
			mask.WriteStart(buffer)
			buffer.WriteString(t.Log)
			mask.WriteEnd(buffer)
		}
	}
}

// GetBytes return the text sequence control with `\x1b[` colored control sequence.
func (t *LogTextCtx) GetBytes() []byte {
	if t == nil {
		return []byte{}
	}
	buffer := &bytes.Buffer{}
	t.WriteBytes(buffer, t.Color)
	return buffer.Bytes()
}

// WriteRawBytes write the text sequence with
// no color control sequence into io.StringWriter.
func (t *LogTextCtx) WriteRawBytes(buffer io.StringWriter) {
	if t == nil || buffer == nil {
		return
	}
	// check inner text
	if t.InnerLog != nil && len(t.InnerLog) != 0 {
		for _, ctx := range t.InnerLog {
			ctx.WriteRawBytes(buffer)
		}
	} else {
		if t.Log != "" {
			buffer.WriteString(t.Log)
		}
	}
}

// GetRawBytes return the text sequence by WriteRawBytes().
func (t *LogTextCtx) GetRawBytes() []byte {
	if t == nil {
		return []byte{}
	}
	buffer := &bytes.Buffer{}
	t.WriteRawBytes(buffer)
	return buffer.Bytes()
}

// GetString returns the text sequence generated by LogTextCtx.GetBytes().
func (t *LogTextCtx) GetString() string {
	info := t.GetBytes()
	return *(*string)(unsafe.Pointer(&info))
}

// GetRawString returns the text sequence generated by LogTextCtx.GetRawBytes().
func (t *LogTextCtx) GetRawString() string {
	info := t.GetRawBytes()
	return *(*string)(unsafe.Pointer(&info))
}

// WriteString writes colored or non-colored string to console.
func (t *LogTextCtx) WriteString(to io.StringWriter, colored bool) {
	if colored {
		t.WriteBytes(to, t.Color)
	} else {
		t.WriteRawBytes(to)
	}
}

func LogTexts(texts ...*LogTextCtx) []*LogTextCtx {
	return texts
}

func New() *LogTextCtx {
	return &LogTextCtx{}
}

func (t *LogTextCtx) WithText(str string) *LogTextCtx {
	if t.InnerLog == nil || len(t.InnerLog) == 0 {
		t.Log = str
	} else {
		t.InnerLog = append(t.InnerLog, New().WithText(str))
	}
	return t
}

func (t *LogTextCtx) WithStr(str string) *LogTextCtx { return t.WithText(str) }

func (t *LogTextCtx) WithColor(mask *Color) *LogTextCtx {
	t.Color = mask
	return t
}
func (t *LogTextCtx) WithInner(texts ...*LogTextCtx) *LogTextCtx {
	t.InnerLog = texts
	return t
}
func (t *LogTextCtx) Then(text ...*LogTextCtx) *LogTextCtx {
	if t.Log != "" {
		// move text to inner

		newText := New()
		newText.Log = t.Log
		newText.Color = t.Color
		t.InnerLog = append(t.InnerLog, newText)
		t.Log = ""
		t.Color = nil
	}
	t.InnerLog = append(t.InnerLog, text...)
	return t
}

func (t *LogTextCtx) And(text ...*LogTextCtx) *LogTextCtx { return t.Then(text...) }

func ColorString(info string, mask ...*Color) *LogTextCtx {
	if mask == nil || len(mask) == 0 {
		return &LogTextCtx{Log: info}
	}
	return &LogTextCtx{Log: info, Color: mask[0]}
}

func RainbowString(info string, rgbFrom [3]uint8, rgbTo [3]uint8, isBg ...bool) *LogTextCtx {
	resultLog := New()
	data := []rune(info)
	if len(data) == 0 {
		return resultLog
	}
	var (
		r  = float32(rgbFrom[0])
		g  = float32(rgbFrom[1])
		b  = float32(rgbFrom[2])
		rv = (float32(rgbTo[0]) - r) / float32(len(data))
		gv = (float32(rgbTo[1]) - g) / float32(len(data))
		bv = (float32(rgbTo[2]) - b) / float32(len(data))
	)

	resultLog.InnerLog = make([]*LogTextCtx, 0, len(info))
	for i := 0; i < len(data); i++ {
		resultLog.InnerLog = append(resultLog.InnerLog, ColorString(string(data[i]), NewColor(RGB(uint8(r), uint8(g), uint8(b), isBg...))))
		r += rv
		g += gv
		b += bv
	}
	return resultLog
}