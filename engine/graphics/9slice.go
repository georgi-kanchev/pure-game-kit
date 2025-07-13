package graphics

type NineSlice struct {
	Node
	Slices [9]Node
}

func NewNineSlice(assetIds [9]string) NineSlice {
	var result = NineSlice{Node: NewNode("")}
	for i := range 9 {
		result.Slices[i] = NewNode(assetIds[i])
	}
	return result
}
