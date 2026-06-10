package assets

import (
	_ "embed"

	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/debug"
	"pure-game-kit/packages/utility/file"
	"pure-game-kit/packages/utility/storage"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type FontId uint8

func LoadFont(pngPath string, jsonPath string) FontId {
	if !file.Exists(pngPath) {
		debug.LogError("Failed to find PNG file: \"", pngPath, "\"")
		return 0
	}

	var fontData = &internal.FontJSON{}
	storage.FromJSON(file.LoadText(jsonPath), fontData)

	var tex = rl.LoadTexture(pngPath)
	if tex.Width == 0 {
		debug.LogError("Failed to load png file: \"", pngPath, "\"")
		return 0
	}
	internal.NextImageId++
	var id = internal.NextImageId
	internal.Images[int32(id)] = internal.ImageData{Texture: tex, CropWidth: float32(tex.Width), CropHeight: float32(tex.Height)}

	ImageId(id).SetSmoothness(true) // bilinear filtering required for MSDF
	var font = internal.LoadFont(fontData, int32(id))
	return FontId(font)
}

func (f FontId) Unload() {
	var font, has = internal.Fonts[uint8(f)]
	if !has {
		return
	}
	ImageId(font.AtlasId).Unload()
	delete(internal.Fonts, uint8(f))
}
func (f FontId) SymbolArea(symbol rune, lineHeight float32) (offsetX, offsetY, width, height float32) {
	var font, has = internal.Fonts[uint8(f)]
	if !has {
		font = internal.Fonts[0]
	}
	var g = font.Chars[symbol]
	var x, y = g.PlaneBounds.Left * lineHeight, g.PlaneBounds.Top * lineHeight
	var w, h = (g.PlaneBounds.Right - g.PlaneBounds.Left) * lineHeight, (g.PlaneBounds.Top - g.PlaneBounds.Bottom) * lineHeight

	if symbol == ' ' {
		w, h = lineHeight/3, lineHeight
	}

	return x, y, w, h
}
func (f FontId) EmbedImage(symbol rune, imageId ImageId) {
	var font, has = internal.Fonts[uint8(f)]
	if !has {
		font = internal.Fonts[0]
	}
	var img = internal.Images[int32(imageId)]
	var aspect = img.CropWidth / img.CropHeight
	var g = font.Chars[symbol]
	g.EmbededImageId = int32(imageId)
	g.Advance = aspect
	g.PlaneBounds = internal.Bounds{Left: 0, Top: font.Ascender, Right: aspect, Bottom: font.Ascender + 1}
	g.AtlasBounds = internal.Bounds{Left: img.CropX, Top: img.CropY, Right: img.CropX + img.CropWidth, Bottom: img.CropY + img.CropHeight}
	font.Chars[symbol] = g
}
