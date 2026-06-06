package internal

import rl "github.com/gen2brain/raylib-go/raylib"

type TileLayer struct {
	Image         *rl.Image
	Texture       rl.Texture2D
	Columns, Rows int

	LastDirtyTime   float32
	CellsWithPoints map[int]struct{}
	Objects         [][6]float32 // was ObjectPoints []float32 — each Shape represents one Tiled object
}
type TileAtlas struct {
	ImageId       int32
	TileSize      int
	ShapesPerTile map[uint16][][6]float32 // was PointsPerTile map[uint16][]float32
}

//=================================================================

var TileLayers = make(map[uint8]*TileLayer)
var TileAtlases = make(map[uint8]*TileAtlas)
var TileLayerNextId, TileAtlasNextId uint8
