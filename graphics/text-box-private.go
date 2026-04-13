package graphics

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/internal"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/number"
	txt "pure-game-kit/utility/text"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type symbol struct {
	Angle                 float32
	Value                 string
	Bounds, Rect, TexRect rl.Rectangle
	Texture               rl.Texture2D

	TopCrop float32

	Color, BackColor, OutlineColor, ShadowColor     uint
	Weight, OutlineWeight, ShadowWeight, ShadowBlur byte
	Underline, Strikethrough                        bool
}
type symbolState struct {
	Color, BackColor, OutlineColor, ShadowColor     uint
	Weight, OutlineWeight, ShadowWeight, ShadowBlur byte
	Underline, Strikethrough, HasAsset              bool
	AssetId                                         string
}

func (t *TextBox) formatSymbols() ([]string, []symbol) {
	var state = textBoxCache{
		t.Text, t.FontId, t.Tint, t.WordWrap,
		t.Width, t.Height,
		t.AlignmentX, t.AlignmentY,
		t.LineHeight, t.SymbolGap, t.LineGap,
	}
	if t.cache == state {
		return t.cacheChars, t.cacheSymbols
	}

	t.cacheSymbols = t.cacheSymbols[:0]
	t.cacheChars = t.cacheChars[:0]

	var wrapped = t.Text
	if t.WordWrap {
		wrapped = t.TextWrap(t.Text)
	}
	var lines = txt.SplitLines(wrapped)
	var font = t.font()
	var gapX = t.gapSymbols()
	var gapY = t.gapLines()
	var textHeight = (t.LineHeight+gapY)*float32(len(lines)) - gapY
	var alignX, alignY = number.Limit(t.AlignmentX, 0, 1), number.Limit(t.AlignmentY, 0, 1)
	var curState = symbolState{
		Color: t.Tint, OutlineColor: palette.Black, ShadowColor: palette.Black, Weight: 1, OutlineWeight: 1, ShadowWeight: 1,
		ShadowBlur: 1,
	}

	var reading = false
	if t.cacheBuilder == nil {
		t.cacheBuilder = txt.NewBuilder()
	}
	t.cacheBuilder.Clear()
	var currentLineStr strings.Builder

	for l, line := range lines {
		var emptyLine = line == ""
		if emptyLine {
			line = " "
		}

		var tagless = internal.RemoveTags(line)
		var lineWidth, _ = t.measure(font, tagless, gapX)
		var assetCount = txt.CountOccurrences(tagless, string(placeholderCharAsset))
		if assetCount > 0 {
			var placeholderWidth, _ = t.measure(font, string(placeholderCharAsset), gapX)
			lineWidth += (t.LineHeight - placeholderWidth) * float32(assetCount)
		}

		var curX = (t.Width - lineWidth) * alignX
		var curY = float32(l)*(t.LineHeight+gapY) + (t.Height-textHeight)*alignY
		var lineStartIdx = len(t.cacheSymbols) // track where this line starts in the symbol slice to build the string later

		for _, c := range line {
			if t.readTag(&reading, c, t.cacheBuilder, &curState) {
				continue
			}

			var symb symbol
			var charSize float32

			// Allocation Warning: string(c) still allocates.
			// Consider if your symbol.Value can be a rune.
			var char = condition.If(emptyLine, "", string(c))

			if curState.HasAsset {
				charSize = t.LineHeight
				var rect = rl.NewRectangle(curX, curY, charSize, charSize)
				var tex, src, rot, flip = internal.AssetData(curState.AssetId)
				internal.EditAssetRects(&src, &rect, t.Angle, rot, flip)
				symb = symbol{Texture: tex, Rect: rect, Bounds: rect, TexRect: src}
				curState.HasAsset = false
			} else {
				charSize, _ = t.measure(font, char, 0)
				symb = t.createSymbol(font, curX, curY, c)
			}

			symb.Bounds.Width, symb.Angle, symb.Value = charSize, 0, char
			symb.Color, symb.BackColor = curState.Color, curState.BackColor
			symb.OutlineColor, symb.ShadowColor = curState.OutlineColor, curState.ShadowColor
			symb.Weight, symb.OutlineWeight = curState.Weight, curState.OutlineWeight
			symb.ShadowWeight, symb.ShadowBlur = curState.ShadowWeight, curState.ShadowBlur
			symb.Underline, symb.Strikethrough = curState.Underline, curState.Strikethrough

			if !t.cropSymbol(symb, gapX, gapY) {
				t.cacheSymbols = append(t.cacheSymbols, symb)
			}

			curX += condition.If(charSize > 0, charSize+gapX, 0)
		}

		currentLineStr.Reset()
		for i := lineStartIdx; i < len(t.cacheSymbols); i++ {
			currentLineStr.WriteString(t.cacheSymbols[i].Value)
		}
		t.cacheChars = append(t.cacheChars, currentLineStr.String())
	}

	t.cache = state
	return t.cacheChars, t.cacheSymbols
}
func (t *TextBox) createSymbol(f rl.Font, x, y float32, c rune) symbol {
	var scaleFactor, padding = float32(t.LineHeight) / float32(f.BaseSize), float32(f.CharsPadding)
	var glyph, atlasRec = rl.GetGlyphInfo(f, int32(c)), rl.GetGlyphAtlasRec(f, int32(c))
	var tx, ty = atlasRec.X - padding, atlasRec.Y - padding
	var tw, th = atlasRec.Width + 2.0*padding, atlasRec.Height + 2.0*padding
	var rx = x + (float32(glyph.OffsetX)-padding)*scaleFactor
	var ry = y + (float32(glyph.OffsetY)-padding)*scaleFactor
	var rw = (atlasRec.Width + 2.0*padding) * scaleFactor
	var rh = (atlasRec.Height + 2.0*padding) * scaleFactor
	var src = rl.NewRectangle(tx, ty, tw, th)
	var dst = rl.NewRectangle(rx, ry, rw, rh)
	var bds = rl.Rectangle{X: x, Y: y, Height: t.LineHeight}
	return symbol{Texture: f.Texture, Rect: dst, TexRect: src, Bounds: bds}
}
func (t *TextBox) cropSymbol(symb symbol, gapX, gapY float32) (skip bool) {
	var rx, ry = symb.Rect.X, symb.Rect.Y
	var bx, by = symb.Bounds.X, symb.Bounds.Y
	var outsideHor = bx+symb.Bounds.Width+gapX < 0 || bx > t.Width
	var outsideVer = by+symb.Bounds.Height+gapY < 0 || by > t.Height
	skip = outsideHor || outsideVer

	var onEdgeLeft = !skip && rx < 0
	var onEdgeRight = !skip && rx+symb.Rect.Width > t.Width
	var onEdgeTop = !skip && ry < 0
	var onEdgeBottom = !skip && ry+symb.Rect.Height > t.Height

	var onEdgeLeftBounds = !skip && bx < 0
	var onEdgeRightBounds = !skip && bx+symb.Bounds.Width+gapX > t.Width
	var onEdgeTopBounds = !skip && by < 0
	var onEdgeBottomBounds = !skip && by+symb.Bounds.Height+gapY > t.Height

	if onEdgeLeft {
		var ratio = -rx / symb.Rect.Width
		symb.Rect.Width -= symb.Rect.Width * ratio
		symb.Rect.X, symb.Rect.Y = 0, ry
		symb.TexRect.X += symb.TexRect.Width * ratio
		symb.TexRect.Width -= symb.TexRect.Width * ratio
	}
	if onEdgeRight {
		var overflow = rx + symb.Rect.Width - t.Width
		var ratio = overflow / symb.Rect.Width
		symb.Rect.Width -= symb.Rect.Width * ratio
		symb.TexRect.Width -= symb.TexRect.Width * ratio
	}
	if onEdgeTop {
		var ratio = -ry / symb.Rect.Height
		symb.Rect.Height -= symb.Rect.Height * ratio
		symb.Rect.Height = max(symb.Rect.Height, 0)
		symb.Rect.X, symb.Rect.Y = rx, 0
		symb.TexRect.Y += symb.TexRect.Height * ratio
		symb.TexRect.Height -= symb.TexRect.Height * ratio
	}
	if onEdgeBottom {
		var overflow = ry + symb.Rect.Height - t.Height
		var ratio = overflow / symb.Rect.Height
		symb.Rect.Height -= symb.Rect.Height * ratio
		symb.Rect.Height = max(symb.Rect.Height, 0)
		symb.TexRect.Height -= symb.TexRect.Height * ratio
	}

	if onEdgeLeftBounds {
		var ratio = -bx / symb.Bounds.Width
		symb.Bounds.Width -= symb.Bounds.Width * ratio
		symb.Bounds.X, symb.Bounds.Y = 0, by
	}
	if onEdgeRightBounds {
		var overflow = bx + symb.Bounds.Width + gapX - t.Width
		var ratio = overflow / symb.Bounds.Width
		symb.Bounds.Width -= symb.Bounds.Width * ratio
	}
	if onEdgeTopBounds {
		var ratio = -by / symb.Bounds.Height
		var boundsCut = symb.Bounds.Height * ratio
		symb.Bounds.Height -= boundsCut
		symb.Bounds.X, symb.Bounds.Y = bx, 0
		symb.TopCrop += boundsCut
	}
	if onEdgeBottomBounds {
		var overflow = by + symb.Bounds.Height + gapY - t.Height
		var ratio = overflow / symb.Bounds.Height
		symb.Bounds.Height -= symb.Bounds.Height * ratio
	}

	return skip
}
func (t *TextBox) measure(font rl.Font, text string, spacing float32) (width, height float32) {
	if font.Texture.ID == 0 || len(text) == 0 {
		return 0, 0
	}

	var tempByteCounter int
	var byteCounter int
	var textWidth float32
	var maxTextWidth float32
	var scaleFactor = t.LineHeight / float32(font.BaseSize)
	var lineGap = t.gapLines()
	var textHeight = t.LineHeight // initial height is just one line

	for _, letter := range text {
		byteCounter++

		if letter != '\n' {
			var glyph = rl.GetGlyphInfo(font, int32(letter))
			var rec = rl.GetGlyphAtlasRec(font, int32(letter))

			if glyph.AdvanceX > 0 {
				textWidth += float32(glyph.AdvanceX)
			} else {
				textWidth += float32(rec.Width) + float32(glyph.OffsetX)
			}
		} else {
			if maxTextWidth < textWidth {
				maxTextWidth = textWidth
			}

			byteCounter = 0
			textWidth = 0
			textHeight += (t.LineHeight + lineGap) // custom height logic: adds a full line + custom gap
		}

		if tempByteCounter < byteCounter {
			tempByteCounter = byteCounter
		}
	}

	if maxTextWidth < textWidth {
		maxTextWidth = textWidth
	}

	// calculate final width matching raylib's scaling math
	var finalWidth = maxTextWidth*scaleFactor + float32(tempByteCounter-1)*spacing
	return finalWidth, textHeight
}

