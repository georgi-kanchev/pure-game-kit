package example

import (
	"fmt"
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/utility/color"
	"pure-game-kit/window"
)

func Keyboard() {
	var cam = graphics.NewCamera(1)
	var text = "Type..."
	var font = assets.LoadDefaultFont()

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.PivotX, cam.PivotY = 0, 0
		text += keyboard.Input()

		if (keyboard.IsKeyPressedOnce(key.Backspace) || keyboard.IsKeyHeld(key.Backspace)) &&
			len(text) > 0 {
			var runes = []rune(text)
			text = string(runes[:len(runes)-1])
		}
		if keyboard.IsKeyPressedOnce(key.Enter) || keyboard.IsKeyHeld(key.Enter) {
			text += "\n"
		}

		if keyboard.IsAnyKeyPressedOnce() {
			fmt.Printf("%v\n", "hello, world")
		}

		cam.DrawText(font, text, 0, 0, 200, color.White)
	}
}
