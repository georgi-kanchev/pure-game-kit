package geometry

import (
	"math"
)

type ShapeGrid struct {
	cells                 map[[2]int][]Shape
	cellWidth, cellHeight int
}

func NewShapeGrid(cellWidth, cellHeight int) ShapeGrid {
	return ShapeGrid{cellWidth: cellWidth, cellHeight: cellHeight, cells: make(map[[2]int][]Shape)}
}

func (shapeGrid *ShapeGrid) SetShapesAtCell(x, y int, shapes ...Shape) {
	var key = [2]int{x, y}
	shapeGrid.cells[key] = []Shape{}
	shapeGrid.cells[key] = append(shapeGrid.cells[key], shapes...)
}

func (shapeGrid *ShapeGrid) AllShapes() []Shape {
	var result = []Shape{}
	for k := range shapeGrid.cells {
		result = append(result, shapeGrid.ShapesAtCell(k[0], k[1])...)
	}
	return result
}
func (shapeGrid *ShapeGrid) ShapesAtCell(x, y int) []Shape {
	// this makes a copy on purpose, the original shape values shouldn't change
	// also the whole result slice is a copy so it cannot be extended by the user without calling SetShapesAtCell()

	var shapes, has = shapeGrid.cells[[2]int{x, y}]
	var w, h = shapeGrid.cellWidth, shapeGrid.cellHeight
	if has {
		var result = make([]Shape, len(shapes))
		for i := range shapes {
			var shape = shapes[i]
			shape.X += float32(x*w) + (float32(w) * 0.5)
			shape.Y += float32(y*h) + (float32(h) * 0.5)
			result[i] = shape
		}

		return result
	}
	return []Shape{}
}
func (shapeGrid *ShapeGrid) ShapesAtPoint(x, y float32) []Shape {
	var w, h = float32(shapeGrid.cellWidth), float32(shapeGrid.cellHeight)
	if w == 0 || h == 0 {
		return []Shape{}
	}
	var i, j = int(math.Floor(float64(x / w))), int(math.Floor(float64(y / h)))
	return shapeGrid.ShapesAtCell(i, j)
}
func (shapeGrid *ShapeGrid) ShapesAroundLine(line Line) []Shape {
	var w, h = float32(shapeGrid.cellWidth), float32(shapeGrid.cellHeight)
	if w == 0 || h == 0 {
		return []Shape{}
	}

	var result []Shape
	var x0, y0, x1, y1 = line.Ax / w, line.Ay / h, line.Bx / w, line.By / h
	var ix0, iy0 = int(math.Floor(float64(x0))), int(math.Floor(float64(y0)))
	var ix1, iy1 = int(math.Floor(float64(x1))), int(math.Floor(float64(y1)))
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
		tMaxX = math.MaxFloat32
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
		tMaxY = math.MaxFloat32
	}

	for { // Traverse until reaching the target cell
		result = append(result, shapeGrid.ShapesAtCell(ix0, iy0)...)
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
func (shapeGrid *ShapeGrid) ShapesAroundShape(shape *Shape) []Shape {
	var w, h = float32(shapeGrid.cellWidth), float32(shapeGrid.cellHeight)
	if w == 0 || h == 0 {
		return []Shape{}
	}

	var corners = shape.CornerPoints()
	var result = []Shape{}

	for i := 1; i < len(corners); i++ {
		var line = NewLine(corners[i-1][0], corners[i-1][1], corners[i][0], corners[i][1])
		result = append(result, shapeGrid.ShapesAroundLine(line)...)
	}
	return result
}

func (shapeGrid *ShapeGrid) CrossPointsWithLine(line Line) [][2]float32 {
	var shapes = shapeGrid.ShapesAroundLine(line)
	var result = [][2]float32{}
	for _, s := range shapes {
		result = append(result, s.internalCrossPointsWithLine(s.CornerPoints(), line)...)
	}

	return result
}
func (shapeGrid *ShapeGrid) CrossPointsWithShape(target *Shape) [][2]float32 {
	var shapes = shapeGrid.ShapesAroundShape(target)
	var targetCorners = target.CornerPoints()
	var result = [][2]float32{}
	for _, s := range shapes {
		result = append(result, s.internalCrossPointsWithShape(s.CornerPoints(), targetCorners)...)
	}

	return result
}

func (shapeGrid *ShapeGrid) IsCrossingLine(line Line) bool {
	var shapes = shapeGrid.ShapesAroundLine(line)
	for _, s := range shapes {
		if s.internalIsCrossingLine(s.CornerPoints(), line) {
			return true
		}
	}

	return false
}
func (shapeGrid *ShapeGrid) IsCrossingShape(target *Shape) bool {
	var shapes = shapeGrid.ShapesAroundShape(target)
	var targetCorners = target.CornerPoints()
	for _, s := range shapes {
		if s.internalIsCrossingShape(s.CornerPoints(), targetCorners) {
			return true
		}
	}

	return false
}
