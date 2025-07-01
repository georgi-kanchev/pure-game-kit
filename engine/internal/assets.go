package internal

import rl "github.com/gen2brain/raylib-go/raylib"

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
