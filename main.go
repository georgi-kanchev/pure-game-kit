package main

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/window"
)

func main() {
	window.Create("game", false, false)
	var view = graphics.NewView(1)
	var obj = graphics.NewImage(0, 0, assets.LoadImage("examples/data/flail.PNG"))

	obj.Width *= 4
	obj.Height *= 4
	for window.KeepOpen() {
		view.DrawObjects(&obj)
	}
}
