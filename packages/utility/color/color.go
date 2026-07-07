// Helper functions that operate on an RGBA color, represented by a single uint number.
// It can be constructed in different ways and then manipulated.
package color

import (
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/random"
)

func Channels(color uint) (r, g, b, a uint8) { return colorToRGBA(color) }
func RGB(r, g, b uint8) uint                 { return colorFromRGBA(r, g, b, 255) }
func RGBA(r, g, b, a uint8) uint             { return colorFromRGBA(r, g, b, a) }
func TagHex(hex string) uint {
	if len(hex) > 0 && hex[0] == '#' {
		hex = hex[1:]
	}
	var r, g, b, a uint8 = 0, 0, 0, 255
	var length = len(hex)
	if length >= 6 {
		r = parseHexPair(hex[0], hex[1])
		g = parseHexPair(hex[2], hex[3])
		b = parseHexPair(hex[4], hex[5])
		if length >= 8 {
			a = parseHexPair(hex[6], hex[7])
		}
	}
	return RGBA(r, g, b, a)
}
func TagRGBA(str string) uint {
	if len(str) >= 5 && str[0:5] == "#rgba" {
		str = str[5:] // strip "#rgba" prefix if present
	} else if len(str) >= 4 && str[0:4] == "#rgb" {
		str = str[4:] // strip "#rgb" prefix if present
	}

	var r, g, b, a uint8 = 0, 0, 0, 255
	var i = 0
	if i < len(str) && str[i] == '(' {
		i++ // skip opening parenthesis if present
	}

	r, i = parseNextChannel(str, i)
	g, i = parseNextChannel(str, i)
	b, i = parseNextChannel(str, i)

	if i < len(str) && str[i] != ')' { // check if there's an alpha channel left before closing parenthesis
		a, i = parseAlphaChannel(str, i)
	}
	return RGBA(r, g, b, a)
}

func RandomBright() uint {
	var r = random.Range[uint8](127, 255)
	var g = random.Range[uint8](127, 255)
	var b = random.Range[uint8](127, 255)
	return RGB(r, g, b)
}
func RandomDark() uint {
	var r = random.Range[uint8](0, 127)
	var g = random.Range[uint8](0, 127)
	var b = random.Range[uint8](0, 127)
	return RGB(r, g, b)
}
func Random() uint {
	var r = random.Range[uint8](0, 255)
	var g = random.Range[uint8](0, 255)
	var b = random.Range[uint8](0, 255)
	return RGB(r, g, b)
}

func Darken(color uint, progress float32) uint {
	var r, g, b, a = colorToRGBA(color)

	r = uint8(number.Map(progress, 0, 1, float32(r), 0))
	g = uint8(number.Map(progress, 0, 1, float32(g), 0))
	b = uint8(number.Map(progress, 0, 1, float32(b), 0))
	return RGBA(r, g, b, a)
}

func Brighten(color uint, progress float32) uint {
	var r, g, b, a = colorToRGBA(color)

	r = uint8(number.Map(progress, 0, 1, float32(r), 255))
	g = uint8(number.Map(progress, 0, 1, float32(g), 255))
	b = uint8(number.Map(progress, 0, 1, float32(b), 255))
	return RGBA(r, g, b, a)
}

func Fade(color uint, target uint, progress float32) uint {
	var r, g, b, a = colorToRGBA(color)
	var tr, tg, tb, ta = colorToRGBA(target)

	r = uint8(number.Map(progress, 0, 1, float32(r), float32(tr)))
	g = uint8(number.Map(progress, 0, 1, float32(g), float32(tg)))
	b = uint8(number.Map(progress, 0, 1, float32(b), float32(tb)))
	a = uint8(number.Map(progress, 0, 1, float32(a), float32(ta)))
	return RGBA(r, g, b, a)
}

func FadeOut(color uint, progress float32) uint {
	var r, g, b, a = colorToRGBA(color)
	a = uint8(number.Map(progress, 0, 1, float32(a), 0))
	return RGBA(r, g, b, a)
}
func FadeIn(color uint, progress float32) uint {
	var r, g, b, a = colorToRGBA(color)
	a = uint8(number.Map(progress, 0, 1, float32(a), 255))
	return RGBA(r, g, b, a)
}
func Tint(color uint, tint uint) uint {
	var r, g, b, a = colorToRGBA(color)
	var tr, tg, tb, ta = colorToRGBA(tint)

	r = uint8((uint(r) * uint(tr)) / 255)
	g = uint8((uint(g) * uint(tg)) / 255)
	b = uint8((uint(b) * uint(tb)) / 255)
	a = uint8((uint(a) * uint(ta)) / 255)
	return RGBA(r, g, b, a)
}
func Opposite(color uint) uint {
	var r, g, b, a = colorToRGBA(color)
	return RGBA(255-r, 255-g, 255-b, a)
}

// private ========================================================

func colorFromRGBA(r, g, b, a uint8) uint {
	return uint(r)<<24 | uint(g)<<16 | uint(b)<<8 | uint(a)
}
func colorToRGBA(value uint) (r, g, b, a uint8) {
	r = uint8(value >> 24)
	g = uint8(value >> 16)
	b = uint8(value >> 8)
	a = uint8(value)
	return
}

func parseHexPair(c1, c2 byte) uint8 {
	return (parseHexChar(c1) << 4) | parseHexChar(c2)
}
func parseHexChar(c byte) uint8 {
	switch {
	case c >= '0' && c <= '9':
		return c - '0'
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10
	}
	return 0
}
func parseNextChannel(str string, i int) (uint8, int) {
	for i < len(str) && (str[i] == ' ' || str[i] == ',') {
		i++ // skip spaces or separators
	}

	var val uint = 0
	for i < len(str) && str[i] >= '0' && str[i] <= '9' {
		val = val*10 + uint(str[i]-'0')
		i++
	}
	if val > 255 {
		val = 255
	}
	return uint8(val), i
}
func parseAlphaChannel(str string, i int) (uint8, int) {
	for i < len(str) && (str[i] == ' ' || str[i] == ',') {
		i++ // skip spaces or separators
	}

	if i < len(str) && str[i] == '1' { // handle standard "0" or "1" quickly if there is no decimal
		if i+1 < len(str) && str[i+1] == '.' { // check if it's exactly 1 or 1.0
			// pass to decimal logic below
		} else {
			return 255, i + 1
		}
	} else if i < len(str) && str[i] == '0' && (i+1 >= len(str) || str[i+1] != '.') {
		return 0, i + 1
	}

	var hasDot = false // parse float manually (e.g., "0.79" or ".79")
	var alphaVal, divisor uint = 0, 1

	for i < len(str) && str[i] != ')' && str[i] != ' ' && str[i] != ',' {
		if str[i] == '.' {
			hasDot, i = true, i+1
			continue
		}
		if str[i] >= '0' && str[i] <= '9' {
			alphaVal = alphaVal*10 + uint(str[i]-'0')
			if hasDot {
				divisor *= 10
			}
		}
		i++
	}
	if divisor == 1 && alphaVal == 1 {
		return 255, i
	}

	var finalA = (alphaVal * 255) / divisor // map float 0.0-1.0 to 0-255 scale safely
	if finalA > 255 {
		finalA = 255
	}
	return uint8(finalA), i
}
