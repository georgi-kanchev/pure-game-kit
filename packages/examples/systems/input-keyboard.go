package example

import (
	"fmt"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/window"
)

func Keyboard() {
	var cam = graphics.NewCamera(1)
	var text = "Type..."

	for window.KeepOpen() {
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
		cam.DrawText(text, x, y, 200)
	}
}
