package tiled

import (
	"pure-game-kit/data/path"
	"pure-game-kit/debug"
	"pure-game-kit/internal"
	"pure-game-kit/tiled/property"
	"pure-game-kit/utility/text"
)

type Tile struct {
	Properties map[string]any
	Objects    []*Object

	IsAnimating bool

	OwnerTileset *Tileset
}

func newTile(tilesetId string, tileId uint32, owner *Tileset) *Tile {
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

	var result = Tile{OwnerTileset: owner}
	result.initProperties(data, tileData)
	result.initObjects(tileData)
	return &result
}

//=================================================================

func (t *Tile) initProperties(tilesetData *internal.Tileset, data *internal.TilesetTile) {
	var w, h = tilesetData.TileWidth, tilesetData.TileHeight

	t.Properties = make(map[string]any)
	t.Properties[property.TileId] = data.Id
	t.Properties[property.TileClass] = data.Class
	t.Properties[property.TileProbability] = text.ToNumber[float32](defaultValueText(data.Probability, "1"))

	if data.Image != nil {
		w, h = data.Image.Width, data.Image.Height
		t.Properties[property.TileImage] = path.New(path.Folder(tilesetData.AssetId), data.Image.Source)
	}

	t.Properties[property.TileWidth] = w
	t.Properties[property.TileHeight] = h

	for _, prop := range data.Properties {
		t.Properties[prop.Name] = parseProperty(prop, t.OwnerTileset.Project)
	}
}
func (t *Tile) initObjects(data *internal.TilesetTile) {
	if len(data.CollisionLayers) > 0 {
		var objs = data.CollisionLayers[0].Objects
		t.Objects = make([]*Object, len(objs))
		for i, obj := range objs {
			t.Objects[i] = newObject(obj, t, nil)
		}
	}
}
