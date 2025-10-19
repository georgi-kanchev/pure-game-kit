package graphics

import "pure-game-kit/internal"

type Box struct {
	Sprite
	EdgeLeft, EdgeRight, EdgeTop, EdgeBottom float32
}

func NewBox(assetId string, x, y float32) Box {
	var result = Box{Sprite: NewSprite(assetId, x, y)}
	var slices, has = internal.Boxes[assetId]

	if !has {
		slices = [9]string{}
	}

	var _, uh = internal.AssetSize(slices[1])
	var lw, _ = internal.AssetSize(slices[3])
	var rw, _ = internal.AssetSize(slices[4])
	var _, dh = internal.AssetSize(slices[6])

	result.EdgeLeft = float32(lw)
	result.EdgeRight = float32(rw)
	result.EdgeBottom = float32(uh)
	result.EdgeTop = float32(dh)
	return result
}
