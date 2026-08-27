package geometry

import (
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/number"
)

type Area internal.Area

func NewArea(x, y, width, height float32) Area { return Area{X: x, Y: y, Width: width, Height: height} }

//=================================================================

func (a Area) ContainsPoint(x, y float32) bool {
	return x > a.X-a.Width/2 && x < a.X+a.Width/2 && y > a.Y-a.Height/2 && y < a.Y+a.Height/2
}
func (a Area) Overlaps(target Area) bool {
	return number.Absolute(a.X-target.X) < (a.Width+target.Width)/2 &&
		number.Absolute(a.Y-target.Y) < (a.Height+target.Height)/2
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
func (a Area) Inside(target Area) Area {
	var newX, newY float32
	if a.Width >= target.Width {
		newX = target.X // hard center if it doesn't fit
	} else {
		newX = number.Limit(a.X, target.X-target.Width/2+a.Width/2, target.X+target.Width/2-a.Width/2)
	}

	if a.Height >= target.Height {
		newY = target.Y // hard center if it doesn't fit
	} else {
		newY = number.Limit(a.Y, target.Y-target.Height/2+a.Height/2, target.Y+target.Height/2-a.Height/2)
	}
	return NewArea(newX, newY, a.Width, a.Height)
}
func (a Area) Outside(target Area, forceHorizontal, forceVertical bool) Area {
	if !a.Overlaps(target) {
		return a
	}

	var minDistanceX, minDistanceY = a.Width/2 + target.Width/2, a.Height/2 + target.Height/2
	var currentDiffX, currentDiffY = a.X - target.X, a.Y - target.Y
	var overlapX = minDistanceX - number.Absolute(currentDiffX)
	var overlapY = minDistanceY - number.Absolute(currentDiffY)
	var newX, newY = a.X, a.Y
	var pushHorizontal bool // determine which axis to resolve along
	if forceHorizontal && !forceVertical {
		pushHorizontal = true
	} else if forceVertical && !forceHorizontal {
		pushHorizontal = false
	} else {
		pushHorizontal = overlapX < overlapY // default to least resistance if both true or both false
	}

	if pushHorizontal {
		if currentDiffX >= 0 {
			newX = target.X + minDistanceX
		} else {
			newX = target.X - minDistanceX
		}
	} else {
		if currentDiffY >= 0 {
			newY = target.Y + minDistanceY
		} else {
			newY = target.Y - minDistanceY
		}
	}
	return NewArea(newX, newY, a.Width, a.Height)
}

func (a Area) ToShape() Shape {
	return NewRectangle(a.X, a.Y, a.Width, a.Height, 0)
}
