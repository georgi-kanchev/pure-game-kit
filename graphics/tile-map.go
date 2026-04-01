package graphics

import (
	"pure-game-kit/internal"
	"pure-game-kit/utility/collection"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/random"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type TileMap struct {
	Quad
	TileSetId, TileLayerId string

	hash           uint64
	cacheAllPoints []float32
	lastDirtyTime  float32
}

type Tile struct {
	// data in bits (32)
	// bits 31..31 = flip					(0 or 1)
	// bits 30..29 = rotations 				(4 values: 0, 90, 180, 270)
	// bits 28..25 = animation frame count 	(0 to 15)
	// bits 24..21 = animation offset		(0 to 15)
	// bits 20..16 = animation frames/s		(0 to 31)
	// bits 15..00 = tile id				(0 to 65535)

	Id       uint16
	Rotation byte // 90 degree turns, ranged 0..3 (possible values: 0, 90, 180, 270)
	Flip     bool

	FrameCount  byte // Ranged 0..15 (sequential tile count in the atlas)
	FrameOffset byte // Ranged 0..15
	FrameSpeed  byte // Ranged 0..31
}

func NewTileMap(tileSetId, tileLayerId string) *TileMap {
	var tileMap = &TileMap{Quad: *NewQuad(0, 0), TileSetId: tileSetId, TileLayerId: tileLayerId}
	var atlas = internal.TileSets[tileSetId]
	var data = internal.TileLayers[tileLayerId]
	if atlas != nil && data != nil && data.Image != nil {
		tileMap.Width = float32(data.Image.Width * int32(atlas.TileWidth))
		tileMap.Height = float32(data.Image.Height * int32(atlas.TileHeight))
	}

	return tileMap
}

func NewTile(id uint16) *Tile {
	return &Tile{Id: id}
}
func NewTileOriented(id uint16, rotation byte, flip bool) *Tile {
	return &Tile{Id: id, Rotation: rotation, Flip: flip}
}
func NewTileAnimated(id uint16, frameCount, frameOffset, frameSpeed byte) *Tile {
	return &Tile{Id: id, FrameCount: frameCount, FrameSpeed: frameSpeed, FrameOffset: frameOffset}
}

//=================================================================

func (tm *TileMap) SetTile(column, row int, tile *Tile) {
	tm.SetTileArea(column, row, 1, 1, tile)
}
func (tm *TileMap) SetTileArea(column, row, width, height int, tile *Tile) {
	var data = internal.TileLayers[tm.TileLayerId]
	var tileSet = internal.TileSets[tm.TileSetId]
	if data == nil || tileSet == nil {
		return
	}

	var packed = uint32(tile.Id&0xFFFF) |
		uint32(tile.FrameSpeed&0x1F)<<16 |
		uint32(tile.FrameOffset&0x0F)<<21 |
		uint32(tile.FrameCount&0x0F)<<25 |
		uint32(tile.Rotation&0x03)<<29

	if tile.Flip {
		packed |= (1 << 31)
	}

	var r = uint8((packed >> 24) & 0xFF)
	var g = uint8((packed >> 16) & 0xFF)
	var b = uint8((packed >> 8) & 0xFF)
	var a = uint8((packed >> 0) & 0xFF)
	var colr = rl.NewColor(r, g, b, a)
	var rect = rl.NewRectangle(float32(column), float32(row), float32(width), float32(height))
	var w, h = tm.Size()
	var _, cellHasPts = tileSet.PointsPerTile[tile.Id]

	for i := row; i < row+height; i++ {
		for j := column; j < column+width; j++ {
			var prevTile = tm.TileAtCell(j, i)
			var _, prevCellHasPts = tileSet.PointsPerTile[prevTile.Id]
			if !prevCellHasPts && !cellHasPts {
				continue
			}

			var index1D = number.Indexes2DToIndex1D(j, i, w, h)
			data.LastDirtyTime = internal.Runtime

			if cellHasPts {
				data.CellsWithPoints[index1D] = struct{}{}
			} else {
				delete(data.CellsWithPoints, index1D)
			}
		}
	}

	rl.ImageDrawRectangle(data.Image, int32(column), int32(row), int32(width), int32(height), colr)
	rl.UpdateTextureRec(*data.Texture, rect, collection.SameItems(width*height, colr))
}

