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
	ImageCrop Area // Zero value = original asset image (or crop)

	// text ===========================================================

	Text        string
	TextFontId  assets.FontId
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
	return Object{Shape: rect, TextFontId: fontId, Text: txt.New(text...), Effects: Effects(internal.DefaultEffects)}
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

func (o *Object) IsShape() bool {
	return o.ImageId == 0 && o.TextFontId == 0 && o.TileLayerId == 0
}
func (o *Object) IsSprite() bool {
	return o.ImageId != 0
}
func (o *Object) IsTextbox() bool {
	return o.TextFontId != 0
}
func (o *Object) IsTilemap() bool {
	return o.TileLayerId != 0
}

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

func (o *Object) TextRecalculate() {
	o.textBatches = nil
}

// private ========================================================

func (o *Object) lineEndAndWidth(fromIndex int) (endIndex int, width float32) {
	if fromIndex >= len(o.Text) {
		return fromIndex, 0
	}

	var lineHeight = o.Effects.TextLineHeight
	var scale = lineHeight / 255
	var gapX = o.Effects.TextSymbolGap * scale
	var font = internal.Fonts[uint8(o.TextFontId)]
	var x, totalWidth float32
	var prevGlyph internal.Glyph

	for i, r := range o.Text[fromIndex:] {
		if r == '\n' {
			return fromIndex + i, max(totalWidth, x)
		}

		x += prevGlyph.Kernings[r] * lineHeight
		var offsetX, _, w, _ = o.TextFontId.SymbolArea(r, lineHeight)
		var glyph = font.Chars[r]

		if o.Effects.TextWordWrap && r == ' ' {
			var wX, wTotal float32
			var wPrev internal.Glyph
			for _, wr := range o.Text[fromIndex+i+1:] {
				if wr == ' ' || wr == '\n' {
					break
				}
				var wOffX, _, wW, _ = o.TextFontId.SymbolArea(wr, lineHeight)
				var wGlyph = font.Chars[wr]
				wX += wPrev.Kernings[wr] * lineHeight
				wPrev, wTotal = wGlyph, max(wX+wOffX+wW, wTotal)
				wX += wGlyph.Advance*lineHeight + gapX
			}
			if x+glyph.Advance*lineHeight+gapX+max(wTotal, wX) > o.Width {
				return fromIndex + i, max(totalWidth, x)
			}
		}
		x, prevGlyph, totalWidth = x+(glyph.Advance*lineHeight+gapX), glyph, max(x+offsetX+w, totalWidth)
	}
	return len(o.Text), max(totalWidth, x)
}
