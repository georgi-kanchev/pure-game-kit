package tiled

import (
	"pure-game-kit/data/path"
	"pure-game-kit/debug"
	"pure-game-kit/internal"
	"pure-game-kit/tiled/property"
)

type Map struct {
	Project    *Project
	Properties map[string]any
	Tilesets   map[*Tileset]uint32
	Layers     []*Layer
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
	m.Tilesets = make(map[*Tileset]uint32, len(data.Tilesets))

	for i, t := range data.Tilesets {
		var tileset = newTileset(path.New(data.Directory, t.Source), m.Project)
		m.Tilesets[tileset] = data.FirstTileIds[i]
	}
}
func (m *Map) initLayers(directory string, layers *internal.Layers) {
	for _, layer := range layers.LayersTiles {
		m.Layers = append(m.Layers, newLayerTiles(layer, m))
	}
	for _, layer := range layers.LayersObjects {
		m.Layers = append(m.Layers, newLayerObjects(layer, m))
	}
	for _, layer := range layers.LayersImages {
		m.Layers = append(m.Layers, newLayerImage(directory, layer, m))
	}
	for _, group := range layers.LayersGroups {
		m.Layers = append(m.Layers, newLayerGroup(group, m))
		m.initLayers(directory, &group.Layers)
	}
}
