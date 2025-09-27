package number

import "math"

const MaxInt = 2147483647
const MinInt = -MaxInt

func Biggest(number, target float32, other ...float32) float32 {
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
func BiggestInt(number, target int, other ...int) int {
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
func Smallest(number, target float32, other ...float32) float32 {
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
func SmallestInt(number, target int, other ...int) int {
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

func Limit(number, a, b float32) float32 {
	if a > b {
		a, b = b, a
	}
	return Biggest(a, Smallest(number, b))
}
func LimitInt(number, a, b int) int {
	return int(Limit(float32(number), float32(a), float32(b)))
}

func WrapRange(number, a, b float32) float32 {
	if a > b {
		a, b = b, a
	}
	var d = b - a
	if d < 0.001 {
		return a
	}
	return float32(math.Mod(math.Mod(float64(number-a), float64(d))+float64(d), float64(d))) + a
}
func WrapRangeInt(number, a, b int) int {
	return int(WrapRange(float32(number), float32(a), float32(b)))
}

func Wrap(number, target float32) float32 {
	if target == 0 {
		return 0
	}
	return float32(math.Mod(math.Mod(float64(number), float64(target))+float64(target), float64(target)))
}
func WrapInt(number, target int) int {
	if target == 0 {
		return 0
	}
	return ((number % target) + target) % target
}

func Snap(number, interval float32) float32 {
	if math.IsNaN(float64(interval)) || math.IsInf(float64(number), 0) || math.Abs(float64(interval)) < 0.001 {
		return number
	}
	var remainder = float32(math.Mod(float64(number), float64(interval)))
	var halfway = interval / 2.0
	if remainder < halfway {
		return number - remainder
	}
	return number + (interval - remainder)
}
func SnapInt(number, interval int) int {
	if interval == 0 {
		return number // avoid divide-by-zero
	}
	var remainder = number % interval
	var halfway = interval / 2
	if remainder < 0 {
		remainder += interval // handle negatives consistently
	}

	if remainder < halfway {
		return number - remainder
	}
	return number + (interval - remainder)
}

func Map(number, fromA, fromB, toA, toB float32) float32 {
	if math.Abs(float64(fromB-fromA)) < 0.001 {
		return (toA + toB) / 2
	}
	var value = ((number-fromA)/(fromB-fromA))*(toB-toA) + toA
	if math.IsNaN(float64(value)) || math.IsInf(float64(value), 0) {
		return toA
	}
	return value
}
func MapInt(number, fromA, fromB, toA, toB int) int {
	return int(Map(float32(number), float32(fromA), float32(fromB), float32(toA), float32(toB)))
}

func IsBetween(number, a, b float32, includeA, includeB bool) bool {
	if a > b {
		a, b = b, a
	}
	var l = a < number
	if includeA {
		l = a <= number
	}
	var u = b > number
	if includeB {
		u = b >= number
	}
	return l && u
}
func IsBetweenInt(number, a, b int, includeA, includeB bool) bool {
	if a > b {
		a, b = b, a
	}
	var l = a < number
	if includeA {
		l = a <= number
	}
	var u = b > number
	if includeB {
		u = b >= number
	}
	return l && u
}

func IsWithin(number, target, distance float32) bool {
	return IsBetween(number, target-distance, target+distance, true, true)
}
func IsWithinInt(number, target, distance int) bool {
	return IsBetweenInt(number, target-distance, target+distance, true, true)
}

func Average(number float32, others ...float32) float32 {
	var sum = number
	for _, n := range others {
		sum += n
	}
	return sum / float32(1+len(others))
}
func AverageInt(number int, others ...int) int {
	var sum = number
	for _, n := range others {
		sum += n
	}
	return sum / (1 + len(others))
}

func Indexes2DToIndex1D(x, y, width, height int) int {
	var result = x*width + y
	var max = width * height
	if result < 0 {
		return 0
	} else if result > max {
		return max
	}
	return result
}
func Index1DToIndexes2D(index, width, height int) (int, int) {
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

func Absolute(number float32) float32 {
	if number < 0 {
		return -number
	}
	return number
}
func AbsoluteInt(number int) int {
	if number < 0 {
		return -number
	}
	return number
}
func Unsign(number float32) float32 {
	return Absolute(number)
}
func UnsignInt(number int) int {
	return AbsoluteInt(number)
}

func Power(number, power float32) float32 {
	return float32(math.Pow(float64(number), float64(power)))
}
func PowerInt(number, power int) int {
	return int(math.Pow(float64(number), float64(power)))
}

func AlmostEquals(number, target, tolerance float32) bool {
	return float32(Unsign(number-target)) <= tolerance
}

func Animate(value, target, rate float32) float32 {
	value += (target - value) * (1.0 - Power(2.0, -rate))

	if AlmostEquals(value, target, 0.001) {
		return target
	}

	return value
}

//=================================================================
// float only

func DivisionRemainder(number, target float32) float32 {
	return float32(math.Mod(float64(number), float64(target)))
}

func SquareRoot(number float32) float32 {
	return float32(math.Sqrt(float64(number)))
}

func Sine(number float32) float32 {
	return float32(math.Sin(float64(number)))
}
func Cosine(number float32) float32 {
	return float32(math.Cos(float64(number)))
}

func Infinity() float32 {
	return float32(math.Inf(1))
}
func NegativeInfinity() float32 {
	return float32(math.Inf(-1))
}

func Distribute(amount int, a, b float32) []float32 {
	if amount <= 0 {
		return []float32{}
	}

	var result = make([]float32, amount)
	var size = b - a
	var spacing = size / float32(amount+1)

	for i := 1; i <= int(amount); i++ {
		result[i-1] = a + float32(i)*spacing
	}

	return result
}
func Precision(number float32) int {
	for i := range 9 {
		if math.Abs(float64(number)-math.Round(float64(number))) < 1e-6 {
			return i
		}
		number *= 10
	}
	return 0
}

// negative precision ignores it
func Round(number float32, precision int) float32 {
	if precision < 0 {
		return float32(math.Round(float64(number)))
	}
	var pow = math.Pow(10, float64(precision))
	return float32(math.Round(float64(number)*pow) / pow)
}

// negative precision ignores it
func RoundUp(number float32, precision int) float32 {
	if precision < 0 {
		return float32(math.Ceil(float64(number)))
	}
	var pow = math.Pow(10, float64(precision))
	return float32(math.Ceil(float64(number)*pow) / pow)
}

// negative precision ignores it
func RoundDown(number float32, precision int) float32 {
	if precision < 0 {
		return float32(math.Floor(float64(number)))
	}
	var pow = math.Pow(10, float64(precision))
	return float32(math.Floor(float64(number)*pow) / pow)
}

func IsNaN(number float32) bool {
	return math.IsNaN(float64(number))
}
func NaN() float32 {
	return float32(math.NaN())
}
