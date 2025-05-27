package utility

type Color struct {
	R, G, B, A byte
}

func (color Color) ToDark() Color {
	return Color{0, 0, 0, 0}
}
func (color Color) ToDarkProgress(unit float32) Color {
	return Color{0, 0, 0, 0}
}

// public Color ToDark(float unit = 0.5f)
// {
//     var red = (byte)Map(unit, 0, 1, R, 0);
//     var green = (byte)Map(unit, 0, 1, G, 0);
//     var blue = (byte)Map(unit, 0, 1, B, 0);
//     return new(red, green, blue);
// }

func NewRGB(r, g, b byte) Color {
	return Color{r, g, b, 255}
}
func NewRGBA(r, g, b, a byte) Color {
	return Color{r, g, b, a}
}

// public Color ToDark(float unit = 0.5f)
// {
//     var red = (byte)Map(unit, 0, 1, R, 0);
//     var green = (byte)Map(unit, 0, 1, G, 0);
//     var blue = (byte)Map(unit, 0, 1, B, 0);
//     return new(red, green, blue);
// }
// public Color ToColor(Color color, float unit = 0.5f)
// {
//     var red = (byte)Map(unit, 0, 1, R, color.r);
//     var green = (byte)Map(unit, 0, 1, G, color.g);
//     var blue = (byte)Map(unit, 0, 1, B, color.b);
//     var alpha = (byte)Map(unit, 0, 1, A, color.a);
//     return new(red, green, blue, alpha);
// }
// public Color ToBright(float unit = 0.5f)
// {
//     var red = (byte)Map(unit, 0, 1, R, 255);
//     var green = (byte)Map(unit, 0, 1, G, 255);
//     var blue = (byte)Map(unit, 0, 1, B, 255);
//     return new(red, green, blue);
// }
// public Color ToTransparent(float unit = 0.5f)
// {
//     return new(R, G, B, (byte)Map(unit, 0, 1, A, 0));
// }
// public Color ToOpaque(float unit = 0.5f)
// {
//     return new(R, G, B, (byte)Map(unit, 0, 1, A, 255));
// }
// public Color ToOpposite()
// {
//     return new((byte)(255 - R), (byte)(255 - G), (byte)(255 - B));
// }