func (t *TextBox) font() rl.Font {
	var font, hasFont = internal.Fonts[t.FontId]
	var defaultFont, hasDefault = internal.Fonts[""]

	if !hasFont && hasDefault {
		font = defaultFont
		hasFont = true // fallback to engine default
	}

	if !hasFont {
		var fallback = rl.GetFontDefault()
		font = fallback // fallback to raylib default
	}
	return font
}
func (t *TextBox) gapSymbols() float32 {
	return t.SymbolGap * t.LineHeight / 5
}
func (t *TextBox) gapLines() float32 {
	return t.LineGap * t.LineHeight / 5
}

func (t *TextBox) readTag(reading *bool, char rune, cur *txt.Builder, s *symbolState) bool {
	if !*reading && char == '{' {
		*reading = true
	}

	if *reading {
		cur.WriteSymbol(char)
		if char == '}' {
			*reading = false
			// Tag is finished, fall through to process it
		} else {
			return true
		}
	}

	var tag = cur.ToText()
	if tag == "" || !txt.StartsWith(tag, "{") {
		return false
	}

	cur.Clear()
	tag = txt.Remove(tag, "{", "}")

	// Empty tag {} resets all formatting to defaults
	if tag == "" {
		s.reset(t.Tint)
		return false
	}

	var parts = txt.Split(tag, "=")
	var name = parts[0]
	var value string
	if len(parts) > 1 {
		value = parts[1]
	}

	switch name {
	case "color":
		s.Color = parseCol(value, t.Tint)
	case "backColor":
		s.BackColor = parseCol(value, 0)
	case "outlineColor":
		s.OutlineColor = parseCol(value, palette.Black)
	case "shadowColor":
		s.ShadowColor = parseCol(value, palette.Black)
	case "_":
		s.Underline = !s.Underline
	case "-":
		s.Strikethrough = !s.Strikethrough
	case "shadowBlur":
		s.ShadowBlur = byte(parseNum(value, 0))
	case "weight", "outlineWeight", "shadowWeight":
		var val byte = 1
		switch value {
		case "thin":
			val = 0
		case "regular":
			val = 1
		case "semiBold":
			val = 2
		case "bold":
			val = 3
		default:
			val = byte(parseNum(value, 1))
		}

		if name == "weight" {
			s.Weight = val
		}
		if name == "outlineWeight" {
			s.OutlineWeight = val
		}
		if name == "shadowWeight" {
			s.ShadowWeight = val
		}
	case "asset":
		s.AssetId = value
		s.HasAsset = true
	default:
		// Handle unknown tags or custom string values if necessary
	}

	return false
}
func (s *symbolState) reset(defaultTint uint) {
	s.Color = defaultTint
	s.BackColor = 0
	s.OutlineColor = palette.Black
	s.ShadowColor = palette.Black
	s.Weight = 1
	s.OutlineWeight = 1
	s.ShadowWeight = 1
	s.ShadowBlur = 1
	s.Underline = false
	s.Strikethrough = false
	s.AssetId = ""
	s.HasAsset = false
}

//=================================================================

func parseCol(value string, defaultValue uint) uint {
	if value == "" {
		return defaultValue
	}

	var rgba = txt.Split(value, " ")
	if len(rgba) == 4 {
		var r, g = txt.ToNumber[byte](rgba[0]), txt.ToNumber[byte](rgba[1])
		var b, a = txt.ToNumber[byte](rgba[2]), txt.ToNumber[byte](rgba[3])
		return color.RGBA(r, g, b, a)
	}
	return defaultValue
}
func parseNum(value string, defaultValue float32) float32 {
	if value == "" {
		return defaultValue
	}

	var result = txt.ToNumber[float32](value)
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
func getOrDefault(curValues map[string]any, name string, defaultValue any) any {
	var val, has = curValues[name]
	if has {
		return val
	}
	return defaultValue
}
