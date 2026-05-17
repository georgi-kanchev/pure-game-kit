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

// --- texCoord2 -----------------------------------------------------------
// texCoord2.x = TextureWidth(12) + TextureHeight(12) = 24 bits

func PackU2(texWidth, texHeight uint16) float32 {
	w := uint32(texWidth&0xFFF) << 12 // bits 23-12
	h := uint32(texHeight & 0xFFF)     // bits 11-0
	return pack24(w | h)
}

// texCoord2.y = BorderColor(24) = 24 bits

func PackV2(borderColor uint) float32 {
	return pack24(uint32(borderColor) & 0xFFFFFF)
}

// --- normal --------------------------------------------------------------
// normal.x = Gamma(6) + Saturation(6) + Contrast(6) + Brightness(6) = 24 bits

func PackNormalX(gamma, saturation, contrast, brightness float32) float32 {
	g := uint32(unitTo6Bit(gamma)) << 18      // bits 23-18
	s := uint32(unitTo6Bit(saturation)) << 12 // bits 17-12
	c := uint32(unitTo6Bit(contrast)) << 6    // bits 11-6
	b := uint32(unitTo6Bit(brightness))       // bits 5-0
	return pack24(g | s | c | b)
}

// normal.y = Grayscale(6) + Inversion(6) + BlurX(6) + BlurY(6) = 24 bits

func PackNormalY(grayscale, inversion float32, blurX, blurY uint8) float32 {
	g := uint32(unitTo6Bit(grayscale)) << 18 // bits 23-18
	i := uint32(unitTo6Bit(inversion)) << 12 // bits 17-12
	x := uint32(blurX&0x3F) << 6             // bits 11-6
	y := uint32(blurY & 0x3F)                // bits 5-0
	return pack24(g | i | x | y)
}

// normal.z = DepthZ(11) + BorderSize(11) + Type(2) = 24 bits

func PackNormalZ(depthZ float32, borderSize uint16, objType uint8) float32 {
	depthZ = number.Limit(depthZ, 0, 1)
	d := uint32(uint16(depthZ*2047.0)) << 13 // bits 23-13 (11 bits)
	b := uint32(borderSize&0x7FF) << 2       // bits 12-2  (11 bits)
	t := uint32(objType & 0x3)               // bits 1-0   (2 bits)
	return pack24(d | b | t)
}

// --- Shape tangent -------------------------------------------------------
// Shape tangent.x = Roundness(32)  — full float, no packing
// Shape tangent.y = PixelSize(32)  — full float, no packing

// --- Sprite tangent ------------------------------------------------------
// Sprite tangent.x = OutlineColor(24)

func PackTangentXSprite(outlineColor uint) float32 {
	return pack24(uint32(outlineColor) & 0xFFFFFF)
}

// Sprite tangent.y = SilhouetteColor(24)

func PackTangentYSprite(silhouetteColor uint) float32 {
	return pack24(uint32(silhouetteColor) & 0xFFFFFF)
}

// Sprite tangent.z = OutlineSize(32)  — full float, no packing

// Sprite tangent.w = Roundness(16) + PixelSize(8) = 24 bits

func PackTangentWSprite(roundness float32, pixelSize uint8) float32 {
	r := uint32(unitTo16Bit(roundness)) << 8 // bits 23-8
	p := uint32(pixelSize)                   // bits 7-0
	return pack24(r | p)
}

// --- Tilemap tangent -----------------------------------------------------
// Tilemap tangent.x = OutlineColor(24)

func PackTangentXTilemap(outlineColor uint) float32 {
	return pack24(uint32(outlineColor) & 0xFFFFFF)
}

// Tilemap tangent.y = SilhouetteColor(24)

func PackTangentYTilemap(silhouetteColor uint) float32 {
	return pack24(uint32(silhouetteColor) & 0xFFFFFF)
}

// Tilemap tangent.z = TileColumns(10) + TileRows(10) + PixelSize(4) = 24 bits

func PackTangentZTilemap(tileCols, tileRows uint16, pixelSize uint8) float32 {
	c := uint32(tileCols&0x3FF) << 14 // bits 23-14
	r := uint32(tileRows&0x3FF) << 4  // bits 13-4
	p := uint32(pixelSize & 0xF)      // bits 3-0
	return pack24(c | r | p)
}

// Tilemap tangent.w = OutlineSize(8) + TileSize(8) + Roundness(8) = 24 bits

func PackTangentWTilemap(outlineSize, tileSize, roundness uint8) float32 {
	o := uint32(outlineSize) << 16 // bits 23-16
	t := uint32(tileSize) << 8     // bits 15-8
	r := uint32(roundness)         // bits 7-0
	return pack24(o | t | r)
}

// --- Text tangent --------------------------------------------------------
// Text tangent.x = OutlineColor(24)

func PackTangentXText(outlineColor uint) float32 {
	return pack24(uint32(outlineColor) & 0xFFFFFF)
}

// Text tangent.y = ShadowColor(24)

func PackTangentYText(shadowColor uint) float32 {
	return pack24(uint32(shadowColor) & 0xFFFFFF)
}

// Text tangent.z = Weight(6) + OutlineWeight(6) + ShadowWeight(6) + ShadowBlur(6) = 24 bits

func PackTangentZText(weight, outlineWeight, shadowWeight, shadowBlur uint8) float32 {
	w := uint32(weight&0x3F) << 18        // bits 23-18
	o := uint32(outlineWeight&0x3F) << 12 // bits 17-12
	s := uint32(shadowWeight&0x3F) << 6   // bits 11-6
	b := uint32(shadowBlur & 0x3F)        // bits 5-0
	return pack24(w | o | s | b)
}

// Text tangent.w = TextShadowX(6) + TextShadowY(6) + Roundness(8) + PixelSize(4) = 24 bits
// TextShadowX/Y are signed; the raw two's-complement byte is stored in 6 bits.

func PackTangentWText(textShadowX, textShadowY int8, roundness, pixelSize uint8) float32 {
	x := uint32(uint8(textShadowX)&0x3F) << 18 // bits 23-18
	y := uint32(uint8(textShadowY)&0x3F) << 12 // bits 17-12
	r := uint32(roundness) << 4                // bits 11-4
	p := uint32(pixelSize & 0xF)              // bits 3-0
	return pack24(x | y | r | p)
}

// =================================================================

func unitToByte(value float32) uint8  { return uint8(number.Limit(value, 0, 1) * 255.0) }
func unitTo6Bit(value float32) uint8  { return uint8(number.Limit(value, 0, 1) * 63.0) }
func unitTo16Bit(value float32) uint16 { return uint16(number.Limit(value, 0, 1) * 65535.0) }
