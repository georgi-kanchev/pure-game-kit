package assets

import (
	"pure-kit/engine/data/file"
	"pure-kit/engine/internal"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadTextures(filePath string) []string {
	tryCreateWindow()

	var result = []string{}
	var id, absolutePath = getIdPath(filePath)
	var tex, has = internal.Textures[id]

	if !file.IsExisting(absolutePath) {
		return result
	}

	if has { // hot reloading?
		rl.UnloadTexture(*tex)
	}

	var texture = rl.LoadTexture(absolutePath)

	if texture.Width != 0 {
		internal.Textures[id] = &texture
		result = append(result, id)
	}

	return result
}
func UnloadTextures(textureId string) {
	var tex, has = internal.Textures[textureId]
	delete(internal.Textures, textureId)
	RemoveTextureAtlases(textureId)

	if has {
		rl.UnloadTexture(*tex)
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
		var rect = internal.AtlasRect{
			CellX: cx, CellY: cy, CountX: cw, CountY: ch, AtlasId: textureId, Rotations: rotations, Flip: flip}
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
		var texRect = internal.AtlasRect{
			AtlasId: atlasId, CellX: cellX, CellY: cellY, CountX: countX, CountY: countY,
			Rotations: rotations, Flip: flip}
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
	delete(internal.Atlases, atlasId)

	for i, a := range internal.AtlasRects {
		if a.AtlasId == atlasId {
			delete(internal.AtlasRects, i)
		}
	}
}
func RemoveTextureAtlasTiles(tileId string) {
	delete(internal.AtlasRects, tileId)
}
func RemoveTextureBoxes(nineSliceId string) {
	delete(internal.Boxes, nineSliceId)
}
