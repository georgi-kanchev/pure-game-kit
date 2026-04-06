package geometry

import (
	"pure-game-kit/internal"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/point"
)

type Shape struct {
	X, Y, Angle, ScaleX, ScaleY,

	minX, minY, maxX, maxY,
	gridX, gridY float32

	corners []float32
}

// Shapes separated by NaN pair.
func NewShapes(corners ...float32) []*Shape {
	var curPts []float32
	var result []*Shape
	for i := 0; i < len(corners); i += 2 {
		var x, y = corners[i], corners[i+1]
		if number.IsNaN(x) || number.IsNaN(y) {
			if len(curPts) > 0 {
				result = append(result, NewShapeCorners(curPts...))
				curPts = []float32{}
			}
			continue
		}
		curPts = append(curPts, x, y)
	}
	if len(curPts) > 0 {
		result = append(result, NewShapeCorners(curPts...))
	}
	return result
}
func NewShapeCorners(corners ...float32) *Shape {
	if len(corners) == 0 {
		return &Shape{}
	}

	var n = len(corners) // close shape if it already isn't (compare first pair to last pair)
	if n >= 4 && (corners[0] != corners[n-2] || corners[1] != corners[n-1]) {
		corners = append(corners, corners[0], corners[1])
	}

	return &Shape{ScaleX: 1, ScaleY: 1, corners: corners}
}
func NewShapeSides(radius float32, sides int) *Shape {
	var corners []float32
	var step float32 = 360.0 / float32(sides)
	for i := range sides {
		var x, y = point.MoveAtAngle(0, 0, step*float32(i)-90, radius)
		corners = append(corners, x, y)
	}

	if len(corners) >= 2 {
		corners = append(corners, corners[0], corners[1])
	}

	return &Shape{ScaleX: 1, ScaleY: 1, corners: corners}
}
func NewShapeQuad(width, height, pivotX, pivotY float32) *Shape {
	var offX, offY = width * pivotX, height * pivotY
	var corners = []float32{
		-offX, -offY,
		width - offX, -offY,
		width - offX, height - offY,
		-offX, height - offY,
		-offX, -offY,
	}
	return &Shape{ScaleX: 1, ScaleY: 1, corners: corners}
}
func NewShapeQuadRounded(width, height, radius, pivotX, pivotY float32, segments int) *Shape {
	var maxRadius = width / 2
	if height/2 < maxRadius {
		maxRadius = height / 2
	}
	if radius > maxRadius {
		radius = maxRadius
	}

	const pi = 3.14159
	var offX, offY = width * pivotX, height * pivotY
	var corners []float32
	var addCorner = func(cx, cy, startAngle float32) {
		for i := 0; i <= segments; i++ {
			var angle = startAngle + (float32(i)/float32(segments))*(pi/2)
			var px = cx + radius*number.Cosine(angle)
			var py = cy + radius*number.Sine(angle)
			corners = append(corners, px-offX, py-offY)
		}
	}
	addCorner(width-radius, radius, -pi/2)    // top right
	addCorner(width-radius, height-radius, 0) // bottom right
	addCorner(radius, height-radius, pi/2)    // bottom left
	addCorner(radius, radius, pi)             // top left

	// Close the shape
	corners = append(corners, corners[0], corners[1])
	return &Shape{ScaleX: 1, ScaleY: 1, corners: corners}
}
func NewShapeEllipse(width, height float32, segments int) *Shape {
	if segments < 3 {
		segments = 3
	}

	var rx, ry = width / 2, height / 2
	// Pre-allocate space for x,y pairs plus closing pair
	var corners = make([]float32, 0, (segments+1)*2)
	var step = 360.0 / float32(segments)

	for i := 0; i < segments; i++ {
		var cx, cy = point.MoveAtAngle(0, 0, float32(i)*step, 1)
		corners = append(corners, cx*rx, cy*ry)
	}
	corners = append(corners, corners[0], corners[1])

	return &Shape{ScaleX: 1, ScaleY: 1, corners: corners}
}

