/*
Provides multiple 2D noise functions, each with a different behavior. They accept an x & y coordinate
to sample from and a scale for the noise itself to zoom in/out. They also accept multiple optional seeds,
making it possible to generate the same noise value when providing the same seeds in the same order
(usually loop indexes or grid coordinates etc). After computation, they return a 0 to 1 ranged value.
The usege is very similar to indexing a pixel color on a coordinate of an infinite gray scaled image.
Useful for terrain generation, controlled randomness, repeating pattern effects & many more.
*/
package noise

import (
	"math"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/random"
	"time"
)

func Perlin(x, y, scale float32, seeds ...float32) float32 {
	var intSeed = floatToIntSeed(random.CombineSeeds(seeds...))
	var gradients = [8][2]float32{
		{1, 0}, {-1, 0}, {0, 1}, {0, -1},
		{1, 1}, {-1, 1}, {1, -1}, {-1, -1},
	}
	var fade = func(t float32) float32 { return t * t * t * (t*(t*6-15) + 10) }
	var lerp = func(a, b, t float32) float32 { return a + t*(b-a) }
	var dot = func(ix, iy int, x, y float32) float32 {
		var hash = random.CombineSeeds(intSeed, ix, iy)
		var g = gradients[hash%8]
		var dx = x - float32(ix)
		var dy = y - float32(iy)
		return dx*g[0] + dy*g[1]
	}

	x *= scale
	y *= scale

	var x0, y0 = int(number.RoundDown(x)), int(number.RoundDown(y))
	var x1, y1 = x0 + 1, y0 + 1
	var sx, sy = fade(x - float32(x0)), fade(y - float32(y0))
	var n00, n10 = dot(x0, y0, x, y), dot(x1, y0, x, y)
	var n01, n11 = dot(x0, y1, x, y), dot(x1, y1, x, y)
	var ix0, ix1 = lerp(n00, n10, sx), lerp(n01, n11, sx)
	var value = lerp(ix0, ix1, sy)

	return (value + 1) / 2 // normalized from [-1, 1] to [0, 1]
}
func OpenSimplex(x, y, scale float32, seeds ...float32) float32 {
	var intSeed = floatToIntSeed(random.CombineSeeds(seeds...))
	const stretch2D = -0.211324865405187 // (1/sqrt(2+1)-1)/2
	const squish2D = 0.366025403784439   // (sqrt(2+1)-1)/2
	const norm2D = 47.0

	var gradients = [8][2]float32{
		{5, 2}, {2, 5}, {-5, 2}, {-2, 5},
		{5, -2}, {2, -5}, {-5, -2}, {-2, -5},
	}
	var getPerm = func(i, j int) int {
		var subSeed = random.CombineSeeds(intSeed, i, j)
		var r = random.Range(0, 255, float32(subSeed))
		return int(r) & 255
	}

	x *= scale
	y *= scale

	var s = (x + y) * float32(stretch2D)
	var xf, yf = x + s, y + s
	var xi, yi = int(number.RoundDown(xf)), int(number.RoundDown(yf))
	var t = float32(xi+yi) * float32(squish2D)
	var X0, Y0 = float32(xi) - t, float32(yi) - t
	var x0, y0 = x - X0, y - Y0
	var x1, y1 float32
	var i1, j1 int

	if x0 > y0 {
		i1, j1 = 1, 0
	} else {
		i1, j1 = 0, 1
	}

	x1 = x0 - float32(i1) + float32(squish2D)
	y1 = y0 - float32(j1) + float32(squish2D)
	var x2, y2 = x0 - 1 + 2*float32(squish2D), y0 - 1 + 2*float32(squish2D)

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

	return number.Limit(value/norm2D+0.5, 0, 1)
}
func Worley(x, y, scale float32, seeds ...float32) float32 {
	x *= scale
	y *= scale

	var xi, yi = int(number.RoundDown(x)), int(number.RoundDown(y))
	var minDist = number.Infinity()
	var intSeed = floatToIntSeed(random.CombineSeeds(seeds...))

	for dy := int(-1); dy <= 1; dy++ {
		for dx := int(-1); dx <= 1; dx++ {
			var ix, iy = xi + dx, yi + dy
			var cellSeed = random.CombineSeeds(intSeed, ix, iy)
			var fx, fy = random.Range(0.0, 1.0, float32(cellSeed)), random.Range(0.0, 1.0, float32(cellSeed+1))
			var cx, cy = float32(ix) + float32(fx), float32(iy) + float32(fy)
			var dx, dy = cx - x, cy - y
			var dist = dx*dx + dy*dy

			if dist < minDist {
				minDist = dist
			}
		}
	}

	// Normalize distance: closer = lower, scale to [0, 1]
	const sqrt2 = 1.4142 // max possible distance from center to feature point
	var result = float32(number.SquareRoot(minDist) / sqrt2)
	return number.Limit(result, 0, 1)
}
func Voronoi(x, y, scale float32, seeds ...float32) float32 {
	x *= scale
	y *= scale

	var seed = random.CombineSeeds(seeds...)
	var intSeed = floatToIntSeed(seed)
	var xi, yi = int(number.RoundDown(x)), int(number.RoundDown(y))
	var minDist = number.Infinity()
	var closestFeature [2]int

	for dy := int(-1); dy <= 1; dy++ {
		for dx := int(-1); dx <= 1; dx++ {
			var ix, iy = xi + dx, yi + dy
			var cellSeed = random.CombineSeeds(intSeed, ix, iy)
			var fx, fy = random.Range(0.0, 1.0, float32(cellSeed)), random.Range(0.0, 1.0, float32(cellSeed+1))
			var cx, cy = float32(ix) + float32(fx), float32(iy) + float32(fy)
			var dx, dy = cx - x, cy - y
			var dist = dx*dx + dy*dy

			if dist < minDist {
				minDist = dist
				closestFeature = [2]int{ix, iy}
			}
		}
	}

	var regionID = random.CombineSeeds(floatToIntSeed(seed), closestFeature[0], closestFeature[1])
	var regionVal = random.Range(0.0, 1.0, float32(regionID))
	return float32(regionVal)
}
func Value(x, y, scale float32, seeds ...float32) float32 {
	x *= scale
	y *= scale

	var intSeed = floatToIntSeed(random.CombineSeeds(seeds...))
	var lerp = func(a, b, t float32) float32 { return a + t*(b-a) }
	var smoothstep = func(t float32) float32 { return t * t * (3 - 2*t) }
	var xi, yi = int(number.RoundDown(x)), int(number.RoundDown(y))
	var xf, yf = x - float32(xi), y - float32(yi)
	var seed00, seed10 = random.CombineSeeds(intSeed, xi, yi), random.CombineSeeds(intSeed, xi+1, yi)
	var seed01, seed11 = random.CombineSeeds(intSeed, xi, yi+1), random.CombineSeeds(intSeed, xi+1, yi+1)
	var v00, v10 = random.Range(0.0, 1.0, float32(seed00)), random.Range(0.0, 1.0, float32(seed10))
	var v01, v11 = random.Range(0.0, 1.0, float32(seed01)), random.Range(0.0, 1.0, float32(seed11))
	var u, v = smoothstep(xf), smoothstep(yf)
	var i1, i2 = lerp(float32(v00), float32(v10), u), lerp(float32(v01), float32(v11), u)

	return lerp(i1, i2, v)
}
func ValueCubic(x, y, scale float32, seeds ...float32) float32 {
	x *= scale
	y *= scale

	var intSeed = floatToIntSeed(random.CombineSeeds(seeds...))
	var xi, yi = int(number.RoundDown(x)), int(number.RoundDown(y))
	var xf, yf = x - float32(xi), y - float32(yi)
	var vals [4][4]float32
	var interpRow [4]float32
	var cubicInterpolate = func(v0, v1, v2, v3, t float32) float32 {
		var P = (v3 - v2) - (v0 - v1)
		var Q = (v0 - v1) - P
		var R = v2 - v0
		var S = v1
		return P*t*t*t + Q*t*t + R*t + S
	}

	for gy := -1; gy <= 2; gy++ {
		for gx := -1; gx <= 2; gx++ {
			var cellSeed = random.CombineSeeds(intSeed, xi+int(gx), yi+int(gy))
			vals[gx+1][gy+1] = float32(random.Range(0.0, 1.0, float32(cellSeed)))
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

//=================================================================
// private

func floatToIntSeed(seed float32) int {
	var intSeed int = 0
	if seed != seed { // is NaN
		intSeed = int(time.Now().UnixNano())
	} else {
		intSeed = int(math.Float32bits(seed))
	}
	return intSeed
}
