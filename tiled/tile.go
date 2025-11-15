package tiled

import (
	"pure-game-kit/data/path"
	"pure-game-kit/debug"
	"pure-game-kit/graphics"
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

func (tile *Tile) Sprite() *graphics.Sprite {
	var atlasId, hasAtlas = tile.OwnerTileset.Properties[property.TilesetAtlasId]
	if hasAtlas {
		atlasId = path.New(atlasId.(string), text.New(tile.Properties[property.TileId]))
	} else {
		atlasId = tile.Properties[property.TileImage]
	}
	var width = tile.Properties[property.TileWidth].(int)
	var height = tile.Properties[property.TileHeight].(int)
	var sprite = graphics.NewSprite(atlasId.(string), 0, 0)
	sprite.Width, sprite.Height = float32(width), float32(height)
	sprite.PivotX, sprite.PivotY = 0, 0
	return sprite
}

//=================================================================

func newTile(tilesetId string, tileId uint32, owner *Tileset) *Tile {
	var data, _ = internal.TiledTilesets[tilesetId]
	if data == nil {
		debug.LogError("Failed to create tile: \"", tilesetId, "/", tileId, "\"\n",
			"No data is loaded with this tileset id.")
		return nil
	}

	var tileData = data.MappedTiles[tileId]
	if tileData == nil { // requested tile has no data, so create default one
		var result = Tile{OwnerTileset: owner}
		var tileWidth = owner.Properties[property.TilesetTileWidth]
		var tileHeight = owner.Properties[property.TilesetTileHeight]
		result.Properties = map[string]any{
			property.TileId:          tileId,
			property.TileWidth:       tileWidth,
			property.TileHeight:      tileHeight,
			property.TileProbability: 1,
		}
		return &result
	}

	var result = Tile{OwnerTileset: owner}
	result.initProperties(data, tileData)
	result.initObjects(tileData)
	return &result
}

//=================================================================

func (tile *Tile) initProperties(tilesetData *internal.Tileset, data *internal.TilesetTile) {
	var w, h = tilesetData.TileWidth, tilesetData.TileHeight

	tile.Properties = make(map[string]any)
	tile.Properties[property.TileId] = data.Id
	tile.Properties[property.TileClass] = data.Class
	tile.Properties[property.TileProbability] = text.ToNumber[float32](defaultValueText(data.Probability, "1"))

	if data.Image != nil {
		w, h = data.Image.Width, data.Image.Height
		tile.Properties[property.TileImage] = data.TextureId
	}

	tile.Properties[property.TileWidth] = w
	tile.Properties[property.TileHeight] = h

	for _, prop := range data.Properties {
		tile.Properties[prop.Name] = parseProperty(prop, tile.OwnerTileset.Project)
	}
}
func (tile *Tile) initObjects(data *internal.TilesetTile) {
	if len(data.CollisionLayers) > 0 {
		var objs = data.CollisionLayers[0].Objects
		tile.Objects = make([]*Object, len(objs))
		for i, obj := range objs {
			tile.Objects[i] = newObject(obj, tile, nil)
		}
	}
}
