package main

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/utility/debug"
	"pure-game-kit/packages/window"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
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

		obj.X, obj.Y = rl.GetMousePosition().X, rl.GetMousePosition().Y

		view.DrawObjects(&obj)

		obj.X += 200
		view.DrawObjects(&obj)
	}
}
