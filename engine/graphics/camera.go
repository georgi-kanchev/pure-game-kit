package graphics

import (
	"pure-kit/engine/geometry/point"
	"pure-kit/engine/internal"
	"pure-kit/engine/utility/angle"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/window"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Camera struct {
	ScreenX, ScreenY, ScreenWidth, ScreenHeight,
	MaskX, MaskY, MaskWidth, MaskHeight int
	X, Y, Angle, Zoom, PivotX, PivotY float32

	// Makes sequencial Draw calls faster.
	// All of the drawing to the camera can be batched, as long as the other parameters don't change.
	// Make sure to turn off batching after done drawing the batch,
	// otherwise the other camera parameters will never take effect visually again.
	//
	// 	// recommended
	// 	camera.Angle = 45
	// 	camera.Batch = true
	// 	camera.Draw...
	// 	camera.Draw...
	// 	camera.Batch = false
	// 	camera.X = 300
	//
	//	// not recommended
	// 	camera.Batch = true
	// 	camera.Draw...
	// 	camera.Angle = 45
	// 	camera.X = 300
	// 	camera.Draw...
	// 	camera.Batch = false
	Batch bool
}

func NewCamera(zoom float32) *Camera {
	var cam = Camera{Zoom: zoom, PivotX: 0.5, PivotY: 0.5}
	cam.SetScreenAreaToWindow()
	return &cam
}

//=================================================================
// setters

func (camera *Camera) DragAndZoom() {
	var delta = rl.GetMouseDelta()
	var scroll = rl.GetMouseWheelMove()

	if rl.IsMouseButtonDown(rl.MouseButtonMiddle) {
		var rad = angle.ToRadians(-camera.Angle)
		var sin, cos = number.Sine(rad), number.Cosine(rad)
		camera.X -= (delta.X*cos - delta.Y*sin) / camera.Zoom
		camera.Y -= (delta.X*sin + delta.Y*cos) / camera.Zoom
	}
	if scroll > 0 {
		camera.Zoom *= 1.05
	} else if scroll < 0 {
		camera.Zoom *= 0.95
	}
}

func (camera *Camera) SetScreenArea(screenX, screenY, screenWidth, screenHeight int) {
	camera.ScreenX = screenX
	camera.ScreenY = screenY
	camera.ScreenWidth = screenWidth
	camera.ScreenHeight = screenHeight
	camera.Mask(screenX, screenY, screenWidth, screenHeight)
}
func (camera *Camera) SetScreenAreaToWindow() {
	tryRecreateWindow()
	var w, h = window.Size()
	camera.SetScreenArea(0, 0, w, h)
}
func (camera *Camera) Mask(screenX, screenY, screenWidth, screenHeight int) {
	camera.MaskX, camera.MaskY = screenX, screenY
	camera.MaskWidth, camera.MaskHeight = screenWidth, screenHeight
}

//=================================================================
// getters

func (camera *Camera) IsHovered() bool {
	var mousePos = rl.GetMousePosition()
	return float32(mousePos.X) > float32(camera.ScreenX) &&
		float32(mousePos.Y) > float32(camera.ScreenY) &&
		float32(mousePos.X) < float32(camera.ScreenX+camera.ScreenWidth) &&
		float32(mousePos.Y) < float32(camera.ScreenY+camera.ScreenHeight)
}
func (camera *Camera) MousePosition() (x, y float32) {
	return camera.PointFromScreen(int(rl.GetMouseX()), int(rl.GetMouseY()))
}
func (camera *Camera) Size() (width, height float32) {
	camera.update()
	return float32(camera.ScreenWidth) / camera.Zoom, float32(camera.ScreenHeight) / camera.Zoom
}

