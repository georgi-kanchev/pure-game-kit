package assets

import (
	"bytes"
	"encoding/binary"
	"pure-game-kit/data/path"
	"pure-game-kit/data/storage"
	"pure-game-kit/execution/condition"
	"pure-game-kit/internal"
	"pure-game-kit/utility/flag"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/text"
)

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

		layer[0].Objects[i] = &newObj
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

func tryCacheLayerTileIds(data *internal.Map, layers *internal.Layers) {
	for _, layer := range layers.LayersTiles {
		tryCacheTileIds(data, layer)
	}
	for _, group := range layers.LayersGroups {
		tryCacheLayerTileIds(data, &group.Layers)
	}
}
func tryCacheTileIds(mapData *internal.Map, layer *internal.LayerTiles) {
	if layer.Tiles != nil {
		return // already cached
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
		return
	}

	var rows = text.Split(tileData, "\n")
	layer.Tiles = make([]uint32, mapData.Width*mapData.Height)
	var usedTilesets = usedTilesets(mapData)

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
			usedTilesets[i].FirstTileId = tileset.FirstTileId
		}
	}
	return usedTilesets
}
func currentTileset(usedTilesets []*internal.Tileset, tile uint32) *internal.Tileset {
	var curTileset = usedTilesets[0]
	for i := len(usedTilesets) - 1; i >= 0; i-- {
		if usedTilesets[i] != nil && tile > usedTilesets[i].FirstTileId {
			curTileset = usedTilesets[i]
			break
		}
	}
	return curTileset
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
