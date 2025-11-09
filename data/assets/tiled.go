package assets

import (
	"pure-game-kit/data/path"
	"pure-game-kit/data/storage"
	"pure-game-kit/execution/condition"
	"pure-game-kit/internal"
	"pure-game-kit/utility/collection"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/text"
)

func LoadTiledTileset(filePath string) string {
	var tileset *internal.Tileset
	var w, h = 0, 0

	storage.FromFileXML(filePath, &tileset)
	if tileset == nil { // error is in storage
		return ""
	}

	tileset.AssetId = filePath

	if tileset.Image.Source != "" {
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
		tileset.MappedTiles[tile.Id] = &tile

		if len(tile.CollisionLayers) > 0 { // detect templates
			tryTemplate(tile.CollisionLayers, path.Folder(filePath))
		}

		if tileset.Image.Source == "" && tile.Image.Source != "" {
			tile.TextureId = LoadTexture(path.New(path.Folder(filePath), tile.Image.Source))
		}

		if len(tile.Animation.Frames) == 0 {
			continue
		} // animated tiles below

		if tileset.Image.Source == "" { // tiles are separate images, not in atlas
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
		tileset.AnimatedTiles = append(tileset.AnimatedTiles, &tile)
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
func LoadTiledWorld(filePath string) (tilemapIds []string) {
	var resultIds = []string{}
	var world *internal.World

	storage.FromFileJSON(filePath, &world)
	if world == nil { // error is in storage
		return resultIds
	}

	world.Directory = path.Folder(filePath)
	world.Name = path.RemoveExtension(path.LastPart(filePath))

	for _, m := range world.Maps {
		var mapPath = path.New(world.Directory, m.FileName)
		var name = path.RemoveExtension(path.LastPart(mapPath))

		if collection.Contains(resultIds, name) {
			continue
		}

		var mapIds = LoadTiledMap(mapPath)
		if len(mapIds) == 0 {
			continue
		}

		resultIds = append(resultIds, mapIds)
		var mp, _ = internal.TiledMaps[mapIds]
		mp.WorldX, mp.WorldY = float32(m.X), float32(m.Y)

		if mp.LayersObjects != nil {
			tryTemplate(mp.LayersObjects, world.Directory)
		}

		for _, grp := range mp.Groups {
			tryTemplate(grp.LayersObjects, world.Directory)
		}
	}

	return resultIds
}
func LoadTiledMap(filePath string) string {
	var mapData *internal.Map

	storage.FromFileXML(filePath, &mapData)
	if mapData == nil { // error is in storage
		return ""
	}

	mapData.Name = path.LastPart(path.RemoveExtension(filePath))
	mapData.Directory = path.Folder(filePath)
	internal.TiledMaps[filePath] = mapData

	for _, t := range mapData.Tilesets {
		LoadTiledTileset(path.New(mapData.Directory, t.Source))
	}

	return filePath
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

func ReloadAllTiledMaps() {
	for id := range internal.TiledMaps {
		LoadTiledMap(id)
	}
}
func ReloadAllTiledTilesets() {
	for id := range internal.TiledTilesets {
		LoadTiledTileset(id)
	}
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

// =================================================================
// private
var cachedTemplates = map[string]*internal.Template{}

func tryTemplate(layer []*internal.LayerObjects, directory string) {
	var objs = layer[0].Objects
	for i, o := range objs {
		if o.Template == "" {
			continue
		}

		var path = path.New(directory, o.Template)
		var template, _ = cachedTemplates[path]
		if template == nil {
			storage.FromFileXML(path, &template)
			cachedTemplates[path] = template
		}

		var newObj = template.Object
		newObj.X, newObj.Y = o.X, o.Y
		newObj.Width = condition.If(o.Width != 0, o.Width, newObj.Width)
		newObj.Height = condition.If(o.Height != 0, o.Height, newObj.Height)
		newObj.Rotation = condition.If(o.Rotation != 0, o.Rotation, newObj.Rotation)
		newObj.Name = condition.If(o.Name != "", o.Name, newObj.Name)
		newObj.Class = condition.If(o.Class != "", o.Class, newObj.Class)
		newObj.Visible = condition.If(o.Visible != "", o.Visible, newObj.Visible)
		newObj.Polygon.Points = condition.If(o.Polygon.Points != "", o.Polygon.Points, newObj.Polygon.Points)
		newObj.Polyline.Points = condition.If(o.Polyline.Points != "", o.Polyline.Points, newObj.Polyline.Points)

		for _, prop := range o.Properties {
			var has, p = hasProp(prop.Name, newObj.Properties)
			if has { // obj overwrites a template property
				prop.Value = p.Value
			} else { // obj adds a new property not present in the template
				newObj.Properties = append(newObj.Properties, prop)
			}
		}

		layer[0].Objects[i] = newObj
	}
}

func hasProp(name string, props []internal.Property) (bool, internal.Property) {
	for _, prop := range props {
		if prop.Name == name {
			return true, prop
		}
	}
	return false, internal.Property{}
}
