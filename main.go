package main

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/engine"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
)

func main() {
	engine.Initialize("pure-game-kit", 0, false, false)

	var view = graphics.NewView(1)
	var obj = graphics.NewObject(0, 0)

	engine.Run(func() {
		if keyboard.IsKeyJustPressed(key.A) {
			obj.ImageId = assets.LoadImage("examples/data/flail.PNG")
		}

		view.DrawObjects(&obj)
	})
	// example.Audio()
}
