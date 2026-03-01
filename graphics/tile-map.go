package graphics

import (
	"pure-game-kit/internal"
	"pure-game-kit/utility/collection"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type TileMap struct {
	Node
	TileAtlasId, TileDataId string

	Effects *Effects
}

type Tile struct {
	// data in bits (32)
	// bits 31..31 = flip					(0 or 1)
	// bits 30..29 = rotations 				(4 values: 0, 90, 180, 270)
	// bits 28..26 = animation frame count 	(0 to 7)
	// bits 25..23 = animation frames/s		(0 to 7)
	// bits 22..20 = animation offset		(0 to 7)
	// bits 19..16 = custom flags			(0 to 16 or 4 flags, for gameplay)
	// bits 15..00 = tile id				(0 to 65535)

	Id       uint16
	Rotation byte // 90 degree turns, ranged 0..3 (possible values: 0, 90, 180, 270)
	Flip     bool

	AnimationFrameCount byte // Ranged 0..7 (sequential tile id count in the atlas)
	AnimationSpeed      byte // Ranged 0..7
	AnimationOffset     byte // Ranged 0..7

	CustomFlags [4]bool // Useful for holding per-tile gameplay data.
}

func NewTileMap(tileAtlasId, tileDataId string) *TileMap {
	var tileMap = &TileMap{Node: *NewNode(0, 0), TileAtlasId: tileAtlasId, TileDataId: tileDataId}
	var atlas = internal.TileAtlases[tileAtlasId]
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
func NewTileAnimated(id uint16, frameCount, speed, offset byte) *Tile {
	return &Tile{Id: id, AnimationFrameCount: frameCount, AnimationSpeed: speed, AnimationOffset: offset}
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

	var packed = uint32(tile.Id) |
		uint32(tile.Rotation%4)<<29 |
		uint32(tile.AnimationFrameCount&0x07)<<26 |
		uint32(tile.AnimationSpeed&0x07)<<23 |
		uint32(tile.AnimationOffset&0x07)<<20

	if tile.Flip {
		packed |= 0x80000000
	}

	for i := 0; i < 4; i++ {
		if tile.CustomFlags[i] {
			packed |= (1 << (16 + i))
		}
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

	var tile = &Tile{
		Id:                  uint16(packed & 0xFFFF),
		Rotation:            byte((packed >> 29) & 0x03),
		Flip:                (packed & 0x80000000) != 0,
		AnimationFrameCount: byte((packed >> 26) & 0x07),
		AnimationSpeed:      byte((packed >> 23) & 0x07),
		AnimationOffset:     byte((packed >> 20) & 0x07),
	}

	for i := range 4 {
		tile.CustomFlags[i] = (packed & (1 << (16 + i))) != 0
	}

	return tile
}
