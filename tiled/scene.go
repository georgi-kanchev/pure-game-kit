package tiled

import (
	"pure-game-kit/data/path"
	"pure-game-kit/debug"
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/internal"
	"pure-game-kit/tiled/property"
)

type Scene struct {
	Project              *Project
	Properties           map[string]any
	Tilesets             []*Tileset
	TilesetsFirstTileIds []uint32
	Layers               []*Layer

	mapId string
}

func NewScene(mapId string, project *Project) *Scene {
	var data, _ = internal.TiledMaps[mapId]
	if data == nil {
		debug.LogError("Failed to create map: \"", mapId, "\"\nNo data is loaded with this map id.")
		return nil
	}

	var result = &Scene{Project: project, mapId: mapId}
	result.Recreate()
	return result
}

func (s *Scene) Recreate() {
	var data, _ = internal.TiledMaps[s.mapId]
	if data == nil {
		return
	}

	s.Layers = []*Layer{}
	s.initProperties(data)
	s.initTilesets(data)
	s.initLayers(data.Directory, &data.Layers, nil)
	s.sortLayers(data)
}

//=================================================================

func (s *Scene) FindLayersBy(property string, value any) []*Layer {
	var result = []*Layer{}
	for _, layer := range s.Layers {
		var curValue, has = layer.Properties[property]
		if has && value == curValue {
			result = append(result, layer)
		}
	}
	return result
}

func (s *Scene) ExtractTileMaps() []*graphics.TileMap {
	var result = []*graphics.TileMap{}
	for _, layer := range s.Layers {
		result = append(result, layer.ExtractTileMaps()...)
	}
	return result
}
func (s *Scene) ExtractSprites() []*graphics.Sprite {
	var result = []*graphics.Sprite{}
	for _, layer := range s.Layers {
		result = append(result, layer.ExtractSprites()...)
	}
	return result
}
func (s *Scene) ExtractTextBoxes() []*graphics.TextBox {
	var result = []*graphics.TextBox{}
	for _, layer := range s.Layers {
		result = append(result, layer.ExtractTextBoxes()...)
	}
	return result
}
func (s *Scene) ExtractShapeGrids() []*geometry.ShapeGrid {
	var result = []*geometry.ShapeGrid{}
	for _, layer := range s.Layers {
		result = append(result, layer.ExtractShapeGrid())
	}
	return result
}
func (s *Scene) ExtractShapes() []*geometry.Shape {
	var result = []*geometry.Shape{}
	for _, layer := range s.Layers {
		result = append(result, layer.ExtractShapes()...)
	}
	return result
}
func (s *Scene) ExtractLines() [][2]float32 {
	var result = [][2]float32{}
	for _, layer := range s.Layers {
		result = append(result, layer.ExtractLines()...)
	}
	return result
}
func (s *Scene) ExtractPoints() [][2]float32 {
	var result = [][2]float32{}
	for _, layer := range s.Layers {
		result = append(result, layer.ExtractPoints()...)
	}
	return result
}

//=================================================================

func (s *Scene) Draw(camera *graphics.Camera) {
	for _, layer := range s.Layers {
		layer.Draw(camera)
	}
}

//=================================================================
// private

func (s *Scene) initProperties(data *internal.Map) {
	s.Properties = make(map[string]any)
	s.Properties[property.MapName] = data.Name
	s.Properties[property.MapClass] = data.Class
	s.Properties[property.MapTileWidth] = data.TileWidth
	s.Properties[property.MapTileHeight] = data.TileHeight
	s.Properties[property.MapColumns] = data.Width
	s.Properties[property.MapRows] = data.Height
	s.Properties[property.MapInfinite] = data.Infinite
	s.Properties[property.MapParallaxX] = data.ParallaxOriginX
	s.Properties[property.MapParallaxY] = data.ParallaxOriginY
	s.Properties[property.MapBackgroundColor] = data.BackgroundColor
	s.Properties[property.MapWorldX] = data.WorldX
	s.Properties[property.MapWorldY] = data.WorldY

	for _, prop := range data.Properties {
		s.Properties[prop.Name] = parseProperty(prop, s.Project)
	}
}
func (s *Scene) initTilesets(data *internal.Map) {
	s.Tilesets = make([]*Tileset, len(data.Tilesets))
	s.TilesetsFirstTileIds = make([]uint32, len(data.Tilesets))

	for i, t := range data.Tilesets {
		s.Tilesets[i] = newTileset(path.New(data.Directory, t.Source), s.Project)
		s.TilesetsFirstTileIds[i] = data.FirstTileIds[i]
	}
}
func (s *Scene) initLayers(directory string, layers *internal.Layers, ownerGroup *Layer) {
	for _, group := range layers.LayersGroups {
		var layer = newLayerGroup(group, s, ownerGroup)
		s.Layers = append(s.Layers, layer)
		s.initLayers(directory, &group.Layers, layer)
	}
	for _, layer := range layers.LayersImages {
		s.Layers = append(s.Layers, newLayerImage(directory, layer, s, ownerGroup))
	}
	for _, layer := range layers.LayersObjects {
		s.Layers = append(s.Layers, newLayerObjects(layer, s, ownerGroup))
	}
	for _, layer := range layers.LayersTiles {
		s.Layers = append(s.Layers, newLayerTiles(layer, s, ownerGroup))
	}
}

func (s *Scene) sortLayers(data *internal.Map) {
	var sortedLayers = []*Layer{}
	for _, id := range data.LayersInOrder {
		for _, layer := range s.Layers {
			var curId = layer.Properties[property.LayerId]
			if curId == id {
				sortedLayers = append(sortedLayers, layer)
				layer.Properties[property.LayerOrder] = len(sortedLayers) - 1
			}
		}
	}
	s.Layers = sortedLayers
}
