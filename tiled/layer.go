package tiled

type Layer struct {
	Properties map[string]any
	TileIds    []uint32
	Objects    []*Object
	TextureId  string
}
