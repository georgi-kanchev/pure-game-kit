package graphics

import (
	"math"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/utility/point"
)

type Node struct {
	X, Y, Angle    float32
	Width, Height  float32
	ScaleX, ScaleY float32
	PivotX, PivotY float32
	Parent         *Node
	Color          uint
}

func NewNode(x, y float32) Node {
	return Node{X: x, Y: y, Width: 100, Height: 100, ScaleX: 1, ScaleY: 1,
		PivotX: 0.5, PivotY: 0.5, Color: color.White}
}

func (node *Node) TransformToCamera() (cX, cY, cAngle, cScaleX, cScaleY float32) {
	var w, h = node.Width, node.Height
	var originPixelX = node.PivotX * float32(w)
	var originPixelY = node.PivotY * float32(h)
	var offsetX = -originPixelX * node.ScaleX
	var offsetY = -originPixelY * node.ScaleY
	var localRad = node.Angle * (math.Pi / 180)
	var sinL, cosL = float32(math.Sin(float64(localRad))), float32(math.Cos(float64(localRad)))
	var originOffsetX = offsetX*cosL - offsetY*sinL
	var originOffsetY = offsetX*sinL + offsetY*cosL
	var localX = node.X + originOffsetX
	var localY = node.Y + originOffsetY

	if node.Parent == nil {
		return localX, localY, node.Angle, node.ScaleX, node.ScaleY
	}

	localX -= offsetX
	localY -= offsetY

	var px, py, pr, psx, psy = node.Parent.TransformToCamera()

	localX *= psx
	localY *= psy

	var parentRad = pr * (math.Pi / 180)
	var sinP, cosP = float32(math.Sin(float64(parentRad))), float32(math.Cos(float64(parentRad)))
	var worldX = localX*cosP - localY*sinP + px
	var worldY = localX*sinP + localY*cosP + py

	cX = worldX
	cY = worldY
	cAngle = pr + node.Angle
	cScaleX = psx * node.ScaleX
	cScaleY = psy * node.ScaleY
	return
}
func (node *Node) TransformFromCamera(cX, cY, cAngle, cScaleX, cScaleY float32) (x, y, angle, scaleX, scaleY float32) {
	if node.Parent != nil {
		cX, cY, cAngle, cScaleX, cScaleY = node.Parent.TransformFromCamera(cX, cY, cAngle, cScaleX, cScaleY)
	}

	var dx = cX - node.X
	var dy = cY - node.Y
	var angleRad = -node.Angle * (math.Pi / 180)
	var sin, cos = float32(math.Sin(float64(angleRad))), float32(math.Cos(float64(angleRad)))
	var localX = dx*cos - dy*sin
	var localY = dx*sin + dy*cos
	var w, h = node.Width, node.Height
	var pivotOffsetX = node.PivotX * w
	var pivotOffsetY = node.PivotY * h

	if node.ScaleX != 0 {
		localX /= node.ScaleX
	}
	if node.ScaleY != 0 {
		localY /= node.ScaleY
	}

	x = localX + pivotOffsetX
	y = localY + pivotOffsetY

	angle = cAngle - node.Angle
	scaleX = cScaleX / node.ScaleX
	scaleY = cScaleY / node.ScaleY
	return
}
func (node *Node) PointToCamera(camera *Camera, x, y float32) (cX, cY float32) {
	var w, h = node.Width, node.Height
	var originPixelX = node.PivotX * float32(w)
	var originPixelY = node.PivotY * float32(h)
	var localX = (x - originPixelX) * node.ScaleX
	var localY = (y - originPixelY) * node.ScaleY
	var localRad = node.Angle * (math.Pi / 180)
	var sinL, cosL = float32(math.Sin(float64(localRad))), float32(math.Cos(float64(localRad)))
	var rotX = localX*cosL - localY*sinL
	var rotY = localX*sinL + localY*cosL
	var worldX = rotX + node.X
	var worldY = rotY + node.Y

	if node.Parent != nil {
		return node.Parent.PointToCamera(camera, worldX, worldY)
	}

	return worldX, worldY
}
func (node *Node) PointFromCamera(camera *Camera, cX, cY float32) (x, y float32) {
	x, y, _, _, _ = node.TransformFromCamera(cX, cY, 0, 1, 1)
	return x, y
}

