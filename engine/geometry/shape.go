package geometry

import (
	"math"
	"pure-kit/engine/geometry/point"
	"pure-kit/engine/utility/angle"
)

type Shape struct {
	X, Y, Angle    float32
	ScaleX, ScaleY float32

	minX, minY float32
	maxX, maxY float32

	corners [][2]float32
}

func NewShapeCorners(corners ...[2]float32) Shape {
	if len(corners) == 0 {
		return Shape{}
	}
	return Shape{ScaleX: 1, ScaleY: 1, corners: append(corners, corners[0])}
}
func NewShapeSides(radius float32, sides int) Shape {
	var corners = [][2]float32{}
	var step float32 = 360.0 / float32(sides)
	for i := range sides {
		var x, y = point.MoveAtAngle(0, 0, step*float32(i)-90, radius)
		corners = append(corners, [2]float32{x, y})
	}

	if len(corners) > 0 {
		corners = append(corners, corners[0])
	}

	return Shape{ScaleX: 1, ScaleY: 1, corners: corners}
}
func NewShapeRectangle(width, height, pivotX, pivotY float32) Shape {
	var offX, offY = width * pivotX, height * pivotY
	var corners = [][2]float32{
		{-offX, -offY},
		{width - offX, -offY},
		{width - offX, height - offY},
		{-offX, height - offY},
		{-offX, -offY},
	}
	return Shape{ScaleX: 1, ScaleY: 1, corners: corners}
}

func (shape *Shape) CornerPoints() [][2]float32 {
	var result = make([][2]float32, len(shape.corners))
	shape.minX, shape.minY = float32(math.Inf(1)), float32(math.Inf(1))
	shape.maxX, shape.maxY = float32(math.Inf(-1)), float32(math.Inf(-1))

	for i := range shape.corners {
		var x, y = shape.corners[i][0], shape.corners[i][1]

		x *= shape.ScaleX
		y *= shape.ScaleY

		var rad = angle.ToRadians(shape.Angle)
		var sin, cos = float32(math.Sin(float64(rad))), float32(math.Cos(float64(rad)))
		var resultX = shape.X + (x*cos - y*sin)
		var resultY = shape.Y + (x*sin + y*cos)

		if shape.minX > resultX {
			shape.minX = resultX
		}
		if shape.minY > resultY {
			shape.minY = resultY
		}
		if shape.maxX < resultX {
			shape.maxX = resultX
		}
		if shape.maxY < resultY {
			shape.maxY = resultY
		}

		result[i] = [2]float32{resultX, resultY}
	}
	return result
}

// all check methods are made with speed in mind, not so much readability
// they should have the least allocations (once per API call for CornerPoints())
// and the least iterations (1 loop for 1 shape and 2 loops for 2 shapes)
// don't be a smartass by "simplifying" and reusing them internally in the future

func (shape *Shape) IsContainingPoint(x, y float32) bool {
	var corners = shape.CornerPoints()

	if !shape.inBoundingBoxPoint(x, y) {
		return false
	}

	return shape.internalIsContainingPoint(corners, x, y)
}

func (shape *Shape) CrossPointsWithLines(lines ...Line) [][2]float32 {
	var corners = shape.CornerPoints()
	var result = [][2]float32{}

	for _, line := range lines {
		if shape.inBoundingBoxLine(line) {
			result = append(result, shape.internalCrossPointsWithLine(corners, line)...)
		}
	}
	return result
}
func (shape *Shape) IsCrossingLines(lines ...Line) bool {
	var corners = shape.CornerPoints()
	for _, line := range lines {
		if shape.inBoundingBoxLine(line) && shape.internalIsCrossingLine(corners, line) {
			return true
		}
	}
	return false
}
func (shape *Shape) IsContainingLines(lines ...Line) bool {
	var corners = shape.CornerPoints()
	for _, line := range lines {
		if !shape.inBoundingBoxLine(line) || !shape.internalIsContainingLine(corners, line) {
			return false
		}
	}
	return true
}
func (shape *Shape) IsOverlappingLines(lines ...Line) bool {
	var corners = shape.CornerPoints()
	for _, line := range lines {
		if shape.inBoundingBoxLine(line) && shape.internalIsOverlappingLine(corners, line) {
			return true
		}
	}
	return false
}

