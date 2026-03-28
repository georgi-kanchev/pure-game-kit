package graphics

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/debug"
	"pure-game-kit/execution/condition"
	"pure-game-kit/internal"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/text"
	tm "pure-game-kit/utility/time"
	"pure-game-kit/utility/time/unit"
)

func (c *Camera) DrawColor(color uint) {
	var tlX, tlY = c.PointFromEdge(0, 0) // Top-Left
	var trX, trY = c.PointFromEdge(1, 0) // Top-Right
	var brX, brY = c.PointFromEdge(1, 1) // Bottom-Right
	var blX, blY = c.PointFromEdge(0, 1) // Bottom-Left
	var renderColor = getColor(color)

	c.begin()
	batch.QueueTriangle(tlX, tlY, trX, trY, brX, brY, renderColor)
	batch.QueueTriangle(tlX, tlY, brX, brY, blX, blY, renderColor)
	batch.Draw()
	c.end()
}
func (c *Camera) DrawGrid(thickness, spacingX, spacingY float32, color uint) {
	if spacingX*c.Zoom < 1 && spacingY*c.Zoom < 1 {
		return // way too dense grid - give up
	}
	c.begin()

	var renderColor = getColor(color)
	var sx, sy, sw, sh = c.area()
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

	for x := left; x <= right; x += spacingX {
		var myThickness = thickness
		if number.DivisionRemainder(x, spacingX*10) == 0 {
			myThickness *= 3
		}
		batch.QueueLine(x, top, x, bottom, myThickness, renderColor)
	}
	for y := top; y <= bottom; y += spacingY {
		var myThickness = thickness
		if number.DivisionRemainder(y, spacingY*10) == 0 {
			myThickness *= 3
		}
		batch.QueueLine(left, y, right, y, myThickness, renderColor)
	}

	if top <= 0 && bottom >= 0 {
		batch.QueueLine(left, 0, right, 0, thickness*6, renderColor)
	}
	if left <= 0 && right >= 0 {
		batch.QueueLine(0, top, 0, bottom, thickness*6, renderColor)
	}

	batch.Draw()
	c.end()
}

//=================================================================

func (c *Camera) DrawLine(ax, ay, bx, by, thickness float32, color uint) {
	c.begin()
	batch.QueueLine(ax, ay, bx, by, thickness, getColor(color))
	batch.Draw()
	c.end()
}

// multiple line paths can be separated by a NaN pair
func (c *Camera) DrawLinesPath(thickness float32, color uint, points ...float32) {
	if thickness == 0 || color == 0 || len(points) == 0 {
		return
	}

	c.begin()
	var col = getColor(color)
	for i := 2; i < len(points); i += 2 {
		var x1, y1 = points[i-2], points[i-1]
		var x2, y2 = points[i], points[i+1]

		if !number.IsNaN(x1) && !number.IsNaN(x2) {
			batch.QueueLine(x1, y1, x2, y2, thickness, col)
		}
	}
	batch.Draw()
	c.end()
}

func (c *Camera) DrawQuad(x, y, width, height, angle float32, color uint) {
	c.begin()
	batch.QueueQuad(x, y, width, height, angle, getColor(color))
	batch.Draw()
	c.end()
}
func (c *Camera) DrawQuadFrame(x, y, width, height, angle, thickness float32, color uint) {
	c.begin()
	var sinRot, cosRot = internal.SinCos(angle)
	var transform = func(px, py float32) (float32, float32) {
		return x + (px*cosRot - py*sinRot), y + (px*sinRot + py*cosRot)
	}

	var absT = thickness
	if absT < 0 {
		absT = -absT
	}

	var h = thickness / 2
	var hx1, hx2, vy1, vy2 float32

	if thickness > 0 {
		hx1, hx2 = -absT, width+absT
		vy1, vy2 = 0, height
	} else {
		hx1, hx2 = 0, width
		vy1, vy2 = absT, height-absT
	}

	var x1, y1 = transform(hx1, -h)
	var x2, y2 = transform(hx2, -h)
	var x3, y3 = transform(hx2, height+h)
	var x4, y4 = transform(hx1, height+h)
	var x5, y5 = transform(width+h, vy1)
	var x6, y6 = transform(width+h, vy2)
	var x7, y7 = transform(-h, vy2)
	var x8, y8 = transform(-h, vy1)
	var col = getColor(color)
	batch.QueueLine(x1, y1, x2, y2, absT, col)
	batch.QueueLine(x3, y3, x4, y4, absT, col)
	batch.QueueLine(x5, y5, x6, y6, absT, col)
	batch.QueueLine(x7, y7, x8, y8, absT, col)
	batch.Draw()
	c.end()
}

func (c *Camera) DrawPoints(radius float32, color uint, points ...float32) {
	c.begin()
	batch.skipStartEnd = true
	batch.skipDraw = true
	for i := 0; i < len(points); i += 2 {
		if i+1 >= len(points) {
			break
		}

		var x, y = points[i], points[i+1]
		c.DrawCircle(x, y, radius, 16, color)
	}
	batch.skipDraw = false
	batch.skipStartEnd = false
	batch.Draw()
	c.end()
}
func (c *Camera) DrawCircle(x, y, radius float32, segments int, color uint) {
	c.DrawArc(x, y, radius*2, radius*2, 1, 0, segments, color)
}
func (c *Camera) DrawArc(x, y, width, height, fill, angle float32, segments int, color uint) {
	c.begin()
	var fillAngle = number.Limit(fill, 0, 1) * 360
	if fillAngle < 360 {
		segments = max(int((fillAngle/360.0)*float32(segments)), 3)
	}

	var radiusH, radiusV = width / 2, height / 2
	var halfPie = fillAngle / 2.0
	var sinRot, cosRot = internal.SinCos(angle)
	var renderColor = getColor(color)
	var s0, c0 = internal.SinCos(halfPie)
	var lx0, ly0 = c0 * radiusH, s0 * radiusV
	var prevX = x + (lx0*cosRot - ly0*sinRot)
	var prevY = y + (lx0*sinRot + ly0*cosRot)

	for i := 1; i <= segments; i++ {
		var t = float32(i) / float32(segments)
		var ang = halfPie - (t * fillAngle)
		var si, co = internal.SinCos(ang)
		var lxi, lyi = co * radiusH, si * radiusV
		var currX = x + (lxi*cosRot - lyi*sinRot)
		var currY = y + (lxi*sinRot + lyi*cosRot)

		batch.QueueTriangle(x, y, prevX, prevY, currX, currY, renderColor)
		prevX, prevY = currX, currY
	}
	if !batch.skipDraw {
		batch.Draw()
	}
	c.end()
}

