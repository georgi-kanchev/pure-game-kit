// A few functions to follow a curve or smooth-out lines.
// Useful for describing a smooth movement or path for a point or a pair of numeric values
// (time & strength for example).
package curve

import (
	"pure-game-kit/utility/number"
)

func Bezier(progress float32, curvePoints ...float32) (x, y float32) {
	var numCoords = len(curvePoints)
	if numCoords == 0 {
		return number.NaN(), number.NaN()
	}
	if numCoords == 2 {
		return curvePoints[0], curvePoints[1]
	}

	var numPoints = numCoords / 2
	var xPoints = make([]float32, numPoints)
	var yPoints = make([]float32, numPoints)

	for i := 0; i < numPoints; i++ {
		xPoints[i] = curvePoints[i*2]
		yPoints[i] = curvePoints[i*2+1]
	}

	for k := 1; k < numPoints; k++ {
		for i := 0; i < numPoints-k; i++ {
			xPoints[i] = (1-progress)*xPoints[i] + progress*xPoints[i+1]
			yPoints[i] = (1-progress)*yPoints[i] + progress*yPoints[i+1]
		}
	}

	return xPoints[0], yPoints[0]
}

func Spline(progress float32, curvePoints ...float32) (x, y float32) {
	if len(curvePoints) < 8 {
		return number.NaN(), number.NaN()
	}

	var numPoints = len(curvePoints) / 2
	var numSegments = numPoints - 3
	var segmentFraction = 1.0 / float32(numSegments)
	var segmentIndex = int(progress / segmentFraction)
	if segmentIndex >= numSegments {
		segmentIndex = numSegments - 1
	}

	var idx = segmentIndex * 2
	var p0x, p0y = curvePoints[idx], curvePoints[idx+1]
	var p1x, p1y = curvePoints[idx+2], curvePoints[idx+3]
	var p2x, p2y = curvePoints[idx+4], curvePoints[idx+5]
	var p3x, p3y = curvePoints[idx+6], curvePoints[idx+7]

	var u = (progress - float32(segmentIndex)*segmentFraction) / segmentFraction
	var u2 = u * u
	var u3 = u2 * u

	var c0 = -0.5*u3 + u2 - 0.5*u
	var c1 = 1.5*u3 - 2.5*u2 + 1.0
	var c2 = -1.5*u3 + 2.0*u2 + 0.5*u
	var c3 = 0.5*u3 - 0.5*u2

	var t0 = c0*p0x + c1*p1x + c2*p2x + c3*p3x
	var t1 = c0*p0y + c1*p1y + c2*p2y + c3*p3y

	return t0, t1
}

func SmoothPath(path ...float32) []float32 {
	if len(path) < 6 {
		return path
	}

	var nextPath []float32
	nextPath = append(nextPath, path[0], path[1])

	for j := 0; j < len(path)-2; j += 2 {
		var p0x, p0y = path[j], path[j+1]
		var p1x, p1y = path[j+2], path[j+3]

		var qx, qy = 0.75*p0x + 0.25*p1x, 0.75*p0y + 0.25*p1y
		var rx, ry = 0.25*p0x + 0.75*p1x, 0.25*p0y + 0.75*p1y

		nextPath = append(nextPath, qx, qy, rx, ry)
	}

	nextPath = append(nextPath, path[len(path)-2], path[len(path)-1])
	return nextPath
}

func SmoothPathSpline(stepsPerSegment int, path ...float32) []float32 {
	if len(path) < 6 {
		return path
	}

	var smoothed []float32
	var extendedPath = append([]float32{path[0], path[1]}, path...)
	extendedPath = append(extendedPath, path[len(path)-2], path[len(path)-1])

	for i := 0; i < len(extendedPath)-6; i += 2 {
		for j := 0; j < stepsPerSegment; j++ {
			var t = float32(j) / float32(stepsPerSegment)
			var px, py = catmullRom(
				extendedPath[i], extendedPath[i+1],
				extendedPath[i+2], extendedPath[i+3],
				extendedPath[i+4], extendedPath[i+5],
				extendedPath[i+6], extendedPath[i+7],
				t,
			)
			smoothed = append(smoothed, px, py)
		}
	}

	smoothed = append(smoothed, path[len(path)-2], path[len(path)-1])
	return smoothed
}

func SmoothPathBezier(stepsPerSegment int, path ...float32) []float32 {
	if len(path) < 6 {
		return path
	}

	var smoothed []float32
	for i := 0; i < len(path)-4; i += 4 {
		var p0x, p0y = path[i], path[i+1]
		var p1x, p1y = path[i+2], path[i+3]
		var p2x, p2y = path[i+4], path[i+5]

		for j := 0; j <= stepsPerSegment; j++ {
			var t = float32(j) / float32(stepsPerSegment)
			var invT = 1 - t
			var x = invT*invT*p0x + 2*invT*t*p1x + t*t*p2x
			var y = invT*invT*p0y + 2*invT*t*p1y + t*t*p2y
			smoothed = append(smoothed, x, y)
		}
	}
	return smoothed
}

func StraightenPath(path ...float32) []float32 {
	if len(path) < 6 {
		return path
	}

	var smoothed = make([]float32, len(path))
	copy(smoothed, path)

	for i := 2; i < len(path)-2; i += 2 {
		smoothed[i] = 0.25*path[i-2] + 0.5*path[i] + 0.25*path[i+2]
		smoothed[i+1] = 0.25*path[i-1] + 0.5*path[i+1] + 0.25*path[i+3]
	}

	return smoothed
}

//=================================================================
// private

func catmullRom(p0x, p0y, p1x, p1y, p2x, p2y, p3x, p3y, t float32) (float32, float32) {
	var tx = 0.5 * ((2 * p1x) +
		(-p0x+p2x)*t +
		(2*p0x-5*p1x+4*p2x-p3x)*t*t +
		(-p0x+3*p1x-3*p2x+p3x)*t*t*t)
	var ty = 0.5 * ((2 * p1y) +
		(-p0y+p2y)*t +
		(2*p0y-5*p1y+4*p2y-p3y)*t*t +
		(-p0y+3*p1y-3*p2y+p3y)*t*t*t)
	return tx, ty
}
