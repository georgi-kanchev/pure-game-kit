package geometry

import (
	"math"
	"pure-kit/engine/geometry/point"
	"pure-kit/engine/utility/angle"
)

type Shape struct {
	X, Y, Angle    float32
	ScaleX, ScaleY float32

	MinX, MinY float32
	MaxX, MaxY float32

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
	shape.MinX, shape.MinY = float32(math.Inf(1)), float32(math.Inf(1))
	shape.MaxX, shape.MaxY = float32(math.Inf(-1)), float32(math.Inf(-1))

	for i := range shape.corners {
		var x, y = shape.corners[i][0], shape.corners[i][1]

		x *= shape.ScaleX
		y *= shape.ScaleY

		var rad = angle.ToRadians(shape.Angle)
		var sin, cos = float32(math.Sin(float64(rad))), float32(math.Cos(float64(rad)))
		var resultX = shape.X + (x*cos - y*sin)
		var resultY = shape.Y + (x*sin + y*cos)

		if shape.MinX > resultX {
			shape.MinX = resultX
		}
		if shape.MinY > resultY {
			shape.MinY = resultY
		}
		if shape.MaxX < resultX {
			shape.MaxX = resultX
		}
		if shape.MaxY < resultY {
			shape.MaxY = resultY
		}

		result[i] = [2]float32{resultX, resultY}
	}
	return result
}

// all check methods are made with speed in mind, not so much readability
// they should have the least allocations (once per API call for getPoints() & CornerPoints())
// and the least iterations (1 loop for 1 shape and 2 loops for 2 shapes)
// don't be a smartass by "simplifying" and reusing them internally in the future

func (shape *Shape) CrossPointsWithLine(line Line) [][2]float32 {
	var corners = shape.CornerPoints()

	if !shape.inBoundingBoxLine(line) {
		return [][2]float32{}
	}

	return shape.internalCrossPointsWithLine(corners, line)
}
func (shape *Shape) CrossPointsWithShape(target *Shape) [][2]float32 {
	var corners = shape.CornerPoints()
	var targetCorners = target.CornerPoints()
	if !shape.inBoundingBoxShape(*target) {
		return [][2]float32{}
	}

	return shape.internalCrossPointsWithShape(corners, targetCorners)
}

func (shape *Shape) IsContainingPoint(x, y float32) bool {
	var corners = shape.CornerPoints()

	if !shape.inBoundingBoxPoint(x, y) {
		return false
	}

	return shape.internalIsContainingPoint(corners, x, y)
}

func (shape *Shape) IsCrossingLine(line Line) bool {
	var corners = shape.CornerPoints()

	if !shape.inBoundingBoxLine(line) {
		return false
	}

	return shape.internalIsCrossingLine(corners, line)
}
func (shape *Shape) IsContainingLine(line Line) bool {
	var corners = shape.CornerPoints()

	if !shape.inBoundingBoxLine(line) {
		return false
	}

	return shape.internalIsContainingLine(corners, line)
}
func (shape *Shape) IsOverlappingLine(line Line) bool {
	var corners = shape.CornerPoints()

	if !shape.inBoundingBoxLine(line) {
		return false
	}

	return shape.internalIsOverlappingLine(corners, line)
}

func (shape *Shape) IsCrossingShape(target *Shape) bool {
	var corners = shape.CornerPoints()
	var targetCorners = target.CornerPoints()

	if !shape.inBoundingBoxShape(*target) {
		return false
	}

	return shape.internalIsCrossingShape(corners, targetCorners)
}

func (shape *Shape) IsContainingShape(target *Shape) bool {
	var corners = shape.CornerPoints()
	var targetCorners = target.CornerPoints()

	if !shape.inBoundingBoxShape(*target) {
		return false
	}

	for i := 1; i < len(targetCorners); i++ {
		var line = NewLine(targetCorners[i-1][0], targetCorners[i-1][1], targetCorners[i][0], targetCorners[i][1])
		if !shape.internalIsContainingLine(corners, line) {
			return false
		}
	}
	return true
}
func (shape *Shape) IsOverlappingShape(target *Shape) bool {
	var corners = shape.CornerPoints()
	var targetCorners = target.CornerPoints()

	if !shape.inBoundingBoxShape(*target) {
		return false
	}

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

// #region private

func (shape *Shape) internalIsContainingPoint(corners [][2]float32, x, y float32) bool {
	var n = len(corners)
	if n < 3 {
		return false
	}

	var inside = false
	for i := range n {
		var j = (i + 1) % n
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
func (shape *Shape) internalCrossPointsWithShape(corners [][2]float32, targetCorners [][2]float32) [][2]float32 {
	var result = [][2]float32{}

	for i := 1; i < len(targetCorners); i++ {
		var line = NewLine(targetCorners[i-1][0], targetCorners[i-1][1], targetCorners[i][0], targetCorners[i][1])
		var pts = shape.internalCrossPointsWithLine(corners, line)
		result = append(result, pts...)
	}
	return result
}

func (shape *Shape) internalIsCrossingShape(corners [][2]float32, targetCorners [][2]float32) bool {
	for i := 1; i < len(targetCorners); i++ {
		var line = NewLine(targetCorners[i-1][0], targetCorners[i-1][1], targetCorners[i][0], targetCorners[i][1])
		if shape.internalIsCrossingLine(corners, line) {
			return true
		}
	}
	return false
}

// these methods are the fastest way to discard a slow check
// they rely on having CornerPoints() called beforehand

func (shape *Shape) inBoundingBoxPoint(x, y float32) bool {
	return x >= shape.MinX && x <= shape.MaxX && y >= shape.MinY && y <= shape.MaxY
}
func (shape *Shape) inBoundingBoxLine(line Line) bool {
	return shape.inBoundingBoxPoint(line.Ax, line.Ay) || shape.inBoundingBoxPoint(line.Bx, line.By)
}
func (shape *Shape) inBoundingBoxShape(target Shape) bool {
	return shape.MinX <= target.MaxX && shape.MaxX >= target.MinX &&
		shape.MinY <= target.MaxY && shape.MaxY >= target.MinY
}

// #endregion
