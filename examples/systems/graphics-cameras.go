package example

import (
	"pure-game-kit/graphics"
	"pure-game-kit/internal"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/point"
	"pure-game-kit/window"
)

func MultipleCameras() {
	var topLeft = graphics.NewCamera(1)
	var topRight = graphics.NewCamera(3)
	var botLeft = graphics.NewCamera(1.5)
	var botRight = graphics.NewCamera(2)
	var ui = graphics.NewCamera(1)

	botLeft.Angle = 25

	// scene objects — created once, updated once per frame
	var orbiterCols = [6]uint{palette.Red, palette.Orange, palette.Yellow, palette.Green, palette.Cyan, palette.Purple}
	var orbiters [6]*graphics.Quad
	for i := range 6 {
		orbiters[i] = graphics.NewQuad(0, 0)
		orbiters[i].Width, orbiters[i].Height = 60, 60
		orbiters[i].Tint = orbiterCols[i]
	}

	var center = graphics.NewQuad(0, 0)
	center.Tint = palette.White

	var cornerCols = [4]uint{palette.Cyan, palette.Magenta, palette.Yellow, palette.Green}
	var cornerPos = [4][2]float32{{-300, -200}, {260, -200}, {-300, 160}, {260, 160}}
	var corners [4]*graphics.Quad
	for i := range 4 {
		corners[i] = graphics.NewQuad(cornerPos[i][0], cornerPos[i][1])
		corners[i].Width, corners[i].Height = 40, 40
		corners[i].Tint = cornerCols[i]
	}

	var quads = []*graphics.Quad{
		center,
		orbiters[0], orbiters[1], orbiters[2], orbiters[3], orbiters[4], orbiters[5],
		corners[0], corners[1], corners[2], corners[3],
	}

	for window.KeepOpen() {
		var ww, wh = window.Size()
		var hw, hh = float32(ww) / 2, float32(wh) / 2

		topLeft.Area = graphics.NewArea(0, 0, hw, hh)
		topRight.Area = graphics.NewArea(hw, 0, hw, hh)
		botLeft.Area = graphics.NewArea(0, hh, hw, hh)
		botRight.Area = graphics.NewArea(hw, hh, hw, hh)

		var t = float32(internal.Runtime)
		botRight.X = number.Cosine(t*0.7) * 120
		botRight.Y = number.Sine(t*0.5) * 80

		// update positions once per frame
		for i := range 6 {
			var angle = t*46 + float32(i)*60
			orbiters[i].X, orbiters[i].Y = point.MoveAtAngle(0, 0, angle, 180)
		}
		var pulse = 60 + number.Sine(t*2)*20
		center.Width, center.Height = pulse*2, pulse*2

		// draw all cameras
		for _, cam := range []*graphics.Camera{topLeft, topRight, botLeft, botRight} {
			cam.DrawColor(palette.DarkGray)
			cam.DrawGrid(1, 50, 50, palette.Gray)
			cam.DrawQuads(quads...)
		}

		// dividers and labels via full-screen ui camera (no Area = no scissor)
		ui.DrawLine(-hw, 0, hw, 0, 2, palette.White)
		ui.DrawLine(0, -hh, 0, hh, 2, palette.White)

		ui.DrawText("overview", -hw+10, -hh+10, 100)
		ui.DrawText("zoom x3", 10, -hh+10, 100)
		ui.DrawText("rotated 25 deg", -hw+10, 10, 100)
		ui.DrawText("panning", 10, 10, 100)

		topLeft.DrawTextDebug(true, false, false, false)
	}
}
