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
)

type Layer struct {
	Properties map[string]any
	TileIds    []uint32  // used by Tile Layers only
	Objects    []*Object // used by Object Layers only

	OwnerMap   *Map
	OwnerGroup *Layer
}

//=================================================================

func (layer *Layer) Lines() [][2]float32 {
	var result = [][2]float32{}
	for i, obj := range layer.Objects {
		if i != 0 {
			result = append(result, [2]float32{number.NaN(), number.NaN()})
		}

		result = append(result, obj.Lines()...)
	}
	return result
}
func (layer *Layer) Points() [][2]float32 {
	var result = [][2]float32{}
	for _, obj := range layer.Objects {
		result = append(result, obj.Points()...)
	}
	return result
}
func (layer *Layer) TextBoxes() []*graphics.TextBox {
	var result = []*graphics.TextBox{}
	for _, obj := range layer.Objects {
		var textBox = obj.TextBox()
		if textBox != nil {
			result = append(result, textBox)
		}
	}
	return result
}
func (layer *Layer) Sprites() []*graphics.Sprite {
	var result = []*graphics.Sprite{}
	var image, hasImage = layer.Properties[property.LayerImage]
	var columns = layer.OwnerMap.Properties[property.MapColumns].(int)
	var rows = layer.OwnerMap.Properties[property.MapRows].(int)
	var worldX, worldY, layerX, layerY = layer.getOffsets()

	if hasImage {
		var imgW = layer.Properties[property.LayerImageWidth].(int)
		var imgH = layer.Properties[property.LayerImageHeight].(int)
		var sprite = graphics.NewSprite(image.(string), worldX+layerX, worldY+layerY)
		sprite.Width, sprite.Height = float32(imgW), float32(imgH)
		sprite.PivotX, sprite.PivotY = 0, 0
		return []*graphics.Sprite{sprite}
	}

	for _, obj := range layer.Objects {
		var sprite = obj.Sprite()
		if sprite != nil {
			result = append(result, sprite)
		}
	}

	for i, tileId := range layer.TileIds {
		var id = flag.TurnOff(tileId, it.Flips)
		if id == 0 {
			continue
		}

		var cellX, cellY = number.Index1DToIndexes2D(i, columns, rows)
		var curTileset, firstId = currentTileset(layer.OwnerMap, id)
		var tile = curTileset.Tiles[id-firstId]
		var sprite = tile.Sprite()
		var tileW = float32(layer.OwnerMap.Properties[property.MapTileWidth].(int))
		var tileH = float32(layer.OwnerMap.Properties[property.MapTileHeight].(int))
		var _, isImage = tile.Properties[property.TileImage]
		var width, height = curTileset.tileRenderSize(sprite.Width, sprite.Height, tileW, tileH)
		var ang, w, h, offX, offY = tileOrientation(tileId, width, height, tileH, isImage)
		var tileRenderSize = curTileset.Properties[property.TilesetRenderSize].(string)

		if tileRenderSize == "grid" {
			sprite.PivotX, sprite.PivotY = 0.5, 0.5
			offX, offY = tileW/2, tileH/2
		}

		sprite.X, sprite.Y = worldX+layerX+float32(cellX)*tileW+offX, worldY+layerY+float32(cellY)*tileH+offY
		sprite.Width, sprite.Height = w, h
		sprite.Angle = ang

		result = append(result, sprite)
	}

	return result
}
func (layer *Layer) Shapes() []*geometry.Shape {
	var result = []*geometry.Shape{}
	for _, object := range layer.Objects {
		result = append(result, object.Shapes()...)
	}
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
func tileOrientation(tileId uint32, w, h, th float32, image bool) (ang, newW, newH, offX, offY float32) {
	var flipH = flag.IsOn(tileId, it.FlipX)
	var flipV = flag.IsOn(tileId, it.FlipY)
	var flipDiag = flag.IsOn(tileId, it.FlipDiag)

	ang = 0.0
	newW, newH = w, h
	offX, offY = 0, condition.If(image, th-h, 0)

	if flipH && !flipV && flipDiag { // rotation 90
		ang = 90
		offX = h
		offY = condition.If(image, th-w, 0)
	} else if flipH && flipV && !flipDiag { // rotation 180
		ang = 180
		offX = w
		offY = condition.If(image, th, h)
	} else if !flipH && flipV && flipDiag { // rotation 270
		ang = 270
		offY = condition.If(image, th, w)
	} else if flipH && !flipV && !flipDiag { // flip x only
		newW = -w
		offX = w
	} else if flipH && flipV && flipDiag { // flip x + rotation 90
		ang = 90
		newW = -w
		offX = h
		offY = condition.If(image, th, w)
	} else if !flipH && flipV && !flipDiag { // flip x + rotation 180
		newH = -h
		offY = condition.If(image, th, h)
	} else if !flipH && !flipV && flipDiag { // flip x + rotation 270
		ang = 270
		newW = -w
		offY = condition.If(image, th-w, 0)
	}
	return ang, newW, newH, offX, offY
}
