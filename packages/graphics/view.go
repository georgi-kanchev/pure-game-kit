package graphics

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/input/mouse"
	"pure-game-kit/packages/input/mouse/button"
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/number"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type View struct {
	X, Y, Zoom, Angle float32

	WindowArea Area // The draw area in window space. Zero value = entire window.
	Mask       Area // In view space. Everything drawn outside of it is cropped. Zero value = no masking.

	//=================================================================

	velocityX, velocityY, dragVelX, dragVelY float32
}

func NewView(zoom float32) View {
	return View{Zoom: zoom}
}

// =================================================================

func (v *View) MouseDragAndZoom() {
	var oldZoom = v.Zoom
	var scroll = mouse.Scroll()

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
	var oldZoom = v.Zoom
	var scroll = mouse.ScrollSmooth()

	if scroll != 0 {
		v.Zoom *= 1 + 0.001*scroll
		var mx, my = v.MousePosition()

		v.X += (mx - v.X) * (v.Zoom/oldZoom - 1)
		v.Y += (my - v.Y) * (v.Zoom/oldZoom - 1)
	}

	const dragFriction, dragStrength = 8.0, 8.0
	var dt = internal.FrameDelta
	var decay = number.Exponential(-dragFriction * dt)
	v.velocityX *= decay
	v.velocityY *= decay

	if mouse.IsButtonPressed(button.Middle) {
		var sin, cos = internal.SinCos(-v.Angle)
		var dx, dy = mouse.CursorDelta()

		dx /= dt
		dy /= dt

		var wdx = (dx*cos - dy*sin) / v.Zoom
		var wdy = (dx*sin + dy*cos) / v.Zoom

		v.velocityX -= wdx * dragStrength * dt
		v.velocityY -= wdy * dragStrength * dt
	}

	v.X += v.velocityX * dt
	v.Y += v.velocityY * dt

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
	var sx2, sy2 = v.PointToScreen(x+width, y+height)
	var sMinX, sMaxX = min(sx1, sx2), max(sx1, sx2)
	var sMinY, sMaxY = min(sy1, sy2), max(sy1, sy2)
	var mx, my, mw, mh = v.area()
	return sMaxX > mx && sMinX < mx+mw && sMaxY > my && sMinY < my+mh
}
func (v *View) IsHovered() bool {
	var mx, my = internal.MouseX, internal.MouseY
	var sx, sy, sw, sh = v.area()
	return mx > sx && my > sy && mx < sx+sw && my < sy+sh
}
func (v *View) MousePosition() (x, y float32) {
	return v.PointFromScreen(internal.MouseX, internal.MouseY)
}
func (v *View) Size() (width, height float32) {
	var _, _, sw, sh = v.area()
	return sw / v.Zoom, sh / v.Zoom
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
	var sx = screenX - float32(internal.WindowWidth/2)
	var sy = screenY - float32(internal.WindowHeight/2)
	return sx, sy
}
func (v *View) PointToScreen(x, y float32) (screenX, screenY float32) {
	var vx = x + float32(internal.WindowWidth/2)
	var vy = y + float32(internal.WindowHeight/2)
	return vx, vy
}
func (v *View) PointFromView(otherView *View, otherX, otherY float32) (myX, myY float32) {
	var screenX, screenY = otherView.PointToScreen(otherX, otherY)
	return v.PointFromScreen(screenX, screenY)
}
func (v *View) PointToView(otherView *View, myX, myY float32) (otherX, otherY float32) {
	return otherView.PointFromView(v, myX, myY)
}
func (v *View) PointFromEdge(edgeX, edgeY float32) (x, y float32) {
	var sx, sy, sw, sh = v.area()
	var scrX, scrY = sx + sw*edgeX, sy + sh*edgeY
	return v.PointFromScreen(scrX, scrY)
}

//=================================================================

