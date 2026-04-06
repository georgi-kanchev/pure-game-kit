package graphics

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/internal"
	"pure-game-kit/utility/collection"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/text"

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

	var result = []symbol{}
	var resultLines = []string{}
	var wrapped = condition.If(t.WordWrap, t.TextWrap(t.Text), t.Text)
	var lines = text.SplitLines(wrapped)
	var curX, curY float32 = 0, 0
	var font = t.font()
	var gapX = t.gapSymbols()
	var gapY = t.gapLines()
	var textHeight = (t.LineHeight+gapY)*float32(len(lines)) - gapY
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

		var tagless = internal.RemoveTags(line)
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
				charSize = t.LineHeight
				var rect = rl.NewRectangle(curX, curY, charSize, charSize)
				var tex, src, rot, flip = internal.AssetData(assetId.(string))
				internal.EditAssetRects(&src, &rect, t.Angle, rot, flip)
				symb = symbol{Texture: tex, Rect: rect, Bounds: rect, TexRect: src}
				delete(curValues, "assetId")
			} else {
				charSize = rl.MeasureTextEx(font, char, t.LineHeight, 0).X
				symb = t.createSymbol(font, curX, curY, c)
			}
			symb.Bounds.Width, symb.Angle, symb.Value = charSize, 0, char
			symb.Color = getOrDefault(curValues, "color", t.Tint).(uint)
			symb.BackColor = getOrDefault(curValues, "backColor", uint(0)).(uint)
			symb.OutlineColor = getOrDefault(curValues, "outlineColor", palette.Black).(uint)
			symb.ShadowColor = getOrDefault(curValues, "shadowColor", palette.Black).(uint)
			symb.Weight = getOrDefault(curValues, "weight", byte(1)).(byte)
			symb.OutlineWeight = getOrDefault(curValues, "outlineWeight", byte(1)).(byte)
			symb.ShadowWeight = getOrDefault(curValues, "shadowWeight", byte(1)).(byte)
			symb.ShadowBlur = getOrDefault(curValues, "shadowBlur", byte(1)).(byte)
			symb.Underline = getOrDefault(curValues, "_", false).(bool)
			symb.Strikethrough = getOrDefault(curValues, "-", false).(bool)

			if !t.cropSymbol(symb, gapX, gapY) {
				result = append(result, symb)
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

	t.cache = state
	t.cacheChars = resultLines
	t.cacheSymbols = result
	return resultLines, result
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

func (t *TextBox) readTag(reading *bool, char rune, cur *text.Builder, curValues map[string]any) (nextChar bool) {
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
		val = condition.If(value == "regular", 1, val)
		val = condition.If(value == "semiBold", 2, val)
		val = condition.If(value == "bold", 3, val)
		curValues[name] = byte(val)
	default:
		curValues[name] = value
	}

	return false
}

//=================================================================

func parseCol(value string, defaultValue uint) uint {
	if value == "" {
		return defaultValue
	}

	var rgba = text.Split(value, " ")
	if len(rgba) == 4 {
		var r, g = text.ToNumber[byte](rgba[0]), text.ToNumber[byte](rgba[1])
		var b, a = text.ToNumber[byte](rgba[2]), text.ToNumber[byte](rgba[3])
		return color.RGBA(r, g, b, a)
	}
	return defaultValue
}
func parseNum(value string, defaultValue float32) float32 {
	if value == "" {
		return defaultValue
	}

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
func getOrDefault(curValues map[string]any, name string, defaultValue any) any {
	var val, has = curValues[name]
	if has {
		return val
	}
	return defaultValue
}
