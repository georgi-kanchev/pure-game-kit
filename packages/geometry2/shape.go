package geometry

import (
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/angle"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/point"
)

type Shape struct{ X, Y, Width, Height, Angle, Roundness float32 }

// NewPoint creates a Shape with zero dimensions.
func NewPoint(x, y float32) Shape {
	return Shape{X: x, Y: y}
}

// NewCircle creates a square shape with maximum roundness.
func NewCircle(x, y, radius float32) Shape {
	return Shape{X: x, Y: y, Width: radius * 2, Height: radius * 2, Roundness: 1.0}
}

// NewRectangle creates a standard rectangle with no rounded corners.
func NewRectangle(x, y, width, height, angle float32) Shape {
	return Shape{X: x, Y: y, Width: width, Height: height, Angle: angle}
}

// NewRoundedRectangle creates a rectangle with a specific rounding factor (0 to 1).
func NewRoundedRectangle(x, y, width, height, angle, roundness float32) Shape {
	return Shape{X: x, Y: y, Width: width, Height: height, Angle: angle, Roundness: roundness}
}

// NewCapsule creates a capsule where (x1, y1) and (x2, y2) are the center points of the circular end-caps.
func NewCapsule(x1, y1, x2, y2, radius float32) Shape {
	var dist = point.DistanceToPoint(x1, y1, x2, y2)
	var ang = angle.BetweenPoints(x1, y1, x2, y2)
	var midX, midY = point.MoveAtAngle(x1, y1, ang, dist*0.5)
	return Shape{X: midX, Y: midY, Width: dist + (radius * 2), Height: radius * 2, Angle: ang, Roundness: 1.0}
}

// NewLine creates a line segment of a specific thickness.
func NewLine(x1, y1, x2, y2, thickness float32) Shape {
	var dist = point.DistanceToPoint(x1, y1, x2, y2)
	var ang = angle.BetweenPoints(x1, y1, x2, y2)
	var midX, midY = point.MoveAtAngle(x1, y1, ang, dist*0.5)
	return Shape{X: midX, Y: midY, Width: dist, Height: thickness, Angle: ang}
}

//=================================================================