//=================================================================

func (s *Shape) CornerPoints() []float32 {
	if s == nil {
		return nil
	}

	var result = make([]float32, len(s.corners))
	s.minX, s.minY = number.Infinity(), number.Infinity()
	s.maxX, s.maxY = number.NegativeInfinity(), number.NegativeInfinity()

	var sin, cos = internal.SinCos(s.Angle)

	for i := 0; i < len(s.corners); i += 2 {
		var x, y = s.corners[i], s.corners[i+1]

		x *= s.ScaleX
		y *= s.ScaleY

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

		result[i], result[i+1] = resultX, resultY
	}
	return result
}
func (s *Shape) Collide(velocityX, velocityY float32, targets ...*Shape) (newVelocityX, newVelocityY float32) {
	for _, target := range targets {
		if s.minX == 0 && s.minY == 0 && s.maxX == 0 && s.maxY == 0 {
			s.CornerPoints()
		}
		if target.minX == 0 && target.minY == 0 && target.maxX == 0 && target.maxY == 0 {
			target.CornerPoints()
		}

		if !s.inBoundingBoxShape(*target) {
			continue
		}

		var cax = s.minX + (s.maxX-s.minX)/2
		var cay = s.minY + (s.maxY-s.minY)/2
		var cbx = target.minX + (target.maxX-target.minX)/2
		var cby = target.minY + (target.maxY-target.minY)/2

		var corners = s.CornerPoints()
		var targetCorners = target.CornerPoints()

		// Remove closing point for SAT calculation
		if len(corners) >= 4 {
			corners = corners[:len(corners)-2]
		}
		if len(targetCorners) >= 4 {
			targetCorners = targetCorners[:len(targetCorners)-2]
		}

		var projectPolygon = func(axisX, axisY float32, points []float32) (min, max float32) {
			var dot float32 = axisX*points[0] + axisY*points[1]
			min, max = dot, dot

			for i := 2; i < len(points); i += 2 {
				var d float32 = axisX*points[i] + axisY*points[i+1]
				if d < min {
					min = d
				} else if d > max {
					max = d
				}
			}
			return
		}

		var willIntersect = true
		var minIntervalDistance = number.Infinity()
		var tAxisX, tAxisY float32

		// Combine edge checking loop
		var totalPoints = len(corners) + len(targetCorners)
		for i := 0; i < totalPoints; i += 2 {
			var x1, y1, x2, y2 float32

			if i < len(corners) {
				x1, y1 = corners[i], corners[i+1]
				var next = (i + 2) % len(corners)
				x2, y2 = corners[next], corners[next+1]
			} else {
				var idx = i - len(corners)
				x1, y1 = targetCorners[idx], targetCorners[idx+1]
				var next = (idx + 2) % len(targetCorners)
				x2, y2 = targetCorners[next], targetCorners[next+1]
			}

			// Normal of the edge
			var axisX, axisY = -(y2 - y1), (x2 - x1)
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
				break
			}

			var absInterval = number.Unsign(iDist)
			if absInterval < minIntervalDistance {
				minIntervalDistance = absInterval
				tAxisX, tAxisY = axisX, axisY

				if (cax-cbx)*tAxisX+(cay-cby)*tAxisY < 0 {
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

func (s *Shape) CrossPointsWithLines(lines ...Line) []float32 {
	var corners = s.CornerPoints()
	var result []float32

	for _, line := range lines {
		if s.inBoundingBoxLine(line) {
			var pts = s.internalCrossPointsWithLine(corners, line)
			result = append(result, pts...)
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

func (s *Shape) CrossPointsWithShapes(shapes ...*Shape) []float32 {
	var corners = s.CornerPoints()
	var result []float32

	for _, target := range shapes {
		var targetCorners = target.CornerPoints()
		if s.inBoundingBoxShape(*target) {
			var pts = s.internalCrossPointsWithShape(corners, targetCorners)
			result = append(result, pts...)
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
