package graphics

import (
	"pure-game-kit/internal"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/point"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (c *Camera) DrawColor(color uint) {
	c.begin()
	c.end()

	var x, y, w, h = c.ScreenX, c.ScreenY, c.ScreenWidth, c.ScreenHeight
	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(h), getColor(color))
}
func (c *Camera) DrawScreenFrame(thickness int, color uint) {
	c.begin()
	c.end()

	var x, y, w, h = c.ScreenX, c.ScreenY, c.ScreenWidth, c.ScreenHeight
	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(thickness), getColor(color))             // top
	rl.DrawRectangle(int32(x+w-thickness), int32(y), int32(thickness), int32(h), getColor(color)) // right
	rl.DrawRectangle(int32(x), int32(y+h-thickness), int32(w), int32(thickness), getColor(color)) // bottom
	rl.DrawRectangle(int32(x), int32(y), int32(thickness), int32(h), getColor(color))             // left
}

func (c *Camera) DrawGrid(thickness, spacingX, spacingY float32, color uint) {
	c.begin()
	skipStartEnd = true
	var sx, sy, sw, sh = c.ScreenX, c.ScreenY, c.ScreenWidth, c.ScreenHeight
	var ulx, uly = c.PointFromScreen(sx, sy)
	var urx, ury = c.PointFromScreen(sx+sw, sy)
	var lrx, lry = c.PointFromScreen(sx+sw, sy+sh)
	var llx, lly = c.PointFromScreen(sx, sy+sh)
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

	var left = number.RoundDown(minX/spacingX) * spacingX
	var right = number.RoundUp(maxX/spacingX) * spacingX
	var top = number.RoundDown(minY/spacingY) * spacingY
	var bottom = number.RoundUp(maxY/spacingY) * spacingY

	for x := left; x <= right; x += spacingX { // vertical
		var myThickness = thickness
		if number.DivisionRemainder(x, spacingX*10) == 0 {
			myThickness *= 3
		}

		c.DrawLine(x, top, x, bottom, myThickness, color)
	}
	for y := top; y <= bottom; y += spacingY { // horizontal
		var myThickness = thickness
		if number.DivisionRemainder(y, spacingY*10) == 0 {
			myThickness *= 3
		}

		c.DrawLine(left, y, right, y, myThickness, color)
	}

	if top <= 0 && bottom >= 0 { // x
		c.DrawLine(left, 0, right, 0, thickness*6, color)
	}
	if left <= 0 && right >= 0 { // y
		c.DrawLine(0, top, 0, bottom, thickness*6, color)
	}
	skipStartEnd = false
	c.end()
}

//=================================================================

func (c *Camera) DrawLine(ax, ay, bx, by, thickness float32, color uint) {
	c.begin()
	rl.DrawLineEx(rl.Vector2{X: ax, Y: ay}, rl.Vector2{X: bx, Y: by}, thickness, getColor(color))
	c.end()
}

// multiple line paths can be separated by a [NaN, NaN]
func (c *Camera) DrawLinesPath(thickness float32, color uint, points ...[2]float32) {
	c.begin()
	var col = getColor(color)
	for i := 1; i < len(points); i++ {
		var p1, p2 = points[i-1], points[i]
		if !number.IsNaN(p1[0]) && !number.IsNaN(p2[0]) {
			batch.QueueLine(p1[0], p1[1], p2[0], p2[1], thickness, col)
		}
	}
	batch.Draw()
	c.end()
}

