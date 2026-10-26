package filter

type Enum uint8

const (
	Point Enum = iota
	Bilinear
	Trilinear
	Anisotropic4x
	Anisotropic8x
	Anisotropic16x
)
