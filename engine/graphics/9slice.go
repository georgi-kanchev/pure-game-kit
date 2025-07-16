package graphics

type NineSlice struct {
	Sprite
	Slices [9]Sprite
}

func NewNineSlice(assetIds [9]string) NineSlice {
	var result = NineSlice{Sprite: NewSprite("")}
	for i := range 9 {
		result.Slices[i] = NewSprite(assetIds[i])
	}
	return result
}
