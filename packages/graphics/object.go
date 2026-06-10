package graphics

import (
	"pure-game-kit/packages/assets"
	geometry "pure-game-kit/packages/geometry"
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/number"
	txt "pure-game-kit/packages/utility/text"
)

type Object struct {
	geometry.Shape

	Mask    Area // In window space.
	Effects Effects

	// image ==========================================================

	ImageId   assets.ImageId
	ImageCrop Area // Zero value = original asset image (or asset crop)

	// text ===========================================================

	// Tags for embeded effects:
	//
	//	Underline:     ✅ // toggle
	//	Crossout:      ❎ // toggle
	//	Weight:        ⏬🔽🔁🔼⏫
	//	Size:          🔇🔈🔉🔊📢
	//	Color:         ⬜⬛🟥🟧🟨🟩🟦🟪🟫
	//	Outline Color: ⚪⚫🔴🟠🟡🟢🔵🟣🟤
	Text       string
	TextFontId assets.FontId

	// Caches the text visuals across frames. Call TextUpdateBatch when visual changes are needed.
	// Useful for a huge static text that changes rarely.
	TextBatch   bool
	textBatches []*internal.Batch

	// tilemap ========================================================

	TileLayerId assets.TileLayerId

	// nine-patch =====================================================

	NinePatchId assets.NinePatchId
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

func NewImage(x, y, scale float32, imageId assets.ImageId) Object {
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
func NewNinePatch(x, y, width, height float32, ninePatchId assets.NinePatchId) Object {
	var eff = Effects(internal.DefaultEffects)
	return Object{Shape: geometry.NewRectangle(x, y, width, height, 0), NinePatchId: ninePatchId, Effects: eff}
}

//=================================================================

func (o *Object) ViewFit(view *View) {
	var x, y = view.PointFromScreen(internal.WindowWidth/2, internal.WindowHeight/2)
	var cw, ch = view.Size()
	var scale = min(cw/o.Width, ch/o.Height)
	o.X, o.Y, o.Angle = x-(0.5)*o.Width*scale, y-(0.5)*o.Height*scale, 0
}
func (o *Object) ViewFill(view *View) {
	var x, y = view.PointFromScreen(internal.WindowWidth/2, internal.WindowHeight/2)
	var cw, ch = view.Size()
	var scale = max(cw/o.Width, ch/o.Height)
	o.X, o.Y, o.Angle = x-(0.5)*o.Width*scale, y-(0.5)*o.Height*scale, 0
}
func (o *Object) ViewStretch(view *View) {
	var x, y = view.PointFromScreen(internal.WindowWidth/2, internal.WindowHeight/2)
	var cw, ch = view.Size()
	var scaleX, scaleY = cw / o.Width, ch / o.Height
	o.X, o.Y, o.Angle = x-(0.5)*o.Width*scaleX, y-(0.5)*o.Height*scaleY, 0
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

// tilemap ========================================================

func (o *Object) TilemapShapes() []geometry.Shape {
	var layer = internal.TileLayers[uint8(o.TileLayerId)]
	if layer == nil {
		return nil
	}

	if layer.Image == nil || layer.Texture.Width == 0 { // is object layer
		var result = make([]geometry.Shape, len(layer.Objects))
		for i, s := range layer.Objects {
			var x, y = o.PointToGlobal(s[0], s[1])
			result[i] = geometry.Shape{X: x, Y: y, Width: s[2], Height: s[3], Angle: o.Angle + s[4], Roundness: s[5]}
		}
		return result
	}

	var w, h = o.TileLayerId.Size()
	var result []geometry.Shape
	for cellIndex1D := range layer.CellsWithPoints {
		var row, column = number.Index1DToIndexes2D(cellIndex1D, w, h)
		result = append(result, o.TilemapShapesAtCell(column, row)...)
	}
	return result
}
func (o *Object) TilemapShapesAtCell(column, row int) []geometry.Shape {
	var tile = o.TileLayerId.TileAtCell(column, row)
	if tile.Id == 0 {
		return nil
	}
	var tw, th = o.TileLayerId.TileSize()
	var shapes = o.TilemapShapesFromTile(tile.Id)
	for i := range shapes {
		var relX, relY = shapes[i].X - tw/2, shapes[i].Y - th/2
		switch tile.Rotations90 {
		case 1:
			relX, relY, shapes[i].Angle = -relY, relX, shapes[i].Angle+90
		case 2:
			relX, relY, shapes[i].Angle = -relX, -relY, shapes[i].Angle+180
		case 3:
			relX, relY, shapes[i].Angle = relY, -relX, shapes[i].Angle+270
		}
		if tile.Flip {
			relX, shapes[i].Angle = -relX, -shapes[i].Angle
		}
		shapes[i].X, shapes[i].Y = relX+tw/2, relY+th/2
		shapes[i].X, shapes[i].Y = o.PointToGlobal(shapes[i].X+float32(column)*tw, shapes[i].Y+float32(row)*th)
	}
	return shapes
}
func (o *Object) TilemapShapesFromTile(tileId uint16) []geometry.Shape {
	var layer = internal.TileLayers[uint8(o.TileLayerId)]
	if layer == nil {
		return nil
	}
	var shapes, has = layer.ShapesPerTile[tileId]
	if !has {
		return nil
	}
	var result = make([]geometry.Shape, len(shapes))
	for i, s := range shapes {
		result[i] = geometry.Shape{X: s[0], Y: s[1], Width: s[2], Height: s[3], Angle: o.Angle + s[4], Roundness: s[5]}
	}
	return result
}
func (o *Object) TilemapPaths() []float32 {
	var layer = internal.TileLayers[uint8(o.TileLayerId)]
	if layer == nil || len(layer.Paths) == 0 {
		return nil
	}
	var result = make([]float32, len(layer.Paths))
	for i := 0; i < len(layer.Paths); i += 2 {
		var px, py = layer.Paths[i], layer.Paths[i+1]
		if number.IsNaN(px) || number.IsNaN(py) {
			result[i], result[i+1] = px, py
			continue
		}
		result[i], result[i+1] = o.PointToGlobal(px, py)
	}
	return result
}

// private ========================================================

func (o *Object) measureLine(fromIndex int, lineHeight float32) (endIndex int, width float32, endLineHeight float32) {
	if fromIndex >= len(o.Text) {
		return fromIndex, 0, lineHeight
	}

	var originalLineHeight = o.Effects.TextLineHeight
	var gapX = o.Effects.TextSymbolGap * originalLineHeight / 255
	var font = internal.Fonts[uint8(o.TextFontId)]
	var x, totalWidth float32
	var prevGlyph internal.Glyph

	for i, r := range o.Text[fromIndex:] {
		if r == '\n' {
			return fromIndex + i, max(totalWidth, x), lineHeight
		}

		var sz = sizes[r]
		if sz != 0 {
			lineHeight = originalLineHeight * sz
			continue
		}

		x += prevGlyph.Kernings[r] * lineHeight
		var offsetX, _, w, _ = o.TextFontId.SymbolArea(r, lineHeight)
		var glyph = font.Chars[r]

		if o.Effects.TextWordWrap && r == ' ' {
			var wX, wTotal float32
			var wPrev internal.Glyph
			var wHeight = lineHeight
			for _, wr := range o.Text[fromIndex+i+1:] {
				if wr == ' ' || wr == '\n' {
					break
				}
				var wsz = sizes[wr]
				if wsz != 0 {
					wHeight = originalLineHeight * wsz
					continue
				}
				var wOffX, _, wW, _ = o.TextFontId.SymbolArea(wr, wHeight)
				var wGlyph = font.Chars[wr]
				wX += wPrev.Kernings[wr] * wHeight
				wPrev, wTotal = wGlyph, max(wX+wOffX+wW, wTotal)
				wX += wGlyph.Advance*wHeight + gapX
			}
			if x+glyph.Advance*lineHeight+gapX+max(wTotal, wX) > o.Width {
				return fromIndex + i, max(totalWidth, x), lineHeight
			}
		}
		x, prevGlyph, totalWidth = x+(glyph.Advance*lineHeight+gapX), glyph, max(x+offsetX+w, totalWidth)
	}
	return len(o.Text), max(totalWidth, x), lineHeight
}
