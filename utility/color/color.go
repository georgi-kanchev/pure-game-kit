/*
Helper functions that operate on an RGBA color, represented by a single uint number.
It can be constructed in different ways and then manipulated.
*/
package color

import (
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/random"
	"pure-game-kit/utility/text"
	"strconv"
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
	var r = random.Range[byte](127, 255)
	var g = random.Range[byte](127, 255)
	var b = random.Range[byte](127, 255)
	return RGB(r, g, b)
}
func RandomDark() uint {
	var r = random.Range[byte](0, 127)
	var g = random.Range[byte](0, 127)
	var b = random.Range[byte](0, 127)
	return RGB(r, g, b)
}
func Random() uint {
	var r = random.Range[byte](0, 255)
	var g = random.Range[byte](0, 255)
	var b = random.Range[byte](0, 255)
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
