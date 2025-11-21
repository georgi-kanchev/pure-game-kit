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

func (layer *Layer) Sprites() []*graphics.Sprite {
	var result = []*graphics.Sprite{}
	var image, hasImage = layer.Properties[property.LayerImage]

	if hasImage {
		var worldX, worldY, layerX, layerY = layer.getOffsets()
		var imgW = layer.Properties[property.LayerImageWidth].(int)
		var imgH = layer.Properties[property.LayerImageHeight].(int)
		var sprite = graphics.NewSprite(image.(string), worldX+layerX, worldY+layerY)
		sprite.Width, sprite.Height = float32(imgW), float32(imgH)
		sprite.PivotX, sprite.PivotY = 0, 0
		return []*graphics.Sprite{sprite}
	}

	if len(layer.Objects) > 0 {
		for _, obj := range layer.Objects {
			var sprite = obj.Sprite()
			if sprite != nil {
				result = append(result, sprite)
			}
		}
		return result
	}

	layer.forEachTile(func(tile *Tile, ang, x, y, w, h, scW, scH float32) {
		var sprite = tile.Sprite()
		sprite.X, sprite.Y = x, y
		sprite.Width, sprite.Height = w, h
		sprite.Angle = ang
		result = append(result, sprite)
	})
	return result
}
func (layer *Layer) TextBoxes() []*graphics.TextBox {
	var result = []*graphics.TextBox{} // tile & image layers don't have textboxes
	for _, obj := range layer.Objects {
		var textBox = obj.TextBox()
		if textBox != nil {
			result = append(result, textBox)
		}
	}
	return result
}
func (layer *Layer) Shapes() []*geometry.Shape {
	var result = []*geometry.Shape{}

	for _, object := range layer.Objects {
		result = append(result, object.Shapes()...)
	}

	layer.forEachTile(func(tile *Tile, ang, x, y, w, h, scW, scH float32) {
		var shapes = tile.Shapes()
		var _, isImage = tile.Properties[property.TileImage]

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
			shape.X, shape.Y = shape.X+x, shape.Y+y
			shape.Angle += ang

			result = append(result, shape)
		}
	})
	return result
}
func (layer *Layer) Lines() [][2]float32 {
	var result = [][2]float32{}
	for i, obj := range layer.Objects {
		if i != 0 {
			result = append(result, [2]float32{number.NaN(), number.NaN()})
		}

		result = append(result, obj.Lines()...)
	}

	layer.forEachTile(func(tile *Tile, ang, x, y, w, h, scW, scH float32) {
		var points = tile.Lines()
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
func (layer *Layer) Points() [][2]float32 {
	var result = [][2]float32{}
	for _, obj := range layer.Objects {
		result = append(result, obj.Points()...)
	}

	layer.forEachTile(func(tile *Tile, ang, x, y, w, h, scW, scH float32) {
		var points = tile.Points()
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

func (layer *Layer) initProperties(data *it.Layer, objs *it.LayerObjects, img *it.LayerImage, dir string) {
	layer.Properties = make(map[string]any)
	layer.Properties[property.LayerName] = data.Id
	layer.Properties[property.LayerClass] = data.Class
	layer.Properties[property.LayerName] = data.Name
	layer.Properties[property.LayerVisible] = data.Visible != "false"
	layer.Properties[property.LayerLocked] = data.Locked
	layer.Properties[property.LayerOpacity] = data.Opacity
	layer.Properties[property.LayerTint] = color.Hex(data.Tint)
	layer.Properties[property.LayerOffsetX] = data.OffsetX
	layer.Properties[property.LayerOffsetY] = data.OffsetY
	layer.Properties[property.LayerParallaxX] = data.ParallaxX
	layer.Properties[property.LayerParallaxY] = data.ParallaxY

	if objs != nil {
		layer.Properties[property.LayerColor] = color.Hex(objs.Color)
		layer.Properties[property.LayerDrawOrder] = objs.DrawOrder
	}

	if img != nil && img.Image != nil {
		layer.Properties[property.LayerImage] = assets.LoadTexture(path.New(dir, img.Image.Source))
		layer.Properties[property.LayerImageWidth] = img.Image.Width
		layer.Properties[property.LayerImageHeight] = img.Image.Height
		layer.Properties[property.LayerTransparentColor] = color.Hex(img.Image.TransparentColor)
		layer.Properties[property.LayerRepeatX] = img.RepeatX
		layer.Properties[property.LayerRepeatY] = img.RepeatY
	}

	for _, prop := range data.Properties {
		layer.Properties[prop.Name] = parseProperty(prop, layer.OwnerMap.Project)
	}
}
func (layer *Layer) initObjects(data *it.LayerObjects) {
	layer.Objects = make([]*Object, len(data.Objects))
	for i, obj := range data.Objects {
		layer.Objects[i] = newObject(obj, nil, layer)
	}
	collection.Reverse(layer.Objects)
}

func (layer *Layer) getOffsets() (worldX, worldY, layerX, layerY float32) {
	worldX = layer.OwnerMap.Properties[property.MapWorldX].(float32)
	worldY = layer.OwnerMap.Properties[property.MapWorldY].(float32)
	layerX = layer.Properties[property.LayerOffsetX].(float32)
	layerY = layer.Properties[property.LayerOffsetY].(float32)
	return
}
func (layer *Layer) forEachTile(action func(tile *Tile, ang, x, y, w, h, scW, scH float32)) {
	if len(layer.TileIds) == 0 {
		return
	}

	var columns = layer.OwnerMap.Properties[property.MapColumns].(int)
	var worldX, worldY, layerX, layerY = layer.getOffsets()
	var rows = layer.OwnerMap.Properties[property.MapRows].(int)
	var cellW = float32(layer.OwnerMap.Properties[property.MapTileWidth].(int))
	var cellH = float32(layer.OwnerMap.Properties[property.MapTileHeight].(int))

	for i, tileId := range layer.TileIds {
		var id = flag.TurnOff(tileId, it.Flips)
		if id == 0 {
			continue
		}

		var cx, cy = number.Index1DToIndexes2D(i, columns, rows)
		var cellX, cellY = float32(cx) * cellW, float32(cy) * cellH
		var curTileset, firstId = currentTileset(layer.OwnerMap, id)
		var renderSize = curTileset.Properties[property.TilesetRenderSize].(string)
		var fillMode = curTileset.Properties[property.TilesetFillMode].(string)
		var tile = curTileset.Tiles[id-firstId]
		var tileW = float32(tile.Properties[property.TileWidth].(int))
		var tileH = float32(tile.Properties[property.TileHeight].(int))
		var _, isImage = tile.Properties[property.TileImage]
		var width, height, ratioW, ratioH = curTileset.tileRenderSize(tileW, tileH, cellW, cellH)
		var ang, w, h, offX, offY = tileOrientation(tileId, width, height, cellH, isImage)
		var scX, scY float32 = 1, 1

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

		action(tile, ang, worldX+layerX+cellX+offX, worldY+layerY+cellY+offY, w, h, scX, scY)
	}
}
