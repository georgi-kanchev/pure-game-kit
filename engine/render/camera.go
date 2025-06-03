package render

import (
	"math"
	"pure-kit/engine/window"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Camera struct {
	ScreenX, ScreenY          int
	ScreenWidth, ScreenHeight int
	X, Y, Angle, Zoom         float32
}

func (camera *Camera) SetScreenArea(screenX, screenY, screenWidth, screenHeight int) {
	camera.ScreenX = screenX
	camera.ScreenY = screenY
	camera.ScreenWidth = screenWidth
	camera.ScreenHeight = screenHeight
}

func (camera *Camera) DrawColor(color uint) {
	camera.update()

	var x, y, w, h = camera.ScreenX, camera.ScreenY, camera.ScreenWidth, camera.ScreenHeight
	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(h), rl.GetColor(color))
}
func (camera *Camera) DrawFrame(size int, color uint) {
	camera.update()

	var x, y, w, h = camera.ScreenX, camera.ScreenY, camera.ScreenWidth, camera.ScreenHeight
	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(size), rl.GetColor(color))        // upper
	rl.DrawRectangle(int32(x+w-size), int32(y), int32(size), int32(h), rl.GetColor(color)) // right
	rl.DrawRectangle(int32(x), int32(y+h-size), int32(w), int32(size), rl.GetColor(color)) // lower
	rl.DrawRectangle(int32(x), int32(y), int32(size), int32(h), rl.GetColor(color))        // left
}
func (camera *Camera) DrawGrid(thickness, spacing float32, color uint) {
	camera.start()

	// Get all 4 world-space corners of the screen
	ulx, uly := camera.CornerUpperLeft(0, 0)
	urx, ury := camera.CornerUpperRight(0, 0)
	llx, lly := camera.CornerLowerLeft(0, 0)
	lrx, lry := camera.CornerLowerRight(0, 0)

	// Compute axis-aligned bounding box (AABB) from the rotated corners
	xs := []float32{ulx, urx, llx, lrx}
	ys := []float32{uly, ury, lly, lry}

	minX, maxX := xs[0], xs[0]
	minY, maxY := ys[0], ys[0]

	for i := 1; i < 4; i++ {
		if xs[i] < minX {
			minX = xs[i]
		}
		if xs[i] > maxX {
			maxX = xs[i]
		}
		if ys[i] < minY {
			minY = ys[i]
		}
		if ys[i] > maxY {
			maxY = ys[i]
		}
	}

	// Snap bounds to grid spacing
	left := float32(math.Floor(float64(minX/spacing))) * spacing
	right := float32(math.Ceil(float64(maxX/spacing))) * spacing
	top := float32(math.Floor(float64(minY/spacing))) * spacing
	bottom := float32(math.Ceil(float64(maxY/spacing))) * spacing

	// Draw vertical lines
	for x := left; x <= right; x += spacing {
		var myThickness = thickness
		if float32(math.Mod(float64(x), float64(spacing)*10)) == 0 {
			myThickness *= 2
		}

		camera.DrawLine(x, top, x, bottom, myThickness, color)
	}

	// Draw horizontal lines
	for y := top; y <= bottom; y += spacing {
		var myThickness = thickness
		if float32(math.Mod(float64(y), float64(spacing)*10)) == 0 {
			myThickness *= 2
		}

		camera.DrawLine(left, y, right, y, myThickness, color)
	}

	// Draw X axis
	if top <= 0 && bottom >= 0 {
		camera.DrawLine(left, 0, right, 0, thickness*3, color)
	}

	// Draw Y axis
	if left <= 0 && right >= 0 {
		camera.DrawLine(0, top, 0, bottom, thickness*3, color)
	}

	camera.stop()
}

