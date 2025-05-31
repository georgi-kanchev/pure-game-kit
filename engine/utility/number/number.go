package utility

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type Noise byte

const (
	NoisePerlin Noise = iota
	NoiseOpenSimplex
	NoiseWorley
	NoiseVoronoi
	NoiseValue
	NoiseValueCubic
)

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

func Seed(keys []int32) int32 {
	hashSeed := func(seed uint64, a int32) uint64 {
		seed ^= uint64(a)
		seed = (seed ^ (seed >> 16)) * 2246822519
		seed = (seed ^ (seed >> 13)) * 3266489917
		seed ^= seed >> 16
		return seed
	}
	seed := uint64(2654435769)

	for _, p := range keys {
		seed = hashSeed(seed, p)
	}
	return int32(seed)
}
func Random(a, b, seed float32) float32 {
	if a == b {
		return a
	}
	if a > b {
		a, b = b, a
	}
	diff := b - a
	intSeed := int32(seed * 2147483647)
	intSeed = (1103515245*intSeed + 12345) % 2147483647
	normalized := float32(intSeed) / 2147483647.0
	return a + normalized*diff
}
func RandomInt(a, b int32, seed float32) int32 {
	return int32(Random(float32(a), float32(b), seed))
}
func GenerateNoise(noise Noise, x, y, scale, seed float32) float32 {
	var intSeed int32 = floatToIntSeed(seed)
	switch noise {
	case NoisePerlin:
		return noisePerlin(x, y, scale, intSeed)
	case NoiseOpenSimplex:
		return noiseOpenSimplex(x, y, scale, intSeed)
	case NoiseWorley:
		return noiseWorley(x, y, scale, intSeed)
	case NoiseVoronoi:
		return noiseVoronoi(x, y, scale, intSeed)
	case NoiseValue:
		return noiseValue(x, y, scale, intSeed)
	case NoiseValueCubic:
		return noiseValueCubic(x, y, scale, intSeed)
	}

	return NaN()
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

func Indexes2DToIndex1D(x, y, width, height int32) int32 {
	result := x*width + y
	max := width * height
	if result < 0 {
		return 0
	} else if result > max {
		return max
	}
	return result
}
func Index1DToIndexes2D(index, width, height int32) (int32, int32) {
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
func floatToIntSeed(seed float32) int32 {
	var intSeed int32 = 0
	if math.IsNaN(float64(seed)) {
		intSeed = int32(time.Now().UnixNano())
	} else {
		intSeed = int32(math.Float32bits(seed))
	}

	return intSeed
}

func noiseOpenSimplex(x, y, scale float32, seed int32) float32 {
	const stretch2D = -0.211324865405187 // (1/sqrt(2+1)-1)/2
	const squish2D = 0.366025403784439   // (sqrt(2+1)-1)/2
	const norm2D = 47.0

	// Gradient directions
	gradients := [8][2]float32{
		{5, 2}, {2, 5}, {-5, 2}, {-2, 5},
		{5, -2}, {2, -5}, {-5, -2}, {-2, -5},
	}

	// Permutation table (generated using Random and Seed)
	getPerm := func(i, j int32) int {
		subSeed := Seed([]int32{seed, i, j})
		r := Random(0, 255, float32(subSeed))
		return int(r) & 255
	}

	// Scale input
	x *= scale
	y *= scale

	// Place input on grid
	s := (x + y) * float32(stretch2D)
	xf := x + s
	yf := y + s

	xi := int32(math.Floor(float64(xf)))
	yi := int32(math.Floor(float64(yf)))

	t := float32(xi+yi) * float32(squish2D)
	X0 := float32(xi) - t
	Y0 := float32(yi) - t
	x0 := x - X0
	y0 := y - Y0

	var x1, y1 float32
	var i1, j1 int32
	if x0 > y0 {
		i1, j1 = 1, 0
	} else {
		i1, j1 = 0, 1
	}

	x1 = x0 - float32(i1) + float32(squish2D)
	y1 = y0 - float32(j1) + float32(squish2D)
	x2 := x0 - 1 + 2*float32(squish2D)
	y2 := y0 - 1 + 2*float32(squish2D)

	contrib := func(dx, dy float32, i, j int32) float32 {
		attsq := 2 - dx*dx - dy*dy
		if attsq > 0 {
			pi := getPerm(i, j) % 8
			grad := gradients[pi]
			att4 := attsq * attsq
			return att4 * att4 * (grad[0]*dx + grad[1]*dy)
		}
		return 0
	}

	value := contrib(x0, y0, xi, yi)
	value += contrib(x1, y1, xi+i1, yi+j1)
	value += contrib(x2, y2, xi+1, yi+1)

	// Normalize to [0, 1]
	return float32(math.Max(0, math.Min(1, float64(value)/norm2D+0.5)))
}
func noisePerlin(x, y, scale float32, seed int32) float32 {
	// Gradients: unit vectors in 8 directions
	gradients := [8][2]float32{
		{1, 0}, {-1, 0}, {0, 1}, {0, -1},
		{1, 1}, {-1, 1}, {1, -1}, {-1, -1},
	}

	// Fade function (smootherstep)
	fade := func(t float32) float32 {
		return t * t * t * (t*(t*6-15) + 10)
	}

	// Linear interpolation
	lerp := func(a, b, t float32) float32 {
		return a + t*(b-a)
	}

	// Dot product between gradient and distance
	dot := func(ix, iy int32, x, y float32) float32 {
		hash := Seed([]int32{seed, ix, iy})
		g := gradients[hash%8]
		dx := x - float32(ix)
		dy := y - float32(iy)
		return dx*g[0] + dy*g[1]
	}

	// Scale input coordinates
	x *= scale
	y *= scale

	x0 := int32(math.Floor(float64(x)))
	y0 := int32(math.Floor(float64(y)))
	x1 := x0 + 1
	y1 := y0 + 1

	sx := fade(x - float32(x0))
	sy := fade(y - float32(y0))

	// Calculate dot products at the four corners
	n00 := dot(x0, y0, x, y)
	n10 := dot(x1, y0, x, y)
	n01 := dot(x0, y1, x, y)
	n11 := dot(x1, y1, x, y)

	// Interpolate the results
	ix0 := lerp(n00, n10, sx)
	ix1 := lerp(n01, n11, sx)
	value := lerp(ix0, ix1, sy)

	// Normalize from [-1, 1] to [0, 1]
	return (value + 1) / 2
}
func noiseWorley(x, y, scale float32, seed int32) float32 {
	// Scale input
	x *= scale
	y *= scale

	xi := int32(math.Floor(float64(x)))
	yi := int32(math.Floor(float64(y)))

	minDist := float32(math.MaxFloat32)

	// Check neighboring cells (3x3 grid)
	for dy := int32(-1); dy <= 1; dy++ {
		for dx := int32(-1); dx <= 1; dx++ {
			ix := xi + dx
			iy := yi + dy

			cellSeed := Seed([]int32{seed, ix, iy})
			fx := Random(0, 1, float32(cellSeed))
			fy := Random(0, 1, float32(cellSeed+1))

			cx := float32(ix) + fx
			cy := float32(iy) + fy

			dx := cx - x
			dy := cy - y
			dist := dx*dx + dy*dy // squared distance

			if dist < minDist {
				minDist = dist
			}
		}
	}

	// Normalize distance: closer = lower, scale to [0, 1]
	// sqrt(2) is max possible distance from center to feature point
	return float32(math.Min(1, math.Max(0, float64(math.Sqrt(float64(minDist))/math.Sqrt2))))
}
func noiseVoronoi(x, y, scale float32, seed int32) float32 {
	// Scale coordinates
	x *= scale
	y *= scale

	xi := int32(math.Floor(float64(x)))
	yi := int32(math.Floor(float64(y)))

	minDist := float32(math.MaxFloat32)
	var closestFeature [2]int32

	// Search 3x3 neighborhood
	for dy := int32(-1); dy <= 1; dy++ {
		for dx := int32(-1); dx <= 1; dx++ {
			ix := xi + dx
			iy := yi + dy

			cellSeed := Seed([]int32{seed, ix, iy})
			fx := Random(0, 1, float32(cellSeed))
			fy := Random(0, 1, float32(cellSeed+1))

			cx := float32(ix) + fx
			cy := float32(iy) + fy

			dx := cx - x
			dy := cy - y
			dist := dx*dx + dy*dy

			if dist < minDist {
				minDist = dist
				closestFeature = [2]int32{ix, iy}
			}
		}
	}

	// Create a stable ID for the region based on its grid coords
	regionID := Seed([]int32{seed, closestFeature[0], closestFeature[1]})
	regionVal := Random(0, 1, float32(regionID))

	return regionVal
}
func noiseValue(x, y, scale float32, seed int32) float32 {
	// Linear interpolation
	lerp := func(a, b, t float32) float32 {
		return a + t*(b-a)
	}
	// Smoothstep easing
	smoothstep := func(t float32) float32 {
		return t * t * (3 - 2*t)
	}

	// Scale input coords
	x *= scale
	y *= scale

	xi := int32(math.Floor(float64(x)))
	yi := int32(math.Floor(float64(y)))

	xf := x - float32(xi)
	yf := y - float32(yi)

	// Generate seeds for corners using your Seed function
	seed00 := Seed([]int32{seed, xi, yi})
	seed10 := Seed([]int32{seed, xi + 1, yi})
	seed01 := Seed([]int32{seed, xi, yi + 1})
	seed11 := Seed([]int32{seed, xi + 1, yi + 1})

	// Get random values from your Random function
	v00 := Random(0, 1, float32(seed00))
	v10 := Random(0, 1, float32(seed10))
	v01 := Random(0, 1, float32(seed01))
	v11 := Random(0, 1, float32(seed11))

	// Smooth interpolation weights
	u := smoothstep(xf)
	v := smoothstep(yf)

	// Bilinear interpolate
	i1 := lerp(v00, v10, u)
	i2 := lerp(v01, v11, u)

	return lerp(i1, i2, v)
}
func noiseValueCubic(x, y, scale float32, seed int32) float32 {
	// Cubic interpolation helper
	cubicInterpolate := func(v0, v1, v2, v3, t float32) float32 {
		P := (v3 - v2) - (v0 - v1)
		Q := (v0 - v1) - P
		R := v2 - v0
		S := v1
		return P*t*t*t + Q*t*t + R*t + S
	}

	// Scale input
	x *= scale
	y *= scale

	xi := int32(math.Floor(float64(x)))
	yi := int32(math.Floor(float64(y)))

	xf := x - float32(xi)
	yf := y - float32(yi)

	// Sample 4x4 grid of values
	var vals [4][4]float32
	for gy := -1; gy <= 2; gy++ {
		for gx := -1; gx <= 2; gx++ {
			cellSeed := Seed([]int32{seed, xi + int32(gx), yi + int32(gy)})
			vals[gx+1][gy+1] = Random(0, 1, float32(cellSeed))
		}
	}

	// Interpolate along x for each row
	var interpRow [4]float32
	for i := 0; i < 4; i++ {
		interpRow[i] = cubicInterpolate(
			vals[0][i], vals[1][i], vals[2][i], vals[3][i], xf)
	}

	// Interpolate the results along y
	result := cubicInterpolate(interpRow[0], interpRow[1], interpRow[2], interpRow[3], yf)

	// Clamp result to [0,1] due to interpolation overshoot
	if result < 0 {
		return 0
	}
	if result > 1 {
		return 1
	}
	return result
}
