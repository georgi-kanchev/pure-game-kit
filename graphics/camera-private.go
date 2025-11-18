package graphics

import (
	"pure-game-kit/internal"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/point"
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
	reusableSprite.Color = color

	camera.DrawSprites(reusableSprite)
}

func beginShader(t *TextBox, thick float32) {
	var sh = internal.ShaderText

	if sh.ID != 0 {
		var smoothness = []float32{t.Smoothness * t.LineHeight / 5}
		rl.BeginShaderMode(sh)
		rl.SetShaderValue(sh, rl.GetShaderLocation(sh, "smoothness"), smoothness, rl.ShaderUniformFloat)
		setShaderThick(thick)
	}
}
func setShaderThick(thick float32) {
	var sh = internal.ShaderText

	if sh.ID != 0 {
		var thickness = []float32{thick}
		thickness[0] = number.Limit(thickness[0], 0, 0.999)
		rl.SetShaderValue(sh, rl.GetShaderLocation(sh, "thickness"), thickness, rl.ShaderUniformFloat)
	}
}
func endShader() {
	var sh = internal.ShaderText

	if sh.ID != 0 {
		rl.EndShaderMode()
	}
}

//=================================================================
// primitives

func triangulate(points [][2]float32) [][3][2]float32 {
	var n = len(points)
	if n < 3 {
		return nil
	}

	var triangles [][3][2]float32
	var verts = make([]int, n)
	for i := range n {
		verts[i] = i
	}

	var ccw = area(points) > 0
	for len(verts) > 3 {
		var earFound = false

		for i := 0; i < len(verts); i++ {
			var prev = verts[(i+len(verts)-1)%len(verts)]
			var curr = verts[i]
			var next = verts[(i+1)%len(verts)]
			var p1 = points[prev]
			var p2 = points[curr]
			var p3 = points[next]

			if !isEar(points, verts, prev, curr, next, ccw) {
				continue
			}

			triangles = append(triangles, [3][2]float32{p1, p2, p3})
			verts = append(verts[:i], verts[i+1:]...)
			earFound = true
			break
		}

		if !earFound {
			break // If no ear found, polygon might be degenerate or self-intersecting
		}
	}

	if len(verts) == 3 {
		triangles = append(triangles, [3][2]float32{
			points[verts[0]],
			points[verts[1]],
			points[verts[2]],
		})
	}

	return triangles
}
func area(points [][2]float32) float32 {
	var a float32
	for i := range points {
		var j = (i + 1) % len(points)
		a += points[i][0]*points[j][1] - points[j][0]*points[i][1]
	}
	return a / 2
}
func isClockwise(points [3][2]float32) bool {
	var p = points
	var area = (p[1][0]-p[0][0])*(p[2][1]-p[0][1]) - (p[2][0]-p[0][0])*(p[1][1]-p[0][1])
	return area < 0 // negative => clockwise
}
func isEar(points [][2]float32, verts []int, i1, i2, i3 int, ccw bool) bool {
	var p1, p2, p3 = points[i1], points[i2], points[i3]
	var cross = (p2[0]-p1[0])*(p3[1]-p1[1]) - (p2[1]-p1[1])*(p3[0]-p1[0])
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
		if pointInTriangle(points[vi], p1, p2, p3) {
			return false
		}
	}

	return true
}
func pointInTriangle(p, a, b, c [2]float32) bool {
	var v0 = [2]float32{c[0] - a[0], c[1] - a[1]}
	var v1 = [2]float32{b[0] - a[0], b[1] - a[1]}
	var v2 = [2]float32{p[0] - a[0], p[1] - a[1]}
	var dot00 = v0[0]*v0[0] + v0[1]*v0[1]
	var dot01 = v0[0]*v1[0] + v0[1]*v1[1]
	var dot02 = v0[0]*v2[0] + v0[1]*v2[1]
	var dot11 = v1[0]*v1[0] + v1[1]*v1[1]
	var dot12 = v1[0]*v2[0] + v1[1]*v2[1]
	var invDenom = 1 / (dot00*dot11 - dot01*dot01)
	var u = (dot11*dot02 - dot01*dot12) * invDenom
	var v = (dot00*dot12 - dot01*dot02) * invDenom

	return (u >= 0) && (v >= 0) && (u+v <= 1)
}

