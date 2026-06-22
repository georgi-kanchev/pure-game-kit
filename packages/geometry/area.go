package geometry

import "pure-game-kit/packages/internal"

type Area internal.Area

func NewArea(x, y, width, height float32) Area { return Area{X: x, Y: y, Width: width, Height: height} }

func (a Area) ContainsPoint(x, y float32) bool {
	return x > a.X && x < a.X+a.Width && y > a.Y && y < a.Y+a.Height
}
func (a Area) Overlaps(target Area) bool {
	return a.X < target.X+target.Width && a.X+a.Width > target.X && a.Y < target.Y+target.Height && a.Y+a.Height > target.Y
}
