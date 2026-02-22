package graphics

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/internal"
	"pure-game-kit/utility/collection"
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
	cacheSymbols []*symbol
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
	var ph = string(internal.Placeholder)

	for w := range words {
		var word = words[w]

		if w < len(words)-1 {
			word += " " // split removes spaces, add it for all words but last one
		}

		var trimWord = txt.Remove(txt.Trim(word), ph)
		var wordSize, _ = t.TextMeasure(trimWord)
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

func (t *TextBox) formatSymbols(cam *Camera) ([]string, []*symbol) {
	var curHash = random.Hash(t)
	if t.hash == curHash {
		return t.cacheChars, t.cacheSymbols
	}

	var result = []*symbol{}
	var resultLines = []string{}
	var wrapped = t.TextWrap(t.Text)
	var lines = txt.SplitLines(wrapped)
	var curX, curY float32 = 0, 0
	var font = t.font()
	var gapX = t.gapSymbols()
	var textHeight = (t.LineHeight+t.gapLines())*float32(len(lines)) - t.gapLines()
	var alignX, alignY = number.Limit(t.AlignmentX, 0, 1), number.Limit(t.AlignmentY, 0, 1)
	var curValues = make(map[string]any)
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
		var skip = false // avoids interrupting the tags if partially culled

		curX = (w - lineWidth) * alignX
		curY = float32(l)*(t.LineHeight+t.gapLines()) + (t.Height-textHeight)*alignY

		var outsideLeftTopOrBottom = curX < 0 || curY < 0 || curY+t.LineHeight-1 > t.Height
		if outsideLeftTopOrBottom {
			skip = true
		}

		for _, c := range line {
			if t.readTag(&reading, c, curTag, curValues) {
				continue
			}

			var col = getOrDefault(curValues, "color", t.Tint).(uint)
			var underline = getOrDefault(curValues, "underline", float32(0)).(float32)
			var symb symbol
			var charSize float32

			var assetId, has = curValues["assetId"]
			if has {
				var x, y = t.PointToCamera(cam, curX, curY)
				charSize = t.LineHeight * 0.9
				y += t.LineHeight * 0.05
				var rect = rl.NewRectangle(x, y, charSize, t.LineHeight*0.9)
				symb = symbol{AssetId: assetId.(string), Angle: t.Angle, Color: col, Rect: rect,
					UnderlineSize: underline, Value: "@"}
				delete(curValues, "assetId")
			} else {
				var char = condition.If(emptyLine, "", string(c))
				charSize = rl.MeasureTextEx(*font, char, t.LineHeight, 0).X
				symb = t.createSymbol(font, cam, curX, curY, c)

				symb.Width, symb.Angle = charSize, t.Angle
				symb.Color, symb.Value = col, char
				symb.UnderlineSize = underline
			}
			if curX+charSize > w { // outside right
				skip = true // rare cases but happens with single symbol & small width
			}

			if !skip {
				result = append(result, &symb)
			}

			lineIndex = number.Limit(lineIndex, 0, len(resultLines))
			if lineIndex == len(resultLines) {
				resultLines = append(resultLines, "")
			}

			resultLines[lineIndex] += symb.Value
			curX += charSize + gapX
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

func (t *TextBox) readTag(reading *bool, char rune, cur *txt.Builder, curValues map[string]any) (nextChar bool) {
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
	if tag != "" {
		cur.Clear()
	}

	if !text.StartsWith(tag, "{") {
		return false
	}

	tag = text.Remove(tag, "{", "}")

	if tag == "" {
		collection.MapClear(curValues)
		return false
	}

	var parts = text.Split(tag, "=")
	var name, value = parts[0], ""
	if len(parts) > 1 {
		value = parts[1]
	}

	switch name {
	case "color":
		curValues[name] = parseCol(value, t.Tint)
	case "underline":
		curValues[name] = parseNum(value, 0)
	default:
		curValues[name] = value
	}

	return false
}

func (t *TextBox) createSymbol(f *rl.Font, cam *Camera, x, y float32, c rune) symbol {
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

	var symbol = symbol{Thickness: t.Thickness, Rect: dst, TexRect: src, X: x, Y: y}
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

func parseCol(value string, defaultValue uint) uint {
	var rgba = text.Split(value, " ")
	if len(rgba) == 4 {
		var r, g = text.ToNumber[byte](rgba[0]), text.ToNumber[byte](rgba[1])
		var b, a = text.ToNumber[byte](rgba[2]), text.ToNumber[byte](rgba[3])
		return color.RGBA(r, g, b, a)
	}
	return defaultValue
}
func parseNum(value string, defaultValue float32) float32 {
	var result = text.ToNumber[float32](value)
	if number.IsNaN(result) {
		return defaultValue
	}
	return result
}
func removeTags(text string) string {
	return regexp.MustCompile(`{.*?}`).ReplaceAllString(text, "")
}

func getOrDefault(curValues map[string]any, name string, defaultValue any) any {
	var val, has = curValues[name]
	if has {
		return val
	}
	return defaultValue
}
