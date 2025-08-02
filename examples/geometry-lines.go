package example

import (
	"pure-kit/engine/execution/condition"
	"pure-kit/engine/geometry/line"
	"pure-kit/engine/graphics"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
)

func GeometryLines() {
	var cam = graphics.NewCamera(1)
	var lineA = line.New(0, 0, 400, 400)
	var lineB = line.New(-400, 400, 0, 0)

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		lineB.Bx, lineB.By = cam.MousePosition()

		var lineColor = condition.If(lineA.IsCrossing(lineB), color.Red, color.Green)
		var ax, ay = lineA.CrossPoint(lineB)
		var bx, by = lineA.ClosestPointTo(lineB.Bx, lineB.By)
		var pointColor = condition.If(lineB.IsLeftOfPoint(lineA.Bx, lineA.By), color.Blue, color.Yellow)

		cam.DrawLine(lineA.Ax, lineA.Ay, lineA.Bx, lineA.By, 5, color.White)
		cam.DrawLine(lineB.Ax, lineB.Ay, lineB.Bx, lineB.By, 5, lineColor)
		cam.DrawCircle(ax, ay, 15, color.Cyan)
		cam.DrawCircle(bx, by, 10, color.Magenta)
		cam.DrawCircle(lineA.Bx, lineA.By, 5, pointColor)
	}
}
