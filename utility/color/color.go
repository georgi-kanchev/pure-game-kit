package color

import (
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/random"
	"pure-game-kit/utility/text"
	"strconv"
)

var (
	Transparent = RGBA(0, 0, 0, 0)   // RGB(A): 0 0 0 (0)
	Black       = RGB(0, 0, 0)       // RGB(A): 0 0 0 (255)
	White       = RGB(255, 255, 255) // RGB(A): 255 255 255 (255)

	LightGray = RGB(191, 191, 191) // RGB(A): 191 191 191 (255)
	Gray      = RGB(127, 127, 127) // RGB(A): 127 127 127 (255)
	DarkGray  = RGB(64, 64, 64)    // RGB(A): 64 64 64 (255)

	LightRed = RGB(255, 127, 127) // RGB(A): 255 127 127 (255)
	Red      = RGB(255, 0, 0)     // RGB(A): 255 0 0 (255)
	DarkRed  = RGB(127, 0, 0)     // RGB(A): 127 0 0 (255)

	LightGreen = RGB(127, 255, 127) // RGB(A): 127 255 127 (255)
	Green      = RGB(0, 255, 0)     // RGB(A): 0 255 0 (255)
	DarkGreen  = RGB(0, 127, 0)     // RGB(A): 0 127 0 (255)

	LightBlue = RGB(127, 127, 255) // RGB(A): 127 127 255 (255)
	Blue      = RGB(0, 0, 255)     // RGB(A): 0 0 255 (255)
	DarkBlue  = RGB(0, 0, 127)     // RGB(A): 0 0 127 (255)

	LightYellow = RGB(255, 255, 127) // RGB(A): 255 255 127 (255)
	Yellow      = RGB(255, 255, 0)   // RGB(A): 255 255 0 (255)
	DarkYellow  = RGB(127, 127, 0)   // RGB(A): 127 127 0 (255)

	LightMagenta = RGB(255, 127, 255) // RGB(A): 255 127 255 (255)
	Magenta      = RGB(255, 0, 255)   // RGB(A): 255 0 255 (255)
	DarkMagenta  = RGB(127, 0, 127)   // RGB(A): 127 0 127 (255)

	LightCyan = RGB(127, 255, 255) // RGB(A): 127 255 255 (255)
	Cyan      = RGB(0, 255, 255)   // RGB(A): 0 255 255 (255)
	DarkCyan  = RGB(0, 127, 127)   // RGB(A): 0 127 127 (255)

	Pink        = RGB(255, 105, 180) // RGB(A): 255 105 180 (255)
	Violet      = RGB(143, 0, 255)   // RGB(A): 143 0 255 (255)
	Purple      = RGB(75, 0, 130)    // RGB(A): 75 0 130 (255)
	Orange      = RGB(255, 127, 80)  // RGB(A): 255 127 80 (255)
	Brown       = RGB(150, 105, 25)  // RGB(A): 150 105 25 (255)
	Azure       = RGB(0, 127, 255)   // RGB(A): 0 127 255 (255)
	ForestGreen = RGB(34, 139, 34)   // RGB(A): 34 139 34 (255)
	SkyBlue     = RGB(135, 206, 235) // RGB(A): 135 206 235 (255)
	Gold        = RGB(255, 215, 0)   // RGB(A): 255 215 0 (255)
	Silver      = RGB(192, 192, 192) // RGB(A): 192 192 192 (255)
	Bronze      = RGB(205, 127, 50)  // RGB(A): 205 127 50 (255)
	Beige       = RGB(245, 245, 220) // RGB(A): 245 245 220 (255)
	Cream       = RGB(255, 253, 208) // RGB(A): 255 253 208 (255)
	Tan         = RGB(210, 180, 140) // RGB(A): 210 180 140 (255)
	Olive       = RGB(128, 128, 0)   // RGB(A): 128 128 0 (255)
	Teal        = RGB(0, 128, 128)   // RGB(A): 0 128 128 (255)
	Turquoise   = RGB(64, 224, 208)  // RGB(A): 64 224 208 (255)
	Indigo      = RGB(75, 0, 130)    // RGB(A): 75 0 130 (255)
	Maroon      = RGB(128, 0, 0)     // RGB(A): 128 0 0 (255)
	Navy        = RGB(0, 0, 128)     // RGB(A): 0 0 128 (255)
	Lime        = RGB(191, 255, 0)   // RGB(A): 191 255 0 (255)
	Mint        = RGB(189, 252, 201) // RGB(A): 189 252 201 (255)
	SeaGreen    = RGB(46, 139, 87)   // RGB(A): 46 139 87 (255)
	Coral       = RGB(255, 127, 80)  // RGB(A): 255 127 80 (255)
	Salmon      = RGB(250, 128, 114) // RGB(A): 250 128 114 (255)
	Crimson     = RGB(220, 20, 60)   // RGB(A): 220 20 60 (255)
)

