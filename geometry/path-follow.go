package geometry

import (
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/point"
)

// calculates the minimal subsections of the routes to traverse to reach target
//
// start and target point can be anywhere (not necessarily on the paths)
//
// multiple paths can be separated by [NaN, NaN]
//
// separate paths need to be sharing points to connect - disconnected paths are discarded
func FollowPaths(startX, startY, targetX, targetY float32, paths ...[2]float32) [][2]float32 {
	var allNodes = createNodes(paths)
	var sx, sy, startNode, _ = closestPointOnPath(startX, startY, allNodes)
	var tx, ty, target1, target2 = closestPointOnPath(targetX, targetY, allNodes)
	var _, path = walk(startNode, target1, target2, sx, sy, targetX, targetY)
	var result = [][2]float32{{startX, startY}, {sx, sy}}
	if len(path) == 0 { // no path was found, perhaps the end target is a disconnected one
		return result
	}

	result = append(result, path...)
	result = append(result, [2]float32{tx, ty}, [2]float32{targetX, targetY})
	return result
}

//=================================================================
// private

type n struct {
	X, Y                float32
	PathIndex, UniqueId int
	Neighbors           []*n

	SearchPrev *n
	SearchDist float32
	SearchSeen bool
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
		var newN = &n{X: point[0], Y: point[1], PathIndex: pathIndex, UniqueId: i}
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

func walk(start *n, target1 *n, target2 *n, startX, startY, targetX, targetY float32) (distance float32, path [][2]float32) {
	var queue = []*n{start}
	var hack = false
	start.SearchSeen = true

	for _, nb := range start.Neighbors {
		if start.PathIndex == nb.PathIndex && nb.UniqueId < start.UniqueId {
			hack = true
			break
		}
	}
	if !hack {
		start.X, start.Y = startX, startY
	}

	for len(queue) > 0 {
		var cur = queue[0]
		queue = queue[1:]

		for _, nb := range cur.Neighbors {
			if (cur == target1 && nb == target2) ||
				(cur == target2 && nb == target1) {
				var line = NewLine(cur.X, cur.Y, nb.X, nb.Y)
				var x, y = line.ClosestToPoint(targetX, targetY)
				var dist = cur.SearchDist + point.DistanceToPoint(cur.X, cur.Y, x, y)
				var path = reconstructFromNode(cur)
				return dist, path // we found the target
			}

			if nb.SearchSeen {
				continue
			}

			nb.SearchSeen = true
			nb.SearchPrev = cur
			nb.SearchDist = cur.SearchDist + point.DistanceToPoint(cur.X, cur.Y, nb.X, nb.Y)

			queue = append(queue, nb)
		}
	}

	return number.Infinity(), nil
}

func reconstructFromNode(n *n) [][2]float32 {
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

func closestPointOnPath(startX, startY float32, nodes []*n) (closestX, closestY float32, n1, n2 *n) {
	var bestDist = number.Infinity()
	for i := 1; i < len(nodes); i++ {
		var p1, p2 = nodes[i-1], nodes[i]
		if p1.PathIndex != p2.PathIndex {
			continue // ignore path connections
		}

		var line = NewLine(p1.X, p1.Y, p2.X, p2.Y)
		var curClosestX, curClosestY = line.ClosestToPoint(startX, startY)
		var dist = point.DistanceToPoint(startX, startY, curClosestX, curClosestY)
		if dist < bestDist {
			bestDist = dist
			closestX, closestY = curClosestX, curClosestY
			n1, n2 = p1, p2
		}
	}
	return closestX, closestY, n1, n2
}
