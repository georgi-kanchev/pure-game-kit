package tiled

import (
	"pure-game-kit/data/path"
	"pure-game-kit/debug"
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/internal"
	"pure-game-kit/tiled/property"
	"pure-game-kit/utility/collection"
)

type Map struct {
	Project              *Project
	Properties           map[string]any
	Tilesets             []*Tileset
	TilesetsFirstTileIds []uint32
	Layers               []*Layer
}

func NewMap(mapId string, project *Project) *Map {
	var data, _ = internal.TiledMaps[mapId]
	if data == nil {
		debug.LogError("Failed to create map: \"", mapId, "\"\nNo data is loaded with this map id.")
		return nil
	}

	var result = Map{Project: project, Layers: []*Layer{}}
	result.initProperties(data)
	result.initTilesets(data)
	result.initLayers(data.Directory, &data.Layers, nil)
	return &result
}

//=================================================================

func (Map *Map) TextBoxes() []*graphics.TextBox {
	var result = []*graphics.TextBox{}
	for _, layer := range Map.Layers {
		result = append(result, layer.TextBoxes()...)
	}
	return result
}
func (Map *Map) Sprites() []*graphics.Sprite {
	var result = []*graphics.Sprite{}
	for _, layer := range Map.Layers {
		result = append(result, layer.Sprites()...)
	}
	return result
}
func (Map *Map) Shapes() []*geometry.Shape {
	var result = []*geometry.Shape{}
	for _, layer := range Map.Layers {
		result = append(result, layer.Shapes()...)
	}
	return result
}
func (Map *Map) Lines() [][2]float32 {
	var result = [][2]float32{}
	for _, layer := range Map.Layers {
		result = append(result, layer.Lines()...)
	}
	return result
}
func (Map *Map) Points() [][2]float32 {
	var result = [][2]float32{}
	for _, layer := range Map.Layers {
		result = append(result, layer.Points()...)
	}
	return result
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

	collection.Reverse(Map.Layers)
}
