package internal

type NinePatch struct {
	ImageId int32

	Left, Right, Top, Bottom float32 // edge sizes in pixels

	TileLeft, TileRight, TileTop, TileBottom, TileCenter bool
}

var NinePatches = make(map[uint8]NinePatch)
var NinePatchNextId uint8
