package geometry

import (
	"pure-game-kit/packages/utility/number"
	"sync"
)

func (s *ShapeGrid) FindPath(startX, startY, targetX, targetY float32, minimizePoints bool, result *[]float32) {
	s.findPath(startX, startY, targetX, targetY, minimizePoints, false, result)
}
func (s *ShapeGrid) FindPathDiagonally(startX, startY, targetX, targetY float32, minimizePoints bool, result *[]float32) {
	s.findPath(startX, startY, targetX, targetY, minimizePoints, true, result)
}

// private ========================================================

var dirs4 = [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
var dirs8 = [8][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
var pathCtxPool = sync.Pool{
	New: func() any { // Pool to reuse pathing contexts. After warming up, this guarantees 0 allocations.
		return &PathContext{nodes: make([]node, 0, 1024), open: make([]int, 0, 1024), bestG: make(map[[2]int]float32, 1024)}
	},
}

type node struct {
	x, y   int
	g, h   float32
	parent int // Index referencing the pool, -1 for none
}
type PathContext struct { // holds all reusable memory buffers
	nodes []node
	open  []int // custom min-heap of indices into `nodes`
	bestG map[[2]int]float32
}

func (s *ShapeGrid) findPath(stx, sty, tarx, tary float32, minPts, diag bool, result *[]float32) {
	if s == nil {
		return
	}

	*result = (*result)[:0] // ensure the result slice starts empty but retains capacity

	var w, h = float32(s.chunkSize), float32(s.chunkSize)
	stx, sty, tarx, tary = stx/w, sty/h, tarx/w, tary/h
	var sx, sy = int(number.RoundDown(stx, 0)), int(number.RoundDown(sty, 0))
	var tx, ty = int(number.RoundDown(tarx, 0)), int(number.RoundDown(tary, 0))
	var startKey, targetKey = [2]int{sx, sy}, [2]int{tx, ty}
	var _, startBlocked = s.cells[startKey]
	if startBlocked {
		return
	}
	var _, targetBlocked = s.cells[targetKey]
	if targetBlocked {
		return
	}

	var ctx = pathCtxPool.Get().(*PathContext) // acquire pooled context for zero-allocation memory buffers
	defer pathCtxPool.Put(ctx)
	ctx.reset()

	ctx.bestG[startKey] = 0
	pushNode(ctx, node{x: sx, y: sy, g: 0, h: heuristic(sx, sy, tx, ty), parent: -1})

	var dirs = dirs4[:]
	if diag {
		dirs = dirs8[:]
	}

	for i := 0; len(ctx.open) > 0 && i < 9999; i++ {
		var currentIdx = popNode(ctx)
		var current = &ctx.nodes[currentIdx]
		if current.g > ctx.bestG[[2]int{current.x, current.y}] {
			continue
		}

		if current.x == tx && current.y == ty {
			for curIdx := currentIdx; curIdx != -1; curIdx = ctx.nodes[curIdx].parent { // backtrack path using indices
				n := &ctx.nodes[curIdx]
				*result = append(*result, (float32(n.x)+0.5)*w, (float32(n.y)+0.5)*h)
			}

			for j := 0; j < len(*result)/4; j++ { // reverse result slice in-place
				var idx1, idx2 = j * 2, len(*result) - 2 - (j * 2)
				(*result)[idx1], (*result)[idx2] = (*result)[idx2], (*result)[idx1]
				(*result)[idx1+1], (*result)[idx2+1] = (*result)[idx2+1], (*result)[idx1+1]
			}

			if minPts {
				removeRedundantPointsInPlace(result)
			}
			return
		}

		for _, dir := range dirs {
			var nx, ny = current.x + dir[0], current.y + dir[1]
			var key = [2]int{nx, ny}
			var _, blocked = s.cells[key]
			if blocked {
				continue
			}

			var moveCost = float32(1.0)
			if dir[0] != 0 && dir[1] != 0 {
				moveCost = 1.5
			}

			var newG = current.g + moveCost
			var prevG, seen = ctx.bestG[key]
			if seen && newG >= prevG {
				continue
			}

			ctx.bestG[key] = newG
			pushNode(ctx, node{x: nx, y: ny, g: newG, h: heuristic(nx, ny, tx, ty), parent: currentIdx})
		}
	}
}

func (ctx *PathContext) reset() {
	ctx.nodes, ctx.open = ctx.nodes[:0], ctx.open[:0]
	clear(ctx.bestG)
}

func pushNode(ctx *PathContext, n node) {
	ctx.nodes = append(ctx.nodes, n)
	var idx = len(ctx.nodes) - 1
	ctx.open = append(ctx.open, idx)
	up(ctx, len(ctx.open)-1)
}
func popNode(ctx *PathContext) int {
	var n = len(ctx.open) - 1
	ctx.open[0], ctx.open[n] = ctx.open[n], ctx.open[0]
	down(ctx, 0, n)
	var idx = ctx.open[n]
	ctx.open = ctx.open[:n]
	return idx
}
func up(ctx *PathContext, j int) {
	for {
		var i = (j - 1) / 2 // parent
		if i == j || !less(ctx, j, i) {
			break
		}
		ctx.open[i], ctx.open[j] = ctx.open[j], ctx.open[i]
		j = i
	}
}
func down(ctx *PathContext, i0, n int) {
	var i = i0
	for {
		var j1 = 2*i + 1
		if j1 >= n || j1 < 0 {
			break
		}
		var j, j2 = j1, j1 + 1
		if j2 < n && less(ctx, j2, j1) {
			j = j2
		}
		if !less(ctx, j, i) {
			break
		}
		ctx.open[i], ctx.open[j] = ctx.open[j], ctx.open[i]
		i = j
	}
}
func less(ctx *PathContext, i, j int) bool {
	var a, b = &ctx.nodes[ctx.open[i]], &ctx.nodes[ctx.open[j]]
	return (a.g + a.h) < (b.g + b.h)
}
func heuristic(ax, ay, bx, by int) float32 {
	var dx, dy = float32(ax - bx), float32(ay - by)
	return number.SquareRoot(dx*dx + dy*dy)
}
func removeRedundantPointsInPlace(pts *[]float32) { // in-place redundancy removal - modifies underlying array without cloning
	var p = *pts
	if len(p) < 6 {
		return
	}

	var writeIdx = 2
	for i := 2; i < len(p)-2; i += 2 {
		var ax, ay, bx, by, cx, cy = p[writeIdx-2], p[writeIdx-1], p[i], p[i+1], p[i+2], p[i+3]
		var cross = (bx-ax)*(cy-ay) - (by-ay)*(cx-ax)
		if cross > -0.001 && cross < 0.001 {
			continue // Collinear: skip `b` by not advancing writeIdx
		}

		p[writeIdx], p[writeIdx+1] = bx, by
		writeIdx += 2
	}

	p[writeIdx], p[writeIdx+1] = p[len(p)-2], p[len(p)-1]
	writeIdx += 2 // always keep the final target coordinate
	*pts = p[:writeIdx]
}
