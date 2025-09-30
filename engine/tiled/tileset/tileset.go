package tileset

import (
	"pure-kit/engine/execution/condition"
	"pure-kit/engine/geometry"
	"pure-kit/engine/geometry/point"
	"pure-kit/engine/internal"
	p "pure-kit/engine/tiled/property"
	"pure-kit/engine/utility/flag"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/text"
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
	case p.TilesetTileCount:
		return text.New(data.TileCount)
	case p.TilesetColumns:
		return text.New(data.Columns)
	case p.TilesetSpacing:
		return text.New(data.Spacing)
	case p.TilesetOffsetX:
		return text.New(data.Offset.X)
	case p.TilesetOffsetY:
		return text.New(data.Offset.Y)
	case p.TilesetAtlasId:
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
func TileAnimationTileIds(tilesetId string, tileId int) (frameTileIds []int) {
	var tile = getTile(tilesetId, tileId)
	if tile == nil {
		return []int{}
	}

	var result = make([]int, len(tile.Animation.Frames))
	for i, frame := range tile.Animation.Frames {
		result[i] = frame.TileId
	}
	return result
}
func TileAnimationDurations(tilesetId string, tileId int) (frameDurations []float32) {
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

func TileObjectProperty(tilesetId string, tileId int, objectNameClassOrId, property string) string {
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

	for _, prop := range obj.Properties {
		if prop.Name == property {
			return prop.Value
		}
	}
	return ""
}
func TileObjectShapes(tilesetId string, tileId int, objectNameOrClass string) []*geometry.Shape {
	var shapes = []*geometry.Shape{}
	var tile = getTile(tilesetId, tileId)
	if tile == nil {
		return shapes
	}

	var objs = tile.CollisionLayers[0].Objects
	var tileset = internal.TiledTilesets[tilesetId]
	for _, obj := range objs {
		if objectNameOrClass != "" && obj.Name != objectNameOrClass && obj.Class != objectNameOrClass {
			continue
		}
		var ptsData = ""
		if obj.PolygonTile.Points != "" {
			ptsData = obj.PolygonTile.Points
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
				var x, y = text.ToNumber(xy[0]), text.ToNumber(xy[1])
				x, y = point.RotateAroundPoint(x, y, 0, 0, obj.Rotation)
				corners = append(corners, [2]float32{x, y})
			}
		}
		var shape = geometry.NewShapeCorners(corners...)
		shape.X = obj.X - float32(tileset.TileWidth)/2
		shape.Y = obj.Y - float32(tileset.TileHeight)/2
		shapes = append(shapes, shape)
	}

	return shapes
}
func TileObjectPoints(tilesetId string, tileId int, objectNameOrClass string) [][2]float32 {
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
		if objectNameOrClass == "" || obj.Name == objectNameOrClass || obj.Class == objectNameOrClass {
			points = append(points, [2]float32{obj.X, obj.Y})
		}
	}
	return points
}

//=================================================================
// private

func getTile(tilesetId string, tileId int) *internal.TilesetTile {
	var data, has = internal.TiledTilesets[tilesetId]
	if has {
		return data.MappedTiles[tileId]
	}

	return nil
}
func getObj(tilesetId string, tileId int, objectNameOrId string) *internal.LayerObject {
	var tile = getTile(tilesetId, tileId)
	if tile == nil {
		return nil
	}

	var layer = tile.CollisionLayers[0]
	for _, obj := range layer.Objects {
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
