package geometry

import (
	"pure-game-kit/packages/utility/number"
)

type ShapeGrid struct {
	cells    map[[2]int][]Shape
	cellSize float32
}

func NewShapeGrid(cellSize float32) *ShapeGrid {
	return &ShapeGrid{cellSize: cellSize, cells: make(map[[2]int][]Shape)}
}

//=================================================================

func (s *ShapeGrid) SetAtCell(x, y int, shapes ...Shape) {
	var key = [2]int{x, y}
	s.cells[key] = []Shape{}
	s.cells[key] = append(s.cells[key], shapes...)
}

//=================================================================

func (s *ShapeGrid) All() []Shape {
	var result = []Shape{}
	for k := range s.cells {
		result = append(result, s.AtCell(k[0], k[1])...)
	}
	return result
}
func (s *ShapeGrid) AtCell(x, y int) []Shape {
	var shapes, has = s.cells[[2]int{x, y}]
	if has {
		return shapes
	}
	return []Shape{}
}
func (s *ShapeGrid) Neighbors(shape Shape) []Shape {
	if s.cellSize <= 0 {
		return nil
	}
	var minX, minY, w, h = shape.Bounds()
	var result, maxX, maxY = []Shape{}, minX + w, minY + h
	var startX, startY = int(number.RoundDown(minX / s.cellSize)), int(number.RoundDown(minY / s.cellSize))
	var endX, endY = int(number.RoundDown(maxX / s.cellSize)), int(number.RoundDown(maxY / s.cellSize))
	for cx := startX; cx <= endX; cx++ {
		for cy := startY; cy <= endY; cy++ {
			result = append(result, s.AtCell(cx, cy)...)
		}
	}
	return result
}

// Diagonals take 1.5 cells distance-wise. This way, range calculations are rounded & have no weird left-overs.
// This quirk makes regular 2D distances incorrect, instead use:
//
//	shapeGrid.RangeDistance(...)
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
func (s *ShapeGrid) RangeDistance(x, y, targetX, targetY int) float32 {
	var dx, dy = number.Absolute(targetX - x), number.Absolute(targetY - y)
	var diag = number.Minimum(dx, dy)
	var straight = number.Maximum(dx, dy) - diag
	return float32(diag)*1.5 + float32(straight)
}
