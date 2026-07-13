package graphics

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/execution/condition"
	geometry "pure-game-kit/packages/geometry"
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/angle"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/debug"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/point"
	"strconv"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (v *View) DrawColor(color uint) {
	obj.X, obj.Y, obj.Roundness, obj.Angle, obj.Effects.Tint, obj.Effects.FillColor = v.X, v.Y, 0, v.Angle, palette.White, color
	obj.TextFontId, obj.Text, obj.Effects.TextLineHeight, obj.Effects.TextColor, obj.ImageId = 0, "", 0, 0, 0
	obj.Width, obj.Height = v.Size()
	v.DrawObject(obj)
}
func (v *View) DrawGrid(thickness, spacingX, spacingY float32, color uint) {
	if spacingX*v.Zoom < 1 && spacingY*v.Zoom < 1 {
		return // way too dense — skip
	}

	var minX, minY, w, h = v.Bounds()
	var maxX, maxY = minX + w, minY + h

	var left, right = number.RoundDown(minX/spacingX) * spacingX, number.RoundUp(maxX/spacingX) * spacingX
	var top, bottom = number.RoundDown(minY/spacingY) * spacingY, number.RoundUp(maxY/spacingY) * spacingY

	for x := left; x <= right; x += spacingX { // vertical
		var t = thickness
		if number.DivisionRemainder(x, spacingX*10) == 0 {
			t *= 3
		}
		v.DrawShape(x, (top+bottom)/2, bottom-top, t, 90, 1, color, geometry.Area{})
	}

	for y := top; y <= bottom; y += spacingY { // horizontal
		var t = thickness
		if number.DivisionRemainder(y, spacingY*10) == 0 {
			t *= 3
		}
		v.DrawShape((left+right)/2, y, right-left, t, 0, 1, color, geometry.Area{})
	}

	if top <= 0 && bottom >= 0 {
		v.DrawShape((left+right)/2, 0, right-left, thickness*6, 0, 1, color, geometry.Area{})
	}
	if left <= 0 && right >= 0 {
		v.DrawShape(0, (top+bottom)/2, bottom-top, thickness*6, 90, 1, color, geometry.Area{})
	}
}
func (v *View) DrawPath(points []float32, thickness float32, color uint, mask geometry.Area) {
	if len(points) < 4 {
		return
	}

	for i := 2; i < len(points); i += 2 {
		var x1, y1, x2, y2 = points[i-2], points[i-1], points[i], points[i+1]
		var isNaN = number.IsNaN(x1) || number.IsNaN(y1) || number.IsNaN(x2) || number.IsNaN(y2)
		var isPoint = x1 == x2 && y1 == y2
		if isNaN || isPoint {
			continue
		}
		var dist = point.DistanceToPoint(x1, y1, x2, y2)
		var midX, midY, ang = (x1 + x2) / 2, (y1 + y2) / 2, angle.BetweenPoints(x1, y1, x2, y2)
		if i == 2 {
			v.DrawShape(x1, y1, thickness, thickness, 0, 1, color, mask)
		}
		v.DrawShape(x2, y2, thickness, thickness, 0, 1, color, mask)
		v.DrawShape(midX, midY, dist, thickness, ang, 0, color, mask)
	}
}
func (v *View) DrawShape(x, y, width, height, angle, roundness float32, color uint, mask geometry.Area) {
	obj.X, obj.Y, obj.Width, obj.Height, obj.Roundness = x, y, width, height, roundness
	obj.Angle, obj.ImageId, obj.Effects.Tint, obj.Effects.FillColor = angle, 0, palette.White, color
	obj.TextFontId, obj.Text, obj.Effects.TextLineHeight, obj.Effects.TextColor, obj.Mask = 0, "", 0, 0, mask
	v.DrawObject(obj)
}
func (v *View) DrawImage(x, y, width, height, angle float32, imageId assets.ImageId, tint uint, mask geometry.Area) {
	obj.X, obj.Y, obj.Width, obj.Height, obj.Roundness = x, y, width, height, 0
	obj.Angle, obj.ImageId, obj.Effects.Tint, obj.Effects.FillColor, obj.Mask = angle, 0, tint, 0, mask
	obj.TextFontId, obj.Text, obj.Effects.TextLineHeight, obj.Effects.TextColor, obj.ImageId = 0, "", 0, 0, imageId
	v.DrawObject(obj)
}
func (v *View) DrawText(text string, x, y, lineHeight float32, fontId assets.FontId, color uint, mask geometry.Area) {
	obj.Effects, obj.Mask = Effects(internal.DefaultEffects), mask
	obj.Width, obj.Height, obj.Effects.FillColor, obj.Roundness = 9999, 9999, 0, 0
	obj.Angle, obj.ImageId, obj.Effects.Tint, obj.Effects.FillColor = v.Angle, 0, palette.White, 0
	obj.TextFontId, obj.Text, obj.Effects.TextLineHeight, obj.Effects.TextColor = fontId, text, lineHeight/v.Zoom, color

	x, y = point.MoveAtAngle(x, y, obj.Angle, obj.Width/2)
	obj.X, obj.Y = point.MoveAtAngle(x, y, obj.Angle+90, obj.Height/2)
	v.DrawObject(obj)
}
func (v *View) DrawObject(object *Object) {
	internal.ViewArea = internal.Area(v.windowArea())
	internal.ViewX, internal.ViewY, internal.ViewZoom, internal.ViewAngle = v.X, v.Y, v.Zoom, v.Angle

	var o = object
	if o == nil || !v.IsAreaVisible(o.Bounds()) || (o.Width == 0 && o.Height == 0) {
		return
	}

	if o.Effects.TextBatch && o.textBatches != nil { // use cache only for batched textboxes
		internal.ReadyBatches = append(internal.ReadyBatches, o.textBatches...)
		return
	}

	var prevImageId = o.ImageId
	if o.Text != "" {
		if o.Effects.TextBatch {
			internal.IsRecording = true
			internal.CurrentBatchRecord = make([]*internal.Batch, 0)
		}
		if prevImageId == 0 { // shapes can use any texture but any non-0 TextFontId will break the batch, so force it
			o.ImageId = assets.ImageId(o.TextFontId)
		}
	}

	var mask = internal.Area(o.Mask)
	if o.Mask != (geometry.Area{}) {
		mask.X += float32(internal.WindowWidth) / 2
		mask.Y += float32(internal.WindowHeight) / 2
	}

	var eff = (*internal.Effects)(&o.Effects)

	if o.TileLayerId != 0 {
		var layer = internal.TileLayers[uint8(o.TileLayerId)]
		if layer.Image != nil {
			var tex = internal.Images[layer.ImageId]
			var src = rl.NewRectangle(tex.CropX, tex.CropY, tex.CropWidth, tex.CropHeight)
			var dst = rl.NewRectangle(o.X-o.Width/2, o.Y-o.Height/2, o.Width, o.Height)
			var cols, rows, sz = uint16(layer.Image.Width), uint16(layer.Image.Height), uint8(layer.TileSize)
			internal.Queue(tex.Texture, layer.Texture, src, dst, o.Angle, 0, mask, eff, 3, sz, cols, rows)
		}
		return
	}
	if o.ImageId != 0 {
		var tex = internal.Images[int32(o.ImageId)]
		if tex.Top != 0 || tex.Left != 0 || tex.Right != 0 || tex.Bottom != 0 {
			v.queueNinePatch(o.X, o.Y, o.Width, o.Height, o.Angle, o.Roundness, int32(o.ImageId), eff, mask)
		} else {
			v.queueQuad(o.X, o.Y, o.Width, o.Height, o.Angle, o.Roundness, int32(o.ImageId), o.ImageCrop, eff, mask)
		}
	} else {
		v.queueQuad(o.X, o.Y, o.Width, o.Height, o.Angle, o.Roundness, int32(o.ImageId), o.ImageCrop, eff, mask)
	}

	if o.Text != "" || o.Effects.TextIsInput { // empty input text needs a cursor position
		v.queueText(o, mask)
		if o.Effects.TextBatch {
			internal.CloseBatch()
			o.textBatches = internal.CurrentBatchRecord
			internal.ReadyBatches = append(internal.ReadyBatches, internal.CurrentBatchRecord...)
			internal.IsRecording = false
			for _, b := range internal.CurrentBatchRecord {
				b.IsMeshDirty = true
			}
		}
	}
	o.ImageId = prevImageId
}
func (v *View) DrawDebugInfo(detailed bool) {
	if condition.TrueEvery(0.2, 0xdeadc0de) {
		v.debugBuffer = v.debugBuffer[:0]
		v.debugBuffer = appendThousands(v.debugBuffer, uint64(internal.FPS))
		v.debugBuffer = append(v.debugBuffer, " FPS\n"...)

		if detailed {
			var targetFPS = internal.WindowTargetFPS
			if internal.WindowVsync {
				targetFPS = byte(rl.GetMonitorRefreshRate(rl.GetCurrentMonitor()))
			}
			if targetFPS == 0 {
				v.debugBuffer = append(v.debugBuffer, "unlimited target FPS"...)
			} else {
				v.debugBuffer = strconv.AppendInt(v.debugBuffer, int64(targetFPS), 10)
				v.debugBuffer = append(v.debugBuffer, " target FPS"...)
			}
			if internal.WindowVsync {
				v.debugBuffer = append(v.debugBuffer, " = monitor hz | vsync on\n"...)
			} else {
				v.debugBuffer = append(v.debugBuffer, " | vsync off\n"...)
			}

			var frameTargetMicroSec = 1000000.0 / float64(targetFPS)
			v.debugBuffer = strconv.AppendFloat(v.debugBuffer, float64(internal.EngineBusyMicroSec)/1000, 'f', 3, 32)
			v.debugBuffer = append(v.debugBuffer, "ms"...)
			if targetFPS > 0 {
				v.debugBuffer = append(v.debugBuffer, " ("...)
				var percent = (float64(internal.EngineBusyMicroSec) / frameTargetMicroSec) * 100
				v.debugBuffer = strconv.AppendFloat(v.debugBuffer, percent, 'f', 0, 32)
				v.debugBuffer = append(v.debugBuffer, "%)"...)
			}
			v.debugBuffer = append(v.debugBuffer, " engine busy (draw + internal + idle)"...)
			v.debugBuffer = append(v.debugBuffer, "\n"...)
			v.debugBuffer = strconv.AppendFloat(v.debugBuffer, float64(internal.GameBusyMicroSec)/1000, 'f', 3, 32)
			v.debugBuffer = append(v.debugBuffer, "ms"...)
			if targetFPS > 0 {
				v.debugBuffer = append(v.debugBuffer, " ("...)
				var percent = (float64(internal.GameBusyMicroSec) / frameTargetMicroSec) * 100
				v.debugBuffer = strconv.AppendFloat(v.debugBuffer, percent, 'f', 0, 32)
				v.debugBuffer = append(v.debugBuffer, "%)"...)
			}
			v.debugBuffer = append(v.debugBuffer, " game busy"...)
			v.debugBuffer = append(v.debugBuffer, "\n"...)

			v.debugBuffer = appendThousands(v.debugBuffer, uint64(internal.DrawCalls))
			v.debugBuffer = append(v.debugBuffer, " draw calls | "...)
			v.debugBuffer = appendThousands(v.debugBuffer, uint64(internal.QuadQueues))
			v.debugBuffer = append(v.debugBuffer, " quads\n"...)

			v.debugBuffer = strconv.AppendInt(v.debugBuffer, int64(internal.Runtime/3600), 10)
			v.debugBuffer = append(v.debugBuffer, "h "...)
			v.debugBuffer = strconv.AppendInt(v.debugBuffer, int64(internal.Runtime/60), 10)
			v.debugBuffer = append(v.debugBuffer, "m "...)
			v.debugBuffer = strconv.AppendInt(v.debugBuffer, int64(internal.Runtime)%60, 10)
			v.debugBuffer = append(v.debugBuffer, "s runtime\n\n"...)

			v.debugBuffer = appendThousands(v.debugBuffer, uint64(internal.NextImageId+1))
			v.debugBuffer = append(v.debugBuffer, " images | "...)
			v.debugBuffer = appendThousands(v.debugBuffer, uint64(-internal.NextImageCropId))
			v.debugBuffer = append(v.debugBuffer, " crops\n"...)
			v.debugBuffer = appendThousands(v.debugBuffer, uint64(len(internal.Fonts)))
			v.debugBuffer = append(v.debugBuffer, " fonts | "...)
			v.debugBuffer = appendThousands(v.debugBuffer, uint64(len(internal.Translations)))
			v.debugBuffer = append(v.debugBuffer, " translations\n"...)
			v.debugBuffer = appendThousands(v.debugBuffer, uint64(len(internal.Sounds)))
			v.debugBuffer = append(v.debugBuffer, " sounds | "...)
			v.debugBuffer = appendThousands(v.debugBuffer, uint64(len(internal.Music)))
			v.debugBuffer = append(v.debugBuffer, " music\n"...)
			v.debugBuffer = appendThousands(v.debugBuffer, uint64(len(internal.TileLayers)))
			v.debugBuffer = append(v.debugBuffer, " tile layers\n"...)
			v.debugBuffer = appendThousands(v.debugBuffer, uint64(len(internal.Layouts)))
			v.debugBuffer = append(v.debugBuffer, " GUI layouts\n"...)
		}
	}

	const size float32 = 30
	var tlx, tly = v.PointFromScreen(5, 5)
	var x, y = point.MoveAtAngle(tlx, tly, v.Angle+90, (size*15)/v.Zoom)
	var str = unsafe.String(unsafe.SliceData(v.debugBuffer), len(v.debugBuffer))
	v.DrawText(str, tlx, tly, size, 0, palette.White, geometry.Area{})

	if detailed {
		v.DrawText(debug.MemoryUsage(), x, y, size, 0, palette.White, geometry.Area{})
	}
	// rl.DrawText(str, 10, 15, 32, rl.White)
	// rl.DrawText(debug.MemoryUsage(), 10, 400, 32, rl.White)

}

