package internal

import (
	"math"
	"pure-game-kit/packages/utility/number"
)

// texCoord2.x = TextureWidth(12) + TextureHeight(12)
// texCoord2.y = BorderColor(6,6,6,6)
//
// normal.x = Gamma(6) + Saturation(6) + Contrast(6) + Brightness(6)
// normal.y = Roundness(10) + PixelSize(4) + BlurX(5) + BlurY(5)
// normal.z = DepthZ(11) + BorderSize(11) + Type(2)
//
// Shape:
//  tangent = free
//
// Sprite:
//  tangent.x = OutlineColor(6,6,6,6)
//  tangent.y = SilhouetteColor(6,6,6,6)
//  tangent.z = OutlineSize(32)
//  tangent.w = free
//
// Text:
//  tangent.x = OutlineColor(6,6,6,6)
//  tangent.y = ShadowColor(6,6,6,6)
//  tangent.z = Weight(8) + OutlineWeight(8) + ShadowWeight(8)
//  tangent.w = TextShadowX(8) + TextShadowY(8) + ShadowBlur(8)
//
// Tilemap:
//  tangent.x = OutlineColor(6,6,6,6)
//  tangent.y = SilhouetteColor(6,6,6,6)
//  tangent.z = TileColumns(12) + TileRows(12)
//  tangent.w = OutlineSize(16) + TileSize(8)

const typeShape, typeSprite, typeText, typeTilemap uint8 = 0, 1, 2, 3

const floatSafe = 0x3F000000 // required to preserve 24 bits correctly from the 32 bits

func pack24(bits uint32) float32 {
	return math.Float32frombits(floatSafe | bits)
}

func packU2(texWidth, texHeight uint16) float32 {
	var w = uint32(texWidth&0xFFF) << 12 // bits 23-12
	var h = uint32(texHeight & 0xFFF)    // bits 11-0
	return pack24(w | h)
}
func packV2(borderColor uint) float32 {
	return packColor24(borderColor)
}

func packColor24(c uint) float32 {
	r := uint32(uint8(c>>24)>>2) << 18
	g := uint32(uint8(c>>16)>>2) << 12
	b := uint32(uint8(c>>8)>>2) << 6
	a := uint32(uint8(c) >> 2)
	return pack24(r | g | b | a)
}

//=================================================================

func packNormalX(gamma, saturation, contrast, brightness int8) float32 {
	// Quantize int8 [-128,127] to 6-bit [0,63]; 0 → 31 (center, ~old 0.5)
	var g = uint32((uint16(int16(gamma)+128)*63)/255) << 18      // bits 23-18
	var s = uint32((uint16(int16(saturation)+128)*63)/255) << 12 // bits 17-12
	var c = uint32((uint16(int16(contrast)+128)*63)/255) << 6    // bits 11-6
	var b = uint32((uint16(int16(brightness)+128)*63)/255)       // bits 5-0
	return pack24(g | s | c | b)
}
func packNormalY(roundness float32, pixelSize, blurX, blurY uint8) float32 {
	var r = uint32(unitTo10Bit(roundness)) << 14 // bits 23-14
	var p = uint32(pixelSize&0xF) << 10          // bits 13-10
	var x = uint32(blurX&0x1F) << 5              // bits 9-5
	var y = uint32(blurY & 0x1F)                 // bits 4-0
	return pack24(r | p | x | y)
}
func packNormalZ(depthZ float32, borderSize float32, objType uint8) float32 {
	depthZ = number.Limit(depthZ, 0, 1)
	var d = uint32(uint16(depthZ*2047.0)) << 13            // bits 23-13 (11 bits)
	var b = uint32(uint16(int16(borderSize*4))&0x7FF) << 2 // bits 12-2  (11 bits, signed, step 0.25)
	var t = uint32(objType & 0x3)                          // bits 1-0   (2 bits)
	return pack24(d | b | t)
}

//=================================================================

func packTangentXSprite(outlineColor uint) float32 {
	return packColor24(outlineColor)
}
func packTangentYSprite(silhouetteColor uint) float32 {
	return packColor24(silhouetteColor)
}
func packTangentWSprite() float32 {
	return 0 // tangent.w is free for sprites
}

//=================================================================

func packTangentXTilemap(outlineColor uint) float32 {
	return packColor24(outlineColor)
}
func packTangentYTilemap(silhouetteColor uint) float32 {
	return packColor24(silhouetteColor)
}
func packTangentZTilemap(tileCols, tileRows uint16) float32 {
	var c = uint32(tileCols&0xFFF) << 12 // bits 23-12
	var r = uint32(tileRows & 0xFFF)     // bits 11-0
	return pack24(c | r)
}
func packTangentWTilemap(outlineSize uint16, tileSize uint8) float32 {
	var o = uint32(outlineSize&0xFFFF) << 8 // bits 23-8
	var t = uint32(tileSize)                // bits 7-0
	return pack24(o | t)
}

//=================================================================

func packTangentXText(outlineColor uint) float32 {
	return packColor24(outlineColor)
}
func packTangentYText(shadowColor uint) float32 {
	return packColor24(shadowColor)
}
func packTangentZText(weight int8, outlineWeight uint8, shadowWeight int8) float32 {
	var w = uint32(uint8(weight)) << 16       // bits 23-16 (signed → unsigned via two's complement)
	var o = uint32(outlineWeight) << 8        // bits 15-8
	var s = uint32(uint8(shadowWeight))       // bits 7-0  (signed → unsigned via two's complement)
	return pack24(w | o | s)
}
func packTangentWText(textShadowX, textShadowY int8, shadowBlur uint8) float32 {
	var x = uint32(uint8(textShadowX)) << 16 // bits 23-16
	var y = uint32(uint8(textShadowY)) << 8  // bits 15-8
	var b = uint32(shadowBlur)               // bits 7-0
	return pack24(x | y | b)
}

// =================================================================

func unitTo10Bit(value float32) uint16 { return uint16(number.Limit(value, 0, 1) * 1023.0) }
