package graphics

import (
	"math"
	"pure-kit/engine/internal"
)

type Node struct {
	X, Y, Angle, ScaleX, ScaleY      float32
	PivotX, PivotY, RepeatX, RepeatY float32
	AssetID                          string
	Parent                           *Node
	Tint                             uint
}

func NewNode(assetId string) Node {
	return Node{AssetID: assetId,
		ScaleX: 1, ScaleY: 1, RepeatX: 1, RepeatY: 1, PivotX: 0.5, PivotY: 0.5, Tint: math.MaxUint32}
}
func NewNodesGrid(tiles map[[2]int]string, cellWidth, cellHeight int, parent *Node) []Node {
	var result = []Node{}
	for k, v := range tiles {
		var node = NewNode(v)
		node.Parent = parent
		node.X = float32(k[0] * cellWidth)
		node.Y = float32(k[1] * cellHeight)
		result = append(result, node)
	}
	return result
}

func (node *Node) Size() (width, height float32) {
	var texture, fullTexture = internal.Textures[node.AssetID]
	width, height = 0, 0

	if fullTexture {
		return float32(texture.Width), float32(texture.Height)
	}

	var texRect, has = internal.AtlasRects[node.AssetID]
	if !has {
		return
	}

	var atlas = texRect.Atlas
	return float32(atlas.CellWidth) * texRect.CountX, float32(atlas.CellHeight) * texRect.CountY
}
func (node *Node) Global() (x, y, angle, scaleX, scaleY float32) {
	// Get texture size for origin offset
	var texWidth, texHeight = node.Size()

	// Step 1: Local origin offset, in local space
	originPixelX := node.PivotX * float32(texWidth)
	originPixelY := node.PivotY * float32(texHeight)
	offsetX := -originPixelX * node.ScaleX
	offsetY := -originPixelY * node.ScaleY

	// Step 2: Rotate origin offset by *local* rotation
	localRad := node.Angle * (math.Pi / 180)
	sinL, cosL := float32(math.Sin(float64(localRad))), float32(math.Cos(float64(localRad)))
	originOffsetX := offsetX*cosL - offsetY*sinL
	originOffsetY := offsetX*sinL + offsetY*cosL

	// Step 3: Local position, adjusted by origin
	localX := node.X + originOffsetX
	localY := node.Y + originOffsetY

	if node.Parent == nil {
		// Parent has no influence over the origin math
		return localX, localY, node.Angle, node.ScaleX, node.ScaleY
	}

	localX -= offsetX
	localY -= offsetY

	// Step 4: Get parent's global transform
	px, py, pr, psx, psy := node.Parent.Global()

	// Step 5: Apply parent scale to this node’s position
	localX *= psx
	localY *= psy

	// Step 6: Rotate local position by parent rotation
	parentRad := pr * (math.Pi / 180)
	sinP, cosP := float32(math.Sin(float64(parentRad))), float32(math.Cos(float64(parentRad)))
	worldX := localX*cosP - localY*sinP + px
	worldY := localX*sinP + localY*cosP + py

	// Final transform
	x = worldX
	y = worldY
	angle = pr + node.Angle
	scaleX = psx * node.ScaleX
	scaleY = psy * node.ScaleY
	return
}
