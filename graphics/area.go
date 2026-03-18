package graphics

type Area struct{ X, Y, Width, Height float32 }

func NewArea(x, y, width, height float32) *Area {
	return &Area{X: x, Y: y, Width: width, Height: height}
}
