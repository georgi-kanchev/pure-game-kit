package tiled

import (
	"pure-game-kit/data/path"
	"pure-game-kit/debug"
	"pure-game-kit/execution/condition"
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/internal"
	"pure-game-kit/tiled/property"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/number"
	"slices"
)

type Tileset struct {
	Project    *Project
	Properties map[string]any
	Tiles      map[uint32]*Tile
}

func (t *Tileset) FindTileBy(property string, value any) []*Tile {
	var result = []*Tile{}
	for _, tile := range t.Tiles {
		var curValue, has = tile.Properties[property]
		if has && value == curValue {
			result = append(result, tile)
		}
	}
	return result
}

func (t *Tileset) ExtractSprites() []*graphics.Sprite {
	var sprites = []*graphics.Sprite{}
	t.forEachTile(true, func(tile *Tile, x, y, w, h, scW, scH float32, sprite *graphics.Sprite) {
		sprite.Width, sprite.Height = w, h
		sprite.X, sprite.Y = x, y-h
		sprites = append(sprites, sprite)
	})
	return sprites
}
func (t *Tileset) ExtractShapes() []*geometry.Shape {
	var result = []*geometry.Shape{}
	t.forEachTile(false, func(tile *Tile, x, y, w, h, scW, scH float32, sprite *graphics.Sprite) {
		var shapes = tile.ExtractShapes()
		for _, shape := range shapes {
			shape.X, shape.Y = shape.X*scW+x, shape.Y*scH+y-h
			shape.ScaleX, shape.ScaleY = scW, scH
			result = append(result, shape)
		}
	})
	return result
}
func (t *Tileset) ExtractLines() [][2]float32 {
	var result = [][2]float32{}
	t.forEachTile(false, func(tile *Tile, x, y, w, h, scW, scH float32, sprite *graphics.Sprite) {
		var lines = tile.ExtractLines()
		for i := range lines {
			lines[i][0], lines[i][1] = lines[i][0]*scW+x, lines[i][1]*scH+y-h
			result = append(result, lines[i])
		}
	})
	return result
}
func (t *Tileset) ExtractPoints() [][2]float32 {
	var result = [][2]float32{}
	t.forEachTile(false, func(tile *Tile, x, y, w, h, scW, scH float32, sprite *graphics.Sprite) {
		var points = tile.ExtractPoints()
		for i := range points {
			points[i][0], points[i][1] = points[i][0]*scW+x, points[i][1]*scH+y-h
			result = append(result, points[i])
		}
	})
	return result
}

//=================================================================

func (t *Tileset) Draw(camera *graphics.Camera) {
	draw(camera, t.ExtractSprites(), nil, t.ExtractShapes(),
		t.ExtractPoints(), t.ExtractLines(), palette.White)
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

func (t *Tileset) initProperties(data *internal.Tileset) {
	var rows = 0
	if data.Columns != 0 {
		rows = data.TileCount / data.Columns
	}

	t.Properties = make(map[string]any)
	t.Properties[property.TilesetName] = data.Name
	t.Properties[property.TilesetClass] = data.Class
	t.Properties[property.TilesetTileWidth] = data.TileWidth
	t.Properties[property.TilesetTileHeight] = data.TileHeight
	t.Properties[property.TilesetColumns] = data.Columns
	t.Properties[property.TilesetRows] = rows
	t.Properties[property.TilesetSpacing] = data.Spacing
	t.Properties[property.TilesetOffsetX] = 0
	t.Properties[property.TilesetOffsetY] = 0
	t.Properties[property.TilesetRenderSize] = data.TileRenderSize
	t.Properties[property.TilesetFillMode] = data.FillMode

	if data.Offset != nil {
		t.Properties[property.TilesetOffsetX] = data.Offset.X
		t.Properties[property.TilesetOffsetY] = data.Offset.Y
	}

	if data.Image != nil {
		t.Properties[property.TilesetAtlasId] = path.New(path.Folder(data.AssetId), data.Image.Source)
	}

	for _, prop := range data.Properties {
		t.Properties[prop.Name] = parseProperty(prop, t.Project)
	}
}
func (t *Tileset) initTiles(data *internal.Tileset) {
	t.Tiles = make(map[uint32]*Tile, data.TileCount)

	for _, tile := range data.Tiles {
		t.Tiles[tile.Id] = newTile(data.AssetId, tile.Id, t)
	}

	if data.Columns != 0 {
		var cols, rows = data.Columns, data.TileCount / data.Columns
		for i := range cols {
			for j := range rows {
				var tileId = uint32(number.Indexes2DToIndex1D(i, j, rows, data.Columns))
				if t.Tiles[tileId] == nil {
					t.Tiles[tileId] = newTile(data.AssetId, tileId, t)
				}
			}
		}
	}
}

func (t *Tileset) forEachTile(isSprite bool,
	action func(t *Tile, x, y, w, h, scW, scH float32, s *graphics.Sprite)) {
	var columns = t.Properties[property.TilesetColumns].(int)
	var x, y float32 = 0, 0
	var keys = []uint32{}
	var tileW = float32(t.Properties[property.TilesetTileWidth].(int))
	var tileH = float32(t.Properties[property.TilesetTileHeight].(int))
	var renderSize = t.Properties[property.TilesetRenderSize].(string)
	var fillMode = t.Properties[property.TilesetFillMode].(string)

	for k := range t.Tiles {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	for i, id := range keys {
		var tile = t.Tiles[id]
		var width = float32(tile.Properties[property.TileWidth].(int))
		var height = float32(tile.Properties[property.TileHeight].(int))
		var sprite *graphics.Sprite
		var scX, scY = tileW / width, tileH / height
		var ratioW, ratioH float32 = 1, 1

		if renderSize == "grid" && fillMode == "" {
			width, height = tileW, tileH
		}

		if isSprite {
			sprite = tile.ExtractSprite()
			width, height = sprite.Width, sprite.Height
		}

		width, height, ratioW, ratioH = t.tileRenderSize(width, height, tileW, tileH)

		if renderSize == "grid" && fillMode == "preserve-aspect-fit" {
			scX, scY = scX*ratioW, scY*ratioH
		}

		if renderSize != "grid" {
			scX, scY = 1, 1
		}

		if columns != 0 {
			x += width
			if i%columns == 0 {
				x = 0
				y += height
			}
		}

		action(tile, x, y, width, height, scX, scY, sprite)

		if columns == 0 {
			x += width
		}
	}
}
func (t *Tileset) tileRenderSize(w, h, tileW, tileH float32) (newW, newH, ratioW, ratioH float32) {
	ratioW, ratioH = 1, 1
	var renderSize = t.Properties[property.TilesetRenderSize].(string)
	var fillMode = t.Properties[property.TilesetFillMode].(string)

	ratioH = condition.If(w > h, h/w, ratioH)
	ratioW = condition.If(w <= h, w/h, ratioW)

	if renderSize == "grid" {
		w, h = tileW, tileH

		if fillMode == "preserve-aspect-fit" {
			w *= ratioW
			h *= ratioH
		}
	}
	return w, h, ratioW, ratioH
}
