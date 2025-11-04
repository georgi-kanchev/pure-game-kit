package number

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Number interface{ Float | Integer }
type Float interface{ ~float32 | ~float64 }
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func Format[T Number](number T, divideThousands bool) string {
	var str string

	switch v := any(number).(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		str = fmt.Sprintf("%d", v)
	case float32:
		str = strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		str = strconv.FormatFloat(v, 'f', -1, 64)
	default:
		str = fmt.Sprint(v) // Fallback just in case
	}

	var parts = strings.SplitN(str, ".", 2)
	var intPart = parts[0]
	var result = ""
	var n = len(intPart)
	for i, c := range intPart {
		if i > 0 && (n-i)%3 == 0 && divideThousands {
			result += " "
		}
		result += string(c)
	}

	if len(parts) == 2 {
		result += "."
		result += parts[1]
	}

	return result
}

func Limit[T Number](number, a, b T) T {
	if a > b {
		a, b = b, a
	}
	return Biggest(a, Smallest(number, b))
}
func Map[T Number](number, fromA, fromB, toA, toB T) T {
	var value T
	var deltaFrom = fromB - fromA
	if math.Abs(float64(deltaFrom)) < 0.001 {
		value = (toA + toB) / 2
		return value
	}

	value = ((number-fromA)/deltaFrom)*(toB-toA) + toA
	if math.IsNaN(float64(value)) || math.IsInf(float64(value), 0) {
		value = toA
	}
	return value
}
func Wrap[T Number](number, a, b T) T {
	if a > b {
		a, b = b, a
	}

	switch any(number).(type) {
	case float32, float64:
		var d = float64(b - a)
		if d < 0.001 {
			return a
		}
		var num = float64(number - a)
		var wrapped = math.Mod(math.Mod(num, d)+d, d) + float64(a)
		return T(wrapped)
	default: // integer types
		var d = int64(b - a)
		if d == 0 {
			return a
		}
		var num = int64(number - a)
		var wrapped = ((num % d) + d) % d
		return T(wrapped + int64(a))
	}
}

func Biggest[T Number](number, target T, other ...T) T {
	var max = number
	if target > max {
		max = target
	}
	for _, v := range other {
		if v > max {
			max = v
		}
	}
	return max
}
func Smallest[T Number](number, target T, other ...T) T {
	var min = number
	if target < min {
		min = target
	}
	for _, v := range other {
		if v < min {
			min = v
		}
	}
	return min
}
func Average[T Number](number T, others ...T) T {
	var sum T = number
	var i int
	for i = 0; i < len(others); i++ {
		var n = others[i]
		sum += n
	}
	return sum / T(1+len(others))
}
func Absolute[T Number](number T) T {
	if number < 0 {
		return -number
	}
	return number
}
func Unsign[T Number](number T) T {
	return Absolute(number)
}
func Snap[T Number](number, interval T) T {
	switch any(number).(type) {
	case float32, float64:
		var n = float64(number)
		var i = float64(interval)
		if math.IsNaN(i) || math.IsInf(n, 0) || math.Abs(i) < 0.001 {
			return number
		}
		var remainder = math.Mod(n, i)
		var halfway = i / 2.0
		if remainder < halfway {
			return T(n - remainder)
		}
		return T(n + (i - remainder))
	default: // integer types
		var num = int64(number)
		var intv = int64(interval)
		if intv == 0 {
			return number
		}
		var remainder = num % intv
		if remainder < 0 {
			remainder += intv
		}
		var halfway = intv / 2
		if remainder < halfway {
			return T(num - remainder)
		}
		return T(num + (intv - remainder))
	}
}
func Power[T Number](number, power T) T {
	return T(math.Pow(float64(number), float64(power)))
}
func SquareRoot[T Number](number T) T {
	return T(math.Sqrt(float64(number)))
}

func IsBetween[T Number](number, a, b T, includeA, includeB bool) bool {
	var lower, upper bool
	if a > b {
		var tmp = a
		a = b
		b = tmp
	}

	lower = a < number
	if includeA {
		lower = a <= number
	}

	upper = b > number
	if includeB {
		upper = b >= number
	}

	return lower && upper
}
func IsWithin[T Number](number, target, distance T) bool {
	var start = target - distance
	var end = target + distance
	return IsBetween(number, start, end, true, true)
}

