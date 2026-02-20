package example

import (
	"fmt"
	"pure-game-kit/execution/condition"
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/time"
	"pure-game-kit/window"
)

func Shapes() {
	var cam = graphics.NewCamera(1)
	var shape = geometry.NewShapeCorners([][2]float32{
		{0, 0},
		{50, -20},
		{100, 0},
		{0, 100},
		{50, 120},
		{100, 100},
	}...)
	var triangle = geometry.NewShapeCorners([2]float32{}, [2]float32{100, 100}, [2]float32{-100, 100})
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

	var stars = [][2]float32{
		{600 + 300, 100}, {600 + 350, 200}, {600 + 450, 200}, {600 + 370, 260},
		{600 + 400, 360}, {600 + 300, 300}, {600 + 200, 360}, {600 + 230, 260},
		{600 + 150, 200}, {600 + 250, 200}, {600 + 300, 100},
		{number.NaN(), number.NaN()},
		{300, 100}, {350, 200}, {450, 200}, {370, 260},
		{400, 360}, {300, 300}, {200, 360}, {230, 260},
		{150, 200}, {250, 200}, {300, 100},
	}

	window.FrameRateLimit = 0

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()

		squircle.Angle += time.FrameDelta() * 60
		shape.Angle += time.FrameDelta() * 60
		var mx, my = cam.MousePosition()
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

		cam.DrawShapesFast(color.Darken(colCircle, 0.5), hexagonPts...)
		cam.DrawLinesPath(8, colCircle, hexagonPts...)
		cam.DrawPoints(4, colCircle, hexagonPts...)

		cam.DrawShapesFast(color.Darken(colRect, 0.5), rectPts...)
		cam.DrawLinesPath(8, colRect, rectPts...)
		cam.DrawPoints(4, colRect, rectPts...)

		cam.DrawShapesFast(palette.Gray, triPts...)
		cam.DrawLinesPath(8, palette.White, triPts...)
		cam.DrawPoints(4, palette.White, triPts...)

		cam.DrawLinesPath(8, colShape, shPts...)
		cam.DrawPoints(4, colShape, shPts...)

		// not DrawShapesFast because stars are concave
		cam.DrawShapes(color.Darken(palette.Violet, 0.5), stars...)
		cam.DrawLinesPath(8, palette.Violet, stars...)
		cam.DrawPoints(4, palette.Violet, stars...)

		cam.DrawShapesFast(palette.DarkGreen, ellPts...)
		cam.DrawLinesPath(8, palette.Green, ellPts...)
		cam.DrawPoints(4, palette.Green, ellPts...)

		cam.DrawShapesFast(palette.Magenta, roundPts...)
		cam.DrawLinesPath(8, palette.DarkMagenta, roundPts...)
		cam.DrawPoints(4, palette.DarkMagenta, roundPts...)

		for _, v := range crossPoints {
			cam.DrawCircle(v[0], v[1], 16, palette.Green)
		}

		fmt.Printf("time.FrameRate(): %v\n", time.FrameRate())
	}
}