// private ========================================================

type line struct {
	start, end int
	width      float32
}

var lines []line
var obj = &Object{} // used for primitive drawing
var colors = map[rune]uint{'⬜': palette.White, '⬛': palette.Black, '🟥': palette.Red, '🟧': palette.Orange,
	'🟨': palette.Yellow, '🟩': palette.Green, '🟦': palette.Blue, '🟪': palette.Purple, '🟫': palette.Brown}
var outlineColors = map[rune]uint{'⚪': palette.White, '⚫': palette.Black, '🔴': palette.Red, '🟠': palette.Orange,
	'🟡': palette.Yellow, '🟢': palette.Green, '🔵': palette.Blue, '🟣': palette.Purple, '🟤': palette.Brown}
var shades = map[rune]float32{'🌑': -0.8, '🌒': -0.6, '🌓': -0.4, '🌔': -0.2, '🌘': 0.2, '🌗': 0.4, '🌖': 0.6, '🌕': 0.8}
var weights = map[rune]int8{'⏬': -100, '🔽': -50, '🔁': 0, '🔼': 50, '⏫': 100}
var sizes = map[rune]float32{'🔇': 0.5, '🔈': 0.75, '🔉': 1.0, '🔊': 1.25, '📢': 1.5}

func (v *View) queueText(o *Object, mask internal.Area) {
	var eff = (*internal.Effects)(&o.Effects)
	var scale = eff.TextLineHeight / 255
	var gapX, gapY = eff.TextSymbolGap * scale, eff.TextLineGap * scale
	var fontData, has = internal.Fonts[uint8(o.TextFontId)]
	if !has { // fallback
		fontData = internal.Fonts[0]
	}
	var atlasTex = internal.Images[fontData.AtlasId].Texture
	var sin, cos = internal.SinCos(o.Angle)
	var shadeCol, shadeOutCol, contentWidth, contentHeight float32
	var w, h = o.Width, o.Height
	var txt = o.Text
	if eff.TextIsInput && txt == "" {
		txt = " " // empty input string should have a cursor position
	}

	lines = lines[:0]
	o.textCursorPos = o.textCursorPos[:0]
	var currentLineHeight = eff.TextLineHeight
	for i := 0; i < len(txt); {
		var lineStart = i
		var end, width, endHeight = o.measureLine(i, currentLineHeight)
		lines = append(lines, line{lineStart, end, width})
		contentWidth, contentHeight = max(contentWidth, width), contentHeight+endHeight*fontData.LineHeight
		currentLineHeight = endHeight
		i = end
		if i < len(txt) {
			contentHeight += gapY
			i++ // skip newline
		}
	}
	var y = o.Y - h/2 - fontData.Ascender*eff.TextLineHeight + eff.TextAlignY*(h-contentHeight)

	var baseLineHeight = eff.TextLineHeight

	for _, ln := range lines {
		var x = (o.X - w/2) + eff.TextAlignX*(w-ln.width)
		var prevGlyph internal.Glyph

		for _, r := range txt[ln.start:ln.end] {
			if !o.Effects.TextIsInput && o.embedEffect(r, eff, &shadeCol, &shadeOutCol, baseLineHeight) {
				continue // tag symbol applies to effects and gets skipped
			}

			var glyph = fontData.Chars[r]
			var kerning, _ = prevGlyph.Kernings[r]
			x += kerning * eff.TextLineHeight

			var src, dst = getGlyphSrcDst(o, r, glyph, x, y, cos, sin, 0)

			if o.Effects.TextIsInput {
				var midX, midY = (x - gapX/2) - o.X, (y + (fontData.Ascender+fontData.Descender)*eff.TextLineHeight/2) - o.Y
				midX, midY = o.X+midX*cos-midY*sin, o.Y+midX*sin+midY*cos
				o.textCursorPos = append(o.textCursorPos, midX)
			}

			if glyph.EmbededImageId != 0 {
				if dst.Width <= 0 { // fully clipped by the textbox
					x += glyph.Advance*eff.TextLineHeight + gapX
					prevGlyph = glyph
					continue
				}
				var prevFill, prevOut = eff.FillColor, eff.OutlineColor
				var x, y = dst.X + dst.Width/2, dst.Y - dst.Height/2
				var area = geometry.NewArea(src.X+src.Width/2, src.Y+src.Height/2, src.Width, -src.Height)
				eff.FillColor, eff.OutlineColor = 0, 0
				v.queueQuad(x, y, dst.Width, dst.Height, o.Angle, 0, glyph.EmbededImageId, area, eff, mask)
				eff.FillColor, eff.OutlineColor = prevFill, prevOut
			} else {
				if r != ' ' && r != '\n' {
					internal.Queue(atlasTex, rl.Texture2D{}, src, dst, o.Angle, 0, mask, eff, internal.KindText, 0, 0, 0)
				}
				if eff.TextUnderline {
					var underscore = fontData.Chars[internal.Underline]
					var src2, dst2 = getGlyphSrcDst(o, internal.Underline, underscore, x, y, cos, sin, dst.Width)
					internal.Queue(atlasTex, rl.Texture2D{}, src2, dst2, o.Angle, 0, mask, eff, internal.KindText, 0, 0, 0)
				}
				if eff.TextCrossout {
					var dash = fontData.Chars[internal.Crossout]
					var src2, dst2 = getGlyphSrcDst(o, internal.Crossout, dash, x, y, cos, sin, dst.Width)
					internal.Queue(atlasTex, rl.Texture2D{}, src2, dst2, o.Angle, 0, mask, eff, internal.KindText, 0, 0, 0)
				}
			}
			x += glyph.Advance*eff.TextLineHeight + gapX
			prevGlyph = glyph
		}

		if o.Effects.TextIsInput {
			var endX, endY = (x - gapX/2) - o.X, (y + (fontData.Ascender+fontData.Descender)*eff.TextLineHeight/2) - o.Y
			endX, endY = o.X+endX*cos-endY*sin, o.Y+endX*sin+endY*cos
			o.textCursorPos = append(o.textCursorPos, endX)
		}
		y += eff.TextLineHeight*fontData.LineHeight + gapY
	}
	o.textWidth, o.textHeight = contentWidth, contentHeight
}
func (v *View) queueQuad(x, y, w, h, a, r float32, imageId int32, crop geometry.Area, eff *internal.Effects, mask internal.Area) {
	var tex = internal.Images[imageId]
	var prevFill = eff.FillColor
	var kind uint8
	if imageId == 0 || tex.Texture.Width == 0 {
		imageId = 0 // fallback to default texture
		tex = internal.Images[imageId]
	} else {
		kind = internal.KindSprite
	}
	if crop == (geometry.Area{}) {
		crop = geometry.NewArea(tex.CropX+tex.CropWidth/2, tex.CropY+tex.CropHeight/2, tex.CropWidth, tex.CropHeight)
	}
	var src = rl.NewRectangle(crop.X-crop.Width/2, crop.Y-crop.Height/2, crop.Width, crop.Height)
	var dst = rl.NewRectangle(x-w/2, y-h/2, w, h)
	internal.Queue(tex.Texture, rl.Texture2D{}, src, dst, a, r, mask, eff, kind, 0, 0, 0)
	eff.FillColor = prevFill
}
func (v *View) queueNinePatch(x, y, w, h, a, r float32, imageId int32, eff *internal.Effects, mask internal.Area) {
	var img = internal.Images[imageId]
	var tw, th, top, left, right, bottom = img.CropWidth, img.CropHeight, img.Top, img.Left, img.Right, img.Bottom
	var sx = [4]float32{img.CropX, img.CropX + left, img.CropX + tw - right, img.CropX + tw}
	var sy = [4]float32{img.CropY, img.CropY + top, img.CropY + th - bottom, img.CropY + th}
	if w < 0 { // flip negative width
		w, sx[0], sx[1], sx[2], sx[3] = -w, sx[3], sx[2], sx[1], sx[0]
	}
	if h < 0 { // flip negative height
		h, sy[0], sy[1], sy[2], sy[3] = -h, sy[3], sy[2], sy[1], sy[0]
	}
	if w <= left+right && left+right > 0 { // scale down edges when there isn't enough space to fit them at natural size
		var s = w / (left + right)
		left, right = left*s, right*s
	}
	if h <= top+bottom && top+bottom > 0 {
		var s = h / (top + bottom)
		top, bottom = top*s, bottom*s
	}
	var dx = [4]float32{x - w/2, x - w/2 + left, x + w/2 - right, x + w/2}
	var dy = [4]float32{y - h/2, y - h/2 + top, y + h/2 - bottom, y + h/2}

	for j := range 3 {
		for i := range 3 {
			var qw, qh, su, sv = dx[i+1] - dx[i], dy[j+1] - dy[j], sx[i+1] - sx[i], sy[j+1] - sy[j]
			if qw > 0 && qh > 0 && su != 0 && sv != 0 {
				var src, dst = rl.NewRectangle(sx[i], sy[j], su, sv), rl.NewRectangle(dx[i], dy[j], qw, qh)
				internal.Queue(img.Texture, rl.Texture2D{}, src, dst, a, r, mask, eff, internal.KindSprite, 0, 0, 0)
			}
		}
	}
}

