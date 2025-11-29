package example

import (
	"math"
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/utility/color"
	"pure-game-kit/window"
)

func PathFollowing() {
	var cam = graphics.NewCamera(2)
	var paths = [][2]float32{
		{0, 0}, {50, 0}, {60, 25}, {50, 50}, {0, 50}, {0, 0}, // main loop
		{float32(math.NaN()), float32(math.NaN())},
		{50, 0}, {100, 0}, {150, 25}, {175, 35}, // branch A
		{float32(math.NaN()), float32(math.NaN())},
		{60, 25}, {100, 50}, {150, 60}, {200, 60}, {225, 75}, // branch B
		{float32(math.NaN()), float32(math.NaN())},
		{250, 0}, {300, 100}, // isolated diagonal segment (likely ignored unless close)
	}

	var start = [2]float32{16, -48}

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.MouseDragAndZoomSmooth()
		cam.DrawGrid(1, 32, 32, color.DarkGray)

		var mx, my = cam.MousePosition()
		var target = [2]float32{mx, my}
		var result = geometry.FollowPath(start, target, paths...)

		cam.DrawLinesPath(1, color.White, paths...)
		cam.DrawLinesPath(1, color.Green, result...)
		cam.DrawPoints(2, color.Green, start, target)
	}
}
