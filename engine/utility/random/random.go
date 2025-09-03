package random

import (
	"math"
	"time"
)

func Seed(keys ...float32) float32 {
	var ints = make([]int, len(keys))
	for i := range ints {
		ints[i] = floatToIntSeed(keys[i])
	}
	return float32(SeedInts(ints...))
}
func SeedInts(keys ...int) int {
	var hashSeed = func(seed uint64, a int) uint64 {
		seed ^= uint64(a)
		seed = (seed ^ (seed >> 16)) * 2246822519
		seed = (seed ^ (seed >> 13)) * 3266489917
		seed ^= seed >> 16
		return seed
	}
	var seed = uint64(2654435769)

	for _, p := range keys {
		seed = hashSeed(seed, p)
	}
	return int(seed)
}

func Range(a, b, seed float32) float32 {
	if a == b {
		return a
	}
	if a > b {
		a, b = b, a
	}

	if seed != seed { // seed is NaN
		seed = float32(time.Now().UnixNano()%1e9) / 1e9 // value in [0,1)
	}

	var diff = b - a
	var intSeed = int(seed * 2147483647)
	intSeed = (1103515245*intSeed + 12345) % 2147483647
	var normalized = float32(intSeed) / 2147483647.0
	return a + normalized*diff
}
func RangeInt(a, b int, seed float32) int {
	return int(Range(float32(a), float32(b), seed))
}

func HasChance(percent float32) bool {
	return HasChanceSeeded(percent, float32(math.NaN()))
}
func HasChanceSeeded(percent, seed float32) bool {
	if percent <= 0 {
		return false
	}
	percent = float32(math.Min(100, float64(percent)))
	var n = Range(0, 100, seed)
	return n <= percent
}

func Shuffle[T any](items ...T) {
	ShuffleSeeded(float32(math.NaN()), items...)
}
func ShuffleSeeded[T any](seed float32, items ...T) {
	for i := len(items) - 1; i > 0; i-- {
		var j = RangeInt(0, i, seed)
		items[i], items[j] = items[j], items[i]
	}
}

func ChooseMultiple[T any](count int, items ...T) []T {
	return chooseMultipleInternal(items, count, float32(math.NaN()))
}
func ChooseMultipleSeeded[T any](count int, seed float32, items ...T) []T {
	return chooseMultipleInternal(items, count, seed)
}

func ChooseOne[T any](items ...T) *T {
	return singlePointer(chooseMultipleInternal(items, 1, float32(math.NaN())))
}
func ChooseOneSeeded[T any](seed float32, items ...T) *T {
	return singlePointer(chooseMultipleInternal(items, 1, seed))
}

