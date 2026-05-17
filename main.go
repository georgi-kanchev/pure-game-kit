package main

import (
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/window"
)

func main() {
	window.Create("game", false, false)
	var view = graphics.NewView(1)
	var obj = graphics.NewObject(0, 0)

	obj.Roundness = 1

	for window.KeepOpen() {
		view.DrawObjects(&obj)
	}
}
