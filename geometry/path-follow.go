package geometry

import (
	"container/heap"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/point"
)

// calculates the minimal subsections of the routes to traverse to reach target
//
// start and target point can be anywhere (not necessarily on the paths)
//
// multiple paths can be separated by [NaN, NaN]
func FollowPaths(startX, startY, targetX, targetY float32, paths ...float32) []float32 {
	var allNodes = createNodes(paths)
	var sx, sy, startNode = closestPointOnPath(startX, startY, allNodes)
	var tx, ty, targetNode = closestPointOnPath(targetX, targetY, allNodes)
	var path = startNode.walk(targetNode)
	var result = []float32{startX, startY, sx, sy}
	if len(path) == 0 { // no path was found, perhaps the end target is a disconnected one
		return result
	}

	result = append(result, path...)
	result = append(result, tx, ty, targetX, targetY)
	result = remove180Turns(result)
	return result
}

//=================================================================
// private

type n struct { // node
	X, Y      float32
	PathIndex int
	Neighbors []*n

	SearchPrev *n
	SearchDist float32
	SearchSeen bool
}
type npq []*n                     // node priority queue
func (pq npq) Len() int           { return len(pq) }
func (pq npq) Less(i, j int) bool { return pq[i].SearchDist < pq[j].SearchDist }
func (pq npq) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq *npq) Push(x any)        { *pq = append(*pq, x.(*n)) }
func (pq *npq) Pop() any {
	old := *pq
	n := old[len(old)-1]
	*pq = old[:len(old)-1]
	return n
}

func createNodes(points []float32) []*n {
	var result = []*n{}
	var pathIndex = 0
	for i := 0; i < len(points); i += 2 {
		var px, py = points[i], points[i+1]
		if number.IsNaN(px) || number.IsNaN(py) {
			pathIndex++
			continue
		}
		var prevNode *n
		if len(result) > 0 {
			prevNode = result[len(result)-1]
		}
		var newN = &n{X: px, Y: py, PathIndex: pathIndex, SearchDist: number.Infinity()}
		result = append(result, newN)

		if prevNode != nil && prevNode.PathIndex == newN.PathIndex {
			newN.Neighbors = append(newN.Neighbors, prevNode)
			prevNode.Neighbors = append(prevNode.Neighbors, newN)
		}
	}
	for _, node := range result {
		for _, target := range result {
			if node != target && node.X == target.X && node.Y == target.Y {
				node.Neighbors = append(node.Neighbors, target)
			}
		}
	}
	return result
}
func (start *n) walk(end *n) []float32 {
	var all []*n
	var stack []*n = []*n{start}
	for len(stack) > 0 {
		var cur = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if !cur.SearchSeen {
			cur.SearchSeen = true
			all = append(all, cur)
			stack = append(stack, cur.Neighbors...)
		}
	}
	for _, node := range all {
		node.SearchDist = number.Infinity()
		node.SearchPrev = nil
		node.SearchSeen = false
	}

	start.SearchDist = 0

	var prioQueue = &npq{}
	heap.Init(prioQueue)
	heap.Push(prioQueue, start)

	for prioQueue.Len() > 0 {
		var cur = heap.Pop(prioQueue).(*n)
		if cur.SearchSeen {
			continue
		}
		cur.SearchSeen = true

		if cur == end {
			return cur.reconstructPath()
		}

		for _, nb := range cur.Neighbors {
			if nb.SearchSeen {
				continue
			}
			var dist = cur.SearchDist + point.DistanceToPoint(cur.X, cur.Y, nb.X, nb.Y)
			if dist < nb.SearchDist {
				nb.SearchDist = dist
				nb.SearchPrev = cur
				heap.Push(prioQueue, nb)
			}
		}
	}

	return nil
}
func (n *n) reconstructPath() []float32 {
	var result []float32
	var unique = map[[2]float32]any{}
	for n != nil {
		var p = [2]float32{n.X, n.Y}
		var _, has = unique[p]
		if !has {
			result = append([]float32{p[0], p[1]}, result...)
		}
		unique[p] = nil
		n = n.SearchPrev
	}

	return result
}
func closestPointOnPath(x, y float32, nodes []*n) (closestX, closestY float32, newNode *n) {
	var bestDist = number.Infinity()
	var bestP1, bestP2 *n

	for i := 1; i < len(nodes); i++ {
		var p1, p2 = nodes[i-1], nodes[i]
		if p1.PathIndex != p2.PathIndex {
			continue // ignore path connections
		}

		var line = NewLine(p1.X, p1.Y, p2.X, p2.Y)
		var curX, curY = line.ClosestToPoint(x, y)
		var dist = point.DistanceToPoint(x, y, curX, curY)

		if dist < bestDist {
			bestDist = dist
			closestX, closestY = curX, curY
			bestP1, bestP2 = p1, p2
		}
	}

	if bestP1 != nil && bestP2 != nil {
		if closestX == bestP1.X && closestY == bestP1.Y {
			return closestX, closestY, bestP1
		}
		if closestX == bestP2.X && closestY == bestP2.Y {
			return closestX, closestY, bestP2
		}

		var proj = &n{X: closestX, Y: closestY, PathIndex: bestP1.PathIndex}
		proj.Neighbors = []*n{bestP1, bestP2}
		bestP1.Neighbors = append(bestP1.Neighbors, proj)
		bestP2.Neighbors = append(bestP2.Neighbors, proj)
		return closestX, closestY, proj
	}

	return closestX, closestY, nil
}
func remove180Turns(path []float32) []float32 {
	var nCount = len(path)
	if nCount < 6 {
		return path
	}

	var result = make([]float32, 0, nCount)
	result = append(result, path[0], path[1])
	for i := 2; i < nCount-2; i += 2 {
		var ax, ay = result[len(result)-2], result[len(result)-1]
		var bx, by = path[i], path[i+1]
		var cx, cy = path[i+2], path[i+3]

		var l1, l2 = NewLine(ax, ay, bx, by), NewLine(bx, by, cx, cy)
		var ang1, ang2 = l1.Angle(), l2.Angle()
		var diff = number.Wrap(ang2-ang1, 0, 360)
		if !number.IsWithin(diff, 180, 0.1) {
			result = append(result, bx, by)
		}
	}
	result = append(result, path[nCount-2], path[nCount-1])
	return result
}