func (v *View) DrawImage(x, y, width, height, angle float32, imageId assets.ImageId, tint uint) {
	object.X, object.Y, object.Roundness = x, y, 0
	object.Width, object.Height = width, height
	object.Angle, object.ImageId = angle, imageId
	object.Effects.Tint, object.Effects.FillColor = tint, 0
	v.DrawObjects(object)
}
func (v *View) DrawText(x, y, lineHeight float32, fontId assets.FontId, color uint, text string) {
	object.Effects = Effects(internal.DefaultEffects)
	object.Text, object.Effects.FillColor, object.Roundness = text, 0, 0
	object.TextFontId, object.Effects.TextLineHeight, object.Angle = fontId, lineHeight, 0
	object.X, object.Y, object.Width, object.Height = x+object.Width/2, y+object.Height/2, 99999, 99999
	v.DrawObjects(object)
}
func (v *View) DrawObjects(objects ...*Object) {
	for _, o := range objects {
		if o == nil || !v.IsAreaVisible(o.Bounds()) {
			continue
		}

		if o.TextBatch && o.textBatches != nil { // use cache only for batched textboxes
			internal.ReadyBatches = append(internal.ReadyBatches, o.textBatches...)
			continue
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
		v.queueShapeOrSprite(o.X, o.Y, o.Width, o.Height, o.Angle, o.Roundness, int32(o.ImageId), o.ImageCrop, eff, mask)

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
}

// private ========================================================

type line struct {
	start, end int
	width      float32
}

var lines []line
var object = &Object{} // used for primitive drawing
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
				v.queueShapeOrSprite(x, y, dst.Width, dst.Height, o.Angle, 0, glyph.EmbededImageId, area, eff, mask)
				eff.FillColor, eff.OutlineColor = prevFill, prevOut
			} else {
				if r != ' ' && r != '\n' {
					internal.Queue(atlasTex, src, dst, o.Angle, 0, mask, eff, internal.KindText)
				}
				if eff.TextUnderline {
					var src2, dst2 = getGlyphSrcDst(o, internal.Underline, fontData.Chars[internal.Underline], x, y, cos, sin, dst.Width)
					internal.Queue(atlasTex, src2, dst2, o.Angle, 0, mask, eff, internal.KindText)
				}
				if eff.TextCrossout {
					var src2, dst2 = getGlyphSrcDst(o, internal.Crossout, fontData.Chars[internal.Crossout], x, y, cos, sin, dst.Width)
					internal.Queue(atlasTex, src2, dst2, o.Angle, 0, mask, eff, internal.KindText)
				}
			}
			x += glyph.Advance*eff.TextLineHeight + gapX
			prevGlyph = glyph
		}

		y += eff.TextLineHeight*fontData.LineHeight + gapY
	}
}
func (v *View) queueShapeOrSprite(x, y, w, h, a, r float32, imageId int32, crop Area, eff *internal.Effects, mask internal.Area) {
	var tex = internal.Images[imageId]
	var prevFill = eff.FillColor
	var kind uint8
	if imageId == 0 || tex.Texture.Width == 0 {
		imageId = 0 // fallback to default texture
		tex = internal.Images[imageId]
		if eff.TextLineHeight == 0 && eff.TextColor == 0 {
			eff.FillColor = eff.Tint
		}
	} else {
		kind = internal.KindSprite
	}
	if crop == (Area{}) {
		crop = NewArea(tex.CropX, tex.CropY, tex.CropWidth, tex.CropHeight)
	}
	var src = rl.NewRectangle(crop.X, crop.Y, crop.Width, crop.Height)
	var dst = rl.NewRectangle(x-w/2, y-h/2, w, h)
	internal.Queue(tex.Texture, src, dst, a, r, mask, eff, kind)
	eff.FillColor = prevFill
}

func (v *View) area() (x, y, w, h float32) {
	if v.WindowArea == (Area{}) {
		return 0, 0, float32(internal.WindowWidth), float32(internal.WindowHeight)
	}
	return v.WindowArea.X, v.WindowArea.Y, v.WindowArea.Width, v.WindowArea.Height
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
