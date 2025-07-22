package graphics

import "pure-kit/engine/internal"

type NineSlice struct {
	Sprite
	// [0]upperLeft [1]upper [2]upperRight [3]right [4]lowerRight [5]lower [6] lowerLeft [7]left
	SliceIds   [8]string
	SliceSizes [4]float32 // [0]upper [1]right [2]lower [3]left
	// [0]upperLeft [1]upper [2]upperRight [3]right [4]lowerRight [5]lower [6] lowerLeft [7]left
	SliceFlipX [8]bool // [0]upper [1]right [2]lower [3]left
	// [0]upperLeft [1]upper [2]upperRight [3]right [4]lowerRight [5]lower [6] lowerLeft [7]left
	SliceFlipY [8]bool // [0]upper [1]right [2]lower [3]left
}

func NewNineSlice(assetId string, x, y float32, assetIds [8]string) NineSlice {
	var result = NineSlice{Sprite: NewSprite("", 0, 0)}
	var _, uh = internal.AssetSize(assetIds[1])
	var rw, _ = internal.AssetSize(assetIds[3])
	var _, dh = internal.AssetSize(assetIds[5])
	var lw, _ = internal.AssetSize(assetIds[7])

	result.AssetId = assetId
	result.X, result.Y = x, y
	result.SliceIds = assetIds
	result.SliceSizes[0] = float32(uh)
	result.SliceSizes[1] = float32(rw)
	result.SliceSizes[2] = float32(dh)
	result.SliceSizes[3] = float32(lw)
	return result
}
