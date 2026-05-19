package main

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/time"
	"pure-game-kit/packages/window"
)

func main() {
	window.Create("game", false, false)
	var view = graphics.NewView(1)
	var obj = graphics.NewObject(0, 0)

	obj.Roundness = 1

	obj.ImageId = assets.LoadImage("examples/data/desert-0.png")
	obj.Effects = graphics.NewEffects()
	var w, h = obj.ImageId.Size()
	obj.Width, obj.Height = float32(w)*4, float32(h)*4

	obj.Effects.BorderSize = -2
	obj.Effects.BorderColor = palette.Red

	for window.KeepOpen() {
		var _ = (number.Sine(time.Running()) + 1) / 2
		//byte(number.Map(loop, 0, 1, 0, 15))
		view.DrawObjects(&obj)
	}
}
