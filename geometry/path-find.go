package geometry

import (
	"container/heap"
	"math"
	"pure-game-kit/utility/collection"
	"pure-game-kit/utility/number"
)

func (s *ShapeGrid) FindPath(startX, startY, targetX, targetY float32, minimizePoints bool) [][2]float32 {
	return s.findPath(startX, startY, targetX, targetY, minimizePoints, false)
}
func (s *ShapeGrid) FindPathDiagonally(startX, startY, targetX, targetY float32, minimizePoints bool) [][2]float32 {
	return s.findPath(startX, startY, targetX, targetY, minimizePoints, true)
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
var directionsDiag = [][2]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}

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

func (s *ShapeGrid) findPath(stx, sty, tarx, tary float32, minPts, diag bool) [][2]float32 {
	var w, h = float32(s.cellWidth), float32(s.cellHeight)
	stx, sty, tarx, tary = stx/w, sty/h, tarx/w, tary/h
	var sx, sy = int(number.RoundDown(stx, 0)), int(number.RoundDown(sty, 0))
	var tx, ty = int(number.RoundDown(tarx, 0)), int(number.RoundDown(tary, 0))
	var open = &priorityQueue{}
	var startNode = &node{x: sx, y: sy, g: 0, h: heuristic(sx, sy, tx, ty)}
	var visited = map[[2]int]*node{}

	heap.Init(open)
	heap.Push(open, startNode)
	visited[[2]int{sx, sy}] = startNode

	var _, startBlocked = s.cells[[2]int{sx, sy}]
	var _, targetBlocked = s.cells[[2]int{tx, ty}]
	if startBlocked || targetBlocked {
		return [][2]float32{}
	}

	currentDirs := directions
	if diag {
		currentDirs = append(currentDirs, directionsDiag...)
	}

	for i := 0; open.Len() > 0 && i < 9999; i++ {
		var current = heap.Pop(open).(*node)
		if current.x == tx && current.y == ty {
			var result = make([][2]float32, 0)
			for cur := current; cur != nil; cur = cur.parent {
				result = append(result, [2]float32{float32(cur.x) + 0.5, float32(cur.y) + 0.5})
			}

			for i := range result {
				result[i][0] *= w
				result[i][1] *= h
			}

			for i := 0; i < len(result)/2; i++ {
				result[i], result[len(result)-1-i] = result[len(result)-1-i], result[i]
			}

			if minPts {
				result = removeRedundantPoints(result)
			}
			return result
		}

		for _, dir := range currentDirs {
			var nx, ny = current.x + dir[0], current.y + dir[1]
			var key = [2]int{nx, ny}
			var _, blocked = s.cells[key]
			if blocked {
				continue
			}

			moveCost := float32(1.0)
			if dir[0] != 0 && dir[1] != 0 {
				moveCost = 1.5
			}

			var newG = current.g + moveCost
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

	return [][2]float32{}
}
