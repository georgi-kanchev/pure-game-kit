package utility

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	point "pure-tile-kit/engine/utility/point"

	easings "github.com/gen2brain/raylib-go/easings"
)

type Ease byte
type Curve byte

const (
	EaseLinear Ease = iota
	EaseSine
	EaseQuad
	EaseCubic
	EaseExpo
	EaseCirc
	EaseBack
	EaseElastic
	EaseBounce
)
const (
	CurveIn Curve = iota
	CurveOut
	CurveInOut
)

func AnimateEase(unit float32, ease Ease, curve Curve) float32 {
	switch ease {
	case EaseLinear:
		return easings.LinearNone(unit, 0, 1, 1)

	case EaseSine:
		switch curve {
		case CurveIn:
			return easings.SineIn(unit, 0, 1, 1)
		case CurveOut:
			return easings.SineOut(unit, 0, 1, 1)
		case CurveInOut:
			return easings.SineInOut(unit, 0, 1, 1)
		}

	case EaseQuad:
		switch curve {
		case CurveIn:
			return easings.QuadIn(unit, 0, 1, 1)
		case CurveOut:
			return easings.QuadOut(unit, 0, 1, 1)
		case CurveInOut:
			return easings.QuadInOut(unit, 0, 1, 1)
		}

	case EaseCubic:
		switch curve {
		case CurveIn:
			return easings.CubicIn(unit, 0, 1, 1)
		case CurveOut:
			return easings.CubicOut(unit, 0, 1, 1)
		case CurveInOut:
			return easings.CubicInOut(unit, 0, 1, 1)
		}

	case EaseExpo:
		switch curve {
		case CurveIn:
			return easings.ExpoIn(unit, 0, 1, 1)
		case CurveOut:
			return easings.ExpoOut(unit, 0, 1, 1)
		case CurveInOut:
			return easings.ExpoInOut(unit, 0, 1, 1)
		}

	case EaseCirc:
		switch curve {
		case CurveIn:
			return easings.CircIn(unit, 0, 1, 1)
		case CurveOut:
			return easings.CircOut(unit, 0, 1, 1)
		case CurveInOut:
			return easings.CircInOut(unit, 0, 1, 1)
		}

	case EaseBack:
		switch curve {
		case CurveIn:
			return easings.BackIn(unit, 0, 1, 1)
		case CurveOut:
			return easings.BackOut(unit, 0, 1, 1)
		case CurveInOut:
			return easings.BackInOut(unit, 0, 1, 1)
		}

	case EaseElastic:
		switch curve {
		case CurveIn:
			return easings.ElasticIn(unit, 0, 1, 1)
		case CurveOut:
			return easings.ElasticOut(unit, 0, 1, 1)
		case CurveInOut:
			return easings.ElasticInOut(unit, 0, 1, 1)
		}

	case EaseBounce:
		switch curve {
		case CurveIn:
			return easings.BounceIn(unit, 0, 1, 1)
		case CurveOut:
			return easings.BounceOut(unit, 0, 1, 1)
		case CurveInOut:
			return easings.BounceInOut(unit, 0, 1, 1)
		}
	}

	return NaN()
}
func AnimateBezier(unit float32, curvePoints []point.F) point.F {
	if len(curvePoints) == 0 {
		return point.F{X: float32(math.NaN()), Y: float32(math.NaN())}
	}
	if len(curvePoints) == 1 {
		return curvePoints[0]
	}

	numPoints := len(curvePoints)
	xPoints := make([]float32, numPoints)
	yPoints := make([]float32, numPoints)

	for i := range numPoints {
		xPoints[i] = curvePoints[i].X
		yPoints[i] = curvePoints[i].Y
	}

	for k := 1; k < numPoints; k++ {
		for i := range numPoints - k {
			xPoints[i] = (1-unit)*xPoints[i] + unit*xPoints[i+1]
			yPoints[i] = (1-unit)*yPoints[i] + unit*yPoints[i+1]
		}
	}

	return point.F{X: xPoints[0], Y: yPoints[0]}
}
func AnimateSpline(unit float32, curvePoints []point.F) point.F {
	if len(curvePoints) < 4 {
		return point.F{X: float32(math.NaN()), Y: float32(math.NaN())}
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

	t0 := c0*p0.X + c1*p1.X + c2*p2.X + c3*p3.X
	t1 := c0*p0.Y + c1*p1.Y + c2*p2.Y + c3*p3.Y

	return point.F{X: t0, Y: t1}
}

func Limit(number, a, b float32) float32 {
	if a > b {
		a, b = b, a
	}
	return float32(math.Max(float64(a), math.Min(float64(number), float64(b))))
}
func LimitInt(number, a, b int32) int32 {
	return int32(Limit(float32(number), float32(a), float32(b)))
}

func WrapRange(number, a, b float32) float32 {
	if a > b {
		a, b = b, a
	}
	d := b - a
	if d < 0.001 {
		return a
	}
	return float32(math.Mod(math.Mod(float64(number-a), float64(d))+float64(d), float64(d))) + a
}
func WrapRangeInt(number, a, b int32) int32 {
	return int32(WrapRange(float32(number), float32(a), float32(b)))
}

func Wrap(number, target float32) float32 {
	if target == 0 {
		return 0
	}
	return float32(math.Mod(math.Mod(float64(number), float64(target))+float64(target), float64(target)))
}
func WrapInt(number, target int32) int32 {
	if target == 0 {
		return 0
	}
	return ((number % target) + target) % target
}

func Snap(number, interval float32) float32 {
	if math.IsNaN(float64(interval)) || math.IsInf(float64(number), 0) || math.Abs(float64(interval)) < 0.001 {
		return number
	}
	remainder := float32(math.Mod(float64(number), float64(interval)))
	halfway := interval / 2.0
	if remainder < halfway {
		return number - remainder
	}
	return number + (interval - remainder)
}

func Map(number float32, fromA, fromB, toA, toB float32) float32 {
	if math.Abs(float64(fromB-fromA)) < 0.001 {
		return (toA + toB) / 2
	}
	value := ((number-fromA)/(fromB-fromA))*(toB-toA) + toA
	if math.IsNaN(float64(value)) || math.IsInf(float64(value), 0) {
		return toA
	}
	return value
}
func MapInt(number int32, fromA, fromB, toA, toB int32) int32 {
	return int32(Map(float32(number), float32(fromA), float32(fromB), float32(toA), float32(toB)))
}

func PadZeros(number float32, amountOfZeros int) string {
	if amountOfZeros == 0 {
		return strconv.FormatFloat(float64(number), 'f', -1, 32)
	}
	if amountOfZeros < 0 {
		width := -amountOfZeros
		return fmt.Sprintf("%0*d", width, int(number))
	}
	return fmt.Sprintf("%.*f", amountOfZeros, number)
}
func PadZerosInt(number int, amountOfZeros int) string {
	if amountOfZeros <= 0 {
		return fmt.Sprintf("%d", number)
	}
	return fmt.Sprintf("%0*d", amountOfZeros, number)
}

func Distribute(amount int32, a, b float32) []float32 {
	if amount <= 0 {
		return []float32{}
	}

	result := make([]float32, amount)
	size := b - a
	spacing := size / float32(amount+1)

	for i := 1; i <= int(amount); i++ {
		result[i-1] = a + float32(i)*spacing
	}

	return result
}
func Precision(number float32) int {
	s := strconv.FormatFloat(float64(number), 'f', -1, 32)
	parts := strings.Split(s, ".")
	if len(parts) == 2 {
		return len(parts[1])
	}
	return 0
}

func IsBetween(number float32, a, b float32, includeA, includeB bool) bool {
	if a > b {
		a, b = b, a
	}
	l := a < number
	if includeA {
		l = a <= number
	}
	u := b > number
	if includeB {
		u = b >= number
	}
	return l && u
}
func IsBetweenInt(number int32, a, b int32, includeA, includeB bool) bool {
	if a > b {
		a, b = b, a
	}
	l := a < number
	if includeA {
		l = a <= number
	}
	u := b > number
	if includeB {
		u = b >= number
	}
	return l && u
}

func IsWithin(number, target, distance float32) bool {
	return IsBetween(number, target-distance, target+distance, true, true)
}
func IsWithinInt(number, target, distance int32) bool {
	return IsBetweenInt(number, target-distance, target+distance, true, true)
}

func Random(a, b, seed float32) float32 {
	if a == b {
		return a
	}
	if a > b {
		a, b = b, a
	}
	diff := b - a

	var intSeed int32
	if math.IsNaN(float64(seed)) {
		intSeed = int32(time.Now().UnixNano())
	} else {
		intSeed = int32(math.Float32bits(seed))
	}

	intSeed = (1103515245*intSeed + 12345) % 2147483647
	normalized := float32(intSeed&0x7FFFFFFF) / 2147483648.0

	return a + normalized*diff
}
func RandomInt(a, b int32, seed float32) int32 {
	return int32(Random(float32(a), float32(b), seed))
}

func HasChance(percent, seed float32) bool {
	if percent <= 0.000001 {
		return false
	}
	percent = Limit(percent, 0, 100)
	n := Random(0, 100, seed)
	return n <= percent
}
func HasChanceInt(percent int32, seed float32) bool {
	return HasChance(float32(percent), seed)
}

func IsNaN(number float32) bool {
	return math.IsNaN(float64(number))
}
func NaN() float32 {
	return float32(math.NaN())
}

func ToSeed(number int32, parameters []int32) int32 {
	seed := uint64(2654435769)

	seed = hashSeed(seed, number)
	for _, p := range parameters {
		seed = hashSeed(seed, p)
	}
	return int32(seed)
}

func ToIndex1D(x, y, width, height int32) int32 {
	result := x*width + y
	max := width * height
	if result < 0 {
		return 0
	} else if result > max {
		return max
	}
	return result
}
func ToIndexes2D(index, width, height int32) (int32, int32) {
	max := width * height
	if index < 0 {
		index = 0
	} else if index > max {
		index = max
	}
	x := index % width
	y := index / width
	return x, y
}

// private

func hashSeed(seed uint64, a int32) uint64 {
	seed ^= uint64(a)
	seed = (seed ^ (seed >> 16)) * 2246822519
	seed = (seed ^ (seed >> 13)) * 3266489917
	seed ^= seed >> 16
	return seed
}
