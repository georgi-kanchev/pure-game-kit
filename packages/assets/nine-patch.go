package assets

import "pure-game-kit/packages/internal"

type NinePatchId uint8

func LoadNinePatch(imageId ImageId, left, right, top, bottom float32, tileLeft, tileRight, tileTop, tileBottom, tileCenter bool) NinePatchId {
	if imageId == 0 {
		return 0
	}

	internal.NinePatchNextId++
	var id = internal.NinePatchNextId
	internal.NinePatches[id] = internal.NinePatch{ImageId: int32(imageId), Left: left, Right: right, Top: top, Bottom: bottom,
		TileLeft: tileLeft, TileRight: tileRight, TileTop: tileTop, TileBottom: tileBottom, TileCenter: tileCenter}
	return NinePatchId(id)
}

func (n NinePatchId) Unload() {
	if n == 0 {
		return
	}
	delete(internal.NinePatches, uint8(n))
}
