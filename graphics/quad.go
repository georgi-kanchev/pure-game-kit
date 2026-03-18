package graphics

import (
	"pure-game-kit/internal"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/number"
)

type Quad struct {
	Area
	Angle          float32
	ScaleX, ScaleY float32
	PivotX, PivotY float32
	Tint           uint
	Mask           *Area
	Effects        *Effects
}

func NewQuad(x, y float32) *Quad {
	return &Quad{Area: Area{X: x, Y: y, Width: 100, Height: 100}, ScaleX: 1, ScaleY: 1,
		PivotX: 0.5, PivotY: 0.5, Tint: palette.White}
}

//=================================================================

func (q *Quad) CameraFit(camera *Camera) {
	var sx, sy, sw, sh = camera.area()
	var x, y = camera.PointFromScreen(sx+sw/2, sy+sh/2)
	var cw, ch = camera.Size()
	var scale = min(cw/q.Width, ch/q.Height)

	q.X = x - (0.5-q.PivotX)*q.Width*scale
	q.Y = y - (0.5-q.PivotY)*q.Height*scale
	q.ScaleX, q.ScaleY = scale, scale
	q.Angle = 0
}
func (q *Quad) CameraFill(camera *Camera) {
	var sx, sy, sw, sh = camera.area()
	var x, y = camera.PointFromScreen(sx+sw/2, sy+sh/2)
	var cw, ch = camera.Size()
	var scale = max(cw/q.Width, ch/q.Height)

	q.X = x - (0.5-q.PivotX)*q.Width*scale
	q.Y = y - (0.5-q.PivotY)*q.Height*scale
	q.ScaleX, q.ScaleY = scale, scale
	q.Angle = 0
}
func (q *Quad) CameraStretch(camera *Camera) {
	var sx, sy, sw, sh = camera.area()
	var x, y = camera.PointFromScreen(sx+sw/2, sy+sh/2)
	var cw, ch = camera.Size()
	var scaleX, scaleY = cw / q.Width, ch / q.Height

	q.X = x - (0.5-q.PivotX)*q.Width*scaleX
	q.Y = y - (0.5-q.PivotY)*q.Height*scaleY
	q.ScaleX, q.ScaleY = scaleX, scaleY
	q.Angle = 0
}

//=================================================================

func (q *Quad) Bounds() (x, y, width, height float32) {
	var x1, y1 = q.CornerTopLeft()
	var x2, y2 = q.CornerTopRight()
	var x3, y3 = q.CornerBottomRight()
	var x4, y4 = q.CornerBottomLeft()
	var minX = number.Smallest(x1, x2, x3, x4)
	var minY = number.Smallest(y1, y2, y3, y4)
	var maxX = number.Biggest(x1, x2, x3, x4)
	var maxY = number.Biggest(y1, y2, y3, y4)
	return minX, minY, maxX - minX, maxY - minY
}
func (q *Quad) PointToLocal(x, y float32) (localX, localY float32) {
	if q.ScaleX == 0 || q.ScaleY == 0 {
		return 0, 0
	}

	var dx, dy = x - q.X, y - q.Y
	var sinL, cosL = internal.SinCos(-q.Angle)
	var rotX = (dx*cosL - dy*sinL) / q.ScaleX
	var rotY = (dx*sinL + dy*cosL) / q.ScaleY
	localX = rotX + q.PivotX*q.Width
	localY = rotY + q.PivotY*q.Height
	return localX, localY
}
func (q *Quad) PointToGlobal(localX, localY float32) (x, y float32) {
	var locX = (localX - (q.PivotX * q.Width)) * q.ScaleX
	var locY = (localY - (q.PivotY * q.Height)) * q.ScaleY
	var sinL, cosL = internal.SinCos(q.Angle)
	x = (locX*cosL - locY*sinL) + q.X
	y = (locX*sinL + locY*cosL) + q.Y
	return x, y
}
func (q *Quad) ContainsPoint(cx, cy float32) bool {
	var x, y = q.PointToLocal(cx, cy)
	return x >= 0 && y >= 0 && x < q.Width && y < q.Height
}

func (q *Quad) CornerTopLeft() (x, y float32)     { return q.PointToGlobal(0, 0) }
func (q *Quad) CornerTopRight() (x, y float32)    { return q.PointToGlobal(q.Width, 0) }
func (q *Quad) CornerBottomRight() (x, y float32) { return q.PointToGlobal(q.Width, q.Height) }
func (q *Quad) CornerBottomLeft() (x, y float32)  { return q.PointToGlobal(0, q.Height) }