func (node *Node) ContainsPoint(camera *Camera, cX, cY float32) bool {
	var x, y = node.PointFromCamera(camera, cX, cY)
	var w, h = node.Width, node.Height
	return x >= 0 && y >= 0 && x < w && y < h
}

func (node *Node) MousePosition(camera *Camera) (x, y float32) {
	x, y = camera.MousePosition()
	return node.PointFromCamera(camera, x, y)
}
func (node *Node) IsHovered(camera *Camera) bool {
	var x, y = camera.MousePosition()
	return node.ContainsPoint(camera, x, y)
}

func (node *Node) CameraFit(camera *Camera) {
	var x, y = camera.PointFromScreen(
		camera.ScreenX+camera.ScreenWidth/2,
		camera.ScreenY+camera.ScreenHeight/2,
	)
	var w, h = node.Width, node.Height
	var cw, ch = camera.Size()
	var scale = min(cw/w, ch/h)

	node.X = x - (0.5-node.PivotX)*w*scale
	node.Y = y - (0.5-node.PivotY)*h*scale
	node.ScaleX, node.ScaleY = scale, scale
	node.Angle = 0
}
func (node *Node) CameraFill(camera *Camera) {
	var x, y = camera.PointFromScreen(
		camera.ScreenX+camera.ScreenWidth/2,
		camera.ScreenY+camera.ScreenHeight/2,
	)
	var w, h = node.Width, node.Height
	var cw, ch = camera.Size()
	var scale = max(cw/w, ch/h)

	node.X = x - (0.5-node.PivotX)*w*scale
	node.Y = y - (0.5-node.PivotY)*h*scale
	node.ScaleX, node.ScaleY = scale, scale
	node.Angle = 0
}
func (node *Node) CameraStretch(camera *Camera) {
	var x, y = camera.PointFromScreen(
		camera.ScreenX+camera.ScreenWidth/2,
		camera.ScreenY+camera.ScreenHeight/2,
	)
	var w, h = node.Width, node.Height
	var cw, ch = camera.Size()
	var scaleX, scaleY = cw / w, ch / h

	node.X = x - (0.5-node.PivotX)*w*scaleX
	node.Y = y - (0.5-node.PivotY)*h*scaleY
	node.ScaleX, node.ScaleY = scaleX, scaleY
	node.Angle = 0
}

func (node *Node) CornerUpperLeft() (x, y float32)  { return node.getCorner(upperLeft) }
func (node *Node) CornerUpperRight() (x, y float32) { return node.getCorner(upperRight) }
func (node *Node) CornerLowerRight() (x, y float32) { return node.getCorner(lowerRight) }
func (node *Node) CornerLowerLeft() (x, y float32)  { return node.getCorner(lowerLeft) }

// #region private

type corner byte

const (
	upperLeft corner = iota
	upperRight
	lowerRight
	lowerLeft
)

func (node *Node) getCorner(corner corner) (x, y float32) {
	var width, height = node.Width, node.Height
	var nx, ny, na, _, _ = node.TransformToCamera()
	var offX, offY = -width * node.PivotX, -height * node.PivotY
	if corner == upperRight || corner == lowerRight {
		offX = width * (1 - node.PivotX)
	}
	if corner == lowerLeft || corner == lowerRight {
		offY = height * (1 - node.PivotY)
	}
	x, y = point.MoveAt(nx, ny, na, offX)
	x, y = point.MoveAt(x, y, na+90, offY)
	return
}

// #endregion
