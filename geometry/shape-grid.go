package geometry

import (
	"pure-game-kit/utility/number"
)

type ShapeGrid struct {
	cells                 map[[2]int][]*Shape
	cellWidth, cellHeight int
}

func NewShapeGrid(cellWidth, cellHeight int) *ShapeGrid {
	return &ShapeGrid{cellWidth: cellWidth, cellHeight: cellHeight, cells: make(map[[2]int][]*Shape)}
}

//=================================================================

func (shapeGrid *ShapeGrid) SetAtCell(x, y int, shapes ...*Shape) {
	var key = [2]int{x, y}
	var w, h = shapeGrid.cellWidth, shapeGrid.cellHeight
	shapeGrid.cells[key] = []*Shape{}

	for _, shape := range shapes {
		shape.gridX = float32(x*w) + (float32(w) * 0.5)
		shape.gridY = float32(y*h) + (float32(h) * 0.5)
	}
	shapeGrid.cells[key] = append(shapeGrid.cells[key], shapes...)
}

//=================================================================

func (shapeGrid *ShapeGrid) Cell(shape *Shape) (cellX, cellY int) {
	var w, h = float32(shapeGrid.cellWidth), float32(shapeGrid.cellHeight)
	var x, y = shape.gridX / w, shape.gridY / h
	return int(x), int(y)
}

func (shapeGrid *ShapeGrid) All() []*Shape {
	var result = []*Shape{}
	for k := range shapeGrid.cells {
		result = append(result, shapeGrid.AtCell(k[0], k[1])...)
	}
	return result
}
func (shapeGrid *ShapeGrid) AtCell(x, y int) []*Shape {
	var shapes, has = shapeGrid.cells[[2]int{x, y}]
	if has {
		return shapes
	}
	return []*Shape{}
}
func (shapeGrid *ShapeGrid) AtPoint(x, y float32) []*Shape {
	var w, h = float32(shapeGrid.cellWidth), float32(shapeGrid.cellHeight)
	if w == 0 || h == 0 {
		return []*Shape{}
	}
	var i, j = number.RoundDown(x / w), number.RoundDown(y / h)
	return shapeGrid.AtCell(int(i), int(j))
}
func (shapeGrid *ShapeGrid) AroundLine(line Line) []*Shape {
	var w, h = float32(shapeGrid.cellWidth), float32(shapeGrid.cellHeight)
	if w == 0 || h == 0 {
		return []*Shape{}
	}

	var result []*Shape
	var x0, y0, x1, y1 = line.Ax / w, line.Ay / h, line.Bx / w, line.By / h
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
		result = append(result, shapeGrid.AtCell(ix0, iy0)...)
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
func (shapeGrid *ShapeGrid) AroundShape(shape *Shape) []*Shape {
	var w, h = float32(shapeGrid.cellWidth), float32(shapeGrid.cellHeight)
	if w == 0 || h == 0 {
		return []*Shape{}
	}

	var corners = shape.CornerPoints()
	var result = []*Shape{}

	for i := 1; i < len(corners); i++ {
		var line = NewLine(corners[i-1][0], corners[i-1][1], corners[i][0], corners[i][1])
		result = append(result, shapeGrid.AroundLine(line)...)
	}
	return result
}
