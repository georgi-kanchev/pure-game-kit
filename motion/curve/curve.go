/*
A few functions to follow a curve or smooth-out lines.
Useful for describing a smooth movement or path for a point or a pair of numeric values
(time & strength for example).
*/
package curve

import (
	"pure-game-kit/utility/number"
)

func Bezier(progress float32, curvePoints [][2]float32) (x, y float32) {
	if len(curvePoints) == 0 {
		return number.NaN(), number.NaN()
	}
	if len(curvePoints) == 1 {
		return curvePoints[0][0], curvePoints[0][1]
	}

	var numPoints = len(curvePoints)
	var xPoints = make([]float32, numPoints)
	var yPoints = make([]float32, numPoints)

	for i := range numPoints {
		xPoints[i] = curvePoints[i][0]
		yPoints[i] = curvePoints[i][1]
	}

	for k := 1; k < numPoints; k++ {
		for i := range numPoints - k {
			xPoints[i] = (1-progress)*xPoints[i] + progress*xPoints[i+1]
			yPoints[i] = (1-progress)*yPoints[i] + progress*yPoints[i+1]
		}
	}

	return xPoints[0], yPoints[0]
}
func Spline(progress float32, curvePoints [][2]float32) (x, y float32) {
	if len(curvePoints) < 4 {
		return number.NaN(), number.NaN()
	}

	var numSegments = len(curvePoints) - 3
	var segmentFraction = 1.0 / float32(numSegments)
	var segmentIndex = int(progress / segmentFraction)
	if segmentIndex >= numSegments {
		segmentIndex = numSegments - 1
	}

	var p0 = curvePoints[segmentIndex]
	var p1 = curvePoints[segmentIndex+1]
	var p2 = curvePoints[segmentIndex+2]
	var p3 = curvePoints[segmentIndex+3]
	var u = (progress - float32(segmentIndex)*segmentFraction) / segmentFraction
	var u2 = u * u
	var u3 = u2 * u
	var c0 = -0.5*u3 + u2 - 0.5*u
	var c1 = 1.5*u3 - 2.5*u2 + 1.0
	var c2 = -1.5*u3 + 2.0*u2 + 0.5*u
	var c3 = 0.5*u3 - 0.5*u2
	var t0 = c0*p0[0] + c1*p1[0] + c2*p2[0] + c3*p3[0]
	var t1 = c0*p0[1] + c1*p1[1] + c2*p2[1] + c3*p3[1]

	return t0, t1
}

func SmoothPath(path [][2]float32) [][2]float32 {
	if len(path) < 3 {
		return path
	}

	var refined = path
	var nextPath [][2]float32
	nextPath = append(nextPath, refined[0])

	for j := 0; j < len(refined)-1; j++ {
		var p0, p1 = refined[j], refined[j+1]
		var q = [2]float32{0.75*p0[0] + 0.25*p1[0], 0.75*p0[1] + 0.25*p1[1]}
		var r = [2]float32{0.25*p0[0] + 0.75*p1[0], 0.25*p0[1] + 0.75*p1[1]}

		nextPath = append(nextPath, q, r)
	}
	nextPath = append(nextPath, refined[len(refined)-1])
	refined = nextPath
	return refined
}
func SmoothPathSpline(path [][2]float32, stepsPerSegment int) [][2]float32 {
	if len(path) < 3 {
		return path
	}

	var smoothed [][2]float32
	var extendedPath = append([][2]float32{path[0]}, path...)
	extendedPath = append(extendedPath, path[len(path)-1])

	for i := 0; i < len(extendedPath)-3; i++ {
		for j := range stepsPerSegment {
			var t = float32(j) / float32(stepsPerSegment)
			var point = catmullRom(extendedPath[i], extendedPath[i+1], extendedPath[i+2], extendedPath[i+3], t)
			smoothed = append(smoothed, point)
		}
	}

	smoothed = append(smoothed, path[len(path)-1])
	return smoothed
}
func SmoothPathBezier(path [][2]float32, stepsPerSegment int) [][2]float32 {
	if len(path) < 3 {
		return path
	}

	var smoothed [][2]float32
	for i := 0; i < len(path)-2; i += 2 {
		var p0, p1, p2 = path[i], path[i+1], path[i+2]
		for j := 0; j <= stepsPerSegment; j++ {
			var t = float32(j) / float32(stepsPerSegment)
			var invT = 1 - t
			var x = invT*invT*p0[0] + 2*invT*t*p1[0] + t*t*p2[0]
			var y = invT*invT*p0[1] + 2*invT*t*p1[1] + t*t*p2[1]
			smoothed = append(smoothed, [2]float32{x, y})
		}
	}
	return smoothed
}

func StraightenPath(path [][2]float32) [][2]float32 {
	if len(path) < 3 {
		return path
	}

	smoothed := make([][2]float32, len(path))
	copy(smoothed, path)

	// Weights: 0.25 (prev), 0.50 (curr), 0.25 (next)
	for i := 1; i < len(path)-1; i++ {
		smoothed[i][0] = 0.25*path[i-1][0] + 0.5*path[i][0] + 0.25*path[i+1][0]
		smoothed[i][1] = 0.25*path[i-1][1] + 0.5*path[i][1] + 0.25*path[i+1][1]
	}

	return smoothed
}

//=================================================================
// private

func catmullRom(p0, p1, p2, p3 [2]float32, t float32) [2]float32 {
	var res [2]float32
	for i := range 2 {
		res[i] = 0.5 * ((2 * p1[i]) +
			(-p0[i]+p2[i])*t +
			(2*p0[i]-5*p1[i]+4*p2[i]-p3[i])*t*t +
			(-p0[i]+3*p1[i]-3*p2[i]+p3[i])*t*t*t)
	}
	return res
}
