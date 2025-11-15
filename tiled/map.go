package tiled

import (
	"pure-game-kit/data/path"
	"pure-game-kit/debug"
	"pure-game-kit/internal"
	"pure-game-kit/tiled/property"
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
	result.initLayers(data.Directory, &data.Layers)
	return &result
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
func (Map *Map) initLayers(directory string, layers *internal.Layers) {
	for _, layer := range layers.LayersTiles {
		Map.Layers = append(Map.Layers, newLayerTiles(layer, Map))
	}
	for _, layer := range layers.LayersObjects {
		Map.Layers = append(Map.Layers, newLayerObjects(layer, Map))
	}
	for _, layer := range layers.LayersImages {
		Map.Layers = append(Map.Layers, newLayerImage(directory, layer, Map))
	}
	for _, group := range layers.LayersGroups {
		Map.Layers = append(Map.Layers, newLayerGroup(group, Map))
		Map.initLayers(directory, &group.Layers)
	}
}
