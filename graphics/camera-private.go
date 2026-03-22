package graphics

import (
	"image/color"
	"pure-game-kit/internal"
	col "pure-game-kit/utility/color"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/number"
	"pure-game-kit/window"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//=================================================================
// objects

var defaultTextPack = &symbol{Color: palette.White, Weight: 1, OutlineColor: palette.Black, OutlineWeight: 3}

//=================================================================
// primitives

func triangulate(points []float32) []float32 {
	var n = len(points) / 2
	if n < 3 {
		return nil
	}

	var triangles []float32
	var verts = make([]int, n)
	for i := range n {
		verts[i] = i
	}

	ccw := area(points) > 0
	for len(verts) > 3 {
		var earFound = false
		for i := 0; i < len(verts); i++ {
			var prev = verts[(i+len(verts)-1)%len(verts)]
			var curr = verts[i]
			var next = verts[(i+1)%len(verts)]

			if !isEar(points, verts, prev, curr, next, ccw) {
				continue
			}

			triangles = append(triangles,
				points[prev*2], points[prev*2+1],
				points[curr*2], points[curr*2+1],
				points[next*2], points[next*2+1],
			)
			verts = append(verts[:i], verts[i+1:]...)
			earFound = true
			break
		}
		if !earFound {
			break
		}
	}

	if len(verts) == 3 {
		for _, v := range verts {
			triangles = append(triangles, points[v*2], points[v*2+1])
		}
	}

	return triangles
}
func area(points []float32) float32 {
	var a float32
	var n = len(points) / 2
	for i := range n {
		var j = (i + 1) % n
		a += points[i*2]*points[j*2+1] - points[j*2]*points[i*2+1]
	}
	return a / 2
}
func isEar(points []float32, verts []int, i1, i2, i3 int, ccw bool) bool {
	var p1x, p1y = points[i1*2], points[i1*2+1]
	var p2x, p2y = points[i2*2], points[i2*2+1]
	var p3x, p3y = points[i3*2], points[i3*2+1]
	var cross = (p2x-p1x)*(p3y-p1y) - (p2y-p1y)*(p3x-p1x)
	if ccw && cross <= 0 {
		return false
	}
	if !ccw && cross >= 0 {
		return false
	}

	for _, vi := range verts {
		if vi == i1 || vi == i2 || vi == i3 {
			continue
		}
		if pointInTriangle(points[vi*2], points[vi*2+1], p1x, p1y, p2x, p2y, p3x, p3y) {
			return false
		}
	}
	return true
}
func pointInTriangle(px, py, ax, ay, bx, by, cx, cy float32) bool {
	var v0x, v0y = cx - ax, cy - ay
	var v1x, v1y = bx - ax, by - ay
	var v2x, v2y = px - ax, py - ay
	var dot00 = v0x*v0x + v0y*v0y
	var dot01 = v0x*v1x + v0y*v1y
	var dot02 = v0x*v2x + v0y*v2y
	var dot11 = v1x*v1x + v1y*v1y
	var dot12 = v1x*v2x + v1y*v2y
	var invDenom = 1 / (dot00*dot11 - dot01*dot01)
	var u = (dot11*dot02 - dot01*dot12) * invDenom
	var v = (dot00*dot12 - dot01*dot02) * invDenom
	return (u >= 0) && (v >= 0) && (u+v <= 1)
}
func isClockwiseFlat(tri []float32) bool {
	return (tri[2]-tri[0])*(tri[5]-tri[1])-(tri[4]-tri[0])*(tri[3]-tri[1]) < 0
}
func separateShapes(points [][2]float32) (flatPoints []float32, ptsCountsPerShape []int) {
	var currentCount int
	for _, p := range points {
		if number.IsNaN(p[0]) || number.IsNaN(p[1]) {
			if currentCount > 0 {
				ptsCountsPerShape = append(ptsCountsPerShape, currentCount)
				currentCount = 0
			}
			continue
		}
		flatPoints = append(flatPoints, p[0], p[1])
		currentCount++
	}
	if currentCount > 0 {
		ptsCountsPerShape = append(ptsCountsPerShape, currentCount)
	}
	return flatPoints, ptsCountsPerShape
}
func isConvex(pts []float32, count int) bool {
	if count < 3 {
		return true // a point or a line block (1 or 2 points) is trivially convex
	}

	var gotPositive, gotNegative bool

	for i := range count {
		var p0, p1, p2 = i * 2, ((i + 1) % count) * 2, ((i + 2) % count) * 2
		var dx1, dy1 = pts[p1] - pts[p0], pts[p1+1] - pts[p0+1]
		var dx2, dy2 = pts[p2] - pts[p1], pts[p2+1] - pts[p1+1]
		var crossProduct = dx1*dy2 - dy1*dx2

		if crossProduct > 0 {
			gotPositive = true
		} else if crossProduct < 0 {
			gotNegative = true
		}

		if gotPositive && gotNegative {
			return false // found both left & right turn - definitively concave
		}
	}

	return true
}

//=================================================================
// camera

var rlCam = rl.Camera2D{}

var debugStr string

const placeholderCharAsset = '@'

func (c *Camera) area() (x, y, w, h float32) {
	if c.Area == nil {
		var ww, wh = window.Size()
		return 0, 0, float32(ww), float32(wh)
	}
	return c.Area.X, c.Area.Y, c.Area.Width, c.Area.Height
}
func (c *Camera) mask() (x, y, w, h float32) {
	if c.Mask == nil {
		var ww, wh = window.Size()
		return 0, 0, float32(ww), float32(wh)
	}
	return c.Mask.X, c.Mask.Y, c.Mask.Width, c.Mask.Height
}

// call before draw to update camera but use screen space instead of camera space
func (c *Camera) update() {
	tryRecreateWindow()

	var sx, sy, sw, sh = c.area()
	rlCam.Target.X = float32(c.X)
	rlCam.Target.Y = float32(c.Y)
	rlCam.Rotation = float32(c.Angle)
	rlCam.Zoom = float32(c.Zoom)
	rlCam.Offset.X = sx + sw/2
	rlCam.Offset.Y = sy + sh/2
}

// call before draw to update camera and use camera space
func (c *Camera) begin() {
	c.update()
	if skipStartEnd {
		return
	}

	rl.BeginMode2D(rlCam)
	mask = c.Mask

	if c.Area != nil {
		var mx, my, mw, mh = c.area()
		rl.BeginScissorMode(int32(mx), int32(my), int32(mw), int32(mh))
	}

	rl.EnableDepthTest()
	c.Effects.updateUniforms(1, 1, nil, nil, true)

	if c.Blend != 0 {
		rl.BeginBlendMode(rl.BlendMode(c.Blend))
	}
}

// call after draw to get back to using screen space
func (c *Camera) end() {
	if skipStartEnd {
		return
	}

	if c.Blend != 0 {
		rl.EndBlendMode()
	}

	rl.DisableDepthTest()
	if c.Area != nil {
		rl.EndScissorMode()
	}
	rl.EndMode2D()
}

//=================================================================
// other

func tryRecreateWindow() {
	if internal.WindowReady {
		return
	}

	if !rl.IsWindowReady() {
		window.Recreate()
	}
}

func getColor(value uint) color.RGBA {
	var r, g, b, a = col.Channels(value)
	return color.RGBA{R: r, G: g, B: b, A: a}
}
func packSymbolColor(s *symbol) rl.Color {
	var packLayer = func(c rl.Color) uint8 {
		var r = (c.R >> 6) & 0x03
		var g = (c.G >> 6) & 0x03
		var b = (c.B >> 6) & 0x03
		var a = (c.A >> 6) & 0x03
		return (r << 6) | (g << 4) | (b << 2) | a
	}

	var thick, out, sh, shSmooth byte = s.Weight, s.OutlineWeight, s.ShadowWeight, s.ShadowBlur
	var r = packLayer(getColor(s.Color))
	var g = packLayer(getColor(s.OutlineColor))
	var b = packLayer(getColor(s.ShadowColor))
	var a = ((thick & 0x03) << 6) | ((out & 0x03) << 4) | ((sh & 0x03) << 2) | (shSmooth & 0x03)
	return rl.NewColor(r, g, b, a)
}
