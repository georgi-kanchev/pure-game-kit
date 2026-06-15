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

	var ctx = routeCtxPool.Get().(*routeContext)
	defer routeCtxPool.Put(ctx)
	ctx.reset()

	createNodes(ctx, paths)
	var originalNodeCount = len(ctx.nodes)

	// restrict segment searches to the original geometry to prevent ghost segments
	var sx, sy, startNodeIdx, sP1, sP2 = closestPointOnPath(ctx, startX, startY, originalNodeCount)
	var tx, ty, targetNodeIdx, tP1, tP2 = closestPointOnPath(ctx, targetX, targetY, originalNodeCount)

	*result = append(*result, startX, startY, sx, sy)

	if startNodeIdx == -1 || targetNodeIdx == -1 {
		return // disconnected or empty paths
	}

	if startNodeIdx >= originalNodeCount && targetNodeIdx >= originalNodeCount {
		if sP1 != -1 && sP1 == tP1 && sP2 == tP2 { // if both points projected as new nodes onto the exact same original edge,
			ctx.addTwoWayEdge(startNodeIdx, targetNodeIdx) // they are disconnected from each other - wire them directly to prevent U-turn
		}
	}

	var pathFound = walk(ctx, startNodeIdx, targetNodeIdx, result)
	if !pathFound {
		return // no valid path found
	}

	*result = append(*result, tx, ty, targetX, targetY)
	removeSequentialDuplicates(result) // strip sequentially stacked duplicates before running angle calculations
	remove180TurnsInPlace(result)
}

// private ========================================================

type routeNode struct {
	X, Y, SearchDist                 float32
	PathIndex, FirstEdge, SearchPrev int
	SearchSeen                       bool
}
type routeEdge struct{ To, Next int }
type routeContext struct {
	nodes []routeNode
	edges []routeEdge
	open  []int
}

var routeCtxPool = sync.Pool{
	New: func() any {
		return &routeContext{nodes: make([]routeNode, 0, 512), edges: make([]routeEdge, 0, 1024), open: make([]int, 0, 512)}
	},
}

func (ctx *routeContext) reset() {
	ctx.nodes = ctx.nodes[:0]
	ctx.edges = ctx.edges[:0]
	ctx.open = ctx.open[:0]
}
func (ctx *routeContext) addTwoWayEdge(a, b int) {
	var e1 = len(ctx.edges)
	ctx.edges = append(ctx.edges, routeEdge{To: b, Next: ctx.nodes[a].FirstEdge})
	ctx.nodes[a].FirstEdge = e1
	var e2 = len(ctx.edges)
	ctx.edges = append(ctx.edges, routeEdge{To: a, Next: ctx.nodes[b].FirstEdge})
	ctx.nodes[b].FirstEdge = e2
}

func pushOpen(ctx *routeContext, idx int) {
	ctx.open = append(ctx.open, idx)
	up2(ctx, len(ctx.open)-1)
}
func popOpen(ctx *routeContext) int {
	var n = len(ctx.open) - 1
	ctx.open[0], ctx.open[n] = ctx.open[n], ctx.open[0]
	down2(ctx, 0, n)
	var idx = ctx.open[n]
	ctx.open = ctx.open[:n]
	return idx
}
func up2(ctx *routeContext, j int) {
	for {
		var i = (j - 1) / 2
		if i == j || !less2(ctx, j, i) {
			break
		}
		ctx.open[i], ctx.open[j] = ctx.open[j], ctx.open[i]
		j = i
	}
}
func down2(ctx *routeContext, i0, n int) {
	var i = i0
	for {
		var j1 = 2*i + 1
		if j1 >= n || j1 < 0 {
			break
		}
		var j = j1
		var j2 = j1 + 1
		if j2 < n && less2(ctx, j2, j1) {
			j = j2
		}
		if !less2(ctx, j, i) {
			break
		}
		ctx.open[i], ctx.open[j] = ctx.open[j], ctx.open[i]
		i = j
	}
}
func less2(ctx *routeContext, i, j int) bool {
	return ctx.nodes[ctx.open[i]].SearchDist < ctx.nodes[ctx.open[j]].SearchDist
}

// ----------------------------------------------------------------------

