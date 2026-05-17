package main

import (
	example "pure-game-kit/examples"
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/utility/debug"
	"pure-game-kit/packages/window"
)

func main() {
	example.Audio()

	window.Create("Test", false, false)
	var view = graphics.NewView(1)
	var obj = graphics.NewObject(0, 0)

	for window.KeepOpen() {
		if keyboard.IsKeyJustPressed(key.A) {
			obj.ImageId = assets.LoadImage("examples/data/flail.PNG")
		}

		if keyboard.IsKeyJustPressed(key.S) {
			print(debug.MemoryUsage())
		}

		obj.X, obj.Y = view.MousePosition()

		view.DrawObjects(&obj)

		obj.X += 200
		view.DrawObjects(&obj)
	}
}
