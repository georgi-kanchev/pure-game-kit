package tilemap

import (
	"pure-kit/engine/data/path"
	"pure-kit/engine/execution/flow"
	"pure-kit/engine/geometry"
	"pure-kit/engine/graphics"
	"pure-kit/engine/internal"
	p "pure-kit/engine/tiled/property"
	"pure-kit/engine/tiled/tileset"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/text"
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

func LayerTiles(mapId, layerNameOrId string) []*graphics.Sprite {
	var mapData, has = internal.TiledMaps[mapId]
	if !has {
		return []*graphics.Sprite{}
	}

	var wantedLayer = findLayerTiles(mapData, layerNameOrId)
	if wantedLayer == nil {
		return []*graphics.Sprite{}
	}

	var result = make([]*graphics.Sprite, 0, mapData.Width*mapData.Height)
	var usedTilesets = usedTilesets(mapData)
	var tileIds = getTileIds(mapData, usedTilesets, wantedLayer)

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

// note:
//
//	objectNameOrClass == "" // include all objects within all tiles
func LayerTilesShapeGrid(mapId, layerNameOrId, objectNameOrClass string) *geometry.ShapeGrid {
	var mapData, has = internal.TiledMaps[mapId]
	if !has {
		return &geometry.ShapeGrid{}
	}

	var wantedLayer = findLayerTiles(mapData, layerNameOrId)
	if wantedLayer == nil {
		return &geometry.ShapeGrid{}
	}

	var result = geometry.NewShapeGrid(mapData.TileWidth, mapData.TileHeight)
	var tilesets = usedTilesets(mapData)
	var tileIds = getTileIds(mapData, tilesets, wantedLayer)

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

		tryAnimateTile(text.New(curTileset.AtlasId, "/", id, "-shapes"), curTileset, id, func(tileId int) {
			result.SetAtCell(x, y, tileset.TileShapes(curTileset.AtlasId, tileId, objectNameOrClass)...)
		})

		result.SetAtCell(x, y, tileset.TileShapes(curTileset.AtlasId, id, objectNameOrClass)...)
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

func findLayerTiles(data *internal.Map, layerNameOrId string) *internal.LayerTiles {
	var wantedLayer *internal.LayerTiles
	for _, layer := range data.LayersTiles {
		if layer.Name == layerNameOrId || layer.Id == int(text.ToNumber(layerNameOrId)) {
			wantedLayer = &layer
			break
		}
	}
	for _, group := range data.Groups {
		for _, layer := range group.LayersTiles {
			if layer.Name == layerNameOrId || layer.Id == int(text.ToNumber(layerNameOrId)) {
				wantedLayer = layer
				break
			}
		}
	}
	return wantedLayer
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
	if objTile == nil || objTile.Animation == nil {
		return
	}

	var animIds = tileset.TileAnimationTileIds(curTileset.AtlasId, tilesetTile)
	if len(animIds) == 0 {
		return
	}

	var animDurs = tileset.TileAnimationDurations(curTileset.AtlasId, tilesetTile)
	var steps = []flow.Step{}
	for stepIndex := range animIds {
		steps = append(steps, flow.Do(func() { onFrameChange(animIds[stepIndex]) }))
		steps = append(steps, flow.WaitForDelay(animDurs[stepIndex])) // frame delay
	}
	steps = append(steps, flow.Do(func() { flow.GoToStep(name, 0) })) // loop forever

	flow.NewSequence(name, steps...)
	flow.Start(name)
}

func color(hex string) uint {
	var trimmed = hex[1:]

	if len(trimmed) == 6 {
		trimmed += "FF"
	} else if len(trimmed) != 8 {
		return 0
	}

	return text.ToUint(trimmed)
}
