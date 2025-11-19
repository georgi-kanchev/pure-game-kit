package tiled

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/data/path"
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/internal"
	"pure-game-kit/tiled/property"
	"pure-game-kit/utility/collection"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/number"
)

type Layer struct {
	Properties map[string]any
	TileIds    []uint32  // used by Tile Layers only
	Objects    []*Object // used by Object Layers only

	OwnerMap *Map
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
	var image = layer.Properties[property.LayerImage].(string)
	if image != "" {
		var worldX, worldY, layerX, layerY = layer.getOffsets()
		var imgW = layer.Properties[property.LayerImageWidth].(int)
		var imgH = layer.Properties[property.LayerImageHeight].(int)
		var sprite = graphics.NewSprite(image, worldX+layerX, worldY+layerY)
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

func newLayerTiles(data *internal.LayerTiles, owner *Map) *Layer {
	var layer = Layer{TileIds: collection.Clone(data.Tiles), OwnerMap: owner}
	layer.initProperties(&data.Layer, nil, nil, "")
	return &layer
}
func newLayerObjects(data *internal.LayerObjects, owner *Map) *Layer {
	var layer = Layer{OwnerMap: owner}
	layer.initProperties(&data.Layer, data, nil, "")
	layer.initObjects(data)
	return &layer
}
func newLayerImage(directory string, data *internal.LayerImage, owner *Map) *Layer {
	var layer = Layer{OwnerMap: owner}
	layer.initProperties(&data.Layer, nil, data, directory)
	return &layer
}
func newLayerGroup(data *internal.LayerGroup, owner *Map) *Layer {
	var layer = Layer{OwnerMap: owner}
	layer.initProperties(&data.Layer, nil, nil, "")
	return &layer
}

//=================================================================

func (layer *Layer) initProperties(
	data *internal.Layer, objs *internal.LayerObjects, img *internal.LayerImage, dir string) {
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
func (layer *Layer) initObjects(data *internal.LayerObjects) {
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
