package tiled

import (
	"pure-game-kit/data/path"
	"pure-game-kit/debug"
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/internal"
	"pure-game-kit/tiled/property"
)

type Map struct {
	Project              *Project
	Properties           map[string]any
	Tilesets             []*Tileset
	TilesetsFirstTileIds []uint32
	Layers               []*Layer

	assetId string
}

func NewMap(mapId string, project *Project) *Map {
	var data, _ = internal.TiledMaps[mapId]
	if data == nil {
		debug.LogError("Failed to create map: \"", mapId, "\"\nNo data is loaded with this map id.")
		return nil
	}

	var result = &Map{Project: project, assetId: mapId}
	result.Recreate()
	return result
}

func (m *Map) Recreate() {
	var data, _ = internal.TiledMaps[m.assetId]
	if data == nil {
		return
	}

	m.Layers = []*Layer{}
	m.initProperties(data)
	m.initTilesets(data)
	m.initLayers(data.Directory, &data.Layers, nil)
	m.sortLayers(data)
}

//=================================================================

func (m *Map) FindLayersBy(property string, value any) []*Layer {
	var result = []*Layer{}
	for _, layer := range m.Layers {
		var curValue, has = layer.Properties[property]
		if has && value == curValue {
			result = append(result, layer)
		}
	}
	return result
}

func (m *Map) ExtractSprites() []*graphics.Sprite {
	var result = []*graphics.Sprite{}
	for _, layer := range m.Layers {
		result = append(result, layer.ExtractSprites()...)
	}
	return result
}
func (m *Map) ExtractTextBoxes() []*graphics.TextBox {
	var result = []*graphics.TextBox{}
	for _, layer := range m.Layers {
		result = append(result, layer.ExtractTextBoxes()...)
	}
	return result
}
func (m *Map) ExtractShapeGrids() []*geometry.ShapeGrid {
	var result = []*geometry.ShapeGrid{}
	for _, layer := range m.Layers {
		result = append(result, layer.ExtractShapeGrid())
	}
	return result
}
func (m *Map) ExtractShapes() []*geometry.Shape {
	var result = []*geometry.Shape{}
	for _, layer := range m.Layers {
		result = append(result, layer.ExtractShapes()...)
	}
	return result
}
func (m *Map) ExtractLines() [][2]float32 {
	var result = [][2]float32{}
	for _, layer := range m.Layers {
		result = append(result, layer.ExtractLines()...)
	}
	return result
}
func (m *Map) ExtractPoints() [][2]float32 {
	var result = [][2]float32{}
	for _, layer := range m.Layers {
		result = append(result, layer.ExtractPoints()...)
	}
	return result
}

//=================================================================

func (m *Map) Draw(camera *graphics.Camera) {
	for _, layer := range m.Layers {
		layer.Draw(camera)
	}
}

//=================================================================
// private

func (m *Map) initProperties(data *internal.Map) {
	m.Properties = make(map[string]any)
	m.Properties[property.MapName] = data.Name
	m.Properties[property.MapClass] = data.Class
	m.Properties[property.MapTileWidth] = data.TileWidth
	m.Properties[property.MapTileHeight] = data.TileHeight
	m.Properties[property.MapColumns] = data.Width
	m.Properties[property.MapRows] = data.Height
	m.Properties[property.MapInfinite] = data.Infinite
	m.Properties[property.MapParallaxX] = data.ParallaxOriginX
	m.Properties[property.MapParallaxY] = data.ParallaxOriginY
	m.Properties[property.MapBackgroundColor] = data.BackgroundColor
	m.Properties[property.MapWorldX] = data.WorldX
	m.Properties[property.MapWorldY] = data.WorldY

	for _, prop := range data.Properties {
		m.Properties[prop.Name] = parseProperty(prop, m.Project)
	}
}
func (m *Map) initTilesets(data *internal.Map) {
	m.Tilesets = make([]*Tileset, len(data.Tilesets))
	m.TilesetsFirstTileIds = make([]uint32, len(data.Tilesets))

	for i, t := range data.Tilesets {
		m.Tilesets[i] = newTileset(path.New(data.Directory, t.Source), m.Project)
		m.TilesetsFirstTileIds[i] = data.FirstTileIds[i]
	}
}
func (m *Map) initLayers(directory string, layers *internal.Layers, ownerGroup *Layer) {
	for _, group := range layers.LayersGroups {
		var layer = newLayerGroup(group, m, ownerGroup)
		m.Layers = append(m.Layers, layer)
		m.initLayers(directory, &group.Layers, layer)
	}
	for _, layer := range layers.LayersImages {
		m.Layers = append(m.Layers, newLayerImage(directory, layer, m, ownerGroup))
	}
	for _, layer := range layers.LayersObjects {
		m.Layers = append(m.Layers, newLayerObjects(layer, m, ownerGroup))
	}
	for _, layer := range layers.LayersTiles {
		m.Layers = append(m.Layers, newLayerTiles(layer, m, ownerGroup))
	}
}

func (m *Map) sortLayers(data *internal.Map) {
	var sortedLayers = []*Layer{}
	for _, id := range data.LayersInOrder {
		for _, layer := range m.Layers {
			var curId = layer.Properties[property.LayerId]
			if curId == id {
				sortedLayers = append(sortedLayers, layer)
				layer.Properties[property.LayerOrder] = len(sortedLayers) - 1
			}
		}
	}
	m.Layers = sortedLayers
}
