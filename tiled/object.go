package tiled

import (
	"pure-game-kit/data/path"
	"pure-game-kit/execution/condition"
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/internal"
	"pure-game-kit/tiled/property"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/flag"
	"pure-game-kit/utility/is"
	"pure-game-kit/utility/point"
	"pure-game-kit/utility/text"
)

type Object struct {
	Properties map[string]any
	Corners    [][2]float32 // used to describe either shapes, lines or points

	OwnerTile  *Tile
	OwnerLayer *Layer
}

func (object *Object) Sprite() *graphics.Sprite {
	var objType = object.Properties[property.ObjectType]
	if objType != "tile" {
		return nil
	}

	var pivotY float32 = 0
	var assetId = ""
	var x, y, w, h, ang = object.getArea()
	var tileId = object.Properties[property.ObjectTileId]
	var flipX, flipY = flag.IsOn(tileId.(uint32), internal.FlipX), flag.IsOn(tileId.(uint32), internal.FlipY)
	var flipOffX, flipOffY = condition.If(flipX, w, 0), condition.If(flipY, -h, 0)
	var worldX, worldY, offsetX, offsetY = object.getOffsets()
	var tile = object.getTile()

	if tile.OwnerTileset != nil {
		var asset, hasAsset = tile.OwnerTileset.Properties[property.TilesetAtlasId]
		pivotY = 1

		if hasAsset {
			assetId = path.New(asset.(string), text.New(tile.Properties[property.TileId]))
		} else {
			assetId = tile.Properties[property.TileImage].(string)
		}
	}

	var sprite = graphics.NewSprite(assetId, worldX+offsetX+x+flipOffX, worldY+offsetY+y+flipOffY)
	sprite.Width, sprite.Height = w, h
	sprite.ScaleX = condition.If(flipX, float32(-1), 1)
	sprite.ScaleY = condition.If(flipY, float32(-1), 1)
	sprite.PivotX, sprite.PivotY = 0, pivotY
	sprite.Angle = ang
	return sprite
}
func (object *Object) TextBox() *graphics.TextBox {
	var objType = object.Properties[property.ObjectType]
	if objType != "text" {
		return nil
	}

	var font = object.Properties[property.ObjectTextFont].(string)
	var txt = object.Properties[property.ObjectText]
	var x, y, w, h, ang = object.getArea()
	var bold = object.Properties[property.ObjectTextBold].(bool)

	for id := range internal.Fonts {
		var name = path.LastPart(path.RemoveExtension(id))
		if text.LowerCase(font) == text.LowerCase(name) {
			font = id // searching font in loaded fonts by name (case insensitive)
		}
	}

	var textBox = graphics.NewTextBox(font, x, y, txt)
	textBox.Width, textBox.Height = w, h
	textBox.Angle = ang
	textBox.AlignmentX = object.Properties[property.ObjectTextAlignX].(float32)
	textBox.AlignmentY = object.Properties[property.ObjectTextAlignY].(float32)
	textBox.WordWrap = object.Properties[property.ObjectTextWordWrap].(bool)
	textBox.Color = object.Properties[property.ObjectTextColor].(uint)
	textBox.Thickness = condition.If(bold, float32(0.8), 0.5)
	textBox.LineHeight = float32(object.Properties[property.ObjectTextFontSize].(int))
	textBox.PivotX, textBox.PivotY = 0, 0
	return textBox
}
func (object *Object) Shapes() []*geometry.Shape {
	var result = []*geometry.Shape{}
	var objType = object.Properties[property.ObjectType]
	if is.OneOf(objType, "text", "point", "line") {
		return result
	}

	var x, y, _, h, ang = object.getArea()
	var worldX, worldY, offsetX, offsetY = object.getOffsets()

	if objType == "tile" {
		result = append(result, object.getTile().Shapes()...)
		for _, shape := range result {
			shape.X += worldX + offsetX + x
			shape.Y += worldY + offsetY + y - h
		}
		return result
	}

	var shape = geometry.NewShapeCorners(object.Corners...)
	shape.Angle = ang
	shape.X, shape.Y = worldX+offsetX+x, worldY+offsetY+y
	result = append(result, shape)
	return result
}
func (object *Object) Lines() [][2]float32 {
	var result = [][2]float32{}
	var objType = object.Properties[property.ObjectType]
	if objType != "line" && objType != "tile" {
		return result
	}

	var x, y, _, h, ang = object.getArea()
	var worldX, worldY, offsetX, offsetY = object.getOffsets()

	if objType == "tile" {
		result = object.getTile().Lines()
		for i := range result {
			result[i][0] += worldX + offsetX + x
			result[i][1] += worldY + offsetY + y - h

			var rotX, rotY = point.RotateAroundPoint(result[i][0], result[i][1], x, y, ang)
			result[i][0] = rotX
			result[i][1] = rotY
		}
		return result
	}

	for i, pt := range object.Corners {
		var xy = [2]float32{worldX + offsetX + x + pt[0], worldY + offsetY + y + pt[1]}

		if i > 0 {
			xy[0], xy[1] = point.RotateAroundPoint(xy[0], xy[1], result[0][0], result[0][1], ang)
		}
		result = append(result, xy)
	}

	return result
}
func (object *Object) Points() [][2]float32 {
	var result = [][2]float32{}
	var objType = object.Properties[property.ObjectType]
	if objType != "point" && objType != "tile" {
		return result
	}

	var x, y, _, h, _ = object.getArea()
	var worldX, worldY, offsetX, offsetY = object.getOffsets()

	if objType == "tile" {
		var result = object.getTile().Points()
		for i := range result {
			result[i][0] += worldX + offsetX + x
			result[i][1] += worldY + offsetY + y - h
		}
		return result
	}

	return [][2]float32{{worldX + offsetX + x, worldY + offsetY + y}}
}

