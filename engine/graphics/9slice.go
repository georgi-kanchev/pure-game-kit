package graphics

import "pure-kit/engine/internal"

type NineSlice struct {
	Sprite
	EdgeLeft, EdgeRight, EdgeTop, EdgeBottom float32
}

func NewNineSlice(assetId string, x, y float32) NineSlice {
	var result = NineSlice{Sprite: NewSprite(assetId, x, y)}
	var slices, has = internal.NineSlices[assetId]

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
