package example

import (
	"fmt"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/motion"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/window"
)

func Tweens() {
	var fade = motion.NewTween(0)
	var view = graphics.NewView(1)

	window.Create("example - tweens", false, false)
	for window.KeepOpen() {
		if keyboard.IsKeyJustPressed(key.A) {
			fade.GoTo(1.0, nil, 255).GoTo(1.0, nil, 0).Wait(0.5).GoTo(1.0, nil, 128)
		}
		if keyboard.IsKeyJustPressed(key.S) {
			fade.ForceReset()
		}
		if keyboard.IsKeyJustPressed(key.D) {
			fade.ForceFinish()
		}

		if fade.IsJustFinished() {
			fmt.Printf("finished\n")
		}

		var value = fade.CurrentValues()[0]
		fmt.Printf("value: %v\n", value)
		view.DrawColor(color.RGB(byte(value), 0, 0))
	}
}
