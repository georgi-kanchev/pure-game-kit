package geometry

import (
	"pure-game-kit/packages/utility/angle"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/point"
	"sync"
)

// Calculates the minimal subsections of the routes to traverse to reach the target point.
// Start point and target point can be anywhere (not necessarily on the paths).
// Pass a pre-allocated result slice to achieve zero allocations.
func FollowPaths(startX, startY, targetX, targetY float32, paths []float32, result *[]float32) {
	*result = (*result)[:0]

	ctx := routeCtxPool.Get().(*RouteContext)
	defer routeCtxPool.Put(ctx)
	ctx.reset()

	createNodes(ctx, paths)

	sx, sy, startNodeIdx := closestPointOnPath(ctx, startX, startY)
	tx, ty, targetNodeIdx := closestPointOnPath(ctx, targetX, targetY)

	*result = append(*result, startX, startY, sx, sy)

	if startNodeIdx == -1 || targetNodeIdx == -1 {
		return // Disconnected or empty paths
	}

	if pathFound := walk(ctx, startNodeIdx, targetNodeIdx, result); !pathFound {
		return // No valid path found
	}

	*result = append(*result, tx, ty, targetX, targetY)
	remove180TurnsInPlace(result)
}

// private ========================================================

type routeNode struct {
	X, Y       float32
	PathIndex  int
	FirstEdge  int // Index into RouteContext.edges, -1 if none
	SearchDist float32
	SearchPrev int
	SearchSeen bool
}

type routeEdge struct {
	To   int
	Next int // Next edge for the source node, -1 if end of list
}

// RouteContext holds all reusable memory buffers
type RouteContext struct {
	nodes []routeNode
	edges []routeEdge
	open  []int // Custom min-heap
}

func (ctx *RouteContext) reset() {
	ctx.nodes = ctx.nodes[:0]
	ctx.edges = ctx.edges[:0]
	ctx.open = ctx.open[:0]
}

func (ctx *RouteContext) addTwoWayEdge(a, b int) {
	e1 := len(ctx.edges)
	ctx.edges = append(ctx.edges, routeEdge{To: b, Next: ctx.nodes[a].FirstEdge})
	ctx.nodes[a].FirstEdge = e1

	e2 := len(ctx.edges)
	ctx.edges = append(ctx.edges, routeEdge{To: a, Next: ctx.nodes[b].FirstEdge})
	ctx.nodes[b].FirstEdge = e2
}

var routeCtxPool = sync.Pool{
	New: func() any {
		return &RouteContext{
			nodes: make([]routeNode, 0, 512),
			edges: make([]routeEdge, 0, 1024),
			open:  make([]int, 0, 512),
		}
	},
}

// --- Custom Typed Min-Heap ---

func pushOpen(ctx *RouteContext, idx int) {
	ctx.open = append(ctx.open, idx)
	up2(ctx, len(ctx.open)-1)
}

func popOpen(ctx *RouteContext) int {
	n := len(ctx.open) - 1
	ctx.open[0], ctx.open[n] = ctx.open[n], ctx.open[0]
	down2(ctx, 0, n)
	idx := ctx.open[n]
	ctx.open = ctx.open[:n]
	return idx
}

func up2(ctx *RouteContext, j int) {
	for {
		i := (j - 1) / 2
		if i == j || !less2(ctx, j, i) {
			break
		}
		ctx.open[i], ctx.open[j] = ctx.open[j], ctx.open[i]
		j = i
	}
}

func down2(ctx *RouteContext, i0, n int) {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 {
			break
		}
		j := j1
		if j2 := j1 + 1; j2 < n && less2(ctx, j2, j1) {
			j = j2
		}
		if !less2(ctx, j, i) {
			break
		}
		ctx.open[i], ctx.open[j] = ctx.open[j], ctx.open[i]
		i = j
	}
}

func less2(ctx *RouteContext, i, j int) bool {
	return ctx.nodes[ctx.open[i]].SearchDist < ctx.nodes[ctx.open[j]].SearchDist
}

// ----------------------------------------------------------------------

func createNodes(ctx *RouteContext, points []float32) {
	pathIndex := 0
	prevIdx := -1

	for i := 0; i < len(points); i += 2 {
		px, py := points[i], points[i+1]
		if number.IsNaN(px) || number.IsNaN(py) {
			pathIndex++
			prevIdx = -1
			continue
		}

		nIdx := len(ctx.nodes)
		ctx.nodes = append(ctx.nodes, routeNode{
			X: px, Y: py, PathIndex: pathIndex, FirstEdge: -1,
		})

		if prevIdx != -1 && ctx.nodes[prevIdx].PathIndex == pathIndex {
			ctx.addTwoWayEdge(prevIdx, nIdx)
		}
		prevIdx = nIdx
	}

	// Link coincident nodes across different paths
	for i := 0; i < len(ctx.nodes); i++ {
		for j := i + 1; j < len(ctx.nodes); j++ {
			if ctx.nodes[i].X == ctx.nodes[j].X && ctx.nodes[i].Y == ctx.nodes[j].Y {
				ctx.addTwoWayEdge(i, j)
			}
		}
	}
}

