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

// multiple line paths can be separated by a [NaN, NaN]
func (c *Camera) DrawLinesPath(thickness float32, color uint, points ...[2]float32) {
	if thickness == 0 || color == 0 || len(points) == 0 {
		return
	}

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
	if width == 0 || height == 0 {
		return
	}
	c.begin()

	if width < 0 {
		x += width
		width *= -1
	}
	if height < 0 {
		y += height
		height *= -1
	}

	var x0, y0, x1, y1 float32     // Outer
	var ix0, iy0, ix1, iy1 float32 // Inner

	if thickness < 0 {
		var t = -thickness
		x0, y0, x1, y1 = 0, 0, width, height
		ix0, iy0, ix1, iy1 = t, t, width-t, height-t
	} else {
		x0, y0, x1, y1 = -thickness, -thickness, width+thickness, height+thickness
		ix0, iy0, ix1, iy1 = 0, 0, width, height
	}

	var sinRot, cosRot = internal.SinCos(angle)
	var renderColor = getColor(color)
	var transform = func(px, py float32) (float32, float32) {
		return x + (px*cosRot - py*sinRot), y + (px*sinRot + py*cosRot)
	}

	var drawRect = func(px1, py1, px2, py2, px3, py3, px4, py4 float32) {
		var v1x, v1y = transform(px1, py1)
		var v2x, v2y = transform(px2, py2)
		var v3x, v3y = transform(px3, py3)
		var v4x, v4y = transform(px4, py4)
		batch.QueueTriangle(v1x, v1y, v2x, v2y, v3x, v3y, renderColor)
		batch.QueueTriangle(v1x, v1y, v3x, v3y, v4x, v4y, renderColor)
	}

	drawRect(x0, y0, x1, y0, x1, iy0, x0, iy0)     // Top
	drawRect(x0, iy1, x1, iy1, x1, y1, x0, y1)     // Bottom
	drawRect(x0, iy0, ix0, iy0, ix0, iy1, x0, iy1) // Left
	drawRect(ix1, iy0, x1, iy0, x1, iy1, ix1, iy1) // Right
	batch.Draw()
	c.end()
}
func (c *Camera) DrawQuad(x, y, width, height, angle float32, color uint) {
	c.begin()
	var sinRot, cosRot = internal.SinCos(angle)
	var renderColor = getColor(color)
	var transform = func(px, py float32) (float32, float32) {
		return x + (px*cosRot - py*sinRot), y + (px*sinRot + py*cosRot)
	}

	var v1x, v1y = transform(0, 0)          // Top-Left
	var v2x, v2y = transform(width, 0)      // Top-Right
	var v3x, v3y = transform(width, height) // Bottom-Right
	var v4x, v4y = transform(0, height)     // Bottom-Left

	batch.QueueTriangle(v1x, v1y, v2x, v2y, v3x, v3y, renderColor)
	batch.QueueTriangle(v1x, v1y, v3x, v3y, v4x, v4y, renderColor)
	batch.Draw()
	c.end()
}
func (c *Camera) DrawQuadRounded(x, y, width, height, radius, angle float32, color uint) {
	c.begin()

	var maxR = width / 2
	if height/2 < maxR {
		maxR = height / 2
	}
	radius = min(radius, maxR)

	var renderColor = getColor(color)
	var sinRot, cosRot = internal.SinCos(angle)
	var transform = func(px, py float32) (float32, float32) {
		return x + (px*cosRot - py*sinRot), y + (px*sinRot + py*cosRot)
	}

	var x1, y1 = transform(0, radius)
	var x2, y2 = transform(width, radius)
	var x3, y3 = transform(width, height-radius)
	var x4, y4 = transform(0, height-radius)
	batch.QueueTriangle(x1, y1, x2, y2, x3, y3, renderColor)
	batch.QueueTriangle(x1, y1, x3, y3, x4, y4, renderColor)

	var x5, y5 = transform(radius, 0)
	var x6, y6 = transform(width-radius, 0)
	var x7, y7 = transform(width-radius, height)
	var x8, y8 = transform(radius, height)
	batch.QueueTriangle(x5, y5, x6, y6, x7, y7, renderColor)
	batch.QueueTriangle(x5, y5, x7, y7, x8, y8, renderColor)

	var segments = 8
	var corners = [4][3]float32{
		{radius, radius, 180},                // Top Left
		{width - radius, radius, 270},        // Top Right
		{width - radius, height - radius, 0}, // Bottom Right
		{radius, height - radius, 90},        // Bottom Left
	}

	for _, corn := range corners {
		var cx, cy = corn[0], corn[1]
		var startAng = corn[2]
		var s0, c0 = internal.SinCos(startAng)
		var px, py = transform(cx+c0*radius, cy-s0*radius)

		for i := 1; i <= segments; i++ {
			var t = float32(i) / float32(segments)
			var currAng = startAng - (t * 90)
			var si, co = internal.SinCos(currAng)
			var ctx, cty = transform(cx+co*radius, cy-si*radius)
			var centX, centY = transform(cx, cy)

			batch.QueueTriangle(centX, centY, px, py, ctx, cty, renderColor)
			px, py = ctx, cty
		}
	}
	batch.Draw()
	c.end()
}

func (c *Camera) DrawPoints(radius float32, color uint, points ...[2]float32) {
	c.begin()
	batch.skipStartEnd = true
	for _, pt := range points {
		c.DrawCircle(pt[0], pt[1], radius, 16, color)
	}
	batch.skipStartEnd = false
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
	batch.Draw()
	c.end()
}

// works with convex + concave (non-self-intersacting) shapes
//
// multiple shapes can be separated by a [NaN, NaN]
func (c *Camera) DrawShapes(color uint, points ...[2]float32) {
	c.begin()

	var flatPoints, ptsCountsPerShape = separateShapes(points)
	var offset = 0
	var renderColor = getColor(color)

	c.Effects.updateUniforms(1, 1, nil, nil, false)

	for _, count := range ptsCountsPerShape {
		var shape = flatPoints[offset : offset+(count*2)]
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
	drawTexture.AssetId, drawTexture.Tint = assetId, color
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
