package graphics

import (
	"pure-kit/engine/internal"
	"pure-kit/engine/utility/text"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type TextBox struct {
	Node
	GapSymbols, GapLines              float32
	Thickness, Smoothness, LineHeight float32
	Value                             string
}

func NewTextBox(fontId string, x, y float32, value ...any) TextBox {
	var node = NewNode(fontId, x, y)
	var textBox = TextBox{Node: node, Value: text.New(value...), LineHeight: 100, Thickness: 0.5, GapSymbols: 0.2}
	var font = textBox.font()
	var measure = rl.MeasureTextEx(*font, textBox.Value, textBox.LineHeight, textBox.gapSymbols())
	textBox.Width, textBox.Height = measure.X, measure.Y
	return textBox
}

func (node *TextBox) Size() (width, height float32) {
	return node.Width, node.Height
}

// #region private

func (node *TextBox) font() *rl.Font {
	var font, hasFont = internal.Fonts[node.AssetId]
	if !hasFont {
		var defaultFont = rl.GetFontDefault()
		font = &defaultFont
	}
	return font
}
func (node *TextBox) gapSymbols() float32 {
	return node.GapSymbols * node.LineHeight / 5
}

// #endregion
