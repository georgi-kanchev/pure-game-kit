package assets

import (
	"pure-game-kit/internal"
	"pure-game-kit/utility/collection"
	"pure-game-kit/utility/number"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadedTileSetIds() []string {
	return collection.MapKeys(internal.TileSets)
}
func LoadedTileDataIds() []string {
	return collection.MapKeys(internal.TileLayers)
}

func LoadTileSet(imageFilePath string, tileWidth, tileHeight int) string {
	var textureId = LoadTexture(imageFilePath)
	var atlas = &internal.TileSet{TextureId: textureId, TileWidth: tileWidth, TileHeight: tileHeight,
		PointsPerTile: make(map[uint16][]float32)}
	internal.TileSets[textureId] = atlas
	return textureId
}
func LoadTileData(id string, columns, rows int) string {
	tryCreateWindow()

	columns = number.Limit(columns, 1, 2048)
	rows = number.Limit(rows, 1, 2048)

	var data = &internal.TileLayer{Image: rl.GenImageColor(columns, rows, rl.Blank), CellsWithPoints: make(map[int]struct{})}
	var tex = rl.LoadTextureFromImage(data.Image)
	rl.SetTextureFilter(tex, rl.FilterPoint)
	data.Texture = &tex
	internal.TileLayers[id] = data
	return id
}

func UnloadTileSet(tileSetId string) {
	var atlas, has = internal.TileSets[tileSetId]

	if has {
		UnloadTexture(atlas.TextureId)
		delete(internal.TileSets, tileSetId)
	}
}
func UnloadTileData(tileMapId string) {
	var _, has = internal.TileLayers[tileMapId]

	if has && !isDefault(tileMapId) {
		delete(internal.TileLayers, tileMapId)
	}
}
func UnloadAllTileSets() {
	for id := range internal.TileSets {
		UnloadTileSet(id)
	}
}
func UnloadAllTileData() {
	for id := range internal.TileLayers {
		UnloadTileData(id)
	}
}
