package tileset

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/data/path"
	"pure-game-kit/execution/condition"
	"pure-game-kit/geometry"
	"pure-game-kit/internal"
	p "pure-game-kit/tiled/property"
	"pure-game-kit/utility/flag"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/point"
	"pure-game-kit/utility/text"
)

func Property(tilesetId, property string) string {
	var data, has = internal.TiledTilesets[tilesetId]
	if !has {
		return ""
	}

	switch property {
	case p.TilesetName:
		return data.Name
	case p.TilesetClass:
		return data.Class
	case p.TilesetTileWidth:
		return text.New(data.TileWidth)
	case p.TilesetTileHeight:
		return text.New(data.TileHeight)
	case p.TilesetColumns:
		return text.New(data.Columns)
	case p.TilesetSpacing:
		return text.New(data.Spacing)
	case p.TilesetOffsetX:
		return text.New(data.Offset.X)
	case p.TilesetOffsetY:
		return text.New(data.Offset.Y)
	case p.TilesetAtlasId:
		return data.AssetId
	}

	// for _, v := range data.Properties {
	// 	if v.Name == property {
	// 		return v.Value
	// 	}
	// }
	return ""
}

func TileProperty(tilesetId string, tileId uint32, property string) string {
	var tile = getTile(tilesetId, tileId)
	if tile == nil {
		return ""
	}

	// for _, prop := range tile.Properties {
	// 	if prop.Name == property {
	// 		return prop.Value
	// 	}
	// }
	return ""
}
func TileAnimationTileIds(tilesetId string, tileId uint32) (frameTileIds []uint32) {
	var tile = getTile(tilesetId, tileId)
	if tile == nil {
		return []uint32{}
	}

	var result = make([]uint32, len(tile.Animation.Frames))
	for i, frame := range tile.Animation.Frames {
		result[i] = frame.TileId
	}
	return result
}
func TileAnimationDurations(tilesetId string, tileId uint32) (frameDurations []float32) {
	var tile = getTile(tilesetId, tileId)
	if tile == nil {
		return []float32{}
	}

	var result = make([]float32, len(tile.Animation.Frames))
	for i, frame := range tile.Animation.Frames {
		result[i] = float32(frame.Duration) / 1000 // ms -> sec
	}
	return result
}
func TileAnimate(tilesetId string, tileId uint32, animate bool) {
	var tile = getTile(tilesetId, tileId)
	if tile != nil {
		var tileset, _ = internal.TiledTilesets[tilesetId]
		tile.IsAnimating = animate

		if !animate { // disabling animation resets the tile to original one
			var w, h = tileset.Columns, tileset.TileCount / tileset.Columns
			var x, y = number.Index1DToIndexes2D(tile.Id, uint32(w), uint32(h))
			var rectId = path.New(tileset.AssetId, text.New(tile.Id))
			assets.SetTextureAtlasTile(tileset.AssetId, rectId, float32(x), float32(y), 1, 1, 0, false)
		}
	}
}
func TileIsAnimated(tilesetId string, tileId uint32) bool {
	var tile = getTile(tilesetId, tileId)
	return tile != nil && tile.Update != nil
}

func TileObjectProperty(tilesetId string, tileId uint32, objectNameClassOrId, property string) string {
	var obj = getObj(tilesetId, tileId, objectNameClassOrId)
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
		return text.New(flag.TurnOff(id, internal.FlipY))
	}

	// for _, prop := range obj.Properties {
	// 	if prop.Name == property {
	// 		return prop.Value
	// 	}
	// }
	return ""
}
func TileObjectShapes(tilesetId string, tileId uint32, objectNameClassOrId string) []*geometry.Shape {
	var result = []*geometry.Shape{}
	var tile = getTile(tilesetId, tileId)
	if tile == nil {
		return result
	}

	var objs = tile.CollisionLayers[0].Objects
	var tileset = internal.TiledTilesets[tilesetId]
	for _, obj := range objs {
		if objectNameClassOrId != "" && obj.Name != objectNameClassOrId && obj.Class != objectNameClassOrId &&
			text.New(obj.Id) != objectNameClassOrId {
			continue
		}
		var ptsData = ""
		if obj.Polyline.Points != "" {
			ptsData = obj.Polyline.Points
		}
		if obj.Polygon.Points != "" {
			ptsData = obj.Polygon.Points
		}
		if ptsData == "" {
			var w, h = obj.Width, obj.Height
			ptsData = text.New(0, ",", 0, " ", w, ",", 0, " ", w, ",", h, " ", 0, ",", h)
		}
		var corners = [][2]float32{}
		var pts = text.Split(ptsData, " ")
		for _, pt := range pts {
			var xy = text.Split(pt, ",")
			if len(xy) == 2 {
				var x, y = text.ToNumber[float32](xy[0]), text.ToNumber[float32](xy[1])
				x, y = point.RotateAroundPoint(x, y, 0, 0, obj.Rotation)
				corners = append(corners, [2]float32{x, y})
			}
		}
		var shape = geometry.NewShapeCorners(corners...)
		shape.X = obj.X - float32(tileset.TileWidth)/2
		shape.Y = obj.Y - float32(tileset.TileHeight)/2
		result = append(result, shape)
	}

	return result
}
func TileObjectPoints(tilesetId string, tileId uint32, objectNameClassOrId string) [][2]float32 {
	var points = [][2]float32{}
	var tile = getTile(tilesetId, tileId)
	if tile == nil {
		return points
	}
	var objs = tile.CollisionLayers[0].Objects
	for _, obj := range objs {
		var isPoint = obj.Width == 0 && obj.Height == 0 && obj.Polygon.Points == ""
		if !isPoint {
			continue
		}
		if objectNameClassOrId == "" || obj.Name == objectNameClassOrId || obj.Class == objectNameClassOrId ||
			text.New(obj.Id) == objectNameClassOrId {
			points = append(points, [2]float32{obj.X, obj.Y})
		}
	}
	return points
}

//=================================================================
// private

func getTile(tilesetId string, tileId uint32) *internal.TilesetTile {
	var data, has = internal.TiledTilesets[tilesetId]
	if has {
		return data.MappedTiles[tileId]
	}

	return nil
}
func getObj(tilesetId string, tileId uint32, objectNameClassOrId string) *internal.LayerObject {
	var tile = getTile(tilesetId, tileId)
	if tile == nil {
		return nil
	}

	var layer = tile.CollisionLayers[0]
	for _, obj := range layer.Objects {
		if obj.Name == objectNameClassOrId || obj.Class == objectNameClassOrId ||
			text.New(obj.Id) == objectNameClassOrId {
			return &obj
		}
	}
	return nil
}
