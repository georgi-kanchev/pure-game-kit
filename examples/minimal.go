package example

import (
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/window"
)

func Minimal() {
	window.Create("example - minimal", false, false)
	var view = graphics.NewView(1)
	window.SetTargetFPS(0)
	for window.KeepOpen() {
		view.DrawDebugInfo(false)
	}
}
