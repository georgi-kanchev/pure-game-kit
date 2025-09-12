package tiled

import (
	"pure-kit/engine/internal"
	"pure-kit/engine/utility/text"
)

const (
	TilesetName       = "name"
	TilesetClass      = "class"
	TilesetTileWidth  = "width"
	TilesetTileHeight = "height"
	TilesetTileCount  = "count"
	TilesetColumns    = "columns"
	TilesetSpacing    = "spacing"
	TilesetOffsetX    = "offsetX"
	TilesetOffsetY    = "offsetY"
	TilesetAtlasId    = "atlasId"
)

func TilesetProperty(tilesetId, property string) string {
	var data, has = internal.TiledTilesets[tilesetId]
	if !has {
		return ""
	}

	switch property {
	case TilesetName:
		return data.Name
	case TilesetClass:
		return data.Class
	case TilesetTileWidth:
		return text.New(data.TileWidth)
	case TilesetTileHeight:
		return text.New(data.TileHeight)
	case TilesetTileCount:
		return text.New(data.TileCount)
	case TilesetColumns:
		return text.New(data.Columns)
	case TilesetSpacing:
		return text.New(data.Spacing)
	case TilesetOffsetX:
		return text.New(data.TileOffset.X)
	case TilesetOffsetY:
		return text.New(data.TileOffset.Y)
	case TilesetAtlasId:
		return data.AtlasId
	}

	for _, v := range data.Properties {
		if v.Name == property {
			return v.Value
		}
	}
	return ""
}
