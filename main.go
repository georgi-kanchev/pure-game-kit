package main

import (
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/time"
	"pure-game-kit/packages/window"
)

func main() {
	window.Create("split screen", false, true)

	var views = [4]graphics.View{
		graphics.NewView(1),
		graphics.NewView(1),
		graphics.NewView(1),
		graphics.NewView(1),
	}

	window.SetTargetFPS(60)

	for window.KeepOpen() {
		var w, h = window.Size()
		var hw, hh = float32(w) / 2, float32(h) / 2

		// top-left
		views[0].WindowArea = graphics.NewArea(0, 0, hw, hh)
		views[0].DrawColor(palette.DarkRed)
		views[0].DrawText(0, 0, 100, 0, palette.Red, "Hello, World!")
		views[0].Angle = number.Map(number.Sine(time.Running()), -1, 1, 0, 180)

		// top-right
		views[1].WindowArea = graphics.NewArea(hw, 0, hw, hh)
		views[1].DrawColor(palette.DarkGreen)
		views[1].DrawText(0, 0, 100, 0, palette.Red, "Hello, World!")

		// bottom-left
		views[2].WindowArea = graphics.NewArea(0, hh, hw, hh)
		views[2].DrawColor(palette.DarkBlue)
		views[2].DrawText(0, 0, 100, 0, palette.Red, "Hello, World!")

		// bottom-right
		views[3].WindowArea = graphics.NewArea(hw, hh, hw, hh)
		views[3].DrawColor(palette.DarkMagenta)
		views[3].DrawText(0, 0, 100, 0, palette.Red, "Hello, World!")
	}
}
