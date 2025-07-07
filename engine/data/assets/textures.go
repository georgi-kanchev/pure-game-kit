package assets

import (
	"path/filepath"
	"pure-kit/engine/data/file"
	"pure-kit/engine/data/folder"
	"pure-kit/engine/internal"
	"pure-kit/engine/window"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadTextures(filePaths ...string) []string {
	if !rl.IsWindowReady() {
		window.Recreate()
	}

	var result = []string{}
	for _, path := range filePaths {
		var absolutePath = filepath.Join(folder.PathOfExecutable(), path)
		path = strings.ReplaceAll(path, file.Extension(path), "")
		var tex, has = internal.Textures[path]

		if !file.Exists(absolutePath) {
			continue
		}

		if has { // hot reloading?
			rl.UnloadTexture(*tex)
		}

		var texture = rl.LoadTexture(absolutePath)

		if texture.Width != 0 {
			internal.Textures[path] = &texture
			result = append(result, path)
		}
	}

	return result
}

func LoadTextureArea(textureId, areaId string, x, y, width, height int) string {
	var _, has = internal.Textures[textureId]

	if has && areaId != "" {
		var atlas = internal.Atlas{TextureId: textureId, CellWidth: 1, CellHeight: 1, Gap: 0}
		var cx, cy, cw, ch = float32(x), float32(y), float32(width), float32(height)
		var rect = internal.AtlasRect{CellX: cx, CellY: cy, CountX: cw, CountY: ch, AtlasId: textureId}
		internal.Atlases[textureId] = atlas
		internal.AtlasRects[areaId] = rect
		return textureId
	}

	return ""
}
func LoadTextureAtlas(textureId string, cellWidth, cellHeight, cellGap int) string {
	var _, has = internal.Textures[textureId]

	if has {
		var atlas = internal.Atlas{TextureId: textureId, CellWidth: cellWidth, CellHeight: cellHeight, Gap: cellGap}
		internal.Atlases[textureId] = atlas
		return textureId
	}

	return ""
}
func LoadTextureAtlasTiles(atlasId string, startCellX, startCellY int, tileIds ...string) []string {
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
func LoadTextureAtlasTile(atlasId, tileId string, cellX, cellY, countX, countY float32) string {
	var _, has = internal.Atlases[atlasId]

	if has && tileId != "" {
		var texRect = internal.AtlasRect{AtlasId: atlasId, CellX: cellX, CellY: cellY, CountX: countX, CountY: countY}
		internal.AtlasRects[tileId] = texRect
		return tileId
	}

	return ""
}

func UnloadTextures(textureIds ...string) {
	for _, v := range textureIds {
		var tex, has = internal.Textures[v]
		delete(internal.Textures, v)
		UnloadTextureAtlases(v)

		if has {
			rl.UnloadTexture(*tex)
		}
	}
}
func UnloadTextureAreas(areaIds ...string) {
	UnloadTextureAtlasTiles(areaIds...)
}
func UnloadTextureAtlases(atlasIds ...string) {
	for _, v := range atlasIds {
		delete(internal.Atlases, v)

		for i, a := range internal.AtlasRects {
			if a.AtlasId == v {
				delete(internal.AtlasRects, i)
			}
		}
	}
}
func UnloadTextureAtlasTiles(tileIds ...string) {
	for _, v := range tileIds {
		delete(internal.AtlasRects, v)
	}
}
