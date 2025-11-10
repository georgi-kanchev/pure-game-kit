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
	Tilesets   []*Tileset
	Layers     []*Layer
}

func NewMap(mapId string, project *Project) *Map {
	var data, _ = internal.TiledMaps[mapId]
	if data == nil {
		debug.LogError("Failed to create map: \"", mapId, "\"\nNo data is loaded with this map id.")
		return nil
	}

	var result = Map{}
	result.initProperties(data, project)
	result.initTilesets(data, project)
	return &result
}

//=================================================================
// private

func (m *Map) initProperties(data *internal.Map, project *Project) {
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
		m.Properties[prop.Name] = parseProperty(prop, project)
	}
}
func (m *Map) initTilesets(data *internal.Map, project *Project) {
	m.Tilesets = make([]*Tileset, len(data.Tilesets))

	for i, t := range data.Tilesets {
		m.Tilesets[i] = NewTileset(path.New(data.Directory, t.Source), project)
	}
}
