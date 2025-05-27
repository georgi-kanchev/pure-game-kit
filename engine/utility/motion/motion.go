package motion

import (
	"math"

	"github.com/gen2brain/raylib-go/easings"
)

func CurveBezier(unit float32, curvePoints [][2]float32) (float32, float32) {
	if len(curvePoints) == 0 {
		return float32(math.NaN()), float32(math.NaN())
	}
	if len(curvePoints) == 1 {
		return curvePoints[0][0], curvePoints[0][1]
	}

	numPoints := len(curvePoints)
	xPoints := make([]float32, numPoints)
	yPoints := make([]float32, numPoints)

	for i := range numPoints {
		xPoints[i] = curvePoints[i][0]
		yPoints[i] = curvePoints[i][1]
	}

	for k := 1; k < numPoints; k++ {
		for i := range numPoints - k {
			xPoints[i] = (1-unit)*xPoints[i] + unit*xPoints[i+1]
			yPoints[i] = (1-unit)*yPoints[i] + unit*yPoints[i+1]
		}
	}

	return xPoints[0], yPoints[0]
}
func CurveSpline(unit float32, curvePoints [][2]float32) (float32, float32) {
	if len(curvePoints) < 4 {
		return float32(math.NaN()), float32(math.NaN())
	}

	numSegments := len(curvePoints) - 3
	segmentFraction := 1.0 / float32(numSegments)
	segmentIndex := int(unit / segmentFraction)
	if segmentIndex >= numSegments {
		segmentIndex = numSegments - 1
	}

	p0 := curvePoints[segmentIndex]
	p1 := curvePoints[segmentIndex+1]
	p2 := curvePoints[segmentIndex+2]
	p3 := curvePoints[segmentIndex+3]

	u := (unit - float32(segmentIndex)*segmentFraction) / segmentFraction
	u2 := u * u
	u3 := u2 * u

	c0 := -0.5*u3 + u2 - 0.5*u
	c1 := 1.5*u3 - 2.5*u2 + 1.0
	c2 := -1.5*u3 + 2.0*u2 + 0.5*u
	c3 := 0.5*u3 - 0.5*u2

	t0 := c0*p0[0] + c1*p1[0] + c2*p2[0] + c3*p3[0]
	t1 := c0*p0[1] + c1*p1[1] + c2*p2[1] + c3*p3[1]

	return t0, t1
}

func EaseLinear(unit float32) float32 {
	return easings.LinearNone(unit, 0, 1, 1)
}
func EaseSineIn(unit float32) float32 {
	return easings.SineIn(unit, 0, 1, 1)
}
func EaseSineOut(unit float32) float32 {
	return easings.SineOut(unit, 0, 1, 1)
}
func EaseSineInOut(unit float32) float32 {
	return easings.SineInOut(unit, 0, 1, 1)
}
func EaseCircIn(unit float32) float32 {
	return easings.CircIn(unit, 0, 1, 1)
}
func EaseCircOut(unit float32) float32 {
	return easings.CircOut(unit, 0, 1, 1)
}
func EaseCircInOut(unit float32) float32 {
	return easings.CircInOut(unit, 0, 1, 1)
}
func EaseCubicIn(unit float32) float32 {
	return easings.CubicIn(unit, 0, 1, 1)
}
func EaseCubicOut(unit float32) float32 {
	return easings.CubicOut(unit, 0, 1, 1)
}
func EaseCubicInOut(unit float32) float32 {
	return easings.CubicInOut(unit, 0, 1, 1)
}
func EaseQuadIn(unit float32) float32 {
	return easings.QuadIn(unit, 0, 1, 1)
}
func EaseQuadOut(unit float32) float32 {
	return easings.QuadOut(unit, 0, 1, 1)
}
func EaseQuadInOut(unit float32) float32 {
	return easings.QuadInOut(unit, 0, 1, 1)
}
func EaseExpoIn(unit float32) float32 {
	return easings.ExpoIn(unit, 0, 1, 1)
}
func EaseExpoOut(unit float32) float32 {
	return easings.ExpoOut(unit, 0, 1, 1)
}
func EaseExpoInOut(unit float32) float32 {
	return easings.ExpoInOut(unit, 0, 1, 1)
}
func EaseBackIn(unit float32) float32 {
	return easings.BackIn(unit, 0, 1, 1)
}
func EaseBackOut(unit float32) float32 {
	return easings.BackOut(unit, 0, 1, 1)
}
func EaseBackInOut(unit float32) float32 {
	return easings.BackInOut(unit, 0, 1, 1)
}
func EaseBounceIn(unit float32) float32 {
	return easings.BounceIn(unit, 0, 1, 1)
}
func EaseBounceOut(unit float32) float32 {
	return easings.BounceOut(unit, 0, 1, 1)
}
func EaseBounceInOut(unit float32) float32 {
	return easings.BounceInOut(unit, 0, 1, 1)
}
func EaseElasticIn(unit float32) float32 {
	return easings.ElasticIn(unit, 0, 1, 1)
}
func EaseElasticOut(unit float32) float32 {
	return easings.ElasticOut(unit, 0, 1, 1)
}
func EaseElasticInOut(unit float32) float32 {
	return easings.ElasticInOut(unit, 0, 1, 1)
}