func separateShapes(points [][2]float32) [][][2]float32 {
	var result [][][2]float32
	var current = [][2]float32{}

	for _, p := range points {
		if number.IsNaN(p[0]) || number.IsNaN(p[1]) {
			if len(current) > 0 { // finish current shape and start a new one
				result = append(result, current)
				current = [][2]float32{}
			}
			continue
		}

		current = append(current, p)
	}

	if len(current) > 0 { // add the last shape if it has points
		result = append(result, current)
	}

	return result
}

//=================================================================
// camera

var rlCam = rl.Camera2D{}
var maskX, maskY, maskW, maskH int

// call before draw to update camera but use screen space instead of camera space
func (camera *Camera) update() {
	tryRecreateWindow()

	rlCam.Target.X = float32(camera.X)
	rlCam.Target.Y = float32(camera.Y)
	rlCam.Rotation = float32(camera.Angle)
	rlCam.Zoom = float32(camera.Zoom)
	rlCam.Offset.X = float32(camera.ScreenX) + float32(camera.ScreenWidth)*float32(camera.PivotX)
	rlCam.Offset.Y = float32(camera.ScreenY) + float32(camera.ScreenHeight)*float32(camera.PivotY)

	var mx = number.Biggest(camera.MaskX, camera.ScreenX)
	var my = number.Biggest(camera.MaskY, camera.ScreenY)
	var maxW = camera.ScreenX + camera.ScreenWidth - mx
	var maxH = camera.ScreenY + camera.ScreenHeight - my
	var mw = number.Smallest(camera.MaskWidth, maxW)
	var mh = number.Smallest(camera.MaskHeight, maxH)

	maskX, maskY, maskW, maskH = mx, my, mw, mh
}

// call before draw to update camera and use camera space
func (camera *Camera) begin() {
	camera.update()
	if camera.Batch {
		return
	}

	rl.BeginMode2D(rlCam)
	rl.BeginScissorMode(int32(maskX), int32(maskY), int32(maskW), int32(maskH))
}

// call after draw to get back to using screen space
func (camera *Camera) end() {
	if camera.Batch {
		return
	}

	rl.EndScissorMode()
	rl.EndMode2D()
}

func (camera *Camera) isAreaVisible(x, y, width, height, angle float32) bool {
	var tlx, tly = x, y
	var trx, try = point.MoveAtAngle(tlx, tly, angle, width)
	var brx, bry = point.MoveAtAngle(trx, try, angle+90, height)
	var blx, bly = point.MoveAtAngle(tlx, tly, angle+90, height)
	var stlx, stly = camera.PointToScreen(tlx, tly)
	var strx, stry = camera.PointToScreen(trx, try)
	var sbrx, sbry = camera.PointToScreen(brx, bry)
	var sblx, sbly = camera.PointToScreen(blx, bly)
	var mtlx, mtly = camera.MaskX, camera.MaskY
	var mbrx, mbry = camera.MaskX + camera.MaskWidth, camera.MaskY + camera.MaskHeight
	var minX = number.Smallest(stlx, strx, sbrx, sblx)
	var maxX = number.Biggest(stlx, strx, sbrx, sblx)
	var minY = number.Smallest(stly, stry, sbry, sbly)
	var maxY = number.Biggest(stly, stry, sbry, sbly)

	return maxY > mtly && minY < mbry && maxX > mtlx && minX < mbrx
}

func tryRecreateWindow() {
	if internal.WindowReady {
		return
	}

	if !rl.IsWindowReady() {
		window.Recreate()
	}
}
