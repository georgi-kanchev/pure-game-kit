package geometry

import "pure-game-kit/utility/number"

func (s *Shape) internalIsContainingPoint(corners []float32, x, y float32) bool {
	var l = len(corners)
	if l < 6 { // A closed shape needs at least 3 points + 1 closing point (8 floats)
		return false
	}

	var inside = false
	// Step by 2 to process x,y pairs
	for i := 0; i < l-2; i += 2 {
		var xi, yi = corners[i], corners[i+1]
		var xj, yj = corners[i+2], corners[i+3]

		// check if edge (vi->vj) crosses horizontal ray from (x,y) to +inf X
		if (yi >= y && yj < y) || (yi < y && yj >= y) {
			// x coordinate of intersection of the edge with line y = py
			var xIntersect = (xj-xi)*(y-yi)/(yj-yi) + xi
			if x < xIntersect {
				inside = !inside
			}
		}
	}
	return inside
}

func (s *Shape) internalCrossPointsWithLine(corners []float32, line Line) []float32 {
	var result []float32

	for i := 2; i < len(corners); i += 2 {
		var curLine = NewLine(corners[i-2], corners[i-1], corners[i], corners[i+1])
		var cx, cy = line.CrossPointWithLine(curLine)

		if !number.IsNaN(cx) && !number.IsNaN(cy) {
			result = append(result, cx, cy)
		}
	}
	return result
}
func (s *Shape) internalIsCrossingLine(corners []float32, line Line) bool {
	for i := 2; i < len(corners); i += 2 {
		var curLine = NewLine(corners[i-2], corners[i-1], corners[i], corners[i+1])
		if line.IsCrossingLine(curLine) {
			return true
		}
	}
	return false
}
func (s *Shape) internalIsContainingLine(corners []float32, line Line) bool {
	var containsA = s.internalIsContainingPoint(corners, line.Ax, line.Ay)
	var containsB = s.internalIsContainingPoint(corners, line.Bx, line.By)
	return containsA && containsB && !s.internalIsCrossingLine(corners, line)
}
func (s *Shape) internalIsOverlappingLine(corners []float32, line Line) bool {
	var containsA = s.internalIsContainingPoint(corners, line.Ax, line.Ay)
	var containsB = s.internalIsContainingPoint(corners, line.Bx, line.By)
	var crossing = s.internalIsCrossingLine(corners, line)
	return containsA || containsB || crossing
}

func (s *Shape) internalCrossPointsWithShape(corners, targetCorners []float32) []float32 {
	var result []float32

	for i := 2; i < len(targetCorners); i += 2 {
		var line = NewLine(targetCorners[i-2], targetCorners[i-1], targetCorners[i], targetCorners[i+1])
		var pts = s.internalCrossPointsWithLine(corners, line)
		result = append(result, pts...)
	}
	return result
}
func (s *Shape) internalIsCrossingShape(corners, targetCorners []float32) bool {
	for i := 2; i < len(targetCorners); i += 2 {
		var line = NewLine(targetCorners[i-2], targetCorners[i-1], targetCorners[i], targetCorners[i+1])
		if s.internalIsCrossingLine(corners, line) {
			return true
		}
	}
	return false
}
func (s *Shape) internalIsContainingShapes(corners, targetCorners []float32) bool {
	for i := 2; i < len(targetCorners); i += 2 {
		var line = NewLine(targetCorners[i-2], targetCorners[i-1], targetCorners[i], targetCorners[i+1])
		if !s.internalIsContainingLine(corners, line) {
			return false
		}
	}
	return true
}
func (s *Shape) internalIsOverlappingShape(corners, targetCorners []float32, target *Shape) bool {
	// crossing + target inside shape checks
	for i := 2; i < len(targetCorners); i += 2 {
		var line = NewLine(targetCorners[i-2], targetCorners[i-1], targetCorners[i], targetCorners[i+1])
		if s.internalIsOverlappingLine(corners, line) {
			return true
		}
	}

	// skipping crossing & straight to shape inside target check
	for i := 2; i < len(corners); i += 2 {
		var line = NewLine(corners[i-2], corners[i-1], corners[i], corners[i+1])
		if target.internalIsContainingLine(targetCorners, line) {
			return true
		}
	}

	return false
}

//=================================================================
// these methods are the fastest way to discard a slow check
// they rely on having CornerPoints() called beforehand

func (s *Shape) inBoundingBoxPoint(x, y float32) bool {
	return x >= s.minX && x <= s.maxX && y >= s.minY && y <= s.maxY
}
func (s *Shape) inBoundingBoxLine(line Line) bool {
	return s.inBoundingBoxPoint(line.Ax, line.Ay) || s.inBoundingBoxPoint(line.Bx, line.By)
}
func (s *Shape) inBoundingBoxShape(target Shape) bool {
	return s.minX <= target.maxX && s.maxX >= target.minX &&
		s.minY <= target.maxY && s.maxY >= target.minY
}
