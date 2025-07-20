package graphics

import (
	"bytes"
	"pure-kit/engine/internal"
	"pure-kit/engine/utility/text"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type TextBox struct {
	Node
	AlignmentX, AlignmentY, Thickness, Smoothness,
	SymbolGap, LineHeight, LineGap float32
	Value            string
	WordWrap         bool
	EmbeddedColors   []uint
	EmbeddedAssetIds []string
}

func NewTextBox(fontId string, x, y float32, value ...any) TextBox {
	var node = NewNode(fontId, x, y)
	var textBox = TextBox{
		Node: node, Value: text.New(value...), LineHeight: 100, Thickness: 0.5, SymbolGap: 0.2, WordWrap: true}
	var font = textBox.font()
	var measure = rl.MeasureTextEx(*font, textBox.Value, textBox.LineHeight, textBox.gapSymbols())
	textBox.Width, textBox.Height = measure.X, measure.Y
	return textBox
}

func (textBox *TextBox) Size() (width, height float32) {
	return textBox.Width, textBox.Height
}

func (textBox *TextBox) WrapValue(value string) string {
	var font = textBox.font()
	var words = strings.Split(value, " ")
	var curX, curY = textBox.X, textBox.Y
	var gapSymbols = textBox.gapSymbols()
	var buffer = bytes.NewBufferString("")

	for w := range words {
		var word = words[w]
		if w < len(words)-1 {
			word += " " // split removes spaces, add it for all words but last one
		}

		var wordLength = text.Length(word)
		var wordSize = rl.MeasureTextEx(*font, strings.Trim(word, " "), textBox.LineHeight, gapSymbols)
		var wordEndOfBox = curX+wordSize.X > textBox.Width
		var firstWord = w == 0

		if !firstWord && textBox.WordWrap && wordEndOfBox {
			curX = 0
			curY += textBox.LineHeight + textBox.gapLines()
			buffer.WriteRune('\n')
		}

		for c := range wordLength {
			var char = rune(word[c])
			var charSize = rl.MeasureTextEx(*font, string(char), textBox.LineHeight, 0)
			var charEndOfBoxX = curX+charSize.X > textBox.Width
			var charFirst = c == 0 && firstWord

			if !charFirst && char != ' ' && (char == '\n' || charEndOfBoxX) {
				curX = 0
				curY += textBox.LineHeight + textBox.gapLines()

				if char != '\n' {
					buffer.WriteRune('\n')
				}
			}

			buffer.WriteRune(char)
			curX += charSize.X + gapSymbols
		}
	}

	var result = buffer.String()
	result = strings.ReplaceAll(result, " \n", "\n")

	return result
}

// #region private

func (textBox *TextBox) font() *rl.Font {
	var font, hasFont = internal.Fonts[textBox.AssetId]
	if !hasFont {
		var defaultFont = rl.GetFontDefault()
		font = &defaultFont
	}
	return font
}
func (textBox *TextBox) gapSymbols() float32 {
	return textBox.SymbolGap * textBox.LineHeight / 5
}
func (textBox *TextBox) gapLines() float32 {
	return textBox.LineGap * textBox.LineHeight / 5
}

// #endregion
