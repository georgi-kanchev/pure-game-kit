package main

import (
	"pure-tile-kit/engine/window"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	for window.KeepOpen() {
		if rl.IsKeyPressed(rl.KeyA) {
			window.Recreate()
		}
	}
}
