package assets

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/angle"
	"pure-game-kit/packages/utility/collection"
	"pure-game-kit/packages/utility/file"
	"pure-game-kit/packages/utility/is"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/path"
	"pure-game-kit/packages/utility/point"
	"pure-game-kit/packages/utility/storage"
	"pure-game-kit/packages/utility/text"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type TileLayerId uint8
type TileAtlasId uint8

func LoadTileAtlas(pngPath string, tileSize int) TileAtlasId {
	internal.TileAtlasNextId++
	var id, imageId = internal.TileAtlasNextId, LoadImage(pngPath)
	rl.SetTextureFilter(internal.Images[int32(imageId)].Texture, rl.FilterPoint)
	var atlas = &internal.TileAtlas{
		ImageId: int32(imageId), TileSize: tileSize, ShapesPerTile: make(map[uint16][][6]float32)}
	internal.TileAtlases[id] = atlas
	return TileAtlasId(id)
}
func LoadTileLayer(columns, rows int) TileLayerId {
	internal.TileLayerNextId++
	columns, rows = number.Limit(columns, 1, 2048), number.Limit(rows, 1, 2048)

	var id = internal.TileLayerNextId
	var data = &internal.TileLayer{Image: rl.GenImageColor(columns, rows, rl.Blank), CellsWithPoints: make(map[int]struct{})}
	var tex = rl.LoadTextureFromImage(data.Image)
	rl.SetTextureFilter(tex, rl.FilterPoint)
	data.Texture = tex
	internal.TileLayers[id] = data
	return TileLayerId(id)
}

func LoadTiledLayers(tmxPath string) (atlasId TileAtlasId, layerIds []TileLayerId) {
	var tileAtlas, tiled = loadTiled(tmxPath)
	var result, dir = make(map[int]TileLayerId), path.Folder(tmxPath)
	if tileAtlas != nil {
		var w, _ = tileAtlas.TileWidth, tileAtlas.TileHeight
		dir = path.New(dir, tileAtlas.Source)
		dir = path.New(path.Folder(dir), tileAtlas.Image.Source)
		atlasId = LoadTileAtlas(dir, w)
	}
	loadLayersRecursively(result, tmxPath, atlasId, tiled, &tiled.layers)
	for _, id := range tiled.LayerIdsInOrder {
		layerIds = append(layerIds, result[id])
	}
	return atlasId, layerIds
}

// private ========================================================

