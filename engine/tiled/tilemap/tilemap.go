package tilemap

import (
	"path"
	"pure-kit/engine/data/file"
	"pure-kit/engine/graphics"
	"pure-kit/engine/internal"
	"pure-kit/engine/utility/text"
	"strconv"
	"strings"
)

const (
	PropertyName            = "name"
	PropertyClass           = "class"
	PropertyTileColumns     = "columns"
	PropertyTileRows        = "rows"
	PropertyTileWidth       = "tileWidth"
	PropertyTileHeight      = "tileHeight"
	PropertyParallaxX       = "parallaxX"
	PropertyParallaxY       = "parallaxY"
	PropertyInfinite        = "infinite"
	PropertyBackgroundColor = "backgroundColor"
)

func Property(mapId, property string) string {
	var data, has = internal.TiledMaps[mapId]
	if !has {
		return ""
	}

	switch property {
	case PropertyName:
		return data.Name
	case PropertyClass:
		return data.Class
	case PropertyTileWidth:
		return text.New(data.TileWidth)
	case PropertyTileHeight:
		return text.New(data.TileHeight)
	case PropertyTileColumns:
		return text.New(data.Width)
	case PropertyTileRows:
		return text.New(data.Height)
	case PropertyParallaxX:
		return text.New(data.ParallaxOriginX)
	case PropertyParallaxY:
		return text.New(data.ParallaxOriginY)
	case PropertyInfinite:
		return text.New(data.Infinite)
	case PropertyBackgroundColor:
		return text.New(color(data.BackgroundColor))
	}

	for _, v := range data.Properties {
		if v.Name == property {
			return v.Value
		}
	}
	return ""
}

func LayerTiles(mapId, layerName string) []*graphics.Sprite {
	var data, has = internal.TiledMaps[mapId]
	if !has {
		return []*graphics.Sprite{}
	}

	var wantedLayer *internal.LayerTiles
	for _, layer := range data.LayersTiles {
		if layer.Name == layerName {
			wantedLayer = &layer
			break
		}
	}
	for _, group := range data.Groups {
		for _, layer := range group.LayersTiles {
			if layer.Name == layerName {
				wantedLayer = &layer
				break
			}
		}
	}
	if wantedLayer == nil {
		return []*graphics.Sprite{}
	}

	var result = make([]*graphics.Sprite, 0, data.Width*data.Height)
	var tileData = strings.Trim(wantedLayer.TileData.Tiles, "\n")
	var rows = strings.Split(tileData, "\n")
	var usedTilesets = make([]*internal.Tileset, len(data.Tilesets))

	for i, tileset := range data.Tilesets {
		if tileset.Source != "" {
			var tilesetId = path.Join(data.Directory, tileset.Source)
			tilesetId = path.Base(strings.ReplaceAll(tilesetId, file.Extension(tilesetId), ""))
			tilesetId = path.Join(data.Directory, tilesetId)
			usedTilesets[i] = internal.TiledTilesets[tilesetId]
			if usedTilesets[i] != nil {
				usedTilesets[i].FirstTileId = tileset.FirstTileId
			}
			continue
		}

		usedTilesets[i] = &tileset
	}

	for i := 0; i < data.Height; i++ {
		var row = strings.Trim(rows[i], ",")
		var columns = strings.Split(row, ",")
		for j := 0; j < data.Width; j++ {
			var tile = int(text.ToNumber(columns[j]))
			var curTileset = usedTilesets[0]

			if tile == 0 {
				continue
			}

			for i := len(usedTilesets) - 1; i >= 0; i-- {
				if usedTilesets[i] != nil && tile > usedTilesets[i].FirstTileId {
					curTileset = usedTilesets[i]
					break
				}
			}

			if curTileset == nil {
				continue
			}

			var tileId = text.New(curTileset.AtlasId, "[", tile-curTileset.FirstTileId, "]")
			var x = float32(j)*float32(data.TileWidth) + data.WorldX
			var y = float32(i)*float32(data.TileHeight) + data.WorldY
			var sprite = graphics.NewSprite(tileId, x, y)

			sprite.Width, sprite.Height = float32(data.TileWidth), float32(data.TileHeight)
			sprite.PivotX, sprite.PivotY = 0, 0
			result = append(result, &sprite)
		}
	}

	return result
}

//=================================================================
// private

func color(hex string) uint {
	var trimmed = strings.TrimPrefix(hex, "#")

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
