package main

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
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

	for window.KeepOpen() {
		obj.Effects.Gamma = (number.Sine(time.Running()) + 1) / 2
		view.DrawObjects(&obj)
	}
}