type tiled struct {
	Width     int        `xml:"width,attr"`
	Height    int        `xml:"height,attr"`
	TileAtlas *tileAtlas `xml:"tileset"`
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
		Width    float32       `xml:"width,attr"`
		Height   float32       `xml:"height,attr"`
		X        float32       `xml:"x,attr"`
		Y        float32       `xml:"y,attr"`
		Rotation float32       `xml:"rotation,attr"`
		Polygon  *objectPoints `xml:"polygon"`
		Polyline *objectPoints `xml:"polyline"`
		Ellipse  *struct{}     `xml:"ellipse"`
		Point    *struct{}     `xml:"point"`
		Capsule  *struct{}     `xml:"capsule"`
	} `xml:"object"`
}
type layerAny struct {
	XMLName   xml.Name
	Id        int         `xml:"id,attr"`
	SubLayers []*layerAny `xml:",any"`
}
type objectPoints struct {
	Points string `xml:"points,attr"`
}
type tileAtlas struct {
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

	Points [][6]float32
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

func loadTiled(tmxFilePath string) (*tileAtlas, *tiled) {
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
		if t.Objects != nil {
			t.Points = loadLayerObjects(t.Objects)
		}
	}
	return atlas, tiled
}
func loadLayersRecursively(result map[int]TileLayerId, tmxFilePath string, atlasId TileAtlasId, tiled *tiled, layers *layers) {
	for _, layer := range layers.LayersTiles {
		result[layer.Id] = loadLayerTiles(atlasId, tiled, layer)
	}
	for _, layer := range layers.LayersObjects {
		internal.TileLayerNextId++
		var id = internal.TileLayerNextId
		internal.TileLayers[id] = &internal.TileLayer{Objects: loadLayerObjects(layer)}
		result[layer.Id] = TileLayerId(id)
	}
	for _, group := range layers.LayersGroups {
		loadLayersRecursively(result, tmxFilePath, atlasId, tiled, group)
	}
}
func loadLayerTiles(atlasId TileAtlasId, tiled *tiled, layer *layerTiles) TileLayerId {
	var tileData = text.Trim(layer.TileData.Tiles)
	var tiles, csv = make([]uint32, tiled.Width*tiled.Height), layer.TileData.Encoding == "csv"
	var data *internal.TileLayer
	var atlas = internal.TileAtlases[uint8(atlasId)]
	var dataId = LoadTileLayer(tiled.Width, tiled.Height)
	data = internal.TileLayers[uint8(dataId)]

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

	var pixels = make([]rl.Color, data.Texture.Width*data.Texture.Height)
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
			var tile = tiled.TileAtlas.TilesLookUp[id]
			if tile != nil && tile.Objects != nil && atlas != nil {
				var tilePts, has = atlas.ShapesPerTile[uint16(tile.Id)]
				if !has {
					tilePts = loadLayerObjects(tile.Objects)
					atlas.ShapesPerTile[uint16(tile.Id)] = tilePts
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
func loadLayerObjects(layer *layerObjects) [][6]float32 {
	var result [][6]float32
	for _, o := range layer.Objects {

		if o.Polygon != nil {
			var pts = parsePointPairs(o.Polygon.Points)
			result = append(result, edgesToLineShapes(o.X, o.Y, o.Rotation, pts, true)...)
			continue
		}
		if o.Polyline != nil {
			var pts = parsePointPairs(o.Polyline.Points)
			result = append(result, edgesToLineShapes(o.X, o.Y, o.Rotation, pts, false)...)
			continue
		}

		var cx, cy = o.X + o.Width/2, o.Y + o.Height/2

		if o.Point != nil {
			result = append(result, [6]float32{o.X, o.Y, 0, 0, 0, 1}) // point = zero-size circle
		} else if o.Ellipse != nil {
			result = append(result, [6]float32{cx, cy, o.Width, o.Height, o.Rotation, 1})
		} else if o.Capsule != nil {
			result = append(result, [6]float32{cx, cy, o.Width, o.Height, o.Rotation, 1})
		} else { // assume rectangle
			result = append(result, [6]float32{cx, cy, o.Width, o.Height, o.Rotation, 0})
		}
	}
	return result
}

func parsePointPairs(data string) []float32 {
	var pts []float32
	var trimmed = text.Trim(data)
	if trimmed == "" {
		return pts
	}
	for _, pt := range text.Split(trimmed, " ") {
		var xy = text.Split(pt, ",")
		if len(xy) == 2 {
			pts = append(pts, text.ToNumber[float32](xy[0]), text.ToNumber[float32](xy[1]))
		}
	}
	return pts
}

func edgesToLineShapes(originX, originY, rotation float32, pts []float32, closed bool) [][6]float32 {
	if len(pts) < 4 {
		return nil
	}
	var n = len(pts) / 2
	var edges = n
	if !closed {
		edges = n - 1
	}
	var result = make([][6]float32, 0, edges)
	for i := range edges {
		var next = (i + 1) % n
		var x1, y1 = pts[i*2], pts[i*2+1]
		var x2, y2 = pts[next*2], pts[next*2+1]
		x1, y1 = point.RotateAroundPoint(originX+x1, originY+y1, originX, originY, rotation)
		x2, y2 = point.RotateAroundPoint(originX+x2, originY+y2, originX, originY, rotation)
		var mx, my = (x1 + x2) / 2, (y1 + y2) / 2
		var dx, dy = x2 - x1, y2 - y1
		var length = number.SquareRoot(dx*dx + dy*dy)
		var ang = angle.BetweenPoints(x1, y1, x2, y2)
		result = append(result, [6]float32{mx, my, length, 0, ang, 0})
	}
	return result
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