func Distribute[T Number](amount int, a, b T) []T {
	if amount <= 0 {
		return []T{}
	}

	var size = b - a
	var spacing = size / T(amount+1)
	var result = make([]T, int(amount))

	for i := 1; i <= amount; i++ {
		result[int(i-1)] = a + T(i)*spacing
	}
	return result
}

func ValueMinimum[T Number]() T {
	var v T
	switch any(v).(type) {
	case int:
		return any(math.MinInt).(T)
	case int8:
		return any(math.MinInt8).(T)
	case int16:
		return any(math.MinInt16).(T)
	case int32:
		return any(math.MinInt32).(T)
	case int64:
		return any(math.MinInt64).(T)
	case uint:
		return any(0).(T)
	case uint8:
		return any(0).(T)
	case uint16:
		return any(0).(T)
	case uint32:
		return any(0).(T)
	case uint64:
		return any(0).(T)
	case float32:
		return any(-math.MaxFloat32).(T)
	case float64:
		return any(-math.MaxFloat64).(T)
	default:
		return 0
	}
}
func ValueMaximum[T Number]() T {
	var v T
	switch any(v).(type) {
	case int:
		return any(math.MaxInt).(T)
	case int8:
		return any(math.MaxInt8).(T)
	case int16:
		return any(math.MaxInt16).(T)
	case int32:
		return any(math.MaxInt32).(T)
	case int64:
		return any(math.MaxInt64).(T)
	case uint:
		return any(uint(math.MaxUint)).(T)
	case uint8:
		return any(math.MaxUint8).(T)
	case uint16:
		return any(math.MaxUint16).(T)
	case uint32:
		return any(math.MaxUint32).(T)
	case uint64:
		return any(uint64(math.MaxUint64)).(T)
	case float32:
		return any(float32(math.MaxFloat32)).(T)
	case float64:
		return any(math.MaxFloat64).(T)
	default:
		return 0
	}
}

//=================================================================
// float only

func Animate[T Float](value, target, rate T) T {
	var result T
	var factor float64 = 1.0 - math.Pow(2.0, -float64(rate))
	var delta float64 = float64(target - value)

	result = T(float64(value) + delta*factor)

	if IsWithin(float64(result), float64(target), float64(0.001)) {
		return target
	}

	return result
}

func DivisionRemainder[T Float](number, target T) T {
	return T(math.Mod(float64(number), float64(target)))
}
func Sine[T Float](number T) T {
	return T(math.Sin(float64(number)))
}
func Cosine[T Float](number T) T {
	return T(math.Cos(float64(number)))
}
func Precision[T Float](number T) int {
	for i := range 9 {
		if math.Abs(float64(number)-math.Round(float64(number))) < 1e-6 {
			return i
		}
		number *= 10
	}
	return 0
}

// negative precision ignores it
func Round[T Float](number T, precision int) T {
	if precision < 0 {
		return T(math.Round(float64(number)))
	}

	var pow = math.Pow(10, float64(precision))
	return T(math.Round(float64(number)*pow) / pow)
}

// negative precision ignores it
func RoundUp[T Float](number T, precision int) T {
	if precision < 0 {
		return T(math.Ceil(float64(number)))
	}
	var pow = math.Pow(10, float64(precision))
	return T(math.Ceil(float64(number)*pow) / pow)
}

// negative precision ignores it
func RoundDown[T Float](number T, precision int) T {
	if precision < 0 {
		return T(math.Floor(float64(number)))
	}
	var pow = math.Pow(10, float64(precision))
	return T(math.Floor(float64(number)*pow) / pow)
}

func Infinity() float32 {
	return float32(math.Inf(1))
}
func NegativeInfinity() float32 {
	return float32(math.Inf(-1))
}
func IsNaN(number float32) bool {
	return math.IsNaN(float64(number))
}
func NaN() float32 {
	return float32(math.NaN())
}

//=================================================================
// int only

func Indexes2DToIndex1D[T Integer](x, y, width, height T) T {
	var result = x*width + y
	var max = width * height
	if result < 0 {
		return 0
	} else if result > max {
		return max
	}
	return result
}
func Index1DToIndexes2D[T Integer](index, width, height T) (T, T) {
	var max = width * height
	if index < 0 {
		index = 0
	} else if index > max {
		index = max
	}
	var x = index % width
	var y = index / width
	return x, y
}
