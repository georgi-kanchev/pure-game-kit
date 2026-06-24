package geometry

import (
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/number"
)

type Area internal.Area

func NewArea(x, y, width, height float32) Area { return Area{X: x, Y: y, Width: width, Height: height} }

func (a Area) ContainsPoint(x, y float32) bool {
	return x > a.X-a.Width/2 && x < a.X+a.Width/2 && y > a.Y-a.Height/2 && y < a.Y+a.Height/2
}
func (a Area) Overlaps(target Area) bool {
	return number.Absolute(a.X-target.X) < (a.Width+target.Width)/2 && number.Absolute(a.Y-target.Y) < (a.Height+target.Height)/2
}
func (a Area) Intersect(target Area) Area {
	if target == (Area{}) {
		return a
	}
	if !a.Overlaps(target) {
		return NewArea(number.NaN(), number.NaN(), number.NaN(), number.NaN())
	}
	var minX, maxX = max(a.X-a.Width/2, target.X-target.Width/2), min(a.X+a.Width/2, target.X+target.Width/2)
	var minY, maxY = max(a.Y-a.Height/2, target.Y-target.Height/2), min(a.Y+a.Height/2, target.Y+target.Height/2)
	var newWidth, newHeight = maxX - minX, maxY - minY
	return NewArea(minX+newWidth/2, minY+newHeight/2, newWidth, newHeight)
}
