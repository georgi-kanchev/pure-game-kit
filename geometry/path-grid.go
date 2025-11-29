package geometry

import (
	"container/heap"
	"math"
	"pure-game-kit/utility/collection"
	"pure-game-kit/utility/number"
)

func (shapeGrid *ShapeGrid) FindPath(start, target [2]float32, turnFactor int, minimizePoints bool) [][2]float32 {
	var w, h = float32(shapeGrid.cellWidth), float32(shapeGrid.cellHeight)
	start[0], start[1], target[0], target[1] = start[0]/w, start[1]/h, target[0]/w, target[1]/h
	var sx, sy = int(number.RoundDown(start[0], 0)), int(number.RoundDown(start[1], 0))
	var tx, ty = int(number.RoundDown(target[0], 0)), int(number.RoundDown(target[1], 0))
	var open = &priorityQueue{}
	var startNode = &node{x: sx, y: sy, g: 0, h: heuristic(sx, sy, tx, ty)}
	var visited = map[[2]int]*node{}

	heap.Init(open)
	heap.Push(open, startNode)
	visited[[2]int{sx, sy}] = startNode

	var _, startBlocked = shapeGrid.cells[[2]int{sx, sy}]
	var _, targetBlocked = shapeGrid.cells[[2]int{tx, ty}]
	if startBlocked || targetBlocked {
		return [][2]float32{}
	}

	for i := 0; open.Len() > 0 && i < 9999; i++ {
		var current = heap.Pop(open).(*node)
		if current.x == tx && current.y == ty {
			var result = make([][2]float32, 0)
			for cur := current; cur != nil; cur = cur.parent {
				result = append(result, [2]float32{float32(cur.x) + 0.5, float32(cur.y) + 0.5})
			} // offset to the center of the cell by adding 0.5

			for i := range result { // convert to world coordinates
				result[i][0] *= w
				result[i][1] *= h
			}

			for i := 0; i < len(result)/2; i++ { // reverse path
				result[i], result[len(result)-1-i] = result[len(result)-1-i], result[i]
			}

			result = shapeGrid.smoothZigzag(result, turnFactor, minimizePoints)
			result = removeRedundantPoints(result)
			return result
		}

		for _, dir := range directions {
			var nx, ny = current.x + dir[0], current.y + dir[1]
			var key = [2]int{nx, ny}
			var _, blocked = shapeGrid.cells[key]
			if blocked {
				continue // unwalkable if present in map
			}

			var newG = current.g + 1
			var old, seen = visited[key]

			if !seen {
				var n = &node{x: nx, y: ny, g: newG, h: heuristic(nx, ny, tx, ty), parent: current}
				visited[key] = n
				heap.Push(open, n)
				continue
			}

			if newG < old.g {
				old.g = newG
				old.parent = current
				heap.Fix(open, old.index)
			}
		}
	}

	return [][2]float32{} // no path found
}

//=================================================================
// private

type priorityQueue []*node

type node struct {
	x, y   int
	g, h   float32
	parent *node
	index  int // for heap
}

var directions = [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

func (pq priorityQueue) Len() int {
	return len(pq)
}
func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].f() < pq[j].f()
}
func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}
func (pq *priorityQueue) Push(x any) {
	var n = x.(*node)
	n.index = len(*pq)
	*pq = append(*pq, n)
}
func (pq *priorityQueue) Pop() any {
	var old = *pq
	var n = len(old)
	var item = old[n-1]
	item.index = -1
	*pq = old[:n-1]
	return item
}

func (n *node) f() float32 {
	return n.g + n.h
}

func heuristic(ax, ay, bx, by int) float32 {
	var dx, dy = float32(ax - bx), float32(ay - by)
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}

func (shapeGrid *ShapeGrid) smoothZigzag(points [][2]float32, turnFactor int, minimizePoints bool) [][2]float32 {
	var pts = collection.Clone(points)
	if len(pts) < 2 {
		return pts
	}

	var w, h = shapeGrid.cellWidth, shapeGrid.cellHeight
	for i := 1; i < len(pts)-1; i++ {
		var horDiff = number.Absolute(pts[i-1][0] - pts[i+1][0])
		var verDiff = number.Absolute(pts[i-1][1] - pts[i+1][1])

		if horDiff > float32(turnFactor*w) || verDiff > float32(turnFactor*h) {
			continue
		}

		if minimizePoints {
			pts = append(pts[:i], pts[i+1:]...)
			i--
			continue
		}
		// Smooth the corner by averaging with neighbors
		pts[i][0] = (pts[i-1][0] + pts[i][0] + pts[i+1][0]) / 3
		pts[i][1] = (pts[i-1][1] + pts[i][1] + pts[i+1][1]) / 3
	}
	return pts
}
func removeRedundantPoints(points [][2]float32) [][2]float32 {
	var pts = collection.Clone(points)
	if len(pts) < 3 {
		return pts
	}

	for i := 1; i < len(pts)-1; i++ {
		var a, b, c = pts[i-1], pts[i], pts[i+1]
		var cross = (b[0]-a[0])*(c[1]-a[1]) - (b[1]-a[1])*(c[0]-a[0])

		if cross > -0.001 && cross < 0.001 {
			pts = append(pts[:i], pts[i+1:]...)
			i--
		}
	}

	return pts
}
