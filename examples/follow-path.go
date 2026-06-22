package example

import (
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/mouse"
	"pure-game-kit/packages/input/mouse/button"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/window"
)

func FollowPath() {
	window.Create("example - follow path", false, true)
	var view = graphics.NewView(2)
	var paths = []float32{}
	var startX, startY float32 = 16, -48
	var p1 = []float32{0, 0, 50, 0, 60, 25, 50, 50, 0, 50, 0, 0}
	var p2 = []float32{50, 0, 100, 0, 150, 25, 175, 35}
	var p3 = []float32{60, 25, 100, 50, 150, 60, 200, 60, 225, 75}
	var p4 = []float32{200, 60, 360, 65}
	var p5 = []float32{225, 75, 240, 90, 260, 110, 280, 115, 300, 130}
	var p6 = []float32{100, 120, 150, 120, 180, 130, 150, 140, 100, 140, 100, 120}
	var p7 = []float32{280, 115, 310, 90, 340, 70, 360, 65}
	var p8 = []float32{360, 65, 375, 80, 390, 100, 400, 125, 405, 150}
	paths = append(paths, p1...) // main loop
	paths = append(paths, number.NaN(), number.NaN())
	paths = append(paths, p2...) // branch A
	paths = append(paths, number.NaN(), number.NaN())
	paths = append(paths, p3...) // branch B
	paths = append(paths, number.NaN(), number.NaN())
	paths = append(paths, p4...) // isolated diagonal segment (likely ignored unless close)
	paths = append(paths, number.NaN(), number.NaN())
	paths = append(paths, p5...) // gentle zig-zag going upward
	paths = append(paths, number.NaN(), number.NaN())
	paths = append(paths, p6...) // short horizontal loop
	paths = append(paths, number.NaN(), number.NaN())
	paths = append(paths, p7...) // downward slanted branch
	paths = append(paths, number.NaN(), number.NaN())
	paths = append(paths, p8...) // tight curved branch
	paths = append(paths, number.NaN(), number.NaN())

	var randomColors = []uint{}
	for range 8 {
		randomColors = append(randomColors, color.Random())
	}

	window.SetTargetFPS(0)

	var result = make([]float32, 0)

	for window.KeepOpen() {
		view.MouseDragAndZoomSmoothly()
		view.DrawGrid(1, 32, 32, palette.DarkGray)

		var mx, my = view.MousePosition()
		geometry.FollowPaths(startX, startY, mx, my, paths, &result)

		if mouse.IsButtonJustPressed(button.Left) {
			// geometry.FollowPaths(startX, startY, mx, my, paths)
			startX, startY = mx, my
		}

		// cam.DrawLinesPath(3, color.Red, paths...)
		view.DrawPath(p1, 5, randomColors[0], geometry.Area{})
		view.DrawPath(p2, 5, randomColors[1], geometry.Area{})
		view.DrawPath(p3, 5, randomColors[2], geometry.Area{})
		view.DrawPath(p4, 5, randomColors[3], geometry.Area{})
		view.DrawPath(p5, 5, randomColors[4], geometry.Area{})
		view.DrawPath(p6, 5, randomColors[5], geometry.Area{})
		view.DrawPath(p7, 5, randomColors[6], geometry.Area{})
		view.DrawPath(p8, 5, randomColors[7], geometry.Area{})
		view.DrawPath(result, 2, palette.Red, geometry.Area{})

		view.DrawDebugInfo(true)
	}
}
