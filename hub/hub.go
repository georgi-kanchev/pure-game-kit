package main

import (
	number "pure-tile-kit/engine/utility/number"
	point "pure-tile-kit/engine/utility/point"
	"pure-tile-kit/engine/window"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	number.AnimateSpline(0.5, []point.F{{X: 1, Y: 4}, {X: 0.3, Y: 0.5}})

	for window.KeepOpen() {
		if rl.IsKeyPressed(rl.KeyA) {
			window.Recreate()
		}
	}
}