func createNodes(ctx *routeContext, points []float32) {
	var pathIndex, prevIdx = 0, -1
	for i := 0; i < len(points); i += 2 {
		var px, py = points[i], points[i+1]
		if number.IsNaN(px) || number.IsNaN(py) {
			pathIndex, prevIdx = pathIndex+1, -1
			continue
		}

		var nIdx = len(ctx.nodes)
		ctx.nodes = append(ctx.nodes, routeNode{X: px, Y: py, PathIndex: pathIndex, FirstEdge: -1})
		if prevIdx != -1 && ctx.nodes[prevIdx].PathIndex == pathIndex {
			ctx.addTwoWayEdge(prevIdx, nIdx)
		}
		prevIdx = nIdx
	}

	for i := 0; i < len(ctx.nodes); i++ {
		for j := i + 1; j < len(ctx.nodes); j++ {
			if ctx.nodes[i].X == ctx.nodes[j].X && ctx.nodes[i].Y == ctx.nodes[j].Y {
				ctx.addTwoWayEdge(i, j)
			}
		}
	}
}
func closestPointOnPath(ctx *routeContext, x, y float32, maxNodes int) (closestX, closestY float32, bestIdx, bestP1, bestP2 int) {
	var bestDist, bP1, bP2 = number.Infinity(), -1, -1
	bestIdx = -1

	for i := 1; i < maxNodes; i++ {
		var p1, p2 = i - 1, i
		if ctx.nodes[p1].PathIndex != ctx.nodes[p2].PathIndex {
			continue
		}

		var n1, n2 = &ctx.nodes[p1], &ctx.nodes[p2]
		var curX, curY = NewLine(n1.X, n1.Y, n2.X, n2.Y, 0.1).ClosestPointToEdge(x, y)
		var dist = point.DistanceToPoint(x, y, curX, curY)
		if dist < bestDist {
			bestDist, closestX, closestY, bP1, bP2 = dist, curX, curY, p1, p2
		}
	}

	if bP1 != -1 && bP2 != -1 {
		var n1, n2 = &ctx.nodes[bP1], &ctx.nodes[bP2]
		if closestX == n1.X && closestY == n1.Y {
			return closestX, closestY, bP1, bP1, bP2
		}
		if closestX == n2.X && closestY == n2.Y {
			return closestX, closestY, bP2, bP1, bP2
		}

		var projIdx = len(ctx.nodes)
		ctx.nodes = append(ctx.nodes, routeNode{X: closestX, Y: closestY, PathIndex: n1.PathIndex, FirstEdge: -1})
		ctx.addTwoWayEdge(bP1, projIdx)
		ctx.addTwoWayEdge(bP2, projIdx)
		return closestX, closestY, projIdx, bP1, bP2
	}
	return closestX, closestY, -1, -1, -1
}
func walk(ctx *routeContext, startIdx, endIdx int, result *[]float32) bool {
	for i := range ctx.nodes {
		ctx.nodes[i].SearchDist, ctx.nodes[i].SearchPrev, ctx.nodes[i].SearchSeen = number.Infinity(), -1, false
	}

	ctx.nodes[startIdx].SearchDist = 0
	pushOpen(ctx, startIdx)

	for len(ctx.open) > 0 {
		var curIdx = popOpen(ctx)
		var cur = &ctx.nodes[curIdx]
		if cur.SearchSeen {
			continue
		}

		cur.SearchSeen = true
		if curIdx == endIdx {
			reconstructPath(ctx, endIdx, result)
			return true
		}

		for eIdx := cur.FirstEdge; eIdx != -1; eIdx = ctx.edges[eIdx].Next {
			var nbIdx = ctx.edges[eIdx].To
			var nb = &ctx.nodes[nbIdx]
			if nb.SearchSeen {
				continue
			}

			var dist = cur.SearchDist + point.DistanceToPoint(cur.X, cur.Y, nb.X, nb.Y)
			if dist < nb.SearchDist {
				nb.SearchDist, nb.SearchPrev = dist, curIdx
				pushOpen(ctx, nbIdx)
			}
		}
	}
	return false
}
func reconstructPath(ctx *routeContext, endIdx int, result *[]float32) {
	var startWrite, curIdx = len(*result), endIdx
	for curIdx != -1 { // blindly append all coords; duplicates are removed cleanly in a later pass
		var n = &ctx.nodes[curIdx]
		*result = append(*result, n.X, n.Y)
		curIdx = n.SearchPrev
	}

	var seg = (*result)[startWrite:] // reverse the newly appended segment in-place
	for i := 0; i < len(seg)/4; i++ {
		idx1, idx2 := i*2, len(seg)-2-(i*2)
		seg[idx1], seg[idx2] = seg[idx2], seg[idx1]
		seg[idx1+1], seg[idx2+1] = seg[idx2+1], seg[idx1+1]
	}
}
func removeSequentialDuplicates(path *[]float32) {
	var p = *path
	if len(p) < 4 {
		return
	}

	var writeIdx = 2
	for i := 2; i < len(p); i += 2 {
		if p[i] == p[writeIdx-2] && p[i+1] == p[writeIdx-1] {
			continue // Skip matching duplicates
		}
		p[writeIdx], p[writeIdx+1], writeIdx = p[i], p[i+1], writeIdx+2
	}
	*path = p[:writeIdx]
}
func remove180TurnsInPlace(path *[]float32) {
	var p = *path
	var nCount = len(p)
	if nCount < 6 {
		return
	}

	var writeIdx = 2
	for i := 2; i < nCount-2; i += 2 {
		var ax, ay, bx, by, cx, cy = p[writeIdx-2], p[writeIdx-1], p[i], p[i+1], p[i+2], p[i+3]
		var ang1, ang2 = angle.BetweenPoints(ax, ay, bx, by), angle.BetweenPoints(bx, by, cx, cy)
		var diff = number.Wrap(ang2-ang1, 0, 360)
		if !number.IsWithin(diff, 180, 0.1) {
			p[writeIdx], p[writeIdx+1] = bx, by
			writeIdx += 2
		}
	}

	p[writeIdx], p[writeIdx+1], writeIdx = p[nCount-2], p[nCount-1], writeIdx+2
	*path = p[:writeIdx]
}
