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

func SetTile(tileDataId string, column, row int, tile uint16, rotation int, flip bool) {
	SetTileArea(tileDataId, column, row, 1, 1, tile, rotation, flip)
}

func SetTileArea(tileDataId string, column, row, width, height int, tile uint16, rotation int, flip bool) {
	var data = internal.TileDatas[tileDataId]
	if data == nil {
		return
	}

	// bits 0..15 = tile id
	// bit 29 & 30 = rotation
	// bit 31 = flip
	var packed = uint32(tile) | uint32(rotation%4)<<29
	if flip {
		packed |= 0x80000000
	}

	var r = uint8((packed >> 24) & 0xFF)
	var g = uint8((packed >> 16) & 0xFF)
	var b = uint8((packed >> 8) & 0xFF)
	var a = uint8((packed >> 0) & 0xFF)

	var colr = rl.NewColor(r, g, b, a)
	var rect = rl.NewRectangle(float32(column), float32(row), float32(width), float32(height))

	rl.ImageDrawRectangle(data.Image, int32(column), int32(row), int32(width), int32(height), colr)
	rl.UpdateTextureRec(*data.Texture, rect, collection.SameItems(width*height, colr))
}

func Tile(tileDataId string, column, row int) (tile uint16, rotation int, flipX bool) {
	var data = internal.TileDatas[tileDataId]
	if data == nil {
		return 0, 0, false
	}

	var c = rl.GetImageColor(*data.Image, int32(column), int32(row))
	var packed = uint32(c.R)<<24 | uint32(c.G)<<16 | uint32(c.B)<<8 | uint32(c.A)

	tile = uint16(packed & 0xFFFF)
	rotation = int((packed >> 29) & 0x03)
	flipX = (packed & 0x80000000) != 0
	return
}
