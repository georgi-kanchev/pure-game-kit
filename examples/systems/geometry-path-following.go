package example

import (
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/input/mouse"
	"pure-game-kit/input/mouse/button"
	"pure-game-kit/utility/color"
	"pure-game-kit/window"
)

func PathFollowing() {
	var cam = graphics.NewCamera(2)
	var p1 = [][2]float32{{0, 0}, {50, 0}, {60, 25}, {50, 50}, {0, 50}, {-16, 16}}
	// var p2 = [][2]float32{{50, 0}, {100, 0}, {150, 25}, {175, 35}}
	// var p3 = [][2]float32{{60, 25}, {100, 50}, {150, 60}, {200, 60}, {225, 75}}
	// var p4 = [][2]float32{{250, 0}, {300, 100}}
	var paths = [][2]float32{}
	paths = append(paths, p1...) // main loop
	// paths = append(paths, [2]float32{float32(math.NaN()), float32(math.NaN())})
	// paths = append(paths, p2...) // branch A
	// paths = append(paths, [2]float32{float32(math.NaN()), float32(math.NaN())})
	// paths = append(paths, p3...) // branch B
	// paths = append(paths, [2]float32{float32(math.NaN()), float32(math.NaN())})
	// paths = append(paths, p4...) // isolated diagonal segment (likely ignored unless close)

	var startX, startY float32 = 16, -48

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.MouseDragAndZoomSmooth()
		cam.DrawGrid(1, 32, 32, color.DarkGray)

		cam.DrawCone(0, 0, 200, 360, 0, color.Red)

		var mx, my = cam.MousePosition()
		var result = geometry.FollowPaths(startX, startY, mx, my, paths...)

		if keyboard.IsKeyJustPressed(key.A) {
			geometry.FollowPaths(startX, startY, mx, my, paths...)
		}

		if mouse.IsButtonJustPressed(button.Left) {
			startX, startY = mx, my
		}
		// cam.DrawLinesPath(2, color.White, paths...)

		cam.DrawLinesPath(3, color.Red, p1...)
		// cam.DrawLinesPath(1, color.Blue, p2...)
		// cam.DrawLinesPath(1, color.Magenta, p3...)
		// cam.DrawLinesPath(1, color.Orange, p4...)

		cam.DrawLinesPath(1, color.Green, result...)
	}
}
