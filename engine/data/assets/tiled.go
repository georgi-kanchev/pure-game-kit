package assets

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"pure-kit/engine/data/path"
	"pure-kit/engine/internal"
	"pure-kit/engine/utility/collection"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/text"
)

func LoadTiledTileset(tsxFilePath string) []string {
	var resultIds = []string{}
	file, err := os.Open(tsxFilePath)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return resultIds
	}
	defer file.Close()

	var tileset *internal.Tileset
	var err2 = xml.NewDecoder(file).Decode(&tileset)
	if err2 != nil {
		return resultIds
	}

	var name = path.LastElement(tsxFilePath)
	name = path.RemoveExtension(name)
	var id = path.New(path.Folder(tsxFilePath), name)
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
	worldFile, err := os.Open(worldFilePath)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return resultIds
	}
	defer worldFile.Close()

	var world *internal.World
	var err2 = json.NewDecoder(worldFile).Decode(&world)
	if err2 != nil {
		return resultIds
	}

	var name = path.LastElement(worldFilePath)
	world.Directory = path.Folder(worldFilePath)
	world.Name = path.RemoveExtension(name)

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
	file, err := os.Open(tmxFilePath)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return resultIds
	}
	defer file.Close()

	var mapData *internal.Map
	var error = xml.NewDecoder(file).Decode(&mapData)
	if error != nil {
		return resultIds
	}

	var name = path.LastElement(tmxFilePath)
	name = path.RemoveExtension(name)
	mapData.Name = name
	mapData.Directory = path.Folder(tmxFilePath)
	var id = path.New(mapData.Directory, name)
	internal.TiledMaps[id] = mapData
	resultIds = append(resultIds, id)

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
