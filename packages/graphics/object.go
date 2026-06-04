package graphics

import (
	"pure-game-kit/packages/assets"
	geometry "pure-game-kit/packages/geometry2"
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/color/palette"
	txt "pure-game-kit/packages/utility/text"
)

type Object struct {
	geometry.Shape

	Mask    Area
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
func NewTilemap(atlasImageId assets.ImageId, tileLayerId assets.TileLayerId) Object {
	return Object{Shape: geometry.NewRectangle(0, 0, 100, 100, 0), TileLayerId: tileLayerId, Effects: Effects(internal.DefaultEffects)}
}

//=================================================================

func (o *Object) ViewFit(view *View) {
	var sx, sy, sw, sh = view.area()
	var x, y = view.PointFromScreen(sx+sw/2, sy+sh/2)
	var cw, ch = view.Size()
	var scale = min(cw/o.Width, ch/o.Height)

	o.X = x - (0.5)*o.Width*scale
	o.Y = y - (0.5)*o.Height*scale
	o.Angle = 0
}
func (o *Object) ViewFill(view *View) {
	var sx, sy, sw, sh = view.area()
	var x, y = view.PointFromScreen(sx+sw/2, sy+sh/2)
	var cw, ch = view.Size()
	var scale = max(cw/o.Width, ch/o.Height)

	o.X = x - (0.5)*o.Width*scale
	o.Y = y - (0.5)*o.Height*scale
	o.Angle = 0
}
func (o *Object) ViewStretch(view *View) {
	var sx, sy, sw, sh = view.area()
	var x, y = view.PointFromScreen(sx+sw/2, sy+sh/2)
	var cw, ch = view.Size()
	var scaleX, scaleY = cw / o.Width, ch / o.Height

	o.X = x - (0.5)*o.Width*scaleX
	o.Y = y - (0.5)*o.Height*scaleY
	o.Angle = 0
}

//=================================================================

func (o *Object) PointToLocal(x, y float32) (localX, localY float32) {
	var dx, dy = x - o.X, y - o.Y
	var sinL, cosL = internal.SinCos(-o.Angle)
	var rotX = (dx*cosL - dy*sinL)
	var rotY = (dx*sinL + dy*cosL)
	localX = rotX + 0.5*o.Width
	localY = rotY + 0.5*o.Height
	return localX, localY
}
func (o *Object) PointToGlobal(localX, localY float32) (x, y float32) {
	var locX = (localX - (0.5 * o.Width))
	var locY = (localY - (0.5 * o.Height))
	var sinL, cosL = internal.SinCos(o.Angle)
	x = (locX*cosL - locY*sinL) + o.X
	y = (locX*sinL + locY*cosL) + o.Y
	return x, y
}
func (o *Object) ContainsPoint(x, y float32) bool {
	var lx, ly = o.PointToLocal(x, y)
	return lx >= 0 && ly >= 0 && lx < o.Width && ly < o.Height
}
func (o *Object) PointFromEdge(edgeX, edgeY float32) (x, y float32) {
	return o.PointToGlobal(o.Width*edgeX, o.Height*edgeY)
}

//=================================================================

// Works only when TextIsBatched is true. Needs to be called when the textbox visuals require a change.
// The update happens upon draw so calling this multiple times per frame is fine.
func (o *Object) TextUpdateBatch() {
	o.textBatches = nil
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

		if r == '\t' {
			var nextX = float32(int(x/originalLineHeight*2)+1) * (originalLineHeight * 2)
			if nextX > o.Width {
				return fromIndex + i, max(totalWidth, x), lineHeight
			}
			x, totalWidth, prevGlyph = nextX, max(totalWidth, nextX), internal.Glyph{}
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