//=================================================================

func (tm *TileMap) TileAtCell(column, row int) *Tile {
	var data = internal.TileLayers[tm.TileLayerId]
	if data == nil {
		return nil
	}

	var c = rl.GetImageColor(*data.Image, int32(column), int32(row))
	var packed = uint32(c.R)<<24 | uint32(c.G)<<16 | uint32(c.B)<<8 | uint32(c.A)

	return &Tile{
		Id:          uint16(packed & 0xFFFF),
		FrameSpeed:  byte((packed >> 16) & 0x1F),
		FrameOffset: byte((packed >> 21) & 0x0F),
		FrameCount:  byte((packed >> 25) & 0x0F),
		Rotation:    byte((packed >> 29) & 0x03),
		Flip:        (packed >> 31) == 1,
	}
}
func (tm *TileMap) CellAtPoint(x, y float32) (column, row int) {
	var localX, localY = tm.PointToLocal(x, y)
	var tw, th = tm.SizeTile()
	return int(localX / tw), int(localY / th)
}

func (tm *TileMap) Points() []float32 {
	var data = internal.TileLayers[tm.TileLayerId]
	if data == nil {
		return nil
	}

	var hash = random.Hash(tm)
	var isStructDirty = tm.hash != hash
	defer func() { tm.hash = hash }()

	if data.Image == nil || data.Texture == nil { // is object layer
		if !isStructDirty {
			return tm.cacheAllPoints
		}

		var copy = collection.Copy(data.ObjectPoints)
		for p := 0; p < len(copy); p += 2 {
			copy[p], copy[p+1] = tm.PointToGlobal(copy[p], copy[p+1])
		}
		tm.cacheAllPoints = copy
		return copy
	}

	var isTileDataDirty = tm.lastDirtyTime != data.LastDirtyTime
	if !isTileDataDirty && !isStructDirty {
		return tm.cacheAllPoints
	}

	var w, h = tm.Size()
	var result = make([]float32, 0, 32)
	for cellIndex1D := range data.CellsWithPoints {
		var row, column = number.Index1DToIndexes2D(cellIndex1D, w, h)
		result = append(result, tm.PointsAtCell(column, row)...)
	}
	tm.cacheAllPoints = result
	tm.lastDirtyTime = data.LastDirtyTime
	return result
}
func (tm *TileMap) PointsAtCell(column, row int) []float32 {
	var tile = tm.TileAtCell(column, row)
	if tile == nil {
		return nil
	}
	var ptsPerTile = tm.PointsFromTile(tile.Id)
	var result = make([]float32, 0, len(ptsPerTile))
	var tw, th = tm.SizeTile()
	for i := 1; i < len(ptsPerTile); i += 2 {
		var x, y = tm.PointToGlobal(ptsPerTile[i-1]+float32(column)*tw, ptsPerTile[i]+float32(row)*th)
		result = append(result, x, y)
	}
	return result
}
func (tm *TileMap) PointsFromTile(tileId uint16) []float32 {
	var tileSet = internal.TileSets[tm.TileSetId]
	if tileSet == nil {
		return nil
	}
	return collection.Copy(tileSet.PointsPerTile[tileId])
}

func (tm *TileMap) Size() (columns, rows int) {
	return internal.AssetSize(tm.TileLayerId)
}
func (tm *TileMap) SizeTile() (width, height float32) {
	var tileSet = internal.TileSets[tm.TileSetId]
	if tileSet == nil {
		return number.NaN(), number.NaN()
	}
	return float32(tileSet.TileWidth), float32(tileSet.TileHeight)
}
func (tm *TileMap) SizeTileSet() (columns, rows int) {
	var tileSet = internal.TileSets[tm.TileSetId]
	if tileSet == nil {
		return 0, 0
	}
	var tw, th = internal.AssetSize(tm.TileSetId)
	return tw / tileSet.TileWidth, th / tileSet.TileHeight
}
