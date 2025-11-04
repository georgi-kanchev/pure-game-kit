package curves

import (
	"math"
	"pure-game-kit/utility/number"
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

func EaseLinear(progress float32) float32 { return progress }

func EaseSineIn(progress float32) float32  { return float32(1 - math.Cos(float64(progress)*math.Pi/2)) }
func EaseSineOut(progress float32) float32 { return float32(math.Sin(float64(progress) * math.Pi / 2)) }
func EaseSineInOut(progress float32) float32 {
	return float32(-0.5 * (math.Cos(math.Pi*float64(progress)) - 1))
}

func EaseCircIn(progress float32) float32 {
	return float32(1 - math.Sqrt(1-math.Pow(float64(progress), 2)))
}
func EaseCircOut(progress float32) float32 {
	return float32(math.Sqrt(1 - math.Pow(float64(progress-1), 2)))
}
func EaseCircInOut(progress float32) float32 {
	p2 := progress * 2
	if p2 < 1 {
		return float32(-0.5 * (math.Sqrt(1-math.Pow(float64(p2), 2)) - 1))
	}
	p2 -= 2
	return float32(0.5 * (math.Sqrt(1-math.Pow(float64(p2), 2)) + 1))
}

func EaseCubicIn(progress float32) float32  { return progress * progress * progress }
func EaseCubicOut(progress float32) float32 { progress--; return progress*progress*progress + 1 }
func EaseCubicInOut(progress float32) float32 {
	p2 := progress * 2
	if p2 < 1 {
		return 0.5 * p2 * p2 * p2
	}
	p2 -= 2
	return 0.5 * (p2*p2*p2 + 2)
}

func EaseQuadIn(progress float32) float32  { return progress * progress }
func EaseQuadOut(progress float32) float32 { return -progress * (progress - 2) }
func EaseQuadInOut(progress float32) float32 {
	p2 := progress * 2
	if p2 < 1 {
		return 0.5 * p2 * p2
	}
	p2--
	return -0.5 * (p2*(p2-2) - 1)
}

func EaseExpoIn(progress float32) float32 {
	if progress == 0 {
		return 0
	}
	return float32(math.Pow(2, 10*(float64(progress)-1)))
}
func EaseExpoOut(progress float32) float32 {
	if progress == 1 {
		return 1
	}
	return float32(1 - math.Pow(2, -10*float64(progress)))
}
func EaseExpoInOut(progress float32) float32 {
	if progress == 0 {
		return 0
	}
	if progress == 1 {
		return 1
	}
	p2 := progress * 2
	if p2 < 1 {
		return float32(0.5 * math.Pow(2, 10*(float64(p2)-1)))
	}
	return float32(0.5 * (2 - math.Pow(2, -10*(float64(p2)-1))))
}

func EaseBackIn(progress float32) float32 {
	const s float32 = 1.70158
	return progress * progress * ((s+1)*progress - s)
}
func EaseBackOut(progress float32) float32 {
	const s float32 = 1.70158
	progress--
	return progress*progress*((s+1)*progress+s) + 1
}
func EaseBackInOut(progress float32) float32 {
	const s float32 = 1.70158 * 1.525
	p2 := progress * 2
	if p2 < 1 {
		return 0.5 * (p2 * p2 * ((s+1)*p2 - s))
	}
	p2 -= 2
	return 0.5 * (p2*p2*((s+1)*p2+s) + 2)
}

func EaseBounceOut(progress float32) float32 {
	if progress < 1/2.75 {
		return 7.5625 * progress * progress
	} else if progress < 2/2.75 {
		progress -= 1.5 / 2.75
		return 7.5625*progress*progress + 0.75
	} else if progress < 2.5/2.75 {
		progress -= 2.25 / 2.75
		return 7.5625*progress*progress + 0.9375
	}
	progress -= 2.625 / 2.75
	return 7.5625*progress*progress + 0.984375
}

func EaseBounceIn(progress float32) float32 {
	return 1 - EaseBounceOut(1-progress)
}
func EaseBounceInOut(progress float32) float32 {
	if progress < 0.5 {
		return EaseBounceIn(progress*2) * 0.5
	}
	return EaseBounceOut(progress*2-1)*0.5 + 0.5
}

func EaseElasticIn(progress float32) float32 {
	if progress == 0 || progress == 1 {
		return progress
	}
	return float32(-math.Pow(2, 10*float64(progress-1)) * math.Sin((float64(progress-1.1) * 5 * math.Pi)))
}
func EaseElasticOut(progress float32) float32 {
	if progress == 0 || progress == 1 {
		return progress
	}
	return float32(math.Pow(2, -10*float64(progress))*math.Sin((float64(progress-0.1)*5*math.Pi)) + 1)
}
func EaseElasticInOut(progress float32) float32 {
	if progress == 0 || progress == 1 {
		return progress
	}
	p2 := progress * 2
	if p2 < 1 {
		return float32(-0.5 * (math.Pow(2, 10*(float64(p2)-1)) * math.Sin((float64(p2-1.1) * 5 * math.Pi))))
	}
	p2--
	return float32(math.Pow(2, -10*float64(p2))*math.Sin((float64(p2-0.1)*5*math.Pi))*0.5 + 1)
}
