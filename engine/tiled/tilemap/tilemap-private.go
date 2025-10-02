package tilemap

import (
	"pure-kit/engine/data/path"
	"pure-kit/engine/internal"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/text"
	"strconv"
)

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
			return layer, nil, nil, &layer.Layer
		}
	}
	for _, layer := range layerObjs {
		if layerHas(&layer.Layer, layerNameOrId) {
			return nil, layer, nil, &layer.Layer
		}
	}
	for _, layer := range layerImgs {
		if layerHas(&layer.Layer, layerNameOrId) {
			return nil, nil, layer, &layer.Layer
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
func forEachTile(mapId, layerNameOrId string,
	do func(x, y, id int, layer *internal.LayerTiles, curTileset *internal.Tileset)) bool {
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

		do(x, y, id, tiles, curTileset)
	}
	return true
}
func getObj(mapId, layerNameOrId, objectNameClassOrId string) *internal.LayerObject {
	var mapData, _ = internal.TiledMaps[mapId]
	if mapData == nil {
		return nil
	}
	var _, objs, _, _ = findLayer(mapData, layerNameOrId)
	if objs == nil {
		return nil
	}

	for _, obj := range objs.Objects {
		if obj.Name == objectNameClassOrId || obj.Class == objectNameClassOrId ||
			text.New(obj.Id) == objectNameClassOrId {
			return &obj
		}
	}
	return nil
}

func col(hex string) uint {
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
