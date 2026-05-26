package assets

import (
	_ "embed"

	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/debug"
	"pure-game-kit/packages/utility/file"
	"pure-game-kit/packages/utility/storage"
)

type FontId uint8

func LoadFont(pngPath string, jsonPath string) FontId {
	if !file.Exists(pngPath) {
		debug.LogError("Failed to find PNG file: \"", pngPath, "\"")
		return 0
	}

	var fontData = &internal.FontJSON{}
	storage.FromJSON(file.LoadText(jsonPath), fontData)
	var atlas = int32(LoadImage(pngPath))
	ImageId(atlas).SetSmoothness(true) // bilinear filtering required for MSDF
	var font = internal.LoadFont(fontData, atlas, false)
	return FontId(font)
}

func (f FontId) UnloadFont() {
	var font, has = internal.Fonts[uint8(f)]
	if !has {
		return
	}
	ImageId(font.AtlasId).UnloadImage()
	delete(internal.Fonts, uint8(f))
}

func (f FontId) SymbolArea(symbol rune, lineHeight float32) (offsetX, offsetY, width, height float32) {
	var font = internal.Fonts[uint8(f)]
	var g = font.Chars[symbol]
	var x, y = g.PlaneBounds.Left * lineHeight, g.PlaneBounds.Top * lineHeight
	var w, h = (g.PlaneBounds.Right - g.PlaneBounds.Left) * lineHeight, (g.PlaneBounds.Top - g.PlaneBounds.Bottom) * lineHeight

	if symbol == ' ' {
		w, h = lineHeight/3, lineHeight
	}

	return x, y, w, h
}
