package graphics

import (
	"math"
	"pure-kit/engine/internal"
)

type Node struct {
	X, Y, Angle, ScaleX, ScaleY      float32
	PivotX, PivotY, RepeatX, RepeatY float32
	AssetId                          string
	Parent                           *Node
	Tint                             uint
}

func NewNode(assetId string) Node {
	return Node{AssetId: assetId,
		ScaleX: 1, ScaleY: 1, RepeatX: 1, RepeatY: 1, PivotX: 0.5, PivotY: 0.5, Tint: math.MaxUint32}
}

func (node *Node) Size() (width, height float32) {
	var w, h = internal.AssetSize(node.AssetId)
	return float32(w), float32(h)
}

func (node *Node) ToCamera() (gX, gY, gAngle, gScaleX, gScaleY float32) {
	var texWidth, texHeight = node.Size()

	originPixelX := node.PivotX * float32(texWidth)
	originPixelY := node.PivotY * float32(texHeight)
	offsetX := -originPixelX * node.ScaleX
	offsetY := -originPixelY * node.ScaleY

	localRad := node.Angle * (math.Pi / 180)
	sinL, cosL := float32(math.Sin(float64(localRad))), float32(math.Cos(float64(localRad)))
	originOffsetX := offsetX*cosL - offsetY*sinL
	originOffsetY := offsetX*sinL + offsetY*cosL

	localX := node.X + originOffsetX
	localY := node.Y + originOffsetY

	if node.Parent == nil {
		return localX, localY, node.Angle, node.ScaleX, node.ScaleY
	}

	localX -= offsetX
	localY -= offsetY

	px, py, pr, psx, psy := node.Parent.ToCamera()

	localX *= psx
	localY *= psy

	parentRad := pr * (math.Pi / 180)
	sinP, cosP := float32(math.Sin(float64(parentRad))), float32(math.Cos(float64(parentRad)))
	worldX := localX*cosP - localY*sinP + px
	worldY := localX*sinP + localY*cosP + py

	gX = worldX
	gY = worldY
	gAngle = pr + node.Angle
	gScaleX = psx * node.ScaleX
	gScaleY = psy * node.ScaleY
	return
}
func (node *Node) FromCamera(gX, gY, gAngle, gScaleX, gScaleY float32) (x, y, angle, scaleX, scaleY float32) {
	if node.Parent != nil {
		gX, gY, gAngle, gScaleX, gScaleY = node.Parent.FromCamera(gX, gY, gAngle, gScaleX, gScaleY)
	}

	// Undo translation
	dx := gX - node.X
	dy := gY - node.Y

	// Undo rotation
	angleRad := -node.Angle * (math.Pi / 180)
	sin, cos := float32(math.Sin(float64(angleRad))), float32(math.Cos(float64(angleRad)))

	localX := dx*cos - dy*sin
	localY := dx*sin + dy*cos

	// Undo scale
	if node.ScaleX != 0 {
		localX /= node.ScaleX
	}
	if node.ScaleY != 0 {
		localY /= node.ScaleY
	}

	// Undo pivot (in local space)
	texWidth, texHeight := node.Size()
	pivotOffsetX := node.PivotX * texWidth
	pivotOffsetY := node.PivotY * texHeight

	x = localX + pivotOffsetX
	y = localY + pivotOffsetY

	angle = gAngle - node.Angle
	scaleX = gScaleX / node.ScaleX
	scaleY = gScaleY / node.ScaleY
	return
}
func (node *Node) PointToCamera(camera *Camera, x, y float32) (cX, cY float32) {
	// Start with local point (x, y) in local node space.
	// Adjust for pivot:
	texWidth, texHeight := node.Size()
	originPixelX := node.PivotX * float32(texWidth)
	originPixelY := node.PivotY * float32(texHeight)

	// Local offset relative to the pivot.
	localX := (x - originPixelX) * node.ScaleX
	localY := (y - originPixelY) * node.ScaleY

	// Rotate local point by node's angle.
	localRad := node.Angle * (math.Pi / 180)
	sinL, cosL := float32(math.Sin(float64(localRad))), float32(math.Cos(float64(localRad)))
	rotX := localX*cosL - localY*sinL
	rotY := localX*sinL + localY*cosL

	// Translate by node position.
	worldX := rotX + node.X
	worldY := rotY + node.Y

	// Recurse up the parent chain.
	if node.Parent != nil {
		parentX, parentY := node.Parent.PointToCamera(camera, worldX, worldY)
		return parentX, parentY
	}

	// If no parent, this is the camera position.
	return worldX, worldY
}
func (node *Node) PointFromCamera(camera *Camera, cX, cY float32) (x, y float32) {
	x, y, _, _, _ = node.FromCamera(cX, cY, 0, 1, 1)
	return x, y
}

func (node *Node) MousePosition(camera *Camera) (x, y float32) {
	x, y = camera.MousePosition()
	return node.PointFromCamera(camera, x, y)
}
func (node *Node) IsHovered(camera *Camera) bool {
	mx, my := node.MousePosition(camera)
	w, h := node.Size()
	return mx >= 0 && my >= 0 && mx < w && my < h
}

func (node *Node) Fit(camera *Camera) {
	var x, y = camera.PointFromScreen(
		camera.ScreenX+camera.ScreenWidth/2,
		camera.ScreenY+camera.ScreenHeight/2,
	)
	var w, h = node.Size()
	var cw, ch = camera.Size()
	var scale = min(cw/w, ch/h)

	node.X = x - (0.5-node.PivotX)*w*scale
	node.Y = y - (0.5-node.PivotY)*h*scale
	node.ScaleX, node.ScaleY = scale, scale
	node.Angle = 0
}
func (node *Node) Fill(camera *Camera) {
	var x, y = camera.PointFromScreen(
		camera.ScreenX+camera.ScreenWidth/2,
		camera.ScreenY+camera.ScreenHeight/2,
	)
	var w, h = node.Size()
	var cw, ch = camera.Size()
	var scale = max(cw/w, ch/h)

	node.X = x - (0.5-node.PivotX)*w*scale
	node.Y = y - (0.5-node.PivotY)*h*scale
	node.ScaleX, node.ScaleY = scale, scale
	node.Angle = 0
}
func (node *Node) Stretch(camera *Camera) {
	var x, y = camera.PointFromScreen(
		camera.ScreenX+camera.ScreenWidth/2,
		camera.ScreenY+camera.ScreenHeight/2,
	)
	var w, h = node.Size()
	var cw, ch = camera.Size()
	var scaleX, scaleY = cw / w, ch / h

	node.X = x - (0.5-node.PivotX)*w*scaleX
	node.Y = y - (0.5-node.PivotY)*h*scaleY
	node.ScaleX, node.ScaleY = scaleX, scaleY
	node.Angle = 0
}
