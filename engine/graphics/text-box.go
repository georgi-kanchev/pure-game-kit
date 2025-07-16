package graphics

import (
	"pure-kit/engine/internal"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/utility/point"
	"pure-kit/engine/utility/text"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type TextBox struct {
	X, Y, Width, Height, Angle, GapSymbols, GapLines  float32
	PivotX, PivotY, Thickness, Smoothness, ValueScale float32
	FontId, Value                                     string
	Color                                             uint
}

func (textBox TextBox) Rectangle() (x, y, width, height float32) {
	x, y = textBox.X, textBox.Y
	x, y = point.MoveAt(x, y, textBox.Angle, -textBox.Width*textBox.PivotX)
	x, y = point.MoveAt(x, y, textBox.Angle+90, -textBox.Height*textBox.PivotY)
	return x, y, textBox.Width, textBox.Height
}

func NewTextBox(fontId string, x, y float32, value ...any) TextBox {
	return TextBox{
		FontId:     fontId,
		Value:      text.New(value...),
		X:          x,
		Y:          y,
		Width:      200,
		Height:     100,
		ValueScale: 1,
		Color:      color.White,
		Thickness:  0.5,
		GapSymbols: 0.2}
}

// #region private

func (textBox TextBox) height() float32 {
	return textBox.ValueScale * float32(textBox.font().BaseSize)
}
func (textBox TextBox) font() *rl.Font {
	var font, hasFont = internal.Fonts[textBox.FontId]
	if !hasFont {
		var defaultFont = rl.GetFontDefault()
		font = &defaultFont
	}
	return font
}
func (textBox TextBox) pivot() rl.Vector2 {
	var _, _, w, h = textBox.Rectangle()
	return rl.Vector2{X: textBox.PivotX * w, Y: textBox.PivotY * h}
}
func (textBox TextBox) gapSymbols() float32 {
	return textBox.GapSymbols * textBox.height() / 5
}

// #endregion
