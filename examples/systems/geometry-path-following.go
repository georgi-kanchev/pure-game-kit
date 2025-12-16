package example

import (
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/input/mouse"
	"pure-game-kit/input/mouse/button"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/number"
	"pure-game-kit/window"
)

func PathFollowing() {
	var cam = graphics.NewCamera(2)
	var paths = [][2]float32{}
	var startX, startY float32 = 16, -48
	var p1 = [][2]float32{{0, 0}, {50, 0}, {60, 25}, {50, 50}, {0, 50}, {0, 0}}
	var p2 = [][2]float32{{50, 0}, {100, 0}, {150, 25}, {175, 35}}
	var p3 = [][2]float32{{60, 25}, {100, 50}, {150, 60}, {200, 60}, {225, 75}}
	var p4 = [][2]float32{{200, 60}, {360, 65}}
	var p5 = [][2]float32{{225, 75}, {240, 90}, {260, 110}, {280, 115}, {300, 130}}
	var p6 = [][2]float32{{100, 120}, {150, 120}, {180, 130}, {150, 140}, {100, 140}, {100, 120}}
	var p7 = [][2]float32{{280, 115}, {310, 90}, {340, 70}, {360, 65}}
	var p8 = [][2]float32{{360, 65}, {375, 80}, {390, 100}, {400, 125}, {405, 150}}
	paths = append(paths, p1...) // main loop
	paths = append(paths, [2]float32{float32(number.NaN()), float32(number.NaN())})
	paths = append(paths, p2...) // branch A
	paths = append(paths, [2]float32{float32(number.NaN()), float32(number.NaN())})
	paths = append(paths, p3...) // branch B
	paths = append(paths, [2]float32{float32(number.NaN()), float32(number.NaN())})
	paths = append(paths, p4...) // isolated diagonal segment (likely ignored unless close)
	paths = append(paths, [2]float32{float32(number.NaN()), float32(number.NaN())})
	paths = append(paths, p5...) // gentle zig-zag going upward
	paths = append(paths, [2]float32{float32(number.NaN()), float32(number.NaN())})
	paths = append(paths, p6...) // short horizontal loop
	paths = append(paths, [2]float32{float32(number.NaN()), float32(number.NaN())})
	paths = append(paths, p7...) // downward slanted branch
	paths = append(paths, [2]float32{float32(number.NaN()), float32(number.NaN())})
	paths = append(paths, p8...) // tight curved branch
	paths = append(paths, [2]float32{float32(number.NaN()), float32(number.NaN())})

	var randomColors = []uint{}
	for range 8 {
		randomColors = append(randomColors, color.Random())
	}

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.MouseDragAndZoomSmooth()
		cam.DrawGrid(1, 32, 32, palette.DarkGray)

		var mx, my = cam.MousePosition()
		var result = geometry.FollowPaths(startX, startY, mx, my, paths...)

		if mouse.IsButtonJustPressed(button.Left) {
			geometry.FollowPaths(startX, startY, mx, my, paths...)
			startX, startY = mx, my
		}

		// cam.DrawLinesPath(3, color.Red, paths...)
		cam.DrawLinesPath(5, randomColors[0], p1...)
		cam.DrawLinesPath(5, randomColors[1], p2...)
		cam.DrawLinesPath(5, randomColors[2], p3...)
		cam.DrawLinesPath(5, randomColors[3], p4...)
		cam.DrawLinesPath(5, randomColors[4], p5...)
		cam.DrawLinesPath(5, randomColors[5], p6...)
		cam.DrawLinesPath(5, randomColors[6], p7...)
		cam.DrawLinesPath(5, randomColors[7], p8...)
		cam.DrawLinesPath(2, palette.Red, result...)
	}
}
