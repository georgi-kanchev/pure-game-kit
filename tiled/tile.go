package tiled

import (
	"pure-game-kit/data/path"
	"pure-game-kit/debug"
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/internal"
	"pure-game-kit/tiled/property"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/text"
)

type Tile struct {
	Properties map[string]any
	Objects    []*Object

	OwnerTileset *Tileset
}

func (t *Tile) FindObjectsBy(property string, value any) []*Object {
	var result = []*Object{}
	for _, obj := range t.Objects {
		var curValue, has = obj.Properties[property]
		if has && value == curValue {
			result = append(result, obj)
		}
	}
	return result
}

func (t *Tile) ExtractSprite() *graphics.Sprite {
	var atlasId, hasAtlas = t.OwnerTileset.Properties[property.TilesetAtlasId]
	if hasAtlas {
		atlasId = path.New(atlasId.(string), text.New(t.Properties[property.TileId]))
	} else {
		atlasId = t.Properties[property.TileImage]
	}
	var width = t.Properties[property.TileWidth].(int)
	var height = t.Properties[property.TileHeight].(int)
	var sprite = graphics.NewSprite(atlasId.(string), 0, 0)
	sprite.Width, sprite.Height = float32(width), float32(height)
	sprite.PivotX, sprite.PivotY = 0, 0

	return sprite
}
func (t *Tile) ExtractShapes() []*geometry.Shape {
	var result = []*geometry.Shape{}
	for _, obj := range t.Objects {
		result = append(result, obj.ExtractShapes()...)
	}
	return result
}
func (t *Tile) ExtractLines() [][2]float32 {
	var result = [][2]float32{}
	for i, obj := range t.Objects {
		if i != 0 {
			result = append(result, [2]float32{number.NaN(), number.NaN()})
		}

		result = append(result, obj.ExtractLines()...)
	}
	return result
}
func (t *Tile) ExtractPoints() [][2]float32 {
	var result = [][2]float32{}
	for _, obj := range t.Objects {
		result = append(result, obj.ExtractPoints()...)
	}
	return result
}

//=================================================================

func (t *Tile) Draw(camera *graphics.Camera) {
	draw(camera, []*graphics.Sprite{t.ExtractSprite()}, nil,
		t.ExtractShapes(), t.ExtractPoints(), t.ExtractLines(), palette.White)
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

func (t *Tile) initProperties(tilesetData *internal.Tileset, data *internal.TilesetTile) {
	var w, h = tilesetData.TileWidth, tilesetData.TileHeight

	t.Properties = make(map[string]any)
	t.Properties[property.TileId] = data.Id
	t.Properties[property.TileClass] = data.Class
	t.Properties[property.TileProbability] = text.ToNumber[float32](defaultValueText(data.Probability, "1"))

	if data.Image != nil {
		w, h = data.Image.Width, data.Image.Height
		t.Properties[property.TileImage] = data.TextureId
	}

	t.Properties[property.TileWidth] = w
	t.Properties[property.TileHeight] = h

	for _, prop := range data.Properties {
		t.Properties[prop.Name] = parseProperty(prop, t.OwnerTileset.Project)
	}
}
func (t *Tile) initObjects(data *internal.TilesetTile) {
	if len(data.CollisionLayers) == 0 {
		return
	}

	var objs = data.CollisionLayers[0].Objects
	t.Objects = make([]*Object, len(objs))
	for i, obj := range objs {
		t.Objects[i] = newObject(obj, t, nil)
	}
}
