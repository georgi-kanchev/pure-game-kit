package geometry

import (
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/number"
)

type ShapeGrid struct {
	cells    map[[2]int][]*Shape
	cellSize int
}

func NewShapeGrid(cellSize int) *ShapeGrid {
	return &ShapeGrid{cellSize: cellSize, cells: make(map[[2]int][]*Shape)}
}

//=================================================================

func (s *ShapeGrid) SetAtCell(x, y int, shapes ...*Shape) {
	var key = [2]int{x, y}
	s.cells[key] = []*Shape{}
	s.cells[key] = append(s.cells[key], shapes...)
}

//=================================================================

func (s *ShapeGrid) All() []*Shape {
	var result = []*Shape{}
	for k := range s.cells {
		result = append(result, s.AtCell(k[0], k[1])...)
	}
	return result
}
func (s *ShapeGrid) AtCell(x, y int) []*Shape {
	var shapes, has = s.cells[[2]int{x, y}]
	if has {
		return shapes
	}
	return []*Shape{}
}
func (s *ShapeGrid) AtPoint(x, y float32) []*Shape {
	var w, h = float32(s.cellSize), float32(s.cellSize)
	if w == 0 || h == 0 {
		return []*Shape{}
	}
	var i, j = number.RoundDown(x / w), number.RoundDown(y / h)
	return s.AtCell(int(i), int(j))
}
func (s *ShapeGrid) AroundLine(line *Shape) []*Shape {
	var w, h = float32(s.cellSize), float32(s.cellSize)
	if w == 0 || h == 0 {
		return []*Shape{}
	}

	var sin, cos = internal.SinCos(line.Angle)
	var hw = line.Width * 0.5
	var ax, ay, bx, by = line.X - cos*hw, line.Y - sin*hw, line.X + cos*hw, line.Y + sin*hw
	var result []*Shape
	var x0, y0, x1, y1 = ax / w, ay / h, bx / w, by / h
	var ix0, iy0 = int(number.RoundDown(x0)), int(number.RoundDown(y0))
	var ix1, iy1 = int(number.RoundDown(x1)), int(number.RoundDown(y1))
	var dx, dy = x1 - x0, y1 - y0
	var stepX, stepY int
	var tMaxX, tMaxY, tDeltaX, tDeltaY float32

	if dx > 0 {
		stepX = 1
		tMaxX = (float32(ix0+1) - x0) / dx
		tDeltaX = 1 / dx
	} else if dx < 0 {
		stepX = -1
		tMaxX = (x0 - float32(ix0)) / -dx
		tDeltaX = 1 / -dx
	} else {
		stepX = 0
		tMaxX = number.Infinity()
	}

	if dy > 0 {
		stepY = 1
		tMaxY = (float32(iy0+1) - y0) / dy
		tDeltaY = 1 / dy
	} else if dy < 0 {
		stepY = -1
		tMaxY = (y0 - float32(iy0)) / -dy
		tDeltaY = 1 / -dy
	} else {
		stepY = 0
		tMaxY = number.Infinity()
	}

	for { // Traverse until reaching the target cell
		result = append(result, s.AtCell(ix0, iy0)...)
		if ix0 == ix1 && iy0 == iy1 {
			break
		}
		if tMaxX < tMaxY {
			ix0 += stepX
			tMaxX += tDeltaX
		} else {
			iy0 += stepY
			tMaxY += tDeltaY
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
