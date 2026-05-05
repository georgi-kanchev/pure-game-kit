package graphics

type NinePatch struct {
	Quad
	BoxId     string
	EdgeScale float32
}

func NewNinePatch(boxId string, x, y float32) *NinePatch {
	var result = &NinePatch{Quad: *NewQuad(x, y), BoxId: boxId, EdgeScale: 1}
	return result
}