func (shape *Shape) CrossPointsWithShapes(shapes ...*Shape) [][2]float32 {
	var corners = shape.CornerPoints()
	var result = [][2]float32{}

	for _, target := range shapes {
		var targetCorners = target.CornerPoints()
		if shape.inBoundingBoxShape(*target) {
			result = append(result, shape.internalCrossPointsWithShape(corners, targetCorners)...)
		}
	}
	return result
}
func (shape *Shape) IsCrossingShapes(shapes ...*Shape) bool {
	var corners = shape.CornerPoints()
	for _, target := range shapes {
		var targetCorners = target.CornerPoints()
		if shape.inBoundingBoxShape(*target) && shape.internalIsCrossingShape(corners, targetCorners) {
			return true
		}
	}
	return false
}
func (shape *Shape) IsContainingShapes(shapes ...*Shape) bool {
	var corners = shape.CornerPoints()
	for _, target := range shapes {
		var targetCorners = target.CornerPoints()
		if !shape.inBoundingBoxShape(*target) || !shape.internalIsContainingShapes(corners, targetCorners) {
			return false
		}
	}
	return true
}
func (shape *Shape) IsOverlappingShapes(shapes ...*Shape) bool {
	var corners = shape.CornerPoints()

	for _, target := range shapes {
		var targetCorners = target.CornerPoints()
		if shape.inBoundingBoxShape(*target) && shape.internalIsOverlappingShape(corners, targetCorners, target) {
			return true
		}
	}
	return false
}

func (shape *Shape) Collide(velocityX, velocityY float32, targets ...*Shape) (newVelocityX, newVelocityY float32) {
	for _, target := range targets {
		if !shape.inBoundingBoxShape(*target) {
			continue
		}

		var cax = shape.minX + (shape.maxX-shape.minX)/2
		var cay = shape.minY + (shape.maxY-shape.minY)/2
		var cbx = target.minX + (target.maxX-target.minX)/2
		var cby = target.minY + (target.maxY-target.minY)/2
		var corners = shape.CornerPoints()
		var targetCorners = target.CornerPoints()

		corners = corners[:len(corners)-1]
		targetCorners = targetCorners[:len(targetCorners)-1]

		var projectPolygon = func(axisX, axisY float32, points [][2]float32) (min, max float32) {
			var dot float32 = axisX*points[0][0] + axisY*points[0][1]
			min = dot
			max = dot

			for i := range points {
				var d float32 = axisX*points[i][0] + axisY*points[i][1]
				if d < min {
					min = d
				} else if d > max {
					max = d
				}
			}
			return
		}
		var computeEdges = func(points [][2]float32) [][2]float32 {
			var edges = make([][2]float32, len(points))
			var count int = len(points)
			for i := range count {
				var p1, p2 = points[i], points[(i+1)%count]
				edges[i][0], edges[i][1] = p2[0]-p1[0], p2[1]-p1[1]
			}
			return edges
		}

		var willIntersect = true
		var edgesA = computeEdges(corners)
		var edgesB = computeEdges(targetCorners)
		var edgeCountA = len(edgesA)
		var edgeCountB = len(edgesB)
		var minIntervalDistance = float32(math.Inf(1))
		var tAxisX, tAxisY float32

		for edgeIndex := 0; edgeIndex < edgeCountA+edgeCountB; edgeIndex++ {
			var edgeX, edgeY float32
			if edgeIndex < edgeCountA {
				edgeX, edgeY = edgesA[edgeIndex][0], edgesA[edgeIndex][1]
			} else {
				edgeX, edgeY = edgesB[edgeIndex-edgeCountA][0], edgesB[edgeIndex-edgeCountA][1]
			}

			var axisX, axisY = -edgeY, edgeX
			var axisLen = float32(math.Hypot(float64(axisX), float64(axisY)))
			if axisLen != 0 {
				axisX /= axisLen
				axisY /= axisLen
			} else {
				continue
			}

			var minA, maxA = projectPolygon(axisX, axisY, corners)
			var minB, maxB = projectPolygon(axisX, axisY, targetCorners)

			var velocityProjection = axisX*velocityX + axisY*velocityY
			if velocityProjection < 0 {
				minA += velocityProjection
			} else {
				maxA += velocityProjection
			}

			var iDist = minA - maxB
			if minA < minB {
				iDist = minB - maxA
			}

			if iDist > 0 {
				willIntersect = false
			}

			if !willIntersect {
				break
			}

			var absInterval = float32(math.Abs(float64(iDist)))
			if absInterval < minIntervalDistance {
				minIntervalDistance = absInterval
				tAxisX, tAxisY = axisX, axisY

				var dx = cax - cbx
				var dy = cay - cby
				if dx*tAxisX+dy*tAxisY < 0 {
					tAxisX, tAxisY = -tAxisX, -tAxisY
				}
			}
		}

		if willIntersect {
			velocityX += tAxisX * minIntervalDistance
			velocityY += tAxisY * minIntervalDistance
		}

	}
	return velocityX, velocityY
}

// #region private

