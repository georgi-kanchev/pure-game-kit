package tiled

import (
	"pure-game-kit/data/path"
	"pure-game-kit/execution/condition"
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/internal"
	"pure-game-kit/tiled/property"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/color/palette"
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

func (o *Object) ExtractSprite() *graphics.Sprite {
	var objType = o.Properties[property.ObjectType]
	if objType != "tile" {
		return nil
	}

	var pivotY float32 = 0
	var assetId = ""
	var x, y, w, h, ang = o.getArea()
	var tileId = o.Properties[property.ObjectTileId]
	var flipX, flipY = flag.IsOn(tileId.(uint32), internal.FlipX), flag.IsOn(tileId.(uint32), internal.FlipY)
	var flipOffX, flipOffY = condition.If(flipX, w, 0), condition.If(flipY, -h, 0)
	var worldX, worldY, offsetX, offsetY = o.getOffsets()
	var tile = o.getTile()

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
func (o *Object) ExtractTextBox() *graphics.TextBox {
	var objType = o.Properties[property.ObjectType]
	if objType != "text" {
		return nil
	}

	var font = o.Properties[property.ObjectTextFont].(string)
	var txt = o.Properties[property.ObjectText]
	var x, y, w, h, ang = o.getArea()
	var bold = o.Properties[property.ObjectTextBold].(bool)

	for id := range internal.Fonts {
		var name = path.LastPart(path.RemoveExtension(id))
		if text.ToLowerCase(font) == text.ToLowerCase(name) {
			font = id // searching font in loaded fonts by name (case insensitive)
		}
	}

	var textBox = graphics.NewTextBox(font, x, y, txt)
	textBox.Width, textBox.Height = w, h
	textBox.Angle = ang
	textBox.AlignmentX = o.Properties[property.ObjectTextAlignX].(float32)
	textBox.AlignmentY = o.Properties[property.ObjectTextAlignY].(float32)
	textBox.WordWrap = o.Properties[property.ObjectTextWordWrap].(bool)
	textBox.Tint = o.Properties[property.ObjectTextColor].(uint)
	textBox.Thickness = condition.If(bold, float32(0.8), 0.5)
	textBox.LineHeight = float32(o.Properties[property.ObjectTextFontSize].(int))
	textBox.PivotX, textBox.PivotY = 0, 0
	return textBox
}
func (o *Object) ExtractShapes() []*geometry.Shape {
	var result = []*geometry.Shape{}
	var objType = o.Properties[property.ObjectType]
	if is.OneOf(objType, "text", "point", "line") {
		return result
	}

	var x, y, _, h, ang = o.getArea()
	var worldX, worldY, offsetX, offsetY = o.getOffsets()

	if objType == "tile" {
		result = append(result, o.getTile().ExtractShapes()...)
		for _, shape := range result {
			shape.X += worldX + offsetX + x
			shape.Y += worldY + offsetY + y - h
		}
		return result
	}

	var shape = geometry.NewShapeCorners(o.Corners...)
	shape.Angle = ang
	shape.X, shape.Y = worldX+offsetX+x, worldY+offsetY+y
	result = append(result, shape)
	return result
}
func (o *Object) ExtractLines() [][2]float32 {
	var result = [][2]float32{}
	var objType = o.Properties[property.ObjectType]
	if objType != "line" && objType != "tile" {
		return result
	}

	var x, y, _, h, ang = o.getArea()
	var worldX, worldY, offsetX, offsetY = o.getOffsets()

	if objType == "tile" {
		result = o.getTile().ExtractLines()
		for i := range result {
			result[i][0] += worldX + offsetX + x
			result[i][1] += worldY + offsetY + y - h

			var rotX, rotY = point.RotateAroundPoint(result[i][0], result[i][1], x, y, ang)
			result[i][0] = rotX
			result[i][1] = rotY
		}
		return result
	}

	for i, pt := range o.Corners {
		var xy = [2]float32{worldX + offsetX + x + pt[0], worldY + offsetY + y + pt[1]}

		if i > 0 {
			xy[0], xy[1] = point.RotateAroundPoint(xy[0], xy[1], result[0][0], result[0][1], ang)
		}
		result = append(result, xy)
	}

	return result
}
func (o *Object) ExtractPoints() [][2]float32 {
	var result = [][2]float32{}
	var objType = o.Properties[property.ObjectType]
	if objType != "point" && objType != "tile" {
		return result
	}

	var x, y, _, h, _ = o.getArea()
	var worldX, worldY, offsetX, offsetY = o.getOffsets()

	if objType == "tile" {
		var result = o.getTile().ExtractPoints()
		for i := range result {
			result[i][0] += worldX + offsetX + x
			result[i][1] += worldY + offsetY + y - h
		}
		return result
	}

	return [][2]float32{{worldX + offsetX + x, worldY + offsetY + y}}
}

//=================================================================

func (o *Object) Draw(camera *graphics.Camera) {
	var sprs = []*graphics.Sprite{o.ExtractSprite()}
	var txts = []*graphics.TextBox{o.ExtractTextBox()}
	draw(camera, sprs, txts, o.ExtractShapes(), o.ExtractPoints(), o.ExtractLines(), palette.White)
}

//=================================================================

var aligns = map[string]float32{"left": 0, "center": 0.5, "right": 1, "top": 0, "bottom": 1, "justify": 0}

func newObject(data *internal.LayerObject, ownerTile *Tile, ownerLayer *Layer) *Object {
	var result = Object{OwnerTile: ownerTile, OwnerLayer: ownerLayer}
	result.initProperties(data)
	result.initCorners(data)
	return &result
}

func (o *Object) initProperties(data *internal.LayerObject) {
	o.Properties = make(map[string]any)
	o.Properties[property.ObjectId] = data.Id
	o.Properties[property.ObjectClass] = data.Class
	o.Properties[property.ObjectTemplate] = data.Template
	o.Properties[property.ObjectName] = data.Name
	o.Properties[property.ObjectVisible] = data.Visible != "false"
	o.Properties[property.ObjectLocked] = data.Locked
	o.Properties[property.ObjectX] = data.X
	o.Properties[property.ObjectY] = data.Y
	o.Properties[property.ObjectWidth] = data.Width
	o.Properties[property.ObjectHeight] = data.Height
	o.Properties[property.ObjectRotation] = data.Rotation
	o.Properties[property.ObjectTileId] = data.Gid
	o.Properties[property.ObjectFlipX] = flag.IsOn(data.Gid, internal.FlipX)
	o.Properties[property.ObjectFlipY] = flag.IsOn(data.Gid, internal.FlipY)

	if data.Text != nil {
		o.Properties[property.ObjectText] = data.Text.Value
		o.Properties[property.ObjectTextFont] = data.Text.FontFamily
		o.Properties[property.ObjectTextFontSize] = data.Text.FontSize
		o.Properties[property.ObjectTextBold] = data.Text.Bold
		o.Properties[property.ObjectTextItalic] = data.Text.Italic
		o.Properties[property.ObjectTextStrikeout] = data.Text.Strikeout
		o.Properties[property.ObjectTextUnderline] = data.Text.Underline
		o.Properties[property.ObjectTextAlignX] = aligns[data.Text.AlignX]
		o.Properties[property.ObjectTextAlignY] = aligns[data.Text.AlignY]
		o.Properties[property.ObjectTextColor] = color.Hex(data.Text.Color)
		o.Properties[property.ObjectTextWordWrap] = data.Text.WordWrap
		o.Properties[property.ObjectType] = "text"
	} else if data.Gid > 0 {
		o.Properties[property.ObjectType] = "tile"
	} else if data.Ellipse != nil {
		o.Properties[property.ObjectType] = "ellipse"
	} else if data.Point != nil {
		o.Properties[property.ObjectType] = "point"
	} else if data.Polyline != nil {
		o.Properties[property.ObjectType] = "line"
	} else if data.Polygon != nil {
		o.Properties[property.ObjectType] = "polygon"
	} else {
		o.Properties[property.ObjectType] = "rectangle"
	}

	var owner *Project = nil
	if o.OwnerLayer != nil {
		owner = o.OwnerLayer.OwnerMap.Project
	} else if o.OwnerTile == nil {
		owner = o.OwnerTile.OwnerTileset.Project
	}

	for _, prop := range data.Properties {
		o.Properties[prop.Name] = parseProperty(prop, owner)
	}
}
func (o *Object) initCorners(data *internal.LayerObject) {
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

	o.Corners = corners
}

func (o *Object) getTile() *Tile {
	var tileId = o.Properties[property.ObjectTileId]
	var id = flag.TurnOff(tileId.(uint32), internal.Flips)
	var curTileset *Tileset = nil
	var firstId uint32 = 1
	if o.OwnerLayer != nil {
		curTileset, firstId = currentTileset(o.OwnerLayer.OwnerMap, id)
	} else if o.OwnerTile == nil {
		curTileset = o.OwnerTile.OwnerTileset
	}

	if id == 0 {
		return nil
	}
	return curTileset.Tiles[id-firstId]
}
func (o *Object) getOffsets() (worldX, worldY, layerX, layerY float32) {
	if o.OwnerLayer != nil {
		worldX, worldY, layerX, layerY = o.OwnerLayer.getOffsets()
	} else if o.OwnerTile == nil {
		layerX = o.OwnerTile.OwnerTileset.Properties[property.TilesetOffsetX].(float32)
		layerY = o.OwnerTile.OwnerTileset.Properties[property.TilesetOffsetY].(float32)
	}
	return
}
func (o *Object) getArea() (x, y, w, h, a float32) {
	x = o.Properties[property.ObjectX].(float32)
	y = o.Properties[property.ObjectY].(float32)
	w = o.Properties[property.ObjectWidth].(float32)
	h = o.Properties[property.ObjectHeight].(float32)
	a = o.Properties[property.ObjectRotation].(float32)
	return
}
