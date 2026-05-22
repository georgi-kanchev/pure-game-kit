package assets

import (
	_ "embed"

	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/debug"
	"pure-game-kit/packages/utility/file"
	"pure-game-kit/packages/utility/storage"
)

type FontId byte

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
	var font, has = internal.Fonts[byte(f)]
	if !has {
		return
	}
	ImageId(font.AtlasId).UnloadImage()
	delete(internal.Fonts, byte(f))
}
