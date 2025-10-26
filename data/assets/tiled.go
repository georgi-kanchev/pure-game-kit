package assets

import (
	"pure-game-kit/data/path"
	"pure-game-kit/data/storage"
	"pure-game-kit/execution/condition"
	"pure-game-kit/execution/flow"
	"pure-game-kit/internal"
	"pure-game-kit/utility/collection"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/text"
)

func LoadTiledTileset(filePath string) string {
	var tileset *internal.Tileset
	var id = path.RemoveExtension(filePath)
	var w, h = 0, 0

	filePath = internal.MakeAbsolutePath(filePath)
	storage.FromFileXML(filePath, &tileset)
	if tileset == nil {
		return ""
	}

	if tileset.Image.Source != "" {
		var textureIds = LoadTextures(path.New(path.Folder(filePath), tileset.Image.Source))
		var atlasId = SetTextureAtlas(textureIds[0], tileset.TileWidth, tileset.TileHeight, tileset.Spacing)

		tileset.AtlasId = atlasId
		w, h = tileset.Columns, tileset.TileCount/tileset.Columns

		for id := range w * h {
			var x, y = number.Index1DToIndexes2D(id, w, h)
			var rectId = text.New(atlasId, "/", id)
			SetTextureAtlasTile(atlasId, rectId, float32(x), float32(y), 1, 1, 0, false)
		}
	}

	internal.TiledTilesets[id] = tileset
	tileset.MappedTiles = map[int]*internal.TilesetTile{}
	for _, tile := range tileset.Tiles {
		tileset.MappedTiles[tile.Id] = &tile

		if len(tile.CollisionLayers) > 0 { // detect templates
			tryTemplate(tile.CollisionLayers, path.Folder(filePath))
		}

		if tileset.Image.Source == "" && tile.Image.Source != "" {
			tile.TextureId = LoadTextures(path.New(path.Folder(filePath), tile.Image.Source))[0]
		}

		if len(tile.Animation.Frames) == 0 {
			continue
		} // animated tiles

		if tileset.Image.Source == "" { // tiles are images, not in atlas
			w, h = tile.Image.Width, tile.Image.Height
		}

		var frame = 0
		var id = text.Remove(filePath, path.LastPart(filePath))
		var name = text.New(id, "/", tile.Id)
		var seq = flow.NewSequence()

		tile.Sequence = seq
		seq.SetSteps(true,
			flow.NowDoLoop(number.ValueMaximum[int](), func(int) {
				var timer = seq.CurrentStepTimer()
				var dur = float32(tile.Animation.Frames[frame].Duration) / 1000 // ms -> sec
				if timer > float32(dur) {
					var newId = tile.Animation.Frames[frame].TileId
					var x, y = number.Index1DToIndexes2D(newId, w, h) // new tile id coords
					SetTextureAtlasTile(id, name, float32(x), float32(y), 1, 1, 0, false)
					seq.GoToNextStep()
					frame++
					frame = frame % len(tile.Animation.Frames)
				}
			}),
			flow.NowDo(seq.Run))
	}

	return id
}
func LoadTiledWorld(filePath string) (tilemapIds []string) {
	filePath = internal.MakeAbsolutePath(filePath)
	var resultIds = []string{}
	var world *internal.World

	storage.FromFileJSON(filePath, &world)
	if world == nil {
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
	var absolutePath = internal.MakeAbsolutePath(filePath)
	var mapData *internal.Map

	storage.FromFileXML(absolutePath, &mapData)
	if mapData == nil {
		return ""
	}

	var root = path.Folder(path.Executable()) + path.Divider()
	var id = text.Replace(path.RemoveExtension(text.Remove(filePath, root)), "\\", "/")
	mapData.Name = path.LastPart(path.RemoveExtension(filePath))
	mapData.Directory = path.Folder(absolutePath)
	internal.TiledMaps[id] = mapData

	for _, t := range mapData.Tilesets {
		LoadTiledTileset(text.Remove(path.New(mapData.Directory, t.Source), root))
	}

	return id
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
