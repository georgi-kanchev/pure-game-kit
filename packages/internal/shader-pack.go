package internal

import (
	"math"
	"pure-game-kit/packages/utility/number"
)

const TypeShape, TypeSprite, TypeText, TypeTilemap byte = 0, 1, 2, 3

// floatSafe sets bits 31-24 = 0x3F so the float32 exponent is always
// in [0x7E, 0x7F] — a valid normal number, never denormal (0x00) or NaN (0xFF).
// Without this, 24-bit data in bits 23-0 can create denormals when bit 23 = 0.
const floatSafe = 0x3F000000

func pack24(bits uint32) float32 {
	return math.Float32frombits(floatSafe | bits)
}

func packU2(texWidth, texHeight uint16) float32 {
	var w = uint32(texWidth&0xFFF) << 12 // bits 23-12
	var h = uint32(texHeight & 0xFFF)    // bits 11-0
	return pack24(w | h)
}
func packV2(borderColor uint) float32 {
	return pack24(uint32(borderColor) & 0xFFFFFF)
}

//=================================================================

func packNormalX(gamma, saturation, contrast, brightness float32) float32 {
	var g = uint32(unitTo6Bit(gamma)) << 18      // bits 23-18
	var s = uint32(unitTo6Bit(saturation)) << 12 // bits 17-12
	var c = uint32(unitTo6Bit(contrast)) << 6    // bits 11-6
	var b = uint32(unitTo6Bit(brightness))       // bits 5-0
	return pack24(g | s | c | b)
}
func packNormalY(grayscale, inversion float32, blurX, blurY uint8) float32 {
	var g = uint32(unitTo6Bit(grayscale)) << 18 // bits 23-18
	var i = uint32(unitTo6Bit(inversion)) << 12 // bits 17-12
	var x = uint32(blurX&0x3F) << 6             // bits 11-6
	var y = uint32(blurY & 0x3F)                // bits 5-0
	return pack24(g | i | x | y)
}
func packNormalZ(depthZ float32, borderSize uint16, objType uint8) float32 {
	depthZ = number.Limit(depthZ, 0, 1)
	var d = uint32(uint16(depthZ*2047.0)) << 13 // bits 23-13 (11 bits)
	var b = uint32(borderSize&0x7FF) << 2       // bits 12-2  (11 bits)
	var t = uint32(objType & 0x3)               // bits 1-0   (2 bits)
	return pack24(d | b | t)
}

//=================================================================

func packTangentXSprite(outlineColor uint) float32 {
	return pack24(uint32(outlineColor) & 0xFFFFFF)
}
func packTangentYSprite(silhouetteColor uint) float32 {
	return pack24(uint32(silhouetteColor) & 0xFFFFFF)
}
func packTangentWSprite(roundness float32, pixelSize uint8) float32 {
	var r = uint32(unitTo16Bit(roundness)) << 8 // bits 23-8
	var p = uint32(pixelSize)                   // bits 7-0
	return pack24(r | p)
}

//=================================================================

func packTangentXTilemap(outlineColor uint) float32 {
	return pack24(uint32(outlineColor) & 0xFFFFFF)
}
func packTangentYTilemap(silhouetteColor uint) float32 {
	return pack24(uint32(silhouetteColor) & 0xFFFFFF)
}
func packTangentZTilemap(tileCols, tileRows uint16, pixelSize uint8) float32 {
	var c = uint32(tileCols&0x3FF) << 14 // bits 23-14
	var r = uint32(tileRows&0x3FF) << 4  // bits 13-4
	var p = uint32(pixelSize & 0xF)      // bits 3-0
	return pack24(c | r | p)
}
func packTangentWTilemap(outlineSize, tileSize, roundness uint8) float32 {
	var o = uint32(outlineSize) << 16 // bits 23-16
	var t = uint32(tileSize) << 8     // bits 15-8
	var r = uint32(roundness)         // bits 7-0
	return pack24(o | t | r)
}

//=================================================================

func packTangentXText(outlineColor uint) float32 {
	return pack24(uint32(outlineColor) & 0xFFFFFF)
}
func packTangentYText(shadowColor uint) float32 {
	return pack24(uint32(shadowColor) & 0xFFFFFF)
}
func packTangentZText(weight, outlineWeight, shadowWeight, shadowBlur uint8) float32 {
	var w = uint32(weight&0x3F) << 18        // bits 23-18
	var o = uint32(outlineWeight&0x3F) << 12 // bits 17-12
	var s = uint32(shadowWeight&0x3F) << 6   // bits 11-6
	var b = uint32(shadowBlur & 0x3F)        // bits 5-0
	return pack24(w | o | s | b)
}
func packTangentWText(textShadowX, textShadowY int8, roundness, pixelSize uint8) float32 {
	var x = uint32(uint8(textShadowX)&0x3F) << 18 // bits 23-18
	var y = uint32(uint8(textShadowY)&0x3F) << 12 // bits 17-12
	var r = uint32(roundness) << 4                // bits 11-4
	var p = uint32(pixelSize & 0xF)               // bits 3-0
	return pack24(x | y | r | p)
}

// =================================================================

func unitTo6Bit(value float32) uint8   { return uint8(number.Limit(value, 0, 1) * 63.0) }
func unitTo16Bit(value float32) uint16 { return uint16(number.Limit(value, 0, 1) * 65535.0) }
