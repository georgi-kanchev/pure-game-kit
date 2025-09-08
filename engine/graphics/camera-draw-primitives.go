package graphics

import (
	"pure-kit/engine/geometry/point"
	"pure-kit/engine/internal"
	"pure-kit/engine/utility/number"

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
func (camera *Camera) DrawCircle(x, y, radius float32, color uint) {
	camera.begin()
	rl.DrawCircle(int32(x), int32(y), radius, rl.GetColor(color))
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
func (camera *Camera) DrawRectangle(x, y, width, height, angle float32, color uint) {
	if !camera.isAreaVisible(x, y, width, height, angle) {
		return
	}

	camera.begin()
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

	rl.DrawRectanglePro(rect, rl.Vector2{X: 0, Y: 0}, angle, rl.GetColor(color))
	camera.end()
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
