package assets

import (
	"maps"
	"pure-game-kit/data/file"
	"pure-game-kit/debug"
	"pure-game-kit/internal"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadTexture(filePath string) string {
	tryCreateWindow()

	var _, has = internal.Textures[filePath]
	if has {
		return filePath
	}

	if !file.IsExisting(filePath) {
		debug.LogError("Failed to find image file: \"", filePath, "\"")
		return ""
	}

	var texture = rl.LoadTexture(filePath)
	if texture.Width == 0 {
		debug.LogError("Failed to load image file: \"", filePath, "\"")
		return ""
	}
	internal.Textures[filePath] = &texture
	return filePath
}
func UnloadTexture(textureId string) {
	var tex, has = internal.Textures[textureId]
	if has && !isDefault(textureId) {
		delete(internal.Textures, textureId)
		rl.UnloadTexture(*tex)
	}
}

func ReloadAllTextures() {
	var loaded = maps.Keys(internal.Textures)
	UnloadAllTextures()
	for id := range loaded {
		LoadTexture(id)
	}
}
func UnloadAllTextures() {
	for id := range internal.Textures {
		UnloadTexture(id)
	}
}

func SetTextureSmoothness(textureId string, smooth bool) {
	var tex, has = internal.Textures[textureId]
	if has && smooth {
		rl.SetTextureFilter(*tex, rl.FilterBilinear)
	}
	if has && !smooth {
		rl.SetTextureFilter(*tex, rl.FilterPoint)
	}
}
func SetTextureArea(textureId, areaId string, x, y, width, height, rotations int, flip bool) string {
	var _, has = internal.Textures[textureId]

	if has && areaId != "" {
		var atlas = internal.Atlas{TextureId: textureId, CellWidth: 1, CellHeight: 1}
		var cx, cy, cw, ch = float32(x), float32(y), float32(width), float32(height)
		var rect = internal.AtlasRect{CellX: cx, CellY: cy, CountX: cw, CountY: ch,
			AtlasId: textureId, Rotations: rotations, Flip: flip}

		internal.Atlases[textureId] = atlas
		internal.AtlasRects[areaId] = rect
		return textureId
	}

	return ""
}
func SetTextureAtlas(textureId string, cellWidth, cellHeight, cellGap int) string {
	var _, has = internal.Textures[textureId]

	if has {
		var atlas = internal.Atlas{TextureId: textureId, CellWidth: cellWidth, CellHeight: cellHeight, Gap: cellGap}
		internal.Atlases[textureId] = atlas
		return textureId
	}

	return ""
}
func SetTextureAtlasTiles(atlasId string, startCellX, startCellY int, tileIds ...string) []string {
	var atlas, has = internal.Atlases[atlasId]
	var tex, _ = internal.Textures[atlas.TextureId]
	var tileCountX = int(tex.Width / (int32(atlas.CellWidth + atlas.Gap)))
	var tileCountY = int(tex.Height / (int32(atlas.CellHeight + atlas.Gap)))
	var index = 0
	var result = []string{}

	if !has {
		return result
	}

	for i := startCellY; i < tileCountY; i++ {
		for j := startCellX; j < tileCountX; j++ {
			if index >= len(tileIds) {
				return result
			}

			result = append(result, tileIds[index])

			if tileIds[index] == "" {
				index++
				continue
			}

			var cx, cy = float32(j), float32(i)
			var texRect = internal.AtlasRect{AtlasId: atlasId, CellX: cx, CellY: cy, CountX: 1, CountY: 1}
			internal.AtlasRects[tileIds[index]] = texRect
			index++
		}
	}

	return result
}
func SetTextureAtlasTile(atlasId, tileId string, cellX, cellY, countX, countY float32, rotations int, flip bool) string {
	var _, has = internal.Atlases[atlasId]

	if has && tileId != "" {
		var texRect = internal.AtlasRect{AtlasId: atlasId,
			CellX: cellX, CellY: cellY, CountX: countX, CountY: countY,
			Rotations: rotations, Flip: flip,
		}
		internal.AtlasRects[tileId] = texRect
		return tileId
	}

	return ""
}
func SetTextureBox(boxId string, assetIds [9]string) string {
	internal.Boxes[boxId] = assetIds
	return boxId
}

func RemoveTextureAreas(areaId string) {
	RemoveTextureAtlasTiles(areaId)
}
func RemoveTextureAtlases(atlasId string) {
	if isDefault(atlasId) {
		return
	}

	delete(internal.Atlases, atlasId)

	for i, a := range internal.AtlasRects {
		if a.AtlasId == atlasId {
			delete(internal.AtlasRects, i)
		}
	}
}
func RemoveTextureAtlasTiles(tileId string) {
	if !isDefault(tileId) {
		delete(internal.AtlasRects, tileId)
	}
}
func RemoveTextureBoxes(boxId string) {
	if !isDefault(boxId) {
		delete(internal.Boxes, boxId)
	}
}
