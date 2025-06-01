package color

import (
	"math"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/random"
)

type Color struct{ R, G, B, A byte }

var (
	Transparent = RGBA(0, 0, 0, 0)
	Black       = RGB(0, 0, 0)
	Gray        = RGB(127, 127, 127)
	White       = RGB(255, 255, 255)
	Red         = RGB(255, 0, 0)
	Green       = RGB(0, 255, 0)
	Blue        = RGB(0, 0, 255)
	Pink        = RGB(255, 105, 180)
	Magenta     = RGB(255, 0, 255)
	Violet      = RGB(143, 0, 255)
	Purple      = RGB(75, 0, 130)
	Yellow      = RGB(255, 255, 0)
	Orange      = RGB(255, 127, 80)
	Brown       = RGB(150, 105, 25)
	Cyan        = RGB(0, 255, 255)
	Azure       = RGB(0, 127, 255)
)

func (color Color) Darken(progress float32) Color {
	r := byte(number.Map(progress, 0, 1, float32(color.R), 0))
	g := byte(number.Map(progress, 0, 1, float32(color.G), 0))
	b := byte(number.Map(progress, 0, 1, float32(color.B), 0))
	return Color{r, g, b, color.A}
}
func (color Color) To(target Color, progress float32) Color {
	r := byte(number.Map(progress, 0, 1, float32(color.R), float32(target.R)))
	g := byte(number.Map(progress, 0, 1, float32(color.G), float32(target.G)))
	b := byte(number.Map(progress, 0, 1, float32(color.B), float32(target.B)))
	a := byte(number.Map(progress, 0, 1, float32(color.A), float32(target.A)))
	return Color{r, g, b, a}
}
func (color Color) Brighten(progress float32) Color {
	r := byte(number.Map(progress, 0, 1, float32(color.R), 255))
	g := byte(number.Map(progress, 0, 1, float32(color.G), 255))
	b := byte(number.Map(progress, 0, 1, float32(color.B), 255))
	return Color{r, g, b, color.A}
}
func (color Color) FadeOut(progress float32) Color {
	a := byte(number.Map(progress, 0, 1, float32(color.A), 0))
	return Color{color.R, color.G, color.B, a}
}
func (color Color) FadeIn(progress float32) Color {
	a := byte(number.Map(progress, 0, 1, float32(color.A), 255))
	return Color{color.R, color.G, color.B, a}
}
func (color Color) Opposite() Color {
	return Color{255 - color.R, 255 - color.G, 255 - color.B, color.A}
}

func RGB(r, g, b byte) Color     { return Color{r, g, b, 255} }
func RGBA(r, g, b, a byte) Color { return Color{r, g, b, a} }
func RandomBright() Color {
	r := randomByteRange(127, 255)
	g := randomByteRange(127, 255)
	b := randomByteRange(127, 255)
	return RGB(r, g, b)
}
func RandomDark() Color {
	r := randomByteRange(0, 127)
	g := randomByteRange(0, 127)
	b := randomByteRange(0, 127)
	return RGB(r, g, b)
}
func Random() Color {
	r := randomByteRange(0, 255)
	g := randomByteRange(0, 255)
	b := randomByteRange(0, 255)
	return RGB(r, g, b)
}

// region private

func randomByteRange(min, max byte) byte {
	return byte(random.Range(float32(min), float32(max)+1, float32(math.NaN())))
}

// endregion
