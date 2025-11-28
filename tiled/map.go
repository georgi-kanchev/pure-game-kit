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

func (Map *Map) Recreate() {
	var data, _ = internal.TiledMaps[Map.assetId]
	if data == nil {
		return
	}

	Map.Layers = []*Layer{}
	Map.initProperties(data)
	Map.initTilesets(data)
	Map.initLayers(data.Directory, &data.Layers, nil)
	Map.sortLayers(data)
}

//=================================================================

func (Map *Map) FindLayersBy(property string, value any) []*Layer {
	var result = []*Layer{}
	for _, layer := range Map.Layers {
		var curValue, has = layer.Properties[property]
		if has && value == curValue {
			result = append(result, layer)
		}
	}
	return result
}

func (Map *Map) ExtractSprites() []*graphics.Sprite {
	var result = []*graphics.Sprite{}
	for _, layer := range Map.Layers {
		result = append(result, layer.ExtractSprites()...)
	}
	return result
}
func (Map *Map) ExtractTextBoxes() []*graphics.TextBox {
	var result = []*graphics.TextBox{}
	for _, layer := range Map.Layers {
		result = append(result, layer.ExtractTextBoxes()...)
	}
	return result
}
func (Map *Map) ExtractShapes() []*geometry.Shape {
	var result = []*geometry.Shape{}
	for _, layer := range Map.Layers {
		result = append(result, layer.ExtractShapes()...)
	}
	return result
}
func (Map *Map) ExtractLines() [][2]float32 {
	var result = [][2]float32{}
	for _, layer := range Map.Layers {
		result = append(result, layer.ExtractLines()...)
	}
	return result
}
func (Map *Map) ExtractPoints() [][2]float32 {
	var result = [][2]float32{}
	for _, layer := range Map.Layers {
		result = append(result, layer.ExtractPoints()...)
	}
	return result
}

//=================================================================

func (Map *Map) Draw(camera *graphics.Camera) {
	for _, layer := range Map.Layers {
		layer.Draw(camera)
	}
}

//=================================================================
// private

func (Map *Map) initProperties(data *internal.Map) {
	Map.Properties = make(map[string]any)
	Map.Properties[property.MapName] = data.Name
	Map.Properties[property.MapClass] = data.Class
	Map.Properties[property.MapTileWidth] = data.TileWidth
	Map.Properties[property.MapTileHeight] = data.TileHeight
	Map.Properties[property.MapColumns] = data.Width
	Map.Properties[property.MapRows] = data.Height
	Map.Properties[property.MapInfinite] = data.Infinite
	Map.Properties[property.MapParallaxX] = data.ParallaxOriginX
	Map.Properties[property.MapParallaxY] = data.ParallaxOriginY
	Map.Properties[property.MapBackgroundColor] = data.BackgroundColor
	Map.Properties[property.MapWorldX] = data.WorldX
	Map.Properties[property.MapWorldY] = data.WorldY

	for _, prop := range data.Properties {
		Map.Properties[prop.Name] = parseProperty(prop, Map.Project)
	}
}
func (Map *Map) initTilesets(data *internal.Map) {
	Map.Tilesets = make([]*Tileset, len(data.Tilesets))
	Map.TilesetsFirstTileIds = make([]uint32, len(data.Tilesets))

	for i, t := range data.Tilesets {
		Map.Tilesets[i] = newTileset(path.New(data.Directory, t.Source), Map.Project)
		Map.TilesetsFirstTileIds[i] = data.FirstTileIds[i]
	}
}
func (Map *Map) initLayers(directory string, layers *internal.Layers, ownerGroup *Layer) {
	for _, group := range layers.LayersGroups {
		var layer = newLayerGroup(group, Map, ownerGroup)
		Map.Layers = append(Map.Layers, layer)
		Map.initLayers(directory, &group.Layers, layer)
	}
	for _, layer := range layers.LayersImages {
		Map.Layers = append(Map.Layers, newLayerImage(directory, layer, Map, ownerGroup))
	}
	for _, layer := range layers.LayersObjects {
		Map.Layers = append(Map.Layers, newLayerObjects(layer, Map, ownerGroup))
	}
	for _, layer := range layers.LayersTiles {
		Map.Layers = append(Map.Layers, newLayerTiles(layer, Map, ownerGroup))
	}
}

func (Map *Map) sortLayers(data *internal.Map) {
	var sortedLayers = []*Layer{}
	for _, id := range data.LayersInOrder {
		for _, layer := range Map.Layers {
			var curId = layer.Properties[property.LayerId]
			if curId == id {
				sortedLayers = append(sortedLayers, layer)
				layer.Properties[property.LayerOrder] = len(sortedLayers) - 1
			}
		}
	}
	Map.Layers = sortedLayers
}
