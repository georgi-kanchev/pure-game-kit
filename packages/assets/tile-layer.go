package assets

import (
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/collection"
	col "pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/number"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// tile data in bits (32)
// bits 31..31 = flip					(0 or 1)
// bits 30..29 = rotations 				(4 values: 0, 90, 180, 270)
// bits 28..25 = animation frame count 	(0 to 15)
// bits 24..21 = animation offset		(0 to 15)
// bits 20..16 = animation frames/s		(0 to 31)
// bits 15..00 = tile id				(0 to 65535)

type Tile struct {
	Id          uint16
	Rotations90 byte // 90 degree turns, ranged 0..3 (possible values: 0, 90, 180, 270)
	Flip        bool

	FrameCount  byte // Ranged 0..15 (sequential tile count in the atlas)
	FrameOffset byte // Ranged 0..15
	FrameSpeed  byte // Ranged 0..31
}
type TileLayerId uint8

func NewTile(id uint16) Tile {
	return Tile{Id: id}
}
func NewTileOriented(id uint16, rotations90 byte, flip bool) Tile {
	return Tile{Id: id, Rotations90: rotations90, Flip: flip}
}
func NewTileAnimated(id uint16, frameCount, frameOffset, frameSpeed byte) Tile {
	return Tile{Id: id, FrameCount: frameCount, FrameSpeed: frameSpeed, FrameOffset: frameOffset}
}

func LoadTileLayer(columns, rows, tileSize int, imageId ImageId) TileLayerId {
	columns, rows = number.Limit(columns, 1, 2048), number.Limit(rows, 1, 2048)

	var id = len(internal.TileLayers) + 1
	var data = &internal.TileLayer{Columns: columns, Rows: rows,
		Image: rl.GenImageColor(columns, rows, rl.Blank), CellsWithPoints: make(map[int]struct{}),
		TileSize: tileSize, ImageId: int32(imageId), ShapesPerTile: make(map[uint16][][6]float32)}
	var tex = rl.LoadTextureFromImage(data.Image)
	rl.SetTextureFilter(tex, rl.FilterPoint)
	data.Texture = tex
	internal.TileLayers[uint8(id)] = data
	return TileLayerId(id)
}

//=================================================================

func (l TileLayerId) SetTile(column, row int, tile Tile) {
	l.SetTileArea(column, row, 1, 1, tile)
}
func (l TileLayerId) SetTileArea(column, row, width, height int, tile Tile) {
	var layer = internal.TileLayers[uint8(l)]
	if layer == nil {
		return
	}

	var packed = newTilePacked(tile)
	var r, g = uint8((packed >> 24) & 0xFF), uint8((packed >> 16) & 0xFF)
	var b, a = uint8((packed >> 8) & 0xFF), uint8((packed >> 0) & 0xFF)
	var colr = rl.NewColor(r, g, b, a)
	var rect = rl.NewRectangle(float32(column), float32(row), float32(width), float32(height))
	var columns, rows = layer.Columns, layer.Rows
	var _, cellHasPts = layer.ShapesPerTile[tile.Id]

	for i := row; i < row+height; i++ {
		for j := column; j < column+width; j++ {
			var prevTile = l.TileAtCell(j, i)
			var _, prevCellHasPts = layer.ShapesPerTile[prevTile.Id]
			if !prevCellHasPts && !cellHasPts {
				continue
			}

			var index1D = number.Indexes2DToIndex1D(j, i, columns, rows)
			if cellHasPts {
				layer.CellsWithPoints[index1D] = struct{}{}
			} else {
				delete(layer.CellsWithPoints, index1D)
			}
		}
	}

	rl.ImageDrawRectangle(layer.Image, int32(column), int32(row), int32(width), int32(height), colr)
	rl.UpdateTextureRec(layer.Texture, rect, *collection.NewListOfItem(width*height, colr).ToSlice())
}
func (l TileLayerId) SetAtlasId(atlasId ImageId) {
	var layer = internal.TileLayers[uint8(l)]
	if layer == nil {
		return
	}
	layer.ImageId = int32(atlasId)
}

//=================================================================

func (l TileLayerId) TileAtCell(column, row int) Tile {
	var layer = internal.TileLayers[uint8(l)]
	if layer == nil {
		return Tile{}
	}

	var c = rl.GetImageColor(*layer.Image, int32(column), int32(row))
	col.RGBA(c.R, c.G, c.B, c.A)
	var packed = uint32(c.R)<<24 | uint32(c.G)<<16 | uint32(c.B)<<8 | uint32(c.A)
	return newTileUnpacked(packed)
}
func (l TileLayerId) Size() (columns, rows int) {
	var layer = internal.TileLayers[uint8(l)]
	if layer == nil {
		return 0, 0
	}
	return layer.Columns, layer.Rows
}
func (l TileLayerId) TileSize() (width, height float32) {
	var layer = internal.TileLayers[uint8(l)]
	if layer == nil {
		return number.NaN(), number.NaN()
	}
	return float32(layer.TileSize), float32(layer.TileSize)
}
func (l TileLayerId) AtlasSize() (columns, rows int) {
	var layer = internal.TileLayers[uint8(l)]
	if layer == nil {
		return 0, 0
	}
	var tex = internal.Images[layer.ImageId]
	return int(tex.CropWidth) / layer.TileSize, int(tex.CropHeight) / layer.TileSize
}
func (l TileLayerId) TileCount() int {
	var w, h = l.AtlasSize()
	return w * h
}
func (l TileLayerId) AtlasId() ImageId {
	var layer = internal.TileLayers[uint8(l)]
	if layer == nil {
		return 0
	}
	return ImageId(layer.ImageId)
}

// private ========================================================

func newTileUnpacked(packed uint32) Tile {
	return Tile{Id: uint16(packed & 0xFFFF), Rotations90: byte((packed >> 29) & 0x03), Flip: (packed >> 31) == 1,
		FrameSpeed: byte((packed >> 16) & 0x1F), FrameOffset: byte((packed >> 21) & 0x0F), FrameCount: byte((packed >> 25) & 0x0F)}
}
func newTilePacked(tile Tile) uint32 {
	var flipBit uint32
	if tile.Flip {
		flipBit = 1 << 31
	}
	return uint32(tile.Id&0xFFFF) | uint32(tile.FrameSpeed&0x1F)<<16 | uint32(tile.FrameOffset&0x0F)<<21 |
		uint32(tile.FrameCount&0x0F)<<25 | uint32(tile.Rotations90&0x03)<<29 | flipBit
}
