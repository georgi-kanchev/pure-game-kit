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
	Tiles      map[uint32]*Tile
}

func (tileset *Tileset) Sprites() []*graphics.Sprite {
	var sprites = []*graphics.Sprite{}
	tileset.forEachTile(true, func(tile *Tile, x, y, w, h float32, sprite *graphics.Sprite) {
		sprite.X, sprite.Y = x, y-sprite.Height
		sprites = append(sprites, sprite)
	})
	return sprites
}
func (tileset *Tileset) Shapes() []*geometry.Shape {
	var result = []*geometry.Shape{}
	tileset.forEachTile(false, func(tile *Tile, x, y, w, h float32, sprite *graphics.Sprite) {
		var shapes = tile.Shapes()
		for _, shape := range shapes {
			shape.X += x
			shape.Y += y - h
			result = append(result, shape)
		}
	})
	return result
}
func (tileset *Tileset) Lines() [][2]float32 {
	var result = [][2]float32{}
	tileset.forEachTile(false, func(tile *Tile, x, y, w, h float32, sprite *graphics.Sprite) {
		var lines = tile.Lines()
		for i := range lines {
			lines[i][0] += x
			lines[i][1] += y - h
			result = append(result, lines[i])
		}
	})
	return result
}
func (tileset *Tileset) Points() [][2]float32 {
	var result = [][2]float32{}
	tileset.forEachTile(false, func(tile *Tile, x, y, w, h float32, sprite *graphics.Sprite) {
		var points = tile.Points()
		for i := range points {
			points[i][0] += x
			points[i][1] += y - h
			result = append(result, points[i])
		}
	})
	return result
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
	tileset.Tiles = make(map[uint32]*Tile, data.TileCount)

	for _, tile := range data.Tiles {
		tileset.Tiles[tile.Id] = newTile(data.AssetId, tile.Id, tileset)
	}

	if data.Columns != 0 {
		var cols, rows = data.Columns, data.TileCount / data.Columns
		for i := range cols {
			for j := range rows {
				var tileId = uint32(number.Indexes2DToIndex1D(i, j, rows, data.Columns))
				if tileset.Tiles[tileId] == nil {
					tileset.Tiles[tileId] = newTile(data.AssetId, tileId, tileset)
				}
			}
		}
	}
}

func (tileset *Tileset) forEachTile(isSprite bool, action func(t *Tile, x, y, w, h float32, s *graphics.Sprite)) {
	var columns = tileset.Properties[property.TilesetColumns].(int)
	var x, y float32 = 0, 0
	var keys = []uint32{}

	for k := range tileset.Tiles {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	for i, id := range keys {
		var tile = tileset.Tiles[id]
		var width = float32(tile.Properties[property.TileWidth].(int))
		var height = float32(tile.Properties[property.TileHeight].(int))
		var sprite *graphics.Sprite

		if isSprite {
			sprite = tile.Sprite()
			width, height = sprite.Width, sprite.Height
		}

		x += width
		if i%columns == 0 {
			x = 0
			y += height
		}

		action(tile, x, y, width, height, sprite)
	}
}
