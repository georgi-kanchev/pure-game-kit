package example

import (
	"pure-game-kit/packages/execution/condition"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/window"
)

func Lines() {
	var view = graphics.NewView(1)
	var lineA = geometry.NewLine(0, 0, 400, 400)
	var lineB = geometry.NewLine(-400, 400, 0, 0)

	for window.KeepOpen() {
		lineB.Bx, lineB.By = view.MousePosition()

		var lineColor = condition.If(lineA.IsCrossingLine(lineB), palette.Red, palette.Green)
		var ax, ay = lineA.CrossPointWithLine(lineB)
		var bx, by = lineA.ClosestToPoint(lineB.Bx, lineB.By)
		var pointColor = condition.If(lineB.IsLeftOfPoint(lineA.Bx, lineA.By), palette.Blue, palette.Yellow)

		view.DrawLine(lineA.Ax, lineA.Ay, lineA.Bx, lineA.By, 5, palette.White)
		view.DrawLine(lineB.Ax, lineB.Ay, lineB.Bx, lineB.By, 5, lineColor)
		view.DrawCircle(ax, ay, 15, 8, palette.Cyan)
		view.DrawCircle(bx, by, 10, 8, palette.Magenta)
		view.DrawCircle(lineA.Bx, lineA.By, 5, 8, pointColor)
	}
}
