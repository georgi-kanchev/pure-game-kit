package geometry

import (
	"pure-kit/engine/geometry/point"
	"pure-kit/engine/utility/angle"
	"pure-kit/engine/utility/number"
)

type Line struct{ Ax, Ay, Bx, By float32 }

func NewLine(ax, ay, bx, by float32) Line {
	return Line{Ax: ax, Ay: ay, Bx: bx, By: by}
}

//=================================================================
// getters

func (line *Line) Angle() float32 {
	return angle.BetweenPoints(line.Ax, line.Ay, line.Bx, line.By)
}
func (line *Line) Normal() float32 {
	return number.Wrap(line.Angle()-90, 360)
}
func (line *Line) Length() float32 {
	return point.DistanceToPoint(line.Ax, line.Ay, line.Bx, line.By)
}

func (line *Line) CrossPointWithLine(target Line) (x, y float32) {
	var dx1 = line.Bx - line.Ax
	var dy1 = line.By - line.Ay
	var dx2 = target.Bx - target.Ax
	var dy2 = target.By - target.Ay
	var det = dx1*dy2 - dy1*dx2

	if det > -0.001 && det < 0.001 { // Lines are parallel or duplicate
		return number.NaN(), number.NaN()
	}

	var s = ((line.Ay-target.Ay)*dx2 - (line.Ax-target.Ax)*dy2) / det
	var t = ((line.Ay-target.Ay)*dx1 - (line.Ax-target.Ax)*dy1) / det

	if s < 0 || s > 1 || t < 0 || t > 1 { // Intersection not within both segments
		return number.NaN(), number.NaN()
	}

	var ix = line.Ax + s*dx1
	var iy = line.Ay + s*dy1

	return ix, iy
}
func (line *Line) ClosestToPoint(targetX, targetY float32) (x, y float32) {
	var ax, ay = line.Ax, line.Ay
	var bx, by = line.Bx, line.By
	var apx, apy = targetX - ax, targetY - ay
	var abx, aby = bx - ax, by - ay
	var magnitude = abx*abx + aby*aby

	if magnitude == 0 { // Line is just a point
		return ax, ay
	}

	var dot = apx*abx + apy*aby
	var distance = dot / magnitude

	if distance < 0 {
		return ax, ay
	}
	if distance > 1 {
		return bx, by
	}

	var cx = ax + abx*distance
	var cy = ay + aby*distance
	return cx, cy
}

func (line *Line) IsCrossingLine(target Line) bool {
	var ax1, ay1, bx1, by1 = line.Ax, line.Ay, line.Bx, line.By
	var ax2, ay2, bx2, by2 = target.Ax, target.Ay, target.Bx, target.By
	var d1 = (bx2-ax2)*(ay1-ay2) - (by2-ay2)*(ax1-ax2)
	var d2 = (bx2-ax2)*(by1-ay2) - (by2-ay2)*(bx1-ax2)
	var d3 = (bx1-ax1)*(ay2-ay1) - (by1-ay1)*(ax2-ax1)
	var d4 = (bx1-ax1)*(by2-ay1) - (by1-ay1)*(bx2-ax1)
	return d1*d2 < 0 && d3*d4 < 0
}
func (line *Line) IsLeftOfPoint(x, y float32) bool {
	return (line.Bx-line.Ax)*(y-line.Ay)-(line.By-line.Ay)*(x-line.Ax) < 0
}
