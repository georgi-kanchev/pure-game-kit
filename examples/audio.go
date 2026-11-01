package example

import (
	"fmt"
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/audio"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/window"
)

func Audio() {
	window.Create("examples - audio", false, false)
	var wood = assets.LoadSound(3, "examples/data/wood.mp3")
	var hammer = assets.LoadSound(3, "examples/data/hammer.ogg")
	var variation = assets.LoadVariations(wood, hammer)
	var audio = audio.New(variation)

	for window.KeepOpen() {
		if keyboard.IsKeyJustPressed(key.A) {
			audio.Play()
			fmt.Printf("play\n")
		}

		if keyboard.IsKeyJustPressed(key.S) {
			audio.Volume = 0.5
			audio.ApplyProperties()
			fmt.Printf("applied volume 50%%\n")
		}

		if audio.IsJustFinished() {
			fmt.Printf("just finished: %v\n", "true")
		}

		if keyboard.IsKeyJustPressed(key.Space) {
			if audio.IsPlaying() {
				audio.Pause()
				fmt.Printf("paused\n")
			} else {
				audio.Play()
				fmt.Printf("resumed\n")
			}
		}
	}
}
