package tilemap

import (
	"pure-kit/engine/data/path"
	"pure-kit/engine/execution/condition"
	"pure-kit/engine/execution/flow"
	"pure-kit/engine/geometry"
	"pure-kit/engine/graphics"
	"pure-kit/engine/internal"
	p "pure-kit/engine/tiled/property"
	"pure-kit/engine/tiled/tileset"
	"pure-kit/engine/utility/flag"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/text"
	"strconv"
)

func Property(mapId, property string) string {
	var data, has = internal.TiledMaps[mapId]
	if !has {
		return ""
	}

	switch property {
	case p.MapName:
		return data.Name
	case p.MapClass:
		return data.Class
	case p.MapTileWidth:
		return text.New(data.TileWidth)
	case p.MapTileHeight:
		return text.New(data.TileHeight)
	case p.MapColumns:
		return text.New(data.Width)
	case p.MapRows:
		return text.New(data.Height)
	case p.MapParallaxX:
		return text.New(data.ParallaxOriginX)
	case p.MapParallaxY:
		return text.New(data.ParallaxOriginY)
	case p.MapInfinite:
		return text.New(data.Infinite)
	case p.MapBackgroundColor:
		return text.New(color(data.BackgroundColor))
	}

	for _, v := range data.Properties {
		if v.Name == property {
			return v.Value
		}
	}
	return ""
}

func LayerProperty(mapId, layerNameOrId, property string) string {
	var mapData, _ = internal.TiledMaps[mapId]
	var _, objs, img, layer = findLayer(mapData, layerNameOrId)
	if mapData == nil || layer == nil {
		return ""
	}

	switch property {
	case p.LayerName:
		return layer.Name
	case p.LayerClass:
		return layer.Class
	case p.LayerVisible:
		return condition.If(layer.Visible == "", "true", "false")
	case p.LayerLocked:
		return text.New(layer.Locked)
	case p.LayerOpacity:
		return text.New(layer.Opacity)
	case p.LayerTint:
		return text.New(color(layer.Tint))
	case p.LayerOffsetX:
		return text.New(layer.OffsetX)
	case p.LayerOffsetY:
		return text.New(layer.OffsetY)
	case p.LayerParallaxX:
		return text.New(layer.ParallaxX)
	case p.LayerParallaxY:
		return text.New(layer.ParallaxY)
	//=================================================================
	case p.LayerColor:
		return text.New(color(objs.Color))
	case p.LayerDrawOrder:
		return text.New(objs.DrawOrder)
	//=================================================================
	case p.LayerImage:
		return text.New(img.Image.Source)
	case p.LayerTransparentColor:
		return text.New(img.Image.TransparentColor)
	case p.LayerRepeatX:
		return text.New(img.RepeatX)
	case p.LayerRepeatY:
		return text.New(img.RepeatY)
	case p.LayerWidth:
		return text.New(img.Image.Width)
	case p.LayerHeight:
		return text.New(img.Image.Height)
	}

	for _, prop := range layer.Properties {
		if prop.Name == property {
			return prop.Value
		}
	}
	return ""
}
func LayerTileId(mapId, layerNameOrId string, cellX, cellY int) int {
	var mapData, _ = internal.TiledMaps[mapId]
	var wantedLayer, _, _, _ = findLayer(mapData, layerNameOrId)
	if mapData == nil || wantedLayer == nil {
		return 0
	}

	cellX -= int(mapData.WorldX) / mapData.TileWidth
	cellY -= int(mapData.WorldY) / mapData.TileHeight
	var cellIndex = number.Indexes2DToIndex1D(cellY, cellX, mapData.Width, mapData.Height)
	var tilesets = usedTilesets(mapData)
	var tilesIds = getTileIds(mapData, tilesets, wantedLayer)
	var tileId = tilesIds[cellIndex]
	return tileId - 1 // 0 in map means empty but 0 is actually a valid tile in the tileset
}
func LayerSprites(mapId, layerNameOrId string) []*graphics.Sprite {
	var mapData, _ = internal.TiledMaps[mapId]
	var tiles, objs, img, _ = findLayer(mapData, layerNameOrId)
	if mapData == nil {
		return []*graphics.Sprite{}
	}

	if img != nil {
		var assetId = path.New(mapData.Directory, path.RemoveExtension(img.Image.Source))
		var sprite = graphics.NewSprite(assetId, mapData.WorldX+img.OffsetX, mapData.WorldY+img.OffsetY)
		sprite.PivotX, sprite.PivotY = 0, 0
		return []*graphics.Sprite{&sprite}
	}

	if objs != nil {
		var result = []*graphics.Sprite{}
		var usedTilesets = usedTilesets(mapData)
		for _, obj := range objs.Objects {
			if obj.Gid == 0 {
				continue
			}

			var id = flag.TurnOff(obj.Gid, internal.FlipX)
			id = flag.TurnOff(id, internal.FlipY)
			var curTileset = currentTileset(usedTilesets, id)
			var assetId = text.New(curTileset.AtlasId, "[", id-curTileset.FirstTileId, "]")
			var sprite = graphics.NewSprite(assetId, mapData.WorldX+obj.X, mapData.WorldY+obj.Y)
			sprite.X += float32(mapData.TileWidth) / 2
			sprite.Y -= float32(mapData.TileHeight) / 2
			sprite.Width, sprite.Height = float32(mapData.TileWidth), float32(mapData.TileHeight)
			sprite.ScaleX = condition.If(flag.IsOn(obj.Gid, internal.FlipX), float32(-1), 1)
			sprite.ScaleY = condition.If(flag.IsOn(obj.Gid, internal.FlipY), float32(-1), 1)

			result = append(result, &sprite)
		}
		return result
	}

	if tiles == nil {
		return []*graphics.Sprite{}
	}

	var result = make([]*graphics.Sprite, 0, mapData.Width*mapData.Height)
	var usedTilesets = usedTilesets(mapData)
	var tileIds = getTileIds(mapData, usedTilesets, tiles)

	for index, tile := range tileIds {
		var curTileset = currentTileset(usedTilesets, tile)
		if curTileset == nil {
			continue
		}

		var id = tile - curTileset.FirstTileId
		var tileId = text.New(curTileset.AtlasId, "[", id, "]")
		var j, i = number.Index1DToIndexes2D(index, mapData.Width, mapData.Height)
		var x = float32(j)*float32(mapData.TileWidth) + mapData.WorldX
		var y = float32(i)*float32(mapData.TileHeight) + mapData.WorldY
		var sprite = graphics.NewSprite(tileId, x, y)

		tryAnimateTile(text.New(curTileset.AtlasId, "/", id), curTileset, id, func(tileId int) {
			sprite.AssetId = text.New(curTileset.AtlasId, "[", tileId, "]")
		})
		sprite.Width, sprite.Height = float32(mapData.TileWidth), float32(mapData.TileHeight)
		sprite.PivotX, sprite.PivotY = 0, 0
		result = append(result, &sprite)
	}

	return result
}

