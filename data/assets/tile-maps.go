package assets

import (
	"pure-game-kit/internal"
	"pure-game-kit/utility/collection"
	"pure-game-kit/utility/number"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadedTileAtlasIds() []string {
	return collection.MapKeys(internal.TileAtlases)
}
func LoadedTileDataIds() []string {
	return collection.MapKeys(internal.TileDatas)
}

func LoadTileAtlas(imageFilePath string, tileWidth, tileHeight int) string {
	var textureId = LoadTexture(imageFilePath)
	var atlas = &internal.TileAtlas{TextureId: textureId, TileWidth: tileWidth, TileHeight: tileHeight}
	internal.TileAtlases[textureId] = atlas
	return textureId
}
func LoadTileData(id string, columns, rows int) string {
	tryCreateWindow()

	columns = number.Limit(columns, 1, 2048)
	rows = number.Limit(rows, 1, 2048)

	var data = &internal.TileData{Image: rl.GenImageColor(columns, rows, rl.Blank)}
	var tex = rl.LoadTextureFromImage(data.Image)
	rl.SetTextureFilter(tex, rl.FilterPoint)
	data.Texture = &tex
	internal.TileDatas[id] = data
	return id
}

func UnloadTileAtlas(tileAtlasId string) {
	var atlas, has = internal.TileAtlases[tileAtlasId]

	if has {
		UnloadTexture(atlas.TextureId)
		delete(internal.TileAtlases, tileAtlasId)
	}
}
func UnloadTileData(tileMapId string) {
	var _, has = internal.TileDatas[tileMapId]

	if has && !isDefault(tileMapId) {
		delete(internal.TileDatas, tileMapId)
	}
}
func UnloadAllTileAtlases() {
	for id := range internal.TileAtlases {
		UnloadTileAtlas(id)
	}
}
func UnloadAllTileData() {
	for id := range internal.TileDatas {
		UnloadTileData(id)
	}
}