func (object *Object) Draw(camera *graphics.Camera) {
	var sprs = []*graphics.Sprite{object.Sprite()}
	var txts = []*graphics.TextBox{object.TextBox()}
	draw(camera, sprs, txts, object.Shapes(), object.Points(), object.Lines(), color.White)
}

//=================================================================

var aligns = map[string]float32{"left": 0, "center": 0.5, "right": 1, "top": 0, "bottom": 1, "justify": 0}

func newObject(data *internal.LayerObject, ownerTile *Tile, ownerLayer *Layer) *Object {
	var result = Object{OwnerTile: ownerTile, OwnerLayer: ownerLayer}
	result.initProperties(data)
	result.initCorners(data)
	return &result
}

func (object *Object) initProperties(data *internal.LayerObject) {
	object.Properties = make(map[string]any)
	object.Properties[property.ObjectId] = data.Id
	object.Properties[property.ObjectClass] = data.Class
	object.Properties[property.ObjectTemplate] = data.Template
	object.Properties[property.ObjectName] = data.Name
	object.Properties[property.ObjectVisible] = data.Visible != "false"
	object.Properties[property.ObjectLocked] = data.Locked
	object.Properties[property.ObjectX] = data.X
	object.Properties[property.ObjectY] = data.Y
	object.Properties[property.ObjectWidth] = data.Width
	object.Properties[property.ObjectHeight] = data.Height
	object.Properties[property.ObjectRotation] = data.Rotation
	object.Properties[property.ObjectTileId] = data.Gid
	object.Properties[property.ObjectFlipX] = flag.IsOn(data.Gid, internal.FlipX)
	object.Properties[property.ObjectFlipY] = flag.IsOn(data.Gid, internal.FlipY)

	if data.Text != nil {
		object.Properties[property.ObjectText] = data.Text.Value
		object.Properties[property.ObjectTextFont] = data.Text.FontFamily
		object.Properties[property.ObjectTextFontSize] = data.Text.FontSize
		object.Properties[property.ObjectTextBold] = data.Text.Bold
		object.Properties[property.ObjectTextItalic] = data.Text.Italic
		object.Properties[property.ObjectTextStrikeout] = data.Text.Strikeout
		object.Properties[property.ObjectTextUnderline] = data.Text.Underline
		object.Properties[property.ObjectTextAlignX] = aligns[data.Text.AlignX]
		object.Properties[property.ObjectTextAlignY] = aligns[data.Text.AlignY]
		object.Properties[property.ObjectTextColor] = color.Hex(data.Text.Color)
		object.Properties[property.ObjectTextWordWrap] = data.Text.WordWrap
		object.Properties[property.ObjectType] = "text"
	} else if data.Gid > 0 {
		object.Properties[property.ObjectType] = "tile"
	} else if data.Ellipse != nil {
		object.Properties[property.ObjectType] = "ellipse"
	} else if data.Point != nil {
		object.Properties[property.ObjectType] = "point"
	} else if data.Polyline != nil {
		object.Properties[property.ObjectType] = "line"
	} else if data.Polygon != nil {
		object.Properties[property.ObjectType] = "polygon"
	} else {
		object.Properties[property.ObjectType] = "rectangle"
	}

	var owner *Project = nil
	if object.OwnerLayer != nil {
		owner = object.OwnerLayer.OwnerMap.Project
	} else if object.OwnerTile == nil {
		owner = object.OwnerTile.OwnerTileset.Project
	}

	for _, prop := range data.Properties {
		object.Properties[prop.Name] = parseProperty(prop, owner)
	}
}
func (object *Object) initCorners(data *internal.LayerObject) {
	var ptsData = ""
	if data.Polyline != nil {
		ptsData = data.Polyline.Points
	}
	if data.Polygon != nil {
		ptsData = data.Polygon.Points
	}
	if ptsData == "" {
		var w, h = data.Width, data.Height
		if data.Point != nil {
			// no ptsData for a single point
		} else if data.Ellipse != nil {
			const segments = 16
			var rx, ry = w / 2, h / 2
			var step = 360.0 / float32(segments)

			for i := range segments {
				var cx, cy = point.MoveAtAngle(0, 0, float32(i)*step, 1)
				var x, y = (cx + 1) * rx, (cy + 1) * ry // shift from center-based to tiled's top-left-based
				var value = text.New(x, ",", y, " ")
				ptsData += value
			}
		} else { // assume it's a rectangle
			ptsData = text.New(0, ",", 0, " ", w, ",", 0, " ", w, ",", h, " ", 0, ",", h)
		}
	}

	var corners = [][2]float32{}
	var pts = text.Split(text.Trim(ptsData), " ")
	for _, pt := range pts {
		var xy = text.Split(pt, ",")
		if len(xy) == 2 {
			var x, y = text.ToNumber[float32](xy[0]), text.ToNumber[float32](xy[1])
			corners = append(corners, [2]float32{x, y})
		}
	}

	object.Corners = corners
}

