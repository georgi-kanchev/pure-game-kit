package graphics

import (
	"pure-game-kit/packages/assets"
	geometry "pure-game-kit/packages/geometry"
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/collection"
	col "pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/number"
	txt "pure-game-kit/packages/utility/text"
)

type Effects internal.Effects

type Object struct {
	geometry.Shape

	Mask    geometry.Area // In view space.
	Effects Effects

	// sprite =========================================================

	ImageId   assets.ImageId
	ImageCrop geometry.Area // Zero value = original asset image (or asset crop)

	// text ===========================================================

	// Tags for embeded effects:
	//
	//	Underline:     ✅ // toggle
	//	Crossout:      ❎ // toggle
	//	Weight:        ⏬🔽🔁🔼⏫
	//	Size:          🔇🔈🔉🔊📢
	//	Color:         ⬜⬛🟥🟧🟨🟩🟦🟪🟫
	//	Outline Color: ⚪⚫🔴🟠🟡🟢🔵🟣🟤
	//	Darken:        🌔🌓🌒🌑 // 20%, 40%, 60%, 80% - use before any color to darken it
	//	Brighten:      🌘🌗🌖🌕 // 20%, 40%, 60%, 80% - use before any color to brighten it
	Text       string
	TextFontId assets.FontId

	textBatches           []*internal.Batch
	textCursorPos         []float32
	textWidth, textHeight float32

	// tilemap ========================================================

	TileLayerId assets.TileLayerId
}

func NewShapePoint(x, y float32) Object {
	var eff = Effects(internal.DefaultEffects)
	eff.BorderSize, eff.FillColor = 5, palette.LightGray
	return Object{Shape: geometry.NewPoint(x, y), Effects: eff}
}
func NewShapeCircle(x, y, radius float32) Object {
	var eff = Effects(internal.DefaultEffects)
	eff.BorderSize, eff.FillColor = 5, palette.LightGray
	return Object{Shape: geometry.NewCircle(x, y, radius), Effects: eff}
}
func NewShapeRectangle(x, y, width, height, angle float32) Object {
	var eff = Effects(internal.DefaultEffects)
	eff.BorderSize, eff.FillColor = 5, palette.LightGray
	return Object{Shape: geometry.NewRectangle(x, y, width, height, angle), Effects: eff}
}
func NewShapeRoundedRectangle(x, y, width, height, angle, roundness float32) Object {
	var eff = Effects(internal.DefaultEffects)
	eff.BorderSize, eff.FillColor = 5, palette.LightGray
	return Object{Shape: geometry.NewRoundedRectangle(x, y, width, height, angle, roundness), Effects: eff}
}
func NewShapeCapsule(x1, y1, x2, y2, radius float32) Object {
	var eff = Effects(internal.DefaultEffects)
	eff.BorderSize, eff.FillColor = 5, palette.LightGray
	return Object{Shape: geometry.NewCapsule(x1, y1, x2, y2, radius), Effects: eff}
}
func NewShapeLine(x1, y1, x2, y2, thickness float32) Object {
	var eff = Effects(internal.DefaultEffects)
	eff.BorderSize, eff.FillColor = 5, palette.LightGray
	return Object{Shape: geometry.NewLine(x1, y1, x2, y2, thickness), Effects: eff}
}

func NewSprite(x, y, scale float32, imageId assets.ImageId) Object {
	var _, _, w, h = imageId.CropArea()
	var eff = Effects(internal.DefaultEffects)
	return Object{Shape: geometry.NewRectangle(x, y, float32(w)*scale, float32(h)*scale, 0), ImageId: imageId, Effects: eff}
}
func NewTextbox(x, y, width, height float32, fontId assets.FontId, text ...any) Object {
	var rect = geometry.NewRectangle(x, y, width, height, 0)
	var eff = Effects(internal.DefaultEffects)
	eff.FillColor = palette.DarkGray
	return Object{Shape: rect, TextFontId: fontId, Text: txt.New(text...), Effects: eff}
}
func NewTilemap(scale float32, layerId assets.TileLayerId) Object {
	var tilemap = Object{TileLayerId: layerId, Effects: Effects(internal.DefaultEffects)}
	var layer = internal.TileLayers[uint8(layerId)]
	if layer != nil {
		tilemap.Width = float32(layer.Columns) * float32(layer.TileSize) * scale
		tilemap.Height = float32(layer.Rows) * float32(layer.TileSize) * scale
	}
	return tilemap
}

//=================================================================

