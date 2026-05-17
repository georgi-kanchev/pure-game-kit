package main

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/window"
)

func main() {
	window.Create("game", false, false)
	var view = graphics.NewView(1)
	var obj = graphics.NewObject(0, 0)

	obj.Roundness = 1

	obj.ImageId = assets.LoadImage("examples/data/flail.PNG")
	obj.Effects = graphics.NewEffects()
	obj.Width, obj.Height = 500, 500

	for window.KeepOpen() {
		view.DrawObjects(&obj)
	}
}
