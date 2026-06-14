package graphics

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/execution/condition"
	"pure-game-kit/packages/input/mouse"
	"pure-game-kit/packages/input/mouse/button"
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

type View struct {
	X, Y, Zoom, Angle float32

	WindowArea Area // The drawing area in window space. Zero value = entire window.

	//=================================================================

	velocityX, velocityY, dragVelX, dragVelY float32

	debugBuffer []byte
}

func NewView(zoom float32) View { return View{Zoom: zoom} }

// =================================================================

func (v *View) MouseDragAndZoom() {
	var oldZoom, scroll = v.Zoom, mouse.Scroll()

	if scroll != 0 {
		v.Zoom *= 1 + 0.05*scroll

		var mx, my = v.MousePosition()
		v.X += (mx - v.X) * (v.Zoom/oldZoom - 1)
		v.Y += (my - v.Y) * (v.Zoom/oldZoom - 1)
	}

	if mouse.IsButtonPressed(button.Middle) {
		var dx, dy = mouse.CursorDelta()
		var sin, cos = internal.SinCos(-v.Angle)
		v.X -= (dx*cos - dy*sin) / v.Zoom
		v.Y -= (dx*sin + dy*cos) / v.Zoom
	}
}
func (v *View) MouseDragAndZoomSmoothly() {
	var oldZoom, scroll = v.Zoom, mouse.ScrollSmooth()

	if scroll != 0 {
		v.Zoom *= 1 + 0.001*scroll

		var mx, my = v.MousePosition()
		v.X += (mx - v.X) * (v.Zoom/oldZoom - 1)
		v.Y += (my - v.Y) * (v.Zoom/oldZoom - 1)
	}

	const dragFriction, dragStrength = 8.0, 8.0
	var dt = internal.FrameDelta
	var decay = number.Exponential(-dragFriction * dt)
	v.velocityX, v.velocityY = v.velocityX*decay, v.velocityY*decay

	if mouse.IsButtonPressed(button.Middle) {
		var sin, cos = internal.SinCos(-v.Angle)
		var dx, dy = mouse.CursorDelta()

		dx, dy = dx/dt, dy/dt
		var wdx, wdy = (dx*cos - dy*sin) / v.Zoom, (dx*sin + dy*cos) / v.Zoom
		v.velocityX, v.velocityY = v.velocityX-(wdx*dragStrength*dt), v.velocityY-(wdy*dragStrength*dt)
	}

	v.X, v.Y = v.X+(v.velocityX*dt), v.Y+v.velocityY*dt

	if number.Absolute(v.velocityX) < 0.0001 {
		v.velocityX = 0
	}
	if number.Absolute(v.velocityY) < 0.0001 {
		v.velocityY = 0
	}
}

//=================================================================

func (v *View) IsAreaVisible(x, y, width, height float32) bool {
	var sx1, sy1 = v.PointToScreen(x, y)
	var sx2, sy2 = v.PointToScreen(x+width, y)
	var sx3, sy3 = v.PointToScreen(x, y+height)
	var sx4, sy4 = v.PointToScreen(x+width, y+height)
	var minX, maxX = min(min(sx1, sx2), min(sx3, sx4)), max(max(sx1, sx2), max(sx3, sx4))
	var minY, maxY = min(min(sy1, sy2), min(sy3, sy4)), max(max(sy1, sy2), max(sy3, sy4))
	return maxX >= 0 && minX <= internal.WindowWidth && maxY >= 0 && minY <= internal.WindowHeight
}
func (v *View) IsHovered() bool {
	return internal.MouseX > 0 && internal.MouseY > 0 && internal.MouseX < internal.WindowWidth && internal.MouseY < internal.WindowHeight
}
func (v *View) MousePosition() (x, y float32) {
	return v.PointFromScreen(internal.MouseX, internal.MouseY)
}
func (v *View) Size() (width, height float32) {
	var windowArea = v.windowArea()
	return windowArea.Width / v.Zoom, windowArea.Height / v.Zoom
}
func (v *View) Bounds() (x, y, width, height float32) {
	var x1, y1 = v.PointFromEdge(0, 0)
	var x2, y2 = v.PointFromEdge(1, 0)
	var x3, y3 = v.PointFromEdge(1, 1)
	var x4, y4 = v.PointFromEdge(0, 1)
	var minX, minY = number.Minimum(x1, x2, x3, x4), number.Minimum(y1, y2, y3, y4)
	var maxX, maxY = number.Maximum(x1, x2, x3, x4), number.Maximum(y1, y2, y3, y4)
	return minX, minY, maxX - minX, maxY - minY
}

