package tiled

import (
	"pure-game-kit/debug"
	"pure-game-kit/internal"
	"pure-game-kit/tiled/property"
	"pure-game-kit/utility/text"
)

type Tile struct {
	Project    *Project
	Properties map[string]any
	Objects    []*Object

	IsAnimating bool
}

func NewTile(tilesetId string, tileId uint32, project *Project) *Tile {
	var data, _ = internal.TiledTilesets[tilesetId]
	if data == nil {
		debug.LogError("Failed to create tile: \"", tilesetId, "/", tileId, "\"\n",
			"No data is loaded with this tileset id.")
		return nil
	}

	var tileData = data.MappedTiles[tileId]
	if tileData == nil {
		debug.LogError("Failed to create tile: \"", tilesetId, "/", tileId, "\"\n",
			"The tileset contains no data with this tile id.")
		return nil
	}

	var result = Tile{}
	result.initProperties(data, tileData, project)
	return &result
}

//=================================================================
// private

func (t *Tile) initProperties(tilesetData *internal.Tileset, data *internal.TilesetTile, project *Project) {
	var w, h = tilesetData.TileWidth, tilesetData.TileHeight
	if data.Image != nil {
		w, h = data.Image.Width, data.Image.Height
	}

	t.Properties = make(map[string]any)
	t.Properties[property.TileId] = data.Id
	t.Properties[property.TileClass] = data.Class
	t.Properties[property.TileProbability] = text.ToNumber[float32](defaultText(data.Probability, "1"))
	t.Properties[property.TileWidth] = w
	t.Properties[property.TileHeight] = h

	for _, prop := range data.Properties {
		t.Properties[prop.Name] = parseProperty(prop, project)
	}
}
