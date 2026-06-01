package graphics

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/input/mouse"
	"pure-game-kit/packages/input/mouse/button"
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/number"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type View struct {
	X, Y, Zoom, Angle float32

	Area Area // The draw area in window space. Zero value = entire window.
	Mask Area // In view space. Everything drawn outside of it is cropped. Zero value = no masking.

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

func (v *View) DrawObjects(objects ...*Object) {
	for _, o := range objects {
		if o == nil || !v.IsAreaVisible(o.Bounds()) {
			continue
		}

		if o.textBatches != nil { // use cache if available
			internal.ReadyBatches = append(internal.ReadyBatches, o.textBatches...)
			continue
		}

		var textBatches []*internal.Batch
		var tbPtr *[]*internal.Batch
		var prevImageId = o.ImageId
		if o.Text != "" {
			textBatches = make([]*internal.Batch, 0)
			tbPtr = &textBatches
			if prevImageId == 0 { // shapes can use any texture but any non-0 TextFontId will break the batch, so force it
				o.ImageId = assets.ImageId(o.TextFontId)
			}
		}

		var tex = internal.Images[int32(o.ImageId)]
		var crop = o.ImageCrop
		if crop == (Area{}) {
			crop = NewArea(tex.CropX, tex.CropY, tex.CropWidth, tex.CropHeight)
		}
		var src = rl.NewRectangle(crop.X, crop.Y, crop.Width, crop.Height)
		var dst = rl.NewRectangle(o.X-o.Width/2, o.Y-o.Height/2, o.Width, o.Height)
		var mask = internal.Area(o.Mask)
		if o.Mask != (Area{}) {
			mask.X += float32(internal.WindowWidth) / 2
			mask.Y += float32(internal.WindowHeight) / 2
		}
		var kind uint8
		if o.ImageId != 0 {
			kind = internal.KindSprite
		}
		internal.Queue(tex.Texture, src, dst, o.Angle, o.Roundness, mask, (*internal.Effects)(&o.Effects), kind, tbPtr)

		if o.Text != "" {
			v.queueText(o, mask, tbPtr)
			internal.CloseBatch(tbPtr)
			o.textBatches = textBatches
			internal.ReadyBatches = append(internal.ReadyBatches, textBatches...)
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

func (v *View) area() (x, y, w, h float32) {
	if v.Area == (Area{}) {
		return 0, 0, float32(internal.WindowWidth), float32(internal.WindowHeight)
	}
	return v.Area.X, v.Area.Y, v.Area.Width, v.Area.Height
}
func (v *View) queueText(o *Object, mask internal.Area, batches *[]*internal.Batch) {
	var eff = (*internal.Effects)(&o.Effects)
	var scale = eff.TextLineHeight / 255
	var gapX, gapY = eff.TextSymbolGap * scale, eff.TextLineGap * scale
	var fontData, has = internal.Fonts[uint8(o.TextFontId)]
	if !has { // fallback
		fontData = internal.Fonts[0]
	}
	var atlasTex = internal.Images[fontData.AtlasId].Texture
	var sin, cos = internal.SinCos(o.Angle)
	var fullLineHeight = eff.TextLineHeight * fontData.LineHeight
	var height float32
	lines = lines[:0]
	for i := 0; i < len(o.Text); {
		var end, width = o.lineEndAndWidth(i)
		lines = append(lines, line{i, end, width})
		height += fullLineHeight
		i = end
		if i < len(o.Text) {
			height += gapY
			i++
		}
	}
	var y = o.Y - o.Height/2 - fontData.Ascender*eff.TextLineHeight + eff.TextAlignY*(o.Height-height)

	for _, ln := range lines {
		var x = (o.X - o.Width/2) + eff.TextAlignX*(o.Width-ln.width)
		var prevGlyph internal.Glyph

		for _, r := range o.Text[ln.start:ln.end] {
			var glyph = fontData.Chars[r]
			var kerning, _ = prevGlyph.Kernings[r]
			x += kerning * eff.TextLineHeight

			var src, dst = getGlyphSrcDst(o, r, glyph, x, y, cos, sin, 0)
			if r != ' ' && r != '\n' {
				internal.Queue(atlasTex, src, dst, o.Angle, 0, mask, eff, internal.KindText, batches)
			}
			if eff.TextUnderline {
				var src2, dst2 = getGlyphSrcDst(o, internal.Underline, fontData.Chars[internal.Underline], x, y, cos, sin, dst.Width)
				internal.Queue(atlasTex, src2, dst2, o.Angle, 0, mask, eff, internal.KindText, batches)
			}
			if eff.TextCrossout {
				var src2, dst2 = getGlyphSrcDst(o, internal.Crossout, fontData.Chars[internal.Crossout], x, y, cos, sin, dst.Width)
				internal.Queue(atlasTex, src2, dst2, o.Angle, 0, mask, eff, internal.KindText, batches)
			}
			x += glyph.Advance*eff.TextLineHeight + gapX
			prevGlyph = glyph
		}

		y += eff.TextLineHeight*fontData.LineHeight + gapY
	}
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
