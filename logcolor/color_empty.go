package logcolor

import "io"

type colorEmpty struct{}

var EmptyColor = &colorEmpty{}

func (c *colorEmpty) Code() string {
	return ""
}

func (c *colorEmpty) String() string {
	return ""
}

func (c *colorEmpty) Write(_ io.StringWriter) {
	return
}

func (c *colorEmpty) IsEmpty() bool {
	return true
}

func (c *colorEmpty) ToC16() ColorMask {
	r := newBasic()
	return &r
}

func (c *colorEmpty) ToC256() ColorMask {
	r := newHundred()
	return &r
}

func (c *colorEmpty) ToCRGB() ColorMask {
	r := newRgb()
	return &r
}

func (c *colorEmpty) MergeFrom(older ColorMask) ColorMask {
	if older == nil {
		return c
	}
	return older
}
