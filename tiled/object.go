package tiled

import (
	"pure-game-kit/data/path"
	"pure-game-kit/execution/condition"
	"pure-game-kit/graphics"
	"pure-game-kit/internal"
	"pure-game-kit/tiled/property"
	"pure-game-kit/utility/flag"
	"pure-game-kit/utility/point"
	"pure-game-kit/utility/text"
)

type Object struct {
	Properties map[string]any
	Points     [][2]float32

	OwnerTile  *Tile
	OwnerLayer *Layer
}

func (object *Object) Sprite() *graphics.Sprite {
	var tileId = object.Properties[property.ObjectTileId]
	var worldX, worldY float32 = 0, 0
	var offsetX, offsetY float32 = 0, 0
	var curTileset *Tileset = nil
	var firstId uint32 = 1
	var assetId = ""
	var x = object.Properties[property.ObjectX].(float32)
	var y = object.Properties[property.ObjectY].(float32)
	var w = object.Properties[property.ObjectWidth].(float32)
	var h = object.Properties[property.ObjectHeight].(float32)
	var originX, originY float32 = w / 2, -h / 2
	var ang = object.Properties[property.ObjectRotation].(float32)
	var id = flag.TurnOff(tileId.(uint32), internal.FlipX)
	id = flag.TurnOff(id, internal.FlipY)

	if object.OwnerLayer != nil {
		curTileset, firstId = currentTileset(object.OwnerLayer.OwnerMap.Tilesets, id)
		worldX = object.OwnerLayer.OwnerMap.Properties[property.MapWorldX].(float32)
		worldY = object.OwnerLayer.OwnerMap.Properties[property.MapWorldY].(float32)
	} else if object.OwnerTile == nil {
		curTileset = object.OwnerTile.OwnerTileset
		offsetX = object.OwnerTile.OwnerTileset.Properties[property.LayerOffsetX].(float32)
		offsetY = object.OwnerTile.OwnerTileset.Properties[property.LayerOffsetY].(float32)
	}

	if curTileset != nil {
		var asset, hasAsset = curTileset.Properties[property.TilesetAtlasId]
		originX, originY = w/2, h/2

		if hasAsset {
			assetId = path.New(asset.(string), text.New(id-firstId))
		} else {
			for _, tile := range curTileset.Tiles {
				var id = tile.Properties[property.TileId].(uint32)
				if id == tileId.(uint32)-firstId {
					assetId = tile.Properties[property.TileImage].(string)
					break
				}
			}
		}
	}

	var sprite = graphics.NewSprite(assetId, worldX+x, worldY+y)
	sprite.X += originX + offsetX
	sprite.Y = sprite.Y - originY + offsetY
	sprite.Width, sprite.Height = w, h
	sprite.ScaleX = condition.If(flag.IsOn(tileId.(uint32), internal.FlipX), float32(-1), 1)
	sprite.ScaleY = condition.If(flag.IsOn(tileId.(uint32), internal.FlipY), float32(-1), 1)
	sprite.Angle = ang
	return sprite
}

//=================================================================

func newObject(data *internal.LayerObject, ownerTile *Tile, ownerLayer *Layer) *Object {
	var result = Object{OwnerTile: ownerTile, OwnerLayer: ownerLayer}
	result.initProperties(data)
	result.initPoints(data)
	return &result
}

//=================================================================

func (t *Object) initProperties(data *internal.LayerObject) {
	t.Properties = make(map[string]any)
	t.Properties[property.ObjectId] = data.Id
	t.Properties[property.ObjectClass] = data.Class
	t.Properties[property.ObjectTemplate] = data.Template
	t.Properties[property.ObjectName] = data.Name
	t.Properties[property.ObjectVisible] = data.Visible != "false"
	t.Properties[property.ObjectLocked] = data.Locked
	t.Properties[property.ObjectX] = data.X
	t.Properties[property.ObjectY] = data.Y
	t.Properties[property.ObjectWidth] = data.Width
	t.Properties[property.ObjectHeight] = data.Height
	t.Properties[property.ObjectRotation] = data.Rotation
	t.Properties[property.ObjectTileId] = data.Gid
	t.Properties[property.ObjectFlipX] = flag.IsOn(data.Gid, internal.FlipX)
	t.Properties[property.ObjectFlipY] = flag.IsOn(data.Gid, internal.FlipY)

	var owner *Project = nil
	if t.OwnerLayer != nil {
		owner = t.OwnerLayer.OwnerMap.Project
	} else if t.OwnerTile == nil {
		owner = t.OwnerTile.OwnerTileset.Project
	}

	for _, prop := range data.Properties {
		t.Properties[prop.Name] = parseProperty(&prop, owner)
	}
}
func (t *Object) initPoints(data *internal.LayerObject) {
	var ptsData = ""
	if data.Polyline.Points != "" {
		ptsData = data.Polyline.Points
	}
	if data.Polygon.Points != "" {
		ptsData = data.Polygon.Points
	}
	if ptsData == "" {
		var w, h = data.Width, data.Height
		if data.Ellipse != nil {
			const segments = 16
			var rx, ry = w / 2, h / 2
			var step = 360.0 / float32(segments)
			var firstValue = ""

			for i := range segments {
				var cx, cy = point.MoveAtAngle(0, 0, float32(i)*step, 1)
				var value = text.New(cx*rx, ",", cy*ry, " ")
				ptsData += value

				if i == 0 {
					firstValue = value
				}
			}
			ptsData += firstValue
		} else { // rectangle
			ptsData = text.New(0, ",", 0, " ", w, ",", 0, " ", w, ",", h, " ", 0, ",", h)
		}
	}

	var points = [][2]float32{}
	var pts = text.Split(text.Trim(ptsData), " ")
	for _, pt := range pts {
		var xy = text.Split(pt, ",")
		if len(xy) == 2 {
			var x, y = text.ToNumber[float32](xy[0]), text.ToNumber[float32](xy[1])
			x, y = point.RotateAroundPoint(x, y, 0, 0, data.Rotation)
			points = append(points, [2]float32{x, y})
		}
	}

	t.Points = points
}
