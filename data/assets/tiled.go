package assets

import (
	"bytes"
	"encoding/binary"
	"pure-game-kit/data/file"
	"pure-game-kit/data/path"
	"pure-game-kit/data/storage"
	"pure-game-kit/internal"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/point"
	"pure-game-kit/utility/text"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadTiledPoints(tmxFilePath string, layerName string) [][2]float32 {
	var tiled *tiled
	var fileContent = file.LoadText(tmxFilePath)
	storage.FromXML(fileContent, &tiled)
	if tiled == nil {
		return nil // error is in storage
	}

	return loadLayerObjectsRecursively(layerName, &tiled.layers)
}
func LoadTiledData(tmxFilePath string) []string {
	tryCreateWindow()

	var tiled *tiled
	var fileContent = file.LoadText(tmxFilePath)
	storage.FromXML(fileContent, &tiled)
	if tiled == nil {
		return nil // error is in storage
	}

	return loadLayerTilesRecursively(tmxFilePath, tiled.Width, tiled.Height, &tiled.layers)
}

//=================================================================
// private

type tiled struct {
	Width  int `xml:"width,attr"`
	Height int `xml:"height,attr"`
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

const flipX, flipY, flipDiag uint32 = 0x80000000, 0x40000000, 0x20000000
const flips = flipX | flipY | flipDiag

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

func loadLayerTilesRecursively(tmxFilePath string, w, h int, layers *layers) []string {
	var result []string
	for _, layer := range layers.LayersTiles {
		result = append(result, loadLayerTiles(tmxFilePath, w, h, layer))
	}
	for _, group := range layers.LayersGroups {
		result = append(result, loadLayerTilesRecursively(tmxFilePath, w, h, group)...)
	}
	return result
}
func loadLayerTiles(tmxFilePath string, w, h int, layer *layerTiles) string {
	var tileData = text.Trim(layer.TileData.Tiles)
	var tiles = make([]uint32, w*h)
	var csv = layer.TileData.Encoding == "csv"
	var dataId = LoadTileData(path.New(tmxFilePath, layer.Name), w, h)

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

	for i := range h {
		var columns []string
		if csv {
			var row = rows[i]
			if text.EndsWith(row, ",") {
				row = text.Part(row, 0, text.Length(row)-1)
			}
			columns = text.Split(row, ",")
		}

		for j := range w {
			var index = number.Indexes2DToIndex1D(i, j, w, h)
			if csv {
				tiles[index] = text.ToNumber[uint32](columns[j])
			}

			if tiles[index] == 0 {
				continue
			}

			var raw = tiles[index]
			var packed = flipTable[raw>>29] | ((raw & 0x1FFFFFFF) - 1)
			var data = internal.TileDatas[dataId]
			var r = uint8((packed >> 24) & 0xFF)
			var g = uint8((packed >> 16) & 0xFF)
			var b = uint8((packed >> 8) & 0xFF)
			var a = uint8((packed >> 0) & 0xFF)
			var colr = rl.NewColor(r, g, b, a)
			var rect = rl.NewRectangle(float32(j), float32(i), float32(1), float32(1))

			rl.ImageDrawPixel(data.Image, int32(j), int32(i), colr)
			rl.UpdateTextureRec(*data.Texture, rect, []rl.Color{colr})
		}
	}
	return dataId
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

func loadLayerObjectsRecursively(name string, layers *layers) [][2]float32 {
	var result [][2]float32
	for _, layer := range layers.LayersObjects {
		result = append(result, loadLayerObjects(name, layer)...)
	}
	for _, group := range layers.LayersGroups {
		result = append(result, loadLayerObjectsRecursively(name, group)...)
	}
	return result
}
func loadLayerObjects(name string, layer *layerObjects) [][2]float32 {
	if layer.Name != name {
		return nil
	}

	var result [][2]float32
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

		var corners = [][2]float32{}
		var pts = text.Split(text.Trim(data), " ")
		for _, pt := range pts {
			var xy = text.Split(pt, ",")
			if len(xy) == 2 {
				var x, y = text.ToNumber[float32](xy[0]), text.ToNumber[float32](xy[1])
				x, y = point.RotateAroundPoint(o.X+x, o.Y+y, o.X, o.Y, o.Rotation)
				corners = append(corners, [2]float32{x, y})
			}
		}
		corners = append(corners, [2]float32{number.NaN(), number.NaN()})
		result = append(result, corners...)
	}
	return result
}
