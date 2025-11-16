package assets

import (
	"maps"
	"pure-game-kit/data/path"
	"pure-game-kit/data/storage"
	"pure-game-kit/internal"
	"pure-game-kit/utility/collection"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/text"
)

func LoadTiledProject(filePath string) string {
	tryCreateWindow()

	var _, has = internal.TiledProjects[filePath]
	if has {
		return filePath
	}

	var data *internal.Project
	storage.FromFileJSON(filePath, &data)
	if data == nil {
		return "" // error is in storage
	}

	internal.TiledProjects[filePath] = data
	return filePath
}
func LoadTiledMapsFromWorld(filePath string) (mapIds []string) {
	tryCreateWindow()

	var resultIds = []string{}
	var world *internal.World

	storage.FromFileJSON(filePath, &world)
	if world == nil {
		return resultIds // error is in storage
	}

	world.Directory = path.Folder(filePath)
	world.Name = path.RemoveExtension(path.LastPart(filePath))

	for _, m := range world.Maps {
		var mapId = LoadTiledMap(path.New(world.Directory, m.FileName))

		if !collection.Contains(resultIds, mapId) {
			resultIds = append(resultIds, mapId)
		}

		var mp, _ = internal.TiledMaps[mapId]
		mp.WorldX, mp.WorldY = float32(m.X), float32(m.Y)
	}

	return resultIds
}
func LoadTiledMap(filePath string) string {
	tryCreateWindow()

	var _, has = internal.TiledMaps[filePath]
	if has {
		return filePath
	}

	var mapData *internal.Map
	storage.FromFileXML(filePath, &mapData)
	if mapData == nil {
		return "" // error is in storage
	}

	mapData.Name = path.LastPart(path.RemoveExtension(filePath))
	mapData.Directory = path.Folder(filePath)
	mapData.FirstTileIds = make([]uint32, len(mapData.Tilesets))
	internal.TiledMaps[filePath] = mapData

	for i, t := range mapData.Tilesets {
		LoadTiledTileset(path.New(mapData.Directory, t.Source))

		// the tileset has no concept of first tile ids, it's a map concept
		// even though it's a tileset field (because of map embedded tilesets)
		mapData.FirstTileIds[i] = t.FirstTileId // so store it in map
		t.FirstTileId = 0                       // and zero it out in tileset to prevent any confusion
	}

	tryCacheLayerTileIds(mapData, &mapData.Layers)

	if mapData.LayersObjects != nil {
		tryTemplate(mapData.LayersObjects, mapData.Directory)
	}
	for _, grp := range mapData.LayersGroups {
		tryTemplate(grp.LayersObjects, mapData.Directory)
	}

	return filePath
}
func LoadTiledTileset(filePath string) string {
	tryCreateWindow()

	var _, has = internal.TiledTilesets[filePath]
	if has {
		return filePath
	}

	var tileset *internal.Tileset
	var w, h = 0, 0

	storage.FromFileXML(filePath, &tileset)
	if tileset == nil {
		return "" // error is in storage
	}

	tileset.AssetId = filePath

	if tileset.Image != nil {
		var textureId = LoadTexture(path.New(path.Folder(filePath), tileset.Image.Source))
		var atlasId = SetTextureAtlas(textureId, tileset.TileWidth, tileset.TileHeight, tileset.Spacing)

		w, h = tileset.Columns, tileset.TileCount/tileset.Columns

		for id := range w * h {
			var x, y = number.Index1DToIndexes2D(id, w, h)
			var rectId = path.New(atlasId, text.New(id))
			SetTextureAtlasTile(atlasId, rectId, float32(x), float32(y), 1, 1, 0, false)
		}
	}

	internal.TiledTilesets[tileset.AssetId] = tileset
	tileset.MappedTiles = map[uint32]*internal.TilesetTile{}
	for _, tile := range tileset.Tiles {
		tileset.MappedTiles[tile.Id] = tile

		if len(tile.CollisionLayers) > 0 { // detect templates
			tryTemplate(tile.CollisionLayers, path.Folder(filePath))
		}

		if tileset.Image == nil && tile.Image != nil {
			tile.TextureId = LoadTexture(path.New(path.Folder(filePath), tile.Image.Source))
		}

		if tile.Animation == nil {
			continue
		} // animated tiles below

		if tileset.Image == nil { // tiles are separate images, not in atlas
			w, h = tile.Image.Width, tile.Image.Height
		}

		var frame = 0
		var atlasId = path.New(path.Folder(tileset.AssetId), tileset.Image.Source)
		var tileId = path.New(atlasId, text.New(tile.Id))
		var totalAnimDuration float32
		for _, f := range tile.Animation.Frames {
			totalAnimDuration += float32(f.Duration) / 1000
		}

		var tileTime = totalAnimDuration
		tileset.AnimatedTiles = append(tileset.AnimatedTiles, tile)
		tile.IsAnimating = true
		tile.Update = func() {
			if !tile.IsAnimating {
				return
			}

			var dur = float32(tile.Animation.Frames[frame].Duration) / 1000 // ms -> sec
			tileTime += internal.Delta
			if tileTime > float32(dur) {
				tileTime = 0
				var newId = tile.Animation.Frames[frame].TileId
				var x, y = number.Index1DToIndexes2D(newId, uint32(w), uint32(h)) // new tile id coords
				SetTextureAtlasTile(atlasId, tileId, float32(x), float32(y), 1, 1, 0, false)
				frame++
				frame = frame % len(tile.Animation.Frames)
			}
		}
	}

	return filePath
}

func UnloadTiledMapsFromWorld(worldFilePath string) {
	var ids = LoadTiledMapsFromWorld(worldFilePath)
	for _, id := range ids {
		UnloadTiledMap(id)
	}
}
func UnloadTiledMap(tilemapId string) {
	delete(internal.TiledMaps, tilemapId)
}
func UnloadTiledTileset(tilesetId string) {
	var tileset, has = internal.TiledTilesets[tilesetId]
	if !has {
		return
	}

	for _, v := range tileset.Tiles {
		if v.TextureId != "" {
			UnloadTexture(v.TextureId)
		}
	}
	if tileset.AssetId != "" {
		UnloadTexture(tileset.AssetId)
	}

	delete(internal.TiledTilesets, tilesetId)
}
func UnloadTiledProject(projectId string) {
	delete(internal.TiledProjects, projectId)
}

func UnloadAllTiledMaps() {
	for id := range internal.TiledMaps {
		UnloadTiledMap(id)
	}
}
func UnloadAllTiledTilesets() {
	for id := range internal.TiledTilesets {
		UnloadTiledTileset(id)
	}
}
func UnloadAllTiledProjects() {
	for id := range internal.TiledProjects {
		UnloadTiledProject(id)
	}
}

func ReloadAllTiledMaps() {
	var loaded = maps.Keys(internal.TiledMaps)
	UnloadAllTiledMaps()
	for id := range loaded {
		LoadTiledMap(id)
	}
}
func ReloadAllTiledTilesets() {
	var loaded = maps.Keys(internal.TiledTilesets)
	UnloadAllTiledTilesets()
	for id := range loaded {
		LoadTiledTileset(id)
	}
}
func ReloadAllTiledProjects() {
	var loaded = maps.Keys(internal.TiledProjects)
	UnloadAllTiledProjects()
	for id := range loaded {
		LoadTiledProject(id)
	}
}
