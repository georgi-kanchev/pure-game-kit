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

func PlaySound(soundId string) {
	var sound, has = internal.Sounds[soundId]
	if has {
		rl.PlaySound(*sound)
	}
}

func LoadTexturesFromFiles(filePaths ...string) []string {
	if !rl.IsWindowReady() {
		window.Recreate()
	}

	var result = []string{}
	for _, path := range filePaths {
		var absolutePath = filepath.Join(folder.PathOfExecutable(), path)
		path = strings.ReplaceAll(path, file.Extension(path), "")
		var _, has = internal.Textures[path]

		if has || !file.Exists(absolutePath) {
			continue
		}

		var texture = rl.LoadTexture(absolutePath)

		if texture.Width != 0 {
			internal.Textures[path] = &texture
			result = append(result, path)
		}
	}

	return result
}
func LoadAreaFromTexture(textureId, areaId string, x, y, width, height int) {
	var tex, has = internal.Textures[textureId]

	if has {
		var atlas = internal.Atlas{Texture: tex, CellWidth: 1, CellHeight: 1, Gap: 0}
		var rect = internal.AtlasRect{
			CellX: float32(x), CellY: float32(y), CountX: float32(width), CountY: float32(height), Atlas: &atlas}
		internal.Atlases[textureId] = atlas
		internal.AtlasRects[areaId] = rect
	}
}

func LoadAtlasFromTexture(textureId string, cellWidth, cellHeight, cellGap int) {
	var tex, has = internal.Textures[textureId]

	if has {
		var atlas = internal.Atlas{Texture: tex, CellWidth: cellWidth, CellHeight: cellHeight, Gap: cellGap}
		internal.Atlases[textureId] = atlas
	}
}
func LoadTilesFromAtlas(atlasId string, startCellX, startCellY int, tileIds ...string) {
	var atlas, has = internal.Atlases[atlasId]

	if !has {
		return
	}

	var tileCountX = int(atlas.Texture.Width / (int32(atlas.CellWidth + atlas.Gap)))
	var tileCountY = int(atlas.Texture.Height / (int32(atlas.CellHeight + atlas.Gap)))
	var index = 0

	for i := startCellY; i < tileCountY; i++ {
		for j := startCellX; j < tileCountX; j++ {
			if index >= len(tileIds) {
				return
			}

			var texRect = internal.AtlasRect{
				Atlas: &atlas, CellX: float32(j), CellY: float32(i), CountX: 1, CountY: 1}
			internal.AtlasRects[tileIds[index]] = texRect
			index++
		}
	}
}
func LoadTileFromAtlas(atlasId, tileId string, cellX, cellY, countX, countY float32) {
	var atlas, has = internal.Atlases[atlasId]

	if has {
		var texRect = internal.AtlasRect{Atlas: &atlas, CellX: cellX, CellY: cellY, CountX: countX, CountY: countY}
		internal.AtlasRects[tileId] = texRect
	}
}

func UnloadTextures(textureIds ...string) {
	for _, v := range textureIds {
		var tex, has = internal.Textures[v]

		if !has {
			continue
		}

		delete(internal.Textures, v)
		delete(internal.Atlases, v)

		for k, v := range internal.AtlasRects {
			if v.Atlas.Texture == tex {
				delete(internal.AtlasRects, k)
			}
		}

		rl.UnloadTexture(*tex)
	}
}

func UnloadSounds(soundIds ...string) {
	for _, v := range soundIds {
		var sound, has = internal.Sounds[v]

		if has {
			delete(internal.Sounds, v)
			rl.UnloadSound(*sound)
		}
	}
}
