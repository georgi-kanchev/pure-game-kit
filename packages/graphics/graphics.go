// Strictly tied to the window, drawing on it through a view and converting between the two coordinate systems.
// The view's drawing consists of two categories: primitives and objects. While using the assets for drawing,
// the graphical objects are still very lightweight and exist independently of them.
// The concept map of the package looks like this:
//   - View - draws objects
//   - Quad - no asset, flat color (useful for batching shapes)
//   - ├ Sprite - texture asset
//   - ├ NinePatch - box asset
//   - ├ TextBox - font asset
//   - └ TileMap - tile set asset + tile data asset
package graphics

import (
	"pure-game-kit/packages/internal"
)

type Area struct{ X, Y, Width, Height float32 }
type Effects internal.Effects

func NewArea(x, y, width, height float32) Area {
	return Area{X: x, Y: y, Width: width, Height: height}
}
