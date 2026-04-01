package assets

import (
	"bytes"
	"encoding/binary"
	"pure-game-kit/data/file"
	"pure-game-kit/data/path"
	"pure-game-kit/data/storage"
	"pure-game-kit/execution/condition"
	"pure-game-kit/internal"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/point"
	"pure-game-kit/utility/text"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadTiledLayers(tmxFilePath string) (tileSetId string, tileLayerIds []string) {
	tryCreateWindow()

	var tileset, tiled = loadTiled(tmxFilePath)
	var dir = path.Folder(tmxFilePath)
	var w, h = tileset.TileWidth, tileset.TileHeight
	dir = path.New(dir, tileset.Source)
	dir = path.New(path.Folder(dir), tileset.Image.Source)
	tileSetId = LoadTileSet(dir, w, h)
	tileLayerIds = loadLayersRecursively(tmxFilePath, tileSetId, tiled, &tiled.layers)
	return tileSetId, tileLayerIds
}

//=================================================================
// private

type tiled struct {
	Width   int      `xml:"width,attr"`
	Height  int      `xml:"height,attr"`
	TileSet *tileSet `xml:"tileset"`
	layers
}
type layers struct {
	LayersGroups  []*layers       `xml:"group"`
	LayersTiles   []*layerTiles   `xml:"layer"`
	LayersObjects []*layerObjects `xml:"objectgroup"`
}
type layerTiles struct {
	Name     string `xml:"name,attr"`
	TileData struct {
		Encoding    string `xml:"encoding,attr"`
		Compression string `xml:"compression,attr"`
		Tiles       string `xml:",chardata"`
	} `xml:"data"`
}
type layerObjects struct {
	Name    string `xml:"name,attr"`
	Objects []*struct {
		Width    float32       `xml:"width,attr"`
		Height   float32       `xml:"height,attr"`
		X        float32       `xml:"x,attr"`
		Y        float32       `xml:"y,attr"`
		Rotation float32       `xml:"rotation,attr"`
		Polygon  *objectPoints `xml:"polygon"`
		Polyline *objectPoints `xml:"polyline"`
		Ellipse  *struct{}     `xml:"ellipse"`
		Point    *struct{}     `xml:"point"`
	} `xml:"object"`
}
type objectPoints struct {
	Points string `xml:"points,attr"`
}
type tileSet struct {
	Source string `xml:"source,attr"`
	Image  *struct {
		Source string `xml:"source,attr"`
	} `xml:"image"`
	TileWidth  int     `xml:"tilewidth,attr"`
	TileHeight int     `xml:"tileheight,attr"`
	Tiles      []*tile `xml:"tile"`

	TilesLookUp map[uint32]*tile
}
type tile struct {
	Id        uint32        `xml:"id,attr"`
	Objects   *layerObjects `xml:"objectgroup"`
	Animation *struct {
		Frames []*tileFrame `xml:"frame"`
	} `xml:"animation"`

	Points []float32
}
type tileFrame struct {
	TileId   uint32 `xml:"tileid,attr"`
	Duration int    `xml:"duration,attr"`
}

const flips = 0x80000000 | 0x40000000 | 0x20000000 // flipX | flipY | flipDiag

var flipTable = [8]uint32{ // Index: [X Y D]
	(0 << 31) | (0 << 29), // 0: 000 | default
	(1 << 31) | (1 << 29), // 1: 001 | flip x + rotation 270
	(1 << 31) | (2 << 29), // 2: 010 | flip x + rotation 180
	(0 << 31) | (3 << 29), // 3: 011 | rotation 270
	(1 << 31) | (0 << 29), // 4: 100 | flip x only
	(0 << 31) | (1 << 29), // 5: 101 | rotation 90
	(0 << 31) | (2 << 29), // 6: 110 | rotation 180
	(1 << 31) | (3 << 29), // 7: 111 | flip x + rotation 90
}

func loadTiled(tmxFilePath string) (*tileSet, *tiled) {
	var tiled *tiled
	var mapContent = file.LoadText(tmxFilePath)
	storage.FromXML(mapContent, &tiled)
	if tiled == nil {
		return nil, nil // error is in storage
	}

	if tiled.TileSet == nil {
		return nil, tiled
	}

	var tileSet = tiled.TileSet
	var dir = path.Folder(tmxFilePath)
	if tileSet.Image == nil {
		storage.FromXML(file.LoadText(path.New(dir, tileSet.Source)), &tileSet)
		if tileSet == nil {
			return nil, nil // error is in storage
		}
	}

	tileSet.TilesLookUp = map[uint32]*tile{}
	for _, t := range tileSet.Tiles {
		tileSet.TilesLookUp[t.Id] = t

		if t.Objects != nil {
			t.Points = loadLayerObjects(t.Objects)
		}
	}

	return tileSet, tiled
}
func loadLayersRecursively(tmxFilePath, tileSetId string, tiled *tiled, layers *layers) []string {
	var result []string
	for _, layer := range layers.LayersTiles {
		var id = loadLayerTiles(tmxFilePath, tileSetId, tiled, layer)
		result = append(result, id)
	}
	for _, layer := range layers.LayersObjects {
		var dataId = path.New(tmxFilePath, layer.Name)
		internal.TileLayers[dataId] = &internal.TileLayer{ObjectPoints: loadLayerObjects(layer)}
		result = append(result, dataId)
	}
	for _, group := range layers.LayersGroups {
		result = append(result, loadLayersRecursively(tmxFilePath, tileSetId, tiled, group)...)
	}
	return result
}
func loadLayerTiles(tmxFilePath, tileSetId string, tiled *tiled, layer *layerTiles) string {
	var tileData = text.Trim(layer.TileData.Tiles)
	var tiles = make([]uint32, tiled.Width*tiled.Height)
	var csv = layer.TileData.Encoding == "csv"
	var dataId = ""
	var data *internal.TileLayer
	var tileSet = internal.TileSets[tileSetId]

	dataId = LoadTileData(path.New(tmxFilePath, layer.Name), tiled.Width, tiled.Height)
	data = internal.TileLayers[dataId]

	if layer.TileData.Encoding == "base64" {
		var b64 = text.FromBase64(text.Trim(tileData))
		switch layer.TileData.Compression {
		case "gzip":
			tiles = tilesFromBytes(storage.DecompressGZIP([]byte(b64)))
		case "zlib":
			tiles = tilesFromBytes(storage.DecompressZLIB([]byte(b64)))
		}
	}

	var rows []string
	if csv {
		rows = text.Split(tileData, "\n")
	}

	for i := range tiled.Height {
		var columns []string
		if csv {
			var row = rows[i]
			if text.EndsWith(row, ",") {
				row = text.Part(row, 0, text.Length(row)-1)
			}
			columns = text.Split(row, ",")
		}

		for j := range tiled.Width {
			var index = number.Indexes2DToIndex1D(i, j, tiled.Width, tiled.Height)
			if csv {
				tiles[index] = text.ToNumber[uint32](columns[j])
			}

			if tiles[index] == 0 {
				continue
			}

			var raw = tiles[index]
			var id = ((raw & 0x1FFFFFFF) - 1)
			var tile = tiled.TileSet.TilesLookUp[id]
			if tile != nil && tile.Objects != nil && tileSet != nil {
				var tilePts, has = tileSet.PointsPerTile[uint16(tile.Id)]
				if !has {
					tilePts = loadLayerObjects(tile.Objects)
					tileSet.PointsPerTile[uint16(tile.Id)] = tilePts
				}

				var cellIndex1D = number.Indexes2DToIndex1D(j, i, tiled.Width, tiled.Height)
				data.CellsWithPoints[cellIndex1D] = struct{}{}
				data.LastDirtyTime = internal.Runtime
			}

			var frameCount, frameSpeed, animOffset uint32
			if tile != nil && tile.Animation != nil && len(tile.Animation.Frames) > 0 {
				var totalDuration = 0
				for _, f := range tile.Animation.Frames {
					totalDuration += f.Duration
				}
				var avgDuration = float32(totalDuration) / float32(len(tile.Animation.Frames))
				var targetFPS = 1000.0 / avgDuration
				var s = condition.If(targetFPS <= 1.0, targetFPS*10.0, ((targetFPS-1.0)/0.45)+10.0)
				frameSpeed = uint32(number.Limit(int(s), 0, 31))
				frameCount = uint32(number.Limit(len(tile.Animation.Frames)-1, 0, 15))
				animOffset = uint32((j ^ i) % 16)
			}

			var finalTile uint32 = 0
			finalTile |= (uint32(flipTable[raw>>29]) << 29)
			finalTile |= (frameCount << 25)
			finalTile |= (animOffset << 21)
			finalTile |= (frameSpeed << 16)
			finalTile |= (id & 0xFFFF)

			var r = uint8((finalTile >> 24) & 0xFF)
			var g = uint8((finalTile >> 16) & 0xFF)
			var b = uint8((finalTile >> 8) & 0xFF)
			var a = uint8((finalTile >> 0) & 0xFF)
			var colr = rl.NewColor(r, g, b, a)
			var rect = rl.NewRectangle(float32(j), float32(i), float32(1), float32(1))

			rl.ImageDrawPixel(data.Image, int32(j), int32(i), colr)
			rl.UpdateTextureRec(*data.Texture, rect, []rl.Color{colr})
		}
	}
	return dataId
}
func loadLayerObjects(layer *layerObjects) []float32 {
	var result []float32
	for _, o := range layer.Objects {
		var data = ""
		if o.Polyline != nil {
			data = o.Polyline.Points
		}
		if o.Polygon != nil {
			data = o.Polygon.Points
		}
		if data == "" {
			if o.Point != nil {
				data = text.New(0, ",", 0)
			} else if o.Ellipse != nil {
				const segments = 32
				var rx, ry = o.Width / 2, o.Height / 2
				var step = 360.0 / float32(segments)

				for i := range segments {
					var cx, cy = point.MoveAtAngle(0, 0, float32(i)*step, 1)
					var x, y = (cx + 1) * rx, (cy + 1) * ry // shift from center-based to tiled's top-left-based
					var value = text.New(x, ",", y, " ")
					data += value
				}
			} else { // assume it's a rectangle
				data = text.New(0, ",", 0, " ", o.Width, ",", 0, " ", o.Width, ",", o.Height, " ", 0, ",", o.Height)
			}
		}

		var corners []float32
		var pts = text.Split(text.Trim(data), " ")
		for _, pt := range pts {
			var xy = text.Split(pt, ",")
			if len(xy) == 2 {
				var x, y = text.ToNumber[float32](xy[0]), text.ToNumber[float32](xy[1])
				x, y = point.RotateAroundPoint(o.X+x, o.Y+y, o.X, o.Y, o.Rotation)
				corners = append(corners, x, y)
			}
		}
		corners = append(corners, corners[0], corners[1])
		corners = append(corners, number.NaN(), number.NaN())
		result = append(result, corners...)
	}
	return result
}

func tilesFromBytes(data []byte) []uint32 {
	if len(data)%4 != 0 {
		return nil
	}

	var numElements = len(data) / 4
	var result = make([]uint32, numElements)
	var reader = bytes.NewReader(data)
	var err = binary.Read(reader, binary.LittleEndian, &result)
	if err != nil {
		return nil
	}
	return result
}
