package main

import (
	"pure-tile-kit/engine/window"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	window.Color.R = 40

	for window.KeepOpen() {
		if rl.IsKeyPressed(rl.KeyQ) {
			window.SetState(window.StateFullscreen)
		}
		if rl.IsKeyPressed(rl.KeyW) {
			window.SetState(window.StateWindowed)
		}
	}
}