func (o *Object) ViewFit(view *View) {
	var x, y = view.PointFromScreen(internal.WindowWidth/2, internal.WindowHeight/2)
	var cw, ch = view.Size()
	var scale = min(cw/o.Width, ch/o.Height)
	o.Width, o.Height = o.Width*scale, o.Height*scale
	o.X, o.Y, o.Angle = x, y, 0
}
func (o *Object) ViewFill(view *View) {
	var x, y = view.PointFromScreen(internal.WindowWidth/2, internal.WindowHeight/2)
	var cw, ch = view.Size()
	var scale = max(cw/o.Width, ch/o.Height)
	o.Width, o.Height = o.Width*scale, o.Height*scale
	o.X, o.Y, o.Angle = x, y, 0
}
func (o *Object) ViewStretch(view *View) {
	var x, y = view.PointFromScreen(internal.WindowWidth/2, internal.WindowHeight/2)
	var cw, ch = view.Size()
	var scaleX, scaleY = cw / o.Width, ch / o.Height
	o.Width, o.Height = o.Width*scaleX, o.Height*scaleY
	o.X, o.Y, o.Angle = x, y, 0
}

//=================================================================

func (o *Object) PointToLocal(x, y float32) (localX, localY float32) {
	var dx, dy = x - o.X, y - o.Y
	var sinL, cosL = internal.SinCos(-o.Angle)
	return (dx*cosL - dy*sinL) + 0.5*o.Width, (dx*sinL + dy*cosL) + 0.5*o.Height
}
func (o *Object) PointToGlobal(localX, localY float32) (x, y float32) {
	var locX, locY = localX - (0.5 * o.Width), localY - (0.5 * o.Height)
	var sinL, cosL = internal.SinCos(o.Angle)
	return (locX*cosL - locY*sinL) + o.X, (locX*sinL + locY*cosL) + o.Y
}
func (o *Object) ContainsPoint(x, y float32) bool {
	var lx, ly = o.PointToLocal(x, y)
	return lx >= 0 && ly >= 0 && lx < o.Width && ly < o.Height
}
func (o *Object) PointFromEdge(edgeX, edgeY float32) (x, y float32) {
	return o.PointToGlobal(o.Width*edgeX, o.Height*edgeY)
}

// text ===========================================================

// Works only when TextIsBatched is true. Needs to be called when the textbox visuals require a change.
// The update happens upon draw so calling this multiple times per frame is fine.
func (o *Object) TextUpdateBatch() {
	o.textBatches = nil
}
func (o *Object) TextCursorPositionAt(index int) float32 {
	if !o.Effects.TextIsInput || index < 0 || index >= len(o.textCursorPos) {
		return number.NaN()
	}
	return o.textCursorPos[index]
}
func (o *Object) TextSize() (width, height float32) {
	return o.textWidth, o.textHeight
}

// tilemap ========================================================

func (o *Object) TilemapShapes() []geometry.Shape {
	var result = []geometry.Shape{}
	var layer = internal.TileLayers[uint8(o.TileLayerId)]
	if layer == nil {
		return result
	}

	if layer.Image == nil || layer.Texture.Width == 0 { // is object layer
		for _, s := range layer.Objects {
			var x, y = o.PointToGlobal(s[0], s[1])
			var rect = geometry.NewRoundedRectangle(x, y, s[2], s[3], o.Angle+s[4], s[5])
			result = collection.Add(result, rect)
		}
		return result
	}

	var w, h = o.TileLayerId.Size()
	for cellIndex1D := range layer.CellsWithPoints {
		var row, column = number.Index1DToIndexes2D(cellIndex1D, w, h)
		result = collection.Add(result, o.TilemapShapesAtCell(column, row)...)
	}
	return result
}
func (o *Object) TilemapShapesAtCell(column, row int) []geometry.Shape {
	var result = []geometry.Shape{}
	var tile = o.TileLayerId.TileAtCell(column, row)
	if tile.Id == 0 {
		return result
	}
	var tw, th = o.TileLayerId.TileSize()

	result = collection.Add(result, o.TilemapShapesFromTile(tile.Id)...)
	for i, s := range result {
		var relX, relY = s.X - tw/2, s.Y - th/2
		switch tile.Rotations90 {
		case 1:
			relX, relY, s.Angle = -relY, relX, s.Angle+90
		case 2:
			relX, relY, s.Angle = -relX, -relY, s.Angle+180
		case 3:
			relX, relY, s.Angle = relY, -relX, s.Angle+270
		}
		if tile.Flip {
			relX, s.Angle = -relX, -s.Angle
		}

		s.X, s.Y = relX+tw/2, relY+th/2
		s.X, s.Y = o.PointToGlobal(s.X+float32(column)*tw, s.Y+float32(row)*th)
		result[i] = s
	}
	return result
}
func (o *Object) TilemapShapesFromTile(tileId uint16) []geometry.Shape {
	var result = []geometry.Shape{}
	var layer = internal.TileLayers[uint8(o.TileLayerId)]
	if layer == nil {
		return result
	}
	var shapes, has = layer.ShapesPerTile[tileId]
	if !has {
		return result
	}
	for _, s := range shapes {
		var rect = geometry.NewRoundedRectangle(s[0], s[1], s[2], s[3], o.Angle+s[4], s[5])
		result = collection.Add(result, rect)
	}
	return result
}
func (o *Object) TilemapPaths() []float32 {
	var result = []float32{}
	var layer = internal.TileLayers[uint8(o.TileLayerId)]
	if layer == nil || len(layer.Paths) == 0 {
		return result
	}

	for i := 0; i < len(layer.Paths)-1; i += 2 {
		var px, py = layer.Paths[i], layer.Paths[i+1]
		if number.IsNaN(px) || number.IsNaN(py) {
			result = collection.Add(result, px, py)
			continue
		}
		var x, y = o.PointToGlobal(px, py)
		result = collection.Add(result, x, y)
	}
	return result
}

