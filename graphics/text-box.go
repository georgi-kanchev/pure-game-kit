package graphics

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/internal"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/random"
	"pure-game-kit/utility/text"
	txt "pure-game-kit/utility/text"
	"regexp"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type TextBox struct {
	Node
	Text, FontId string
	WordWrap     bool
	AlignmentX, AlignmentY,
	Thickness, Smoothness,
	LineHeight, SymbolGap, LineGap float32

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

	var replaced, originals = internal.ReplaceStrings(text, '{', '}', internal.Placeholder)
	var words = txt.Split(replaced, " ")
	var curX, curY float32 = 0, 0
	var buffer = txt.NewBuilder()
	var tagIndex = 0
	var width = t.Width

	for w := range words {
		var word = words[w]

		if w < len(words)-1 {
			word += " " // split removes spaces, add it for all words but last one
		}

		var wordSize, _ = t.TextMeasure(txt.Trim(word))
		var wordEndOfBox = curX+wordSize > width
		var wordFirst = w == 0
		var wordNewLine = !wordFirst && t.WordWrap && wordEndOfBox

		if wordNewLine {
			curX = 0
			curY += t.LineHeight + t.gapLines()
			buffer.WriteSymbol('\n')
		}

		for i, c := range word {
			var char = string(c)
			var charSize, _ = t.TextMeasure(char)
			charSize = condition.If(c == internal.Placeholder, 0, charSize)
			var charEndOfBoxX = curX+charSize > width
			var charFirst = i == 0 && wordFirst
			var charNewLine = !charFirst && char != " " && (char == "\n" || charEndOfBoxX)

			if charNewLine {
				curX = 0
				curY += t.LineHeight + t.gapLines()

				if char != "\n" {
					buffer.WriteSymbol('\n')
				}
			}

			if c == internal.Placeholder {
				char = "{" + originals[tagIndex] + "}"
				tagIndex++
			}
			buffer.WriteText(char)
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

	var s = symbols[symbolIndex]
	return s.X, s.Y, s.Width, t.LineHeight, s.Angle
}

//=================================================================
// private

type symbol struct {
	Angle, Thickness float32
	Value, AssetId   string
	Rect, TexRect    rl.Rectangle
	Color            uint

	UnderlineSize float32
	X, Y, Width   float32
}

func (t *TextBox) formatSymbols(cam *Camera) ([]string, []symbol) {
	var curHash = random.Hash(t)
	if t.hash == curHash {
		return t.cacheChars, t.cacheSymbols
	}

	var result = []symbol{}
	var resultLines = []string{}
	var wrapped = t.TextWrap(t.Text)
	var lines = txt.SplitLines(wrapped)
	var curX, curY float32 = 0, 0
	var font = t.font()
	var gapX = t.gapSymbols()
	var textHeight = (t.LineHeight+t.gapLines())*float32(len(lines)) - t.gapLines()
	var alignX, alignY = number.Limit(t.AlignmentX, 0, 1), number.Limit(t.AlignmentY, 0, 1)
	var curColor = t.Tint
	var curUnderline float32
	var lineIndex = 0
	var reading = false
	var curTag = text.NewBuilder()
	var w = t.Width

	for l, line := range lines {
		var emptyLine = line == ""
		if emptyLine {
			line = " " // empty lines shouldn't be skipped
		}

		var tagless = removeTags(line)
		var lineWidth, _ = t.TextMeasure(tagless)
		var skip = false // replaces 'continue' to avoid skipping the offset calculations

		curX = (w - lineWidth) * alignX
		curY = float32(l)*(t.LineHeight+t.gapLines()) + (t.Height-textHeight)*alignY

		var outsideLeftTopOrBottom = curX < 0 || curY < 0 || curY+t.LineHeight-1 > t.Height
		if outsideLeftTopOrBottom {
			skip = true
		}

		for _, c := range line {
			if skip {
				break
			}

			if t.readTag(&reading, c, curTag, &curColor, &curUnderline) {
				continue
			}

			var char = condition.If(emptyLine, "", string(c))
			var charSize = rl.MeasureTextEx(*font, char, t.LineHeight, 0)
			var symbol = t.createSymbol(font, cam, curX, curY, charSize.X, t.Angle, curUnderline, c, char, curColor)
			var outsideRight = curX+charSize.X > w

			if outsideRight {
				skip = true
				break // rare cases but happens with single symbol & small width
			}

			result = append(result, symbol)

			lineIndex = number.Limit(lineIndex, 0, len(resultLines))
			if lineIndex == len(resultLines) {
				resultLines = append(resultLines, "")
			}

			resultLines[lineIndex] += symbol.Value
			curX += charSize.X + gapX
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

func (t *TextBox) readTag(reading *bool, char rune, cur *txt.Builder, col *uint, un *float32) (nextChar bool) {
	if !*reading && char == '{' {
		*reading = true
	}

	if *reading {
		cur.WriteSymbol(char)

		if char == '}' {
			*reading = false
		}

		return true
	}

	var tag = cur.ToText()
	if tag == "{#}" {
		*col = t.Tint
	} else if text.StartsWith(tag, "{#") {
		*col = parseCol(tag, "#", t.Tint)
	} else if tag == "{_}" {
		*un = 0
	} else if text.StartsWith(tag, "{_") {
		*un = parseNum(tag, "_", 0)
	}

	if tag != "" {
		cur.Clear()
	}
	return false
}

func (t *TextBox) createSymbol(f *rl.Font, cam *Camera, x, y, w, a, un float32, c rune, ch string, col uint) symbol {
	var scaleFactor, padding = float32(t.LineHeight) / float32(f.BaseSize), float32(f.CharsPadding)
	var glyph, atlasRec = rl.GetGlyphInfo(*f, int32(c)), rl.GetGlyphAtlasRec(*f, int32(c))
	var tx, ty = atlasRec.X - padding, atlasRec.Y - padding
	var tw, th = atlasRec.Width + 2.0*padding, atlasRec.Height + 2.0*padding
	var rx = x + (float32(glyph.OffsetX)-padding)*scaleFactor
	var ry = y + (float32(glyph.OffsetY)-padding)*scaleFactor
	var rw = (atlasRec.Width + 2.0*padding) * scaleFactor
	var rh = (atlasRec.Height + 2.0*padding) * scaleFactor
	var src, dst = rl.NewRectangle(tx, ty, tw, th), rl.NewRectangle(rx, ry, rw, rh)
	dst.X, dst.Y = t.PointToCamera(cam, dst.X, dst.Y)
	x, y = t.PointToCamera(cam, x, y)
	var symbol = symbol{
		Angle: a, Thickness: t.Thickness, Value: ch, Color: col, Rect: dst, TexRect: src, UnderlineSize: un,
		X: x, Y: y, Width: w}
	return symbol
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

func parseCol(tag, symbol string, defaultValue uint) uint {
	tag = text.Remove(tag, "{"+symbol, "}")
	var rgba = text.Split(tag, " ")
	if len(rgba) == 4 {
		var r, g = text.ToNumber[byte](rgba[0]), text.ToNumber[byte](rgba[1])
		var b, a = text.ToNumber[byte](rgba[2]), text.ToNumber[byte](rgba[3])
		return color.RGBA(r, g, b, a)
	}
	return defaultValue
}
func parseNum(tag, symbol string, defaultValue float32) float32 {
	tag = text.Remove(tag, "{"+symbol, "}")
	var result = text.ToNumber[float32](tag)
	if number.IsNaN(result) {
		return defaultValue
	}
	return result
}
func removeTags(text string) string {
	return regexp.MustCompile(`{.*?}`).ReplaceAllString(text, "")
}
