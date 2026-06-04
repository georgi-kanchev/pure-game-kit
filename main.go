package main

import (
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/window"
)

func main() {
	window.Create("game", false, true)
	var view = graphics.NewView(1)

	window.SetTargetFPS(60)
	for window.KeepOpen() {
		_ = view
	}
}
