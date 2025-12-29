package tiled

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/data/path"
	"pure-game-kit/execution/condition"
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	it "pure-game-kit/internal"
	"pure-game-kit/tiled/property"
	"pure-game-kit/utility/collection"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/flag"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/point"
)

type Layer struct {
	Properties map[string]any
	TileIds    []uint32  // used by Tile Layers only
	Objects    []*Object // used by Object Layers only

	OwnerMap   *Map
	OwnerGroup *Layer
}

//=================================================================

func (l *Layer) FindObjectsBy(property string, value any) []*Object {
	var result = []*Object{}
	for _, obj := range l.Objects {
		var curValue, has = obj.Properties[property]
		if has && value == curValue {
			result = append(result, obj)
		}
	}
	return result
}

func (l *Layer) ExtractSprites() []*graphics.Sprite {
	var result = []*graphics.Sprite{}
	var image, hasImage = l.Properties[property.LayerImage]
	var tint = l.Properties[property.LayerTint].(uint)

	if hasImage {
		var worldX, worldY, layerX, layerY = l.getOffsets()
		var imgW = l.Properties[property.LayerImageWidth].(int)
		var imgH = l.Properties[property.LayerImageHeight].(int)
		var sprite = graphics.NewSprite(image.(string), worldX+layerX, worldY+layerY)
		sprite.Color = tint
		sprite.Width, sprite.Height = float32(imgW), float32(imgH)
		sprite.PivotX, sprite.PivotY = 0, 0
		return []*graphics.Sprite{sprite}
	}

	if len(l.Objects) > 0 {
		for _, obj := range l.Objects {
			var sprite = obj.ExtractSprite()
			if sprite != nil {
				sprite.Color = tint
				result = append(result, sprite)
			}
		}
		return result
	}

	l.forEachTile(func(tile *Tile, ang, x, y, w, h, scW, scH float32, cellX, cellY int) {
		var sprite = tile.ExtractSprite()
		sprite.X, sprite.Y = x, y
		sprite.Width, sprite.Height = w, h
		sprite.Angle = ang
		sprite.Color = tint
		result = append(result, sprite)
	})
	return result
}
func (l *Layer) ExtractTextBoxes() []*graphics.TextBox {
	var result = []*graphics.TextBox{} // tile & image layers don't have textboxes
	for _, obj := range l.Objects {
		var textBox = obj.ExtractTextBox()
		if textBox != nil {
			result = append(result, textBox)
		}
	}
	return result
}
func (l *Layer) ExtractShapeGrid() *geometry.ShapeGrid {
	var tileW = l.OwnerMap.Properties[property.MapTileWidth].(int)
	var tileH = l.OwnerMap.Properties[property.MapTileHeight].(int)
	var result = geometry.NewShapeGrid(tileW, tileH)
	var cellW = float32(l.OwnerMap.Properties[property.MapTileWidth].(int))
	var cellH = float32(l.OwnerMap.Properties[property.MapTileHeight].(int))

	l.forEachTile(func(tile *Tile, ang, x, y, w, h, scW, scH float32, cellX, cellY int) {
		var shapes = tile.ExtractShapes()
		if len(shapes) == 0 {
			return
		}
		var _, isImage = tile.Properties[property.TileImage]
		var cx, cy = float32(cellX) * cellW, float32(cellY) * cellH
		var offX, offY = x - cx, y - cy

		for _, shape := range shapes {
			if isImage {
				shape.ScaleX, shape.ScaleY = scW, scH
			}

			if w < 0 {
				shape.X *= -1
				shape.ScaleX *= -1
				shape.Angle = -shape.Angle
			}
			if h < 0 {
				shape.Y *= -1
				shape.ScaleY *= -1
				shape.Angle = -shape.Angle
			}

			shape.X, shape.Y = point.RotateAroundPoint(shape.X*scW, shape.Y*scH, 0, 0, ang)
			shape.X, shape.Y = shape.X+offX-cellW/2, shape.Y+offY-cellH/2
			shape.Angle += ang
		}
		result.SetAtCell(cellX, cellY, shapes...)
	})
	return result
}
func (l *Layer) ExtractShapes() []*geometry.Shape {
	var result = []*geometry.Shape{}
	for _, object := range l.Objects {
		result = append(result, object.ExtractShapes()...)
	}
	return result
}
func (l *Layer) ExtractLines() [][2]float32 {
	var result = [][2]float32{}
	for i, obj := range l.Objects {
		if i != 0 {
			result = append(result, [2]float32{number.NaN(), number.NaN()})
		}

		result = append(result, obj.ExtractLines()...)
	}

	l.forEachTile(func(tile *Tile, ang, x, y, w, h, scW, scH float32, cellX, cellY int) {
		var points = tile.ExtractLines()
		for _, pt := range points {
			pt[0], pt[1] = pt[0]*scW, pt[1]*scH
			pt[0], pt[1] = pt[0]*condition.If(w < 0, float32(-1), 1), pt[1]*condition.If(h < 0, float32(-1), 1)
			pt[0], pt[1] = point.RotateAroundPoint(pt[0], pt[1], 0, 0, ang)
			pt[0], pt[1] = pt[0]+x, pt[1]+y
			result = append(result, pt)
		}
	})
	return result
}
func (l *Layer) ExtractPoints() [][2]float32 {
	var result = [][2]float32{}
	for _, obj := range l.Objects {
		result = append(result, obj.ExtractPoints()...)
	}

	l.forEachTile(func(tile *Tile, ang, x, y, w, h, scW, scH float32, cellX, cellY int) {
		var points = tile.ExtractPoints()
		for _, pt := range points {
			pt[0], pt[1] = pt[0]*scW, pt[1]*scH
			pt[0], pt[1] = pt[0]*condition.If(w < 0, float32(-1), 1), pt[1]*condition.If(h < 0, float32(-1), 1)
			pt[0], pt[1] = point.RotateAroundPoint(pt[0], pt[1], 0, 0, ang)
			pt[0], pt[1] = pt[0]+x, pt[1]+y
			result = append(result, pt)
		}
	})
	return result
}

