package internal

import rl "github.com/gen2brain/raylib-go/raylib"

type TileLayer struct {
	Image         *rl.Image
	Texture       rl.Texture2D
	Columns, Rows int // present because Object layers don't have an Image data to get size from

	CellsWithPoints map[int]struct{} // hash set, 0 byte per value, only check if key is present
	Objects         [][6]float32
}
type TileAtlas struct {
	ImageId       int32
	TileSize      int
	ShapesPerTile map[uint16][][6]float32
}

//=================================================================

var TileLayers = make(map[uint8]*TileLayer)
var TileAtlases = make(map[uint8]*TileAtlas)
var TileLayerNextId, TileAtlasNextId uint8
