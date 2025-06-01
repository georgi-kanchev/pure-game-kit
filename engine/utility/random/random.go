package random

import (
	"math"
	"time"
)

func Seed(keys []float32) float32 {
	var ints = make([]int32, len(keys))
	for i := range ints {
		ints[i] = floatToIntSeed(keys[i])
	}
	return float32(SeedInts(ints))
}
func SeedInts(keys []int32) int32 {
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

func HasChance(percent, seed float32) bool {
	if percent <= 0 {
		return false
	}
	percent = float32(math.Min(100, float64(percent)))
	n := Random(0, 100, seed)
	return n <= percent
}

func NoisePerlin(x, y, scale, seed float32) float32 {
	var intSeed = floatToIntSeed(seed)
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
		hash := SeedInts([]int32{intSeed, ix, iy})
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
func NoiseOpenSimplex(x, y, scale, seed float32) float32 {
	var intSeed = floatToIntSeed(seed)
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
		subSeed := SeedInts([]int32{intSeed, i, j})
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
func NoiseWorley(x, y, scale, seed float32) float32 {
	var instSeed = floatToIntSeed(seed)
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

			cellSeed := SeedInts([]int32{instSeed, ix, iy})
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
func NoiseVoronoi(x, y, scale, seed float32) float32 {
	var intSeed = floatToIntSeed(seed)
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

			cellSeed := SeedInts([]int32{intSeed, ix, iy})
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
	regionID := SeedInts([]int32{floatToIntSeed(seed), closestFeature[0], closestFeature[1]})
	regionVal := Random(0, 1, float32(regionID))

	return regionVal
}
func NoiseValue(x, y, scale, seed float32) float32 {
	var intSeed = floatToIntSeed(seed)
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
	seed00 := SeedInts([]int32{intSeed, xi, yi})
	seed10 := SeedInts([]int32{intSeed, xi + 1, yi})
	seed01 := SeedInts([]int32{intSeed, xi, yi + 1})
	seed11 := SeedInts([]int32{intSeed, xi + 1, yi + 1})

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
func NoiseValueCubic(x, y, scale, seed float32) float32 {
	var intSeed = floatToIntSeed(seed)
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
			cellSeed := SeedInts([]int32{intSeed, xi + int32(gx), yi + int32(gy)})
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

// region private

func floatToIntSeed(seed float32) int32 {
	var intSeed int32 = 0
	if math.IsNaN(float64(seed)) {
		intSeed = int32(time.Now().UnixNano())
	} else {
		intSeed = int32(math.Float32bits(seed))
	}

	return intSeed
}

// endregion
