package graphics

import (
	"pure-game-kit/packages/input/mouse"
	"pure-game-kit/packages/input/mouse/button"
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/number"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type View struct {
	X, Y, Zoom, Angle float32

	Area    Area // The draw area in window space. Zero value = entire window.
	Mask    Area // In view space. Everything drawn outside of it is cropped. Zero value = no masking.
	Effects *Effects
	Blend   int

	//=================================================================

	velocityX, velocityY, dragVelX, dragVelY float32
}

func NewView(zoom float32) *View {
	return &View{Zoom: zoom}
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

		var tex = internal.Images[int32(o.ImageId)]
		var crop = o.ImageCropArea
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
		var eff = (*internal.Effects)(o.Effects)
		var kind uint8
		if o.ImageId != 0 {
			kind = internal.KindSprite // sprite
		}
		internal.QueueTexture(tex.Texture, src, dst, o.Angle, o.Roundness, getColor(o.Color), mask, eff, kind)

		if o.Text != "" {
			v.queueText(o, mask, eff)
		}
	}
}

// private ========================================================

func (v *View) area() (x, y, w, h float32) {
	if v.Area == (Area{}) {
		return 0, 0, float32(internal.WindowWidth), float32(internal.WindowHeight)
	}
	return v.Area.X, v.Area.Y, v.Area.Width, v.Area.Height
}

func (v *View) queueText(o *Object, mask internal.Area, eff *internal.Effects) {
	var lineHeight float32 = 40
	var scale = lineHeight / 255
	var c = palette.White
	var gapX, gapY float32
	if eff != nil {
		lineHeight = eff.TextLineHeight
		scale = lineHeight / 255
		c = eff.TextColor
		gapX, gapY = eff.TextSymbolGap*scale, eff.TextLineGap*scale
	}

	var fontData = internal.Fonts[uint8(o.TextFontId)]
	var atlasTex = internal.Images[fontData.AtlasId].Texture
	var x = o.X - o.Width/2
	var y = o.Y - o.Height/2 - fontData.Ascender*lineHeight
	var col = getColor(c)
	var prevGlyph internal.Glyph
	var sinA, cosA = internal.SinCos(o.Angle)
	for _, r := range o.Text {
		var glyph = fontData.Chars[r]
		var kerning, _ = prevGlyph.Kernings[r]
		x += kerning * lineHeight

		if r == ' ' {
			x += lineHeight/3 + gapX
			prevGlyph = glyph
			continue
		} else if r == '\n' {
			x = o.X - o.Width/2
			y += lineHeight*fontData.LineHeight + gapY
			continue
		}

		var plane, atlas = glyph.PlaneBounds, glyph.AtlasBounds
		var srcW, srcH = atlas.Right - atlas.Left, atlas.Bottom - atlas.Top
		var srcX, srcY = atlas.Left, atlas.Top
		var dstX, dstY = x + plane.Left*lineHeight, y + plane.Top*lineHeight
		var dstW, dstH = (plane.Right - plane.Left) * lineHeight, (plane.Top - plane.Bottom) * lineHeight

		var physTop, physBot = dstY, dstY - dstH
		var tbLeft, tbTop = o.X - o.Width/2, o.Y - o.Height/2
		var tbRight, tbBot = o.X + o.Width/2, o.Y + o.Height/2
		var clipLeft, clipRight = max(dstX, tbLeft), min(dstX+dstW, tbRight)
		var clipTop, clipBot = max(physTop, tbTop), min(physBot, tbBot)
		if clipLeft >= clipRight || clipTop >= clipBot {
			x += glyph.Advance * lineHeight
			prevGlyph = glyph
			continue
		}
		var clippedW, clippedH = clipRight - clipLeft, clipBot - clipTop
		var origH = physBot - physTop
		srcX += (clipLeft - dstX) / dstW * srcW
		srcY += (clipTop - physTop) / origH * srcH
		srcW, srcH = srcW*(clippedW/dstW), srcH*(clippedH/origH)
		dstX, dstY = clipLeft, clipTop
		dstW, dstH = clippedW, -clippedH // restore negative convention

		var dx, dy = (clipLeft + clippedW/2) - o.X, ((clipTop + clipBot) / 2) - o.Y
		dstX = o.X + dx*cosA - dy*sinA - dstW/2
		dstY = o.Y + dx*sinA + dy*cosA + dstH/2

		var dst = rl.NewRectangle(dstX, dstY, dstW, dstH)
		var src = rl.NewRectangle(srcX, srcY, srcW, srcH)
		internal.QueueTexture(atlasTex, src, dst, o.Angle, o.Roundness, col, mask, eff, 2)
		x += glyph.Advance*lineHeight + gapX
		prevGlyph = glyph
	}
}
