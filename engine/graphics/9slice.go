package graphics

type NineSlice struct {
	Sprite
	Slices [9]Sprite
}

func NewNineSlice(assetIds [9]string) NineSlice {
	var result = NineSlice{Sprite: NewSprite("", 0, 0)}
	for i := range 9 {
		result.Slices[i] = NewSprite(assetIds[i], 0, 0)
	}
	return result
}
