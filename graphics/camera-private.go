package graphics

import (
	"image/color"
	"pure-game-kit/internal"
	col "pure-game-kit/utility/color"
	"pure-game-kit/utility/number"
	"pure-game-kit/window"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//=================================================================
// objects

var reusableSprite = NewSprite("", 0, 0)

func drawBoxPart(camera *Camera, parent *Node, x, y, w, h float32, id string, color uint) {
	reusableSprite.AssetId = id
	reusableSprite.X, reusableSprite.Y = x, y
	reusableSprite.Parent = parent
	reusableSprite.Width, reusableSprite.Height = w, h
	reusableSprite.ScaleX, reusableSprite.ScaleY = 1, 1
	reusableSprite.Tint = color

	camera.DrawSprites(reusableSprite)
}

//=================================================================
// primitives

func triangulate(points []float32) []float32 {
	n := len(points) / 2
	if n < 3 {
		return nil
	}

	var triangles []float32
	var verts = make([]int, n)
	for i := 0; i < n; i++ {
		verts[i] = i
	}

	ccw := area(points) > 0
	for len(verts) > 3 {
		earFound := false
		for i := 0; i < len(verts); i++ {
			prev := verts[(i+len(verts)-1)%len(verts)]
			curr := verts[i]
			next := verts[(i+1)%len(verts)]

			if !isEar(points, verts, prev, curr, next, ccw) {
				continue
			}

			// Add triangle vertices to flat slice
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
	n := len(points) / 2
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		a += points[i*2]*points[j*2+1] - points[j*2]*points[i*2+1]
	}
	return a / 2
}
func isEar(points []float32, verts []int, i1, i2, i3 int, ccw bool) bool {
	p1x, p1y := points[i1*2], points[i1*2+1]
	p2x, p2y := points[i2*2], points[i2*2+1]
	p3x, p3y := points[i3*2], points[i3*2+1]

	cross := (p2x-p1x)*(p3y-p1y) - (p2y-p1y)*(p3x-p1x)
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
	v0x, v0y := cx-ax, cy-ay
	v1x, v1y := bx-ax, by-ay
	v2x, v2y := px-ax, py-ay

	dot00 := v0x*v0x + v0y*v0y
	dot01 := v0x*v1x + v0y*v1y
	dot02 := v0x*v2x + v0y*v2y
	dot11 := v1x*v1x + v1y*v1y
	dot12 := v1x*v2x + v1y*v2y

	invDenom := 1 / (dot00*dot11 - dot01*dot01)
	u := (dot11*dot02 - dot01*dot12) * invDenom
	v := (dot00*dot12 - dot01*dot02) * invDenom

	return (u >= 0) && (v >= 0) && (u+v <= 1)
}
func isClockwiseFlat(tri []float32) bool {
	area := (tri[2]-tri[0])*(tri[5]-tri[1]) - (tri[4]-tri[0])*(tri[3]-tri[1])
	return area < 0
}

func separateShapes(points [][2]float32) (flatPoints []float32, shapeCounts []int) {
	var currentCount int
	for _, p := range points {
		if number.IsNaN(p[0]) || number.IsNaN(p[1]) {
			if currentCount >= 3 {
				shapeCounts = append(shapeCounts, currentCount)
			}
			currentCount = 0
			continue
		}
		flatPoints = append(flatPoints, p[0], p[1])
		currentCount++
	}
	if currentCount >= 3 {
		shapeCounts = append(shapeCounts, currentCount)
	}
	return flatPoints, shapeCounts
}

//=================================================================
// camera

var rlCam = rl.Camera2D{}
var maskX, maskY, maskW, maskH int

// call before draw to update camera but use screen space instead of camera space
func (c *Camera) update() {
	tryRecreateWindow()

	rlCam.Target.X = float32(c.X)
	rlCam.Target.Y = float32(c.Y)
	rlCam.Rotation = float32(c.Angle)
	rlCam.Zoom = float32(c.Zoom)
	rlCam.Offset.X = float32(c.ScreenX) + float32(c.ScreenWidth/2)
	rlCam.Offset.Y = float32(c.ScreenY) + float32(c.ScreenHeight/2)

	var mx = number.Biggest(c.MaskX, c.ScreenX)
	var my = number.Biggest(c.MaskY, c.ScreenY)
	var maxW = c.ScreenX + c.ScreenWidth - mx
	var maxH = c.ScreenY + c.ScreenHeight - my
	var mw = number.Smallest(c.MaskWidth, maxW)
	var mh = number.Smallest(c.MaskHeight, maxH)

	maskX, maskY, maskW, maskH = mx, my, mw, mh
}

// call before draw to update camera and use camera space
func (c *Camera) begin() {
	c.update()
	if skipStartEnd {
		return
	}

	rl.BeginMode2D(rlCam)
	rl.BeginScissorMode(int32(maskX), int32(maskY), int32(maskW), int32(maskH))

	if c.Blend != 0 {
		rl.BeginBlendMode(rl.BlendMode(c.Blend))
	}
	if c.Effects != nil {
		rl.BeginShaderMode(internal.Shader)
		rl.EnableDepthTest()
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
	if c.Effects != nil {
		rl.DisableDepthTest()
		rl.EndShaderMode()
	}

	rl.EndScissorMode()
	rl.EndMode2D()
}

func (c *Camera) isAreaVisible(x, y, width, height, angle float32) bool {
	c.update()
	// optimized for speed
	var sinA, cosA = internal.SinCos(angle)
	var sinB, cosB = internal.SinCos(angle + 90)
	var tlx, tly = x, y
	var trx, try = tlx + cosA*width, tly + sinA*width
	var blx, bly = tlx + cosB*height, tly + sinB*height
	var brx, bry = trx + cosB*height, try + sinB*height
	var tx, ty = float32(rlCam.Target.X), float32(rlCam.Target.Y)
	var zoom = float32(rlCam.Zoom)
	var sinR, cosR = internal.SinCos(rlCam.Rotation)
	var offX, offY = float32(rlCam.Offset.X), float32(rlCam.Offset.Y)
	var pointToScreen = func(px, py float32) (float32, float32) { // inlined to skip cam.update() on each call
		px -= tx
		py -= ty
		px *= zoom
		py *= zoom
		var rx, ry = px*cosR - py*sinR, px*sinR + py*cosR
		return rx + offX, ry + offY
	}
	var stlx, stly = pointToScreen(tlx, tly)
	var strx, stry = pointToScreen(trx, try)
	var sbrx, sbry = pointToScreen(brx, bry)
	var sblx, sbly = pointToScreen(blx, bly)
	var minX = number.Smallest(stlx, strx, sbrx, sblx)
	var maxX = number.Biggest(stlx, strx, sbrx, sblx)
	var minY = number.Smallest(stly, stry, sbry, sbly)
	var maxY = number.Biggest(stly, stry, sbry, sbly)
	var mtlx, mtly = float32(c.MaskX), float32(c.MaskY)
	var mbrx, mbry = mtlx + float32(c.MaskWidth), mtly + float32(c.MaskHeight)
	return maxY > mtly && minY < mbry && maxX > mtlx && minX < mbrx
}

//=================================================================
// other

func getColor(value uint) color.RGBA {
	var r, g, b, a = col.Channels(value)
	return color.RGBA{R: r, G: g, B: b, A: a}
}

func tryRecreateWindow() {
	if internal.WindowReady {
		return
	}

	if !rl.IsWindowReady() {
		window.Recreate()
	}
}