func LayerObjectProperty(mapId, layerNameOrId, objectNameOrId, property string) string {
	var obj = getObj(mapId, layerNameOrId, objectNameOrId)
	if obj == nil {
		return ""
	}

	switch property {
	case p.ObjectName:
		return obj.Name
	case p.ObjectClass:
		return obj.Class
	case p.ObjectTemplate:
		return obj.Template
	case p.ObjectVisible:
		return condition.If(obj.Visible == "", "true", "false")
	case p.ObjectLocked:
		return text.New(obj.Locked)
	case p.ObjectX:
		return text.New(obj.X)
	case p.ObjectY:
		return text.New(obj.Y)
	case p.ObjectWidth:
		return text.New(obj.Width)
	case p.ObjectHeight:
		return text.New(obj.Height)
	case p.ObjectRotation:
		return text.New(obj.Rotation)
	case p.ObjectFlipX:
		return text.New(flag.IsOn(obj.Gid, internal.FlipX))
	case p.ObjectFlipY:
		return text.New(flag.IsOn(obj.Gid, internal.FlipY))
	case p.ObjectTileId:
		var id = flag.TurnOff(obj.Gid, internal.FlipX)
		id = flag.TurnOff(id, internal.FlipY)
		var mapData, _ = internal.TiledMaps[mapId]
		var current = currentTileset(usedTilesets(mapData), id)
		return text.New(id - current.FirstTileId)
	}

	for _, prop := range obj.Properties {
		if prop.Name == property {
			return prop.Value
		}
	}
	return ""
}

