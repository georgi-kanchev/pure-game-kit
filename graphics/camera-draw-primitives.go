package graphics

import (
	"pure-game-kit/internal"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/point"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (camera *Camera) DrawColor(color uint) {
	camera.begin()
	camera.end()

	var x, y, w, h = camera.ScreenX, camera.ScreenY, camera.ScreenWidth, camera.ScreenHeight
	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(h), rl.GetColor(color))
}
func (camera *Camera) DrawScreenFrame(thickness int, color uint) {
	camera.begin()
	camera.end()

	var x, y, w, h = camera.ScreenX, camera.ScreenY, camera.ScreenWidth, camera.ScreenHeight
	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(thickness), rl.GetColor(color))             // top
	rl.DrawRectangle(int32(x+w-thickness), int32(y), int32(thickness), int32(h), rl.GetColor(color)) // right
	rl.DrawRectangle(int32(x), int32(y+h-thickness), int32(w), int32(thickness), rl.GetColor(color)) // bottom
	rl.DrawRectangle(int32(x), int32(y), int32(thickness), int32(h), rl.GetColor(color))             // left
}
func (camera *Camera) DrawGrid(thickness, spacingX, spacingY float32, color uint) {
	camera.begin()
	camera.Batch = true
	var sx, sy, sw, sh = camera.ScreenX, camera.ScreenY, camera.ScreenWidth, camera.ScreenHeight
	var ulx, uly = camera.PointFromScreen(sx, sy)
	var urx, ury = camera.PointFromScreen(sx+sw, sy)
	var lrx, lry = camera.PointFromScreen(sx+sw, sy+sh)
	var llx, lly = camera.PointFromScreen(sx, sy+sh)
	var xs = []float32{ulx, urx, llx, lrx}
	var ys = []float32{uly, ury, lly, lry}
	var minX, maxX = xs[0], xs[0]
	var minY, maxY = ys[0], ys[0]

	for i := 1; i < 4; i++ {
		if xs[i] < minX {
			minX = xs[i]
		}
		if xs[i] > maxX {
			maxX = xs[i]
		}
		if ys[i] < minY {
			minY = ys[i]
		}
		if ys[i] > maxY {
			maxY = ys[i]
		}
	}

	var left = number.RoundDown(minX/spacingX, -1) * spacingX
	var right = number.RoundUp(maxX/spacingX, -1) * spacingX
	var top = number.RoundDown(minY/spacingY, -1) * spacingY
	var bottom = number.RoundUp(maxY/spacingY, -1) * spacingY

	for x := left; x <= right; x += spacingX { // vertical
		var myThickness = thickness
		if number.DivisionRemainder(x, spacingX*10) == 0 {
			myThickness *= 2
		}

		camera.DrawLine(x, top, x, bottom, myThickness, color)
	}
	for y := top; y <= bottom; y += spacingY { // horizontal
		var myThickness = thickness
		if number.DivisionRemainder(y, spacingY*10) == 0 {
			myThickness *= 2
		}

		camera.DrawLine(left, y, right, y, myThickness, color)
	}

	if top <= 0 && bottom >= 0 { // x
		camera.DrawLine(left, 0, right, 0, thickness*4, color)
	}
	if left <= 0 && right >= 0 { // y
		camera.DrawLine(0, top, 0, bottom, thickness*4, color)
	}
	camera.Batch = false
	camera.end()
}

