package graphics

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/internal"
	"pure-game-kit/utility/number"
	txt "pure-game-kit/utility/text"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type textBoxCache struct {
	text                           string
	fontId                         string
	tint                           uint
	wordWrap                       bool
	width, height                  float32
	alignmentX, alignmentY         float32
	lineHeight, symbolGap, lineGap float32
}

type TextBox struct {
	Quad
	Text, FontId string
	WordWrap     bool
	AlignmentX, AlignmentY,
	LineHeight, SymbolGap, LineGap,
	ShadowOffsetX, ShadowOffsetY float32

	cache        textBoxCache
	cacheChars   []string
	cacheSymbols []symbol
	cacheWrap    string
}

func NewTextBox(fontId string, x, y float32, text ...any) *TextBox {
	var quad = NewQuad(x, y)
	var textBox = &TextBox{
		FontId: fontId, Quad: *quad, Text: txt.New(text...), LineHeight: 100, SymbolGap: 0.2, WordWrap: true,
	}
	var font = textBox.font()
	var measure = rl.MeasureTextEx(font, textBox.Text, textBox.LineHeight, textBox.gapSymbols())
	textBox.Width, textBox.Height = measure.X, measure.Y
	return textBox
}

//=================================================================

// Does not wrap the text - use TextWrap(...) beforehand if intended.
func (t *TextBox) TextMeasure(text string) (width, height float32) {
	return t.measure(t.font(), text)
}
func (t *TextBox) TextWrap(text string) string {
	var state = textBoxCache{
		t.Text, t.FontId, t.Tint, t.WordWrap,
		t.Width, t.Height,
		t.AlignmentX, t.AlignmentY,
		t.LineHeight, t.SymbolGap, t.LineGap,
	}
	if t.cache == state {
		return t.cacheWrap
	}

	var replaced, originals = internal.ReplaceStrings(text, '{', '}', internal.Placeholder)
	var words = txt.Split(replaced, " ")
	var curX, curY float32 = 0, 0
	var buffer = txt.NewBuilder()
	var tagIndex = 0
	var ph = string(internal.Placeholder)
	var gapY = t.gapLines()

	for w := range words {
		var word = words[w]

		if w < len(words)-1 {
			word += " " // split removes spaces, add it for all words but last one
		}

		var trimWord = txt.Remove(txt.Trim(word), ph)
		var wordSize, _ = t.TextMeasure(trimWord)

		if txt.ContainsAll(trimWord, string(placeholderCharAsset)) {
			wordSize += t.LineHeight
		}

		var wordEndOfBox = curX+wordSize > t.Width+1
		var wordFirst = w == 0
		var wordNewLine = !wordFirst && t.WordWrap && wordEndOfBox

		if wordNewLine {
			curX = 0
			curY += t.LineHeight + gapY

			buffer.WriteSymbol('\n')
		}

		for i, c := range word {
			var char = string(c)
			var charSize, _ = t.TextMeasure(char)
			charSize = condition.If(c == internal.Placeholder, 0, charSize)
			charSize = condition.If(c == placeholderCharAsset, t.LineHeight, charSize)
			var charEndOfBoxX = charSize > 0 && curX+charSize > t.Width+1
			var charFirst = i == 0 && wordFirst
			var charNewLine = !charFirst && char != " " && (char == "\n" || charEndOfBoxX)

			if charEndOfBoxX { // outside right
				continue // rare cases but happens with single symbol & small width
			}

			if charNewLine {
				curX = 0
				curY += t.LineHeight + gapY

				if char != "\n" {
					buffer.WriteSymbol('\n')
				}
			}

			if c == internal.Placeholder {
				char = "{" + originals[tagIndex] + "}"
				tagIndex++
			}
			buffer.WriteText(char)
			curX += condition.If(charSize > 0, charSize+t.gapSymbols(), 0)
		}
	}
	var result = buffer.ToText()
	result = txt.Replace(result, " \n", "\n")
	t.cache = state
	t.cacheWrap = result
	return result
}
func (t *TextBox) TextLines() []string {
	var lines, _ = t.formatSymbols()
	return lines
}
func (t *TextBox) TextSymbol(symbolIndex int) (x, y, width, height, angle float32) {
	var _, symbols = t.formatSymbols()
	if symbolIndex < 0 || symbolIndex >= len(symbols) {
		return number.NaN(), number.NaN(), number.NaN(), number.NaN(), number.NaN()
	}

	var s = symbols[symbolIndex]
	var gx, gy = t.PointToGlobal(s.Bounds.X, s.Bounds.Y)
	return gx, gy, s.Bounds.Width * t.ScaleX, t.LineHeight * t.ScaleY, t.Angle
}
