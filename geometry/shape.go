package geometry

import (
	"pure-game-kit/utility/angle"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/point"
)

type Shape struct {
	X, Y, Angle, ScaleX, ScaleY,
	minX, minY, maxX, maxY,
	gridX, gridY float32

	corners [][2]float32
}

func NewShapeCorners(corners ...[2]float32) *Shape {
	if len(corners) == 0 {
		return &Shape{}
	}

	if corners[0] != corners[len(corners)-1] {
		corners = append(corners, corners[0])
	} // close shape if it already isn't

	return &Shape{ScaleX: 1, ScaleY: 1, corners: corners}
}
func NewShapeSides(radius float32, sides int) *Shape {
	var corners = [][2]float32{}
	var step float32 = 360.0 / float32(sides)
	for i := range sides {
		var x, y = point.MoveAtAngle(0, 0, step*float32(i)-90, radius)
		corners = append(corners, [2]float32{x, y})
	}

	if len(corners) > 0 {
		corners = append(corners, corners[0])
	}

	return &Shape{ScaleX: 1, ScaleY: 1, corners: corners}
}
func NewShapeRectangle(width, height, pivotX, pivotY float32) *Shape {
	var offX, offY = width * pivotX, height * pivotY
	var corners = [][2]float32{
		{-offX, -offY},
		{width - offX, -offY},
		{width - offX, height - offY},
		{-offX, height - offY},
		{-offX, -offY},
	}
	return &Shape{ScaleX: 1, ScaleY: 1, corners: corners}
}
func NewShapeEllipse(width, height float32, segments int) *Shape {
	if segments < 3 {
		segments = 3
	}

	var rx, ry = width / 2, height / 2
	var corners = make([][2]float32, 0, segments+1)
	var step = 360.0 / float32(segments)

	for i := 0; i < segments; i++ {
		var cx, cy = point.MoveAtAngle(0, 0, float32(i)*step, 1)
		corners = append(corners, [2]float32{cx * rx, cy * ry})
	}
	corners = append(corners, corners[0])

	return &Shape{ScaleX: 1, ScaleY: 1, corners: corners}
}

//=================================================================

func (s *Shape) CornerPoints() [][2]float32 {
	if s == nil {
		return nil
	}

	var result = make([][2]float32, len(s.corners))
	s.minX, s.minY = number.Infinity(), number.Infinity()
	s.maxX, s.maxY = number.NegativeInfinity(), number.NegativeInfinity()

	for i := range s.corners {
		var x, y = s.corners[i][0], s.corners[i][1]

		x *= s.ScaleX
		y *= s.ScaleY

		var rad = angle.ToRadians(s.Angle)
		var sin, cos = number.Sine(rad), number.Cosine(rad)
		var resultX = s.gridX + s.X + (x*cos - y*sin)
		var resultY = s.gridY + s.Y + (x*sin + y*cos)

		if s.minX > resultX {
			s.minX = resultX
		}
		if s.minY > resultY {
			s.minY = resultY
		}
		if s.maxX < resultX {
			s.maxX = resultX
		}
		if s.maxY < resultY {
			s.maxY = resultY
		}

		result[i] = [2]float32{resultX, resultY}
	}
	return result
}

