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
func NewNodesGrid(tiles map[[2]float32]string, cellWidth, cellHeight int, parent *Node) []Node {
	var result = []Node{}
	for k, v := range tiles {
		var node = NewNode(v)
		node.Parent = parent
		node.X = float32(k[0] * float32(cellWidth))
		node.Y = float32(k[1] * float32(cellHeight))
		result = append(result, node)
	}
	return result
}

func (node *Node) MousePosition(camera *Camera) (x, y float32) {
	x, y = camera.MousePosition()
	x, y, _, _, _ = node.FromGlobal(x, y, 0, 1, 1)
	return x, y
}
func (node *Node) Size() (width, height float32) {
	var w, h = internal.AssetSize(node.AssetId)
	return float32(w), float32(h)
}
func (node *Node) ToGlobal() (gX, gY, gAngle, gScaleX, gScaleY float32) {
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

	px, py, pr, psx, psy := node.Parent.ToGlobal()

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
func (node *Node) FromGlobal(gX, gY, gAngle, gScaleX, gScaleY float32) (x, y, angle, scaleX, scaleY float32) {
	if node.Parent != nil {
		gX, gY, gAngle, gScaleX, gScaleY = node.Parent.FromGlobal(gX, gY, gAngle, gScaleX, gScaleY)
	}

	texWidth, texHeight := node.Size()
	pivotOffsetX := node.PivotX * texWidth * node.ScaleX
	pivotOffsetY := node.PivotY * texHeight * node.ScaleY

	angleRad := -node.Angle * (math.Pi / 180)
	sin, cos := float32(math.Sin(float64(angleRad))), float32(math.Cos(float64(angleRad)))

	dx := gX - node.X
	dy := gY - node.Y

	localX := dx*cos - dy*sin
	localY := dx*sin + dy*cos

	if node.ScaleX != 0 {
		localX /= node.ScaleX
	}
	if node.ScaleY != 0 {
		localY /= node.ScaleY
	}

	x = localX + pivotOffsetX
	y = localY + pivotOffsetY

	angle = gAngle - node.Angle
	scaleX = gScaleX / node.ScaleX
	scaleY = gScaleY / node.ScaleY
	return
}