func (v *View) PointFromScreen(screenX, screenY float32) (x, y float32) {
	var wa = v.windowArea()
	x, y = screenX-(wa.X+wa.Width/2), screenY-(wa.Y+wa.Height/2)
	if v.Zoom != 0 {
		x, y = x/v.Zoom, y/v.Zoom
	}
	if v.Angle != 0 {
		var sin, cos = internal.SinCos(v.Angle)
		x, y = x*cos-y*sin, x*sin+y*cos
	}
	return x + v.X, y + v.Y
}
func (v *View) PointToScreen(x, y float32) (screenX, screenY float32) {
	x, y = x-v.X, y-v.Y
	if v.Angle != 0 {
		var sin, cos = internal.SinCos(v.Angle)
		x, y = x*cos+y*sin, -x*sin+y*cos
	}
	if v.Zoom != 0 {
		x, y = x*v.Zoom, y*v.Zoom
	}
	var wa = v.windowArea()
	return x + (wa.X + wa.Width/2), y + (wa.Y + wa.Height/2)
}
func (v *View) PointFromView(otherView *View, otherX, otherY float32) (myX, myY float32) {
	return v.PointFromScreen(otherView.PointToScreen(otherX, otherY))
}
func (v *View) PointToView(otherView *View, myX, myY float32) (otherX, otherY float32) {
	return otherView.PointFromView(v, myX, myY)
}
func (v *View) PointFromEdge(edgeX, edgeY float32) (x, y float32) {
	var wa = v.windowArea()
	return v.PointFromScreen(wa.X+wa.Width*edgeX, wa.Y+wa.Height*edgeY)
}

//=================================================================

