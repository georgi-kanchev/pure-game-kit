package graphics

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
