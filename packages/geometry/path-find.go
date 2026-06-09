package geometry

import (
	"container/heap"
	"math"
	"pure-game-kit/packages/utility/collection"
	"pure-game-kit/packages/utility/number"
)

func (s *ShapeGrid) FindPath(startX, startY, targetX, targetY float32, minimizePoints bool) []float32 {
	return s.findPath(startX, startY, targetX, targetY, minimizePoints, false)
}
func (s *ShapeGrid) FindPathDiagonally(startX, startY, targetX, targetY float32, minimizePoints bool) []float32 {
	return s.findPath(startX, startY, targetX, targetY, minimizePoints, true)
}

// private ========================================================

type priorityQueue []*node

type node struct {
	x, y   int
	g, h   float32
	parent *node
}

var directions = [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
var directionsDiag = [][2]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}

func (pq priorityQueue) Len() int           { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool { return pq[i].f() < pq[j].f() }
func (pq priorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq *priorityQueue) Push(x any)        { *pq = append(*pq, x.(*node)) }
func (pq *priorityQueue) Pop() any {
	var old = *pq
	var item = old[len(old)-1]
	*pq = old[:len(old)-1]
	return item
}
func (n *node) f() float32 { return n.g + n.h }

func heuristic(ax, ay, bx, by int) float32 {
	var dx, dy = float32(ax - bx), float32(ay - by)
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}
func removeRedundantPoints(points []float32) []float32 {
	var pts = collection.Clone(points)
	if len(pts) < 6 {
		return pts
	}

	for i := 2; i < len(pts)-2; i += 2 {
		var ax, ay, bx, by, cx, cy = pts[i-2], pts[i-1], pts[i], pts[i+1], pts[i+2], pts[i+3]
		var cross = (bx-ax)*(cy-ay) - (by-ay)*(cx-ax)
		if cross > -0.001 && cross < 0.001 {
			pts = append(pts[:i], pts[i+2:]...)
			i -= 2
		}
	}
	return pts
}
func (s *ShapeGrid) findPath(stx, sty, tarx, tary float32, minPts, diag bool) []float32 {
	if s == nil {
		return nil
	}

	var w, h = float32(s.chunkSize), float32(s.chunkSize)
	stx, sty, tarx, tary = stx/w, sty/h, tarx/w, tary/h
	var sx, sy = int(number.RoundDown(stx, 0)), int(number.RoundDown(sty, 0))
	var tx, ty = int(number.RoundDown(tarx, 0)), int(number.RoundDown(tary, 0))
	var startKey = [2]int{sx, sy}
	var targetKey = [2]int{tx, ty}

	var _, startBlocked = s.cells[startKey]
	var _, targetBlocked = s.cells[targetKey]
	if startBlocked || targetBlocked {
		return []float32{}
	}

	var open = &priorityQueue{}
	var bestG = map[[2]int]float32{startKey: 0}
	heap.Init(open)
	heap.Push(open, &node{x: sx, y: sy, g: 0, h: heuristic(sx, sy, tx, ty)})

	var currentDirs = directions
	if diag {
		currentDirs = append(currentDirs, directionsDiag...)
	}

	for i := 0; open.Len() > 0 && i < 9999; i++ {
		var current = heap.Pop(open).(*node)

		if current.g > bestG[[2]int{current.x, current.y}] {
			continue // Skip stale entries (better g was already found for this cell)
		}

		if current.x == tx && current.y == ty {
			var result = make([]float32, 0)
			for cur := current; cur != nil; cur = cur.parent {
				result = append(result, (float32(cur.x)+0.5)*w, (float32(cur.y)+0.5)*h)
			}

			for i := 0; i < len(result)/4; i++ {
				var idx1, idx2 = i * 2, len(result) - 2 - (i * 2)
				result[idx1], result[idx2], result[idx1+1], result[idx2+1] = result[idx2], result[idx1], result[idx2+1], result[idx1+1]
			}

			if minPts {
				result = removeRedundantPoints(result)
			}
			return result
		}

		for _, dir := range currentDirs {
			var nx, ny = current.x + dir[0], current.y + dir[1]
			var key = [2]int{nx, ny}
			if _, blocked := s.cells[key]; blocked {
				continue
			}

			moveCost := float32(1.0)
			if dir[0] != 0 && dir[1] != 0 {
				moveCost = 1.5
			}

			var newG = current.g + moveCost
			if prevG, seen := bestG[key]; seen && newG >= prevG {
				continue
			}

			bestG[key] = newG
			heap.Push(open, &node{x: nx, y: ny, g: newG, h: heuristic(nx, ny, tx, ty), parent: current})
		}
	}
	return []float32{}
}
