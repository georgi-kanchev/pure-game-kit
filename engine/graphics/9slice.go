package graphics

import "pure-kit/engine/internal"

type NineSlice struct {
	Sprite
	// [0]upperLeft [1]upper [2]upperRight [3]left [4]right [5]lowerLeft [6]lower [7]lowerRight
	SliceIds   [8]string
	SliceSizes [4]float32 // [0]upper [1]left [2]right [3]lower
	// [0]upperLeft [1]upper [2]upperRight [3]left [4]right [5]lowerLeft [6]lower [7]lowerRight
	SliceFlipX [8]bool // [0]upper [1]left [2]right [3]lower
	// [0]upperLeft [1]upper [2]upperRight [3]left [4]right [5]lowerLeft [6]lower [7]lowerRight
	SliceFlipY [8]bool // [0]upper [1]left [2]right [3]lower
}

func NewNineSlice(assetId string, x, y float32, assetIds [8]string) NineSlice {
	var result = NineSlice{Sprite: NewSprite("", 0, 0)}
	var _, uh = internal.AssetSize(assetIds[1])
	var lw, _ = internal.AssetSize(assetIds[3])
	var rw, _ = internal.AssetSize(assetIds[4])
	var _, dh = internal.AssetSize(assetIds[6])

	result.AssetId = assetId
	result.X, result.Y = x, y
	result.SliceIds = assetIds
	result.SliceSizes[0] = float32(uh)
	result.SliceSizes[1] = float32(lw)
	result.SliceSizes[2] = float32(rw)
	result.SliceSizes[3] = float32(dh)
	return result
}
