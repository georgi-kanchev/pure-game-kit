package internal

import (
	"math"
	"pure-game-kit/packages/utility/number"
)

const TypeShape, TypeSprite, TypeText, TypeTilemap byte = 0, 1, 2, 3

func PackU2(texWidth, texHeight uint16, objType uint8) float32 {
	var w = uint32(texWidth&0xFFF) << 20   // 12 bits
	var h = uint32(texHeight&0xFFF) << 8   // 12 bits
	var t = uint32(objType)                // 8 bits
	return math.Float32frombits(w | h | t) // = 32 bits
}
func PackV2(borderColor uint, roundness uint8) float32 {
	var b = uint32(borderColor) & 0xFFFFFF00 // 24 bits
	var r = uint32(roundness)                // 8 bits
	return math.Float32frombits(b | r)       // = 32 bits
}

func PackNormalX(gamma, saturation, contrast, brightness float32) float32 {
	var g = uint32(gamma*63.0) << 18      // 6 bits
	var s = uint32(saturation*63.0) << 12 // 6 bits
	var c = uint32(contrast*63.0) << 6    // 6 bits
	var b = uint32(brightness * 63.0)     // 6 bits
	return float32(g | s | c | b)         // = 24 bits
}
func PackNormalY(grayscale, inversion float32, blurX, blurY uint8) float32 {
	var g = uint32(unitToByte(grayscale)) << 18 // 8 bits
	var i = uint32(unitToByte(inversion)) << 12 // 8 bits
	var x = uint32(blurX) << 6                  // 8 bits
	var y = uint32(blurY)                       // 8 bits
	return math.Float32frombits(g | i | x | y)  // = 32 bits
}
func PackNormalZ(depthZ float32, borderSize uint16, pixelSize uint8) float32 {
	depthZ = number.Limit(depthZ, 0, 1)
	var d = uint32(uint16(depthZ*4095.0)) << 20 // 12 bits
	var b = uint32(borderSize&0xFFF) << 8       // 12 bits
	var p = uint32(pixelSize)                   // 8 bits
	return math.Float32frombits(d | b | p)      // = 32 bits
}

func PackTangentXSpriteOrTilemap(outlineColor uint, outlineSize uint8) float32 {
	var c = uint32(outlineColor) & 0xFFFFFF00 // 24 bits
	var o = uint32(outlineSize)               // 8 bits
	return math.Float32frombits(c | o)        // = 32 bits
}
func PackTangentYSpriteOrTilemap(silhouetteColor uint) float32 {
	return math.Float32frombits(uint32(silhouetteColor)) // 32 bits
}
func PackTangentZTilemap(tileCols, tileRows uint16, tileSize uint8) float32 {
	var c = uint32(tileCols&0xFFF) << 20   // 12 bits
	var r = uint32(tileRows&0xFFF) << 8    // 12 bits
	var s = uint32(tileSize)               // 8 bits
	return math.Float32frombits(c | r | s) // = 32 bits
}
func PackTangentXText(outlineColor uint, textShadowX int8) float32 {
	var c = uint32(outlineColor) & 0xFFFFFF00 // 24 bits
	var s = uint32(uint8(textShadowX))        // 8 bits
	return math.Float32frombits(c | s)        // = 32 bits
}
func PackTangentYText(shadowColor uint, textShadowY int8) float32 {
	var c = uint32(shadowColor) & 0xFFFFFF00 // 24 bits
	var s = uint32(uint8(textShadowY))       // 8 bits
	return math.Float32frombits(c | s)       // = 32 bits
}
func PackTangentZText(weight, outlineWeight, shadowWeight, shadowBlur uint8) float32 {
	var w = uint32(weight) << 24               // 8 bits
	var o = uint32(outlineWeight) << 16        // 8 bits
	var e = uint32(shadowWeight) << 8          // 8 bits
	var b = uint32(shadowBlur)                 // 8 bits
	return math.Float32frombits(w | o | e | b) // = 32 bits
}

//=================================================================

func unitToByte(value float32) uint8 { return uint8(number.Limit(value, 0, 1) * 255.0) }
