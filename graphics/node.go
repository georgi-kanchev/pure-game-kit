package graphics

import (
	"pure-game-kit/utility/angle"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/point"
)

type Node struct {
	X, Y, Angle    float32
	Width, Height  float32
	ScaleX, ScaleY float32
	PivotX, PivotY float32
	Parent         *Node
	Color          uint
}

func NewNode(x, y float32) *Node {
	return &Node{X: x, Y: y, Width: 100, Height: 100, ScaleX: 1, ScaleY: 1,
		PivotX: 0.5, PivotY: 0.5, Color: palette.White}
}

//=================================================================

func (n *Node) CameraFit(camera *Camera) {
	var x, y = camera.PointFromScreen(
		camera.ScreenX+camera.ScreenWidth/2,
		camera.ScreenY+camera.ScreenHeight/2,
	)
	var w, h = n.Width, n.Height
	var cw, ch = camera.Size()
	var scale = min(cw/w, ch/h)

	n.X = x - (0.5-n.PivotX)*w*scale
	n.Y = y - (0.5-n.PivotY)*h*scale
	n.ScaleX, n.ScaleY = scale, scale
	n.Angle = 0
}
func (n *Node) CameraFill(camera *Camera) {
	var x, y = camera.PointFromScreen(
		camera.ScreenX+camera.ScreenWidth/2,
		camera.ScreenY+camera.ScreenHeight/2,
	)
	var w, h = n.Width, n.Height
	var cw, ch = camera.Size()
	var scale = max(cw/w, ch/h)

	n.X = x - (0.5-n.PivotX)*w*scale
	n.Y = y - (0.5-n.PivotY)*h*scale
	n.ScaleX, n.ScaleY = scale, scale
	n.Angle = 0
}
func (n *Node) CameraStretch(camera *Camera) {
	var x, y = camera.PointFromScreen(
		camera.ScreenX+camera.ScreenWidth/2,
		camera.ScreenY+camera.ScreenHeight/2,
	)
	var w, h = n.Width, n.Height
	var cw, ch = camera.Size()
	var scaleX, scaleY = cw / w, ch / h

	n.X = x - (0.5-n.PivotX)*w*scaleX
	n.Y = y - (0.5-n.PivotY)*h*scaleY
	n.ScaleX, n.ScaleY = scaleX, scaleY
	n.Angle = 0
}

//=================================================================

func (n *Node) TransformToCamera() (cx, cy, cAngle, cScaleX, cScaleY float32) {
	var w, h = n.Width, n.Height
	var originPixelX = n.PivotX * float32(w)
	var originPixelY = n.PivotY * float32(h)
	var offsetX = -originPixelX * n.ScaleX
	var offsetY = -originPixelY * n.ScaleY
	var localRad = angle.ToRadians(n.Angle)
	var sinL, cosL = number.Sine(localRad), number.Cosine(localRad)
	var originOffsetX = offsetX*cosL - offsetY*sinL
	var originOffsetY = offsetX*sinL + offsetY*cosL
	var localX = n.X + originOffsetX
	var localY = n.Y + originOffsetY

	if n.Parent == nil {
		return localX, localY, n.Angle, n.ScaleX, n.ScaleY
	}

	localX -= offsetX
	localY -= offsetY

	var px, py, pa, psx, psy = n.Parent.TransformToCamera()

	localX *= psx
	localY *= psy

	var parentRad = angle.ToRadians(pa)
	var sinP, cosP = number.Sine(parentRad), number.Cosine(parentRad)
	var worldX = localX*cosP - localY*sinP + px
	var worldY = localX*sinP + localY*cosP + py

	cx = worldX
	cy = worldY
	cAngle = pa + n.Angle
	cScaleX = psx * n.ScaleX
	cScaleY = psy * n.ScaleY
	return
}
func (n *Node) TransformFromCamera(cx, cy, cAngle, cScaleX, cScaleY float32) (x, y, angle, scaleX, scaleY float32) {
	if n.Parent != nil {
		cx, cy, cAngle, cScaleX, cScaleY = n.Parent.TransformFromCamera(cx, cy, cAngle, cScaleX, cScaleY)
	}

	var dx = cx - n.X
	var dy = cy - n.Y
	var angleRad = toRad(-n.Angle)
	var sin, cos = number.Sine(angleRad), number.Cosine(angleRad)
	var localX = dx*cos - dy*sin
	var localY = dx*sin + dy*cos
	var w, h = n.Width, n.Height
	var pivotOffsetX = n.PivotX * w
	var pivotOffsetY = n.PivotY * h

	if n.ScaleX != 0 {
		localX /= n.ScaleX
	}
	if n.ScaleY != 0 {
		localY /= n.ScaleY
	}

	x = localX + pivotOffsetX
	y = localY + pivotOffsetY

	angle = cAngle - n.Angle
	scaleX = cScaleX / n.ScaleX
	scaleY = cScaleY / n.ScaleY
	return
}
func (n *Node) PointToCamera(camera *Camera, x, y float32) (cx, cy float32) {
	var w, h = n.Width, n.Height
	var originPixelX = n.PivotX * float32(w)
	var originPixelY = n.PivotY * float32(h)
	var localX = (x - originPixelX) * n.ScaleX
	var localY = (y - originPixelY) * n.ScaleY
	var localRad = toRad(n.Angle)
	var sinL, cosL = number.Sine(localRad), number.Cosine(localRad)
	var rotX = localX*cosL - localY*sinL
	var rotY = localX*sinL + localY*cosL
	var worldX = rotX + n.X
	var worldY = rotY + n.Y

	if n.Parent != nil {
		return n.Parent.PointToCamera(camera, worldX, worldY)
	}

	return worldX, worldY
}
func (n *Node) PointFromCamera(camera *Camera, cx, cy float32) (x, y float32) {
	x, y, _, _, _ = n.TransformFromCamera(cx, cy, 0, 1, 1)
	return x, y
}

func (n *Node) ContainsPoint(camera *Camera, cx, cy float32) bool {
	var x, y = n.PointFromCamera(camera, cx, cy)
	var w, h = n.Width, n.Height
	return x >= 0 && y >= 0 && x < w && y < h
}

func (n *Node) MousePosition(camera *Camera) (x, y float32) {
	x, y = camera.MousePosition()
	return n.PointFromCamera(camera, x, y)
}
func (n *Node) IsHovered(camera *Camera) bool {
	var x, y = camera.MousePosition()
	return n.ContainsPoint(camera, x, y)
}

func (n *Node) CornerTopLeft() (x, y float32)     { return n.getCorner(topLeft) }
func (n *Node) CornerTopRight() (x, y float32)    { return n.getCorner(topRight) }
func (n *Node) CornerBottomRight() (x, y float32) { return n.getCorner(bottomRight) }
func (n *Node) CornerBottomLeft() (x, y float32)  { return n.getCorner(bottomLeft) }

//=================================================================
// private

type corner byte

const (
	topLeft corner = iota
	topRight
	bottomRight
	bottomLeft
)

func toRad(ang float32) float32 { return angle.ToRadians(ang) }

func (n *Node) getCorner(corner corner) (x, y float32) {
	var width, height = n.Width, n.Height
	var nx, ny, na, _, _ = n.TransformToCamera()
	switch corner {
	case topLeft:
		return nx, ny
	case topRight:
		return point.MoveAtAngle(nx, ny, na, width)
	case bottomRight:
		var trx, try = point.MoveAtAngle(nx, ny, na, width)
		return point.MoveAtAngle(trx, try, na+90, height)
	case bottomLeft:
		return point.MoveAtAngle(nx, ny, na+90, height)
	}
	return
}
