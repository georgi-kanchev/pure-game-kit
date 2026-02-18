package graphics

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/internal"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/random"
	txt "pure-game-kit/utility/text"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type TextBox struct {
	Node
	Text, FontId string
	WordWrap     bool
	AlignmentX, AlignmentY,
	Thickness, Smoothness,
	SymbolGap, LineHeight, LineGap float32
	EmbeddedColorsTag, EmbeddedAssetsTag,
	EmbeddedThicknessesTag rune
	EmbeddedAssetIds    []string
	EmbeddedColors      []uint
	EmbeddedThicknesses []float32

	// Skip advanced feature properties for faster render. Properties used:
	// 	FontId, Text, X, Y, Width, LineHeight, WordWrap, Thickness, SymbolGap, Tint
	Fast bool

	hash         uint32
	cacheChars   []string
	cacheSymbols []symbol
	cacheWrap    string
}

func NewTextBox(fontId string, x, y float32, text ...any) *TextBox {
	var node = NewNode(x, y)
	var textBox = &TextBox{
		FontId: fontId, Node: *node, Text: txt.New(text...), LineHeight: 100,
		Thickness: 0.5, Smoothness: 0.02, SymbolGap: 0.2, WordWrap: true,
		EmbeddedColorsTag: '`', EmbeddedAssetsTag: '^', EmbeddedThicknessesTag: '*',
	}
	var font = textBox.font()
	var measure = rl.MeasureTextEx(*font, textBox.Text, textBox.LineHeight, textBox.gapSymbols())
	textBox.Width, textBox.Height = measure.X, measure.Y
	return textBox
}

//=================================================================

// Does not wrap the text - use TextWrap(...) beforehand if intended.
func (t *TextBox) TextMeasure(text string) (width, height float32) {
	var size = rl.MeasureTextEx(*t.font(), text, t.LineHeight, t.gapSymbols())
	height = float32(txt.CountOccurrences(text, "\n")+1) * (t.LineHeight + t.gapLines())
	return size.X, height // raylib doesn't seem to calculate height correctly
}
func (t *TextBox) TextWrap(text string) string {
	var curHash = random.Hash(t)
	if t.hash == curHash {
		return t.cacheWrap
	}

	var words = txt.Split(text, " ")
	var curX, curY float32 = 0, 0
	var buffer = txt.NewBuilder()

	for w := range words {
		var word = words[w]
		if w < len(words)-1 {
			word += " " // split removes spaces, add it for all words but last one
		}

		var wordSize, _ = t.TextMeasure(txt.Trim(word))
		var wordEndOfBox = curX+wordSize > t.Width
		var firstWord = w == 0

		if !firstWord && t.WordWrap && wordEndOfBox {
			curX = 0
			curY += t.LineHeight + t.gapLines()
			buffer.WriteSymbol('\n')
		}

		for i, c := range word {
			var char = string(c)
			var charSize, _ = t.TextMeasure(string(char))
			var charEndOfBoxX = curX+charSize > t.Width
			var charFirst = i == 0 && firstWord

			if !charFirst && char != " " && (char == "\n" || charEndOfBoxX) {
				curX = 0
				curY += t.LineHeight + t.gapLines()

				if char != "\n" {
					buffer.WriteSymbol('\n')
				}
			}

			buffer.WriteText(char)

			if char == string(t.EmbeddedColorsTag) || char == string(t.EmbeddedThicknessesTag) {
				continue // these tags have 0 width when rendering so wrapping shouldn't be affected by them
			} // however, the assets tag has width and it should

			curX += charSize + t.gapSymbols()
		}
	}

	var result = buffer.ToText()
	result = txt.Replace(result, " \n", "\n")
	t.hash = curHash
	t.cacheWrap = result
	return result
}
func (t *TextBox) TextLines(camera *Camera) []string {
	var lines, _ = t.formatSymbols(camera)
	return lines
}
func (t *TextBox) TextSymbol(camera *Camera, symbolIndex int) (cX, cY, cWidth, cHeight, cAngle float32) {
	var _, symbols = t.formatSymbols(camera)
	if symbolIndex < 0 || symbolIndex >= len(symbols) {
		return number.NaN(), number.NaN(), number.NaN(), number.NaN(), number.NaN()
	}

	var symbol = symbols[symbolIndex]
	cX, cY = t.PointToCamera(camera, symbol.Rect.X, symbol.Rect.Y)
	return cX, cY, symbol.Rect.Width, symbol.Rect.Height, symbol.Angle
}

//=================================================================
// private

type symbol struct {
	Angle, Thickness float32
	Value, AssetId   string
	Rect, TexRect    rl.Rectangle
	Color            uint
}

