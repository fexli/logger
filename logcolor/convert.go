package logcolor

import "math"

var (
	step256 = []BasicColorMask{0x00, 0x5f, 0x87, 0xaf, 0xd7, 0xff}
)

func get256rv(c BasicColorMask) (h BasicColorMask) {
	if c == 0 {
		return 0
	}
	return c*40 + 55
}

func cvtrgbto16(r, g, b BasicColorMask) BasicColorMask {
	var bright, c, k BasicColorMask
	// eco bright-specific
	if r == 0x80 && g == 0x80 && b == 0x80 { // 0x80=128
		bright = 53
	} else if r == 0xff || g == 0xff || b == 0xff { // 0xff=255
		bright = 60
	} // else bright = 0

	if r == g && g == b {
		if r > 0x7f {
			r = 1
		} else {
			r = 0
		}
		if g > 0x7f {
			g = 1
		} else {
			g = 0
		}
		if b > 0x7f {
			b = 1
		} else {
			b = 0
		}

	} else {
		k = (r + g + b) / 3
		if r >= k {
			r = 1
		} else {
			r = 0
		}
		if g >= k {
			g = 1
		} else {
			g = 0
		}
		if b >= k {
			b = 1
		} else {
			b = 0
		}
	}
	if r > 0 {
		c = 1
	}
	if g > 0 {
		c += 2
	}
	if b > 0 {
		c += 4
	}
	return bright + c
}
func cvt256torgb(c BasicColorMask) (r, g, b BasicColorMask) {
	return get256rv((c - 16) / 36), get256rv((c - 16) % 36 / 6), get256rv((c - 16) % 6)
}
func cvt256to16(c BasicColorMask) BasicColorMask {
	if c < 8 {
		return c
	} else if c < 16 {
		return c + highLightExtra - 8
	} else if c < 232 {
		return cvtrgbto16(cvt256torgb(c))
	} else {
		return cvt256grayto16(c)
	}
}

func cvt256grayto16(c BasicColorMask) BasicColorMask {
	if c > 253 {
		return highLightExtra + 7
	} else if c >= 248 {
		return 7
	} else if c >= 237 {
		return highLightExtra
	}
	return 0
}
func abs(a, b BasicColorMask) BasicColorMask {
	if v := int16(a) - int16(b); v < 0 {
		return BasicColorMask(-v)
	} else {
		return BasicColorMask(v)
	}
}
func getrgbclosetstep(c BasicColorMask) BasicColorMask {
	var dst, di BasicColorMask
	for i, u := range step256 {
		if c == u {
			return BasicColorMask(i)
		}
		sm := BasicColorMask(math.Abs(float64(c - u)))
		if di == 0 || sm < dst {
			dst = sm
			di = BasicColorMask(i)
		}
	}
	return di
}
func cvtrgbto256(r, g, b BasicColorMask) BasicColorMask {
	return getrgbclosetstep(r)*36 + getrgbclosetstep(g)*6 + getrgbclosetstep(b) + 16
}
