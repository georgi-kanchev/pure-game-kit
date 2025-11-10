package tiled

import (
	"pure-game-kit/debug"
	"pure-game-kit/internal"
	"pure-game-kit/tiled/property"
)

type Tileset struct {
	Project    *Project
	Properties map[string]any
	Tiles      []*Tile
}

func NewTileset(tilesetId string, project *Project) *Tileset {
	var data, _ = internal.TiledTilesets[tilesetId]
	if data == nil {
		debug.LogError("Failed to create tileset: \"", tilesetId, "\"\nNo data is loaded with this tileset id.")
		return nil
	}

	var result = Tileset{Project: project}
	result.initProperties(data)
	result.initTiles(data)
	return &result
}

//=================================================================
// private

func (t *Tileset) initProperties(data *internal.Tileset) {
	t.Properties = make(map[string]any)
	t.Properties[property.TilesetName] = data.Name
	t.Properties[property.TilesetClass] = data.Class
	t.Properties[property.TilesetTileWidth] = data.TileWidth
	t.Properties[property.TilesetTileHeight] = data.TileHeight
	t.Properties[property.TilesetColumns] = data.Columns
	t.Properties[property.TilesetRows] = data.TileCount / data.Columns
	t.Properties[property.TilesetOffsetX] = data.Offset.X
	t.Properties[property.TilesetOffsetY] = data.Offset.Y
	t.Properties[property.TilesetSpacing] = data.Spacing

	for _, prop := range data.Properties {
		t.Properties[prop.Name] = parseProperty(prop, t.Project)
	}
}
func (t *Tileset) initTiles(data *internal.Tileset) {
	t.Tiles = make([]*Tile, len(data.Tiles))

	for i, tile := range data.Tiles {
		t.Tiles[i] = NewTile(data.AssetId, tile.Id, t.Project)
	}
}
