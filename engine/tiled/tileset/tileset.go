package tileset

import (
	"pure-kit/engine/internal"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/text"
)

const (
	PropertyName       = "name"
	PropertyClass      = "class"
	PropertyTileWidth  = "width"
	PropertyTileHeight = "height"
	PropertyTileCount  = "count"
	PropertyColumns    = "columns"
	PropertySpacing    = "spacing"
	PropertyOffsetX    = "offsetX"
	PropertyOffsetY    = "offsetY"
	PropertyAtlasId    = "atlasId"
)

func Property(tilesetId, property string) string {
	var data, has = internal.TiledTilesets[tilesetId]
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
	case PropertyTileCount:
		return text.New(data.TileCount)
	case PropertyColumns:
		return text.New(data.Columns)
	case PropertySpacing:
		return text.New(data.Spacing)
	case PropertyOffsetX:
		return text.New(data.Offset.X)
	case PropertyOffsetY:
		return text.New(data.Offset.Y)
	case PropertyAtlasId:
		return data.AtlasId
	}

	for _, v := range data.Properties {
		if v.Name == property {
			return v.Value
		}
	}
	return ""
}

func TileProperty(tilesetId string, tileId int, property string) string {
	var tile = getTile(tilesetId, tileId)
	if tile == nil {
		return ""
	}

	for _, prop := range tile.Properties {
		if prop.Name == property {
			return prop.Value
		}
	}
	return ""
}
func TileShapeProperty(tilesetId string, tileId int, shapeNameOrId, property string) string {
	var obj = getObj(tilesetId, tileId, shapeNameOrId)
	if obj == nil {
		return ""
	}

	for _, prop := range obj.Properties {
		if prop.Name == property {
			return prop.Value
		}
	}
	return ""
}
func TileAnimationTileIds(tilesetId string, tileId int) (frameTileIds []int) {
	var tile = getTile(tilesetId, tileId)
	var result = make([]int, len(tile.Animation.Frames))
	for i, frame := range tile.Animation.Frames {
		result[i] = frame.TileId
	}
	return result
}
func TileAnimationDurations(tilesetId string, tileId int) (frameDurations []float32) {
	var tile = getTile(tilesetId, tileId)
	var result = make([]float32, len(tile.Animation.Frames))
	for i, frame := range tile.Animation.Frames {
		result[i] = float32(frame.Duration) / 1000 // ms -> sec
	}
	return result
}

//=================================================================
// private

func getTile(tilesetId string, tileId int) *internal.TilesetTile {
	var data, has = internal.TiledTilesets[tilesetId]
	if !has {
		return nil
	}

	for _, tile := range data.Tiles {
		if tile.Id == tileId {
			return &tile
		}
	}

	return nil
}
func getObj(tilesetId string, tileId int, shapeNameOrId string) *internal.LayerObject {
	var tile = getTile(tilesetId, tileId)
	if tile == nil {
		return nil
	}

	var layer = tile.CollisionLayers[0]
	for _, obj := range layer.Objects {
		if obj.Name == shapeNameOrId {
			return obj
		}
		var id = text.ToNumber(shapeNameOrId)
		if !number.IsNaN(id) && obj.Id == int(id) {
			return obj
		}
	}
	return nil
}