func (c *Camera) DrawQuadFrame(x, y, width, height, angle, thickness float32, color uint) {
	if thickness == 0 {
		return
	}

	c.begin()
	skipStartEnd = true
	defer func() {
		skipStartEnd = false
		c.end()
	}()

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

		c.DrawQuad(x, y, width, thickness, angle, color)
		c.DrawQuad(trx, try, thickness, height, angle, color)
		c.DrawQuad(brx, bry, width, thickness, angle, color)
		c.DrawQuad(x, y, thickness, height, angle, color)
		return
	}

	var x1, y1 = point.MoveAtAngle(x, y, angle-90, thickness)
	var tlx, tly = point.MoveAtAngle(x1, y1, angle-180, thickness)
	var trx, try = point.MoveAtAngle(x1, y1, angle, width)
	var blx, bly = point.MoveAtAngle(tlx, tly, angle+90, height+thickness)

	c.DrawQuad(tlx, tly, width+thickness*2, thickness, angle, color)
	c.DrawQuad(trx, try, thickness, height+thickness*2, angle, color)
	c.DrawQuad(blx, bly, width+thickness*2, thickness, angle, color)
	c.DrawQuad(tlx, tly, thickness, height+thickness*2, angle, color)
}
func (c *Camera) DrawQuad(x, y, width, height, angle float32, colors ...uint) {
	if !c.isAreaVisible(x, y, width, height, angle) {
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

	if len(colors) == 0 {
		colors = append(colors, palette.White)
	}

	if len(colors) == 1 {
		c.begin()
		rl.DrawRectanglePro(rect, rl.Vector2{X: 0, Y: 0}, angle, getColor(colors[0]))
		c.end()
		return // draw regular rect with one provided color
	}

	for len(colors) < 4 { // if fewer than 4 colors, pad with last provided color
		colors = append(colors, colors[len(colors)-1])
	}

	var tl, tr = getColor(colors[0]), getColor(colors[1])
	var br, bl = getColor(colors[2]), getColor(colors[3])
	var prevAng = c.Angle
	c.Angle = angle
	c.begin()
	rect.X, rect.Y = point.RotateAroundPoint(rect.X, rect.Y, c.X, c.Y, -angle)
	rl.DrawRectangleGradientEx(rect, tl, bl, br, tr)
	c.end()
	c.Angle = prevAng
}
func (c *Camera) DrawQuadRounded(x, y, width, height, radius, angle float32, color uint) {
	c.begin()
	skipStartEnd = true
	var maxRadius = width / 2
	if height/2 < maxRadius {
		maxRadius = height / 2
	}
	radius = min(radius, maxRadius)

	var sz = radius * 2
	var vw, vh = width - 2*radius, height
	var vx, vy = point.RotateAroundPoint(x+radius, y, x, y, angle)
	var hw, hh = width, height - 2*radius
	var hx, hy = point.RotateAroundPoint(x, y+radius, x, y, angle)
	var tlx, tly = point.RotateAroundPoint(x+radius, y+radius, x, y, angle)
	var trx, try = point.RotateAroundPoint(x+width-radius, y+radius, x, y, angle)
	var blx, bly = point.RotateAroundPoint(x+radius, y+height-radius, x, y, angle)
	var brx, bry = point.RotateAroundPoint(x+width-radius, y+height-radius, x, y, angle)

	c.DrawQuad(vx, vy, vw, vh, angle, color)
	c.DrawQuad(hx, hy, hw, hh, angle, color)
	c.DrawArc(tlx, tly, sz, sz, 1, 0, 24, color)
	c.DrawArc(trx, try, sz, sz, 1, 0, 24, color)
	c.DrawArc(blx, bly, sz, sz, 1, 0, 24, color)
	c.DrawArc(brx, bry, sz, sz, 1, 0, 24, color)
	skipStartEnd = false
	c.end()
}

func (c *Camera) DrawPoints(radius float32, color uint, points ...[2]float32) {
	c.begin()
	skipStartEnd = true
	for _, pt := range points {
		c.DrawCircle(pt[0], pt[1], radius, color)
	}
	skipStartEnd = false
	c.end()
}
func (c *Camera) DrawCircle(x, y, radius float32, colors ...uint) {
	const segments = 24

	if len(colors) == 0 {
		colors = append(colors, palette.White)
	}

	if len(colors) == 1 {
		c.DrawArc(x, y, radius*2, radius*2, 1, 0, segments, colors[0])
	} else if len(colors) > 1 {
		c.begin()
		var step = float32(360.0 / float32(segments))
		rl.Begin(rl.Triangles)
		for i := range segments {
			var ang1, ang2 = float32(i) * step, float32(i+1) * step
			var p1x, p1y = point.MoveAtAngle(x, y, ang1, radius)
			var p2x, p2y = point.MoveAtAngle(x, y, ang2, radius)
			rl.Color4ub(color.Channels(colors[1]))
			rl.Vertex2f(p2x, p2y)
			rl.Color4ub(color.Channels(colors[1]))
			rl.Vertex2f(p1x, p1y)
			rl.Color4ub(color.Channels(colors[0]))
			rl.Vertex2f(x, y)
		}
		rl.End()
		c.end()
	}
}
func (c *Camera) DrawArc(x, y, width, height, fill, angle float32, segments int, color uint) {
	var fillAngle = number.Limit(fill, 0, 1) * 360
	if fillAngle < 360 {
		segments = max(int((fillAngle/360.0)*float32(segments)), 3)
	}

	var points = make([]rl.Vector2, segments+2)
	var radiusH, radiusV = width / 2, height / 2
	var halfPie = fillAngle / 2.0
	var sinRot, cosRot = internal.SinCos(angle)

	points[0] = rl.Vector2{X: x, Y: y}
	for i := 0; i <= segments; i++ {
		var t = float32(i) / float32(segments)
		var localAngDeg = (halfPie - (t * fillAngle))
		var sin, cos = internal.SinCos(localAngDeg)
		var localX, localY = cos * radiusH, sin * radiusV
		var rotatedX, rotatedY = localX*cosRot - localY*sinRot, localX*sinRot + localY*cosRot
		points[i+1] = rl.Vector2{X: x + rotatedX, Y: y + rotatedY}
	}

	c.begin()
	rl.DrawTriangleFan(points, getColor(color))
	c.end()
}

// works with convex + concave (non-self-intersacting) shapes
//
// multiple shapes can be separated by a [NaN, NaN]
func (c *Camera) DrawShapes(color uint, points ...[2]float32) {
	c.begin()

	var flatPoints, shapeCounts = separateShapes(points)
	var offset = 0
	var renderColor = getColor(color)

	for _, count := range shapeCounts {
		var shape = flatPoints[offset : offset+(count*2)]
		offset += (count * 2)

		// Remove closing point if duplicate
		if shape[0] == shape[len(shape)-2] && shape[1] == shape[len(shape)-1] {
			shape = shape[:len(shape)-2]
		}

		var triangles = triangulate(shape)
		if len(triangles) == 0 {
			continue
		}

		for i := 0; i < len(triangles); i += 6 {
			var x1, y1 = triangles[i+0], triangles[i+1]
			var x2, y2 = triangles[i+2], triangles[i+3]
			var x3, y3 = triangles[i+4], triangles[i+5]

			if !isClockwiseFlat(triangles[i : i+6]) {
				x1, y1, x3, y3 = x3, y3, x1, y1
			}

			batch.QueueTriangle(x1, y1, x2, y2, x3, y3, renderColor)
		}
	}

	batch.Draw()
	c.end()
}

// works with convex shapes only
//
// multiple shapes can be separated by a [NaN, NaN] point
func (c *Camera) DrawShapesFast(color uint, points ...[2]float32) {
	c.begin()

	var flatPoints, shapeCounts = separateShapes(points)
	var offset = 0
	var renderColor = getColor(color)

	for _, count := range shapeCounts {
		var shapeData = flatPoints[offset : offset+(count*2)]
		offset += (count * 2)

		var isReverse = area(shapeData) >= 0
		var x0, y0 = shapeData[0], shapeData[1]

		for i := 1; i < count-1; i++ {
			var x1, y1, x2, y2 float32

			if isReverse {
				var idx1 = (count - i) * 2
				var idx2 = (count - i - 1) * 2
				x1, y1 = shapeData[idx1], shapeData[idx1+1]
				x2, y2 = shapeData[idx2], shapeData[idx2+1]
			} else {
				var idx1 = i * 2
				var idx2 = (i + 1) * 2
				x1, y1 = shapeData[idx1], shapeData[idx1+1]
				x2, y2 = shapeData[idx2], shapeData[idx2+1]
			}

			batch.QueueTriangle(x0, y0, x1, y1, x2, y2, renderColor)
		}
	}

	batch.Draw()
	c.end()
}

//=================================================================

func (c *Camera) DrawTexture(textureId string, x, y, width, height, angle float32, color uint) {
	var texture, has = internal.Textures[textureId]
	if !c.isAreaVisible(x, y, width, height, angle) || !has {
		return
	}

	c.begin()
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

	rl.DrawTexturePro(*texture, rectTexture, rectWorld, rl.Vector2{}, 0, getColor(color))
	c.end()
}
func (c *Camera) DrawText(text string, x, y, height float32) {
	c.DrawTextAdvanced("", text, x, y, height, 0.5, 0, palette.White)
}
func (c *Camera) DrawTextAdvanced(fontId, text string, x, y, height, thickness, symbolGap float32, color uint) {
	c.begin()

	var sh = internal.ShaderText
	var font, has = internal.Fonts[fontId]

	if !has {
		var def, hasDefault = internal.Fonts[""]
		if hasDefault {
			font = def
		} else {
			var fallback = rl.GetFontDefault()
			font = &fallback
		}

	}

	if sh.ID != 0 {
		rl.BeginShaderMode(sh)
		//rl.SetShaderValue(sh, internal.ShaderTextLoc, []float32{thickness, 0.02}, rl.ShaderUniformVec2)
	}

	rl.DrawTextPro(*font, text, rl.Vector2{X: x, Y: y}, rl.Vector2{}, 0, height, symbolGap, getColor(color))

	if sh.ID != 0 {
		rl.EndShaderMode()
	}

	c.end()
}
