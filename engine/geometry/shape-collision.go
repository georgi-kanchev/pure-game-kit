package geometry

import (
	"math"
)

func (shape *Shape) Collide(velocityX, velocityY float32, targets ...*Shape) (newVelocityX, newVelocityY float32) {
	for _, target := range targets {
		if !shape.inBoundingBoxShape(*target) {
			continue
		}

		var cax = shape.MinX + (shape.MaxX-shape.MinX)/2
		var cay = shape.MinY + (shape.MaxY-shape.MinY)/2
		var cbx = target.MinX + (target.MaxX-target.MinX)/2
		var cby = target.MinY + (target.MaxY-target.MinY)/2
		var corners = shape.CornerPoints()
		var targetCorners = target.CornerPoints()

		corners = corners[:len(corners)-1]
		targetCorners = targetCorners[:len(targetCorners)-1]

		var projectPolygon = func(axisX, axisY float32, points [][2]float32) (min, max float32) {
			var dot float32 = axisX*points[0][0] + axisY*points[0][1]
			min = dot
			max = dot

			for i := range points {
				var d float32 = axisX*points[i][0] + axisY*points[i][1]
				if d < min {
					min = d
				} else if d > max {
					max = d
				}
			}
			return
		}
		var computeEdges = func(points [][2]float32) [][2]float32 {
			var edges = make([][2]float32, len(points))
			var count int = len(points)
			for i := range count {
				var p1, p2 = points[i], points[(i+1)%count]
				edges[i][0], edges[i][1] = p2[0]-p1[0], p2[1]-p1[1]
			}
			return edges
		}

		var willIntersect = true
		var edgesA = computeEdges(corners)
		var edgesB = computeEdges(targetCorners)
		var edgeCountA = len(edgesA)
		var edgeCountB = len(edgesB)
		var minIntervalDistance = float32(math.Inf(1))
		var tAxisX, tAxisY float32

		for edgeIndex := 0; edgeIndex < edgeCountA+edgeCountB; edgeIndex++ {
			var edgeX, edgeY float32
			if edgeIndex < edgeCountA {
				edgeX, edgeY = edgesA[edgeIndex][0], edgesA[edgeIndex][1]
			} else {
				edgeX, edgeY = edgesB[edgeIndex-edgeCountA][0], edgesB[edgeIndex-edgeCountA][1]
			}

			var axisX, axisY = -edgeY, edgeX
			var axisLen = float32(math.Hypot(float64(axisX), float64(axisY)))
			if axisLen != 0 {
				axisX /= axisLen
				axisY /= axisLen
			} else {
				continue
			}

			var minA, maxA = projectPolygon(axisX, axisY, corners)
			var minB, maxB = projectPolygon(axisX, axisY, targetCorners)

			var velocityProjection = axisX*velocityX + axisY*velocityY
			if velocityProjection < 0 {
				minA += velocityProjection
			} else {
				maxA += velocityProjection
			}

			var iDist = minA - maxB
			if minA < minB {
				iDist = minB - maxA
			}

			if iDist > 0 {
				willIntersect = false
			}

			if !willIntersect {
				break
			}

			var absInterval = float32(math.Abs(float64(iDist)))
			if absInterval < minIntervalDistance {
				minIntervalDistance = absInterval
				tAxisX, tAxisY = axisX, axisY

				var dx = cax - cbx
				var dy = cay - cby
				if dx*tAxisX+dy*tAxisY < 0 {
					tAxisX, tAxisY = -tAxisX, -tAxisY
				}
			}
		}

		if willIntersect {
			velocityX += tAxisX * minIntervalDistance
			velocityY += tAxisY * minIntervalDistance
		}

	}
	return velocityX, velocityY
}
