package color

import (
	"math"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/random"
)

var (
	Transparent = RGBA(0, 0, 0, 0)   // R[0] 	G[0] 	B[0] 	A[0]
	Black       = RGB(0, 0, 0)       // R[0] 	G[0] 	B[0] 	A[255]
	Gray        = RGB(127, 127, 127) // R[127] 	G[127] 	B[127] 	A[255]
	White       = RGB(255, 255, 255) // R[255] 	G[255] 	B[255] 	A[255]
	Red         = RGB(255, 0, 0)     // R[255] 	G[0] 	B[0] 	A[255]
	Green       = RGB(0, 255, 0)     // R[0] 	G[255] 	B[0] 	A[255]
	Blue        = RGB(0, 0, 255)     // R[0] 	G[0] 	B[255] 	A[255]
	Pink        = RGB(255, 105, 180) // R[255] 	G[105] 	B[108] 	A[255]
	Magenta     = RGB(255, 0, 255)   // R[255] 	G[0] 	B[255] 	A[255]
	Violet      = RGB(143, 0, 255)   // R[143] 	G[0] 	B[255] 	A[255]
	Purple      = RGB(75, 0, 130)    // R[75] 	G[0] 	B[130] 	A[255]
	Yellow      = RGB(255, 255, 0)   // R[255] 	G[255] 	B[0] 	A[255]
	Orange      = RGB(255, 127, 80)  // R[255] 	G[127] 	B[80] 	A[255]
	Brown       = RGB(150, 105, 25)  // R[150] 	G[105] 	B[25] 	A[255]
	Cyan        = RGB(0, 255, 255)   // R[0] 	G[255] 	B[255] 	A[255]
	Azure       = RGB(0, 127, 255)   // R[0] 	G[127] 	B[255] 	A[255]
)

func Channels(color uint) (r, g, b, a byte) { return colorToRGBA(color) }
func RGB(r, g, b byte) uint                 { return colorFromRGBA(r, g, b, 255) }
func RGBA(r, g, b, a byte) uint             { return colorFromRGBA(r, g, b, a) }

func RandomBright() uint {
	r := randomByteRange(127, 255)
	g := randomByteRange(127, 255)
	b := randomByteRange(127, 255)
	return RGB(r, g, b)
}
func RandomDark() uint {
	r := randomByteRange(0, 127)
	g := randomByteRange(0, 127)
	b := randomByteRange(0, 127)
	return RGB(r, g, b)
}
func Random() uint {
	r := randomByteRange(0, 255)
	g := randomByteRange(0, 255)
	b := randomByteRange(0, 255)
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

// region private

func randomByteRange(min, max byte) byte {
	return byte(random.Range(float32(min), float32(max)+1, float32(math.NaN())))
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

// endregion