func (t *TextBox) formatSymbols(cam *Camera) ([]string, []symbol) {
	var curHash = random.Hash(t)
	if t.hash == curHash {
		return t.cacheChars, t.cacheSymbols
	}

	var result = []symbol{}
	var resultLines = []string{}
	var assetTag = string(t.EmbeddedAssetsTag)
	var colorTag = string(t.EmbeddedColorsTag)
	var thickTag = string(t.EmbeddedThicknessesTag)
	var wrapped = t.TextWrap(t.Text)
	var lines = txt.SplitLines(wrapped)
	var _, _, ang, _, _ = t.TransformToCamera()
	var curX, curY float32 = 0, 0
	var font = t.font()
	var textHeight = (t.LineHeight+t.gapLines())*float32(len(lines)) - t.gapLines()
	var curColor = t.Tint
	var curThick = t.Thickness
	var alignX, alignY = number.Limit(t.AlignmentX, 0, 1), number.Limit(t.AlignmentY, 0, 1)
	var colorIndex, assetIndex, thickIndex = 0, 0, 0
	var lastChar = ""
	var lineIndex = 0
	// although some chars are "outside" of the box, they still need to be iterated cuz of colorIndex and assetIndex

	for l, line := range lines {
		var emptyLine = line == ""
		if emptyLine {
			line = " " // empty lines shouldn't be skipped
		}

		var tagless = txt.Remove(line, colorTag, thickTag)
		var lineWidth, _ = t.TextMeasure(tagless)
		var skip = false // replaces 'continue' to avoid skipping the offset calculations

		curX = (t.Width - lineWidth) * alignX
		curY = float32(l)*(t.LineHeight+t.gapLines()) + (t.Height-textHeight)*alignY

		// hide text outside the box left, top & bottom
		if curX < 0 || curY < 0 || curY+t.LineHeight-1 > t.Height {
			skip = true // no need for right cuz text wraps there
		}

		if !skip && lineIndex == 0 && l > 0 && txt.Length(line) > 3 {
			line = "..." + line[3:] // invisible lines before this one, indicate it
		}
		if curY+t.LineHeight*1.5-1 > t.Height && l < len(lines)-1 && txt.Length(line) > 3 {
			line = line[:len(line)-3] + "..." // invisible lines after this one, indicate it
		}

		for _, c := range line {
			var char = condition.If(emptyLine, "", string(c))
			var charSize = rl.MeasureTextEx(*font, char, t.LineHeight, 0)

			if char == "\r" {
				lastChar = char
				continue // use as zerospace character or skip anyway
			}

			if curX+charSize.X > t.Width+1 {
				skip = true
			}

			if char == colorTag {
				if colorIndex < len(t.EmbeddedColors) {
					curColor = t.EmbeddedColors[colorIndex]
					colorIndex++
					continue
				}
				curColor = t.Tint
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
				var scaleFactor = float32(t.LineHeight) / float32(font.BaseSize)
				var glyph = rl.GetGlyphInfo(*font, int32(c))
				var atlasRec = rl.GetGlyphAtlasRec(*font, int32(c))
				var padding = float32(font.CharsPadding)
				var rect = rl.NewRectangle(
					curX+(float32(glyph.OffsetX)-padding)*scaleFactor,
					curY+(float32(glyph.OffsetY)-padding)*scaleFactor,
					(atlasRec.Width+2.0*padding)*scaleFactor,
					(atlasRec.Height+2.0*padding)*scaleFactor)
				var texRect = rl.NewRectangle(
					atlasRec.X-padding,
					atlasRec.Y-padding,
					atlasRec.Width+2.0*padding,
					atlasRec.Height+2.0*padding)

				rect.X, rect.Y = t.PointToCamera(cam, rect.X, rect.Y)

				var symbol = symbol{
					Angle: ang, Thickness: curThick, Value: char, Color: curColor, Rect: rect, TexRect: texRect,
				}
				if isAsset {
					symbol.AssetId = t.EmbeddedAssetIds[assetIndex]
				}

				result = append(result, symbol)

				if lineIndex == len(resultLines) {
					resultLines = append(resultLines, "")
				}

				resultLines[lineIndex] += symbol.Value
				curX += charSize.X + t.gapSymbols()

			}

			if isAsset {
				assetIndex++
			}

			if char != "\n" {
				lastChar = char
			}
		}

		if !skip {
			lineIndex++
		}
	}

	t.hash = curHash
	t.cacheChars = resultLines
	t.cacheSymbols = result
	return resultLines, result
}

func (t *TextBox) font() *rl.Font {
	var font, hasFont = internal.Fonts[t.FontId]
	var defaultFont, hasDefault = internal.Fonts[""]

	if !hasFont && hasDefault {
		font = defaultFont
		hasFont = true // fallback to engine default
	}

	if !hasFont {
		var fallback = rl.GetFontDefault()
		font = &fallback // fallback to raylib default
	}
	return font
}
func (t *TextBox) gapSymbols() float32 {
	return t.SymbolGap * t.LineHeight / 5
}
func (t *TextBox) gapLines() float32 {
	return t.LineGap * t.LineHeight / 5
}
