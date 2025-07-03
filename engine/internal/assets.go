package internal

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type AtlasRect struct {
	CellX, CellY, CountX, CountY float32
	Atlas                        *Atlas
}
type Atlas struct {
	Texture                    *rl.Texture2D
	CellWidth, CellHeight, Gap int
}

var Textures = make(map[string]*rl.Texture2D)
var AtlasRects = make(map[string]AtlasRect)
var Atlases = make(map[string]Atlas)

var TileMaps = make(map[string][]string)

var Sounds = make(map[string]*rl.Sound)
var Music = make(map[string]*rl.Music)

func AssetSize(assetId string) (width, height int) {
	var texture, fullTexture = Textures[assetId]
	width, height = 0, 0

	if fullTexture {
		return int(texture.Width), int(texture.Height)
	}

	var texRect, has = AtlasRects[assetId]
	if !has {
		return
	}

	var atlas = texRect.Atlas
	return atlas.CellWidth * int(texRect.CountX), atlas.CellHeight * int(texRect.CountY)
}
