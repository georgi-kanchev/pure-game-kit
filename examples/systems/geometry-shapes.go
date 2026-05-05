package example

import (
	"pure-game-kit/packages/execution/condition"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/time"
	"pure-game-kit/packages/window"
)

func Shapes() {
	var view = graphics.NewView(1)
	var shape = geometry.NewShapeCorners(
		0, 0,
		50, -20,
		100, 0,
		0, 100,
		50, 120,
		100, 100,
	)
	var triangle = geometry.NewShapeCorners(0, 0, 100, 100, -100, 100)
	var rectangle = geometry.NewShapeQuad(700, 500, 0.5, 0.5)
	var hexagon = geometry.NewShapeSides(500, 6)
	var ellipse = geometry.NewShapeEllipse(200, 500, 16)
	var squircle = geometry.NewShapeQuadRounded(300, 150, 20, 0.5, 0.5, 16)

	shape.ScaleX, shape.ScaleY = 5, 5
	shape.X += 180
	shape.Y -= 200
	rectangle.Angle = 45

	ellipse.X, ellipse.Y = -800, 0
	ellipse.Angle = 20

	squircle.X, squircle.Y = -650, 100

	var stars = []float32{
		600 + 300, 100, 600 + 350, 200, 600 + 450, 200, 600 + 370, 260,
		600 + 400, 360, 600 + 300, 300, 600 + 200, 360, 600 + 230, 260,
		600 + 150, 200, 600 + 250, 200, 600 + 300, 100,
		number.NaN(), number.NaN(),
		300, 100, 350, 200, 450, 200, 370, 260,
		400, 360, 300, 300, 200, 360, 230, 260,
		150, 200, 250, 200, 300, 100,
	}

	for window.KeepOpen() {
		squircle.Angle += time.FrameDelta() * 60
		shape.Angle += time.FrameDelta() * 60
		var mx, my = view.MousePosition()
		var colShape = condition.If(shape.IsOverlappingShapes(triangle), palette.Red, palette.Green)
		var colRect = condition.If(rectangle.IsCrossingShapes(shape), palette.Brown, palette.Cyan)
		var colCircle = condition.If(hexagon.IsContainingShapes(triangle), palette.Yellow, palette.Pink)

		triangle.X, triangle.Y = mx, my

		var crossPoints = hexagon.CrossPointsWithShapes(shape)
		var hexagonPts = hexagon.CornerPoints()
		var rectPts = rectangle.CornerPoints()
		var triPts = triangle.CornerPoints()
		var shPts = shape.CornerPoints()
		var ellPts = ellipse.CornerPoints()
		var roundPts = squircle.CornerPoints()

		view.DrawShapes(color.Darken(colCircle, 0.5), hexagonPts...)
		view.DrawLinesPath(8, colCircle, hexagonPts...)
		view.DrawPoints(4, colCircle, hexagonPts...)

		view.DrawShapes(color.Darken(colRect, 0.5), rectPts...)
		view.DrawLinesPath(8, colRect, rectPts...)
		view.DrawPoints(4, colRect, rectPts...)

		view.DrawShapes(palette.Gray, triPts...)
		view.DrawLinesPath(8, palette.White, triPts...)
		view.DrawPoints(4, palette.White, triPts...)

		view.DrawLinesPath(8, colShape, shPts...)
		view.DrawPoints(4, colShape, shPts...)

		view.DrawShapes(color.Darken(palette.Violet, 0.5), stars...)
		view.DrawLinesPath(8, palette.Violet, stars...)
		view.DrawPoints(4, palette.Violet, stars...)

		view.DrawShapes(palette.DarkGreen, ellPts...)
		view.DrawLinesPath(8, palette.Green, ellPts...)
		view.DrawPoints(4, palette.Green, ellPts...)

		view.DrawShapes(palette.Magenta, roundPts...)
		view.DrawLinesPath(8, palette.DarkMagenta, roundPts...)
		view.DrawPoints(4, palette.DarkMagenta, roundPts...)

		view.DrawPoints(16, palette.Green, crossPoints...)

		view.DrawTextDebug(true, true, true, true)
	}
}
