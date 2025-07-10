package number

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func Limit(number, a, b float32) float32 {
	if a > b {
		a, b = b, a
	}
	return float32(math.Max(float64(a), math.Min(float64(number), float64(b))))
}
func LimitInt(number, a, b int) int {
	return int(Limit(float32(number), float32(a), float32(b)))
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
func MapInt(number int, fromA, fromB, toA, toB int) int {
	return int(Map(float32(number), float32(fromA), float32(fromB), float32(toA), float32(toB)))
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

func Distribute(amount int, a, b float32) []float32 {
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
func IsBetweenInt(number int, a, b int, includeA, includeB bool) bool {
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
func IsWithinInt(number, target, distance int) bool {
	return IsBetweenInt(number, target-distance, target+distance, true, true)
}

func Average(numbers ...float32) float32 {
	var sum float32
	for _, n := range numbers {
		sum += float32(n)
	}
	return sum / float32(len(numbers))
}
func AverageInt(numbers ...int) int {
	var sum int
	for _, n := range numbers {
		sum += int(n)
	}
	return sum / len(numbers)
}

func Indexes2DToIndex1D(x, y, width, height int) int {
	result := x*width + y
	max := width * height
	if result < 0 {
		return 0
	} else if result > max {
		return max
	}
	return result
}
func Index1DToIndexes2D(index, width, height int) (int, int) {
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

func ByteSizeToText(byteSize int64) string {
	const unit = 1024
	if byteSize < unit {
		return fmt.Sprintf("%d B", byteSize)
	}
	div, exp := int64(unit), 0
	for n := byteSize / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float32(byteSize)/float32(div), "KMGTPE"[exp])
}

func IsNaN(number float32) bool { return math.IsNaN(float64(number)) }
func NaN() float32              { return float32(math.NaN()) }
