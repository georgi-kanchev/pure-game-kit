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
func FollowPaths(startX, startY, targetX, targetY float32, paths ...[2]float32) [][2]float32 {
	var allNodes = createNodes(paths)
	var sx, sy, startNode = closestPointOnPath(startX, startY, allNodes)
	var tx, ty, targetNode = closestPointOnPath(targetX, targetY, allNodes)
	var path = startNode.walk(targetNode)
	var result = [][2]float32{{startX, startY}, {sx, sy}}
	if len(path) == 0 { // no path was found, perhaps the end target is a disconnected one
		return result
	}

	result = append(result, path...)
	result = append(result, [2]float32{tx, ty}, [2]float32{targetX, targetY})
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

func createNodes(points [][2]float32) []*n {
	var result = []*n{}
	var pathIndex = 0
	for i, point := range points {
		if number.IsNaN(point[0]) || number.IsNaN(point[1]) {
			pathIndex++
			continue
		}
		var prevNode *n
		if i != 0 {
			prevNode = result[len(result)-1]
		}
		var newN = &n{X: point[0], Y: point[1], PathIndex: pathIndex, SearchDist: number.Infinity()}
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
func (start *n) walk(end *n) [][2]float32 {
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
func (n *n) reconstructPath() [][2]float32 {
	var result [][2]float32
	var unique = map[[2]float32]any{}
	for n != nil {
		var point = [2]float32{n.X, n.Y}
		var _, has = unique[point]
		if !has {
			result = append([][2]float32{point}, result...)
		}
		unique[point] = nil
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
func remove180Turns(path [][2]float32) [][2]float32 {
	var n = len(path)
	if n < 3 {
		return path
	}

	var result = make([][2]float32, 0, n)
	result = append(result, path[0])
	for i := 1; i < n-1; i++ {
		var a, b, c = result[len(result)-1], path[i], path[i+1]
		var l1, l2 = NewLine(a[0], a[1], b[0], b[1]), NewLine(b[0], b[1], c[0], c[1])
		var ang1, ang2 = l1.Angle(), l2.Angle()
		var diff = number.Wrap(ang2-ang1, 0, 360)
		if !number.IsWithin(diff, 180, 0.1) {
			result = append(result, b)
		}
	}
	result = append(result, path[n-1])
	return result
}
