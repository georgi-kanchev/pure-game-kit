package assets

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/collection"
	"pure-game-kit/packages/utility/file"
	"pure-game-kit/packages/utility/is"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/path"
	"pure-game-kit/packages/utility/storage"
	"pure-game-kit/packages/utility/text"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Expects a Tiled map file with a single embedded tileset atlas.
func LoadTileLayersFromTiled(tmxPath string) []TileLayerId {
	var tileAtlas, tiled = loadTiled(tmxPath)
	var result, dir = make(map[int]TileLayerId), path.Folder(tmxPath)
	var imageId ImageId
	var tileSize int
	if tileAtlas != nil {
		var w, _ = tileAtlas.TileWidth, tileAtlas.TileHeight
		dir = path.New(dir, tileAtlas.Source)
		dir = path.New(path.Folder(dir), tileAtlas.Image.Source)
		imageId = LoadImage(dir)
		rl.SetTextureFilter(internal.Images[int32(imageId)].Texture, rl.FilterPoint)
		tileSize = w
	}
	loadLayersRecursively(result, tmxPath, imageId, tileSize, tiled, &tiled.layers)
	var layerIds []TileLayerId
	for _, id := range tiled.LayerIdsInOrder {
		layerIds = append(layerIds, result[id])
	}
	return layerIds
}

// private ========================================================

type tiled struct {
	Width     int    `xml:"width,attr"`
	Height    int    `xml:"height,attr"`
	TileAtlas *atlas `xml:"tileset"`
	layers
	LayerIdsInOrder []int
}
type layers struct {
	LayersGroups  []*layers       `xml:"group"`
	LayersTiles   []*layerTiles   `xml:"layer"`
	LayersObjects []*layerObjects `xml:"objectgroup"`
}
type layersInOrder struct {
	Layers []*layerAny `xml:",any"`
}
type layerTiles struct {
	Id       int    `xml:"id,attr"`
	Name     string `xml:"name,attr"`
	TileData struct {
		Encoding    string `xml:"encoding,attr"`
		Compression string `xml:"compression,attr"`
		Tiles       string `xml:",chardata"`
	} `xml:"data"`
}
type layerObjects struct {
	Id      int    `xml:"id,attr"`
	Name    string `xml:"name,attr"`
	Objects []*struct {
		Width    float32 `xml:"width,attr"`
		Height   float32 `xml:"height,attr"`
		X        float32 `xml:"x,attr"`
		Y        float32 `xml:"y,attr"`
		Rotation float32 `xml:"rotation,attr"`
		Gid      uint16  `xml:"gid,attr"`
		Polygon  *struct {
			Points string `xml:"points,attr"`
		} `xml:"polygon"`
		Polyline *struct {
			Points string `xml:"points,attr"`
		} `xml:"polyline"`
		Ellipse *struct{} `xml:"ellipse"`
		Point   *struct{} `xml:"point"`
		Capsule *struct{} `xml:"capsule"`
	} `xml:"object"`
}
type layerAny struct {
	XMLName   xml.Name
	Id        int         `xml:"id,attr"`
	SubLayers []*layerAny `xml:",any"`
}
type atlas struct {
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
		Frames []*struct {
			TileId   uint32 `xml:"tileid,attr"`
			Duration int    `xml:"duration,attr"`
		} `xml:"frame"`
	} `xml:"animation"`
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