func NoisePerlin(x, y, scale, seed float32) float32 {
	var intSeed = floatToIntSeed(seed)
	var gradients = [8][2]float32{
		{1, 0}, {-1, 0}, {0, 1}, {0, -1},
		{1, 1}, {-1, 1}, {1, -1}, {-1, -1},
	}
	var fade = func(t float32) float32 { return t * t * t * (t*(t*6-15) + 10) }
	var lerp = func(a, b, t float32) float32 { return a + t*(b-a) }
	var dot = func(ix, iy int, x, y float32) float32 {
		var hash = SeedInts(intSeed, ix, iy)
		var g = gradients[hash%8]
		var dx = x - float32(ix)
		var dy = y - float32(iy)
		return dx*g[0] + dy*g[1]
	}

	x *= scale
	y *= scale

	var x0 = int(math.Floor(float64(x)))
	var y0 = int(math.Floor(float64(y)))
	var x1 = x0 + 1
	var y1 = y0 + 1
	var sx = fade(x - float32(x0))
	var sy = fade(y - float32(y0))
	var n00 = dot(x0, y0, x, y)
	var n10 = dot(x1, y0, x, y)
	var n01 = dot(x0, y1, x, y)
	var n11 = dot(x1, y1, x, y)
	var ix0 = lerp(n00, n10, sx)
	var ix1 = lerp(n01, n11, sx)
	var value = lerp(ix0, ix1, sy)

	return (value + 1) / 2 // normalized from [-1, 1] to [0, 1]
}
func NoiseOpenSimplex(x, y, scale, seed float32) float32 {
	var intSeed = floatToIntSeed(seed)
	const stretch2D = -0.211324865405187 // (1/sqrt(2+1)-1)/2
	const squish2D = 0.366025403784439   // (sqrt(2+1)-1)/2
	const norm2D = 47.0

	var gradients = [8][2]float32{
		{5, 2}, {2, 5}, {-5, 2}, {-2, 5},
		{5, -2}, {2, -5}, {-5, -2}, {-2, -5},
	}
	var getPerm = func(i, j int) int {
		var subSeed = SeedInts(intSeed, i, j)
		var r = Range(0, 255, float32(subSeed))
		return int(r) & 255
	}

	x *= scale
	y *= scale

	var s = (x + y) * float32(stretch2D)
	var xf = x + s
	var yf = y + s
	var xi = int(math.Floor(float64(xf)))
	var yi = int(math.Floor(float64(yf)))
	var t = float32(xi+yi) * float32(squish2D)
	var X0 = float32(xi) - t
	var Y0 = float32(yi) - t
	var x0 = x - X0
	var y0 = y - Y0
	var x1, y1 float32
	var i1, j1 int

	if x0 > y0 {
		i1, j1 = 1, 0
	} else {
		i1, j1 = 0, 1
	}

	x1 = x0 - float32(i1) + float32(squish2D)
	y1 = y0 - float32(j1) + float32(squish2D)
	var x2 = x0 - 1 + 2*float32(squish2D)
	var y2 = y0 - 1 + 2*float32(squish2D)

	var contrib = func(dx, dy float32, i, j int) float32 {
		var attsq = 2 - dx*dx - dy*dy
		if attsq > 0 {
			var pi = getPerm(i, j) % 8
			var grad = gradients[pi]
			var att4 = attsq * attsq
			return att4 * att4 * (grad[0]*dx + grad[1]*dy)
		}
		return 0
	}

	var value = contrib(x0, y0, xi, yi)
	value += contrib(x1, y1, xi+i1, yi+j1)
	value += contrib(x2, y2, xi+1, yi+1)

	return float32(math.Max(0, math.Min(1, float64(value)/norm2D+0.5))) // normalized to [0, 1]
}
func NoiseWorley(x, y, scale, seed float32) float32 {
	var instSeed = floatToIntSeed(seed)
	x *= scale
	y *= scale

	var xi = int(math.Floor(float64(x)))
	var yi = int(math.Floor(float64(y)))
	var minDist = float32(math.MaxFloat32)

	for dy := int(-1); dy <= 1; dy++ {
		for dx := int(-1); dx <= 1; dx++ {
			var ix = xi + dx
			var iy = yi + dy
			var cellSeed = SeedInts(instSeed, ix, iy)
			var fx = Range(0, 1, float32(cellSeed))
			var fy = Range(0, 1, float32(cellSeed+1))
			var cx = float32(ix) + fx
			var cy = float32(iy) + fy
			var dx = cx - x
			var dy = cy - y
			var dist = dx*dx + dy*dy

			if dist < minDist {
				minDist = dist
			}
		}
	}

	// Normalize distance: closer = lower, scale to [0, 1]
	// sqrt(2) is max possible distance from center to feature point
	return float32(math.Min(1, math.Max(0, float64(math.Sqrt(float64(minDist))/math.Sqrt2))))
}
func NoiseVoronoi(x, y, scale, seed float32) float32 {
	var intSeed = floatToIntSeed(seed)
	x *= scale
	y *= scale

	var xi = int(math.Floor(float64(x)))
	var yi = int(math.Floor(float64(y)))
	var minDist = float32(math.MaxFloat32)
	var closestFeature [2]int

	for dy := int(-1); dy <= 1; dy++ {
		for dx := int(-1); dx <= 1; dx++ {
			var ix = xi + dx
			var iy = yi + dy
			var cellSeed = SeedInts(intSeed, ix, iy)
			var fx = Range(0, 1, float32(cellSeed))
			var fy = Range(0, 1, float32(cellSeed+1))
			var cx = float32(ix) + fx
			var cy = float32(iy) + fy
			var dx = cx - x
			var dy = cy - y
			var dist = dx*dx + dy*dy

			if dist < minDist {
				minDist = dist
				closestFeature = [2]int{ix, iy}
			}
		}
	}

	var regionID = SeedInts(floatToIntSeed(seed), closestFeature[0], closestFeature[1])
	var regionVal = Range(0, 1, float32(regionID))

	return regionVal
}
func NoiseValue(x, y, scale, seed float32) float32 {
	var intSeed = floatToIntSeed(seed)
	var lerp = func(a, b, t float32) float32 { return a + t*(b-a) }
	var smoothstep = func(t float32) float32 { return t * t * (3 - 2*t) }

	x *= scale
	y *= scale

	var xi = int(math.Floor(float64(x)))
	var yi = int(math.Floor(float64(y)))
	var xf = x - float32(xi)
	var yf = y - float32(yi)
	var seed00 = SeedInts(intSeed, xi, yi)
	var seed10 = SeedInts(intSeed, xi+1, yi)
	var seed01 = SeedInts(intSeed, xi, yi+1)
	var seed11 = SeedInts(intSeed, xi+1, yi+1)
	var v00 = Range(0, 1, float32(seed00))
	var v10 = Range(0, 1, float32(seed10))
	var v01 = Range(0, 1, float32(seed01))
	var v11 = Range(0, 1, float32(seed11))
	var u = smoothstep(xf)
	var v = smoothstep(yf)
	var i1 = lerp(v00, v10, u)
	var i2 = lerp(v01, v11, u)

	return lerp(i1, i2, v)
}
func NoiseValueCubic(x, y, scale, seed float32) float32 {
	var intSeed = floatToIntSeed(seed)
	var cubicInterpolate = func(v0, v1, v2, v3, t float32) float32 {
		var P = (v3 - v2) - (v0 - v1)
		var Q = (v0 - v1) - P
		var R = v2 - v0
		var S = v1
		return P*t*t*t + Q*t*t + R*t + S
	}

	x *= scale
	y *= scale

	var xi = int(math.Floor(float64(x)))
	var yi = int(math.Floor(float64(y)))
	var xf = x - float32(xi)
	var yf = y - float32(yi)
	var vals [4][4]float32
	var interpRow [4]float32

	for gy := -1; gy <= 2; gy++ {
		for gx := -1; gx <= 2; gx++ {
			var cellSeed = SeedInts(intSeed, xi+int(gx), yi+int(gy))
			vals[gx+1][gy+1] = Range(0, 1, float32(cellSeed))
		}
	}

	for i := range 4 {
		interpRow[i] = cubicInterpolate(
			vals[0][i], vals[1][i], vals[2][i], vals[3][i], xf)
	}

	var result = cubicInterpolate(interpRow[0], interpRow[1], interpRow[2], interpRow[3], yf)

	if result < 0 {
		return 0
	}
	if result > 1 {
		return 1
	}
	return result
}

// region private

func floatToIntSeed(seed float32) int {
	var intSeed int = 0
	if math.IsNaN(float64(seed)) {
		intSeed = int(time.Now().UnixNano())
	} else {
		intSeed = int(math.Float32bits(seed))
	}
	return intSeed
}

func chooseMultipleInternal[T any](items []T, count int, seed float32) []T {
	if len(items) == 0 || count <= 0 {
		return []T{}
	}

	var clone = make([]T, len(items))
	copy(clone, items)

	for i := len(clone) - 1; i > 0; i-- {
		var j = RangeInt(0, i, seed)
		clone[i], clone[j] = clone[j], clone[i]
	}

	if count > len(clone) {
		count = len(clone)
	}

	return clone[:count]
}

func singlePointer[T any](items []T) *T {
	if len(items) == 0 {
		return nil
	}
	return &items[0]
}

// endregion