// empty objectNameOrClass includes all objects within all tiles
func LayerShapeGrid(mapId, tileLayerNameOrId, objectNameOrClass string) *geometry.ShapeGrid {
	var mapData, _ = internal.TiledMaps[mapId]
	if mapData == nil {
		return nil
	}

	var result = geometry.NewShapeGrid(mapData.TileWidth, mapData.TileHeight)
	var success = forEachTile(mapId, tileLayerNameOrId, func(x, y, id int, curTileset *internal.Tileset) {
		tryAnimateTile(text.New(curTileset.AtlasId, "/", id, "-shapes"), curTileset, id, func(tileId int) {
			result.SetAtCell(x, y, tileset.TileObjectShapes(curTileset.AtlasId, tileId, objectNameOrClass)...)
		})

		result.SetAtCell(x, y, tileset.TileObjectShapes(curTileset.AtlasId, id, objectNameOrClass)...)
	})
	return condition.If(success, result, nil)
}

// empty objectNameOrClass includes all objects within all tiles
func LayerPoints(mapId, layerNameOrId, objectNameOrClass string) [][2]float32 {
	var result = [][2]float32{}
	var success = forEachTile(mapId, layerNameOrId, func(x, y, id int, curTileset *internal.Tileset) {
		var pts = tileset.TileObjectPoints(curTileset.AtlasId, id, objectNameOrClass)

		if flow.IsExisting(text.New(curTileset.AtlasId, "/", id)) {
			return // skipping any points of animated tiles
		}

		for i := range pts {
			pts[i][0] += float32(x) * float32(curTileset.TileWidth)
			pts[i][1] += float32(y) * float32(curTileset.TileHeight)
		}
		result = append(result, pts...)
	})

	if !success {
		var mapData, _ = internal.TiledMaps[mapId]
		var _, objs, _, _ = findLayer(mapData, layerNameOrId)
		if mapData != nil && objs != nil {
			for _, obj := range objs.Objects {
				if obj.Width == 0 && obj.Height == 0 && obj.Polygon.Points == "" {
					result = append(result, [2]float32{mapData.WorldX + obj.X, mapData.WorldY + obj.Y})
				}
			}
		}

	}

	return result
}

//=================================================================
// private

func getTileIds(mapData *internal.Map, usedTilesets []*internal.Tileset, layer *internal.LayerTiles) []int {
	if layer.Tiles != nil {
		return layer.Tiles // fast return if cached
	} // cache otherwise

	var tileData = text.Trim(layer.TileData.Tiles)
	var rows = text.Split(tileData, "\n")
	layer.Tiles = make([]int, mapData.Width*mapData.Height)

	for i := 0; i < mapData.Height; i++ {
		var row = rows[i]
		if text.EndsWith(row, ",") {
			row = row[:len(row)-1]
		}

		var columns = text.Split(row, ",")
		for j := 0; j < mapData.Width; j++ {
			var tile = int(text.ToNumber(columns[j]))
			if tile == 0 {
				continue
			}

			var curTileset = currentTileset(usedTilesets, tile)
			if curTileset == nil {
				continue
			}

			var index = number.Indexes2DToIndex1D(i, j, mapData.Width, mapData.Height)
			layer.Tiles[index] = tile
		}
	}

	return layer.Tiles
}

