package utility

import "math"

type Color struct {
	R, G, B, A byte
}

func (color Color) ToDark(progress float32) Color {
	r := byte(mapF(progress, 0, 1, float32(color.R), 0))
	g := byte(mapF(progress, 0, 1, float32(color.G), 0))
	b := byte(mapF(progress, 0, 1, float32(color.B), 0))
	return Color{r, g, b, color.A}
}
func (color Color) ToColor(target Color, progress float32) Color {
	r := byte(mapF(progress, 0, 1, float32(color.R), float32(target.R)))
	g := byte(mapF(progress, 0, 1, float32(color.G), float32(target.G)))
	b := byte(mapF(progress, 0, 1, float32(color.B), float32(target.B)))
	a := byte(mapF(progress, 0, 1, float32(color.A), float32(target.A)))
	return Color{r, g, b, a}
}
func (color Color) ToBright(progress float32) Color {
	r := byte(mapF(progress, 0, 1, float32(color.R), 255))
	g := byte(mapF(progress, 0, 1, float32(color.G), 255))
	b := byte(mapF(progress, 0, 1, float32(color.B), 255))
	return Color{r, g, b, color.A}
}
func (color Color) ToTransparent(progress float32) Color {
	a := byte(mapF(progress, 0, 1, float32(color.A), 0))
	return Color{color.R, color.G, color.B, a}
}
func (color Color) ToOpaque(progress float32) Color {
	a := byte(mapF(progress, 0, 1, float32(color.A), 255))
	return Color{color.R, color.G, color.B, a}
}
func (color Color) ToOpposite() Color {
	return Color{255 - color.R, 255 - color.G, 255 - color.B, color.A}
}

func RGB(r, g, b byte) Color {
	return Color{r, g, b, 255}
}
func RGBA(r, g, b, a byte) Color {
	return Color{r, g, b, a}
}
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
func Transparent() Color {
	return Color{0, 0, 0, 0}
}
func Black() Color {
	return RGB(0, 0, 0)
}
func Gray() Color {
	return RGB(127, 127, 127)
}
func White() Color {
	return RGB(255, 255, 255)
}
func Red() Color {
	return RGB(255, 0, 0)
}
func Green() Color {
	return RGB(0, 255, 0)
}
func Blue() Color {
	return RGB(0, 0, 255)
}
func Pink() Color {
	return RGB(255, 105, 180)
}
func Magenta() Color {
	return RGB(255, 0, 255)
}
func Violet() Color {
	return RGB(143, 0, 255)
}
func Purple() Color {
	return RGB(75, 0, 130)
}
func Yellow() Color {
	return RGB(255, 255, 0)
}
func Orange() Color {
	return RGB(255, 127, 80)
}
func Brown() Color {
	return RGB(150, 105, 25)
}
func Cyan() Color {
	return RGB(0, 255, 255)
}
func Azure() Color {
	return RGB(0, 127, 255)
}

// region private
func mapF(number float32, fromA, fromB, toA, toB float32) float32 { // copied from utility/number
	if math.Abs(float64(fromB-fromA)) < 0.001 {
		return (toA + toB) / 2
	}
	value := ((number-fromA)/(fromB-fromA))*(toB-toA) + toA
	if math.IsNaN(float64(value)) || math.IsInf(float64(value), 0) {
		return toA
	}
	return value
}

func random(a, b, seed float32) float32 { // copied from utility/number
	if a == b {
		return a
	}
	if a > b {
		a, b = b, a
	}
	diff := b - a
	intSeed := int32(seed * 2147483647)
	intSeed = (1103515245*intSeed + 12345) % 2147483647
	normalized := float32(intSeed) / 2147483647.0
	return a + normalized*diff
}
func randomByteRange(min, max byte) byte {
	val := random(float32(min), float32(max)+1, float32(math.NaN())) // max+1 to include max
	return byte(val)
}

// endregion
