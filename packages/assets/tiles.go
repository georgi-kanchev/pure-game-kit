package assets

import (
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/collection"
	"pure-game-kit/packages/utility/number"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Tile struct {
	// data in bits (32)
	// bits 31..31 = flip					(0 or 1)
	// bits 30..29 = rotations 				(4 values: 0, 90, 180, 270)
	// bits 28..25 = animation frame count 	(0 to 15)
	// bits 24..21 = animation offset		(0 to 15)
	// bits 20..16 = animation frames/s		(0 to 31)
	// bits 15..00 = tile id				(0 to 65535)

	Id          uint16
	Rotations90 byte // 90 degree turns, ranged 0..3 (possible values: 0, 90, 180, 270)
	Flip        bool

	FrameCount  byte // Ranged 0..15 (sequential tile count in the atlas)
	FrameOffset byte // Ranged 0..15
	FrameSpeed  byte // Ranged 0..31
}

func NewTile(id uint16) Tile {
	return Tile{Id: id}
}
func NewTileOriented(id uint16, rotations90 byte, flip bool) Tile {
	return Tile{Id: id, Rotations90: rotations90, Flip: flip}
}
func NewTileAnimated(id uint16, frameCount, frameOffset, frameSpeed byte) Tile {
	return Tile{Id: id, FrameCount: frameCount, FrameSpeed: frameSpeed, FrameOffset: frameOffset}
}

//=================================================================

type TileLayerId uint8
type TileAtlasId uint8

func LoadTileAtlas(pngPath string, tileWidth, tileHeight int) TileAtlasId {
	internal.TileAtlasNextId++
	var id, imageId = internal.TileAtlasNextId, LoadImage(pngPath)
	var atlas = &internal.TileAtlas{
		ImageId: int32(imageId), TileWidth: tileWidth, TileHeight: tileHeight, PointsPerTile: make(map[uint16][]float32)}
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

//=================================================================

func (l TileLayerId) SetTile(column, row int, tile Tile) {
	l.SetTileArea(column, row, 1, 1, tile)
}
func (l TileLayerId) SetTileArea(column, row, width, height int, tile Tile) {
	var layer = internal.TileLayers[uint8(l)]
	var atlas = internal.TileAtlases[uint8(layer.AtlasId)]
	if layer == nil {
		return
	}

	var packed = uint32(tile.Id&0xFFFF) | uint32(tile.FrameSpeed&0x1F)<<16 | uint32(tile.FrameOffset&0x0F)<<21 |
		uint32(tile.FrameCount&0x0F)<<25 | uint32(tile.Rotations90&0x03)<<29

	if tile.Flip {
		packed |= (1 << 31)
	}

	var r, g = uint8((packed >> 24) & 0xFF), uint8((packed >> 16) & 0xFF)
	var b, a = uint8((packed >> 8) & 0xFF), uint8((packed >> 0) & 0xFF)
	var colr, rect = rl.NewColor(r, g, b, a), rl.NewRectangle(float32(column), float32(row), float32(width), float32(height))
	var w, h = l.Size()
	var _, cellHasPts = atlas.PointsPerTile[tile.Id]

	for i := row; i < row+height; i++ {
		for j := column; j < column+width; j++ {
			var prevTile = l.TileAtCell(j, i)
			var _, prevCellHasPts = atlas.PointsPerTile[prevTile.Id]
			if !prevCellHasPts && !cellHasPts {
				continue
			}

			var index1D = number.Indexes2DToIndex1D(j, i, w, h)
			layer.LastDirtyTime = internal.Runtime

			if cellHasPts {
				layer.CellsWithPoints[index1D] = struct{}{}
			} else {
				delete(layer.CellsWithPoints, index1D)
			}
		}
	}

	rl.ImageDrawRectangle(layer.Image, int32(column), int32(row), int32(width), int32(height), colr)
	rl.UpdateTextureRec(layer.Texture, rect, collection.SameItems(width*height, colr))
}

//=================================================================

func (l TileLayerId) TileAtCell(column, row int) Tile {
	var data = internal.TileLayers[uint8(l)]
	if data == nil {
		return Tile{}
	}

	var c = rl.GetImageColor(*data.Image, int32(column), int32(row))
	var packed = uint32(c.R)<<24 | uint32(c.G)<<16 | uint32(c.B)<<8 | uint32(c.A)

	return Tile{
		Id:          uint16(packed & 0xFFFF),
		FrameSpeed:  byte((packed >> 16) & 0x1F),
		FrameOffset: byte((packed >> 21) & 0x0F),
		FrameCount:  byte((packed >> 25) & 0x0F),
		Rotations90: byte((packed >> 29) & 0x03),
		Flip:        (packed >> 31) == 1,
	}
}

// func (l TileLayerId) Points() []float32 {
// 	var data = internal.TileLayers[uint8(l)]
// 	if data == nil {
// 		return nil
// 	}

// 	var hash = random.HashPrimitives(
// 		tm.AtlasId, tm.LayerId, tm.Color,
// 		tm.X, tm.Y, tm.Width, tm.Height, tm.Angle,
// 		tm.ScaleX, tm.ScaleY, tm.PivotX, tm.PivotY,
// 	)
// 	var isStructDirty = tm.hash != hash
// 	defer func() { tm.hash = hash }()

// 	if data.Image == nil || data.Texture.Width == 0 { // is object layer
// 		if !isStructDirty {
// 			return tm.allPointsCache
// 		}

// 		var copy = collection.Copy(data.ObjectPoints)
// 		for p := 0; p < len(copy); p += 2 {
// 			copy[p], copy[p+1] = tm.PointToGlobal(copy[p], copy[p+1])
// 		}
// 		tm.allPointsCache = copy
// 		return copy
// 	}

// 	var isTileDataDirty = tm.lastDirtyTime != data.LastDirtyTime
// 	if !isTileDataDirty && !isStructDirty {
// 		return tm.allPointsCache
// 	}

//		var w, h = tm.Size()
//		var result = make([]float32, 0, 32)
//		var afterFirst = false
//		for cellIndex1D := range data.CellsWithPoints {
//			var row, column = number.Index1DToIndexes2D(cellIndex1D, w, h)
//			result = append(result, tm.PointsAtCell(column, row)...)
//			if !afterFirst {
//				result = append(result, number.NaN(), number.NaN())
//				afterFirst = true
//			}
//		}
//		tm.allPointsCache = result
//		tm.lastDirtyTime = data.LastDirtyTime
//		return result
//	}
// func (l TileLayerId) PointsAtCell(column, row int) []float32 {
// 	var tile = tm.TileAtCell(column, row)
// 	if tile.Id == 0 {
// 		return nil
// 	}
// 	var ptsPerTile = tm.PointsFromTile(tile.Id)
// 	var result = make([]float32, 0, len(ptsPerTile))
// 	var tw, th = tm.SizeTile()
// 	for i := 1; i < len(ptsPerTile); i += 2 {
// 		var x, y = ptsPerTile[i-1], ptsPerTile[i]
// 		x, y = point.RotateAroundPoint(x, y, tw/2, th/2, float32(tile.Rotations90)*90)
// 		if tile.Flip {
// 			x = tw - x
// 		}
// 		x, y = tm.PointToGlobal(x+float32(column)*tw, y+float32(row)*th)
// 		result = append(result, x, y)
// 	}
// 	return result
// }

// The points are in tile space.
// func (l TileLayerId) PointsFromTile(tileId uint16) []float32 {
// 	var tileSet = internal.TileAtlases[uint8(tm.AtlasId)]
// 	if tileSet == nil {
// 		return nil
// 	}
// 	var pts, has = tileSet.PointsPerTile[tileId]
// 	if !has {
// 		return nil
// 	}
// 	return collection.Copy(pts)
// }

func (l TileLayerId) Size() (columns, rows int) {
	return 1, 1
}
func (l TileLayerId) SizeTile() (width, height float32) {
	var layer = internal.TileLayers[uint8(l)]
	var atlas = internal.TileAtlases[uint8(layer.AtlasId)]
	if atlas == nil {
		return number.NaN(), number.NaN()
	}
	return float32(atlas.TileWidth), float32(atlas.TileHeight)
}
func (l TileLayerId) SizeTileSet() (columns, rows int) {
	var layer = internal.TileLayers[uint8(l)]
	var atlas = internal.TileAtlases[uint8(layer.AtlasId)]
	if atlas == nil {
		return 0, 0
	}
	var tw, th = 1, 1
	return tw / atlas.TileWidth, th / atlas.TileHeight
}

func (l TileLayerId) TileCount() int {
	var w, h = l.SizeTileSet()
	return w * h
}