func (object *Object) getTile() *Tile {
	var tileId = object.Properties[property.ObjectTileId]
	var id = flag.TurnOff(tileId.(uint32), internal.Flips)
	var curTileset *Tileset = nil
	var firstId uint32 = 1
	if object.OwnerLayer != nil {
		curTileset, firstId = currentTileset(object.OwnerLayer.OwnerMap, id)
	} else if object.OwnerTile == nil {
		curTileset = object.OwnerTile.OwnerTileset
	}

	if id == 0 {
		return nil
	}
	return curTileset.Tiles[id-firstId]
}
func (object *Object) getOffsets() (worldX, worldY, layerX, layerY float32) {
	if object.OwnerLayer != nil {
		worldX, worldY, layerX, layerY = object.OwnerLayer.getOffsets()
	} else if object.OwnerTile == nil {
		layerX = object.OwnerTile.OwnerTileset.Properties[property.TilesetOffsetX].(float32)
		layerY = object.OwnerTile.OwnerTileset.Properties[property.TilesetOffsetY].(float32)
	}
	return
}
func (object *Object) getArea() (x, y, w, h, a float32) {
	x = object.Properties[property.ObjectX].(float32)
	y = object.Properties[property.ObjectY].(float32)
	w = object.Properties[property.ObjectWidth].(float32)
	h = object.Properties[property.ObjectHeight].(float32)
	a = object.Properties[property.ObjectRotation].(float32)
	return
}
