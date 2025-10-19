package tileset

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/execution/condition"
	"pure-game-kit/execution/flow"
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
func TileAnimate(tilesetId string, tileId int, animate bool) {
	var tile = getTile(tilesetId, tileId)
	if tile != nil {
		var tileset, _ = internal.TiledTilesets[tilesetId]
		var seq = tile.Sequence.(*flow.Sequence)
		seq.GoToStep(condition.If(animate, 0, -1))

		if !animate { // disabling animation resets the tile to original one
			var w, h = tileset.Columns, tileset.TileCount / tileset.Columns
			var x, y = number.Index1DToIndexes2D(tile.Id, w, h)
			var rectId = text.New(tileset.AtlasId, "/", tile.Id)
			assets.SetTextureAtlasTile(tileset.AtlasId, rectId, float32(x), float32(y), 1, 1, 0, false)
		}
	}
}
func TileIsAnimated(tilesetId string, tileId int) bool {
	var tile = getTile(tilesetId, tileId)
	return tile != nil && tile.Sequence != nil
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
func TileObjectShapes(tilesetId string, tileId int, objectNameClassOrId string) []*geometry.Shape {
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
		result = append(result, shape)
	}

	return result
}
func TileObjectPoints(tilesetId string, tileId int, objectNameClassOrId string) [][2]float32 {
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

func getTile(tilesetId string, tileId int) *internal.TilesetTile {
	var data, has = internal.TiledTilesets[tilesetId]
	if has {
		return data.MappedTiles[tileId]
	}

	return nil
}
func getObj(tilesetId string, tileId int, objectNameClassOrId string) *internal.LayerObject {
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