func loadTiled(tmxFilePath string) (*atlas, *tiled) {
	var tiled *tiled
	var mapContent = file.LoadText(tmxFilePath)
	storage.FromXML(mapContent, &tiled)
	if tiled == nil {
		return nil, nil // error is in storage
	}

	if tiled.TileAtlas == nil {
		return nil, tiled
	}

	var atlas, dir = tiled.TileAtlas, path.Folder(tmxFilePath)
	if atlas.Image == nil {
		storage.FromXML(file.LoadText(path.New(dir, atlas.Source)), &atlas)
		if atlas == nil {
			return nil, nil // error is in storage
		}
	}

	var layersInOrder *layersInOrder
	storage.FromXML(mapContent, &layersInOrder)
	var order = getLayersOrder(layersInOrder.Layers)
	collection.Reverse(order)
	tiled.LayerIdsInOrder = order

	atlas.TilesLookUp = map[uint32]*tile{}
	for _, t := range atlas.Tiles {
		atlas.TilesLookUp[t.Id] = t
	}
	return atlas, tiled
}
func loadLayersRecursively(result map[int]TileLayerId, tmxFilePath string, imageId ImageId, tileSize int, tiled *tiled, layers *layers) {
	for _, layer := range layers.LayersTiles {
		result[layer.Id] = loadLayerTiles(imageId, tileSize, tiled, layer)
	}
	for _, layer := range layers.LayersObjects {
		internal.TileLayerNextId++
		var id = internal.TileLayerNextId
		var objs, pts = loadLayerObjects(layer)
		internal.TileLayers[id] = &internal.TileLayer{
			Objects: objs, Paths: pts, Columns: tiled.Width, Rows: tiled.Height, TileSize: tiled.TileAtlas.TileWidth}
		result[layer.Id] = TileLayerId(id)
	}
	for _, group := range layers.LayersGroups {
		loadLayersRecursively(result, tmxFilePath, imageId, tileSize, tiled, group)
	}
}
func loadLayerTiles(imageId ImageId, tileSize int, tiled *tiled, layer *layerTiles) TileLayerId {
	var tileData = text.Trim(layer.TileData.Tiles)
	var tiles, csv = make([]uint32, tiled.Width*tiled.Height), layer.TileData.Encoding == "csv"
	var dataId = LoadTileLayer(tiled.Width, tiled.Height, tileSize, imageId)
	var data = internal.TileLayers[uint8(dataId)]

	if layer.TileData.Encoding == "base64" {
		var b64 = storage.FromBase64(text.Trim(tileData))
		switch layer.TileData.Compression {
		case "gzip":
			tiles = tilesFromBytes(storage.DecompressGZIP([]byte(b64)))
		case "zlib":
			tiles = tilesFromBytes(storage.DecompressZLIB([]byte(b64)))
		}
	}

	var pixels = make([]rl.Color, data.Texture.Width*data.Texture.Height)
	for i := range tiled.Height {
		var row string
		if csv {
			row = text.SplitIndex(tileData, "\n", i)
			if text.EndsWith(row, ",") {
				row = text.Part(row, 0, text.Length(row)-1)
			}
		}

		for j := range tiled.Width {
			var index = number.Indexes2DToIndex1D(i, j, tiled.Width, tiled.Height)
			if csv {
				tiles[index] = text.ToNumber[uint32](text.SplitIndex(row, ",", j))
			}

			if tiles[index] == 0 {
				continue
			}

			var raw = tiles[index]
			var id = ((raw & 0x1FFFFFFF) - 1)
			var tile = tiled.TileAtlas.TilesLookUp[id]
			if tile != nil && tile.Objects != nil {
				var tileShapes, has = data.ShapesPerTile[uint16(tile.Id)]
				if !has {
					tileShapes, _ = loadLayerObjects(tile.Objects)
					data.ShapesPerTile[uint16(tile.Id)] = tileShapes
				}

				var cellIndex1D = number.Indexes2DToIndex1D(j, i, tiled.Width, tiled.Height)
				data.CellsWithPoints[cellIndex1D] = struct{}{}
			}

			var frameCount, frameSpeed, animOffset uint32
			if tile != nil && tile.Animation != nil && len(tile.Animation.Frames) > 0 {
				var totalDuration = 0
				for _, f := range tile.Animation.Frames {
					totalDuration += f.Duration
				}
				var avgDuration = float32(totalDuration) / float32(len(tile.Animation.Frames))
				var targetFPS = 1000.0 / avgDuration
				var s = targetFPS * 10.0
				if targetFPS > 1.0 {
					s = ((targetFPS - 1.0) / 0.45) + 10.0
				}

				frameSpeed = uint32(number.Limit(int(s), 0, 31))
				frameCount = uint32(number.Limit(len(tile.Animation.Frames)-1, 0, 15))
				animOffset = uint32((j ^ i) % 16)
			}

			if raw > 1000 {
				print()
			}

			var finalTile uint32 = 0
			finalTile |= flipTable[raw>>29]
			finalTile |= frameCount << 25
			finalTile |= animOffset << 21
			finalTile |= frameSpeed << 16
			finalTile |= id & 0xFFFF

			var r, g = uint8((finalTile >> 24) & 0xFF), uint8((finalTile >> 16) & 0xFF)
			var b, a = uint8((finalTile >> 8) & 0xFF), uint8((finalTile >> 0) & 0xFF)
			var col = rl.NewColor(r, g, b, a)
			var index1D = number.Indexes2DToIndex1D(i, j, int(data.Image.Width), int(data.Image.Height))
			pixels[index1D] = col
			rl.ImageDrawPixel(data.Image, int32(j), int32(i), col)
		}
	}
	var rect = rl.NewRectangle(0, 0, float32(data.Texture.Width), float32(data.Texture.Height))
	rl.UpdateTextureRec(data.Texture, rect, pixels)
	return dataId
}
func loadLayerObjects(layer *layerObjects) (shapes [][6]float32, pts []float32) {
	for _, o := range layer.Objects {
		if o.Polygon != nil || o.Polyline != nil { // lines/polygons not handled here
			if pts != nil {
				pts = append(pts, number.NaN(), number.NaN())
			}
			var ptsString string
			if o.Polygon != nil {
				ptsString = o.Polygon.Points
			} else if o.Polyline != nil {
				ptsString = o.Polyline.Points
			}
			var sin, cos float32 = 0, 1
			if o.Rotation != 0 {
				sin, cos = internal.SinCos(o.Rotation)
			}
			var splitCount = text.SplitCount(ptsString, " ")
			for k := 0; k < splitCount; k++ {
				var v = text.SplitIndex(ptsString, " ", k)
				var px, py = text.ToNumber[float32](text.SplitIndex(v, ",", 0)), text.ToNumber[float32](text.SplitIndex(v, ",", 1))
				pts = append(pts, o.X+px*cos-py*sin, o.Y+px*sin+py*cos)
			}
		} else if o.Gid != 0 {
			var cx, cy = pivotToCenter(o.X, o.Y, o.Width/2, -o.Height/2, o.Rotation)
			shapes = append(shapes, [6]float32{cx, cy, o.Width, o.Height, o.Rotation, 0})
		} else if o.Point != nil {
			shapes = append(shapes, [6]float32{o.X, o.Y, 5, 5, 0, 1})
		} else if o.Ellipse != nil || o.Capsule != nil {
			var cx, cy = pivotToCenter(o.X, o.Y, o.Width/2, o.Height/2, o.Rotation)
			shapes = append(shapes, [6]float32{cx, cy, o.Width, o.Height, o.Rotation, 1})
		} else { // assume rectangle
			var cx, cy = pivotToCenter(o.X, o.Y, o.Width/2, o.Height/2, o.Rotation)
			shapes = append(shapes, [6]float32{cx, cy, o.Width, o.Height, o.Rotation, 0})
		}
	}
	return shapes, pts
}

func pivotToCenter(x, y, halfW, halfH, angle float32) (cx, cy float32) {
	// converts a Tiled pivot position (top-left for regular, bottom-center for GID) to a center-center position
	if angle == 0 {
		return x + halfW, y + halfH
	}
	var sin, cos = internal.SinCos(angle)
	return x + halfW*cos - halfH*sin, y + halfW*sin + halfH*cos
}

func tilesFromBytes(data []byte) []uint32 {
	if len(data)%4 != 0 {
		return nil
	}
	var result, reader = make([]uint32, len(data)/4), bytes.NewReader(data)
	var err = binary.Read(reader, binary.LittleEndian, &result)
	if err != nil {
		return nil
	}
	return result
}
func getLayersOrder(layers []*layerAny) []int {
	var result = []int{}
	collection.Reverse(layers)
	for _, layer := range layers {
		if is.OneOf(layer.XMLName.Local, "layer", "objectgroup") {
			result = append(result, layer.Id)
		}
		result = append(result, getLayersOrder(layer.SubLayers)...)
	}
	return result
}
