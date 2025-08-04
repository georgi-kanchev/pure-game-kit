package geometry

import (
	"math"
	"pure-kit/engine/geometry/point"
	"pure-kit/engine/utility/angle"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Shape struct {
	X, Y, Angle    float32
	ScaleX, ScaleY float32
	corners        [][2]float32
}

func NewShape(corners ...[2]float32) Shape {
	if len(corners) == 0 {
		return Shape{}
	}
	return Shape{ScaleX: 1, ScaleY: 1, corners: append(corners, corners[0])}
}
func NewRectangle(width, height, pivotX, pivotY float32) Shape {
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
func NewCircle(radius float32, segments int) Shape {
	var corners = [][2]float32{}
	var step float32 = 360.0 / float32(segments)
	for i := range segments {
		var x, y = point.MoveAtAngle(0, 0, step*float32(i), radius)
		corners = append(corners, [2]float32{x, y})
	}

	if len(corners) > 0 {
		corners = append(corners, corners[0])
	}

	return Shape{ScaleX: 1, ScaleY: 1, corners: corners}
}

func (shape *Shape) CornerPoints() [][2]float32 {
	var result = make([][2]float32, len(shape.corners))
	for i := range shape.corners {
		var x, y = shape.corners[i][0], shape.corners[i][1]

		x *= shape.ScaleX
		y *= shape.ScaleY

		var rad = angle.ToRadians(shape.Angle)
		var sin, cos = float32(math.Sin(float64(rad))), float32(math.Cos(float64(rad)))
		var resultX = shape.X + (x*cos - y*sin)
		var resultY = shape.Y + (x*sin + y*cos)

		result[i] = [2]float32{resultX, resultY}
	}
	return result
}

// all check methods are made with speed in mind, not so much readability
// they should have the least allocations (once per API call for getPoints() & CornerPoints())
// and the least iterations (1 loop for 1 shape and 2 loops for 2 shapes)
// don't be a smartass by "simplifying" and reusing them internally in the future

func (shape *Shape) CrossPointsWithLine(line Line) [][2]float32 {
	return shape.internalCrossPointsWithLine(shape.CornerPoints(), line)
}
func (shape *Shape) CrossPointsWithShape(target *Shape) [][2]float32 {
	var corners = shape.CornerPoints()
	var targetCorners = target.CornerPoints()
	var result = [][2]float32{}

	for i := 1; i < len(targetCorners); i++ {
		var line = NewLine(targetCorners[i-1][0], targetCorners[i-1][1], targetCorners[i][0], targetCorners[i][1])
		var pts = shape.internalCrossPointsWithLine(corners, line)
		result = append(result, pts...)
	}
	return result
}

func (shape *Shape) IsContainingPoint(x, y float32) bool {
	return shape.internalIsContainingPoint(shape.getPoints(shape.CornerPoints()), x, y)
}

func (shape *Shape) IsCrossingLine(line Line) bool {
	return shape.internalIsCrossingLine(shape.CornerPoints(), line)
}
func (shape *Shape) IsContainingLine(line Line) bool {
	var corners = shape.CornerPoints()
	var points = shape.getPoints(corners)
	return shape.internalIsContainingLine(corners, points, line)
}
func (shape *Shape) IsOverlappingLine(line Line) bool {
	var corners = shape.CornerPoints()
	var points = shape.getPoints(corners)
	return shape.internalIsOverlappingLine(corners, points, line)
}

func (shape *Shape) IsCrossingShape(target *Shape) bool {
	var corners = shape.CornerPoints()
	var targetCorners = target.CornerPoints()

	for i := 1; i < len(targetCorners); i++ {
		var line = NewLine(targetCorners[i-1][0], targetCorners[i-1][1], targetCorners[i][0], targetCorners[i][1])
		if shape.internalIsCrossingLine(corners, line) {
			return true
		}
	}
	return false
}
func (shape *Shape) IsContainingShape(target *Shape) bool {
	var corners = shape.CornerPoints()
	var points = shape.getPoints(corners)
	var targetCorners = target.CornerPoints()

	for i := 1; i < len(targetCorners); i++ {
		var line = NewLine(targetCorners[i-1][0], targetCorners[i-1][1], targetCorners[i][0], targetCorners[i][1])
		if !shape.internalIsContainingLine(corners, points, line) {
			return false
		}
	}
	return true
}
func (shape *Shape) IsOverlappingShape(target *Shape) bool {
	var corners = shape.CornerPoints()
	var points = shape.getPoints(corners)
	var targetCorners = target.CornerPoints()
	var targetPoints = target.getPoints(targetCorners)

	// overlap happens when:
	// 		one of shape's corners is within target
	//		one of target's corners is within shape
	// 		or there is a crossing

	for i := 1; i < len(targetCorners); i++ { // crossing + target inside shape checks
		var line = NewLine(targetCorners[i-1][0], targetCorners[i-1][1], targetCorners[i][0], targetCorners[i][1])
		if shape.internalIsOverlappingLine(corners, points, line) {
			return true
		}
	}

	for i := 1; i < len(corners); i++ { // skipping crossing & straight to shape inside target check
		var line = NewLine(corners[i-1][0], corners[i-1][1], corners[i][0], corners[i][1])
		if target.internalIsContainingLine(targetCorners, targetPoints, line) {
			return true
		}
	}

	return false
}

// #region private

func (shape *Shape) internalIsContainingPoint(points []rl.Vector2, x, y float32) bool {
	return rl.CheckCollisionPointPoly(rl.Vector2{X: x, Y: y}, points)
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
func (shape *Shape) internalIsContainingLine(corners [][2]float32, points []rl.Vector2, line Line) bool {
	var containsA = rl.CheckCollisionPointPoly(rl.Vector2{X: line.Ax, Y: line.Ay}, points)
	var containsB = rl.CheckCollisionPointPoly(rl.Vector2{X: line.Bx, Y: line.By}, points)
	return containsA && containsB && !shape.internalIsCrossingLine(corners, line)
}
func (shape *Shape) internalIsOverlappingLine(corners [][2]float32, points []rl.Vector2, line Line) bool {
	var containsA = shape.internalIsContainingPoint(points, line.Ax, line.Ay)
	var containsB = shape.internalIsContainingPoint(points, line.Bx, line.By)
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

func (shape *Shape) getPoints(corners [][2]float32) []rl.Vector2 {
	var result = make([]rl.Vector2, len(shape.corners))
	for i, p := range corners {
		result[i] = rl.Vector2{X: p[0], Y: p[1]}
	}
	return result
}

// #endregion