func (camera *Camera) PointFromScreen(screenX, screenY int) (x, y float32) {
	camera.update()

	var sx, sy = float32(screenX), float32(screenY)
	sx -= float32(rlCam.Offset.X)
	sy -= float32(rlCam.Offset.Y)

	var angle = angle.ToRadians(-rlCam.Rotation)
	var cos, sin = number.Cosine(angle), number.Sine(angle)
	var rotX, rotY = sx*cos - sy*sin, sx*sin + sy*cos

	rotX /= float32(rlCam.Zoom)
	rotY /= float32(rlCam.Zoom)
	rotX += float32(rlCam.Target.X)
	rotY += float32(rlCam.Target.Y)
	return rotX, rotY
}
func (camera *Camera) PointToScreen(x, y float32) (screenX, screenY int) {
	camera.update()

	x -= float32(rlCam.Target.X)
	y -= float32(rlCam.Target.Y)
	x *= float32(rlCam.Zoom)
	y *= float32(rlCam.Zoom)

	var angle = angle.ToRadians(rlCam.Rotation)
	var cos, sin = number.Cosine(angle), number.Sine(angle)
	var rotX, rotY = x*cos - y*sin, x*sin + y*cos

	rotX += float32(rlCam.Offset.X)
	rotY += float32(rlCam.Offset.Y)
	return int(rotX), int(rotY)
}
func (camera *Camera) PointFromCamera(otherCamera *Camera, otherX, otherY float32) (myX, myY float32) {
	var screenX, screenY = otherCamera.PointToScreen(otherX, otherY)
	return camera.PointFromScreen(screenX, screenY)
}
func (camera *Camera) PointToCamera(otherCamera *Camera, myX, myY float32) (otherX, otherY float32) {
	return otherCamera.PointFromCamera(camera, myX, myY)
}
func (camera *Camera) PointFromPivot(pivotX, pivotY float32) (x, y float32) {
	// useful to get edge coordinates
	var prevX, prevY = camera.PivotX, camera.PivotY
	camera.PivotX, camera.PivotY = pivotX, pivotY
	var scrX, scrY = camera.PointToScreen(0, 0)
	camera.PivotX, camera.PivotY = prevX, prevY
	return camera.PointFromScreen(scrX, scrY)
}

//=================================================================
// private

var rlCam = rl.Camera2D{}
var maskX, maskY, maskW, maskH int

// call before draw to update camera but use screen space instead of camera space
func (camera *Camera) update() {
	tryRecreateWindow()

	rlCam.Target.X = float32(camera.X)
	rlCam.Target.Y = float32(camera.Y)
	rlCam.Rotation = float32(camera.Angle)
	rlCam.Zoom = float32(camera.Zoom)
	rlCam.Offset.X = float32(camera.ScreenX) + float32(camera.ScreenWidth)*float32(camera.PivotX)
	rlCam.Offset.Y = float32(camera.ScreenY) + float32(camera.ScreenHeight)*float32(camera.PivotY)

	var mx = number.BiggestInt(camera.MaskX, camera.ScreenX)
	var my = number.BiggestInt(camera.MaskY, camera.ScreenY)
	var maxW = camera.ScreenX + camera.ScreenWidth - mx
	var maxH = camera.ScreenY + camera.ScreenHeight - my
	var mw = number.SmallestInt(camera.MaskWidth, maxW)
	var mh = number.SmallestInt(camera.MaskHeight, maxH)

	maskX, maskY, maskW, maskH = mx, my, mw, mh
}

// call before draw to update camera and use camera space
func (camera *Camera) begin() {
	camera.update()
	if camera.Batch {
		return
	}

	rl.BeginMode2D(rlCam)
	rl.BeginScissorMode(int32(maskX), int32(maskY), int32(maskW), int32(maskH))
}

// call after draw to get back to using screen space
func (camera *Camera) end() {
	if camera.Batch {
		return
	}

	rl.EndScissorMode()
	rl.EndMode2D()
}

func (camera *Camera) isAreaVisible(x, y, width, height, angle float32) bool {
	var tlx, tly = x, y
	var trx, try = point.MoveAtAngle(tlx, tly, angle, width)
	var brx, bry = point.MoveAtAngle(trx, try, angle+90, height)
	var blx, bly = point.MoveAtAngle(tlx, tly, angle+90, height)
	var stlx, stly = camera.PointToScreen(tlx, tly)
	var strx, stry = camera.PointToScreen(trx, try)
	var sbrx, sbry = camera.PointToScreen(brx, bry)
	var sblx, sbly = camera.PointToScreen(blx, bly)
	var mtlx, mtly = camera.MaskX, camera.MaskY
	var mbrx, mbry = camera.MaskX + camera.MaskWidth, camera.MaskY + camera.MaskHeight
	var minX = number.SmallestInt(stlx, strx, sbrx, sblx)
	var maxX = number.BiggestInt(stlx, strx, sbrx, sblx)
	var minY = number.SmallestInt(stly, stry, sbry, sbly)
	var maxY = number.BiggestInt(stly, stry, sbry, sbly)

	return maxY > mtly && minY < mbry && maxX > mtlx && minX < mbrx
}

func tryRecreateWindow() {
	if internal.WindowReady {
		return
	}

	if !rl.IsWindowReady() {
		window.Recreate()
	}
}
