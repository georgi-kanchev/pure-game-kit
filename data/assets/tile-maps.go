package assets

import (
	"pure-game-kit/internal"
	"pure-game-kit/utility/collection"
	"pure-game-kit/utility/number"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadedTileMapIds() []string {
	return collection.MapKeys(internal.TileMaps)
}

func LoadTileMap(id string, columns, rows int) string {
	tryCreateWindow()

	columns = number.Limit(columns, 1, 2048)
	rows = number.Limit(rows, 1, 2048)

	var data struct {
		Image   *rl.Image
		Texture *rl.Texture2D
	}
	data.Image = rl.GenImageColor(columns, rows, rl.Blank)
	var tex = rl.LoadTextureFromImage(data.Image)
	data.Texture = &tex
	internal.TileMaps[id] = data
	return id
}

func UnloadTileMap(tileMapId string) {
	var _, has = internal.TileMaps[tileMapId]

	if has && !isDefault(tileMapId) {
		delete(internal.TileMaps, tileMapId)
	}
}
func UnloadAllTileMaps() {
	for id := range internal.TileMaps {
		UnloadTileMap(id)
	}
}
