package geometry

import "pure-game-kit/packages/internal"

type Area internal.Area

func NewArea(x, y, width, height float32) Area { return Area{X: x, Y: y, Width: width, Height: height} }

func (a Area) ContainsPoint(x, y float32) bool {
	return x > a.X-a.Width/2 && x < a.X+a.Width/2 && y > a.Y-a.Height/2 && y < a.Y+a.Height/2
}
func (a Area) Overlaps(target Area) bool {
	var dx, dy = a.X - target.X, a.Y - target.Y
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	return dx < (a.Width+target.Width)/2 && dy < (a.Height+target.Height)/2
}