//=================================================================

func (layer *Layer) Draw(camera *graphics.Camera) {
	var l = layer
	var col, hasCol = l.Properties[property.LayerColor]
	if !hasCol {
		col = palette.White
	}

	draw(camera, l.ExtractSprites(), l.ExtractTextBoxes(),
		append(l.ExtractShapes(), l.ExtractShapeGrid().All()...), l.ExtractPoints(), l.ExtractLines(), col.(uint))
}

//=================================================================

func newLayerTiles(data *it.LayerTiles, owner *Map, group *Layer) *Layer {
	var layer = Layer{TileIds: collection.Clone(data.Tiles), OwnerMap: owner, OwnerGroup: group}
	layer.initProperties(&data.Layer, nil, nil, "")
	return &layer
}
func newLayerObjects(data *it.LayerObjects, owner *Map, group *Layer) *Layer {
	var layer = Layer{OwnerMap: owner, OwnerGroup: group}
	layer.initProperties(&data.Layer, data, nil, "")
	layer.initObjects(data)
	return &layer
}
func newLayerImage(directory string, data *it.LayerImage, owner *Map, group *Layer) *Layer {
	var layer = Layer{OwnerMap: owner, OwnerGroup: group}
	layer.initProperties(&data.Layer, nil, data, directory)
	return &layer
}
func newLayerGroup(data *it.LayerGroup, owner *Map, group *Layer) *Layer {
	var layer = Layer{OwnerMap: owner, OwnerGroup: group}
	layer.initProperties(&data.Layer, nil, nil, "")
	return &layer
}

//=================================================================

