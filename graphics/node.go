package graphics

import (
	"pure-game-kit/internal"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/number"
)

type Node struct {
	X, Y, Angle    float32
	Width, Height  float32
	ScaleX, ScaleY float32
	PivotX, PivotY float32
	Tint           uint

	renderId int32
}

func NewNode(x, y float32) *Node {
	return &Node{X: x, Y: y, Width: 100, Height: 100, ScaleX: 1, ScaleY: 1,
		PivotX: 0.5, PivotY: 0.5, Tint: palette.White, renderId: -1}
}

//=================================================================

func (n *Node) CameraFit(camera *Camera) {
	var x, y = camera.PointFromScreen(camera.ScreenX+camera.ScreenWidth/2, camera.ScreenY+camera.ScreenHeight/2)
	var cw, ch = camera.Size()
	var scale = min(cw/n.Width, ch/n.Height)

	n.X = x - (0.5-n.PivotX)*n.Width*scale
	n.Y = y - (0.5-n.PivotY)*n.Height*scale
	n.ScaleX, n.ScaleY = scale, scale
	n.Angle = 0
}
func (n *Node) CameraFill(camera *Camera) {
	var x, y = camera.PointFromScreen(camera.ScreenX+camera.ScreenWidth/2, camera.ScreenY+camera.ScreenHeight/2)
	var cw, ch = camera.Size()
	var scale = max(cw/n.Width, ch/n.Height)

	n.X = x - (0.5-n.PivotX)*n.Width*scale
	n.Y = y - (0.5-n.PivotY)*n.Height*scale
	n.ScaleX, n.ScaleY = scale, scale
	n.Angle = 0
}
func (n *Node) CameraStretch(camera *Camera) {
	var x, y = camera.PointFromScreen(camera.ScreenX+camera.ScreenWidth/2, camera.ScreenY+camera.ScreenHeight/2)
	var cw, ch = camera.Size()
	var scaleX, scaleY = cw / n.Width, ch / n.Height

	n.X = x - (0.5-n.PivotX)*n.Width*scaleX
	n.Y = y - (0.5-n.PivotY)*n.Height*scaleY
	n.ScaleX, n.ScaleY = scaleX, scaleY
	n.Angle = 0
}

//=================================================================

func (n *Node) Area() (x, y, width, height float32) {
	var x1, y1 = n.CornerTopLeft()
	var x2, y2 = n.CornerTopRight()
	var x3, y3 = n.CornerBottomRight()
	var x4, y4 = n.CornerBottomLeft()
	var minX = number.Smallest(x1, x2, x3, x4)
	var minY = number.Smallest(y1, y2, y3, y4)
	var maxX = number.Biggest(x1, x2, x3, x4)
	var maxY = number.Biggest(y1, y2, y3, y4)
	return minX, minY, maxX - minX, maxY - minY
}
func (n *Node) PointToLocal(x, y float32) (localX, localY float32) {
	if n.ScaleX == 0 || n.ScaleY == 0 {
		return 0, 0
	}

	var dx, dy = x - n.X, y - n.Y
	var sinL, cosL = internal.SinCos(-n.Angle)
	var rotX = (dx*cosL - dy*sinL) / n.ScaleX
	var rotY = (dx*sinL + dy*cosL) / n.ScaleY
	localX = rotX + n.PivotX*n.Width
	localY = rotY + n.PivotY*n.Height
	return localX, localY
}
func (n *Node) PointToGlobal(localX, localY float32) (x, y float32) {
	var locX = (localX - (n.PivotX * n.Width)) * n.ScaleX
	var locY = (localY - (n.PivotY * n.Height)) * n.ScaleY
	var sinL, cosL = internal.SinCos(n.Angle)
	x = (locX*cosL - locY*sinL) + n.X
	y = (locX*sinL + locY*cosL) + n.Y
	return x, y
}

func (n *Node) ContainsPoint(cx, cy float32) bool {
	var x, y = n.PointToLocal(cx, cy)
	var w, h = n.Width, n.Height
	return x >= 0 && y >= 0 && x < w && y < h
}

func (n *Node) CornerTopLeft() (x, y float32)     { return n.PointToGlobal(0, 0) }
func (n *Node) CornerTopRight() (x, y float32)    { return n.PointToGlobal(n.Width, 0) }
func (n *Node) CornerBottomRight() (x, y float32) { return n.PointToGlobal(n.Width, n.Height) }
func (n *Node) CornerBottomLeft() (x, y float32)  { return n.PointToGlobal(0, n.Height) }
