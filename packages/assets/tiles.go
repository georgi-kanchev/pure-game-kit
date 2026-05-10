package assets

import (
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/number"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadTileSet(imagePath string, tileWidth, tileHeight int) string {
	// var textureId = LoadImage(imagePath)
	// var atlas = &internal.TileSet{ImageId: int32(textureId), TileWidth: tileWidth, TileHeight: tileHeight,
	// 	PointsPerTile: make(map[uint16][]float32)}
	// internal.TileSets[textureId] = atlas
	return ""
}
func LoadTileData(id string, columns, rows int) string {
	columns = number.Limit(columns, 1, 2048)
	rows = number.Limit(rows, 1, 2048)

	var data = &internal.TileLayer{Image: rl.GenImageColor(columns, rows, rl.Blank), CellsWithPoints: make(map[int]struct{})}
	var tex = rl.LoadTextureFromImage(data.Image)
	rl.SetTextureFilter(tex, rl.FilterPoint)
	data.Texture = &tex
	internal.TileLayers[id] = data
	return id
}
