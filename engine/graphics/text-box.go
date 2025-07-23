package graphics

import (
	"bytes"
	"pure-kit/engine/internal"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/symbols"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type TextBox struct {
	Node
	Text     string
	WordWrap bool
	AlignmentX, AlignmentY,
	Thickness, Smoothness,
	SymbolGap, LineHeight, LineGap float32
	EmbeddedColorsTag, EmbeddedAssetsTag,
	EmbeddedThicknessesTag rune
	EmbeddedAssetIds    []string
	EmbeddedColors      []uint
	EmbeddedThicknesses []float32
}

func NewTextBox(fontId string, x, y float32, text ...any) TextBox {
	var node = NewNode(fontId, x, y)
	var textBox = TextBox{
		Node: node, Text: symbols.New(text...), LineHeight: 100,
		Thickness: 0.5, Smoothness: 0.02, SymbolGap: 0.2, WordWrap: true,
		EmbeddedColorsTag: '`', EmbeddedAssetsTag: '^', EmbeddedThicknessesTag: '*',
	}
	var font = textBox.font()
	var measure = rl.MeasureTextEx(*font, textBox.Text, textBox.LineHeight, textBox.gapSymbols())
	textBox.Width, textBox.Height = measure.X, measure.Y
	return textBox
}

func (textBox *TextBox) Size() (width, height float32) {
	return textBox.Width, textBox.Height
}

func (textBox *TextBox) TextWrap(text string) string {
	var font = textBox.font()
	var words = strings.Split(text, " ")
	var curX, curY = textBox.X, textBox.Y
	var gapSymbols = textBox.gapSymbols()
	var buffer = bytes.NewBufferString("")

	for w := range words {
		var word = words[w]
		if w < len(words)-1 {
			word += " " // split removes spaces, add it for all words but last one
		}

		var wordLength = symbols.Count(word)
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
func (textBox *TextBox) TextSymbols() string {
	var lines = textBox.TextLines()
	var result = ""
	for _, v := range lines {
		result += v
	}
	return result
}
func (textBox *TextBox) TextLines() []string {
	var lines, _ = textBox.formatSymbols()
	return lines
}
func (textBox *TextBox) TextSymbol(camera *Camera, symbolIndex int) (cX, cY, cWidth, cHeight, cAngle float32) {
	var _, symbols = textBox.formatSymbols()
	if symbolIndex < 0 || symbolIndex >= len(symbols) {
		return
	}

	var symbol = symbols[symbolIndex]
	cX, cY = textBox.PointToCamera(camera, symbol.X, symbol.Y)
	return cX, cY, symbol.Width, symbol.Height, symbol.Angle
}

// #region private

type symbol struct {
	X, Y, Angle, Width, Height,
	Thickness float32
	Value, AssetId string
	Font           *rl.Font
	Color          rl.Color
}

func (t *TextBox) formatSymbols() ([]string, []symbol) {
	var result = []symbol{}
	var resultLines = []string{}
	var assetTag = string(t.EmbeddedAssetsTag)
	var colorTag = string(t.EmbeddedColorsTag)
	var thickTag = string(t.EmbeddedThicknessesTag)
	var wrapped = t.TextWrap(t.Text)
	var lines = strings.Split(wrapped, "\n")
	var _, _, ang, _, _ = t.TransformToCamera()
	var curX, curY float32 = 0, 0
	var font = t.font()
	var textHeight = (t.LineHeight+t.gapLines())*float32(len(lines)) - t.gapLines()
	var curColor = rl.GetColor(t.Color)
	var curThick = t.Thickness
	var alignX, alignY = number.Limit(t.AlignmentX, 0, 1), number.Limit(t.AlignmentY, 0, 1)
	var colorIndex, assetIndex, thickIndex = 0, 0, 0
	var lastChar = ""
	// although some chars are "outside" of the box, they still need to be iterated cuz of colorIndex and assetIndex

	for l, line := range lines {
		var tagless = strings.ReplaceAll(line, colorTag, "")
		tagless = strings.ReplaceAll(tagless, thickTag, "")
		var lineSize = rl.MeasureTextEx(*font, tagless, t.LineHeight, t.gapSymbols())
		var lineLength = symbols.Count(line)
		var skip = false // replaces 'continue' to avoid skipping the offset calculations

		curX = (t.Width - lineSize.X) * alignX
		curY = float32(l)*(t.LineHeight+t.gapLines()) + (t.Height-textHeight)*alignY

		// hide text outside the box left, top & bottom
		if curX < 0 || curY < 0 || curY+t.LineHeight-1 > t.Height {
			skip = true // no need for right cuz text wraps there
		}

		for c := range lineLength {
			var char = string(line[c])
			var charSize = rl.MeasureTextEx(*font, char, t.LineHeight, 0)

			if line[c] == '\r' {
				lastChar = char
				continue // use as zerospace character or skip anyway
			}

			if curX+charSize.X > t.Width {
				skip = true
			}

			if char == colorTag {
				if colorIndex < len(t.EmbeddedColors) {
					curColor = rl.GetColor(t.EmbeddedColors[colorIndex])
					colorIndex++
					continue
				}
				curColor = rl.GetColor(t.Color)
				continue
			}

			if char == thickTag {
				if thickIndex < len(t.EmbeddedThicknesses) {
					curThick = t.EmbeddedThicknesses[thickIndex]
					thickIndex++
					continue

				}
				curThick = t.Thickness
				continue
			}

			var isAsset = char == assetTag && lastChar != assetTag && assetIndex < len(t.EmbeddedAssetIds)

			if !skip {
				var symbol = symbol{
					X: curX, Y: curY,
					Width: charSize.X, Height: t.LineHeight,
					Angle:     ang,
					Thickness: curThick,
					Value:     char,
					Color:     curColor,
					Font:      font,
				}
				if isAsset {
					symbol.AssetId = t.EmbeddedAssetIds[assetIndex]
				}

				result = append(result, symbol)

				if l == len(resultLines) {
					resultLines = append(resultLines, "")
				}
				resultLines[l] += symbol.Value

				curX += charSize.X + t.gapSymbols()
			}

			if isAsset {
				assetIndex++
			}

			if char != "\n" {
				lastChar = char
			}
		}
	}

	return resultLines, result
}

func (textBox *TextBox) font() *rl.Font {
	var font, hasFont = internal.Fonts[textBox.AssetId]
	var defaultFont, hasDefault = internal.Fonts[""]

	if !hasFont && hasDefault {
		font = defaultFont
		hasFont = true // fallback to engine default
	}

	if !hasFont {
		var defaultFont = rl.GetFontDefault()
		font = &defaultFont // fallback to raylib default
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