func (v *View) windowArea() geometry.Area {
	if v.WindowArea == (geometry.Area{}) {
		return geometry.NewArea(internal.WindowWidth/2, internal.WindowHeight/2, internal.WindowWidth, internal.WindowHeight)
	}
	return v.WindowArea
}
func getGlyphSrcDst(o *Object, r rune, glyph internal.Glyph, x, y, cos, sin, newWidth float32) (src, dst rl.Rectangle) {
	var symbol = o.TextFontId.SymbolArea(r, o.Effects.TextLineHeight)
	var offsetX, offsetY, dstW, dstH = symbol.X, symbol.Y, symbol.Width, symbol.Height
	if newWidth != 0 {
		dstW = newWidth
	}
	var atlas, dstX, dstY = glyph.AtlasBounds, x + offsetX, y + offsetY
	var srcX, srcY, srcW, srcH = atlas.Left, atlas.Top, atlas.Right - atlas.Left, atlas.Bottom - atlas.Top
	var left, top, right, bot = o.X - o.Width/2, o.Y - o.Height/2, o.X + o.Width/2, o.Y + o.Height/2
	var clipL, clipR, clipT, clipB = max(dstX, left), min(dstX+dstW, right), max(dstY, top), min((dstY - dstH), bot)
	if clipL >= clipR || clipT >= clipB {
		return rl.Rectangle{}, rl.Rectangle{}
	}

	var clippedW, clippedH, origH = clipR - clipL, clipB - clipT, (dstY - dstH) - dstY
	var dx, dy = (clipL + clippedW/2) - o.X, ((clipT + clipB) / 2) - o.Y
	srcX, srcY = srcX+((clipL-dstX)/dstW*srcW), srcY+((clipT-dstY)/origH*srcH)
	srcW, srcH = srcW*(clippedW/dstW), srcH*(clippedH/origH)
	dstW, dstH = clippedW, -clippedH
	dstX, dstY = o.X+dx*cos-dy*sin-dstW/2, o.Y+dx*sin+dy*cos+dstH/2
	return rl.NewRectangle(srcX, srcY, srcW, srcH), rl.NewRectangle(dstX, dstY, dstW, dstH)
}

func appendThousands(buf []byte, n uint64) []byte {
	var tmp [32]byte
	var s = strconv.AppendUint(tmp[:0], n, 10)
	var length = len(s)
	for i, c := range s {
		if i > 0 && (length-i)%3 == 0 {
			buf = append(buf, ' ')
		}
		buf = append(buf, c)
	}
	return buf
}
