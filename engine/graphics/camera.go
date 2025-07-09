package graphics

import (
	"math"
	"pure-kit/engine/window"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Camera struct {
	ScreenX, ScreenY, ScreenWidth, ScreenHeight int
	X, Y, Angle, Zoom, PivotX, PivotY           float32
}

func NewCamera(zoom float32) Camera {
	var cam = Camera{Zoom: zoom, PivotX: 0.5, PivotY: 0.5}
	cam.SetScreenAreaToWindow()
	return cam
}

func (camera *Camera) DragAndZoom() {
	var delta = rl.GetMouseDelta()
	var scroll = rl.GetMouseWheelMove()

	if rl.IsMouseButtonDown(rl.MouseButtonMiddle) {
		camera.X -= delta.X / camera.Zoom
		camera.Y -= delta.Y / camera.Zoom
	}
	if scroll > 0 {
		camera.Zoom *= 1.05
	} else if scroll < 0 {
		camera.Zoom *= 0.95
	}
}
func (camera *Camera) IsHovered() bool {
	var mousePos = rl.GetMousePosition()
	return mousePos.X > float32(camera.ScreenX) &&
		mousePos.Y > float32(camera.ScreenY) &&
		mousePos.X < float32(camera.ScreenX+camera.ScreenWidth) &&
		mousePos.Y < float32(camera.ScreenY+camera.ScreenHeight)
}

func (camera *Camera) SetScreenArea(screenX, screenY, screenWidth, screenHeight int) {
	camera.ScreenX = screenX
	camera.ScreenY = screenY
	camera.ScreenWidth = screenWidth
	camera.ScreenHeight = screenHeight
}
func (camera *Camera) SetScreenAreaToWindow() {
	tryRecreateWindow()

	var w, h = window.Size()
	camera.ScreenX = 0
	camera.ScreenY = 0
	camera.ScreenWidth = w
	camera.ScreenHeight = h
}

func (camera *Camera) Size() (width, height float32) {
	camera.update()
	return float32(camera.ScreenWidth) / camera.Zoom, float32(camera.ScreenHeight) / camera.Zoom
}
func (camera *Camera) MousePosition() (x, y float32) {
	return camera.PointFromScreen(int(rl.GetMouseX()), int(rl.GetMouseY()))
}
func (camera *Camera) PointFromScreen(screenX, screenY int) (x, y float32) {
	camera.begin()

	var sx = float32(screenX)
	var sy = float32(screenY)

	sx -= rlCam.Offset.X
	sy -= rlCam.Offset.Y

	var angle = -rlCam.Rotation * rl.Deg2rad
	var cos = float32(math.Cos(float64(angle)))
	var sin = float32(math.Sin(float64(angle)))

	var rotX = sx*cos - sy*sin
	var rotY = sx*sin + sy*cos

	rotX /= rlCam.Zoom
	rotY /= rlCam.Zoom

	rotX += rlCam.Target.X
	rotY += rlCam.Target.Y

	camera.end()
	return rotX, rotY
}
func (camera *Camera) PointToScreen(x, y float32) (screenX, screenY int) {
	camera.begin()

	x -= rlCam.Target.X
	y -= rlCam.Target.Y

	x *= rlCam.Zoom
	y *= rlCam.Zoom

	var angle = rlCam.Rotation * rl.Deg2rad
	var cos = float32(math.Cos(float64(angle)))
	var sin = float32(math.Sin(float64(angle)))
	var rotX = x*cos - y*sin
	var rotY = x*sin + y*cos

	rotX += rlCam.Offset.X
	rotY += rlCam.Offset.Y

	camera.end()
	return int(rotX), int(rotY)
}
func (camera *Camera) PointFromCamera(otherCamera *Camera, otherX, otherY float32) (myX, myY float32) {
	var screenX, screenY = otherCamera.PointToScreen(otherX, otherY)
	return camera.PointFromScreen(screenX, screenY)
}
func (camera *Camera) PointToCamera(otherCamera *Camera, myX, myY float32) (otherX, otherY float32) {
	return otherCamera.PointFromCamera(camera, myX, myY)
}

// region private

var rlCam = rl.Camera2D{}

// call before draw to update camera but use screen space instead of camera space
func (camera *Camera) update() {
	camera.begin()
	camera.end()
}

// call before draw to update camera and use camera space
func (camera *Camera) begin() {
	tryRecreateWindow()

	rlCam.Target.X = camera.X
	rlCam.Target.Y = camera.Y
	rlCam.Rotation = camera.Angle
	rlCam.Zoom = camera.Zoom
	rlCam.Offset.X = float32(camera.ScreenX) + float32(camera.ScreenWidth)*camera.PivotX
	rlCam.Offset.Y = float32(camera.ScreenY) + float32(camera.ScreenHeight)*camera.PivotY
	rl.BeginMode2D(rlCam)
	rl.BeginScissorMode(
		int32(camera.ScreenX), int32(camera.ScreenY),
		int32(camera.ScreenWidth), int32(camera.ScreenHeight))
}

// call after draw to get back to using screen space
func (camera *Camera) end() {
	rl.EndScissorMode()
	rl.EndMode2D()
}

func tryRecreateWindow() {
	if !rl.IsWindowReady() {
		window.Recreate()
	}
}

// endregion
