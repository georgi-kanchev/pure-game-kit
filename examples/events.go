package example

import (
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/debug"
	"pure-game-kit/packages/window"
)

func Events() {
	window.Create("game", true, true)
	var view = graphics.NewView(1)
	var textObj = graphics.NewTextbox(0, 0, 2500, 1500, 0, "Hello, World!")
	debug.ProfileAllocations(10)
	for window.KeepOpen() {
		textObj.Text = debug.MemoryUsage()

		view.DrawObjects(&textObj)
	}
}
