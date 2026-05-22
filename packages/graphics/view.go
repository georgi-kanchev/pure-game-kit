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
		var eff *internal.Effects
		if o.Effects != nil {
			eff = (*internal.Effects)(o.Effects)
		}
		var kind byte
		if o.charValue != 0 {
			kind = 2 // text
			var td = internal.TextDraw{ShadowColor: o.TextShadowColor, Weight: o.TextWeight,
				ShadowBlur: byte(o.TextShadowBlur), ShadowX: int8(o.TextShadowOffsetX), ShadowY: int8(o.TextShadowOffsetY)}
			internal.QueueText(tex.Texture, src, dst, o.Angle, getColor(o.Color), internal.Area(o.Mask), eff, td)
			continue
		} else if o.ImageId != 0 {
			kind = 1 // sprite
		}
		internal.QueueTexture(tex.Texture, src, dst, o.Angle, getColor(o.Color), internal.Area(o.Mask), eff, kind)

		if o.Text != "" {
			var fontData = internal.Fonts[byte(o.TextFontId)]
			chars = chars[:0]
			var x, y = o.X, o.Y
			for _, r := range o.Text {
				var symbol = NewImage(0, 0, 0)
				var char = fontData.Chars[r]
				var dst = char.PlaneBounds
				var src = char.AtlasBounds
				var w, h = float32(src.Right - src.Left), float32(src.Bottom - src.Top)
				symbol.ImageId = assets.ImageId(fontData.AtlasId)
				symbol.ImageCropArea = Area{X: float32(src.Left), Y: float32(src.Top), Width: w, Height: h}
				symbol.X = x + (float32(dst.Left) * o.TextLineHeight)
				symbol.Y = y + (float32(dst.Top) * o.TextLineHeight)
				symbol.Width = (float32(char.PlaneBounds.Right) - float32(char.PlaneBounds.Left)) * o.TextLineHeight
				symbol.Height = (float32(char.PlaneBounds.Bottom) - float32(char.PlaneBounds.Top)) * o.TextLineHeight
				symbol.charValue = r
				symbol.Color = o.TextColor
				symbol.TextColor = o.TextColor
				symbol.TextShadowColor = o.TextShadowColor
				symbol.TextWeight = o.TextWeight
				symbol.TextShadowBlur = o.TextShadowBlur
				symbol.TextShadowOffsetX = o.TextShadowOffsetX
				symbol.TextShadowOffsetY = o.TextShadowOffsetY
				x += float32(char.Advance)*o.TextLineHeight + 10
				chars = append(chars, symbol)
			}
			for _, c := range chars {
				v.DrawObjects(&c)
			}
		}
	}
}

// private ========================================================

var chars []Object

func (v *View) area() (x, y, w, h float32) {
	if v.Area == (Area{}) {
		return 0, 0, float32(internal.WindowWidth), float32(internal.WindowHeight)
	}
	return v.Area.X, v.Area.Y, v.Area.Width, v.Area.Height
}
