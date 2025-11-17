package tiled

import (
	"pure-game-kit/data/path"
	"pure-game-kit/debug"
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/internal"
	"pure-game-kit/tiled/property"
	"pure-game-kit/utility/number"
	"slices"
)

type Tileset struct {
	Project    *Project
	Properties map[string]any
	Tiles      []*Tile
}

func (tileset *Tileset) Sprites() []*graphics.Sprite {
	var sprites = []*graphics.Sprite{}
	var columns = tileset.Properties[property.TilesetColumns].(int)
	var x, y float32 = 0, 0
	for i, tile := range tileset.Tiles {
		var sprite = tile.Sprite()
		x += sprite.Width
		if i%columns == 0 {
			x = 0
			y += sprite.Height
		}

		sprite.X, sprite.Y = x, y-sprite.Height
		sprites = append(sprites, sprite)
	}
	return sprites
}

func (tileset *Tileset) Shapes() []*geometry.Shape {
	var shapes = []*geometry.Shape{}
	var columns = tileset.Properties[property.TilesetColumns].(int)
	var x, y float32 = 0, 0
	for i, tile := range tileset.Tiles {
		var width = tile.Properties[property.TileWidth].(int)
		var height = tile.Properties[property.TileHeight].(int)
		x += float32(width)
		if i%columns == 0 {
			x = 0
			y += float32(height)
		}

		for _, obj := range tile.Objects {
			var shape = obj.Shape()
			shape.X, shape.Y = x, y-float32(height)
			shapes = append(shapes, shape)
		}
	}
	return shapes
}

//=================================================================

func newTileset(tilesetId string, project *Project) *Tileset {
	var data, _ = internal.TiledTilesets[tilesetId]
	if data == nil {
		debug.LogError("Failed to create tileset: \"", tilesetId, "\"\nNo data is loaded with this tileset id.")
		return nil
	}

	if project != nil {
		var cache, hasCache = project.UniqueTilesets[tilesetId]
		if hasCache { // maps in the same project will try to reuse tilesets instead of load them
			return cache
		}
	}

	var result = Tileset{Project: project}
	result.initProperties(data)
	result.initTiles(data)

	if project != nil {
		project.UniqueTilesets[tilesetId] = &result
	}

	return &result
}

//=================================================================

func (tileset *Tileset) initProperties(data *internal.Tileset) {
	var rows = 0
	if data.Columns != 0 {
		rows = data.TileCount / data.Columns
	}

	tileset.Properties = make(map[string]any)
	tileset.Properties[property.TilesetName] = data.Name
	tileset.Properties[property.TilesetClass] = data.Class
	tileset.Properties[property.TilesetTileWidth] = data.TileWidth
	tileset.Properties[property.TilesetTileHeight] = data.TileHeight
	tileset.Properties[property.TilesetColumns] = data.Columns
	tileset.Properties[property.TilesetRows] = rows
	tileset.Properties[property.TilesetSpacing] = data.Spacing
	tileset.Properties[property.TilesetOffsetX] = 0
	tileset.Properties[property.TilesetOffsetY] = 0

	if data.Offset != nil {
		tileset.Properties[property.TilesetOffsetX] = data.Offset.X
		tileset.Properties[property.TilesetOffsetY] = data.Offset.Y
	}

	if data.Image != nil {
		tileset.Properties[property.TilesetAtlasId] = path.New(path.Folder(data.AssetId), data.Image.Source)
	}

	for _, prop := range data.Properties {
		tileset.Properties[prop.Name] = parseProperty(prop, tileset.Project)
	}
}
func (tileset *Tileset) initTiles(data *internal.Tileset) {
	var tiles = map[uint32]*Tile{}
	var keys = make([]uint32, 0, data.TileCount)
	tileset.Tiles = make([]*Tile, data.TileCount)

	for _, tile := range data.Tiles {
		keys = append(keys, tile.Id)
		tiles[tile.Id] = newTile(data.AssetId, tile.Id, tileset)
	}

	if data.Columns != 0 {
		var cols, rows = data.Columns, data.TileCount / data.Columns
		for i := range cols {
			for j := range rows {
				var tileId = uint32(number.Indexes2DToIndex1D(i, j, rows, data.Columns))
				if tiles[tileId] == nil {
					keys = append(keys, tileId)
					tiles[tileId] = newTile(data.AssetId, tileId, tileset)
				}
			}
		}
	}

	slices.Sort(keys)

	var i = 0
	for _, key := range keys {
		tileset.Tiles[i] = tiles[key]
		i++
	}
}
