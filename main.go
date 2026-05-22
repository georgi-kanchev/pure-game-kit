package main

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/debug"
	"pure-game-kit/packages/window"
)

func main() {
	window.Create("game", false, false)
	var view = graphics.NewView(1)
	var font = assets.LoadFont("tools/msdf-atlas-gen/Libre.png", "tools/msdf-atlas-gen/Libre.json")
	var obj = graphics.NewTextbox(0, 0, 100, 100, font, "Lemon the Cat")

	obj.Width *= 4
	obj.Height *= 4
	obj.X -= 500
	obj.TextLineHeight = 200
	obj.Color = palette.Gray

	for window.KeepOpen() {
		if keyboard.IsKeyJustPressed(key.F5) {
			print(debug.MemoryUsage())
		}

		view.DrawObjects(&obj)
	}
}