func (v *View) DrawColor(color uint) {
	obj.X, obj.Y, obj.Roundness, obj.Angle, obj.Effects.Tint, obj.Effects.FillColor = v.X, v.Y, 0, v.Angle, color, 0
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
		v.DrawShape(x, (top+bottom)/2, bottom-top, t, 90, 1, color, Area{})
	}

	for y := top; y <= bottom; y += spacingY { // horizontal
		var t = thickness
		if number.DivisionRemainder(y, spacingY*10) == 0 {
			t *= 3
		}
		v.DrawShape((left+right)/2, y, right-left, t, 0, 1, color, Area{})
	}

	if top <= 0 && bottom >= 0 {
		v.DrawShape((left+right)/2, 0, right-left, thickness*6, 0, 1, color, Area{})
	}
	if left <= 0 && right >= 0 {
		v.DrawShape(0, (top+bottom)/2, bottom-top, thickness*6, 90, 1, color, Area{})
	}
}
func (v *View) DrawShape(x, y, width, height, angle, roundness float32, color uint, mask Area) {
	obj.X, obj.Y, obj.Width, obj.Height, obj.Roundness = x, y, width, height, roundness
	obj.Angle, obj.ImageId, obj.Effects.Tint, obj.Effects.FillColor = angle, 0, palette.White, color
	obj.TextFontId, obj.Text, obj.Effects.TextLineHeight, obj.Effects.TextColor, obj.Mask = 0, "", 0, 0, mask
	v.DrawObject(obj)
}
func (v *View) DrawPath(points []float32, thickness float32, color uint, mask Area) {
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
func (v *View) DrawImage(x, y, width, height, angle float32, imageId assets.ImageId, tint uint, mask Area) {
	obj.X, obj.Y, obj.Width, obj.Height, obj.Roundness = x, y, width, height, 0
	obj.Angle, obj.ImageId, obj.Effects.Tint, obj.Effects.FillColor, obj.Mask = angle, 0, tint, 0, mask
	obj.TextFontId, obj.Text, obj.Effects.TextLineHeight, obj.Effects.TextColor, obj.ImageId = 0, "", 0, 0, imageId
	v.DrawObject(obj)
}
func (v *View) DrawText(x, y, lineHeight float32, fontId assets.FontId, color uint, text string, mask Area) {
	obj.Effects = Effects(internal.DefaultEffects)
	obj.Width, obj.Height, obj.Effects.FillColor, obj.Roundness = 9999, 9999, 0, 0
	obj.Angle, obj.ImageId, obj.Effects.Tint, obj.Effects.FillColor = v.Angle, 0, palette.White, 0
	obj.TextFontId, obj.Text, obj.Effects.TextLineHeight, obj.Effects.TextColor = fontId, text, lineHeight/v.Zoom, color
	obj.Mask = mask

	x, y = point.MoveAtAngle(x, y, obj.Angle, obj.Width/2)
	obj.X, obj.Y = point.MoveAtAngle(x, y, obj.Angle+90, obj.Height/2)
	v.DrawObject(obj)
}
func (v *View) DrawObject(object *Object) {
	internal.ViewArea = internal.Area(v.windowArea())
	internal.ViewX, internal.ViewY, internal.ViewZoom, internal.ViewAngle = v.X, v.Y, v.Zoom, v.Angle

	var o = object
	if o == nil || !v.IsAreaVisible(o.Bounds()) {
		return
	}

	if o.TextBatch && o.textBatches != nil { // use cache only for batched textboxes
		internal.ReadyBatches = append(internal.ReadyBatches, o.textBatches...)
		return
	}

	var prevImageId = o.ImageId
	if o.Text != "" {
		if o.TextBatch {
			internal.IsRecording = true
			internal.CurrentBatchRecord = make([]*internal.Batch, 0)
		}
		if prevImageId == 0 { // shapes can use any texture but any non-0 TextFontId will break the batch, so force it
			o.ImageId = assets.ImageId(o.TextFontId)
		}
	}

	var mask = internal.Area(o.Mask)
	if o.Mask != (Area{}) {
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
			var cols, rows = uint16(layer.Image.Width), uint16(layer.Image.Height)
			internal.Queue(tex.Texture, layer.Texture, src, dst, o.Angle, 0, mask, eff, 3, uint8(layer.TileSize), cols, rows)
		}
		return
	}
	v.queueQuad(o.X, o.Y, o.Width, o.Height, o.Angle, o.Roundness, int32(o.ImageId), o.ImageCrop, eff, mask)

	if o.Text != "" {
		v.queueText(o, mask)
		if o.TextBatch {
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
	var tlx, tly = v.PointFromScreen(10, 10)
	const size float32 = 30

	if condition.TrueEvery(0.1, 0xdeadc0de) {
		v.debugBuffer = v.debugBuffer[:0]
		v.debugBuffer = strconv.AppendFloat(v.debugBuffer, float64(internal.FPS), 'f', 0, 32)
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
			v.debugBuffer = append(v.debugBuffer, " draw calls\n"...)

			v.debugBuffer = strconv.AppendInt(v.debugBuffer, int64(internal.Runtime/3600), 10)
			v.debugBuffer = append(v.debugBuffer, "h "...)
			v.debugBuffer = strconv.AppendInt(v.debugBuffer, int64(internal.Runtime/60), 10)
			v.debugBuffer = append(v.debugBuffer, "m "...)
			v.debugBuffer = strconv.AppendInt(v.debugBuffer, int64(internal.Runtime)%60, 10)
			v.debugBuffer = append(v.debugBuffer, "s runtime\n\n"...)

			v.debugBuffer = appendThousands(v.debugBuffer, uint64(internal.NextImageId+1))
			v.debugBuffer = append(v.debugBuffer, " images | "...)
			v.debugBuffer = appendThousands(v.debugBuffer, uint64(internal.NextImageCropId))
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

	if detailed {
		v.DrawText(tlx, tly+(size*14)/v.Zoom, size, 0, palette.White, debug.MemoryUsage(), Area{})
	}
	var str = unsafe.String(unsafe.SliceData(v.debugBuffer), len(v.debugBuffer))
	v.DrawText(tlx, tly, size, 0, palette.White, str, Area{})
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
	var contentHeight float32
	var originalLineHeight = eff.TextLineHeight
	lines = lines[:0]
	var currentLineHeight = originalLineHeight
	for i := 0; i < len(o.Text); {
		var lineStart = i
		var end, width, endHeight = o.measureLine(i, currentLineHeight)
		lines = append(lines, line{lineStart, end, width})
		contentHeight += endHeight * fontData.LineHeight
		currentLineHeight = endHeight
		i = end
		if i < len(o.Text) {
			contentHeight += gapY
			i++ // skip newline
		}
	}
	var y = o.Y - o.Height/2 - fontData.Ascender*eff.TextLineHeight + eff.TextAlignY*(o.Height-contentHeight)

	for _, ln := range lines {
		var x = (o.X - o.Width/2) + eff.TextAlignX*(o.Width-ln.width)
		var prevGlyph internal.Glyph

		for _, r := range o.Text[ln.start:ln.end] {
			if embedEffect(r, eff, originalLineHeight) {
				continue // tag symbol applies to effects and gets skipped
			}

			var glyph = fontData.Chars[r]
			var kerning, _ = prevGlyph.Kernings[r]
			x += kerning * eff.TextLineHeight

			var src, dst = getGlyphSrcDst(o, r, glyph, x, y, cos, sin, 0)
			if glyph.EmbededImageId != 0 {
				if dst.Width <= 0 { // fully clipped by the textbox
					x += glyph.Advance*eff.TextLineHeight + gapX
					prevGlyph = glyph
					continue
				}
				var prevFill, prevOut = eff.FillColor, eff.OutlineColor
				var x, y = dst.X + dst.Width/2, dst.Y + dst.Height/2
				var area = NewArea(src.X, src.Y, src.Width, src.Height)
				eff.FillColor, eff.OutlineColor = 0, 0
				v.queueQuad(x, y, dst.Width, dst.Height, o.Angle, 0, glyph.EmbededImageId, area, eff, mask)
				eff.FillColor, eff.OutlineColor = prevFill, prevOut
			} else {
				if r != ' ' && r != '\n' {
					internal.Queue(atlasTex, rl.Texture2D{}, src, dst, o.Angle, 0, mask, eff, internal.KindText, 0, 0, 0)
				}
				if eff.TextUnderline {
					var src2, dst2 = getGlyphSrcDst(o, internal.Underline, fontData.Chars[internal.Underline], x, y, cos, sin, dst.Width)
					internal.Queue(atlasTex, rl.Texture2D{}, src2, dst2, o.Angle, 0, mask, eff, internal.KindText, 0, 0, 0)
				}
				if eff.TextCrossout {
					var src2, dst2 = getGlyphSrcDst(o, internal.Crossout, fontData.Chars[internal.Crossout], x, y, cos, sin, dst.Width)
					internal.Queue(atlasTex, rl.Texture2D{}, src2, dst2, o.Angle, 0, mask, eff, internal.KindText, 0, 0, 0)
				}
			}
			x += glyph.Advance*eff.TextLineHeight + gapX
			prevGlyph = glyph
		}

		y += eff.TextLineHeight*fontData.LineHeight + gapY
	}
}
func (v *View) queueQuad(x, y, w, h, a, r float32, imageId int32, crop Area, eff *internal.Effects, mask internal.Area) {
	var tex = internal.Images[imageId]
	var prevFill = eff.FillColor
	var kind uint8
	if imageId == 0 || tex.Texture.Width == 0 {
		imageId = 0 // fallback to default texture
		tex = internal.Images[imageId]
	} else {
		kind = internal.KindSprite
	}
	if crop == (Area{}) {
		crop = NewArea(tex.CropX, tex.CropY, tex.CropWidth, tex.CropHeight)
	}
	var src = rl.NewRectangle(crop.X, crop.Y, crop.Width, crop.Height)
	var dst = rl.NewRectangle(x-w/2, y-h/2, w, h)
	internal.Queue(tex.Texture, rl.Texture2D{}, src, dst, a, r, mask, eff, kind, 0, 0, 0)
	eff.FillColor = prevFill
}

func (v *View) windowArea() Area {
	if v.WindowArea == (Area{}) {
		return NewArea(0, 0, internal.WindowWidth, internal.WindowHeight)
	}
	return v.WindowArea
}
func getGlyphSrcDst(o *Object, r rune, glyph internal.Glyph, x, y, cos, sin, newWidth float32) (src, dst rl.Rectangle) {
	var offsetX, offsetY, dstW, dstH = o.TextFontId.SymbolArea(r, o.Effects.TextLineHeight)
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
func embedEffect(r rune, effect *internal.Effects, originalLineHeight float32) (success bool) {
	if r == '✅' {
		effect.TextUnderline = !effect.TextUnderline
		return true
	}
	if r == '❎' {
		effect.TextCrossout = !effect.TextCrossout
		return true
	}

	var color = colors[r]
	var outlineColor = outlineColors[r]
	var weight, hasWeights = weights[r]
	var size = sizes[r]
	if color != 0 {
		effect.TextColor = color
		return true
	} else if outlineColor != 0 {
		effect.OutlineColor = outlineColor
		return true
	} else if hasWeights {
		effect.TextWeight = weight
		return true
	} else if size != 0 {
		effect.TextLineHeight = originalLineHeight * size
		return true
	}

	return false
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
