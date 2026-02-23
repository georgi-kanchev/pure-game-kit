package graphics

import (
	"pure-game-kit/input/mouse"
	"pure-game-kit/input/mouse/button"
	"pure-game-kit/internal"
	"pure-game-kit/utility/number"
	"pure-game-kit/window"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Camera struct {
	ScreenX, ScreenY, ScreenWidth, ScreenHeight,
	MaskX, MaskY, MaskWidth, MaskHeight, Blend int
	X, Y, Angle, Zoom float32

	Effects *Effects

	//=================================================================

	velocityX, velocityY, dragVelX, dragVelY float32
}

func NewCamera(zoom float32) *Camera {
	var cam = Camera{Zoom: zoom}
	cam.SetScreenAreaToWindow()

	if batch == nil {
		batch = &Batch{}
		batch.Init(16)
	}
	return &cam
}

// =================================================================

func (c *Camera) MouseDragAndZoom() {
	c.Zoom *= 1 + 0.05*mouse.Scroll()

	if mouse.IsButtonPressed(button.Middle) {
		var dx, dy = mouse.CursorDelta()
		var sin, cos = internal.SinCos(-c.Angle)
		c.X -= (dx*cos - dy*sin) / c.Zoom
		c.Y -= (dx*sin + dy*cos) / c.Zoom
	}
}
func (c *Camera) MouseDragAndZoomSmoothly() {
	c.Zoom *= 1 + 0.001*mouse.ScrollSmooth()

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

func (c *Camera) SetScreenArea(screenX, screenY, screenWidth, screenHeight int) {
	c.ScreenX = screenX
	c.ScreenY = screenY
	c.ScreenWidth = screenWidth
	c.ScreenHeight = screenHeight
	c.Mask(screenX, screenY, screenWidth, screenHeight)
}
func (c *Camera) SetScreenAreaToWindow() {
	tryRecreateWindow()
	var w, h = window.Size()
	c.SetScreenArea(0, 0, w, h)
}
func (c *Camera) Mask(screenX, screenY, screenWidth, screenHeight int) {
	c.MaskX, c.MaskY = screenX, screenY
	c.MaskWidth, c.MaskHeight = screenWidth, screenHeight
}

//=================================================================

func (c *Camera) IsAreaVisible(x, y, width, height float32) bool {
	c.update()
	var sx1, sy1 = c.PointToScreen(x, y)
	var sx2, sy2 = c.PointToScreen(x+width, y+height)
	var sMinX, sMaxX = min(sx1, sx2), max(sx1, sx2)
	var sMinY, sMaxY = min(sy1, sy2), max(sy1, sy2)
	var maskR, maskB = c.MaskX + c.MaskWidth, c.MaskY + c.MaskHeight
	return sMaxX > c.MaskX && sMinX < maskR && sMaxY > c.MaskY && sMinY < maskB
}
func (c *Camera) IsHovered() bool {
	var mousePos = rl.GetMousePosition()
	return float32(mousePos.X) > float32(c.ScreenX) &&
		float32(mousePos.Y) > float32(c.ScreenY) &&
		float32(mousePos.X) < float32(c.ScreenX+c.ScreenWidth) &&
		float32(mousePos.Y) < float32(c.ScreenY+c.ScreenHeight)
}
func (c *Camera) MousePosition() (x, y float32) {
	return c.PointFromScreen(int(rl.GetMouseX()), int(rl.GetMouseY()))
}
func (c *Camera) Size() (width, height float32) {
	c.update()
	return float32(c.ScreenWidth) / c.Zoom, float32(c.ScreenHeight) / c.Zoom
}

func (c *Camera) PointFromScreen(screenX, screenY int) (x, y float32) {
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
func (c *Camera) PointToScreen(x, y float32) (screenX, screenY int) {
	c.update()

	x -= float32(rlCam.Target.X)
	y -= float32(rlCam.Target.Y)
	x *= float32(rlCam.Zoom)
	y *= float32(rlCam.Zoom)

	var sin, cos = internal.SinCos(rlCam.Rotation)
	var rotX, rotY = x*cos - y*sin, x*sin + y*cos

	rotX += float32(rlCam.Offset.X)
	rotY += float32(rlCam.Offset.Y)
	return int(rotX), int(rotY)
}
func (c *Camera) PointFromCamera(otherCamera *Camera, otherX, otherY float32) (myX, myY float32) {
	var screenX, screenY = otherCamera.PointToScreen(otherX, otherY)
	return c.PointFromScreen(screenX, screenY)
}
func (c *Camera) PointToCamera(otherCamera *Camera, myX, myY float32) (otherX, otherY float32) {
	return otherCamera.PointFromCamera(c, myX, myY)
}
func (c *Camera) PointFromEdge(edgeX, edgeY float32) (x, y float32) {
	var scrX, scrY = float32(c.ScreenWidth) * edgeX, float32(c.ScreenHeight) * edgeY
	return c.PointFromScreen(int(scrX), int(scrY))
}
