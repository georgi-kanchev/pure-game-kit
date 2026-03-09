package graphics

import (
	"pure-game-kit/internal"
	"pure-game-kit/utility/collection"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type TileMap struct {
	Node
	TileSetId, TileDataId string
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

func NewTileMap(tileSetId, tileDataId string) *TileMap {
	var tileMap = &TileMap{Node: *NewNode(0, 0), TileSetId: tileSetId, TileDataId: tileDataId}
	var atlas = internal.TileSets[tileSetId]
	var data = internal.TileDatas[tileDataId]
	if atlas != nil && data != nil {
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
	var data = internal.TileDatas[tm.TileDataId]
	if data == nil {
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

	rl.ImageDrawRectangle(data.Image, int32(column), int32(row), int32(width), int32(height), colr)
	rl.UpdateTextureRec(*data.Texture, rect, collection.SameItems(width*height, colr))
}

//=================================================================

func (tm *TileMap) TileAt(column, row int) *Tile {
	var data = internal.TileDatas[tm.TileDataId]
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
