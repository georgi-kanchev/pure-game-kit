package curves

import (
	"pure-game-kit/utility/number"

	"github.com/gen2brain/raylib-go/easings"
)

func TraceBezier(progress float32, curvePoints [][2]float32) (x, y float32) {
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
func TraceSpline(progress float32, curvePoints [][2]float32) (x, y float32) {
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

func EaseLinear(progress float32) float32 { return easings.LinearNone(progress, 0, 1, 1) }

func EaseSineIn(progress float32) float32    { return easings.SineIn(progress, 0, 1, 1) }
func EaseSineOut(progress float32) float32   { return easings.SineOut(progress, 0, 1, 1) }
func EaseSineInOut(progress float32) float32 { return easings.SineInOut(progress, 0, 1, 1) }

func EaseCircIn(progress float32) float32    { return easings.CircIn(progress, 0, 1, 1) }
func EaseCircOut(progress float32) float32   { return easings.CircOut(progress, 0, 1, 1) }
func EaseCircInOut(progress float32) float32 { return easings.CircInOut(progress, 0, 1, 1) }

func EaseCubicIn(progress float32) float32    { return easings.CubicIn(progress, 0, 1, 1) }
func EaseCubicOut(progress float32) float32   { return easings.CubicOut(progress, 0, 1, 1) }
func EaseCubicInOut(progress float32) float32 { return easings.CubicInOut(progress, 0, 1, 1) }

func EaseQuadIn(progress float32) float32    { return easings.QuadIn(progress, 0, 1, 1) }
func EaseQuadOut(progress float32) float32   { return easings.QuadOut(progress, 0, 1, 1) }
func EaseQuadInOut(progress float32) float32 { return easings.QuadInOut(progress, 0, 1, 1) }

func EaseExpoIn(progress float32) float32    { return easings.ExpoIn(progress, 0, 1, 1) }
func EaseExpoOut(progress float32) float32   { return easings.ExpoOut(progress, 0, 1, 1) }
func EaseExpoInOut(progress float32) float32 { return easings.ExpoInOut(progress, 0, 1, 1) }

func EaseBackIn(progress float32) float32    { return easings.BackIn(progress, 0, 1, 1) }
func EaseBackOut(progress float32) float32   { return easings.BackOut(progress, 0, 1, 1) }
func EaseBackInOut(progress float32) float32 { return easings.BackInOut(progress, 0, 1, 1) }

func EaseBounceIn(progress float32) float32    { return easings.BounceIn(progress, 0, 1, 1) }
func EaseBounceOut(progress float32) float32   { return easings.BounceOut(progress, 0, 1, 1) }
func EaseBounceInOut(progress float32) float32 { return easings.BounceInOut(progress, 0, 1, 1) }

func EaseElasticIn(progress float32) float32    { return easings.ElasticIn(progress, 0, 1, 1) }
func EaseElasticOut(progress float32) float32   { return easings.ElasticOut(progress, 0, 1, 1) }
func EaseElasticInOut(progress float32) float32 { return easings.ElasticInOut(progress, 0, 1, 1) }