// Negative result means inside, zero means on the edge, positive means outside.
func (s Shape) DistanceToPoint(x, y float32) float32 {
	var px, py = x - s.X, y - s.Y
	var sinR, cosR = internal.SinCos(-s.Angle)
	var lx = px*cosR - py*sinR
	var ly = px*sinR + py*cosR
	var hx, hy = s.Width * 0.5, s.Height * 0.5
	var r = s.roundness() * min(hx, hy)
	var qx = number.Absolute(lx) - (hx - r)
	var qy = number.Absolute(ly) - (hy - r)
	return number.SquareRoot(max(qx, 0)*max(qx, 0)+max(qy, 0)*max(qy, 0)) + min(max(qx, qy), 0) - r
}
func (s Shape) ContainsPoint(x, y float32) bool {
	return s.DistanceToPoint(x, y) <= 0
}
func (s Shape) ClosestPointToEdge(x, y float32) (edgeX, edgeY float32) {
	var px, py = x - s.X, y - s.Y
	var sinR, cosR = internal.SinCos(-s.Angle)
	var lx = px*cosR - py*sinR
	var ly = px*sinR + py*cosR
	var hx, hy = s.Width * 0.5, s.Height * 0.5
	var r = s.roundness() * min(hx, hy)
	var cx = max(-(hx - r), min(hx-r, lx))
	var cy = max(-(hy - r), min(hy-r, ly))
	var dx, dy = lx - cx, ly - cy
	var dist = number.SquareRoot(dx*dx + dy*dy)

	var bx, by float32
	if dist > 1e-6 {
		bx = cx + r*dx/dist
		by = cy + r*dy/dist
	} else { // inside inner box: push to nearest flat face
		var dPosX, dNegX = hx - lx, hx + lx
		var dPosY, dNegY = hy - ly, hy + ly
		var minD = min(min(dPosX, dNegX), min(dPosY, dNegY))
		switch minD {
		case dPosX:
			bx, by = hx, ly
		case dNegX:
			bx, by = -hx, ly
		case dPosY:
			bx, by = lx, hy
		default:
			bx, by = lx, -hy
		}
	}

	return bx*cosR + by*sinR + s.X, -bx*sinR + by*cosR + s.Y
}
func (s Shape) Raycast(x, y, angle, length float32) (hitX, hitY float32) {
	var lx, ly = point.MoveAtAngle(x, y, angle, length*0.5)
	if !s.Overlaps(Shape{X: lx, Y: ly, Width: length, Angle: angle}) {
		return number.NaN(), number.NaN()
	}
	var dy, dx = internal.SinCos(angle)
	var t = float32(0)
	for range 64 {
		var cx, cy = x + t*dx, y + t*dy
		var dist = s.DistanceToPoint(cx, cy)
		if dist <= 1e-4 {
			return s.ClosestPointToEdge(cx, cy)
		}
		t += dist
	}
	return number.NaN(), number.NaN()
}
func (s Shape) Bounds() (minX, minY, maxX, maxY float32) {
	var sinR, cosR = internal.SinCos(s.Angle)
	sinR, cosR = number.Absolute(sinR), number.Absolute(cosR)
	var hx, hy = s.Width * 0.5, s.Height * 0.5
	var r = s.roundness() * min(hx, hy)
	var extentX = (hx-r)*cosR + (hy-r)*sinR + r
	var extentY = (hx-r)*sinR + (hy-r)*cosR + r
	return s.X - extentX, s.Y - extentY, s.X + extentX, s.Y + extentY
}
func (s Shape) Overlaps(target Shape) bool {
	// AABB broadphase
	var sMinX, sMinY, sMaxX, sMaxY = s.Bounds()
	var oMinX, oMinY, oMaxX, oMaxY = target.Bounds()
	if sMaxX < oMinX || oMaxX < sMinX || sMaxY < oMinY || oMaxY < sMinY {
		return false
	}

	var dx, dy = target.X - s.X, target.Y - s.Y
	var sSin, sCos = internal.SinCos(s.Angle)
	var oSin, oCos = internal.SinCos(target.Angle)

	var projX = number.Absolute(dx*sCos + dy*sSin) // s local X
	if projX > s.support(sCos, sSin, sCos, sSin)+target.support(sCos, sSin, oCos, oSin) {
		return false
	}
	var projY = number.Absolute(dx*-sSin + dy*sCos) // s local Y
	if projY > s.support(-sSin, sCos, sCos, sSin)+target.support(-sSin, sCos, oCos, oSin) {
		return false
	}
	var projOx = number.Absolute(dx*oCos + dy*oSin) // other local X
	if projOx > s.support(oCos, oSin, sCos, sSin)+target.support(oCos, oSin, oCos, oSin) {
		return false
	}

	var projOy = number.Absolute(dx*-oSin + dy*oCos) // other local Y
	if projOy > s.support(-oSin, oCos, sCos, sSin)+target.support(-oSin, oCos, oCos, oSin) {
		return false
	}

	// corner axis
	var pAx, pAy = s.nearestInnerBoxPoint(target.X, target.Y, sCos, sSin)
	var pBx, pBy = target.nearestInnerBoxPoint(s.X, s.Y, oCos, oSin)
	var cax, cay = pBx - pAx, pBy - pAy
	var l = number.SquareRoot(cax*cax + cay*cay)
	if l > 1e-6 {
		cax, cay = cax/l, cay/l
		var proj = number.Absolute(dx*cax + dy*cay)
		if proj > s.support(cax, cay, sCos, sSin)+target.support(cax, cay, oCos, oSin) {
			return false
		}
	}
	return true
}
func (s Shape) Collide(target Shape) Shape {
	if !s.Overlaps(target) {
		return target
	}

	var dx, dy = target.X - s.X, target.Y - s.Y
	var sSin, sCos = internal.SinCos(s.Angle)
	var oSin, oCos = internal.SinCos(target.Angle)
	var minDepth = number.ValueBiggest[float32]()
	var minAx0, minAx1 float32

	{ // s local X
		var ax0, ax1 = sCos, sSin
		var d = dx*ax0 + dy*ax1
		if d < 0 {
			d, ax0, ax1 = -d, -ax0, -ax1
		}
		var depth = s.support(ax0, ax1, sCos, sSin) + target.support(ax0, ax1, oCos, oSin) - d
		if depth < minDepth {
			minDepth, minAx0, minAx1 = depth, ax0, ax1
		}
	}
	{ // s local Y
		var ax0, ax1 = -sSin, sCos
		var d = dx*ax0 + dy*ax1
		if d < 0 {
			d, ax0, ax1 = -d, -ax0, -ax1
		}
		var depth = s.support(ax0, ax1, sCos, sSin) + target.support(ax0, ax1, oCos, oSin) - d
		if depth < minDepth {
			minDepth, minAx0, minAx1 = depth, ax0, ax1
		}
	}
	{ // other local X
		var ax0, ax1 = oCos, oSin
		var d = dx*ax0 + dy*ax1
		if d < 0 {
			d, ax0, ax1 = -d, -ax0, -ax1
		}
		var depth = s.support(ax0, ax1, sCos, sSin) + target.support(ax0, ax1, oCos, oSin) - d
		if depth < minDepth {
			minDepth, minAx0, minAx1 = depth, ax0, ax1
		}
	}
	{ // other local Y
		var ax0, ax1 = -oSin, oCos
		var d = dx*ax0 + dy*ax1
		if d < 0 {
			d, ax0, ax1 = -d, -ax0, -ax1
		}
		var depth = s.support(ax0, ax1, sCos, sSin) + target.support(ax0, ax1, oCos, oSin) - d
		if depth < minDepth {
			minDepth, minAx0, minAx1 = depth, ax0, ax1
		}
	}

	// corner axis
	var pAx, pAy = s.nearestInnerBoxPoint(target.X, target.Y, sCos, sSin)
	var pBx, pBy = target.nearestInnerBoxPoint(s.X, s.Y, oCos, oSin)
	var ax0, ax1 = pBx - pAx, pBy - pAy
	var l = number.SquareRoot(ax0*ax0 + ax1*ax1)
	if l > 1e-6 {
		ax0, ax1 = ax0/l, ax1/l
		var d = dx*ax0 + dy*ax1
		if d < 0 {
			d, ax0, ax1 = -d, -ax0, -ax1
		}
		var depth = s.support(ax0, ax1, sCos, sSin) + target.support(ax0, ax1, oCos, oSin) - d
		if depth < minDepth {
			minDepth, minAx0, minAx1 = depth, ax0, ax1
		}
	}

	target.X += minAx0 * minDepth
	target.Y += minAx1 * minDepth
	return target
}

// private ========================================================

func (s Shape) support(ax0, ax1, cosR, sinR float32) float32 {
	var hx, hy = s.Width * 0.5, s.Height * 0.5
	var r = s.roundness() * min(hx, hy)
	var dX = number.Absolute(ax0*cosR + ax1*sinR)
	var dY = number.Absolute(-ax0*sinR + ax1*cosR)
	return (hx-r)*dX + (hy-r)*dY + r
}
func (s Shape) nearestInnerBoxPoint(px, py, cosR, sinR float32) (float32, float32) {
	var rx, ry = px - s.X, py - s.Y
	var lx = rx*cosR + ry*sinR
	var ly = -rx*sinR + ry*cosR
	var hx, hy = s.Width * 0.5, s.Height * 0.5
	var r = s.roundness() * min(hx, hy)
	lx = max(-(hx - r), min(hx-r, lx))
	ly = max(-(hy - r), min(hy-r, ly))
	return s.X + lx*cosR - ly*sinR, s.Y + lx*sinR + ly*cosR
}

func (s Shape) roundness() float32 {
	return max(min(s.Roundness, 1), 0)
}
