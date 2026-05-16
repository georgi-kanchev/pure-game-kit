package main

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/window"
)

func main() {
	// window.Create("Test", false, false)
	var view = graphics.NewView(1)
	var obj = graphics.NewObject(0, 0)

	for window.KeepOpen() {
		if keyboard.IsKeyJustPressed(key.A) {
			obj.ImageId = assets.LoadImage("examples/data/flail.PNG")
		}

		view.DrawObjects(&obj)
	}
	// example.Audio()
}
