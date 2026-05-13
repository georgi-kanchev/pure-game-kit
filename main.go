package main

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/audio"
	"pure-game-kit/packages/engine"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
)

func main() {
	engine.Initialize("pure-game-kit", 60, 120, false, false)

	var sound = assets.LoadSound("examples/data/hammer.ogg", 3)
	var audio = audio.New(sound)

	engine.Run(func() {
		if keyboard.IsKeyJustPressed(key.A) {
			audio.Play()
		}

		if audio.IsJustFinished() {
			print("done")
		}
	})
}