func (camera *Camera) DrawLine(ax, ay, bx, by, thickness float32, color uint) {
	camera.start()
	rl.DrawLineEx(rl.Vector2{X: ax, Y: ay}, rl.Vector2{X: bx, Y: by}, thickness, rl.GetColor(color))
	camera.stop()
}
func (camera *Camera) DrawRectangle(x, y, width, height float32, color uint) {
	camera.start()
	rl.DrawRectangle(int32(x), int32(y), int32(width), int32(height), rl.GetColor(color))
	camera.stop()
}
func (camera *Camera) DrawCircle(x, y, radius float32, color uint) {
	camera.start()
	rl.DrawCircle(int32(x), int32(y), radius, rl.GetColor(color))
	camera.stop()
}
func (camera *Camera) DrawTileMap() {

}

func (camera *Camera) PointFromScreen(screenX, screenY int) (x, y float32) {
	camera.start()

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

	camera.stop()
	return rotX, rotY
}
func (camera *Camera) PointToScreen(x, y float32) (screenX, screenY int) {
	camera.start()

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

	camera.stop()
	return int(rotX), int(rotY)
}
func (camera *Camera) PointFromCamera(otherCamera *Camera, otherX, otherY float32) (myX, myY float32) {
	var screenX, screenY = otherCamera.PointToScreen(otherX, otherY)
	return camera.PointFromScreen(screenX, screenY)
}
func (camera *Camera) PointToCamera(otherCamera *Camera, myX, myY float32) (otherX, otherY float32) {
	return otherCamera.PointFromCamera(camera, myX, myY)
}

func (camera *Camera) Size() (width, height float32) {
	camera.update()
	return float32(camera.ScreenWidth) / camera.Zoom, float32(camera.ScreenHeight) / camera.Zoom
}
func (camera *Camera) CornerUpperLeft(offsetX, offsetY float32) (x, y float32) {
	camera.update()

	var sx, sy, z = camera.ScreenX, camera.ScreenY, camera.Zoom
	x, y = camera.PointFromScreen(sx+int(offsetX*z), sy+int(offsetY*z))
	return x, y
}
func (camera *Camera) CornerUpperRight(offsetX, offsetY float32) (x, y float32) {
	camera.update()

	var sx, sy, z = camera.ScreenX + camera.ScreenWidth, camera.ScreenY, camera.Zoom
	x, y = camera.PointFromScreen(sx+int(offsetX*z), sy+int(offsetY*z))
	return x, y
}
func (camera *Camera) CornerLowerLeft(offsetX, offsetY float32) (x, y float32) {
	camera.update()

	var sx, sy, z = camera.ScreenX, camera.ScreenY + camera.ScreenHeight, camera.Zoom
	x, y = camera.PointFromScreen(sx+int(offsetX*z), sy+int(offsetY*z))
	return x, y
}
func (camera *Camera) CornerLowerRight(offsetX, offsetY float32) (x, y float32) {
	camera.update()

	var sx, sy, z = camera.ScreenX + camera.ScreenWidth, camera.ScreenY + camera.ScreenHeight, camera.Zoom
	x, y = camera.PointFromScreen(sx+int(offsetX*z), sy+int(offsetY*z))
	return x, y
}

// region private

var rlCam = rl.Camera2D{}

// call before draw to update camera but use screen space instead of camera space
func (camera *Camera) update() {
	camera.start()
	camera.stop()
}

// call before draw to update camera and use camera space
func (camera *Camera) start() {
	if !rl.IsWindowReady() {
		window.Recreate()
	}

	rlCam.Target.X = camera.X
	rlCam.Target.Y = camera.Y
	rlCam.Rotation = camera.Angle
	rlCam.Zoom = camera.Zoom
	rlCam.Offset.X = float32(camera.ScreenX) + float32(camera.ScreenWidth)/2
	rlCam.Offset.Y = float32(camera.ScreenY) + float32(camera.ScreenHeight)/2
	rl.BeginMode2D(rlCam)
	rl.BeginScissorMode(
		int32(camera.ScreenX), int32(camera.ScreenY),
		int32(camera.ScreenWidth), int32(camera.ScreenHeight))
}

// call after draw to get back to using screen space
func (camera *Camera) stop() {
	rl.EndScissorMode()
	rl.EndMode2D()
}

// endregion
