package example

import (
	"pure-kit/engine/data/assets"
	"pure-kit/engine/graphics"
	"pure-kit/engine/input/keyboard"
	"pure-kit/engine/input/keyboard/key"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
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

		cam.DrawText(font, text, 0, 0, 200, color.White)
	}
}