// works with convex + concave (non-self-intersacting) shapes
//
// multiple shapes can be separated by a NaN pair
func (c *Camera) DrawShapes(color uint, points ...float32) {
	c.begin()

	var ptsCountsPerShape = separateShapes(points)
	var offset = 0
	var renderColor = getColor(color)

	c.Effects.updateUniforms(1, 1, nil, nil, false)

	for _, count := range ptsCountsPerShape {
		for offset < len(points) && number.IsNaN(points[offset]) {
			offset += 2
		}

		var shape = points[offset : offset+(count*2)]
		offset += count * 2

		if count > 2 && shape[0] == shape[len(shape)-2] && shape[1] == shape[len(shape)-1] {
			shape = shape[:len(shape)-2]
			count--
		}

		if count < 3 {
			continue
		}

		if isConvex(shape, count) { // fast path - triangle fan (convex only)
			var isReverse = area(shape) >= 0
			var x0, y0 = shape[0], shape[1]

			for i := 1; i < count-1; i++ {
				var x1, y1, x2, y2 float32

				if isReverse {
					var idx1 = (count - i) * 2
					var idx2 = (count - i - 1) * 2
					x1, y1 = shape[idx1], shape[idx1+1]
					x2, y2 = shape[idx2], shape[idx2+1]
				} else {
					var idx1 = i * 2
					var idx2 = (i + 1) * 2
					x1, y1 = shape[idx1], shape[idx1+1]
					x2, y2 = shape[idx2], shape[idx2+1]
				}

				batch.QueueTriangle(x0, y0, x1, y1, x2, y2, renderColor)
			}
			continue
		}

		var triangles = triangulate(shape) // slow path - ear clipping/triangulation (concave only)
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

//=================================================================

func (c *Camera) DrawTexture(assetId string, x, y, scaleX, scaleY, angle float32, color uint) {
	var w, h = assets.Size(assetId)
	drawTexture.TextureId, drawTexture.Tint = assetId, color
	drawTexture.X, drawTexture.Y = x, y
	drawTexture.PivotX, drawTexture.PivotY = 0, 0
	drawTexture.Width, drawTexture.Height = float32(w), float32(h)
	drawTexture.Angle = angle
	drawTexture.ScaleX, drawTexture.ScaleY = scaleX, scaleY
	c.DrawSprites(drawTexture)
}
func (c *Camera) DrawText(text string, x, y, height float32) {
	c.DrawTextAdvanced("", text, x, y, height, 0, 0, 0, palette.White)
}
func (c *Camera) DrawTextAdvanced(fontId, text string, x, y, lineHeight, angle, symbolGap, lineGap float32, color uint) {
	drawText.FontId, drawText.Tint = fontId, color
	drawText.X, drawText.Y = x, y
	drawText.PivotX, drawText.PivotY = 0, 0
	drawText.Text, drawText.Angle = text, angle
	drawText.SymbolGap, drawText.LineGap = symbolGap, lineGap
	drawText.Width, drawText.Height = 99999, 99999
	drawText.WordWrap, drawText.LineHeight = false, lineHeight
	c.DrawTextBoxes(drawText)
}
func (c *Camera) DrawTextDebug(fps, time, assets, memory bool) {
	if condition.TrueEvery(0.15, ";;;debug") {
		debugStr = ""
		if fps {
			debugStr += text.New("FPS ", int(internal.FPS), " (", int(internal.AverageFPS), ")\n\n")
		}
		if time {
			debugStr += text.New(
				"Time: \n",
				"Running = ", tm.AsClock12(internal.Runtime, ":", unit.Hour|unit.Timer, false), "\n",
				"Frame Busy = ", number.Round(internal.FrameTime*1000, 3), "ms ",
				"(", number.Round((internal.FrameTime/internal.DeltaTime)*100), "%)\n",
				"Frame Idle = ", number.Round((internal.DeltaTime-internal.FrameTime)*1000, 3), "ms ",
				"(", number.Round(((internal.DeltaTime-internal.FrameTime)/internal.DeltaTime)*100), "%)\n",
				"Frame Total = ", number.Round((internal.DeltaTime)*1000, 3), "ms ",
				"\n\n")
		}
		if assets {
			debugStr += text.New("Assets: \n",
				"Textures = ", len(internal.Textures), "\n",
				"Fonts = ", len(internal.Fonts), "\n",
				"Sounds = ", len(internal.Sounds), "\n",
				"Music = ", len(internal.Music), "\n",
				"Tile Data = ", len(internal.TileDatas), "\n\n")
		}
		if memory {
			debugStr += debug.MemoryUsage()
		}
	}

	var tlx, tly = c.PointFromEdge(0, 0)
	c.DrawText(debugStr, tlx+10/c.Zoom, tly, 40/c.Zoom)
}
