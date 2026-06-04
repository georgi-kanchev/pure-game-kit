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
		views[0].DrawText(0, 0, 100, 0, palette.Red, "Hello, World!")
		views[0].Angle = number.Map(number.Sine(time.Running()), -1, 1, 0, 5)

		views[1].DrawText(0, 0, 100, 0, palette.Red, "Hello, World!")
	}
}
