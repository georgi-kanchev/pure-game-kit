package internal

import rl "github.com/gen2brain/raylib-go/raylib"

type TileLayer struct {
	Image   *rl.Image
	Texture rl.Texture2D
	AtlasId uint8

	LastDirtyTime   float32
	CellsWithPoints map[int]struct{}
	ObjectPoints    []float32
}
type TileAtlas struct {
	ImageId               int32
	TileWidth, TileHeight int
	PointsPerTile         map[uint16][]float32
}

//=================================================================

var TileLayers = make(map[uint8]*TileLayer)
var TileAtlases = make(map[uint8]*TileAtlas)
var TileLayerNextId, TileAtlasNextId uint8