// private ========================================================

func (o *Object) measureLine(fromIndex int, lineHeight float32) (endIndex int, width float32, endLineHeight float32) {
	if fromIndex >= len(o.Text) {
		return fromIndex, 0, lineHeight
	}

	var originalLineHeight = o.Effects.TextLineHeight
	var scale = originalLineHeight / 255
	var gapX = o.Effects.TextSymbolGap * scale
	var font = internal.Fonts[uint8(o.TextFontId)]
	var x, totalWidth float32
	var prevGlyph internal.Glyph

	for i, r := range o.Text[fromIndex:] {
		if !o.Effects.TextIsInput && r == '\n' {
			return fromIndex + i, totalWidth, lineHeight
		}

		var sz = sizes[r]
		if sz != 0 {
			lineHeight = originalLineHeight * sz
			continue
		}

		x += prevGlyph.Kernings[r] * lineHeight
		var glyph = font.Chars[r]

		if o.Effects.TextWordWrap && r == ' ' {
			var wX, wTotal float32
			var wPrev internal.Glyph
			var wHeight = lineHeight
			for _, wr := range o.Text[fromIndex+i+1:] {
				if wr == ' ' || (!o.Effects.TextIsInput && wr == '\n') {
					break
				}
				var wsz = sizes[wr]
				if wsz != 0 {
					wHeight = originalLineHeight * wsz
					continue
				}
				var symbol = o.TextFontId.SymbolArea(wr, wHeight)
				var wGlyph = font.Chars[wr]
				wX += wPrev.Kernings[wr] * wHeight
				wPrev, wTotal = wGlyph, max(wX+symbol.X+symbol.Width, wTotal)
				wX += wGlyph.Advance*wHeight + gapX
			}
			if x+glyph.Advance*lineHeight+gapX+max(wTotal, wX) > o.Width {
				return fromIndex + i, totalWidth, lineHeight
			}
		}
		x, prevGlyph, totalWidth = x+(glyph.Advance*lineHeight+gapX), glyph, max(x+glyph.Advance*lineHeight, totalWidth)
	}
	return len(o.Text), totalWidth, lineHeight
}
func (o *Object) embedEffect(r rune, effect *internal.Effects, shadeCol, shadeOutCol *float32, baseLineHeight float32) (success bool) {
	if r == '✅' {
		effect.TextUnderline = !effect.TextUnderline
		return true
	}
	if r == '❎' {
		effect.TextCrossout = !effect.TextCrossout
		return true
	}

	var color = colors[r]
	var outlineColor = outlineColors[r]
	var weight, hasWeights = weights[r]
	var size = sizes[r]
	var sh, hasShade = shades[r]
	if color != 0 {
		if *shadeCol < 0 {
			color = col.Darken(color, number.Absolute(*shadeCol))
		} else {
			color = col.Brighten(color, *shadeCol)
		}
		*shadeCol = 0    // reset shade
		*shadeOutCol = 0 // reset shade
		effect.TextColor = color
		return true
	} else if outlineColor != 0 {
		if *shadeOutCol < 0 {
			outlineColor = col.Darken(outlineColor, number.Absolute(*shadeOutCol))
		} else {
			outlineColor = col.Brighten(outlineColor, *shadeOutCol)
		}
		*shadeCol = 0    // reset shade
		*shadeOutCol = 0 // reset shade
		effect.OutlineColor = outlineColor
		return true
	} else if hasWeights {
		effect.TextWeight = weight
		return true
	} else if size != 0 {
		effect.TextLineHeight = baseLineHeight * size
		return true
	} else if hasShade {
		*shadeCol = sh
		*shadeOutCol = sh
	}
	return false
}
