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

	obj.ImageId = assets.LoadImage("examples/data/flail.PNG")
	obj.Effects = graphics.NewEffects()
	var w, h = obj.ImageId.Size()
	obj.Width, obj.Height = float32(w)*4, float32(h)*4

	obj.Effects.OutlineSize = 2
	obj.Effects.OutlineColor = palette.Red

	// obj.Effects.BorderColor = palette.Green

	for window.KeepOpen() {
		var loop = (number.Sine(time.Running()) + 1) / 2
		obj.Effects.BorderSize = number.Map(loop, 0, 1, 0, 40)
		view.DrawObjects(&obj)
	}
}