func Channels(color uint) (r, g, b, a byte) { return colorToRGBA(color) }
func RGB(r, g, b byte) uint                 { return colorFromRGBA(r, g, b, 255) }
func RGBA(r, g, b, a byte) uint             { return colorFromRGBA(r, g, b, a) }
func Hex(hex string) uint {
	var r, g, b, a uint64
	hex = text.Remove(hex, "#")
	if text.Length(hex) >= 6 {
		r, _ = strconv.ParseUint(hex[0:2], 16, 8)
		g, _ = strconv.ParseUint(hex[2:4], 16, 8)
		b, _ = strconv.ParseUint(hex[4:6], 16, 8)
		a = 255
	} else if text.Length(hex) == 8 {
		a, _ = strconv.ParseUint(hex[6:8], 16, 8)
	}
	return RGBA(byte(r), byte(g), byte(b), byte(a))
}

func RandomBright() uint {
	var r = randomByteRange(127, 255)
	var g = randomByteRange(127, 255)
	var b = randomByteRange(127, 255)
	return RGB(r, g, b)
}
func RandomDark() uint {
	var r = randomByteRange(0, 127)
	var g = randomByteRange(0, 127)
	var b = randomByteRange(0, 127)
	return RGB(r, g, b)
}
func Random() uint {
	var r = randomByteRange(0, 255)
	var g = randomByteRange(0, 255)
	var b = randomByteRange(0, 255)
	return RGB(r, g, b)
}

func Darken(color uint, progress float32) uint {
	var r, g, b, a = colorToRGBA(color)

	r = byte(number.Map(progress, 0, 1, float32(r), 0))
	g = byte(number.Map(progress, 0, 1, float32(g), 0))
	b = byte(number.Map(progress, 0, 1, float32(b), 0))
	return RGBA(r, g, b, a)
}

func Brighten(color uint, progress float32) uint {
	var r, g, b, a = colorToRGBA(color)

	r = byte(number.Map(progress, 0, 1, float32(r), 255))
	g = byte(number.Map(progress, 0, 1, float32(g), 255))
	b = byte(number.Map(progress, 0, 1, float32(b), 255))
	return RGBA(r, g, b, a)
}

func Fade(color uint, target uint, progress float32) uint {
	var r, g, b, a = colorToRGBA(color)
	var tr, tg, tb, ta = colorToRGBA(target)

	r = byte(number.Map(progress, 0, 1, float32(r), float32(tr)))
	g = byte(number.Map(progress, 0, 1, float32(g), float32(tg)))
	b = byte(number.Map(progress, 0, 1, float32(b), float32(tb)))
	a = byte(number.Map(progress, 0, 1, float32(a), float32(ta)))
	return RGBA(r, g, b, a)
}

func FadeOut(color uint, progress float32) uint {
	var r, g, b, a = colorToRGBA(color)
	a = byte(number.Map(progress, 0, 1, float32(a), 0))
	return RGBA(r, g, b, a)
}
func FadeIn(color uint, progress float32) uint {
	var r, g, b, a = colorToRGBA(color)
	a = byte(number.Map(progress, 0, 1, float32(a), 255))
	return RGBA(r, g, b, a)
}
func Opposite(color uint) uint {
	var r, g, b, a = colorToRGBA(color)
	return RGBA(255-r, 255-g, 255-b, a)
}

//=================================================================
// private

func randomByteRange(min, max byte) byte {
	return byte(random.Range(float32(min), float32(max)+1, number.NaN()))
}

func colorFromRGBA(r, g, b, a byte) uint {
	return uint(r)<<24 | uint(g)<<16 | uint(b)<<8 | uint(a)
}
func colorToRGBA(value uint) (r, g, b, a uint8) {
	r = uint8(value >> 24)
	g = uint8(value >> 16)
	b = uint8(value >> 8)
	a = uint8(value)
	return
}
