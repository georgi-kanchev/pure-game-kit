package graphics

import "pure-game-kit/internal"

type NinePatch struct {
	Quad

	BoxId string
	EdgeLeft, EdgeRight,
	EdgeTop, EdgeBottom float32
}

func NewNinePatch(boxId string, x, y float32) *NinePatch {
	var result = &NinePatch{Quad: *NewQuad(x, y), BoxId: boxId}
	var slices, has = internal.Boxes[boxId]

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