func (shape *Shape) internalIsContainingPoint(corners [][2]float32, x, y float32) bool {
	var l = len(corners)
	if l < 3 {
		return false
	}

	var inside = false
	for i := range l {
		var j = (i + 1) % l
		var vi = corners[i]
		var vj = corners[j]

		// check if edge (vi->vj) crosses horizontal ray from (px,py) to +inf X
		if (vi[1] >= y && vj[1] < y) || (vi[1] < y && vj[1] >= y) {
			// x coordinate of intersection of the edge with line y = py
			var xIntersect = (vj[0]-vi[0])*(y-vi[1])/(vj[1]-vi[1]) + vi[0]
			if x < xIntersect {
				inside = !inside
			}
		}
	}
	return inside
}

func (shape *Shape) internalCrossPointsWithLine(corners [][2]float32, line Line) [][2]float32 {
	var result = [][2]float32{}

	for i := 1; i < len(corners); i++ {
		var curLine = NewLine(corners[i-1][0], corners[i-1][1], corners[i][0], corners[i][1])
		var cx, cy = line.CrossPointWithLine(curLine)

		if !math.IsNaN(float64(cx)) && !math.IsNaN(float64(cy)) {
			result = append(result, [2]float32{cx, cy})
		}
	}
	return result
}
func (shape *Shape) internalIsCrossingLine(corners [][2]float32, line Line) bool {
	for i := 1; i < len(corners); i++ {
		var curLine = NewLine(corners[i-1][0], corners[i-1][1], corners[i][0], corners[i][1])
		if line.IsCrossingLine(curLine) {
			return true
		}
	}
	return false
}
func (shape *Shape) internalIsContainingLine(corners [][2]float32, line Line) bool {
	var containsA = shape.internalIsContainingPoint(corners, line.Ax, line.Ay)
	var containsB = shape.internalIsContainingPoint(corners, line.Bx, line.By)
	return containsA && containsB && !shape.internalIsCrossingLine(corners, line)
}
func (shape *Shape) internalIsOverlappingLine(corners [][2]float32, line Line) bool {
	var containsA = shape.internalIsContainingPoint(corners, line.Ax, line.Ay)
	var containsB = shape.internalIsContainingPoint(corners, line.Bx, line.By)
	var crossing = shape.internalIsCrossingLine(corners, line)
	return containsA || containsB || crossing
}

func (shape *Shape) internalCrossPointsWithShape(corners, targetCorners [][2]float32) [][2]float32 {
	var result = [][2]float32{}

	for i := 1; i < len(targetCorners); i++ {
		var line = NewLine(targetCorners[i-1][0], targetCorners[i-1][1], targetCorners[i][0], targetCorners[i][1])
		var pts = shape.internalCrossPointsWithLine(corners, line)
		result = append(result, pts...)
	}
	return result
}
func (shape *Shape) internalIsCrossingShape(corners, targetCorners [][2]float32) bool {
	for i := 1; i < len(targetCorners); i++ {
		var line = NewLine(targetCorners[i-1][0], targetCorners[i-1][1], targetCorners[i][0], targetCorners[i][1])
		if shape.internalIsCrossingLine(corners, line) {
			return true
		}
	}
	return false
}
func (shape *Shape) internalIsContainingShapes(corners, targetCorners [][2]float32) bool {
	for i := 1; i < len(targetCorners); i++ {
		var line = NewLine(targetCorners[i-1][0], targetCorners[i-1][1], targetCorners[i][0], targetCorners[i][1])
		if !shape.internalIsContainingLine(corners, line) {
			return false
		}
	}
	return true
}
func (shape *Shape) internalIsOverlappingShape(corners, targetCorners [][2]float32, target *Shape) bool {
	// overlap happens when:
	// 		one of shape's corners is within target
	//		one of target's corners is within shape
	// 		or there is a crossing

	for i := 1; i < len(targetCorners); i++ { // crossing + target inside shape checks
		var line = NewLine(targetCorners[i-1][0], targetCorners[i-1][1], targetCorners[i][0], targetCorners[i][1])
		if shape.internalIsOverlappingLine(corners, line) {
			return true
		}
	}

	for i := 1; i < len(corners); i++ { // skipping crossing & straight to shape inside target check
		var line = NewLine(corners[i-1][0], corners[i-1][1], corners[i][0], corners[i][1])
		if target.internalIsContainingLine(targetCorners, line) {
			return true
		}
	}

	return false
}

// these methods are the fastest way to discard a slow check
// they rely on having CornerPoints() called beforehand

func (shape *Shape) inBoundingBoxPoint(x, y float32) bool {
	return x >= shape.minX && x <= shape.maxX && y >= shape.minY && y <= shape.maxY
}
func (shape *Shape) inBoundingBoxLine(line Line) bool {
	return shape.inBoundingBoxPoint(line.Ax, line.Ay) || shape.inBoundingBoxPoint(line.Bx, line.By)
}
func (shape *Shape) inBoundingBoxShape(target Shape) bool {
	return shape.minX <= target.maxX && shape.maxX >= target.minX &&
		shape.minY <= target.maxY && shape.maxY >= target.minY
}

// #endregion
