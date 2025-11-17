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
)

type Layer struct {
	Properties map[string]any
	TileIds    []uint32  // used by Tile Layers only
	Objects    []*Object // used by Object Layers only

	OwnerMap *Map
}

//=================================================================

func (layer *Layer) Sprites() []*graphics.Sprite {
	var result = []*graphics.Sprite{}
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
	for _, obj := range layer.Objects {
		result = append(result, obj.Shape())
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

func (t *Layer) initProperties(
	data *internal.Layer, objs *internal.LayerObjects, img *internal.LayerImage, dir string) {
	t.Properties = make(map[string]any)
	t.Properties[property.LayerName] = data.Id
	t.Properties[property.LayerClass] = data.Class
	t.Properties[property.LayerName] = data.Name
	t.Properties[property.LayerVisible] = data.Visible != "false"
	t.Properties[property.LayerLocked] = data.Locked
	t.Properties[property.LayerOpacity] = data.Opacity
	t.Properties[property.LayerTint] = color.Hex(data.Tint)
	t.Properties[property.LayerOffsetX] = data.OffsetX
	t.Properties[property.LayerOffsetY] = data.OffsetY
	t.Properties[property.LayerParallaxX] = data.ParallaxX
	t.Properties[property.LayerParallaxY] = data.ParallaxY

	if objs != nil {
		t.Properties[property.LayerColor] = color.Hex(objs.Color)
		t.Properties[property.LayerDrawOrder] = objs.DrawOrder
	}

	if img != nil && img.Image != nil {
		t.Properties[property.LayerImage] = assets.LoadTexture(path.New(dir, img.Image.Source))
		t.Properties[property.LayerTransparentColor] = color.Hex(img.Image.TransparentColor)
		t.Properties[property.LayerRepeatX] = img.RepeatX
		t.Properties[property.LayerRepeatY] = img.RepeatY
	}

	for _, prop := range data.Properties {
		t.Properties[prop.Name] = parseProperty(prop, t.OwnerMap.Project)
	}
}
func (t *Layer) initObjects(data *internal.LayerObjects) {
	t.Objects = make([]*Object, len(data.Objects))
	for i, obj := range data.Objects {
		t.Objects[i] = newObject(obj, nil, t)
	}
}
