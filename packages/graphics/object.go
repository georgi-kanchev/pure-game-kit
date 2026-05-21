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
	Color uint

	Mask    Area
	Effects *Effects

	// image ==========================================================

	ImageId       assets.ImageId
	ImageCropArea Area // Zero value = entire image

	// text ===========================================================

	Text     string
	TextFont assets.FontId

	TextAlignX, TextAlignY, TextLineHeight, TextSymbolGap, TextLineGap float32

	TextWordWrap, TextUnderline, TextCrossout bool

	TextWeight, TextShadowSize, TextShadowBlur,
	TextShadowOffsetX, TextShadowOffsetY float32

	TextBackColor, TextShadowColor uint

	// tilemap ========================================================

	TileLayerId assets.TileLayerId

	// private ========================================================

	cache     textCache // tracks when to regenerate chars
	chars     []Object
	lineCount int
	charValue rune
}

func NewShapePoint(x, y float32, color uint) Object {
	return Object{Shape: geometry.NewPoint(x, y), Color: color}
}
func NewShapeCircle(x, y, radius float32, color uint) Object {
	return Object{Shape: geometry.NewCircle(x, y, radius), Color: color}
}
func NewShapeRectangle(x, y, width, height, angle float32, color uint) Object {
	return Object{Shape: geometry.NewRectangle(x, y, width, height, angle), Color: color}
}
func NewShapeRoundedRectangle(x, y, width, height, angle, roundness float32, color uint) Object {
	return Object{Shape: geometry.NewRoundedRectangle(x, y, width, height, angle, roundness), Color: color}
}
func NewShapeCapsule(x1, y1, x2, y2, radius float32, color uint) Object {
	return Object{Shape: geometry.NewCapsule(x1, y1, x2, y2, radius), Color: color}
}
func NewShapeLine(x1, y1, x2, y2, thickness float32, color uint) Object {
	return Object{Shape: geometry.NewLine(x1, y1, x2, y2, thickness), Color: color}
}

func NewImage(x, y float32, imageId assets.ImageId) Object {
	var w, h = imageId.Size()
	return Object{Shape: geometry.NewRectangle(x, y, float32(w), float32(h), 0), ImageId: imageId, Color: palette.White}
}
func NewTextbox(x, y, width, height float32, fontId assets.FontId, text ...any) Object {
	var rect = geometry.NewRectangle(x, y, width, height, 0)
	return Object{
		Shape: rect, TextFont: fontId, Text: txt.New(text), Color: palette.White, TextLineHeight: 100, TextWordWrap: true}
}
func NewTilemap(atlasImageId assets.ImageId, tileLayerId assets.TileLayerId) Object {
	return Object{Shape: geometry.NewRectangle(0, 0, 100, 100, 0), TileLayerId: tileLayerId, Color: palette.White}
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
	return o.ImageId == 0 && o.TextFont == 0 && o.TileLayerId == 0
}
func (o *Object) IsSprite() bool {
	return o.ImageId != 0
}
func (o *Object) IsTextbox() bool {
	return o.TextFont != 0
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

// text ===========================================================

func (o *Object) TextMeasure(text string) (width, height float32) {
	o.tryRegenerateText()
	return 0, 0
}
func (o *Object) TextLineCount() int {
	o.tryRegenerateText()
	return o.lineCount
}
func (o *Object) TextSymbol(index int) Object {
	o.tryRegenerateText()
	return Object{}
}

// private ========================================================

type textCache struct {
	text string
	font assets.FontId
	width, height, alignX, alignY,
	lineHeight, symbolGap, lineGap float32
	wordWrap bool
}

func (o *Object) tryRegenerateText() {
	var w, h, ax, ay = o.Width, o.Height, o.TextAlignX, o.TextAlignY
	var state = textCache{o.Text, o.TextFont, w, h, ax, ay, o.TextLineHeight, o.TextSymbolGap, o.TextLineGap, o.TextWordWrap}
	if state == o.cache {
		return
	}
	o.cache = state

	o.chars = o.chars[:]
	var fontData = internal.Fonts[byte(o.TextFont)]
	for i, r := range o.Text {
		var symbol = Object{TextFont: assets.FontId(fontData.AtlasId), charValue: r}
		symbol.X = float32(30 * i)
		o.chars = append(o.chars, symbol)
	}
}

func (o *Object) setDefaults() {
	o.TextLineHeight = 100
	o.TextWordWrap = true
}