func (l *Layer) initProperties(data *it.Layer, objs *it.LayerObjects, img *it.LayerImage, dir string) {
	l.Properties = make(map[string]any)
	l.Properties[property.LayerId] = data.Id
	l.Properties[property.LayerClass] = data.Class
	l.Properties[property.LayerName] = data.Name
	l.Properties[property.LayerVisible] = data.Visible != "false"
	l.Properties[property.LayerLocked] = data.Locked
	l.Properties[property.LayerOpacity] = data.Opacity
	l.Properties[property.LayerTint] = condition.If(data.Tint == "", palette.White, color.Hex(data.Tint))
	l.Properties[property.LayerOffsetX] = data.OffsetX
	l.Properties[property.LayerOffsetY] = data.OffsetY
	l.Properties[property.LayerParallaxX] = data.ParallaxX
	l.Properties[property.LayerParallaxY] = data.ParallaxY

	if objs != nil {
		l.Properties[property.LayerColor] = color.Hex(objs.Color)
		l.Properties[property.LayerDrawOrder] = objs.DrawOrder
	}

	if img != nil && img.Image != nil {
		l.Properties[property.LayerImage] = assets.LoadTexture(path.New(dir, img.Image.Source))
		l.Properties[property.LayerImageWidth] = img.Image.Width
		l.Properties[property.LayerImageHeight] = img.Image.Height
		l.Properties[property.LayerTransparentColor] = color.Hex(img.Image.TransparentColor)
		l.Properties[property.LayerRepeatX] = img.RepeatX
		l.Properties[property.LayerRepeatY] = img.RepeatY
	}

	for _, prop := range data.Properties {
		l.Properties[prop.Name] = parseProperty(prop, l.OwnerMap.Project)
	}
}
func (l *Layer) initObjects(data *it.LayerObjects) {
	l.Objects = make([]*Object, len(data.Objects))
	for i, obj := range data.Objects {
		l.Objects[i] = newObject(obj, nil, l)
		l.Objects[i].Properties[property.ObjectOrder] = i
	}
}

func (l *Layer) getOffsets() (worldX, worldY, layerX, layerY float32) {
	worldX = l.OwnerMap.Properties[property.MapWorldX].(float32)
	worldY = l.OwnerMap.Properties[property.MapWorldY].(float32)
	layerX = l.Properties[property.LayerOffsetX].(float32)
	layerY = l.Properties[property.LayerOffsetY].(float32)
	return
}
func (l *Layer) forEachTile(action func(tile *Tile, ang, x, y, w, h, scW, scH float32, cellX, cellY int)) {
	if len(l.TileIds) == 0 {
		return
	}

	var columns = l.OwnerMap.Properties[property.MapColumns].(int)
	var worldX, worldY, layerX, layerY = l.getOffsets()
	var rows = l.OwnerMap.Properties[property.MapRows].(int)
	var cellW = float32(l.OwnerMap.Properties[property.MapTileWidth].(int))
	var cellH = float32(l.OwnerMap.Properties[property.MapTileHeight].(int))

	for i, tileId := range l.TileIds {
		var id = flag.TurnOff(tileId, it.Flips)
		if id == 0 {
			continue
		}

		var cx, cy = number.Index1DToIndexes2D(i, columns, rows)
		var cellX, cellY = float32(cx) * cellW, float32(cy) * cellH
		var curTileset, firstId = currentTileset(l.OwnerMap, id)
		var renderSize = curTileset.Properties[property.TilesetRenderSize].(string)
		var fillMode = curTileset.Properties[property.TilesetFillMode].(string)
		var tile = curTileset.Tiles[id-firstId]
		var tileW = float32(tile.Properties[property.TileWidth].(int))
		var tileH = float32(tile.Properties[property.TileHeight].(int))
		var _, isImage = tile.Properties[property.TileImage]
		var width, height, ratioW, ratioH = curTileset.tileRenderSize(tileW, tileH, cellW, cellH)
		var ang, w, h, offX, offY = tileOrientation(tileId, width, height, cellH, isImage)
		var scX, scY float32 = 1, 1

		cellX, cellY = cellX+worldX+layerX, cellY+worldY+layerY
		cx += int((worldX + layerX) / cellW)
		cy += int((worldY + layerY) / cellH)

		if renderSize == "grid" {
			scX, scY = cellW/tileW, cellH/tileH

			if fillMode == "preserve-aspect-fit" {
				scX, scY = scX*ratioW, scY*ratioH

				switch ang {
				case 0, 180:
					offY -= cellH/2 - number.Smallest(width, height)/2
				case 90, 270:
					offX += cellW/2 - number.Smallest(width, height)/2
				}
			}
		}

		action(tile, ang, cellX+offX, cellY+offY, w, h, scX, scY, cx, cy)
	}
}
