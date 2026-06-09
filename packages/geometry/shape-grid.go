package geometry

import (
	"pure-game-kit/packages/utility/number"
)

type ShapeGrid struct {
	shapes    []Shape
	cells     map[[2]int][]int
	chunkSize float32
}

// For fastest results, use 2 times the max size object in chunk size.
// Similarly sized objects are preferred. Use multiple grids otherwise.
func NewShapeGrid(chunkSize float32) ShapeGrid {
	return ShapeGrid{chunkSize: chunkSize, cells: make(map[[2]int][]int)}
}

//=================================================================

func (s *ShapeGrid) AddShapes(shapes ...Shape) {
	var startIndex = len(s.shapes)
	s.shapes = append(s.shapes, shapes...)

	for i, sh := range shapes {
		var actualIndex = startIndex + i
		var startX, startY, endX, endY = s.cellsUnderShape(sh)

		for cy := startY; cy <= endY; cy++ {
			for cx := startX; cx <= endX; cx++ {
				var cenX, cenY = (float32(cx) * s.chunkSize) + (s.chunkSize * 0.5), (float32(cy) * s.chunkSize) + (s.chunkSize * 0.5)
				var cellShape = NewRectangle(cenX, cenY, s.chunkSize, s.chunkSize, 0)
				if sh.Overlaps(cellShape) {
					var key = [2]int{cx, cy}
					s.cells[key] = append(s.cells[key], actualIndex)
				}
			}
		}
	}
}
func (s *ShapeGrid) RemoveShapes(shapes ...Shape) {
	for _, target := range shapes {
		var startX, startY, endX, endY = s.cellsUnderShape(target)
		var targetIdx = -1 // Track the index to zero out

		for cy := startY; cy <= endY; cy++ {
			for cx := startX; cx <= endX; cx++ {
				var cenX, cenY = (float32(cx) * s.chunkSize) + (s.chunkSize * 0.5), (float32(cy) * s.chunkSize) + (s.chunkSize * 0.5)
				var cellShape = NewRectangle(cenX, cenY, s.chunkSize, s.chunkSize, 0)
				if !target.Overlaps(cellShape) {
					continue
				}

				var key = [2]int{cx, cy}
				var indices, has = s.cells[key]
				if !has {
					continue
				}

				var newCell []int
				for _, idx := range indices {
					if s.shapes[idx] == target {
						targetIdx = idx // We found the shape's underlying index
					} else {
						newCell = append(newCell, idx) // Keep all other indices
					}
				}

				if len(newCell) == 0 {
					delete(s.cells, key)
				} else {
					s.cells[key] = newCell
				}
			}
		}

		if targetIdx != -1 { // Zero out the shape in the backing array
			s.shapes[targetIdx] = Shape{}
		}
	}
}

//=================================================================

func (s *ShapeGrid) All() []Shape {
	var result = make([]Shape, 0, len(s.shapes))
	for _, sh := range s.shapes {
		if sh != (Shape{}) { // Filter out the removed shapes
			result = append(result, sh)
		}
	}
	return result
}
func (s *ShapeGrid) AtCell(x, y int) []Shape {
	var result []Shape
	var shapes, has = s.cells[[2]int{x, y}]
	if has {
		for _, index := range shapes {
			result = append(result, s.shapes[index])
		}
	}
	return result
}
func (s *ShapeGrid) Neighbors(shape Shape) []Shape {
	if s.chunkSize <= 0 {
		return nil
	}

	var startX, startY, endX, endY = s.cellsUnderShape(shape)
	var uniqueIndices = make(map[int]struct{})

	for cy := startY; cy <= endY; cy++ {
		for cx := startX; cx <= endX; cx++ {
			// var cenX, cenY = (float32(cx) * s.chunkSize) + (s.chunkSize * 0.5), (float32(cy) * s.chunkSize) + (s.chunkSize * 0.5)
			// var cellShape = NewRectangle(cenX, cenY, s.chunkSize, s.chunkSize, 0)
			// if shape.Overlaps(cellShape) {
			var indices, has = s.cells[[2]int{cx, cy}]
			if has {
				for _, idx := range indices {
					uniqueIndices[idx] = struct{}{}
				}
			}
			// }
		}
	}

	var result []Shape
	for idx := range uniqueIndices {
		result = append(result, s.shapes[idx])
	}
	return result
}

// Diagonals take 1.5 cells distance-wise. This way, range calculations are rounded & have no weird left-overs.
// This quirk makes regular 2D distances incorrect, in such cases use s.RangeDistance().
func (s *ShapeGrid) RangeCells(startX, startY int, maxDistance float32, diagonals bool) [][2]int {
	type state struct {
		x, y          int
		remainingDist float32
	}
	var visited = make(map[[2]int]float32)
	var queue = []state{{startX, startY, maxDistance + 0.1}}

	for len(queue) > 0 && s != nil {
		var curr = queue[0]
		queue = queue[1:]
		var currPos = [2]int{curr.x, curr.y}

		if val, exists := visited[currPos]; exists && val >= curr.remainingDist {
			continue
		}
		visited[currPos] = curr.remainingDist

		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				if dx == 0 && dy == 0 {
					continue
				}

				var nextX, nextY = curr.x + dx, curr.y + dy
				var nextPos = [2]int{nextX, nextY}

				if shapes, blocked := s.cells[nextPos]; blocked && len(shapes) > 0 {
					continue
				}

				var cost float32 = 1.0
				if dx != 0 && dy != 0 {
					cost = 1.5

					if !diagonals {
						var s1, b1 = s.cells[[2]int{curr.x + dx, curr.y}]
						var s2, b2 = s.cells[[2]int{curr.x, curr.y + dy}]
						if (b1 && len(s1) > 0) || (b2 && len(s2) > 0) {
							continue
						}
					}
				}

				var nextRemaining = curr.remainingDist - cost
				if nextRemaining >= 0 {
					queue = append(queue, state{nextX, nextY, nextRemaining})
				}
			}
		}
	}

	var result = make([][2]int, 0, len(visited))
	for pos := range visited {
		result = append(result, pos)
	}
	return result
}

// A distorted distance that accounts for diagonals taking 1.5 cells. See s.RangeCells().
func (s *ShapeGrid) RangeDistance(x, y, targetX, targetY int) float32 {
	var dx, dy = number.Absolute(targetX - x), number.Absolute(targetY - y)
	var diag = number.Minimum(dx, dy)
	var straight = number.Maximum(dx, dy) - diag
	return float32(diag)*1.5 + float32(straight)
}

// private ========================================================

func (s *ShapeGrid) cellsUnderShape(shape Shape) (startX, startY, endX, endY int) {
	var minX, minY, w, h = shape.Bounds()
	var maxX, maxY = minX + w, minY + h
	startX, startY = int(number.RoundDown(minX/s.chunkSize)), int(number.RoundDown(minY/s.chunkSize))
	endX, endY = int(number.RoundDown(maxX/s.chunkSize)), int(number.RoundDown(maxY/s.chunkSize))
	return startX, startY, endX, endY
}
