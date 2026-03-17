package graphics

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/internal"
	"pure-game-kit/utility/collection"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/random"
	"pure-game-kit/utility/text"
	txt "pure-game-kit/utility/text"
	"regexp"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type TextBox struct {
	Quad
	Text, FontId string
	WordWrap     bool
	AlignmentX, AlignmentY,
	LineHeight, SymbolGap, LineGap,
	ShadowOffsetX, ShadowOffsetY float32

	// Skip advanced feature properties for faster render. Properties used:
	// 	FontId, Text, X, Y, Width, LineHeight, WordWrap, Thickness, SymbolGap, Tint
	Fast bool

	hash         uint32
	cacheChars   []string
	cacheSymbols []*symbol
	cacheWrap    string
}

func NewTextBox(fontId string, x, y float32, text ...any) *TextBox {
	var quad = NewQuad(x, y)
	var textBox = &TextBox{
		FontId: fontId, Quad: *quad, Text: txt.New(text...), LineHeight: 100, SymbolGap: 0.2, WordWrap: true,
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
	var ph = string(internal.Placeholder)
	var gapY = t.gapLines()

	for w := range words {
		var word = words[w]

		if w < len(words)-1 {
			word += " " // split removes spaces, add it for all words but last one
		}

		var trimWord = txt.Remove(txt.Trim(word), ph)
		var wordSize, _ = t.TextMeasure(trimWord)

		if !t.Fast && txt.Contains(trimWord, string(placeholderCharAsset)) {
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
			charSize = condition.If(!t.Fast && c == placeholderCharAsset, t.LineHeight, charSize)
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

			if !t.Fast && c == internal.Placeholder {
				char = "{" + originals[tagIndex] + "}"
				tagIndex++
			}
			buffer.WriteText(char)
			curX += condition.If(charSize > 0, charSize+t.gapSymbols(), 0)
		}
	}
	var result = buffer.ToText()
	result = txt.Replace(result, " \n", "\n")
	t.hash = curHash
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
	return s.Bounds.X, s.Bounds.Y, s.Bounds.Width, t.LineHeight, s.Angle
}

//=================================================================
// private

type symbol struct {
	Angle                 float32
	Value                 string
	Bounds, Rect, TexRect rl.Rectangle
	Texture               *rl.Texture2D

	TopCrop float32

	Color, BackColor, OutlineColor, ShadowColor     uint
	Weight, OutlineWeight, ShadowWeight, ShadowBlur byte
	Underline, Strikethrough                        bool
}

func (t *TextBox) formatSymbols() ([]string, []*symbol) {
	var curHash = random.Hash(t)
	if t.hash == curHash {
		return t.cacheChars, t.cacheSymbols
	}

	var result = []*symbol{}
	var resultLines = []string{}
	var wrapped = condition.If(t.WordWrap, t.TextWrap(t.Text), t.Text)
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

	for l, line := range lines {
		var emptyLine = line == ""
		if emptyLine {
			line = " " // empty lines shouldn't be skipped
		}

		var tagless = removeTags(line)
		var lineWidth, _ = t.TextMeasure(tagless)

		var assetCount = text.CountOccurrences(tagless, string(placeholderCharAsset))
		if assetCount > 0 { // account embedded assets in line width
			var placeholderWidth, _ = t.TextMeasure(string(placeholderCharAsset))
			lineWidth += (t.LineHeight - placeholderWidth) * float32(assetCount)
		}

		curX = (t.Width - lineWidth) * alignX
		curY = float32(l)*(t.LineHeight+t.gapLines()) + (t.Height-textHeight)*alignY

		for _, c := range line {
			if t.readTag(&reading, c, curTag, curValues) {
				continue
			}

			var symb symbol
			var charSize float32
			var char = condition.If(emptyLine, "", string(c))
			var assetId, has = curValues["assetId"]

			if has {
				var x, y = t.PointToGlobal(curX, curY)
				charSize = t.LineHeight
				var rect = rl.NewRectangle(x, y, charSize, charSize)
				var tex, src, rot, flip = internal.AssetData(assetId.(string))
				internal.EditAssetRects(&src, &rect, t.Angle, rot, flip)
				rect.Width *= t.ScaleX
				rect.Height *= t.ScaleY
				symb = symbol{Texture: tex, Rect: rect, Bounds: rect, TexRect: src}
				delete(curValues, "assetId")
			} else {
				charSize = rl.MeasureTextEx(*font, char, t.LineHeight, 0).X
				symb = t.createSymbol(font, curX, curY, c)
			}
			symb.Bounds.Width, symb.Angle, symb.Value = charSize*t.ScaleX, t.Angle, char
			symb.Color = getOrDefault(curValues, "color", t.Tint).(uint)
			symb.BackColor = getOrDefault(curValues, "backColor", uint(0)).(uint)
			symb.OutlineColor = getOrDefault(curValues, "outlineColor", palette.Black).(uint)
			symb.ShadowColor = getOrDefault(curValues, "shadowColor", palette.Black).(uint)
			symb.Weight = getOrDefault(curValues, "weight", byte(1)).(byte)
			symb.OutlineWeight = getOrDefault(curValues, "outlineWeight", byte(2)).(byte)
			symb.ShadowWeight = getOrDefault(curValues, "shadowWeight", byte(1)).(byte)
			symb.ShadowBlur = getOrDefault(curValues, "shadowBlur", byte(0)).(byte)
			symb.Underline = getOrDefault(curValues, "_", false).(bool)
			symb.Strikethrough = getOrDefault(curValues, "-", false).(bool)

			if !t.cropSymbol(&symb) {
				result = append(result, &symb)
			}

			lineIndex = number.Limit(lineIndex, 0, len(resultLines))
			if lineIndex == len(resultLines) {
				resultLines = append(resultLines, "")
			}

			resultLines[lineIndex] += symb.Value
			curX += condition.If(charSize > 0, charSize+gapX, 0)
		}

		lineIndex++
	}

	t.hash = curHash
	t.cacheChars = resultLines
	t.cacheSymbols = result
	return resultLines, result
}
func (t *TextBox) createSymbol(f *rl.Font, x, y float32, c rune) symbol {
	var scaleFactor, padding = float32(t.LineHeight) / float32(f.BaseSize), float32(f.CharsPadding)
	var glyph, atlasRec = rl.GetGlyphInfo(*f, int32(c)), rl.GetGlyphAtlasRec(*f, int32(c))
	var tx, ty = atlasRec.X - padding, atlasRec.Y - padding
	var tw, th = atlasRec.Width + 2.0*padding, atlasRec.Height + 2.0*padding
	var rx = x + (float32(glyph.OffsetX)-padding)*scaleFactor
	var ry = y + (float32(glyph.OffsetY)-padding)*scaleFactor
	var rw = (atlasRec.Width + 2.0*padding) * scaleFactor
	var rh = (atlasRec.Height + 2.0*padding) * scaleFactor
	var src, dst = rl.NewRectangle(tx, ty, tw, th), rl.NewRectangle(rx, ry, rw, rh)
	var bds = rl.Rectangle{}
	dst.X, dst.Y = t.PointToGlobal(dst.X, dst.Y)
	bds.X, bds.Y = t.PointToGlobal(x, y)
	dst.Width *= t.ScaleX
	dst.Height *= t.ScaleY
	bds.Height = t.LineHeight * t.ScaleY

	var symbol = symbol{Texture: &f.Texture, Rect: dst, TexRect: src, Bounds: bds}
	return symbol
}
func (t *TextBox) cropSymbol(symb *symbol) (skip bool) {
	var rx, ry = t.PointToLocal(symb.Rect.X, symb.Rect.Y)
	var bx, by = t.PointToLocal(symb.Bounds.X, symb.Bounds.Y)
	var outsideHor = bx+symb.Bounds.Width/t.ScaleX+t.gapSymbols() < 0 || bx > t.Width
	var outsideVer = by+symb.Bounds.Height/t.ScaleY+t.gapLines() < 0 || by > t.Height
	skip = outsideHor || outsideVer

	var onEdgeLeft = !skip && rx < 0
	var onEdgeRight = !skip && rx+symb.Rect.Width/t.ScaleX > t.Width
	var onEdgeTop = !skip && ry < 0
	var onEdgeBottom = !skip && ry+symb.Rect.Height/t.ScaleY > t.Height

	var onEdgeLeftBounds = !skip && bx < 0
	var onEdgeRightBounds = !skip && bx+symb.Bounds.Width/t.ScaleX+t.gapSymbols() > t.Width
	var onEdgeTopBounds = !skip && by < 0
	var onEdgeBottomBounds = !skip && by+symb.Bounds.Height/t.ScaleY+t.gapLines() > t.Height

	if onEdgeLeft {
		var lx, ly = t.PointToLocal(symb.Rect.X, symb.Rect.Y)
		var ratio = -lx / (symb.Rect.Width / t.ScaleX)
		var rectCut = symb.Rect.Width * ratio
		var texCut = symb.TexRect.Width * ratio
		symb.Rect.Width -= rectCut
		symb.Rect.X, symb.Rect.Y = t.PointToGlobal(0, ly)
		symb.TexRect.X += texCut
		symb.TexRect.Width -= texCut
	}
	if onEdgeRight {
		var lx, _ = t.PointToLocal(symb.Rect.X, symb.Rect.Y)
		var rightEdge = lx + (symb.Rect.Width / t.ScaleX)
		var overflow = rightEdge - t.Width
		var ratio = overflow / (symb.Rect.Width / t.ScaleX)
		symb.Rect.Width -= symb.Rect.Width * ratio
		symb.TexRect.Width -= symb.TexRect.Width * ratio
	}
	if onEdgeTop {
		var lx, ly = t.PointToLocal(symb.Rect.X, symb.Rect.Y)
		var ratio = -ly / (symb.Rect.Height / t.ScaleY)
		var rectCut = symb.Rect.Height * ratio
		var texCut = symb.TexRect.Height * ratio
		symb.Rect.Height -= rectCut
		symb.Rect.Height = max(symb.Rect.Height, 0)
		symb.Rect.X, symb.Rect.Y = t.PointToGlobal(lx, 0)
		symb.TexRect.Y += texCut
		symb.TexRect.Height -= texCut
	}
	if onEdgeBottom {
		var _, ly = t.PointToLocal(symb.Rect.X, symb.Rect.Y)
		var bottomEdge = ly + (symb.Rect.Height / t.ScaleY)
		var overflow = bottomEdge - t.Height
		var ratio = overflow / (symb.Rect.Height / t.ScaleY)
		symb.Rect.Height -= symb.Rect.Height * ratio
		symb.Rect.Height = max(symb.Rect.Height, 0)
		symb.TexRect.Height -= symb.TexRect.Height * ratio
	}

	if onEdgeLeftBounds {
		var lx, ly = t.PointToLocal(symb.Bounds.X, symb.Bounds.Y)
		var ratio = -lx / (symb.Bounds.Width / t.ScaleX)
		var boundsCut = symb.Bounds.Width * ratio
		symb.Bounds.X, symb.Bounds.Y = t.PointToGlobal(0, ly)
		symb.Bounds.Width -= boundsCut
	}
	if onEdgeRightBounds {
		var lx, _ = t.PointToLocal(symb.Bounds.X, symb.Bounds.Y)
		var rightEdge = lx + (symb.Bounds.Width / t.ScaleX) + t.gapSymbols()
		var overflow = rightEdge - t.Width
		var ratio = overflow / (symb.Bounds.Width / t.ScaleX)
		symb.Bounds.Width -= symb.Bounds.Width * ratio
	}
	if onEdgeTopBounds {
		var lx, ly = t.PointToLocal(symb.Bounds.X, symb.Bounds.Y)
		var ratio = -ly / (symb.Bounds.Height / t.ScaleY)
		var boundsCut = symb.Bounds.Height * ratio
		symb.Bounds.X, symb.Bounds.Y = t.PointToGlobal(lx, 0)
		symb.Bounds.Height -= boundsCut
		symb.TopCrop += boundsCut
	}
	if onEdgeBottomBounds {
		var _, ly = t.PointToLocal(symb.Bounds.X, symb.Bounds.Y)
		var bottomEdge = ly + (symb.Bounds.Height / t.ScaleY) + t.gapLines()
		var overflow = bottomEdge - t.Height
		var ratio = overflow / (symb.Bounds.Height / t.ScaleY)
		symb.Bounds.Height -= symb.Bounds.Height * ratio
	}

	return skip
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
	case "backColor":
		curValues[name] = parseCol(value, 0)
	case "outlineColor", "shadowColor":
		curValues[name] = parseCol(value, palette.Black)
	case "_", "-": // underline, strikethrough
		toggleBool(curValues, name)
	case "shadowBlur":
		curValues[name] = byte(parseNum(value, 0))
	case "weight", "outlineWeight", "shadowWeight":
		var val = 1
		val = condition.If(value == "thin", 0, val)
		val = condition.If(value == "semiBold", 2, val)
		val = condition.If(value == "bold", 3, val)
		curValues[name] = byte(val)
	default:
		curValues[name] = value
	}

	return false
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
func toggleBool(curValues map[string]any, name string) {
	if _, has := curValues[name]; has {
		delete(curValues, name)
	} else {
		curValues[name] = true
	}
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
