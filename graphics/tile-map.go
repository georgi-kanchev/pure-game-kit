package graphics

import "pure-game-kit/internal"

type TileMap struct {
	Node

	TileAtlasId, TileDataId string

	effects *Effects
}

func NewTileMap(tileAtlasId, tileDataId string) *TileMap {
	var tileMap = &TileMap{
		Node: *NewNode(0, 0), TileAtlasId: tileAtlasId, TileDataId: tileDataId, effects: NewEffects(),
	}
	var atlas = internal.TileAtlases[tileAtlasId]
	var data = internal.TileDatas[tileDataId]
	if atlas != nil && data != nil {
		tileMap.Width = float32(data.Image.Width * int32(atlas.TileWidth))
		tileMap.Height = float32(data.Image.Height * int32(atlas.TileHeight))
	}

	return tileMap
}