func findLayer(data *internal.Map, layerNameOrId string) (
	*internal.LayerTiles, *internal.LayerObjects, *internal.LayerImage, *internal.Layer) {
	if data == nil {
		return nil, nil, nil, nil
	}

	var layerTiles = data.LayersTiles
	var layerObjs = data.LayersObjects
	var layerImgs = data.LayersImages

	for _, group := range data.Groups {
		layerTiles = append(layerTiles, group.LayersTiles...)
		layerObjs = append(layerObjs, group.LayersObjects...)
		layerImgs = append(layerImgs, group.LayersImages...)
	}

	for _, layer := range layerTiles {
		if layerHas(&layer.Layer, layerNameOrId) {
			return &layer, nil, nil, &layer.Layer
		}
	}
	for _, layer := range layerObjs {
		if layerHas(&layer.Layer, layerNameOrId) {
			return nil, &layer, nil, &layer.Layer
		}
	}
	for _, layer := range layerImgs {
		if layerHas(&layer.Layer, layerNameOrId) {
			return nil, nil, &layer, &layer.Layer
		}
	}

	return nil, nil, nil, nil
}
func layerHas(layer *internal.Layer, layerNameOrId string) bool {
	return layer.Name == layerNameOrId || layer.Id == int(text.ToNumber(layerNameOrId))

}
func usedTilesets(data *internal.Map) []*internal.Tileset {
	var usedTilesets = make([]*internal.Tileset, len(data.Tilesets))

	for i, tileset := range data.Tilesets {
		if tileset.Source != "" {
			var tilesetId = path.New(data.Directory, tileset.Source)
			tilesetId = path.RemoveExtension(path.LastElement(tilesetId))
			tilesetId = path.New(data.Directory, tilesetId)
			usedTilesets[i] = internal.TiledTilesets[tilesetId]
			if usedTilesets[i] != nil {
				usedTilesets[i].FirstTileId = tileset.FirstTileId
			}
			continue
		}

		usedTilesets[i] = &tileset
	}
	return usedTilesets
}
func currentTileset(usedTilesets []*internal.Tileset, tile int) *internal.Tileset {
	var curTileset = usedTilesets[0]
	for i := len(usedTilesets) - 1; i >= 0; i-- {
		if usedTilesets[i] != nil && tile > usedTilesets[i].FirstTileId {
			curTileset = usedTilesets[i]
			break
		}
	}
	return curTileset
}
func tryAnimateTile(name string, curTileset *internal.Tileset, tilesetTile int, onFrameChange func(tileId int)) {
	var objTile = curTileset.MappedTiles[tilesetTile]
	if objTile == nil || len(objTile.Animation.Frames) == 0 {
		return
	}

	var animIds = tileset.TileAnimationTileIds(curTileset.AtlasId, tilesetTile)
	if len(animIds) == 0 {
		return
	}

	var animDurs = tileset.TileAnimationDurations(curTileset.AtlasId, tilesetTile)
	var steps = []flow.Step{}
	for stepIndex := range animIds {
		steps = append(steps, flow.NowDo(func() { onFrameChange(animIds[stepIndex]) }))
		steps = append(steps, flow.NowWaitForDelay(animDurs[stepIndex])) // frame delay
	}
	steps = append(steps, flow.NowDo(func() { flow.GoToStep(name, 0) })) // loop forever

	flow.NewSequence(name, true, steps...)
}
func forEachTile(mapId, layerNameOrId string, do func(x, y, id int, curTileset *internal.Tileset)) bool {
	var mapData, _ = internal.TiledMaps[mapId]
	var tiles, _, _, _ = findLayer(mapData, layerNameOrId)
	if mapData == nil || tiles == nil {
		return false
	}

	var tilesets = usedTilesets(mapData)
	var tileIds = getTileIds(mapData, tilesets, tiles)

	for i, id := range tileIds {
		if id == 0 {
			continue
		}
		id-- // 0 in map means empty but 0 is actually a valid tile in the tileset

		var curTileset = currentTileset(tilesets, id)
		id -= curTileset.FirstTileId - 1 // same as id
		var tile, _ = curTileset.MappedTiles[id]
		if tile == nil || len(tile.CollisionLayers) == 0 {
			continue
		}

		var x, y = number.Index1DToIndexes2D(i, mapData.Width, mapData.Height)
		x += int(mapData.WorldX) / mapData.TileWidth
		y += int(mapData.WorldY) / mapData.TileHeight

		do(x, y, id, curTileset)
	}
	return true
}
func getObj(mapId, layerNameOrId, objectNameOrId string) *internal.LayerObject {
	var mapData, _ = internal.TiledMaps[mapId]
	if mapData == nil {
		return nil
	}
	var _, objs, _, _ = findLayer(mapData, layerNameOrId)
	if objs == nil {
		return nil
	}

	for _, obj := range objs.Objects {
		if obj.Name == objectNameOrId {
			return &obj
		}
		var id = text.ToNumber(objectNameOrId)
		if !number.IsNaN(id) && obj.Id == int(id) {
			return &obj
		}
	}
	return nil
}

func color(hex string) uint {
	var trimmed = hex[1:]

	if len(trimmed) == 6 {
		trimmed += "FF"
	} else if len(trimmed) != 8 {
		return 0
	}

	var value, err = strconv.ParseUint(trimmed, 16, 32)
	if err != nil {
		return 0
	}

	return uint(value)
}
