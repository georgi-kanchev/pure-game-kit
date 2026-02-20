package example

import (
	"fmt"
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/window"
)

func Keyboard() {
	var cam = graphics.NewCamera(1)
	var text = "Type..."
	var font = assets.LoadDefaultFont()

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		text += keyboard.Input()

		if (keyboard.IsKeyJustPressed(key.Backspace) || keyboard.IsKeyHeld(key.Backspace)) &&
			len(text) > 0 {
			var runes = []rune(text)
			text = string(runes[:len(runes)-1])
		}
		if keyboard.IsKeyJustPressed(key.Enter) || keyboard.IsKeyHeld(key.Enter) {
			text += "\n"
		}

		if keyboard.IsAnyKeyJustPressed() {
			fmt.Printf("%v\n", "hello, world")
		}

		var x, y = cam.PointFromScreen(0, 0)
		cam.DrawTextAdvanced(font, text, x, y, 200, 0.5, 0, palette.White)
	}
}
