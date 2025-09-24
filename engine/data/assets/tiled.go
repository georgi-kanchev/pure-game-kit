package assets

import (
	"pure-kit/engine/data/path"
	"pure-kit/engine/data/storage"
	"pure-kit/engine/internal"
	"pure-kit/engine/utility/collection"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/text"
)

func LoadTiledTileset(tsxFilePath string) []string {
	var resultIds = []string{}
	var tileset *internal.Tileset
	var id = path.RemoveExtension(tsxFilePath)

	storage.FromFileXML(tsxFilePath, &tileset)
	if tileset == nil {
		return resultIds
	}

	resultIds = append(resultIds, id)
	internal.TiledTilesets[id] = tileset

	var texturePath = path.New(path.Folder(tsxFilePath), tileset.Image.Source)
	var textureIds = LoadTextures(texturePath)
	if len(textureIds) == 0 {
		return resultIds
	}

	var atlasId = SetTextureAtlas(textureIds[0], tileset.TileWidth, tileset.TileHeight, tileset.Spacing)
	var w, h = tileset.Columns, tileset.TileCount / tileset.Columns
	tileset.AtlasId = atlasId

	for id := range w * h {
		var x, y = number.Index1DToIndexes2D(id, w, h)
		var rectId = text.New(atlasId, "[", id, "]")
		SetTextureAtlasTile(atlasId, rectId, float32(x), float32(y), 1, 1, 0, false)
	}

	tileset.MappedTiles = map[int]internal.TilesetTile{}
	for _, tile := range tileset.Tiles {
		tileset.MappedTiles[tile.Id] = tile
	}

	return resultIds
}
func LoadTiledWorld(worldFilePath string) (tilemapIds []string) {
	var resultIds = []string{}
	var world *internal.World

	storage.FromFileJSON(worldFilePath, &world)
	if world == nil {
		return resultIds
	}

	world.Directory = path.Folder(worldFilePath)
	world.Name = path.RemoveExtension(path.LastElement(worldFilePath))

	for _, m := range world.Maps {
		var mapPath = path.New(world.Directory, m.FileName)
		var name = path.RemoveExtension(path.LastElement(mapPath))

		if collection.Contains(resultIds, name) {
			continue
		}

		var mapIds = LoadTiledMap(mapPath)
		if len(mapIds) == 0 {
			continue
		}

		resultIds = append(resultIds, mapIds...)
		for _, mapId := range mapIds {
			var mp, _ = internal.TiledMaps[mapId]
			mp.WorldX, mp.WorldY = float32(m.X), float32(m.Y)
		}
	}

	return resultIds
}
func LoadTiledMap(tmxFilePath string) []string {
	var resultIds = []string{}
	var mapData *internal.Map
	var name = path.RemoveExtension(tmxFilePath)

	storage.FromFileXML(tmxFilePath, &mapData)
	if mapData == nil {
		return resultIds
	}

	mapData.Name = path.LastElement(name)
	mapData.Directory = path.Folder(tmxFilePath)
	internal.TiledMaps[name] = mapData
	resultIds = append(resultIds, name)

	return resultIds
}

func UnloadTiledWorld(worldFilePath string) {
	var ids = LoadTiledWorld(worldFilePath)
	for _, id := range ids {
		UnloadTiledMap(id)
	}
}
func UnloadTiledMap(tilemapId string) {
	delete(internal.TiledMaps, tilemapId)
}
func UnloadTiledTileset(tilesetId string) {
	delete(internal.TiledTilesets, tilesetId)
}