func closestPointOnPath(ctx *RouteContext, x, y float32) (closestX, closestY float32, bestIdx int) {
	bestDist := number.Infinity()
	bestP1, bestP2 := -1, -1
	bestIdx = -1

	for i := 1; i < len(ctx.nodes); i++ {
		p1, p2 := i-1, i
		if ctx.nodes[p1].PathIndex != ctx.nodes[p2].PathIndex {
			continue // ignore path connections
		}

		n1, n2 := &ctx.nodes[p1], &ctx.nodes[p2]

		// Note: Assuming NewLine returns a value struct and doesn't allocate.
		curX, curY := NewLine(n1.X, n1.Y, n2.X, n2.Y, 0.1).ClosestPointToEdge(x, y)
		dist := point.DistanceToPoint(x, y, curX, curY)

		if dist < bestDist {
			bestDist, closestX, closestY, bestP1, bestP2 = dist, curX, curY, p1, p2
		}
	}

	if bestP1 != -1 && bestP2 != -1 {
		n1, n2 := &ctx.nodes[bestP1], &ctx.nodes[bestP2]
		if closestX == n1.X && closestY == n1.Y {
			return closestX, closestY, bestP1
		}
		if closestX == n2.X && closestY == n2.Y {
			return closestX, closestY, bestP2
		}

		// Insert new projection node
		projIdx := len(ctx.nodes)
		ctx.nodes = append(ctx.nodes, routeNode{
			X: closestX, Y: closestY, PathIndex: n1.PathIndex, FirstEdge: -1,
		})
		ctx.addTwoWayEdge(bestP1, projIdx)
		ctx.addTwoWayEdge(bestP2, projIdx)
		return closestX, closestY, projIdx
	}

	return closestX, closestY, -1
}

func walk(ctx *RouteContext, startIdx, endIdx int, result *[]float32) bool {
	// Initialize traversal states
	for i := range ctx.nodes {
		ctx.nodes[i].SearchDist = number.Infinity()
		ctx.nodes[i].SearchPrev = -1
		ctx.nodes[i].SearchSeen = false
	}

	ctx.nodes[startIdx].SearchDist = 0
	pushOpen(ctx, startIdx)

	for len(ctx.open) > 0 {
		curIdx := popOpen(ctx)
		cur := &ctx.nodes[curIdx]

		if cur.SearchSeen {
			continue
		}
		cur.SearchSeen = true

		if curIdx == endIdx {
			reconstructPath(ctx, endIdx, result)
			return true
		}

		for eIdx := cur.FirstEdge; eIdx != -1; eIdx = ctx.edges[eIdx].Next {
			nbIdx := ctx.edges[eIdx].To
			nb := &ctx.nodes[nbIdx]

			if nb.SearchSeen {
				continue
			}

			dist := cur.SearchDist + point.DistanceToPoint(cur.X, cur.Y, nb.X, nb.Y)
			if dist < nb.SearchDist {
				nb.SearchDist = dist
				nb.SearchPrev = curIdx
				pushOpen(ctx, nbIdx)
			}
		}
	}
	return false
}

func reconstructPath(ctx *RouteContext, endIdx int, result *[]float32) {
	startWrite := len(*result)
	curIdx := endIdx

	for curIdx != -1 {
		n := &ctx.nodes[curIdx]

		// O(N) Uniqueness check instead of allocating map[[2]float32]any
		isDup := false
		for i := startWrite; i < len(*result); i += 2 {
			if (*result)[i] == n.X && (*result)[i+1] == n.Y {
				isDup = true
				break
			}
		}

		if !isDup {
			*result = append(*result, n.X, n.Y)
		}
		curIdx = n.SearchPrev
	}

	// Reverse the newly appended segment in-place
	seg := (*result)[startWrite:]
	for i := 0; i < len(seg)/4; i++ {
		idx1, idx2 := i*2, len(seg)-2-(i*2)
		seg[idx1], seg[idx2] = seg[idx2], seg[idx1]
		seg[idx1+1], seg[idx2+1] = seg[idx2+1], seg[idx1+1]
	}
}

func remove180TurnsInPlace(path *[]float32) {
	p := *path
	nCount := len(p)
	if nCount < 6 {
		return
	}

	writeIdx := 2
	for i := 2; i < nCount-2; i += 2 {
		ax, ay := p[writeIdx-2], p[writeIdx-1]
		bx, by := p[i], p[i+1]
		cx, cy := p[i+2], p[i+3]

		ang1, ang2 := angle.BetweenPoints(ax, ay, bx, by), angle.BetweenPoints(bx, by, cx, cy)
		diff := number.Wrap(ang2-ang1, 0, 360)

		if !number.IsWithin(diff, 180, 0.1) {
			p[writeIdx], p[writeIdx+1] = bx, by
			writeIdx += 2
		}
	}

	p[writeIdx], p[writeIdx+1] = p[nCount-2], p[nCount-1]
	writeIdx += 2

	*path = p[:writeIdx]
}