func (camera *Camera) DrawLine(ax, ay, bx, by, thickness float32, color uint) {
	camera.begin()
	rl.DrawLineEx(rl.Vector2{X: ax, Y: ay}, rl.Vector2{X: bx, Y: by}, thickness, rl.GetColor(color))
	camera.end()
}
func (camera *Camera) DrawLinesPath(thickness float32, color uint, points ...[2]float32) {
	camera.begin()
	for i := 1; i < len(points); i++ {
		rl.DrawLineEx(
			rl.Vector2{X: points[i-1][0], Y: points[i-1][1]},
			rl.Vector2{X: points[i][0], Y: points[i][1]}, thickness, rl.GetColor(color))
	}
	camera.end()
}
func (camera *Camera) DrawCircle(x, y, radius float32, colors ...uint) {
	camera.begin()
	if len(colors) == 1 {
		rl.DrawCircle(int32(x), int32(y), radius, rl.GetColor(colors[0]))
	} else if len(colors) == 2 {
		rl.DrawCircleGradient(int32(x), int32(y), radius, rl.GetColor(colors[0]), rl.GetColor(colors[1]))
	}
	camera.end()
}
func (camera *Camera) DrawEllipse(x, y, width, height float32, color uint) {
	camera.begin()
	rl.DrawEllipse(int32(x), int32(y), width/2, height/2, rl.GetColor(color))
	camera.end()
}
func (camera *Camera) DrawFrame(x, y, width, height, angle, thickness float32, color uint) {
	if thickness == 0 {
		return
	}

	camera.begin()
	camera.Batch = true
	defer func() { camera.Batch = false }()

	if width < 0 {
		x, y = point.MoveAtAngle(x, y, angle+180, -width)
		width *= -1
	}
	if height < 0 {
		x, y = point.MoveAtAngle(x, y, angle+270, -height)
		height *= -1
	}

	if thickness < 0 {
		thickness *= -1
		var trx, try = point.MoveAtAngle(x, y, angle, width-thickness)
		var brx, bry = point.MoveAtAngle(x, y, angle+90, height-thickness)

		camera.DrawRectangle(x, y, width, thickness, angle, color)
		camera.DrawRectangle(trx, try, thickness, height, angle, color)
		camera.DrawRectangle(brx, bry, width, thickness, angle, color)
		camera.DrawRectangle(x, y, thickness, height, angle, color)
		return
	}

	var x1, y1 = point.MoveAtAngle(x, y, angle-90, thickness)
	var tlx, tly = point.MoveAtAngle(x1, y1, angle-180, thickness)
	var trx, try = point.MoveAtAngle(x1, y1, angle, width)
	var blx, bly = point.MoveAtAngle(tlx, tly, angle+90, height+thickness)

	camera.DrawRectangle(tlx, tly, width+thickness*2, thickness, angle, color)
	camera.DrawRectangle(trx, try, thickness, height+thickness*2, angle, color)
	camera.DrawRectangle(blx, bly, width+thickness*2, thickness, angle, color)
	camera.DrawRectangle(tlx, tly, thickness, height+thickness*2, angle, color)
	camera.end()
}
func (camera *Camera) DrawRectangle(x, y, width, height, angle float32, colors ...uint) {
	if !camera.isAreaVisible(x, y, width, height, angle) {
		return
	}

	var rect = rl.Rectangle{X: x, Y: y, Width: width, Height: height}

	// raylib doesn't seem to have negative width/height???
	if rect.Width < 0 && rect.Height > 0 {
		rect.X, rect.Y = point.MoveAtAngle(rect.X, rect.Y, angle+180, -rect.Width)
		rect.Width *= -1
	}
	if rect.Height < 0 && rect.Width > 0 {
		rect.X, rect.Y = point.MoveAtAngle(rect.X, rect.Y, angle+270, -rect.Height)
		rect.Height *= -1
	}

	if len(colors) == 1 {
		camera.begin()
		rl.DrawRectanglePro(rect, rl.Vector2{X: 0, Y: 0}, angle, rl.GetColor(colors[0]))
		camera.end()
		return // draw regular rect with one provided color
	}

	for len(colors) < 4 { // if fewer than 4 colors, pad with last provided color
		colors = append(colors, colors[len(colors)-1])
	}

	var tl, tr = rl.GetColor(colors[0]), rl.GetColor(colors[1])
	var br, bl = rl.GetColor(colors[2]), rl.GetColor(colors[3])
	var prevAng = camera.Angle
	camera.Angle = angle
	camera.begin()
	rect.X, rect.Y = point.RotateAroundPoint(rect.X, rect.Y, camera.X, camera.Y, -angle)
	rl.DrawRectangleGradientEx(rect, tl, bl, br, tr)
	camera.end()
	camera.Angle = prevAng
}
func (camera *Camera) DrawShape(color uint, points ...[2]float32) {
	camera.begin()
	defer camera.end()
	if len(points) < 3 {
		return
	}

	var triangles = triangulate(points)
	if len(triangles) == 0 {
		return
	}

	for _, tri := range triangles {
		var v1 = rl.NewVector2(tri[0][0], tri[0][1])
		var v2 = rl.NewVector2(tri[1][0], tri[1][1])
		var v3 = rl.NewVector2(tri[2][0], tri[2][1])
		rl.DrawTriangle(v3, v2, v1, rl.GetColor(color))
	}
}

func (camera *Camera) DrawTexture(textureId string, x, y, width, height, angle float32, color uint) {
	if !camera.isAreaVisible(x, y, width, height, angle) {
		return
	}

	var texture, has = internal.Textures[textureId]
	if !has {
		return
	}

	camera.begin()
	var texX, texY float32 = 0.0, 0.0
	var w, h = width, height
	var texW, texH = texture.Width, texture.Height
	var rectTexture = rl.Rectangle{X: texX, Y: texY, Width: float32(texW), Height: float32(texH)}
	var rectWorld = rl.Rectangle{X: x, Y: y, Width: float32(w), Height: float32(h)}

	if rectWorld.Width < 0 {
		rectWorld.X, rectWorld.Y = point.MoveAtAngle(rectWorld.X, rectWorld.Y, angle+180, -rectWorld.Width)
		rectTexture.Width *= -1
	}
	if rectWorld.Height < 0 {
		rectWorld.X, rectWorld.Y = point.MoveAtAngle(rectWorld.X, rectWorld.Y, angle+270, -rectWorld.Height)
		rectTexture.Height *= -1
	}

	rl.DrawTexturePro(*texture, rectTexture, rectWorld, rl.Vector2{}, 0, rl.GetColor(color))
	camera.end()
}
func (camera *Camera) DrawText(fontId, text string, x, y, height float32, color uint) {
	camera.begin()

	var sh = internal.ShaderText
	var font, has = internal.Fonts[fontId]

	if !has {
		var defaultFont = rl.GetFontDefault()
		font = &defaultFont
	}

	if sh.ID != 0 {
		rl.BeginShaderMode(sh)
		rl.SetShaderValue(sh, rl.GetShaderLocation(sh, "smoothness"), []float32{0}, rl.ShaderUniformFloat)
		rl.SetShaderValue(sh, rl.GetShaderLocation(sh, "thickness"), []float32{0.5}, rl.ShaderUniformFloat)
	}

	rl.DrawTextPro(*font, text, rl.Vector2{X: x, Y: y}, rl.Vector2{}, 0, height, 0, rl.GetColor(color))

	if sh.ID != 0 {
		rl.EndShaderMode()
	}

	camera.end()
}

//=================================================================
// private

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
