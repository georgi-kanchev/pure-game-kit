package assets

import (
	"pure-kit/engine/data/path"
	"pure-kit/engine/data/storage"
	"pure-kit/engine/execution/condition"
	"pure-kit/engine/internal"
	"pure-kit/engine/utility/collection"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/text"
)

func LoadTiledTileset(tsxFilePath string) string {
	var tileset *internal.Tileset
	var id = path.RemoveExtension(tsxFilePath)

	storage.FromFileXML(tsxFilePath, &tileset)
	if tileset == nil {
		return ""
	}

	var texturePath = path.New(path.Folder(tsxFilePath), tileset.Image.Source)
	var textureIds = LoadTextures(texturePath)
	if len(textureIds) == 0 {
		return ""
	}

	var atlasId = SetTextureAtlas(textureIds[0], tileset.TileWidth, tileset.TileHeight, tileset.Spacing)
	var w, h = tileset.Columns, tileset.TileCount / tileset.Columns

	tileset.AtlasId = atlasId
	internal.TiledTilesets[id] = tileset

	for id := range w * h {
		var x, y = number.Index1DToIndexes2D(id, w, h)
		var rectId = text.New(atlasId, "[", id, "]")
		SetTextureAtlasTile(atlasId, rectId, float32(x), float32(y), 1, 1, 0, false)
	}

	tileset.MappedTiles = map[int]*internal.TilesetTile{}
	for _, tile := range tileset.Tiles {
		tileset.MappedTiles[tile.Id] = &tile

		// detect templates
		if len(tile.CollisionLayers) > 0 {
			tryTemplate(tile.CollisionLayers, path.Folder(tsxFilePath))
		}
	}

	return id
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

		resultIds = append(resultIds, mapIds)
		var mp, _ = internal.TiledMaps[mapIds]
		mp.WorldX, mp.WorldY = float32(m.X), float32(m.Y)

		tryTemplate(mp.LayersObjects, world.Directory)
		for _, grp := range mp.Groups {
			tryTemplate(grp.LayersObjects, world.Directory)
		}
	}

	return resultIds
}
func LoadTiledMap(tmxFilePath string) string {
	var mapData *internal.Map
	var name = path.RemoveExtension(tmxFilePath)

	storage.FromFileXML(tmxFilePath, &mapData)
	if mapData == nil {
		return ""
	}

	mapData.Name = path.LastElement(name)
	mapData.Directory = path.Folder(tmxFilePath)
	internal.TiledMaps[name] = mapData
	return name
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

// =================================================================
// private
var cachedTemplates = map[string]*internal.Template{}

func tryTemplate(layer []internal.LayerObjects, directory string) {
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
		newObj.PolygonTile.Points =
			condition.If(o.PolygonTile.Points != "", o.PolygonTile.Points, newObj.PolygonTile.Points)

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
