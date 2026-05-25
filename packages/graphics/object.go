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

	ImageId   assets.ImageId
	ImageCrop Area // Zero value = original asset image (or crop)

	// text ===========================================================

	Text       string
	TextFontId assets.FontId

	// tilemap ========================================================

	TileLayerId assets.TileLayerId
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

func NewImage(x, y, scale float32, imageId assets.ImageId) Object {
	var _, _, w, h = imageId.CropArea()
	return Object{Shape: geometry.NewRectangle(x, y, float32(w)*scale, float32(h)*scale, 0), ImageId: imageId, Color: palette.White}
}
func NewTextbox(x, y, width, height float32, fontId assets.FontId, text ...any) Object {
	var rect = geometry.NewRectangle(x, y, width, height, 0)
	return Object{Shape: rect, TextFontId: fontId, Text: txt.New(text...)}
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

// text ===========================================================

func (o *Object) TextMeasure(text string) (width, height float32) {
	return 0, 0
}
func (o *Object) TextLineCount() int {
	return 0
}
func (o *Object) TextSymbol(index int) Object {
	return Object{}
}
