package graphics

import (
	"pure-game-kit/input/mouse"
	"pure-game-kit/input/mouse/button"
	"pure-game-kit/internal"
	"pure-game-kit/utility/number"
)

type Camera struct {
	X, Y, Zoom, Angle float32

	Area    Area // The draw area in window space. Zero value = entire window.
	Mask    Area // In camera space. Everything drawn outside of it is cropped. Zero value = no masking.
	Effects Effects
	Blend   int

	//=================================================================

	velocityX, velocityY, dragVelX, dragVelY float32
}

func NewCamera(zoom float32) *Camera {
	tryRecreateWindow()
	var cam = &Camera{Zoom: zoom, Effects: NewEffects()}
	if batch == nil {
		batch = &batchData{}
		batch.Init(16)
	}
	return cam
}

// =================================================================

func (c *Camera) MouseDragAndZoom() {
	var oldZoom = c.Zoom
	var scroll = mouse.Scroll()

	if scroll != 0 {
		c.Zoom *= 1 + 0.05*scroll
		var mx, my = c.MousePosition()

		c.X += (mx - c.X) * (c.Zoom/oldZoom - 1)
		c.Y += (my - c.Y) * (c.Zoom/oldZoom - 1)
	}

	if mouse.IsButtonPressed(button.Middle) {
		var dx, dy = mouse.CursorDelta()
		var sin, cos = internal.SinCos(-c.Angle)
		c.X -= (dx*cos - dy*sin) / c.Zoom
		c.Y -= (dx*sin + dy*cos) / c.Zoom
	}
}
func (c *Camera) MouseDragAndZoomSmoothly() {
	var oldZoom = c.Zoom
	var scroll = mouse.ScrollSmooth()

	if scroll != 0 {
		c.Zoom *= 1 + 0.001*scroll
		var mx, my = c.MousePosition()

		c.X += (mx - c.X) * (c.Zoom/oldZoom - 1)
		c.Y += (my - c.Y) * (c.Zoom/oldZoom - 1)
	}

	const dragFriction, dragStrength = 8.0, 8.0
	var dt = internal.DeltaTime
	var decay = number.Exponential(-dragFriction * dt)
	c.velocityX *= decay
	c.velocityY *= decay

	if mouse.IsButtonPressed(button.Middle) {
		var sin, cos = internal.SinCos(-c.Angle)
		var dx, dy = mouse.CursorDelta()

		dx /= dt
		dy /= dt

		var wdx = (dx*cos - dy*sin) / c.Zoom
		var wdy = (dx*sin + dy*cos) / c.Zoom

		c.velocityX -= wdx * dragStrength * dt
		c.velocityY -= wdy * dragStrength * dt
	}

	c.X += c.velocityX * dt
	c.Y += c.velocityY * dt

	if number.Absolute(c.velocityX) < 0.0001 {
		c.velocityX = 0
	}
	if number.Absolute(c.velocityY) < 0.0001 {
		c.velocityY = 0
	}
}

//=================================================================

func (c *Camera) IsAreaVisible(x, y, width, height float32) bool {
	c.update()
	var sx1, sy1 = c.PointToScreen(x, y)
	var sx2, sy2 = c.PointToScreen(x+width, y+height)
	var sMinX, sMaxX = min(sx1, sx2), max(sx1, sx2)
	var sMinY, sMaxY = min(sy1, sy2), max(sy1, sy2)
	var mx, my, mw, mh = c.area()
	return sMaxX > mx && sMinX < mx+mw && sMaxY > my && sMinY < my+mh
}
func (c *Camera) IsHovered() bool {
	var mx, my = internal.MouseX, internal.MouseY
	var sx, sy, sw, sh = c.area()
	return mx > sx && my > sy && mx < sx+sw && my < sy+sh
}
func (c *Camera) MousePosition() (x, y float32) {
	return c.PointFromScreen(internal.MouseX, internal.MouseY)
}
func (c *Camera) Size() (width, height float32) {
	c.update()
	var _, _, sw, sh = c.area()
	return sw / c.Zoom, sh / c.Zoom
}
func (c *Camera) Bounds() (x, y, width, height float32) {
	var x1, y1 = c.PointFromEdge(0, 0)
	var x2, y2 = c.PointFromEdge(1, 0)
	var x3, y3 = c.PointFromEdge(1, 1)
	var x4, y4 = c.PointFromEdge(0, 1)
	var minX, minY = number.Smallest(x1, x2, x3, x4), number.Smallest(y1, y2, y3, y4)
	var maxX, maxY = number.Biggest(x1, x2, x3, x4), number.Biggest(y1, y2, y3, y4)
	return minX, minY, maxX - minX, maxY - minY
}

func (c *Camera) PointFromScreen(screenX, screenY float32) (x, y float32) {
	c.update()

	var sx, sy = float32(screenX), float32(screenY)
	sx -= float32(rlCam.Offset.X)
	sy -= float32(rlCam.Offset.Y)

	var sin, cos = internal.SinCos(-rlCam.Rotation)
	var rotX, rotY = sx*cos - sy*sin, sx*sin + sy*cos

	rotX /= float32(rlCam.Zoom)
	rotY /= float32(rlCam.Zoom)
	rotX += float32(rlCam.Target.X)
	rotY += float32(rlCam.Target.Y)
	return rotX, rotY
}
func (c *Camera) PointToScreen(x, y float32) (screenX, screenY float32) {
	c.update()

	x -= float32(rlCam.Target.X)
	y -= float32(rlCam.Target.Y)
	x *= float32(rlCam.Zoom)
	y *= float32(rlCam.Zoom)

	var sin, cos = internal.SinCos(rlCam.Rotation)
	var rotX, rotY = x*cos - y*sin, x*sin + y*cos

	rotX += float32(rlCam.Offset.X)
	rotY += float32(rlCam.Offset.Y)
	return rotX, rotY
}
func (c *Camera) PointFromCamera(otherCamera *Camera, otherX, otherY float32) (myX, myY float32) {
	var screenX, screenY = otherCamera.PointToScreen(otherX, otherY)
	return c.PointFromScreen(screenX, screenY)
}
func (c *Camera) PointToCamera(otherCamera *Camera, myX, myY float32) (otherX, otherY float32) {
	return otherCamera.PointFromCamera(c, myX, myY)
}
func (c *Camera) PointFromEdge(edgeX, edgeY float32) (x, y float32) {
	var _, _, sw, sh = c.area()
	var scrX, scrY = sw * edgeX, sh * edgeY
	return c.PointFromScreen(scrX, scrY)
}
