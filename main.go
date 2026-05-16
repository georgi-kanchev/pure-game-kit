package main

import (
	"pure-game-kit/packages/engine"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	engine.Initialize("pure-game-kit", 60, 120, false, false)

	engine.Run(func() {
		rl.LoadTexture("examples/data/flail.PNG")
	})
	// example.Audio()
}
