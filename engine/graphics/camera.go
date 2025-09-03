package graphics

import (
	"math"
	"pure-kit/engine/geometry/point"
	"pure-kit/engine/internal"
	"pure-kit/engine/window"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Camera struct {
	ScreenX, ScreenY, ScreenWidth, ScreenHeight int
	X, Y, Angle, Zoom, PivotX, PivotY           float32

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

	maskX, maskY, maskW, maskH int
}

func NewCamera(zoom float32) *Camera {
	var cam = Camera{Zoom: zoom, PivotX: 0.5, PivotY: 0.5}
	cam.SetScreenAreaToWindow()
	return &cam
}

func (camera *Camera) DragAndZoom() {
	var delta = rl.GetMouseDelta()
	var scroll = rl.GetMouseWheelMove()

	if rl.IsMouseButtonDown(rl.MouseButtonMiddle) {
		camera.X -= float32(delta.X) / camera.Zoom
		camera.Y -= float32(delta.Y) / camera.Zoom
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
	camera.maskX, camera.maskY = screenX, screenY
	camera.maskW, camera.maskH = screenWidth, screenHeight
}

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

	var sx = float32(screenX)
	var sy = float32(screenY)

	sx -= float32(rlCam.Offset.X)
	sy -= float32(rlCam.Offset.Y)

	var angle = -rlCam.Rotation * rl.Deg2rad
	var cos = float32(math.Cos(float64(angle)))
	var sin = float32(math.Sin(float64(angle)))

	var rotX = sx*cos - sy*sin
	var rotY = sx*sin + sy*cos

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

	var angle = rlCam.Rotation * rl.Deg2rad
	var cos = float32(math.Cos(float64(angle)))
	var sin = float32(math.Sin(float64(angle)))
	var rotX = x*cos - y*sin
	var rotY = x*sin + y*cos

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

// region private

var rlCam = rl.Camera2D{}
var maskX, maskY, maskW, maskH int32

// call before draw to update camera but use screen space instead of camera space
func (camera *Camera) update() {
	tryRecreateWindow()

	rlCam.Target.X = float32(camera.X)
	rlCam.Target.Y = float32(camera.Y)
	rlCam.Rotation = float32(camera.Angle)
	rlCam.Zoom = float32(camera.Zoom)
	rlCam.Offset.X = float32(camera.ScreenX) + float32(camera.ScreenWidth)*float32(camera.PivotX)
	rlCam.Offset.Y = float32(camera.ScreenY) + float32(camera.ScreenHeight)*float32(camera.PivotY)

	var mx = int(math.Max(float64(camera.maskX), float64(camera.ScreenX)))
	var my = int(math.Max(float64(camera.maskY), float64(camera.ScreenY)))
	var maxW = camera.ScreenX + camera.ScreenWidth - mx
	var maxH = camera.ScreenY + camera.ScreenHeight - my
	var mw = int32(math.Min(float64(camera.maskW), float64(maxW)))
	var mh = int32(math.Min(float64(camera.maskH), float64(maxH)))

	maskX, maskY, maskW, maskH = int32(mx), int32(my), mw, mh
}

// call before draw to update camera and use camera space
func (camera *Camera) begin() {
	camera.update()
	if camera.Batch {
		return
	}

	rl.BeginMode2D(rlCam)
	rl.BeginScissorMode(maskX, maskY, maskW, maskH)
}

// call after draw to get back to using screen space
func (camera *Camera) end() {
	if camera.Batch {
		return
	}

	rl.EndScissorMode()
	rl.EndMode2D()
}

func (camera *Camera) isAreaVisible(x, y, width, height, pivotX, pivotY, angle float32) bool {
	var tlx, _ = point.MoveAtAngle(x, y, angle, width*pivotX)
	var _, tly = point.MoveAtAngle(x, y, angle+90, height*pivotY)
	var trx, try = point.MoveAtAngle(tlx, tly, angle, width)
	var brx, bry = point.MoveAtAngle(trx, try, angle+90, height)
	var blx, bly = point.MoveAtAngle(tlx, tly, angle+90, height)
	var stlx, stly = camera.PointToScreen(tlx, tly)
	var strx, stry = camera.PointToScreen(trx, try)
	var sbrx, sbry = camera.PointToScreen(brx, bry)
	var sblx, sbly = camera.PointToScreen(blx, bly)
	var mtlx, mtly = camera.maskX, camera.maskY
	var mbrx, mbry = camera.maskX + camera.maskW, camera.maskY + camera.maskH
	var minX = int(math.Min(math.Min(float64(stlx), float64(strx)), math.Min(float64(sbrx), float64(sblx))))
	var maxX = int(math.Max(math.Max(float64(stlx), float64(strx)), math.Max(float64(sbrx), float64(sblx))))
	var minY = int(math.Min(math.Min(float64(stly), float64(stry)), math.Min(float64(sbry), float64(sbly))))
	var maxY = int(math.Max(math.Max(float64(stly), float64(stry)), math.Max(float64(sbry), float64(sbly))))

	return !(maxX < mtlx || minX > mbrx || maxY < mtly || minY > mbry)
}

func tryRecreateWindow() {
	if internal.WindowReady {
		return
	}

	if !rl.IsWindowReady() {
		window.Recreate()
	}
}

// endregion
