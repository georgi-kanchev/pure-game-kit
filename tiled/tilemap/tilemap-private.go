package tilemap

import (
	"bytes"
	"encoding/binary"
	"pure-game-kit/data/path"
	"pure-game-kit/data/storage"
	"pure-game-kit/internal"
	"pure-game-kit/utility/flag"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/text"
	"strconv"
)

func getTileIds(mapData *internal.Map, usedTilesets []*internal.Tileset, layer *internal.LayerTiles) []uint32 {
	if layer.Tiles != nil {
		return layer.Tiles // fast return if cached
	} // cache otherwise

	var tileData = text.Trim(layer.TileData.Tiles)

	if layer.TileData.Encoding == "base64" {
		var b64 = text.FromBase64(text.Trim(tileData))
		var data = []byte{}

		switch layer.TileData.Compression {
		case "gzip":
			data = storage.DecompressGZIP([]byte(b64))
		case "zlib":
			data = storage.DecompressZLIB([]byte(b64))
		}

		layer.Tiles = bytesToTiles(data)
		return layer.Tiles
	}

	var rows = text.Split(tileData, "\n")
	layer.Tiles = make([]uint32, mapData.Width*mapData.Height)

	for i := 0; i < mapData.Height; i++ {
		var row = rows[i]
		if text.EndsWith(row, ",") {
			row = row[:len(row)-1]
		}

		var columns = text.Split(row, ",")
		for j := 0; j < mapData.Width; j++ {
			var tile = text.ToNumber[uint32](columns[j])
			if tile == 0 {
				continue
			}

			var unoriented = flag.TurnOff(tile, internal.Flips)
			var curTileset = currentTileset(usedTilesets, unoriented)
			if curTileset == nil {
				continue
			}

			var index = number.Indexes2DToIndex1D(i, j, mapData.Width, mapData.Height)
			layer.Tiles[index] = uint32(tile)
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

	for _, group := range data.LayersGroups {
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
	return layer.Name == layerNameOrId || layer.Id == text.ToNumber[uint32](layerNameOrId)

}
func usedTilesets(data *internal.Map) []*internal.Tileset {
	var usedTilesets = make([]*internal.Tileset, len(data.Tilesets))

	for i, tileset := range data.Tilesets {
		if tileset.Source == "" {
			continue
		}

		var tilesetId = path.New(data.Directory, tileset.Source)
		usedTilesets[i] = internal.TiledTilesets[tilesetId]

		if usedTilesets[i] != nil {
			// usedTilesets[i].FirstTileId = tileset.FirstTileId
		}
	}
	return usedTilesets
}
func currentTileset(usedTilesets []*internal.Tileset, tile uint32) *internal.Tileset {
	var curTileset = usedTilesets[0]
	for i := len(usedTilesets) - 1; i >= 0; i-- {
		if usedTilesets[i] != nil && tile > 0 { //usedTilesets[i].FirstTileId {
			curTileset = usedTilesets[i]
			break
		}
	}
	return curTileset
}
func forEachTile(mapId, layerNameOrId string,
	do func(x, y int, id uint32, layer *internal.LayerTiles, curTileset *internal.Tileset)) bool {
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
		//id -= curTileset.FirstTileId - 1 // same as id
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
			return obj
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

func bytesToTiles(data []byte) []uint32 {
	if len(data)%4 != 0 {
		return nil
	}

	var numElements = len(data) / 4
	var result = make([]uint32, numElements)
	var reader = bytes.NewReader(data)
	var err = binary.Read(reader, binary.LittleEndian, &result)
	if err != nil {
		return nil
	}
	return result
}

func getTileOrientation(tileId uint32, w, h float32) (angle float32, newW, newH float32) {
	var flipH = flag.IsOn(tileId, internal.FlipX)
	var flipV = flag.IsOn(tileId, internal.FlipY)
	var flipDiag = flag.IsOn(tileId, internal.FlipDiag)

	angle = 0.0
	newW, newH = w, h

	if flipH && !flipV && flipDiag {
		angle = 90
	} else if flipH && flipV && !flipDiag {
		angle = 180
	} else if !flipH && flipV && flipDiag {
		angle = 270
	} else if flipH && !flipV && !flipDiag {
		newW = -w
	} else if flipH && flipV && flipDiag {
		newW = -w
		angle = 90
	} else if !flipH && flipV && !flipDiag {
		newW = -w
		angle = 180
	} else if !flipH && !flipV && flipDiag {
		newW = -w
		angle = 270
	}
	return angle, newW, newH
}

func tileRenderSize(w, h float32, mapData *internal.Map, tileset *internal.Tileset) (float32, float32) {
	var ratioW, ratioH float32 = 1, 1
	if w > h {
		ratioH = h / w
	} else {
		ratioW = w / h
	}

	if tileset.TileRenderSize == "grid" {
		w, h = float32(mapData.TileWidth), float32(mapData.TileHeight)
	}
	if tileset.FillMode == "preserve-aspect-fit" {
		w *= ratioW
		h *= ratioH
	}
	return w, h
}