func (s *Shape) Collide(velocityX, velocityY float32, targets ...*Shape) (newVelocityX, newVelocityY float32) {
	for _, target := range targets {
		var corners [][2]float32
		var targetCorners [][2]float32
		if s.minX == 0 && s.minY == 0 && s.maxX == 0 && s.maxY == 0 {
			corners = s.CornerPoints()
		}
		if target.minX == 0 && target.minY == 0 && target.maxX == 0 && target.maxY == 0 {
			corners = target.CornerPoints()
		}

		if !s.inBoundingBoxShape(*target) {
			continue
		}

		var cax = s.minX + (s.maxX-s.minX)/2
		var cay = s.minY + (s.maxY-s.minY)/2
		var cbx = target.minX + (target.maxX-target.minX)/2
		var cby = target.minY + (target.maxY-target.minY)/2
		corners = s.CornerPoints()
		targetCorners = target.CornerPoints()
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
		var minIntervalDistance = number.Infinity()
		var tAxisX, tAxisY float32

		for edgeIndex := 0; edgeIndex < edgeCountA+edgeCountB; edgeIndex++ {
			var edgeX, edgeY float32
			if edgeIndex < edgeCountA {
				edgeX, edgeY = edgesA[edgeIndex][0], edgesA[edgeIndex][1]
			} else {
				edgeX, edgeY = edgesB[edgeIndex-edgeCountA][0], edgesB[edgeIndex-edgeCountA][1]
			}

			var axisX, axisY = -edgeY, edgeX
			var axisLen = number.SquareRoot(axisX*axisX + axisY*axisY)
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

			var absInterval = number.Unsign(iDist)
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

//=================================================================
// all check methods are made with speed in mind, not so much readability
// they should have the least iterations and allocations (once per API call for CornerPoints())
// don't be a smartass by "simplifying" and reusing them internally in the future

func (s *Shape) IsContainingPoint(x, y float32) bool {
	var corners = s.CornerPoints()

	if !s.inBoundingBoxPoint(x, y) {
		return false
	}

	return s.internalIsContainingPoint(corners, x, y)
}

func (s *Shape) CrossPointsWithLines(lines ...Line) [][2]float32 {
	var corners = s.CornerPoints()
	var result = [][2]float32{}

	for _, line := range lines {
		if s.inBoundingBoxLine(line) {
			result = append(result, s.internalCrossPointsWithLine(corners, line)...)
		}
	}
	return result
}
func (s *Shape) IsCrossingLines(lines ...Line) bool {
	var corners = s.CornerPoints()
	for _, line := range lines {
		if s.inBoundingBoxLine(line) && s.internalIsCrossingLine(corners, line) {
			return true
		}
	}
	return false
}
func (s *Shape) IsContainingLines(lines ...Line) bool {
	var corners = s.CornerPoints()
	for _, line := range lines {
		if !s.inBoundingBoxLine(line) || !s.internalIsContainingLine(corners, line) {
			return false
		}
	}
	return true
}
func (s *Shape) IsOverlappingLines(lines ...Line) bool {
	var corners = s.CornerPoints()
	for _, line := range lines {
		if s.inBoundingBoxLine(line) && s.internalIsOverlappingLine(corners, line) {
			return true
		}
	}
	return false
}

func (s *Shape) CrossPointsWithShapes(shapes ...*Shape) [][2]float32 {
	var corners = s.CornerPoints()
	var result = [][2]float32{}

	for _, target := range shapes {
		var targetCorners = target.CornerPoints()
		if s.inBoundingBoxShape(*target) {
			result = append(result, s.internalCrossPointsWithShape(corners, targetCorners)...)
		}
	}
	return result
}
func (s *Shape) IsCrossingShapes(shapes ...*Shape) bool {
	var corners = s.CornerPoints()
	for _, target := range shapes {
		var targetCorners = target.CornerPoints()
		if s.inBoundingBoxShape(*target) && s.internalIsCrossingShape(corners, targetCorners) {
			return true
		}
	}
	return false
}
func (s *Shape) IsContainingShapes(shapes ...*Shape) bool {
	var corners = s.CornerPoints()
	for _, target := range shapes {
		var targetCorners = target.CornerPoints()
		if !s.inBoundingBoxShape(*target) || !s.internalIsContainingShapes(corners, targetCorners) {
			return false
		}
	}
	return true
}
func (s *Shape) IsOverlappingShapes(shapes ...*Shape) bool {
	var corners = s.CornerPoints()

	for _, target := range shapes {
		var targetCorners = target.CornerPoints()
		if s.inBoundingBoxShape(*target) && s.internalIsOverlappingShape(corners, targetCorners, target) {
			return true
		}
	}
	return false
}
