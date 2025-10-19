package example

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/utility/color"
	"pure-game-kit/window"
)

func Lines() {
	var cam = graphics.NewCamera(1)
	var lineA = geometry.NewLine(0, 0, 400, 400)
	var lineB = geometry.NewLine(-400, 400, 0, 0)

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		lineB.Bx, lineB.By = cam.MousePosition()

		var lineColor = condition.If(lineA.IsCrossingLine(lineB), color.Red, color.Green)
		var ax, ay = lineA.CrossPointWithLine(lineB)
		var bx, by = lineA.ClosestToPoint(lineB.Bx, lineB.By)
		var pointColor = condition.If(lineB.IsLeftOfPoint(lineA.Bx, lineA.By), color.Blue, color.Yellow)

		cam.DrawLine(lineA.Ax, lineA.Ay, lineA.Bx, lineA.By, 5, color.White)
		cam.DrawLine(lineB.Ax, lineB.Ay, lineB.Bx, lineB.By, 5, lineColor)
		cam.DrawCircle(ax, ay, 15, color.Cyan)
		cam.DrawCircle(bx, by, 10, color.Magenta)
		cam.DrawCircle(lineA.Bx, lineA.By, 5, pointColor)
	}
}
