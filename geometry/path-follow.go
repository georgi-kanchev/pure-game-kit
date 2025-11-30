package geometry

import (
	"pure-game-kit/execution/condition"
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
	var result = [][2]float32{{startX, startY}, {sx, sy}}
	result = append(result, startWalking(startNode, target1, target2, targetX, targetY)...)
	result = append(result, [2]float32{tx, ty}, [2]float32{targetX, targetY})
	return result
}

//=================================================================
// private

type n struct {
	X, Y                float32
	PathIndex, UniqueId int
	Visited             bool // to stop infinite loops
	Connections         []*n // can be self-joining loop or other paths joining in
	Up, Down            *n
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
		var newN = &n{X: point[0], Y: point[1], PathIndex: pathIndex, UniqueId: i, Down: prevNode}
		result = append(result, newN)

		if prevNode != nil {
			prevNode.Up = newN
		}
	}
	for _, node := range result {
		for _, target := range result {
			if node.UniqueId != target.UniqueId && node.X == target.X && node.Y == target.Y {
				node.Connections = append(node.Connections, target)
			}
		}
	}
	return result
}

func startWalking(startNode, target1, target2 *n, targetX, targetY float32) [][2]float32 {
	var upDist, up = walk(startNode, target1, target2, true, targetX, targetY)
	var downDist, down = walk(startNode, target1, target2, false, targetX, targetY)
	return condition.If(upDist < downDist, up, down)
}
func walk(node *n, target1 *n, target2 *n, up bool, targetX, targetY float32) (dist float32, pts [][2]float32) {
	var totalDistance float32
	var points = [][2]float32{}
	var prev = node
	var cur = condition.If(up, node.Up, node.Down)

	if !up { // add starting node if walking down, otherwise it's skipped
		points = append(points, [2]float32{prev.X, prev.Y})
	}

	if cur == nil { // our first step is a deadend
		totalDistance = number.Infinity()
	} // this is not the path, so indicate it with infinite distance

	for cur != nil {
		if (prev == target1 && cur == target2) ||
			(prev == target2 && cur == target1) {
			var line = NewLine(prev.X, prev.Y, cur.X, cur.Y)
			var x, y = line.ClosestToPoint(targetX, targetY)
			totalDistance += point.DistanceToPoint(prev.X, prev.Y, x, y)
			points = append(points, [2]float32{x, y})
			break // target was found
		}

		totalDistance += point.DistanceToPoint(prev.X, prev.Y, cur.X, cur.Y)
		points = append(points, [2]float32{cur.X, cur.Y})
		prev = cur
		cur = condition.If(up, cur.Up, cur.Down)

		if cur == nil { // we reached a deadend
			totalDistance = number.Infinity()
			break // this is not the path, so indicate it with infinite distance
		}
	}

	return totalDistance, points
}

func closestPointOnPath(startX, startY float32, nodes []*n) (closestX, closestY float32, n1, n2 *n) {
	var bestDist = number.Infinity()
	for i := 1; i < len(nodes); i++ {
		var p1, p2 = nodes[i-1], nodes[i]
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
