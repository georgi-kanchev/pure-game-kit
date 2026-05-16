package main

import (
	"fmt"
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/engine"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
)

func main() {
	engine.Initialize("pure-game-kit", 60, 0, false, false)

	var view = graphics.NewView(1)
	var obj = graphics.NewObject(0, 0)

	var loadFlail = engine.NewWork(func() {
		obj.ImageId = assets.LoadImage("examples/data/flail.PNG")
	})
	engine.Run(func() {
		if keyboard.IsKeyJustPressed(key.A) {
			loadFlail.Start()
			fmt.Printf("work started\n")
		}

		if loadFlail.IsWorking() {
			fmt.Printf("working...\n")
		}

		if loadFlail.IsJustFinished() {
			fmt.Printf("work done!\n")
		}

		view.DrawObjects(&obj)
	})
	// example.Audio()
}
