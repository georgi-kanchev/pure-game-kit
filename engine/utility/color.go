package color

type Color struct {
	R, G, B, A byte
}

func NewRGB(r, g, b byte) Color {
	return Color{r, g, b, 255}
}
func NewRGBA(r, g, b, a byte) Color {
	return Color{r, g, b, a}
}
