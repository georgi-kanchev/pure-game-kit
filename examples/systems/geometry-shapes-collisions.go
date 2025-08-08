package example

import (
	"pure-kit/engine/geometry"
	"pure-kit/engine/graphics"
	"pure-kit/engine/input/keyboard"
	"pure-kit/engine/input/keyboard/key"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/utility/direction"
	"pure-kit/engine/utility/seconds"
	"pure-kit/engine/window"
)

func Collisions() {
	var cam = graphics.NewCamera(1.5)
	var shape = geometry.NewShapeSides(100, 4)
	var shape2 = geometry.NewShapeSides(300, 3)
	var shape3 = geometry.NewShapeSides(420, 5)
	var shape4 = geometry.NewShapeCorners(
		[2]float32{10, 10},
		[2]float32{150, -50},
		[2]float32{180, 100},
		[2]float32{120, 180},
		[2]float32{40, 160},
	)

	shape2.Angle, shape3.Angle = 40, 33
	shape3.X, shape4.X = 700, -700
	shape4.ScaleX, shape4.ScaleY = 2.5, 2.5

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()

		var dirX, dirY float32 = 0, 0
		var step = seconds.FrameDelta() * 600

		if keyboard.IsKeyPressed(key.A) {
			dirX -= 1
		}
		if keyboard.IsKeyPressed(key.D) {
			dirX += 1
		}
		if keyboard.IsKeyPressed(key.W) {
			dirY -= 1
		}
		if keyboard.IsKeyPressed(key.S) {
			dirY += 1
		}
		shape4.Angle++

		dirX, dirY = direction.Normalize(dirX, dirY)
		dirX, dirY = shape.Collide(dirX*step, dirY*step, &shape2, &shape3, &shape4)
		shape.X += dirX
		shape.Y += dirY
		cam.DrawLinesPath(8, color.Red, shape2.CornerPoints()...)
		cam.DrawLinesPath(8, color.Red, shape3.CornerPoints()...)
		cam.DrawLinesPath(8, color.Red, shape4.CornerPoints()...)
		cam.DrawLinesPath(8, color.Green, shape.CornerPoints()...)
	}
}
