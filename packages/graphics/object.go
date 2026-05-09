package graphics

import (
	"pure-game-kit/packages/assets"
	geometry "pure-game-kit/packages/geometry2"
	"pure-game-kit/packages/internal"
)

type Object struct {
	geometry.Shape
	ScaleX, ScaleY float32
	PivotX, PivotY float32
	Color          uint

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

	// private ========================================================

	cache     textCache // tracks when to regenerate chars
	chars     []Object
	lineCount int
	charValue rune
}

func NewObject(x, y float32) Object {
	return Object{Shape: geometry.NewRectangle(x, y, 100, 100, 0), TextLineHeight: 100, TextWordWrap: true}
}

//=================================================================

func (o *Object) ViewFit(view *View) {
	var sx, sy, sw, sh = view.area()
	var x, y = view.PointFromScreen(sx+sw/2, sy+sh/2)
	var cw, ch = view.Size()
	var scale = min(cw/o.Width, ch/o.Height)

	o.X = x - (0.5-o.PivotX)*o.Width*scale
	o.Y = y - (0.5-o.PivotY)*o.Height*scale
	o.ScaleX, o.ScaleY = scale, scale
	o.Angle = 0
}
func (o *Object) ViewFill(view *View) {
	var sx, sy, sw, sh = view.area()
	var x, y = view.PointFromScreen(sx+sw/2, sy+sh/2)
	var cw, ch = view.Size()
	var scale = max(cw/o.Width, ch/o.Height)

	o.X = x - (0.5-o.PivotX)*o.Width*scale
	o.Y = y - (0.5-o.PivotY)*o.Height*scale
	o.ScaleX, o.ScaleY = scale, scale
	o.Angle = 0
}
func (o *Object) ViewStretch(view *View) {
	var sx, sy, sw, sh = view.area()
	var x, y = view.PointFromScreen(sx+sw/2, sy+sh/2)
	var cw, ch = view.Size()
	var scaleX, scaleY = cw / o.Width, ch / o.Height

	o.X = x - (0.5-o.PivotX)*o.Width*scaleX
	o.Y = y - (0.5-o.PivotY)*o.Height*scaleY
	o.ScaleX, o.ScaleY = scaleX, scaleY
	o.Angle = 0
}

//=================================================================

func (o *Object) PointToLocal(x, y float32) (localX, localY float32) {
	if o.ScaleX == 0 || o.ScaleY == 0 {
		return 0, 0
	}

	var dx, dy = x - o.X, y - o.Y
	var sinL, cosL = internal.SinCos(-o.Angle)
	var rotX = (dx*cosL - dy*sinL) / o.ScaleX
	var rotY = (dx*sinL + dy*cosL) / o.ScaleY
	localX = rotX + o.PivotX*o.Width
	localY = rotY + o.PivotY*o.Height
	return localX, localY
}
func (o *Object) PointToGlobal(localX, localY float32) (x, y float32) {
	var locX = (localX - (o.PivotX * o.Width)) * o.ScaleX
	var locY = (localY - (o.PivotY * o.Height)) * o.ScaleY
	var sinL, cosL = internal.SinCos(o.Angle)
	x = (locX*cosL - locY*sinL) + o.X
	y = (locX*sinL + locY*cosL) + o.Y
	return x, y
}
func (o *Object) ContainsPoint(cx, cy float32) bool {
	var x, y = o.PointToLocal(cx, cy)
	return x >= 0 && y >= 0 && x < o.Width && y < o.Height
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
	var fontData = internal.Fonts2[byte(o.TextFont)]
	for i, r := range o.Text {
		var symbol = Object{TextFont: assets.FontId(fontData.AtlasId), charValue: r}
		symbol.X = float32(30 * i)
		o.chars = append(o.chars, symbol)
	}
}
